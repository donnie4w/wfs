// Copyright (c) 2023, donnie <donnie4w@gmail.com>
// All rights reserved.
// Use of t source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// github.com/donnie4w/wfs

package tc

import (
	"fmt"
	"github.com/donnie4w/gofer/base58"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/donnie4w/go-logger/logger"
	"github.com/donnie4w/gofer/gosignal"
	. "github.com/donnie4w/gofer/hashmap"
	"github.com/donnie4w/gofer/image"
	tldbKs "github.com/donnie4w/gofer/keystore"
	goutil "github.com/donnie4w/gofer/util"
	"github.com/donnie4w/tlnet"
	. "github.com/donnie4w/wfs/keystore"
	"github.com/donnie4w/wfs/stub"
	"github.com/donnie4w/wfs/sys"
	"github.com/donnie4w/wfs/util"
)

func init() {
	sys.Serve.Put(4, adminservice)
	sys.Serve.Put(5, clientservice)
	sys.WsClient = newAdminLocalWsClient
}

type adminService struct {
	isClose bool
	tlAdmin *tlnet.Tlnet
}

var adminservice = &adminService{false, tlnet.NewTlnet()}

func (t *adminService) Serve() (err error) {
	images = &image.Image{}
	if sys.Conf.Resample == 1 {
		images.ResizeFilter = image.Lanczos
	} else {
		images.ResizeFilter = image.MitchellNetravali
	}

	if strings.TrimSpace(sys.DEBUGADDR) != "" {
		go tlDebug()
		<-time.After(500 * time.Millisecond)
	}
	t.addSignalEvent()
	initAccount()
	initStore()
	if strings.TrimSpace(sys.WEBADDR) != "" {
		err = t._serve(strings.TrimSpace(sys.WEBADDR), sys.Conf.Admin_Ssl_crt, sys.Conf.Admin_Ssl_crt_key)
	}
	return
}

func (t *adminService) Close() (err error) {
	defer util.Recover()
	if strings.TrimSpace(sys.WEBADDR) != "" {
		t.isClose = true
		err = t.tlAdmin.Close()
	}
	return
}

func (s *adminService) addSignalEvent() {
	gosignal.ListenSignalEvent(func(sig os.Signal) {
		sys.FmtLog("services closing down")
		sys.Wfs.Close()
		sys.FmtLog("services stopped")
		os.Exit(0)
	}, syscall.SIGTERM, syscall.SIGINT, syscall.Signal(0xa))
}

func (t *adminService) _serve(addr string, serverCrt, serverKey string) (err error) {
	defer util.Recover()
	if addr, err = util.ParseAddr(addr); err != nil {
		return
	}
	sys.WEBADDR = addr
	t.tlAdmin.Handle("/login", loginHandler)
	t.tlAdmin.Handle("/init", initHandler)
	t.tlAdmin.Handle("/lang", langHandler)
	t.tlAdmin.Handle("/", initHandler)
	t.tlAdmin.Handle("/about", aboutHandler)
	t.tlAdmin.Handle("/bootstrap.css", cssHandler)
	t.tlAdmin.Handle("/bootstrap.min.js", jsHandler)
	t.tlAdmin.HandleWithFilter("/r/", loginFilter(), readHandler)
	t.tlAdmin.HandleWithFilter("/file", loginFilter(), fileHtml)
	t.tlAdmin.HandleWithFilter("/fragment", loginFilter(), fragmentHtml)
	t.tlAdmin.HandleWithFilter("/defrag", loginFilter(), defragData)
	t.tlAdmin.HandleWithFilter("/filedata", loginFilter(), fileDataHandler)
	t.tlAdmin.HandleWithFilter("/monitor", loginFilter(), monitorHtml)
	t.tlAdmin.HandleWebSocketBindConfig("/monitorData", mntHandler, mntConfig())
	t.tlAdmin.HandleWithFilter("/append/", authFilter(), appendHandler)
	t.tlAdmin.HandleWithFilter("/delete/", authFilter(), deleteHandler)
	t.tlAdmin.HandleWithFilter("/rename", loginFilter(), renameHandler)
	t.tlAdmin.HandleWebSocketBindConfig("/export", exportHandler, wsConfig())
	t.tlAdmin.HandleWebSocketBindConfig("/exportincr", exportIncrHandler, wsConfig())
	t.tlAdmin.HandleWebSocketBindConfig("/exportfile", exportFileHandler, wsConfig())
	t.tlAdmin.HandleWebSocketBindConfig("/import", importHandler, wsConfig())
	t.tlAdmin.HandleWebSocketBindConfig("/importfile", importfileHandler, wsConfig())

	if serverCrt != "" && serverKey != "" {
		sys.FmtLog("webAdmin start tls [", addr, "]")
		err = t.tlAdmin.HttpsStart(addr, serverCrt, serverKey)
	}
	if !t.isClose {
		sys.FmtLog("webAdmin start [", addr, "]")
		err = t.tlAdmin.HttpStart(addr)
	}
	if !t.isClose && err != nil {
		fmt.Println("webAdmin start failed:", err.Error())
		sys.Wfs.Close()
		os.Exit(1)
	}
	if t.isClose {
		err = nil
	}
	return
}

var sessionMap = NewMapL[string, *tldbKs.UserBean]()

func loginFilter() (f *tlnet.Filter) {
	defer util.Recover()
	f = tlnet.NewFilter()
	f.AddIntercept(".*?", func(hc *tlnet.HttpContext) bool {
		if len(Admin.AdminList()) > 0 {
			if !isLogin(hc) {
				hc.Redirect("/login")
				return true
			}
		} else {
			hc.Redirect("/init")
			return true
		}
		return false
	})

	f.AddIntercept(`[^\s]+`, func(hc *tlnet.HttpContext) bool {
		if hc.PostParamTrimSpace("atype") != "" && !isAdmin(hc) {
			hc.ResponseString(resultHtml("Permission Denied"))
			return true
		}
		return false
	})
	return
}

func authFilter() (f *tlnet.Filter) {
	defer util.Recover()
	f = tlnet.NewFilter()
	f.AddIntercept(".*?", func(hc *tlnet.HttpContext) bool {
		if isLogin(hc) {
			return false
		}
		if authAccount(hc) {
			return false
		}
		hc.ResponseBytes(http.StatusUnauthorized, nil)
		return true
	})
	return
}

func authAccount(hc *tlnet.HttpContext) bool {
	name := hc.Request().Header.Get("username")
	pwd := hc.Request().Header.Get("password")
	if _r, ok := Admin.GetAdmin(name); ok {
		if strings.EqualFold(_r.Pwd, goutil.Md5Str(pwd)) || strings.EqualFold(_r.Pwd, pwd) {
			return true
		}
	}
	return false
}

func getSessionid() string {
	return fmt.Sprint("t", goutil.CRC32(goutil.Int64ToBytes(sys.UUID)))
}

func getLangId() string {
	return fmt.Sprint("l", goutil.CRC32(goutil.Int64ToBytes(sys.UUID)))
}

func isLogin(hc *tlnet.HttpContext) (isLogin bool) {
	if len(Admin.AdminList()) > 0 {
		if _r, err := hc.GetCookie(getSessionid()); err == nil && sessionMap.Has(_r) {
			isLogin = true
		}
	}
	return
}

func isAdmin(hc *tlnet.HttpContext) (_r bool) {
	if c, err := hc.GetCookie(getSessionid()); err == nil {
		if u, ok := sessionMap.Get(c); ok {
			_r = u.Type == 1
		}
	}
	return
}

func langHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	lang := hc.GetParamTrimSpace("lang")
	if lang == "en" || lang == "zh" {
		hc.SetCookie(getLangId(), lang, "/", 86400)
	}
	hc.Redirect("/")
}

func getLang(hc *tlnet.HttpContext) LANG {
	if lang, err := hc.GetCookie(getLangId()); err == nil {
		if lang == "zh" {
			return ZH
		} else if lang == "en" {
			return EN
		}
	}
	return ZH
}

func cssHandler(hc *tlnet.HttpContext) {
	hc.Writer().Header().Add("Content-Type", "text/html")
	textTplByText(cssContent(), nil, hc)
}

func jsHandler(hc *tlnet.HttpContext) {
	hc.Writer().Header().Add("Content-Type", "text/html")
	textTplByText(jsContent(), nil, hc)
}

func aboutHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	tplToHtml(getLang(hc), ABOUT, sys.VERSION, hc)
}

/***********************************************************************/
func initHandler(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	if len(Admin.AdminList()) > 0 && !isLogin(hc) {
		hc.Redirect("/login")
		return
	}
	if _type := hc.GetParam("type"); _type != "" {
		isadmin := isAdmin(hc)
		if _type == "1" {
			if name, pwd, _type := hc.PostParamTrimSpace("adminName"), hc.PostParamTrimSpace("adminPwd"), hc.PostParamTrimSpace("adminType"); name != "" && pwd != "" {
				if n := len(Admin.AdminList()); (n > 0 && isadmin) || n == 0 {
					alterType := false
					if t, err := strconv.Atoi(_type); err == nil {
						if _r, err := hc.GetCookie(getSessionid()); err == nil && sessionMap.Has(_r) {
							if u, ok := sessionMap.Get(_r); ok && u.Name == name && t != int(u.Type) {
								alterType = true
							}
						}
						if !alterType {
							Admin.PutAdmin(name, pwd, int8(t))
						}
					}
				} else {
					goto DENIED
				}
			}
		} else if _type == "2" && isLogin(hc) {
			if isadmin {
				if name := hc.PostParamTrimSpace("adminName"); name != "" {
					if u, ok := Admin.GetAdmin(name); ok && u.Type == 1 {
						i, j := 0, 0
						for _, s := range Admin.AdminList() {
							if _u, _ := Admin.GetAdmin(s); _u.Type == 1 {
								i++
							} else if _u.Type == 2 {
								j++
							}
						}
						if j > 0 && i == 1 {
							hc.ResponseString(resultHtml("failed,There cannot be only Observed users"))
							return
						}
					}
					Admin.DelAdmin(name)
					sessionMap.Range(func(k string, v *tldbKs.UserBean) bool {
						if v.Name == name {
							sessionMap.Del(k)
						}
						return true
					})
				}
			} else {
				goto DENIED
			}
		}
		hc.Redirect("/init")
		return
	} else {
		initHtml(hc)
		return
	}
DENIED:
	hc.ResponseString(resultHtml("Permission Denied"))
}

func loginHandler(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	if hc.PostParamTrimSpace("type") == "1" {
		name, pwd := hc.PostParamTrimSpace("name"), hc.PostParamTrimSpace("pwd")
		if _r, ok := Admin.GetAdmin(name); ok {
			if strings.EqualFold(_r.Pwd, goutil.Md5Str(pwd)) {
				sid := goutil.Md5Str(fmt.Sprint(time.Now().UnixNano()))
				sessionMap.Put(sid, _r)
				hc.SetCookie(getSessionid(), sid, "/", 86400)
				hc.Redirect("/")
				return
			}
		}
		hc.ResponseString(resultHtml("Login Failed"))
		return
	}
	loginHtml(hc)
}

func initHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	_isAdmin := isAdmin(hc)
	show, init, sc := "", false, _isAdmin
	if len(Admin.AdminList()) == 0 {
		show, init, sc = "no user is created for admin, create a management user first", true, true
	}
	av := &AdminView{Show: show, Init: init, ShowCreate: sc}
	if isLogin(hc) {
		m := make(map[string]string, 0)
		for _, s := range Admin.AdminList() {
			if u, ok := Admin.GetAdmin(s); ok {
				if _isAdmin && u.Type == 1 {
					m[s] = "Admin"
				} else if u.Type == 2 {
					m[s] = "Observed"
				}
			}
		}
		av.AdminUser = m
	}
	tplToHtml(getLang(hc), INIT, av, hc)
}

func loginHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	tplToHtml(getLang(hc), LOGIN, []byte{}, hc)
}

func initAccount() {
	if sys.Conf.AdminUserName != nil && sys.Conf.AdminPassword != nil {
		Admin.PutAdmin(*sys.Conf.AdminUserName, *sys.Conf.AdminPassword, 1)
	} else if sys.Conf.Init && len(Admin.AdminList()) == 0 {
		Admin.PutAdmin(sys.DefaultAccount[0], sys.DefaultAccount[1], 1)
	}
}

func initStore() {
	Admin.PutOther("webaddr", sys.WEBADDR)
}

func readHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	uri := hc.Request().RequestURI
	if bs, cy, err := getData(uri[2:]); err == nil {
		if cy != "" {
			hc.Writer().Header().Add("Content-Type", cy)
		}
		hc.ResponseBytes(0, bs)
	} else {
		hc.Writer().WriteHeader(404)
	}
}

func fileHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	if hc.PostParam("limit") == "0" {
		hc.ResponseString(fmt.Sprint(`{"limit":`, sys.DataMaxsize, `}`))
	} else {
		tplToHtml(getLang(hc), FILE, nil, hc)
	}
}

func fragmentHtml(hc *tlnet.HttpContext) {
	defer func() {
		if err := recover(); err != nil {
			hc.ResponseString(resultHtml("server error:", err))
		}
	}()
	fbs := make([]*FragmentBean, 0)
	filepath.WalkDir(sys.WFSDATA+"/wfsfile", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if i, b := base58.DecodeForInt64([]byte(d.Name())); b && util.CheckNodeId(int64(i)) && !sys.IsEmptyBigFile(path) {
				fi, _ := d.Info()
				fb := &FragmentBean{Name: d.Name(), FileSize: fi.Size(), Time: fi.ModTime().Format(time.DateTime)}
				if fa, err := sys.FragAnalysis(d.Name()); err == nil {
					fb.FragmentSize = fa.FileSize - fa.ActualSize + fa.RmSize
					fb.Status = 1
				} else if err.Equal(sys.ERR_DEFRAG_FORBID) {
					fb.Status = 2
				}
				fbs = append(fbs, fb)
			}
		}
		return nil
	})
	sort.Slice(fbs, func(i, j int) bool { return fbs[i].Time > fbs[j].Time })
	tplToHtml(getLang(hc), FRAGMENT, fbs, hc)
}

var defragStatus = false

func defragData(hc *tlnet.HttpContext) {
	if defragStatus {
		hc.ResponseString(`{"status":false,"desc":"The system is defragmenting"}`)
		return
	}
	node := hc.PostParamTrimSpace("node")
	if node != "" {
		defragStatus = true
		if err := sys.Defrag(node); err == nil {
			hc.ResponseString(`{"status":true}`)
		} else {
			hc.ResponseString(`{"status":false,"desc":"` + err.WfsError().GetInfo() + `"}`)
		}
		defragStatus = false
	}
}

func fileDataHandler(hc *tlnet.HttpContext) {
	searchType := hc.PostParamTrimSpace("searchType")
	if searchType == "1" {
		pageNumber, _ := strconv.Atoi(hc.PostParamTrimSpace("pageNumber"))
		pagecount, _ := strconv.Atoi(hc.PostParamTrimSpace("pagecount"))
		lastId, _ := strconv.Atoi(hc.PostParamTrimSpace("lastId"))
		if fp := searchByPage(pageNumber, pagecount, int64(lastId)); fp != nil {
			hc.ResponseBytes(0, goutil.JsonEncode(fp))
		}
	} else if searchType == "2" {
		prev := hc.PostParamTrimSpace("prevName")
		if fp := searchByPrev(prev); fp != nil {
			hc.ResponseBytes(0, goutil.JsonEncode(fp))
		}
	}
}

func webclientInfo() (proxy string, protocol string, port int) {
	conf := sys.GetConfg()
	if conf != nil && util.IsURL(conf.ImgViewingRevProxy) {
		return conf.ImgViewingRevProxy, "", 0
	}
	protocol = "http://"
	if sys.Conf.Ssl_crt != "" && sys.Conf.Ssl_crt_key != "" {
		protocol = "https://"
	}
	port = sys.LISTEN
	return
}

func searchByPage(pagenum, pagecount int, lastId int64) (fp *FilePage) {
	if pagenum < 1 {
		pagenum = 1
	}
	if nextId := sys.Seq() - int64((pagenum-1)*pagecount); nextId > 0 {
		if lastId > 0 && nextId >= lastId {
			nextId = lastId - 1
		}
		if pbs := sys.SearchLimit(nextId, int64(pagecount)); pbs != nil {
			fp = &FilePage{TotalNum: int(sys.Count()), CurrentNum: pagenum, FS: make([]*FileBean, 0)}
			fp.RevProxy, fp.CliProtocol, fp.ClientPort = webclientInfo()
			for _, pb := range pbs {
				fb := &FileBean{Name: pb.Path, Size: len(pb.Body), Time: util.TimestrampFormat(pb.Timestramp), Id: int(pb.Id)}
				fp.FS = append(fp.FS, fb)
			}
		}
	}
	return
}

func searchByPrev(prev string) (fp *FilePage) {
	if prev != "" {
		if pbs := sys.SearchLike(prev); pbs != nil {
			fp = &FilePage{FS: make([]*FileBean, 0), TotalNum: len(pbs), CurrentNum: 1}
			fp.RevProxy, fp.CliProtocol, fp.ClientPort = webclientInfo()
			for _, pb := range pbs {
				fb := &FileBean{Name: pb.Path, Size: len(pb.Body), Time: util.TimestrampFormat(pb.Timestramp), Id: int(pb.Id)}
				fp.FS = append(fp.FS, fb)
			}
			sort.Slice(fp.FS, func(i, j int) bool { return fp.FS[i].Id > fp.FS[j].Id })
		}
	}
	return
}

func monitorHtml(hc *tlnet.HttpContext) {
	tplToHtml(getLang(hc), MONITOR, nil, hc)
}

func mntConfig() (wc *tlnet.WebsocketConfig) {
	wc = &tlnet.WebsocketConfig{}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
		if !isLogin(hc) {
			hc.WS.Close()
			return
		}
	}
	return
}

func mntHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	s := string(hc.WS.Read())
	if t, err := strconv.Atoi(s); err == nil {
		if t < 1 {
			t = 1
		}
		for hc.WS.Error == nil {
			if j, err := monitorToJson(); err == nil {
				hc.WS.Send(j)
			}
			<-time.After(time.Duration(t) * time.Second)
		}
	}
}

func wsConfig() (wc *tlnet.WebsocketConfig) {
	wc = &tlnet.WebsocketConfig{}
	wc.OnOpen = func(hc *tlnet.HttpContext) {
		if !authAccount(hc) {
			hc.WS.Send([]byte{1})
			<-time.After(2 * time.Second)
			hc.WS.Close()
		}
	}
	return
}

func exportHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	bs := hc.WS.Read()
	if len(bs) == 1 && bs[0] == 1 {
		sys.Export(func(bean *stub.SnapshotBean) bool {
			if bs, err := bean.ToBytes(); err == nil && hc.WS.Error == nil {
				if err = hc.WS.Send(bs); err == nil {
					return true
				}
			}
			return false
		})
		hc.WS.Send([]byte{0})
	}
}

func exportIncrHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	bs := hc.WS.Read()
	if len(bs) == 16 {
		start := goutil.BytesToInt64(bs[:8])
		limit := goutil.BytesToInt64(bs[8:16])
		sys.ExportByINCR(start, limit, func(beans *stub.SnapshotBeans) bool {
			if beans != nil {
				if bs, err := beans.ToBytes(); err == nil && hc.WS.Error == nil {
					if err = hc.WS.Send(bs); err == nil {
						return true
					}
				}
			}
			return false
		})
		hc.WS.Send([]byte{0})
	}
}

func exportFileHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	bs := hc.WS.Read()
	if len(bs) == 16 {
		start := goutil.BytesToInt64(bs[:8])
		limit := goutil.BytesToInt64(bs[8:16])
		sys.ExportFile(start, limit, func(beans *stub.SnapshotFile) bool {
			if beans != nil {
				if bs, err := beans.ToBytes(); err == nil && hc.WS.Error == nil {
					if err = hc.WS.Send(bs); err == nil {
						return true
					}
				}
			}
			return false
		})
		hc.WS.Send([]byte{0})
	}
}

var importcover = true

func importHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	bs := hc.WS.Read()

	if len(bs) == 1 {
		switch bs[0] {
		case 1:
			hc.WS.Send([]byte{0})
		case 2:
			importcover = true
		case 3:
			importcover = false
		}
		return
	}

	if sb, err := stub.BytesToSnapshotBean(bs); err == nil {
		sys.Import(sb, importcover)
	} else {
		logger.Error(err)
		hc.WS.Send([]byte{2})
		hc.WS.Close()
	}

}

func importfileHandler(hc *tlnet.HttpContext) {
	defer util.Recover()
	bs := hc.WS.Read()

	if len(bs) == 1 {
		switch bs[0] {
		case 1:
			hc.WS.Send([]byte{0})
		case 2:
			importcover = true
		case 3:
			importcover = false
		}
		return
	}

	if ssf, err := stub.BytesToSnapshotFile(bs); err == nil {
		sys.ImportFile(ssf)
	} else {
		logger.Error(err)
		hc.WS.Send([]byte{2})
		hc.WS.Close()
	}

}
