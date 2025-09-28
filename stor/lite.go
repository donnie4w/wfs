// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package stor

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/donnie4w/wfs/stub"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
)

type sqliteDB struct {
	db     *sql.DB
	path   string
	config *sys.SQLiteConfig
}

func newSQLiteDB(config *sys.DBConfig) *sqliteDB {
	return &sqliteDB{
		path:   config.Path,
		config: config.SQLite,
	}
}

func (s *sqliteDB) Open() error {
	dir := s.path
	if strings.Contains(s.path, ".db") {
		dir = filepath.Dir(s.path)
	}
	if err := os.MkdirAll(dir, 0777); err != nil {
		return err
	}

	dsn := s.buildDSN()

	var err error
	s.db, err = sql.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	if err := s.db.Ping(); err != nil {
		return err
	}

	s.db.SetMaxOpenConns(50)
	s.db.SetMaxIdleConns(25)
	s.db.SetConnMaxLifetime(0)

	return s.createTable()
}

func (s *sqliteDB) buildDSN() string {
	params := []string{
		fmt.Sprintf("_journal_mode=%s", s.config.JournalMode),
		fmt.Sprintf("_synchronous=%s", s.config.SyncMode),
		fmt.Sprintf("_cache_size=%d", s.config.CacheSize),
		fmt.Sprintf("_busy_timeout=%d", s.config.BusyTimeout),
	}

	if s.config.ForeignKeys {
		params = append(params, "_foreign_keys=on")
	} else {
		params = append(params, "_foreign_keys=off")
	}

	return s.path + "?" + strings.Join(params, "&")
}

func (s *sqliteDB) createTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS kv_store (
			key BLOB PRIMARY KEY,
			value BLOB NOT NULL
		)
	`)
	if err != nil {
		return err
	}

	pragmas := []string{
		fmt.Sprintf("PRAGMA mmap_size = %d;", s.config.MMapSize),
	}

	if s.config.AutoVacuum != "" {
		pragmas = append(pragmas, fmt.Sprintf("PRAGMA auto_vacuum = %s;", s.config.AutoVacuum))
	}

	for _, pragma := range pragmas {
		if _, err := s.db.Exec(pragma); err != nil {
			return err
		}
	}

	return nil
}

func (s *sqliteDB) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *sqliteDB) Has(key []byte) bool {
	var exists int
	err := s.db.QueryRow("SELECT 1 FROM kv_store WHERE key = ?", key).Scan(&exists)
	return err == nil && exists == 1
}

func (s *sqliteDB) Put(key, value []byte) error {
	_, err := s.db.Exec(`
		INSERT OR REPLACE INTO kv_store (key, value) 
		VALUES (?, ?)
	`, key, value)
	return err
}

func (s *sqliteDB) Get(key []byte) ([]byte, error) {
	var value []byte
	err := s.db.QueryRow("SELECT value FROM kv_store WHERE key = ?", key).Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return value, err
}

func (s *sqliteDB) GetString(key []byte) (string, error) {
	value, err := s.Get(key)
	if err != nil {
		return "", err
	}
	return string(value), nil
}

func (s *sqliteDB) Del(key []byte) error {
	_, err := s.db.Exec("DELETE FROM kv_store WHERE key = ?", key)
	return err
}

func (s *sqliteDB) BatchPut(kvmap map[*[]byte][]byte) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO kv_store (key, value) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for k, v := range kvmap {
		if _, err := stmt.Exec(*k, v); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *sqliteDB) Batch(put map[*[]byte][]byte, del [][]byte) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	if del != nil && len(del) > 0 {
		delStmt, err := tx.Prepare("DELETE FROM kv_store WHERE key = ?")
		if err != nil {
			tx.Rollback()
			return err
		}
		defer delStmt.Close()

		for _, key := range del {
			if _, err := delStmt.Exec(key); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if put != nil && len(put) > 0 {
		putStmt, err := tx.Prepare("INSERT OR REPLACE INTO kv_store (key, value) VALUES (?, ?)")
		if err != nil {
			tx.Rollback()
			return err
		}
		defer putStmt.Close()

		for k, v := range put {
			if _, err := putStmt.Exec(*k, v); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	return tx.Commit()
}

func (s *sqliteDB) GetLike(prefix []byte) (map[string][]byte, error) {
	rows, err := s.db.Query(`
		SELECT key, value FROM kv_store 
		WHERE key >= ? AND key < ?
		ORDER BY key`,
		prefix, s.nextPrefix(prefix))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]byte)
	for rows.Next() {
		var key []byte
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[string(key)] = value
	}

	return result, rows.Err()
}

func (s *sqliteDB) nextPrefix(prefix []byte) []byte {
	result := make([]byte, len(prefix))
	copy(result, prefix)

	for i := len(result) - 1; i >= 0; i-- {
		if result[i] < 0xFF {
			result[i]++
			return result[:i+1]
		}
	}

	return nil
}

func (s *sqliteDB) GetKeys() ([]string, error) {
	rows, err := s.db.Query("SELECT key FROM kv_store ORDER BY key")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []string
	for rows.Next() {
		var key []byte
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, string(key))
	}

	return keys, rows.Err()
}

func (s *sqliteDB) GetKeysPrefix(prefix []byte) ([][]byte, error) {
	rows, err := s.db.Query(`
		SELECT key FROM kv_store 
		WHERE key >= ? AND key < ?
		ORDER BY key`,
		prefix, s.nextPrefix(prefix))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys [][]byte
	for rows.Next() {
		var key []byte
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)
		keys = append(keys, keyCopy)
	}

	return keys, rows.Err()
}

func (s *sqliteDB) GetValuesPrefix(prefix []byte) ([][]byte, error) {
	rows, err := s.db.Query(`
		SELECT value FROM kv_store 
		WHERE key >= ? AND key < ?
		ORDER BY key`,
		prefix, s.nextPrefix(prefix))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values [][]byte
	for rows.Next() {
		var value []byte
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		valueCopy := make([]byte, len(value))
		copy(valueCopy, value)
		values = append(values, valueCopy)
	}

	return values, rows.Err()
}

func (s *sqliteDB) GetKeysPrefixLimit(prefix []byte, limit int) ([][]byte, error) {
	rows, err := s.db.Query(`
		SELECT key FROM kv_store 
		WHERE key >= ? AND key < ?
		ORDER BY key
		LIMIT ?`,
		prefix, s.nextPrefix(prefix), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys [][]byte
	for rows.Next() {
		var key []byte
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keyCopy := make([]byte, len(key))
		copy(keyCopy, key)
		keys = append(keys, keyCopy)
	}

	return keys, rows.Err()
}

func (s *sqliteDB) GetValuesPrefixLimit(prefix []byte, limit int) ([][]byte, error) {
	rows, err := s.db.Query(`
		SELECT value FROM kv_store 
		WHERE key >= ? AND key < ?
		ORDER BY key
		LIMIT ?`,
		prefix, s.nextPrefix(prefix), limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var values [][]byte
	for rows.Next() {
		var value []byte
		if err := rows.Scan(&value); err != nil {
			return nil, err
		}
		valueCopy := make([]byte, len(value))
		copy(valueCopy, value)
		values = append(values, valueCopy)
	}

	return values, rows.Err()
}

func (s *sqliteDB) GetIterLimit(prefix string, limit string) (map[string][]byte, error) {
	rows, err := s.db.Query(`
		SELECT key, value FROM kv_store 
		WHERE key >= ? AND key < ?
		ORDER BY key`,
		[]byte(prefix), []byte(limit))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]byte)
	for rows.Next() {
		var key []byte
		var value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		result[string(key)] = value
	}

	return result, rows.Err()
}

func (s *sqliteDB) SnapshotToStream(prefix []byte, streamfunc func(bean *stub.SnapshotBean) bool) error {
	defer util.Recover()

	var query string
	var args []interface{}

	if prefix != nil {
		query = `
			SELECT key, value FROM kv_store 
			WHERE key >= ? AND key < ?
			ORDER BY key`
		args = []interface{}{prefix, s.nextPrefix(prefix)}
	} else {
		query = "SELECT key, value FROM kv_store ORDER BY key"
		args = []interface{}{}
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	rows, err := tx.Query(query, args...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var key, value []byte
		if err := rows.Scan(&key, &value); err != nil {
			return err
		}

		bean := &stub.SnapshotBean{}
		bean.Copy(key, value)

		if !streamfunc(bean) {
			break
		}
	}

	return rows.Err()
}

func (s *sqliteDB) LoadSnapshotBeans(beans ...*stub.SnapshotBean) error {
	defer util.Recover()

	if beans == nil || len(beans) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO kv_store (key, value) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, bean := range beans {
		if bean == nil {
			continue
		}
		if _, err := stmt.Exec(bean.Key, bean.Value); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s *sqliteDB) LoadSnapshotBean(bean *stub.SnapshotBean) error {
	defer util.Recover()

	if bean == nil {
		return nil
	}

	return s.Put(bean.Key, bean.Value)
}
