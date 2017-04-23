/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package conf

import (
	"flag"
)

var _version_ = "0.0.2"

var CF = new(confbean)
var Cmd = new(cmdBean)

type confbean struct {
	Port              int
	Chiplimit         int32
	MaxDataSize       int64
	FileData          string
	MaxFileSize       int64
	Readonly          bool
	Keepalive         int
	Pprof             int
	ReadPerSecond     int64
	ServerReadTimeout int
	Bind              string
	Compress          bool
}
type cmdBean struct {
	SetWeight   string
	AddSlave    string
	RemoveSlave string
	Slavelist   bool
	CutOff      bool
	Ping        string
}

func ParseFlag() {
	CF.ParseFlag()
	Cmd.ParseFlag()
}

func (this *confbean) ParseFlag() {
	flag.IntVar(&this.Port, "p", 3434, "service port")
	flag.Int64Var(&this.MaxDataSize, "MD", int64(1)<<30, "max data size")
	flag.Int64Var(&this.MaxFileSize, "max", int64(1)<<24, "max uploadfile size")
	flag.StringVar(&this.FileData, "fileData", "data", "filedata path")
	flag.BoolVar(&this.Readonly, "readonly", false, "only be read, cannot be deleted or stored")
	flag.IntVar(&this.Keepalive, "keepalive", 0, "keepalive time(second)")
	flag.IntVar(&this.Pprof, "pprof", 0, "pprof port")
	flag.Int64Var(&this.ReadPerSecond, "rps", 60, "read per second")
	flag.IntVar(&this.ServerReadTimeout, "srt", 5, "server read timeout")
	flag.StringVar(&this.Bind, "bind", "0.0.0.0", "")
	flag.BoolVar(&this.Compress, "compress", false, "compress file")
}

func (this *cmdBean) ParseFlag() {
	flag.StringVar(&this.SetWeight, "setweight", "", "weight : -setweight nodename1:10")
	flag.StringVar(&this.AddSlave, "addslave", "", "add a slave node(nodename:ip:port) ： -addslave nodename1:192.168.1.100:3434")
	flag.StringVar(&this.RemoveSlave, "removeslave", "", "remove a slave node(nodename) ： -removeslave nodename1")
	flag.StringVar(&this.Ping, "ping", "", "-ping 127.0.0.1:4545")
	flag.BoolVar(&this.Slavelist, "slavelist", false, "get slave list")
	flag.BoolVar(&this.CutOff, "cut", false, "cut off the data file")
}

func IsCmd() bool {
	return Cmd.AddSlave != "" || Cmd.CutOff || Cmd.SetWeight != "" || Cmd.Slavelist || Cmd.RemoveSlave != "" || Cmd.Ping != ""
}
