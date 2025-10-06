package stcp

import (
	"fmt"
	"net"

	"github.com/pinealctx/neptune/syncx/pipe/q"
)

// QSendConn queue send connection sender
// send queue based connection sender
// user can put bytes to send queue async then send bytes in queue one by one
type QSendConn struct {
	*BasicConnIO
	// send queue -- actually the queue is bytes
	sendQ *q.Q[[]byte]
}

// NewQSendConnHandler : new queue send connection handler
// sendQSize : send queue size, if sendQSize is 0, the send queue has unlimited capacity
//
//	otherwise, the send queue has limited capacity.
//	if the send queue is full, Put2Queue will return error
func NewQSendConnHandler(conn net.Conn, sendQSize int, reader IConnReader) *QSendConn {
	h := &QSendConn{
		BasicConnIO: NewBasicConnIO(conn, reader),
		sendQ:       q.NewQ[[]byte](sendQSize),
	}
	h.metaInfo.Store(NewBasicMetaInfo(conn))
	return h
}

// Close closes connection handler (required, goroutine-safe, re-entrant)
// This method can be called directly via IConnSender/IConnIO interface to gracefully shutdown
// the connection and trigger the associated ConnHandler.Exit() through the goroutine defer chain
// Multiple calls are safe and will not panic
func (x *QSendConn) Close() error {
	return x.CloseX(func() error {
		// close send queue - this will cause LoopSend() to exit
		x.sendQ.Close()
		// close connection - this will cause loopReceive() to exit
		return x.conn.Close()
	})
}

// PutMsg send bytes async(put to send queue) (optional, goroutine-safe, re-entrant)
func (x *QSendConn) PutMsg(msg IMsg) error {
	if msg == nil {
		return fmt.Errorf("QSendConn.PutMsg: nil message")
	}
	bsMsg, ok := msg.(BytesMsg)
	if !ok {
		return fmt.Errorf("QSendConn.PutMsg: invalid message name: %s, expected BytesMsg", msg.Name())
	}
	if len(bsMsg.Bs) == 0 {
		return fmt.Errorf("QSendConn.PutMsg: empty BytesMsg")
	}
	return x.sendQ.Push(bsMsg.Bs)
}

// PopMsgBytes pop message bytes to send (optional, goroutine-safe, re-entrant)
func (x *QSendConn) PopMsgBytes() ([]byte, error) {
	return x.sendQ.Pop()
}
