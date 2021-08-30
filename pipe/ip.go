package pipe

import "context"

type Rsp struct {
	//result
	r interface{}
	//error
	err error
}

//GenProc : generic proc
type GenProc struct {
	//context
	ctx context.Context
	//request
	req interface{}
	//response chan
	rspChan chan Rsp
	//slot index
	sIndex int
}

func NewGenProc(ctx context.Context, slotIndex int, req interface{}) *GenProc {
	return &GenProc{
		ctx:     ctx,
		req:     req,
		rspChan: make(chan Rsp, 1),
		sIndex:  slotIndex,
	}
}

func (m *GenProc) SetRsp(rsp interface{}, err error) {
	m.rspChan <- Rsp{
		r:   rsp,
		err: err,
	}
}

func (m *GenProc) Rsp() (interface{}, error) {
	select {
	case <-m.ctx.Done():
		return nil, m.ctx.Err()
	case r := <-m.rspChan:
		return r.r, r.err
	}
}
