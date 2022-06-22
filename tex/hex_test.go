package tex

import (
	"math"
	"testing"
)

func TestI64Hex(t *testing.T) {
	var x int64 = math.MaxInt64
	testI64(t, x)
	x = -1
	testI64(t, x)

	x = math.MaxInt32
	testI64(t, x)
}

func testI64(t *testing.T, i int64) {
	var h1 = I64Hex(i)
	var h2 = I64HexV2(i)
	var h3 = U64Hex(uint64(i))
	var h4 = U64HexV2(uint64(i))

	t.Log("hex", h1)
	t.Log("base32", h2)
	t.Log("hex", h3)
	t.Log("base32", h4)

	var (
		ri  int64
		ru  uint64
		err error
	)
	ri, err = HexI64(h1)
	if err != nil || ri != i {
		panic(err)
	}

	ri, err = HexI64V2(h2)
	if err != nil || ri != i {
		panic(err)
	}

	ru, err = HexU64(h3)
	if err != nil || ru != uint64(i) {
		panic(err)
	}

	ru, err = HexU64V2(h4)
	if err != nil || ru != uint64(i) {
		panic(err)
	}
}
