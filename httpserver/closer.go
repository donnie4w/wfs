package httpserver

import (
	"fmt"
	"sync"
	"time"

	. "wfs/conf"
)

type Closer interface {
	Close() error
}

//var pool = NewPool()

type Pool struct {
	pool1 *ClosePool
	pool2 *ClosePool
	b     bool
	mu    *sync.Mutex
}

func (this *Pool) Add(c Closer) {
	this.mu.Lock()
	defer this.mu.Unlock()
	if this.b {
		this.pool1.Add(c)
	} else {
		this.pool2.Add(c)
	}
}

func (this *Pool) _ticker() {
	for {
		if this.b {
			this.pool2._check()
			this.b = false
		} else {
			this.pool1._check()
			this.b = true
		}
		time.Sleep(time.Duration(CF.Keepalive) * time.Second)
	}
}

func NewPool() (p *Pool) {
	p = &Pool{NewClosePool(), NewClosePool(), true, new(sync.Mutex)}
	go p._ticker()
	return
}

type ClosePool struct {
	pool map[Closer]int64
	mu   *sync.RWMutex
}

func NewClosePool() (cp *ClosePool) {
	cp = &ClosePool{make(map[Closer]int64, 0), new(sync.RWMutex)}
	return
}

func (this *ClosePool) Add(c Closer) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.pool[c] = time.Now().Unix()
}

func (this *ClosePool) close(c Closer) {
	fmt.Println("close", c)
	c.Close()
	delete(this.pool, c)
}

func (this *ClosePool) _check() {
	this.mu.RLock()
	defer this.mu.RUnlock()
	for k, v := range this.pool {
		if time.Now().Unix()-v >= int64(CF.Keepalive) {
			this.close(k)
		}
	}
}
