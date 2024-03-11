// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package sys

import (
	"crypto/tls"
	"strings"

	"github.com/donnie4w/wfs/util"
	"golang.org/x/net/websocket"
)

type WS struct {
	ws     *websocket.Conn
	err    error
	closed bool
}

func (t *WS) Send(bs []byte) (err error) {
	defer util.Recover()
	if err = websocket.Message.Send(t.ws, bs); err != nil {
		t.err = err
	}
	return
}

func (t *WS) Receive(procData func([]byte) bool) (err error) {
	defer util.Recover()
	for {
		if t.err != nil || t.closed {
			break
		}
		var byt []byte
		if err = websocket.Message.Receive(t.ws, &byt); err == nil {
			if !procData(byt) {
				break
			}
		} else {
			t.err = err
			break
		}
	}
	return
}

func (t *WS) Close() error {
	defer util.Recover()
	t.closed = true
	return t.ws.Close()
}

func NewWfsClient(server, origin, name, password string) (_r *WS, err error) {
	var config *websocket.Config
	if config, err = websocket.NewConfig(server, origin); err == nil {
		if strings.HasPrefix(server, "wss://") {
			config.TlsConfig = &tls.Config{InsecureSkipVerify: true}
		}
		config.Header.Add("username", name)
		config.Header.Add("password", password)
		if conn, err := websocket.DialConfig(config); err == nil {
			_r = &WS{ws: conn}
		} else {
			return nil, err
		}
	}
	return
}
