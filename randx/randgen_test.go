package randx

import (
	"testing"
)

func TestRandGen_Gen(t *testing.T) {
	g := NewRandGen([]byte("0123456789"), 6)
	t.Log(g.Gen())
	t.Log(g.Gen())
	t.Log(g.Gen())
	t.Log(g.Gen())
	t.Log(g.Gen())
	t.Log(g.Gen())
	t.Log(g.Gen())
}
