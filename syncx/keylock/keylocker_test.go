package keylock

import (
	"sync"
	"testing"
	"time"
)

// About 400ns
func TestKeyLocker_Crazy1(t *testing.T) {
	var x = NewKeyLocker()
	testLockerCrazy(t, x)
}

// About 80ns
func TestKeyLocker_Crazy2(t *testing.T) {
	var x = NewKeyLockeGrp()
	testLockerCrazy(t, x)
}

// About 120ns
func TestKeyLocker_Crazy3(t *testing.T) {
	var x = NewXHashKeyLockeGrp()
	testLockerCrazy(t, x)
}

func testLockerCrazy(t *testing.T, x Locker) {
	t.Helper()
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
