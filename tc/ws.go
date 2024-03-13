// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package tc

import (
	"fmt"
	"strings"

	"github.com/donnie4w/wfs/keystore"
	"github.com/donnie4w/wfs/sys"
)

func newAdminLocalWsClient(tls bool, pid, opaddr, requri, name, pwd string) (ws *sys.WS, err error) {
	local := false
	if opaddr == "" {
		local = true
		if pid == "" {
			pid = sys.WFSDATA
		}
		if err = keystore.LoadAdmin(pid + "/logs"); err != nil {
			return nil, fmt.Errorf("wfs data directory cannot be parsed:%s", pid)
		}
		if v, ok := keystore.Admin.GetOther("webaddr"); ok {
			if strings.HasPrefix(v, ":") {
				opaddr = "127.0.0.1" + v
			} else {
				opaddr = v
			}
		}
		if sys.Conf.Admin_Ssl_crt != "" && sys.Conf.Admin_Ssl_crt_key != "" {
			tls = true
		} else {
			tls = false
		}

		ss := keystore.Admin.AdminList()
		if len(ss) > 0 {
			name = ss[0]
			if ub, ok := keystore.Admin.GetAdmin(name); ok {
				pwd = ub.Pwd
			}
		}
	}
	proto := "ws://"
	if tls {
		proto = "wss://"
	}

	if ws, err = sys.NewWfsClient(proto+opaddr+requri, "http://127.0.0.1", name, pwd); err != nil {
		if !local {
			err = fmt.Errorf("wfs service cannot be connected, host[%s],tls[%t],username[%s],password[%s]", opaddr, tls, name, pwd)
		}
	}
	return
}
