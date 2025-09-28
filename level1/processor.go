// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package level1

import (
	"context"
	"strings"

	goutil "github.com/donnie4w/gofer/util"
	. "github.com/donnie4w/wfs/keystore"
	. "github.com/donnie4w/wfs/stub"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
)

type processhandle struct {
}

var processor = &processhandle{}

func ctx2CliContext(ctx context.Context) *pcontext {
	return ctx.Value("CliContext").(*pcontext)
}

func (t *processhandle) Ping(ctx context.Context) (_r int8, _err error) {
	defer util.Recover()
	mux := ctx2CliContext(ctx).mux
	defer mux.Unlock()
	mux.Lock()
	_r = 1
	return
}

func (t *processhandle) Append(ctx context.Context, wf *WfsFile) (_r *WfsAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_err = sys.ERR_AUTH.Error()
		return
	}
	_r = &WfsAck{Ok: true}
	if wf.Name != "" && len(wf.Data) > 0 {
		if int64(len(wf.Data)) > sys.DataMaxsize {
			_r.Ok, _r.Error = false, sys.ERR_OVERSIZE.WfsError()
			return
		}
		compress := sys.CompressType
		if wf.Compress != nil {
			compress = int32(*wf.Compress)
		}
		if _, err := sys.AppendData(wf.Name, wf.Data, compress); err != nil {
			_r.Ok, _r.Error = false, err.WfsError()
		}
	}
	return
}

func (t *processhandle) Delete(ctx context.Context, path string) (_r *WfsAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_err = sys.ERR_AUTH.Error()
		return
	}
	_r = &WfsAck{Ok: true}
	if path != "" {
		if err := sys.DelData(path); err != nil {
			_r.Ok, _r.Error = false, err.WfsError()
		}
	}
	return
}

func (t *processhandle) Get(ctx context.Context, path string) (_r *WfsData, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_err = sys.ERR_AUTH.Error()
		return
	}
	_r = &WfsData{}
	if path != "" {
		_r.Data = sys.GetData(path)
	}
	return
}

func (t *processhandle) Auth(ctx context.Context, wa *WfsAuth) (_r *WfsAck, _err error) {
	defer util.Recover()
	mux := ctx2CliContext(ctx).mux
	mux.Lock()
	defer mux.Unlock()
	_r = &WfsAck{Ok: true}
	if wa.Name != nil && wa.Pwd != nil && auth(*wa.Name, *wa.Pwd) {
		ctx2CliContext(ctx).isAuth = true
	} else {
		_r.Ok, _r.Error = false, sys.ERR_NOPASS.WfsError()
	}
	return
}

func (t *processhandle) Rename(ctx context.Context, path string, newpath string) (_r *WfsAck, _err error) {
	defer util.Recover()
	cc := ctx2CliContext(ctx)
	cc.mux.Lock()
	defer cc.mux.Unlock()
	if noAuthAndClose(cc) {
		_err = sys.ERR_AUTH.Error()
		return
	}
	_r = &WfsAck{Ok: true}
	if err := sys.Modify(path, newpath); err != nil {
		_r.Ok = false
		_r.Error = err.WfsError()
	}
	return
}

func auth(name, pwd string) (b bool) {
	if _r, ok := Admin.GetAdmin(name); ok {
		b = strings.EqualFold(_r.Pwd, goutil.Md5Str(pwd))
	}
	return
}

func noAuthAndClose(cc *pcontext) (b bool) {
	if !cc.isAuth {
		cc.tt.Close()
		b = true
	}
	return
}
