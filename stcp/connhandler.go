package stcp

import (
	"net"
	"sync"

	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
)

// ConnHandler connection handler
type ConnHandler struct {
	// connection reader(function)
	connReader ConnReaderFunc
	// connection sender(handler)
	connSender IConnSender
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

// NewConnRunner : new connection runner
func NewConnRunner(connReader ConnReaderFunc, connSender IConnSender) *ConnHandler {
	return &ConnHandler{
		connReader: connReader,
		connSender: connSender,
		startHooks: make([]ConnStartEvent, 0, 1),
		exitHooks:  make([]ConnExitEvent, 0, 1),
	}
}

// GetConnSender get connection sender
func (x *ConnHandler) GetConnSender() IConnSender {
	return x.connSender
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
					ulog.Error("loopReceive.recover", zap.Any("panic", r), zap.Object("metaInfo", x.connSender.MetaInfo()),
						zap.Stack("stack"))
				}
			}()
			defer x.Exit()
			conn := x.connSender.Conn()
			x.loopReceive(conn)
		}()

		go func() {
			defer func() {
				r := recover()
				if r != nil {
					ulog.Error("x.connSender.LoopSend.recover", zap.Any("panic", r), zap.Object("metaInfo", x.connSender.MetaInfo()),
						zap.Stack("stack"))
				}
			}()
			defer x.Exit()

			x.connSender.LoopSend()
		}()
		for _, hook := range x.startHooks {
			hook(x.connSender)
		}
	})
}

// Exit : exit connection handler
func (x *ConnHandler) Exit() {
	x.exitOnce.Do(func() {
		err := x.connSender.Close()
		if err != nil {
			ulog.Error("x.connSender.Close", zap.Error(err), zap.Object("metaInfo", x.connSender.MetaInfo()))
		}
		for _, hook := range x.exitHooks {
			hook(x.connSender)
		}
	})
}

// loopReceive loop receive
// x.connReader is the function to read from connection
// when x.connReader returns error, loopReceive will exit
func (x *ConnHandler) loopReceive(conn net.Conn) {
	for {
		err := x.connReader(x.connSender, conn)
		if err != nil {
			ulog.Info("loopReceive.connReader", zap.Object("metaInfo", x.connSender.MetaInfo()), zap.Error(err))
			break
		}
	}
}
