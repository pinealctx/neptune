package randx

import "testing"

func TestRandInt(t *testing.T) {
	var count = 0
	for i := 0; i < 10000; i++ {
		var x = RandInt()
		if x < 0 {
			count++
		}
	}
	t.Log(count)
}

func TestRandInt32(t *testing.T) {
	var count = 0
	for i := 0; i < 10000; i++ {
		var x = RandInt32()
		if x < 0 {
			count++
		}
	}
	t.Log(count)
}

func TestRandInt64(t *testing.T) {
	var count = 0
	for i := 0; i < 10000; i++ {
		var x = RandInt64()
		if x < 0 {
			count++
		}
	}
	t.Log(count)
}
