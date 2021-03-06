package keylock

type Locker interface {
	//Lock write lock
	Lock(key interface{})
	//Unlock write unlock
	Unlock(key interface{})
	//RLock read lock
	RLock(key interface{})
	//RULock read unlock
	RULock(key interface{})
}
