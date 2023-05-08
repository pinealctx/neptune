package cache

// LRUFacade an interface to define a LRU cache
type LRUFacade interface {
	// Get returns a value from the cache, and marks the entry as most recently used.
	Get(key interface{}) (v Value, ok bool)
	// Peek returns a value from the cache without changing the LRU order.
	Peek(key interface{}) (v Value, ok bool)
	// Exist : return true if key in map
	Exist(key interface{}) bool
	// Set sets a value in the cache.
	Set(key interface{}, value Value)
	// Delete removes an entry from the cache, and returns if the entry existed.
	Delete(key interface{}) bool
}

// MapFacade an interface to define a Map
type MapFacade interface {
	//Set : set key-value
	Set(key interface{}, value interface{})
	//Get : get value
	Get(key interface{}) (interface{}, bool)
	//Exist : return true if key in map
	Exist(key interface{}) bool
	//Delete : delete a key
	Delete(key interface{})
}
