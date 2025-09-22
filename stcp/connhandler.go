package stcp

import (
	"net"
	"sync"
	"time"

	"github.com/pinealctx/neptune/syncx/pipe/q"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ReadTimeout : Timeout for read/write
var (
	// write timeout
	writeTimeout *atomic.Duration
)

func init() {
	writeTimeout = atomic.NewDuration(time.Second * 5)
}

func SetWriteTimeout(d time.Duration) {
	writeTimeout.Store(d)
}

// MetaInfo meta info for logging
type MetaInfo interface {
	zapcore.ObjectMarshaler
}

// BasicMetaInfo basic meta info
type BasicMetaInfo struct {
	RemoteAddr string
}

func (m *BasicMetaInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("remoteAddr", m.RemoteAddr)
	return nil
}

// ConnStartEvent on connection start
type ConnStartEvent func(metaInfo MetaInfo)

// ConnExitEvent on connection exit
type ConnExitEvent func(conn net.Conn, metaInfo MetaInfo)

// IncomingHook incoming hook function
type IncomingHook func(handler *ConnHandler, conn net.Conn, metaInfo MetaInfo) error

// ConnHandler connection handler
type ConnHandler struct {
	// connection
	conn net.Conn
	// send queue -- actually the queue is bytes
	sendQ *q.Q[[]byte]
	//start once
	startOnce sync.Once
	//exit once
	exitOnce sync.Once

	// meta info
	metaInfo atomic.Value

	// hook functions
	// start hook
	startHooks []ConnStartEvent
	// exit hook
	exitHooks []ConnExitEvent
	// incoming hook
	incomingHook IncomingHook
}

// NewConnHandler : new connection handler
func NewConnHandler(conn net.Conn, qSize int, incomingHook IncomingHook) *ConnHandler {
	basicMeta := &BasicMetaInfo{
		RemoteAddr: conn.RemoteAddr().String(),
	}
	x := &ConnHandler{
		conn:         conn,
		sendQ:        q.NewQ[[]byte](qSize),
		incomingHook: incomingHook,
	}
	x.metaInfo.Store(basicMeta)
	return x
}

// NewConnHandlerV2 : new connection handler with exit hooks
func NewConnHandlerV2(conn net.Conn, qSize int, incomingHook IncomingHook, exitHook ConnExitEvent) *ConnHandler {
	x := NewConnHandler(conn, qSize, incomingHook)
	x.exitHooks = []ConnExitEvent{exitHook}
	return x
}

// NewConnHandlerV3 : new connection handler with start and exit hooks
func NewConnHandlerV3(conn net.Conn, qSize int, incomingHook IncomingHook,
	startHook ConnStartEvent, exitHook ConnExitEvent) *ConnHandler {
	x := NewConnHandler(conn, qSize, incomingHook)
	x.startHooks = []ConnStartEvent{startHook}
	x.exitHooks = []ConnExitEvent{exitHook}
	return x
}

// AddStartHook add start hook
func (x *ConnHandler) AddStartHook(hook ConnStartEvent) {
	x.startHooks = append(x.startHooks, hook)
}

// AddExitHook add exit hook
func (x *ConnHandler) AddExitHook(hook ConnExitEvent) {
	x.exitHooks = append(x.exitHooks, hook)
}

// Start : start connection handler
func (x *ConnHandler) Start() {
	x.startOnce.Do(func() {
		go x.loopReceive()
		go x.loopSend()
		for _, hook := range x.startHooks {
			hook(x.metaInfo.Load().(MetaInfo))
		}
	})
}

// SendAsync send bytes async
func (x *ConnHandler) SendAsync(bs []byte) error {
	return x.sendQ.Push(bs)
}

// SetMetaInfo set meta info
func (x *ConnHandler) SetMetaInfo(m MetaInfo) {
	x.metaInfo.Store(m)
}

// CLose : close connection
func (x *ConnHandler) CLose() {
	x.quit()
}

// loop receive
func (x *ConnHandler) loopReceive() {
	defer func() {
		r := recover()
		if r != nil {
			ulog.Error("ConnHandler.panic.in.loopReceive",
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)),
				zap.Any("panic", r), zap.Stack("stack"))
		}
	}()
	defer x.quit()

	for {
		err := x.incomingHook(x, x.conn, x.metaInfo.Load().(MetaInfo))
		if err != nil {
			ulog.Info("ConnHandler.incoming.hook.failed", zap.Error(err),
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
	}
}

// loopSend loop to send bytes
func (x *ConnHandler) loopSend() {
	defer func() {
		r := recover()
		if r != nil {
			ulog.Error("ConnHandler.panic.in.loopSend",
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)),
				zap.Any("panic", r), zap.Stack("stack"))
		}
	}()
	defer x.quit()

	for {
		sendBytes, err := x.sendQ.Pop()
		if err != nil {
			// queue closed
			ulog.Info("ConnHandler.sendQ.closed", zap.Error(err),
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
		err = x.send(sendBytes, writeTimeout.Load())
		if err != nil {
			ulog.Info("ConnHandler.send.failed", zap.Error(err),
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
	}
}

// quit : quit connection
func (x *ConnHandler) quit() {
	x.exitOnce.Do(func() {
		x.sendQ.Close()
		err := x.conn.Close()
		metaInfo := x.metaInfo.Load().(MetaInfo)
		if err != nil {
			ulog.Error("ConnHandler.close.conn.error", zap.Object("metaInfo", metaInfo),
				zap.Error(err))
		}
		for _, hook := range x.exitHooks {
			hook(x.conn, metaInfo)
		}
	})
}

// send bytes
func (x *ConnHandler) send(bs []byte, timeout time.Duration) error {
	err := x.conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		ulog.Info("ConnHandler.set.write.conn.deadline",
			zap.Error(err), zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
		return err
	}
	_, err = x.conn.Write(bs)
	return err
}
