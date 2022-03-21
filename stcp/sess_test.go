package stcp

import (
	"github.com/pinealctx/neptune/ulog"
	"go.uber.org/zap"
	"sync"
	"testing"
)

func TestSession_RemoteAddr(t *testing.T) {
	var s = NewSession(nil, nil)
	t.Log(s.RemoteAddr())
	t.Log(s.RemoteZap())
	ulog.Error("addr", s.RemoteZap())
	ulog.Error("sess", s.KeyZaps()...)
}

func TestSession_Recover(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	var b = NewSessionMgr(nil)
	var s = NewSession(b, nil)
	go func() {
		defer wg.Done()
		defer s.recovery()
		defer s.quit()
		panic("go routine 1 panic")
	}()
	go func() {
		defer wg.Done()
		defer s.recovery()
		defer s.quit()
		panic("go routine 2 panic")
	}()
	wg.Wait()
}

func TestSession_Recover1(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer recoverTest1()
		panic("go routine 1 panic")
	}()
	go func() {
		defer wg.Done()
		defer recoverTest1()
		panic("go routine 2 panic")
	}()
	wg.Wait()
}

func TestSession_Recover2(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer recoverTest2()
		panic("go routine 1 panic")
	}()
	go func() {
		defer wg.Done()
		defer recoverTest2()
		panic("go routine 2 panic")
	}()
	wg.Wait()
}

func recoverTest1() {
	var r = recover()
	if r != nil {
		//has panic
		ulog.Error("session.panic", zap.Any("panic", r),
			zap.Stack("stack"))
	}
}

func recoverTest2() {
	defer func() {
		var r = recover()
		if r != nil {
			//has panic
			ulog.Error("session.panic", zap.Any("panic", r),
				zap.Stack("stack"))
		}
	}()
	ulog.Error("recover")
}
