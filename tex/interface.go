package tex

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ToBytes 将基本类型转换为[]byte
func ToBytes(v interface{}) []byte {
	switch vt := v.(type) {
	case string:
		return []byte(vt)
	case []byte:
		return vt
	case JsByte:
		return vt
	default:
		return nil
	}
}

// ToString 将基本类型转换为字符串
func ToString(v interface{}) string {

	//检查整数
	var iv, bInt = tryNum2Int(v)
	if bInt {
		return strconv.Itoa(iv)
	}

	//检查字符串
	bString := true
	var sv string
	switch val := v.(type) {
	case string:
		sv = val
	case []byte:
		sv = string(val)
	default:
		bString = false
	}
	if bString {
		return sv
	}

	//检查浮点数
	bFloat := true
	var fv float64
	switch val := v.(type) {
	case float32:
		fv = float64(val)
	case float64:
		fv = val
	default:
		bFloat = false
	}
	if bFloat {
		return strconv.FormatFloat(fv, 'g', 3, 64)
	}

	return fmt.Sprint(v)
}

// ToStringList 将基本类型转换为字符串list
func ToStringList(v interface{}) []string {
	switch val := v.(type) {
	case []string:
		return val
	case [][]byte:
		rv := make([]string, len(val))
		for i := range val {
			rv[i] = string(val[i])
		}
		return rv
	case []interface{}:
		rv := make([]string, len(val))
		for i := range val {
			rv[i] = ToString(val[i])
		}
		return rv
	default:
		return make([]string, 0)
	}
}

// ToBool 将基本类型转换为bool
func ToBool(v interface{}) bool {
	switch val := v.(type) {
	case bool:
		return val
	case int8:
		return val != 0
	case uint8:
		return val != 0
	case int16:
		return val != 0
	case uint16:
		return val != 0
	case int32:
		return val != 0
	case uint32:
		return val != 0
	case int64:
		return val != 0
	case uint64:
		return val != 0
	case int:
		return val != 0
	case uint:
		return val != 0
	case JsInt64:
		return val != 0
	case JsUInt64:
		return val != 0
	case time.Duration:
		return val != 0
	case Duration:
		return val != 0
	case *Duration:
		return *val != 0
	case float32:
		return val != 0
	case float64:
		return val != 0
	case string:
		return strings.ToUpper(val) == `TRUE`
	}
	return false
}

// ToInt 将基本类型转换为整型
func ToInt(v interface{}) int {
	var iv int
	switch val := v.(type) {
	case int8:
		iv = int(val)
	case uint8:
		iv = int(val)
	case int16:
		iv = int(val)
	case uint16:
		iv = int(val)
	case int32:
		iv = int(val)
	case uint32:
		iv = int(val)
	case int64:
		iv = int(val)
	case uint64:
		iv = int(val)
	case int:
		iv = val
	case uint:
		iv = int(val)
	case JsInt64:
		iv = int(val)
	case JsUInt64:
		iv = int(val)
	case time.Duration:
		iv = int(val)
	case Duration:
		iv = int(val)
	case *Duration:
		iv = int(*val)
	case float32:
		iv = int(val)
	case float64:
		iv = int(val)
	case string:
		iv, _ = strconv.Atoi(val)
	default:
		iv = 0
	}
	return iv
}

func ToUInt(v interface{}) uint {
	return uint(ToInt(v))
}

// ToInt32 将基本类型转换为32位整型
func ToInt32(v interface{}) int32 {
	var iv int32
	switch val := v.(type) {
	case int8:
		iv = int32(val)
	case uint8:
		iv = int32(val)
	case int16:
		iv = int32(val)
	case uint16:
		iv = int32(val)
	case int32:
		iv = val
	case uint32:
		iv = int32(val)
	case int64:
		iv = int32(val)
	case uint64:
		iv = int32(val)
	case int:
		iv = int32(val)
	case uint:
		iv = int32(val)
	case JsInt64:
		iv = int32(val)
	case JsUInt64:
		iv = int32(val)
	case time.Duration:
		iv = int32(val)
	case Duration:
		iv = int32(val)
	case *Duration:
		iv = int32(*val)
	case float32:
		iv = int32(val)
	case float64:
		iv = int32(val)
	case string:
		/* #nosec */
		iiv, _ := strconv.Atoi(val)
		/* #nosec */
		//WTF, convert int to uint32, it works.
		iv = int32(iiv)
	default:
		iv = 0
	}
	return iv
}

// ToInt32 将基本类型转换为32位无符号整型
func ToUInt32(v interface{}) uint32 {
	var uv uint32
	switch val := v.(type) {
	case int8:
		uv = uint32(val)
	case uint8:
		uv = uint32(val)
	case int16:
		uv = uint32(val)
	case uint16:
		uv = uint32(val)
	case int32:
		uv = uint32(val)
	case uint32:
		uv = val
	case int64:
		uv = uint32(val)
	case uint64:
		uv = uint32(val)
	case int:
		uv = uint32(val)
	case uint:
		uv = uint32(val)
	case JsInt64:
		uv = uint32(val)
	case JsUInt64:
		uv = uint32(val)
	case time.Duration:
		uv = uint32(val)
	case Duration:
		uv = uint32(val)
	case *Duration:
		uv = uint32(*val)
	case float32:
		uv = uint32(val)
	case float64:
		uv = uint32(val)
	case string:
		/* #nosec */
		iiv, _ := strconv.Atoi(val)
		uv = uint32(iiv)
	default:
		uv = 0
	}
	return uv
}

// ToInt64 将基本类型转换为64位整型
func ToInt64(v interface{}) int64 {
	var iv int64
	switch val := v.(type) {
	case int8:
		iv = int64(val)
	case uint8:
		iv = int64(val)
	case int16:
		iv = int64(val)
	case uint16:
		iv = int64(val)
	case int32:
		iv = int64(val)
	case uint32:
		iv = int64(val)
	case int64:
		iv = val
	case uint64:
		iv = int64(val)
	case int:
		iv = int64(val)
	case uint:
		iv = int64(val)
	case JsInt64:
		iv = int64(val)
	case JsUInt64:
		iv = int64(val)
	case time.Duration:
		iv = int64(val)
	case Duration:
		iv = int64(val)
	case *Duration:
		iv = int64(*val)
	case float32:
		iv = int64(val)
	case float64:
		iv = int64(val)
	case string:
		/* #nosec */
		iiv, _ := strconv.Atoi(val)
		iv = int64(iiv)
	default:
		iv = 0
	}
	return iv
}

func ToUInt64(v interface{}) uint64 {
	return uint64(ToInt64(v))
}

// ToJsInt64 将基本类型转换为json int64整型
func ToJsInt64(v interface{}) JsInt64 {
	switch val := v.(type) {
	case JsInt64:
		return val
	case JsUInt64:
		return JsInt64(val)
	default:
		var i64 = ToInt64(v)
		return JsInt64(i64)
	}
}

// ToJsUInt64 将基本类型转换为json uint64整型
func ToJsUInt64(v interface{}) JsUInt64 {
	switch val := v.(type) {
	case JsUInt64:
		return val
	case JsInt64:
		return JsUInt64(val)
	default:
		var i64 = ToInt64(v)
		return JsUInt64(i64)
	}
}

// ToFloat64 将基本类型转换为浮点型
func ToFloat64(v interface{}) float64 {
	var iv float64
	switch val := v.(type) {
	case int8:
		iv = float64(val)
	case uint8:
		iv = float64(val)
	case int16:
		iv = float64(val)
	case uint16:
		iv = float64(val)
	case int32:
		iv = float64(val)
	case uint32:
		iv = float64(val)
	case int64:
		iv = float64(val)
	case uint64:
		iv = float64(val)
	case int:
		iv = float64(val)
	case uint:
		iv = float64(val)
	case JsInt64:
		iv = float64(val)
	case JsUInt64:
		iv = float64(val)
	case time.Duration:
		iv = float64(val)
	case Duration:
		iv = float64(val)
	case *Duration:
		iv = float64(*val)
	case float32:
		iv = float64(val)
	case float64:
		iv = val
	case string:
		iv, _ = strconv.ParseFloat(val, 64)
	default:
		iv = 0
	}
	return iv
}

// ToTime 转换为time.Time
func ToTime(v interface{}) (time.Time, bool) {
	t, ok := v.(time.Time)
	return t, ok
}

// ToDuration 转换为time.Duration
func ToDuration(v interface{}) (time.Duration, bool) {
	switch val := v.(type) {
	case time.Duration:
		return val, true
	case Duration:
		return time.Duration(val), true
	case *Duration:
		return time.Duration(*val), true
	case float32:
		return time.Duration(val), true
	case float64:
		return time.Duration(val), true
	}
	var num, ok = tryNum2Int64(v)
	return time.Duration(num), ok
}

// try to convert number type to int
func tryNum2Int(v interface{}) (int, bool) {
	var (
		iv int
		ok = true
	)
	switch val := v.(type) {
	case int8:
		iv = int(val)
	case uint8:
		iv = int(val)
	case int16:
		iv = int(val)
	case uint16:
		iv = int(val)
	case int32:
		iv = int(val)
	case uint32:
		iv = int(val)
	case int64:
		iv = int(val)
	case uint64:
		iv = int(val)
	case int:
		iv = val
	case uint:
		iv = int(val)
	case JsInt64:
		iv = int(val)
	case JsUInt64:
		iv = int(val)
	case time.Duration:
		iv = int(val)
	case Duration:
		iv = int(val)
	default:
		ok = false
	}
	return iv, ok
}

// try to convert number type to int
func tryNum2Int64(v interface{}) (int64, bool) {
	var (
		iv int64
		ok = true
	)
	switch val := v.(type) {
	case int8:
		iv = int64(val)
	case uint8:
		iv = int64(val)
	case int16:
		iv = int64(val)
	case uint16:
		iv = int64(val)
	case int32:
		iv = int64(val)
	case uint32:
		iv = int64(val)
	case int64:
		iv = val
	case uint64:
		iv = int64(val)
	case int:
		iv = int64(val)
	case uint:
		iv = int64(val)
	case JsInt64:
		iv = int64(val)
	case JsUInt64:
		iv = int64(val)
	case time.Duration:
		iv = int64(val)
	case Duration:
		iv = int64(val)
	default:
		ok = false
	}
	return iv, ok
}
