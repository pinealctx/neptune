package main

import (
	"context"
	"fmt"
	"github.com/pinealctx/neptune/pipe/grpcexample/pb"
	"hash/crc32"
)

type Transfer struct {
}

func (t Transfer) SlotSayHello(ctx context.Context, slotIndex int, req *pb.Halo) (*pb.Halo, error) {
	var r = &pb.Halo{
		Msg: fmt.Sprintf("echo:%s", req.Msg),
	}
	return r, nil
}

func (t Transfer) SlotHalo(req *pb.Halo) int {
	return int(crc32.ChecksumIEEE([]byte(req.Msg)))
}

func init() {
	DefaultSlotTransfer = Transfer{}
}
