package tex

import (
	"encoding/binary"
)

//U128BitMap use 2 uint64
//Low and High
type U128BitMap struct {
	Low  U64BitMap
	High U64BitMap
}

func (u *U128BitMap) Set(i byte) {
	if i <= 63 {
		u.Low.Set(i)
		return
	}
	if i >= 64 && i <= 127 {
		u.High.Set(i - 64)
	}
}

func (u *U128BitMap) Len() int {
	return u.Low.Len() + u.High.Len()
}

func (u *U128BitMap) ReverseLen() int {
	return 128 - u.Len()
}

func (u *U128BitMap) Full() bool {
	return u.Low.Full() && u.High.Full()
}

func (u *U128BitMap) Values() []byte {
	var l = u.Len()
	if l == 0 {
		return nil
	}
	var buf = make([]byte, 0, l)
	buf = append(buf, u.Low.Values()...)
	buf = append(buf, u.High.ValuesX(64)...)
	return buf
}

func (u *U128BitMap) Left() []byte {
	var l = u.ReverseLen()
	if l == 0 {
		return nil
	}
	var buf = make([]byte, 0, l)
	buf = append(buf, u.Low.Left()...)
	buf = append(buf, u.High.LeftX(64)...)
	return buf
}

func (u *U128BitMap) Bytes() []byte {
	if u.Len() < 16 {
		return u.Values()
	}
	var buf [16]byte
	binary.LittleEndian.PutUint64(buf[:], uint64(u.Low))
	binary.LittleEndian.PutUint64(buf[8:], uint64(u.High))
	return buf[:]
}

func (u *U128BitMap) FromBytes(buf []byte) {
	var l = len(buf)
	if l < 16 {
		u.fromValues(buf)
		return
	}
	if l == 16 {
		u.Low = U64BitMap(binary.LittleEndian.Uint64(buf))
		u.High = U64BitMap(binary.LittleEndian.Uint64(buf[8:]))
	}
}

func (u *U128BitMap) fromValues(buf []byte) {
	var c byte
	for i := range buf {
		c = buf[i]
		u.Set(c)
	}
}
