package randx

import (
	"testing"
)

func TestRandIntSecure(t *testing.T) {
	var count = 0
	for i := 0; i < 10000; i++ {
		var x = RandIntSecure()
		if x < 0 {
			count++
		}
	}
	t.Log(count)
}

func TestRandInt32Secure(t *testing.T) {
	var count = 0
	for i := 0; i < 10000; i++ {
		var x = RandInt32Secure()
		if x < 0 {
			count++
		}
	}
	t.Log(count)
}

func TestRandInt64Secure(t *testing.T) {
	var count = 0
	for i := 0; i < 10000; i++ {
		var x = RandInt64Secure()
		if x < 0 {
			count++
		}
	}
	t.Log(count)
}
