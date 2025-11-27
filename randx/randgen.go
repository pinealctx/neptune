package randx

// RandGen is a random string generator
type RandGen struct {
	// randChars is the set of characters to use for generating random strings
	randChars []byte
	// length is the length of the generated random strings
	length int
}

// NewRandGen creates a new RandGen instance
func NewRandGen(randChars []byte, length int) *RandGen {
	return &RandGen{
		randChars: randChars,
		length:    length,
	}
}

// Gen : generates a random string
func (rg *RandGen) Gen() string {
	b := make([]byte, rg.length)
	size := len(rg.randChars)
	for i := 0; i < rg.length; i++ {
		b[i] = rg.randChars[IntNSecure(size)]
	}
	return string(b)
}
