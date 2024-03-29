package semap

import (
	"context"
	"sync"
)

// SemMapper define interface to Acquire/Release semaphore map container
type SemMapper interface {
	//AcquireRead acquire for read
	AcquireRead(ctx context.Context, key interface{}) (*Weighted, error)
	//ReleaseRead release read lock
	ReleaseRead(key interface{}, w *Weighted)
	//AcquireWrite acquire for write
	AcquireWrite(ctx context.Context, key interface{}) (*Weighted, error)
	//ReleaseWrite release write lock
	ReleaseWrite(key interface{}, w *Weighted)
}

// SemMap semaphore map
type SemMap struct {
	m       map[interface{}]*Weighted
	mux     *sync.Mutex
	rwRatio int
}

// NewSemMap new semaphore map
func NewSemMap(opts ...Option) SemMapper {
	var o = RangeOption(opts...)
	return newSemMap(o.rwRatio)
}

// newSemMap new semaphore map
// rwRatio : read/write ratio, for example, if it's 10, means that 10 read go routine can enter at same time.
// if one write go routine enters, no read go routine can enter.
func newSemMap(rwRatio int) *SemMap {
	var m = &SemMap{}
	m.mux = &sync.Mutex{}
	m.m = make(map[interface{}]*Weighted)
	m.rwRatio = rwRatio
	return m
}

// AcquireRead acquire for read
func (s *SemMap) AcquireRead(ctx context.Context, key interface{}) (*Weighted, error) {
	return s.acquire(ctx, key, 1)
}

// ReleaseRead release read lock
func (s *SemMap) ReleaseRead(key interface{}, w *Weighted) {
	s.release(key, w, 1)
}

// AcquireWrite acquire for write
func (s *SemMap) AcquireWrite(ctx context.Context, key interface{}) (*Weighted, error) {
	return s.acquire(ctx, key, s.rwRatio)
}

// ReleaseWrite release write lock
func (s *SemMap) ReleaseWrite(key interface{}, w *Weighted) {
	s.release(key, w, s.rwRatio)
}

// acquire : acquire lock
func (s *SemMap) acquire(ctx context.Context, key interface{}, n int) (*Weighted, error) {
	var err error
	s.mux.Lock()
	var w, ok = s.m[key]
	if ok {
		err = w.acquire(ctx, s.mux, n)
		if err != nil {
			return nil, err
		}
		return w, nil
	}
	w = newWeighted(s.rwRatio)
	s.m[key] = w
	err = w.acquire(ctx, s.mux, n)
	if err != nil {
		return nil, err
	}
	return w, nil
}

// release : release lock
func (s *SemMap) release(key interface{}, w *Weighted, n int) {
	s.mux.Lock()
	defer s.mux.Unlock()
	var empty = w.release(n)
	if empty {
		delete(s.m, key)
		return
	}
}
