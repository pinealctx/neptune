package timex

import (
	"testing"
	"time"
)

func TestTimerStop1(t *testing.T) {
	procTest(t, testTimerStop1)
}

func testTimerStop1(t *testing.T) {
	t.Helper()
	var x = NewTimer(time.Second * 2)
	var tt, ok = <-x.C
	t.Log(tt, ok)
	x.Stop()
	x.Stop()
}

func TestTimerStop2(t *testing.T) {
	procTest(t, testTimerStop2)
}

func testTimerStop2(t *testing.T) {
	t.Helper()
	var x = NewTimer(time.Second * 2)
	x.Stop()
	x.Stop()
}

func TestTimerReset1(t *testing.T) {
	procTest(t, testTimerReset1)
}

func testTimerReset1(t *testing.T) {
	t.Helper()
	var x = NewTimer(time.Second * 2)
	var tt, ok = <-x.C
	t.Log(tt, ok)
	x.Reset(time.Second * 3)
	tt, ok = <-x.C
	t.Log(tt, ok)
	x.Stop()
}

func TestTimerReset2(t *testing.T) {
	procTest(t, testTimerReset2)
}

func testTimerReset2(t *testing.T) {
	t.Helper()
	var x = NewTimer(time.Second * 2)
	var tt, ok = <-x.C
	t.Log(tt, ok)
	x.Reset(time.Second * 3)
	x.Reset(time.Second * 2)
	x.Reset(time.Second * 1)
	tt, ok = <-x.C
	t.Log(tt, ok)
	x.Stop()
}

func TestTimerReset3(t *testing.T) {
	procTest(t, testTimerReset3)
}

func testTimerReset3(t *testing.T) {
	t.Helper()
	var x = NewTimer(time.Second * 2)
	x.Stop()
	x.Reset(time.Second * 3)
	x.Reset(time.Second * 2)
	x.Reset(time.Second * 1)
	var tt, ok = <-x.C
	t.Log(tt, ok)
	x.Stop()
}

func procTest(t *testing.T, f func(*testing.T)) {
	t.Helper()
	t.Log(time.Now(), "begin")
	f(t)
	t.Log(time.Now(), "end")
}
