package bitmap1024

import "testing"

func TestNewU32BitTipFromU32(t *testing.T) {
	x1 := NewU32BitTipFromU32(0)
	for i := uint32(3); i < 1000; i++ {
		err := x1.SetU32(i)
		if err != nil {
			panic(err)
		}
	}

	x2 := NewU32BitTipFromU32(2047)
	for i := uint32(1035); i < 2000; i++ {
		err := x2.SetU32(i)
		if err != nil {
			panic(err)
		}
	}

	var xs = U32BitTips{x1.Reverse(), x2.Reverse()}
	var r = xs.GetNAsU32(100)
	t.Log("len:", len(r), "cap:", cap(r))
	t.Log(r)
	r = xs.RGetNAsU32(100)
	t.Log("len:", len(r), "cap:", cap(r))
	t.Log(r)
}

func TestNewU32BitTipIter(t *testing.T) {
	x1 := NewU32BitTipFromU32(0)
	for i := uint32(3); i < 1000; i++ {
		err := x1.SetU32(i)
		if err != nil {
			panic(err)
		}
	}

	x2 := NewU32BitTipFromU32(2047)
	for i := uint32(1035); i < 1980; i++ {
		err := x2.SetU32(i)
		if err != nil {
			panic(err)
		}
	}

	var xs = U32BitTips{x1.Reverse(), x2.Reverse()}
	var r = xs.GetNAsU32(100)
	t.Log("len:", len(r), "cap:", cap(r))
	t.Log(r)
	r = xs.RGetNAsU32(100)
	t.Log("len:", len(r), "cap:", cap(r))
	t.Log(r)
}

