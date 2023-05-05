package compress

import (
	"github.com/golang/snappy"
)

var (
	Snappy = &_Snappy{}
)

// snappy compress
type _Snappy struct {
}

func (s *_Snappy) Compress(src []byte) []byte {
	if len(src) == 0 {
		return nil
	}
	return snappy.Encode(nil, src)
}

func (s *_Snappy) DeCompress(src []byte) ([]byte, error) {
	if len(src) == 0 {
		return nil, nil
	}
	return snappy.Decode(nil, src)
}

func (s *_Snappy) CompressWithPrefix(src []byte, prefix []byte) []byte {
	var srcSize = len(src)
	if srcSize == 0 {
		return prefix
	}
	var prefixSize = len(prefix)
	if prefixSize == 0 {
		return s.Compress(src)
	}
	var totalSize = snappy.MaxEncodedLen(srcSize) + prefixSize
	var totalBuf = make([]byte, totalSize)
	copy(totalBuf[:prefixSize], prefix)
	var encBuf = snappy.Encode(totalBuf[prefixSize:], src)
	var encSize = len(encBuf)
	var actualSize = prefixSize + encSize
	if totalSize >= actualSize {
		//enough
		return totalBuf[:actualSize]
	} else {
		var extendBuf = make([]byte, actualSize)
		copy(extendBuf[:prefixSize], prefix)
		copy(extendBuf[prefixSize:], encBuf)
		return extendBuf
	}
}
