package tex

import (
	"encoding/binary"
	"math/bits"
)

var (
	//u64 mask table
	u64Tab [64]U64BitMap
	//64 sequence bytes
	seq64Buf [64]byte
)

type U64BitMap uint64

func (b *U64BitMap) Set(i byte) {
	if i <= 63 {
		*b |= u64Tab[i]
	}
}

func (b U64BitMap) Len() int {
	if b.Full() {
		return 64
	}
	return bits.OnesCount64(uint64(b))
}

func (b U64BitMap) ReverseLen() int {
	return 64 - b.Len()
}

func (b U64BitMap) Full() bool {
	return b == ^U64BitMap(0)
}

func (b U64BitMap) Values() []byte {
	if b.Full() {
		return seq64Buf[:]
	}
	return b.values(0)
}

func (b U64BitMap) Left() []byte {
	var r = ^b
	return r.Values()
}

func (b U64BitMap) ValuesX(plus byte) []byte {
	if b.Full() {
		var c = seq64Buf
		for i := 0; i < 64; i++ {
			c[i] += plus
		}
		return c[:]
	}
	return b.values(plus)
}

func (b U64BitMap) LeftX(plus byte) []byte {
	var r = ^b
	return r.ValuesX(plus)
}

func (b U64BitMap) Bytes() []byte {
	if b.Len() < 8 {
		return b.Values()
	}
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(b))
	return buf[:]
}

func (b *U64BitMap) FromBytes(buf []byte) {
	var l = len(buf)
	if l < 8 {
		b.fromValues(buf)
		return
	}
	if l == 8 {
		*b = U64BitMap(binary.LittleEndian.Uint64(buf))
	}
}

func (b *U64BitMap) fromValues(buf []byte) {
	var c byte
	for i := range buf {
		c = buf[i]
		b.Set(c)
	}
}

func (b U64BitMap) values(plus byte) []byte {
	var l = b.Len()
	if l == 0 {
		return nil
	}
	var s = make([]byte, 0, l)
	var c = 0
	for i := byte(0); i < 64; i++ {
		if b&u64Tab[i] != 0 {
			s = append(s, i+plus)
			c++
			if c == l {
				break
			}
		}
	}
	return s
}

func init() {
	for i := uint64(0); i < 64; i++ {
		u64Tab[i] = 1 << i
	}
	for i := byte(0); i < 64; i++ {
		seq64Buf[i] = i
	}
}
