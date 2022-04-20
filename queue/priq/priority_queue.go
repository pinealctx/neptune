// Package priq 带优先级功的队列，可以进行select操作，并且有容量限制
package priq

import (
	"container/heap"
	"errors"
	"sync"
)

var ErrQueueIsFull = errors.New("queue.is.full")

type IEntry interface {
	GetPriority() int
}

type wrapEntry struct {
	entry IEntry
	seq   int64
}

type EntryList []*wrapEntry

func (e EntryList) Len() int {
	return len(e)
}

// Less 根据less从小到大排列，排在前面的先Pop（也就是优先级更高）
func (e EntryList) Less(i, j int) bool {
	pi := e[i].entry.GetPriority()
	pj := e[j].entry.GetPriority()
	// 同样的优先级，早入队的，优先级更高
	if pi == pj {
		return e[i].seq < e[j].seq
	} else {
		return pi > pj
	}
}

func (e EntryList) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e *EntryList) Push(x interface{}) {
	*e = append(*e, x.(*wrapEntry))
}

func (e *EntryList) Pop() interface{} {
	head := (*e)[len(*e)-1]
	// 避免可能的引用内存不释放
	(*e)[len(*e)-1] = nil
	*e = (*e)[:len(*e)-1]

	return head
}

type PriQueue struct {
	mu       sync.Mutex
	entries  EntryList
	capacity int
	signal   chan struct{}
	curSeq   int64
}

func NewPriQueue(capability int) *PriQueue {
	p := &PriQueue{}
	p.capacity = capability
	p.signal = make(chan struct{}, 1)

	return p
}

// Pop 不阻塞，如果返回nil代表当前队列是空的
func (pq *PriQueue) Pop() IEntry {

	pq.mu.Lock()
	if len(pq.entries) == 0 {
		pq.mu.Unlock()
		return nil
	}
	e := heap.Pop(&pq.entries).(*wrapEntry)
	needSignal := len(pq.entries) > 0
	// mu只锁entries
	pq.mu.Unlock()
	if needSignal {
		pq.tyrSignal()
	}

	return e.entry
}

func (pq *PriQueue) Len() int {

	pq.mu.Lock()
	defer pq.mu.Unlock()

	return len(pq.entries)
}

// Push 不阻塞
func (pq *PriQueue) Push(e IEntry) error {

	pq.mu.Lock()
	if len(pq.entries) >= pq.capacity {
		pq.mu.Unlock()
		return ErrQueueIsFull
	}
	pq.curSeq++
	heap.Push(&pq.entries, &wrapEntry{
		entry: e,
		seq:   pq.curSeq,
	})
	pq.mu.Unlock()
	pq.tyrSignal()

	return nil
}

// WaitCh 对 WaitCh() 返回的channel进行select，然后再使用Pop
func (pq *PriQueue) WaitCh() <-chan struct{} {
	return pq.signal
}

func (pq *PriQueue) tyrSignal() {
	select {
	case pq.signal <- struct{}{}:
	default:
	}
}
