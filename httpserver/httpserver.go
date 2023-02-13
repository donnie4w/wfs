/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package httpserver

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	_ "net/http/pprof"
	"os"
	"strconv"
	"time"
	. "wfs/conf"
	"wfs/storge"

	"github.com/donnie4w/simplelog/logging"
	"github.com/donnie4w/tlnet"
)

type CmdType string

const (
	setWeight   CmdType = "setWeight"
	addSlave    CmdType = "addSlave"
	cutOff      CmdType = "cutOff"
	removeSlave CmdType = "removeSlave"
	slavelist   CmdType = "slavelist"
	ping        CmdType = "ping"
)
const (
	_200 = "200" //ok
	_403 = "403" //forbidden
	_404 = "404" //no found
	_500 = "500" //err
)

//func getData(key string) (bs []byte, err error) {
//	bs, shardname, err := storge.GetData(key)
//	if err == nil {
//		if bs == nil && len(shardname) > 0 {
//		}
//	}
//	return
//}

// func read(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
func read(hc *tlnet.HttpContext) {
	//	uri := r.RequestURI
	//	uri = uri[3:]
	//	name := uri
	//	arg := ""
	//	if strings.Contains(uri, "?") {
	//		index := strings.Index(uri, "?")
	//		name = uri[:index]
	//		arg = uri[index:]
	//	}
	//	w.Header().Set("Content-Type", "image/jpg")
	//	bs, err := getData(name)
	//	if err == nil {
	//		if strings.HasPrefix(arg, "?imageView2") {
	//			spec := NewSpec(bs, arg)
	//			w.Write(spec.GetData())
	//		} else {
	//			w.Write(bs)
	//		}
	//	} else {
	//		fmt.Println("err:", err.Error())
	//		w.Write([]byte("404"))
	//	}

	// bs, err := GetData(r.RequestURI)
	bs, err := GetData(hc.Request().RequestURI)
	if err == nil {
		// w.Write(bs)
		hc.ResponseBytes(0, bs)
	} else {
		// w.Write([]byte("404"))
		hc.Writer().WriteHeader(404)
		// hc.ResponseString("404")
	}

	//	if CF.Keepalive > 0 {
	//		c, ok := w.(http.Hijacker)
	//		if ok {
	//			a, _, err := c.Hijack()
	//			if err == nil {
	//				go pool.Add(a)
	//			} else {
	//				fmt.Println("err:", err.Error())
	//			}
	//		}
	//	}
}

// func upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
// 	file, handler, err := r.FormFile("file")
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer file.Close()
// 	buf := new(bytes.Buffer)
// 	by := make([]byte, 512)
// 	for {
// 		if n, err := file.Read(by); err == nil {
// 			buf.Write(by[:n])
// 			if int64(buf.Len()) > CF.MaxFileSize {
// 				fmt.Fprint(w, "file too large")
// 				return
// 			}
// 		} else {
// 			break
// 		}
// 	}
// 	bs := buf.Bytes()
// 	if len(bs) > 0 {
// 		name := ""
// 		uri := r.RequestURI
// 		if len(uri) > 2 {
// 			name = uri[3:]
// 		}
// 		if name == "" {
// 			name = handler.Filename
// 		}
// 		//	err := storge.AppendData(bs, name, "")
// 		err := AppendData(bs, name, "")
// 		if err == nil {
// 			fmt.Fprint(w, "ok:", len(bs), " ", name)
// 		} else {
// 			fmt.Fprint(w, err.Error())
// 		}
// 	}
// }

func upload(hc *tlnet.HttpContext) {
	defer myRecover()
	hc.MaxBytesReader(CF.MaxFileSize)
	file, handler, err := hc.FormFile("file")
	if err != nil {
		// logging.Error(err)
		hc.ResponseString(err.Error())
		return
	}
	defer file.Close()
	// buf := new(bytes.Buffer)
	// by := make([]byte, 512)
	// for {
	// 	if n, err := file.Read(by); err == nil {
	// 		buf.Write(by[:n])
	// 		if int64(buf.Len()) > CF.MaxFileSize {
	// 			fmt.Fprint(w, "file too large")
	// 			return
	// 		}
	// 	} else {
	// 		break
	// 	}
	// }
	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	bs := buf.Bytes()
	if len(bs) > 0 {
		name := ""
		uri := hc.Request().RequestURI
		if len(uri) > 2 {
			name = uri[3:]
		}
		if name == "" {
			name = handler.Filename
		}
		//	err := storge.AppendData(bs, name, "")
		err := AppendData(bs, name, "")
		if err == nil {
			// fmt.Fprint(w, "ok:", len(bs), " ", name)
			hc.ResponseString(fmt.Sprint("ok:", len(bs), " ", name))
		} else {
			// fmt.Fprint(w, err.Error())
			hc.ResponseString(err.Error())
		}
	}
}

// func del(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
func del(hc *tlnet.HttpContext) {
	var name string
	uri := hc.Request().RequestURI
	if len(uri) > 2 {
		name = uri[3:]
	}
	if name != "" {
		//		err := storge.DelData(name)
		err := DelData(name)
		if err == nil {
			// fmt.Fprintf(w, "delete ok, %s!\n", name)
			hc.ResponseString(fmt.Sprint("delete ok,", name))
		} else {
			// fmt.Fprintf(w, "delete error, %s!\n", name, " | ", err.Error())
			hc.ResponseString(fmt.Sprint("delete error,", name, " | ", err.Error()))
		}
	}
}

//func domain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	appkey := ps.ByName("appkey")
//	domain := ps.ByName("domain")
//	fmt.Fprintf(w, "domain create ok:, %s\n", appkey, " ", domain)
//}

//func ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	fmt.Fprintf(w, "")
//}

// func check(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
func check(hc *tlnet.HttpContext) {
	r := hc.Request()
	if r.Body != nil {
		defer r.Body.Close()
		bs, _ := ioutil.ReadAll(r.Body)
		name := string(bs)
		b := storge.Exsit(name)
		if b {
			// w.Write([]byte{1})
			hc.ResponseString("1")
		} else {
			// w.Write([]byte{0})
			hc.ResponseString("0")
		}
	}
}

/********************************************************************************/
func savePort(port int) {
	f, err := os.OpenFile(fmt.Sprint(CF.FileData, "/CURRENT"), os.O_RDWR|os.O_CREATE, 0644)
	defer f.Close()
	if err == nil {
		_, err = f.WriteString(fmt.Sprint(port))
		if err != nil {
			os.Exit(1)
		}
	} else {
		os.Exit(1)
	}
}

func readPort() (port int) {
	bs, err := ioutil.ReadFile(fmt.Sprint(CF.FileData, "/CURRENT"))
	if err == nil {
		port, _ = strconv.Atoi(string(bs))
	}
	return
}

/********************************************************************************/
func Start() {
	if IsCmd() {
		CF.Port = readPort()
		switch {
		case Cmd.AddSlave != "":
			wfsCmd(addSlave, Cmd.AddSlave)
		case Cmd.CutOff:
			wfsCmd(cutOff, "1")
		case Cmd.SetWeight != "":
			wfsCmd(setWeight, Cmd.SetWeight)
		case Cmd.RemoveSlave != "":
			wfsCmd(removeSlave, Cmd.RemoveSlave)
		case Cmd.Slavelist:
			wfsCmd(slavelist, "1")
		case Cmd.Ping != "":
			wfsCmd(ping, Cmd.Ping)
		}
	} else {

		if CF.Pprof > 0 {
			go http.ListenAndServe(fmt.Sprint(CF.Bind, ":", CF.Pprof), nil)
		}

		storge.Init()
		// router := httprouter.New()
		// router.POST("/thrift", thandler)
		// //		router.POST("/ping", ping)
		// router.GET("/r/*.r", read)
		// router.POST("/c", check)
		// router.POST("/u/*.r", upload)
		// router.POST("/u", upload)
		// router.DELETE("/d/*.r", del)

		// fmt.Println(CF.Port)
		// srv := &http.Server{
		// 	ReadTimeout: time.Duration(CF.ServerReadTimeout) * time.Second,
		// 	Addr:        fmt.Sprint(CF.Bind, ":", CF.Port),
		// 	Handler:     router,
		// }
		// savePort(CF.Port)
		// Init()
		// fmt.Println("wfs start,listen:", CF.Port)
		// err := srv.ListenAndServe()

		savePort(CF.Port)
		Init()
		tl := tlnet.NewTlnet()
		tl.AddHandlerFunc("/thrift", nil, thandler)
		tl.POST("/u/", upload)
		tl.GET("/r/", read)
		tl.POST("/c", check)
		tl.POST("/u", upload)
		tl.DELETE("/d/", del)
		tl.ReadTimeout(time.Duration(CF.ServerReadTimeout) * time.Second)
		logging.Info("wfs start,listen:", CF.Port)
		err := tl.HttpStart(fmt.Sprint(CF.Bind, ":", CF.Port))
		if err != nil {
			logging.Error("httpserver start error:", err.Error())
			os.Exit(1)
		}
	}
}

func myRecover() {
	if er := recover(); er != nil {
		logging.Error(er)
	}
}

//----------------------------------------------------------------------------------------------

// curl -F "file=@015.jpg" "http://127.0.0.1:3434/u/015.jpg"
// curl -X DELETE "http://127.0.0.1:3434/d/aaa/1.jpg"
// curl -I -X GET "http://127.0.0.1:3434/r/aaa/1.jpg"
