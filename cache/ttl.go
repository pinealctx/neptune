package cache

import (
	"container/list"
	"github.com/pinealctx/neptune/idgen/random"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"sync"
	"time"
)

const (
	uuidRetryLimit = 10
)

var (
	// ErrTTLKeyNotFound When key not found will return this error
	ErrTTLKeyNotFound = status.Error(codes.NotFound, "ttl.key.not.found")
	// ErrTTLKeyTimeout When key timeout will return this error
	ErrTTLKeyTimeout = status.Error(codes.Unavailable, "ttl.key.timeout")
	// ErrTTLInternal When gen uuid retry limit will return this error
	ErrTTLInternal = status.Error(codes.Internal, "uuid.retry.limit")

	// now Gen now unix timestamp
	now = func() int64 { return time.Now().Unix() }
)

type TTLNode struct {
	key      interface{}
	value    interface{}
	deadline int64
}

type TTLCache struct {
	size    int
	ttl     int64
	eleList *list.List
	eleHash map[interface{}]*list.Element
	sync.RWMutex
}

func NewTTLCache(size int, ttl int64) *TTLCache {
	return &TTLCache{
		size:    size,
		ttl:     ttl,
		eleList: list.New(),
		eleHash: make(map[interface{}]*list.Element),
	}
}

// Set key value to list
func (t *TTLCache) Set(key, value interface{}) {
	t.Lock()
	defer t.Unlock()
	t.set(key, value)
}

// SetWithUUID This is an extension method, when setting the value,
// the UUID key will be automatically generated to ensure that the key is not repeated
func (t *TTLCache) SetWithUUID(value interface{}) (string, error) {
	var uuid = random.MD5UUID()
	t.Lock()
	defer t.Unlock()
	var ok bool
	var tryCount int
	for {
		_, ok = t.eleHash[uuid]
		if !ok {
			t.set(uuid, value)
			return uuid, nil
		}
		if tryCount > uuidRetryLimit {
			return "", ErrTTLInternal
		}
		tryCount++
		uuid = random.MD5UUID()
	}
}

// Get value by key from list
func (t *TTLCache) Get(key interface{}, remove bool) (interface{}, error) {
	t.Lock()
	defer t.Unlock()
	return t.get(key, remove)
}

// Remove key value from list
func (t *TTLCache) Remove(key interface{}) (interface{}, error) {
	t.Lock()
	defer t.Unlock()
	var node, ok = t.removeKey(key)
	if !ok {
		return nil, ErrTTLKeyNotFound
	}
	if now() > node.deadline {
		return node.value, ErrTTLKeyTimeout
	}
	return node.value, nil
}

// Clear cache
func (t *TTLCache) Clear() {
	t.Lock()
	defer t.Unlock()
	t.eleHash = make(map[interface{}]*list.Element)
	t.eleList.Init()
}

// Len get hash and list length
func (t *TTLCache) Len() (hashLen, listLen int) {
	t.RLock()
	defer t.RUnlock()
	hashLen, listLen = len(t.eleHash), t.eleList.Len()
	return
}

// remove Remove element
func (t *TTLCache) remove(ele *list.Element, node *TTLNode) {
	if ele != nil {
		t.eleList.Remove(ele)
		delete(t.eleHash, node.key)
	}
}

// removeTail Remove tail element
func (t *TTLCache) removeTail() *TTLNode {
	var ele = t.eleList.Back()
	if ele == nil {
		return nil
	}
	var node = ele.Value.(*TTLNode)
	t.remove(ele, node)
	return node
}

// removeKey Remove element by key
func (t *TTLCache) removeKey(key interface{}) (*TTLNode, bool) {
	var ele, ok = t.eleHash[key]
	if !ok {
		return nil, false
	}
	var node = ele.Value.(*TTLNode)
	t.remove(ele, node)
	return node, true
}

// set key value to list
func (t *TTLCache) set(key, value interface{}) {
	var ele, ok = t.eleHash[key]
	if ok {
		t.eleList.MoveToFront(ele)
		var node = ele.Value.(*TTLNode)
		node.value = value
		node.deadline = t.deadline()
		return
	}
	var node = &TTLNode{key: key, value: value, deadline: t.deadline()}
	ele = t.eleList.PushFront(node)
	if t.eleList.Len() > t.size {
		t.removeTail()
	}
	t.eleHash[key] = ele
}

// get Fetch value by key.
func (t *TTLCache) get(key interface{}, remove bool) (interface{}, error) {
	var ele, ok = t.eleHash[key]
	if !ok {
		return nil, ErrTTLKeyNotFound
	}
	var node = ele.Value.(*TTLNode)
	if now() > node.deadline {
		t.remove(ele, node)
		return node.value, ErrTTLKeyTimeout
	}
	if remove {
		t.remove(ele, node)
		return node.value, nil
	}
	t.eleList.MoveToFront(ele)
	return node.value, nil
}

// deadline Calc deadline time
func (t *TTLCache) deadline() int64 {
	return now() + t.ttl
}
