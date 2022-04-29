package cache

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"time"
)

type ttlRdsCache struct {
	cmd    redis.Cmdable
	prefix string
	ttl    int64
}

func NewTTLRdsCache(cmd redis.Cmdable, prefix string, ttl int64) TTLCache {
	return &ttlRdsCache{
		cmd:    cmd,
		prefix: prefix,
		ttl:    ttl,
	}
}

func (t *ttlRdsCache) Set(ctx context.Context, key string, value []byte, fns ...SetOptFn) error {
	var o = &setOption{ttl: t.ttl}
	for _, fn := range fns {
		fn(o)
	}
	key = t.key(key)
	if o.mustNotExist {
		var ok, err = t.cmd.SetNX(ctx, key, value, time.Duration(o.ttl)).Result()
		if err != nil {
			return err
		}
		if !ok {
			return ErrTTLKeyExists
		}
		return nil
	}
	var ex = time.Duration(o.ttl)
	if o.keepTTL {
		ex = redis.KeepTTL
	}
	var _, err = t.cmd.Set(ctx, key, value, ex).Result()
	return err
}

func (t *ttlRdsCache) Get(ctx context.Context, key string, fns ...GetOptFn) ([]byte, error) {
	var o = &getOption{ttl: t.ttl}
	for _, fn := range fns {
		fn(o)
	}
	key = t.key(key)
	var getFn = t.cmd.Get
	if o.removeAfterGet {
		getFn = t.cmd.GetDel
	}
	var v, err = getFn(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, ErrTTLKeyNotFound
		}
		return nil, err
	}
	if o.updateTTL {
		err = t.cmd.Expire(ctx, key, time.Duration(o.ttl)).Err()
		if err != nil {
			return nil, err
		}
	}
	return v, nil
}

func (t *ttlRdsCache) Remove(ctx context.Context, key string) error {
	key = t.key(key)
	var _, err = t.cmd.Del(ctx, key).Result()
	return err
}

func (t *ttlRdsCache) Clear(ctx context.Context) {
	var scanCmd = t.cmd.Scan(ctx, 0, t.prefix+"*", 0)
	var err = scanCmd.Err()
	if err != nil {
		ulog.Error("ttlRdsCache.Clear.Scan.error", zap.String("prefix", t.prefix), zap.Error(err))
		return
	}
	var iter = scanCmd.Iterator()
	for iter.Next(ctx) {
		err = t.cmd.Del(ctx, iter.Val()).Err()
		if err != nil {
			ulog.Error("ttlRdsCache.Clear.Del.error", zap.String("key", iter.Val()), zap.Error(err))
			return
		}
	}
}

func (t *ttlRdsCache) key(k string) string {
	return t.prefix + k
}
