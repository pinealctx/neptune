package etcd

import (
	"context"
	"fmt"
	"path"
	"strings"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
)

// DirEvent dir event
type DirEvent struct {
	Events   []*clientv3.Event
	Err      error
	Revision int64
}

// AddEvent add event
func (w *DirEvent) AddEvent(e *clientv3.Event) {
	w.Events = append(w.Events, e)
}

// DebugInfo debug info
func (w *DirEvent) DebugInfo() string {
	var buffer strings.Builder
	var head = fmt.Sprintf("len event:%d, error:%+v, revision:%+v",
		len(w.Events), w.Err, w.Revision)
	buffer.WriteString(head)
	buffer.WriteString(":")
	for _, event := range w.Events {
		buffer.WriteString("type:")
		buffer.WriteString(event.Type.String())
		buffer.WriteString("key:")
		buffer.WriteString(string(event.Kv.Key))
		buffer.WriteString(",")
		buffer.WriteString("val:")
		buffer.WriteString(string(event.Kv.Value))
	}
	return buffer.String()
}

// WatchDir watch dir
// ctx -- cancel controller, it's a parent context for watch
// nodePath -- path
// timeout -- get dir timeout
// ignoreEmpty -- if ignore, even get nothing from dir, watch still continue to work
// return first get dir content/dir watch event/current watch session cancel
func (c *Client) WatchDir(ctx context.Context,
	nodePath string, timeout time.Duration, ignoreEmpty bool) (DirRet, <-chan DirEvent, context.CancelFunc) {

	//get dir
	var gCtx, cancel = context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var dRet = c.GetDir(gCtx, nodePath, false)
	if dRet.Err != nil {
		// even node is empty, it's can still watch
		var canIgnore = ignoreEmpty && IsNotFoundErr(dRet.Err)
		if !canIgnore {
			return dRet, nil, nil
		}
	}

	var watchPath = path.Join(c.root, nodePath) + "/"
	var outCtx, outCancel = context.WithCancel(ctx)
	var watcher = c.eCli.Watch(clientv3.WithRequireLeader(outCtx),
		watchPath, clientv3.WithPrefix())
	if watcher == nil {
		dRet.Err = genErr(watchPath, WatchFail)
		return dRet, nil, outCancel
	}

	// Create the notifications channel, send updates to it.
	var nChan = make(chan DirEvent, 1)
	go func() {
		// loop receive notify
		dirLoop(outCtx, nChan, watcher, nodePath)
	}()
	return dRet, nChan, outCancel
}

// dir watch loop
func dirLoop(ctx context.Context,
	notifyChan chan<- DirEvent, watcher clientv3.WatchChan, nodePath string) {
	var err error

	defer close(notifyChan)

	for {
		select {
		case <-ctx.Done():
			// This includes context cancellation errors.
			notifyChan <- DirEvent{Err: convertErr(nodePath, ctx.Err())}
			return
		case wRsp, ok := <-watcher:
			if !ok {
				notifyChan <- DirEvent{Err: genErr(nodePath, WatchClosed)}
				return
			}
			if wRsp.Canceled {
				// Final notification.
				err = wRsp.Err()
				notifyChan <- DirEvent{Err: convertErr(nodePath, convertErr(nodePath, err))}
				return
			}

			err = wRsp.Err()
			if err != nil {
				//watch error
				notifyChan <- DirEvent{Err: convertErr(nodePath, err)}
			}

			if len(wRsp.Events) == 0 {
				notifyChan <- DirEvent{Err: genErr(nodePath, WatchUnexpected)}
				return
			}

			if procDirEvents(wRsp, nodePath, notifyChan) {
				return
			}
		}

		//loop for again
	}
}

// proc dir events, if not put and delete, need exist current watch
// return true -- something wrong, need exit current session
// return false -- done current event work, move on
func procDirEvents(wRsp clientv3.WatchResponse, nodePath string, notifyChan chan<- DirEvent) bool {
	var watchDir = DirEvent{}
	watchDir.Revision = wRsp.Header.Revision

	for _, ev := range wRsp.Events {
		switch ev.Type {
		case mvccpb.PUT, mvccpb.DELETE:
			watchDir.AddEvent(ev)
			notifyChan <- watchDir
		default:
			watchDir.Err = genErr(nodePath, WatchUnexpected)
			notifyChan <- watchDir
			return true
		}
	}
	return false
}
