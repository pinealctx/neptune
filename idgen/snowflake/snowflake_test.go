package snowflake

import (
	"testing"
	"time"
)

func TestSetup(t *testing.T) {
	var ts = time.Unix(0, _epoch*MsDivNs)
	t.Log(ts)
	var t1 = time.Date(2021, 1, 1, 0, 0, 0, 0, time.Local)
	t.Log(t1)
	t.Log(t1.UnixNano() / MsDivNs)
	Setup(UseEpoch(time.Now()))
	ts = time.Unix(0, _epoch*MsDivNs)
	t.Log(ts)
	var timeShift = _nodeBits + StepBits
	t.Log((1 << timeShift) - 1)
	t.Logf("%x", (1<<timeShift)-1)
	var i = int64(time.Hour * 24 / time.Millisecond)
	t.Log(i)

	var now = time.Now()
	var tp = now.Unix()
	ts = time.Unix(tp, 0)
	t.Log(ts)
	ts = now
	t.Log(ts.Year(), ts.Month(), ts.Day(), ts.Hour(), ts.Minute(), ts.Second(), ts.Nanosecond()/int(time.Millisecond))
}

func TestNode_Generate(t *testing.T) {
	var node, err = NewMonoNode(1)
	if err != nil {
		panic(err)
	}
	var id = node.Generate()
	t.Log(id)
	var cn = CnStyle(id)
	t.Log(cn)
	var nid int64
	nid, err = FromChStyle(cn)
	if err != nil {
		panic(err)
	}
	t.Log(nid)
	if nid != id {
		t.Fail()
	}
	var ts, nv, step = IDParseEx(id)
	t.Log(ts, nv, step)
}

func TestNode_GenerateB(t *testing.T) {
	var node, err = NewMonoNode(0)
	if err != nil {
		panic(err)
	}
	var t1 = time.Now()
	var dc = time.Duration(10000)
	for i := time.Duration(0); i < dc; i++ {
		var id = node.Generate()
		//t.Log(id)
		var cn = CnStyle(id)
		//t.Log(cn)
		var nid int64
		nid, err = FromChStyle(cn)
		if err != nil {
			panic(err)
		}
		//t.Log(nid)
		if nid != id {
			t.Fail()
			return
		}
	}
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/dc)
}

func TestTimeIDRange(t *testing.T) {
	var node, err = NewMonoNode(0)
	if err != nil {
		panic(err)
	}
	var id = node.Generate()
	var now = time.Now()
	t.Log(now)
	var ts = now.UnixNano()
	t.Logf("%064b\n", ts)
	t.Logf("%064b\n", ts/MsDivNs)
	var min, max = TimeIDRange(now)
	t.Logf("%064b\n", id)
	t.Logf("%064b\n", min)
	t.Logf("%064b\n", max)

	var ms, n1, step1 = IDParse(id)
	t.Logf("%064b\n", ms)
	t.Logf("%064b\n", n1)
	t.Logf("%064b\n", step1)

	var t2, n2, s2 = IDParseEx(id)
	t.Log(t2, n2, s2)

	t2, n2, s2 = IDParseEx(min)
	t.Log(t2, n2, s2)
	t2, n2, s2 = IDParseEx(max)
	t.Log(t2, n2, s2)

	var timeShift = _nodeBits + StepBits
	var reMax int64 = (1 << timeShift) - 1
	t.Logf("%064b\n", reMax)
}

func TestTimeBetweenID(t *testing.T) {
	var node, err = NewMonoNode(0)
	if err != nil {
		panic(err)
	}
	var id = node.Generate()
	var now = time.Now()
	t.Log(now)
	t.Logf("%064b\n", id)
	var min, max = TimeBetweenID(now, now.Add(time.Second*10))
	t.Logf("%064b\n", min)
	t.Logf("%064b\n", max)

	var t2, n2, s2 = IDParseEx(min)
	t.Log(t2, n2, s2)
	t2, n2, s2 = IDParseEx(max)
	t.Log(t2, n2, s2)
}
