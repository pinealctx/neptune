package slicex

import "testing"

func TestOverwrite2FrontInPlace128Bytes(t *testing.T) {
	testOverwrite2FrontInPlace128Bytes(t, 10, -1)
	testOverwrite2FrontInPlace128Bytes(t, 10, 64)
	testOverwrite2FrontInPlace128Bytes(t, 64, 0)
	testOverwrite2FrontInPlace128Bytes(t, 64, 32)
	testOverwrite2FrontInPlace128Bytes(t, 64, 64)
}

func testOverwrite2FrontInPlace128Bytes(t *testing.T, startPos, length int) {
	t.Helper()
	x := make([]byte, 128)
	for i := 0; i < 128; i++ {
		x[i] = byte(i)
	}
	t.Logf("Before: %v", x)
	Overwrite2FrontInPlace(x, startPos, length)
	t.Logf("After: %v", x)
}

func BenchmarkOverwrite2FrontInPlace1KBytes1(b *testing.B) {
	// exceed move form second bytes to the end
	x := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		x[i] = byte(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Overwrite2FrontInPlace(x, 1, -1)
	}
}

func BenchmarkOverwrite2FrontInPlace1KBytes2(b *testing.B) {
	// exceed move form second bytes to the end
	x := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		x[i] = byte(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Overwrite2FrontInPlace(x, 1, 512)
	}
}

func BenchmarkOverwrite2FrontInPlace1KBytes3(b *testing.B) {
	// exceed move form second bytes to the end
	x := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		x[i] = byte(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Overwrite2FrontInPlace(x, 128, 128)
	}
}

func BenchmarkOverwrite2FrontInPlace1KBytes4(b *testing.B) {
	// exceed move form second bytes to the end
	x := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		x[i] = byte(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Overwrite2FrontInPlace(x, 896, 128)
	}
}

func BenchmarkOverwrite2FrontInPlace1KBytes5(b *testing.B) {
	// exceed move form second bytes to the end
	x := make([]byte, 1024)
	for i := 0; i < 1024; i++ {
		x[i] = byte(i)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Overwrite2FrontInPlace(x, 888, 128)
	}
}
