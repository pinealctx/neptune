package timex

import "time"

type Timer struct {
	*time.Timer
}

// NewTimer creates a new Timer that will send
// the current time on its channel after at least duration d.
func NewTimer(d time.Duration) Timer {
	return Timer{time.NewTimer(d)}
}

// Stop : stop timer
// In time.Timer, Stop can not ensure drain the timer channel.
// Caller should use external code.
//
//	if !t.Stop() {
//		<-t.C
//	}
//
// Actually, it the t.C be read from others, the outside caller should avoid drain the timer channel again.
// For instance, time handler be called so that t.C is read then it's empty.
//
// This code will be blocked here forever because t.C is empty now.
//
//		if !t.Stop() {
//			<-t.C
//		}
//
//	 The function use "select with default" to avoid this case.
func (x Timer) Stop() {
	if !x.Timer.Stop() {
		select {
		case <-x.C:
		default:
		}
	}
}

// Reset changes the timer to expire after duration d.
// This function first call "Stop" to make sure t.C drained
func (x Timer) Reset(d time.Duration) {
	x.Stop()
	x.Timer.Reset(d)
}
