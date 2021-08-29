package alioss

import (
	"bytes"
	"context"
	aliOss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/pinealctx/neptune/store/oss"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"net/http"
)

//OneBucketOss all oss store into one bucket
//所有数据保存在同一个bucket中的方式
type OneBucketOss struct {
	bucket *aliOss.Bucket
}

//NewOneBucketOss new
func NewOneBucketOss(endP, buckName, accKey, secret string) (*OneBucketOss, error) {
	var cli, err = aliOss.New(endP, accKey, secret)
	if err != nil {
		return nil, err
	}
	var bucket *aliOss.Bucket
	bucket, err = cli.Bucket(buckName)
	if err != nil {
		return nil, err
	}
	return &OneBucketOss{
		bucket: bucket,
	}, nil
}

//Save save
func (o *OneBucketOss) Save(_ context.Context, key string, data []byte, acl oss.ACLType) error {
	return o.save(key, data, acl)
}

//Delete delete
func (o *OneBucketOss) Delete(_ context.Context, key string) error {
	return o.bucket.DeleteObject(key)
}

//DeleteMulti delete multi keys
func (o *OneBucketOss) DeleteMulti(_ context.Context, keys []string) ([]string, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	var r, err = o.bucket.DeleteObjects(keys)
	if err != nil {
		return nil, err
	}
	return r.DeletedObjects, err
}

//Get return data/error
func (o *OneBucketOss) Get(_ context.Context, key string) ([]byte, error) {
	var buf, err = o.get(key)
	if err != nil {
		var ossErr, ok = err.(aliOss.ServiceError)
		if ok {
			if ossErr.StatusCode == http.StatusNotFound && ossErr.Code == "NoSuchKey" {
				return nil, status.Errorf(codes.NotFound, "oss.key.not.exist:%+v", key)
			}
		}
		return nil, err
	}
	return buf, nil
}

//save raw
func (o *OneBucketOss) save(key string, data []byte, acl oss.ACLType) error {
	var reader = bytes.NewReader(data)
	if acl == oss.PublicRead {
		return o.bucket.PutObject(key, reader, aliOss.ObjectACL(aliOss.ACLPublicRead))
	}
	if acl == oss.Private {
		return o.bucket.PutObject(key, reader, aliOss.ObjectACL(aliOss.ACLPrivate))
	}
	return o.bucket.PutObject(key, reader)
}

//get raw
func (o *OneBucketOss) get(key string) ([]byte, error) {
	var reader, err = o.bucket.GetObject(key)
	if err != nil {
		return nil, err
	}
	defer func() { _ = reader.Close() }()
	return ioutil.ReadAll(reader)
}
