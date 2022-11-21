package randx

import (
	crd "crypto/rand"
	"encoding/binary"
	"io"
	/* #nosec */
	"math/rand"
	"time"
)

// RandBetween : generate rand between [a, b], return int
func RandBetween(min int, max int) int {
	if min > max {
		panic("invalid min and max of RandBetween")
	}
	if min == max {
		//no choice
		return min
	}
	var r = feedNowRand()
	return r.Intn(max-min+1) + min
}

// SimpleRandBetween : generate rand between [a, b], use same global rand source: rand.Rand
func SimpleRandBetween(min int, max int) int {
	if min > max {
		panic("invalid min and max of QRandBetween")
	}
	if min == max {
		//no choice
		return min
	}
	/* #nosec */
	rand.Seed(time.Now().UnixNano())
	/* #nosec */
	return rand.Intn(max-min+1) + min
}

// RandBetweenSecure : generate rand between [a, b], return int
func RandBetweenSecure(min int, max int) int {
	if min > max {
		panic("invalid min and max of SecRandBetween")
	}
	if min == max {
		//no choice
		return min
	}
	var buf [8]byte
	var _, err = io.ReadFull(crd.Reader, buf[:])
	if err != nil {
		return RandBetween(min, max)
	}
	var seed = int64(binary.LittleEndian.Uint64(buf[:]))
	/* #nosec */
	var r = rand.New(rand.NewSource(seed))
	return r.Intn(max-min+1) + min
}
