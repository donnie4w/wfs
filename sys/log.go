// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
//

package sys

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/donnie4w/go-logger/logger"
)

var log = logger.NewLogger().SetFormat(logger.FORMAT_DATE | logger.FORMAT_TIME).SetLevel(logger.LEVEL_INFO)

func FmtLog(v ...any) {
	info := fmt.Sprint(v...)
	a, b := "", ""
	ll := 80
	if ll >= len(info) {
		for i := 0; i < (ll-len(info))/2; i++ {
			a = a + "="
		}
		b = a
		if ll > len(info)+len(a)*2 {
			b = a + "="
		}
	}
	log.Info(a, info, b)
}

func blankLine() {
	log.Write([]byte("\n"))
}

func wfslogo() {
	_r := `
	====================================================================
	==========      ==    ===    ==    =======    ======      ==========
	==========      ==   == ==   ==    ==         ==          ==========
	==========      ==  ==   ==  ==    ======     ======      ==========
	==========      == ==     == ==    ==             ==      ==========
	==========      ====       ====    ==         ======      ==========
	====================================================================
	`
	log.Info(_r)
}

func writePid(dir string) {
	path := dir + "/logs/wfs.pid"
	os.MkdirAll(filepath.Dir(path), 0777)
	if f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666); err == nil {
		defer f.Close()
		f.WriteString(fmt.Sprint(os.Getpid()))
	}
}
