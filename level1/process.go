// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package level1

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/gothrift/thrift"
	. "github.com/donnie4w/wfs/stub"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
)

func init() {
	sys.Serve.Put(3, server)
}

var _transportFactory = thrift.NewTBufferedTransportFactory(1 << 13)
var _tcompactProtocolFactory = thrift.NewTCompactProtocolFactoryConf(&thrift.TConfiguration{})

var server = &service{}

type service struct {
	isClose         bool
	serverTransport thrift.TServerTransport
}

func (t *service) _server(_addr string, processor thrift.TProcessor, TLS bool, serverCrt, serverKey string) (err error) {
	if TLS {
		cfg := &tls.Config{}
		var cert tls.Certificate
		if cert, err = tls.LoadX509KeyPair(serverCrt, serverKey); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
			t.serverTransport, err = thrift.NewTSSLServerSocketTimeout(_addr, cfg, sys.SocketTimeout)
		}
	} else {
		t.serverTransport, err = thrift.NewTServerSocketTimeout(_addr, sys.SocketTimeout)
	}

	if err == nil && t.serverTransport != nil {
		server := thrift.NewTSimpleServer4(processor, t.serverTransport, nil, nil)
		if err = server.Listen(); err == nil {
			s := fmt.Sprint("operating services start[", _addr, "]")
			if TLS {
				s = fmt.Sprint("operating services start tls[", _addr, "]")
			}
			sys.FmtLog(s)
			for {
				if _transport, err := server.ServerTransport().Accept(); err == nil {
					go func() {
						defer util.Recover()
						cc := newCliContext(_transport)
						defer cc.close()
						defaultCtx := context.WithValue(context.Background(), "CliContext", cc)
						if inputTransport, err := _transportFactory.GetTransport(_transport); err == nil {
							inputProtocol := _tcompactProtocolFactory.GetProtocol(inputTransport)
							for {
								ok, err := processor.Process(defaultCtx, inputProtocol, inputProtocol)
								if errors.Is(err, thrift.ErrAbandonRequest) {
									break
								}
								if errors.As(err, new(thrift.TTransportException)) && err != nil {
									break
								}
								if !ok {
									break
								}
							}
						}
					}()
				}
			}
		}
	}
	if !t.isClose && err != nil {
		fmt.Println("operating services start failed:", err)
		sys.Wfs.Close()
		os.Exit(1)
	}
	return
}

func (t *service) Close() (err error) {
	defer util.Recover()
	if strings.TrimSpace(sys.OPADDR) != "" {
		t.isClose = true
		err = t.serverTransport.Close()
	}
	return
}

func (t *service) Serve() (err error) {
	if strings.TrimSpace(sys.OPADDR) != "" {
		tls := false
		if sys.Conf.Admin_Ssl_crt != "" && sys.Conf.Admin_Ssl_crt_key != "" {
			tls = true
		}
		err = t._server(strings.TrimSpace(sys.OPADDR), NewWfsIfaceProcessor(processor), tls, sys.Conf.Admin_Ssl_crt, sys.Conf.Admin_Ssl_crt_key)
	} else {
		sys.FmtLog("no operating services")
	}
	return
}

type pcontext struct {
	Id       int64
	isAuth   bool
	tt       thrift.TTransport
	mux      *sync.Mutex
	_isClose bool
}

func newCliContext(tt thrift.TTransport) (cc *pcontext) {
	cc = &pcontext{goutil.UUID64(), false, tt, &sync.Mutex{}, false}
	return
}

func (t *pcontext) close() {
	defer util.Recover()
	defer t.mux.Unlock()
	t.mux.Lock()
	if !t._isClose {
		t._isClose = true
		t.tt.Close()
	}
}
