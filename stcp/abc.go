package stcp

import (
	"net"
	"time"

	"go.uber.org/atomic"
	"go.uber.org/zap/zapcore"
)

var (
	// write timeout
	writeTimeout *atomic.Duration
)

// MetaInfo meta info for logging
type MetaInfo interface {
	zapcore.ObjectMarshaler
	GetRemoteAddr() string
}

// KeyIntBytesPair key uint32 and bytes pair
type KeyIntBytesPair struct {
	Key uint32
	Val []byte
}

// KeyStrBytesPair key string and bytes pair
type KeyStrBytesPair struct {
	Key string
	Val []byte
}

// IConnSender connection handler interface
// all methods are goroutine safe
// user must implement this interface
// Conn/LoopSend/SetMetaInfo/MetaInfo/Close are required
// Put2Queue/Put2SendMap/Put2SendSMap/Put2SendMaps/Put2SendSMaps are optional, but at least one of them should be implemented, others can be no-op
type IConnSender interface {
	// Conn : get connection (required)
	Conn() net.Conn
	// LoopSend : loop to send (required)
	LoopSend()
	// SetMetaInfo set meta info (required)
	SetMetaInfo(m MetaInfo)
	// MetaInfo get meta info (required)
	MetaInfo() MetaInfo
	// Close : close connection handler
	Close() error

	// Put2Queue put bytes to send queue (optional)
	Put2Queue(bs []byte) error
	// Put2SendMap put bytes to send map (optional)
	Put2SendMap(key uint32, bs []byte) error
	// Put2SendSMap put bytes to send map (optional)
	Put2SendSMap(key string, bs []byte) error
	// Put2SendMaps put multiple key uint32 and bytes pairs to send map (optional)
	Put2SendMaps(pairs []KeyIntBytesPair) error
	// Put2SendSMaps put multiple key string and bytes pairs to send map (optional)
	Put2SendSMaps(pairs []KeyStrBytesPair) error
}

// ConnStartEvent on connection start
type ConnStartEvent func(handler IConnSender)

// ConnExitEvent on connection exit
type ConnExitEvent func(handler IConnSender)

// ConnReaderFunc connection reader function
type ConnReaderFunc func(handler IConnSender, conn net.Conn) error

// ConnSenderFactory connection sender factory
type ConnSenderFactory func(conn net.Conn) IConnSender

// BasicMetaInfo basic meta info
type BasicMetaInfo struct {
	RemoteAddr string
}

// NewBasicMetaInfo new basic meta info
func NewBasicMetaInfo(conn net.Conn) *BasicMetaInfo {
	return &BasicMetaInfo{
		RemoteAddr: conn.RemoteAddr().String(),
	}
}

// MarshalLogObject marshal log object
func (m *BasicMetaInfo) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("remoteAddr", m.RemoteAddr)
	return nil
}

// GetRemoteAddr get remote address
func (m *BasicMetaInfo) GetRemoteAddr() string {
	return m.RemoteAddr
}

// SendBytes2Conn send bytes to connection
// Utility function
func SendBytes2Conn(conn net.Conn, bs []byte) error {
	err := conn.SetWriteDeadline(time.Now().Add(writeTimeout.Load()))
	if err != nil {
		return err
	}
	_, err = conn.Write(bs)
	return err
}

// SetWriteTimeout set write timeout
// Utility function, it's a global setting
// If not set, default is 5 seconds
func SetWriteTimeout(d time.Duration) {
	writeTimeout.Store(d)
}

func init() {
	writeTimeout = atomic.NewDuration(time.Second * 5)
}
