// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package tc

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/donnie4w/tlnet"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
)

func appendHandler(hc *tlnet.HttpContext) {
	defer util.Recover()

	file, handler, err := hc.FormFile("file")
	if err != nil {
		hc.ResponseString(err.Error())
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)
	bs := buf.Bytes()

	if int64(len(bs)) > sys.DataMaxsize {
		hc.ResponseString(`{"status":false, "desc":"` + sys.ERR_OVERSIZE.WfsError().GetInfo() + `"}`)
		return
	}

	name := hc.PostParam("filename")
	if name == "" {
		if uri := hc.Request().RequestURI; len(uri) > 8 {
			name = uri[8:]
		} else {
			name = handler.Filename
		}
	}

	if len(bs) > 0 && name != "" {
		if err := sys.AppendData(name, bs, sys.CompressType); err == nil {
			hc.ResponseString(`{"status":true, "name":"` + name + `","size":` + strconv.Itoa(len(bs)) + `}`)
		} else {
			hc.ResponseString(`{"status":false, "desc":"` + err.WfsError().GetInfo() + `"}`)
		}
	}
}

func deleteHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	if !strings.EqualFold(hc.ReqInfo.Method, http.MethodDelete) {
		hc.ResponseBytes(http.StatusMethodNotAllowed, nil)
		return
	}

	name := hc.PostParam("filename")
	if name == "" {
		if uri := hc.Request().RequestURI; len(uri) > 8 {
			name = uri[8:]
		}
	}

	if name != "" {
		if err := sys.DelData(name); err == nil {
			hc.ResponseString(`{"status":true, "name":"` + name + `"}`)
		} else {
			hc.ResponseString(`{"status":false, "desc":"` + err.WfsError().GetInfo() + `"}`)
		}
	}
}
