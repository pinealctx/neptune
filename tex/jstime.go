package tex

import (
	"strconv"
	"time"
)

// JsUnixTime convert time to js timestamp(unix) with json
type JsUnixTime time.Time

// MarshalJSON marshal json
func (i JsUnixTime) MarshalJSON() ([]byte, error) {
	buf := []byte(strconv.FormatInt(time.Time(i).Unix(), 10))
	newBuf := make([]byte, 0, len(buf)+2)
	newBuf = append(newBuf, '"')
	newBuf = append(newBuf, buf...)
	newBuf = append(newBuf, '"')
	return newBuf, nil
}

// UnmarshalJSON unmarshal json
func (i *JsUnixTime) UnmarshalJSON(b []byte) error {
	lb := len(b)
	if lb <= 2 {
		return ErrInvalidInt64Js
	}

	strBuf := string(b[1 : lb-1])
	t, err := strconv.Atoi(strBuf)
	if err != nil {
		return err
	}
	*i = JsUnixTime(time.Unix(int64(t), 0).Local())
	return nil
}

// JsNanoTime convert time to js timestamp(unix nano) with json
type JsNanoTime time.Time

// MarshalJSON marshal json
func (i JsNanoTime) MarshalJSON() ([]byte, error) {
	buf := []byte(strconv.FormatInt(time.Time(i).UnixNano(), 10))
	newBuf := make([]byte, 0, len(buf)+2)
	newBuf = append(newBuf, '"')
	newBuf = append(newBuf, buf...)
	newBuf = append(newBuf, '"')
	return newBuf, nil
}

// UnmarshalJSON unmarshal json
func (i *JsNanoTime) UnmarshalJSON(b []byte) error {
	lb := len(b)
	if lb <= 2 {
		return ErrInvalidInt64Js
	}

	strBuf := string(b[1 : lb-1])
	t, err := strconv.Atoi(strBuf)
	if err != nil {
		return err
	}
	*i = JsNanoTime(time.Unix(0, int64(t)).Local())
	return nil
}
