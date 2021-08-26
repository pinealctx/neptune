package etcd

import (
	"context"
	"sync"
	"time"

	"github.com/pinealctx/neptune/ulog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

const (
	//ChanBufferSize chan buffer size
	ChanBufferSize = 10
)

type Watcher struct {
	//client
	cli *Client

	//context control
	ctx    context.Context
	cancel context.CancelFunc

	//run once
	startOnce sync.Once
	closeOnce sync.Once

	//watch chan
	dirChan chan DirEvent

	//path
	path string

	//time config
	timeout time.Duration
	backoff time.Duration
}

//NewWatcher new watcher
func NewWatcher(cli *Client, path string, timeout, backoff time.Duration) *Watcher {
	var ctx, cancel = context.WithCancel(context.Background())
	return &Watcher{
		ctx:    ctx,
		cancel: cancel,
		cli:    cli,

		dirChan: make(chan DirEvent, ChanBufferSize),

		path:    path,
		timeout: timeout,
		backoff: backoff,
	}
}

//DirChan dir chan
func (w *Watcher) DirChan() <-chan DirEvent {
	return w.dirChan
}

//StartWatchDir loop watch dir -- if chan closed, which means go routine out
func (w *Watcher) StartWatchDir() {
	w.startOnce.Do(func() {
		go w.startWatchDir(true)
	})
}

//StartWatchDirWhenExist loop watch dir -- if chan closed, which means go routine out
//only in case there is any key under the watch path
func (w *Watcher) StartWatchDirWhenExist() {
	w.startOnce.Do(func() {
		go w.startWatchDir(false)
	})
}

//Stop stop loop watch
func (w *Watcher) Stop() {
	w.closeOnce.Do(func() {
		var ok bool
		w.cancel()
		for {
			//exhaust all channel
			_, ok = <-w.dirChan
			if !ok {
				return
			}
		}
	})
}

//start watch dir
func (w *Watcher) startWatchDir(ignoreEmpty bool) {
	var (
		dRet   DirRet
		cancel context.CancelFunc
		inited bool

		//watch chan
		iChan <-chan DirEvent
	)

	//close channel finally
	defer close(w.dirChan)

	for {
		if !inited {
			inited = true
		} else {
			//select ctx done
			select {
			case <-w.ctx.Done():
				return
			default:
			}
			time.Sleep(w.backoff)
		}

		//check cancel,if not nil,call cancel in case of re-setup watch
		if cancel != nil {
			cancel()
		}

		//setup watch
		dRet, iChan, cancel = w.cli.WatchDir(w.ctx, w.path, w.timeout, ignoreEmpty)

		if !w.watchDirMoveOn(dRet, ignoreEmpty) {
			ulog.Error("cli.watch.dir.setup.error", zap.String("path", w.path), zap.Error(dRet.Err))
			continue
		}

		if iChan == nil {
			ulog.Error("cli.watch.dir.setup.error.watchKv.nil", zap.String("path", w.path))
			continue
		}

		if cancel == nil {
			//never gonna happen
			ulog.Error("cli.watch.dir.setup.error.cancel.nil", zap.String("path", w.path))
			continue
		}

		if w.procFirstRet(dRet) {
			return
		}

		if w.procSession(iChan) {
			return
		}
	}
}

//loop dir session
//if ctx.Done, return true, means watch is done.
//if return false, just quit current session, continue to next session
func (w *Watcher) procSession(iChan <-chan DirEvent) bool {
	var (
		dEvent DirEvent
		ok     bool
	)
	for {
		select {
		case dEvent, ok = <-iChan:
			if !ok {
				//break current for loop, re-setup watch
				ulog.Error("cli.re-setup.watch.dir", zap.String("path", w.path))
				return false
			}
			if dEvent.Err != nil {
				//error occur, record it
				ulog.Error("cli.watch.dir.chan.error", zap.String("path", w.path), zap.Error(dEvent.Err))
				continue
			}
			w.dirChan <- dEvent
		case <-w.ctx.Done():
			return true
		}
	}

}

//proc first dir result
//if dRet.Err is not nil, means get dir is empty, just return to next to watch.
//if dRet.Err is nil, put the result to channel
//if ctx.Done, return true, means watch is done.
//if return false, just quit current session, continue to next session
func (w *Watcher) procFirstRet(dRet DirRet) bool {
	if dRet.Err != nil {
		if !IsNotFoundErr(dRet.Err) {
			ulog.Error("cli.watch.dir.setup.crazy.error", zap.String("path", w.path), zap.Error(dRet.Err))
		}
		ulog.Info("cli.watch.dir.setup.error.no.children", zap.String("path", w.path), zap.Error(dRet.Err))
		select {
		case <-w.ctx.Done():
			return true
		default:
		}
		return false
	}

	select {
	//put first get
	case w.dirChan <- first2Event(dRet):
	case <-w.ctx.Done():
		return true
	}
	return false
}

//handle cli.WatchDir error
//return
//true  --- move on
//false --- restart from the loop(continue)
func (w *Watcher) watchDirMoveOn(dRet DirRet, ignoreEmpty bool) bool {
	if dRet.Err != nil {
		if !ignoreEmpty || !IsNotFoundErr(dRet.Err) {
			return false
		}
	}
	return true
}

//first get dir to event
func first2Event(dirRet DirRet) DirEvent {
	var event = DirEvent{
		Err:      dirRet.Err,
		Revision: dirRet.Revision,
	}
	var size = len(dirRet.KVS)
	if size == 0 {
		return event
	}
	event.Events = make([]*clientv3.Event, size)
	for i := 0; i < size; i++ {
		event.Events[i] = &clientv3.Event{
			Type: mvccpb.PUT,
			Kv: &mvccpb.KeyValue{
				Key:   dirRet.KVS[i].Key,
				Value: dirRet.KVS[i].Value,
			},
		}
	}
	return event
}
