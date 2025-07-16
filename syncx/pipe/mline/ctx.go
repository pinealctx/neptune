package mline

import "context"

// CallFn : call function，回调函数
// Input:
// ctx -- context
// sIndex -- 表示在处理的go routine数组中对应的index，如果每个go routine有独立的缓存/内存，那此index可以用于定位相关的内存/缓存。
// param -- call param
// Output:
// r -- 回调函数调用后的返回值
// err - 回调函数调用失败后返回error
type CallFn func(ctx context.Context, sIndex int, req any) (rsp any, err error)

// CallCtx : call context, function and param
type CallCtx struct {
	call      CallFn
	param     any
	hashIndex int
}

// NewCallCtx ： new call context
// hashIndex : 一个与请求相关的散列值，例如如果用user_id作为散列分组，可以直接传user_id，如果是一个字符串，可以将此字符串CRC32散列一下。
// 传入的hashIndex应该是稳定的散列方式，例如在请求A中使用了CRC32作为散列方式，那所以相关的请求都应该用CRC32，
// 断不可在别的请求中使用类似xxhash这样别的散列方式。
// call : 回调函数
// param : 函数参数
func NewCallCtx(hashIndex int, call CallFn, param any) *CallCtx {
	return &CallCtx{
		call:      call,
		param:     param,
		hashIndex: hashIndex,
	}
}

// AsyncR : async call result.
type AsyncR struct {
	//result
	r any
	//error
	err error
}

// AsyncCtx : async call context
type AsyncCtx struct {
	//context
	ctx context.Context
	//async call function:
	call CallFn
	//call param
	param any
	//return chan
	rChan chan AsyncR
}

// newAsyncCtx : new async call context
// ctx -- context
// sIndex -- slot index
// call -- async call function
// param -- async call param
func newAsyncCtx(ctx context.Context, call CallFn, param any) *AsyncCtx {
	return &AsyncCtx{
		ctx:   ctx,
		call:  call,
		param: param,
		rChan: make(chan AsyncR, 1),
	}
}

// SetR : set return
func (m *AsyncCtx) SetR(r any, err error) {
	m.rChan <- AsyncR{
		r:   r,
		err: err,
	}
}

// R : get response with wait
func (m *AsyncCtx) R() (any, error) {
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	case rc := <-m.rChan:
		return rc.r, rc.err
	}
}
