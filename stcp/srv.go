package stcp

import (
	"net"
	"time"

	"github.com/pinealctx/neptune/timex"
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/atomic"
	"go.uber.org/zap"
)

// ServerAcceptCnf server start config
type ServerAcceptCnf struct {
	Address        string         `json:"address"`
	AcceptDelay    timex.Duration `json:"acceptDelay"`
	AcceptMaxDelay timex.Duration `json:"acceptMaxDelay"`
	AcceptMaxRetry int            `json:"acceptMaxRetry"`
	MaxConn        int32          `json:"maxConn"`
}

// DefaultServerAcceptCnf : get default start cnf
func DefaultServerAcceptCnf() *ServerAcceptCnf {
	return &ServerAcceptCnf{
		AcceptDelay:    timex.NewDuration(5 * time.Microsecond),
		AcceptMaxDelay: timex.NewDuration(200 * time.Millisecond),
		AcceptMaxRetry: 100,
		MaxConn:        65535,
	}
}

// CommServerCnf common server config
type CommServerCnf struct {
	*ServerAcceptCnf
	MaxSendQSize int `json:"maxSendQSize"`
}

// DefaultCommServerCnf : get default common server cnf
func DefaultCommServerCnf() *CommServerCnf {
	return &CommServerCnf{
		ServerAcceptCnf: DefaultServerAcceptCnf(),
		MaxSendQSize:    512,
	}
}

// ITemporary interface
// to replace net.Error interface to avoid go lint check
type ITemporary interface {
	Temporary() bool
}

// ConnHandlerGenerator connection handler generator
type ConnHandlerGenerator func(conn net.Conn) IConnHandler

// BasicServer : basic tcp server frame
type BasicServer struct {
	ln                   net.Listener
	connectionCount      atomic.Int32
	svrCnf               *ServerAcceptCnf
	connHandlerGenerator ConnHandlerGenerator
	startHooker          ConnStartEvent
	exitHooker           ConnExitEvent
}

// NewBasicServer :
func NewBasicServer(cnf *ServerAcceptCnf, gen ConnHandlerGenerator) *BasicServer {
	return &BasicServer{
		svrCnf:               cnf,
		connHandlerGenerator: gen,
	}
}

// Address :
func (x *BasicServer) Address() string {
	return x.svrCnf.Address
}

// ConnectionCount :
func (x *BasicServer) ConnectionCount() int32 {
	return x.connectionCount.Load()
}

// Run : loop start server
// errChan : a channel to receive error
// cnf : server start config
func (s *CommonServer) Run(errChan chan<- error) {
	go func() {
		err := s.Start()
		errChan <- err
	}()
}

// Close : close server
func (s *CommonServer) Close() error {
	return s.ln.Close()
}

// Start : loop start server with config
func (s *CommonServer) Start() error {
	var err = s.startListen()
	if err != nil {
		return err
	}
	return s.loopAccept()
}

// SetStartHook : set start hook
func (x *BasicServer) SetStartHook(hook ConnStartEvent) {
	x.startHooker = hook
}

// SetExitHook : set exit hook
func (x *BasicServer) SetExitHook(hook ConnExitEvent) {
	x.exitHooker = hook
}

// loop to accept
func (x *BasicServer) loopAccept() error {
	var conn net.Conn
	var err error
	var errNet ITemporary
	var ok bool

	var accDelay time.Duration
	var accRetryCount int

	accCnfDelay := x.svrCnf.AcceptDelay.Value()
	accCnfMaxDelay := x.svrCnf.AcceptMaxDelay.Value()
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
		if accRetryCount >= x.svrCnf.AcceptMaxRetry {
			return err
		}
		//Temporary error
		if accDelay <= 0 {
			accDelay = accCnfDelay
		} else {
			accDelay *= 2
		}
		if accDelay >= accCnfMaxDelay {
			accDelay = accCnfMaxDelay
		}
		time.Sleep(accDelay)
		return nil
	}

	var outErr error

	ulog.Info("tcp.server.accept.loop.start")
	for {
		conn, err = x.ln.Accept()
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

		if x.connectionCount.Inc() >= x.svrCnf.MaxConn {
			var e = conn.Close()
			ulog.Error("touch.max.conn.num", zap.Error(e), zap.Int32("currentConn", x.connectionCount.Load()))
		} else {
			//处理connection
			connHandler := NewConnHandlerRunner(x.connHandlerGenerator(conn))
			connHandler.AddStartHook(x.startHook)
			connHandler.AddExitHook(x.exitHook)
			connHandler.Start()
		}
	}
}

// startHook : when connection start
func (x *BasicServer) startHook(metaInfo MetaInfo) {
	ulog.Info("connection.start", zap.Object("metaInfo", metaInfo),
		zap.Int32("currentConn", x.connectionCount.Load()))
	if x.startHooker != nil {
		x.startHooker(metaInfo)
	}
}

// exitHook : when connection exit
func (x *BasicServer) exitHook(conn net.Conn, metaInfo MetaInfo) {
	x.connectionCount.Dec()
	ulog.Info("connection.exit", zap.Object("metaInfo", metaInfo),
		zap.Int32("currentConn", x.connectionCount.Load()))
	if x.exitHooker != nil {
		x.exitHooker(conn, metaInfo)
	}
}

// start to listen
func (x *BasicServer) startListen() error {
	var err error
	x.ln, err = net.Listen("tcp", x.svrCnf.Address)
	if err != nil {
		ulog.Error("start.tcp.server",
			zap.Error(err), zap.String("listen", x.svrCnf.Address))
		return err
	}
	ulog.Info("tcp.server.started", zap.String("listen", x.svrCnf.Address))
	return nil
}

// CommonServer : common tcp server frame
type CommonServer struct {
	*BasicServer
	incomingHook CommonIncomingHook
	maxSendQSize int
}

// NewCommonTCPSrv : new common tcp server
func NewCommonTCPSrv(cnf *CommServerCnf, incomingHook CommonIncomingHook) *CommonServer {
	x := &CommonServer{
		incomingHook: incomingHook,
		maxSendQSize: cnf.MaxSendQSize,
	}
	basicServer := NewBasicServer(cnf.ServerAcceptCnf, x.connHandlerGenerator)
	x.BasicServer = basicServer
	return x
}

// common connection handler generator
func (s *CommonServer) connHandlerGenerator(conn net.Conn) IConnHandler {
	return NewCommonConnHandler(conn, s.maxSendQSize, s.incomingHook)
}
