package stcp

import (
	"fmt"
	"net"
	"sync"

	"github.com/pinealctx/neptune/syncx/pipe/q"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// QSender queue send connection sender
// send queue based connection sender
// user can put bytes to send queue async
// and LoopSend will send bytes in queue one by one
type QSender struct {
	// close once
	closeOnce sync.Once
	// connection
	conn net.Conn
	// meta info
	metaInfo atomic.Value
	// send queue -- actually the queue is bytes
	sendQ *q.Q[[]byte]
}

// NewQSendConnHandler : new queue send connection handler
// sendQSize : send queue size, if sendQSize is 0, the send queue has unlimited capacity
//
//	otherwise, the send queue has limited capacity.
//	if the send queue is full, Put2Queue will return error
func NewQSendConnHandler(conn net.Conn, sendQSize int) *QSender {
	h := &QSender{
		conn:  conn,
		sendQ: q.NewQ[[]byte](sendQSize),
	}
	h.metaInfo.Store(NewBasicMetaInfo(conn))
	return h
}

// Conn returns the underlying network connection (required, goroutine-safe)
func (x *QSender) Conn() net.Conn {
	return x.conn
}

// SetMetaInfo sets meta info for logging (required, goroutine-safe, re-entrant)
func (x *QSender) SetMetaInfo(m MetaInfo) {
	x.metaInfo.Store(m)
}

// MetaInfo gets meta info for logging (required, goroutine-safe)
func (x *QSender) MetaInfo() MetaInfo {
	v := x.metaInfo.Load()
	if v == nil {
		return nil
	}
	// nolint : forcetypeassert // I know the type is exactly here
	return v.(MetaInfo)
}

// Close closes connection handler (required, goroutine-safe, re-entrant)
// This method can be called directly via IConnSender interface to gracefully shutdown
// the connection and trigger the associated ConnHandler.Exit() through the goroutine defer chain
// Multiple calls are safe and will not panic
func (x *QSender) Close() error {
	var err error
	x.closeOnce.Do(func() {
		// close send queue - this will cause LoopSend() to exit
		x.sendQ.Close()
		// close connection - this will cause loopReceive() to exit
		err = x.conn.Close()
	})
	return err
}

// Put2Queue send bytes async(put to send queue) (optional, goroutine-safe, re-entrant)
func (x *QSender) Put2Queue(bs []byte) error {
	return x.sendQ.Push(bs)
}

// Put2SendMap put bytes to send map (optional, goroutine-safe, re-entrant)
// Note: QSender treats this as equivalent to Put2Queue, ignoring the key
func (x *QSender) Put2SendMap(_ uint32, bs []byte) error {
	return x.sendQ.Push(bs)
}

// Put2SendSMap put bytes to send map (optional, goroutine-safe, re-entrant)
// Note: QSender treats this as equivalent to Put2Queue, ignoring the key
func (x *QSender) Put2SendSMap(_ string, bs []byte) error {
	return x.sendQ.Push(bs)
}

// Put2SendMaps put multiple key uint32 and bytes pairs to send map (optional, goroutine-safe, re-entrant)
// Note: QSender does not support batch operations
func (x *QSender) Put2SendMaps(_ []KeyIntBytesPair) error {
	ulog.Error("QSender.Put2SendMaps.not.supported")
	return fmt.Errorf("QSender.Put2SendMaps.not.supported")
}

// Put2SendSMaps put multiple key string and bytes pairs to send map (optional, goroutine-safe, re-entrant)
// Note: QSender does not support batch operations
func (x *QSender) Put2SendSMaps(_ []KeyStrBytesPair) error {
	ulog.Error("QSender.Put2SendSMaps.not.supported")
	return fmt.Errorf("QSender.Put2SendSMaps.not.supported")
}

// loopSend is the internal sending loop (required, NOT goroutine-safe)
// WARNING: This method is ONLY called by ConnHandler internally.
// NEVER call this method from external code.
func (x *QSender) loopSend() {
	for {
		sendBytes, err := x.sendQ.Pop()
		if err != nil {
			// queue closed
			// nolint : forcetypeassert // I know the type is exactly here
			ulog.Info("QSender.sendQ.closed", zap.Error(err), zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
		err = SendBytes2Conn(x.conn, sendBytes)
		if err != nil {
			// nolint : forcetypeassert // I know the type is exactly here
			ulog.Info("QSender.send.failed", zap.Error(err), zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
	}
}
