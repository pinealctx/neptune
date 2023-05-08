package tiny

import (
	"github.com/pinealctx/neptune/remap"
)

// WideLRUCache use LruCache group array as a wide lru cache
type WideLRUCache struct {
	ls       []*LRUCache
	calKeyFn func(key interface{}) int
	rehash   *remap.ReMap
}

// NeWideLRU new wide lru cache
func NeWideLRU(capacity int64, opts ...remap.Option) LRU {
	return newWideLRUCache(capacity, false, opts...)
}

// NewWideXHashLRU new wide lru cache use xxhash as group
func NewWideXHashLRU(capacity int64, opts ...remap.Option) LRU {
	return newWideLRUCache(capacity, true, opts...)
}

// newWideLRUCache new wide lru cache group
func newWideLRUCache(capacity int64, useXHash bool, opts ...remap.Option) LRU {
	var w = &WideLRUCache{}
	w.rehash = remap.NewReMap(opts...)
	var numbs = w.rehash.Numbs()
	w.ls = make([]*LRUCache, numbs)
	var pSize = capacity/int64(numbs) + 1
	for i := uint64(0); i < numbs; i++ {
		w.ls[i] = NewLRUCache(pSize)
	}
	if useXHash {
		w.calKeyFn = w.rehash.XHashIndex
	} else {
		w.calKeyFn = w.rehash.SimpleIndex
	}
	return w
}

// Get returns a value from the cache, and marks the entry as most recently used.
func (w *WideLRUCache) Get(key interface{}) (v interface{}, ok bool) {
	return w.calculateKey(key).Get(key)
}

// Peek returns a value from the cache without changing the LRU order.
func (w *WideLRUCache) Peek(key interface{}) (v interface{}, ok bool) {
	return w.calculateKey(key).Peek(key)
}

// Exist : return true if key in map
func (w *WideLRUCache) Exist(key interface{}) bool {
	return w.calculateKey(key).Exist(key)
}

// Set sets a value in the cache.
func (w *WideLRUCache) Set(key interface{}, value interface{}) {
	w.calculateKey(key).Set(key, value)
}

// Delete removes an entry from the cache, and returns if the entry existed.
func (w *WideLRUCache) Delete(key interface{}) bool {
	return w.calculateKey(key).Delete(key)
}

// calculate key
func (w *WideLRUCache) calculateKey(key interface{}) *LRUCache {
	var i = w.calKeyFn(key)
	return w.ls[i]
}
