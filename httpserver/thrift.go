/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package httpserver

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"git.apache.org/thrift.git/lib/go/thrift"
	. "github.com/donnie4w/wfs/conf"
	. "github.com/donnie4w/wfs/httpserver/protocol"
	"github.com/donnie4w/wfs/storge"
	"github.com/julienschmidt/httprouter"
)

func thandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("thandler err:", err)
		}
	}()
	//	ss := strings.Split(r.RemoteAddr, ":")
	//	fmt.Println("ip:", ss[0])
	if "POST" == r.Method {
		protocolFactory := thrift.NewTCompactProtocolFactory()
		transport := thrift.NewStreamTransport(r.Body, w)
		inProtocol := protocolFactory.GetProtocol(transport)
		outProtocol := protocolFactory.GetProtocol(transport)
		processor := NewIWfsProcessor(&ServiceImpl{})
		processor.Process(inProtocol, outProtocol)
	}
}

type ServiceImpl struct {
}

func (t *ServiceImpl) WfsPost(wf *WfsFile) (r *WfsAck, err error) {
	r = NewWfsAck()
	status := int32(200)
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
			status = 500
		}
		r.Status = &status
	}()
	bs := wf.GetFileBody()
	if int64(len(bs)) > CF.MaxFileSize {
		return r, errors.New(fmt.Sprint("too large:", len(bs)))
	}
	err = storge.AppendData(bs, wf.GetName(), wf.GetFileType())
	if err != nil {
		status = 500
	}
	return
}

// 拉取
//
// Parameters:
//  - Name
func (t *ServiceImpl) WfsRead(uri string) (r *WfsFile, err error) {
	r = NewWfsFile()
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
		}
	}()
	name := uri
	arg := ""
	if strings.Contains(uri, "?") {
		index := strings.Index(uri, "?")
		name = uri[:index]
		arg = uri[index:]
	}
	bs, err := storge.GetData(name)
	if err == nil {
		if strings.HasPrefix(arg, "?imageView2") {
			spec := NewSpec(bs, arg)
			r.FileBody = spec.GetData()
		} else {
			r.FileBody = bs
		}
	}
	return
}

// 删除
//
// Parameters:
//  - Name
func (t *ServiceImpl) WfsDel(name string) (r *WfsAck, err error) {
	r = NewWfsAck()
	status := int32(200)
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
			status = 500
		}
		r.Status = &status
	}()
	err = storge.DelData(name)
	if err != nil {
		status = 500
	}
	return
}
