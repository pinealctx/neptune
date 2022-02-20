package async

import (
	"context"
	"testing"
	"time"
)

func TestProcChan_Call(t *testing.T) {
	var runner = NewProcChan(WithName("test1"))
	go runner.Run()
	var size = DefaultQSize * 100
	var xs = make([]*_incTX, size)
	for i := 0; i < size; i++ {
		xs[i] = &_incTX{x: i}
	}
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
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
		_, _ = xs[i].Do(nil)
	}
	t2 = time.Now()
	dur = t2.Sub(t1)
	t.Log("sync use time:", dur, "average:", dur/time.Duration(size))
}

type _incChan struct {
	x  int
	rc chan int
}

func (i *_incChan) do() {
	i.x++
	i.rc <- i.x
}

func (i *_incChan) r() int {
	var x = <-i.rc
	return x
}

func (i *_incChan) doCtx(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
		i.x++
		close(i.rc)
	}
}

func (i *_incChan) rCtx(ctx context.Context) int {
	select {
	case <-ctx.Done():
		return 0
	case <-i.rc:
		return i.x
	}
}

func TestChan_Call(t *testing.T) {
	var size = DefaultQSize * 100
	var cc = make(chan *_incChan, DefaultQSize)
	go func() {
		var i *_incChan
		for {
			i = <-cc
			i.do()
		}
	}()
	var xs = make([]*_incChan, size)
	for i := 0; i < size; i++ {
		xs[i] = &_incChan{x: i, rc: make(chan int)}
	}

	var t1 = time.Now()
	for i := 0; i < size; i++ {
		cc <- xs[i]
		xs[i].r()
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("use time:", dur, "average:", dur/time.Duration(size))
}

func TestChan_CallCtx(t *testing.T) {
	var size = DefaultQSize * 100
	var cc = make(chan *_incChan, DefaultQSize)
	var ctx, cancel = context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	go func() {
		var i *_incChan
		for {
			i = <-cc
			i.doCtx(ctx)
		}
	}()
	var xs = make([]*_incChan, size)
	for i := 0; i < size; i++ {
		xs[i] = &_incChan{x: i, rc: make(chan int)}
	}

	var t1 = time.Now()
	for i := 0; i < size; i++ {
		cc <- xs[i]
		xs[i].rCtx(ctx)
	}
	var t2 = time.Now()
	var dur = t2.Sub(t1)
	t.Log("use time:", dur, "average:", dur/time.Duration(size))
}
