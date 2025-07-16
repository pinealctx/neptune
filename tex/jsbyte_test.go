package tex

import (
	"encoding/json"
	"testing"

	"github.com/pinealctx/neptune/jsonx"
)

type BT struct {
	A JsByte `json:"a"`
}

func TestByteJs_MarshalJSON(t *testing.T) {
	var x = map[string]any{
		`a`: JsByte{0, 128, 255},
	}
	s, err := jsonx.JSONMarshal(x)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	if err != nil {
		t.Error(err)
	}
	t.Log(`x1 is `, x)

	s, err = json.Marshal(x)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	if err != nil {
		t.Error(err)
	}
	t.Log(`x2 is `, x)

	var y = BT{A: []byte{0, 128, 255}}

	s, err = json.Marshal(y)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &y)
	if err != nil {
		t.Error(err)
	}
	t.Log(`y1 is `, y)

	err = jsonx.JSONUnmarshal(s, &x)
	if err != nil {
		t.Error(err)
	}
	t.Log(`x1 is `, x)

	s, err = jsonx.JSONMarshal(y)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	if err != nil {
		t.Error(err)
	}
	t.Log(`x2 is `, x)

	y = BT{}
	s, err = jsonx.JSONMarshal(y)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	if err != nil {
		t.Error(err)
	}
	t.Log(`x is `, x)
}

func TestJsByte_MarshalJSON(t *testing.T) {
	var x map[string]any
	var y = BT{}
	var z BT

	var s, err = jsonx.JSONMarshal(y)
	t.Log(string(s), err)

	err = jsonx.JSONUnmarshal(s, &x)
	t.Log(`x is `, x, err)

	err = jsonx.JSONUnmarshal(s, &z)
	t.Log(`z is `, z, err)

	y = BT{A: []byte{}}
	s, err = jsonx.JSONMarshal(y)
	t.Log(string(s), err)

	err = jsonx.JSONUnmarshal(s, &x)
	t.Log(`x is `, x, err)

	err = jsonx.JSONUnmarshal(s, &z)
	t.Log(`z is `, z, err)

	y = BT{A: []byte{0}}
	s, err = jsonx.JSONMarshal(y)
	t.Log(string(s), err)

	err = jsonx.JSONUnmarshal(s, &x)
	t.Log(`x is `, x, err)

	err = jsonx.JSONUnmarshal(s, &z)
	t.Log(`z is `, z, err)

	y = BT{A: []byte{1}}
	s, err = jsonx.JSONMarshal(y)
	t.Log(string(s), err)

	err = jsonx.JSONUnmarshal(s, &x)
	t.Log(`x is `, x, err)

	err = jsonx.JSONUnmarshal(s, &z)
	t.Log(`z is `, z, err)

	y = BT{A: []byte{0, 0}}
	s, err = jsonx.JSONMarshal(y)
	t.Log(string(s), err)

	err = jsonx.JSONUnmarshal(s, &x)
	t.Log(`x is `, x, err)

	err = jsonx.JSONUnmarshal(s, &z)
	t.Log(`z is `, z, err)
}
