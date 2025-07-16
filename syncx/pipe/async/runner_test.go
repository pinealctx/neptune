package async

import (
	"context"
	"testing"
	"time"
)

type _incT struct {
	x int
}

func (i *_incT) add(_ context.Context, _ int) (int, error) {
	i.x++
	return i.x, nil
}

type _incTX struct {
	x int
}

func (i *_incTX) Do(_ context.Context) (any, error) {
	i.x++
	return i.x, nil
}

func TestRunnerQ_Call(t *testing.T) {
	var runner = NewRunnerQ(WithName("test1"))
	go runner.Run()
	var size = DefaultQSize * 100
	var xs = make([]_incT, size)
	for i := 0; i < size; i++ {
		xs[i].x = i
	}
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var t1 = time.Now()
	for i := 0; i < size; i++ {
		var r, err = runner.AsyncCall(ctx, xs[i].add, 0)
		if err != nil {
			panic(err)
		}
		var ix = r.(int)
		if ix != i+1 {
			panic("not.equals")
		}
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("use time:", dur, "average:", dur/time.Duration(size))

	t1 = time.Now()
	for i := 0; i < size; i++ {
		_, _ = xs[i].add(context.TODO(), 0)
	}
	t2 = time.Now()
	dur = t2.Sub(t1)
	t.Log("sync use time:", dur, "average:", dur/time.Duration(size))
}

func TestRunnerQ_Delegate(t *testing.T) {
	var runner = NewRunnerQ(WithName("test1"))
	go runner.Run()
	var size = DefaultQSize * 100
	var xs = make([]*_incTX, size)
	for i := 0; i < size; i++ {
		xs[i] = &_incTX{x: i}
	}
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var t1 = time.Now()
	for i := 0; i < size; i++ {
		var r, err = runner.AsyncDelegate(ctx, xs[i].Do)
		if err != nil {
			panic(err)
		}
		var ix = r.(int)
		if ix != i+1 {
			panic("not.equals")
		}
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("use time:", dur, "average:", dur/time.Duration(size))

	t1 = time.Now()
	for i := 0; i < size; i++ {
		_, _ = xs[i].Do(context.TODO())
	}
	t2 = time.Now()
	dur = t2.Sub(t1)
	t.Log("sync use time:", dur, "average:", dur/time.Duration(size))
}

func TestRunnerQ_Proc(t *testing.T) {
	var runner = NewRunnerQ(WithName("test1"))
	go runner.Run()
	var size = DefaultQSize * 100
	var xs = make([]*_incTX, size)
	for i := 0; i < size; i++ {
		xs[i] = &_incTX{x: i}
	}
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	var t1 = time.Now()
	for i := 0; i < size; i++ {
		var r, err = runner.AsyncProc(ctx, xs[i])
		if err != nil {
			panic(err)
		}
		var ix = r.(int)
		if ix != i+1 {
			panic("not.equals")
		}
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("use time:", dur, "average:", dur/time.Duration(size))

	t1 = time.Now()
	for i := 0; i < size; i++ {
		_, _ = xs[i].Do(context.TODO())
	}
	t2 = time.Now()
	dur = t2.Sub(t1)
	t.Log("sync use time:", dur, "average:", dur/time.Duration(size))
}
