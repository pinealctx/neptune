package tree

import (
	"github.com/pinealctx/neptune/ds/tree/btree"
	"sync"
)

//IterFn : btree iter function
type IterFn func(pivot btree.Item, iterator btree.ItemIterator)

//FilterFn : filter a node value is in condition or not
type FilterFn func(n btree.Item) bool

type BTree struct {
	t *btree.BTree
	rw *sync.RWMutex
}

func NewBTree() *BTree {
	var b = &BTree{}
	b.t = btree.New(2)
	b.rw = &sync.RWMutex{}
	return b
}

func (b *BTree) Insert(v btree.Item) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.t.ReplaceOrInsert(v)
}

func (b *BTree) Update(oldV btree.Item, newV btree.Item) {
	b.rw.Lock()
	defer b.rw.Unlock()
	var e = b.t.Delete(oldV)
	if e == nil {
		return
	}
	b.t.ReplaceOrInsert(newV)
}

func (b *BTree) UpdateOrInsert(oldV btree.Item, newV btree.Item) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.t.Delete(oldV)
	b.t.ReplaceOrInsert(newV)
}

func (b *BTree) Delete(k btree.Item) {
	b.rw.Lock()
	defer b.rw.Unlock()
	b.t.Delete(k)
}

func (b *BTree) Get(k btree.Item) btree.Item {
	b.rw.RLock()
	defer b.rw.RUnlock()
	return b.t.Get(k)
}

func (b *BTree) AscendGte(k btree.Item, filter func(btree.Item) bool, n int) []btree.Item {
	return b.iter(k, b.t.AscendGreaterOrEqual, filter, n)
}

func (b *BTree) AscendGt(k btree.Item, filter func(btree.Item) bool, n int) []btree.Item {
	return b.iter(k, b.t.AscendGreater, filter, n)
}

func (b *BTree) DescendLte(k btree.Item, filter func(btree.Item) bool, n int) []btree.Item {
	return b.iter(k, b.t.DescendLessOrEqual, filter, n)
}

func (b *BTree) DescendLt(k btree.Item, filter func(btree.Item) bool, n int) []btree.Item {
	return b.iter(k, b.t.DescendLess, filter, n)
}

func (b *BTree) iter(k btree.Item, iterFn IterFn, filter FilterFn, n int) []btree.Item {
	if n == 0 {
		return nil
	}
	var ns = make([]btree.Item, 0, n)
	var c = 0
	var fn = func(v btree.Item) bool {
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