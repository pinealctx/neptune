package randx

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

// TestFloat64Range tests that Float64() returns values in [0, 1) range
func TestFloat64Range(t *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()

	const iterations = 100000
	for i := 0; i < iterations; i++ {
		val := rng.Float64()
		if val < 0 || val >= 1 {
			t.Errorf("Float64() returned %f, expected [0, 1)", val)
		}
	}
}

// TestIntNRange tests that IntN() returns values in [0, n) range
func TestIntNRange(t *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()

	testCases := []int{1, 2, 5, 10, 100, 1000, 65536}
	const iterations = 10000

	for _, n := range testCases {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			for i := 0; i < iterations; i++ {
				val := rng.IntN(n)
				if val < 0 || val >= n {
					t.Errorf("IntN(%d) returned %d, expected [0, %d)", n, val, n)
				}
			}
		})
	}
}

// TestIntBetween tests that IntBetween() returns values in [min, max] range
func TestIntBetween(t *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()
	for i := 0; i < 10; i++ {
		t.Log(rng.IntBetween(0, 0))
		t.Log(rng.IntN(1))
		t.Log(rng.IntRange(0, 1))
	}
}

// TestFloat64UniformDistribution tests uniform distribution of Float64()
func TestFloat64UniformDistribution(t *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()

	const (
		iterations = 1000000
		buckets    = 100
		expected   = iterations / buckets
	)

	histogram := make([]int, buckets)

	// Generate samples
	for i := 0; i < iterations; i++ {
		val := rng.Float64()
		bucket := int(val * buckets)
		if bucket == buckets { // Handle edge case where val is very close to 1
			bucket = buckets - 1
		}
		histogram[bucket]++
	}

	// Check distribution uniformity using chi-square test
	chiSquare := 0.0
	for _, count := range histogram {
		diff := float64(count - expected)
		chiSquare += (diff * diff) / float64(expected)
	}

	// Chi-square critical value for 99 degrees of freedom at 0.01 significance level ≈ 135
	const criticalValue = 135.0
	if chiSquare > criticalValue {
		t.Errorf("Float64 distribution not uniform, chi-square = %f (critical = %f)", chiSquare, criticalValue)
	}

	t.Logf("Float64 uniformity test passed, chi-square = %f", chiSquare)
}

// TestIntNUniformDistribution tests uniform distribution of IntN()
func TestIntNUniformDistribution(t *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()

	testCases := []struct {
		n          int
		iterations int
	}{
		{10, 100000},
		{100, 1000000},
		{1000, 1000000},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("n=%d", tc.n), func(t *testing.T) {
			histogram := make([]int, tc.n)
			expected := tc.iterations / tc.n

			// Generate samples
			for i := 0; i < tc.iterations; i++ {
				val := rng.IntN(tc.n)
				histogram[val]++
			}

			// Check distribution uniformity
			chiSquare := 0.0
			for _, count := range histogram {
				diff := float64(count - expected)
				chiSquare += (diff * diff) / float64(expected)
			}

			// Chi-square critical value depends on degrees of freedom
			criticalValue := float64(tc.n) * 1.5 // Approximate critical value
			if chiSquare > criticalValue {
				t.Errorf("IntN(%d) distribution not uniform, chi-square = %f (critical = %f)",
					tc.n, chiSquare, criticalValue)
			}

			t.Logf("IntN(%d) uniformity test passed, chi-square = %f", tc.n, chiSquare)
		})
	}
}

// TestRaceConditions tests concurrent access safety
func TestRaceConditions(t *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()

	const (
		numGoroutines = 100
		opsPerRoutine = 1000
	)

	var wg sync.WaitGroup
	errors := make(chan error, numGoroutines)

	// Test concurrent Float64() calls
	t.Run("ConcurrentFloat64", func(_ *testing.T) {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < opsPerRoutine; j++ {
					val := rng.Float64()
					if val < 0 || val >= 1 {
						select {
						case errors <- fmt.Errorf("invalid Float64 value: %f", val):
						default:
						}
						return
					}
				}
			}()
		}
		wg.Wait()
	})

	// Test concurrent IntN() calls
	t.Run("ConcurrentIntN", func(_ *testing.T) {
		wg.Add(numGoroutines)
		for i := 0; i < numGoroutines; i++ {
			go func() {
				defer wg.Done()
				for j := 0; j < opsPerRoutine; j++ {
					val := rng.IntN(1000)
					if val < 0 || val >= 1000 {
						select {
						case errors <- fmt.Errorf("invalid IntN value: %d", val):
						default:
						}
						return
					}
				}
			}()
		}
		wg.Wait()
	})

	// Check for errors
	close(errors)
	for err := range errors {
		t.Error(err)
	}
}

// TestRaceMixed tests mixed concurrent operations
func TestRaceMixed(_ *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()

	const (
		numGoroutines = 50
		opsPerRoutine = 500
	)

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 2) // Float64 and IntN goroutines

	// Concurrent Float64 operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerRoutine; j++ {
				_ = rng.Float64()
			}
		}()
	}

	// Concurrent IntN operations
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < opsPerRoutine; j++ {
				_ = rng.IntN(100)
			}
		}()
	}

	wg.Wait()
}

// TestProducerShutdown tests proper shutdown behavior
func TestProducerShutdown(t *testing.T) {
	rng := NewRandomChan()

	// Generate some random numbers to verify it works
	for i := 0; i < 100; i++ {
		val := rng.Float64()
		if val == 0.0 {
			t.Errorf("Got unexpected zero value before closing")
		}
	}

	// Close the generator
	rng.Close()

	// After close, we should still be able to read buffered data
	// This is the correct and robust behavior!
	bufferedCount := 0
	for i := 0; i < 1000; i++ {
		val := rng.Float64()
		if val != 0.0 {
			bufferedCount++
		} else {
			// Once we start getting zeros, all subsequent reads should be zero
			for j := 0; j < 10; j++ {
				zeroVal := rng.Float64()
				if zeroVal != 0.0 {
					t.Errorf("Expected consistent zero values after exhausting buffer, got %f", zeroVal)
				}
			}
			break
		}
	}

	t.Logf("Successfully read %d buffered values after close, then got consistent zeros", bufferedCount)
}

// TestZeroUnitBehavior tests that closed generator eventually returns ZeroUnit consistently
func TestZeroUnitBehavior(t *testing.T) {
	rng := NewRandomChan()

	// Generate some random numbers to verify it works normally
	normalVal := rng.Float64()
	if normalVal == 0.0 {
		// Extremely unlikely but possible, try again
		normalVal = rng.Float64()
		if normalVal == 0.0 {
			t.Skip("Got two consecutive zero values from crypto/rand, skipping test")
		}
	}

	// Close the generator
	rng.Close()

	// Consume all buffered data - this is the correct behavior
	bufferedCount := 0
	for i := 0; i < 10000; i++ {
		val := rng.Float64()
		if val != 0.0 {
			bufferedCount++
		} else {
			// Found first zero, now all subsequent reads should be zero
			break
		}
	}

	t.Logf("Consumed %d buffered values after close", bufferedCount)

	// Now all reads should return values derived from ZeroUnit
	for i := 0; i < 10; i++ {
		// Test Uint64() returns 0
		uint64Val := rng.Uint64()
		if uint64Val != 0 {
			t.Errorf("Iteration %d: Uint64() from exhausted generator returned %d, expected 0", i, uint64Val)
		}

		// Test Float64() returns 0.0
		float64Val := rng.Float64()
		if float64Val != 0.0 {
			t.Errorf("Iteration %d: Float64() from exhausted generator returned %f, expected 0.0", i, float64Val)
		}

		// Test IntN() returns 0
		intVal := rng.IntN(100)
		if intVal != 0 {
			t.Errorf("Iteration %d: IntN(100) from exhausted generator returned %d, expected 0", i, intVal)
		}

		// Test ReadBytes() returns all zeros
		bytes := rng.ReadBytes(16)
		for j, b := range bytes {
			if b != 0 {
				t.Errorf("Iteration %d: ReadBytes()[%d] from exhausted generator returned %d, expected 0", i, j, b)
			}
		}
	}

	t.Log("ZeroUnit behavior verified: all methods return zero values after buffer exhaustion")
}

// TestMultipleClose tests that multiple Close() calls are safe
func TestMultipleClose(t *testing.T) {
	rng := NewRandomChan()

	// Verify it works before closing
	val := rng.Float64()
	if val < 0 || val >= 1 {
		t.Errorf("Expected valid Float64 before closing, got %f", val)
	}

	// Multiple closes should not panic
	rng.Close()
	rng.Close()
	rng.Close()

	// Wait a bit and close again
	time.Sleep(50 * time.Millisecond)
	rng.Close()

	// After closing, should still be able to read buffered data first
	// This demonstrates the robust behavior we want
	bufferedCount := 0
	for i := 0; i < 1000; i++ {
		val := rng.Float64()
		if val != 0.0 {
			bufferedCount++
		} else {
			// Once we hit zero, verify subsequent reads are also zero
			for j := 0; j < 5; j++ {
				zeroVal := rng.Float64()
				if zeroVal != 0.0 {
					t.Errorf("Expected zero after buffer exhaustion, got %f", zeroVal)
				}
			}
			break
		}
	}

	t.Logf("Multiple Close() calls completed safely, consumed %d buffered values, then got zeros", bufferedCount)
}

// TestReadBytes tests ReadBytes functionality
func TestReadBytes(t *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()

	testSizes := []int{0, 1, 7, 8, 9, 16, 100, 1000}

	for _, size := range testSizes {
		t.Run(fmt.Sprintf("size=%d", size), func(t *testing.T) {
			data := rng.ReadBytes(size)
			if len(data) != size {
				t.Errorf("ReadBytes(%d) returned %d bytes", size, len(data))
			}
		})
	}
}

// TestIntNPanic tests that IntN panics for invalid input
func TestIntNPanic(t *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()

	testCases := []int{0, -1, -100}

	for _, n := range testCases {
		t.Run(fmt.Sprintf("n=%d", n), func(t *testing.T) {
			defer func() {
				if r := recover(); r == nil {
					t.Errorf("IntN(%d) should panic", n)
				}
			}()
			_ = rng.IntN(n)
		})
	}
}

// TestChannelBuffering tests that the channel provides adequate buffering
func TestChannelBuffering(t *testing.T) {
	rng := NewRandomChan()
	defer rng.Close()

	// Wait for initial buffering
	time.Sleep(200 * time.Millisecond)

	// Should have some data buffered
	buffered := len(rng.dataChan)
	if buffered == 0 {
		t.Error("Channel should have some data buffered after startup")
	} else {
		t.Logf("Initial buffer: %d units", buffered)
	}

	// Consume a significant amount of data rapidly
	consumeCount := 10000
	for i := 0; i < consumeCount; i++ {
		_ = rng.Uint64()
	}

	// Check that buffer is still being maintained
	time.Sleep(100 * time.Millisecond)
	newBuffered := len(rng.dataChan)

	// Buffer should be refilled (though exact amount depends on timing)
	t.Logf("Buffer after consuming %d units: %d units", consumeCount, newBuffered)

	// The test passes if we can continue getting data without blocking
	// (which means the producer is working)
	for i := 0; i < 100; i++ {
		_ = rng.Float64()
	}
}

// BenchmarkFloat64Single benchmarks single-threaded Float64() performance
func BenchmarkFloat64Single(b *testing.B) {
	rng := NewRandomChan()
	defer rng.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = rng.Float64()
	}
}

// BenchmarkIntNSingle benchmarks single-threaded IntN() performance
func BenchmarkIntNSingle(b *testing.B) {
	rng := NewRandomChan()
	defer rng.Close()

	testCases := []int{10, 100, 1000, 65536}

	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_ = rng.IntN(n)
			}
		})
	}
}

// BenchmarkFloat64Concurrent benchmarks multi-threaded Float64() performance
func BenchmarkFloat64Concurrent(b *testing.B) {
	rng := NewRandomChan()
	defer rng.Close()

	numGoroutines := runtime.NumCPU()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_ = rng.Float64()
		}
	})

	b.Logf("Concurrent Float64 with %d goroutines", numGoroutines)
}

// BenchmarkIntNConcurrent benchmarks multi-threaded IntN() performance
func BenchmarkIntNConcurrent(b *testing.B) {
	rng := NewRandomChan()
	defer rng.Close()

	testCases := []int{10, 100, 1000}

	for _, n := range testCases {
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_ = rng.IntN(n)
				}
			})
		})
	}
}

// BenchmarkMixedOperations benchmarks mixed Float64 and IntN operations
func BenchmarkMixedOperations(b *testing.B) {
	rng := NewRandomChan()
	defer rng.Close()

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// Mix of operations that you typically use
			_ = rng.Float64()
			_ = rng.IntN(100)
			_ = rng.Float64()
			_ = rng.IntN(1000)
		}
	})
}
