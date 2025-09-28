// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package stor

import (
	"github.com/donnie4w/wfs/stub"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	levelutil "github.com/syndtr/goleveldb/leveldb/util"
	"os"
)

type ldb struct {
	db       *leveldb.DB
	dir_name string
	config   *sys.LevelDBConfig
}

func newLevelDB(config *sys.DBConfig) *ldb {
	return &ldb{
		dir_name: config.Path,
		config:   config.LevelDB,
	}
}

func (t *ldb) Open() (err error) {
	if err = os.MkdirAll(t.dir_name, 0777); err != nil {
		return
	}

	options := &opt.Options{
		Filter:                 filter.NewBloomFilter(t.config.BloomFilterBits),
		OpenFilesCacheCapacity: t.config.OpenFilesCache,
		BlockCacheCapacity:     t.config.BlockCacheSize,
		WriteBuffer:            t.config.WriteBufferSize,
	}

	t.db, err = leveldb.OpenFile(t.dir_name, options)
	if _, corrupted := err.(*errors.ErrCorrupted); corrupted {
		t.db, err = leveldb.RecoverFile(t.dir_name, nil)
	}
	return
}

func (t *ldb) Close() (err error) {
	if t.db != nil {
		return t.db.Close()
	}
	return nil
}

func (t *ldb) Has(key []byte) (b bool) {
	if t.db == nil {
		return false
	}
	b, _ = t.db.Has(key, nil)
	return
}

func (t *ldb) Put(key, value []byte) (err error) {
	if t.db == nil {
		return os.ErrInvalid
	}
	return t.db.Put(key, value, &opt.WriteOptions{Sync: sys.SYNC})
}

func (t *ldb) Get(key []byte) (value []byte, err error) {
	if t.db == nil {
		return nil, os.ErrInvalid
	}
	return t.db.Get(key, nil)
}

func (t *ldb) GetString(key []byte) (value string, err error) {
	v, err := t.Get(key)
	return string(v), err
}

func (t *ldb) Del(key []byte) (err error) {
	if t.db == nil {
		return os.ErrInvalid
	}
	return t.db.Delete(key, &opt.WriteOptions{Sync: sys.SYNC})
}

func (t *ldb) BatchPut(kvmap map[*[]byte][]byte) (err error) {
	if t.db == nil {
		return os.ErrInvalid
	}
	batch := new(leveldb.Batch)
	for k, v := range kvmap {
		batch.Put(*k, v)
	}
	return t.db.Write(batch, &opt.WriteOptions{Sync: sys.SYNC})
}

func (t *ldb) Batch(put map[*[]byte][]byte, del [][]byte) (err error) {
	if t.db == nil {
		return os.ErrInvalid
	}
	batch := new(leveldb.Batch)
	if put != nil {
		for k, v := range put {
			batch.Put(*k, v)
		}
	}
	if del != nil {
		for _, v := range del {
			batch.Delete(v)
		}
	}
	return t.db.Write(batch, &opt.WriteOptions{Sync: sys.SYNC})
}

func (t *ldb) GetLike(prefix []byte) (datamap map[string][]byte, err error) {
	if t.db == nil {
		return nil, os.ErrInvalid
	}
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

func (t *ldb) GetKeys() (bys []string, err error) {
	if t.db == nil {
		return nil, os.ErrInvalid
	}
	iter := t.db.NewIterator(nil, nil)
	defer iter.Release()
	bys = make([]string, 0)
	for iter.Next() {
		bys = append(bys, string(iter.Key()))
	}
	err = iter.Error()
	return
}

func (t *ldb) GetKeysPrefix(prefix []byte) (bys [][]byte, err error) {
	if t.db == nil {
		return nil, os.ErrInvalid
	}
	iter := t.db.NewIterator(levelutil.BytesPrefix(prefix), nil)
	defer iter.Release()
	bys = make([][]byte, 0)
	for iter.Next() {
		bs := make([]byte, len(iter.Key()))
		copy(bs, iter.Key())
		bys = append(bys, bs)
	}
	err = iter.Error()
	return
}

func (t *ldb) GetValuesPrefix(prefix []byte) (bys [][]byte, err error) {
	if t.db == nil {
		return nil, os.ErrInvalid
	}
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

func (t *ldb) GetKeysPrefixLimit(prefix []byte, limit int) (bys [][]byte, err error) {
	if t.db == nil {
		return nil, os.ErrInvalid
	}
	iter := t.db.NewIterator(levelutil.BytesPrefix(prefix), nil)
	defer iter.Release()
	bys = make([][]byte, 0)
	i := 0
	for iter.Next() {
		if i >= limit {
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

func (t *ldb) GetValuesPrefixLimit(prefix []byte, limit int) (bys [][]byte, err error) {
	if t.db == nil {
		return nil, os.ErrInvalid
	}
	iter := t.db.NewIterator(levelutil.BytesPrefix(prefix), nil)
	defer iter.Release()
	bys = make([][]byte, 0)
	i := 0
	for iter.Next() {
		if i >= limit {
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

func (t *ldb) GetIterLimit(prefix string, limit string) (datamap map[string][]byte, err error) {
	if t.db == nil {
		return nil, os.ErrInvalid
	}
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
	if t.db == nil {
		return os.ErrInvalid
	}
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
	if t.db == nil {
		return os.ErrInvalid
	}
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
	if t.db == nil {
		return os.ErrInvalid
	}
	defer util.Recover()
	if bean != nil {
		err = t.Put(bean.Key, bean.Value)
	}
	return
}
