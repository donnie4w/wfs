// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package sys

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"strconv"
	"sync/atomic"

	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/simplelog/logging"
)

func init() {
	Wfs = &server{}
}

func praseflag() {
	flag.StringVar(&DEBUGADDR, "debug", "", "debug address")
	flag.StringVar(&WFSJSON, "c", "wfs.json", "configuration file of wfs in json")
	flag.StringVar(&ORIGIN, "origin", "", "origin for websocket")
	flag.BoolVar(&LOGDEBUG, "log", false, "debug log on or off")
	flag.StringVar(&Service, "s", "", "services command")
	flag.StringVar(&Pid, "p", "", "path of wfs data or pid of wfs")
	flag.IntVar(&GOGC, "gc", -1, "a collection is triggered when the ratio of freshly allocated data")
	flag.Parse()
	parsec()

	if Conf.Mode != nil {
		Mode = *Conf.Mode
	}

	if Conf.Sync != nil {
		SYNC = *Conf.Sync
	}
	if Conf.WebAddr != "" {
		WEBADDR = Conf.WebAddr
	}
	if Conf.Listen > 0 {
		LISTEN = Conf.Listen
	}
	if Conf.Opaddr != "" {
		OPADDR = Conf.Opaddr
	}
	if Conf.DataMaxsize > 0 {
		DataMaxsize = Conf.DataMaxsize * KB
	}

	if Conf.Memlimit > 0 {
		Memlimit = Conf.Memlimit
	}

	if Conf.FileSize > 0 {
		FileSize = Conf.FileSize * MB
		if DataMaxsize > FileSize {
			DataMaxsize = FileSize
		}
	}

	if Conf.WfsData != nil {
		WFSDATA = *Conf.WfsData
	}

	if Conf.Compress != nil {
		CompressType = *Conf.Compress
	}

	if Conf.MaxSide > 0 {
		MaxSide = Conf.MaxSide
	}

	if Conf.MaxSigma > 0 {
		MaxSigma = Conf.MaxSigma
	}

	if Conf.MaxPixel > 0 {
		MaxPixel = Conf.MaxPixel
	}

	if Service == "stop" {
		if Pid == "" {
			if bs, err := goutil.ReadFile(WFSDATA + "/logs/wfs.pid"); err == nil {
				Pid = string(bs)
			}
		} else {
			if _, err := strconv.Atoi(Pid); err != nil {
				if bs, err := goutil.ReadFile(Pid + "/logs/wfs.pid"); err == nil {
					Pid = string(bs)
				}
			}
		}
		sendTerminated(Pid)
		os.Exit(0)
	}

	flag.Usage = usage
	flag.Usage()

	debug.SetMemoryLimit(Memlimit * MB)
	debug.SetGCPercent(GOGC)
	wfsmkdir(WFSDATA)
	log.SetRollingFile(WFSDATA+"/logs", "wfs.log", 100, logging.MB)
	if LOGDEBUG {
		logging.SetFormat(logging.FORMAT_DATE|logging.FORMAT_TIME|logging.FORMAT_SHORTFILENAME).SetRollingFile(WFSDATA+"/logs", "wfs.log", 100, logging.MB)
	} else {
		logging.SetLevel(logging.LEVEL_OFF)
	}
}

func usage() {
	exename := "wfs"
	if runtime.GOOS == "windows" {
		exename = "wfs.exe"
	}
	fmt.Fprintln(os.Stderr, `wfs version: wfs/`+VERSION+`
Usage: `+exename+`	
	-c: configuration file  e.g:  -c wfs.json
`)
}

func parsec() {
	if defaultConf != "" {
		Conf, _ = goutil.JsonDecode[*ConfBean]([]byte(defaultConf))
	} else {
		Conf = GetConfg()
	}
	if Conf == nil {
		FmtLog("empty config")
		Conf = &ConfBean{}
	}
}

func GetConfg() (conf *ConfBean) {
	if bs, err := goutil.ReadFile(WFSJSON); err == nil {
		conf, _ = goutil.JsonDecode[*ConfBean](bs)
	}
	return
}

func wfsmkdir(dir string) (err error) {
	if err = os.MkdirAll(dir+"/logs", 0777); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	if err = os.MkdirAll(dir+"/wfsdb", 0777); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	if err = os.MkdirAll(dir+"/wfsfile", 0777); err != nil {
		logging.Error(err)
		os.Exit(1)
	}
	return
}

type stat struct {
	creq  int64
	cpros int64
	tx    int64
	ibs   int64
	obs   int64
}

func (t *stat) CReq() int64 {
	return t.creq
}
func (t *stat) CReqDo() {
	atomic.AddInt64(&t.creq, 1)
}
func (t *stat) CReqDone() {
	atomic.AddInt64(&t.creq, -1)
}

func (t *stat) CPros() int64 {
	return t.cpros
}
func (t *stat) CProsDo() {
	atomic.AddInt64(&t.cpros, 1)
}
func (t *stat) CProsDone() {
	atomic.AddInt64(&t.cpros, -1)
}

func (t *stat) Tx() int64 {
	return t.tx
}
func (t *stat) TxDo() {
	atomic.AddInt64(&t.tx, 1)
}
func (t *stat) TxDone() {
	atomic.AddInt64(&t.tx, -1)
}

func (t *stat) Ibs() int64 {
	return t.ibs
}
func (t *stat) Ib(i int64) {
	atomic.AddInt64(&t.ibs, i)
}
func (t *stat) Obs() int64 {
	return t.obs
}
func (t *stat) Ob(i int64) {
	atomic.AddInt64(&t.obs, i)
}
