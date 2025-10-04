package stcp

import (
	"net"
	"sync"

	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
)

// ConnHandler connection handler
type ConnHandler[T any] struct {
	// read processor
	readProcessor ReadProcessor[T]
	// connection io interface
	iConnIO IConnIO[T]
	//start once
	startOnce sync.Once
	// exit once
	exitOnce sync.Once
	// hook functions
	// start hook
	startHooks []ConnStartEvent[T]
	// exit hook
	exitHooks []ConnExitEvent[T]
}

// NewConnHandler : new connection handler
func NewConnHandler[T any](readerProcessor ReadProcessor[T], iConnIO IConnIO[T]) *ConnHandler[T] {
	return &ConnHandler[T]{
		readProcessor: readerProcessor,
		iConnIO:       iConnIO,
		startHooks:    make([]ConnStartEvent[T], 0, 1),
		exitHooks:     make([]ConnExitEvent[T], 0, 1),
	}
}

// GetIConn get connection interface
func (x *ConnHandler[T]) GetIConn() IConnIO[T] {
	return x.iConnIO
}

// AddStartHook add start hook
func (x *ConnHandler[T]) AddStartHook(hook ConnStartEvent[T]) {
	x.startHooks = append(x.startHooks, hook)
}

// AddExitHook add exit hook
func (x *ConnHandler[T]) AddExitHook(hook ConnExitEvent[T]) {
	x.exitHooks = append(x.exitHooks, hook)
}

// Start : start connection handler
func (x *ConnHandler[T]) Start() {
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

			x.iConnIO.loopSend()
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
func (x *ConnHandler[T]) Exit() {
	x.exitOnce.Do(func() {
		err := x.iConnIO.Close()
		if err != nil {
			ulog.Error("x.handler.Close", zap.Error(err), zap.Object("metaInfo", x.iConnIO.MetaInfo()))
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
func (x *ConnHandler[T]) loopReceive(conn net.Conn) {
	for {
		buf, err := x.iConnIO.ReadFrame(conn)
		if err != nil {
			ulog.Info("loopReceive.connReader", zap.Object("metaInfo", x.iConnIO.MetaInfo()), zap.Error(err))
			break
		}
		err = x.readProcessor(x.iConnIO, buf)
		if err != nil {
			ulog.Info("loopReceive.readProcessor", zap.Object("metaInfo", x.iConnIO.MetaInfo()), zap.Error(err))
			break
		}
	}
}
