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
		AcceptMaxDelay: timex.NewDuration(50 * time.Millisecond),
		AcceptMaxRetry: 0,
		MaxConn:        65535,
	}
}

// TcpServer tcp server
type TcpServer struct {
	acceptCnf         *ServerAcceptCnf
	startHooker       ConnStartEvent
	exitHooker        ConnExitEvent
	connReader        ConnReaderFunc
	connSenderFactory ConnSenderFactory

	connCount atomic.Int32
	ln        net.Listener
}

// NewTcpServer : new tcp server
func NewTcpServer(cnf *ServerAcceptCnf, connReader ConnReaderFunc, connSenderFactory ConnSenderFactory) *TcpServer {
	return &TcpServer{
		acceptCnf:         cnf,
		connReader:        connReader,
		connSenderFactory: connSenderFactory,
	}
}

// Address : get listen address
func (x *TcpServer) Address() string {
	return x.acceptCnf.Address
}

// ConnCount : get current connection count
func (x *TcpServer) ConnCount() int32 {
	return x.connCount.Load()
}

// SetStartHooker : set connection start hooker
func (x *TcpServer) SetStartHooker(hooker ConnStartEvent) {
	x.startHooker = hooker
}

// SetExitHooker : set connection exit hooker
func (x *TcpServer) SetExitHooker(hooker ConnExitEvent) {
	x.exitHooker = hooker
}

// Run : run server
// errChan : error channel
// if server exit, the error will be sent to errChan
func (x *TcpServer) Run(errChan chan<- error) {
	go func() {
		err := x.start()
		errChan <- err
	}()
}

// Close : close the server listener
func (x *TcpServer) Close() error {
	if x.ln != nil {
		return x.ln.Close()
	}
	return nil
}

// start : start server
func (x *TcpServer) start() error {
	var err error
	x.ln, err = net.Listen("tcp", x.acceptCnf.Address)
	if err != nil {
		return err
	}
	return x.loopAccept()
}

// loop to accept connection
func (x *TcpServer) loopAccept() error {
	// configuration values
	accMinDelay := x.acceptCnf.AcceptDelay.Value()
	accMaxDelay := x.acceptCnf.AcceptMaxDelay.Value()
	maxConnCount := x.acceptCnf.MaxConn
	maxRetry := x.acceptCnf.AcceptMaxRetry

	// net and error variables
	var conn net.Conn
	var err error
	var netErr net.Error
	var ok bool

	// retry and delay variables
	accDelay := time.Duration(0)
	retryCount := 0

	ulog.Info("TcpServer.loopAccept.start", zap.String("address", x.acceptCnf.Address))
	for {
		conn, err = x.ln.Accept()
		if err != nil {
			// check if it's a temporary error that we should retry
			netErr, ok = err.(net.Error)
			if !ok {
				return err
			}
			// nolint: staticcheck // SA1019: net.Error.Temporary is deprecated: Temporary is deprecated: see https://golang.org/doc/go1.17#net
			if !netErr.Temporary() {
				return err
			}

			// check retry count
			if maxRetry > 0 {
				retryCount++
				if retryCount >= maxRetry {
					return err
				}
			}

			// It's a temporary error, try to accept again
			if accDelay < accMinDelay {
				accDelay = accMinDelay
			} else if accDelay > accMaxDelay {
				accDelay = accMaxDelay
			} else {
				accDelay *= 2
			}
			time.Sleep(accDelay)
			continue
		}
		// reset
		accDelay = 0
		retryCount = 0

		curConnCount := x.connCount.Inc()
		if curConnCount > maxConnCount {
			err = conn.Close()
			x.connCount.Dec()
			ulog.Error("TcpServer.loopAccept.close.too.many", zap.Int32("currentConnCount", curConnCount), zap.Error(err))
		} else {
			connRunner := NewConnRunner(x.connReader, x.connSenderFactory(conn))
			connRunner.AddStartHook(x.connStartHook)
			connRunner.AddExitHook(x.connExitHook)
			connRunner.Start()
		}
	}
}

// connStartHook : when connection start
func (x *TcpServer) connStartHook(connSender IConnSender) {
	ulog.Info("connection.start", zap.Object("metaInfo", connSender.MetaInfo()), zap.Int32("currentConn", x.connCount.Load()))
	if x.startHooker != nil {
		x.startHooker(connSender)
	}
}

// connExitHook : when connection exit
func (x *TcpServer) connExitHook(connSender IConnSender) {
	curConnCount := x.connCount.Dec()
	ulog.Info("connection.exit", zap.Object("metaInfo", connSender.MetaInfo()), zap.Int32("currentConn", curConnCount))
	if x.exitHooker != nil {
		x.exitHooker(connSender)
	}
}
