package bitmap1024

import (
	"testing"
	"time"

	"github.com/pinealctx/neptune/bitmap1024/internal"
)

func TestBit1024_SparseBenchIter(t *testing.T) {
	var x64s = make([]Bit1024, 64)
	for i := int16(0); i < 64; i++ {
		x64s[i] = NewBit1024()
		x64s[i].SetI16(1 << i)
	}
	t.Log()
	testBit1024IterSparse(t, "one sparse", x64s)

	for i := int16(0); i < 64; i++ {
		for j := int16(0); j < 64; j++ {
			x64s[i].SetI16(1 << j)
		}
	}
	t.Log("16 sparse on each uint64")
	testBit1024IterSparse(t, "16 sparse on each uint64", x64s)

	for i := int16(0); i < 64; i++ {
		for j := int16(0); j < 1024; j++ {
			x64s[i].SetI16(j)
		}
	}
	testBit1024IterSparse(t, "full iter", x64s)
}

func testBit1024IterSparse(t *testing.T, name string, x64s []Bit1024) {
	t.Helper()
	var (
		t1    time.Time
		t2    time.Time
		d     time.Duration
		is    []int16
		count = time.Duration(10000061)
	)

	t.Log("")
	t.Log(name)

	internal.SetSparseMagic(0)
	t1 = time.Now()
	for i := time.Duration(0); i < count; i++ {
		is = x64s[i%64].GetNAsI16(100)
	}
	t2 = time.Now()
	d = t2.Sub(t1)
	t.Log("use:", d, "average:", d/count)
	t.Log("len:", len(is), "cap:", cap(is), "is:", is)

	t1 = time.Now()
	for i := time.Duration(0); i < count; i++ {
		is = x64s[i%64].RGetNAsI16(100)
	}
	t2 = time.Now()
	d = t2.Sub(t1)
	t.Log("reverse use:", d, "average:", d/count)
	t.Log("reverse len:", len(is), "cap:", cap(is), "is:", is)

	//use sparse method
	internal.SetSparseMagic(9)
	t1 = time.Now()
	for i := time.Duration(0); i < count; i++ {
		is = x64s[i%64].GetNAsI16(100)
	}
	t2 = time.Now()
	d = t2.Sub(t1)
	t.Log("sparse mode", "use:", d, "average:", d/count)
	t.Log("len:", len(is), "cap:", cap(is), "is:", is)

	t1 = time.Now()
	for i := time.Duration(0); i < count; i++ {
		is = x64s[i%64].RGetNAsI16(100)
	}
	t2 = time.Now()
	d = t2.Sub(t1)
	t.Log("sparse mode", "use:", d, "average:", d/count)
	t.Log("reverse len:", len(is), "cap:", cap(is), "is:", is)
}
