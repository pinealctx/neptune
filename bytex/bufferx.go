package bytex

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
)

const (
	defaultByteBuff = 1024
)

var (
	ErrByteBufferEmpty = errors.New("byte.buffer.empty")
	ErrReadWrongNum    = errors.New("byte.buffer.wrong.num")
	ErrSizeLimit       = errors.New("byte.buffer.size.limit")
)

//IBufferX buffer interface
type IBufferX interface {
	//Len : length
	Len() int
	//Read specific p, if length is not enough, return error
	Read(p []byte) error
	//ReadN n bytes, if length is not enough, return error
	ReadN(n int) ([]byte, error)
	//Write to buffer
	Write(p []byte)
	//Bytes left bytes
	Bytes() []byte
	//Reset reset cursor
	Reset()

	//ReWrite a buffer
	ReWrite(pos int, p []byte)
	//ReWriteU32 rewrite with a specific pos
	ReWriteU32(pos int, v uint32)

	//ReadU8 read byte
	ReadU8() (byte, error)
	//WriteU8 write byte
	WriteU8(byte)

	//ReadBool read bool
	ReadBool() (bool, error)
	//WriteBool write bool
	WriteBool(bool)

	//ReadLimitString read string
	ReadLimitString(limit uint32) (string, error)
	//WriteLimitString write string
	WriteLimitString(limit uint32, val string) error

	//ReadString read string
	ReadString() (string, error)
	//WriteString write string
	WriteString(val string)

	//ReadU16 read uint16
	ReadU16() (uint16, error)
	//WriteU16 write uint16
	WriteU16(uint16)

	//ReadI16 read int16
	ReadI16() (int16, error)
	//WriteI16 write int16
	WriteI16(int16)

	//ReadU32 read uint32
	ReadU32() (uint32, error)
	//WriteU32 write uint32
	WriteU32(uint32)

	//ReadI32 read int32
	ReadI32() (int32, error)
	//WriteI32 write int32
	WriteI32(int32)

	//ReadU64 read uint64
	ReadU64() (uint64, error)
	//WriteU64 write uint64
	WriteU64(uint64)

	//ReadI64 read int64
	ReadI64() (int64, error)
	//WriteI64 write int64
	WriteI64(int64)

	//ReadF64 read float64
	ReadF64() (float64, error)
	//WriteF64 write float64
	WriteF64(float64)
}

//BufferX buffer implement
type BufferX struct {
	buffer *bytes.Buffer
}

//Len : buffer len
func (b *BufferX) Len() int {
	return b.buffer.Len()
}

//read to a buffer
func (b *BufferX) Read(p []byte) error {
	var l = len(p)
	if l == 0 {
		return nil
	}
	var size, err = b.buffer.Read(p)
	if err != nil {
		return err
	}
	if size != l {
		return ErrByteBufferEmpty
	}
	return nil
}

//ReadN read n length buffer
func (b *BufferX) ReadN(n int) ([]byte, error) {
	if n <= 0 {
		return nil, ErrReadWrongNum
	}
	var buf = make([]byte, n)
	var err = b.Read(buf)
	return buf, err
}

//ZReadN read n length buffer - no copy
func (b *BufferX) ZReadN(n int) ([]byte, error) {
	if n < 0 {
		return nil, ErrReadWrongNum
	}
	var data = b.buffer.Next(n)
	if len(data) != n {
		return nil, ErrByteBufferEmpty
	}
	return data, nil
}

//write to a buffer
func (b *BufferX) Write(p []byte) {
	_, _ = b.buffer.Write(p)
}

//ReWrite buffer
func (b *BufferX) ReWrite(pos int, p []byte) {
	var buf = b.buffer.Bytes()
	copy(buf[pos:], p)
}

//ReadLimitString read limit size string
func (b *BufferX) ReadLimitString(limit uint32) (string, error) {
	var n, err = b.ReadU32()
	if err != nil {
		return "", err
	}

	if n > limit {
		return "", ErrSizeLimit
	}

	var size = int(n)
	var data = b.buffer.Next(size)
	if len(data) != size {
		return "", ErrByteBufferEmpty
	}
	return string(data), nil
}

//WriteLimitString write limit size string
func (b *BufferX) WriteLimitString(limit uint32, val string) error {
	var size = uint32(len(val))
	if size > limit {
		return ErrSizeLimit
	}
	b.WriteU32(size)
	_, _ = b.buffer.WriteString(val)
	return nil
}

//ReadString read string
func (b *BufferX) ReadString() (string, error) {
	var n, err = b.ReadU32()
	if err != nil {
		return "", err
	}
	var size = int(n)
	var data = b.buffer.Next(size)
	if len(data) != size {
		return "", ErrByteBufferEmpty
	}
	return string(data), nil
}

//WriteString write string
func (b *BufferX) WriteString(val string) {
	var size = uint32(len(val))
	b.WriteU32(size)
	_, _ = b.buffer.WriteString(val)
}

//Bytes all buffer bytes
func (b *BufferX) Bytes() []byte {
	return b.buffer.Bytes()
}

//Reset clear all buffer bytes
func (b *BufferX) Reset() {
	b.buffer.Reset()
}

//ReadBool read a bool
func (b *BufferX) ReadBool() (bool, error) {
	var x, err = b.ReadU8()
	if err != nil {
		return false, err
	}
	return x != 0, nil
}

//WriteBool write a bool
func (b *BufferX) WriteBool(v bool) {
	if v {
		_ = b.buffer.WriteByte(1)
	} else {
		_ = b.buffer.WriteByte(0)
	}
}

//ReadU8 read a byte
func (b *BufferX) ReadU8() (byte, error) {
	return b.buffer.ReadByte()
}

//WriteU8 write a byte
func (b *BufferX) WriteU8(v byte) {
	_ = b.buffer.WriteByte(v)
}

//ReadU16 read uint16
func (b *BufferX) ReadU16() (uint16, error) {
	var u16buf [2]byte
	var err = b.Read(u16buf[:])
	if err != nil {
		return 0, err
	}
	var u16 = binary.LittleEndian.Uint16(u16buf[:])
	return u16, err
}

//WriteU16 write uint16
func (b *BufferX) WriteU16(v uint16) {
	var u16buf [2]byte
	binary.LittleEndian.PutUint16(u16buf[:], v)
	b.Write(u16buf[:])
}

//ReadI16 read int16
func (b *BufferX) ReadI16() (int16, error) {
	var u16, err = b.ReadU16()
	return int16(u16), err
}

//WriteI16 write int16
func (b *BufferX) WriteI16(v int16) {
	b.WriteU16(uint16(v))
}

//ReadU32 read uint32
func (b *BufferX) ReadU32() (uint32, error) {
	var u32buf [4]byte
	var err = b.Read(u32buf[:])
	if err != nil {
		return 0, err
	}
	var u32 = binary.LittleEndian.Uint32(u32buf[:])
	return u32, err
}

//WriteU32 write uint32
func (b *BufferX) WriteU32(v uint32) {
	var u32buf [4]byte
	binary.LittleEndian.PutUint32(u32buf[:], v)
	b.Write(u32buf[:])
}

//ReWriteU32 rewrite uint32
func (b *BufferX) ReWriteU32(pos int, v uint32) {
	var u32buf [4]byte
	binary.LittleEndian.PutUint32(u32buf[:], v)
	b.ReWrite(pos, u32buf[:])
}

//ReadI32 read int32
func (b *BufferX) ReadI32() (int32, error) {
	var u32, err = b.ReadU32()
	return int32(u32), err
}

//WriteI32 write int32
func (b *BufferX) WriteI32(v int32) {
	b.WriteU32(uint32(v))
}

//ReadU64 read uint64
func (b *BufferX) ReadU64() (uint64, error) {
	var u64buf [8]byte
	var err = b.Read(u64buf[:])
	if err != nil {
		return 0, err
	}
	var u64 = binary.LittleEndian.Uint64(u64buf[:])
	return u64, err
}

//WriteU64 write uint64
func (b *BufferX) WriteU64(v uint64) {
	var u64buf [8]byte
	binary.LittleEndian.PutUint64(u64buf[:], v)
	b.Write(u64buf[:])
}

//ReadI64 read int64
func (b *BufferX) ReadI64() (int64, error) {
	var u64, err = b.ReadU64()
	return int64(u64), err
}

//WriteI64 write int64
func (b *BufferX) WriteI64(v int64) {
	b.WriteU64(uint64(v))
}

//ReadVarU64 read variant uint64
func (b *BufferX) ReadVarU64() (uint64, error) {
	return binary.ReadUvarint(b.buffer)
}

//WriteVarU64 write variant uint64
func (b *BufferX) WriteVarU64(v uint64) {
	var u64buf [12]byte
	var n = binary.PutUvarint(u64buf[:], v)
	b.Write(u64buf[:n])
}

//ReadVarI64 read variant int64
func (b *BufferX) ReadVarI64() (int64, error) {
	return binary.ReadVarint(b.buffer)
}

//WriteVarI64 write variant int64
func (b *BufferX) WriteVarI64(v int64) {
	var i64buf [12]byte
	var n = binary.PutVarint(i64buf[:], v)
	b.Write(i64buf[:n])
}

//ReadVarU32 read variant uint32
func (b *BufferX) ReadVarU32() (uint32, error) {
	var v, err = binary.ReadUvarint(b.buffer)
	return uint32(v), err
}

//WriteVarU32 write variant uint32
func (b *BufferX) WriteVarU32(v uint32) {
	var u64buf [12]byte
	var n = binary.PutUvarint(u64buf[:], uint64(v))
	b.Write(u64buf[:n])
}

//ReadVarI32 read variant int32
func (b *BufferX) ReadVarI32() (int32, error) {
	var v, err = binary.ReadVarint(b.buffer)
	return int32(v), err
}

//WriteVarI32 write variant int32
func (b *BufferX) WriteVarI32(v int32) {
	var i64buf [12]byte
	var n = binary.PutVarint(i64buf[:], int64(v))
	b.Write(i64buf[:n])
}

//ReadF64 read float64
func (b *BufferX) ReadF64() (float64, error) {
	var u64, err = b.ReadU64()
	return math.Float64frombits(u64), err
}

//WriteF64 write float64
func (b *BufferX) WriteF64(v float64) {
	b.WriteU64(math.Float64bits(v))
}

//NewReadableBufferX new buffer from existed bytes to read.
//Use existed bytes to fill buffer, the buffer is always used as read stream.
func NewReadableBufferX(data []byte) *BufferX {
	var buffer = bytes.NewBuffer(data)
	var bufferX = &BufferX{buffer: buffer}
	return bufferX
}

//NewBufferX new buffer with a default size
//as unpack
func NewBufferX() *BufferX {
	var data = make([]byte, defaultByteBuff)
	var buffer = bytes.NewBuffer(data)
	buffer.Reset()
	var bufferX = &BufferX{buffer: buffer}
	return bufferX
}

//NewSizedBufferX : new buffer with specific size
func NewSizedBufferX(size int) *BufferX {
	var data = make([]byte, size)
	var buffer = bytes.NewBuffer(data)
	buffer.Reset()
	var bufferX = &BufferX{buffer: buffer}
	return bufferX
}
