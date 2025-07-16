package tex

import (
	"encoding/json"
	"testing"

	"github.com/pinealctx/neptune/jsonx"
)

type T struct {
	A int64   `json:"a"`
	B JsInt64 `json:"b"`
}

func TestInt64Json_MarshalJSON(t *testing.T) {
	var x = map[string]any{
		`a`: int64(9223372036854775806),
		`b`: JsInt64(9223372036854775806),
	}
	s, err := jsonx.JSONMarshal(x)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(`x1 is `, x)

	s, err = json.Marshal(x)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(`x2 is `, x)

	var y = T{
		9223372036854775806, 9223372036854775806,
	}
	s, err = json.Marshal(y)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &y)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(`y1 is `, y)

	err = jsonx.JSONUnmarshal(s, &x)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(`x1 is `, x)

	s, err = jsonx.JSONMarshal(y)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(`x2 is `, x)

}

func Test_JsInt64UnmarshalJSON(t *testing.T) {
	var j = `{"b":"12121212"}`
	var s T
	var err = jsonx.JSONUnmarshal([]byte(j), &s)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(s)
}
