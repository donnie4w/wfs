/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package conf

import (
	"flag"
)

var _version_ = "0.0.1"

var CF = new(confbean)

type confbean struct {
	Port        int
	Chiplimit   int32
	MaxDataSize int64
	Data        string
	MaxFileSize int64
}

func ParseFlag() {
	CF.ParseFlag()
}

func (this *confbean) ParseFlag() {
	flag.IntVar(&this.Port, "p", 3434, "service port")
	flag.Int64Var(&this.MaxDataSize, "MD", int64(1)<<30, "max data size")
	flag.Int64Var(&this.MaxFileSize, "max", int64(1)<<24, "max uploadfile size")
	flag.StringVar(&this.Data, "data", "data/fsdb", "data path")
}
