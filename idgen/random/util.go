package random

import (
	cRand "crypto/rand"
	"encoding/binary"
	"encoding/hex"
	"hash"
	"io"
	/* #nosec */
	"math/rand"
	"strings"
	"time"
)

type Intn func(n int) int

//生成随机字符串
func GenNonceStr(baseStr string, length int) string {
	//随机数种子设置
	/* #nosec */
	var s = rand.NewSource(time.Now().UnixNano())
	/* #nosec */
	var r = rand.New(s)

	return genNonceStr(baseStr, length, r.Intn)
}

//Security rand
func SecGenNonceStr(baseStr string, length int) string {
	//随机数种子设置
	var buf [8]byte
	var _, err = io.ReadFull(cRand.Reader, buf[:])
	var seed int64
	if err != nil {
		seed = time.Now().UnixNano()
	} else {
		seed = int64(binary.LittleEndian.Uint64(buf[:]))
		if seed < 0 {
			seed = -seed
		}
	}

	/* #nosec */
	var s = rand.NewSource(seed)
	/* #nosec */
	var r = rand.New(s)
	return genNonceStr(baseStr, length, r.Intn)
}

func genNonceStr(baseStr string, length int, fn Intn) string {
	var strBuilder strings.Builder
	var bSize = len(baseStr)

	var index int
	for i := 0; i < length; i++ {
		index = fn(bSize - 1)
		strBuilder.WriteByte(baseStr[index])
	}

	return strBuilder.String()
}

//生成相关的hash
func writeHex(h hash.Hash, buf []byte) string {
	_, _ = h.Write(buf)
	return hex.EncodeToString(h.Sum(nil))
}
