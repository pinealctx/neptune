package strvali

import (
	"regexp"
	"strings"
)

var (
	emailReg = regexp.MustCompile(`^\w[-\w.+]*@([A-Za-z0-9][-_A-Za-z0-9]*\.)+[A-Za-z]{2,14}$`)
)

// IsValidEmail checks if the provided email string is a valid email format
func IsValidEmail(email string) bool {
	return emailReg.MatchString(email)
}

// NormalizeEmail converts email to lowercase
// Email addresses are case-insensitive according to RFC 5321
func NormalizeEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}
