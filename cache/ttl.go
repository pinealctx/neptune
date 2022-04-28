package cache

import (
	"container/list"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"math"
	"sync"
	"time"
)

var (
	// ErrTTLKeyExists When key exists will return this error.
	ErrTTLKeyExists = status.Error(codes.AlreadyExists, "ttl.key.exists")
	// ErrTTLKeyNotFound When key not found will return this error.
	ErrTTLKeyNotFound = status.Error(codes.NotFound, "ttl.key.not.found")
	// ErrTTLKeyTimeout When key timeout will return this error.
	ErrTTLKeyTimeout = status.Error(codes.Unavailable, "ttl.key.timeout")

	// now Gen now unix timestamp
	now = func() int64 { return time.Now().Unix() }
)

type ttlNode struct {
	key      interface{}
	value    interface{}
	deadline int64
}

type setOption struct {
	ttl          int64
	mustNotExist bool
}

type SetOptFn func(*setOption)

func WithTTL(ttl int64) SetOptFn {
	return func(option *setOption) {
		option.ttl = ttl
	}
}

func WithMustNotExist() SetOptFn {
	return func(option *setOption) {
		option.mustNotExist = true
	}
}

type getOption struct {
	removeAfterGet bool
	update         bool
}

type GetOptFn func(*getOption)

func WithRemoveAfterGet() GetOptFn {
	return func(option *getOption) {
		option.removeAfterGet = true
	}
}

func WithUpdate() GetOptFn {
	return func(option *getOption) {
		option.update = true
	}
}

type TTLCache interface {
	// Set key value into ttl cache.
	Set(key, value interface{}, fns ...SetOptFn) error
	// Get value by key from ttl cache.
	Get(key interface{}, fns ...GetOptFn) (interface{}, error)
	// Remove value by key.
	Remove(key interface{})
	// Clear cache.
	Clear()
	// Len get ttl cache length.
	Len() (hashLen, listLen int)
}

type ttlCache struct {
	size    int
	ttl     int64
	eleList *list.List
	eleHash map[interface{}]*list.Element
	sync.RWMutex
}

// NewTTLCache New ttl cache
func NewTTLCache(size int, ttl int64) TTLCache {
	return &ttlCache{
		size:    size,
		ttl:     ttl,
		eleList: list.New(),
		eleHash: make(map[interface{}]*list.Element),
	}
}

// Set key value to list.
func (t *ttlCache) Set(key, value interface{}, fns ...SetOptFn) error {
	t.Lock()
	defer t.Unlock()
	return t.set(key, value, fns...)
}

// Get value by key from list.
func (t *ttlCache) Get(key interface{}, fns ...GetOptFn) (interface{}, error) {
	t.Lock()
	defer t.Unlock()
	return t.get(key, fns...)
}

// Remove key value from list.
func (t *ttlCache) Remove(key interface{}) {
	t.Lock()
	defer t.Unlock()
	var ele, ok = t.eleHash[key]
	if ok {
		t.remove(ele, ele.Value.(*ttlNode))
	}
}

// Clear cache.
func (t *ttlCache) Clear() {
	t.Lock()
	defer t.Unlock()
	t.eleHash = make(map[interface{}]*list.Element)
	t.eleList.Init()
}

// Len get hash and list length.
func (t *ttlCache) Len() (hashLen, listLen int) {
	t.RLock()
	defer t.RUnlock()
	hashLen, listLen = len(t.eleHash), t.eleList.Len()
	return
}

// remove Remove element.
func (t *ttlCache) remove(ele *list.Element, node *ttlNode) {
	if ele != nil {
		t.eleList.Remove(ele)
		delete(t.eleHash, node.key)
	}
}

// removeTail Remove tail element.
func (t *ttlCache) removeTail() *ttlNode {
	var ele = t.eleList.Back()
	if ele == nil {
		return nil
	}
	var node = ele.Value.(*ttlNode)
	t.remove(ele, node)
	return node
}

// set key value to list.
func (t *ttlCache) set(key, value interface{}, fns ...SetOptFn) error {
	var o = &setOption{ttl: t.ttl}
	for _, fn := range fns {
		fn(o)
	}
	var ele, ok = t.eleHash[key]
	if ok {
		if o.mustNotExist {
			return ErrTTLKeyExists
		}
		t.eleList.MoveToFront(ele)
		var node = ele.Value.(*ttlNode)
		node.value = value
		node.deadline = deadline(o.ttl)
		return nil
	}
	var node = &ttlNode{key: key, value: value, deadline: deadline(o.ttl)}
	ele = t.eleList.PushFront(node)
	if t.eleList.Len() > t.size {
		t.removeTail()
	}
	t.eleHash[key] = ele
	return nil
}

// get Fetch value by key.
func (t *ttlCache) get(key interface{}, fns ...GetOptFn) (interface{}, error) {
	var o = &getOption{}
	for _, fn := range fns {
		fn(o)
	}

	var ele, ok = t.eleHash[key]
	if !ok {
		return nil, ErrTTLKeyNotFound
	}
	var node = ele.Value.(*ttlNode)
	if now() > node.deadline {
		t.remove(ele, node)
		return node.value, ErrTTLKeyTimeout
	}
	if o.removeAfterGet {
		t.remove(ele, node)
		return node.value, nil
	}
	if o.update {
		t.eleList.MoveToFront(ele)
	}
	return node.value, nil
}

func deadline(ttl int64) int64 {
	if ttl < 0 {
		return math.MaxInt64
	}
	return now() + ttl
}
