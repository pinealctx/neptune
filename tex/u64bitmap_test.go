package tex

import (
	"testing"
)

func TestU64BitMap_Set(t *testing.T) {
	var x U64BitMap

	t.Logf("%064b", x)
	t.Log("len:", x.Len())
	t.Log("values:", x.Values())
	t.Log("left:", x.Left())
	t.Log("full:", x.Full())
	logXFromBytes(t, x.Bytes())

	for i := byte(0); i < 64; i++ {
		x.Set(i)
		t.Logf("%064b", x)
		t.Log("len:", x.Len())
		t.Log("values:", x.Values())
		t.Log("left:", x.Left())
		t.Log("full:", x.Full())
		logXFromBytes(t, x.Bytes())
	}
}

func TestNilAppend(t *testing.T) {
	var x []byte
	var y []byte

	var z = append(x, y...)
	t.Log(x == nil)
	t.Log(y == nil)
	t.Log(z == nil)
}

func logXFromBytes(t *testing.T, buf []byte) {
	t.Log("bytes:", buf)
	var ub U64BitMap
	ub.FromBytes(buf)
	t.Logf("clone %064b", ub)
	t.Log("")
	t.Log("")
}
