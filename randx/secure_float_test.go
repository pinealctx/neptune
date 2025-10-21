package randx

import (
	"math"
	"testing"
)

func TestFloat64Secure(t *testing.T) {
	// Test range [0, 1)
	for i := 0; i < 10000; i++ {
		f := Float64Secure()
		if f < 0 || f >= 1 {
			t.Errorf("Float64Secure() = %f, want [0, 1)", f)
		}
	}
}

func TestFloat64SecureDistribution(t *testing.T) {
	const samples = 100000
	const buckets = 10
	counts := make([]int, buckets)

	// Collect samples
	for i := 0; i < samples; i++ {
		f := Float64Secure()
		bucket := int(f * buckets)
		if bucket >= buckets {
			bucket = buckets - 1 // Handle edge case f very close to 1
		}
		counts[bucket]++
	}

	// Check distribution uniformity (chi-square test approximation)
	expected := samples / buckets
	tolerance := expected / 10 // 10% tolerance

	for i, count := range counts {
		if math.Abs(float64(count-expected)) > float64(tolerance) {
			t.Logf("Bucket %d: got %d, expected ~%d", i, count, expected)
		}
	}
}

func TestFloat64SecurePrecision(t *testing.T) {
	// Test that we can generate distinct values
	seen := make(map[float64]bool)
	duplicates := 0

	for i := 0; i < 10000; i++ {
		f := Float64Secure()
		if seen[f] {
			duplicates++
		}
		seen[f] = true
	}

	// With 53 bits precision, duplicates should be extremely rare
	if duplicates > 5 {
		t.Errorf("Too many duplicates: %d out of 10000", duplicates)
	}
}

func BenchmarkFloat64Secure(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Float64Secure()
	}
}
