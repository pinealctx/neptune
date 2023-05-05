package tree

import (
	"github.com/pinealctx/neptune/ds/tree/btree"
	"sync"
)

// Node : redirect btree.Item
type Node = btree.Item

// FilterFn : filter a node value is in condition or not
type FilterFn func(n Node) bool

// BTree : wrapped b-tree
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

// Insert : insert node to btree
func (b *BTree) Insert(v Node) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.t.ReplaceOrInsert(v)
}

// Update : update old node to given new node.
// If old node not exist, the newV could not be updated.
// Return : if old node not exist, return false, else return true.
func (b *BTree) Update(oldV Node, newV Node) bool {
	b.rw.Lock()
	defer b.rw.Unlock()
	var e = b.t.Delete(oldV)
	if e == nil {
		return false
	}
	b.t.ReplaceOrInsert(newV)
	return true
}

// UpdateOrInsert : if ole node exists, update old node to new node, else insert new node to btree.
// The new node will always be inserted or replaced.
// Return : bool
// If ole node not exist, return false, it indicates that ole node not found but new node be inserted.
// Else return ture.
func (b *BTree) UpdateOrInsert(oldV Node, newV Node) bool {
	b.rw.Lock()
	defer b.rw.Unlock()
	var e = b.t.Delete(oldV)
	b.t.ReplaceOrInsert(newV)
	return e != nil
}

// Delete : delete node, actually, figure out a node which related node sort fields match.
// Return : bool
// If deleted node exist, return true, else return false
func (b *BTree) Delete(k Node) bool {
	b.rw.Lock()
	defer b.rw.Unlock()
	var e = b.t.Delete(k)
	return e != nil
}

// Get : get node by key
func (b *BTree) Get(k Node) Node {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.t.Get(k)
}

// AscendGte : ascend get nodes(>=k).
// k : anchor key
// filter : filter a node
// n : the max length of nodes to get
func (b *BTree) AscendGte(k Node, filter FilterFn, n int) []Node {
	return b.iterWalk(k, b.t.AscendGreaterOrEqual, filter, n)
}

// AscendGt : ascend get nodes(>k).
// k : anchor key
// filter : filter a node
// n : the max length of nodes to get
func (b *BTree) AscendGt(k Node, filter FilterFn, n int) []Node {
	return b.iterWalk(k, b.t.AscendGreater, filter, n)
}

// DescendLte : descend get nodes(<=k).
// k : anchor key
// filter : filter a node
// n : the max length of nodes to get
func (b *BTree) DescendLte(k Node, filter FilterFn, n int) []Node {
	return b.iterWalk(k, b.t.DescendLessOrEqual, filter, n)
}

// DescendLt : descend get nodes(<k).
// k : anchor key
// filter : filter a node
// n : the max length of nodes to get
func (b *BTree) DescendLt(k Node, filter FilterFn, n int) []Node {
	return b.iterWalk(k, b.t.DescendLess, filter, n)
}

// _gBtreeIterWrap : btree iter function, wrap google btree iterator function
// pivot : anchor key
// iterator : the iter continue or not function
type _gBtreeIterWrap func(pivot btree.Item, iterator btree.ItemIterator)

// iterWalk : iter node and gather filtered nodes.
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
