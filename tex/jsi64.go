package tex

import (
	"errors"
	"strconv"
)

var (
	jsonBrace         = []byte(`""`)
	ErrInvalidInt64Js = errors.New(`int64 invalid string`)
)

// JsInt64
// json could not support large number
// use string to replace the number if a field is int64
type JsInt64 int64

// MarshalJSON marshal json
func (i JsInt64) MarshalJSON() ([]byte, error) {
	buf := []byte(strconv.FormatInt(int64(i), 10))
	newBuf := make([]byte, 0, len(buf)+2)
	newBuf = append(newBuf, '"')
	newBuf = append(newBuf, buf...)
	newBuf = append(newBuf, '"')
	return newBuf, nil
}

// UnmarshalJSON unmarshal json
func (i *JsInt64) UnmarshalJSON(b []byte) error {
	lb := len(b)
	if lb == 0 {
		return ErrInvalidInt64Js
	}

	if b[0] == '"' && b[lb-1] == '"' {
		strBuf := string(b[1 : lb-1])
		if strBuf == "" {
			*i = 0
			return nil
		}
		t, err := strconv.Atoi(strBuf)
		if err != nil {
			return err
		}
		*i = JsInt64(t)
		return nil
	}

	strBuf := string(b)
	t, err := strconv.Atoi(strBuf)
	if err != nil {
		return err
	}
	*i = JsInt64(t)
	return nil
}
