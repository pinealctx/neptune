package tex

import (
	"reflect"
	"testing"

	"github.com/vmihailenco/msgpack/v5"

	"github.com/pinealctx/neptune/jsonx"
)

type X struct {
	A int `msgpack:"a"`
}

func TestMapMarshal(t *testing.T) {
	var m = make(map[string]any)
	m["a"] = []string{"1", "2", "3"}
	m["b"] = []int64{1, 2, 3}
	m["c"] = []X{{A: 1}, {A: 2}}
	var buf, err = msgpack.Marshal(m)
	if err != nil {
		panic(err)
	}
	var x map[string]any
	err = msgpack.Unmarshal(buf, &x)
	if err != nil {
		panic(err)
	}
	for k, v := range x {
		t.Log("k", k)
		t.Log("v", v)
		t.Log("v type:", reflect.TypeOf(v))
	}

	buf, err = jsonx.JSONFastMarshal(m)
	if err != nil {
		panic(err)
	}
	var y map[string]any
	err = jsonx.JSONFastUnmarshal(buf, &y)
	if err != nil {
		panic(err)
	}
	for k, v := range y {
		t.Log("k", k)
		t.Log("v", v)
		t.Log("v type:", reflect.TypeOf(v))
	}
}

func TestMapArray(t *testing.T) {
	var m = make(map[string]any)
	m["a"] = []string{"1", "2", "3"}
	m["b"] = []int64{1, 2, 3}
	m["c"] = []X{{A: 1}, {A: 2}}
	var buf, err = msgpack.Marshal(m)
	if err != nil {
		panic(err)
	}
	var x map[string]any
	err = msgpack.Unmarshal(buf, &x)
	if err != nil {
		panic(err)
	}

	var (
		i64s []int64
		ss   []string
		ok   bool
	)

	ss, ok = MapVal2StringList(x, "a")
	if !ok {
		panic(ss)
	}
	t.Log(ss)
	t.Log(reflect.TypeOf(ss))

	i64s, ok = MapVal2Int64List(x, "b")
	if !ok {
		panic(i64s)
	}
	t.Log(i64s)
	t.Log(reflect.TypeOf(i64s))

	buf, err = jsonx.JSONFastMarshal(m)
	if err != nil {
		panic(err)
	}
	var y map[string]any
	err = jsonx.JSONFastUnmarshal(buf, &y)
	if err != nil {
		panic(err)
	}

	ss, ok = MapVal2StringList(y, "a")
	if !ok {
		panic(ss)
	}
	t.Log(ss)
	t.Log(reflect.TypeOf(ss))

	i64s, ok = MapVal2Int64List(y, "b")
	if !ok {
		panic(i64s)
	}
	t.Log(i64s)
	t.Log(reflect.TypeOf(i64s))
}
