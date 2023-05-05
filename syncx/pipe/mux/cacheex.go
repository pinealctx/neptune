package mux

import (
	"github.com/pinealctx/neptune/cache"
)

// FacadeMap : extend cache.Map to interface CacheFacade
type FacadeMap struct {
	cache.Map
}

func NewFacadeMap() CacheFacade {
	var m = &FacadeMap{}
	m.Init()
	return m
}

func (m *FacadeMap) Peek(key interface{}) (interface{}, bool) {
	return m.Get(key)
}

type _wrapper struct {
	v interface{}
}

func (w _wrapper) Size() int {
	var sizeV, ok = w.v.(cache.Value)
	if ok {
		return sizeV.Size()
	}
	return 1
}

// FacadeLRU : extend LRUCache to interface CacheFacade
type FacadeLRU struct {
	cache.LRUCache
}

// NewFacadeLRU : creates a new empty cache with the given capacity.
func NewFacadeLRU(capacity int64) CacheFacade {
	var m = &FacadeLRU{}
	m.Init(capacity)
	return m
}

// Peek : only useful in lru cache, peek means no update LRU order.
func (m *FacadeLRU) Peek(key interface{}) (interface{}, bool) {
	var w, ok = m.LRUCache.Peek(key)
	if !ok {
		return nil, false
	}
	return w.(_wrapper).v, true
}

// Get : get from cache, in lru cache, key order also be updated.
func (m *FacadeLRU) Get(key interface{}) (interface{}, bool) {
	var w, ok = m.LRUCache.Get(key)
	if !ok {
		return nil, false
	}
	return w.(_wrapper).v, true
}

// Set : set to cache
func (m *FacadeLRU) Set(key interface{}, value interface{}) {
	m.LRUCache.Set(key, _wrapper{v: value})
}

// Delete : delete key from cache
func (m *FacadeLRU) Delete(key interface{}) {
	m.LRUCache.Delete(key)
}
