// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package tc

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/donnie4w/gofer/image"
	"github.com/donnie4w/tlnet"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
)

type clientService struct {
	isClose bool
	tln     *tlnet.Tlnet
}

var clientservice = &clientService{false, tlnet.NewTlnet()}
var images *image.Image

func (t *clientService) Serve() (err error) {
	if sys.LISTEN > 0 {
		err = t._serve(fmt.Sprint(":", sys.LISTEN), sys.Conf.Ssl_crt, sys.Conf.Ssl_crt_key)
	}
	return
}

func (t *clientService) Close() (err error) {
	defer util.Recover()
	if sys.LISTEN > 0 {
		t.isClose = true
		err = t.tln.Close()
	}
	return
}

func (t *clientService) _serve(addr string, serverCrt, serverKey string) (err error) {
	defer util.Recover()
	if addr, err = util.ParseAddr(addr); err != nil {
		return
	}
	t.tln.Handle("/", loadHandler)
	if serverCrt != "" && serverKey != "" {
		sys.FmtLog("client start tls [", addr, "]")
		err = t.tln.HttpStartTLS(addr, serverCrt, serverKey)
	}
	if !t.isClose {
		sys.FmtLog("client listen [", addr, "]")
		err = t.tln.HttpStart(addr)
	}
	if !t.isClose && err != nil {
		fmt.Println("client start failed:", err.Error())
		sys.Wfs.Close()
		os.Exit(1)
	}
	if t.isClose {
		err = nil
	}
	return
}

func loadHandler(hc *tlnet.HttpContext) {
	if bs, _ := getData(hc.Request().RequestURI); bs != nil {
		hc.ResponseBytes(0, bs)
	} else {
		hc.Writer().WriteHeader(404)
	}
}

func getData(uri string) (retbs []byte, err sys.ERROR) {
	if len(uri) > 1 {
		return getDataByName(uri[1:])
	}
	return
}

func getDataByName(uri1 string) (bs []byte, err sys.ERROR) {
	path, argstr := uri1, ""
	if index := strings.Index(uri1, "?"); index > 0 {
		path = uri1[:index]
		argstr = uri1[index:]
	}
	if decoded, err := url.QueryUnescape(path); err == nil {
		path = decoded
	}
	bs = sys.GetData(path)
	if bs != nil {
		if isImageMode(argstr) {
			if iv2, err := parseUriToImagemode(argstr); err == nil {
				if bss, err := images.Encode(bs, iv2.width, iv2.height, image.Mode(iv2.mode), iv2.getOptions()); err == nil {
					bs = bss
				}
			}
		}
	} else {
		if sys.Conf.SLASH && uri1[0] != '/' {
			return getDataByName("/" + uri1)
		}
		err = sys.ERR_NOTEXSIT
	}
	return
}

func isImageMode(s string) bool {
	i := strings.Index(s, "/")
	if i <= 0 {
		return false
	}
	switch s[:i] {
	case sys.IMAGEMODE, sys.IMAGEVIEW, sys.IMAGEVIEW2:
		return true
	default:
		return false
	}
}
