package randx

import (
	/* #nosec */
	"math/rand"
	"time"
)

// RandInt random generate int, rand seed is current nano timestamp
func RandInt() int {
	return int(RandUint64())
}

// RandUint random generate uint, rand seed is current nano timestamp
func RandUint() uint {
	return uint(RandUint64())
}

// RandInt64 random generate int64, rand seed is current nano timestamp
func RandInt64() int64 {
	return int64(RandUint64())
}

// RandUint64 random generate uint64, rand seed is current nano timestamp
func RandUint64() uint64 {
	var r = feedNowRand()
	return r.Uint64()
}

// RandInt32 random generate int32, rand seed is current nano timestamp
func RandInt32() int32 {
	return int32(RandUint32())
}

// RandUint32 random generate uint32, rand seed is current nano timestamp
func RandUint32() uint32 {
	var r = feedNowRand()
	return r.Uint32()
}

// feed current nano timestamp as rand seed
func feedNowRand() *rand.Rand {
	/* #nosec */
	s := rand.NewSource(time.Now().UnixNano())
	/* #nosec */
	return rand.New(s)
}
