package stcp

import (
	"net"
	"time"

	"go.uber.org/zap"

	"github.com/pinealctx/neptune/ulog"
)

// IConnMgr interface
// 连接管理接口
type IConnMgr interface {
	//ConnCount 当前连接数
	ConnCount() int32
	//Do handle connection
	Do(conn net.Conn)
	//SetLogger setup logger
	SetLogger(logger *ulog.Logger)
}

// ITemporary interface
// to replace net.Error interface to avoid go lint check
type ITemporary interface {
	Temporary() bool
}

// Server : tcp server frame
type Server struct {
	ln      net.Listener
	ch      IConnMgr //connection manager
	address string
}

// Address :
func (s *Server) Address() string {
	return s.address
}

// NewTCPSrv :
func NewTCPSrv(address string, ch IConnMgr) *Server {
	return &Server{
		address: address,
		ch:      ch,
	}
}

// NewTCPSrvX : use a simple IConnMgr
func NewTCPSrvX(address string, rh ISession, opts ...MOption) *Server {
	var ch = NewSessionMgr(rh, opts...)
	return NewTCPSrv(address, ch)
}

// LoopStart :
func (s *Server) LoopStart(opts ...Option) error {
	var cnf = defaultStartOpt()
	for _, opt := range opts {
		opt(cnf)
	}

	s.ch.SetLogger(cnf.logger)
	var err = s.startListen(cnf)
	if err != nil {
		return err
	}

	return s.loopAccept(cnf)
}

// Start : loop start server will go loop state
// Use a channel to receive error
func (s *Server) Start(opts ...Option) <-chan error {
	var eh = make(chan error, 1)
	var err error
	go func() {
		err = s.LoopStart(opts...)
		if err != nil {
			eh <- err
		}
	}()
	return eh
}

// Close : close server
func (s *Server) Close() error {
	return s.ln.Close()
}

// start to listen
func (s *Server) startListen(cnf *_SrvStartOpt) error {
	var err error
	s.ln, err = net.Listen("tcp", s.address)
	if err != nil {
		cnf.Logger().Error("start.tcp.server",
			zap.Error(err), zap.String("listen", s.address))
		return err
	}
	return nil
}

// loop to accept
func (s *Server) loopAccept(cnf *_SrvStartOpt) error {
	var conn net.Conn
	var err error
	var errNet ITemporary
	var ok bool

	var accDelay time.Duration
	var accRetryCount int

	var handleErr = func() error {
		errNet, ok = err.(ITemporary)
		if !ok {
			return err
		}
		if !errNet.Temporary() {
			return err
		}

		//setup retry
		accRetryCount++
		if accRetryCount >= cnf.acceptMaxRetry {
			return err
		}
		//Temporary error
		if accDelay <= 0 {
			accDelay = cnf.acceptDelay
		} else {
			accDelay *= 2
		}
		if accDelay >= cnf.acceptMaxDelay {
			accDelay = cnf.acceptMaxDelay
		}
		time.Sleep(accDelay)
		return nil
	}

	var outErr error

	for {
		conn, err = s.ln.Accept()
		if err != nil {
			outErr = handleErr()
			if outErr != nil {
				cnf.Logger().Error("tcp.server.accept", zap.Error(err))
				return outErr
			}
			continue
		}
		accDelay = 0
		accRetryCount = 0

		if s.ch.ConnCount() >= cnf.maxConn {
			var e = conn.Close()
			cnf.Logger().Error("touch.max.conn.num", zap.Error(e))
		} else {
			//处理connection
			s.ch.Do(conn)
		}
	}
}

// Option server start option
type Option func(o *_SrvStartOpt)

// WithMaxConn : setup max conn number
func WithMaxConn(r int32) Option {
	return func(o *_SrvStartOpt) {
		o.maxConn = r
	}
}

// WithLogger : setup logger
func WithLogger(l *ulog.Logger) Option {
	return func(o *_SrvStartOpt) {
		o.logger = l
	}
}

// WithAccDelay : setup acceptDelay
func WithAccDelay(t time.Duration) Option {
	return func(o *_SrvStartOpt) {
		o.acceptDelay = t
	}
}

// WithAccMaxDelay : setup acceptMaxDelay
func WithAccMaxDelay(t time.Duration) Option {
	return func(o *_SrvStartOpt) {
		o.acceptMaxDelay = t
	}
}

// WithAccMaxRetry : setup acceptMaxRetry
func WithAccMaxRetry(r int) Option {
	return func(o *_SrvStartOpt) {
		o.acceptMaxRetry = r
	}
}

// server start config
type _SrvStartOpt struct {
	acceptDelay    time.Duration
	acceptMaxDelay time.Duration
	acceptMaxRetry int
	maxConn        int32
	logger         *ulog.Logger
}

func (o *_SrvStartOpt) Logger() *ulog.Logger {
	if o.logger == nil {
		return ulog.GetDefaultLogger()
	}
	return o.logger
}

// get default start cnf
func defaultStartOpt() *_SrvStartOpt {
	return &_SrvStartOpt{
		acceptDelay:    5 * time.Microsecond,
		acceptMaxDelay: 200 * time.Millisecond,
		acceptMaxRetry: 100,
		maxConn:        65535,
	}
}
