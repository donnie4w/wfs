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
	if bs, cy, _ := getData(hc.Request().RequestURI); bs != nil {
		if cy != "" {
			hc.Writer().Header().Add("Content-Type", cy)
		}
		hc.ResponseBytes(0, bs)
	} else {
		hc.Writer().WriteHeader(404)
	}
}

func getData(uri string) (retbs []byte, t string, err sys.ERROR) {
	if len(uri) > 1 {
		return getDataByName(uri[1:])
	}
	return
}

func getDataByName(uri1 string) (bs []byte, ct string, err sys.ERROR) {
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
		if argstr != "" {
			m, o := getmode(argstr)
			if isImageMode(m) {
				if iv2, err := parseUriToImagemode(argstr); err == nil {
					if bss, err := images.Encode(bs, iv2.width, iv2.height, image.Mode(iv2.mode), iv2.getOptions()); err == nil {
						bs = bss
					}
				}
			} else {
				ct, _ = contentType(m, o)
			}
		} else if index := strings.LastIndex(path, "."); index > 0 {
			if suffix := path[index:]; len(suffix) > 1 {
				ct, _ = contentType(suffix, "")
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

func getmode(s string) (string, string) {
	i := strings.Index(s, "/")
	if i <= 0 {
		return s, ""
	} else {
		return s[:i], s[i+1:]
	}
}

func isImageMode(m string) bool {
	switch m {
	case sys.IMAGEMODE, sys.IMAGEVIEW, sys.IMAGEVIEW2:
		return true
	default:
		return false
	}
}

func contentType(m string, o string) (r string, b bool) {
	if m != "" {
		switch strings.ToLower(m[1:]) {
		case "js":
			r = "application/javascript"
		case "css":
			r = "text/css"
		case "markdown", "md":
			r = "text/markdown"
		case "json":
			r = "application/json"
		case "xml":
			r = "application/xml"
		case "xsl":
			r = "text/xsl"
		case "wsdl":
			r = "application/wsdl+xml"
		case "xsd":
			r = "application/xsd+xml"
		case "rss":
			r = "application/rss+xml"
		case "word", "doc":
			r = "application/msword"
		case "stream":
			r = "application/octet-stream"
		case "plain", "txt", "text":
			r = "text/plain"
		case "html":
			r = "text/html"
		case "jpeg", "jpg":
			r = "image/jpeg"
		case "png":
			r = "image/png"
		case "gif":
			r = "image/gif"
		case "mpeg":
			r = "audio/mpeg"
		case "tiff", "tif":
			r = "image/tiff"
		case "webp":
			r = "image/webp"
		case "ogg":
			r = "audio/ogg"
		case "wav":
			r = "audio/wav"
		case "mp3":
			r = "audio/mpeg"
		case "flac":
			r = "audio/flac"
		case "mp4", "m4a":
			r = "video/mp4"
		case "quicktime", "mov", "qt":
			r = "video/quicktime"
		case "webm":
			r = "video/webm"
		case "mixed":
			r = "multipart/mixed"
		case "docx":
			r = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
		case "utf8", "utf-8", "utf-16", "gbk", "gb18030", "iso-8859-1", "latin1", "big5", "windows-1251", "windows-1252", "shift_jis", "euc-kr", "us-ascii":
			r = charset(m[1:])
			b = true
		default:
			if len(m) > 14 && strings.EqualFold(m[0:14], "?content-type=") {
				r = m[14:]
				if o != "" {
					r = r + "/" + o
				}
				b = true
			}
		}
	}
	if o != "" && !b {
		if r != "" {
			return r + ";" + charset(o), b
		} else {
			return charset(o), b
		}
	}
	return
}

func charset(o string) string {
	if strings.HasPrefix(o, "charset=") {
		return o
	} else {
		return "charset=" + strings.ToUpper(o)
	}
}
