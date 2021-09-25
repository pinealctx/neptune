package cache

import (
	"github.com/pinealctx/neptune/remap"
	"math"
	"sync"
	"testing"
	"time"
)

func Test100KLRU(t *testing.T) {
	test100KLRU(t, 2000000, 50, 10)
	t.Log("")
	test100KLRU(t, 2000000, 9, 1)
	t.Log("")

	test100KLRU(t, 20000, 50, 10)
	t.Log("")
	test100KLRU(t, 20000, 9, 1)
}

func test100KLRU(t *testing.T, size int64, rc int, wc int) {
	var lru = NewSingleLRUCache(size)
	test100KPassLRU(t, "one", lru, rc, wc)

	lru = NewWideXHashLRUCache(size)
	test100KPassLRU(t, "73", lru, rc, wc)

	lru = NewWideXHashLRUCache(size, remap.WithPrime(7))
	test100KPassLRU(t, "7", lru, rc, wc)

	lru = NewWideXHashLRUCache(size, remap.WithPrime(13))
	test100KPassLRU(t, "13", lru, rc, wc)

	lru = NewWideXHashLRUCache(size, remap.WithPrime(31))
	test100KPassLRU(t, "31", lru, rc, wc)

	lru = NewWideXHashLRUCache(size, remap.WithPrime(211))
	test100KPassLRU(t, "211", lru, rc, wc)

	lru = NewWideXHashLRUCache(size, remap.WithPrime(251))
	test100KPassLRU(t, "251", lru, rc, wc)

	lru = NewWideXHashLRUCache(size, remap.WithPrime(509))
	test100KPassLRU(t, "509", lru, rc, wc)
}

func test100KPassLRU(t *testing.T, name string, lru LRUFacade, rc int, wc int) {
	var count = 100000
	var wg sync.WaitGroup
	wg.Add(rc + wc)

	var minI64 = math.MinInt64
	var minI32 = math.MinInt32
	var minI16 = math.MinInt16
	var minI8 = math.MinInt8
	var n1 = -1
	var max uint64 = math.MaxUint64
	var mMax = max + 1
	var nMax = uint64(-1 * int64(max))

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
				lru.Set(j, _I(j))
			}
			for j := 0; j < c; j++ {
				lru.Set(us[j], _I(us[j]))
			}
		}()
	}

	for i := 0; i < rc; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				var v, ok = lru.Get(j)
				if ok {
					if j != int(v.(_I)) {
						panic(j)
					}
				}
			}
			for j := 0; j < c; j++ {
				var v, ok = lru.Get(us[j])
				if ok {
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

type _I uint64

func (i _I) Size() int {
	return 1
}
