package oss

import (
	"io"
)

const (
	//Private 不公开访问权限
	Private ACLType = 0
	//PublicRead 可公开读的访问权限
	PublicRead ACLType = 1
	//DefaultACL 缺省权限(由全局配置决定，例如OSS整体的bucket权限)
	DefaultACL ACLType = 2
)

type ACLType int32

var (
	//StoreIns global instance
	StoreIns *StoreWrapper
)

type aclOption struct {
	o ACLType
}

// ACLOption ACL option function
type ACLOption func(o *aclOption)

// UsePrivateACL use private ACL
func UsePrivateACL() ACLOption {
	return func(o *aclOption) {
		o.o = Private
	}
}

// UsePublicACL use private ACL
func UsePublicACL() ACLOption {
	return func(o *aclOption) {
		o.o = PublicRead
	}
}

// IOssStore store interface
type IOssStore interface {
	//Save save k-v
	Save(key string, data []byte, acl ACLType) error

	//SaveWithReader : save with io.Reader
	SaveWithReader(key string, reader io.Reader, acl ACLType) error

	//SaveWithReadCloser : save with io.ReadCloser
	SaveWithReadCloser(key string, readCloser io.ReadCloser, acl ACLType) error

	//Delete delete k
	Delete(key string) error

	//DeleteMulti delete multi keys
	DeleteMulti(keys []string) ([]string, error)

	//Get : get v from k
	Get(key string) (data []byte, err error)
}

// StoreWrapper store container
type StoreWrapper struct {
	storeIns IOssStore
}

func (c *StoreWrapper) SetStore(i IOssStore) {
	c.storeIns = i
}

// Save save k-v
func (c *StoreWrapper) Save(key string, data []byte, opts ...ACLOption) error {
	var o = &aclOption{o: DefaultACL}
	for _, opt := range opts {
		opt(o)
	}
	return c.storeIns.Save(key, data, o.o)
}

// SaveWithReader : save with io.Reader
func (c *StoreWrapper) SaveWithReader(key string, reader io.Reader, opts ...ACLOption) error {
	var o = &aclOption{o: DefaultACL}
	for _, opt := range opts {
		opt(o)
	}
	return c.storeIns.SaveWithReader(key, reader, o.o)
}

// SaveWithReadCloser : save with io.ReadCloser
func (c *StoreWrapper) SaveWithReadCloser(key string, readCloser io.ReadCloser, opts ...ACLOption) error {
	var o = &aclOption{o: DefaultACL}
	for _, opt := range opts {
		opt(o)
	}
	return c.storeIns.SaveWithReadCloser(key, readCloser, o.o)
}

// Delete delete k
func (c *StoreWrapper) Delete(key string) error {
	return c.storeIns.Delete(key)
}

// DeleteMulti delete multi
func (c *StoreWrapper) DeleteMulti(keys []string) ([]string, error) {
	return c.storeIns.DeleteMulti(keys)
}

// Get : get v from k
func (c *StoreWrapper) Get(key string) (data []byte, err error) {
	return c.storeIns.Get(key)
}

func init() {
	StoreIns = &StoreWrapper{}
}
