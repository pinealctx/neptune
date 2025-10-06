package stcp

import (
	"net"
	"sync"

	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
)

// ConnHandler connection handler
type ConnHandler struct {
	// read processor
	readProcessor ReadProcessor
	// connection io interface
	iConnIO IConnIO
	//start once
	startOnce sync.Once
	// exit once
	exitOnce sync.Once
	// hook functions
	// start hook
	startHooks []ConnStartEvent
	// exit hook
	exitHooks []ConnExitEvent
}

// NewConnHandler : new connection handler
func NewConnHandler(readerProcessor ReadProcessor, iConnIO IConnIO) *ConnHandler {
	return &ConnHandler{
		readProcessor: readerProcessor,
		iConnIO:       iConnIO,
		startHooks:    make([]ConnStartEvent, 0, 1),
		exitHooks:     make([]ConnExitEvent, 0, 1),
	}
}

// GetIConn get connection interface
func (x *ConnHandler) GetIConn() IConnIO {
	return x.iConnIO
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
		go func() {
			defer func() {
				r := recover()
				if r != nil {
					ulog.Error("ConnHandler.loopReceive.recover", zap.Any("panic", r), zap.Object("metaInfo", x.iConnIO.MetaInfo()),
						zap.Stack("stack"))
				}
			}()
			defer x.Exit()

			conn := x.iConnIO.Conn()
			x.loopReceive(conn)
		}()

		go func() {
			defer func() {
				r := recover()
				if r != nil {
					ulog.Error("ConnHandler.loopSend.recover", zap.Any("panic", r), zap.Object("metaInfo", x.iConnIO.MetaInfo()),
						zap.Stack("stack"))
				}
			}()
			defer x.Exit()

			x.loopSend()
		}()

		for _, hook := range x.startHooks {
			func() {
				defer func() {
					r := recover()
					if r != nil {
						ulog.Error("ConnHandler.startHook.recover", zap.Any("panic", r), zap.Object("metaInfo", x.iConnIO.MetaInfo()),
							zap.Stack("stack"))
					}
				}()

				hook(x.iConnIO)
			}()
		}
	})
}

// Exit : exit connection handler
func (x *ConnHandler) Exit() {
	x.exitOnce.Do(func() {
		err := x.iConnIO.Close()
		if err != nil {
			ulog.Error("ConnHandler.iConnIO.Close", zap.Error(err), zap.Object("metaInfo", x.iConnIO.MetaInfo()))
		}
		for _, hook := range x.exitHooks {
			func() {
				defer func() {
					r := recover()
					if r != nil {
						ulog.Error("ConnHandler.exitHook.recover", zap.Any("panic", r), zap.Object("metaInfo", x.iConnIO.MetaInfo()),
							zap.Stack("stack"))
					}
				}()
				hook(x.iConnIO)
			}()
		}
	})
}

// loopReceive loop receive
// WARNING: This method is ONLY called by ConnHandler internally.
// NEVER call this method from external code - it will cause undefined behavior.
// This method runs in its own dedicated goroutine managed by ConnHandler.
func (x *ConnHandler) loopReceive(conn net.Conn) {
	for {
		buf, err := x.iConnIO.ReadFrame(conn)
		if err != nil {
			ulog.Info("ConnHandler.loopReceive.connReader", zap.Object("metaInfo", x.iConnIO.MetaInfo()), zap.Error(err))
			break
		}
		err = x.readProcessor(x.iConnIO, buf)
		if err != nil {
			ulog.Info("ConnHandler.loopReceive.readProcessor", zap.Object("metaInfo", x.iConnIO.MetaInfo()), zap.Error(err))
			break
		}
	}
}

// loopSend is the internal sending loop (required, NOT goroutine-safe)
// WARNING: This method is ONLY called by ConnHandler internally.
// NEVER call this method from external code.
func (x *ConnHandler) loopSend() {

	for {
		sendBytes, err := x.iConnIO.PopMsgBytes()
		if err != nil {
			// queue closed
			// nolint : forcetypeassert // I know the type is exactly here
			ulog.Info("ConnHandler.sendQ.closed", zap.Error(err), zap.Object("metaInfo", x.iConnIO.MetaInfo()))
			return
		}
		err = x.iConnIO.sendBytes2Conn(sendBytes)
		if err != nil {
			// nolint : forcetypeassert // I know the type is exactly here
			ulog.Info("ConnHandler.send.failed", zap.Error(err), zap.Object("metaInfo", x.iConnIO.MetaInfo()))
			return
		}
	}
}
