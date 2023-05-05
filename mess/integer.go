package mess

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"github.com/btcsuite/btcutil/base58"
)

var (
	// Full permutation of [0, 1, 2, 3]
	/*
	   1. [0, 1, 2, 3]
	   2. [0, 1, 3, 2]
	   3. [0, 2, 1, 3]
	   4. [0, 2, 3, 1]
	   5. [0, 3, 1, 2]
	   6. [0, 3, 2, 1]

	   7. [1, 0, 2, 3]
	   8. [1, 0, 3, 2]
	   9. [1, 2, 0, 3]
	   10. [1, 2, 3, 0]
	   11. [1, 3, 0, 2]
	   12. [1, 3, 2, 0]

	   13. [2, 0, 1, 3]
	   14. [2, 0, 3, 1]
	   15. [2, 1, 0, 3]
	   16. [2, 1, 3, 0]
	   17. [2, 3, 0, 1]
	   18. [2, 3, 1, 0]

	   19. [3, 0, 1, 2]
	   20. [3, 0, 2, 1]
	   21. [3, 1, 0, 2]
	   22. [3, 1, 2, 0]
	   23. [3, 2, 0, 1]
	   24. [3, 2, 1, 0]
	*/
	permutation = [24][4]int{
		{0, 1, 2, 3}, {0, 1, 3, 2}, {0, 2, 1, 3}, {0, 2, 3, 1}, {0, 3, 1, 2}, {0, 3, 2, 1},
		{1, 0, 2, 3}, {1, 0, 3, 2}, {1, 2, 0, 3}, {1, 2, 3, 0}, {1, 3, 0, 2}, {1, 3, 2, 0},
		{2, 0, 1, 3}, {2, 0, 3, 1}, {2, 1, 0, 3}, {2, 1, 3, 0}, {2, 3, 0, 1}, {2, 3, 1, 0},
		{3, 0, 1, 2}, {3, 0, 2, 1}, {3, 1, 0, 2}, {3, 1, 2, 0}, {3, 2, 0, 1}, {3, 2, 1, 0},
	}
)

type IntCypher struct {
	key           []byte
	iv            []byte
	block         cipher.Block
	EncU32        [24]func(uint32) uint32
	DecU32        [24]func(uint32) uint32
	EncU32ToStr   [24]func(uint32) string
	DecStrToU32   [24]func(string) uint32
	EncU32ToStrEx [24]func(uint32) string
	DecStrToU32Ex [24]func(string) uint32
}

func NewIntCypher(key []byte, iv []byte) *IntCypher {
	var block, err = aes.NewCipher(key)
	if err != nil {
		panic(err)
	}
	var x = &IntCypher{
		key:   key,
		iv:    iv,
		block: block,
	}
	for i := byte(0); i < 24; i++ {
		x.EncU32[i] = x.createCryptNumFunc(x.encryptU32, i)
		x.DecU32[i] = x.createCryptNumFunc(x.decryptU32, i)
		x.EncU32ToStr[i] = x.createCryptStrFunc(x.encryptU32ToStr, i)
		x.DecStrToU32[i] = x.createDecryptStrFunc(x.decryptStrToU32, i)
		x.EncU32ToStrEx[i] = x.createCryptStrFunc(x.encryptU32ToStrEx, i)
		x.DecStrToU32Ex[i] = x.createDecryptStrFunc(x.decryptStrToU32Ex, i)
	}
	return x
}

func (x *IntCypher) encryptU32(number uint32, v byte) uint32 {
	var data [4]byte
	binary.BigEndian.PutUint32(data[:], number)
	var stream = cipher.NewCTR(x.block, x.iv)
	stream.XORKeyStream(data[:], data[:])

	// re-permute data with v, v is index of permutation
	var result [4]byte
	for i := 0; i < 4; i++ {
		result[i] = data[permutation[v][i]]
	}
	return binary.BigEndian.Uint32(result[:])
}

func (x *IntCypher) decryptU32(number uint32, v byte) uint32 {
	var data [4]byte
	binary.BigEndian.PutUint32(data[:], number)

	// re-permute data with v, v is index of permutation
	var result [4]byte
	for i := 0; i < 4; i++ {
		result[permutation[v][i]] = data[i]
	}

	var stream = cipher.NewCTR(x.block, x.iv)
	stream.XORKeyStream(result[:], result[:])
	return binary.BigEndian.Uint32(result[:])
}

func (x *IntCypher) encryptU32ToStr(number uint32, v byte) string {
	var data [4]byte
	binary.BigEndian.PutUint32(data[:], number)
	var stream = cipher.NewCTR(x.block, x.iv)
	stream.XORKeyStream(data[:], data[:])

	// re-permute data with v, v is index of permutation
	var result [4]byte
	for i := 0; i < 4; i++ {
		result[i] = data[permutation[v][i]]
	}
	return base58.Encode(result[:])
}

func (x *IntCypher) decryptStrToU32(str string, v byte) uint32 {
	var data = base58.Decode(str)
	if len(data) != 4 {
		return 0
	}
	// re-permute data with v, v is index of permutation
	var result [4]byte
	for i := 0; i < 4; i++ {
		result[permutation[v][i]] = data[i]
	}

	var stream = cipher.NewCTR(x.block, x.iv)
	stream.XORKeyStream(result[:], result[:])
	return binary.BigEndian.Uint32(result[:])
}

func (x *IntCypher) encryptU32ToStrEx(number uint32, v byte) string {
	var data [4]byte
	binary.BigEndian.PutUint32(data[:], number)
	var stream = cipher.NewCTR(x.block, x.iv)
	stream.XORKeyStream(data[:], data[:])

	// re-permute data with v, v is index of permutation
	var result [4]byte
	for i := 0; i < 4; i++ {
		result[i] = data[permutation[v][i]]
	}
	var rs = base58.Encode(result[:])
	// switch first and second char
	if len(rs) < 2 {
		return rs
	}

	var cs = []byte(rs)
	cs[0], cs[1] = cs[1], cs[0]
	return string(cs)
}

func (x *IntCypher) decryptStrToU32Ex(str string, v byte) uint32 {
	// switch first and second char
	if len(str) >= 2 {
		var cs = []byte(str)
		cs[0], cs[1] = cs[1], cs[0]
		str = string(cs)
	}
	var data = base58.Decode(str)
	if len(data) != 4 {
		return 0
	}
	// re-permute data with v, v is index of permutation
	var result [4]byte
	for i := 0; i < 4; i++ {
		result[permutation[v][i]] = data[i]
	}

	var stream = cipher.NewCTR(x.block, x.iv)
	stream.XORKeyStream(result[:], result[:])
	return binary.BigEndian.Uint32(result[:])
}

func (x *IntCypher) createCryptNumFunc(fn func(uint32, byte) uint32, i byte) func(uint32) uint32 {
	return func(n uint32) uint32 {
		return fn(n, i)
	}
}

func (x *IntCypher) createCryptStrFunc(fn func(uint32, byte) string, i byte) func(uint32) string {
	return func(n uint32) string {
		return fn(n, i)
	}
}

func (x *IntCypher) createDecryptStrFunc(fn func(string, byte) uint32, i byte) func(string) uint32 {
	return func(s string) uint32 {
		return fn(s, i)
	}
}
