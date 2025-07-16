package tree

import (
	"sync"
	"testing"
	"time"

	"github.com/pinealctx/neptune/ds/tree/btree"
)

type X struct {
	a int
}

func (x *X) Less(b Node) bool {
	return x.a < b.(*X).a
}

// TestBtreeCmp :
// 此用例测试结果说明，在index变化后，不能直接插入，需要先删除，再插入
func TestBtreeCmp(t *testing.T) {
	var b = btree.New(2)
	var x = &X{1}
	insertX(t, b, &X{9})
	insertX(t, b, &X{8})
	insertX(t, b, &X{7})
	insertX(t, b, &X{6})
	insertX(t, b, &X{-1})
	insertX(t, b, &X{-2})
	insertX(t, b, &X{-3})
	insertX(t, b, x)
	t.Log("len init:", b.Len())
	x.a = 200
	t.Log("len set", b.Len())
	insertX(t, b, x)
	t.Log("len after:", b.Len())

	b.Ascend(func(i Node) bool {
		t.Logf("out: v:%+v, p:%p\n", i.(*X).a, i.(*X))
		return true
	})
}

type Int int

func (i Int) Less(b Node) bool {
	return i < b.(Int)
}

// 测试多个go routine不加锁读
// 一个go routine写的情况
// 事实证明，it's a joke, no cow here.
func TestBtreeNoLock(t *testing.T) {
	var b = btree.New(2)
	var wg sync.WaitGroup

	var c = 3
	wg.Add(c)
	var t1 = time.Now()
	//1万次
	var count = 1 * 10000
	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			b.ReplaceOrInsert(Int(i % 1023))
		}
	}()

	for i := 0; i < c-1; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				b.Ascend(func(_ Node) bool {
					return true
				})
			}
		}()
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/time.Duration(count*c))
}

// 测试多个go routine加读锁读
// 一个go routine写加写锁的情况
// 遍历1023节点长度的数耗时在9us左右 -- 因为都使用了写锁
func TestBtreeLock(t *testing.T) {
	var b = btree.New(2)
	var wg sync.WaitGroup
	var l sync.RWMutex

	var c = 3
	wg.Add(c)
	var t1 = time.Now()
	//1万次
	var count = 1 * 10000
	go func() {
		defer func() {
			wg.Done()
		}()
		for i := 0; i < count; i++ {
			l.Lock()
			b.ReplaceOrInsert(Int(i % 1023))
			l.Unlock()
		}
	}()

	for i := 0; i < c-1; i++ {
		go func(_ int) {
			defer func() {
				wg.Done()
			}()
			for j := 0; j < count; j++ {
				l.RLock()
				b.Ascend(func(_ Node) bool {
					return true
				})
				l.RUnlock()
			}
		}(i)
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/time.Duration(count*c))
}

// 测试多个go routine加写锁锁读
// 一个go routine写加写锁的情况
// 遍历1023节点长度的数耗时在9us左右 -- 因为都使用了写锁
func TestBtreeWLock(t *testing.T) {
	var b = btree.New(2)
	var wg sync.WaitGroup
	var l sync.Mutex

	var c = 3
	wg.Add(c)
	var t1 = time.Now()
	//1万次
	var count = 1 * 10000
	go func() {
		defer func() {
			wg.Done()
		}()
		for i := 0; i < count; i++ {
			l.Lock()
			b.ReplaceOrInsert(Int(i % 1023))
			l.Unlock()
		}
	}()

	for i := 0; i < c-1; i++ {
		go func(_ int) {
			defer func() {
				wg.Done()
			}()
			for j := 0; j < count; j++ {
				l.Lock()
				b.Ascend(func(_ Node) bool {
					return true
				})
				l.Unlock()
			}
		}(i)
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/time.Duration(count*c))
}

// 测试多个go routine读不加锁
// 遍历1023节点长度的数耗时在4us左右 -- 这是因为多核
func TestBtreeNoLockRead(t *testing.T) {
	var b = btree.New(2)
	var wg sync.WaitGroup

	var c = 3
	wg.Add(c - 1)
	var t1 = time.Now()
	//1万次
	var count = 1 * 10000

	for i := 0; i < count; i++ {
		b.ReplaceOrInsert(Int(i % 1023))
	}

	for i := 0; i < c-1; i++ {
		go func(_ int) {
			defer func() {
				wg.Done()
			}()
			for j := 0; j < count; j++ {
				b.Ascend(func(_ Node) bool {
					return true
				})
			}
		}(i)
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/time.Duration(count*c))
}

// 测试多个go routine读只加加锁
// 遍历1023节点长度的数耗时在4us左右 -- 这是因为多核
func TestBtreeLockRead(t *testing.T) {
	var b = btree.New(2)
	var wg sync.WaitGroup
	var l sync.RWMutex

	var c = 3
	wg.Add(c - 1)
	var t1 = time.Now()
	//1万次
	var count = 1 * 10000

	for i := 0; i < count; i++ {
		b.ReplaceOrInsert(Int(i % 1023))
	}

	for i := 0; i < c-1; i++ {
		go func(_ int) {
			defer func() {
				wg.Done()
			}()
			for j := 0; j < count; j++ {
				l.RLock()
				b.Ascend(func(_ Node) bool {
					return true
				})
				l.RUnlock()
			}
		}(i)
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/time.Duration(count*c))
}

// 测试纯遍历耗时
// 遍历1023节点长度的数耗时在11us左右
func TestBtreeWalk(t *testing.T) {
	var b = btree.New(2)

	//10万次
	var count = 10 * 10000

	for i := 0; i < count; i++ {
		b.ReplaceOrInsert(Int(i % 1023))
	}
	var t1 = time.Now()
	for i := 0; i < count; i++ {
		b.Ascend(func(_ Node) bool {
			return true
		})
	}

	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("tree len:", b.Len(), "use time:", d, "average:", d/time.Duration(count))
}

func insertX(t *testing.T, b *btree.BTree, x *X) {
	t.Helper()
	var y = b.ReplaceOrInsert(x)
	if y != nil {
		t.Log("already:", y.(*X).a)
	}
}
