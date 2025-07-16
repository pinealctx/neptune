package vcode

import (
	"time"

	"github.com/pinealctx/neptune/idgen/random"
)

type vCache struct {
	counterTime time.Time
	setTime     time.Time
	sendCount   int
	verifyCount int
	code        string
	hash        string
}

func newSenderCache(now time.Time) *vCache {
	return &vCache{counterTime: now}
}

func (c *vCache) refresh(now time.Time) {
	c.counterTime = now
	c.sendCount = 0
}

func (c *vCache) updateSend(code string, now time.Time) {
	c.code = code
	c.hash = random.MD5UUID()
	c.setTime = now
	c.sendCount++
	c.verifyCount = 0
}

func (c *vCache) updateVerify() {
	c.verifyCount++
}

func (c *vCache) Size() int {
	return 1
}
