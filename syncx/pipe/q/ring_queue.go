package q

import (
	"sync"
)

// RingQ represents a thread-safe ring buffer queue with fixed capacity using pre-allocated slice
// The slice is allocated at creation time with the specified capacity
type RingQ[T any] struct {
	// items holds the queue data in a circular buffer
	items []T
	// head is the index of the first item
	head int
	// tail is the index where next item will be inserted
	tail int
	// count is the current number of items in queue
	count int
	// capacity is the maximum number of items the queue can hold
	capacity int
	// closed indicates if the queue is closed
	closed bool

	// lock protects all queue operations
	lock sync.Mutex
	// cond is used to signal waiting goroutines
	cond sync.Cond
}

// NewRingQ creates a new slice-based ring queue with fixed capacity
// The slice is pre-allocated with the specified capacity
func NewRingQ[T any](capacity int) *RingQ[T] {
	if capacity <= 0 {
		panic("slice queue capacity must be positive")
	}

	q := &RingQ[T]{
		items:    make([]T, capacity), // Pre-allocate the slice
		capacity: capacity,
	}
	q.cond.L = &q.lock
	return q
}

// Push adds an item to the end of the queue
// Returns ErrQueueFull if queue is at capacity
func (q *RingQ[T]) Push(item T) error {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return ErrClosed
	}
	if q.count >= q.capacity {
		return ErrQueueFull
	}

	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.capacity
	q.count++
	q.cond.Signal()
	return nil
}

// Pop removes and returns an item from the front of the queue
// Blocks if queue is empty until an item is available or queue is closed
func (q *RingQ[T]) Pop() (T, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	// Wait until queue has data or is closed
	for q.count == 0 && !q.closed {
		q.cond.Wait()
	}

	var zero T
	// After wait loop exits, either we have data OR queue is closed
	// If queue is closed, return error
	if q.closed {
		return zero, ErrClosed
	}

	// We have data (since !q.closed and loop exited)
	item := q.items[q.head]
	q.items[q.head] = zero // Clear to avoid memory leaks
	q.head = (q.head + 1) % q.capacity
	q.count--
	return item, nil
}

// Close closes the queue, all subsequent operations will return ErrClosed
func (q *RingQ[T]) Close() {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return
	}
	q.closed = true
	q.cond.Broadcast()
}

// Len returns the current number of items in the queue
func (q *RingQ[T]) Len() int {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count
}

// Cap returns the maximum capacity of the queue
func (q *RingQ[T]) Cap() int {
	return q.capacity
}

// IsClosed returns true if the queue is closed
func (q *RingQ[T]) IsClosed() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.closed
}

// IsFull returns true if the queue is at capacity
func (q *RingQ[T]) IsFull() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count >= q.capacity
}

// IsEmpty returns true if the queue has no items
func (q *RingQ[T]) IsEmpty() bool {
	q.lock.Lock()
	defer q.lock.Unlock()
	return q.count == 0
}

// Reset clears all items from the queue (useful for reusing the queue)
func (q *RingQ[T]) Reset() {
	q.lock.Lock()
	defer q.lock.Unlock()

	if q.closed {
		return
	}

	// Clear all items to avoid memory leaks
	var zero T
	for i := 0; i < q.capacity; i++ {
		q.items[i] = zero
	}

	q.head = 0
	q.tail = 0
	q.count = 0
}
