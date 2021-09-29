package bitmap1024

import (
	"encoding/binary"
	"fmt"
)

const (
	L16  = 16
	B64  = 64
	L128 = 128
)

//Bit1024 : 16*8 = 128bytes
type Bit1024 []Bit64

//NewBit1024 : new bit 1024
func NewBit1024() Bit1024 {
	return make([]Bit64, L16)
}

//SetI32 : set bit
func (b Bit1024) SetI32(i int32) {
	var index = i / B64
	if index >= 0 && index < L16 {
		var mod = byte(i % B64)
		b[index].Set(mod)
	}
}

//UnsetI32 : unset bit
func (b Bit1024) UnsetI32(i int32) {
	var index = i / B64
	if index >= 0 && index < L16 {
		var mod = byte(i % B64)
		b[index].Unset(mod)
	}
}

//SetI16 : set bit
func (b Bit1024) SetI16(i int16) {
	var index = i / B64
	if index >= 0 && index < L16 {
		var mod = byte(i % B64)
		b[index].Set(mod)
	}
}

//UnsetI16 : unset bit
func (b Bit1024) UnsetI16(i int16) {
	var index = i / B64
	if index >= 0 && index < L16 {
		var mod = byte(i % B64)
		b[index].Unset(mod)
	}
}

//Len : length of set bit
func (b Bit1024) Len() int {
	var c int
	for i := 0; i < L16; i++ {
		c += b[i].Len()
	}
	return c
}

func (b Bit1024) NLen() int {
	return 1024 - b.Len()
}

//Reverse : reverse
func (b Bit1024) Reverse() Bit1024 {
	var c = make([]Bit64, L16)
	for i := 0; i < L16; i++ {
		c[i] = b[i].Reverse()
	}
	return c
}

//OrThenReverse : or operation then reverse
func (b Bit1024) OrThenReverse(c Bit1024) Bit1024 {
	var d = make([]Bit64, L16)
	for i := 0; i < L16; i++ {
		d[i] = b[i].Or(c[i]).Reverse()
	}
	return d
}

//And : and
func (b Bit1024) And(c Bit1024) Bit1024 {
	var d = make([]Bit64, L16)
	for i := 0; i < L16; i++ {
		d[i] = b[i].And(c[i])
	}
	return d
}

//Or : or
func (b Bit1024) Or(c Bit1024) Bit1024 {
	var d = make([]Bit64, L16)
	for i := 0; i < L16; i++ {
		d[i] = b[i].Or(c[i])
	}
	return d
}

//Marshal : marshal bitmap
func (b Bit1024) Marshal() []byte {
	var n = b.Len()
	if n == 0 {
		return nil
	}
	if n < 64 {
		var buf = make([]byte, n*2)
		var s = b.GetNAsI16(n)
		for i := 0; i < n; i++ {
			binary.LittleEndian.PutUint16(buf[i*2:], uint16(s[i]))
		}
		return buf
	}
	var buf = make([]byte, L128)
	for i := 0; i < L16; i++ {
		binary.LittleEndian.PutUint64(buf[i*8:], uint64(b[i]))
	}
	return buf
}

//Unmarshal : unmarshal bitmap
func (b Bit1024) Unmarshal(buf []byte) error {
	var n = len(buf)
	if n == 0 {
		return nil
	}
	if n > L128 {
		return fmt.Errorf("bit.1024.out.of.range:%+v", n)
	}
	if n%2 != 0 {
		return fmt.Errorf("bit.1024.invalid.length:%+v", n)
	}

	if n < L128 {
		var en = n / 2
		for i := 0; i < en; i++ {
			var i16 = int16(binary.LittleEndian.Uint16(buf[i*2:]))
			if i16 < 0 || i16 > 1023 {
				return fmt.Errorf("bit.1024.invalid.element:%+v", i16)
			}
			b.SetI16(i16)
		}
		return nil
	}

	//n must be 128
	for i := 0; i < L16; i++ {
		var b64 = Bit64(binary.LittleEndian.Uint64(buf[i*8:]))
		b[i] = b64
	}
	return nil
}

//Equal : equal other
func (b Bit1024) Equal(c Bit1024) bool {
	for i := 0; i < L16; i++ {
		if b[i] != c[i] {
			return false
		}
	}
	return true
}

//GetNAsI64 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit1024) GetNAsI64(n int) []int64 {
	return b.getNAsI64(n, false)
}

//RGetNAsI64 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit1024) RGetNAsI64(n int) []int64 {
	return b.getNAsI64(n, true)
}

//GetNAsI32 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit1024) GetNAsI32(n int) []int32 {
	return b.getNAsI32(n, false)
}

//RGetNAsI32 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit1024) RGetNAsI32(n int) []int32 {
	return b.getNAsI32(n, true)
}

//GetNAsI16 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit1024) GetNAsI16(n int) []int16 {
	return b.getNAsI16(n, false)
}

//RGetNAsI16 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit1024) RGetNAsI16(n int) []int16 {
	return b.getNAsI16(n, true)
}

//IterAsI64 : iterate bitmap
//input: s -- an int64 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit1024) IterAsI64(s []int64, pos int, add int64, n int) int {
	var (
		iterN  = 0
		left   = n
		cursor = pos
		eIterN int
	)
	for i := int64(0); i < L16; i++ {
		if iterN >= n {
			break
		}
		eIterN = b[i].IterAsI64(s, cursor, B64*i+add, left)
		iterN += eIterN
		cursor += eIterN
		left = n - iterN
	}
	return iterN
}

//RIterAsI64 : reverse iterate bitmap
//input: s -- an int64 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit1024) RIterAsI64(s []int64, pos int, add int64, n int) int {
	var (
		iterN  = 0
		left   = n
		cursor = pos
		eIterN int
	)
	for i := int64(L16 - 1); i >= 0; i-- {
		if iterN >= n {
			break
		}
		eIterN = b[i].RIterAsI64(s, cursor, B64*i+add, left)
		iterN += eIterN
		cursor += eIterN
		left = n - iterN
	}
	return iterN
}

//IterAsI32 : iterate bitmap
//input: s -- an int32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit1024) IterAsI32(s []int32, pos int, add int32, n int) int {
	var (
		iterN  = 0
		left   = n
		cursor = pos
		eIterN int
	)
	for i := int32(0); i < L16; i++ {
		if iterN >= n {
			break
		}
		eIterN = b[i].IterAsI32(s, cursor, B64*i+add, left)
		iterN += eIterN
		cursor += eIterN
		left = n - iterN
	}
	return iterN
}

//RIterAsI32 : reverse iterate bitmap
//input: s -- an int32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit1024) RIterAsI32(s []int32, pos int, add int32, n int) int {
	var (
		iterN  = 0
		left   = n
		cursor = pos
		eIterN int
	)
	for i := int32(L16 - 1); i >= 0; i-- {
		if iterN >= n {
			break
		}
		eIterN = b[i].RIterAsI32(s, cursor, B64*i+add, left)
		iterN += eIterN
		cursor += eIterN
		left = n - iterN
	}
	return iterN
}

//IterAsU32 : iterate bitmap
//input: s -- an uint32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit1024) IterAsU32(s []uint32, pos int, add uint32, n int) int {
	var (
		iterN  = 0
		left   = n
		cursor = pos
		eIterN int
	)
	for i := uint32(0); i < L16; i++ {
		if iterN >= n {
			break
		}
		eIterN = b[i].IterAsU32(s, cursor, B64*i+add, left)
		iterN += eIterN
		cursor += eIterN
		left = n - iterN
	}
	return iterN
}

//RIterAsU32 : reverse iterate bitmap
//input: s -- an uint32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit1024) RIterAsU32(s []uint32, pos int, add uint32, n int) int {
	var (
		iterN  = 0
		left   = n
		cursor = pos
		eIterN int
	)
	for i := int32(L16 - 1); i >= 0; i-- {
		if iterN >= n {
			break
		}
		eIterN = b[i].RIterAsU32(s, cursor, B64*uint32(i)+add, left)
		iterN += eIterN
		cursor += eIterN
		left = n - iterN
	}
	return iterN
}

//IterAsI16 : iterate bitmap
//input: s -- an int32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit1024) IterAsI16(s []int16, pos int, add int16, n int) int {
	var (
		iterN  = 0
		left   = n
		cursor = pos
		eIterN int
	)
	for i := int16(0); i < L16; i++ {
		if iterN >= n {
			break
		}
		eIterN = b[i].IterAsI16(s, cursor, B64*i+add, left)
		iterN += eIterN
		cursor += eIterN
		left = n - iterN
	}
	return iterN
}

//RIterAsI16 : reverse iterate bitmap
//input: s -- an int32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit1024) RIterAsI16(s []int16, pos int, add int16, n int) int {
	var (
		iterN  = 0
		left   = n
		cursor = pos
		eIterN int
	)
	for i := int16(L16 - 1); i >= 0; i-- {
		if iterN >= n {
			break
		}
		eIterN = b[i].RIterAsI16(s, cursor, B64*i+add, left)
		iterN += eIterN
		cursor += eIterN
		left = n - iterN
	}
	return iterN
}

//getNAsI64 : iterate bitmap
//input: n -- iter number
//reverse -- reverse iter or not
func (b Bit1024) getNAsI64(n int, reverse bool) []int64 {
	var s = make([]int64, n)
	var iterN int
	if reverse {
		iterN = b.RIterAsI64(s, 0, 0, n)
	} else {
		iterN = b.IterAsI64(s, 0, 0, n)
	}
	if iterN == 0 {
		return nil
	}
	return s[:iterN]
}

//getNAsI32 : iterate bitmap
//input: n -- iter number
//reverse -- reverse iter or not
func (b Bit1024) getNAsI32(n int, reverse bool) []int32 {
	var s = make([]int32, n)
	var iterN int
	if reverse {
		iterN = b.RIterAsI32(s, 0, 0, n)
	} else {
		iterN = b.IterAsI32(s, 0, 0, n)
	}
	if iterN == 0 {
		return nil
	}
	return s[:iterN]
}

//getNAsI16 : iterate bitmap
//input: n -- iter number
//reverse -- reverse iter or not
func (b Bit1024) getNAsI16(n int, reverse bool) []int16 {
	var s = make([]int16, n)
	var iterN int
	if reverse {
		iterN = b.RIterAsI16(s, 0, 0, n)
	} else {
		iterN = b.IterAsI16(s, 0, 0, n)
	}
	if iterN == 0 {
		return nil
	}
	return s[:iterN]
}
