package mux

import (
	"context"
	"sync"
)

//CacheGen : cache facade generator
type CacheGen func() CacheFacade

//WorkerGrp : worker group
type WorkerGrp struct {
	//mux size
	muxSize int
	//queue deep size
	deepSize int

	//multi workers
	ws []*Worker

	//wait group
	wg *sync.WaitGroup

	//go routine exit signal chan
	exitSignal chan struct{}
	//stop once
	stopOnce sync.Once
}

//NewWorkGrp : new work group
func NewWorkGrp(cg CacheGen, opts ...Option) *WorkerGrp {
	var w = buildWorkGrp(opts...)
	for i := 0; i < w.muxSize; i++ {
		var ca = cg()
		w.ws[i] = NewWorker(w.deepSize, w.wg, ca)
	}
	return w
}

//NewWorkGrpWithMapCache : new work group bind with map cache.
func NewWorkGrpWithMapCache(opts ...Option) *WorkerGrp {
	var w = buildWorkGrp(opts...)
	for i := 0; i < w.muxSize; i++ {
		w.ws[i] = NewWorker(w.deepSize, w.wg, NewFacadeMap())
	}
	return w
}

//NewWorkGrpWithLRU : new work group bind with lru cache.
func NewWorkGrpWithLRU(lruCap int64, opts ...Option) *WorkerGrp {
	var w = buildWorkGrp(opts...)
	for i := 0; i < w.muxSize; i++ {
		w.ws[i] = NewWorker(w.deepSize, w.wg, NewFacadeLRU(lruCap))
	}
	return w
}

//DoGet : get from cache first if not load from db
func (w *WorkerGrp) DoGet(ctx context.Context, loadFn RenewDataFn, k Hashed2Int) (interface{}, error) {
	return w.ws[w.locHash(k)].DoGet(ctx, loadFn, k)
}

//DoAdd : add item
func (w *WorkerGrp) DoAdd(ctx context.Context, addFn RenewDataFn, k Hashed2Int, data interface{}) (interface{}, error) {
	return w.ws[w.locHash(k)].DoAdd(ctx, addFn, k, data)
}

//DoUpdate : update item
func (w *WorkerGrp) DoUpdate(ctx context.Context,
	loadFn RenewDataFn, updFn UpdateDataFn, k Hashed2Int, data interface{}) (interface{}, error) {
	return w.ws[w.locHash(k)].DoUpdate(ctx, loadFn, updFn, k, data)
}

//DoDelete : delete item
func (w *WorkerGrp) DoDelete(ctx context.Context, deleteFn DeleteFn, k Hashed2Int) (interface{}, error) {
	return w.ws[w.locHash(k)].DoDelete(ctx, deleteFn, k)
}

//DoUpdOrAddIfNull :
//1. load.
//2. update if existed.
//3. add if not existed.
func (w *WorkerGrp) DoUpdOrAddIfNull(ctx context.Context,
	loadFn RenewDataFn, updFn UpdateDataFn, addFn RenewDataFn, isNotFoundFn IsNotFoundFn,
	k Hashed2Int, data interface{}) (interface{}, error) {
	return w.ws[w.locHash(k)].DoUpdOrAddIfNull(ctx, loadFn, updFn, addFn, isNotFoundFn, k, data)
}

//DoUpsertThenLoad :
//1. upsert.
//2. update cache if cache hit.
//3. load cache if cache miss.
func (w *WorkerGrp) DoUpsertThenLoad(ctx context.Context,
	upsertFn UpdateDataFn, loadFn RenewDataFn, k Hashed2Int, data interface{}) (interface{}, error) {
	return w.ws[w.locHash(k)].DoUpsertThenLoad(ctx, upsertFn, loadFn, k, data)
}

//DoUpsertThenRenewInCache :
//1. upsert.
//2. update cache if cache hit.
func (w *WorkerGrp) DoUpsertThenRenewInCache(ctx context.Context,
	upsertFn UpdateDataFn, k Hashed2Int, data interface{}) (interface{}, error) {
	return w.ws[w.locHash(k)].DoUpsertThenRenewInCache(ctx, upsertFn, k, data)
}

//MuxSize : get mux size
func (w *WorkerGrp) MuxSize() int {
	return w.muxSize
}

//DeepSize : get deep size
func (w *WorkerGrp) DeepSize() int {
	return w.deepSize
}

//Start : start all work go routine
func (w *WorkerGrp) Start() {
	for i := 0; i < w.muxSize; i++ {
		w.ws[i].Start()
	}
}

//Stop : stop
func (w *WorkerGrp) Stop() {
	w.stopOnce.Do(w.stop)
}

//WaitStop : wait stop
func (w *WorkerGrp) WaitStop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-w.exitSignal:
		return nil
	}
}

//stop all works
func (w *WorkerGrp) stop() {
	for i := 0; i < w.muxSize; i++ {
		w.ws[i].Stop()
	}
	//a go routine to wait all children done then signal it.
	go w.signalExit()
}

//signal all work go routine exit
func (w *WorkerGrp) signalExit() {
	w.wg.Wait()
	w.exitSignal <- struct{}{}
}

// localize hash to worker
func (w *WorkerGrp) locHash(k Hashed2Int) int {
	var hashNum = k.HashedInt()
	if hashNum < 0 {
		hashNum = -hashNum
	}
	hashNum %= w.muxSize
	return hashNum
}

//build work group build
func buildWorkGrp(opts ...Option) *WorkerGrp {
	var o = &_Option{
		muxSize:  DefaultMuxSize,
		deepSize: DefaultDeepSize,
	}
	for _, opt := range opts {
		opt(o)
	}
	var w = &WorkerGrp{}
	w.muxSize, w.deepSize = o.muxSize, o.deepSize

	w.wg = &sync.WaitGroup{}
	w.exitSignal = make(chan struct{}, 1)
	w.wg.Add(w.muxSize)

	w.ws = make([]*Worker, w.muxSize)
	return w
}
