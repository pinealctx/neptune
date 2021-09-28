package bitmap1024

import (
	"fmt"
	"math"
)

/*
An int64 combine with [Uint32 1024]
The max value is 4398046510079(4万亿), it's a big number, can be used as some id, such as message id.
*/

const (
	C1K = 1024
)

//BigU32s BigU32 list
type BigU32s []*BigU32

//Reverse : clone and reverse bit
func (b BigU32s) Reverse() BigU32s{
	var l = len(b)
	if l == 0 {
		return nil
	}
	var c = make(BigU32s, l)
	for i := 0; i < l; i++ {
		c[i] = &BigU32{}
		c[i].Start = b[i].Start
		c[i].B1024 = b[i].B1024.Reverse()
	}
	return c
}

//GetNAsI64 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b BigU32s) GetNAsI64(n int) []int64 {
	return b.getNAsI64(n, false)
}

//RGetNAsI64 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b BigU32s) RGetNAsI64(n int) []int64 {
	return b.getNAsI64(n, true)
}

//getNAsI64 : iterate bitmap
//input: n -- iter count
//reverse -- reverse iter or not
//return: actually iter list
func (b BigU32s) getNAsI64(n int, reverse bool) []int64 {
	var l = len(b)
	if l == 0 {
		return nil
	}
	var s = make([]int64, n)
	var (
		iterN  = 0
		left   = n
		pos    = 0
		eIterN int
	)

	for i := 0; i < l; i++ {
		if iterN >= n {
			break
		}

		if reverse {
			eIterN = b[i].RIterAsI64(s, pos, left)
		} else {
			eIterN = b[i].IterAsI64(s, pos, left)
		}
		iterN += eIterN
		pos += eIterN
		left = n - iterN
	}
	return s[:iterN]
}

//BigU32 combine uint32 and 1024
type BigU32 struct {
	Start uint32
	B1024 Bit1024
}

//NewBigU32 : new big u32
func NewBigU32() *BigU32 {
	return &BigU32{
		B1024: NewBit1024(),
	}
}

//NewBigU32FromData : new big u32 from uint32 and []byte
func NewBigU32FromData(start uint32, buf []byte) (*BigU32, error) {
	var b = &BigU32{}
	b.Start = start
	b.B1024 = NewBit1024()
	var err = b.B1024.Unmarshal(buf)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//NewBigU32FromI64 : new big u32 from int64
func NewBigU32FromI64(i64 int64) (*BigU32, error) {
	if i64 < 0 || i64 >= math.MaxUint32*1024 {
		return nil, fmt.Errorf("big.u32.unsupport.i64:%+v", i64)
	}
	var u32 = uint32(i64 / C1K)
	var mod = int16(i64 % C1K)
	var b = &BigU32{
		Start: u32,
		B1024: NewBit1024(),
	}
	b.B1024.SetI16(mod)
	return b, nil
}

//Reverse : reverse bit
func (b *BigU32) Reverse() *BigU32{
	return &BigU32{
		Start: b.Start,
		B1024: b.B1024.Reverse(),
	}
}

//SetI64 set int64
func (b *BigU32) SetI64(i64 int64) error {
	if i64 < 0 || i64 >= math.MaxUint32*1024 {
		return fmt.Errorf("big.u32.set.unsupport.i64:%+v", i64)
	}
	var u32 = uint32(i64 / C1K)
	if u32 != b.Start {
		return fmt.Errorf("big.u32.set.invalid.start:%+v -- %+v", u32, i64)
	}
	var mod = int16(i64 % C1K)
	b.B1024.SetI16(mod)
	return nil
}

//GetNAsI64 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b *BigU32) GetNAsI64(n int) []int64 {
	return b.getNAsI64(n, false)
}

//RGetNAsI64 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b *BigU32) RGetNAsI64(n int) []int64 {
	return b.getNAsI64(n, true)
}

//IterAsI64 : iterate bitmap
//input: s -- an int64 array, pos -- slice position, n -- iter count
//return:
//actually iter count
func (b *BigU32) IterAsI64(s []int64, pos int, n int) int {
	return b.B1024.IterAsI64(s, pos, int64(b.Start*C1K), n)
}

//RIterAsI64 : reverse iterate bitmap
//input: s -- an int64 array, pos -- slice position, n -- iter count
//return:
//actually iter count
func (b *BigU32) RIterAsI64(s []int64, pos int, n int) int {
	return b.B1024.RIterAsI64(s, pos, int64(b.Start*C1K), n)
}

//getNAsI64 : iterate bitmap
//input: n -- iter number
//reverse -- reverse iter or not
func (b *BigU32) getNAsI64(n int, reverse bool) []int64 {
	var s = make([]int64, n)
	var iterN int
	if reverse {
		iterN = b.RIterAsI64(s, 0, n)
	} else {
		iterN = b.IterAsI64(s, 0, n)
	}
	if iterN == 0 {
		return nil
	}
	return s[:iterN]
}
