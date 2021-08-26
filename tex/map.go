package tex

import (
	"time"
)

//KeyInMap : Map中是否包含了key
func KeyInMap(m map[string]interface{}, k string) bool {
	_, ok := m[k]
	return ok
}

//MapVal2Bool : 将map中的value转换为bool
func MapVal2Bool(m map[string]interface{}, k string) bool {
	x, ok := m[k]
	if !ok {
		return false
	}
	return ToBool(x)
}

//MapVal2Int : 将map中的value转换为int
func MapVal2Int(m map[string]interface{}, k string) int {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToInt(x)
}

//MapVal2Int64 : 将map中的value转换为int64
func MapVal2Int64(m map[string]interface{}, k string) int64 {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToInt64(x)
}

//MapVal2JsInt64 : 将map中的value转换为 json int64
func MapVal2JsInt64(m map[string]interface{}, k string) JsInt64 {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToJsInt64(x)
}

//MapVal2Int32 : 将map中的value转换为int32
func MapVal2Int32(m map[string]interface{}, k string) int32 {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToInt32(x)
}

//MapVal2Float64 : 将map中的value转换为float64
func MapVal2Float64(m map[string]interface{}, k string) float64 {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToFloat64(x)
}

//MapVal2String : 将map中的value转换为string
func MapVal2String(m map[string]interface{}, k string) string {
	x, ok := m[k]
	if !ok {
		return ""
	}
	return ToString(x)
}

//MapVal2Bytes : 将map中的value转换为[]byte
func MapVal2Bytes(m map[string]interface{}, k string) []byte {
	x, ok := m[k]
	if !ok {
		return nil
	}
	return ToBytes(x)
}

//MapVal2StringList : 将map中的value转换为string list
func MapVal2StringList(m map[string]interface{}, k string) []string {
	x, ok := m[k]
	if !ok {
		return make([]string, 0)
	}
	return ToStringList(x)
}

//MapVal2Time : 将map中的value转换为time
func MapVal2Time(m map[string]interface{}, k string) (time.Time, bool) {
	x, ok := m[k]
	if !ok {
		return time.Time{}, false
	}
	return ToTime(x)
}
