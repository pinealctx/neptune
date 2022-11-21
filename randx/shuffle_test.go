package randx

import (
	"bytes"
	"sort"
	"testing"
)

type btsT []byte

func (x btsT) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x btsT) Len() int {
	return len(x)
}

func TestShuffle(t *testing.T) {
	var x = sort.IntSlice{
		0, 1, 2, 3, 4, 5,
	}
	Shuffle(x)
	t.Log(x)
}

func TestShuffleX(t *testing.T) {
	var x = btsT{0, 1, 2}
	for i := 0; i < 8; i++ {
		var y = make(btsT, 3)
		copy(y, x)
		Shuffle(y)
		t.Log(y)
	}
}

func TestShuffleY(t *testing.T) {
	var x1 = btsT{0, 1, 2}
	var x2 = btsT{0, 2, 1}
	var x3 = btsT{1, 0, 2}
	var x4 = btsT{1, 2, 0}
	var x5 = btsT{2, 0, 1}
	var x6 = btsT{2, 1, 0}
	var c1, c2, c3, c4, c5, c6 int
	for i := 0; i < 100000; i++ {
		var y = make(btsT, 3)
		copy(y, x1)
		Shuffle(y)
		if bytes.Compare(x1, y) == 0 {
			c1++
		}
		if bytes.Compare(x2, y) == 0 {
			c2++
		}
		if bytes.Compare(x3, y) == 0 {
			c3++
		}
		if bytes.Compare(x4, y) == 0 {
			c4++
		}
		if bytes.Compare(x5, y) == 0 {
			c5++
		}
		if bytes.Compare(x6, y) == 0 {
			c6++
		}
	}
	t.Log(c1 + c2 + c3 + c4 + c5 + c6)
	t.Log(float64(c1) / 100000)
	t.Log(float64(c2) / 100000)
	t.Log(float64(c3) / 100000)
	t.Log(float64(c4) / 100000)
	t.Log(float64(c5) / 100000)
	t.Log(float64(c6) / 100000)
}
