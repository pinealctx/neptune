package pipe

const (
	//DefaultSlotSize slot size
	//素数
	DefaultSlotSize = 509
	//DefaultQSize default pipe size, total current request can be pipe is 509*1024*8 = 4169728
	DefaultQSize = 1024 * 8
)

// shunt option
type _Option struct {
	//slot size
	slotSize int
	//queue size in each slot
	qSize int
}

// Option shunt option function
type Option func(o *_Option)

// WithSlotSize setup slot size
func WithSlotSize(slotSize int) Option {
	return func(o *_Option) {
		o.slotSize = slotSize
	}
}

// WithQSize setup queue size in each slot
func WithQSize(qSize int) Option {
	return func(o *_Option) {
		o.qSize = qSize
	}
}

// GetOption : return slot size and q size
func GetOption(opts ...Option) (slotSize int, qSize int) {
	var o = &_Option{
		slotSize: DefaultSlotSize,
		qSize:    DefaultQSize,
	}
	for _, opt := range opts {
		opt(o)
	}
	return o.slotSize, o.qSize
}
