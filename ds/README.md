## ds一些数据结构

### ds/tree
封装google/btree成了线程安全的btree，并简化了函数。

主要的函数是新建树，增/删/改/查，按指定位置和方向遍历树。

```go

// NewBTree : new
func NewBTree() *BTree

//Insert : insert node to btree
func (b *BTree) Insert(v Node) 

//Update : update old node to given new node.
//If old node not exist, the newV could not be updated.
//Return : if old node not exist, return false, else return true.
func (b *BTree) Update(oldV Node, newV Node) bool 

//UpdateOrInsert : if ole node exists, update old node to new node, else insert new node to btree.
//The new node will always be inserted or replaced.
//Return : bool
//If ole node not exist, return false, it indicates that ole node not found but new node be inserted.
//Else return ture.
func (b *BTree) UpdateOrInsert(oldV Node, newV Node) bool 

//Delete : delete node, actually, figure out a node which related node sort fields match.
//Return : bool
//If deleted node exist, return true, else return false
func (b *BTree) Delete(k Node) bool 

//Get : get node by key
func (b *BTree) Get(k Node) Node 

//AscendGte : ascend get nodes(>=k).
//k : anchor key
//filter : filter a node
//n : the max length of nodes to get
func (b *BTree) AscendGte(k Node, filter FilterFn, n int) []Node 

//AscendGt : ascend get nodes(>k).
//k : anchor key
//filter : filter a node
//n : the max length of nodes to get
func (b *BTree) AscendGt(k Node, filter FilterFn, n int) []Node 

//DescendLte : descend get nodes(<=k).
//k : anchor key
//filter : filter a node
//n : the max length of nodes to get
func (b *BTree) DescendLte(k Node, filter FilterFn, n int) []Node 

//DescendLt : descend get nodes(<k).
//k : anchor key
//filter : filter a node
//n : the max length of nodes to get
func (b *BTree) DescendLt(k Node, filter FilterFn, n int) []Node
```