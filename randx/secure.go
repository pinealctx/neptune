package randx

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

// RandIntSecure random generate int, read data from linux /dev/urandom
// actually, it's almost a random id
func RandIntSecure() int {
	var v, err = randReadU64()
	if err != nil {
		return RandInt()
	}
	return int(v)
}

// RandUintSecure random generate uint, read data from linux /dev/urandom
// actually, it's almost a random id
func RandUintSecure() uint {
	var v, err = randReadU64()
	if err != nil {
		return RandUint()
	}
	return uint(v)
}

// RandInt64Secure random generate int64, read data from linux /dev/urandom
// actually, it's almost a random id
func RandInt64Secure() int64 {
	var v, err = randReadU64()
	if err != nil {
		return RandInt64()
	}
	return int64(v)
}

// RandUint64Secure random generate uint64, read data from linux /dev/urandom
// actually, it's almost a random id
func RandUint64Secure() uint64 {
	var v, err = randReadU64()
	if err != nil {
		return RandUint64()
	}
	return v
}

// RandUint32Secure random generate uint32, read data from linux /dev/urandom
// actually, it's almost a random id
func RandUint32Secure() uint32 {
	var v, err = randReadU32()
	if err != nil {
		return RandUint32()
	}
	return v
}

// RandInt32Secure random generate int32, read data from linux /dev/urandom
// actually, it's almost a random id
func RandInt32Secure() int32 {
	var v, err = randReadU32()
	if err != nil {
		return RandInt32()
	}
	return int32(v)
}

// read from linux /dev/urandom to wrap a random uint64
func randReadU64() (uint64, error) {
	var buf [8]byte
	var _, err = io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf[:]), nil
}

// read from linux /dev/urandom to wrap a random uint32
func randReadU32() (uint32, error) {
	var buf [4]byte
	var _, err = io.ReadFull(rand.Reader, buf[:])
	if err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf[:]), nil
}
