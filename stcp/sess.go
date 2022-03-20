package stcp

import (
	"github.com/pinealctx/neptune/syncx/pipe/q"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net"
	"sync"
	"time"
)

//NetIO : net io
type NetIO interface {
	//Set set session related value
	Set(v interface{})
	//Get get session related value
	Get() interface{}
	//SetRemoteAddr :
	SetRemoteAddr(addr string)
	//RemoteAddr :
	RemoteAddr() string
	//Send : send bytes
	Send(bs []byte) error
	//Read : read specific bytes
	Read(bs []byte) error
	//Close : close session
	Close()
	//Logger : get uber logger
	Logger() *ulog.Logger
	//RemoteInfo : for uber log
	RemoteInfo() zap.Field
	//AllInfo : for uber log, all info
	AllInfo() zap.Field
	//KeyOut : for uber log, key info
	KeyOut() zap.Field
}

//ISessDebug : 做调试用的信息
type ISessDebug interface {
	//KeyOut 可打印的关键信息
	KeyOut() zapcore.ObjectMarshaler
	//All Debug info
	All() zapcore.ObjectMarshaler
}

//Session session
type Session struct {
	//连接
	conn net.Conn
	//所带字段
	value atomic.Value
	//base handler
	b *SessionMgr
	//发送q
	//queue -- actually the queue is bytes
	sendQ *q.Q
	//remote addr
	remoteAddr atomic.String
	//start once
	startOnce sync.Once
	//exit once
	exitOnce sync.Once
}

//NewSession :
func NewSession(b *SessionMgr, conn net.Conn) *Session {
	return &Session{
		b:     b,
		conn:  conn,
		sendQ: q.NewQ(),
	}
}

//Start :
func (s *Session) Start() {
	s.startOnce.Do(func() {
		s.b.count.Inc()
		go s.loopSend()
		go s.loopReceive()
	})
}

//Set :
func (s *Session) Set(v interface{}) {
	s.value.Store(v)
}

//Get :
func (s *Session) Get() interface{} {
	return s.value.Load()
}

//SetRemoteAddr :
func (s *Session) SetRemoteAddr(addr string) {
	s.remoteAddr.Store(addr)
}

//RemoteAddr :
func (s *Session) RemoteAddr() string {
	return absRemoteAddr(s.remoteAddr, s.conn)
}

//Send : send bytes, put bytes to queue, not send directly
func (s *Session) Send(bs []byte) error {
	return s.sendQ.AddReq(bs)
}

//Read : read specific bytes
func (s *Session) Read(bs []byte) error {
	var _, err = io.ReadFull(s.conn, bs)
	return err
}

//Close : close session
//use close send msg queue to exit loop read/write go routine
func (s *Session) Close() {
	s.sendQ.Close()
}

//Logger : get logger
func (s *Session) Logger() *ulog.Logger {
	return s.b.Logger()
}

//RemoteInfo : for uber log
func (s *Session) RemoteInfo() zap.Field {
	return zap.String("session.Addr", s.RemoteAddr())
}

//AllInfo : for uber log
func (s *Session) AllInfo() zap.Field {
	return absSessionInfo(s.value, true)
}

//KeyOut : for uber log
func (s *Session) KeyOut() zap.Field {
	return absSessionInfo(s.value, false)
}

//loop send
func (s *Session) loopSend() {
	var (
		err   error
		qItem interface{}
		bs    []byte
		ok    bool
	)
	defer s.recovery()
	defer s.quit()

	for {
		qItem, err = s.sendQ.PopAnyway()
		if err != nil {
			s.b.Logger().Debug("quit.in.send.q", zap.Error(err), s.KeyOut())
			return
		}
		bs, ok = qItem.([]byte)
		if !ok || len(bs) == 0 {
			s.b.Logger().Error("invalid.send.q", zap.Any("q", qItem))
			return
		}
		err = s.send(bs)
		if err != nil {
			s.loggerSendReadErr("connection.send.fail", err)
			return
		}
	}
}

//loop receive
func (s *Session) loopReceive() {
	defer s.recovery()
	defer s.quit()

	for {
		var err = s.conn.SetReadDeadline(time.Now().Add(s.b.readTimeout))
		if err != nil {
			s.b.Logger().Error("set.read.conn.deadline", zap.Error(err), s.RemoteInfo())
			return
		}
		err = s.b.rh.Read(s)
		if err != nil {
			s.loggerSendReadErr("connection.read.fail", err)
			return
		}
	}
}

//send buffer
func (s *Session) send(buf []byte) error {
	var err = s.conn.SetWriteDeadline(time.Now().Add(s.b.writeTimeout))
	if err != nil {
		s.b.Logger().Error("set.write.conn.deadline", zap.Error(err), s.RemoteInfo())
		return err
	}
	_, err = s.conn.Write(buf)
	return err
}

//quit :
func (s *Session) quit() {
	s.exitOnce.Do(func() {
		s.b.rh.OnExit(s)
		s.b.count.Dec()
		s.sendQ.Close()
		if s.conn != nil {
			var err = s.conn.Close()
			if err != nil {
				s.b.Logger().Error("close.conn", zap.Error(err), s.RemoteInfo())
			}
		}
	})
}

//recovery :
func (s *Session) recovery() {
	var r = recover()
	if r != nil {
		//has panic
		s.b.Logger().Error("session.panic", zap.Any("panic", r),
			zap.Stack("stack"))
	}
}

//log send/read error
func (s *Session) loggerSendReadErr(msg string, err error) {
	if s.b.Logger().Level() >= zapcore.WarnLevel {
		s.b.Logger().Warn(msg, zap.Error(err), s.KeyOut(), s.RemoteInfo())
	}
}

//abs remote address
func absRemoteAddr(rdAddr atomic.String, conn net.Conn) string {
	var rd = rdAddr.Load()
	if rd == "" {
		if conn == nil {
			return ""
		}
		var ra = conn.RemoteAddr()
		if ra == nil {
			return ""
		}
		return ra.String()
	}
	return rd
}

//abs session info to debug info
func absSessionInfo(value atomic.Value, all bool) zap.Field {
	var info = value.Load()
	if info == nil {
		return zap.Bool("sessionInfo.value.empty", true)
	}
	var sessInfo, ok = info.(ISessDebug)
	if ok {
		if sessInfo == nil {
			return zap.Bool("sessionInfo.value.empty", false)
		}
		if all {
			return zap.Object("sessionInfo.All", sessInfo.All())
		}
		return zap.Object("sessionInfo", sessInfo.KeyOut())
	} else {
		return zap.Any("unknown.sessionInfo", info)
	}
}
