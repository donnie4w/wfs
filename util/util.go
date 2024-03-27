// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
//

package util

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net/url"
	"os"
	"runtime/debug"
	"strings"
	"time"

	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/simplelog/logging"
	"google.golang.org/protobuf/proto"
)

func ParseAddr(addr string) (_r string, err error) {
	if _r = addr; !strings.Contains(_r, ":") {
		if goutil.MatchString("^[0-9]{4,5}$", addr) {
			_r = ":" + _r
		} else {
			err = errors.New("error format :" + addr)
		}
	}
	return
}

func PEncode(m proto.Message) ([]byte, error) {
	return proto.Marshal(m)
}

func PDecode(bs []byte, m proto.Message) (err error) {
	return proto.Unmarshal(bs, m)
}

func Encode(e any) (by []byte, err error) {
	buf := new(bytes.Buffer)
	enc := gob.NewEncoder(buf)
	err = enc.Encode(e)
	by = buf.Bytes()
	return
}

func Decode(buf []byte, e any) (err error) {
	decoder := gob.NewDecoder(bytes.NewReader(buf))
	err = decoder.Decode(e)
	return
}

func TimestrampFormat(tt int64) string {
	return time.Unix(0, tt).Format(time.DateTime)
}

func Recover() {
	if err := recover(); err != nil {
		logging.Error(string(debug.Stack()))
	}
}

func IsFileExist(path string) (_r bool) {
	if path != "" {
		_, err := os.Stat(path)
		_r = err == nil || os.IsExist(err)
	}
	return
}

func IsURL(str string) bool {
	u, err := url.ParseRequestURI(str)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

func CreateNodeId() int64 {
	bs := goutil.Int64ToBytes(goutil.RandId())
	b8 := goutil.CRC8(bs[:7])
	bs[7] = b8
	return goutil.BytesToInt64(bs)
}

func CheckNodeId(nodeId int64) bool {
	bs := goutil.Int64ToBytes(int64(nodeId))
	b8 := goutil.CRC8(bs[:7])
	return b8 == bs[7]
}

func OpenFile(fname string, flag int, perm os.FileMode) (file *os.File, err error) {
	if file, err = os.OpenFile(fname, flag, perm); err == nil {
		if fi, err := file.Stat(); err == nil {
			if fi.Mode().Perm() != perm {
				file.Chmod(perm)
			}
		}
	}
	return
}