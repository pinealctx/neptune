package btree

// AscendGreater calls the iterator for every value in the tree within
// the range (pivot, last], until iterator returns false.
func (t *BTree) AscendGreater(pivot Item, iterator ItemIterator) {
	if t.root == nil {
		return
	}
	t.root.iterate(ascend, pivot, nil, false, false, iterator)
}

// DescendLess calls the iterator for every value in the tree within the range
// (pivot, first], until iterator returns false.
func (t *BTree) DescendLess(pivot Item, iterator ItemIterator) {
	if t.root == nil {
		return
	}
	t.root.iterate(descend, pivot, nil, false, false, iterator)
}
