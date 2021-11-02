package timex

import (
	"time"
)

//TodayBegin : today time begin
//获取今天起始时间
func TodayBegin() time.Time {
	var now = time.Now()
	return DayBegin(now)
}

//TodayDeltaBegin : one day begin with specific delta with today
//指定天数差异，可以获取早于今天n天的一天起始时间或晚于今天n天的一天起始时间
//其中由n为正或负来决定
//例如，昨天的起始时间可以由TodayDeltaBegin(-1)来获取，
//明天的起始时间可以由TodayDeltaBegin(1)来获取。
func TodayDeltaBegin(n int) time.Time {
	var now = time.Now()
	return DayDeltaBegin(now, n)
}

//TodayDeltaDayBegins : a series day begin list which have been specified delta day delta list
//与TodayDeltaBegin类似，不同的是通过指定一系列的天数差异来获取不同差值天数的起始时间
//例如，TodayDeltaDayBegins(-2, -1, 1, 2)可以分别获取前天/昨天/明天/后天各自的起始时间
func TodayDeltaDayBegins(diffs ...int) []time.Time {
	var now = time.Now()
	return DayDeltaBegins(now, diffs...)
}

//DayBegin : input a time to figure its day begin
//指定一个时间，获取此时间所在的那天起始时间
func DayBegin(at time.Time) time.Time {
	var year, month, day = at.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, at.Location())
}

//DayDeltaBegin : input a time to figure its delta day begin
//指定一个时间和天数差异，获取早于此时间n天或晚于此时间n天的那天的起始时间
//其中由n为正或负来决定
//例如，比输入时间早一天的起始时间可以由DayDeltaBegin(t, -1)来获取，
//比输入时间晚一天的起始时间可以由DayDeltaBegin(t, 1)来获取。
func DayDeltaBegin(at time.Time, n int) time.Time {
	var year, month, day = at.Date()
	return time.Date(year, month, day+n, 0, 0, 0, 0, at.Location())
}

//DayDeltaBegins : input a time to calculate a delta list day begins
//与DayDeltaBegin类似，不同的是通过指定一系列的天数差异来获取不同差值天数的起始时间
//例如，DayDeltaBegins(t, -2, -1, 1, 2)可以分别获取比t早2天/早1天/晚1天/晚2天的起始时间
func DayDeltaBegins(at time.Time, diffs ...int) []time.Time {
	var size = len(diffs)
	if size == 0 {
		return nil
	}
	var ts = make([]time.Time, size)
	var year, month, day = at.Date()

	for i := 0; i < size; i++ {
		ts[i] = time.Date(year, month, day+diffs[i], 0, 0, 0, 0, at.Location())
	}
	return ts
}

//DayDelta : input a time to calculate a delta day
func DayDelta(at time.Time, n int) time.Time {
	var year, month, day = at.Date()
	var hour, minute, second = at.Clock()
	var nsec = at.Nanosecond()
	return time.Date(year, month, day+n, hour, minute, second, nsec, at.Location())
}
