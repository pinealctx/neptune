package timex

import (
	"sync"
	"sync/atomic"
	"time"
)

var (
	_globalQuickTime *QuickTime
)

func init() {
	_globalQuickTime = NewQuickTime(time.Millisecond)
}

// Now returns the current time from the global QuickTime instance
func Now() time.Time {
	return _globalQuickTime.Now()
}

// StopGlobalQuickTime stops the global QuickTime instance
func StopGlobalQuickTime() {
	_globalQuickTime.Stop()
}

// QuickTime provides a high-performance current time retrieval mechanism
// by updating the current time at regular intervals using a background goroutine.
// This approach minimizes the overhead of frequent system calls to get the current time,
// making it suitable for high-throughput applications where time retrieval is frequent.
type QuickTime struct {
	currentTime atomic.Pointer[time.Time]
	ticker      *time.Ticker
	done        chan struct{}
	wg          sync.WaitGroup
	interval    time.Duration
	stopOnce    sync.Once
}

func NewQuickTime(interval time.Duration) *QuickTime {
	if interval <= 0 {
		interval = time.Millisecond
	}
	x := &QuickTime{
		ticker:   time.NewTicker(interval),
		done:     make(chan struct{}),
		interval: interval,
	}
	now := time.Now()
	x.currentTime.Store(&now)
	x.wg.Add(1)
	go x.run()
	return x
}

func (x *QuickTime) Now() time.Time {
	t := x.currentTime.Load()
	if t != nil {
		return *t
	}
	return time.Now()
}

func (x *QuickTime) Stop() {
	x.stopOnce.Do(func() {
		close(x.done)
		x.ticker.Stop()
		x.wg.Wait()
	})
}

func (x *QuickTime) run() {
	defer x.wg.Done()
	for {
		select {
		case <-x.ticker.C:
			now := time.Now()
			x.currentTime.Store(&now)
		case <-x.done:
			return
		}
	}
}
