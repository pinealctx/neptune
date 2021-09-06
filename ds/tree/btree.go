package tree

import (
	"github.com/pinealctx/neptune/ds/tree/btree"
	"sync"
)

// Node : redirect btree.Item
type Node = btree.Item

//FilterFn : filter a node value is in condition or not
type FilterFn func(n Node) bool

//BTree : wrapped b-tree
type BTree struct {
	t  *btree.BTree
	rw *sync.RWMutex
}

// NewBTree : new
func NewBTree() *BTree {
	var b = &BTree{}
	b.t = btree.New(2)
	b.rw = &sync.RWMutex{}
	return b
}

//Insert : insert node to btree
func (b *BTree) Insert(v Node) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.t.ReplaceOrInsert(v)
}


//Update : update old node to given new node.
func (b *BTree) Update(oldV Node, newV Node) {
	b.rw.Lock()
	defer b.rw.Unlock()
	var e = b.t.Delete(oldV)
	if e == nil {
		return
	}
	b.t.ReplaceOrInsert(newV)
}

//UpdateOrInsert : if ole node exists, update old node to new node, else insert new node to btree.
func (b *BTree) UpdateOrInsert(oldV Node, newV Node) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.t.Delete(oldV)
	b.t.ReplaceOrInsert(newV)
}

//Delete : delete node, actually, figure out a node which related node sort fields match.
func (b *BTree) Delete(k Node) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.t.Delete(k)
}

func (b *BTree) Get(k Node) Node {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.t.Get(k)
}

//AscendGte : ascend get nodes(>=k).
//k : anchor key
//filter : filter a node
//n : the max length of nodes to get
func (b *BTree) AscendGte(k Node, filter FilterFn, n int) []Node {
	return b.iterWalk(k, b.t.AscendGreaterOrEqual, filter, n)
}

//AscendGt : ascend get nodes(>k).
//k : anchor key
//filter : filter a node
//n : the max length of nodes to get
func (b *BTree) AscendGt(k Node, filter FilterFn, n int) []Node {
	return b.iterWalk(k, b.t.AscendGreater, filter, n)
}

//DescendLte : descend get nodes(<=k).
//k : anchor key
//filter : filter a node
//n : the max length of nodes to get
func (b *BTree) DescendLte(k Node, filter FilterFn, n int) []Node {
	return b.iterWalk(k, b.t.DescendLessOrEqual, filter, n)
}

//DescendLt : descend get nodes(<k).
//k : anchor key
//filter : filter a node
//n : the max length of nodes to get
func (b *BTree) DescendLt(k Node, filter FilterFn, n int) []Node {
	return b.iterWalk(k, b.t.DescendLess, filter, n)
}

//_gBtreeIterWrap : btree iter function, wrap google btree iterator function
//pivot : anchor key
//iterator : the iter continue or not function
type _gBtreeIterWrap func(pivot btree.Item, iterator btree.ItemIterator)

//iterWalk : iter node and gather filtered nodes.
func (b *BTree) iterWalk(k Node, iterFn _gBtreeIterWrap, filter FilterFn, n int) []Node {
	if n == 0 {
		return nil
	}
	var ns = make([]Node, 0, n)
	var c = 0
	var fn = func(v Node) bool {
		if c >= n {
			return false
		}
		if filter(v) {
			ns = append(ns, v)
			c++
		}
		return true
	}

	b.rw.RLock()
	defer b.rw.RUnlock()
	iterFn(k, fn)
	return ns
}
