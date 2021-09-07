package mux

import (
	"encoding/binary"
	"hash/crc32"
)

//Hashed2Int : can be hashed to int
type Hashed2Int interface {
	//HashedInt : the hashed int
	HashedInt() int
}

/*---basic type to implement HashedInt-------*/

/*--------basic hashed int-------*/
//基本的整型使用强制类型转换就可以作为其散列值
//前提是这些数值没有太多的规律，比较杂乱，如果是雪花算法产生的id，推荐使用后面的用CRC32进行散列的类型。

//Byte : byte implement HashedInt
type Byte byte

func (v Byte) HashedInt() int {
	return int(v)
}

//Int8 : byte implement HashedInt
type Int8 int8

func (v Int8) HashedInt() int {
	return int(v)
}

//Int16 : int16 implement HashedInt
type Int16 int16

func (v Int16) HashedInt() int {
	return int(v)
}

//UInt16 : uint16 implement HashedInt
type UInt16 uint16

func (v UInt16) HashedInt() int {
	return int(v)
}

//Int32 : int32 implement HashedInt, use crc32 to hash
type Int32 int32

func (v Int32) HashedInt() int {
	return int(v)
}

//UInt32 : uint32 implement HashedInt, use crc32 to hash
type UInt32 uint32

func (v UInt32) HashedInt() int {
	return int(v)
}

//Int64 : int64 implement HashedInt, use crc32 to hash
type Int64 int64

func (v Int64) HashedInt() int {
	return int(v)
}

//UInt64 : uint64 implement HashedInt, use crc32 to hash
type UInt64 uint64

func (v UInt64) HashedInt() int {
	return int(v)
}

//Int : int implement HashedInt
type Int int

func (v Int) HashedInt() int {
	return int(v)
}

//UInt : uint implement HashedInt
type UInt uint

func (v UInt) HashedInt() int {
	return int(v)
}

/*--------crc 32 hashed int-------*/
//对于一些有规律的整数，需要使用CRC来散列，例如雪花算法产生的id，后面步长相关的尾数大部分都为0。

//Int32CRC : int32 implement HashedInt, use crc32 to hash
type Int32CRC int32

func (v Int32CRC) HashedInt() int {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	return int(crc32.ChecksumIEEE(b[:]))
}

//UInt32CRC : uint32 implement HashedInt, use crc32 to hash
type UInt32CRC uint32

func (v UInt32CRC) HashedInt() int {
	var b [4]byte
	binary.LittleEndian.PutUint32(b[:], uint32(v))
	return int(crc32.ChecksumIEEE(b[:]))
}

//Int64CRC : int64 implement HashedInt, use crc32 to hash
type Int64CRC int64

func (v Int64CRC) HashedInt() int {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	return int(crc32.ChecksumIEEE(b[:]))
}

//UInt64CRC : uint64 implement HashedInt, use crc32 to hash
type UInt64CRC uint64

func (v UInt64CRC) HashedInt() int {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	return int(crc32.ChecksumIEEE(b[:]))
}

//IntCRC : int implement HashedInt, use crc32 to hash
type IntCRC int

func (v IntCRC) HashedInt() int {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	return int(crc32.ChecksumIEEE(b[:]))
}

//UIntCRC : uint implement HashedInt, use crc32 to hash
type UIntCRC uint

func (v UIntCRC) HashedInt() int {
	var b [8]byte
	binary.LittleEndian.PutUint64(b[:], uint64(v))
	return int(crc32.ChecksumIEEE(b[:]))
}

/*--------crc 32 hashed string and bytes-------*/
//字符串与byte数值需要使用crc散列

//String : string implement HashedInt, use crc32 to hash
type String string

func (v String) HashedInt() int {
	return int(crc32.ChecksumIEEE([]byte(v)))
}

//Bytes : bytes implement HashedInt, use crc32 to hash
type Bytes []byte

func (v Bytes) HashedInt() int {
	return int(crc32.ChecksumIEEE(v))
}
