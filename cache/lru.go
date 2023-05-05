/*
Copyright 2017 Google Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package cache implements a LRU cache.
//
// The implementation borrows heavily from SmallLRUCache
// (originally by Nathan Schrenk). The object maintains a doubly-linked list of
// elements. When an element is accessed, it is promoted to the head of the
// list. When space is needed, the element at the tail of the list
// (the least recently used element) is evicted.
package cache

import (
	"container/list"
	"fmt"
	"sync"
)

// LRUFacade an interface to define a LRU cache
type LRUFacade interface {
	// Get returns a value from the cache, and marks the entry as most recently used.
	Get(key interface{}) (v Value, ok bool)
	// Peek returns a value from the cache without changing the LRU order.
	Peek(key interface{}) (v Value, ok bool)
	//Exist : return true if key in map
	Exist(key interface{}) bool
	// Set sets a value in the cache.
	Set(key interface{}, value Value)
	// Delete removes an entry from the cache, and returns if the entry existed.
	Delete(key interface{}) bool
}

// LRUCache is a typical LRU cache implementation.  If the cache
// reaches the capacity, the least recently used item is deleted from
// the cache. Note the capacity is not the number of items, but the
// total sum of the Size() of each item.
type LRUCache struct {
	mu sync.Mutex

	// list & table contain *entry objects.
	list  *list.List
	table map[interface{}]*list.Element

	size      int64
	capacity  int64
	evictions int64
}

// Value is the interface values that go into LRUCache need to satisfy
type Value interface {
	// Size returns how big this value is. If you want to just track
	// the cache by number of objects, you may return the size as 1.
	Size() int
}

// Item is what is stored in the cache
type Item struct {
	Key   interface{}
	Value Value
}

type entry struct {
	key   interface{}
	value Value
	size  int64
}

// NewSingleLRUCache create a single lru cache
func NewSingleLRUCache(capacity int64) LRUFacade {
	return NewLRUCache(capacity)
}

// NewLRUCache creates a new empty cache with the given capacity.
func NewLRUCache(capacity int64) *LRUCache {
	var c = &LRUCache{}
	c.Init(capacity)
	return c
}

// Init : init memory
func (lru *LRUCache) Init(capacity int64) {
	lru.list = list.New()
	lru.table = make(map[interface{}]*list.Element)
	lru.capacity = capacity
}

// Get returns a value from the cache, and marks the entry as most
// recently used.
func (lru *LRUCache) Get(key interface{}) (v Value, ok bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return nil, false
	}
	lru.list.MoveToFront(element)
	return element.Value.(*entry).value, true
}

// Peek returns a value from the cache without changing the LRU order.
func (lru *LRUCache) Peek(key interface{}) (v Value, ok bool) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return nil, false
	}
	return element.Value.(*entry).value, true
}

// Exist : return true if key in map
func (lru *LRUCache) Exist(key interface{}) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	var _, ok = lru.table[key]
	return ok
}

// Set sets a value in the cache.
func (lru *LRUCache) Set(key interface{}, value Value) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element := lru.table[key]; element != nil {
		lru.updateInPlace(element, value)
	} else {
		lru.addNew(key, value)
	}
}

// SetAndGetRemoved sets a value in the cache and returns the removed value list
func (lru *LRUCache) SetAndGetRemoved(key interface{}, value Value) (removedValueList []Value) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	if element := lru.table[key]; element != nil {
		return lru.updateInPlaceAndGetRemoved(element, value)
	} else {
		return lru.addNewAndGetRemoved(key, value)
	}
}

// SetIfAbsent will set the value in the cache if not present. If the
// value exists in the cache, we don't set it.
func (lru *LRUCache) SetIfAbsent(key interface{}, value Value) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	if element := lru.table[key]; element != nil {
		lru.list.MoveToFront(element)
	} else {
		lru.addNew(key, value)
	}
}

// Delete removes an entry from the cache, and returns if the entry existed.
func (lru *LRUCache) Delete(key interface{}) bool {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	element := lru.table[key]
	if element == nil {
		return false
	}

	lru.list.Remove(element)
	delete(lru.table, key)
	lru.size -= element.Value.(*entry).size
	return true
}

// Clear will clear the entire cache.
func (lru *LRUCache) Clear() {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.list.Init()
	lru.table = make(map[interface{}]*list.Element)
	lru.size = 0
}

// SetCapacity will set the capacity of the cache. If the capacity is
// smaller, and the current cache size exceed that capacity, the cache
// will be shrank.
func (lru *LRUCache) SetCapacity(capacity int64) {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	lru.capacity = capacity
	lru.checkCapacity()
}

// Stats returns a few stats on the cache.
func (lru *LRUCache) Stats() (length, size, capacity, evictions int64) {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return int64(lru.list.Len()), lru.size, lru.capacity, lru.evictions
}

// StatsJSON returns stats as a JSON object in a string.
func (lru *LRUCache) StatsJSON() string {
	if lru == nil {
		return "{}"
	}
	l, s, c, e := lru.Stats()
	return fmt.Sprintf("{\"Length\": %v, \"Size\": %v, \"Capacity\": %v, \"Evictions\": %v}", l, s, c, e)
}

// Length returns how many elements are in the cache
func (lru *LRUCache) Length() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return int64(lru.list.Len())
}

// Size returns the sum of the objects' Size() method.
func (lru *LRUCache) Size() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.size
}

// Capacity returns the cache maximum capacity.
func (lru *LRUCache) Capacity() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.capacity
}

// Evictions returns the eviction count.
func (lru *LRUCache) Evictions() int64 {
	lru.mu.Lock()
	defer lru.mu.Unlock()
	return lru.evictions
}

// Keys returns all the keys for the cache, ordered from most recently
// used to last recently used.
func (lru *LRUCache) Keys() []interface{} {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	keys := make([]interface{}, 0, lru.list.Len())
	for e := lru.list.Front(); e != nil; e = e.Next() {
		keys = append(keys, e.Value.(*entry).key)
	}
	return keys
}

// Items returns all the values for the cache, ordered from most recently
// used to last recently used.
func (lru *LRUCache) Items() []Item {
	lru.mu.Lock()
	defer lru.mu.Unlock()

	items := make([]Item, 0, lru.list.Len())
	for e := lru.list.Front(); e != nil; e = e.Next() {
		v := e.Value.(*entry)
		items = append(items, Item{Key: v.key, Value: v.value})
	}
	return items
}

func (lru *LRUCache) updateInPlace(element *list.Element, value Value) {
	valueSize := int64(value.Size())
	sizeDiff := valueSize - element.Value.(*entry).size
	element.Value.(*entry).value = value
	element.Value.(*entry).size = valueSize
	lru.size += sizeDiff
	lru.list.MoveToFront(element)
	lru.checkCapacity()
}

func (lru *LRUCache) updateInPlaceAndGetRemoved(element *list.Element, value Value) []Value {
	valueSize := int64(value.Size())
	sizeDiff := valueSize - element.Value.(*entry).size
	element.Value.(*entry).value = value
	element.Value.(*entry).size = valueSize
	lru.size += sizeDiff
	lru.list.MoveToFront(element)
	return lru.checkCapacityAndGetRemoved()
}

func (lru *LRUCache) addNew(key interface{}, value Value) {
	newEntry := &entry{key, value, int64(value.Size())}
	element := lru.list.PushFront(newEntry)
	lru.table[key] = element
	lru.size += newEntry.size
	lru.checkCapacity()
}

func (lru *LRUCache) addNewAndGetRemoved(key interface{}, value Value) []Value {
	newEntry := &entry{key, value, int64(value.Size())}
	element := lru.list.PushFront(newEntry)
	lru.table[key] = element
	lru.size += newEntry.size
	return lru.checkCapacityAndGetRemoved()
}

func (lru *LRUCache) checkCapacity() {
	// Partially duplicated from Delete
	for lru.size > lru.capacity {
		delElem := lru.list.Back()
		delValue := delElem.Value.(*entry)
		lru.list.Remove(delElem)
		delete(lru.table, delValue.key)
		lru.size -= delValue.size
		lru.evictions++
	}
}

func (lru *LRUCache) checkCapacityAndGetRemoved() (removedValueList []Value) {
	for lru.size > lru.capacity {
		delElem := lru.list.Back()
		delValue := delElem.Value.(*entry)
		lru.list.Remove(delElem)
		delete(lru.table, delValue.key)
		lru.size -= delValue.size
		lru.evictions++
		removedValueList = append(removedValueList, delValue.value)
	}

	return
}
