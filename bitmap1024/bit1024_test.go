package bitmap1024

import (
	"testing"
	"time"
)

func TestBit1024_SetI16(t *testing.T) {
	var x = NewBit1024()
	for i := int16(0); i < 1024; i++ {
		x.SetI16(i)
		t.Log("len:", x.Len(), "NLen:", x.NLen())
	}
}

func TestBit1024_Iter(t *testing.T) {
	var a = NewBit1024()
	var b = NewBit1024()
	a.SetI16(0)
	a.SetI16(1)
	a.SetI16(2)
	a.SetI16(3)
	a.SetI16(4)
	a.SetI16(5)
	a.SetI16(6)
	a.SetI16(7)
	a.SetI16(8)
	a.SetI16(9)
	a.SetI16(10)

	b.SetI16(1023)
	b.SetI16(1022)
	b.SetI16(1021)
	b.SetI16(1020)
	b.SetI16(1019)
	b.SetI16(1018)
	b.SetI16(1017)
	b.SetI16(1016)

	t.Log(1023 - 1016 + 1 + 10 + 1)
	t.Log(1023 - 1016 + 1 + 10 + 1 - 1024)
	var c = a.OrThenReverse(b)
	var x = make([]int16, 30)
	c.IterAsI16(x, 0, 0, 30)
	t.Log(x)
	t.Log(len(x))
	x = make([]int16, 30)
	c.RIterAsI16(x, 0, 0, 30)
	t.Log(x)
	t.Log(len(x))

	t.Log(c.GetNAsI16(100))
	t.Log(c.RGetNAsI16(100))

	t.Log(len(c.GetNAsI16(1024)))
	t.Log(len(c.RGetNAsI16(1024)))
}

func TestBit1024_BenchSetUnset(t *testing.T) {
	var x = NewBit1024()
	var count = time.Duration(10000000)
	var t1 = time.Now()
	for i := time.Duration(0); i < count; i++ {
		x.SetI16(int16(i) % 1024)
		x.UnsetI16(int16(i) % 1024)
	}
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use:", d, "average:", d/count)
}

func TestBit1024_BenchIter(t *testing.T) {
	var x = NewBit1024()
	var count = time.Duration(10000000)
	var is []int16
	var t1 = time.Now()
	for i := time.Duration(0); i < count; i++ {
		x.SetI16(int16(i) % 1024)
		x.UnsetI16(int16(i+1) % 1024)
		is = x.GetNAsI16(100)
	}
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use:", d, "average:", d/count)
	t.Log("islen:", len(is), "cap:", cap(is), "is:", is)
}

func TestBit1024_BenchReverseIter(t *testing.T) {
	var x = NewBit1024()
	var count = time.Duration(10000000)
	var is []int16
	var t1 = time.Now()
	for i := time.Duration(0); i < count; i++ {
		x.SetI16(int16(i) % 1024)
		x.UnsetI16(int16(i+1) % 1024)
		is = x.RGetNAsI16(100)
	}
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use:", d, "average:", d/count)
	t.Log("islen:", len(is), "cap:", cap(is), "is:", is)
}

func TestBit1024_BenchMake(t *testing.T) {
	var count = time.Duration(1000000)
	var c []int64
	var t1 = time.Now()
	for i := time.Duration(0); i < count; i++ {
		c = make([]int64, 100)
	}
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use:", d, "average:", d/count)
	t.Log(len(c))
}

func TestBit1024_BenchMarshal(t *testing.T) {
	var x = NewBit1024()
	var buf []byte
	var count = time.Duration(16 * 1024)
	var t1 = time.Now()
	for k := time.Duration(0); k < count; k++ {
		for i := int16(0); i < 1024; i++ {
			x.SetI16(i)
			buf = x.Marshal()
			var y = NewBit1024()
			var err = y.Unmarshal(buf)
			if err != nil {
				panic(err)
			}
			if !x.Equal(y) {
				panic(y)
			}
		}
	}
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/(count*1024))
}

func TestBit1024_ReverseBenchMarshal(t *testing.T) {
	var x = NewBit1024()
	var buf []byte
	var count = time.Duration(16 * 1024)
	var t1 = time.Now()
	for k := time.Duration(0); k < count; k++ {
		x.SetI16(int16(k % 1024))
		for i := int16(0); i < 1024; i++ {
			buf = x.Marshal()
			var y = NewBit1024()
			var err = y.Unmarshal(buf)
			if err != nil {
				panic(err)
			}
			if !x.Equal(y) {
				panic(y)
			}
		}
	}
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/(count*1024))
}
