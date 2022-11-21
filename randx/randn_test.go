package randx

import "testing"

func TestRandBetween(t *testing.T) {
	testRandBetween(t, 0, 9)
	testRandBetween(t, 1, 10)
	testRandBetween(t, -4, 5)
}

func testRandBetween(t *testing.T, min, max int) {
	var count = 0
	for i := 0; i < 10000; i++ {
		var x = RandBetween(min, max)
		if x == 1 {
			count++
		}
	}
	t.Log(float64(count) / 10000)
}

func TestSimpleRandBetween(t *testing.T) {
	testSimpleRandBetween(t, 0, 9)
	testSimpleRandBetween(t, 1, 10)
	testSimpleRandBetween(t, -4, 5)
}

func testSimpleRandBetween(t *testing.T, min, max int) {
	var count = 0
	for i := 0; i < 10000; i++ {
		var x = SimpleRandBetween(min, max)
		if x == 1 {
			count++
		}
	}
	t.Log(float64(count) / 10000)
}

func TestRandBetweenSecure(t *testing.T) {
	testRandBetweenSecure(t, 0, 9)
	testRandBetweenSecure(t, 1, 10)
	testRandBetweenSecure(t, -4, 5)
}

func testRandBetweenSecure(t *testing.T, min, max int) {
	var count = 0
	for i := 0; i < 10000; i++ {
		var x = RandBetweenSecure(min, max)
		if x == 1 {
			count++
		}
	}
	t.Log(float64(count) / 10000)
}
