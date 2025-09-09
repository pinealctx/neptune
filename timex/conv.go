package timex

import "time"

// NanosToMillis : Helper functions for time conversion
func NanosToMillis(nanos int64) int64 {
	// nolint:gosec // Ignore int64 to uint64 conversion
	return nanos / int64(time.Millisecond)
}

// MillisToNanos : Helper functions for time conversion
func MillisToNanos(millis int64) int64 {
	// nolint:gosec // Ignore int64 to uint64 conversion
	return millis * int64(time.Millisecond)
}
