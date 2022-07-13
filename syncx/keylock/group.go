package keylock

import "github.com/pinealctx/neptune/remap"

type KeyLockerGrp struct {
	ls       []*KeyLocker
	calKeyFn func(key interface{}) int
	rehash   *remap.ReMap
}

//Lock write lock
func (w *KeyLockerGrp) Lock(key interface{}) {
	w.calculateKey(key).Lock(key)
}

//Unlock write unlock
func (w *KeyLockerGrp) Unlock(key interface{}) {
	w.calculateKey(key).Unlock(key)
}

//RLock read lock
func (w *KeyLockerGrp) RLock(key interface{}) {
	w.calculateKey(key).RLock(key)
}

//RULock read unlock
func (w *KeyLockerGrp) RULock(key interface{}) {
	w.calculateKey(key).RULock(key)
}

//calculate key
func (w *KeyLockerGrp) calculateKey(key interface{}) *KeyLocker {
	var i = w.calKeyFn(key)
	return w.ls[i]
}

//NewKeyLockeGrp new wide key locker group
func NewKeyLockeGrp(opts ...remap.Option) Locker {
	return newKeyLockeGrp(false, opts...)
}

//NewXHashKeyLockeGrp new wide key locker group
func NewXHashKeyLockeGrp(opts ...remap.Option) Locker {
	return newKeyLockeGrp(true, opts...)
}

//newKeyLockeGrp new wide key locker group
func newKeyLockeGrp(useXHash bool, opts ...remap.Option) Locker {
	var w = &KeyLockerGrp{}
	w.rehash = remap.NewReMap(opts...)
	var numbs = w.rehash.Numbs()
	w.ls = make([]*KeyLocker, numbs)
	for i := uint64(0); i < numbs; i++ {
		w.ls[i] = NewKeyLockerInstance()
	}
	if useXHash {
		w.calKeyFn = w.rehash.XHashIndex
	} else {
		w.calKeyFn = w.rehash.SimpleIndex
	}
	return w
}
