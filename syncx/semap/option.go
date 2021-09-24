package semap

const (
	//DefaultSize default map size
	DefaultSize = 1024 * 16
	//DefaultRWRatio default read/write ratio
	DefaultRWRatio = 10
)

//remap option: use a prime number as group number
type _Option struct {
	prime   uint64
	size    int
	rwRatio int
}

//Option : option function
type Option func(o *_Option)

//WithPrime : setup prime number
func WithPrime(prime uint64) Option {
	return func(o *_Option) {
		o.prime = prime
	}
}

//WithSize : setup map size
func WithSize(size int) Option {
	return func(o *_Option) {
		o.size = size
	}
}

//WithRwRatio : setup read write ratio
func WithRwRatio(rwRatio int) Option {
	return func(o *_Option) {
		o.rwRatio = rwRatio
	}
}

//RangeOption : range option
func RangeOption(opts ...Option) *_Option {
	var o = &_Option{
		size:    DefaultSize,
		rwRatio: DefaultRWRatio,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
