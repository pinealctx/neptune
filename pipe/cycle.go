package pipe

import (
	"context"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"sync"
)

//QHandler : queue handler
type QHandler func(ctx context.Context, index int, req interface{}) (interface{}, error)

//Cycle : multi-queue cycle
type Cycle struct {
	//slot size
	slotSize int
	//queue size in each slot
	qSizeInSlot int

	//multi queues
	qs []*Q

	//queue handler
	qHandler QHandler

	//wait group
	wg *sync.WaitGroup

	//go routine exit chan
	exitChan chan struct{}
	//stop once
	stopOnce sync.Once
}

//NewCycle : new cycle
func NewCycle(qHandler QHandler, opts ...ShOption) *Cycle {
	//option
	var o = &_ShOption{
		slotSize:    DefaultSlotSize,
		qSizeInSlot: DefaultQInSlotSize,
	}
	for _, opt := range opts {
		opt(o)
	}
	//new shunt then init
	return newCycle(qHandler, o.slotSize, o.qSizeInSlot)
}

//NewRunCycle : new cycle and run
func NewRunCycle(qHandler QHandler, opts ...ShOption) *Cycle {
	var c = NewCycle(qHandler, opts...)
	c.Run()
	return c
}

//newCycle : new cycle with size
func newCycle(qHandler QHandler, slotSize int, qSizeInSlot int) *Cycle {
	var c = &Cycle{}
	c.slotSize, c.qSizeInSlot = slotSize, qSizeInSlot
	c.wg = &sync.WaitGroup{}
	c.exitChan = make(chan struct{}, 1)
	c.wg.Add(c.slotSize)

	c.qs = make([]*Q, c.slotSize)
	for i := 0; i < c.slotSize; i++ {
		c.qs[i] = NewQ(WithQReqSize(c.qSizeInSlot))
	}
	c.qHandler = qHandler
	return c
}

//SizeOfSlot : get slot size
func (c *Cycle) SizeOfSlot() int {
	return c.slotSize
}

//SizeOfQInSlot : get queue size in each slot
func (c *Cycle) SizeOfQInSlot() int {
	return c.qSizeInSlot
}

//Call : wrap call
func (c *Cycle) Call(ctx context.Context, slotIndex int, req interface{}) (interface{}, error) {
	var proc, err = c.addMsg(ctx, slotIndex, req)
	if err != nil {
		return nil, err
	}
	return proc.Rsp()
}

//CallPrior : wrap prior call
func (c *Cycle) CallPrior(ctx context.Context, slotIndex int, req interface{}) (interface{}, error) {
	var proc, err = c.addPriorMsg(ctx, slotIndex, req)
	if err != nil {
		return nil, err
	}
	return proc.Rsp()
}

//Run : run all queue msg handler
func (c *Cycle) Run() {
	for i := 0; i < c.slotSize; i++ {
		go c.popLoop(i)
	}
}

//Stop : stop
func (c *Cycle) Stop() {
	c.stopOnce.Do(c.stop)
}

//WaitStop : wait stop
func (c *Cycle) WaitStop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.exitChan:
		return nil
	}
}

//NormalizeSlotIndex slot index
func (c *Cycle) NormalizeSlotIndex(index int) int {
	return normalizeSlotIndex(index, c.slotSize)
}

//addMsg : add msg
func (c *Cycle) addMsg(ctx context.Context, slotIndex int, req interface{}) (*GenProc, error) {
	slotIndex = c.NormalizeSlotIndex(slotIndex)
	var proc = NewGenProc(ctx, slotIndex, req)
	var err = convertQueueErr(c.qs[slotIndex].AddReq(proc))
	return proc, err
}

//addPriorMsg : add prior msg
func (c *Cycle) addPriorMsg(ctx context.Context, slotIndex int, req interface{}) (*GenProc, error) {
	slotIndex = c.NormalizeSlotIndex(slotIndex)
	var proc = NewGenProc(ctx, slotIndex, req)
	var err = convertQueueErr(c.qs[slotIndex].AddPriorReq(proc))
	return proc, err
}

//pop msg loop
func (c *Cycle) popLoop(index int) {
	var (
		err  error
		item interface{}
		proc *GenProc
		rsp  interface{}
		ok   bool

		mq = c.qs[index]
	)

	defer c.wg.Done()
	for {

		item, err = mq.PopAnyway()
		if err != nil {
			ulog.Debug("msg.proc.module.item.quit",
				zap.Int("index", index),
				zap.Error(err))
			return
		}
		proc, ok = item.(*GenProc)
		if !ok {
			ulog.Error("msg.proc.module.item.invalid.req",
				zap.Int("index", index),
				zap.Reflect("inputMsg", item))
			return
		}
		if c.qHandler == nil {
			proc.SetRsp(nil, ErrNoHandler)
			ulog.Error("msg.proc.module.no.handler")
			return
		}
		rsp, err = c.qHandler(proc.ctx, index, proc.req)
		if err != nil {
			proc.SetRsp(nil, err)
		} else {
			proc.SetRsp(rsp, nil)
		}

	}
}

//stop work
func (c *Cycle) stop() {
	for i := 0; i < c.slotSize; i++ {
		c.qs[i].Close()
	}
	//a go routine to wait all children done then signal it.
	go c.signalDone()
}

//signal all children go routine done
func (c *Cycle) signalDone() {
	c.wg.Wait()
	c.exitChan <- struct{}{}
}

//normalizeSlotIndex slot index
func normalizeSlotIndex(index int, slotSize int) int {
	if index < 0 {
		index = -index
	}
	index %= slotSize
	return index
}

//convert msg queue error
func convertQueueErr(err error) error {
	if err == nil {
		return nil
	}
	if err == ErrReqQFull {
		return ErrQueueFull
	}
	if err == ErrClosed {
		return ErrQueueClosed
	}
	return err
}
