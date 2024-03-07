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
	Opaddr             string    `json:"opaddr"`
	WebAddr            string    `json:"webaddr"`
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
