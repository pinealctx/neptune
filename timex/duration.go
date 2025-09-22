package timex

import (
	"encoding/json"
	"fmt"
	"time"
)

// Duration is a wrapper of time.Duration, extend json marshal/unmarshal
type Duration time.Duration

// NewDuration creates a new Duration
func NewDuration(d time.Duration) Duration {
	return Duration(d)
}

// MarshalJSON implements the json.Marshaler interface.
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v any
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return fmt.Errorf("invalid duration: %v", v)
	}
}

// String : stringer
func (d Duration) String() string {
	return time.Duration(d).String()
}

// Value : get time.Duration value
func (d Duration) Value() time.Duration {
	return time.Duration(d)
}
