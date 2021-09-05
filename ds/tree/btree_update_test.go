package tree

import "testing"

func TestBTree_Update(t *testing.T) {
	var bt = NewBTree()
	for i := 1; i <= 5; i++ {
		var x = &xCmp{
			x: i,
			y: i,
		}
		bt.Insert(x)
	}

	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.Update(&xCmp{x:1}, &xCmp{x:-1, y:-1})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.Update(&xCmp{x:2}, &xCmp{x:-2, y:-2})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.Update(&xCmp{x:3}, &xCmp{x:-3, y:-3})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.Update(&xCmp{x:4}, &xCmp{x:-4, y:-4})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.Update(&xCmp{x:5}, &xCmp{x:-5, y:-5})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.Update(&xCmp{x:6}, &xCmp{x:-6, y:-6})
	testAscendGt(t, bt, nil, 6)
	t.Log("")
}

func TestBTree_UpdateOrInsert(t *testing.T) {
	var bt = NewBTree()
	for i := 1; i <= 5; i++ {
		var x = &xCmp{
			x: i,
			y: i,
		}
		bt.Insert(x)
	}

	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.UpdateOrInsert(&xCmp{x:1}, &xCmp{x:-1, y:-1})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.UpdateOrInsert(&xCmp{x:2}, &xCmp{x:-2, y:-2})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.UpdateOrInsert(&xCmp{x:3}, &xCmp{x:-3, y:-3})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.UpdateOrInsert(&xCmp{x:4}, &xCmp{x:-4, y:-4})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.UpdateOrInsert(&xCmp{x:5}, &xCmp{x:-5, y:-5})
	testAscendGt(t, bt, nil, 6)
	t.Log("")

	bt.UpdateOrInsert(&xCmp{x:6}, &xCmp{x:-6, y:-6})
	testAscendGt(t, bt, nil, 6)
	t.Log("")
}
