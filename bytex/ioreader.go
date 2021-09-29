package bytex

import (
	"encoding/binary"
	"io"
	"math"
)

//IReaderX reader interface
type IReaderX interface {
	//Read : read specific p, if length is not enough, return error
	Read(p []byte) error
	//ReadN : read n bytes, if length is not enough, return error
	ReadN(n int) ([]byte, error)

	ZReadN(n int) ([]byte, error)

	//ReadByte : read byte
	ReadByte() (byte, error)
	//ReadBool : read bool
	ReadBool() (bool, error)

	//ReadLimitString : read string
	ReadLimitString(limit uint32) (string, error)

	//ReadString : read string
	ReadString() (string, error)

	//ReadU16 : read uint16
	ReadU16() (uint16, error)

	//ReadI16 : read int16
	ReadI16() (int16, error)

	//ReadU32 : read uint32
	ReadU32() (uint32, error)

	//ReadI32 : read int32
	ReadI32() (int32, error)

	//ReadU64 : read uint64
	ReadU64() (uint64, error)

	//ReadI64 : read int64
	ReadI64() (int64, error)

	//ReadF64 : read float64
	ReadF64() (float64, error)
}

//ReaderX buffer implement
type ReaderX struct {
	reader io.Reader
}

//Read : read to a buffer
func (b *ReaderX) Read(p []byte) error {
	var l = len(p)
	if l == 0 {
		return nil
	}
	var size, err = b.reader.Read(p)
	if err != nil {
		return err
	}
	if size != l {
		return ErrByteBufferEmpty
	}
	return nil
}

//ReadN read n length buffer
func (b *ReaderX) ReadN(n int) ([]byte, error) {
	if n <= 0 {
		return nil, ErrReadWrongNum
	}
	var buf = make([]byte, n)
	var err = b.Read(buf)
	return buf, err
}

//ZReadN read n length buffer - no copy
func (b *ReaderX) ZReadN(n int) ([]byte, error) {
	return b.ReadN(n)
}

//ReadLimitString read limit size string
func (b *ReaderX) ReadLimitString(limit uint32) (string, error) {
	var n, err = b.ReadU32()
	if err != nil {
		return "", err
	}

	if n > limit {
		return "", ErrSizeLimit
	}

	var size = int(n)
	var data []byte
	data, err = b.ZReadN(size)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

//ReadString read string
func (b *ReaderX) ReadString() (string, error) {
	var n, err = b.ReadU32()
	if err != nil {
		return "", err
	}
	var size = int(n)
	var data []byte
	data, err = b.ZReadN(size)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

//ReadBool read a bool
func (b *ReaderX) ReadBool() (bool, error) {
	var x, err = b.ReadByte()
	if err != nil {
		return false, err
	}
	return x != 0, nil
}

//ReadByte read a byte
func (b *ReaderX) ReadByte() (byte, error) {
	var p [1]byte
	var err = b.Read(p[:])
	if err != nil {
		return 0, err
	}
	return p[0], nil
}

//ReadU16 read uint16
func (b *ReaderX) ReadU16() (uint16, error) {
	var u16buf [2]byte
	var err = b.Read(u16buf[:])
	if err != nil {
		return 0, err
	}
	var u16 = binary.LittleEndian.Uint16(u16buf[:])
	return u16, err
}

//ReadI16 read int16
func (b *ReaderX) ReadI16() (int16, error) {
	var u16, err = b.ReadU16()
	return int16(u16), err
}

//ReadU32 read uint32
func (b *ReaderX) ReadU32() (uint32, error) {
	var u32buf [4]byte
	var err = b.Read(u32buf[:])
	if err != nil {
		return 0, err
	}
	var u32 = binary.LittleEndian.Uint32(u32buf[:])
	return u32, err
}

//ReadI32 read int32
func (b *ReaderX) ReadI32() (int32, error) {
	var u32, err = b.ReadU32()
	return int32(u32), err
}

//ReadU64 read uint64
func (b *ReaderX) ReadU64() (uint64, error) {
	var u64buf [8]byte
	var err = b.Read(u64buf[:])
	if err != nil {
		return 0, err
	}
	var u64 = binary.LittleEndian.Uint64(u64buf[:])
	return u64, err
}

//ReadI64 read int64
func (b *ReaderX) ReadI64() (int64, error) {
	var u64, err = b.ReadU64()
	return int64(u64), err
}

//ReadF64 read float64
func (b *ReaderX) ReadF64() (float64, error) {
	var u64, err = b.ReadU64()
	return math.Float64frombits(u64), err
}

//NewReaderX new with io.reader
func NewReaderX(reader io.Reader) *ReaderX {
	return &ReaderX{reader: reader}
}
