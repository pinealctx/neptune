package randx

var (
	_lowercaseChars = []byte("abcdefghijklmnopqrstuvwxyz")
	_uppercaseChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	_digitChars     = []byte("0123456789")
)

// PassGen is a password generator
type PassGen struct {
	lowercaseChars []byte
	uppercaseChars []byte
	digitChars     []byte
	specialChars   []byte
	allChars       []byte
}

// NewPassGen creates a new PassGen instance
func NewPassGen(specialChars []byte) *PassGen {
	x := &PassGen{
		lowercaseChars: _lowercaseChars,
		uppercaseChars: _uppercaseChars,
		digitChars:     _digitChars,
		specialChars:   specialChars,
	}
	l1, l2, l3, l4 := len(x.lowercaseChars), len(x.uppercaseChars), len(x.digitChars), len(x.specialChars)
	x.allChars = make([]byte, 0, l1+l2+l3+l4)
	x.allChars = append(x.allChars, x.lowercaseChars...)
	x.allChars = append(x.allChars, x.uppercaseChars...)
	x.allChars = append(x.allChars, x.digitChars...)
	x.allChars = append(x.allChars, x.specialChars...)
	return x
}

// GenerateRandomPassword generates a random password of specified length
// Requirements:
// - length must be at least 4
// - must contain at least 1 lowercase letter
// - must contain at least 1 uppercase letter
// - must contain at least 1 digit
// - must contain at least 1 special character
func (x *PassGen) GenerateRandomPassword(length int) string {
	if length < 4 {
		panic("password length must be at least 4")
	}

	result := make([]byte, length)

	// Step 1: Ensure at least one character from each required category
	result[0] = x.lowercaseChars[IntNSecure(len(x.lowercaseChars))]
	result[1] = x.uppercaseChars[IntNSecure(len(x.uppercaseChars))]
	result[2] = x.digitChars[IntNSecure(len(x.digitChars))]
	result[3] = x.specialChars[IntNSecure(len(x.specialChars))]

	// Step 2: Fill remaining positions with random characters from all categories
	for i := 4; i < length; i++ {
		randomIndex := IntNSecure(len(x.allChars))
		result[i] = x.allChars[randomIndex]
	}

	// Step 3: Shuffle the entire password to avoid predictable patterns
	shufflePassword(result)
	return string(result)
}

// shufflePassword shuffles the password characters using Fisher-Yates algorithm
func shufflePassword(password []byte) {
	n := len(password)
	for i := n - 1; i > 0; i-- {
		j := IntNSecure(i + 1)
		password[i], password[j] = password[j], password[i]
	}
}
