package stcp

import (
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/atomic"
	"net"
	"time"
)

//IMsgReader message handler
//消息读取接口
type IMsgReader interface {
	//Read : read
	Read(s *Session) error
}

//SessionMgr :
type SessionMgr struct {
	count        atomic.Int32
	logger       *ulog.Logger
	writeTimeout time.Duration
	readTimeout  time.Duration
	rh           IMsgReader
}

//NewSessionMgr :
func NewSessionMgr(rh IMsgReader, opts ...MOption) *SessionMgr {
	var cnf = defaultSessMgrOpt()
	for _, opt := range opts {
		opt(cnf)
	}
	return &SessionMgr{
		writeTimeout: cnf.writeTimeout,
		readTimeout:  cnf.readTimeout,
		rh:           rh,
	}
}

//SetLogger :
func (m *SessionMgr) SetLogger(logger *ulog.Logger) {
	m.logger = logger
}

//ConnCount 当前连接数
func (m *SessionMgr) ConnCount() int32 {
	return m.count.Load()
}

//Do :
func (m *SessionMgr) Do(conn net.Conn) {
	var session = NewSession(m, conn)
	session.Start()
}

//Logger : get logger
func (m *SessionMgr) Logger() *ulog.Logger {
	if m.logger == nil {
		return ulog.GetDefaultLogger()
	}
	return m.logger
}

//MOption session mgr option
type MOption func(o *_SessMgrOpt)

//WithWriteTimeout : setup write timeout
func WithWriteTimeout(t time.Duration) MOption {
	return func(o *_SessMgrOpt) {
		o.writeTimeout = t
	}
}

//WithReadTimeout : setup read timeout
func WithReadTimeout(t time.Duration) MOption {
	return func(o *_SessMgrOpt) {
		o.readTimeout = t
	}
}

//session mgr config
type _SessMgrOpt struct {
	writeTimeout time.Duration
	readTimeout  time.Duration
}

//get default session mgr config
func defaultSessMgrOpt() *_SessMgrOpt {
	return &_SessMgrOpt{
		writeTimeout: 8 * time.Second,  //8 second
		readTimeout:  20 * time.Second, //20 second
	}
}
