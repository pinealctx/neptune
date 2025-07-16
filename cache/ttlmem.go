package cache

import (
	"container/list"
	"context"
	"math"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TTLCache interface {
	// Set key value into ttl cache.
	Set(ctx context.Context, key string, value []byte, fns ...SetOptFn) error
	// Get value by key from ttl cache.
	Get(ctx context.Context, key string, fns ...GetOptFn) ([]byte, error)
	// Remove value by key.
	Remove(ctx context.Context, key string) error
	// Clear cache.
	Clear(ctx context.Context)
}

var (
	// ErrTTLKeyExists When key exists will return this error.
	ErrTTLKeyExists = status.Error(codes.AlreadyExists, "ttl.key.exists")
	// ErrTTLKeyNotFound When key not found will return this error.
	ErrTTLKeyNotFound = status.Error(codes.NotFound, "ttl.key.not.found")

	// now Gen now unix timestamp
	now = func() int64 { return time.Now().Unix() }
)

type ttlNode struct {
	key      string
	value    []byte
	deadline int64
}

type setOption struct {
	ttl          int64
	mustNotExist bool
	keepTTL      bool
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

func WithKeepTTL() SetOptFn {
	return func(option *setOption) {
		option.keepTTL = true
	}
}

type getOption struct {
	ttl            int64
	removeAfterGet bool
	updateTTL      bool
}

type GetOptFn func(*getOption)

func WithRemoveAfterGet() GetOptFn {
	return func(option *getOption) {
		option.removeAfterGet = true
	}
}

func WithUpdateTTL(ttl int64) GetOptFn {
	return func(option *getOption) {
		option.updateTTL = true
		if ttl != 0 {
			option.ttl = ttl
		}
	}
}

type ttlMemCache struct {
	size    int
	ttl     int64
	eleList *list.List
	eleHash map[string]*list.Element
	sync.RWMutex
}

// NewTTLMemCache New ttl cache
func NewTTLMemCache(size int, ttl int64) TTLCache {
	return &ttlMemCache{
		size:    size,
		ttl:     ttl,
		eleList: list.New(),
		eleHash: make(map[string]*list.Element),
	}
}

// Set key value to list.
func (t *ttlMemCache) Set(_ context.Context, key string, value []byte, fns ...SetOptFn) error {
	t.Lock()
	defer t.Unlock()
	return t.set(key, value, fns...)
}

// Get value by key from list.
func (t *ttlMemCache) Get(_ context.Context, key string, fns ...GetOptFn) ([]byte, error) {
	t.Lock()
	defer t.Unlock()
	return t.get(key, fns...)
}

// Remove key value from list.
func (t *ttlMemCache) Remove(_ context.Context, key string) error {
	t.Lock()
	defer t.Unlock()
	var ele, ok = t.eleHash[key]
	if ok {
		t.remove(ele, ele.Value.(*ttlNode))
	}
	return nil
}

// Clear cache.
func (t *ttlMemCache) Clear(_ context.Context) {
	t.Lock()
	defer t.Unlock()
	t.eleHash = make(map[string]*list.Element)
	t.eleList.Init()
}

// remove Remove element.
func (t *ttlMemCache) remove(ele *list.Element, node *ttlNode) {
	if ele != nil {
		t.eleList.Remove(ele)
		delete(t.eleHash, node.key)
	}
}

// removeTail Remove tail element.
func (t *ttlMemCache) removeTail() *ttlNode {
	var ele = t.eleList.Back()
	if ele == nil {
		return nil
	}
	var node = ele.Value.(*ttlNode)
	t.remove(ele, node)
	return node
}

// set key value to list.
func (t *ttlMemCache) set(key string, value []byte, fns ...SetOptFn) error {
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
		if !o.keepTTL {
			node.deadline = deadline(o.ttl)
		}
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
func (t *ttlMemCache) get(key string, fns ...GetOptFn) ([]byte, error) {
	var o = &getOption{ttl: t.ttl}
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
		return node.value, ErrTTLKeyNotFound
	}
	if o.removeAfterGet {
		t.remove(ele, node)
		return node.value, nil
	}
	if o.updateTTL {
		node.deadline = deadline(o.ttl)
	}
	t.eleList.MoveToFront(ele)
	return node.value, nil
}

func deadline(ttl int64) int64 {
	if ttl <= 0 {
		return math.MaxInt64
	}
	return now() + ttl
}
