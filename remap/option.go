package remap

const (
	//DefaultPrime 缺省使用73 -- 致敬谢耳朵最喜欢的素数
	DefaultPrime uint64 = 73
)

//remap option: use a prime number as group number
type _Option struct {
	prime uint64
}

//Option : option function
type Option func(o *_Option)

//WithPrime : setup prime number
func WithPrime(prime uint64) Option {
	return func(o *_Option) {
		o.prime = prime
	}
}
