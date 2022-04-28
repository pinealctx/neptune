package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTTLCache_Normal(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	_ = c.Set("123", "111")
	var v, err = c.Get("123")
	assert.Equal(t, nil, err)
	assert.Equal(t, "111", v)
}

func TestTTLCache_NotFound(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	var v, err = c.Get("123")
	assert.Equal(t, ErrTTLKeyNotFound, err)
	assert.Equal(t, nil, v)
}

func TestTTLCache_GetThenRemove(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	_ = c.Set("123", "111")
	var v, err = c.Get("123", WithRemoveAfterGet())
	v, err = c.Get("123")
	assert.Equal(t, ErrTTLKeyNotFound, err)
	assert.Equal(t, nil, v)
}

func TestTTLCache_Remove(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	_ = c.Set("123", "111")
	c.Remove("123")
	var v, err = c.Get("123")
	assert.Equal(t, ErrTTLKeyNotFound, err)
	assert.Equal(t, nil, v)
}

func TestTTLCache_Timeout(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	_ = c.Set("123", "111")
	now = func() int64 {
		return time.Now().Unix() + 2000
	}
	var v, err = c.Get("123")
	assert.Equal(t, ErrTTLKeyTimeout, err)
	assert.Equal(t, "111", v)
}
