/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package httpserver

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//	_ "net/http/pprof"
	"os"
	"strings"

	"github.com/donnie4w/go-logger/logger"
	. "github.com/donnie4w/wfs/conf"
	"github.com/donnie4w/wfs/storge"
	"github.com/julienschmidt/httprouter"
)

func getData(key string) (bs []byte, err error) {
	bs, err = storge.GetData(key)
	return
}

func read(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	uri := r.RequestURI
	uri = uri[3:]
	name := uri
	arg := ""
	if strings.Contains(uri, "?") {
		index := strings.Index(uri, "?")
		name = uri[:index]
		arg = uri[index:]
	}
	//	fmt.Println("key===>", name)
	//	w.Header().Set("Content-Type", "text/html")
	//	w.Header().Set("Content-Type", "image/jpg")
	bs, err := getData(name)
	if err == nil {
		if strings.HasPrefix(arg, "?imageView2") {
			spec := NewSpec(bs, arg)
			w.Write(spec.GetData())
		} else {
			w.Write(bs)
		}
	} else {
		w.Write([]byte("404"))
	}

}

func upload(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	file, handler, err := r.FormFile("file")
	defer file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	bs, err := ioutil.ReadAll(file)
	if err == nil && bs != nil {
		if int64(len(bs)) > CF.MaxFileSize {
			fmt.Fprint(w, "too large:", len(bs))
			return
		} else {
			fmt.Fprint(w, "ok:", len(bs), " ", handler.Filename)
		}

	}
	name := ""
	uri := r.RequestURI
	if len(uri) > 2 {
		name = uri[3:]
	}
	if name == "" {
		name = handler.Filename
	}
	//	fmt.Println("name:", name)
	//	contentType := r.Header.Get("Content-Type")
	storge.AppendData(bs, name, "")
}

func delete(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ""
	uri := r.RequestURI
	if len(uri) > 2 {
		name = uri[3:]
	}
	if name != "" {
		err := storge.DelData(name)
		if err == nil {
			fmt.Fprintf(w, "delete ok, %s!\n", name)
		} else {
			fmt.Fprintf(w, "delete error, %s!\n", name, " | ", err.Error())
		}
	}
}

func domain(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	appkey := ps.ByName("appkey")
	domain := ps.ByName("domain")
	fmt.Fprintf(w, "domain create ok:, %s\n", appkey, " ", domain)
}

func ping(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
}

func check(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	name := ps.ByName("name")
	b := storge.Exsit(name)
	if b {
		w.Write([]byte{1})
	}
}

//func _pprof(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
//	pprof.Profile(w, r)
//}

func Start() {
	storge.Init()
	router := httprouter.New()

	router.POST("/thrift", thandler)
	router.POST("/ping", ping)
	router.GET("/r/*.r", read)
	router.POST("/c", check)
	router.POST("/u/*.r", upload)
	router.POST("/u", upload)
	router.DELETE("/d/*.r", delete)
	//	go http.ListenAndServe(fmt.Sprint(":", 5555), nil)
	err := http.ListenAndServe(fmt.Sprint(":", CF.Port), router)
	if err != nil {
		logger.Error("httpserver start error:", err.Error())
		os.Exit(1)
	}
}

// curl -F "file=@015.jpg" "http://127.0.0.1:3434/u/015.jpg"
// curl -X DELETE "http://127.0.0.1:3434/d/aaa/1.jpg"
// curl -I -X GET "http://127.0.0.1:3434/r/aaa/1.jpg"
