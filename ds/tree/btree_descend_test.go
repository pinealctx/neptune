package tree

import (
	"testing"
)

func TestBTree_DescendLtFromNil(t *testing.T) {
	testBTreeDescendLtFrom(t, nil)
}

func TestBTree_DescendLtFrom0(t *testing.T) {
	testBTreeDescendLtNum(t, 0)
}

func TestBTree_DescendLtFrom1(t *testing.T) {
	testBTreeDescendLtNum(t, 1)
}

func TestBTree_DescendLtFrom2(t *testing.T) {
	testBTreeDescendLtNum(t, 2)
}

func TestBTree_DescendLtFrom3(t *testing.T) {
	testBTreeDescendLtNum(t, 3)
}

func TestBTree_DescendLtFrom4(t *testing.T) {
	testBTreeDescendLtNum(t, 4)
}

func TestBTree_DescendLtFrom5(t *testing.T) {
	testBTreeDescendLtNum(t, 5)
}

func TestBTree_DescendLteFromNil(t *testing.T) {
	testBTreeDescendLteFrom(t, nil)
}

func TestBTree_DescendLteFrom0(t *testing.T) {
	testBTreeDescendLteNum(t, 0)
}

func TestBTree_DescendLteFrom1(t *testing.T) {
	testBTreeDescendLteNum(t, 1)
}

func TestBTree_DescendLteFrom2(t *testing.T) {
	testBTreeDescendLteNum(t, 2)
}

func TestBTree_DescendLteFrom3(t *testing.T) {
	testBTreeDescendLteNum(t, 3)
}

func TestBTree_DescendLteFrom4(t *testing.T) {
	testBTreeDescendLteNum(t, 4)
}

func TestBTree_DescendLteFrom5(t *testing.T) {
	testBTreeDescendLteNum(t, 5)
}

func testBTreeDescendLtNum(t *testing.T, k int) {
	testBTreeDescendLtFrom(t, &xCmp{x: k})
}

func testBTreeDescendLtFrom(t *testing.T, k Node) {
	testBTreeIter(t, testDescendLt, k)
}

func testBTreeDescendLteNum(t *testing.T, k int) {
	testBTreeDescendLteFrom(t, &xCmp{x: k})
}

func testBTreeDescendLteFrom(t *testing.T, k Node) {
	testBTreeIter(t, testDescendLte, k)
}

func testDescendLte(t *testing.T, b *BTree, k Node, n int) {
	testIter(t, k, b.DescendLte, n)
}

func testDescendLt(t *testing.T, b *BTree, k Node, n int) {
	testIter(t, k, b.DescendLt, n)
}
