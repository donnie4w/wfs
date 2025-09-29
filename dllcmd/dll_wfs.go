// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package main

import "C"
import (
	"encoding/json"
	"github.com/donnie4w/go-logger/logger"
	_ "github.com/donnie4w/wfs/keystore"
	_ "github.com/donnie4w/wfs/level1"
	_ "github.com/donnie4w/wfs/stor"
	. "github.com/donnie4w/wfs/sys"
	_ "github.com/donnie4w/wfs/tc"
	"sync"
	"unsafe"
)

/*
#include <stdlib.h>
*/
import "C"

var (
	initMutex     sync.RWMutex
	isInitialized bool
	initError     string
	initOnce      sync.Once
)

// ServiceConfig 服务配置结构体
type ServiceConfig struct {
	EnableHTTP   *bool `json:"http"`   // 是否开启HTTP服务
	EnableThrift *bool `json:"thrift"` // 是否开启Thrift服务
	EnableAdmin  *bool `json:"admin"`  // 是否开启Admin服务
}

//export Init
func Init(configJSON *C.char) *C.char {
	initMutex.RLock()
	if isInitialized {
		initMutex.RUnlock()
		return nil
	}
	if initError != "" {
		errMsg := initError
		initMutex.RUnlock()
		return C.CString(errMsg)
	}
	defer initMutex.RUnlock()

	if configJSON != nil {
		jsonStr := C.GoString(configJSON)
		if jsonStr != "" {
			var config ServiceConfig
			if err := json.Unmarshal([]byte(jsonStr), &config); err != nil {
				return C.CString("Failed to parse config JSON: " + err.Error())
			}
			if config.EnableAdmin != nil {
				if !*config.EnableAdmin {
					WEBADDR = ""
				}
			}
			if config.EnableHTTP != nil {
				if !*config.EnableHTTP {
					LISTEN = 0
				}
			}
			if config.EnableThrift != nil {
				if !*config.EnableThrift {
					OPADDR = ""
				}
			}
		}
	}

	initOnce.Do(func() {
		go func() {
			isInitialized = true
			if err := Wfs.Serve(); err != nil {
				initError = "Initialization failed: " + err.Error()
				logger.Error("init wfs server error:", err)
			} else {
				logger.Info("WFS initialization completed")
			}
		}()
	})

	if initError != "" {
		return C.CString(initError)
	}
	return nil
}

//export IsInit
func IsInit() C.int {
	initMutex.RLock()
	defer initMutex.RUnlock()

	if isInitialized {
		return 1
	}
	return 0
}

//export GetInitStatus
func GetInitStatus() *C.char {
	initMutex.RLock()
	defer initMutex.RUnlock()

	if initError != "" {
		return C.CString(initError)
	}
	if isInitialized {
		return C.CString("initialized")
	}
	return C.CString("initializing")
}

//export Append
func Append(name *C.char, data *C.uchar, dataLen C.int, compress C.int) *C.char { // 改为 C.uchar
	initMutex.RLock()
	if !isInitialized {
		initMutex.RUnlock()
		return C.CString("uninitialized")
	}
	initMutex.RUnlock()

	if name == nil || data == nil || dataLen <= 0 {
		return C.CString(ERR_PARAMS.Error().Error())
	}

	goName := C.GoString(name)
	goData := C.GoBytes(unsafe.Pointer(data), dataLen) // 使用unsafe.Pointer转换

	if int64(len(goData)) > DataMaxsize {
		return C.CString(ERR_OVERSIZE.Error().Error())
	}

	if _, err := AppendData(goName, goData, int32(compress)); err != nil {
		return C.CString(err.Error().Error())
	}

	return nil
}

//export Delete
func Delete(path *C.char) *C.char {
	initMutex.RLock()
	if !isInitialized {
		initMutex.RUnlock()
		return C.CString("uninitialized")
	}
	initMutex.RUnlock()

	if path == nil {
		return C.CString(ERR_PARAMS.Error().Error())
	}

	goPath := C.GoString(path)
	if err := DelData(goPath); err != nil {
		return C.CString(err.Error().Error())
	}

	return nil
}

//export Get
func Get(path *C.char, resultLen *C.int) *C.uchar {
	initMutex.RLock()
	if !isInitialized {
		initMutex.RUnlock()
		*resultLen = -1
		return nil
	}
	initMutex.RUnlock()

	if path == nil {
		*resultLen = 0
		return nil
	}

	goPath := C.GoString(path)
	data := GetData(goPath)

	if data == nil || len(data) == 0 {
		*resultLen = 0
		return nil
	}

	*resultLen = C.int(len(data))
	return (*C.uchar)(C.CBytes(data))
}

//export Rename
func Rename(path *C.char, newpath *C.char) *C.char {
	initMutex.RLock()
	if !isInitialized {
		initMutex.RUnlock()
		return C.CString("uninitialized")
	}
	initMutex.RUnlock()

	if path == nil || newpath == nil {
		return C.CString(ERR_PARAMS.Error().Error())
	}

	goPath := C.GoString(path)
	goNewPath := C.GoString(newpath)

	if err := Modify(goPath, goNewPath); err != nil {
		return C.CString(err.Error().Error())
	}

	return nil
}

//export Has
func Has(path *C.char) C.int {
	initMutex.RLock()
	if !isInitialized {
		initMutex.RUnlock()
		return -1 //uninitialized
	}
	initMutex.RUnlock()
	if path == nil {
		return 0
	}
	key := C.GoString(path)
	exists := Contains(key)
	if exists {
		return 1
	}
	return 0
}

//export GetKeys
func GetKeys(fromId C.longlong, limit C.int) *C.char {
	initMutex.RLock()
	if !isInitialized {
		initMutex.RUnlock()
		return C.CString(`{"error":"uninitialized"}`)
	}
	initMutex.RUnlock()

	pbs := SearchLimit(int64(fromId), int64(limit))
	if pbs == nil {
		return C.CString(`{"keys":[]}`)
	}

	type KeyInfo struct {
		Name string `json:"name"`
		ID   int64  `json:"id"`
	}

	keys := make([]KeyInfo, len(pbs))
	for i, pb := range pbs {
		keys[i] = KeyInfo{
			Name: pb.Path,
			ID:   pb.Id,
		}
	}

	result := map[string]interface{}{
		"keys": keys,
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		return C.CString(`{"error":"json marshal error"}`)
	}

	return C.CString(string(jsonData))
}

//export FreeMemory
func FreeMemory(ptr unsafe.Pointer) {
	if ptr != nil {
		C.free(ptr)
	}
}

//export FreeString
func FreeString(s *C.char) {
	if s != nil {
		C.free(unsafe.Pointer(s))
	}
}

//export Close
func Close() {
	initMutex.Lock()
	defer initMutex.Unlock()

	if isInitialized {
		isInitialized = false
		initError = ""
		Wfs.Close()
		logger.Info("WFS closed")
	}
}

func main() {}
