package strvali

import (
	"regexp"
	"strings"
)

var (
	// emailReg is the regular expression for email validation.
	// 1. Local part: Disallow consecutive dots, must start/end with word char.
	// 2. Domain part: Allow underscores (kept as requested), but must end with alphanumeric char.
	// 3. TLD: Length up to 63 (standard DNS label limit).
	emailReg = regexp.MustCompile(`^[\w]+([-+.]\w+)*@([A-Za-z0-9]([-_A-Za-z0-9]*[A-Za-z0-9])?\.)+[A-Za-z]{2,63}$`)
)

// IsValidEmail checks if the provided email string is a valid email format
func IsValidEmail(email string) bool {
	if len(email) > 254 {
		return false
	}
	return emailReg.MatchString(email)
}

// NormalizeEmail converts email to lowercase
// Email addresses are case-insensitive according to RFC 5321
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
