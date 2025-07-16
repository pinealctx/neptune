package nano

import (
	"sync"
	"time"
)

// UnixNanoID use unix nano as id, but it must be increment.
type UnixNanoID struct {
	current int64
	sync.Mutex
}

// NewUnixNanoID new unix nano id
func NewUnixNanoID(current int64) *UnixNanoID {
	return &UnixNanoID{
		current: current,
	}
}

// GenID generate id
func (n *UnixNanoID) GenID() int64 {
	var ts = time.Now().UnixNano()
	return n.GenIDByTS(ts)
}

// GenIDByTS swap, if bigger than current, return it
func (n *UnixNanoID) GenIDByTS(ts int64) int64 {
	n.Lock()
	if ts > n.current {
		n.current = ts
	} else {
		n.current++
		ts = n.current
	}
	n.Unlock()
	return ts
}

// UnixNanoNoLockID lock free -- without lock
// control by caller user
type UnixNanoNoLockID struct {
	current int64
}

// NewUnixNanoNoLockID new unix nano id -- without lock
func NewUnixNanoNoLockID(current int64) *UnixNanoNoLockID {
	return &UnixNanoNoLockID{
		current: current,
	}
}

// GenID generate id
func (n *UnixNanoNoLockID) GenID() int64 {
	var ts = time.Now().UnixNano()
	return n.GenIDByTS(ts)
}

// GenIDByTS swap, if bigger than current, return it
func (n *UnixNanoNoLockID) GenIDByTS(ts int64) int64 {
	if ts > n.current {
		n.current = ts
	} else {
		n.current++
		ts = n.current
	}
	return ts
}
