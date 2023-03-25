/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package db

import (
	"bytes"
	"fmt"

	"github.com/donnie4w/simplelog/logging"

	"os"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

var logger = logging.NewLogger().SetFormat(logging.FORMAT_DATE | logging.FORMAT_TIME)
var _recover_ = "_recover_"

type DB struct {
	db     *leveldb.DB
	dbname string
	bak    bool
	f      *os.File
}

func NewDB(dbname string, backup bool) (db *DB) {
	db = new(DB)
	db.dbname = dbname
	db.bak = backup
	db.openDB()
	if !db._isRecover() {
		db.recover()
	}
	if backup {
		db.openbackupFile()
	}
	return
}

func (this *DB) openDB() {
	o := &opt.Options{
		Filter: filter.NewBloomFilter(10),
	}
	var err error
	this.db, err = leveldb.OpenFile(this.dbname, o)
	if err != nil {
		logger.Error("system start error:", err.Error())
		os.Exit(0)
	}
}

func (this *DB) DBExist(key []byte) (b bool) {
	b, _ = this.db.Has(key, nil)
	return
}

func (this *DB) Put(key, value []byte) (err error) {
	err = this.db.Put(key, value, nil)
	if this.bak {
		this.backup(key, value)
	}
	return
}

func (this *DB) Get(key []byte) (value []byte, err error) {
	value, err = this.db.Get(key, nil)
	return
}

func (this *DB) Del(key []byte) (err error) {
	err = this.db.Delete(key, nil)
	if this.bak {
		this.backup(key, nil)
	}
	return
}

func (this *DB) _isRecover() (b bool) {
	_, err := this.db.Get([]byte(_recover_), nil)
	if err == nil {
		b = true
	}
	return
}

func (this *DB) openbackupFile() (err error) {
	this.f, err = os.OpenFile(fmt.Sprint(this.dbname, "/WFS_BACKUP.LOG"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("open bakcp file error:", err.Error())
		os.Exit(1)
	}
	return
}

func (this *DB) backup(key, value []byte) (err error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	l := len(key)
	l2 := 0
	if value != nil {
		l += len(value)
		l2 = len(key)
	}
	buf.Write(octal2bytes(int32(l)))
	buf.Write(octal2bytes(int32(l2)))
	buf.Write(key)
	if value != nil {
		buf.Write(value)
	}
	this.f.Write(buf.Bytes())
	this.f.Sync()
	return
}

func (this *DB) recover() {
	var err error
	this.f, err = os.OpenFile(fmt.Sprint(this.dbname, "/WFS_BACKUP.LOG"), os.O_RDONLY, 0644)
	if err == nil {
		defer this.f.Close()
		for {
			var key, value []byte
			bs, err := readfile(this.f, 4)
			if err != nil {
				break
			}
			l := int(bytes2Octal(bs))
			bs, err = readfile(this.f, 4)
			l2 := int(bytes2Octal(bs))
			//			fmt.Println("==>", l, " ", l2)
			if l2 > 0 {
				key, err = readfile(this.f, l2)
				if err != nil {
					break
				}
				value, err = readfile(this.f, l-l2)
				if err != nil {
					break
				}
				this.db.Put(key, value, nil)
			} else {
				key, err = readfile(this.f, l)
				if err != nil {
					break
				}
				this.db.Delete(key, nil)
			}
		}
		this.db.Put([]byte(_recover_), []byte{0}, nil)
	}
}
func readfile(f *os.File, n int) (bs []byte, err error) {
	bs = make([]byte, n)
	_, err = f.Read(bs)
	return
}

func (this *DB) GetLike(prefix string) (datamap map[string][]byte) {
	iter := this.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)
	if iter != nil {
		datamap = make(map[string][]byte, 0)
		for iter.Next() {
			datamap[string(iter.Key())], _ = this.Get(iter.Key())
		}
		iter.Release()
	}
	return
}

func (this *DB) GetIterLimit(prefix string, limit string) (datamap map[string][]byte) {
	ran := new(util.Range)
	if prefix != "" {
		ran.Start = []byte(prefix)
	}
	if limit != "" {
		ran.Limit = []byte(limit)
	}
	if prefix == "" && limit == "" {
		ran = nil
	}
	iter := this.db.NewIterator(ran, nil)
	datamap = make(map[string][]byte, 0)
	for iter.Next() {
		data, err := this.db.Get(iter.Key(), nil)
		if err == nil {
			datamap[string(iter.Key())] = data
		}
	}
	iter.Release()
	return
}
func octal2bytes(row int32) (bs []byte) {
	bs = make([]byte, 0)
	for i := 0; i < 4; i++ {
		r := row >> uint((3-i)*4)
		bs = append(bs, byte(r))
	}
	return
}

func bytes2Octal(bb []byte) (value int32) {
	value = int32(0x0000)
	for i, b := range bb {
		ii := uint(b) << uint((3-i)*4)
		value = value | int32(ii)
	}
	return
}
