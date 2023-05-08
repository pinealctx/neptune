package keylock

import (
	"github.com/pinealctx/neptune/remap"
)

type TKeyLockerGrp[T comparable] struct {
	ls       []*TKeyLocker[T]
	calKeyFn func(key interface{}) int
	rehash   *remap.ReMap
}

// Lock write lock
func (w *TKeyLockerGrp[T]) Lock(key T) {
	w.calculateKey(key).Lock(key)
}

// Unlock write unlock
func (w *TKeyLockerGrp[T]) Unlock(key T) {
	w.calculateKey(key).Unlock(key)
}

// RLock read lock
func (w *TKeyLockerGrp[T]) RLock(key T) {
	w.calculateKey(key).RLock(key)
}

// RUnlock read unlock
func (w *TKeyLockerGrp[T]) RUnlock(key T) {
	w.calculateKey(key).RUnlock(key)
}

// calculate key
func (w *TKeyLockerGrp[T]) calculateKey(key T) *TKeyLocker[T] {
	var i = w.calKeyFn(key)
	return w.ls[i]
}

// NewTKeyLockeGrp new wide key locker group
func NewTKeyLockeGrp[T comparable](opts ...remap.Option) TLocker[T] {
	return newTKeyLockeGrp[T](false, opts...)
}

// NewTXHashTKeyLockeGrp new wide key locker group
func NewTXHashTKeyLockeGrp[T comparable](opts ...remap.Option) TLocker[T] {
	return newTKeyLockeGrp[T](true, opts...)
}

// newTKeyLockeGrp new wide key locker group
func newTKeyLockeGrp[T comparable](useXHash bool, opts ...remap.Option) TLocker[T] {
	var w = &TKeyLockerGrp[T]{}
	w.rehash = remap.NewReMap(opts...)
	var numbs = w.rehash.Numbs()
	w.ls = make([]*TKeyLocker[T], numbs)
	for i := uint64(0); i < numbs; i++ {
		w.ls[i] = NewTKeyLockerInstance[T]()
	}
	if useXHash {
		w.calKeyFn = w.rehash.XHashIndex
	} else {
		w.calKeyFn = w.rehash.SimpleIndex
	}
	return w
}
