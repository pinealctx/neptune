package line

import (
	"context"
	"github.com/pinealctx/neptune/pipe"
	"github.com/pinealctx/neptune/pipe/q"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"sync"
)

type _Option struct {
	qSize int
}

//Option : only qSize option
type Option func(option *_Option)

//WithQSize : setup qSize
func WithQSize(qSize int) Option {
	return func(o *_Option) {
		o.qSize = qSize
	}
}

//Line : async runner
type Line struct {
	//queue size
	qSize int

	//queue
	q *q.Q

	//wait group
	wg *sync.WaitGroup

	//start once
	startOnce sync.Once
	//stop once
	stopOnce sync.Once
}

//NewLine : new async line
func NewLine(wg *sync.WaitGroup, opts ...Option) *Line {
	var o = &_Option{
		qSize: pipe.DefaultQSize,
	}
	for _, opt := range opts {
		opt(o)
	}
	return newLine(o.qSize, wg)
}

//newLine : new async line
func newLine(qSize int, wg *sync.WaitGroup) *Line {
	var c = &Line{}
	c.qSize = qSize
	c.q = q.NewQ(q.WithQReqSize(c.qSize))
	c.wg = wg
	return c
}

//QSize : get queue size in each slot
func (c *Line) QSize() int {
	return c.qSize
}

//AsyncCall : wrap call
//ctx -- context.Context
//callCtx -- call context
func (c *Line) AsyncCall(ctx context.Context, callCtx *CallCtx) (interface{}, error) {
	var proc, err = c.addCallCtx(ctx, callCtx)
	if err != nil {
		return nil, err
	}
	return proc.R()
}

//Run : run all queue msg handler
func (c *Line) Run() {
	c.startOnce.Do(func() {
		c.wg.Add(1)
		go c.popLoop()
	})
}

//Stop : stop
func (c *Line) Stop() {
	c.stopOnce.Do(func() {
		c.q.Close()
	})
}

//addCallCtx : add call context
func (c *Line) addCallCtx(ctx context.Context, callCtx *CallCtx) (*AsyncCtx, error) {
	var proc = NewAsyncCtx(ctx, callCtx.Call, callCtx.Param)
	var err = pipe.ConvertQueueErr(c.q.AddReq(proc))
	return proc, err
}

//pop call loop
func (c *Line) popLoop() {
	var (
		err  error
		item interface{}
		ac   *AsyncCtx
		r    interface{}
		ok   bool
	)

	defer c.wg.Done()
	for {

		item, err = c.q.PopAnyway()
		if err != nil {
			ulog.Debug("q.quit.in.line.handler",
				zap.Error(err))
			return
		}
		ac, ok = item.(*AsyncCtx)
		if !ok {
			ulog.Error("invalid.async.line.call.context",
				zap.Reflect("context", item))
			return
		}

		r, err = ac.call(ac.ctx, ac.param)
		if err != nil {
			ac.SetR(nil, err)
		} else {
			ac.SetR(r, nil)
		}
	}
}
