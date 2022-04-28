package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTTLCache_Normal(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	c.Set("123", "111")
	var v, err = c.Get("123", true)
	assert.Equal(t, nil, err)
	assert.Equal(t, "111", v)
}

func TestTTLCache_NotFound(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	var v, err = c.Get("123", true)
	assert.Equal(t, ErrTTLKeyNotFound, err)
	assert.Equal(t, nil, v)
}

func TestTTLCache_GetThenRemove(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	c.Set("123", "111")
	var v, err = c.Get("123", true)
	v, err = c.Get("123", true)
	assert.Equal(t, ErrTTLKeyNotFound, err)
	assert.Equal(t, nil, v)
}

func TestTTLCache_Remove(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	c.Set("123", "111")
	var v, err = c.Remove("123")
	assert.Equal(t, nil, err)
	assert.Equal(t, "111", v)
}

func TestTTLCache_Timeout(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	c.Set("123", "111")
	now = func() int64 {
		return time.Now().Unix() + 2000
	}
	var v, err = c.Get("123", true)
	assert.Equal(t, ErrTTLKeyTimeout, err)
	assert.Equal(t, "111", v)
}

func TestTTLCache_SetWithUUID(t *testing.T) {
	var c = NewTTLCache(10, 1000)
	var k, err = c.SetWithUUID("111")
	assert.Equal(t, nil, err)
	assert.NotEmpty(t, k)
	var v interface{}
	v, err = c.Get(k, true)
	assert.Equal(t, nil, err)
	assert.Equal(t, "111", v)
}
