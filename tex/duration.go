package tex

import (
	"errors"
	"time"
)

var (
	ErrInvalidDurationJs = errors.New(`duration invalid string`)
)

type JsDuration struct {
	time.Duration
}

func NewJsDuration(d time.Duration) JsDuration {
	return JsDuration{Duration: d}
}

//MarshalJSON
//marshal json
func (i JsDuration) MarshalJSON() ([]byte, error) {
	var buf = []byte(i.String())
	var newBuf = make([]byte, 0, len(buf)+2)
	newBuf = append(newBuf, '"')
	newBuf = append(newBuf, buf...)
	newBuf = append(newBuf, '"')
	return newBuf, nil
}

//UnmarshalJSON
//unmarshal json
func (i *JsDuration) UnmarshalJSON(b []byte) error {
	var lb = len(b)
	if lb <= 2 {
		return ErrInvalidDurationJs
	}

	var strBuf = string(b[1 : lb-1])
	var dur, err = time.ParseDuration(strBuf)
	if err != nil {
		return err
	}
	i.Duration = dur
	return nil
}

//UnmarshalTOML
//unmarshal toml
func (i *JsDuration) UnmarshalTOML(v interface{}) error {
	var s, ok = v.(string)
	if !ok {
		return ErrInvalidDurationJs
	}
	var dur, err = time.ParseDuration(s)
	if err != nil {
		return err
	}
	i.Duration = dur
	return nil
}

//From : from duration
func (i *JsDuration) From(d time.Duration) {
	i.Duration = d
}
