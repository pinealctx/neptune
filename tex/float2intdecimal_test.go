package tex

import (
	"strconv"
	"testing"
	"time"

	"github.com/pinealctx/neptune/randx"
)

func TestFloat64ToDecimalCompare1(t *testing.T) {
	count := 2000000
	_min, _max := 0.8, 0.9
	decimalPlaces := 8
	testFloat64ToDecimalCompare(t, decimalPlaces, _min, _max, count)
}

func TestFloat64ToDecimalCompare2(t *testing.T) {
	count := 2000000
	_min, _max := 3000.0, 4000.0
	decimalPlaces := 4
	testFloat64ToDecimalCompare(t, decimalPlaces, _min, _max, count)
}

func TestFloat2IntDecimalCompare1(t *testing.T) {
	count := 2000000
	_min, _max := 3000.0, 4000.0
	decimalPlaces := 4
	testFloat2IntDecimalCompare(t, decimalPlaces, _min, _max, count)
}

func TestFloat2IntDecimalCompare2(t *testing.T) {
	count := 2000000
	_min, _max := 0.8, 0.9
	decimalPlaces := 5
	testFloat2IntDecimalCompare(t, decimalPlaces, _min, _max, count)
}

func TestVerifyFloat2IntDecimal(t *testing.T) {
	testFloat2IntDecimal(t, "0.1", 1)
	testFloat2IntDecimal(t, "0.2", 1)
	testFloat2IntDecimal(t, "0.3", 1)
	testFloat2IntDecimal(t, "0.4", 1)
	testFloat2IntDecimal(t, "0.5", 1)
	testFloat2IntDecimal(t, "0.6", 1)
	testFloat2IntDecimal(t, "0.7", 1)
	testFloat2IntDecimal(t, "0.8", 1)
	testFloat2IntDecimal(t, "0.9", 1)

	testFloat2IntDecimal(t, "0.93929", 5)
	testFloat2IntDecimal(t, "0.93927", 5)
	testFloat2IntDecimal(t, "0.93926", 5)
	testFloat2IntDecimal(t, "0.93921", 5)
	testFloat2IntDecimal(t, "0.93915", 5)

	testFloat2IntDecimal(t, "0.94023", 5)
	testFloat2IntDecimal(t, "0.94025", 5)
	testFloat2IntDecimal(t, "0.94026", 5)
	testFloat2IntDecimal(t, "0.94028", 5)
	testFloat2IntDecimal(t, "0.94029", 5)

	testFloat2IntDecimal(t, "0.93900", 5)
	testFloat2IntDecimal(t, "0.93901", 5)
	testFloat2IntDecimal(t, "0.93902", 5)
	testFloat2IntDecimal(t, "0.93903", 5)
	testFloat2IntDecimal(t, "0.93904", 5)
	testFloat2IntDecimal(t, "0.93905", 5)
	testFloat2IntDecimal(t, "0.93906", 5)
	testFloat2IntDecimal(t, "0.93907", 5)
	testFloat2IntDecimal(t, "0.93908", 5)
	testFloat2IntDecimal(t, "0.93909", 5)

	testFloat2IntDecimal(t, "100000.0", 1)
	testFloat2IntDecimal(t, "2000000.0", 1)
	testFloat2IntDecimal(t, "19000000.0", 1)
	testFloat2IntDecimal(t, "3000000.0", 1)
	testFloat2IntDecimal(t, "5000000.0", 1)

	testFloat2IntDecimal(t, "1000000.0", 1)
	testFloat2IntDecimal(t, "9500000.0", 1)
	testFloat2IntDecimal(t, "10000000.0", 1)

	testFloat2IntDecimal(t, "999999999.0", 1)
	testFloat2IntDecimal(t, "999999999.1", 1)
	testFloat2IntDecimal(t, "999999999.2", 1)
	testFloat2IntDecimal(t, "999999999.3", 1)
	testFloat2IntDecimal(t, "999999999.4", 1)
	testFloat2IntDecimal(t, "999999999.5", 1)
	testFloat2IntDecimal(t, "999999999.6", 1)
	testFloat2IntDecimal(t, "999999999.7", 1)
	testFloat2IntDecimal(t, "999999999.8", 1)
	testFloat2IntDecimal(t, "999999999.9", 1)
}

func testFloat2IntDecimal(t *testing.T, fStr string, decimalPlaces int) {
	t.Helper()
	fv, _ := strconv.ParseFloat(fStr, 64)
	iv1, _ := Float64ToIntDecimalV1(fv, decimalPlaces)
	iv2, _ := Float64ToIntDecimalV2(fv, decimalPlaces)
	iv3, _ := Float64ToIntDecimalV3(fv, decimalPlaces)
	if iv1 != iv2 || iv1 != iv3 {
		t.Errorf("mismatch: %+v\n", fv)
		return
	}
	f1, _ := IntDecimalToFloat64V1(iv1, decimalPlaces)
	f2, _ := IntDecimalToFloat64V2(iv1, decimalPlaces)
	f3, _ := IntDecimalToFloat64V3(iv1, decimalPlaces)
	dc1, _ := Float64ToDecimalV1(f1, decimalPlaces)
	dc2, _ := Float64ToDecimalV1(f2, decimalPlaces)
	dc3, _ := Float64ToDecimalV1(f3, decimalPlaces)
	if !dc1.Equal(dc2) || !dc1.Equal(dc3) {
		t.Errorf("mismatch: %+v\n", fv)
		return
	}
	str1, _ := Float642String(f1, decimalPlaces)
	str2, _ := Float642String(f2, decimalPlaces)
	str3, _ := Float642String(f3, decimalPlaces)
	if str1 != str2 || str1 != str3 {
		t.Errorf("mismatch: %+v\n", fv)
		return
	}
	t.Logf("int:%v, fv:%.22f, f1:%.22f, f2:%.22f, f3:%.22f, str:%v\n", iv1, fv, f1, f2, f3, str1)
}

func testFloat64ToDecimalCompare(t *testing.T, decimalPlaces int, _min, _max float64, count int) {
	t.Helper()
	rng := randx.NewRandomChan()
	defer rng.Close()

	t1 := time.Now()
	for i := 0; i < count; i++ {
		fv := rng.Float64Range(_min, _max)
		v1, _ := Float64ToDecimalV1(fv, decimalPlaces)
		v2, _ := Float64ToDecimalV2(fv, decimalPlaces)
		if !v1.Equal(v2) {
			t.Errorf("mismatch: %+v\n", fv)
			return
		}
	}
	t2 := time.Now()
	dur := t2.Sub(t1)
	avg := dur / time.Duration(count)
	t.Logf("use time:%+v, average:%+v\n", dur, avg)
}

func testFloat2IntDecimalCompare(t *testing.T, decimalPlaces int, _min, _max float64, count int) {
	t.Helper()
	rng := randx.NewRandomChan()
	defer rng.Close()

	t1 := time.Now()
	for i := 0; i < count; i++ {
		fv := rng.Float64Range(_min, _max)
		v1, _ := Float64ToIntDecimalV1(fv, decimalPlaces)
		v2, _ := Float64ToIntDecimalV2(fv, decimalPlaces)
		v3, _ := Float64ToIntDecimalV3(fv, decimalPlaces)
		if v1 != v2 || v1 != v3 {
			t.Errorf("mismatch: %+v\n", fv)
			return
		}
	}
	t2 := time.Now()
	dur := t2.Sub(t1)
	avg := dur / time.Duration(count)
	t.Logf("use time:%+v, average:%+v\n", dur, avg)
}

// BenchmarkFloat64ToIntDecimalV1 benchmarks V1 implementation
func BenchmarkFloat64ToIntDecimalV1(b *testing.B) {
	rng := randx.NewRandomChan()
	defer rng.Close()

	// Pre-generate test data to avoid randomization overhead in benchmark
	testData := make([]float64, 10000)
	for i := range testData {
		testData[i] = rng.Float64Range(0.8, 4000.0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := testData[i%len(testData)]
		_, _ = Float64ToIntDecimalV1(val, 5)
	}
}

// BenchmarkFloat64ToIntDecimalV2 benchmarks V2 implementation
func BenchmarkFloat64ToIntDecimalV2(b *testing.B) {
	rng := randx.NewRandomChan()
	defer rng.Close()

	testData := make([]float64, 10000)
	for i := range testData {
		testData[i] = rng.Float64Range(0.8, 4000.0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := testData[i%len(testData)]
		_, _ = Float64ToIntDecimalV2(val, 5)
	}
}

// BenchmarkFloat64ToIntDecimalV3 benchmarks V3 implementation
func BenchmarkFloat64ToIntDecimalV3(b *testing.B) {
	rng := randx.NewRandomChan()
	defer rng.Close()

	testData := make([]float64, 10000)
	for i := range testData {
		testData[i] = rng.Float64Range(0.8, 4000.0)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		val := testData[i%len(testData)]
		_, _ = Float64ToIntDecimalV3(val, 5)
	}
}

// BenchmarkFloat64ToIntDecimalComparison runs all three versions for comparison
func BenchmarkFloat64ToIntDecimalComparison(b *testing.B) {
	rng := randx.NewRandomChan()
	defer rng.Close()

	testData := make([]float64, 10000)
	for i := range testData {
		testData[i] = rng.Float64Range(0.8, 4000.0)
	}

	b.Run("V1", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			val := testData[i%len(testData)]
			_, _ = Float64ToIntDecimalV1(val, 5)
		}
	})

	b.Run("V2", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			val := testData[i%len(testData)]
			_, _ = Float64ToIntDecimalV2(val, 5)
		}
	})

	b.Run("V3", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			val := testData[i%len(testData)]
			_, _ = Float64ToIntDecimalV3(val, 5)
		}
	})
}

// BenchmarkFloat64ToIntDecimalDifferentRanges tests different value ranges
func BenchmarkFloat64ToIntDecimalDifferentRanges(b *testing.B) {
	rng := randx.NewRandomChan()
	defer rng.Close()

	// Different test scenarios
	scenarios := []struct {
		name     string
		min, max float64
		decimals int
	}{
		{"SmallDecimals", 0.8, 0.9, 8},
		{"LargeNumbers", 3000.0, 4000.0, 4},
		{"MixedRange", 1.0, 1000.0, 6},
		{"HighPrecision", 0.001, 0.999, 10},
	}

	for _, scenario := range scenarios {
		testData := make([]float64, 1000)
		for i := range testData {
			testData[i] = rng.Float64Range(scenario.min, scenario.max)
		}

		b.Run(scenario.name+"/V1", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val := testData[i%len(testData)]
				_, _ = Float64ToIntDecimalV1(val, scenario.decimals)
			}
		})

		b.Run(scenario.name+"/V2", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val := testData[i%len(testData)]
				_, _ = Float64ToIntDecimalV2(val, scenario.decimals)
			}
		})

		b.Run(scenario.name+"/V3", func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				val := testData[i%len(testData)]
				_, _ = Float64ToIntDecimalV3(val, scenario.decimals)
			}
		})
	}
}

// BenchmarkFloat64ToIntDecimalMemoryAllocation tests memory allocation patterns
func BenchmarkFloat64ToIntDecimalMemoryAllocation(b *testing.B) {
	rng := randx.NewRandomChan()
	defer rng.Close()

	testData := make([]float64, 1000)
	for i := range testData {
		testData[i] = rng.Float64Range(0.8, 4000.0)
	}

	b.Run("V1-Allocs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			val := testData[i%len(testData)]
			_, _ = Float64ToIntDecimalV1(val, 5)
		}
	})

	b.Run("V2-Allocs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			val := testData[i%len(testData)]
			_, _ = Float64ToIntDecimalV2(val, 5)
		}
	})

	b.Run("V3-Allocs", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			val := testData[i%len(testData)]
			_, _ = Float64ToIntDecimalV3(val, 5)
		}
	})
}
