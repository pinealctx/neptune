package mux

import (
	"container/list"
	"context"
	"errors"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	//ErrClosed close
	ErrClosed = status.Error(codes.Unavailable, "q.closed")

	//ErrQFull req q full
	ErrQFull = status.Error(codes.ResourceExhausted, "q.full")

	//ErrSync never gonna happen
	ErrSync = errors.New("never.gonna.happen.crazy")
)

// Q actor queue structure define
type Q struct {
	//request queue list
	reqList *list.List

	//stop channel
	stopChan chan struct{}

	//request pipe size max number
	reqMaxNum int

	//is closed
	closed bool

	//queue lock
	lock sync.Mutex
	//queue condition
	cond sync.Cond
}

// NewQ new queue
func NewQ(reqMaxNum int) *Q {
	var actorQ = &Q{
		reqList:  list.New(),
		stopChan: make(chan struct{}),
	}

	if reqMaxNum > 0 {
		actorQ.reqMaxNum = reqMaxNum
	}
	actorQ.cond.L = &actorQ.lock
	return actorQ
}

// AddReqAnyway add normal request to the normal queue end place anyway
// if queue full, sleep then try
func (a *Q) AddReqAnyway(req any, ts time.Duration) error {
	var err error
	for {
		err = a.AddReq(req)
		if err == ErrQFull {
			time.Sleep(ts)
		} else {
			return err
		}
	}
}

// AddReq add normal request to the normal queue end place.
func (a *Q) AddReq(req any) error {
	a.lock.Lock()
	defer a.lock.Unlock()
	if a.closed {
		return ErrClosed
	}
	if a.reqMaxNum > 0 {
		if a.reqList.Len() >= a.reqMaxNum {
			return ErrQFull
		}
	}
	a.reqList.PushBack(req)
	a.cond.Broadcast()
	return nil
}

// AddPriorReq add normal request to the normal queue first place.
func (a *Q) AddPriorReq(req any) error {
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
func (a *Q) Pop() (any, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for a.reqList.Len() == 0 {
		if a.closed {
			return nil, ErrClosed
		}
		a.cond.Wait()
	}
	if a.closed {
		return nil, ErrClosed
	}
	var front = a.reqList.Front()
	if front != nil {
		a.reqList.Remove(front)
		return front.Value, nil
	}
	return nil, ErrSync
}

// PopAnyway consume an item like Pop, but it can consume even the queue is closed.
func (a *Q) PopAnyway() (any, error) {
	a.lock.Lock()
	defer a.lock.Unlock()
	for a.reqList.Len() == 0 {
		if a.closed {
			return nil, ErrClosed
		}
		a.cond.Wait()
	}
	var front = a.reqList.Front()
	if front != nil {
		a.reqList.Remove(front)
		return front.Value, nil
	}
	return nil, ErrSync
}

// Close : close the queue
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

// WaitClose wait close, must call in another go routine
func (a *Q) WaitClose(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-a.stopChan:
		return nil
	}
}

// IsClosed is closed or not
func (a *Q) IsClosed() bool {
	a.lock.Lock()
	defer a.lock.Unlock()
	return a.closed
}
