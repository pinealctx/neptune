package tex

import (
	"errors"
	"strconv"
)

var (
	ErrInvalidUInt64Js = errors.New(`uint64 invalid string`)
)

//JsUInt64
//json could not support large number
//use string to replace the number if a field is uint64
type JsUInt64 uint64

//MarshalJSON
//marshal json
func (i JsUInt64) MarshalJSON() ([]byte, error) {
	buf := []byte(strconv.FormatUint(uint64(i), 10))
	newBuf := make([]byte, 0, len(buf)+2)
	newBuf = append(newBuf, '"')
	newBuf = append(newBuf, buf...)
	newBuf = append(newBuf, '"')
	return newBuf, nil
}

//UnmarshalJSON
//unmarshal json
func (i *JsUInt64) UnmarshalJSON(b []byte) error {
	lb := len(b)
	if lb <= 2 {
		return ErrInvalidUInt64Js
	}

	strBuf := string(b[1 : lb-1])
	t, err := strconv.ParseUint(strBuf, 10, 64)
	if err != nil {
		return err
	}
	*i = JsUInt64(t)
	return nil
}
