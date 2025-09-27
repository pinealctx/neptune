package async

import (
	"context"
	"reflect"
	"testing"
	"time"
)

type wrapFuncT struct {
}

func (w *wrapFuncT) invalidFn() {
}

func (w *wrapFuncT) validFn(_ context.Context, _ int) (int, error) {
	return 0, nil
}

type dWrapFuncT struct {
}

func (w dWrapFuncT) invalidFn() {
}

func (w dWrapFuncT) validFn(_ context.Context, _ int) (int, error) {
	return 0, nil
}

func incAdd(ctx context.Context, inc any) (int, error) {
	var call = reflect.ValueOf(inc)
	var params [1]reflect.Value

	params[0] = reflect.ValueOf(ctx)
	var rets = call.Call(params[:])
	var result = rets[0].Interface()
	var err error
	if !rets[1].IsNil() {
		// nolint : forcetypeassert // I know the type is exactly here
		err = rets[1].Interface().(error)
	}
	// nolint : forcetypeassert // I know the type is exactly here
	return result.(int), err
}

func TestValidateFunction(t *testing.T) {
	var f1 = func() {
	}
	var f2 = func(_ context.Context, _ int) (int, error) {
		return 0, nil
	}
	var w = &wrapFuncT{}
	var dw dWrapFuncT

	var ok bool
	var ty reflect.Type
	ty, ok = validateFn(f1)
	if ok {
		t.Fail()
		return
	}
	t.Log(ty)

	ty, ok = validateFn(f2)
	if !ok {
		t.Fail()
		return
	}
	t.Log(ty)

	ty, ok = validateFn(w.invalidFn)
	if ok {
		t.Fail()
		return
	}
	t.Log(ty)

	ty, ok = validateFn(w.validFn)
	if !ok {
		t.Fail()
		return
	}
	t.Log(ty)

	ty, ok = validateFn(dw.invalidFn)
	if ok {
		t.Fail()
		return
	}
	t.Log(ty)

	ty, ok = validateFn(dw.validFn)
	if !ok {
		t.Fail()
		return
	}
	t.Log(ty)

	var ty1 = reflect.TypeOf(f2)
	var ty2 = reflect.TypeOf(w.validFn)
	t.Log(ty1 == ty2)

	ty1 = reflect.TypeOf(dw.validFn)
	t.Log(ty1 == ty2)
}

func TestReflect(t *testing.T) {
	var size = DefaultQSize * 100
	var xs = make([]*_incTX, size)
	for i := 0; i < size; i++ {
		xs[i] = &_incTX{x: i}
	}
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var t1 = time.Now()
	for i := 0; i < size; i++ {
		var r, err = incAdd(ctx, xs[i].Do)
		if err != nil {
			panic(err)
		}
		if r != i+1 {
			panic("not.equals")
		}
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("use time:", dur, "average:", dur/time.Duration(size))
}
