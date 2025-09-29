// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type MigrationState struct {
	LastKey    []byte `json:"last_key"`
	TotalCount int64  `json:"total_count"`
	StartTime  int64  `json:"start_time"`
}

var batchSize = 10000

func main() {
	length := len(os.Args)
	if length != 3 && length != 4 {
		fmt.Printf("Usage: %s <leveldb_directory> <dest_db_file> <batch_size (default: 10000)>\n", os.Args[0])
		fmt.Println("Example: migration.exe /wfsdata/wfsdb /wfsdata/wfsdb/wfs.db")
		fmt.Println("Example 2: migration.exe /wfsdata/wfsdb /wfsdata/wfsdb/wfs.db 10000")
		os.Exit(1)
	}

	leveldbDir := os.Args[1]
	liteFile := os.Args[2]

	if length == 4 {
		if parsed, err := strconv.Atoi(os.Args[3]); err == nil && parsed > 0 {
			batchSize = parsed
		} else {
			batchSize = 10000
		}
	}

	log.Printf("Starting migration: LevelDB(%s) -> SQLite(%s)", leveldbDir, liteFile)

	if err := migrateData(leveldbDir, liteFile); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Migration completed successfully!")
}

func migrateData(leveldbDir, liteFile string) error {
	// Open LevelDB
	ldb, err := leveldb.OpenFile(leveldbDir, &opt.Options{})
	if err != nil {
		return fmt.Errorf("failed to open LevelDB: %v", err)
	}
	defer ldb.Close()

	// Open SQLite database
	db, err := sql.Open("sqlite3", liteFile)
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Ensure table exists
	if err := createSQLiteTable(db); err != nil {
		return err
	}

	// Load migration state
	state, err := loadMigrationState()
	if err != nil {
		return err
	}

	if state == nil {
		state = &MigrationState{
			StartTime: time.Now().Unix(),
		}
		log.Println("Starting new migration task")
	} else {
		log.Printf("Resuming from checkpoint: %d records migrated", state.TotalCount)
		if len(state.LastKey) > 0 {
			log.Printf("Resuming from key: %s", string(state.LastKey))
		}
	}

	// Perform migration
	return doMigration(ldb, db, state)
}

func createSQLiteTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS kv_store (
			key BLOB PRIMARY KEY,
			value BLOB NOT NULL
		)
	`)
	return err
}

func loadMigrationState() (*MigrationState, error) {
	data, err := os.ReadFile("migration_state.json")
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read state file: %v", err)
	}

	var state MigrationState
	if err := json.Unmarshal(data, &state); err != nil {
		return nil, fmt.Errorf("failed to parse state file: %v", err)
	}

	return &state, nil
}

func saveMigrationState(state *MigrationState) error {
	data, err := json.Marshal(state)
	if err != nil {
		return err
	}
	return os.WriteFile("migration_state.json", data, 0644)
}

func doMigration(ldb *leveldb.DB, db *sql.DB, state *MigrationState) error {
	startTime := time.Now()
	lastSaveTime := startTime

	for {
		keys, lastKey, hasMore, err := queryKeysWithPagination(ldb, state.LastKey, batchSize)
		if err != nil {
			return fmt.Errorf("failed to query keys: %v", err)
		}

		if len(keys) == 0 {
			break
		}

		successCount, err := processBatch(ldb, db, keys)
		if err != nil {
			return fmt.Errorf("failed to process batch: %v", err)
		}

		state.TotalCount += int64(successCount)
		if lastKey != nil {
			state.LastKey = make([]byte, len(lastKey))
			copy(state.LastKey, lastKey)
		}

		elapsed := time.Since(startTime)
		speed := float64(state.TotalCount) / elapsed.Seconds()
		log.Printf("Progress: %d records, Speed: %.0f records/sec, Last key: %s",
			state.TotalCount, speed, truncateKey(string(state.LastKey)))

		// Save state every minute
		if time.Since(lastSaveTime) >= time.Minute {
			if err := saveMigrationState(state); err != nil {
				log.Printf("Warning: failed to save migration state: %v", err)
			} else {
				log.Printf("Migration state saved: %d records migrated", state.TotalCount)
			}
			lastSaveTime = time.Now()
		}

		if !hasMore {
			break
		}
	}

	if err := saveMigrationState(state); err != nil {
		log.Printf("Warning: failed to save final state: %v", err)
	}

	if err := os.WriteFile("migration_complete.txt",
		[]byte(fmt.Sprintf("Migration completed at: %s, Total records: %d",
			time.Now().Format("2006-01-02 15:04:05"), state.TotalCount)), 0644); err != nil {
		log.Printf("Warning: failed to write completion marker: %v", err)
	}

	log.Printf("Migration completed! Total records: %d, Duration: %v", state.TotalCount, time.Since(startTime))
	return nil
}

func queryKeysWithPagination(ldb *leveldb.DB, startKey []byte, pageSize int) (keys [][]byte, lastKey []byte, hasMore bool, err error) {
	iter := ldb.NewIterator(nil, nil)
	defer iter.Release()

	if startKey != nil {
		iter.Seek(startKey)
	}

	for iter.Next() && len(keys) < pageSize {
		keyCopy := make([]byte, len(iter.Key()))
		copy(keyCopy, iter.Key())
		keys = append(keys, keyCopy)
	}

	hasMore = iter.Valid()
	if hasMore {
		lastKey = iter.Key()
	} else if len(keys) > 0 {
		lastKey = keys[len(keys)-1]
	}

	if err := iter.Error(); err != nil {
		return nil, nil, false, err
	}

	return keys, lastKey, hasMore, nil
}

func processBatch(ldb *leveldb.DB, db *sql.DB, keys [][]byte) (int, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("INSERT OR REPLACE INTO kv_store (key, value) VALUES (?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	successCount := 0
	for _, key := range keys {
		value, err := ldb.Get(key, nil)
		if err != nil {
			log.Printf("Warning: failed to get value for key=%s: %v", string(key), err)
			continue
		}
		if _, err := stmt.Exec(key, value); err != nil {
			log.Printf("Warning: failed to insert key=%s into DB: %v", string(key), err)
			continue
		}

		successCount++
	}

	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return successCount, nil
}

func truncateKey(key string) string {
	if len(key) > 50 {
		return key[:47] + "..."
	}
	return key
}
