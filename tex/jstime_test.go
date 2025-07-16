package tex

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pinealctx/neptune/jsonx"
)

type TJT struct {
	A JsUnixTime `json:"a"`
	B JsNanoTime `json:"b"`
}

func TestJsUnixTime_MarshalJSON(t *testing.T) {
	// test JsUnixTime json marshal and unmarshal
	var now = time.Now()
	var j = TJT{
		A: JsUnixTime(now),
		B: JsNanoTime(now),
	}
	t.Log(now)
	t.Log(time.Time(j.A), time.Time(j.B))
	var s, err = json.Marshal(j)
	if err != nil {
		panic(err)
	}
	t.Log(string(s))

	var j2 TJT
	err = json.Unmarshal(s, &j2)
	if err != nil {
		panic(err)
	}
	t.Log(time.Time(j2.A), time.Time(j2.B))
}

func TestJsUnixTime_FastMarshalJSON(t *testing.T) {
	// test JsUnixTime json marshal and unmarshal
	var now = time.Now()
	var j = TJT{
		A: JsUnixTime(now),
		B: JsNanoTime(now),
	}
	t.Log(now)
	t.Log(time.Time(j.A), time.Time(j.B))
	var s, err = jsonx.JSONFastMarshal(j)
	if err != nil {
		panic(err)
	}
	t.Log(string(s))

	var j2 TJT
	err = jsonx.JSONFastUnmarshal(s, &j2)
	if err != nil {
		panic(err)
	}
	t.Log(time.Time(j2.A), time.Time(j2.B))
}
