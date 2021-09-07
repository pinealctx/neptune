package line

import "context"

//CallCtx : call context, function and param
type CallCtx struct {
	Call  func(ctx context.Context, req interface{}) (rsp interface{}, err error)
	Param interface{}
}

//AsyncR : async call result.
type AsyncR struct {
	//result
	r interface{}
	//error
	err error
}

//AsyncCtx : async call context
type AsyncCtx struct {
	//context
	ctx context.Context
	//async call function:
	//ctx -- context
	//param -- call param
	//r -- return param
	//err - return err if failed
	call func(ctx context.Context, param interface{}) (r interface{}, err error)
	//call param
	param interface{}
	//return chan
	rChan chan AsyncR
}

//NewAsyncCtx : new async call context
//ctx -- context
//call -- async call function
//param -- async call param
func NewAsyncCtx(ctx context.Context,
	call func(context.Context, interface{}) (interface{}, error), param interface{}) *AsyncCtx {
	return &AsyncCtx{
		ctx:   ctx,
		call:  call,
		param: param,
		rChan: make(chan AsyncR, 1),
	}
}

//SetR : set return
func (m *AsyncCtx) SetR(r interface{}, err error) {
	m.rChan <- AsyncR{
		r:   r,
		err: err,
	}
}

//R : get response with wait
func (m *AsyncCtx) R() (interface{}, error) {
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	case rc := <-m.rChan:
		return rc.r, rc.err
	}
}
