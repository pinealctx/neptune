package bitmap1024

import "testing"

func TestBit64_Set(t *testing.T) {
	var x Bit64

	t.Logf("%064b", x)
	t.Log("len:", x.Len())
	t.Log("full:", x.Full())

	for i := byte(0); i < 64; i++ {
		x.Set(i)
		t.Logf("%064b", x)
		t.Log("len:", x.Len())
		t.Log("NLen:", x.NLen())
		t.Log("full:", x.Full())
	}

	for i := byte(0); i < 64; i++ {
		x.Unset(i)
		t.Logf("%064b", x)
		t.Log("len:", x.Len())
		t.Log("NLen:", x.NLen())
		t.Log("full:", x.Full())
	}
}

func TestBit64_Iter(t *testing.T) {
	var a Bit64
	a.Set(0)
	a.Set(1)
	a.Set(2)
	a.Set(3)
	a.Set(4)
	a.Set(5)
	a.Set(6)
	a.Set(7)
	a.Set(8)
	a.Set(9)
	a.Set(10)

	a.Set(63)
	a.Set(62)
	a.Set(61)
	a.Set(60)
	a.Set(59)
	a.Set(58)
	a.Set(57)
	a.Set(56)

	a = a.Reverse()

	var x = make([]int64, 0, 30)
	a.IterAsI64(x, 0, 0, 30)
	t.Log(x)
	x = make([]int64, 0, 30)
	a.RIterAsI64(x, 0, 0,30)
	t.Log(x)
}

