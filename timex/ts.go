package timex

import "time"

func UnixNano(nano int64) time.Time {
	return time.Unix(nano/1e9, nano%1e9)
}
