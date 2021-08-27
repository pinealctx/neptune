package tex

import (
	"testing"
)

func TestU128BitMap_Set(t *testing.T) {
	var x U128BitMap

	t.Logf("%064b", x)
	t.Log("len:", x.Len())
	t.Log("values:", x.Values())
	t.Log("left:", x.Left())
	t.Log("full:", x.Full())
	t.Log("bytes:", x.Bytes())
	logxyFromBytes(t, x.Bytes())

	for i := byte(0); i < 128; i++ {
		x.Set(i)
		t.Logf("%064b", x)
		t.Log("len:", x.Len())
		t.Log("values:", x.Values())
		t.Log("left:", x.Left())
		t.Log("full:", x.Full())
		t.Log("bytes:", x.Bytes())
		logxyFromBytes(t, x.Bytes())
	}
}

func logxyFromBytes(t *testing.T, buf []byte) {
	t.Log("bytes:", buf)
	var ub U128BitMap
	ub.FromBytes(buf)
	t.Logf("clone %064b", ub)
	t.Log("")
	t.Log("")
}
