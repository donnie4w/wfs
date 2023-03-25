package httpserver

/**
 * Copyright 2017 wfs Author. All Rights Reserved.
 * email: donnie4w@gmail.com
 */
import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	. "wfs/conf"
	"wfs/storge"

	"github.com/donnie4w/simplelog/logging"
)

var Factory *SlaveFactory

func Init() {
	Factory = NewSlaveFactory()
	Factory._initSlaves()
}

type SlaveBean struct {
	Name        string
	Addr        string
	Weight      int32
	WeightCount int32
}

func NewSlaveBean(name, addr string, weight int32) *SlaveBean {
	return &SlaveBean{Name: name, Addr: addr, Weight: weight, WeightCount: weight}
}

type SlaveFactory struct {
	mu         *sync.RWMutex
	slaveMap   map[string]*SlaveBean
	slavevalid map[string]byte
	slaveSlb   map[string]byte
	readint    int32
}

func NewSlaveFactory() (sf *SlaveFactory) {
	sf = &SlaveFactory{new(sync.RWMutex), make(map[string]*SlaveBean, 0), make(map[string]byte, 0), make(map[string]byte, 0), 0}
	sf.addSlave("master", "")
	go sf.hbtask()
	return
}

func (this *SlaveFactory) _initSlaves() {
	slavemap := storge.SlaveList()
	if slavemap != nil {
		for _, v := range slavemap {
			sb := _Decoder(v)
			this._puts(&sb)
		}
	}
}

func (this *SlaveFactory) lenght() int {
	this.mu.RLock()
	defer this.mu.RUnlock()
	return len(this.slaveMap)
}

func (this *SlaveFactory) addSlave(name, addr string) (err error) {
	if strings.HasPrefix(name, "_") {
		err = errors.New("can't start with '_'")
		return
	}
	if _, ok := this.slaveMap[name]; ok {
		err = errors.New("exist")
		return
	}
	sb, err := this._addSlave(name, addr, 10)
	if err == nil {
		storge.SaveSlave(name, _Encoder(*sb))
	}
	return
}

func (this *SlaveFactory) _addSlave(name, addr string, weight int32) (sb *SlaveBean, err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if name != "master" {
		_, err = this.ping(addr)
		if err != nil {
			logging.Error("ping ", addr, " err:", err)
			return
		}
	} else {
		addr = fmt.Sprint("127.0.0.1:", CF.Port)
	}
	sb = NewSlaveBean(name, addr, weight)
	this.slaveMap[name] = sb
	this.slavevalid[name] = 0
	this.slaveSlb[name] = 0
	return
}

func (this *SlaveFactory) _puts(sb *SlaveBean) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.slaveMap[sb.Name] = sb
	_, err = this.ping(sb.Addr)
	if err == nil {
		this.slavevalid[sb.Name] = 0
	}
	return
}

func (this *SlaveFactory) setWeight(name string, weight int32) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if sl, ok := this.slaveMap[name]; ok {
		sl.Weight = weight
		storge.SaveSlave(name, _Encoder(*sl))
	}
}

func (this *SlaveFactory) remove(name string) (err error) {
	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.slaveMap, name)
	delete(this.slavevalid, name)
	delete(this.slaveSlb, name)
	storge.DelSlave(name)
	return
}

func (this *SlaveFactory) _invalid(name string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.slavevalid, name)
}

func (this *SlaveFactory) getAddrByName(name string) (addr string) {
	this.mu.RLock()
	defer this.mu.RUnlock()
	if sb, ok := this.slaveMap[name]; ok {
		return sb.Addr
	}
	return
}

func (this *SlaveFactory) slavelist() (s string) {
	this.mu.RLock()
	defer this.mu.RUnlock()
	for _, slave := range this.slaveMap {
		_, ok := this.slavevalid[slave.Name]
		valid := "1"
		if !ok {
			valid = "0"
		}
		s = fmt.Sprintln(s, slave.Name, " ", slave.Addr, " ", slave.Weight, " ", valid)
	}
	return
}

func (this *SlaveFactory) ping(addr string) (s string, err error) {
	_, er := wfsPing(addr)
	if er == nil {
		s = _200
	} else {
		s = _404
		err = errors.New(fmt.Sprint(addr, " ping failed"))
	}
	return
}

func (this *SlaveFactory) getSlaveByWeight() (sb *SlaveBean) {
	this.mu.Lock()
	defer this.mu.Unlock()
	atomic.AddInt32(&this.readint, 1)
	atomic.CompareAndSwapInt32(&this.readint, 1<<31-1, 1)
LOOP:
	if len(this.slaveSlb) > 0 {
		c := this.readint % int32(len(this.slaveSlb))
		i := int32(0)
		name := ""
		for k, _ := range this.slaveSlb {
			if i == c {
				name = k
				break
			}
			i++
		}
		delete(this.slaveSlb, name)
		if _, ok := this.slavevalid[name]; !ok {
			goto LOOP
		}
		if _sb, ok := this.slaveMap[name]; ok {
			if _sb.WeightCount <= 0 {
				goto LOOP
			} else {
				atomic.AddInt32(&_sb.WeightCount, -1)
				return _sb
			}
		} else {
			goto LOOP
		}
	} else {
		for k, _ := range this.slavevalid {
			if v, ok := this.slaveMap[k]; ok {
				if ok && v.WeightCount > 0 {
					this.slaveSlb[k] = 0
				}
			}
		}

		if len(this.slaveSlb) > 0 {
			goto LOOP
		} else if len(this.slaveMap) > 0 {
			for k, v := range this.slaveMap {
				v.WeightCount = v.Weight
				if _, ok := this.slavevalid[k]; ok {
					this.slaveSlb[k] = 0
				}
			}
			if len(this.slaveSlb) > 0 {
				goto LOOP
			}
		}
	}
	return
}

func AppendData(bs []byte, name string, fileType string) (err error) {
	sb := Factory.getSlaveByWeight()
	if sb == nil || sb.Name == "master" {
		err = storge.AppendData(bs, name, fileType, "")
	} else {
		_, err = wfsPost(sb.Addr, bs, name, fileType)
		if err == nil {
			err = storge.AppendData(bs, name, fileType, sb.Name)
		} else {
			Factory._invalid(sb.Name)
			err = storge.AppendData(bs, name, fileType, "")
		}
	}
	return
}

func GetData(uri string) (retbs []byte, err error) {
	return getDataByName(uri[3:])
}

func _getIndexFromStr(s, _i string) (i int) {
	for i = len(s) - 1; i >= 0; i-- {
		if s[i:i+1] == _i {
			return
		}
	}
	return
}

func getDataByName(uri string) (retbs []byte, err error) {
	uri3 := uri
	name := uri3
	arg := ""
	if index := _getIndexFromStr(uri3, "?"); index > 0 {
		name = uri3[:index]
		arg = uri3[index:]
	}
	bs, shardname, err := storge.GetData(name)
	if err == nil && bs != nil {
		if strings.HasPrefix(arg, "?imageView2") {
			spec := NewSpec(bs, arg)
			retbs = spec.GetData()
		} else {
			retbs = bs
		}
	} else if len(shardname) > 0 {
		addr := Factory.getAddrByName(shardname)
		if addr != "" {
			wf, er := wfsRead(addr, uri)
			if er == nil {
				retbs = wf.GetFileBody()
			} else {
				err = er
			}
		} else {
			fmt.Println("err:", shardname, " is not exist")
		}
	}
	if err != nil {
		err = errors.New(_404)
	}
	return
}

func DelData(name string) (er error) {
	shardname, err := storge.DelData(name)
	if err == nil {
		if len(shardname) > 0 {
			addr := Factory.getAddrByName(shardname)
			_, er = wfsDel(addr, name)
		}
	}
	return
}

/******************************************************************************/
func _Decoder(data []byte) (sb SlaveBean) {
	var network bytes.Buffer
	_, er := network.Write(data)
	dec := gob.NewDecoder(&network)
	if er == nil {
		er = dec.Decode(&sb)
	}
	return
}

func _Encoder(sb SlaveBean) (bs []byte) {
	var network bytes.Buffer
	enc := gob.NewEncoder(&network)
	enc.Encode(sb)
	bs = network.Bytes()
	return
}

/******************************************************************************/

func (this *SlaveFactory) addHbs(name string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.slavevalid[name] = 0
}

func (this *SlaveFactory) delHbs(name string) {
	this.mu.Lock()
	defer this.mu.Unlock()
	delete(this.slavevalid, name)
}

func (this *SlaveFactory) getAll() (sbs []*SlaveBean) {
	this.mu.Lock()
	defer this.mu.Unlock()
	sbs = make([]*SlaveBean, 0)
	for _, v := range this.slaveMap {
		sbs = append(sbs, v)
	}
	return
}

func (this *SlaveFactory) hbtask() {
	ticker := time.NewTicker(15 * time.Second)
	for {
		select {
		case <-ticker.C:
			sbs := this.getAll()
			for _, v := range sbs {
				_, err := this.ping(v.Addr)
				if err == nil {
					this.addHbs(v.Name)
				} else {
					this.delHbs(v.Name)
				}
			}
		}
	}
}
