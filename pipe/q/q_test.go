package q

import (
	"context"
	"crypto/rand"
	"fmt"
	"sync"
	"testing"
	"time"
)

type _Cond struct {
	lock sync.Mutex
	cond sync.Cond
}

func (c *_Cond) Init() {
	c.cond.L = &c.lock
}

func (c *_Cond) Wait() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cond.Wait()
}

func (c *_Cond) Wakeup() {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.cond.Broadcast()
}

func TestCondCase1(t *testing.T) {
	var x = &_Cond{}
	x.Init()
	x.Wakeup()
	x.Wakeup()
	x.Wakeup()
	x.Wait()
}

func TestCondCase2(t *testing.T) {
	var x = &_Cond{}
	x.Init()
	x.Wakeup()
	x.Wakeup()
	x.Wakeup()
	var wait sync.WaitGroup
	wait.Add(1)
	go func() {
		defer wait.Done()
		x.Wait()
	}()

	go func() {
		x.Wakeup()
	}()

	wait.Wait()
}

func TestActorQ_AddCtrl(t *testing.T) {
	var q = NewQ()
	_ = q.AddReq(-1)
	_ = q.AddReq(-2)
	_ = q.AddReq(-3)
	_ = q.AddReq(-4)
	_ = q.AddReq(-5)
	_ = q.AddCtrl(1)
	_ = q.AddCtrl(2)
	_ = q.AddCtrl(3)
	_ = q.AddCtrl(4)
	_ = q.AddCtrl(5)

	var i, _ = q.Pop()
	t.Log("1 -->", i)
	i, _ = q.Pop()
	t.Log("2 -->", i)
	i, _ = q.Pop()
	t.Log("3 -->", i)
	i, _ = q.Pop()
	t.Log("4 -->", i)
	i, _ = q.Pop()
	t.Log("5 -->", i)
	i, _ = q.Pop()
	t.Log("6 -->", i)
	i, _ = q.Pop()
	t.Log("7 -->", i)
	i, _ = q.Pop()
	t.Log("8 -->", i)
	i, _ = q.Pop()
	t.Log("9 -->", i)
	i, _ = q.Pop()
	t.Log("10 -->", i)
}

func TestActorQ_RandomAdd(t *testing.T) {
	var q = NewQ()
	testActorQ(q, time.Minute*5, false)
}

func TestActorQ_RandomAdd_PopAny(t *testing.T) {
	var q = NewQ()
	testActorQ(q, time.Minute*5, true)
}

func TestActorQ_RandomAddWithCtrlMax(t *testing.T) {
	var q = NewQ(WithQCtrlSize(5))
	testActorQ(q, time.Minute*5, false)
}

func TestActorQ_RandomAddWithReqMax(t *testing.T) {
	var q = NewQ(WithQReqSize(5))
	testActorQ(q, time.Minute*5, false)
}

func TestActorQ_RandomAddWithCtrlReqMax(t *testing.T) {
	var q = NewQ(WithQReqSize(5), WithQCtrlSize(5))
	testActorQ(q, time.Minute*5, true)
}

func TestActorQ_WaitClose(t *testing.T) {
	var q = NewQ()
	testActorNotify(q)
	_ = q.WaitClose(context.Background())
}

func TestActorQ_WaitClear(t *testing.T) {
	var q = NewQ()
	testActorNotify(q)
	_ = q.WaitClear(context.Background())
}

func testActorQ(q *Q, maxWaitTime time.Duration, popAny bool) {
	var startTime = time.Now()

	var waitGrp sync.WaitGroup
	waitGrp.Add(3)

	var fn = q.Pop
	if popAny {
		fn = q.PopAnyway
	}

	go func() {
		defer waitGrp.Done()
		defer q.Close()
		for {
			var buf [1]byte
			var _, err = rand.Read(buf[:])
			if err != nil {
				fmt.Println("produce return by error1:", err)
				return
			}

			switch buf[0] % 4 {
			case 0:
				err = q.AddCtrl(buf[0])
			case 1:
				err = q.AddPriorCtrl(buf[0])
			case 2:
				err = q.AddReq(buf[0])
			case 3:
				err = q.AddPriorReq(buf[0])
			}

			if err != nil {
				fmt.Println("produce return by error2:", err)
				return
			}
			fmt.Println("produce ---> ", buf[0])
			var now = time.Now()
			if now.Sub(startTime) > maxWaitTime {
				fmt.Println("test for 5 minutes, it's over")
				return
			}

			time.Sleep(time.Duration(buf[0]) * time.Microsecond)
		}
	}()

	go func() {
		defer waitGrp.Done()
		defer q.Close()
		for {
			var buf [1]byte
			var _, err = rand.Read(buf[:])
			if err != nil {
				fmt.Println("consume return by error1:", err)
				return
			}
			var i interface{}
			i, err = fn()
			if err != nil {
				fmt.Println("consume return by error2:", err)
				return
			}
			fmt.Println("consumer  --> ", i)
			time.Sleep(time.Duration(buf[0]) * time.Microsecond)
		}
	}()

	go func() {
		defer waitGrp.Done()
		for {
			var buf [1]byte
			var _, err = rand.Read(buf[:])
			if err != nil {
				fmt.Println("try close return by error:", err)
				return
			}
			var closed = q.TryClose()
			fmt.Println("try close-->:", closed)
			if closed {
				return
			}
			time.Sleep(time.Duration(buf[0]) * time.Second / 2)
		}
	}()

	waitGrp.Wait()
}

func testActorNotify(q *Q) {
	var waitGrp sync.WaitGroup
	waitGrp.Add(3)
	go func() {
		var buf [1]byte
		defer waitGrp.Done()
		defer q.Close()
		for i := 0; i < 1000; i++ {
			var _, err = rand.Read(buf[:])
			if err != nil {
				fmt.Println("produce return by error1:", err)
				return
			}
			err = q.AddCtrl(buf[0])
			if err != nil {
				fmt.Println("produce return by error2:", err)
				return
			}
			time.Sleep(time.Microsecond * time.Duration(buf[0]))
		}
	}()

	go func() {
		defer waitGrp.Done()
		defer q.TryClear()
		for {
			var buf [1]byte
			var _, err = rand.Read(buf[:])
			if err != nil {
				fmt.Println("consume return by error1:", err)
				return
			}
			var i interface{}
			i, err = q.PopAnyway()
			if err != nil {
				fmt.Println("consume return by error2:", err)
				return
			}
			fmt.Println("consumer  --> ", i)
			time.Sleep(time.Duration(buf[0]) * time.Microsecond)
		}
	}()

	go func() {
		defer waitGrp.Done()
		for {
			var buf [1]byte
			var _, err = rand.Read(buf[:])
			if err != nil {
				fmt.Println("try close return by error:", err)
				return
			}
			var closed = q.TryClose()
			fmt.Println("try close-->:", closed)
			if closed {
				return
			}
			time.Sleep(time.Duration(buf[0]) * time.Second / 2)
		}
	}()

	waitGrp.Wait()
}
