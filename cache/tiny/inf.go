package tiny

// LRU an interface to define a LRU cache, each v is regard same size
type LRU interface {
	// Get returns a value from the cache, and marks the entry as most recently used.
	Get(key any) (v any, ok bool)
	// Peek returns a value from the cache without changing the LRU order.
	Peek(key any) (v any, ok bool)
	// Exist : return true if key in map
	Exist(key any) bool
	// Set sets a value in the cache.
	Set(key any, value any)
	// Delete removes an entry from the cache, and returns if the entry existed.
	Delete(key any) bool
}
