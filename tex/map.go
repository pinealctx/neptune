package tex

import (
	"time"
)

// KeyInMap : Map中是否包含了key
func KeyInMap(m map[string]interface{}, k string) bool {
	_, ok := m[k]
	return ok
}

// MapVal2Bool : 将map中的value转换为bool
func MapVal2Bool(m map[string]interface{}, k string) bool {
	x, ok := m[k]
	if !ok {
		return false
	}
	return ToBool(x)
}

// MapVal2Int : 将map中的value转换为int
func MapVal2Int(m map[string]interface{}, k string) int {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToInt(x)
}

// MapVal2Int64 : 将map中的value转换为int64
func MapVal2Int64(m map[string]interface{}, k string) int64 {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToInt64(x)
}

// MapVal2JsInt64 : 将map中的value转换为 json int64
func MapVal2JsInt64(m map[string]interface{}, k string) JsInt64 {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToJsInt64(x)
}

// MapVal2Int32 : 将map中的value转换为int32
func MapVal2Int32(m map[string]interface{}, k string) int32 {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToInt32(x)
}

// MapVal2Float64 : 将map中的value转换为float64
func MapVal2Float64(m map[string]interface{}, k string) float64 {
	x, ok := m[k]
	if !ok {
		return 0
	}
	return ToFloat64(x)
}

// MapVal2String : 将map中的value转换为string
func MapVal2String(m map[string]interface{}, k string) string {
	x, ok := m[k]
	if !ok {
		return ""
	}
	return ToString(x)
}

// MapVal2Bytes : 将map中的value转换为[]byte
func MapVal2Bytes(m map[string]interface{}, k string) []byte {
	x, ok := m[k]
	if !ok {
		return nil
	}
	return ToBytes(x)
}

/*//MapVal2StringList : 将map中的value转换为string list
func MapVal2StringList(m map[string]interface{}, k string) []string {
	x, ok := m[k]
	if !ok {
		return make([]string, 0)
	}
	return ToStringList(x)
}
*/

// MapVal2Time : 将map中的value转换为time
func MapVal2Time(m map[string]interface{}, k string) (time.Time, bool) {
	x, ok := m[k]
	if !ok {
		return time.Time{}, false
	}
	return ToTime(x)
}

// MapVal2Duration : 将map中的value转换为Duration
func MapVal2Duration(m map[string]interface{}, k string) (time.Duration, bool) {
	x, ok := m[k]
	if !ok {
		return 0, false
	}
	return ToDuration(x)
}

// MapMerge : merge map
// merge diff to base, return a new cloned map.
func MapMerge(base map[string]interface{}, diff map[string]interface{}) map[string]interface{} {
	var c = MapClone(base)
	mapMerge(c, diff)
	return c
}

// MapClone : clone a new map from old
// Attention: deeply clone only support map[string]interface{}, can not deeply clone pointer or other pointer like type.
// 只支持值为map[string]interface{}类型的深度克隆，不支持类似指针类型或切片类型的深度克隆
func MapClone(base map[string]interface{}) map[string]interface{} {
	var c = make(map[string]interface{})
	mapCopy(base, c)
	return c
}

// mapMerge : merge map
// base -- map[string]interface{}, cannot be nil, if nil will panic
// diff -- map[string]interface{}
func mapMerge(base map[string]interface{}, diff map[string]interface{}) {
	var (
		src interface{}
		mc  map[string]interface{}
		ok  bool
	)
	for k, v := range diff {
		src, ok = base[k]
		if !ok || src == nil {
			base[k] = v
			continue
		}
		mc, ok = src.(map[string]interface{})
		if !ok || mc == nil {
			base[k] = v
			continue
		}

		switch val := v.(type) {
		case map[string]interface{}:
			if val == nil {
				continue
			}
			mapMerge(mc, val)
		default:
			base[k] = val
		}
	}
}

// mapCopy : clone map
// input -- input map to clone.
// clone -- cloned map, clone map must be empty map but not nil.
func mapCopy(input map[string]interface{}, clone map[string]interface{}) {
	var (
		mc map[string]interface{}
		ok bool
	)
	for k, v := range input {
		mc, ok = v.(map[string]interface{})
		if ok {
			var nc = make(map[string]interface{})
			mapCopy(mc, nc)
			clone[k] = nc
		} else {
			clone[k] = v
		}
	}
}
