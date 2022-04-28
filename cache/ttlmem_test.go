package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTTLCache_Normal(t *testing.T) {
	var c = NewTTLMemCache(10, 1000)
	_ = c.Set(context.TODO(), "123", []byte("111"))
	var v, err = c.Get(context.TODO(), "123")
	assert.Equal(t, nil, err)
	assert.Equal(t, []byte("111"), v)
}

func TestTTLCache_NotFound(t *testing.T) {
	var c = NewTTLMemCache(10, 1000)
	var v, err = c.Get(context.TODO(), "123")
	assert.Equal(t, ErrTTLKeyNotFound, err)
	assert.Empty(t, v)
}

func TestTTLCache_GetThenRemove(t *testing.T) {
	var c = NewTTLMemCache(10, 1000)
	_ = c.Set(context.TODO(), "123", []byte("111"))
	var _, _ = c.Get(context.TODO(), "123", WithRemoveAfterGet())
	var v, err = c.Get(context.TODO(), "123")
	assert.Equal(t, ErrTTLKeyNotFound, err)
	assert.Empty(t, v)
}

func TestTTLCache_Remove(t *testing.T) {
	var c = NewTTLMemCache(10, 1000)
	_ = c.Set(context.TODO(), "123", []byte("111"))
	_ = c.Remove(context.TODO(), "123")
	var v, err = c.Get(context.TODO(), "123")
	assert.Equal(t, ErrTTLKeyNotFound, err)
	assert.Empty(t, v)
}

func TestTTLCache_Timeout(t *testing.T) {
	var c = NewTTLMemCache(10, 1000)
	_ = c.Set(context.TODO(), "123", []byte("111"))
	now = func() int64 {
		return time.Now().Unix() + 2000
	}
	var v, err = c.Get(context.TODO(), "123")
	assert.Equal(t, ErrTTLKeyNotFound, err)
	assert.Equal(t, []byte("111"), v)
}
