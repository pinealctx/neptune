package async

import "sync"

const (
	//DefaultQSize default runner queue size
	DefaultQSize = 1024 * 8
)

// runner config option
type optionT struct {
	//queue size
	size int
	//wait group
	wg *sync.WaitGroup
	//name
	name string
}

// Option : only qSize option
type Option func(option *optionT)

// WithQSize : setup qSize
func WithQSize(qSize int) Option {
	return func(o *optionT) {
		o.size = qSize
	}
}

// WithName : setup name
func WithName(name string) Option {
	return func(o *optionT) {
		o.name = name
	}
}

// WithWaitGroup : setup outside wait group controller
// If this value be set, the wait group will be increase when runner loop go routine exits.
func WithWaitGroup(wg *sync.WaitGroup) Option {
	return func(o *optionT) {
		o.wg = wg
	}
}
