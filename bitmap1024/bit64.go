package bitmap1024

import "math/bits"

var (
	//u64 mask table
	u64Tab [64]Bit64
	//64 sequence bytes
	seq64Buf [64]byte
)

//Bit64 uint64 as bitmap
type Bit64 uint64

//Set : set bit
func (b *Bit64) Set(i byte) {
	if i <= 63 {
		*b |= u64Tab[i]
	}
}

//Unset : unset bit
func (b *Bit64) Unset(i byte) {
	if i <= 63 {
		*b &= ^u64Tab[i]
	}
}

//Len : length of set bit
func (b Bit64) Len() int {
	if b.Full() {
		return 64
	}
	return bits.OnesCount64(uint64(b))
}

//NLen : length of not set bit
func (b Bit64) NLen() int {
	return 64 - b.Len()
}

//Full : is full set
func (b Bit64) Full() bool {
	return b == ^Bit64(0)
}

//Reverse : reverse
func (b Bit64) Reverse() Bit64 {
	return ^b
}

//And : and
func (b Bit64) And(c Bit64) Bit64 {
	return b & c
}

//Or : or
func (b Bit64) Or(c Bit64) Bit64 {
	return b | c
}

//GetNAsI64 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit64) GetNAsI64(n int) []int64 {
	return b.getNAsI64(n, false)
}

//RGetNAsI64 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit64) RGetNAsI64(n int) []int64 {
	return b.getNAsI64(n, true)
}

//GetNAsI32 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit64) GetNAsI32(n int) []int32 {
	return b.getNAsI32(n, false)
}

//RGetNAsI32 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit64) RGetNAsI32(n int) []int32 {
	return b.getNAsI32(n, true)
}

//GetNAsI16 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit64) GetNAsI16(n int) []int16 {
	return b.getNAsI16(n, false)
}

//RGetNAsI16 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit64) RGetNAsI16(n int) []int16 {
	return b.getNAsI16(n, true)
}

//GetNAsI8 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit64) GetNAsI8(n int) []int8 {
	return b.getNAsI8(n, false)
}

//RGetNAsI8 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b Bit64) RGetNAsI8(n int) []int8 {
	return b.getNAsI8(n, true)
}

//IterAsI64 : iterate bitmap
//input: s -- an int64 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) IterAsI64(s []int64, pos int, add int64, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := int64(0); i < 64; i++ {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = i + add
			cursor++
			c++
		}
	}
	return c
}

//RIterAsI64 : reverse iterate bitmap
//input: s -- an int64 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) RIterAsI64(s []int64, pos int, add int64, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := int64(63); i >= 0; i-- {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = i + add
			cursor++
			c++
		}
	}
	return c
}

//IterAsI32 : iterate bitmap
//input: s -- an int32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) IterAsI32(s []int32, pos int, add int32, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := int32(0); i < 64; i++ {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = i + add
			cursor++
			c++
		}
	}
	return c
}

//RIterAsI32 : reverse iterate bitmap
//input: s -- an int32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) RIterAsI32(s []int32, pos int, add int32, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := int32(63); i >= 0; i-- {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = i + add
			cursor++
			c++
		}
	}
	return c
}

//IterAsU32 : iterate bitmap
//input: s -- an uint32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) IterAsU32(s []uint32, pos int, add uint32, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := uint32(0); i < 64; i++ {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = i + add
			cursor++
			c++
		}
	}
	return c
}

//RIterAsU32 : reverse iterate bitmap
//input: s -- an uint32 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) RIterAsU32(s []uint32, pos int, add uint32, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := int32(63); i >= 0; i-- {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = uint32(i) + add
			cursor++
			c++
		}
	}
	return c
}

//IterAsI16 : iterate bitmap
//input: s -- an int16 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) IterAsI16(s []int16, pos int, add int16, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := int16(0); i < 64; i++ {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = i + add
			cursor++
			c++
		}
	}
	return c
}

//RIterAsI16 : reverse iterate bitmap
//input: s -- an int16 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) RIterAsI16(s []int16, pos int, add int16, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := int16(63); i >= 0; i-- {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = i + add
			cursor++
			c++
		}
	}
	return c
}

//IterAsI8 : iterate bitmap
//input: s -- an int8 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) IterAsI8(s []int8, pos int, add int8, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := int8(0); i < 64; i++ {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = i + add
			cursor++
			c++
		}
	}
	return c
}

//RIterAsI8 : reverse iterate bitmap
//input: s -- an int8 array, add -- added number, n -- iter count
//return:
//actually iter count
func (b Bit64) RIterAsI8(s []int8, pos int, add int8, n int) int {
	var l = b.Len()
	if l == 0 {
		return 0
	}
	var c = 0
	var cursor = pos
	for i := int8(63); i >= 0; i-- {
		if b&u64Tab[i] != 0 {
			if c >= n || c >= l {
				break
			}
			s[cursor] = i + add
			cursor++
			c++
		}
	}
	return c
}

//getNAsI64 : iterate bitmap
//input: n -- iter number
//reverse -- reverse iter or not
func (b Bit64) getNAsI64(n int, reverse bool) []int64 {
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
func (b Bit64) getNAsI32(n int, reverse bool) []int32 {
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
func (b Bit64) getNAsI16(n int, reverse bool) []int16 {
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

//getNAsI8 : iterate bitmap
//input: n -- iter number
//reverse -- reverse iter or not
func (b Bit64) getNAsI8(n int, reverse bool) []int8 {
	var s = make([]int8, n)
	var iterN int
	if reverse {
		iterN = b.RIterAsI8(s, 0, 0, n)
	} else {
		iterN = b.IterAsI8(s, 0, 0, n)
	}
	if iterN == 0 {
		return nil
	}
	return s[:iterN]
}

func init() {
	for i := uint64(0); i < 64; i++ {
		u64Tab[i] = 1 << i
	}
	for i := byte(0); i < 64; i++ {
		seq64Buf[i] = i
	}
}
