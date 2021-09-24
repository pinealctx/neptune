package semap

import (
	"context"
)

//WideSemMap use SemMap array as wide map
type WideSemMap struct {
	ms       []*SemMap
	calKeyFn func(key interface{}) int
}

//NewWideSemMap new wide semaphore map
func NewWideSemMap(size int, rwRatio int) SemMapper {
	return newWideSemMap(size, rwRatio, SimpleIndex)
}

//NewWideXHashSemMap new wide semaphore map
func NewWideXHashSemMap(size int, rwRatio int) SemMapper {
	return newWideSemMap(size, rwRatio, XHashIndex)
}

//newWideSemMap new wide semaphore map
func newWideSemMap(size int, rwRatio int, fn func(key interface{}) int) SemMapper {
	var w = &WideSemMap{}
	w.ms = make([]*SemMap, numbs)
	for i := uint64(0); i < numbs; i++ {
		w.ms[i] = newSemMap(size, rwRatio)
	}
	w.calKeyFn = fn
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
