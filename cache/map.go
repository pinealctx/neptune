package cache

import (
	"sync"
)

//Map with locked map
type Map struct {
	m    map[interface{}]interface{}
	lock sync.RWMutex
}

func NewMap() *Map {
	var m = &Map{}
	m.Init()
	return m
}

func (m *Map) Init() {
	m.m = make(map[interface{}]interface{})
}

func (m *Map) Set(key interface{}, value interface{}) {
	m.lock.Lock()
	m.m[key] = value
	m.lock.Unlock()
}

func (m *Map) Get(key interface{}) (interface{}, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var v, ok = m.m[key]
	return v, ok
}

func (m *Map) Exist(key interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var _, ok = m.m[key]
	return ok
}

func (m *Map) Delete(key interface{}) {
	m.lock.Lock()
	delete(m.m, key)
	m.lock.Unlock()
}
