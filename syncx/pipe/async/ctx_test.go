package async

import (
	"context"
	"testing"
)

func TestCtxRun(t *testing.T) {
	var f1 = func(_ context.Context, _ int) (int, error) {
		return 10, nil
	}
	var aCtx = newCallCtx(context.Background(), f1, 1)
	//should panic
	aCtx.run()
	var r, err = aCtx.r()
	if err != nil {
		t.Error(err)
		return
	}
	var ri = r.(int)
	if ri != 10 {
		t.Fail()
		return
	}
}

func TestCtxRunPanic(t *testing.T) {
	defer func() {
		var e = recover()
		if e == nil {
			t.Error("should panic")
			t.Fail()
		}
		t.Log("panic catch:", e)
	}()
	var f1 = func(_ context.Context, _ int) (int, error) {
		return 10, nil
	}
	var aCtx = newCallCtx(context.Background(), f1, "1")
	//should panic
	aCtx.run()
	t.Log(aCtx.r())
}
