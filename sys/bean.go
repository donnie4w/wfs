// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package sys

type Server interface {
	Serve() (err error)
	Close() (err error)
}

type ConfBean struct {
	FileSize           int64     `json:"filesize"`
	Opaddr             *string   `json:"opaddr"`
	WebAddr            *string   `json:"webaddr"`
	Listen             int       `json:"listen"`
	Admin_Ssl_crt      string    `json:"admin.ssl_certificate"`
	Admin_Ssl_crt_key  string    `json:"admin.ssl_certificate_key"`
	Ssl_crt            string    `json:"ssl_certificate"`
	Ssl_crt_key        string    `json:"ssl_certificate_key"`
	Memlimit           int64     `json:"memlimit"`
	DataMaxsize        int64     `json:"data.maxsize"`
	Init               bool      `json:"init"`
	Keystore           *string   `json:"keystore"`
	Mode               *int      `json:"mode"`
	Sync               *bool     `json:"sync"`
	Compress           *int32    `json:"compress"`
	WfsData            *string   `json:"data.dir"`
	SLASH              bool      `json:"prefix.slash"`
	MaxSigma           float64   `json:"maxsigma"`
	MaxSide            int       `json:"maxside"`
	MaxPixel           int       `json:"maxpixel"`
	Resample           int8      `json:"resample"`
	ImgViewingRevProxy string    `json:"imgViewingRevProxy"`
	FileHash           *int      `json:"filehash"`
	AdminUserName      *string   `json:"adminusername"`
	AdminPassword      *string   `json:"adminpassword"`
	Restrict           *int      `json:"restrict"`
	DBType             int       `json:"dbtype"`
	DBConfig           *DBConfig `json:"db"`
}

type PathBean struct {
	Id         int64
	Path       string
	Body       []byte
	Timestramp int64
}

type FragBean struct {
	Node       string
	RmSize     int64
	ActualSize int64
	FileSize   int64
}

type openssl struct {
	PrivateBytes []byte
	PublicBytes  []byte
	PublicPath   string
	PrivatePath  string
}

// DBType 数据库类型
type DBType string

const (
	DBTypeLevelDB DBType = "leveldb"
	DBTypeSQLite  DBType = "sqlite"
)

// DBConfig 数据库配置
type DBConfig struct {
	Type    DBType
	Path    string
	LevelDB *LevelDBConfig `json:"leveldb,omitempty"`
	SQLite  *SQLiteConfig  `json:"sqlite,omitempty"`
}

type LevelDBConfig struct {
	BloomFilterBits int `json:"bloom_filter_bits,omitempty"` // 布隆过滤器位数
	OpenFilesCache  int `json:"open_files_cache,omitempty"`  // 打开文件缓存数量
	BlockCacheSize  int `json:"block_cache_size,omitempty"`  // 块缓存大小
	WriteBufferSize int `json:"write_buffer_size,omitempty"` // 写缓冲区大小
}

type SQLiteConfig struct {
	JournalMode string `json:"journal_mode,omitempty"` // WAL, DELETE, etc.
	SyncMode    string `json:"sync_mode,omitempty"`    // NORMAL, FULL, etc.
	CacheSize   int    `json:"cache_size,omitempty"`   // 缓存大小
	MMapSize    int    `json:"mmap_size,omitempty"`    // 内存映射大小
	BusyTimeout int    `json:"busy_timeout,omitempty"` // 繁忙超时(毫秒)
	AutoVacuum  string `json:"auto_vacuum,omitempty"`  // NONE, FULL, INCREMENTAL
	ForeignKeys bool   `json:"foreign_keys,omitempty"` // 外键约束
}

type istat interface {
	CReq() int64
	CReqDo()
	CReqDone()

	CPros() int64
	CProsDo()
	CProsDone()

	Tx() int64
	TxDo()
	TxDone()

	Ibs() int64
	Ib(int64)

	Obs() int64
	Ob(int64)
}
