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

// CommonIncomingHook incoming hook function
type CommonIncomingHook func(handler *CommonConnHandler, conn net.Conn, metaInfo MetaInfo) error

// ConnHandlerRunner connection handler runner
type ConnHandlerRunner struct {
	IConnHandler
	//start once
	startOnce sync.Once
}

// NewConnHandlerRunner : new connection handler runner
func NewConnHandlerRunner(handler IConnHandler) *ConnHandlerRunner {
	return &ConnHandlerRunner{
		IConnHandler: handler,
	}
}

// Start : start connection handler
func (x *ConnHandlerRunner) Start() {
	x.startOnce.Do(func() {
		go x.LoopReceive()
		go x.LoopSend()
		startHooks := x.StartHooks()
		metaInfo := x.MetaInfo()
		for _, hook := range startHooks {
			hook(metaInfo)
		}
	})
}

// IConnHandler connection handler interface
type IConnHandler interface {
	// AddStartHook add start hook
	AddStartHook(hook ConnStartEvent)
	// AddExitHook add exit hook
	AddExitHook(hook ConnExitEvent)
	// LoopReceive : loop receive
	LoopReceive()
	// LoopSend : loop to send bytes
	LoopSend()
	// SetMetaInfo set meta info
	SetMetaInfo(m MetaInfo)
	// MetaInfo get meta info
	MetaInfo() MetaInfo
	// StartHooks get start hooks
	StartHooks() []ConnStartEvent
	// ExitHooks get exit hooks
	ExitHooks() []ConnExitEvent
}

// BasicConnHandler basic connection handler
type BasicConnHandler struct {
	// connection
	Conn net.Conn
	// meta info
	metaInfo atomic.Value
	// hook functions
	// start hook
	startHooks []ConnStartEvent
	// exit hook
	exitHooks []ConnExitEvent
}

// NewBasicConnHandler : new basic connection handler
func NewBasicConnHandler(conn net.Conn) *BasicConnHandler {
	basicMeta := &BasicMetaInfo{
		RemoteAddr: conn.RemoteAddr().String(),
	}
	x := &BasicConnHandler{
		Conn: conn,
	}
	x.metaInfo.Store(basicMeta)
	x.startHooks = make([]ConnStartEvent, 0, 1)
	x.exitHooks = make([]ConnExitEvent, 0, 1)
	return x
}

// SetMetaInfo set meta info
func (x *BasicConnHandler) SetMetaInfo(m MetaInfo) {
	x.metaInfo.Store(m)
}

// MetaInfo get meta info
func (x *BasicConnHandler) MetaInfo() MetaInfo {
	return x.metaInfo.Load().(MetaInfo)
}

// AddStartHook add start hook
func (x *BasicConnHandler) AddStartHook(hook ConnStartEvent) {
	x.startHooks = append(x.startHooks, hook)
}

// AddExitHook add exit hook
func (x *BasicConnHandler) AddExitHook(hook ConnExitEvent) {
	x.exitHooks = append(x.exitHooks, hook)
}

// StartHooks get start hooks
func (x *BasicConnHandler) StartHooks() []ConnStartEvent {
	return x.startHooks
}

// ExitHooks get exit hooks
func (x *BasicConnHandler) ExitHooks() []ConnExitEvent {
	return x.exitHooks
}

// Send bytes to connection with timeout
func (x *BasicConnHandler) Send(bs []byte, timeout time.Duration) error {
	err := x.Conn.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		ulog.Info("BasicConnHandler.set.write.conn.deadline",
			zap.Error(err), zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
		return err
	}
	_, err = x.Conn.Write(bs)
	return err
}

// CommonConnHandler connection handler
type CommonConnHandler struct {
	*BasicConnHandler
	// send queue -- actually the queue is bytes
	sendQ *q.Q[[]byte]
	//exit once
	exitOnce sync.Once
	// incoming hook
	incomingHook CommonIncomingHook
}

// NewCommonConnHandler : new connection handler
func NewCommonConnHandler(conn net.Conn, qSize int, incomingHook CommonIncomingHook) *CommonConnHandler {
	basicHandler := NewBasicConnHandler(conn)
	x := &CommonConnHandler{
		BasicConnHandler: basicHandler,
		sendQ:            q.NewQ[[]byte](qSize),
		incomingHook:     incomingHook,
	}
	return x
}

// SendAsync send bytes async
func (x *CommonConnHandler) SendAsync(bs []byte) error {
	return x.sendQ.Push(bs)
}

// CLose : close connection
func (x *CommonConnHandler) CLose() {
	x.quit()
}

// LoopReceive : loop receive
func (x *CommonConnHandler) LoopReceive() {
	defer func() {
		r := recover()
		if r != nil {
			ulog.Error("ConnHandler.panic.in.LoopReceive",
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)),
				zap.Any("panic", r), zap.Stack("stack"))
		}
	}()
	defer x.quit()

	for {
		err := x.incomingHook(x, x.Conn, x.metaInfo.Load().(MetaInfo))
		if err != nil {
			ulog.Info("ConnHandler.incoming.hook.failed", zap.Error(err),
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
	}
}

// LoopSend : loop to send bytes
func (x *CommonConnHandler) LoopSend() {
	defer func() {
		r := recover()
		if r != nil {
			ulog.Error("ConnHandler.panic.in.LoopSend",
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
		err = x.Send(sendBytes, writeTimeout.Load())
		if err != nil {
			ulog.Info("ConnHandler.send.failed", zap.Error(err),
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
	}
}

// quit : quit connection
func (x *CommonConnHandler) quit() {
	x.exitOnce.Do(func() {
		x.sendQ.Close()
		err := x.Conn.Close()
		metaInfo := x.metaInfo.Load().(MetaInfo)
		if err != nil {
			ulog.Error("ConnHandler.close.conn.error", zap.Object("metaInfo", metaInfo),
				zap.Error(err))
		}
		for _, hook := range x.exitHooks {
			hook(x.Conn, metaInfo)
		}
	})
}
