package cache

import (
	"sync"
)

//MapFacade an interface to define a Map
type MapFacade interface {
	//Set : set key-value
	Set(key interface{}, value interface{})
	//Get : get value
	Get(key interface{}) (interface{}, bool)
	//Exist : return true if key in map
	Exist(key interface{}) bool
	//Delete : delete a key
	Delete(key interface{})
}

//Map with locked map
type Map struct {
	m    map[interface{}]interface{}
	lock sync.RWMutex
}

//NewSingleMap create a single map
func NewSingleMap() MapFacade {
	return NewMap()
}

func NewMap() *Map {
	var m = &Map{}
	m.Init()
	return m
}

func (m *Map) Init() {
	m.m = make(map[interface{}]interface{})
}

//Set : set key-value
func (m *Map) Set(key interface{}, value interface{}) {
	m.lock.Lock()
	m.m[key] = value
	m.lock.Unlock()
}

//Get : get value
func (m *Map) Get(key interface{}) (interface{}, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var v, ok = m.m[key]
	return v, ok
}

//Exist : return true if key in map
func (m *Map) Exist(key interface{}) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var _, ok = m.m[key]
	return ok
}

//Delete : delete a key
func (m *Map) Delete(key interface{}) {
	m.lock.Lock()
	delete(m.m, key)
	m.lock.Unlock()
}
