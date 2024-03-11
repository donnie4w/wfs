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
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

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

	flag.StringVar(&host, "host", "", "host")
	flag.StringVar(&user, "user", "", "user")
	flag.StringVar(&pwd, "pwd", "", "pwd")
	flag.StringVar(&out, "o", "", "path of metadata file")
	flag.BoolVar(&cover, "cover", false, "whether to overwrite the same path data")
	flag.BoolVar(&extls, "tls", false, "use tls")
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

	if Service != "" {
		praseService(Service)
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
	fmt.Fprintln(os.Stderr, `wfs version: wfs/`+VERSION+`
Usage: `+exename()+`	
	-c: configuration file  e.g: `+exename()+`  -c wfs.json
`)
}

func exename() string {
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprint(runtime.GOOS, strings.ReplaceAll(VERSION, ".", ""), "_wfs", ".exe")
	case "darwin":
		return fmt.Sprint("mac", strings.ReplaceAll(VERSION, ".", ""), "_wfs")
	default:
		return fmt.Sprint(runtime.GOOS, strings.ReplaceAll(VERSION, ".", ""), "_wfs")
	}
}

func parsec() {
	if defaultConf != "" {
		Conf, _ = goutil.JsonDecode[*ConfBean]([]byte(defaultConf))
	} else {
		Conf = GetConfg()
	}
	if Conf == nil {
		Conf = &ConfBean{}
	}
}

func praseService(s string) {
	defer os.Exit(0)
	switch s {
	case "stop":
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
	case "export":
		if ws, err := WsClient(extls, host, "/export", user, pwd); err == nil {
			ws.Send([]byte{1})
			if out != "" {
				if goutil.IsFileExist(out) {
					fmt.Println("file already exists:", out)
					return
				}
				os.MkdirAll(filepath.Dir(out), 0777)
			} else {
				out = "wfsmeta" + time.Now().Format("20060102150405")
			}
			if f, err := os.OpenFile(out, os.O_APPEND|os.O_CREATE|os.O_TRUNC, 0666); err == nil {
				defer f.Close()
				start := time.Now().Unix()
				ws.Receive(func(bs []byte) bool {
					if len(bs) == 1 {
						switch bs[0] {
						case 0:
							fmt.Println(time.Now().Format(time.DateTime)+"，export data >>", out, "(", time.Now().Unix()-start, "s)")
						case 1:
							f.Close()
							os.Remove(out)
							fmt.Println(ERR_NOPASS.WfsError().GetInfo())
							os.Exit(0)
						}
						return false
					}
					if _, err = f.Write(goutil.Int16ToBytes(int16(len(bs)))); err != nil {
						fmt.Println(err.Error())
						return false
					}
					if _, err = f.Write(bs); err != nil {
						fmt.Println(err.Error())
						return false
					}
					return true
				})
			} else {
				fmt.Println(err.Error())
			}
		} else {
			fmt.Println("export error:", err.Error())
		}
	case "import":
		if !goutil.IsFileExist(out) {
			fmt.Println("file does not exist:", out)
			return
		}
		var ws *WS
		var err error
		if ws, err = WsClient(extls, host, "/import", user, pwd); err == nil {
			start := time.Now().Unix()
			go ws.Receive(func(bs []byte) bool {
				if len(bs) == 1 {
					switch bs[0] {
					case 0:
						fmt.Println(time.Now().Format(time.DateTime)+"，import data >>", out, "(", time.Now().Unix()-start, "s)")
					case 1:
						fmt.Println(ERR_NOPASS.WfsError().GetInfo())
					case 2:
						fmt.Println("server error")
					}
				}
				os.Exit(0)
				return false
			})
			var c byte = 2
			if !cover {
				c = 3
			}
			if err = ws.Send([]byte{c}); err == nil {
				offset, n := 0, 0
				var f *os.File
				if f, err = os.OpenFile(out, os.O_RDONLY, 0666); err == nil {
					defer f.Close()
					fl, _ := f.Stat()
					for offset < int(fl.Size()) {
						head := make([]byte, 2)
						if n, err = f.ReadAt(head, int64(offset)); err != nil || n != 2 {
							break
						}
						length := goutil.BytesToInt16(head)
						body := make([]byte, length)
						if n, err = f.ReadAt(body, int64(offset+2)); err != nil || n != int(length) {
							break
						}
						if err = ws.Send(body); err != nil {
							break
						}
						offset += 2 + int(length)
					}
					if err != nil {
						goto ERR
					}
					err = ws.Send([]byte{1})
					<-time.After(5 * time.Second)
				}
			}
		ERR:
			if err != nil {
				fmt.Println("import error:", err.Error())
			}
		} else {
			fmt.Println("import error:", err.Error())
		}

	default:
		fmt.Printf("could not find service with: %s\n", s)
		os.Exit(1)
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
