package async

import (
	"context"
	"reflect"
)

// Delegate : proc delegate function
type Delegate func(ctx context.Context) (interface{}, error)

// Proc : proc interface
type Proc interface {
	//Do : do proc
	Do(ctx context.Context) (interface{}, error)
}

// ctx runner
type ctxRunnerI interface {
	//run
	run()
	//get result
	r() (interface{}, error)
}

// callCtxT call function context
// An async wrapped function must follow such template:
// func(ctx context, param AnyType) (result AnyType, err error)
// "AnyType" here means any type
type callCtxT struct {
	//context
	ctx context.Context
	//call function type
	functionType reflect.Type
	//call function value -- pointer
	functionValue reflect.Value
	//call arg -- exclude first arg --- (ctx context)
	arg interface{}
	//wait -- wait chan
	wait chan struct{}
	//result -- exclude last arg -- (err error)
	result interface{}
	//error
	err error
}

// new call context
func newCallCtx(ctx context.Context, fn interface{}, arg interface{}) *callCtxT {
	var fnType, valid = validateFn(fn)
	if !valid {
		panic("new async call in case function is nil")
	}
	return &callCtxT{
		ctx:           ctx,
		functionType:  fnType,
		functionValue: reflect.ValueOf(fn),
		arg:           arg,
		wait:          make(chan struct{}),
	}
}

// r : get result with wait
func (c *callCtxT) r() (interface{}, error) {
	select {
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	case <-c.wait:
		return c.result, c.err
	}
}

// run
func (c *callCtxT) run() {
	var params [2]reflect.Value

	defer close(c.wait)

	select {
	//if context done, return
	case <-c.ctx.Done():
		c.err = c.ctx.Err()
		return
	default:
	}

	params[0] = reflect.ValueOf(c.ctx)
	params[1] = reflect.ValueOf(c.arg)
	var rets = c.functionValue.Call(params[:])
	c.result = rets[0].Interface()
	if !rets[1].IsNil() {
		c.err = rets[1].Interface().(error)
	}
}

// delegateCtxT : proc context, interface
type delegateCtxT struct {
	//context
	ctx context.Context
	//delegate function
	delegate Delegate
	//wait -- wait chan
	wait chan struct{}
	//result -- exclude last arg -- (err error)
	result interface{}
	//error
	err error
}

// new delegate context
func newDelegateCtx(ctx context.Context, delegate Delegate) *delegateCtxT {
	return &delegateCtxT{
		ctx:      ctx,
		delegate: delegate,
		wait:     make(chan struct{}),
	}
}

// r : get result with wait
func (c *delegateCtxT) r() (interface{}, error) {
	select {
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	case <-c.wait:
		return c.result, c.err
	}
}

// run
func (c *delegateCtxT) run() {
	defer close(c.wait)

	select {
	//if context done, return
	case <-c.ctx.Done():
		c.err = c.ctx.Err()
		return
	default:
	}
	c.result, c.err = c.delegate(c.ctx)
}

// procCtxT : proc context, interface
type procCtxT struct {
	//context
	ctx context.Context
	//proc interface
	proc Proc
	//wait -- wait chan
	wait chan struct{}
	//result -- exclude last arg -- (err error)
	result interface{}
	//error
	err error
}

// new proc context
func newProcCtx(ctx context.Context, proc Proc) *procCtxT {
	return &procCtxT{
		ctx:  ctx,
		proc: proc,
		wait: make(chan struct{}),
	}
}

// r : get result with wait
func (c *procCtxT) r() (interface{}, error) {
	select {
	case <-c.ctx.Done():
		return nil, c.ctx.Err()
	case <-c.wait:
		return c.result, c.err
	}
}

// run
func (c *procCtxT) run() {
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
