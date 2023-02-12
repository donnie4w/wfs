/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	. "wfs/conf"

	. "wfs/httpserver/protocol"

	"github.com/apache/thrift/lib/go/thrift"
	// "github.com/julienschmidt/httprouter"
)

// func thandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
func thandler(w http.ResponseWriter, r *http.Request) {
	defer myRecover()
	if "POST" == r.Method {
		protocolFactory := thrift.NewTCompactProtocolFactory()
		transport := thrift.NewStreamTransport(r.Body, w)
		inProtocol := protocolFactory.GetProtocol(transport)
		outProtocol := protocolFactory.GetProtocol(transport)
		processor := NewIWfsProcessor(&ServiceImpl{r.RemoteAddr[:strings.Index(r.RemoteAddr, ":")]})
		processor.Process(context.Background(), inProtocol, outProtocol)
	}
}

type ServiceImpl struct {
	ip string
}

func (t *ServiceImpl) WfsPost(ctx context.Context, wf *WfsFile) (r *WfsAck, err error) {
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
	//	err = storge.AppendData(bs, wf.GetName(), wf.GetFileType())
	err = AppendData(bs, wf.GetName(), wf.GetFileType())
	if err != nil {
		status = 500
	}
	return
}

// Parameters:
//  - Name
func (t *ServiceImpl) WfsRead(ctx context.Context, uri string) (r *WfsFile, err error) {
	r = NewWfsFile()
	defer myRecover()
	//	name := uri
	//	arg := ""
	//	if strings.Contains(uri, "?") {
	//		index := strings.Index(uri, "?")
	//		name = uri[:index]
	//		arg = uri[index:]
	//	}
	//	bs, err := storge.GetData(name)
	//	if err == nil {
	//		if strings.HasPrefix(arg, "?imageView2") {
	//			spec := NewSpec(bs, arg)
	//			r.FileBody = spec.GetData()
	//		} else {
	//			r.FileBody = bs
	//		}
	//	}
	r.FileBody, err = getDataByName(uri)
	return
}

// Parameters:
//  - Name
func (t *ServiceImpl) WfsDel(ctx context.Context, name string) (r *WfsAck, err error) {
	r = NewWfsAck()
	status := int32(200)
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
			status = 500
		}
		r.Status = &status
	}()
	//	err = storge.DelData(name)
	err = DelData(name)
	if err != nil {
		status = 500
	}
	return
}

// Parameters:
//  - Wc
func (t *ServiceImpl) WfsCmd(ctx context.Context, wc *WfsCmd) (r *WfsAck, err error) {
	if t.ip != "127.0.0.1" {
		return nil, errors.New(_403)
	}
	cmdkey := wc.GetCmdKey()
	cmdvalue := wc.GetCmdValue()
	//	fmt.Println("wfscmd:", cmdkey, " , ", cmdvalue)
	r = NewWfsAck()
	ret := _200
	switch CmdType(cmdkey) {
	case setWeight:
		ss := strings.Split(cmdvalue, ":")
		w, er := strconv.Atoi(ss[1])
		if er == nil {
			Factory.setWeight(ss[0], int32(w))
		} else {
			ret = er.Error()
		}
	case addSlave:
		er := Factory.addSlave(cmdvalue[:strings.Index(cmdvalue, ":")], cmdvalue[strings.Index(cmdvalue, ":")+1:])
		if er != nil {
			ret = er.Error()
		}
	case cutOff:
	case removeSlave:
		er := Factory.remove(cmdvalue)
		if er != nil {
			ret = er.Error()
		}
	case slavelist:
		ret = Factory.slavelist()
	case ping:
		ret, _ = Factory.ping(cmdvalue)
	default:
		ret = _404
	}
	r.Desc = &ret
	return
}

func wfsCmd(cmdkey CmdType, cmdvalue string) (er error) {
	wc := NewWfsCmd()
	ck := string(cmdkey)
	wc.CmdKey, wc.CmdValue = &ck, &cmdvalue
	httpPostClient(fmt.Sprint("http://127.0.0.1:", CF.Port, "/thrift"), 15000, func(client *IWfsClient) {
		wa, err := client.WfsCmd(context.Background(), wc)
		if wa != nil {
			fmt.Println(wa.GetDesc())
		} else if err != nil {
			er = err
		}
	})
	return
}

func wfsPost(addr string, bs []byte, filename string, fileType string) (wa *WfsAck, err error) {
	wf := NewWfsFile()
	wf.FileBody, wf.FileType, wf.Name = bs, &fileType, &filename
	httpPostClient(fmt.Sprint("http://", addr, "/thrift"), 5000, func(client *IWfsClient) {
		wa, err = client.WfsPost(context.Background(), wf)
	})
	return
}

func wfsRead(addr string, uri string) (wf *WfsFile, err error) {
	httpPostClient(fmt.Sprint("http://", addr, "/thrift"), 5000, func(client *IWfsClient) {
		wf, err = client.WfsRead(context.Background(), uri)
	})
	return
}

func wfsDel(addr string, name string) (wa *WfsAck, err error) {
	httpPostClient(fmt.Sprint("http://", addr, "/thrift"), 5000, func(client *IWfsClient) {
		wa, err = client.WfsDel(context.Background(), name)
	})
	return
}

func wfsPing(addr string) (wa *WfsAck, err error) {
	wc := NewWfsCmd()
	ck := string(ping)
	wc.CmdKey = &ck
	httpPostClient(fmt.Sprint("http://", addr, "/thrift"), 5000, func(client *IWfsClient) {
		wa, err = client.WfsCmd(context.Background(), wc)
	})
	return
}

func httpPostClient(urlstr string, timeout int64, f func(*IWfsClient)) (err error) {
	defer func() {
		if er := recover(); er != nil {
			err = errors.New(fmt.Sprint(er))
		}
	}()
	protocolFactory := thrift.NewTCompactProtocolFactory()
	if err != nil {
		return err
	}
	// transport := &thrift.NewTHttpClient{nil, parsedURL, bytes.NewBuffer(buf), http.Header{}, timeout}
	transport, _ := thrift.NewTHttpClientWithOptions(urlstr, thrift.THttpClientOptions{&http.Client{Timeout: time.Duration(timeout)}})
	client := NewIWfsClientFactory(transport, protocolFactory)
	if err := transport.Open(); err != nil {
		return err
	}
	defer transport.Close()
	f(client)
	return
}
