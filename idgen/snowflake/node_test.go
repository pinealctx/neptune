package snowflake

import (
	"testing"
	"time"
)

func TestNewHardNode(t *testing.T) {
	var (
		id1, id2 int64
		n1, n2   Node
		err      error
	)
	n1, err = NewMonoNode(0)
	if err != nil {
		panic(err)
	}
	n2, err = NewNode(0, 0)

	id1, id2 = n1.Generate(), n2.Generate()
	t.Log(id1)
	t.Log(id2)
	t.Log(IDParse(id1))
	t.Log(IDParse(id2))
}

func TestNewHardNodeHook1(t *testing.T) {
	var now = time.Now()
	_HookNow = func() time.Time {
		return now
	}
	var node, _ = NewNode(0, 0)
	var size = 4096 * 3
	var ids = make([]int64, size)
	for i := 0; i < size; i++ {
		ids[i] = node.Generate()
	}
	t.Log(IDFields(ids[0]))
	t.Log(IDFields(ids[4095]))
	t.Log(IDFields(ids[4096]))
	t.Log(IDFields(ids[4096*2-1]))
	t.Log(IDFields(ids[4096*2]))
	t.Log(IDFields(ids[4096*3-1]))
}

func TestNewHardNodeHook2(t *testing.T) {
	var now = time.Now()
	_HookNow = func() time.Time {
		return now
	}
	var node, _ = NewNode(0, 0)
	var id = node.Generate()
	t.Log(id)
	t.Log(IDFields(id))

	_HookNow = func() time.Time {
		return now.Add(-time.Second)
	}
	node, _ = NewNode(0, id)

	var size = 4096 * 3
	var ids = make([]int64, size)
	for i := 0; i < size; i++ {
		ids[i] = node.Generate()
	}

	t.Log(IDFields(ids[0]))
	t.Log(IDFields(ids[1]))
	t.Log(IDFields(ids[4094]))
	t.Log(IDFields(ids[4095]))
	t.Log(IDFields(ids[4096]))
	t.Log(IDFields(ids[4096*2-2]))
	t.Log(IDFields(ids[4096*2-1]))
	t.Log(IDFields(ids[4096*2]))
	t.Log(IDFields(ids[4096*3-2]))
	t.Log(IDFields(ids[4096*3-1]))
}

func TestNewHardNodeHook3(t *testing.T) {
	var now = time.Now()
	_HookNow = func() time.Time {
		return now
	}
	var node, _ = NewNode(0, 0)
	var id = node.Generate()
	t.Log(id)
	t.Log(IDFields(id))
	id = node.Generate()
	t.Log(id)
	t.Log(IDFields(id))

	_HookNow = func() time.Time {
		return now.Add(-time.Second)
	}
	node, _ = NewNode(0, id)

	var size = 4096 * 3
	var ids = make([]int64, size)
	for i := 0; i < size; i++ {
		ids[i] = node.Generate()
	}

	t.Log(IDFields(ids[0]))
	t.Log(IDFields(ids[1]))
	t.Log(IDFields(ids[4094]))
	t.Log(IDFields(ids[4095]))
	t.Log(IDFields(ids[4096]))
	t.Log(IDFields(ids[4096*2-2]))
	t.Log(IDFields(ids[4096*2-1]))
	t.Log(IDFields(ids[4096*2]))
	t.Log(IDFields(ids[4096*3-2]))
	t.Log(IDFields(ids[4096*3-1]))
}
