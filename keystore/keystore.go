// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
//

package keystore

import (
	"fmt"
	"github.com/donnie4w/gofer/util"
	"os"
	"strconv"

	"github.com/donnie4w/wfs/sys"

	"github.com/donnie4w/gofer/keystore"
)

func init() {
	sys.KeyStoreInit = Init
}

func Init(dir string) {
	if dir == "" {
		dir, _ = os.Getwd()
	}
	if err := InitAdmin(dir + "/logs"); err != nil {
		fmt.Println("keystore init failed:", err.Error())
		sys.Wfs.Close()
		os.Exit(1)
	}
	sys.FmtLog("wfs", sys.VERSION, " uuid[", sys.UUID, "]")
}

func InitAdmin(dir string) (err error) {
	if keystore.KeyStore, err = keystore.NewKeyStore(dir, "keystore.tdb"); err == nil {
		Admin.Load()
		if v, ok := Admin.GetOther("WFSUUID"); ok {
			id, _ := strconv.ParseUint(v, 10, 64)
			sys.UUID = int64(id)
		} else {
			sys.UUID = util.UUID64()
			Admin.PutOther("WFSUUID", fmt.Sprint(sys.UUID))
		}
	}
	return
}

func LoadAdmin(dir string) (err error) {
	if keystore.KeyStore, err = keystore.LoadKeyStore(dir, "keystore.tdb"); err == nil {
		Admin.Load()
	}
	return
}
