package stringx

import (
	"strings"
)

// Reverse : reverse a string
// "" => ""
// "a" => "a"
// "ab" => "ba"
// "abc" => "cba"
// "abcd" => "dcba"
// "aba" => "aba"
func Reverse(s string) string {
	var rs = []rune(s)
	var size = len(rs)
	for i, j := 0, size-1; i < j; i, j = i+1, j-1 {
		rs[i], rs[j] = rs[j], rs[i]
	}
	return string(rs)
}

// Concat : concat
func Concat(ss ...string) string {
	return ConcatSlice(ss)
}

// ConcatSlice : concat slice
func ConcatSlice(ss []string) string {
	var builder = &strings.Builder{}
	builder.Grow(countStrings(ss))
	for i := range ss {
		_, _ = builder.WriteString(ss[i])
	}
	return builder.String()
}

// Append : append string
func Append(h string, ss ...string) string {
	return AppendSlice(h, ss)
}

// AppendSlice : append string slice to head string
func AppendSlice(h string, ss []string) string {
	var builder = &strings.Builder{}
	builder.Grow(len(h) + countStrings(ss))
	_, _ = builder.WriteString(h)
	for i := range ss {
		_, _ = builder.WriteString(ss[i])
	}
	return builder.String()
}

// ContainAny : contain any sub string
func ContainAny(s string, sbs ...string) bool {
	return ContainsAnySubs(s, sbs)
}

// ContainsAnySubs : contain any sub string
func ContainsAnySubs(s string, sbs []string) bool {
	for _, sb := range sbs {
		if strings.Contains(s, sb) {
			return true
		}
	}
	return false
}

// TrimStringsSpace : trim a string slice space
func TrimStringsSpace(ss []string) []string {
	var size = len(ss)
	if size == 0 {
		return nil
	}
	var ns = make([]string, size)
	for i := 0; i < size; i++ {
		ns[i] = strings.TrimSpace(ss[i])
	}
	return ns
}

func countStrings(ss []string) int {
	var number = 0
	for i := range ss {
		number += len(ss[i])
	}
	return number
}
