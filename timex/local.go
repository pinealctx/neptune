package timex

import (
	"time"
)

var (
	_localDiff time.Duration
)

//LocalDiff local time zone diff
func LocalDiff() time.Duration {
	return _localDiff
}

func init() {
	var t = time.Now()
	var _, diff = t.Zone()
	_localDiff = time.Duration(diff) * time.Second
}
