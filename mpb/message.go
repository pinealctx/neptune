package mpb

import (
	"encoding/binary"
	"github.com/pinealctx/neptune/errorx"
	"github.com/pinealctx/neptune/tex"
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
	"reflect"
)

const (
	// ErrMark error fingerprint
	ErrMark = 0
	// EmptyMark emptypb.Empty fingerprint
	EmptyMark = 1
)

var (
	_defaultMsgPacker = NewMsgPacker()
	_emptyData        = make([]byte, 4)
	_emptyMsg         = &emptypb.Empty{}
)

// FingerprintMsg extends "Fingerprint() uint32" to proto message
type FingerprintMsg interface {
	proto.Message
	Fingerprint() uint32
}

// MsgPacker pack/unpack a protobuf message with a header tag
type MsgPacker struct {
	genFuncMap map[uint32]func() proto.Message
	typeMap    map[reflect.Type]struct{}
}

// NewMsgPacker new MsgPacker instance
func NewMsgPacker() *MsgPacker {
	return &MsgPacker{
		genFuncMap: make(map[uint32]func() proto.Message),
		typeMap:    make(map[reflect.Type]struct{}),
	}
}

// RegisterGenerator register a protobuf generator function with tag
func (x *MsgPacker) RegisterGenerator(genFn func() proto.Message) {
	var exist bool
	var mo = genFn()
	if mo == nil {
		panic("generate func return nil")
	}
	var reflectV = reflect.ValueOf(mo)
	if reflectV.Kind() != reflect.Ptr {
		panic("generate func return not pointer")
	}
	if reflectV.IsNil() {
		panic("generate func return nil pointer")
	}

	var reflectT = reflect.TypeOf(mo)

	var fm, ok = mo.(FingerprintMsg)
	if !ok {
		panic("the message must implement function \"Fingerprint() uint32\" to return its unique fingerprint")
	}

	var fingerprint = fm.Fingerprint()
	_, exist = x.genFuncMap[fingerprint]
	if exist {
		panic("fingerprint already exist")
	}
	_, exist = x.typeMap[reflectT]
	if exist {
		panic("type already exist")
	}

	x.genFuncMap[fingerprint] = genFn
	x.typeMap[reflectT] = struct{}{}
}

// MarshalMsg marshal a protobuf message
// support *emptypb.Empty(it uses a specific tag "EmptyMark")
// and other proto message which extends "Fingerprint() uint32"
func (x *MsgPacker) MarshalMsg(msg proto.Message) ([]byte, error) {
	switch v := msg.(type) {
	case FingerprintMsg:
		return x.marshalMsg(v)
	case *emptypb.Empty:
		return _emptyData, nil
	default:
		return nil, errorx.NewfWithStack("unsupported message:%+v", msg.ProtoReflect().Descriptor().FullName())
	}
}

// UnmarshalMsg unmarshal a proto message from bytes.
// msg -- return msg
// err -- unmarshal error
func (x *MsgPacker) UnmarshalMsg(data []byte) (msg proto.Message, err error) {
	var preProc = unmarshalEmptyMsgOrErr(data)
	if preProc.sysErr != nil {
		return nil, preProc.sysErr
	}
	if preProc.union != nil {
		return preProc.union, nil
	}

	return x.unmarshalRegisteredMsg(preProc.fingerprint, data[4:])
}

// UnmarshalResponse unmarshal to rpc response from bytes.
// msg -- return msg
// msgErr -- return error
// err -- unmarshal error
func (x *MsgPacker) UnmarshalResponse(data []byte) (msg proto.Message, msgErr error, sysErr error) {
	var preProc = unmarshalEmptyMsgOrErr(data)
	if preProc.sysErr != nil {
		return nil, nil, preProc.sysErr
	}
	if preProc.msgErr != nil {
		return nil, preProc.msgErr, nil
	}
	if preProc.emptyMsg != nil {
		return preProc.emptyMsg, nil, nil
	}

	var m, e = x.unmarshalRegisteredMsg(preProc.fingerprint, data[4:])
	if e != nil {
		return nil, nil, e
	}
	return m, nil, nil
}

// UnmarshalRPC unmarshal response to 2 params return
func (x *MsgPacker) UnmarshalRPC(data []byte) (proto.Message, error) {
	return ToResponse(x.UnmarshalResponse(data))
}

func (x *MsgPacker) marshalMsg(msg FingerprintMsg) ([]byte, error) {
	var fingerprint = msg.Fingerprint()
	var _, ok = x.genFuncMap[fingerprint]
	if !ok {
		return nil, errorx.NewfWithStack("not registered message:%+v",
			msg.ProtoReflect().Descriptor().FullName())
	}
	return marshalProtoMsg(msg.Fingerprint(), msg)
}

func (x *MsgPacker) unmarshalRegisteredMsg(fingerprint uint32, data []byte) (proto.Message, error) {
	var gFn, ok = x.genFuncMap[fingerprint]
	if !ok {
		return nil, errorx.NewfWithStack("fingerprint:%x, not.found", fingerprint)
	}

	var m = gFn()
	var e = proto.Unmarshal(data, m)
	if e != nil {
		return nil, errorx.WrapWithStack(e, "unmarshal proto msg")
	}
	return m, nil
}

// RegisterGenerator register a protobuf generator function with tag
func RegisterGenerator(genFn func() proto.Message) {
	_defaultMsgPacker.RegisterGenerator(genFn)
}

// MarshalMsg marshal a protobuf message
// support *emptypb.Empty(it uses a specific tag "EmptyMark")
// and other proto message which extends "Fingerprint() uint32"
func MarshalMsg(msg proto.Message) ([]byte, error) {
	return _defaultMsgPacker.MarshalMsg(msg)
}

// MarshalError err can marshal/unmarshal, it uses a specific tag "ErrMark"
func MarshalError(err error) ([]byte, error) {
	var v, _ = status.FromError(err)
	return marshalProtoMsg(ErrMark, v.Proto())
}

// MarshalEmpty tag an empty message
func MarshalEmpty() []byte {
	return _emptyData
}

// UnmarshalMsg unmarshal a proto message from bytes.
// msg -- return msg
// err -- unmarshal error
func UnmarshalMsg(data []byte) (msg proto.Message, err error) {
	return _defaultMsgPacker.UnmarshalMsg(data)
}

// UnmarshalResponse unmarshal to rpc response from bytes.
// msg -- return msg
// msgErr -- return error
// err -- unmarshal error
func UnmarshalResponse(data []byte) (msg proto.Message, msgErr error, sysErr error) {
	return _defaultMsgPacker.UnmarshalResponse(data)
}

// UnmarshalRPC unmarshal response to 2 params return
func UnmarshalRPC(data []byte) (proto.Message, error) {
	return _defaultMsgPacker.UnmarshalRPC(data)
}

// UnmarshalEmptyResponse unmarshal to rpc response from bytes.
// msg -- return msg
// msgErr -- return error
// err -- unmarshal error
func UnmarshalEmptyResponse(data []byte) (msgErr error, sysErr error) {
	var preProc = unmarshalEmptyMsgOrErr(data)
	if preProc.sysErr != nil {
		return nil, preProc.sysErr
	}
	if preProc.msgErr != nil {
		return preProc.msgErr, nil
	}
	if preProc.emptyMsg == nil {
		return errorx.NewfWithStack("excepted empty msg:%+v", preProc.fingerprint), nil
	}
	return nil, nil
}

// UnmarshalEmptyRPC unmarshal to rpc response from bytes.
// msg -- return msg
// msgErr -- return error
// err -- unmarshal error
func UnmarshalEmptyRPC(data []byte) error {
	return ToErr(UnmarshalEmptyResponse(data))
}

// ToResponse convert 3 params to 2 params
func ToResponse(msg proto.Message, msgErr error, sysErr error) (proto.Message, error) {
	if sysErr != nil {
		return nil, sysErr
	}
	if msgErr != nil {
		return nil, msgErr
	}
	return msg, nil
}

// ToErr combine error
func ToErr(msgErr error, sysErr error) error {
	if sysErr != nil {
		return sysErr
	}
	return msgErr
}

func marshalProtoMsg(fingerprint uint32, msg proto.Message) ([]byte, error) {
	var data, err = proto.Marshal(msg)
	if err != nil {
		return nil, errorx.WrapWithStack(err, "marshal proto msg")
	}
	var size = len(data) + 4
	var buf = tex.NewSizedBuffer(size)
	var header [4]byte
	binary.LittleEndian.PutUint32(header[:], fingerprint)
	_, _ = buf.Write(header[:])
	_, _ = buf.Write(data)
	return buf.Bytes(), nil
}

type emptyOrErrMsgT struct {
	fingerprint uint32
	union       proto.Message
	emptyMsg    *emptypb.Empty
	msgErr      error
	sysErr      error
}

func feedEmptyOrErrMsg(fingerprint uint32, ep *emptypb.Empty, msgErr *spb.Status, sysErr error) *emptyOrErrMsgT {
	var x = &emptyOrErrMsgT{
		fingerprint: fingerprint,
		emptyMsg:    ep,
		sysErr:      sysErr,
	}
	if msgErr != nil {
		x.msgErr = status.FromProto(msgErr).Err()
		x.union = msgErr
	}
	if ep != nil {
		x.union = ep
	}
	return x
}

func unmarshalEmptyMsgOrErr(data []byte) *emptyOrErrMsgT {
	var size = len(data)
	if size < 4 {
		return feedEmptyOrErrMsg(0, nil, nil, errorx.NewWithStack("invalid message length"))
	}
	var fingerprint = binary.LittleEndian.Uint32(data)

	if fingerprint == ErrMark {
		var mErr, e = unmarshalErr(data[4:])
		if e != nil {
			return feedEmptyOrErrMsg(fingerprint, nil, nil, e)
		}
		return feedEmptyOrErrMsg(fingerprint, nil, mErr, nil)
	}

	if fingerprint == EmptyMark {
		return feedEmptyOrErrMsg(fingerprint, _emptyMsg, nil, nil)
	}
	return feedEmptyOrErrMsg(fingerprint, nil, nil, nil)
}

func unmarshalErr(data []byte) (*spb.Status, error) {
	var wrapErr = &spb.Status{}
	var err = proto.Unmarshal(data, wrapErr)
	if err != nil {
		return nil, errorx.WrapWithStack(err, "unmarshal proto error")
	}
	return wrapErr, nil
}

func init() {
	binary.LittleEndian.PutUint32(_emptyData, EmptyMark)
}
