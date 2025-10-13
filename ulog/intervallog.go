package ulog

import (
	"go.uber.org/atomic"
	"go.uber.org/zap"
	"sync"
	"time"
)

// LogFunc : log function type
type LogFunc func(msg string, fields ...zap.Field)

// Key : an interface for key abstraction
type Key[K any] interface {
	// Key extracts the key of type K
	Key() K
}

// CounterLogger is a generic struct that maps keys of any type K to integer counters.
type CounterLogger[K Key[any]] struct {
	debugMap *sync.Map
	infoMap  *sync.Map
	warnMap  *sync.Map
	interval *atomic.Int64
}

// NewCounterLogger creates a new CounterLogger instance.
func NewCounterLogger[K Key[any]](interval int64) CounterLogger[K] {
	return CounterLogger[K]{
		interval: atomic.NewInt64(interval),
		debugMap: &sync.Map{},
		infoMap:  &sync.Map{},
		warnMap:  &sync.Map{},
	}
}

// SetInterval sets the interval for the CounterLogger.
func (x CounterLogger[K]) SetInterval(interval int64) {
	x.interval.Store(interval)
}

// GetInterval gets the interval for the CounterLogger.
func (x CounterLogger[K]) GetInterval() int64 {
	return x.interval.Load()
}

// Debug logs a debug message if the count for the given key exceeds the set interval.
func (x CounterLogger[K]) Debug(k K, msg string, args ...zap.Field) {
	logWithKeyCount(Debug, x.debugMap, k, x.interval.Load(), msg, args...)
}

// Info logs an info message if the count for the given key exceeds the set interval.
func (x CounterLogger[K]) Info(k K, msg string, args ...zap.Field) {
	logWithKeyCount(Info, x.infoMap, k, x.interval.Load(), msg, args...)
}

// Warn logs a warn message if the count for the given key exceeds the set interval.
func (x CounterLogger[K]) Warn(k K, msg string, args ...zap.Field) {
	logWithKeyCount(Warn, x.warnMap, k, x.interval.Load(), msg, args...)
}

// Debug logs a debug message if the count for the given key exceeds the set interval.
func logWithKeyCount[K any](fn LogFunc, kMap *sync.Map, k Key[K], interval int64, msg string, args ...zap.Field) {
	// get from kMap
	key := k.Key()
	cntI, ok := kMap.Load(key)
	if !ok {
		// not found, store and log
		kMap.Store(key, int64(1))
		fn(msg, args...)
		return
	}
	cnt, ok := cntI.(int64)
	if !ok {
		// type assertion failed, store and log
		kMap.Store(key, int64(1))
		fn(msg, args...)
		return
	}
	cnt++
	if cnt >= interval {
		// update and log
		kMap.Store(key, int64(1))
		fn(msg, args...)
	} else {
		// update only
		kMap.Store(key, cnt)
	}
}

// TimeKey : an interface for time abstraction
type TimeKey[K any] interface {
	// Key extracts the key of type K
	Key() K
	// ExtractTime extracts the time.Time value
	ExtractTime() time.Time
}

// TimeLogger is a generic struct that maps keys of any type K to time.Time values.
type TimeLogger[K TimeKey[any]] struct {
	debugMap *sync.Map
	infoMap  *sync.Map
	warnMap  *sync.Map
	duration *atomic.Duration
}

// NewTimeLogger creates a new TimeLogger instance.
func NewTimeLogger[K TimeKey[any]](dur time.Duration) TimeLogger[K] {
	return TimeLogger[K]{
		duration: atomic.NewDuration(dur),
		debugMap: &sync.Map{},
		infoMap:  &sync.Map{},
		warnMap:  &sync.Map{},
	}
}

// SetDuration sets the duration for the TimeLogger.
func (x TimeLogger[K]) SetDuration(dur time.Duration) {
	x.duration.Store(dur)
}

// GetDuration gets the duration for the TimeLogger.
func (x TimeLogger[K]) GetDuration() time.Duration {
	return x.duration.Load()
}

// Debug logs a debug message if the time since the last log for the given key exceeds the set duration.
func (x TimeLogger[K]) Debug(tk TimeKey[K], msg string, args ...zap.Field) {
	logWithTimeKey(Debug, x.debugMap, tk, x.duration.Load(), msg, args...)
}

// Info logs an info message if the time since the last log for the given key exceeds the set duration.
func (x TimeLogger[K]) Info(tk TimeKey[K], msg string, args ...zap.Field) {
	logWithTimeKey(Info, x.infoMap, tk, x.duration.Load(), msg, args...)
}

// Warn logs a warn message if the time since the last log for the given key exceeds the set duration.
func (x TimeLogger[K]) Warn(tk TimeKey[K], msg string, args ...zap.Field) {
	logWithTimeKey(Warn, x.warnMap, tk, x.duration.Load(), msg, args...)
}

// Debug logs a debug message if the time since the last log for the given key exceeds the set duration.
func logWithTimeKey[K any](fn LogFunc, kMap *sync.Map, k TimeKey[K], duration time.Duration, msg string, args ...zap.Field) {
	// get from kMap
	key := k.Key()
	prvTimeI, ok := kMap.Load(key)
	if !ok {
		// not found, store and log
		kMap.Store(key, k.ExtractTime())
		fn(msg, args...)
		return
	}
	prvTime, ok := prvTimeI.(time.Time)
	if !ok {
		// type assertion failed, store and log
		kMap.Store(key, k.ExtractTime())
		fn(msg, args...)
		return
	}
	logTime := k.ExtractTime()
	if logTime.Sub(prvTime) >= duration {
		// update and log
		kMap.Store(key, logTime)
		fn(msg, args...)
	}
}
