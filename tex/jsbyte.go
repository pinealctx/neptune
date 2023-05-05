package tex

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

var (
	ErrInvalidByteJs = errors.New(`byte.invalid.string`)
)

// JsByte
// json could not support readable []byte
// use string to replace the byte array
type JsByte []byte

// MarshalJSON
// marshal json
func (i JsByte) MarshalJSON() ([]byte, error) {
	buf := i.ToJS()
	newBuf := make([]byte, 0, len(buf)+2)
	newBuf = append(newBuf, '"')
	newBuf = append(newBuf, buf...)
	newBuf = append(newBuf, '"')
	return newBuf, nil
}

// UnmarshalJSON
// unmarshal json
func (i *JsByte) UnmarshalJSON(b []byte) error {
	lb := len(b)
	if lb < 2 {
		return ErrInvalidByteJs
	}

	strBuf := string(b[1 : lb-1])
	return i.FromString(strBuf)
}

// FromString
// from string
func (i *JsByte) FromString(strBuf string) error {
	if len(strBuf) == 0 {
		*i = nil
		return nil
	}

	var strNums = strings.Split(strBuf, "/")
	var size = len(strNums)
	if size == 0 {
		*i = nil
		return nil
	}

	*i = make(JsByte, size)
	for j := 0; j < size; j++ {
		t, err := strconv.Atoi(strNums[j])
		if err != nil {
			return err
		}
		(*i)[j] = byte(t)
	}
	return nil
}

// ToJS
// split byte to string
// use '/' split
func (i JsByte) ToJS() []byte {
	var builder = i.splitBuilder()
	return builder.Bytes()
}

// ToString
// split byte to byte
// use '/' split
func (i JsByte) ToString() string {
	var builder = i.splitBuilder()
	return builder.String()
}

// split builder
func (i JsByte) splitBuilder() *bytes.Buffer {
	var builder = bytes.NewBuffer(nil)
	var size = len(i)
	for j := 0; j < size; j++ {
		_, _ = builder.WriteString(strconv.Itoa(int(i[j])))
		if j != size-1 {
			_, _ = builder.WriteString("/")
		}
	}
	return builder
}
