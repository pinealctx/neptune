package jsonx

import (
	"testing"
)

func TestCompressJson1(t *testing.T) {
	var x = map[string]interface{}{
		"a": "int64(9223372036854775806)",
		"b": "JsInt64(9223372036854775806)",
		"c": "JsInt64(9223372036854775806)",
		"d": "JsInt64(9223372036854775806)",
	}
	var buf1, err = JSONFastMarshalSnappy(x)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("len:", len(buf1), " first:", buf1[0])

	var y map[string]interface{}
	err = JSONFastUnmarshalSnappy(buf1, &y)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(y)
}

func TestCompressJson2(t *testing.T) {
	var x = map[string]interface{}{
		"a": "int64(9223372036854775806)",
		"b": "JsInt64(9223372036854775806)",
		"c": "JsInt64(9223372036854775806)",
		"d": "JsInt64(9223372036854775806)",
		"e": "JsInt64(9223372036854775806)",
		"f": "JsInt64(9223372036854775806)",
		"g": "JsInt64(9223372036854775806)",
		"h": "JsInt64(9223372036854775806)",
		"i": "JsInt64(9223372036854775806)",
	}
	var buf1, err = JSONFastMarshalSnappy(x)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log("len:", len(buf1), " first:", buf1[0])

	var y map[string]interface{}
	err = JSONFastUnmarshalSnappy(buf1, &y)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(y)
}
