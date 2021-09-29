## bytex

go内置的bytes.Buffer功能很丰富，但它更多面向的是bytes这样的数据类型，bytex封装了bytes.Buffer，丰富了更多的基本类型操作。

在参考了MT协议情况下，实现了序列化io接口IBufferX。

从io的角度来说，bytes缓冲区分为read io与write io。
```go
//创建一个可读的缓冲区，data是可读的数据，大部分情况来自网络数据，RPC数据，数据库或文件数据
func NewReadableBufferX(data []byte) *BufferX {
	var buffer = bytes.NewBuffer(data)
	var bufferX = &BufferX{buffer: buffer}
	return bufferX
}

//创建一个可写的缓冲区大小为1K，如果写入内容超过1K，此缓冲区会自动增长
func NewBufferX() *BufferX {
	var data = make([]byte, defaultByteBuff)
	var buffer = bytes.NewBuffer(data)
	buffer.Reset()
	var bufferX = &BufferX{buffer: buffer}
	return bufferX
}

//创建一个可写的缓冲区，大小为传入的size，如果写入内容超过此大小，此缓冲区会自动增长
//从效率上来说，传入的size应该与要写入的数据大小匹配，这样不会浪费内存，也不会因为
//在缓冲区不够时新开内存影响效率。
func NewSizedBufferX(size int) *BufferX {
	var data = make([]byte, size)
	var buffer = bytes.NewBuffer(data)
	buffer.Reset()
	var bufferX = &BufferX{buffer: buffer}
	return bufferX
}
```

```go
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
	//重写缓冲区指定位置的数据，最常用的case是有些数据头部可能需要记录整段数据大小。
	//但最开始写入时，大小并不确定，可以先用0填充头部，在写完整段内容有了大小后在重写头部位置的数据。
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
```