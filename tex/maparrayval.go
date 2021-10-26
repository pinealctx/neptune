package tex

//MapVal2BoolList : 将map中的value转换为bool list
func MapVal2BoolList(m map[string]interface{}, k string) ([]bool, bool) {
	var es, ok = mapValIsInterfaceArray(m, k)
	if !ok {
		return nil, false
	}
	var size = len(es)
	if size == 0 {
		return nil, true
	}
	var rs = make([]bool, size)
	for i := 0; i < size; i++ {
		rs[i] = ToBool(es[i])
	}
	return rs, true
}

//MapVal2IntList : 将map中的value转换为int list
func MapVal2IntList(m map[string]interface{}, k string) ([]int, bool) {
	var es, ok = mapValIsInterfaceArray(m, k)
	if !ok {
		return nil, false
	}
	var size = len(es)
	if size == 0 {
		return nil, true
	}
	var rs = make([]int, size)
	for i := 0; i < size; i++ {
		rs[i] = ToInt(es[i])
	}
	return rs, true
}

//MapVal2Int64List : 将map中的value转换为int64 list
func MapVal2Int64List(m map[string]interface{}, k string) ([]int64, bool) {
	var es, ok = mapValIsInterfaceArray(m, k)
	if !ok {
		return nil, false
	}
	var size = len(es)
	if size == 0 {
		return nil, true
	}
	var rs = make([]int64, size)
	for i := 0; i < size; i++ {
		rs[i] = ToInt64(es[i])
	}
	return rs, true
}

//MapVal2JsInt64List : 将map中的value转换为 json int64 list
func MapVal2JsInt64List(m map[string]interface{}, k string) ([]JsInt64, bool) {
	var es, ok = mapValIsInterfaceArray(m, k)
	if !ok {
		return nil, false
	}
	var size = len(es)
	if size == 0 {
		return nil, true
	}
	var rs = make([]JsInt64, size)
	for i := 0; i < size; i++ {
		rs[i] = ToJsInt64(es[i])
	}
	return rs, true
}

//MapVal2Int32List : 将map中的value转换为int32 list
func MapVal2Int32List(m map[string]interface{}, k string) ([]int32, bool) {
	var es, ok = mapValIsInterfaceArray(m, k)
	if !ok {
		return nil, false
	}
	var size = len(es)
	if size == 0 {
		return nil, true
	}
	var rs = make([]int32, size)
	for i := 0; i < size; i++ {
		rs[i] = ToInt32(es[i])
	}
	return rs, true
}

//MapVal2Float64List : 将map中的value转换为float64 list
func MapVal2Float64List(m map[string]interface{}, k string) ([]float64, bool) {
	var es, ok = mapValIsInterfaceArray(m, k)
	if !ok {
		return nil, false
	}
	var size = len(es)
	if size == 0 {
		return nil, true
	}
	var rs = make([]float64, size)
	for i := 0; i < size; i++ {
		rs[i] = ToFloat64(es[i])
	}
	return rs, true
}

//MapVal2StringList : 将map中的value转换为string list
func MapVal2StringList(m map[string]interface{}, k string) ([]string, bool) {
	var es, ok = mapValIsInterfaceArray(m, k)
	if !ok {
		return nil, false
	}
	var size = len(es)
	if size == 0 {
		return nil, true
	}
	var rs = make([]string, size)
	for i := 0; i < size; i++ {
		rs[i] = ToString(es[i])
	}
	return rs, true
}

//is interface array
func mapValIsInterfaceArray(m map[string]interface{}, k string) ([]interface{}, bool) {
	var x, ok = m[k]
	if !ok {
		return nil, false
	}
	var es []interface{}
	es, ok = x.([]interface{})
	if !ok {
		return nil, false
	}
	return es, true
}
