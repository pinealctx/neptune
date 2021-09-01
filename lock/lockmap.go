package lock

import "sync"

type SimpleLockDocker struct {
	lockMap map[interface{}]*sync.Mutex
	locker  sync.Mutex
}

func NewSimpleLockDocker() *SimpleLockDocker {
	return &SimpleLockDocker{
		lockMap: make(map[interface{}]*sync.Mutex),
	}
}

func (d *SimpleLockDocker) GetLock(key interface{}) *sync.Mutex {
	d.locker.Lock()
	defer d.locker.Unlock()
	var locker, ok = d.lockMap[key]
	if ok {
		return locker
	}
	locker = &sync.Mutex{}
	d.lockMap[key] = locker
	return locker
}
