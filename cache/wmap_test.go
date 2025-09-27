package cache

import (
	"math"
	"sync"
	"testing"
	"time"

	"github.com/pinealctx/neptune/remap"
)

func Test100KMap(t *testing.T) {
	test100KMap(t, 50, 10)
	t.Log("")
	test100KMap(t, 9, 1)
	t.Log("")

	test100KMap(t, 50, 10)
	t.Log("")
	test100KMap(t, 9, 1)
}

func test100KMap(t *testing.T, rc int, wc int) {
	t.Helper()
	var mp = NewSingleMap()
	test100KPassMap(t, "one", mp, rc, wc)

	mp = NewWideMap()
	test100KPassMap(t, "73", mp, rc, wc)

	mp = NewWideMap(remap.WithPrime(7))
	test100KPassMap(t, "7", mp, rc, wc)

	mp = NewWideMap(remap.WithPrime(13))
	test100KPassMap(t, "13", mp, rc, wc)

	mp = NewWideMap(remap.WithPrime(31))
	test100KPassMap(t, "31", mp, rc, wc)

	mp = NewWideMap(remap.WithPrime(211))
	test100KPassMap(t, "211", mp, rc, wc)

	mp = NewWideMap(remap.WithPrime(251))
	test100KPassMap(t, "251", mp, rc, wc)

	mp = NewWideMap(remap.WithPrime(509))
	test100KPassMap(t, "509", mp, rc, wc)
}

func test100KPassMap(t *testing.T, name string, mp MapFacade, rc int, wc int) {
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
				mp.Set(j, _I(j))
			}
			for j := 0; j < c; j++ {
				mp.Set(us[j], _I(us[j]))
			}
		}()
	}

	for i := 0; i < rc; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				var v, ok = mp.Get(j)
				if ok {
					// nolint : forcetypeassert // I know the type is exactly here
					if j != int(v.(_I)) {
						panic(j)
					}
				}
			}
			for j := 0; j < c; j++ {
				var v, ok = mp.Get(us[j])
				if ok {
					// nolint : forcetypeassert // I know the type is exactly here
					if us[j] != uint64(v.(_I)) {
						panic(us[j])
					}
				}
			}
		}()
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log(name, "read go routine", rc, "write go routine", wc,
		"use time:", d, "average:", d/time.Duration(count*(rc+wc)))
}
