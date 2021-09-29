package q

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

var (
	//ErrClosed close
	ErrClosed = errors.New("pipe.q.closed")

	//ErrReqQFull req q full
	ErrReqQFull = errors.New("pipe.q.req.full")

	//ErrSync never gonna happen
	ErrSync = errors.New("never.gonna.happen.crazy")
)

//option for queue
type _Option struct {
	reqMaxNum int
}

//Option : option function
type Option func(o *_Option)

//WithSize setup max queue number of request queue
//if max is 0, which means no limit
func WithSize(num int) Option {
	return func(o *_Option) {
		o.reqMaxNum = num
	}
}

//Q actor queue structure define
type Q struct {
	//request queue list
	reqList *list.List
	//request pipe size max number
	reqMaxNum int
	//is closed
	closed bool

	//queue lock
	lock sync.Mutex
	//queue condition
	cond sync.Cond
}

//NewQ new queue
func NewQ(options ...Option) *Q {
	var actorQ = &Q{
		reqList: list.New(),
	}
	var option = &_Option{}
	for _, opt := range options {
		opt(option)
	}
	if option.reqMaxNum > 0 {
		actorQ.reqMaxNum = option.reqMaxNum
	}
	actorQ.cond.L = &actorQ.lock
	return actorQ
}

//AddReqAnyway add normal request to the normal queue end place anyway
//if queue full, sleep then try
func (a *Q) AddReqAnyway(req interface{}, ts time.Duration) error {
	var err error
	for {
		err = a.AddReq(req)
		if err == ErrReqQFull {
			time.Sleep(ts)
		} else {
			return err
		}
	}
}

//AddReq add normal request to the normal queue end place.
func (a *Q) AddReq(req interface{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return ErrClosed
	}
	if a.reqMaxNum > 0 {
		if a.reqList.Len() >= a.reqMaxNum {
			return ErrReqQFull
		}
	}
	a.reqList.PushBack(req)
	a.cond.Broadcast()
	return nil
}

//AddPriorReq add normal request to the normal queue first place.
func (a *Q) AddPriorReq(req interface{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return ErrClosed
	}
	a.reqList.PushFront(req)
	a.cond.Broadcast()
	return nil
}

//Pop consume an item, if list is empty, it's been blocked
func (a *Q) Pop() (interface{}, error) {
	return a.pop(true)
}

//PopAnyway consume an item like Pop, but it can consume even the queue is closed.
func (a *Q) PopAnyway() (interface{}, error) {
	return a.pop(false)
}

//Close : close the queue
func (a *Q) Close() {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return
	}
	a.closed = true
	a.cond.Broadcast()
}

//pop : pop item
//input: checkClose
//if true  --> when queue is closed, pop can will return error even in case someone is in queue.
//if false --> when queue is closed, the queue also can be pop if anyone is in queue.
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
