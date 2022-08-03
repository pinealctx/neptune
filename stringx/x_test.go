package stringx

import "testing"

// "" => ""
// "a" => "a"
// "ab" => "ba"
// "abc" => "cba"
// "abcd" => "dcba"
// "aba" => "aba"
func TestReverse(t *testing.T) {
	var s = ""
	t.Log(Reverse(s), s)
	s = "a"
	t.Log(Reverse(s), s)
	s = "ab"
	t.Log(Reverse(s), s)
	s = "abc"
	t.Log(Reverse(s), s)
	s = "abcd"
	t.Log(Reverse(s), s)
	s = "aba"
	t.Log(Reverse(s), s)
	s = "abcde"
	t.Log(Reverse(s), s)
	s = "abcdef"
	t.Log(Reverse(s), s)
}

func TestConcat(t *testing.T) {
	t.Log(Concat())
	t.Log(Concat(""))
	t.Log(Concat("a"))
	t.Log(Concat("a", "", "b"))

	t.Log(Append("", "b"))
	t.Log(Append("a", "b"))
	t.Log(Append("a", "b", "c"))
	t.Log(Append("a", "b", "", "c", "", "d"))
}
