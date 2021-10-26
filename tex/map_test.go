package tex

import (
	"testing"
)

func TestMapClone(t *testing.T) {
	var m1 = make(map[string]interface{})
	m1["1"] = 1
	m1["2"] = 2
	var r1 = make(map[string]interface{})
	r1["1"] = 1
	r1["2"] = 1
	m1["r"] = r1

	var m2 = MapClone(m1)
	t.Log(m1)
	t.Log(m2)

	m2["1"] = -1
	m2["2"] = -2
	var x = m2["r"].(map[string]interface{})
	x["1"] = -1
	x["2"] = -2
	t.Log(m1)
	t.Log(m2)
}

func TestMapMerge(t *testing.T) {
	var m1 = make(map[string]interface{})
	var m2 = make(map[string]interface{})
	m1["0"] = 0
	m1["1"] = 1
	m1["2"] = 2
	var r1 = make(map[string]interface{})
	r1["1"] = 1
	r1["2"] = 1
	m1["r"] = r1

	m2["1"] = -1
	m2["2"] = -2
	var r2 = make(map[string]interface{})
	r2["1"] = -1
	r2["2"] = -1
	r2["3"] = -3
	r2["4"] = map[string]interface{}{"4": -4}
	m2["r"] = r2

	var m3 = MapMerge(m1, m2)
	t.Log(m1)
	t.Log(m2)
	t.Log(m3)
	t.Log("")

	m3 = MapMerge(nil, m2)
	t.Log(m1)
	t.Log(m2)
	t.Log(m3)
	t.Log("")

	m3 = MapMerge(m1, nil)
	t.Log(m1)
	t.Log(m2)
	t.Log(m3)
	t.Log("")
}
