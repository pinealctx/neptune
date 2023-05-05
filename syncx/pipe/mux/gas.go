package mux

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrDupKey = status.Error(codes.AlreadyExists, "gas.duplicated.key")
)

// OpType : op type
type OpType int

// CacheFacade : cache facade, could be lru cache, map, or both mixed.
type CacheFacade interface {
	//Peek : only useful in lru cache, peek means no update LRU order.
	Peek(key interface{}) (interface{}, bool)
	//Get : get from cache, in lru cache, key order also be updated.
	Get(key interface{}) (interface{}, bool)
	//Set : set to cache
	Set(key interface{}, value interface{})
	//Delete : delete key from cache
	Delete(key interface{})
}

// R : combine interface and error, for return result
type R struct {
	//return value
	r interface{}
	//error
	err error
}

// AsyncC : async cell
type AsyncC struct {
	//context
	ctx context.Context
	//opcode
	op OpCode
	//return value, actually it's chan for async call
	rChan chan R
}

// NewAsync : new async call
func NewAsync(ctx context.Context, op OpCode) *AsyncC {
	return &AsyncC{
		ctx:   ctx,
		op:    op,
		rChan: make(chan R, 1),
	}
}

// SetR : set op result
func (a *AsyncC) SetR(r interface{}, err error) {
	a.rChan <- R{
		r:   r,
		err: err,
	}
}

// R : get result
func (a *AsyncC) R() (interface{}, error) {
	select {
	case <-a.ctx.Done():
		return nil, a.ctx.Err()
	case re := <-a.rChan:
		return re.r, re.err
	}
}

// RenewDataFn : renew data function, including load/add, excluding delete.
// input: ctx->context:can be ignored in case; d->input data
// output: v-> return cache value; err->if failed, return err
type RenewDataFn func(ctx context.Context, d interface{}) (v interface{}, err error)

// UpdateDataFn : update data function.
// input: ctx->context:can be ignored in case; d->input data; e->the existed item.
// output: v-> return cache value; err->if failed, return err.
// actually e is from load function if it's not in cache.
type UpdateDataFn func(ctx context.Context, d interface{}, e interface{}) (v interface{}, err error)

// DeleteFn : delete data function.
// input: ctx->context:can be ignored in case; d->input data.
// output: error->if failed, return err.
type DeleteFn func(ctx context.Context, d interface{}) error

// IsNotFoundFn : to detective an error is "not found" or not.
type IsNotFoundFn func(err error) bool

type OpCode interface {
	GetK() interface{}
}

// OpLoad : wrapped load command
type OpLoad struct {
	//loadFn: load item
	loadFn RenewDataFn
	//k: the key in cache
	k interface{}
}

func NewLoad(l RenewDataFn, k interface{}) OpCode {
	return &OpLoad{
		loadFn: l,
		k:      k,
	}
}

func (o *OpLoad) GetK() interface{} {
	return o.k
}

// OpAdd : wrapped add command
type OpAdd struct {
	//addFn: add item
	addFn RenewDataFn
	//k: the key in cache
	k interface{}
	//data: input data
	data interface{}
}

func NewAdd(a RenewDataFn, k interface{}, data interface{}) OpCode {
	return &OpAdd{
		addFn: a,
		k:     k,
		data:  data,
	}
}

func (o *OpAdd) GetK() interface{} {
	return o.k
}

// OpUpdate : wrapped update command
// Actually : update should after load data
type OpUpdate struct {
	//loadFn : load item
	loadFn RenewDataFn
	//updFn : update item
	updFn UpdateDataFn
	//k: the key in cache
	k interface{}
	//data : input data
	data interface{}
}

func NewUpdate(l RenewDataFn, u UpdateDataFn, k interface{}, data interface{}) OpCode {
	return &OpUpdate{
		loadFn: l,
		updFn:  u,
		k:      k,
		data:   data,
	}
}

func (o *OpUpdate) GetK() interface{} {
	return o.k
}

// OpDelete : wrapped delete command
type OpDelete struct {
	//deleteFn : delete item
	deleteFn DeleteFn
	//k: the key in cache
	k interface{}
}

func NewDelete(d DeleteFn, k interface{}) OpCode {
	return &OpDelete{
		deleteFn: d,
		k:        k,
	}
}

func (o *OpDelete) GetK() interface{} {
	return o.k
}

//OpMixUpdOrAddIfNull : 1.load. 2. update it if existed. 3. add it if not existed.
/*
load;
if not found -> add.
else update.
*/
type OpMixUpdOrAddIfNull struct {
	//loadFn : load item
	loadFn RenewDataFn
	//updFn : update item
	updFn UpdateDataFn
	//addFn : add item
	addFn RenewDataFn
	//isNotFoundFn : if not found error or not
	isNotFoundFn IsNotFoundFn
	//k: the key in cache
	k interface{}
	//data : input data
	data interface{}
}

func NewMixUpdOrAddIfNull(l RenewDataFn, u UpdateDataFn, a RenewDataFn, i IsNotFoundFn,
	k interface{}, data interface{}) OpCode {
	return &OpMixUpdOrAddIfNull{
		loadFn:       l,
		updFn:        u,
		addFn:        a,
		isNotFoundFn: i,
		k:            k,
		data:         data,
	}
}

func (o *OpMixUpdOrAddIfNull) GetK() interface{} {
	return o.k
}

//OpMixUpsertThenLoad : upsert it first, then load it if not in cache.
/*
1. upsert.
2. if in cache, refresh cache.
3. if not in cache, reload from cache.
*/
type OpMixUpsertThenLoad struct {
	//upsertFn : upsert item
	upsertFn UpdateDataFn
	//loadFn : load item
	loadFn RenewDataFn
	//k: the key in cache
	k interface{}
	//data : input data
	data interface{}
}

func NewMixUpsertThenLoad(p UpdateDataFn, l RenewDataFn, k interface{}, data interface{}) OpCode {
	return &OpMixUpsertThenLoad{
		upsertFn: p,
		loadFn:   l,
		k:        k,
		data:     data,
	}
}

func (o *OpMixUpsertThenLoad) GetK() interface{} {
	return o.k
}

//OpMixUpsertThenRenewInCache : upsert it first, then renew in cache item.
/*
1. upsert.
2. if in cache, renew cache.
3. if not in cache, do nothing.
*/
type OpMixUpsertThenRenewInCache struct {
	//upsertFn : upsert item
	upsertFn UpdateDataFn
	//k: the key in cache
	k interface{}
	//data : input data
	data interface{}
}

func NewMixUpsertThenRenewInCache(p UpdateDataFn, k interface{}, data interface{}) OpCode {
	return &OpMixUpsertThenRenewInCache{
		upsertFn: p,
		k:        k,
		data:     data,
	}
}

func (o *OpMixUpsertThenRenewInCache) GetK() interface{} {
	return o.k
}
