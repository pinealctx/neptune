package semap

import (
	"encoding/binary"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"hash/crc32"
	"reflect"
	"sort"
)

//Bs an interface can convert self to bytes
type Bs interface {
	//ToBytes convert self to bytes
	ToBytes() []byte
}

//SimpleIndex figure simple index
func SimpleIndex(i interface{}) int {
	var it uint64
	switch v := i.(type) {
	case byte:
		it = uint64(v)
	case int8:
		it = uint64(v)
	case int16:
		it = uint64(v)
	case uint16:
		it = uint64(v)
	case int32:
		it = uint64(v)
	case uint32:
		it = uint64(v)
	case int64:
		it = uint64(v)
	case uint64:
		it = v
	case int:
		it = uint64(v)
	case uint:
		it = uint64(v)
	default:
		return XHashIndex(i)
	}
	return int(it % numbs)
}

//XHashIndex figure xhash index
func XHashIndex(i interface{}) int {
	return SearchIndex(XXHash(i))
}

//SearchIndex search uint64 index
func SearchIndex(x uint64) int {
	var i = SearchUInt64s(nps, x)
	if i < 0 || i >= int(numbs) {
		return 0
	}
	return i
}

//SearchUInt64s search uint64 array index
func SearchUInt64s(a []uint64, x uint64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

//XXHash : use xxhash hash interface
func XXHash(i interface{}) uint64 {
	switch v := i.(type) {
	case string:
		return xxhash.Sum64String(v)
	default:
		return xxhash.Sum64(ToBytes(i))
	}
}

//ToBytes : convert interface to []byte
func ToBytes(i interface{}) []byte {
	switch v := i.(type) {
	case byte:
		return []byte{v}
	case int8:
		return []byte{byte(v)}
	case int16:
		var buf [2]byte
		binary.LittleEndian.PutUint16(buf[:], uint16(v))
		return buf[:]
	case uint16:
		var buf [2]byte
		binary.LittleEndian.PutUint16(buf[:], v)
		return buf[:]
	case int32:
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], uint32(v))
		return buf[:]
	case uint32:
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], v)
		return buf[:]
	case int64:
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(v))
		return buf[:]
	case uint64:
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], v)
		return buf[:]
	case int:
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(v))
		return buf[:]
	case uint:
		var buf [8]byte
		binary.LittleEndian.PutUint64(buf[:], uint64(v))
		return buf[:]
	case string:
		var crc = crc32.ChecksumIEEE([]byte(v))
		var buf [4]byte
		binary.LittleEndian.PutUint32(buf[:], crc)
		return buf[:]
	case []byte:
		return v
	case Bs:
		return v.ToBytes()
	default:
		panic(fmt.Sprintf("unsupported.type.for.slot:%+v", reflect.TypeOf(i)))
	}
}
