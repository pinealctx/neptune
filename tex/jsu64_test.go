package tex

import (
	"encoding/json"
	"testing"

	"github.com/pinealctx/neptune/jsonx"
)

type TX struct {
	A int64    `json:"a"`
	B JsUInt64 `json:"b"`
}

func TestUInt64Json_MarshalJSON(t *testing.T) {
	var x = map[string]any{
		`a`: int64(9223372036854775806),
		`b`: JsUInt64(9223372036854775806),
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

	var y = TX{
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

func Test_JsUInt64UnmarshalJSON1(t *testing.T) {
	var j = `{"b":"12121212"}`
	var s TX
	var err = jsonx.JSONUnmarshal([]byte(j), &s)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(s)
}

func Test_JsUInt64UnmarshalJSON2(t *testing.T) {
	var j = `{"b":"-12121212"}`
	var s TX
	var err = jsonx.JSONUnmarshal([]byte(j), &s)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(s)
}
