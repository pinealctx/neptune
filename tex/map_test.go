package tex

import "testing"

func TestMapMerge(t *testing.T) {
	var m1 = make(map[string]interface{})
	var m2 = make(map[string]interface{})
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
	r2["4"] = map[string]interface{}{"4":-4}
	m2["r"] = r2

	MapMerge(m1, m2)
	t.Log(m1)
}
