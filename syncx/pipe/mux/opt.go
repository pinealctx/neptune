package mux

const (
	//DefaultMuxSize slot size
	//素数
	DefaultMuxSize = 127
	//DefaultDeepSize default queue deep
	DefaultDeepSize = 1024 * 8
)

// mux option
type _Option struct {
	//mux size
	muxSize int
	//queue deep size
	deepSize int
}

// Option mux option function
type Option func(o *_Option)

// WithSize setup mux size
func WithSize(muxSize int) Option {
	return func(o *_Option) {
		o.muxSize = muxSize
	}
}

// WithDeep setup queue deep
func WithDeep(deepSize int) Option {
	return func(o *_Option) {
		o.deepSize = deepSize
	}
}
