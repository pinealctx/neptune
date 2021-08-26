package tex

import (
	"github.com/json-iterator/go"
	"io/ioutil"
)

const (
	CompressThreshHold              = 256
	CompressSnappy     CompressType = 128
)

type CompressType byte

var (
	jsonStd     = jsoniter.ConfigCompatibleWithStandardLibrary
	jsonDefault = jsoniter.ConfigDefault
	jsonFast    = jsoniter.ConfigFastest
)

var _ = jsonDefault

//标准jsoniter json库 100%兼容
var (
	JSONMarshal   = jsonStd.Marshal
	JSONUnmarshal = jsonStd.Unmarshal
)

//快速jsoniter json库 -- 浮点数会丢失精度，小数点最多后6位
var (
	JSONFastMarshal   = jsonFast.Marshal
	JSONFastUnmarshal = jsonFast.Unmarshal

	JSONFastMarshalSnappy   = fastMarshalSnappy
	JSONFastUnmarshalSnappy = fastUnmarshalSnappy

	TrySnappyCompress = compressJSON
)

var _ = TrySnappyCompress

//LoadJSONFile2Obj
//读取json文件并序列化到传入的指针中
func LoadJSONFile2Obj(fileName string, v interface{}) error {
	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = JSONUnmarshal(buf, v)
	return err
}

//fast json with compress
func fastMarshalSnappy(v interface{}) ([]byte, error) {
	var buf, err = JSONFastMarshal(v)
	if err != nil {
		return nil, err
	}
	return compressJSON(buf), nil
}

//fast json with compress unmarshal
func fastUnmarshalSnappy(data []byte, v interface{}) error {
	var size = len(data)
	if size == 0 {
		return nil
	}
	if CompressType(data[0]) == CompressSnappy {
		var buf, err = Snappy.DeCompress(data[1:])
		if err != nil {
			return err
		}
		return JSONFastUnmarshal(buf, v)
	} else {
		return JSONFastUnmarshal(data, v)
	}
}

//compress json
func compressJSON(buf []byte) []byte {
	var bufSize = len(buf)
	if bufSize > CompressThreshHold {
		var compressBuf, _ = Snappy.CompressWithPrefix(buf, []byte{byte(CompressSnappy)})
		if len(compressBuf) > bufSize {
			return buf
		}
		return compressBuf
	}
	return buf
}
