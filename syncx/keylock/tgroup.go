package keylock

import (
	"github.com/pinealctx/neptune/remap"
	"golang.org/x/exp/slices"
)

type multiKeyT[T comparable] struct {
	index int
	ks    []T
}

// TKeyLockerGrp wide key locker group
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

// Locks write lock
func (w *TKeyLockerGrp[T]) Locks(keys []T) {
	var ms = w.calculateSortedMultiKeys(keys)
	var ws = make([]*wrapLocker, 0, len(keys))
	for _, ks := range ms {
		ws = append(ws, w.ls[ks.index].getWriteLocks(ks.ks)...)
	}
	for _, wr := range ws {
		wr.rwLocker.Lock()
	}
}

// Unlocks write unlock
func (w *TKeyLockerGrp[T]) Unlocks(keys []T) {
	var m = w.calculateSortedMultiKeys(keys)
	for _, ks := range m {
		w.ls[ks.index].Unlocks(ks.ks)
	}
}

// RLocks read lock
func (w *TKeyLockerGrp[T]) RLocks(keys []T) {
	var ms = w.calculateSortedMultiKeys(keys)
	var ws = make([]*wrapLocker, 0, len(keys))
	for _, ks := range ms {
		ws = append(ws, w.ls[ks.index].getReadLocks(ks.ks)...)
	}
	for _, wr := range ws {
		wr.rwLocker.RLock()
	}
}

// RUnlocks read unlock
func (w *TKeyLockerGrp[T]) RUnlocks(keys []T) {
	var m = w.calculateSortedMultiKeys(keys)
	for _, ks := range m {
		w.ls[ks.index].RUnlocks(ks.ks)
	}
}

// calculate key
func (w *TKeyLockerGrp[T]) calculateKey(key T) *TKeyLocker[T] {
	var i = w.calKeyFn(key)
	return w.ls[i]
}

// calculate multi keys
func (w *TKeyLockerGrp[T]) calculateSortedMultiKeys(keys []T) []multiKeyT[T] {
	// calculate each key hash index
	// and group by index sorted
	var m = make(map[int][]T)
	for _, key := range keys {
		var i = w.calKeyFn(key)
		m[i] = append(m[i], key)
	}

	var ms = make([]multiKeyT[T], 0, len(m))
	for i, ks := range m {
		ms = append(ms, multiKeyT[T]{index: i, ks: ks})
	}

	// sort by index
	slices.SortFunc[multiKeyT[T]](ms, func(a, b multiKeyT[T]) bool {
		return a.index < b.index
	})
	return ms
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
