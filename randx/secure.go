package randx

import (
	"crypto/rand"
	"encoding/binary"
	"io"
)

// IntNSecure random generate int, read data from linux /dev/urandom
// in range [0, n)
func IntNSecure(n int) int {
	if n <= 0 {
		panic("invalid argument to IntN")
	}

	// Use rejection sampling for uniform distribution
	mask := uint64(1)
	for mask < uint64(n) {
		mask <<= 1
	}
	mask--

	for {
		val := RandUint64Secure() & mask
		if val < uint64(n) {
			return int(val)
		}
	}
}

// IntRangeSecure returns a random int in [min, max) range
// Convenient wrapper for common range operations
func IntRangeSecure(_min, _max int) int {
	if _min >= _max {
		panic("invalid range: min must be less than max")
	}
	return _min + IntNSecure(_max-_min)
}

// IntBetweenSecure returns a random int in [min, max] range (inclusive)
// Convenient wrapper for common range operations
func IntBetweenSecure(_min, _max int) int {
	if _min == _max {
		return _min
	}
	if _min > _max {
		panic("invalid range: min must be less than or equal to max")
	}
	return _min + IntNSecure(_max-_min+1)
}

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

// Float64RangeSecure returns a random float64 in [min, max) range
// Uses linear transformation to maintain uniform distribution
func Float64RangeSecure(_min, _max float64) float64 {
	if _min >= _max {
		panic("invalid range: min must be less than max")
	}

	return _min + Float64Secure()*(_max-_min)
}

// Float64Secure returns a float64 in [0, 1) range
// Uses 53 bits for proper float64 precision
func Float64Secure() float64 {
	// Use 53 bits for float64 precision to avoid bias
	return float64(RandUint64Secure()>>11) / (1 << 53)
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
