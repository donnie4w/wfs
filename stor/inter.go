// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package stor

import (
	"fmt"
	"sync"

	"github.com/donnie4w/wfs/stub"
	"github.com/donnie4w/wfs/sys"
)

type DB interface {
	Open() error
	Close() error
	Has(key []byte) bool
	Put(key, value []byte) error
	Get(key []byte) ([]byte, error)
	GetString(key []byte) (string, error)
	Del(key []byte) error
	BatchPut(kvmap map[*[]byte][]byte) error
	Batch(put map[*[]byte][]byte, del [][]byte) error
	GetLike(prefix []byte) (map[string][]byte, error)
	GetKeys() ([]string, error)
	GetKeysPrefix(prefix []byte) ([][]byte, error)
	GetValuesPrefix(prefix []byte) ([][]byte, error)
	GetKeysPrefixLimit(prefix []byte, limit int) ([][]byte, error)
	GetValuesPrefixLimit(prefix []byte, limit int) ([][]byte, error)
	GetIterLimit(prefix string, limit string) (map[string][]byte, error)
	SnapshotToStream(prefix []byte, streamfunc func(bean *stub.SnapshotBean) bool) error
	LoadSnapshotBeans(beans ...*stub.SnapshotBean) error
	LoadSnapshotBean(bean *stub.SnapshotBean) error
}

var dbMap = make(map[string]DB, 0)
var dbmux = new(sync.Mutex)

func DefaultConfig(dbType sys.DBType, path string) *sys.DBConfig {
	config := &sys.DBConfig{
		Type: dbType,
		Path: path,
	}

	switch dbType {
	case sys.DBTypeLevelDB:
		config.LevelDB = &sys.LevelDBConfig{
			BloomFilterBits: 10,
			OpenFilesCache:  1 << 10, // 1024
			BlockCacheSize:  sys.DBBUFFER / 2,
			WriteBufferSize: sys.DBBUFFER / 4,
		}
	case sys.DBTypeSQLite:
		config.SQLite = &sys.SQLiteConfig{
			JournalMode: "WAL",
			SyncMode:    "NORMAL",
			CacheSize:   128 * 1 << 20,
			MMapSize:    256 * 1 << 20,
			BusyTimeout: 5000,
			AutoVacuum:  "INCREMENTAL",
			ForeignKeys: true,
		}
	}
	return config
}

func New(config *sys.DBConfig) (DB, error) {
	dbmux.Lock()
	defer dbmux.Unlock()

	if db, ok := dbMap[config.Path]; ok {
		return db, nil
	}

	var db DB
	switch config.Type {
	case sys.DBTypeLevelDB:
		db = newLevelDB(config)
	case sys.DBTypeSQLite:
		db = newSQLiteDB(config)
	default:
		return nil, fmt.Errorf("unsupported database type: %s", config.Type)
	}

	if err := db.Open(); err != nil {
		return nil, err
	}

	dbMap[config.Path] = db
	return db, nil
}

func NewWithDefaults(dbType int) (DB, error) {
	var config *sys.DBConfig
	if dbType == 2 {
		config = DefaultConfig(sys.DBTypeSQLite, sys.WFSDATA+"/wfsdb/wfs.db")
	} else {
		config = DefaultConfig(sys.DBTypeLevelDB, sys.WFSDATA+"/wfsdb")
	}
	return New(config)
}

func CloseAll() error {
	dbmux.Lock()
	defer dbmux.Unlock()

	var lastErr error
	for path, db := range dbMap {
		if err := db.Close(); err != nil {
			lastErr = err
		}
		delete(dbMap, path)
	}
	return lastErr
}
