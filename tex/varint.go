package tex

import (
	"errors"

	"google.golang.org/protobuf/encoding/protowire"
)

var (
	ErrEmptyBuffer        = errors.New("empty.buffer")
	ErrInvalidVarint      = errors.New("invalid.varint.encoding")
	ErrNotRegularEncoding = errors.New("not.regular.varint.encoding")
)

// AppendVarint64 appends a varint-encoded int64 to the given byte slice
func AppendVarint64(buf []byte, value int64) []byte {
	// nolint:gosec // intentional int64->uint64 conversion for varint encoding
	return protowire.AppendVarint(buf, uint64(value))
}

// EncodeVarint64 encodes an int64 value using variable-length encoding (same as protobuf)
// Returns the encoded bytes
func EncodeVarint64(value int64) []byte {
	// nolint:gosec // intentional int64->uint64 conversion for varint encoding
	uValue := uint64(value)
	buf := make([]byte, protowire.SizeVarint(uValue))
	protowire.AppendVarint(buf[:0], uValue)
	return buf
}

// DecodeVarint64 decodes a variable-length encoded int64 from bytes
// Returns the decoded value, the number of bytes consumed, and any error
func DecodeVarint64(buf []byte) (value int64, bytesRead int, err error) {
	u, n, e := DecodeUvarint64(buf)
	if e != nil {
		return 0, n, e
	}
	// nolint:gosec // intentional uint64->int64 conversion for varint decoding
	return int64(u), n, nil
}

// AppendUvarint64 appends a varint-encoded uint64 to the given byte slice
func AppendUvarint64(buf []byte, value uint64) []byte {
	return protowire.AppendVarint(buf, value)
}

// EncodeUvarint64 encodes a uint64 value using variable-length encoding (same as protobuf)
// Returns the encoded bytes
func EncodeUvarint64(value uint64) []byte {
	buf := make([]byte, protowire.SizeVarint(value))
	protowire.AppendVarint(buf[:0], value)
	return buf
}

// DecodeUvarint64 decodes a variable-length encoded uint64 from bytes
// Returns the decoded value, the number of bytes consumed, and any error
func DecodeUvarint64(buf []byte) (value uint64, bytesRead int, err error) {
	if len(buf) == 0 {
		return 0, 0, ErrEmptyBuffer
	}

	value, bytesRead = protowire.ConsumeVarint(buf)
	if bytesRead <= 0 {
		return 0, 0, ErrInvalidVarint
	}

	if bytesRead > protowire.SizeVarint(value) {
		// This indicates that the varint was not encoded in its shortest form
		// which is not allowed in standard varint encoding
		return value, bytesRead, ErrNotRegularEncoding
	}
	return value, bytesRead, nil
}

// VarintSize returns the number of bytes required to encode the given int64 value
func VarintSize(value int64) int {
	// nolint:gosec // intentional int64->uint64 conversion for varint size calculation
	return protowire.SizeVarint(uint64(value))
}

// UvarintSize returns the number of bytes required to encode the given uint64 value
func UvarintSize(value uint64) int {
	return protowire.SizeVarint(value)
}
