package priq

import (
	"fmt"
	"testing"
	"time"
)

type qItem int

func (qi qItem) GetPriority() int {
	return int(qi)
}

func assert(t *testing.T, expect bool) {
	if !expect {
		t.FailNow()
	}
}

func TestPriQueue_Len(t *testing.T) {

	cases := []int{1, -10, 15, 3, 7, 5, 100}
	q := NewPriQueue(len(cases) - 1)
	for i, v := range cases {
		if i < len(cases)-1 {
			err := q.Push(qItem(v))
			assert(t, err == nil)
			assert(t, q.Len() == i+1)
		} else {
			err := q.Push(qItem(v))
			assert(t, err != nil)
		}
	}
	for i := q.Len(); i > 0; i-- {
		q.Pop()
		assert(t, q.Len() == i-1)
	}
}

func TestPriQueue_op(t *testing.T) {

	cases := []int{1, -10, 15, 3, 7, 5, 100, -2}
	q := NewPriQueue(len(cases))
	for _, v := range cases {
		err := q.Push(qItem(v))
		assert(t, err == nil)
	}
	for {
		e := q.Pop()
		if e == nil {
			break
		}
		fmt.Printf("%d,", e.(qItem))
	}
	fmt.Println()
}

func TestPriQueue_Use(t *testing.T) {
	q := NewPriQueue(1000 * 1000 * 1000)
	cases := []int{1, -10, 15, 3, 7, 5, 100, -2}
	for _, v := range cases {
		err := q.Push(qItem(v))
		assert(t, err == nil)
	}

	go func() {
		for {
			select {
			case <-q.WaitCh():
				e := q.Pop()
				if e != nil {
					fmt.Printf("%d\n", e.(qItem))
				}
			}
		}

	}()

	go func() {
		for i := 0; i < 15; i++ {
			_ = q.Push(qItem(i))
			if i%3 == 0 {
				time.Sleep(450 * time.Millisecond)
			}
		}
	}()

	time.Sleep(10 * time.Second)
}

// 500ns级别
func BenchmarkPriQueue_Push(b *testing.B) {
	q := NewPriQueue(1000 * 1000 * 1000)
	// 先填充100万个作为基础
	for i := 0; i < 1000*1000; i++ {
		err := q.Push(qItem(i))
		if err != nil {
			b.FailNow()
		}
	}
	for j := 0; j < b.N; j++ {
		_ = q.Push(qItem(j))
		//if err != nil {
		//	b.FailNow()
		//}
	}
}

// 50ns级别
func BenchmarkPriQueue_Pop(b *testing.B) {
	q := NewPriQueue(1000 * 1000 * 1000)
	// 先填充100万个作为基础
	for i := 0; i < 1000*1000; i++ {
		err := q.Push(qItem(i))
		if err != nil {
			b.FailNow()
		}
	}
	for j := 0; j < b.N; j++ {
		q.Pop()
	}
}

type TOrder struct {
	Pri   int
	Value int
}

func (o *TOrder) GetPriority() int {
	return o.Pri
}

func TestOrder(t *testing.T) {
	q := NewPriQueue(100)
	for i := 0; i < 100; i++ {
		pri := i / 10
		if i > 50 {
			pri = pri - 5
		}
		_ = q.Push(&TOrder{
			Pri:   pri,
			Value: i,
		})
	}

	for {
		e := q.Pop()
		if e == nil {
			break
		}
		o := e.(*TOrder)
		fmt.Printf("%d : %d\n", o.Pri, o.Value)
	}

}
