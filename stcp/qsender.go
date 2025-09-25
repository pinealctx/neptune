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

// Conn : get connection (required)
func (x *QSender) Conn() net.Conn {
	return x.conn
}

// LoopSend : loop to send (required)
func (x *QSender) LoopSend() {
	for {
		sendBytes, err := x.sendQ.Pop()
		if err != nil {
			// queue closed
			ulog.Info("QSender.sendQ.closed", zap.Error(err),
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
		err = SendBytes2Conn(x.conn, sendBytes)
		if err != nil {
			ulog.Info("QSender.send.failed", zap.Error(err),
				zap.Object("metaInfo", x.metaInfo.Load().(MetaInfo)))
			return
		}
	}
}

// SetMetaInfo set meta info (required)
func (x *QSender) SetMetaInfo(m MetaInfo) {
	x.metaInfo.Store(m)
}

// MetaInfo get meta info (required)
func (x *QSender) MetaInfo() MetaInfo {
	v := x.metaInfo.Load()
	if v == nil {
		return nil
	}
	return v.(MetaInfo)
}

// Close : close connection handler
// This method can be called directly via IConnSender interface to gracefully shutdown
// the connection and trigger the associated ConnHandler.Exit() through the goroutine defer chain
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

// Put2Queue send bytes async(put to send queue)
func (x *QSender) Put2Queue(bs []byte) error {
	return x.sendQ.Push(bs)
}

// Put2SendMap put bytes to send map (not supported)
func (x *QSender) Put2SendMap(_ uint32, bs []byte) error {
	return x.sendQ.Push(bs)
}

// Put2SendSMap put bytes to send map (not supported)
func (x *QSender) Put2SendSMap(_ string, bs []byte) error {
	return x.sendQ.Push(bs)
}

// Put2SendMaps put multiple key uint32 and bytes pairs to send map (not supported)
func (x *QSender) Put2SendMaps(_ []KeyIntBytesPair) error {
	ulog.Error("QSender.Put2SendMaps.not.supported")
	return fmt.Errorf("QSender.Put2SendMaps.not.supported")
}

// Put2SendSMaps put multiple key string and bytes pairs to send map (not supported)
func (x *QSender) Put2SendSMaps(_ []KeyStrBytesPair) error {
	ulog.Error("QSender.Put2SendSMaps.not.supported")
	return fmt.Errorf("QSender.Put2SendSMaps.not.supported")
}
