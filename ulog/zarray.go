package ulog

import (
	"time"

	"go.uber.org/zap/zapcore"
)

//Implement basic types array encoding for zap.

type BoolArray []bool

func (bs BoolArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range bs {
		arr.AppendBool(bs[i])
	}
	return nil
}

type ByteStringsArray [][]byte

func (bss ByteStringsArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range bss {
		arr.AppendByteString(bss[i])
	}
	return nil
}

type Complex128Array []complex128

func (nums Complex128Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendComplex128(nums[i])
	}
	return nil
}

type Complex64Array []complex64

func (nums Complex64Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendComplex64(nums[i])
	}
	return nil
}

type DurationArray []time.Duration

func (ds DurationArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ds {
		arr.AppendDuration(ds[i])
	}
	return nil
}

type Float64Array []float64

func (nums Float64Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendFloat64(nums[i])
	}
	return nil
}

type Float32Array []float32

func (nums Float32Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendFloat32(nums[i])
	}
	return nil
}

type IntArray []int

func (nums IntArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt(nums[i])
	}
	return nil
}

type Int64Array []int64

func (nums Int64Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt64(nums[i])
	}
	return nil
}

type Int32Array []int32

func (nums Int32Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt32(nums[i])
	}
	return nil
}

type Int16Array []int16

func (nums Int16Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt16(nums[i])
	}
	return nil
}

type Int8Array []int8

func (nums Int8Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendInt8(nums[i])
	}
	return nil
}

type StringArray []string

func (ss StringArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ss {
		arr.AppendString(ss[i])
	}
	return nil
}

type TimeArray []time.Time

func (ts TimeArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range ts {
		arr.AppendTime(ts[i])
	}
	return nil
}

type UintArray []uint

func (nums UintArray) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint(nums[i])
	}
	return nil
}

type Uint64Array []uint64

func (nums Uint64Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint64(nums[i])
	}
	return nil
}

type Uint32Array []uint32

func (nums Uint32Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint32(nums[i])
	}
	return nil
}

type Uint16Array []uint16

func (nums Uint16Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint16(nums[i])
	}
	return nil
}

type Uint8Array []uint8

func (nums Uint8Array) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	for i := range nums {
		arr.AppendUint8(nums[i])
	}
	return nil
}
