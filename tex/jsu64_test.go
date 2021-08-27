package tex

import (
	"encoding/json"
	"testing"
)

type TX struct {
	A int64    `json:"a"`
	B JsUInt64 `json:"b"`
}

func TestUInt64Json_MarshalJSON(t *testing.T) {
	var x = map[string]interface{}{
		`a`: int64(9223372036854775806),
		`b`: JsUInt64(9223372036854775806),
	}
	s, err := JSONMarshal(x)
	t.Log(string(s))
	t.Log(err)

	err = JSONUnmarshal(s, &x)
	t.Log(`x1 is `, x)

	s, err = json.Marshal(x)
	t.Log(string(s))
	t.Log(err)

	err = JSONUnmarshal(s, &x)
	t.Log(`x2 is `, x)

	var y = TX{
		9223372036854775806, 9223372036854775806,
	}
	s, err = json.Marshal(y)
	t.Log(string(s))
	t.Log(err)

	err = JSONUnmarshal(s, &y)
	t.Log(`y1 is `, y)

	err = JSONUnmarshal(s, &x)
	t.Log(`x1 is `, x)

	s, err = JSONMarshal(y)
	t.Log(string(s))
	t.Log(err)

	err = JSONUnmarshal(s, &x)
	t.Log(`x2 is `, x)

}

func Test_JsUInt64UnmarshalJSON1(t *testing.T) {
	var j = `{"b":"12121212"}`
	var s TX
	var err = JSONUnmarshal([]byte(j), &s)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(s)
}

func Test_JsUInt64UnmarshalJSON2(t *testing.T) {
	var j = `{"b":"-12121212"}`
	var s TX
	var err = JSONUnmarshal([]byte(j), &s)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(s)
}
