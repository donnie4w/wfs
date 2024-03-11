// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package sys

import (
	"fmt"

	. "github.com/donnie4w/wfs/stub"
)

var ERR_OVERSIZE = err(4101, "data oversize")
var ERR_NOPASS = err(4102, "verification fail")
var ERR_PARAMS = err(4103, "parameter error")
var ERR_EXSIT = err(4104, "has exist")
var ERR_AUTH = err(4105, "limited authority")
var ERR_NOTEXSIT = err(4106, "not exist")
var ERR_NEWPATHEXIST = err(4107, "the new path already exists")

var ERR_UNDEFINED = err(5101, "undefined error")
var ERR_DEFRAG_FORBID = err(5102, "operation files cannot be defragmented")
var ERR_DEFRAG_UNDERWAY = err(5103, "defragmentation is underway")
var ERR_STOPSERVICE = err(5104, "service has stopped")
var ERR_FILEAPPEND = err(5105, "append data error")
var ERR_FILECREATE = err(5106, "create file error")

type ERROR interface {
	WfsError() *WfsError
	Error() error
	Equal(ERROR) bool
	Code() int32
}

type wfserror struct {
	code int32
	info string
}

func err(code int32, info string) ERROR {
	return &wfserror{code, info}
}

func (t *wfserror) WfsError() *WfsError {
	return &WfsError{Code: &t.code, Info: &t.info}
}

func (t *wfserror) Error() error {
	return fmt.Errorf("code:%d,info:%s", t.code, t.info)
}

func (t *wfserror) Equal(e ERROR) bool {
	return t.code == e.Code()
}

func (t *wfserror) Code() int32 {
	return t.code
}
