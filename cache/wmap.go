package cache

import (
	"github.com/pinealctx/neptune/remap"
)

// WideMap use Map group array as a wide map
type WideMap struct {
	ms       []*Map
	calKeyFn func(key any) int
	rehash   *remap.ReMap
}

// NewWideMap new wide Map
func NewWideMap(opts ...remap.Option) MapFacade {
	return newWideMap(false, opts...)
}

// NewWideXHashMap new wide Map use xxhash as group
func NewWideXHashMap(opts ...remap.Option) MapFacade {
	return newWideMap(true, opts...)
}

// newWideMap new wide map
func newWideMap(useXHash bool, opts ...remap.Option) MapFacade {
	var w = &WideMap{}
	w.rehash = remap.NewReMap(opts...)
	var numbs = w.rehash.Numbs()
	w.ms = make([]*Map, numbs)
	for i := uint64(0); i < numbs; i++ {
		w.ms[i] = NewMap()
	}
	if useXHash {
		w.calKeyFn = w.rehash.XHashIndex
	} else {
		w.calKeyFn = w.rehash.SimpleIndex
	}
	return w
}

// Set : set key-value
func (w *WideMap) Set(key any, value any) {
	w.calculateKey(key).Set(key, value)
}

// Get : get value
func (w *WideMap) Get(key any) (any, bool) {
	return w.calculateKey(key).Get(key)
}

// Exist : return true if key in map
func (w *WideMap) Exist(key any) bool {
	return w.calculateKey(key).Exist(key)
}

// Delete : delete a key
func (w *WideMap) Delete(key any) {
	w.calculateKey(key).Delete(key)
}

// calculate key
func (w *WideMap) calculateKey(key any) *Map {
	var i = w.calKeyFn(key)
	return w.ms[i]
}
