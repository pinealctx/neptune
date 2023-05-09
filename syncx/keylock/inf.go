package keylock

type Locker interface {
	// Lock write lock
	Lock(key interface{})
	// Unlock write unlock
	Unlock(key interface{})
	// RLock read lock
	RLock(key interface{})
	// RUnlock read unlock
	RUnlock(key interface{})
}

type TLocker[T comparable] interface {
	// Lock write lock
	Lock(key T)
	// Unlock write unlock
	Unlock(key T)
	// RLock read lock
	RLock(key T)
	// RUnlock read unlock
	RUnlock(key T)

	// Locks write lock
	Locks(keys []T)
	// Unlocks write unlock
	Unlocks(keys []T)
	// RLocks read lock
	RLocks(keys []T)
	// RUnlocks read unlock
	RUnlocks(keys []T)
}
