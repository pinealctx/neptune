package genericx

// Number is a constraint that permits any numeric type.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// SignedNumber is a constraint that permits any signed numeric type.
type SignedNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// Integer is a constraint that permits any integer type.
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// SignedInteger is a constraint that permits any signed integer type.
type SignedInteger interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// UnsignedInteger is a constraint that permits any unsigned integer type.
type UnsignedInteger interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

// Min returns the smaller of two values.
func Min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}

// Max returns the larger of two values.
func Max[T Number](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// Abs returns the absolute value of a signed number.
func Abs[T SignedNumber](a T) T {
	if a < 0 {
		return -a
	}
	return a
}
