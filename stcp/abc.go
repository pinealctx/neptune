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
	// ObjectMarshaler marshal log object
	zapcore.ObjectMarshaler
	// GetRemoteAddr get remote address
	GetRemoteAddr() string
}

// IMsg : message interface
type IMsg interface {
	// Name : message name
	Name() string
}

// IConnSender connection sender interface
//
// Thread Safety & Re-entrance Requirements:
// - ALL methods MUST be goroutine-safe and re-entrant
// - Multiple calls to Close() should be safe (may return error but MUST NOT panic)
// - Resource cleanup operations should be idempotent
type IConnSender interface {
	// Conn returns the underlying network connection (required, goroutine-safe)
	Conn() net.Conn

	// SetMetaInfo sets meta info for logging (required, goroutine-safe, re-entrant)
	SetMetaInfo(m MetaInfo)

	// MetaInfo gets meta info for logging (required, goroutine-safe)
	MetaInfo() MetaInfo

	// Close closes the connection (required, goroutine-safe, re-entrant)
	// Multiple calls should be safe, may return error but MUST NOT panic
	Close() error

	// PutMsg put message to send (required, goroutine-safe)
	// return : error if any
	PutMsg(msg IMsg) error

	// PopMsgBytes pop message bytes to send (required, goroutine-safe)
	// return : message bytes and error if any
	PopMsgBytes() ([]byte, error)

	// sendBytes2Conn send bytes to connection, an internal utility function.
	// WARNING: This method is ONLY called by ConnHandler internally.
	// NEVER call this method from external code - it will cause undefined behavior.
	// return : error if any
	sendBytes2Conn(bs []byte) error
}

// IConnReader connection reader interface
type IConnReader interface {
	// ReadFrame reads one frame(an entire message bytes) from connection
	// conn : connection to read from
	// return : read buffer and error if any
	ReadFrame(conn net.Conn) ([]byte, error)
}

// IConnIO connection io interface
//
// Thread Safety & Re-entrance Requirements:
// - Multiple calls to Close() should be safe (may return error but MUST NOT panic)
// - Resource cleanup operations should be idempotent
type IConnIO interface {
	// IConnSender connection sender interface
	IConnSender
	// IConnReader connection reader interface
	IConnReader
}

// ConnStartEvent on connection start
type ConnStartEvent func(iConnIO IConnIO)

// ConnExitEvent on connection exit
type ConnExitEvent func(iConnIO IConnIO)

// ConnReaderFactory connection reader factory
type ConnReaderFactory func(conn net.Conn) IConnReader

// ConnIOFactory connection io factory
type ConnIOFactory func(conn net.Conn) IConnIO

// ReadProcessor read handler logic
// iConnIO : connection io interface
// buffer : read buffer
// return : error if any
// Actually, this is the core function to process the read data
type ReadProcessor func(iConnIO IConnIO, buffer []byte) error

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
