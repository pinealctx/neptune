package remap

import (
	"encoding/binary"
	"fmt"
	"github.com/cespare/xxhash/v2"
	"math"
	"reflect"
	"sort"
)

// Bs an interface can convert self to bytes
type Bs interface {
	//ToBytes convert self to bytes
	ToBytes() []byte
}

// HitGroup an interface can convert self to a fixed uint64 which can target a group id
type HitGroup interface {
	//Hit convert self to a specific uint64 value
	Hit() uint64
}

// ReMap remap meta info
type ReMap struct {
	numbs uint64   //切分的份数,取素数
	nps   []uint64 //切分uint64为211份,从小到大排列
}

// NewReMap : new remap instance
func NewReMap(opts ...Option) *ReMap {
	var r = &ReMap{}
	var o = &_Option{prime: DefaultPrime}
	for _, opt := range opts {
		opt(o)
	}
	r.numbs = o.prime
	var x uint64 = math.MaxUint64
	var y = x / r.numbs
	r.nps = make([]uint64, r.numbs)
	for i := uint64(0); i < r.numbs; i++ {
		r.nps[i] = y * (uint64(i) + 1)
	}
	r.nps[r.numbs-1] = math.MaxUint64
	return r
}

// Numbs 获取切分份数
func (r *ReMap) Numbs() uint64 {
	return r.numbs
}

// SimpleIndex figure simple index
func (r *ReMap) SimpleIndex(i interface{}) int {
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
	case HitGroup:
		it = v.Hit()
	default:
		return r.XHashIndex(i)
	}
	return int(it % r.numbs)
}

// XHashIndex figure xhash index
func (r *ReMap) XHashIndex(i interface{}) int {
	return r.SearchIndex(XXHash(i))
}

// SearchIndex search uint64 index
func (r *ReMap) SearchIndex(x uint64) int {
	var i = SearchUInt64s(r.nps, x)
	if i < 0 || i >= int(r.numbs) {
		return 0
	}
	return i
}

// SearchUInt64s search uint64 array index
func SearchUInt64s(a []uint64, x uint64) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

// XXHash : use xxhash hash interface
func XXHash(i interface{}) uint64 {
	switch v := i.(type) {
	case string:
		return xxhash.Sum64String(v)
	default:
		return xxhash.Sum64(ToBytes(i))
	}
}

// ToBytes : convert interface to []byte
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
		return []byte(v)
	case []byte:
		return v
	case Bs:
		return v.ToBytes()
	default:
		panic(fmt.Sprintf("unsupported.type.for.slot:%+v", reflect.TypeOf(i)))
	}
}
