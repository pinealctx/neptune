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
func RandBetween(minV, maxV int) int {
	if minV > maxV {
		panic("invalid minV and maxV of RandBetween")
	}
	if minV == maxV {
		//no choice
		return minV
	}
	var r = feedNowRand()
	return r.Intn(maxV-minV+1) + minV
}

// SimpleRandBetween : generate rand between [a, b], use same global rand source: rand.Rand
func SimpleRandBetween(minV, maxV int) int {
	if minV > maxV {
		panic("invalid minV and maxV of QRandBetween")
	}
	if minV == maxV {
		//no choice
		return minV
	}
	// nolint : gosec // this random is only used for simple rand
	rand.New(rand.NewSource(time.Now().UnixNano()))
	// nolint : gosec // this random is only used for simple rand
	return rand.Intn(maxV-minV+1) + minV
}

// RandBetweenSecure : generate rand between [a, b], return int
func RandBetweenSecure(minV, maxV int) int {
	if minV > maxV {
		panic("invalid minV and maxV of SecRandBetween")
	}
	if minV == maxV {
		//no choice
		return minV
	}
	var buf [8]byte
	var _, err = io.ReadFull(crd.Reader, buf[:])
	if err != nil {
		return RandBetween(minV, maxV)
	}
	var seed = int64(binary.LittleEndian.Uint64(buf[:]))
	/* #nosec */
	var r = rand.New(rand.NewSource(seed))
	return r.Intn(maxV-minV+1) + minV
}
