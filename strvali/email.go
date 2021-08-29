package strvali

import (
	"regexp"
)

var (
	emailReg = regexp.MustCompile(`^\w[-\w.+]*@([A-Za-z0-9][-_A-Za-z0-9]*\.)+[A-Za-z]{2,14}$`)
)

//IsValidEmail 验证是否为合法的邮箱
func IsValidEmail(email string) bool {
	return emailReg.MatchString(email)
}
