package bitmap1024

import (
	"fmt"
	"math"
)

/*
An uint32 combine with [Uint32 1024]
The max value is max uint32.
*/

const (
	MaxU32TipStart uint32 = math.MaxUint32 / C1K
)

//U32BitTips U32BitTip list
type U32BitTips []*U32BitTip

//Reverse : clone and reverse bit
func (b U32BitTips) Reverse() U32BitTips {
	var l = len(b)
	if l == 0 {
		return nil
	}
	var c = make(U32BitTips, l)
	for i := 0; i < l; i++ {
		c[i] = &U32BitTip{}
		c[i].Start = b[i].Start
		c[i].B1024 = b[i].B1024.Reverse()
	}
	return c
}

//GetNAsU32 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b U32BitTips) GetNAsU32(n int) []uint32 {
	var l = len(b)
	if l == 0 {
		return nil
	}
	var s = make([]uint32, n)
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

		eIterN = b[i].IterAsU32(s, pos, left)
		iterN += eIterN
		pos += eIterN
		left = n - iterN
	}
	return s[:iterN]
}

//RGetNAsU32 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b U32BitTips) RGetNAsU32(n int) []uint32 {
	var l = len(b)
	if l == 0 {
		return nil
	}
	var s = make([]uint32, n)
	var (
		iterN  = 0
		left   = n
		pos    = 0
		eIterN int
	)

	for i := l - 1; i >= 0; i-- {
		if iterN >= n {
			break
		}

		eIterN = b[i].RIterAsU32(s, pos, left)
		iterN += eIterN
		pos += eIterN
		left = n - iterN
	}
	return s[:iterN]
}

//U32BitTip combine uint32 and 1024
type U32BitTip struct {
	Start uint32
	B1024 Bit1024
}

//NewU32BitTip : new big u32
func NewU32BitTip() *U32BitTip {
	return &U32BitTip{
		B1024: NewBit1024(),
	}
}

//NewU32BitTipFromData : new big u32 from uint32 and []byte
func NewU32BitTipFromData(start uint32, buf []byte) (*U32BitTip, error) {
	if start > MaxU32TipStart {
		return nil, fmt.Errorf("u32.unsupport.start:%+v", start)
	}
	var b = &U32BitTip{}
	b.Start = start
	b.B1024 = NewBit1024()
	var err = b.B1024.Unmarshal(buf)
	if err != nil {
		return nil, err
	}
	return b, nil
}

//NewU32BitTipFromU32 : new u32 bit tip from uint32
func NewU32BitTipFromU32(u32 uint32) *U32BitTip {
	var start = u32 / C1K
	var mod = int16(u32 % C1K)
	var b = &U32BitTip{
		Start: start,
		B1024: NewBit1024(),
	}
	b.B1024.SetI16(mod)
	return b
}

//Reverse : reverse bit
func (b *U32BitTip) Reverse() *U32BitTip {
	return &U32BitTip{
		Start: b.Start,
		B1024: b.B1024.Reverse(),
	}
}

//SetU32 set uint32
func (b *U32BitTip) SetU32(u32 uint32) error {
	var start = u32 / C1K
	if start != b.Start {
		return fmt.Errorf("u32.set.invalid.start:%+v -- %+v", u32, u32)
	}
	var mod = int16(u32 % C1K)
	b.B1024.SetI16(mod)
	return nil
}

//GetNAsI64 : iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b *U32BitTip) GetNAsI64(n int) []uint32 {
	return b.getNAsU32(n, false)
}

//RGetNAsI64 : reverse iterate bitmap
//input: n -- iter count
//return: actually iter list
func (b *U32BitTip) RGetNAsI64(n int) []uint32 {
	return b.getNAsU32(n, true)
}

//IterAsU32 : iterate bitmap
//input: s -- an uint32 array, pos -- slice position, n -- iter count
//return:
//actually iter count
func (b *U32BitTip) IterAsU32(s []uint32, pos int, n int) int {
	return b.B1024.IterAsU32(s, pos, b.Start*C1K, n)
}

//RIterAsU32 : reverse iterate bitmap
//input: s -- an uint32 array, pos -- slice position, n -- iter count
//return:
//actually iter count
func (b *U32BitTip) RIterAsU32(s []uint32, pos int, n int) int {
	return b.B1024.RIterAsU32(s, pos, b.Start*C1K, n)
}

//getNAsU32 : iterate bitmap
//input: n -- iter number
//reverse -- reverse iter or not
func (b *U32BitTip) getNAsU32(n int, reverse bool) []uint32 {
	var s = make([]uint32, n)
	var iterN int
	if reverse {
		iterN = b.IterAsU32(s, 0, n)
	} else {
		iterN = b.RIterAsU32(s, 0, n)
	}
	if iterN == 0 {
		return nil
	}
	return s[:iterN]
}
