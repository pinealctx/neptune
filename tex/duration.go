package tex

import (
	"errors"
	"time"
)

var (
	ErrInvalidDuration = errors.New("invalid.duration.string")
)

type Duration time.Duration

func (i Duration) MarshalJSON() ([]byte, error) {
	var ii = (time.Duration)(i)
	var bytes = []byte(ii.String())
	var out = make([]byte, 0, len(bytes)+2)
	out = append(out, '"')
	out = append(out, bytes...)
	out = append(out, '"')
	return out, nil
}

func (i *Duration) UnmarshalJSON(b []byte) error {
	var l = len(b)
	if l <= 2 {
		return ErrInvalidDuration
	}
	var dur, err = time.ParseDuration(string(b[1 : l-1]))
	if err != nil {
		return err
	}
	*i = (Duration)(dur)
	return nil
}

func (i Duration) Duration() time.Duration {
	return time.Duration(i)
}

func (i *Duration) UnmarshalTOML(v interface{}) error {
	var s, ok = v.(string)
	if !ok {
		return ErrInvalidDuration
	}
	var dur, err = time.ParseDuration(s)
	if err != nil {
		return err
	}
	*i = (Duration)(dur)
	return nil
}
