// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package stroge

import (
	"bytes"
	"io"
	"os"
	"sync"

	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	levelutil "github.com/syndtr/goleveldb/leveldb/util"
)

var dbMap map[string]*ldb = make(map[string]*ldb, 0)

var mux = new(sync.Mutex)

type ldb struct {
	db       *leveldb.DB
	dir_name string
}

func New(_dir string) (db *ldb, err error) {
	mux.Lock()
	defer mux.Unlock()
	if err = os.MkdirAll(_dir, 0777); err != nil {
		return
	}
	var ok bool
	if db, ok = dbMap[_dir]; !ok {
		db = &ldb{dir_name: _dir}
		if err = db.Open(); err == nil {
			dbMap[_dir] = db
		} else {
			return
		}
	}
	return
}

func (t *ldb) Open() (err error) {
	options := &opt.Options{
		Filter:                 filter.NewBloomFilter(10),
		OpenFilesCacheCapacity: 1 << 10,
		BlockCacheCapacity:     int(sys.DBBUFFER / 2),
		WriteBuffer:            int(sys.DBBUFFER / 4),
	}
	t.db, err = leveldb.OpenFile(t.dir_name, options)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		t.db, err = leveldb.RecoverFile(t.dir_name, nil)
	}
	return
}

func (t *ldb) Close() (err error) {
	return t.db.Close()
}

func (t *ldb) Has(key []byte) (b bool) {
	b, _ = t.db.Has(key, nil)
	return
}

func (t *ldb) Put(key, value []byte) (err error) {
	return t.db.Put(key, value, &opt.WriteOptions{Sync: sys.SYNC})
}

func (t *ldb) Get(key []byte) (value []byte, err error) {
	value, err = t.db.Get(key, nil)
	return
}

func (t *ldb) GetString(key []byte) (value string, err error) {
	v, er := t.db.Get(key, nil)
	return string(v), er
}

func (t *ldb) Del(key []byte) (err error) {
	return t.db.Delete(key, &opt.WriteOptions{Sync: sys.SYNC})
}

func (t *ldb) BatchPut(kvmap map[*[]byte][]byte) (err error) {
	batch := new(leveldb.Batch)
	for k, v := range kvmap {
		batch.Put(*k, v)
	}
	err = t.db.Write(batch, &opt.WriteOptions{Sync: sys.SYNC})
	return
}

func (t *ldb) Batch(put map[*[]byte][]byte, del [][]byte) (err error) {
	batch := new(leveldb.Batch)
	if put != nil {
		for k, v := range put {
			batch.Put([]byte(*k), v)
		}
	}
	if del != nil {
		for _, v := range del {
			batch.Delete(v)
		}
	}
	err = t.db.Write(batch, &opt.WriteOptions{Sync: sys.SYNC})
	return
}

func (t *ldb) GetLike(prefix []byte) (datamap map[string][]byte, err error) {
	iter := t.db.NewIterator(levelutil.BytesPrefix(prefix), nil)
	if iter != nil {
		defer iter.Release()
		datamap = make(map[string][]byte, 0)
		for iter.Next() {
			bs := make([]byte, len(iter.Value()))
			copy(bs, iter.Value())
			datamap[string(iter.Key())] = bs
		}
		err = iter.Error()
	}
	return
}

// all key
func (t *ldb) GetKeys() (bys []string, err error) {
	iter := t.db.NewIterator(nil, nil)
	defer iter.Release()
	bys = make([]string, 0)
	for iter.Next() {
		bys = append(bys, string(iter.Key()))
	}
	err = iter.Error()
	return
}

// all key by prefix
func (t *ldb) GetKeysPrefix(prefix []byte) (bys [][]byte, err error) {
	iter := t.db.NewIterator(levelutil.BytesPrefix(prefix), nil)
	defer iter.Release()
	bys = make([][]byte, 0)
	for iter.Next() {
		bs := make([]byte, len(iter.Value()))
		copy(bs, iter.Value())
		bys = append(bys, bs)
	}
	err = iter.Error()
	return
}

// all key by prefix
func (t *ldb) GetValuesPrefix(prefix []byte) (bys [][]byte, err error) {
	iter := t.db.NewIterator(levelutil.BytesPrefix(prefix), nil)
	defer iter.Release()
	bys = make([][]byte, 0)
	for iter.Next() {
		bs := make([]byte, len(iter.Value()))
		copy(bs, iter.Value())
		bys = append(bys, bs)
	}
	err = iter.Error()
	return
}

// all key by prefix
func (t *ldb) GetKeysPrefixLimit(prefix []byte, limit int) (bys [][]byte, err error) {
	iter := t.db.NewIterator(levelutil.BytesPrefix(prefix), nil)
	defer iter.Release()
	bys = make([][]byte, 0)
	i := 0
	for iter.Next() {
		if i > limit {
			break
		}
		bs := make([]byte, len(iter.Key()))
		copy(bs, iter.Key())
		bys = append(bys, bs)
		i++
	}
	err = iter.Error()
	return
}

// all value by prefix
func (t *ldb) GetValuesPrefixLimit(prefix []byte, limit int) (bys [][]byte, err error) {
	iter := t.db.NewIterator(levelutil.BytesPrefix(prefix), nil)
	defer iter.Release()
	bys = make([][]byte, 0)
	i := 0
	for iter.Next() {
		if i > limit {
			break
		}
		bs := make([]byte, len(iter.Value()))
		copy(bs, iter.Value())
		bys = append(bys, bs)
		i++
	}
	err = iter.Error()
	return
}

/*
*
Start of the key range, include in the range.
Limit of the key range, not include in the range.
*/
func (t *ldb) GetIterLimit(prefix string, limit string) (datamap map[string][]byte, err error) {
	iter := t.db.NewIterator(&levelutil.Range{Start: []byte(prefix), Limit: []byte(limit)}, nil)
	defer iter.Release()
	datamap = make(map[string][]byte, 0)
	for iter.Next() {
		data, er := t.db.Get(iter.Key(), nil)
		if er == nil {
			datamap[string(iter.Key())] = data
		}
	}
	err = iter.Error()
	return
}

func (t *ldb) Snapshot() (*leveldb.Snapshot, error) {
	return t.db.GetSnapshot()
}

// ////////////////////////////////////////////////////////
type BakStub struct {
	Key   []byte
	Value []byte
}

func (t *BakStub) copy(k, v []byte) {
	t.Key, t.Value = make([]byte, len(k)), make([]byte, len(v))
	copy(t.Key, k)
	copy(t.Value, v)
}

func (t *ldb) BackupToDisk(filename string, prefix []byte) error {
	defer util.Recover()
	snap, err := t.Snapshot()
	if err != nil {
		return err
	}
	defer snap.Release()
	bs := _TraverseSnap(snap, prefix)
	b, e := goutil.Encode(bs)
	if e != nil {
		return e
	}
	f, er := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if er != nil {
		return er
	}
	defer f.Close()
	_, er = f.Write(b)
	return er
}

func RecoverBackup(filename string) (bs []*BakStub) {
	defer util.Recover()
	f, er := os.Open(filename)
	if er == nil {
		defer f.Close()
	} else {
		return
	}
	var buf bytes.Buffer
	_, err := io.Copy(&buf, f)
	if err == nil {
		goutil.Decode[[]*BakStub](buf.Bytes())
	}
	return
}
func (t *ldb) LoadDataFile(filename string) (err error) {
	bs := RecoverBackup(filename)
	for _, v := range bs {
		err = t.Put(v.Key, v.Value)
	}
	return
}

func (t *ldb) LoadBytes(buf []byte) (err error) {
	var bs []*BakStub
	err = util.Decode(buf, &bs)
	if err == nil {
		for _, v := range bs {
			err = t.Put(v.Key, v.Value)
		}
	}
	return
}

func _TraverseSnap(snap *leveldb.Snapshot, prefix []byte) (bs []*BakStub) {
	ran := new(levelutil.Range)
	if prefix != nil {
		ran = levelutil.BytesPrefix(prefix)
	} else {
		ran = nil
	}
	iter := snap.NewIterator(ran, nil)
	defer iter.Release()
	bs = make([]*BakStub, 0)
	for iter.Next() {
		ss := new(BakStub)
		ss.copy(iter.Key(), iter.Value())
		bs = append(bs, ss)
	}
	return
}
