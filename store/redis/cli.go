package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func New(addr string) (redis.Cmdable, error) {
	var client = redis.NewClient(&redis.Options{Addr: addr})
	var _, err = client.Ping(context.Background()).Result()
	return client, err
}

func NewCluster(addr ...string) (redis.Cmdable, error) {
	var client = redis.NewClusterClient(&redis.ClusterOptions{Addrs: addr})
	var _, err = client.Ping(context.Background()).Result()
	return client, err
}
