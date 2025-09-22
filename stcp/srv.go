package stcp

import (
	"net"
	"time"

	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// ServerAcceptCnf server start config
type ServerAcceptCnf struct {
	AcceptDelay    time.Duration `json:"acceptDelay"`
	AcceptMaxDelay time.Duration `json:"acceptMaxDelay"`
	AcceptMaxRetry int           `json:"acceptMaxRetry"`
	MaxSendQSize   int           `json:"maxSendQSize"`
	MaxConn        int32         `json:"maxConn"`
}

// DefaultServerAcceptCnf : get default start cnf
func DefaultServerAcceptCnf() *ServerAcceptCnf {
	return &ServerAcceptCnf{
		AcceptDelay:    5 * time.Microsecond,
		AcceptMaxDelay: 200 * time.Millisecond,
		AcceptMaxRetry: 100,
		MaxSendQSize:   1024,
		MaxConn:        65535,
	}
}

// ITemporary interface
// to replace net.Error interface to avoid go lint check
type ITemporary interface {
	Temporary() bool
}

// Server : tcp server frame
type Server struct {
	ln              net.Listener
	connectionCount atomic.Int32
	incomingHook    IncomingHook
	address         string
}

// NewTCPSrv :
func NewTCPSrv(address string, incomingHook IncomingHook) *Server {
	return &Server{
		address:      address,
		incomingHook: incomingHook,
	}
}

// Address :
func (s *Server) Address() string {
	return s.address
}

// ConnectionCount :
func (s *Server) ConnectionCount() int32 {
	return s.connectionCount.Load()
}

// RunWithOption : loop start server
// errChan : a channel to receive error
// opts : server start options
func (s *Server) RunWithOption(errChan chan<- error, opts ...Option) {
	go func() {
		err := s.LoopStart(opts...)
		errChan <- err
	}()
}

// RunWithCnf : loop start server with config
// errChan : a channel to receive error
// cnf : server start config
func (s *Server) RunWithCnf(errChan chan<- error, cnf *ServerAcceptCnf) {
	go func() {
		err := s.LoopStartX(cnf)
		errChan <- err
	}()
}

// LoopStart : loop start server will go loop state
func (s *Server) LoopStart(opts ...Option) error {
	var cnf = DefaultServerAcceptCnf()
	for _, opt := range opts {
		opt(cnf)
	}
	return s.LoopStartX(cnf)
}

// LoopStartX : loop start server with config
func (s *Server) LoopStartX(cnf *ServerAcceptCnf) error {
	var err = s.startListen()
	if err != nil {
		return err
	}
	return s.loopAccept(cnf)
}

// Close : close server
func (s *Server) Close() error {
	return s.ln.Close()
}

// start to listen
func (s *Server) startListen() error {
	var err error
	s.ln, err = net.Listen("tcp", s.address)
	if err != nil {
		ulog.Error("start.tcp.server",
			zap.Error(err), zap.String("listen", s.address))
		return err
	}
	return nil
}

// loop to accept
func (s *Server) loopAccept(cnf *ServerAcceptCnf) error {
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
		if accRetryCount >= cnf.AcceptMaxRetry {
			return err
		}
		//Temporary error
		if accDelay <= 0 {
			accDelay = cnf.AcceptDelay
		} else {
			accDelay *= 2
		}
		if accDelay >= cnf.AcceptMaxDelay {
			accDelay = cnf.AcceptMaxDelay
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
				ulog.Error("tcp.server.accept", zap.Error(err))
				return outErr
			}
			continue
		}
		accDelay = 0
		accRetryCount = 0

		if s.connectionCount.Inc() >= cnf.MaxConn {
			var e = conn.Close()
			ulog.Error("touch.max.conn.num", zap.Error(e), zap.Int32("currentConn", s.connectionCount.Load()))
		} else {
			//处理connection
			connHandler := NewConnHandlerV3(conn, cnf.MaxSendQSize, s.incomingHook, s.startHook, s.exitHook)
			connHandler.Start()
		}
	}
}

// startHook : when connection start
func (s *Server) startHook(metaInfo MetaInfo) {
	ulog.Info("connection.start", zap.Object("metaInfo", metaInfo),
		zap.Int32("currentConn", s.connectionCount.Load()))
}

// exitHook : when connection exit
func (s *Server) exitHook(_ net.Conn, metaInfo MetaInfo) {
	s.connectionCount.Dec()
	ulog.Info("connection.exit", zap.Object("metaInfo", metaInfo),
		zap.Int32("currentConn", s.connectionCount.Load()))
}

// Option server start option
type Option func(o *ServerAcceptCnf)

// WithMaxConn : setup max conn number
func WithMaxConn(r int32) Option {
	return func(o *ServerAcceptCnf) {
		o.MaxConn = r
	}
}

// WithAccDelay : setup acceptDelay
func WithAccDelay(t time.Duration) Option {
	return func(o *ServerAcceptCnf) {
		o.AcceptDelay = t
	}
}

// WithAccMaxDelay : setup acceptMaxDelay
func WithAccMaxDelay(t time.Duration) Option {
	return func(o *ServerAcceptCnf) {
		o.AcceptMaxDelay = t
	}
}

// WithAccMaxRetry : setup acceptMaxRetry
func WithAccMaxRetry(r int) Option {
	return func(o *ServerAcceptCnf) {
		o.AcceptMaxRetry = r
	}
}

// WithMaxSendQSize : setup max send queue size
func WithMaxSendQSize(s int) Option {
	return func(o *ServerAcceptCnf) {
		o.MaxSendQSize = s
	}
}
