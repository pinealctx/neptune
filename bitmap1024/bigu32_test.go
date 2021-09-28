package bitmap1024

import (
	"math"
	"testing"
)

func TestBigU32(t *testing.T) {
	var x int64 = math.MinInt32 * 1024 - 1
	t.Log(x)
	x = math.MinInt32
	t.Log(x)
	t.Log(x%1024)
	for i := int64(0); i < 50; i++ {
		x -= i
		t.Log(x)
		t.Log(x%1024)
	}
	t.Log(math.MaxUint32)
}
