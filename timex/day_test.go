package timex

import (
	"testing"
	"time"
)

func TestTodayBegin(t *testing.T) {
	t.Log(TodayBegin())
}

func TestDayDeltaBegin(t *testing.T) {
	t.Log(TodayDeltaBegin(-30))
	t.Log(TodayDeltaBegin(-1))
	t.Log(TodayDeltaBegin(0))
	t.Log(TodayDeltaBegin(1))
	t.Log(TodayDeltaBegin(30))
}

func TestDayDeltaBegins(t *testing.T) {
	t.Log(TodayDeltaDayBegins(-7, -3, -1, 0, 1, 3, 7))
	t.Log(TodayDeltaDayBegins(-3, -2, -1, 0, 1, 2, 3))
}

func TestDayDelta(t *testing.T) {
	var now = time.Now()
	t.Log(now)
	t.Log(DayDelta(now, -2))
	t.Log(DayDelta(now, -1))
	t.Log(DayDelta(now, 0))
	t.Log(DayDelta(now, 1))
	t.Log(DayDelta(now, 2))
}
