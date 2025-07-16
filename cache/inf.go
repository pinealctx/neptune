package cache

// LRUFacade an interface to define a LRU cache
type LRUFacade interface {
	// Get returns a value from the cache, and marks the entry as most recently used.
	Get(key any) (v Value, ok bool)
	// Peek returns a value from the cache without changing the LRU order.
	Peek(key any) (v Value, ok bool)
	// Exist : return true if key in map
	Exist(key any) bool
	// Set sets a value in the cache.
	Set(key any, value Value)
	// Delete removes an entry from the cache, and returns if the entry existed.
	Delete(key any) bool
}

// MapFacade an interface to define a Map
type MapFacade interface {
	//Set : set key-value
	Set(key any, value any)
	//Get : get value
	Get(key any) (any, bool)
	//Exist : return true if key in map
	Exist(key any) bool
	//Delete : delete a key
	Delete(key any)
}
