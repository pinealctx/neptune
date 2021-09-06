package raw

import (
	"context"
	"github.com/pinealctx/neptune/pipe"
	"github.com/pinealctx/neptune/pipe/q"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"sync"
)

type SlotAble interface {
	Slot() int
}

//Mux : multi-queue mux
type Mux struct {
	//slot size
	slotSize int
	//queue size in each slot
	qSize int

	//multi queues
	qs []*q.Q

	//wait group
	wg *sync.WaitGroup

	//go routine exit chan
	exitChan chan struct{}
	//stop once
	stopOnce sync.Once
}

//NewMux : new mux
func NewMux(opts ...pipe.Option) *Mux {
	//option
	var slotSize, qSize = pipe.GetOption(opts...)
	//new shunt then init
	return newMux(slotSize, qSize)
}

//NewMuxRun : new mux and run
func NewMuxRun(opts ...pipe.Option) *Mux {
	var c = NewMux(opts...)
	c.Run()
	return c
}

//newMux : new cycle with size
func newMux(slotSize int, qSizeInSlot int) *Mux {
	var c = &Mux{}
	c.slotSize, c.qSize = slotSize, qSizeInSlot
	c.wg = &sync.WaitGroup{}
	c.exitChan = make(chan struct{}, 1)
	c.wg.Add(c.slotSize)

	c.qs = make([]*q.Q, c.slotSize)
	for i := 0; i < c.slotSize; i++ {
		c.qs[i] = q.NewQ(q.WithQReqSize(c.qSize))
	}
	return c
}

//SlotSize : get slot size
func (c *Mux) SlotSize() int {
	return c.slotSize
}

//QSize : get queue size in each slot
func (c *Mux) QSize() int {
	return c.qSize
}

//CallSlot : wrap slot call
//ctx -- context.Context
//sIndex -- slot index
//callCtx -- call context
func (c *Mux) CallSlot(ctx context.Context, sIndex int, callCtx *SlotCallCtx) (interface{}, error) {
	var proc, err = c.addSlotCallCtx(ctx, sIndex, callCtx)
	if err != nil {
		return nil, err
	}
	return proc.R()
}

//Call : wrap call
//ctx -- context.Context
//callCtx -- call context
func (c *Mux) Call(ctx context.Context, callCtx *CallCtx) (interface{}, error) {
	var proc, err = c.addCallCtx(ctx, callCtx)
	if err != nil {
		return nil, err
	}
	return proc.R()
}

//Run : run all queue msg handler
func (c *Mux) Run() {
	for i := 0; i < c.slotSize; i++ {
		go c.popLoop(i)
	}
}

//Stop : stop
func (c *Mux) Stop() {
	c.stopOnce.Do(c.stop)
}

//WaitStop : wait stop
func (c *Mux) WaitStop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-c.exitChan:
		return nil
	}
}

//NormalizeSlotIndex slot index
func (c *Mux) NormalizeSlotIndex(index int) int {
	return pipe.NormalizeSlotIndex(index, c.slotSize)
}

//addSlotCallCtx : add slot call context
func (c *Mux) addSlotCallCtx(ctx context.Context, slotIndex int, callCtx *SlotCallCtx) (*AsyncCtx, error) {
	slotIndex = c.NormalizeSlotIndex(slotIndex)
	var proc = NewAsyncCtx(ctx, slotIndex, callCtx)
	var err = pipe.ConvertQueueErr(c.qs[slotIndex].AddReq(proc))
	return proc, err
}

//addCallCtx : add call context
func (c *Mux) addCallCtx(ctx context.Context, callCtx *CallCtx) (*AsyncCtx, error) {
	var slotIndex = c.NormalizeSlotIndex(callCtx.Param.Slot())
	var proc = NewAsyncCtxM(ctx, slotIndex, callCtx.Call, callCtx.Param)
	var err = pipe.ConvertQueueErr(c.qs[slotIndex].AddReq(proc))
	return proc, err
}

//pop msg loop
func (c *Mux) popLoop(index int) {
	var (
		err  error
		item interface{}
		ac   *AsyncCtx
		r    interface{}
		ok   bool

		mq = c.qs[index]
	)

	defer c.wg.Done()
	for {

		item, err = mq.PopAnyway()
		if err != nil {
			ulog.Debug("q.quit.in.raw.handler",
				zap.Int("index", index),
				zap.Error(err))
			return
		}
		ac, ok = item.(*AsyncCtx)
		if !ok {
			ulog.Error("invalid.async.call.context",
				zap.Int("index", index),
				zap.Reflect("context", item))
			return
		}

		r, err = ac.call(ac.ctx, ac.sIndex, ac.param)
		if err != nil {
			ac.SetR(nil, err)
		} else {
			ac.SetR(r, nil)
		}
	}
}

//stop work
func (c *Mux) stop() {
	for i := 0; i < c.slotSize; i++ {
		c.qs[i].Close()
	}
	//a go routine to wait all children done then signal it.
	go c.signalDone()
}

//signal all children go routine done
func (c *Mux) signalDone() {
	c.wg.Wait()
	c.exitChan <- struct{}{}
}
