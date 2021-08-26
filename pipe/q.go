package pipe

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"
)

var (
	//ErrClosed close
	ErrClosed = errors.New("pipe.q.closed")

	//ErrCtrlQFull ctrl q full
	ErrCtrlQFull = errors.New("pipe.q.ctrl.full")
	//ErrReqQFull req q full
	ErrReqQFull = errors.New("pipe.q.req.full")

	//ErrSync never gonna happen
	ErrSync = errors.New("never.gonna.happen.crazy")
)

//option for actor queue
type _QOption struct {
	ctrlMaxNum int
	reqMaxNum  int
}

//QOption option function
type QOption func(o *_QOption)

//WithQCtrlSize setup max queue number of control queue
//if max is 0, which means no limit
func WithQCtrlSize(num int) QOption {
	return func(o *_QOption) {
		o.ctrlMaxNum = num
	}
}

//WithQReqSize setup max queue number of request queue
//if max is 0, which means no limit
func WithQReqSize(num int) QOption {
	return func(o *_QOption) {
		o.reqMaxNum = num
	}
}

//Q actor queue structure define
type Q struct {
	//control queue list
	ctrlList *list.List
	//request queue list
	reqList *list.List

	//stop channel
	stopChan chan struct{}
	//clear channel
	clearChan chan struct{}

	//control pipe size max number
	ctrlMaxNum int
	//request pipe size max number
	reqMaxNum int

	//is closed
	closed bool
	//is cleared
	cleared bool

	//queue lock
	lock sync.Mutex
	//queue condition
	cond sync.Cond
}

//NewQ new queue
func NewQ(options ...QOption) *Q {
	var actorQ = &Q{
		ctrlList:  list.New(),
		reqList:   list.New(),
		stopChan:  make(chan struct{}),
		clearChan: make(chan struct{}),
	}
	var option = &_QOption{}
	for _, opt := range options {
		opt(option)
	}
	if option.ctrlMaxNum > 0 {
		actorQ.ctrlMaxNum = option.ctrlMaxNum
	}
	if option.reqMaxNum > 0 {
		actorQ.reqMaxNum = option.reqMaxNum
	}
	actorQ.cond.L = &actorQ.lock
	return actorQ
}

//AddCtrlAnyway dd control request to the control queue end place anyway
//if queue full, sleep then try
func (a *Q) AddCtrlAnyway(cmd interface{}, ts time.Duration) error {
	var err error
	for {
		err = a.AddCtrl(cmd)
		if err == ErrCtrlQFull {
			time.Sleep(ts)
		} else {
			return err
		}
	}
}

//AddCtrl add control request to the control queue end place.
func (a *Q) AddCtrl(cmd interface{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return ErrClosed
	}
	if a.ctrlMaxNum > 0 {
		if a.ctrlList.Len() >= a.ctrlMaxNum {
			return ErrCtrlQFull
		}
	}
	a.ctrlList.PushBack(cmd)
	a.cond.Broadcast()
	return nil
}

//AddPriorCtrl add control request to the control queue first place.
func (a *Q) AddPriorCtrl(cmd interface{}) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return ErrClosed
	}
	a.ctrlList.PushFront(cmd)
	a.cond.Broadcast()
	return nil
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
	a.lock.Lock()
	defer a.lock.Unlock()
	for a.ctrlList.Len() == 0 && a.reqList.Len() == 0 {
		if a.closed {
			return nil, ErrClosed
		}
		a.cond.Wait()
	}
	if a.closed {
		return nil, ErrClosed
	}
	//firstly pop ctrl list
	var front = a.ctrlList.Front()
	if front != nil {
		a.ctrlList.Remove(front)
		return front.Value, nil
	}
	front = a.reqList.Front()
	if front != nil {
		a.reqList.Remove(front)
		return front.Value, nil
	}
	return nil, ErrSync
}

//PopAnyway consume an item like Pop, but it can consume even the queue is closed.
func (a *Q) PopAnyway() (interface{}, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for a.ctrlList.Len() == 0 && a.reqList.Len() == 0 {
		if a.closed {
			return nil, ErrClosed
		}
		a.cond.Wait()
	}
	//firstly pop ctrl list
	var front = a.ctrlList.Front()
	if front != nil {
		a.ctrlList.Remove(front)
		return front.Value, nil
	}
	front = a.reqList.Front()
	if front != nil {
		a.reqList.Remove(front)
		return front.Value, nil
	}
	return nil, ErrSync
}

//Close : close the queue
func (a *Q) Close() {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return
	}
	close(a.stopChan)
	a.closed = true
	a.cond.Broadcast()
}

//TryClose try to close a queue in case it's empty.
//otherwise, the queue can not be closed.
func (a *Q) TryClose() bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return a.closed
	}
	if a.ctrlList.Len() == 0 && a.reqList.Len() == 0 {
		close(a.stopChan)
		a.closed = true
		a.cond.Broadcast()
	}
	return a.closed
}

//TryClear try to clear a queue.
//the function should be called by consumer.
//if the queue be set can pop even after closed.
//the consumer handle the last pop item, it can call the function.
//after the function be called, which means the queue is totally clear.
//no producer/no consumer anymore
func (a *Q) TryClear() bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.cleared {
		return a.cleared
	}
	if a.closed {
		if a.ctrlList.Len() == 0 && a.reqList.Len() == 0 {
			close(a.clearChan)
			a.cleared = true
		}
	}
	return a.cleared
}

//WaitClose wait close, must call in another go routine
func (a *Q) WaitClose(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-a.stopChan:
		return nil
	}
}

//WaitClear wait clear, must call in another go routine
//clear must be called after queue closed.
//be caution: if there is no one to call TryClear to clear the queue
//the clear would never happen
func (a *Q) WaitClear(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-a.clearChan:
		return nil
	}
}

//IsClosed is closed or not
func (a *Q) IsClosed() bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.closed
}

//IsCleared is cleared or not
func (a *Q) IsCleared() bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.cleared
}
