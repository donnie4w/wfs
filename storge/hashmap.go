package storge

import (
	"sync"
)

type hashmap struct {
	lock  *sync.RWMutex
	cache map[interface{}]interface{}
}

func NewHashMap() *hashmap {
	return &hashmap{lock: new(sync.RWMutex), cache: make(map[interface{}]interface{}, 0)}
}

func (this *hashmap) Put(key, value interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	this.cache[key] = value
}
func (this *hashmap) Del(key interface{}) {
	this.lock.Lock()
	defer this.lock.Unlock()
	delete(this.cache, key)
}
func (this *hashmap) Get(key interface{}) (v interface{}, ok bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	v, ok = this.cache[key]
	return
}
