package keylock

import (
	"sync"
)

// wrap read/write locker
type wrapLocker struct {
	rwLocker   sync.RWMutex
	readCount  int
	writeCount int
}

// KeyLocker global locker based on key
type KeyLocker struct {
	lockMap map[any]*wrapLocker
	locker  sync.Mutex
}

// NewKeyLocker new key locker
func NewKeyLocker() Locker {
	return &KeyLocker{
		lockMap: make(map[any]*wrapLocker),
	}
}

// NewKeyLockerInstance new key locker instance
func NewKeyLockerInstance() *KeyLocker {
	return &KeyLocker{
		lockMap: make(map[any]*wrapLocker),
	}
}

// Lock write lock
func (d *KeyLocker) Lock(key any) {
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
func (d *KeyLocker) Unlock(key any) {
	var (
		wrLocker *wrapLocker
	)
	d.locker.Lock()
	wrLocker = d.lockMap[key]
	wrLocker.rwLocker.Unlock()
	wrLocker.writeCount--
	d.tryFree(key, wrLocker)
	d.locker.Unlock()
}

// RLock read lock
func (d *KeyLocker) RLock(key any) {
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
func (d *KeyLocker) RUnlock(key any) {
	var (
		wrLocker *wrapLocker
	)
	d.locker.Lock()
	wrLocker = d.lockMap[key]
	wrLocker.rwLocker.RUnlock()
	wrLocker.readCount--
	d.tryFree(key, wrLocker)
	d.locker.Unlock()
}

// try to free a key locker from map
func (d *KeyLocker) tryFree(key any, wrLocker *wrapLocker) {
	if wrLocker.readCount == 0 && wrLocker.writeCount == 0 {
		delete(d.lockMap, key)
	}
}
