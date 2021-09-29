package tex

import (
	"encoding/json"
	"github.com/pinealctx/neptune/jsonx"
	"testing"
	"time"
)

type DurationT struct {
	A int64      `json:"a"`
	B JsDuration `json:"b"`
}

func TestJsDurationMarshalJSON(t *testing.T) {
	var x = map[string]interface{}{
		`a`: int64(9223372036854775806),
		`b`: JsDuration{time.Second * 10},
	}
	s, err := jsonx.JSONMarshal(x)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	t.Log(`x1 is `, x)

	s, err = json.Marshal(x)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	t.Log(`x2 is `, x)

	var y = DurationT{
		9223372036854775806, JsDuration{time.Second * 10},
	}
	s, err = json.Marshal(y)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &y)
	t.Log(`y1 is `, y)

	err = jsonx.JSONUnmarshal(s, &x)
	t.Log(`x1 is `, x)

	s, err = jsonx.JSONMarshal(y)
	t.Log(string(s))
	t.Log(err)

	err = jsonx.JSONUnmarshal(s, &x)
	t.Log(`x2 is `, x)

}

func TestJsDurationUnmarshalJSON(t *testing.T) {
	var j = `{"b":"20s"}`
	var s DurationT
	var err = jsonx.JSONUnmarshal([]byte(j), &s)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(s)
	t.Log(int64(s.B.Duration))
}
