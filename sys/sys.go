// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
package sys

import (
	"time"
)

type server struct{}

func (s *server) Serve() error {
	praseflag()
	blankLine()
	wfslogo()
	KeyStoreInit(WFSDATA)
	writePid(WFSDATA)
	Serve.Ascend(func(_ int, s Server) bool {
		go func() {
			if err := s.Serve(); err != nil {
				FmtLog(err)
			}
		}()
		<-time.After(time.Millisecond << 9)
		return true
	})
	select {}
}

func (s *server) Close() (err error) {
	Serve.Descend(func(_ int, s Server) bool {
		defer func() { recover() }()
		s.Close()
		return true
	})
	return
}
