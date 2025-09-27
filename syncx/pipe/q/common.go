package q

import "errors"

var (
	// ErrClosed is returned when attempting to operate on a closed queue
	ErrClosed = errors.New("pipe.q.closed")

	// ErrQueueFull is returned when attempting to push to a full queue
	ErrQueueFull = errors.New("pipe.q.full")
)

// Queue defines the common interface for all queue implementations
type Queue[T any] interface {
	// Push adds an item to the queue
	Push(item T) error

	// PushBlocking adds an item to the end of the queue
	// Blocks if queue is at capacity until space is available or queue is closed
	PushBlocking(item T) error

	// Pop removes and returns an item from the queue
	// Blocks if queue is empty until an item is available or queue is closed
	Pop() (T, error)

	// Peek returns the item at the front of the queue without removing it
	// Returns zero value if the queue is empty
	Peek() T

	// Close closes the queue
	Close()

	// Len returns the current number of items in the queue
	Len() int

	// Cap returns the maximum capacity of the queue
	Cap() int

	// IsUnlimited returns true if the queue has unlimited capacity
	IsUnlimited() bool

	// IsClosed returns true if the queue is closed
	IsClosed() bool

	// IsFull returns true if the queue is at capacity
	IsFull() bool

	// IsEmpty returns true if the queue has no items
	IsEmpty() bool

	// Reset clears all items from the queue
	Reset()
}
