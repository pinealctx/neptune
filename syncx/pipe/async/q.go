package async

import (
	"container/list"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

var (
	//ErrClosed -- msg queue closed
	ErrClosed = status.Error(codes.Unavailable, "async.queue.closed")
	//ErrFull -- msg queue is full
	ErrFull = status.Error(codes.ResourceExhausted, "async.queue.full")
	//ErrSync it should never gonna happen
	ErrSync = status.Error(codes.OutOfRange, "async.never.gonna.happen.crazy")
)

// Q actor queue structure define
type Q struct {
	//request queue list
	reqList *list.List
	//size queue size, if size is 0, it means it's an infinity queue
	size int
	//is closed
	closed bool

	//queue lock
	lock sync.Mutex
	//queue condition
	cond sync.Cond
}

// NewQ new queue
// size: queue size, if set 0, it means that no limitation of size, it's an infinity queue.
func NewQ(size int) *Q {
	if size < 0 {
		size = 0
	}
	var actorQ = &Q{
		reqList: list.New(),
		size:    size,
	}
	actorQ.cond.L = &actorQ.lock
	return actorQ
}

// IsClosed return queue is closed or not
func (a *Q) IsClosed() bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.closed
}

// Size return queue size
func (a *Q) Size() int {
	return a.size
}

// AddAnyway add normal request to the normal queue end place anyway
// if queue full, sleep then try
func (a *Q) AddAnyway(req interface{}, ts time.Duration) error {
	var err error
	for {
		err = a.Add(req)
		if err == ErrFull {
			time.Sleep(ts)
		} else {
			return err
		}
	}
}

// Add : add normal request to the normal queue end place.
func (a *Q) Add(req interface{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return ErrClosed
	}
	if a.size > 0 {
		if a.reqList.Len() >= a.size {
			return ErrFull
		}
	}
	a.reqList.PushBack(req)
	a.cond.Broadcast()
	return nil
}

// AddPrior add prior request to queue first place.
func (a *Q) AddPrior(req interface{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return ErrClosed
	}
	a.reqList.PushFront(req)
	a.cond.Broadcast()
	return nil
}

// Pop consume an item, if list is empty, it's been blocked
func (a *Q) Pop() (interface{}, error) {
	return a.pop(true)
}

// PopAnyway consume an item like Pop, but it can consume even the queue is closed.
func (a *Q) PopAnyway() (interface{}, error) {
	return a.pop(false)
}

// Close : close the queue
func (a *Q) Close() {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return
	}
	a.closed = true
	a.cond.Broadcast()
}

// pop : pop item
// input: checkClose
// if true  --> when queue is closed, pop can will return error even in case someone is in queue.
// if false --> when queue is closed, the queue also can be pop if anyone is in queue.
func (a *Q) pop(checkClose bool) (interface{}, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for a.reqList.Len() == 0 {
		if a.closed {
			return nil, ErrClosed
		}
		a.cond.Wait()
	}
	if checkClose {
		if a.closed {
			return nil, ErrClosed
		}
	}
	//pop req
	var front = a.reqList.Front()
	if front != nil {
		a.reqList.Remove(front)
		return front.Value, nil
	}
	return nil, ErrSync
}
