package bytex

import (
	"math"
	"testing"
	"time"
)

func TestBufferX_ReWrite1(t *testing.T) {
	var x = NewSizedBufferX(4)
	x.WriteU32(1)
	//used
	x.WriteU32(2)
	x.ReWriteU32(0, 10)
	var v, err = x.ReadU32()
	if err != nil {
		panic(err)
	}
	t.Log(v)
	v, err = x.ReadU32()
	if err != nil {
		panic(err)
	}
	t.Log(v)
}

func TestBufferX_ReWrite2(t *testing.T) {
	var x = NewSizedBufferX(4)
	x.WriteU32(1)
	x.ReWriteU32(0, 10)
	//used
	x.WriteU32(2)
	var v, err = x.ReadU32()
	if err != nil {
		panic(err)
	}
	t.Log(v)
	v, err = x.ReadU32()
	if err != nil {
		panic(err)
	}
	t.Log(v)
}

// 87ns
func TestBufferX_WriteVarI32(t *testing.T) {
	var i []int
	var count time.Duration = 1000000

	var t1 = time.Now()
	for k := time.Duration(0); k < count; k++ {
		i = writeVarI32()
	}
	var t2 = time.Now()
	t.Log(i)
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", (d/count)/time.Duration(len(i)))
}

func writeVarI32() []int {
	var i32 = []int32{
		math.MinInt32,
		math.MinInt16,
		math.MinInt8,
		-1,
		0,
		1,
		math.MaxInt32,
		math.MaxInt16,
		math.MaxInt8,

		//8191
		8191,
		//1万
		100 * 100,
		//10万
		10 * 100 * 100,
		//100万
		100 * 100 * 100,
	}
	var l = len(i32)
	var s = make([]int, l)
	for i := 0; i < l; i++ {
		var x = NewSizedBufferX(6)
		x.WriteVarI32(i32[i])
		s[i] = x.Len()
		var r, _ = x.ReadVarI32()
		if r != i32[i] {
			panic(r)
		}
	}
	return s
}

// 41ns
func TestBufferX_WriteI32(t *testing.T) {
	var i []int
	var count time.Duration = 1000000

	var t1 = time.Now()
	for k := time.Duration(0); k < count; k++ {
		i = writeI32()
	}
	var t2 = time.Now()
	t.Log(i)
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", (d/count)/time.Duration(len(i)))
}

func writeI32() []int {
	var i32 = []int32{
		math.MinInt32,
		math.MinInt16,
		math.MinInt8,
		-1,
		0,
		1,
		math.MaxInt32,
		math.MaxInt16,
		math.MaxInt8,

		//8191
		8191,
		//1万
		100 * 100,
		//10万
		10 * 100 * 100,
		//100万
		100 * 100 * 100,
	}
	var l = len(i32)
	var s = make([]int, l)
	for i := 0; i < l; i++ {
		var x = NewSizedBufferX(6)
		x.WriteI32(i32[i])
		s[i] = x.Len()
		var r, _ = x.ReadI32()
		if r != i32[i] {
			panic(r)
		}
	}
	return s
}
