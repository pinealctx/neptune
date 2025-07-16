package async

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/pinealctx/neptune/ulog"
)

// procChanCtxT : proc context, interface
type procChanCtxT struct {
	//context
	ctx context.Context
	//proc interface
	proc Proc
	//wait -- wait chan
	wait chan struct{}
	//result -- exclude last arg -- (err error)
	result any
	//error
	err error
}

// new proc chan context
func newProcChanCtx(ctx context.Context, proc Proc) *procChanCtxT {
	return &procChanCtxT{
		ctx:  ctx,
		proc: proc,
		wait: make(chan struct{}),
	}
}

// r : get result with wait
func (c *procChanCtxT) r(stopChan <-chan struct{}) (any, error) {
	select {
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	case <-stopChan:
		return nil, ErrClosed
	case <-c.wait:
		return c.result, c.err
	}
}

// run
func (c *procChanCtxT) run() {
	defer close(c.wait)

	select {
	//if context done, return
	case <-c.ctx.Done():
		c.err = c.ctx.Err()
		return
	default:
	}
	c.result, c.err = c.proc.Do(c.ctx)
}

// ProcChan : async proc chan
type ProcChan struct {
	//chan
	ch chan *procChanCtxT

	//wait group
	wg *sync.WaitGroup

	//stopChan
	stopChan chan struct{}

	//size
	size int

	//start once
	startOnce sync.Once
	//stop once
	stopOnce sync.Once

	//set a name
	name string
}

// NewProcChan : new async proc queue
func NewProcChan(opts ...Option) *ProcChan {
	var o = &optionT{
		size: DefaultQSize,
	}
	for _, opt := range opts {
		opt(o)
	}
	var c = &ProcChan{
		size: o.size,
	}
	c.ch = make(chan *procChanCtxT, o.size)
	c.wg = o.wg
	c.stopChan = make(chan struct{})
	c.name = o.name
	return c
}

// Size : get queue size
func (c *ProcChan) Size() int {
	return c.size
}

// AsyncProc : async proc
// ctx -- context.Context
// proc -- proc interface
func (c *ProcChan) AsyncProc(ctx context.Context, proc Proc) (any, error) {
	var procCtx, err = c.addCallCtx(ctx, proc)
	if err != nil {
		return nil, err
	}
	return procCtx.r(c.stopChan)
}

// Run : run all queue msg handler
func (c *ProcChan) Run() {
	c.startOnce.Do(func() {
		if c.wg != nil {
			c.wg.Add(1)
		}
		go c.popLoop()
	})
}

// Stop : stop
func (c *ProcChan) Stop() {
	c.stopOnce.Do(func() {
		close(c.stopChan)
	})
}

// WaitStop : wait runner loop exits
func (c *ProcChan) WaitStop() {
	<-c.stopChan
}

// addCallCtx : add call context
func (c *ProcChan) addCallCtx(ctx context.Context, proc Proc) (*procChanCtxT, error) {
	var procCtx = newProcChanCtx(ctx, proc)
	select {
	case c.ch <- procCtx:
		return procCtx, nil
	case <-c.stopChan:
		return procCtx, ErrClosed
	default:
		return procCtx, ErrFull
	}
}

// pop call loop
func (c *ProcChan) popLoop() {
	var (
		err error
		cc  *procChanCtxT
	)

	defer func() {
		if c.wg != nil {
			c.wg.Done()
		}
	}()

	for {

		select {
		case cc = <-c.ch:
			cc.run()
		case <-c.stopChan:
			ulog.Debug("quit.in.proc.handler",
				zap.String("name", c.name),
				zap.Error(err))
			return
		}
	}
}
