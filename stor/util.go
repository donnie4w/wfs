// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package stor

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	goutil "github.com/donnie4w/gofer/util"
	"runtime"

	"github.com/donnie4w/gofer/base58"
	"github.com/donnie4w/gofer/compress"
	. "github.com/donnie4w/wfs/stub"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
	lru "github.com/hashicorp/golang-lru/v2"
)

func strToInt(s string) (u uint64, b bool) {
	u, b = base58.DecodeForInt64([]byte(s))
	return
}

func intToStr(u uint64) (s string) {
	return string(base58.EncodeForInt64(u))
}

func bytesToWfsFileBean(bs []byte) (wffb *WfsFileBean) {
	wffb = &WfsFileBean{}
	if util.PDecode(bs, wffb) != nil {
		wffb = nil
	}
	return
}

func wfsFileBeanToBytes(b *WfsFileBean) (bs []byte) {
	bs, _ = util.PEncode(b)
	return
}

func bytesToWfsNodeBean(bs []byte) (wnb *WfsNodeBean) {
	wnb = &WfsNodeBean{}
	if util.PDecode(bs, wnb) != nil {
		wnb = nil
	}
	return
}

func wfsNodeBeanToBytes(b *WfsNodeBean) (bs []byte) {
	bs, _ = util.PEncode(b)
	return
}

func bytesToWfsPathBean(bs []byte) (wpb *WfsPathBean) {
	wpb = &WfsPathBean{}
	if util.PDecode(bs, wpb) != nil {
		wpb = nil
	}
	return
}

func wfsPathBeanToBytes(b *WfsPathBean) (bs []byte) {
	bs, _ = util.PEncode(b)
	return
}

func praseCompress(bs []byte, compressType int32) (_r []byte) {
	switch compressType {
	case 1:
		return compress.Snappy(bs)
	case 2:
		_r, _ = compress.Zstd(bs)
	case 3, 4, 5, 6, 7, 8, 9, 10, 11:
		_r, _ = compress.ZlibLevel(bs, int(compressType)-2)
	default:
		_r = bs
	}
	return
}

func praseUncompress(bs []byte, compressType int32) (_r []byte) {
	switch compressType {
	case 1:
		_r, _ = compress.UnSnappy(bs)
	case 2:
		_r, _ = compress.UnZstd(bs)
	case 3, 4, 5, 6, 7, 8, 9, 10, 11:
		_r, _ = compress.UnZlib(bs)
	default:
		_r = bs
	}
	return
}

func fingerprint(bs []byte) []byte {
	switch sys.FileHash {
	case 0:
		return goutil.Int64ToBytes(int64(goutil.CRC64(bs)))
	case 1:
		hash := md5.Sum(bs)
		return hash[:]
	case 2:
		hash := sha1.Sum(bs)
		return hash[:]
	case 3:
		hash := sha256.Sum256(bs)
		return hash[:]
	default:
		return goutil.Int64ToBytes(int64(goutil.CRC64(bs)))
	}
}

func fileoffset() int {
	switch sys.FileHash {
	case 1:
		return 20
	case 2:
		return 24
	case 3:
		return 36
	default:
		return 12
	}
}

func fingerprintLen() int {
	switch sys.FileHash {
	case 1:
		return 16
	case 2:
		return 20
	case 3:
		return 32
	default:
		return 8
	}
}

var cache *lru.Cache[string, []byte]

func initcache() (err error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	freemem := uint64(sys.Memlimit*sys.MB) - m.TotalAlloc
	length := freemem / 128
	if length > 1<<20 {
		length = 1 << 20
	} else if length < 1<<10 {
		length = 1 << 10
	}
	cache, err = lru.New[string, []byte](int(length))
	return
}

func cacheGet(key []byte) (bs []byte, err error) {
	if cache != nil {
		if bs, ok := cache.Get(string(key)); ok {
			return bs, nil
		}
	}
	if bs, err = wfsdb.Get(key); err == nil {
		cachePut(key, bs)
	}
	return
}

func cachePut(key, value []byte) bool {
	if cache == nil {
		return false
	}
	return cache.Add(string(key), value)
}

func cacheDel(key []byte) bool {
	if cache == nil {
		return false
	}
	return cache.Remove(string(key))
}
