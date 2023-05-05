package etcd

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/pkg/v3/transport"
	"go.etcd.io/etcd/client/v3"
	"path"
	"strings"
	"time"
)

const (
	//DefaultDialTimeout default dial timeout
	DefaultDialTimeout = 3 * time.Second

	//IgnoreRevision use -1 as ignore revision check
	IgnoreRevision = -1
)

// etcd cli option
type _Option struct {
	//dial timeout
	dialTimeout time.Duration
	//tls info
	tlsInfo *transport.TLSInfo
}

// Option setup option
type Option func(*_Option)

// WithTLS set tls
func WithTLS(certPath, keyPath, caPath string) Option {
	return func(o *_Option) {
		o.tlsInfo = &transport.TLSInfo{
			CertFile:      certPath,
			KeyFile:       keyPath,
			TrustedCAFile: caPath,
		}
	}
}

// WithDialTimeout set dial timeout
func WithDialTimeout(timeout time.Duration) Option {
	return func(o *_Option) {
		o.dialTimeout = timeout
	}
}

// WRet Write Ret
type WRet struct {
	Revision int64
	Err      error
}

// NodeRet Read Node Ret
type NodeRet struct {
	Data     []byte
	Err      error
	Revision int64
}

// DirRet Read Dir Ret
type DirRet struct {
	KVS      []*mvccpb.KeyValue
	Err      error
	Revision int64
}

// ExceptNotFoundErr dir ret is error except "not found"
func (d *DirRet) ExceptNotFoundErr() bool {
	if d.Err == nil {
		return false
	}
	return !IsNotFoundErr(d.Err)
}

// DebugInfo debug info
func (d *DirRet) DebugInfo() string {
	var buffer strings.Builder
	var head = fmt.Sprintf("len kv:%d, error:%+v, revision:%+v",
		len(d.KVS), d.Err, d.Revision)
	buffer.WriteString(head)
	buffer.WriteString(":\n")
	for _, kv := range d.KVS {
		buffer.WriteString("key:")
		buffer.WriteString(string(kv.Key))
		buffer.WriteString(",")
		buffer.WriteString("val:")
		buffer.WriteString(string(kv.Value))
		buffer.WriteString("\n")
	}
	return buffer.String()
}

// Client etcd client
type Client struct {
	//etcd client
	eCli *clientv3.Client
	//root path
	root string
	//url (endpoint)
	url string
	//option
	option *_Option
}

// NewClient new etcd client
func NewClient(url string, root string, ops ...Option) (*Client, error) {
	if root == "" {
		return nil, ErrEmptyRoot
	}

	if !strings.HasPrefix(root, "/") {
		return nil, ErrInvalidPath
	}

	if strings.TrimPrefix(root, "/") == "" {
		return nil, ErrEmptyRoot
	}

	var cli = &Client{url: url, root: root}
	var option = &_Option{}
	option.dialTimeout = DefaultDialTimeout
	for _, op := range ops {
		op(option)
	}

	var cnf = clientv3.Config{
		Endpoints:   strings.Split(url, ","),
		DialTimeout: option.dialTimeout,
	}

	if option.tlsInfo != nil {
		var tlsConfig, err = option.tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}
		cnf.TLS = tlsConfig
	}

	cli.option = option

	var err error
	cli.eCli, err = clientv3.New(cnf)
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// Create create node with value, if node exist fail
func (c *Client) Create(ctx context.Context, nodePath string, content []byte) WRet {
	nodePath = path.Join(c.root, nodePath)

	var rsp, err = c.eCli.Txn(ctx).
		If(clientv3.Compare(clientv3.Version(nodePath), "=", 0)).
		Then(clientv3.OpPut(nodePath, string(content))).
		Commit()

	if err != nil {
		return WRet{Err: convertErr(nodePath, err)}
	}
	if !rsp.Succeeded {
		return WRet{
			Err:      genErr(nodePath, NodeExist),
			Revision: rsp.Header.Revision,
		}
	}
	return WRet{Revision: rsp.Header.Revision}
}

// Delete delete node
func (c *Client) Delete(ctx context.Context, nodePath string, revision int64) WRet {
	nodePath = path.Join(c.root, nodePath)

	if revision == IgnoreRevision {
		//ignore revision check, delete directly
		var rsp, err = c.eCli.Delete(ctx, nodePath)
		if err != nil {
			return WRet{Err: convertErr(nodePath, err)}
		}
		if rsp.Deleted != 1 {
			return WRet{
				Err:      genErr(nodePath, NodeNotFound),
				Revision: rsp.Header.Revision,
			}
		}
		return WRet{Revision: rsp.Header.Revision}
	}

	// We have to do a transaction. This means: if the
	// node revision is what we expect, delete it,
	// otherwise get the file. If the transaction doesn't
	// succeed, we also ask for the value of the
	// node. That way we'll know if it failed because it
	// didn't exist, or because the revision was wrong.

	//global revision should be add 1 because etcd compare only supports "< or == or >", no "<="
	//here we need "<=" compare, so use globalVersion plus 1 to fit the "<=" case
	var maxVersion = revision + 1
	var rsp, err = c.eCli.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(nodePath), "<", maxVersion)).
		Then(clientv3.OpDelete(nodePath)).
		Else(clientv3.OpGet(nodePath)).
		Commit()
	if err != nil {
		return WRet{Err: convertErr(nodePath, err)}
	}
	if !rsp.Succeeded {
		if len(rsp.Responses) > 0 {
			return WRet{
				Err:      genErr(nodePath, BadVersion),
				Revision: rsp.Header.Revision,
			}
		}
		return WRet{
			Err:      genErr(nodePath, NodeNotFound),
			Revision: rsp.Header.Revision,
		}
	}
	return WRet{Revision: rsp.Header.Revision}

}

// DeleteDir delete dir
func (c *Client) DeleteDir(ctx context.Context, nodePath string, revision int64) WRet {
	nodePath = path.Join(c.root, nodePath) + "/"

	if revision == IgnoreRevision {
		//ignore revision check, delete directly
		var rsp, err = c.eCli.Delete(ctx, nodePath, clientv3.WithPrefix())
		if err != nil {
			return WRet{Err: convertErr(nodePath, err)}
		}
		if rsp.Deleted == 0 {
			return WRet{
				Err:      genErr(nodePath, NodeNotFound),
				Revision: rsp.Header.Revision,
			}
		}
		return WRet{Revision: rsp.Header.Revision}
	}

	//global revision should be add 1 because etcd compare only supports "</==/>", no "<="
	//here we need "<=" compare, so use globalVersion plus 1 to fit the "<=" case
	var maxVersion = revision + 1
	var rsp, err = c.eCli.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(nodePath).WithPrefix(), "<", maxVersion)).
		Then(clientv3.OpDelete(nodePath, clientv3.WithPrefix())).
		Else(clientv3.OpGet(nodePath, clientv3.WithPrefix())).
		Commit()
	if err != nil {
		return WRet{Err: convertErr(nodePath, err)}
	}
	if !rsp.Succeeded {
		if len(rsp.Responses) > 0 {
			if len(rsp.Responses[0].GetResponseRange().Kvs) > 0 {
				return WRet{
					Err:      genErr(nodePath, BadVersion),
					Revision: rsp.Header.Revision,
				}
			}
		}
		return WRet{
			Err:      genErr(nodePath, NodeNotFound),
			Revision: rsp.Header.Revision,
		}
	}
	return WRet{Revision: rsp.Header.Revision}

}

// Put put node data
func (c *Client) Put(ctx context.Context, nodePath string, content []byte, revision int64) WRet {
	nodePath = path.Join(c.root, nodePath)

	if revision == IgnoreRevision {
		//ignore revision check, update directly
		var rsp, err = c.eCli.Put(ctx, nodePath, string(content))
		if err != nil {
			return WRet{Err: convertErr(nodePath, err)}
		}
		return WRet{Revision: rsp.Header.Revision}
	}

	//global revision should be add 1 because etcd compare only supports "< or == or >", no "<="
	//here we need "<=" compare, so use globalVersion plus 1 to fit the "<=" case
	var maxVersion = revision + 1
	var rsp, err = c.eCli.Txn(ctx).
		If(clientv3.Compare(clientv3.ModRevision(nodePath), "<", maxVersion)).
		Then(clientv3.OpPut(nodePath, string(content))).
		Commit()
	if err != nil {
		return WRet{Err: convertErr(nodePath, err)}
	}
	if !rsp.Succeeded {
		return WRet{
			Err:      genErr(nodePath, BadVersion),
			Revision: rsp.Header.Revision,
		}
	}
	return WRet{Revision: rsp.Header.Revision}
}

// Get get node data
func (c *Client) Get(ctx context.Context, nodePath string, opts ...clientv3.OpOption) NodeRet {
	nodePath = path.Join(c.root, nodePath)
	var rsp, err = c.eCli.Get(ctx, nodePath, opts...)
	if err != nil {
		return NodeRet{Err: convertErr(nodePath, err)}
	}
	if len(rsp.Kvs) != 1 {
		return NodeRet{
			Err:      convertErr(nodePath, err),
			Revision: rsp.Header.Revision,
		}
	}
	return NodeRet{
		Data:     rsp.Kvs[0].Value,
		Revision: rsp.Header.Revision,
	}
}

// GetDir common get dir
// get dir - if keyOnly is false, get all dir children data.
func (c *Client) GetDir(ctx context.Context, nodePath string, keyOnly bool) DirRet {
	nodePath = path.Join(c.root, nodePath) + "/"

	var inOpts = make([]clientv3.OpOption, 0, 3)
	inOpts = append(inOpts, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if keyOnly {
		inOpts = append(inOpts, clientv3.WithKeysOnly())
	}
	var rsp, err = c.eCli.Get(ctx,
		nodePath,
		inOpts...)
	if err != nil {
		return DirRet{Err: convertErr(nodePath, err)}
	}
	var count = len(rsp.Kvs)
	if count == 0 {
		return DirRet{
			Err:      genErr(nodePath, NodeNotFound),
			Revision: rsp.Header.Revision,
		}
	}

	var kvList = make([]*mvccpb.KeyValue, 0, count)
	for _, kv := range rsp.Kvs {
		var key = string(kv.Key)
		if !strings.HasPrefix(key, nodePath) {
			return DirRet{
				Err:      genErr(nodePath, BadRsp),
				Revision: rsp.Header.Revision,
			}
		}
		kvList = append(kvList, kv)
	}
	return DirRet{
		KVS:      kvList,
		Revision: rsp.Header.Revision,
	}
}

// Close close
func (c *Client) Close() error {
	if c.eCli != nil {
		return c.eCli.Close()
	}
	return nil
}
