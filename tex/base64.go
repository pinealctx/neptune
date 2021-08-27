package tex

import (
	"database/sql/driver"
	"encoding/base64"
	"fmt"
	"reflect"
)

type Base64Bytes []byte

func (i *Base64Bytes) Scan(value interface{}) error {
	var ds string
	switch v := value.(type) {
	case []byte:
		ds = string(v)
	case string:
		ds = v
	default:
		return fmt.Errorf("unsupported.base64.type:%+v", reflect.TypeOf(value))
	}
	var buf, err = base64.RawStdEncoding.DecodeString(ds)
	if err != nil {
		return err
	}
	*i = buf
	return nil
}

func (i Base64Bytes) Value() (driver.Value, error) {
	return base64.RawStdEncoding.EncodeToString(i), nil
}
