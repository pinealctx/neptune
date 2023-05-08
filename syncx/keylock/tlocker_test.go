package keylock

import (
	"sync"
	"testing"
	"time"
)

// About 400ns
func TestTKeyLocker_Crazy1(t *testing.T) {
	var x = NewTKeyLocker[int64]()
	testTLockerCrazy(t, x)
}

// About 80ns
func TestTKeyLocker_Crazy2(t *testing.T) {
	var x = NewTKeyLockeGrp[int64]()
	testTLockerCrazy(t, x)
}

// About 120ns
func TestTKeyLocker_Crazy3(t *testing.T) {
	var x = NewTXHashTKeyLockeGrp[int64]()
	testTLockerCrazy(t, x)
}

func testTLockerCrazy(t *testing.T, x TLocker[int64]) {
	var wg sync.WaitGroup
	wg.Add(200)

	var count = int64(30000)
	var t1 = time.Now()

	//200 go routine
	//50 write go routine
	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()
			for j := int64(0); j < count; j++ {
				x.Lock(j % 256)
				x.Unlock(j % 256)
			}
		}()
	}

	//100 read go routine
	for i := 0; i < 150; i++ {
		go func() {
			defer wg.Done()
			for j := int64(0); j < count; j++ {
				x.RLock(j % 256)
				x.RUnlock(j % 256)
			}
		}()
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/(200*time.Duration(count)))
}
