package mux

import (
	"context"
	"sync"

	"go.uber.org/zap"

	"github.com/pinealctx/neptune/ulog"
)

// Worker : a go routine to handle queue work
type Worker struct {
	workQ *Q
	wg    *sync.WaitGroup
	ca    CacheFacade
}

func NewWorker(qSize int, wg *sync.WaitGroup, ca CacheFacade) *Worker {
	return &Worker{
		workQ: NewQ(qSize),
		wg:    wg,
		ca:    ca,
	}
}

// DoGet : get from cache first if not load from db
func (w *Worker) DoGet(ctx context.Context, loadFn RenewDataFn, k any) (any, error) {
	var v, ok = w.ca.Get(k)
	if ok {
		//hit in cache
		return v, nil
	}
	return w.asyncCall(ctx, NewLoad(loadFn, k))
}

// DoAdd : add item
func (w *Worker) DoAdd(ctx context.Context, addFn RenewDataFn, k any, data any) (any, error) {
	return w.asyncCall(ctx, NewAdd(addFn, k, data))
}

// DoUpdate : update item
func (w *Worker) DoUpdate(ctx context.Context,
	loadFn RenewDataFn, updFn UpdateDataFn, k any, data any) (any, error) {
	return w.asyncCall(ctx, NewUpdate(loadFn, updFn, k, data))
}

// DoDelete : delete item
func (w *Worker) DoDelete(ctx context.Context, deleteFn DeleteFn, k any) (any, error) {
	return w.asyncCall(ctx, NewDelete(deleteFn, k))
}

// DoUpdOrAddIfNull :
// 1. load.
// 2. update if existed.
// 3. add if not existed.
func (w *Worker) DoUpdOrAddIfNull(ctx context.Context,
	loadFn RenewDataFn, updFn UpdateDataFn, addFn RenewDataFn, isNotFoundFn IsNotFoundFn,
	k any, data any) (any, error) {
	return w.asyncCall(ctx, NewMixUpdOrAddIfNull(loadFn, updFn, addFn, isNotFoundFn, k, data))
}

// DoUpsertThenLoad :
// 1. upsert.
// 2. update cache if cache hit.
// 3. load cache if cache miss.
func (w *Worker) DoUpsertThenLoad(ctx context.Context,
	upsertFn UpdateDataFn, loadFn RenewDataFn, k any, data any) (any, error) {
	return w.asyncCall(ctx, NewMixUpsertThenLoad(upsertFn, loadFn, k, data))
}

// DoUpsertThenRenewInCache :
// 1. upsert.
// 2. update cache if cache hit.
func (w *Worker) DoUpsertThenRenewInCache(ctx context.Context,
	upsertFn UpdateDataFn, k any, data any) (any, error) {
	return w.asyncCall(ctx, NewMixUpsertThenRenewInCache(upsertFn, k, data))
}

// Start : start handler go routine
func (w *Worker) Start() {
	go w.runLoop()
}

// Stop : close queue, not accept input anymore.
func (w *Worker) Stop() {
	w.workQ.Close()
}

// async call
func (w *Worker) asyncCall(ctx context.Context, op OpCode) (any, error) {
	var c = NewAsync(ctx, op)
	var err = w.workQ.AddReq(c)
	if err != nil {
		return nil, err
	}
	return c.R()
}

// loop go routine to handle async call
func (w *Worker) runLoop() {
	var (
		e   any
		c   *AsyncC
		err error
	)
	defer w.wg.Done()

	for {
		e, err = w.workQ.PopAnyway()
		if err != nil {
			ulog.Debug("work.module.quit",
				zap.Error(err))
			return
		}
		c = e.(*AsyncC)
		w.handleAsync(c)
	}
}

// async handler entry
func (w *Worker) handleAsync(c *AsyncC) {
	switch op := c.op.(type) {
	case *OpLoad:
		w.handleLoad(c, op)
	case *OpAdd:
		w.handleAdd(c, op)
	case *OpUpdate:
		w.handleUpdate(c, op)
	case *OpDelete:
		w.handleDelete(c, op)
	case *OpMixUpdOrAddIfNull:
		w.handleMixUpdOrAddIfNull(c, op)
	case *OpMixUpsertThenLoad:
		w.handleMixUpsertThenLoad(c, op)
	case *OpMixUpsertThenRenewInCache:
		w.handleMixUpsertThenRenewInCache(c, op)
	}
}

// handle load
func (w *Worker) handleLoad(c *AsyncC, op *OpLoad) {
	var v, ok = w.ca.Get(op.k)
	if ok {
		//hit in cache
		c.SetR(v, nil)
		return
	}
	var err error
	v, err = op.loadFn(c.ctx, op.k)
	if err != nil {
		//load error
		c.SetR(nil, err)
		return
	}
	//renew cache
	w.ca.Set(op.k, v)
	//set result
	c.SetR(v, nil)
}

// handle add
func (w *Worker) handleAdd(c *AsyncC, op *OpAdd) {
	var _, exist = w.ca.Peek(op.k)
	if exist {
		//key duplicate
		c.SetR(nil, ErrDupKey)
		return
	}
	var v, err = op.addFn(c.ctx, op.data)
	if err != nil {
		//add error
		c.SetR(nil, err)
		return
	}
	//renew cache
	w.ca.Set(op.k, v)
	//set result
	c.SetR(v, nil)
}

// handle update
func (w *Worker) handleUpdate(c *AsyncC, op *OpUpdate) {
	var pre, ok = w.ca.Peek(op.k)
	if ok {
		var v, err = op.updFn(c.ctx, op.data, pre)
		if err != nil {
			//update error
			c.SetR(nil, err)
			return
		}
		//renew cache
		w.ca.Set(op.k, v)
		//set result
		c.SetR(v, nil)
		return
	}
	var v, err = op.loadFn(c.ctx, op.k)
	if err != nil {
		//load error
		c.SetR(nil, err)
		return
	}
	v, err = op.updFn(c.ctx, op.data, v)
	if err != nil {
		//update error
		c.SetR(nil, err)
		return
	}
	//renew cache
	w.ca.Set(op.k, v)
	//set result
	c.SetR(v, nil)
}

// handle delete
func (w *Worker) handleDelete(c *AsyncC, op *OpDelete) {
	var err = op.deleteFn(c.ctx, op.k)
	if err != nil {
		//delete error
		c.SetR(nil, err)
		return
	}
	//renew cache
	w.ca.Delete(op.k)
	//set result
	c.SetR(nil, nil)
}

// handle update if exist else add if not exist.
func (w *Worker) handleMixUpdOrAddIfNull(c *AsyncC, op *OpMixUpdOrAddIfNull) {
	var pre, ok = w.ca.Peek(op.k)
	if ok {
		var v, err = op.updFn(c.ctx, op.data, pre)
		if err != nil {
			//update error
			c.SetR(nil, err)
			return
		}
		//renew cache
		w.ca.Set(op.k, v)
		//set result
		c.SetR(v, nil)
		return
	}
	var v, err = op.loadFn(c.ctx, op.k)
	if err != nil {
		if !op.isNotFoundFn(err) {
			//other error, not "not found" error
			//load error
			c.SetR(nil, err)
			return
		}
		v, err = op.addFn(c.ctx, op.data)
		if err != nil {
			//add error
			c.SetR(nil, err)
			return
		}
		//renew cache
		w.ca.Set(op.k, v)
		//set result
		c.SetR(v, nil)
		return
	}
	v, err = op.updFn(c.ctx, op.data, v)
	if err != nil {
		//update error
		c.SetR(nil, err)
		return
	}
	//renew cache
	w.ca.Set(op.k, v)
	//set result
	c.SetR(v, nil)
}

// handle upsert first the re-new cache
func (w *Worker) handleMixUpsertThenLoad(c *AsyncC, op *OpMixUpsertThenLoad) {
	var pre, ok = w.ca.Peek(op.k)
	if ok {
		var v, err = op.upsertFn(c.ctx, op.data, pre)
		if err != nil {
			//upsert error
			c.SetR(nil, err)
			return
		}
		//renew cache
		w.ca.Set(op.k, v)
		//set result
		c.SetR(v, nil)
		return
	}

	var _, err = op.upsertFn(c.ctx, op.data, nil)
	if err != nil {
		//upsert error
		c.SetR(nil, err)
		return
	}
	var v any
	v, err = op.loadFn(c.ctx, op.k)
	if err != nil {
		//load error
		c.SetR(nil, err)
		return
	}
	//renew cache
	w.ca.Set(op.k, v)
	//set result
	c.SetR(v, nil)
}

// handle upsert first the re-new cache if cache hit
func (w *Worker) handleMixUpsertThenRenewInCache(c *AsyncC, op *OpMixUpsertThenRenewInCache) {
	var pre, ok = w.ca.Peek(op.k)
	if ok {
		var v, err = op.upsertFn(c.ctx, op.data, pre)
		if err != nil {
			//upsert error
			c.SetR(nil, err)
			return
		}
		//renew cache
		w.ca.Set(op.k, v)
		//set result
		c.SetR(v, nil)
		return
	}

	//not in cache, just upsert, do not renew cache.
	var v, err = op.upsertFn(c.ctx, op.data, nil)
	if err != nil {
		//upsert error
		c.SetR(nil, err)
		return
	}
	//set result
	c.SetR(v, nil)
}
