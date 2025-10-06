package stcp

import (
	"fmt"
	"net"
	"sync"

	"go.uber.org/atomic"
)

// BytesMsg : bytes message
type BytesMsg struct {
	Bs []byte
}

// NewBytesMsg : new bytes message
func NewBytesMsg(bs []byte) *BytesMsg {
	return &BytesMsg{Bs: bs}
}

// Name : message name
func (x *BytesMsg) Name() string {
	return "BytesMsg"
}

// BasicConnIO basic connection io
type BasicConnIO struct {
	// close once
	closeOnce sync.Once
	// connection
	conn net.Conn
	// meta info
	metaInfo atomic.Value
	// reader
	reader IConnReader
}

// NewBasicConnIO : new basic connection io
func NewBasicConnIO(conn net.Conn, reader IConnReader) *BasicConnIO {
	x := &BasicConnIO{
		conn:   conn,
		reader: reader,
	}
	x.metaInfo.Store(NewBasicMetaInfo(conn))
	return x
}

// Conn returns the underlying network connection
func (x *BasicConnIO) Conn() net.Conn {
	return x.conn
}

// SetMetaInfo sets meta info for logging (required, goroutine-safe, re-entrant)
func (x *BasicConnIO) SetMetaInfo(m MetaInfo) {
	x.metaInfo.Store(m)
}

// MetaInfo gets meta info for logging (required, goroutine-safe)
func (x *BasicConnIO) MetaInfo() MetaInfo {
	v := x.metaInfo.Load()
	if v == nil {
		return nil
	}
	// nolint : forcetypeassert // I know the type is exactly here
	return v.(MetaInfo)
}

// CloseX closes connection handler (required, goroutine-safe, re-entrant)
// This method can be called directly via IConnSender/IConnIO interface to gracefully shutdown
// the connection and trigger the associated ConnHandler.Exit() through the goroutine defer chain
// Multiple calls are safe and will not panic
// closeFn : custom close function, if nil, use x.conn.Close()
func (x *BasicConnIO) CloseX(closeFn func() error) error {
	var err error
	x.closeOnce.Do(func() {
		err = closeFn()
	})
	return err
}

// ReadFrame reads one frame(an entire message bytes) from connection
func (x *BasicConnIO) ReadFrame(conn net.Conn) ([]byte, error) {
	buf, err := x.reader.ReadFrame(conn)
	if err != nil {
		return nil, fmt.Errorf("BasicConnIO.ReadFrame: %w", err)
	}
	return buf, nil
}

// sendBytes2Conn send bytes to connection, an internal utility function.
// WARNING: This method is ONLY called by ConnHandler internally.
// NEVER call this method from external code - it will cause undefined behavior.
func (x *BasicConnIO) sendBytes2Conn(bs []byte) error {
	return SendBytes2Conn(x.conn, bs)
}
