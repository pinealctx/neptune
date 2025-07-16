package cache

import (
	"sync"
)

// Map with locked map
type Map struct {
	m    map[any]any
	lock sync.RWMutex
}

// NewSingleMap create a single map
func NewSingleMap() MapFacade {
	return NewMap()
}

func NewMap() *Map {
	var m = &Map{}
	m.Init()
	return m
}

func (m *Map) Init() {
	m.m = make(map[any]any)
}

// Set : set key-value
func (m *Map) Set(key any, value any) {
	m.lock.Lock()
	m.m[key] = value
	m.lock.Unlock()
}

// Get : get value
func (m *Map) Get(key any) (any, bool) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var v, ok = m.m[key]
	return v, ok
}

// Exist : return true if key in map
func (m *Map) Exist(key any) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()
	var _, ok = m.m[key]
	return ok
}

// Delete : delete a key
func (m *Map) Delete(key any) {
	m.lock.Lock()
	delete(m.m, key)
	m.lock.Unlock()
}
