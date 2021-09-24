package semap

import (
	"context"
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

const maxSleep = 1 * time.Millisecond

func TestWeighted(t *testing.T) {
	testWeighted(t, NewSemMap(1, 1))
}

func TestWeightedPanic(t *testing.T) {
	sem := NewSemMap(1, 5)
	testWeightedPanic(t, sem)
}

func TestLock(t *testing.T) {
	var sem = NewSemMap(2, 5)
	testLock(sem)
}

func testWeighted(t *testing.T, sem SemMapper) {
	t.Parallel()

	n := runtime.GOMAXPROCS(0)
	loops := 10000 / n
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			hammerWeighted(sem, int64(i), loops)
		}()
	}
	wg.Wait()
}

func testWeightedPanic(t *testing.T, sem SemMapper) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("release of an unacquired weighted semaphore did not panic")
		}
	}()
	var w, err = sem.AcquireRead(context.Background(), 1)
	if err != nil {
		t.Fail()
		panic(err)
	}
	sem.ReleaseWrite(1, w)
}

func testLock(sem SemMapper) {
	var count = 5
	var wg sync.WaitGroup
	wg.Add(5)

	go func() {
		defer wg.Done()
		for i := 0; i < count; i++ {
			var w, err = sem.AcquireWrite(context.Background(), 1)
			if err != nil {
				panic(err)
			}
			fmt.Println("in write", time.Now())
			time.Sleep(time.Second * 5)
			sem.ReleaseWrite(i, w)
		}
	}()

	for i := 0; i < 4; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < count; j++ {
				var w, err = sem.AcquireRead(context.Background(), 1)
				if err != nil {
					panic(err)
				}
				fmt.Println("in read", time.Now())
				time.Sleep(time.Second * 1)
				sem.ReleaseRead(j, w)
			}
		}()
	}
	wg.Wait()
}

func hammerWeighted(sem SemMapper, n int64, loops int) {
	for i := 0; i < loops; i++ {
		var w, err = sem.AcquireWrite(context.Background(), n)
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Duration(rand.Int63n(int64(maxSleep/time.Nanosecond))) * time.Nanosecond)
		sem.ReleaseWrite(n, w)
	}
}
