package ulog

import (
	"go.uber.org/zap"
)

// Zaps a extendable zap fields
type Zaps []zap.Field

func NewBuffedZaps(n int) Zaps {
	return make(Zaps, 0, n)
}

// NewDefaultZaps default constructor
func NewDefaultZaps() Zaps {
	return NewBuffedZaps(4)
}

// NewMidZaps middle size constructor
func NewMidZaps() Zaps {
	return NewBuffedZaps(8)
}

// NewHugeZaps huge size constructor
func NewHugeZaps() Zaps {
	return NewBuffedZaps(32)
}

func (z *Zaps) Append(e zap.Field) *Zaps {
	*z = append(*z, e)
	return z
}

func (z *Zaps) Appends(es ...zap.Field) *Zaps {
	*z = append(*z, es...)
	return z
}
