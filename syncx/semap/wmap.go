package semap

import (
	"context"
	"github.com/pinealctx/neptune/remap"
)

//WideSemMap use SemMap array as wide map
type WideSemMap struct {
	ms       []*SemMap
	calKeyFn func(key interface{}) int
	rehash   *remap.ReMap
}

//NewWideSemMap new wide semaphore map
func NewWideSemMap(opts ...Option) SemMapper {
	var o = RangeOption(opts...)
	return newWideSemMap(o.rwRatio, o.prime, false)
}

//NewWideXHashSemMap new wide semaphore map
func NewWideXHashSemMap(opts ...Option) SemMapper {
	var o = RangeOption(opts...)
	return newWideSemMap(o.rwRatio, o.prime, true)
}

//newWideSemMap new wide semaphore map
func newWideSemMap(rwRatio int, prime uint64, useXHash bool) SemMapper {
	var w = &WideSemMap{}
	if prime > 0 {
		w.rehash = remap.NewReMap(remap.WithPrime(prime))
	} else {
		w.rehash = remap.NewReMap()
	}
	var numbs = w.rehash.Numbs()
	w.ms = make([]*SemMap, numbs)
	for i := uint64(0); i < numbs; i++ {
		w.ms[i] = newSemMap(rwRatio)
	}
	if useXHash {
		w.calKeyFn = w.rehash.XHashIndex
	} else {
		w.calKeyFn = w.rehash.SimpleIndex
	}
	return w
}

//AcquireRead acquire for read
func (s *WideSemMap) AcquireRead(ctx context.Context, key interface{}) (*Weighted, error) {
	return s.calculateKey(key).AcquireRead(ctx, key)
}

//ReleaseRead release read lock
func (s *WideSemMap) ReleaseRead(key interface{}, w *Weighted) {
	s.calculateKey(key).ReleaseRead(key, w)
}

//AcquireWrite acquire for write
func (s *WideSemMap) AcquireWrite(ctx context.Context, key interface{}) (*Weighted, error) {
	return s.calculateKey(key).AcquireWrite(ctx, key)
}

//ReleaseWrite release write lock
func (s *WideSemMap) ReleaseWrite(key interface{}, w *Weighted) {
	s.calculateKey(key).ReleaseWrite(key, w)
}

//calculate key
func (s *WideSemMap) calculateKey(key interface{}) *SemMap {
	var i = s.calKeyFn(key)
	return s.ms[i]
}
