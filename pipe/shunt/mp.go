package shunt

import (
	"context"
	"google.golang.org/protobuf/proto"
)

type OutPutMsg struct {
	msg proto.Message
	err error
}

type MsgProc struct {
	//context
	ctx context.Context
	//input msg
	inputMsg proto.Message
	//output chan
	outputChan chan OutPutMsg
	//slot index
	slotIndex int
	//function index
	fnIndex int
}

func NewMsgProc(ctx context.Context, slotIndex int, fnIndex int, inputMsg proto.Message) *MsgProc {
	return &MsgProc{
		ctx:        ctx,
		inputMsg:   inputMsg,
		outputChan: make(chan OutPutMsg, 1),
		slotIndex:  slotIndex,
		fnIndex:    fnIndex,
	}
}

func (m *MsgProc) SetOutput(outputMsg proto.Message, err error) {
	m.outputChan <- OutPutMsg{
		msg: outputMsg,
		err: err,
	}
}

func (m *MsgProc) GetOutput() (proto.Message, error) {
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	case ret := <-m.outputChan:
		return ret.msg, ret.err
	}
}
