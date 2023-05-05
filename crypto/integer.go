package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
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
	key    []byte
	iv     []byte
	block  cipher.Block
	EncU32 [24]func(uint32) uint32
	DecU32 [24]func(uint32) uint32
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
		x.EncU32[i] = x.createCryptFunc(x.encryptU32, i)
		x.DecU32[i] = x.createCryptFunc(x.decryptU32, i)
	}
	return x
}

func (x *IntCypher) EncryptU64(number uint64) uint64 {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], number)
	var stream = cipher.NewCTR(x.block, x.iv)
	stream.XORKeyStream(data[:], data[:])
	return binary.BigEndian.Uint64(data[:])
}

func (x *IntCypher) DecryptU64(number uint64) uint64 {
	var data [8]byte
	binary.BigEndian.PutUint64(data[:], number)
	var stream = cipher.NewCTR(x.block, x.iv)
	stream.XORKeyStream(data[:], data[:])
	return binary.BigEndian.Uint64(data[:])
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

func (x *IntCypher) createCryptFunc(fn func(uint32, byte) uint32, i byte) func(uint32) uint32 {
	return func(n uint32) uint32 {
		return fn(n, i)
	}
}
