package async

import (
	"context"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"sync"
)

// RunnerQ : async runner
type RunnerQ struct {
	//queue
	q *Q

	//wait group
	wg *sync.WaitGroup

	//stopChan
	stopChan chan struct{}

	//start once
	startOnce sync.Once
	//stop once
	stopOnce sync.Once

	//set a name
	name string
}

// NewRunnerQ : new async runner queue
func NewRunnerQ(opts ...Option) *RunnerQ {
	var o = &optionT{
		size: DefaultQSize,
	}
	for _, opt := range opts {
		opt(o)
	}
	var c = &RunnerQ{}
	c.q = NewQ(o.size)
	c.wg = o.wg
	c.stopChan = make(chan struct{})
	c.name = o.name
	return c
}

// Size : get queue size
func (c *RunnerQ) Size() int {
	return c.q.size
}

// AsyncCall : async call
// ctx -- context.Context
// callCtxT -- call context
func (c *RunnerQ) AsyncCall(fn interface{}, ctx context.Context, arg interface{}) (interface{}, error) {
	var proc, err = c.addCallCtx(ctx, fn, arg)
	if err != nil {
		return nil, err
	}
	return proc.r()
}

// AsyncDelegate : async delegate
// ctx -- context.Context
// delegate -- delegate function
func (c *RunnerQ) AsyncDelegate(ctx context.Context, delegate Delegate) (interface{}, error) {
	var procCtx, err = c.addDelegateCtx(ctx, delegate)
	if err != nil {
		return nil, err
	}
	return procCtx.r()
}

// AsyncProc : async proc
// ctx -- context.Context
// proc -- proc interface
func (c *RunnerQ) AsyncProc(ctx context.Context, proc Proc) (interface{}, error) {
	var procCtx, err = c.addProcCtx(ctx, proc)
	if err != nil {
		return nil, err
	}
	return procCtx.r()
}

// Run : run all queue msg handler
func (c *RunnerQ) Run() {
	c.startOnce.Do(func() {
		if c.wg != nil {
			c.wg.Add(1)
		}
		go c.popLoop()
	})
}

// Stop : stop
func (c *RunnerQ) Stop() {
	c.stopOnce.Do(func() {
		c.q.Close()
	})
}

// WaitStop : wait runner loop exits
func (c *RunnerQ) WaitStop() {
	<-c.stopChan
}

// addCallCtx : add call context
func (c *RunnerQ) addCallCtx(ctx context.Context, fn interface{}, arg interface{}) (*callCtxT, error) {
	var callCtx = newCallCtx(ctx, fn, arg)
	var err = c.q.Add(callCtx)
	return callCtx, err
}

// addDelegateCtx : add delegate context
func (c *RunnerQ) addDelegateCtx(ctx context.Context, delegate Delegate) (*delegateCtxT, error) {
	var delegateCtx = newDelegateCtx(ctx, delegate)
	var err = c.q.Add(delegateCtx)
	return delegateCtx, err
}

// addProcCtx : add proc context
func (c *RunnerQ) addProcCtx(ctx context.Context, proc Proc) (*procCtxT, error) {
	var procCtx = newProcCtx(ctx, proc)
	var err = c.q.Add(procCtx)
	return procCtx, err
}

// pop call loop
func (c *RunnerQ) popLoop() {
	var (
		err  error
		item interface{}
		cc   ctxRunnerI
		ok   bool
	)

	defer func() {
		if c.wg != nil {
			c.wg.Done()
		}
		close(c.stopChan)
	}()

	for {

		item, err = c.q.PopAnyway()
		if err != nil {
			ulog.Debug("quit.in.line.handler",
				zap.String("name", c.name),
				zap.Error(err))
			return
		}
		cc, ok = item.(ctxRunnerI)
		if !ok {
			ulog.Error("invalid.async.line.call.context",
				zap.String("name", c.name),
				zap.Reflect("context", item))
			return
		}
		cc.run()
	}
}
