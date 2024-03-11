// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs
//

package tc

import (
	htmlTpl "html/template"
	"os"
	textTpl "text/template"

	"github.com/donnie4w/simplelog/logging"
	"github.com/donnie4w/tlnet"
)

type TXTYPE int
type LANG int

const (
	ZH LANG = 0
	EN LANG = 1
)

const (
	_ TXTYPE = iota
	LOGIN
	INIT
	ABOUT
	FILE
	FRAGMENT
	MONITOR
)

var mod = 1 //0debug，1release

func tplToHtml(lang LANG, flag TXTYPE, v any, hc *tlnet.HttpContext) {
	dir, _ := os.Getwd()
	switch flag {
	case LOGIN:
		tpl(lang, dir+"/tc/html/login.html", loginText, dir+"/tc/html/loginEn.html", loginEnText, v, hc)
	case INIT:
		tpl(lang, dir+"/tc/html/init.html", initText, dir+"/tc/html/initEn.html", initEnText, v, hc)
	case ABOUT:
		tpl(lang, dir+"/tc/html/about.html", aboutText, dir+"/tc/html/aboutEn.html", aboutEnText, v, hc)
	case FILE:
		tpl(lang, dir+"/tc/html/file.html", fileText, dir+"/tc/html/fileEn.html", fileEnText, v, hc)
	case FRAGMENT:
		tpl(lang, dir+"/tc/html/fragment.html", fragmentText, dir+"/tc/html/fragmentEn.html", fragmentEnText, v, hc)
	case MONITOR:
		tpl(lang, dir+"/tc/html/monitor.html", monitorText, dir+"/tc/html/monitorEn.html", monitorEnText, v, hc)
	}
}

func tpl(lang LANG, tplZHPath, tplZHText, tplENPath, tplENText string, v any, hc *tlnet.HttpContext) {
	if lang == ZH {
		if mod == 0 {
			textTplByPath(tplZHPath, v, hc)
		} else if mod == 1 {
			textTplByText(tplZHText, v, hc)
		}
	} else if lang == EN {
		if mod == 0 {
			textTplByPath(tplENPath, v, hc)
		} else if mod == 1 {
			textTplByText(tplENText, v, hc)
		}
	}
}

func textTplByPath(path string, data any, hc *tlnet.HttpContext) {
	if tp, err := textTpl.ParseFiles(path); err == nil {
		tp.Execute(hc.Writer(), data)
	} else {
		logging.Error(err)
	}
}

func textTplByText(text string, data any, hc *tlnet.HttpContext) {

	tl := textTpl.New("wfsTemplate")
	if _, err := tl.Parse(text); err == nil {
		tl.Execute(hc.Writer(), data)
	} else {
		logging.Error(err)
	}
}

func htmlTplByPath(path string, data any, hc *tlnet.HttpContext) {
	if tp, err := htmlTpl.ParseFiles(path); err == nil {
		tp.Execute(hc.Writer(), data)
	} else {
		logging.Error(err)
	}
}
