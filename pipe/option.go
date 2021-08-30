package pipe

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	//DefaultSlotSize slot size
	//素数
	DefaultSlotSize = 509
	//DefaultQInSlotSize default pipe size, total current request can be pipe is 509*1024*8 = 4169728
	DefaultQInSlotSize = 1024 * 8
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

//shunt option
type _ShOption struct {
	//slot size
	slotSize int
	//queue size in each slot
	qSizeInSlot int
}

//ShOption shunt option function
type ShOption func(o *_ShOption)

//WithSlotSize setup slot size
func WithSlotSize(slotSize int) ShOption {
	return func(o *_ShOption) {
		o.slotSize = slotSize
	}
}

//WithQSizeInSlot setup queue size in each slot
func WithQSizeInSlot(qSizeInSlot int) ShOption {
	return func(o *_ShOption) {
		o.qSizeInSlot = qSizeInSlot
	}
}
