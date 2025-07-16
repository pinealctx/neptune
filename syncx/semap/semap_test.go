package semap

import (
	"context"
	"fmt"
	"math"
	"sync"
	"testing"
	"time"

	"go.uber.org/atomic"
)

func Test1KLock(t *testing.T) {
	var s = NewSemMap(WithRwRatio(10))
	var count = 100000
	var wg sync.WaitGroup
	wg.Add(10)

	var t1 = time.Now()
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			var w, err = s.AcquireWrite(context.Background(), i)
			if err != nil {
				panic(err)
			}
			s.ReleaseWrite(i, w)
		}
	}()

	for i := 0; i < 9; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				var w, err = s.AcquireRead(context.Background(), j)
				if err != nil {
					panic(err)
				}
				s.ReleaseRead(j, w)
			}
		}()
	}
	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/time.Duration(count*10))
}

func Test100KLockWide(t *testing.T) {
	test100KLockRcWc(t, 100, 50, 10)
	t.Log("")
	test100KLockRcWc(t, 10, 9, 1)
	t.Log("")

	test100KLockRcWc(t, 100, 50, 10)
	t.Log("")
	test100KLockRcWc(t, 10, 9, 1)
}

func Test100KLockXHash(t *testing.T) {
	test100KXHashLockRcWc(t, 100, 50, 10)
	t.Log("")
	test100KXHashLockRcWc(t, 10, 9, 1)
	t.Log("")

	test100KXHashLockRcWc(t, 100, 50, 10)
	t.Log("")
	test100KXHashLockRcWc(t, 10, 9, 1)
}

func TestNotAcquired(_ *testing.T) {
	var sem = NewSemMap(WithRwRatio(5))
	testNotAcquired(sem, 200, 50)
	sem = NewWideSemMap(WithRwRatio(5))
	testNotAcquired(sem, 200, 50)
	sem = NewWideXHashSemMap(WithRwRatio(5))
	testNotAcquired(sem, 200, 50)
}

func testNotAcquired(sem SemMapper, rc int, wc int) {
	var wg sync.WaitGroup
	wg.Add(rc + wc)
	var count atomic.Int32

	for i := 0; i < wc; i++ {
		go func() {
			defer wg.Done()
			var ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			var w, err = sem.AcquireWrite(ctx, "key")
			if err != nil {
				fmt.Println("can.not.acquire.write", err)
				count.Add(1)
				return
			}
			defer sem.ReleaseWrite("key", w)
			time.Sleep(time.Second)
		}()
	}

	for i := 0; i < rc; i++ {
		go func() {
			defer wg.Done()
			var ctx, cancel = context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			var w, err = sem.AcquireRead(ctx, "key")
			if err != nil {
				fmt.Println("can.not.acquire.read", err)
				count.Add(1)
				return
			}
			defer sem.ReleaseRead("key", w)
			time.Sleep(time.Second)
		}()
	}
	wg.Wait()
	fmt.Println("failed count:", count.Load())
	//finally, sem can get
	var w, err = sem.AcquireWrite(context.Background(), "key")
	if err != nil {
		panic(err)
	}
	sem.ReleaseWrite("key", w)
}

func test100KLockRcWc(t *testing.T, rwRation int, rc int, wc int) {
	t.Helper()
	var opts = []OptionFn{
		WithRwRatio(rwRation),
	}

	var sem = NewSemMap(opts...)
	test100KLock(t, "one", sem, rc, wc)

	sem = NewWideSemMap(opts...)
	test100KLock(t, "73", sem, rc, wc)

	sem = NewWideSemMap(append(opts, WithPrime(7))...)
	test100KLock(t, "7", sem, rc, wc)

	sem = NewWideSemMap(append(opts, WithPrime(13))...)
	test100KLock(t, "13", sem, rc, wc)

	sem = NewWideSemMap(append(opts, WithPrime(31))...)
	test100KLock(t, "31", sem, rc, wc)

	sem = NewWideSemMap(append(opts, WithPrime(211))...)
	test100KLock(t, "211", sem, rc, wc)

	sem = NewWideSemMap(append(opts, WithPrime(251))...)
	test100KLock(t, "251", sem, rc, wc)

	sem = NewWideSemMap(append(opts, WithPrime(509))...)
	test100KLock(t, "509", sem, rc, wc)
}

func test100KXHashLockRcWc(t *testing.T, rwRation int, rc int, wc int) {
	t.Helper()
	var opts = []OptionFn{
		WithRwRatio(rwRation),
	}
	var sem = NewSemMap(opts...)
	test100KLock(t, "one", sem, rc, wc)

	sem = NewWideXHashSemMap(opts...)
	test100KLock(t, "73", sem, rc, wc)

	sem = NewWideXHashSemMap(append(opts, WithPrime(7))...)
	test100KLock(t, "7", sem, rc, wc)

	sem = NewWideXHashSemMap(append(opts, WithPrime(13))...)
	test100KLock(t, "13", sem, rc, wc)

	sem = NewWideXHashSemMap(append(opts, WithPrime(31))...)
	test100KLock(t, "31", sem, rc, wc)

	sem = NewWideXHashSemMap(append(opts, WithPrime(211))...)
	test100KLock(t, "211", sem, rc, wc)

	sem = NewWideXHashSemMap(append(opts, WithPrime(251))...)
	test100KLock(t, "251", sem, rc, wc)

	sem = NewWideXHashSemMap(append(opts, WithPrime(509))...)
	test100KLock(t, "509", sem, rc, wc)
}

func test100KLock(t *testing.T, name string, sem SemMapper, rc int, wc int) {
	t.Helper()
	var count = 100000
	var wg sync.WaitGroup
	wg.Add(rc + wc)

	var minI64 = math.MinInt64
	var minI32 = math.MinInt32
	var minI16 = math.MinInt16
	var minI8 = math.MinInt8
	var n1 = -1
	var maxV uint64 = math.MaxUint64
	var mMax = maxV + 1
	var nMax = uint64(-1 * int64(maxV))

	var us = []uint64{
		math.MaxUint64,
		uint64(minI64),
		math.MaxInt64,
		uint64(minI32),
		math.MaxInt32,
		uint64(minI16),
		math.MaxInt16,
		uint64(minI8),
		math.MaxInt8,
		uint64(n1),
		0,
		mMax,
		nMax,
	}
	var c = len(us)
	var t1 = time.Now()

	for i := 0; i < wc; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				var w, err = sem.AcquireWrite(context.Background(), j)
				if err != nil {
					panic(err)
				}
				sem.ReleaseWrite(j, w)
			}
			for j := 0; j < c; j++ {
				var w, err = sem.AcquireWrite(context.Background(), us[j])
				if err != nil {
					panic(err)
				}
				sem.ReleaseWrite(us[j], w)
			}
		}()
	}

	for i := 0; i < rc; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				var w, err = sem.AcquireRead(context.Background(), j)
				if err != nil {
					panic(err)
				}
				sem.ReleaseRead(j, w)
			}
			for j := 0; j < c; j++ {
				var w, err = sem.AcquireRead(context.Background(), us[j])
				if err != nil {
					panic(err)
				}
				sem.ReleaseRead(us[j], w)
			}
		}()
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log(name, "read go routine", rc, "write go routine", wc,
		"use time:", d, "average:", d/time.Duration(count*(rc+wc)))
}
