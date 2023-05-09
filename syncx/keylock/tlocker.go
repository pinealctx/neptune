package keylock

import (
	"sync"
)

// TKeyLocker global locker based on key
type TKeyLocker[T comparable] struct {
	lockMap map[T]*wrapLocker
	locker  sync.Mutex
}

// NewTKeyLocker new key locker
func NewTKeyLocker[T comparable]() TLocker[T] {
	return &TKeyLocker[T]{
		lockMap: make(map[T]*wrapLocker),
	}
}

// NewTKeyLockerInstance new key locker instance
func NewTKeyLockerInstance[T comparable]() *TKeyLocker[T] {
	return &TKeyLocker[T]{
		lockMap: make(map[T]*wrapLocker),
	}
}

// Lock write lock
func (d *TKeyLocker[T]) Lock(key T) {
	var (
		wrLocker *wrapLocker
		ok       bool
	)
	d.locker.Lock()
	wrLocker, ok = d.lockMap[key]
	if !ok {
		wrLocker = &wrapLocker{}
		d.lockMap[key] = wrLocker
	}
	wrLocker.writeCount++
	d.locker.Unlock()
	wrLocker.rwLocker.Lock()
}

// Unlock write unlock
func (d *TKeyLocker[T]) Unlock(key T) {
	var wrLocker *wrapLocker

	d.locker.Lock()
	wrLocker = d.lockMap[key]
	wrLocker.rwLocker.Unlock()
	wrLocker.writeCount--
	d.tryFree(key, wrLocker)
	d.locker.Unlock()
}

// RLock read lock
func (d *TKeyLocker[T]) RLock(key T) {
	var (
		wrLocker *wrapLocker
		ok       bool
	)
	d.locker.Lock()
	wrLocker, ok = d.lockMap[key]
	if !ok {
		wrLocker = &wrapLocker{}
		d.lockMap[key] = wrLocker
	}
	wrLocker.readCount++
	d.locker.Unlock()
	wrLocker.rwLocker.RLock()
}

// RUnlock read unlock
func (d *TKeyLocker[T]) RUnlock(key T) {
	var wrLocker *wrapLocker

	d.locker.Lock()
	wrLocker = d.lockMap[key]
	wrLocker.rwLocker.RUnlock()
	wrLocker.readCount--
	d.tryFree(key, wrLocker)
	d.locker.Unlock()
}

// Locks write lock
func (d *TKeyLocker[T]) Locks(keys []T) {
	var ws = d.getWriteLocks(keys)
	for _, wrLocker := range ws {
		wrLocker.rwLocker.Lock()
	}
}

// Unlocks write unlock
func (d *TKeyLocker[T]) Unlocks(keys []T) {
	var wrLocker *wrapLocker

	d.locker.Lock()
	for _, key := range keys {
		wrLocker = d.lockMap[key]
		wrLocker.rwLocker.Unlock()
		wrLocker.writeCount--
		d.tryFree(key, wrLocker)
	}
	d.locker.Unlock()
}

// RLocks read lock
func (d *TKeyLocker[T]) RLocks(keys []T) {
	var ws = d.getReadLocks(keys)
	for _, wrLocker := range ws {
		wrLocker.rwLocker.RLock()
	}
}

// RUnlocks read unlock
func (d *TKeyLocker[T]) RUnlocks(keys []T) {
	var wrLocker *wrapLocker

	d.locker.Lock()
	for _, key := range keys {
		wrLocker = d.lockMap[key]
		wrLocker.rwLocker.RUnlock()
		wrLocker.readCount--
		d.tryFree(key, wrLocker)
	}
	d.locker.Unlock()
}

// get write locks
func (d *TKeyLocker[T]) getWriteLocks(keys []T) []*wrapLocker {
	var (
		wrLocker *wrapLocker
		ok       bool
		ws       []*wrapLocker
	)
	ws = make([]*wrapLocker, len(keys))
	d.locker.Lock()
	for i, key := range keys {
		wrLocker, ok = d.lockMap[key]
		if !ok {
			wrLocker = &wrapLocker{}
			d.lockMap[key] = wrLocker
		}
		wrLocker.writeCount++
		ws[i] = wrLocker
	}
	d.locker.Unlock()
	return ws
}

// get read locks
func (d *TKeyLocker[T]) getReadLocks(keys []T) []*wrapLocker {
	var (
		wrLocker *wrapLocker
		ok       bool
		ws       []*wrapLocker
	)
	ws = make([]*wrapLocker, len(keys))
	d.locker.Lock()
	for i, key := range keys {
		wrLocker, ok = d.lockMap[key]
		if !ok {
			wrLocker = &wrapLocker{}
			d.lockMap[key] = wrLocker
		}
		wrLocker.readCount++
		ws[i] = wrLocker
	}
	d.locker.Unlock()
	return ws
}

// try to free a key locker from map
func (d *TKeyLocker[T]) tryFree(key T, wrLocker *wrapLocker) {
	if wrLocker.readCount == 0 && wrLocker.writeCount == 0 {
		delete(d.lockMap, key)
	}
}
