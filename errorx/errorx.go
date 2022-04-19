// 相对于 github.com/pkg/errors 的改动
// 1)简化了stack信息（不包含完整路径，不然太长了）
// 2)加入了重复包含stack的判断，避免使用时的心智负担
// 3)把函数名字规整了一下，使其意义更明确
//
// 主要功能
//
// 名字上有 f 的是format的版本，有 WithStack 是带堆栈信息的，这里带的stack信息是从产出error的地方开始，最多往回追溯32层的调用栈信息
// 1)生成新的error
// New / Newf / NewWithStack / NewfWithStack
//
// 2)给 error 额外附带stack信息
// WithStack
//
// 3)包装error，并附带额外的信息
// Wrap / Wrapf / WrapWithStack / WrapfWithStack

package errorx

import (
	stderrors "errors"
	"fmt"
	"io"
)

func New(message string) error {
	return stderrors.New(message)
}

func Newf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args)
}

// NewWithStack returns an error with the supplied message.
// NewWithStack also records the stack trace at the point it was called.
func NewWithStack(message string) error {
	return &fundamental{
		msg:   message,
		stack: callers(),
	}
}

// NewfWithStack formats according to a format specifier and returns the string
// as a value that satisfies error.
// NewfWithStack also records the stack trace at the point it was called.
func NewfWithStack(format string, args ...interface{}) error {
	return &fundamental{
		msg:   fmt.Sprintf(format, args...),
		stack: callers(),
	}
}

// fundamental is an error that has a message and a stack, but no caller.
type fundamental struct {
	msg string
	*stack
}

func (f *fundamental) Error() string { return f.msg }

func (f *fundamental) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			io.WriteString(s, f.msg)
			f.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, f.msg)
	case 'q':
		fmt.Fprintf(s, "%q", f.msg)
	}
}

func hasBeenWithStack(err error) bool {
	for err != nil {
		switch err.(type) {
		case *withStack, *fundamental:
			return true
		}
		err = Unwrap(err)
	}
	return false
}



type withStack struct {
	error
	*stack
}


func (w *withStack) Cause() error { return w.error }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withStack) Unwrap() error { return w.error }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Cause())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())

	}
}

// WrapWithStack returns an error annotating err with a stack trace
// at the point WrapWithStack is called, and the supplied message.
// If err is nil, WrapWithStack returns nil.
func WrapWithStack(err error, message string) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   message,
	}
	if hasBeenWithStack(err) {
		return err
	}
	return &withStack{
		err,
		callers(),
	}
}

// WrapfWithStack returns an error annotating err with a stack trace
// at the point WrapfWithStack is called, and the format specifier.
// If err is nil, WrapfWithStack returns nil.
func WrapfWithStack(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
	if hasBeenWithStack(err) {
		return err
	}
	return &withStack{
		err,
		callers(),
	}
}

// Wrap annotates err with a new message.
// If err is nil, WithMessage returns nil.
func Wrap(err error, message string) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   message,
	}
}

// Wrapf annotates err with the format specifier.
// If err is nil, WithMessagef returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withMessage{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
	}
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	if hasBeenWithStack(err) {
		return err
	}
	return &withStack{
		err,
		callers(),
	}
}

type withMessage struct {
	cause error
	msg   string
}

func (w *withMessage) Error() string { return w.msg + ": " + w.cause.Error() }
func (w *withMessage) Cause() error  { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withMessage) Unwrap() error { return w.cause }

func (w *withMessage) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	type causer interface {
		Cause() error
	}

	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			break
		}
		err = cause.Cause()
	}
	return err
}

func GetFullStack(err error) string {
	for err != nil {
		switch errT := err.(type) {
		case *withStack:
			return errT.stack.getFullStackStr()
		case *fundamental:
			return errT.stack.getFullStackStr()
		}
		err = Unwrap(err)
	}
	return ""
}

