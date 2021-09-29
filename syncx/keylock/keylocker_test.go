package keylock

import (
	"sync"
	"testing"
	"time"
)

//About 500ns
func TestKeyLocker_Crazy(t *testing.T) {
	var x = NewKeyLocker()

	var wg sync.WaitGroup
	wg.Add(200)

	var count = time.Duration(30000)
	var t1 = time.Now()

	//200 go routine
	//50 write go routine
	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()
			for j := time.Duration(0); j < count; j++ {
				x.Lock(j % 256)
				x.Unlock(j % 256)
			}
		}()
	}

	//100 read go routine
	for i := 0; i < 150; i++ {
		go func() {
			defer wg.Done()
			for j := time.Duration(0); j < count; j++ {
				x.RLock(j % 256)
				x.RULock(j % 256)
			}
		}()
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/(200*count))
}
