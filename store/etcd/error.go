package etcd

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	OK Code = iota
	//Unavailable server unavailable
	Unavailable
	//Timeout rpc timeout
	Timeout
	//Cancelled rpc cancelled
	Cancelled
	//NodeExist node not exist
	NodeExist
	//NodeNotFound node not found
	NodeNotFound
	//BadVersion bad version check
	BadVersion
	//BadRsp bad response
	BadRsp
	//WatchFail watch fail
	WatchFail
	//WatchUnexpected watch unexpected event
	WatchUnexpected
	//WatchClosed watch closed
	WatchClosed
	//Unknown unknown error
	Unknown
)

var (
	//ErrEmptyRoot root path empty is dangerous
	ErrEmptyRoot = errors.New("etcd.root.path.empty")
	//ErrInvalidPath invalid path
	ErrInvalidPath = errors.New("invalid.etcd.path")
)

// Code error code
type Code int

// Error error define
type Error struct {
	code Code
	msg  string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s:%d", e.msg, e.code)
}

func (e *Error) Code() Code {
	return e.code
}

// ErrCode get error code of etcd
func ErrCode(err error) Code {
	if err == nil {
		return OK
	}
	var codeErr, ok = err.(*Error)
	if ok {
		return codeErr.code
	}
	return Unknown
}

// IsNotFoundErr is not found
func IsNotFoundErr(err error) bool {
	var code = ErrCode(err)
	return code == NodeNotFound
}

// convert error to consul error
func convertErr(nodePath string, err error) error {
	if err == nil {
		return nil
	}

	var typeErr, ok = err.(rpctypes.EtcdError)
	if ok {
		var code = typeErr.Code()
		switch code {
		case codes.NotFound:
			return genErr(nodePath, NodeNotFound)
		case codes.Unavailable:
			return genErr(nodePath, Unavailable)
		case codes.DeadlineExceeded:
			return genErr(nodePath, Timeout)
		}
		return &Error{
			code: Unknown,
			msg:  fmt.Sprintf("%s: etcd err: %+v", nodePath, err),
		}
	}
	var rpcStatus, rok = status.FromError(err)
	if rok {
		var code = rpcStatus.Code()
		switch code {
		case codes.NotFound:
			return genErr(nodePath, NodeNotFound)
		case codes.Canceled:
			return genErr(nodePath, Cancelled)
		case codes.DeadlineExceeded:
			return genErr(nodePath, Timeout)
		}
		return &Error{
			code: Unknown,
			msg:  fmt.Sprintf("%s: rpc err: %+v", nodePath, err),
		}
	}

	switch err {
	case context.Canceled:
		return genErr(nodePath, Cancelled)
	case context.DeadlineExceeded:
		return genErr(nodePath, Timeout)
	}
	return &Error{
		code: Unknown,
		msg:  fmt.Sprintf("%s: unknown err: %+v", nodePath, err),
	}
}

// make error
func genErr(node string, code Code) error {
	var msg string
	switch code {
	case OK:
		return nil
	case Unavailable:
		msg = fmt.Sprintf("%s: unavailable", node)
	case Timeout:
		msg = fmt.Sprintf("%s: timeout", node)
	case Cancelled:
		msg = fmt.Sprintf("%s: cancalled", node)
	case NodeExist:
		msg = fmt.Sprintf("%s: already exists", node)
	case NodeNotFound:
		msg = fmt.Sprintf("%s: not found", node)
	case BadVersion:
		msg = fmt.Sprintf("%s: bad version", node)
	case BadRsp:
		msg = fmt.Sprintf("%s: bad response", node)
	case WatchFail:
		msg = fmt.Sprintf("%s: watch fail", node)
	case WatchUnexpected:
		msg = fmt.Sprintf("%s: watch unexcepted", node)
	case WatchClosed:
		msg = fmt.Sprintf("%s: watch closed", node)
	default:
		msg = fmt.Sprintf("%s: undefined error", node)
	}
	return &Error{
		code: code,
		msg:  msg,
	}
}
