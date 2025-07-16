package async

import (
	"sync"
	"testing"
	"time"
)

func TestQ_Add(t *testing.T) {
	var q = NewQ(1)
	var err = q.Add(1)
	if err != nil {
		t.Fail()
		return
	}
	err = q.Add(2)
	if err != ErrFull {
		t.Fail()
		return
	}
}

func TestQ_Close(t *testing.T) {
	var q = NewQ(1)
	q.Close()
	var err = q.Add(1)
	if err != ErrClosed {
		t.Fail()
		return
	}
	_, err = q.PopAnyway()
	if err != ErrClosed {
		t.Fail()
	}
}

func TestQ_Pop(t *testing.T) {
	var q = NewQ(1024)
	var wg sync.WaitGroup
	wg.Add(4)

	go func() {
		var e error
		defer wg.Done()
		for i := 0; i < 9; i++ {
			e = q.Add(i)
			if e != nil {
				panic(e)
			}
		}
		time.Sleep(time.Second * 10)
		e = q.Add(9)
		if e != nil {
			panic(e)
		}
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 15)
		q.Close()
		q.Close()
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 20)
		var e = q.Add(0)
		if e != ErrClosed {
			panic("not.closed")
		}
		t.Log(time.Now(), e)
	}()

	go func() {
		defer wg.Done()
		time.Sleep(time.Second * 1)
		for {
			var r, err = q.PopAnyway()
			if err != nil {
				t.Log(time.Now(), err)
				return
			}
			t.Log(time.Now(), r)
		}
	}()
	wg.Wait()
}
