// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
//

package sys

import (
	"github.com/donnie4w/wfs/stub"
)

var (
	KeyStoreInit func(string)
	Count        func() int64
	Seq          func() int64
	AppendData   func(string, []byte, int32) (int64, ERROR)
	GetData      func(string) []byte
	DelData      func(string) ERROR
	//Add            func([]byte, []byte) error
	//Del            func([]byte) error
	SearchLike     func(string) []*PathBean
	SearchLimit    func(int64, int64) []*PathBean
	Defrag         func(string) ERROR
	FragAnalysis   func(string) (*FragBean, ERROR)
	Export         func(func(*stub.SnapshotBean) bool) error
	ExportByINCR   func(int64, int64, func(*stub.SnapshotBeans) bool) ERROR
	ExportFile     func(int64, int64, func(*stub.SnapshotFile) bool) ERROR
	ExportByPaths  func([]string, func(snaps *stub.SnapshotBeans) bool) ERROR
	Import         func(*stub.SnapshotBean, bool) error
	ImportFile     func(*stub.SnapshotFile) ERROR
	Modify         func(string, string) ERROR
	WsClient       func(tls bool, pid, opaddr, requri, name, pwd string) (ws *WS, err error)
	IsEmptyBigFile func(string) bool
)
