// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package stroge

import (
	"os"
	"sync"

	"github.com/donnie4w/wfs/stub"
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

func (t *ldb) SnapshotToStream(prefix []byte, streamfunc func(bean *stub.SnapshotBean) bool) (err error) {
	defer util.Recover()
	snap, err := t.db.GetSnapshot()
	if err != nil {
		return err
	}
	defer snap.Release()
	var ran *levelutil.Range
	if prefix != nil {
		ran = levelutil.BytesPrefix(prefix)
	}
	iter := snap.NewIterator(ran, nil)
	defer iter.Release()
	for iter.Next() {
		bean := new(stub.SnapshotBean)
		bean.Copy(iter.Key(), iter.Value())
		if !streamfunc(bean) {
			break
		}
	}
	return
}

func (t *ldb) LoadSnapshotBeans(beans ...*stub.SnapshotBean) (err error) {
	defer util.Recover()
	if beans != nil {
		batch := make(map[*[]byte][]byte, 0)
		for _, bean := range beans {
			batch[&bean.Key] = bean.Value
		}
		err = t.BatchPut(batch)
	}
	return
}

func (t *ldb) LoadSnapshotBean(bean *stub.SnapshotBean) (err error) {
	defer util.Recover()
	if bean != nil {
		err = t.Put(bean.Key, bean.Value)
	}
	return
}