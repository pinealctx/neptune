package semap

const (
	//DefaultRWRatio default read/write ratio
	DefaultRWRatio = 10
)

// remap option: use a prime number as group number
type Option struct {
	prime   uint64
	rwRatio int
}

// OptionFn : option function
type OptionFn func(o *Option)

// WithPrime : setup prime number
func WithPrime(prime uint64) OptionFn {
	return func(o *Option) {
		o.prime = prime
	}
}

// WithRwRatio : setup read write ratio
func WithRwRatio(rwRatio int) OptionFn {
	return func(o *Option) {
		o.rwRatio = rwRatio
	}
}

// RangeOption : range option
func RangeOption(opts ...OptionFn) *Option {
	var o = &Option{
		rwRatio: DefaultRWRatio,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o
}
