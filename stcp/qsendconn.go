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

// QSendConn queue send connection sender
// send queue based connection sender
// user can put bytes to send queue async
// and LoopSend will send bytes in queue one by one
type QSendConn[T any] struct {
	// close once
	closeOnce sync.Once
	// connection
	conn net.Conn
	// meta info
	metaInfo atomic.Value
	// send queue -- actually the queue is bytes
	sendQ *q.Q[[]byte]
	// reader
	reader IConnReader
}

// NewQSendConnHandler : new queue send connection handler
// sendQSize : send queue size, if sendQSize is 0, the send queue has unlimited capacity
//
//	otherwise, the send queue has limited capacity.
//	if the send queue is full, Put2Queue will return error
func NewQSendConnHandler[T any](conn net.Conn, sendQSize int, readerFactory ConnReaderFactory) *QSendConn[T] {
	h := &QSendConn[T]{
		conn:   conn,
		sendQ:  q.NewQ[[]byte](sendQSize),
		reader: readerFactory(conn),
	}
	h.metaInfo.Store(NewBasicMetaInfo(conn))
	return h
}

// NewQSendConnHandlerWithReader : new queue send connection handler with custom reader
// sendQSize : send queue size, if sendQSize is 0, the send queue has unlimited capacity
//
//	otherwise, the send queue has limited capacity.
//	if the send queue is full, Put2Queue will return error
func NewQSendConnHandlerWithReader[T any](conn net.Conn, sendQSize int, reader IConnReader) *QSendConn[T] {
	h := &QSendConn[T]{
		conn:   conn,
		sendQ:  q.NewQ[[]byte](sendQSize),
		reader: reader,
	}
	h.metaInfo.Store(NewBasicMetaInfo(conn))
	return h
}

// Conn returns the underlying network connection (required, goroutine-safe)
func (x *QSendConn[T]) Conn() net.Conn {
	return x.conn
}

// SetMetaInfo sets meta info for logging (required, goroutine-safe, re-entrant)
func (x *QSendConn[T]) SetMetaInfo(m MetaInfo) {
	x.metaInfo.Store(m)
}

// MetaInfo gets meta info for logging (required, goroutine-safe)
func (x *QSendConn[T]) MetaInfo() MetaInfo {
	v := x.metaInfo.Load()
	if v == nil {
		return nil
	}
	// nolint : forcetypeassert // I know the type is exactly here
	return v.(MetaInfo)
}

// Close closes connection handler (required, goroutine-safe, re-entrant)
// This method can be called directly via IConnSender/IConnIO interface to gracefully shutdown
// the connection and trigger the associated ConnHandler.Exit() through the goroutine defer chain
// Multiple calls are safe and will not panic
func (x *QSendConn[T]) Close() error {
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
func (x *QSendConn[T]) Put2Queue(bs []byte) error {
	return x.sendQ.Push(bs)
}

// Put2SendMap put bytes to send map (optional, goroutine-safe, re-entrant)
// Note: treats this as equivalent to Put2Queue, ignoring the key
func (x *QSendConn[T]) Put2SendMap(_ T, bs []byte) error {
	return x.sendQ.Push(bs)
}

// Put2SendMaps put multiple key uint32 and bytes pairs to send map (optional, goroutine-safe, re-entrant)
// Note: does not support batch operations
func (x *QSendConn[T]) Put2SendMaps(_ []KvItem[T]) error {
	ulog.Error("QSendConn.Put2SendMaps.not.supported")
	return fmt.Errorf("QSendConn.Put2SendMaps.not.supported")
}

// ReadFrame reads one frame(an entire message bytes) from connection
func (x *QSendConn[T]) ReadFrame(conn net.Conn) ([]byte, error) {
	buf, err := x.reader.ReadFrame(conn)
	if err != nil {
		return nil, fmt.Errorf("QSendConn.ReadFrame: %w", err)
	}
	return buf, nil
}

// loopSend is the internal sending loop (required, NOT goroutine-safe)
// WARNING: This method is ONLY called by ConnHandler internally.
// NEVER call this method from external code.
// nolint:unused // This method implements IConnSender[T].loopSend() interface
func (x *QSendConn[T]) loopSend() {
	for {
		sendBytes, err := x.sendQ.Pop()
		if err != nil {
			// queue closed
			// nolint : forcetypeassert // I know the type is exactly here
			ulog.Info("QSendConn.sendQ.closed", zap.Error(err), zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
		err = SendBytes2Conn(x.conn, sendBytes)
		if err != nil {
			// nolint : forcetypeassert // I know the type is exactly here
			ulog.Info("QSendConn.send.failed", zap.Error(err), zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
	}
}
