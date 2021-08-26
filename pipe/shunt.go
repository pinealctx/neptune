package pipe

import (
	"context"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"sync"
)

const (
	//DefaultSlotSize slot size
	//素数
	DefaultSlotSize = 509
	//DefaultQInSlotSize default pipe size, total current request can be pipe is 509*1024*8 = 4169728
	DefaultQInSlotSize = 1024*8
)

var (
	//ErrNoMsgHandler -- no message handler
	ErrNoMsgHandler = status.Error(codes.Unimplemented, "no.msg.handler")
	//ErrMsgQueueFull -- msg queue is full
	ErrMsgQueueFull = status.Error(codes.ResourceExhausted, "msg.queue.full")
	//ErrMsgQueueClosed -- msg queue closed
	ErrMsgQueueClosed = status.Error(codes.Unavailable, "msg.queue.closed")
	//ErrInvalidParam -- invalid msg param
	ErrInvalidParam = status.Error(codes.InvalidArgument, "invalid.input.msg")
	//ErrInvalidResult -- invalid result
	ErrInvalidResult = status.Error(codes.Internal, "invalid.output.msg")
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

//ProcMsgFn function
type ProcMsgFn func(ctx context.Context, index int, inputMsg proto.Message) (outputMsg proto.Message, err error)

//ProcMsgInfo -- name and function
type ProcMsgInfo struct {
	fn   ProcMsgFn
	name string
}

//MsgHandler : msg handler
type MsgHandler struct {
	//function map
	//key -- function index
	//value -- message handler
	fnMap map[int]*ProcMsgInfo
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		fnMap: make(map[int]*ProcMsgInfo),
	}
}

//RegisterFn -- register function
func (h *MsgHandler) RegisterFn(name string, fnIndex int, fn ProcMsgFn) {
	h.fnMap[fnIndex] = &ProcMsgInfo{
		fn:   fn,
		name: name,
	}
}

//GetFn -- get function by fn index
func (h *MsgHandler) GetFn(fnIndex int) ProcMsgFn {
	var fnInfo, ok = h.fnMap[fnIndex]
	if !ok {
		return nil
	}
	return fnInfo.fn
}

//GetFunctionName -- get function name by fn index
func (h *MsgHandler) GetFunctionName(fnIndex int) string {
	var fn, ok = h.fnMap[fnIndex]
	if !ok {
		return ""
	}
	return fn.name
}

//Shunt multi-queue: like shunt
type Shunt struct {
	//slot size
	slotSize int
	//queue size in each slot
	qSizeInSlot int

	//multi queues
	qs []*Q

	//msg handler
	msgHandler *MsgHandler

	//wait group
	wg *sync.WaitGroup

	//go routine exit chan
	exitChan chan struct{}
	//stop once
	stopOnce sync.Once
}

//NewShunt : new shunt
func NewShunt(opts ...ShOption) *Shunt {
	//option
	var o = &_ShOption{
		slotSize:    DefaultSlotSize,
		qSizeInSlot: DefaultQInSlotSize,
	}
	for _, opt := range opts {
		opt(o)
	}
	//new shunt then init
	return NewShuntWithSize(o.slotSize, o.qSizeInSlot)
}

//NewShuntWithSize : new shunt with size
func NewShuntWithSize(slotSize int, qSizeInSlot int) *Shunt {
	var shunt = &Shunt{}
	shunt.slotSize, shunt.qSizeInSlot = slotSize, qSizeInSlot
	shunt.wg = &sync.WaitGroup{}
	shunt.exitChan = make(chan struct{}, 1)
	shunt.wg.Add(shunt.slotSize)

	shunt.qs = make([]*Q, shunt.slotSize)
	for i := 0; i < shunt.slotSize; i++ {
		shunt.qs[i] = NewQ(WithQReqSize(shunt.qSizeInSlot))
	}
	shunt.msgHandler = NewMsgHandler()
	return shunt
}

//RegisterMsgHandler : register msg handler
func (s *Shunt) RegisterMsgHandler(name string, fnIndex int, fn ProcMsgFn) {
	s.msgHandler.RegisterFn(name, fnIndex, fn)
}

//SizeOfSlot : get slot size
func (s *Shunt) SizeOfSlot() int {
	return s.slotSize
}

//SizeOfQInSlot : get queue size in each slot
func (s *Shunt) SizeOfQInSlot() int {
	return s.qSizeInSlot
}

//AddMsg : add msg
func (s *Shunt) AddMsg(ctx context.Context, slotIndex int, fnIndex int, inputMsg proto.Message) (*MsgProc, error) {
	slotIndex = s.normalizeSlotIndex(slotIndex)
	var msgProc = NewMsgProc(ctx, slotIndex, fnIndex, inputMsg)
	var err = convertQueueErr(s.qs[slotIndex].AddReq(msgProc))
	return msgProc, err
}

//AddPriorMsg : add prior msg
func (s *Shunt) AddPriorMsg(ctx context.Context, slotIndex int, fnIndex int, inputMsg proto.Message) (*MsgProc, error) {
	slotIndex = s.normalizeSlotIndex(slotIndex)
	var msgProc = NewMsgProc(ctx, slotIndex, fnIndex, inputMsg)
	var err = convertQueueErr(s.qs[slotIndex].AddPriorReq(msgProc))
	return msgProc, err
}

//Run : run all queue msg handler
func (s *Shunt) Run() {
	for i := 0; i < s.slotSize; i++ {
		go s.popLoop(i)
	}
}

//Stop : stop
func (s *Shunt) Stop() {
	s.stopOnce.Do(s.stop)
}

//WaitStop : wait stop
func (s *Shunt) WaitStop(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.exitChan:
		return nil
	}
}

//pop msg loop
func (s *Shunt) popLoop(index int) {
	var (
		err       error
		item      interface{}
		msgProc   *MsgProc
		outPutMsg proto.Message
		msgFn     ProcMsgFn
		ok        bool

		mq = s.qs[index]
	)

	defer s.wg.Done()
	for {
		item, err = mq.PopAnyway()
		if err != nil {
			ulog.Debug("msg.proc.module.item.quit",
				zap.Int("index", index),
				zap.Error(err))
			return
		}
		msgProc, ok = item.(*MsgProc)
		if !ok {
			ulog.Error("msg.proc.module.item.invalid.msg",
				zap.Int("index", index),
				zap.Reflect("inputMsg", item))
			return
		}
		msgFn = s.msgHandler.GetFn(msgProc.fnIndex)
		if msgFn == nil {
			msgProc.SetOutput(nil, ErrNoMsgHandler)
			continue
		}

		outPutMsg, err = msgFn(msgProc.ctx, index, msgProc.inputMsg)
		if err != nil {
			msgProc.SetOutput(nil, err)
		} else {
			msgProc.SetOutput(outPutMsg, nil)
		}
	}
}

//stop work
func (s *Shunt) stop() {
	for i := 0; i < s.slotSize; i++ {
		s.qs[i].Close()
	}
	//a go routine to wait all children done then signal it.
	go s.signalDone()
}

//signal all children go routine done
func (s *Shunt) signalDone() {
	s.wg.Wait()
	s.exitChan <- struct{}{}
}

//normalize slot index
func (s *Shunt) normalizeSlotIndex(index int) int {
	if index < 0 {
		index = -index
	}
	index %= s.slotSize
	return index
}

//convert msg queue error
func convertQueueErr(err error) error {
	if err == nil {
		return nil
	}
	if err == ErrReqQFull {
		return ErrMsgQueueFull
	}
	if err == ErrClosed {
		return ErrMsgQueueClosed
	}
	return err
}
