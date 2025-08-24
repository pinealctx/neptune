package mline

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/pinealctx/neptune/syncx/pipe"
	"github.com/pinealctx/neptune/syncx/pipe/q"
	"github.com/pinealctx/neptune/ulog"
)

// MultiLine : multi-queue handler
type MultiLine struct {
	//slot size
	slotSize int
	//queue size in each slot
	qSize int

	//multi queues
	qs []*q.Q[*AsyncCtx]

	//wait group
	wg *sync.WaitGroup

	//go routine exit chan
	exitChan chan struct{}
	//stop once
	stopOnce sync.Once
}

// NewMultiLine : new multi-queue group
func NewMultiLine(opts ...pipe.Option) *MultiLine {
	//option
	var slotSize, qSize = pipe.GetOption(opts...)
	//new shunt then init
	return newMux(slotSize, qSize)
}

// newMux : new cycle with size
func newMux(slotSize int, qSizeInSlot int) *MultiLine {
	var c = &MultiLine{}
	c.slotSize, c.qSize = slotSize, qSizeInSlot
	c.wg = &sync.WaitGroup{}
	c.exitChan = make(chan struct{}, 1)
	c.wg.Add(c.slotSize)

	c.qs = make([]*q.Q[*AsyncCtx], c.slotSize)
	for i := 0; i < c.slotSize; i++ {
		c.qs[i] = q.NewQ[*AsyncCtx](c.qSize)
	}
	return c
}

// SlotSize : get slot size
func (c *MultiLine) SlotSize() int {
	return c.slotSize
}

// QSize : get queue size in each slot
func (c *MultiLine) QSize() int {
	return c.qSize
}

// IndexOf : get slot index.
// 将散列值映射成处理数组的index，举例来说，如果以user id作为散列值，则整个处理逻辑会用user id的绝对值对处理数组长度取模，取模后的值就是
// 其在数组中的位置。
// Input : i -- a slot key number. 此参数就是分片使用的hash值。
// Output : index of slot.返回此hash值在处理数组中对应的位置。
func (c *MultiLine) IndexOf(i int) int {
	return pipe.NormalizeSlotIndex(i, c.slotSize)
}

// AsyncCall : wrap call
// ctx -- context.Context
// callCtx -- call context
func (c *MultiLine) AsyncCall(ctx context.Context, callCtx *CallCtx) (any, error) {
	var proc, err = c.addCallCtx(ctx, callCtx)
	if err != nil {
		return nil, err
	}
	return proc.R()
}

// Run : run all queue msg handler
func (c *MultiLine) Run() {
	for i := 0; i < c.slotSize; i++ {
		go c.popLoop(i)
	}
}

// Stop : stop
func (c *MultiLine) Stop() {
	c.stopOnce.Do(c.stop)
}

// WaitStop : wait stop
func (c *MultiLine) WaitStop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.exitChan:
		return nil
	}
}

// addCallCtx : add call context
func (c *MultiLine) addCallCtx(ctx context.Context, callCtx *CallCtx) (*AsyncCtx, error) {
	var slotIndex = pipe.NormalizeSlotIndex(callCtx.hashIndex, c.slotSize)
	var proc = newAsyncCtx(ctx, callCtx.call, callCtx.param)
	var err = pipe.ConvertQueueErr(c.qs[slotIndex].Push(proc))
	return proc, err
}

// pop msg loop
func (c *MultiLine) popLoop(index int) {
	var (
		err error
		ac  *AsyncCtx
		r   any

		mq = c.qs[index]
	)

	defer c.wg.Done()
	for {
		ac, err = mq.Pop()
		if err != nil {
			ulog.Debug("q.quit.in.raw.handler",
				zap.Int("index", index),
				zap.Error(err))
			return
		}

		r, err = ac.call(ac.ctx, index, ac.param)
		if err != nil {
			ac.SetR(nil, err)
		} else {
			ac.SetR(r, nil)
		}
	}
}

// stop work
func (c *MultiLine) stop() {
	for i := 0; i < c.slotSize; i++ {
		c.qs[i].Close()
	}
	//a go routine to wait all children done then signal it.
	go c.signalDone()
}

// signal all children go routine done
func (c *MultiLine) signalDone() {
	c.wg.Wait()
	c.exitChan <- struct{}{}
}
