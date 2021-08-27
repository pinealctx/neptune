package oss

import (
	"context"
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
	//OssStoreIns global instance
	OssStoreIns *OssStore
)

type aclOption struct {
	o ACLType
}

//ACLOption ACL option function
type ACLOption func(o *aclOption)

//UsePrivateACL use private ACL
func UsePrivateACL() ACLOption {
	return func(o *aclOption) {
		o.o = Private
	}
}

//UsePublicACL use private ACL
func UsePublicACL() ACLOption {
	return func(o *aclOption) {
		o.o = PublicRead
	}
}

//IOssStore store interface
type IOssStore interface {
	//Save save k-v
	Save(ctx context.Context, key string, data []byte, acl ACLType) error

	//Delete delete k
	Delete(ctx context.Context, key string) error

	//DeleteMulti delete multi keys
	DeleteMulti(ctx context.Context, keys []string) ([]string, error)

	//Get get v from k
	Get(ctx context.Context, key string) (data []byte, err error)
}

//OssStore store container
type OssStore struct {
	storeIns IOssStore
}

func (c *OssStore) SetStore(i IOssStore) {
	c.storeIns = i
}

//Save save k-v
func (c *OssStore) Save(ctx context.Context, key string, data []byte, opts ...ACLOption) error {
	var o = &aclOption{o: DefaultACL}
	for _, opt := range opts {
		opt(o)
	}
	return c.storeIns.Save(ctx, key, data, o.o)
}

//Delete delete k
func (c *OssStore) Delete(ctx context.Context, key string) error {
	return c.storeIns.Delete(ctx, key)
}

//DeleteMulti delete multi
func (c *OssStore) DeleteMulti(ctx context.Context, keys []string) ([]string, error) {
	return c.storeIns.DeleteMulti(ctx, keys)
}

//Get get v from k
func (c *OssStore) Get(ctx context.Context, key string) (data []byte, err error) {
	return c.storeIns.Get(ctx, key)
}

func init() {
	OssStoreIns = &OssStore{}
}
