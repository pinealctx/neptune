package tex

import (
	"fmt"
	"math"
	"testing"
)

func TestVarint64Encoding(t *testing.T) {
	testCases := []int64{
		0, 1, -1, 127, 128, -128, 16383, 16384, -16384,
		2097151, 2097152, -2097152,
		9223372036854775807,  // max int64
		-9223372036854775808, // min int64
	}

	for _, original := range testCases {
		// Test encoding and decoding
		encoded := EncodeVarint64(original)
		decoded, n, err := DecodeVarint64(encoded)

		if err != nil {
			t.Errorf("Failed to decode varint for value %d: %v", original, err)
			continue
		}

		if decoded != original {
			t.Errorf("Varint encoding/decoding mismatch: original=%d, decoded=%d", original, decoded)
		}

		if n != len(encoded) {
			t.Errorf("Bytes consumed mismatch: expected=%d, actual=%d", len(encoded), n)
		}

		// Test size calculation
		expectedSize := VarintSize(original)
		if len(encoded) != expectedSize {
			t.Errorf("Size calculation mismatch for %d: expected=%d, actual=%d",
				original, expectedSize, len(encoded))
		}

		t.Logf("Value %d encoded to %d bytes: %v", original, len(encoded), encoded)
	}
}

func TestUvarint64Encoding(t *testing.T) {
	testCases := []uint64{
		0, 1, 127, 128, 16383, 16384,
		2097151, 2097152,
		18446744073709551615, // max uint64
	}

	for _, original := range testCases {
		// Test encoding and decoding
		encoded := EncodeUvarint64(original)
		decoded, n, err := DecodeUvarint64(encoded)

		if err != nil {
			t.Errorf("Failed to decode uvarint for value %d: %v", original, err)
			continue
		}

		if decoded != original {
			t.Errorf("Uvarint encoding/decoding mismatch: original=%d, decoded=%d", original, decoded)
		}

		if n != len(encoded) {
			t.Errorf("Bytes consumed mismatch: expected=%d, actual=%d", len(encoded), n)
		}

		// Test size calculation
		expectedSize := UvarintSize(original)
		if len(encoded) != expectedSize {
			t.Errorf("Size calculation mismatch for %d: expected=%d, actual=%d",
				original, expectedSize, len(encoded))
		}

		t.Logf("Value %d encoded to %d bytes: %v", original, len(encoded), encoded)
	}
}

func TestEncodeUvarint64Spec(_ *testing.T) {
	xs := []uint64{0, 127, 128, 255, 256, 511, 512, 1023, 1024, 4095, 4096, 16 * 1024, 64 * 1024,
		1024 * 1024, 2*1024*1024 - 1, 2 * 1024 * 1024, 4 * 1024 * 1024}
	for _, x := range xs {
		testEncodeUvarint64(x)
	}
}

func TestMaxUvarint64(_ *testing.T) {
	testEncodeUvarint64(127)
	testEncodeUvarint64(128)
	testEncodeUvarint64(math.MaxUint8 - 1)
	testEncodeUvarint64(math.MaxUint8)
	testEncodeUvarint64(math.MaxUint16 - 1)
	testEncodeUvarint64(math.MaxUint16)
	testEncodeUvarint64(math.MaxUint32 - 1)
	testEncodeUvarint64(math.MaxUint32)
	testEncodeUvarint64(math.MaxUint64 - 1)
	testEncodeUvarint64(math.MaxUint64)
	fmt.Println(uint64(math.MaxUint64))
}

func TestMustValidDecode(t *testing.T) {
	// max 1 byte valid value
	testMustInvalidUvarint64(t, false, []byte{0x7F})
	// max 2 bytes valid value
	testMustInvalidUvarint64(t, false, []byte{0xFF, 0x7F})
	// max 3 bytes valid value
	testMustInvalidUvarint64(t, false, []byte{0xFF, 0xFF, 0x7F})
	// max 4 bytes valid value
	testMustInvalidUvarint64(t, false, []byte{0xFF, 0xFF, 0xFF, 0x7F})
	// max 5 bytes valid value
	testMustInvalidUvarint64(t, false, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x7F})
	// max 6 bytes valid value
	testMustInvalidUvarint64(t, false, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F})
	// max 7 bytes valid value
	testMustInvalidUvarint64(t, false, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F})
	// max 8 bytes valid value
	testMustInvalidUvarint64(t, false, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F})
	// max 9 bytes valid value
	testMustInvalidUvarint64(t, false, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x7F})
}

func TestInvalidEncoding(t *testing.T) {
	// Test nil buffer
	testMustInvalidUvarint64(t, true, nil)
	// Test empty buffer
	testMustInvalidUvarint64(t, true, []byte{})
	// Test invalid varint (incomplete)
	testMustInvalidUvarint64(t, true, []byte{0xFF})
	// Test invalid varint (incomplete), only MSB set
	testMustInvalidUvarint64(t, true, []byte{0x80})
	// Test invalid varint (too long)
	testMustInvalidUvarint64(t, true, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF})
	testMustInvalidUvarint64(t, true, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00})
	testMustInvalidUvarint64(t, true, []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80})
	testMustInvalidUvarint64(t, true, []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})

	testMustInvalidUvarint64(t, true, []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x00})
	testMustInvalidUvarint64(t, false, []byte{0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x80, 0x01})
	testMustInvalidUvarint64(t, false, []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x01})
}

func TestEncodingEfficiency(t *testing.T) {
	// Test that small numbers use fewer bytes
	testCases := map[int64]int{
		0:     1, // 0 should use 1 byte
		127:   1, // 127 should use 1 byte
		128:   2, // 128 should use 2 bytes
		16383: 2, // 16383 should use 2 bytes
		16384: 3, // 16384 should use 3 bytes
	}

	for value, expectedBytes := range testCases {
		encoded := EncodeVarint64(value)
		if len(encoded) != expectedBytes {
			t.Errorf("Value %d: expected %d bytes, got %d bytes",
				value, expectedBytes, len(encoded))
		}
	}
}

func BenchmarkVarint64Encode(b *testing.B) {
	values := []int64{0, 127, 128, 16383, 16384, 2097151, 2097152}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			EncodeVarint64(v)
		}
	}
}

func BenchmarkVarint64Decode(b *testing.B) {
	// Pre-encode some values
	encoded := [][]byte{
		EncodeVarint64(0),
		EncodeVarint64(127),
		EncodeVarint64(128),
		EncodeVarint64(16383),
		EncodeVarint64(16384),
		EncodeVarint64(2097151),
		EncodeVarint64(2097152),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, buf := range encoded {
			// nolint : errcheck // Ignore errors for benchmarking
			DecodeVarint64(buf)
		}
	}
}

func BenchmarkUvarint64Encode(b *testing.B) {
	values := []uint64{0, 127, 128, 16383, 16384, 2097151, 2097152}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, v := range values {
			// nolint : errcheck // Ignore errors for benchmarking
			EncodeUvarint64(v)
		}
	}
}

func BenchmarkUvarint64Decode(b *testing.B) {
	// Pre-encode some values
	encoded := [][]byte{
		EncodeUvarint64(0),
		EncodeUvarint64(127),
		EncodeUvarint64(128),
		EncodeUvarint64(16383),
		EncodeUvarint64(16384),
		EncodeUvarint64(2097151),
		EncodeUvarint64(2097152),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, buf := range encoded {
			// nolint : errcheck // Ignore errors for benchmarking
			DecodeUvarint64(buf)
		}
	}
}

func testEncodeUvarint64(v uint64) {
	buf := EncodeUvarint64(v)
	fmt.Printf("original:%+v,", v)
	fmt.Printf("encode:%+v,", buf)
	decode, n, err := DecodeUvarint64(buf)
	fmt.Printf("decode:%+v,", decode)
	fmt.Printf("size:%+v,", n)
	if err != nil {
		fmt.Printf("error:%+v", err)
	}
	fmt.Printf("\n")
}

func testMustInvalidUvarint64(t *testing.T, invalid bool, bs []byte) {
	t.Helper()
	t.Logf("testMustInvalidUvarint64: invalid=%v, bs=%v", invalid, bs)
	v, n, err := DecodeUvarint64(bs)
	if invalid {
		if err == nil {
			t.Error("Expected error for invalid uvarint encoding")
		}
		t.Logf("error: %v", err)
	} else {
		if err != nil {
			t.Errorf("Did not expect error for valid uvarint encoding: %v", err)
		} else {
			t.Logf("value: %+v, size: %+v", v, n)
		}
	}
}
