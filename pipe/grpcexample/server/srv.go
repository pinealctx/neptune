package main

import (
	"context"
	"github.com/pinealctx/neptune/pipe"
	"github.com/pinealctx/neptune/pipe/grpcexample/pb"
	"google.golang.org/protobuf/proto"
)

const (
	IndexSayHello = 0
)

var (
	DefaultSlotTransfer ISlotTransfer
)

type ISlotTransfer interface {
	SlotSayHello(ctx context.Context, slotIndex int, req *pb.Halo) (*pb.Halo, error)
	SlotHalo(req *pb.Halo) int
}

type SrvRunner struct {
	shunt *pipe.Shunt
	pb.UnsafeHelloServiceServer
}

func NewSrvRunner(opts ...pipe.ShOption) *SrvRunner {
	var s = &SrvRunner{}
	s.Init()
	return s
}

func (s *SrvRunner) Init(opts ...pipe.ShOption) {
	s.shunt = pipe.NewShunt(opts...)

	/*Register msg handler*/
	var slotSayHelloFn = func(ctx context.Context, index int, req proto.Message) (output proto.Message, err error) {
		var msg, ok = req.(*pb.Halo)
		if !ok {
			return nil, pipe.ErrInvalidParam
		}
		return DefaultSlotTransfer.SlotSayHello(ctx, index, msg)
	}
	s.shunt.RegisterMsgHandler("SayHello", IndexSayHello, slotSayHelloFn)
	s.shunt.Run()
}

func (s *SrvRunner) SayHello(ctx context.Context, req *pb.Halo) (*pb.Halo, error) {
	var slot = DefaultSlotTransfer.SlotHalo(req)
	var msgProc, err = s.shunt.AddMsg(ctx, slot, IndexSayHello, req)
	if err != nil {
		return nil, err
	}
	var outMsg proto.Message
	outMsg, err = msgProc.GetOutput()
	if err != nil {
		return nil, err
	}
	var r, ok = outMsg.(*pb.Halo)
	if !ok {
		return nil, pipe.ErrInvalidResult
	}
	return r, nil
}
