package redisx

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	pingTimeout = 5 * time.Second
)

func NewClient(opt *redis.Options) (redis.Cmdable, error) {
	var rds = redis.NewClient(opt)
	var err = pingTest(rds)
	if err != nil {
		return nil, err
	}
	return rds, nil
}

func NewClusterClient(opt *redis.ClusterOptions) (redis.Cmdable, error) {
	var rds = redis.NewClusterClient(opt)
	var err = pingTest(rds)
	if err != nil {
		return nil, err
	}
	return rds, nil
}

func NewClientV2(opt *redis.Options) (*redis.Client, error) {
	var rds = redis.NewClient(opt)
	var err = pingTest(rds)
	if err != nil {
		return nil, err
	}
	return rds, nil
}

func NewClusterClientV2(opt *redis.ClusterOptions) (*redis.ClusterClient, error) {
	var rds = redis.NewClusterClient(opt)
	var err = pingTest(rds)
	if err != nil {
		return nil, err
	}
	return rds, nil
}

func pingTest(cmd redis.Cmdable) error {
	var ctx, cancel = context.WithTimeout(context.Background(), pingTimeout)
	defer cancel()
	var err = cmd.Ping(ctx).Err()
	if err != nil {
		return fmt.Errorf("redis.ping.error: %+v", err)
	}
	return nil
}
