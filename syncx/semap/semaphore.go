package semap

import (
	"container/list"
	"context"
	"sync"
)

type waiter struct {
	n     int
	ready chan<- struct{} // Closed when semaphore acquired.
}

// Weighted provides a way to bound concurrent access to a resource.
// The callers can request access with a given weight.
type Weighted struct {
	size    int
	cur     int
	waiters list.List
}

// newWeighted creates a new weighted semaphore with the given
// maximum combined weight for concurrent access.
func newWeighted(n int) *Weighted {
	w := &Weighted{size: n}
	return w
}

// acquire : acquires the semaphore with a weight of n, blocking until resources
// are available or ctx is done. On success, returns nil. On failure, returns
// ctx.Err() and leaves the semaphore unchanged.
//
// If ctx is already done, Acquire may still succeed without blocking.
func (s *Weighted) acquire(ctx context.Context, mu *sync.Mutex, n int) error {
	if s.size-s.cur >= n && s.waiters.Len() == 0 {
		s.cur += n
		mu.Unlock()
		return nil
	}

	if n > s.size {
		// Don't make other Acquire calls block on one that's doomed to fail.
		mu.Unlock()
		<-ctx.Done()
		return ctx.Err()
	}

	ready := make(chan struct{})
	w := waiter{n: n, ready: ready}
	elem := s.waiters.PushBack(w)
	mu.Unlock()

	select {
	case <-ctx.Done():
		err := ctx.Err()
		mu.Lock()
		select {
		case <-ready:
			// Acquired the semaphore after we were canceled.  Rather than trying to
			// fix up the queue, just pretend we didn't notice the cancelation.
			err = nil
		default:
			isFront := s.waiters.Front() == elem
			s.waiters.Remove(elem)
			// If we're at the front and there're extra tokens left, notify other waiters.
			if isFront && s.size > s.cur {
				s.notifyWaiters()
			}
		}
		mu.Unlock()
		return err

	case <-ready:
		return nil
	}
}

// release : releases the semaphore with a weight of n.
// if waiter list is empty, return true.
func (s *Weighted) release(n int) bool {
	s.cur -= n
	if s.cur < 0 {
		panic("semaphore: released more than held")
	}
	return s.notifyWaiters()
}

// notify other waiters, if waiter list is empty, return true
func (s *Weighted) notifyWaiters() bool {
	for {
		next := s.waiters.Front()
		if next == nil {
			return true // No more waiters blocked.
		}

		// nolint : forcetypeassert // I know the type is exactly here
		w := next.Value.(waiter)
		if s.size-s.cur < w.n {
			// Not enough tokens for the next waiter.  We could keep going (to try to
			// find a waiter with a smaller request), but under load that could cause
			// starvation for large requests; instead, we leave all remaining waiters
			// blocked.
			//
			// Consider a semaphore used as a read-write lock, with N tokens, N
			// readers, and one writer.  Each reader can Acquire(1) to obtain a read
			// lock.  The writer can Acquire(N) to obtain a write lock, excluding all
			// of the readers.  If we allow the readers to jump ahead in the queue,
			// the writer will starve — there is always one token available for every
			// reader.
			break
		}

		s.cur += w.n
		s.waiters.Remove(next)
		close(w.ready)
	}
	return false
}
