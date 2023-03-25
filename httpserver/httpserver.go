/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package httpserver

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"
	. "wfs/conf"
	"wfs/storge"

	"wfs/conf"

	"github.com/donnie4w/simplelog/logging"
	"github.com/donnie4w/tlnet"
)

var logger = conf.Logger

type CmdType string

const (
	setWeight   CmdType = "setWeight"
	addSlave    CmdType = "addSlave"
	cutOff      CmdType = "cutOff"
	removeSlave CmdType = "removeSlave"
	slavelist   CmdType = "slavelist"
	ping        CmdType = "ping"
)
const (
	_200 = "200" //ok
	_403 = "403" //forbidden
	_404 = "404" //no found
	_500 = "500" //err
)

func read(hc *tlnet.HttpContext) {
	bs, err := GetData(hc.Request().RequestURI)
	if err == nil {
		hc.ResponseBytes(0, bs)
	} else {
		hc.Writer().WriteHeader(404)
	}
}

func upload(hc *tlnet.HttpContext) {
	defer myRecover()
	hc.MaxBytesReader(CF.MaxFileSize)
	file, handler, err := hc.FormFile("file")
	if err != nil {
		hc.ResponseString(err.Error())
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	bs := buf.Bytes()
	if len(bs) > 0 {
		name := ""
		uri := hc.Request().RequestURI
		if len(uri) > 2 {
			name = uri[3:]
		}
		if name == "" {
			name = handler.Filename
		}
		err := AppendData(bs, name, "")
		if err == nil {
			hc.ResponseString(fmt.Sprint("ok:", len(bs), " ", name))
		} else {
			hc.ResponseString(err.Error())
		}
	}
}

func del(hc *tlnet.HttpContext) {
	var name string
	uri := hc.Request().RequestURI
	if len(uri) > 2 {
		name = uri[3:]
	}
	if name != "" {
		err := DelData(name)
		if err == nil {
			hc.ResponseString(fmt.Sprint("delete ok,", name))
		} else {
			hc.ResponseString(fmt.Sprint("delete error,", name, " | ", err.Error()))
		}
	}
}

func check(hc *tlnet.HttpContext) {
	r := hc.Request()
	if r.Body != nil {
		defer r.Body.Close()
		bs, _ := ioutil.ReadAll(r.Body)
		name := string(bs)
		b := storge.Exsit(name)
		if b {
			hc.ResponseString("1")
		} else {
			hc.ResponseString("0")
		}
	}
}

/********************************************************************************/
func savePort(port int) {
	f, err := os.OpenFile(fmt.Sprint(CF.FileData, "/CURRENT"), os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	if err == nil {
		_, err = f.WriteString(fmt.Sprint(port))
		if err != nil {
			os.Exit(1)
		}
	} else {
		os.Exit(1)
	}
}

func readPort() (port int) {
	bs, err := ioutil.ReadFile(fmt.Sprint(CF.FileData, "/CURRENT"))
	if err == nil {
		port, _ = strconv.Atoi(string(bs))
	}
	return
}

/********************************************************************************/
func Start() {
	if IsCmd() {
		CF.Port = readPort()
		switch {
		case Cmd.AddSlave != "":
			wfsCmd(addSlave, Cmd.AddSlave)
		case Cmd.CutOff:
			wfsCmd(cutOff, "1")
		case Cmd.SetWeight != "":
			wfsCmd(setWeight, Cmd.SetWeight)
		case Cmd.RemoveSlave != "":
			wfsCmd(removeSlave, Cmd.RemoveSlave)
		case Cmd.Slavelist:
			wfsCmd(slavelist, "1")
		case Cmd.Ping != "":
			wfsCmd(ping, Cmd.Ping)
		}
	} else {

		if CF.Pprof > 0 {
			go http.ListenAndServe(fmt.Sprint(CF.Bind, ":", CF.Pprof), nil)
		}
		storge.Init()
		savePort(CF.Port)
		Init()
		tl := tlnet.NewTlnet()
		tl.AddHandlerFunc("/thrift", nil, thandler)
		tl.POST("/u/", upload)
		tl.GET("/r/", read)
		tl.POST("/c", check)
		tl.POST("/u", upload)
		tl.DELETE("/d/", del)
		tl.ReadTimeout(time.Duration(CF.ServerReadTimeout) * time.Second)
		logger.Info("wfs start,listen:", CF.Port)
		err := tl.HttpStart(fmt.Sprint(CF.Bind, ":", CF.Port))
		if err != nil {
			logger.Error("httpserver start error:", err.Error())
			os.Exit(1)
		}
	}
}

func myRecover() {
	if er := recover(); er != nil {
		logger.Error(er)
	}
}

//----------------------------------------------------------------------------------------------

// curl -F "file=@015.jpg" "http://127.0.0.1:3434/u/015.jpg"
// curl -X DELETE "http://127.0.0.1:3434/d/aaa/1.jpg"
// curl -I -X GET "http://127.0.0.1:3434/r/aaa/1.jpg"
