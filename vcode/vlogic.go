package vcode

import (
	"fmt"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/pinealctx/neptune/cache"
	"github.com/pinealctx/neptune/idgen/random"
)

const (
	numChars = "0123456789"
)

var (
	ErrSendTooFreq          = status.Error(codes.ResourceExhausted, "send.code.freq.limit")
	ErrVerifyCodeRetryLimit = status.Error(codes.ResourceExhausted, "verify.code.retry.limit")
	ErrSendCountLimit       = status.Error(codes.ResourceExhausted, "send.code.count.limit")

	ErrVerifyCodeNotExist     = status.Error(codes.NotFound, "verify.code.not.exist")
	ErrVerifyCodeTimeout      = status.Error(codes.DeadlineExceeded, "verify.code.timeout")
	ErrVerifyCodeNotMatch     = status.Error(codes.Unauthenticated, "verify.code.not.match")
	ErrVerifyCodeHashNotMatch = status.Error(codes.FailedPrecondition, "verify.hash.code.not.match")
)

type smsModule interface {
	SendCode(areaCode, phone, code string) error
}

type CacheModule interface {
	Get(key any) (v any, ok bool)
	Peek(key any) (v any, ok bool)
	Set(key any, value any)
}

type simpleCache struct {
	lru *cache.LRUCache
}

func NewSimpleCache(c int64) CacheModule {
	return &simpleCache{lru: cache.NewLRUCache(c)}
}

func (s simpleCache) Get(key any) (v any, ok bool) {
	return s.lru.Get(key)
}

func (s simpleCache) Peek(key any) (v any, ok bool) {
	return s.lru.Peek(key)
}

func (s simpleCache) Set(key any, value any) {
	var v, ok = value.(cache.Value)
	if ok {
		s.lru.Set(key, v)
	}
}

type VCLogic interface {
	SendSMSCode(areaCode, phone string) (string, error)
	VerifySMSCode(areaCode, phone, code, hash string) error
}

type sender struct {
	*Config
	sms    smsModule
	cacheM CacheModule
}

func NewSimpleLogic(config *Config, sms smsModule, cacheM CacheModule) VCLogic {
	return &sender{
		Config: config,
		cacheM: NewSimpleCache(config.CacheSize),
		sms:    sms,
	}
}

func (s *sender) SendSMSCode(areaCode, phone string) (string, error) {
	var key = fmt.Sprintf("%s-%s", areaCode, phone)
	var now = time.Now()
	var c = s.fetchCache(key, true)
	if c == nil {
		c = newSenderCache(now)
	}
	var err = s.checkSend(c, now)
	if err != nil {
		return "", err
	}
	var code = s.genCode(phone)
	c.updateSend(code, now)
	s.cacheM.Set(key, c)
	if !s.Mock {
		err = s.sms.SendCode(areaCode, phone, code)
	}
	return c.hash, err
}

func (s *sender) VerifySMSCode(areaCode, phone, code, hash string) error {
	var key = fmt.Sprintf("%s%s", areaCode, phone)
	var c = s.fetchCache(key, false)
	if c == nil {
		return ErrVerifyCodeNotExist
	}
	return s.checkVerify(c, code, hash)
}

func (s *sender) fetchCache(key string, peek bool) *vCache {
	var fn = s.cacheM.Get
	if peek {
		fn = s.cacheM.Peek
	}
	var item, ok = fn(key)
	if !ok {
		return nil
	}
	var ret *vCache
	ret, ok = item.(*vCache)
	if !ok {
		return nil
	}
	return ret
}

func (s *sender) genCode(phone string) string {
	if !s.Mock {
		return random.SecGenNonceStr(numChars, s.CodeLen)
	}
	var l = len(phone)
	if l >= s.CodeLen {
		return phone[l-s.CodeLen:]
	}
	var dist = phone
	for i := 0; i < s.CodeLen-l; i++ {
		dist = fmt.Sprintf("0%s", dist)
	}
	return dist
}

func (s *sender) checkSend(c *vCache, now time.Time) error {
	if now.Sub(c.setTime) < s.MinInterval.Duration() {
		return ErrSendTooFreq
	}
	if now.Sub(c.counterTime) > s.CounterDuration.Duration() {
		c.refresh(now)
		return nil
	}
	if c.sendCount > s.MaxCount {
		return ErrSendCountLimit
	}
	return nil
}

func (s *sender) checkVerify(c *vCache, code, hash string) error {
	c.updateVerify()
	if c.verifyCount > s.MaxVerifyCount {
		return ErrVerifyCodeRetryLimit
	}
	if c.code != code {
		return ErrVerifyCodeNotMatch
	}
	if c.hash != hash {
		return ErrVerifyCodeHashNotMatch
	}
	var now = time.Now()
	if now.Sub(c.setTime) > s.TTL.Duration() {
		return ErrVerifyCodeTimeout
	}
	return nil
}
