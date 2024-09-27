// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
//

package sys

import (
	"fmt"
	"time"

	"github.com/donnie4w/gofer/hashmap"
)

const VERSION = "1.0.7"

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

var (
	Serve          = hashmap.NewTreeMap[int, Server](5)
	Wfs            Server
	STARTTIME      = time.Now()
	UUID           int64
	Service        string
	LOGDEBUG       bool
	GOGC           int
	Pid            string
	ORIGIN         string
	DEBUGADDR      string
	WFSJSON        string
	Conf           *ConfBean
	Mode           = 1
	MaxSigma       = float64(20)
	MaxSide        = 10000
	MaxPixel       = 300000000
	CompressType   = int32(1)
	WFSDATA        = "wfsdata"
	IMAGEVIEW2     = "?imageView2"
	IMAGEVIEW      = "?imageView"
	IMAGEMODE      = "?mode"
	MD2HTML        = "?md2html"
	DefaultAccount = [2]string{"admin", "123"}
	WEBADDR        = fmt.Sprint(6<<10 + 2)
	OPADDR         = fmt.Sprint(":", 5<<10+2)
	LISTEN         = 4<<10 + 2
	SYNC           = false
	Restrict       = 95
	DBBUFFER       = 1 << 6 * MB
	InaccurateTime = time.Now().UnixNano()
	ConnectTimeout = 10 * time.Second
	WaitTimeout    = 10 * time.Second
	FileSize       = int64(1000 * MB)
	DataMaxsize    = FileSize / 5
	SocketTimeout  = 10 * time.Second
	OpenSSL        = &openssl{}
	Memlimit       = int64(1 << 10)
	FileHash       = 0
	defaultConf    = ""
	host           = ""
	user           = ""
	pwd            = ""
	out            = ""
	cover          = false
	extls          = false
	start          = int64(0)
	limit          = int64(0)
	efile          = false
	filegz         = false
	metaType       = byte(1)
	fileType       = byte(2)
	useOriginal    = byte(0)
	useZlib        = byte(1)
)
