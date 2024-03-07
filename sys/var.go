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

const VERSION = "1.0.1"

const (
	KB = 1 << 10
	MB = 1 << 20
	GB = 1 << 30
)

var (
	Serve          = hashmap.NewSortMap[int, Server]()
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
	DefaultAccount = [2]string{"admin", "123"}
	WEBADDR        = fmt.Sprint(6<<10 + 2)
	OPADDR         = fmt.Sprint(":", 5<<10+2)
	LISTEN         = 4<<10 + 2
	SYNC           = false
	DBBUFFER       = 1 << 6 * MB
	InaccurateTime = time.Now().UnixNano()
	ConnectTimeout = 10 * time.Second
	WaitTimeout    = 10 * time.Second
	FileSize       = int64(500 * MB)
	DataMaxsize    = int64(100 * MB)
	SocketTimeout  = 10 * time.Second
	OpenSSL        = &openssl{}
	Memlimit       = int64(1 << 10)
	defaultConf    = ""
)
