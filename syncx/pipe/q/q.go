package q

import (
	"container/list"
	"sync"
)

// Q represents a thread-safe queue with dynamic capacity using linked list
// It provides the same interface as SliceQueue but uses container/list for storage
type Q[T any] struct {
	// items holds the queue data using linked list
	items *list.List
	// capacity is the maximum number of items the queue can hold (0 means unlimited)
	capacity int
	// closed indicates if the queue is closed
	closed bool

	// lock protects all queue operations
	lock sync.Mutex
	// condSub is used to signal waiting Pop() operations
	condSub sync.Cond
	// condPub is used to signal waiting PushBlocking() operations
	condPub sync.Cond
}

// NewQ creates a new list-based queue with specified capacity
// If capacity is 0, the queue has unlimited capacity
func NewQ[T any](capacity int) *Q[T] {
	if capacity < 0 {
		panic("queue capacity must be non-negative")
	}

	q := &Q[T]{
		items:    list.New(),
		capacity: capacity,
	}
	q.condSub.L = &q.lock
	q.condPub.L = &q.lock
	return q
}

// Push adds an item to the end of the queue
// Returns ErrQueueFull if queue is at capacity (when capacity > 0)
func (q *Q[T]) Push(item T) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return ErrClosed
	}
	if q.capacity > 0 && q.items.Len() >= q.capacity {
		return ErrQueueFull
	}

	q.items.PushBack(item)
	q.condSub.Signal() // Signal waiting Pop() operations
	return nil
}

// PushBlocking adds an item to the end of the queue
// Blocks if queue is at capacity until space is available or queue is closed
func (q *Q[T]) PushBlocking(item T) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Wait until there's space or queue is closed
	for q.capacity > 0 && q.items.Len() >= q.capacity && !q.closed {
		q.condPub.Wait() // wait for Pop to consume items
	}

	if q.closed {
		return ErrClosed
	}

	q.items.PushBack(item)
	q.condSub.Signal() // Signal waiting Pop() operations
	return nil
}

// Pop removes and returns an item from the front of the queue
// Blocks if queue is empty until an item is available or queue is closed
// Important: If the queue is closed, it immediately returns ErrClosed regardless of whether there are items left.
// This ensures that once a queue is closed, no further data can be consumed, which is useful for graceful shutdown scenarios.
func (q *Q[T]) Pop() (T, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Wait until queue has data or is closed
	for q.items.Len() == 0 && !q.closed {
		q.condSub.Wait() // wait for Push()/PushBlocking() to add items
	}

	var zero T
	// After wait loop exits, either we have data OR queue is closed
	// If queue is closed and no data, return error
	if q.closed {
		return zero, ErrClosed
	}

	// We have data (since !q.closed and loop exited)
	front := q.items.Front()
	q.items.Remove(front)
	q.condPub.Signal() // Signal waiting PushBlocking operations
	// nolint : forcetypeassert // I know the type is exactly here
	return front.Value.(T), nil
}

// Peek returns the item at the front of the queue without removing it
// Returns zero value if the queue is empty
func (q *Q[T]) Peek() T {
	q.lock.Lock()
	defer q.lock.Unlock()

	var zero T
	if q.items.Len() == 0 {
		return zero
	}

	front := q.items.Front()
	// nolint : forcetypeassert // I know the type is exactly here
	return front.Value.(T)
}

// Close closes the queue, all subsequent operations will return ErrClosed
func (q *Q[T]) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return
	}
	q.closed = true
	q.condSub.Broadcast() // Wake up all waiting Pop() operations
	q.condPub.Broadcast() // Wake up all waiting PushBlocking() operations
}

// Len returns the current number of items in the queue
func (q *Q[T]) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.items.Len()
}

// Cap returns the maximum capacity of the queue (0 means unlimited)
func (q *Q[T]) Cap() int {
	return q.capacity
}

// IsUnlimited returns true if the queue has unlimited capacity
func (q *Q[T]) IsUnlimited() bool {
	return q.capacity == 0
}

// IsClosed returns true if the queue is closed
func (q *Q[T]) IsClosed() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.closed
}

// IsFull returns true if the queue is at capacity (always false for unlimited capacity)
func (q *Q[T]) IsFull() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	if q.capacity == 0 {
		return false // unlimited capacity
	}
	return q.items.Len() >= q.capacity
}

// IsEmpty returns true if the queue has no items
func (q *Q[T]) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.items.Len() == 0
}

// Reset clears all items from the queue (useful for reusing the queue)
func (q *Q[T]) Reset() {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return
	}

	// Clear all items
	q.items.Init()
}
