package pipe

import (
	"github.com/pinealctx/neptune/pipe/q"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	//ErrNoHandler -- no message handler
	ErrNoHandler = status.Error(codes.Unimplemented, "no.req.handler")
	//ErrQueueFull -- msg queue is full
	ErrQueueFull = status.Error(codes.ResourceExhausted, "req.queue.full")
	//ErrQueueClosed -- msg queue closed
	ErrQueueClosed = status.Error(codes.Unavailable, "req.queue.closed")
	//ErrInvalidParam -- invalid msg param
	ErrInvalidParam = status.Error(codes.InvalidArgument, "invalid.req.param")
	//ErrInvalidRsp -- invalid result
	ErrInvalidRsp = status.Error(codes.Internal, "invalid.rsp.msg")
)

//ConvertQueueErr convert msg queue error
func ConvertQueueErr(err error) error {
	if err == nil {
		return nil
	}
	if err == q.ErrReqQFull {
		return ErrQueueFull
	}
	if err == q.ErrClosed {
		return ErrQueueClosed
	}
	return err
}

//NormalizeSlotIndex slot index
func NormalizeSlotIndex(index int, slotSize int) int {
	if index < 0 {
		index = -index
	}
	index %= slotSize
	return index
}
