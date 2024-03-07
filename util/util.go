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
	"crypto/tls"
	"encoding/gob"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"runtime/debug"
	"sort"
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
	return time.Unix(0, tt).Format("2006-01-02 15:04:05")
}

func ArraySubElement[K int | int8 | int32 | int64 | string](a []K, k K) (_r []K) {
	if a != nil {
		_r = make([]K, 0)
		for _, v := range a {
			if v != k {
				_r = append(_r, v)
			}
		}
	}
	return
}

func ArraySub[K int | int8 | int32 | int64 | string](a1, a2 []K) (_r []K) {
	_r = make([]K, 0)
	if a1 != nil && a2 != nil {
		m := make(map[K]byte, 0)
		for _, a := range a2 {
			m[a] = 0
		}
		for _, a := range a1 {
			if _, ok := m[a]; !ok {
				_r = append(_r, a)
			}
		}
	} else if a2 == nil {
		return a1
	}
	return
}

func Recover() {
	if err := recover(); err != nil {
		logging.Error(string(debug.Stack()))
	}
}

func ContainStrings(li []string, v string) (b bool) {
	if li == nil {
		return false
	}
	sort.Strings(li)
	idx := sort.SearchStrings(li, v)
	if idx < len(li) {
		b = li[idx] == v
	}
	return
}

func ContainInt[T int64 | uint64 | int | uint | uint32 | int32](li []T, v T) (b bool) {
	if li == nil {
		return false
	}
	sort.Slice(li, func(i, j int) bool { return li[i] < li[j] })
	idx := sort.Search(len(li), func(i int) bool { return li[i] >= v })
	if idx < len(li) {
		b = li[idx] == v
	}
	return
}

func HttpPost(bs []byte, close bool, httpurl string) (_r []byte, err error) {
	tr := &http.Transport{DisableKeepAlives: true}
	if strings.HasPrefix(httpurl, "https:") {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}
	client := http.Client{Transport: tr}
	var req *http.Request
	if req, err = http.NewRequest(http.MethodPost, httpurl, bytes.NewReader(bs)); err == nil {
		if close {
			req.Close = true
		}
		var resp *http.Response
		if resp, err = client.Do(req); err == nil {
			if close {
				defer resp.Body.Close()
			}
			var body []byte
			if body, err = io.ReadAll(resp.Body); err == nil {
				_r = body
			}
		}
	}
	return
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