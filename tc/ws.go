// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package tc

import (
	"strings"

	"github.com/donnie4w/wfs/keystore"
	"github.com/donnie4w/wfs/sys"
)

func init() {

}

func newAdminLocalWsClient(tls bool, opaddr, requri, name, pwd string) (ws *sys.WS, err error) {
	procol := "ws://"
	if tls {
		procol = "wss://"
	}
	if err = keystore.LoadAdmin(sys.WFSDATA + "/logs"); err != nil {
		return
	}
	if opaddr == "" {
		if v, ok := keystore.Admin.GetOther("webaddr"); ok {
			if strings.HasPrefix(v, ":") {
				opaddr = "127.0.0.1" + v
			} else {
				opaddr = v
			}
		}
	}
	if name == "" {
		ss := keystore.Admin.AdminList()
		if len(ss) > 0 {
			name = ss[0]
			if ub, ok := keystore.Admin.GetAdmin(name); ok {
				pwd = ub.Pwd
			}
		}
	}
	ws, err = sys.NewWfsClient(procol+opaddr+requri, "http://127.0.0.1", name, pwd)
	return
}
