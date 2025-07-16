package tree

import (
	"testing"
)

type xCmp struct {
	x int
	y int
}

func (x *xCmp) Less(v Node) bool {
	return x.x < v.(*xCmp).x
}

func TestBTree_Get(t *testing.T) {
	var bt = NewBTree()
	for i := 1; i <= 10; i++ {
		var x = &xCmp{
			x: i,
			y: i,
		}
		bt.Insert(x)
	}
	for i := 1; i <= 10; i++ {
		var k = &xCmp{x: i}
		var v = bt.Get(k).(*xCmp)
		t.Log("i:", i, "v:", v.x, v.y)
	}
}

func TestBTree_Delete(t *testing.T) {
	var bt = NewBTree()
	for i := 1; i <= 10; i++ {
		var x = &xCmp{
			x: i,
			y: i,
		}
		bt.Insert(x)
	}
	t.Log(bt.t.Len())
	for i := 1; i <= 10; i++ {
		var k = &xCmp{x: i}
		bt.Delete(k)
	}
	t.Log(bt.t.Len())
}

func TestBTree_AscendGtFromNil(t *testing.T) {
	testBTreeAscendGtFrom(t, nil)
}

func TestBTree_AscendGtFrom0(t *testing.T) {
	testBTreeAscendGtNum(t, 0)
}

func TestBTree_AscendGtFrom1(t *testing.T) {
	testBTreeAscendGtNum(t, 1)
}

func TestBTree_AscendGtFrom2(t *testing.T) {
	testBTreeAscendGtNum(t, 2)
}

func TestBTree_AscendGtFrom3(t *testing.T) {
	testBTreeAscendGtNum(t, 3)
}

func TestBTree_AscendGtFrom4(t *testing.T) {
	testBTreeAscendGtNum(t, 4)
}

func TestBTree_AscendGtFrom5(t *testing.T) {
	testBTreeAscendGtNum(t, 5)
}

func TestBTree_AscendGtFrom6(t *testing.T) {
	testBTreeAscendGtNum(t, 6)
}

func TestBTree_AscendGtFrom7(t *testing.T) {
	testBTreeAscendGtNum(t, 7)
}

func TestBTree_AscendGtFrom8(t *testing.T) {
	testBTreeAscendGtNum(t, 8)
}

func TestBTree_AscendGteFromNil(t *testing.T) {
	testBTreeAscendGteFrom(t, nil)
}

func TestBTree_AscendGteFrom0(t *testing.T) {
	testBTreeAscendGteNum(t, 0)
}

func TestBTree_AscendGteFrom1(t *testing.T) {
	testBTreeAscendGteNum(t, 1)
}

func TestBTree_AscendGteFrom2(t *testing.T) {
	testBTreeAscendGteNum(t, 2)
}

func TestBTree_AscendGteFrom3(t *testing.T) {
	testBTreeAscendGteNum(t, 3)
}

func TestBTree_AscendGteFrom4(t *testing.T) {
	testBTreeAscendGteNum(t, 4)
}

func TestBTree_AscendGteFrom5(t *testing.T) {
	testBTreeAscendGteNum(t, 5)
}

func testBTreeAscendGtNum(t *testing.T, k int) {
	t.Helper()
	testBTreeAscendGtFrom(t, &xCmp{x: k})
}

func testBTreeAscendGtFrom(t *testing.T, k Node) {
	t.Helper()
	testBTreeIter(t, testAscendGt, k)
}

func testBTreeAscendGteNum(t *testing.T, k int) {
	t.Helper()
	testBTreeAscendGteFrom(t, &xCmp{x: k})
}

func testBTreeAscendGteFrom(t *testing.T, k Node) {
	t.Helper()
	testBTreeIter(t, testAscendGte, k)
}

func testBTreeIter(t *testing.T, fn _testIterProc, k Node) {
	t.Helper()
	var bt = NewBTree()
	for i := 1; i <= 3; i++ {
		var x = &xCmp{
			x: i,
			y: i,
		}
		bt.Insert(x)
	}

	for i := 1; i <= 5; i++ {
		t.Log("iter:", i)
		fn(t, bt, k, i)
		t.Log("")
	}
}

type _testIterProc func(t *testing.T, b *BTree, k Node, n int)

func testAscendGt(t *testing.T, b *BTree, k Node, n int) {
	t.Helper()
	testIter(t, k, b.AscendGt, n)
}

func testAscendGte(t *testing.T, b *BTree, k Node, n int) {
	t.Helper()
	testIter(t, k, b.AscendGte, n)
}

type _testIterFn func(k Node, filter FilterFn, n int) []Node

func testIter(t *testing.T, k Node, ifn _testIterFn, n int) {
	t.Helper()
	var ffn = func(_ Node) bool {
		return true
	}
	var ns = ifn(k, ffn, n)
	for i := range ns {
		var v = ns[i].(*xCmp)
		t.Log("x:", v.x, "y:", v.y)
	}
}
