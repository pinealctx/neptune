package tex

import (
	"database/sql/driver"
	"strconv"
	"time"
)

// UnixNano2Time
// 纳秒时间戳时间
type UnixNano2Time time.Time

// Scan : sql scan
func (s *UnixNano2Time) Scan(value any) error {
	var ts int64
	switch v := value.(type) {
	case int32:
		ts = int64(v)
	case uint32:
		ts = int64(v)
	case int64:
		ts = v
	case uint64:
		ts = int64(v)
	case int:
		ts = int64(v)
	case uint:
		ts = int64(v)
	}
	*s = UnixNano2Time(time.Unix(0, ts))
	return nil
}

// Value : sql value
func (s UnixNano2Time) Value() (driver.Value, error) {
	return time.Time(s).UnixNano(), nil
}

// Unix2Time 秒时间戳时间
type Unix2Time time.Time

// Scan : sql scan
func (s *Unix2Time) Scan(value any) error {
	var ts int64
	switch v := value.(type) {
	case int32:
		ts = int64(v)
	case uint32:
		ts = int64(v)
	case int64:
		ts = v
	case uint64:
		ts = int64(v)
	case int:
		ts = int64(v)
	case uint:
		ts = int64(v)
	}
	*s = Unix2Time(time.Unix(ts, 0))
	return nil
}

// Value : sql value
func (s Unix2Time) Value() (driver.Value, error) {
	return time.Time(s).Unix(), nil
}

// UnixStamp 时间转成秒(只需要秒/Database中的datetime)
type UnixStamp int64

// Scan : sql scan
func (i *UnixStamp) Scan(value any) error {
	var t, ok = value.(time.Time)
	if ok {
		*i = UnixStamp(t.Unix())
	}
	return nil
}

// Value : sql value
func (i UnixStamp) Value() (driver.Value, error) {
	return time.Unix(int64(i), 0), nil
}

// MarshalJSON marshal json
func (i UnixStamp) MarshalJSON() ([]byte, error) {
	buf := []byte(strconv.FormatInt(int64(i), 10))
	newBuf := make([]byte, 0, len(buf)+2)
	newBuf = append(newBuf, '"')
	newBuf = append(newBuf, buf...)
	newBuf = append(newBuf, '"')
	return newBuf, nil
}

// UnmarshalJSON unmarshal json
func (i *UnixStamp) UnmarshalJSON(b []byte) error {
	lb := len(b)
	if lb <= 2 {
		return ErrInvalidInt64Js
	}

	strBuf := string(b[1 : lb-1])
	t, err := strconv.Atoi(strBuf)
	if err != nil {
		return err
	}
	*i = UnixStamp(t)
	return nil
}

// SQLTime2Unix convert timestamp to time.Time in SQL write mode.
// convert time.Time to timestamp when load data from database.
// implement sql.Driver Scan/Value
type SQLTime2Unix int64

// Scan : sql scan
func (i *SQLTime2Unix) Scan(value any) error {
	var t, ok = value.(time.Time)
	if ok {
		*i = SQLTime2Unix(t.Unix())
	}
	return nil
}

// Value : sql value
func (i SQLTime2Unix) Value() (driver.Value, error) {
	return time.Unix(int64(i), 0), nil
}
