package strvali

import (
	"fmt"
	"github.com/nyaruka/phonenumbers"
	"regexp"
)

var (
	cnPhoneReg = regexp.MustCompile(`^1[3456789]\d{9}$`)
)

// IsValidPhoneNum verify phone(xx) with area(+xx) code.
func IsValidPhoneNum(areaCode, phone string) bool {
	if len(areaCode) <= 1 || len(phone) == 0 {
		return false
	}
	if notNumStr(areaCode[1:]) || notNumStr(phone) {
		return false
	}
	if areaCode == "+86" {
		return verifyCNPhone(phone)
	}
	var phoneNum, err = phonenumbers.Parse(fmt.Sprintf("%s%s", areaCode, phone), "")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(phoneNum)
}

func notNumStr(s string) bool {
	for _, k := range s {
		if k < '0' || k > '9' {
			return true
		}
	}
	return false
}

func verifyCNPhone(k string) bool {
	return cnPhoneReg.MatchString(k)
}
