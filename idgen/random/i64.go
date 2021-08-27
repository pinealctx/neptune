package random

import (
	cRand "crypto/rand"
	"encoding/binary"
	"io"
	/* #nosec */
	"math/rand"
	"time"
)

func RandomI64() int64 {
	/* #nosec */
	s := rand.NewSource(time.Now().UnixNano())
	/* #nosec */
	r := rand.New(s)
	return r.Int63()
}

func SecRandomI64() int64 {
	var buf [8]byte
	var _, err = io.ReadFull(cRand.Reader, buf[:])
	if err != nil {
		return RandomI64()
	}

	var v = int64(binary.LittleEndian.Uint64(buf[:]))
	if v < 0 {
		v = -v
	}
	return v
}
