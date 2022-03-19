package stcp

import (
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"io"
	"net"
	"sync"
	"time"
)

//Echo 应答式session
type Echo struct {
	//连接
	conn net.Conn
	//所带字段
	value atomic.Value
	//base handler
	b *EchoMgr
	//start once
	startOnce sync.Once
	//exit once
	exitOnce sync.Once
}

//NewEcho :
func NewEcho(b *EchoMgr, conn net.Conn) *Echo {
	return &Echo{
		b:    b,
		conn: conn,
	}
}

//Start :
func (s *Echo) Start() {
	s.startOnce.Do(func() {
		s.b.count.Inc()
		go s.b.eh.RunEcho(s)
	})
}

//ReleaseCount : when Echo session is disposed, this function should be called to decrease connection counter
//减少引用计数
func (s *Echo) ReleaseCount() {
	s.exitOnce.Do(func() {
		s.b.count.Dec()
	})
}

//Set :
func (s *Echo) Set(v interface{}) {
	s.value.Store(v)
}

//Get :
func (s *Echo) Get() interface{} {
	return s.value.Load()
}

//RemoteAddr :
func (s *Echo) RemoteAddr() string {
	if s.conn == nil {
		return ""
	}
	var ra = s.conn.RemoteAddr()
	if ra == nil {
		return ""
	}
	return ra.String()
}

//KeyOut : for uber log
func (s *Echo) KeyOut() zap.Field {
	return absSessionInfo(s.value, false)
}

//AllInfo : for uber log
func (s *Echo) AllInfo() zap.Field {
	return absSessionInfo(s.value, true)
}

//RemoteInfo : for uber log
func (s *Echo) RemoteInfo() zap.Field {
	return zap.String("session.Addr", s.RemoteAddr())
}

//Send : send bytes, put bytes to queue, not send directly
func (s *Echo) Send(bs []byte) error {
	var err = s.conn.SetWriteDeadline(time.Now().Add(s.b.writeTimeout))
	if err != nil {
		s.Logger().Error("set.write.conn.deadline", zap.Error(err), s.RemoteInfo())
		return err
	}
	_, err = s.conn.Write(bs)
	return err
}

//Read : read specific bytes
func (s *Echo) Read(bs []byte) error {
	var err = s.conn.SetReadDeadline(time.Now().Add(s.b.readTimeout))
	if err != nil {
		s.Logger().Error("set.read.conn.deadline", zap.Error(err), s.RemoteInfo())
		return err
	}
	_, err = io.ReadFull(s.conn, bs)
	return err
}

//Logger : get logger
func (s *Echo) Logger() *ulog.Logger {
	return s.b.Logger()
}

//IEcho echo handler
type IEcho interface {
	//RunEcho :
	RunEcho(s *Echo)
}

//EchoMgr :
type EchoMgr struct {
	count        atomic.Int32
	logger       *ulog.Logger
	writeTimeout time.Duration
	readTimeout  time.Duration
	eh           IEcho
}

//NewEchoMgr :
func NewEchoMgr(eh IEcho, opts ...MOption) *EchoMgr {
	var cnf = defaultSessMgrOpt()
	for _, opt := range opts {
		opt(cnf)
	}
	return &EchoMgr{
		writeTimeout: cnf.writeTimeout,
		readTimeout:  cnf.readTimeout,
		eh:           eh,
	}
}

//SetLogger :
func (m *EchoMgr) SetLogger(logger *ulog.Logger) {
	m.logger = logger
}

//ConnCount 当前连接数
func (m *EchoMgr) ConnCount() int32 {
	return m.count.Load()
}

//Do :
func (m *EchoMgr) Do(conn net.Conn) {
	var echo = NewEcho(m, conn)
	echo.Start()
}

//Logger : get logger
func (m *EchoMgr) Logger() *ulog.Logger {
	if m.logger == nil {
		return ulog.GetDefaultLogger()
	}
	return m.logger
}
