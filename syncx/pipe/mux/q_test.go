package mux

import (
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

func TestActorQ_RandomAdd(t *testing.T) {
	var q = NewQ(0)
	testActorQ(q, time.Minute*5, false)
}

func TestActorQ_RandomAdd_PopAny(t *testing.T) {
	var q = NewQ(0)
	testActorQ(q, time.Minute*5, true)
}

func TestActorQ_RandomAddWithReqMax(t *testing.T) {
	var q = NewQ(5)
	testActorQ(q, time.Minute*5, false)
}

func testActorQ(q *Q, maxWaitTime time.Duration, popAny bool) {
	var startTime = time.Now()

	var waitGrp sync.WaitGroup
	waitGrp.Add(2)

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

			switch buf[0] % 2 {
			case 0:
				err = q.AddReq(buf[0])
			case 1:
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

	waitGrp.Wait()
}
