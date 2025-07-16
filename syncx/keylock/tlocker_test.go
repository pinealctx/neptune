package keylock

import (
	"sort"
	"sync"
	"testing"
	"time"
)

// About 400ns
func TestTKeyLocker_Crazy1(t *testing.T) {
	var x = NewTKeyLocker[int]()
	testTLockerCrazy(t, x)
}

// About 80ns
func TestTKeyLocker_Crazy2(t *testing.T) {
	var x = NewTKeyLockeGrp[int]()
	testTLockerCrazy(t, x)
}

// About 120ns
func TestTKeyLocker_Crazy3(t *testing.T) {
	var x = NewTXHashTKeyLockeGrp[int]()
	testTLockerCrazy(t, x)
}

func TestTKeyLocker_RandCrazyA1(t *testing.T) {
	var x = NewTKeyLocker[int]()
	testTLockerRandCrazyA(t, x)
}

func TestTKeyLocker_RandCrazyA2(t *testing.T) {
	var x = NewTKeyLockeGrp[int]()
	testTLockerRandCrazyA(t, x)
}

func TestTKeyLocker_RandCrazyA3(t *testing.T) {
	var x = NewTXHashTKeyLockeGrp[int]()
	testTLockerRandCrazyA(t, x)
}

func TestTKeyLocker_RandCrazyB1(t *testing.T) {
	var x = NewTKeyLocker[int]()
	testTLockerRandCrazyB(t, x)
}

func TestTKeyLocker_RandCrazyB2(t *testing.T) {
	var x = NewTKeyLockeGrp[int]()
	testTLockerRandCrazyB(t, x)
}

func TestTKeyLocker_RandCrazyB3(t *testing.T) {
	var x = NewTXHashTKeyLockeGrp[int]()
	testTLockerRandCrazyB(t, x)
}

func TestTKeyLocker_RandCrazyC1(t *testing.T) {
	var x = NewTKeyLocker[int]()
	testTLockerRandCrazyC(t, x)
}

func TestTKeyLocker_RandCrazyC2(t *testing.T) {
	var x = NewTKeyLockeGrp[int]()
	testTLockerRandCrazyC(t, x)
}

func TestTKeyLocker_RandCrazyC3(t *testing.T) {
	var x = NewTXHashTKeyLockeGrp[int]()
	testTLockerRandCrazyC(t, x)
}

func testTLockerCrazy(t *testing.T, x TLocker[int]) {
	t.Helper()
	var wg sync.WaitGroup
	wg.Add(200)

	var count = int(30000)
	var t1 = time.Now()

	//200 go routine
	//50 write go routine
	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()
			for j := int(0); j < count; j++ {
				x.Lock(j % 256)
				x.Unlock(j % 256)
			}
		}()
	}

	//150 read go routine
	for i := 0; i < 150; i++ {
		go func() {
			defer wg.Done()
			for j := int(0); j < count; j++ {
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

func testTLockerRandCrazyA(t *testing.T, x TLocker[int]) {
	t.Helper()
	var ks = make([][]int, 32)
	var index = int(0)
	for i := 0; i < 32; i++ {
		ks[i] = make([]int, 8)
		for j := 0; j < 8; j++ {
			ks[i][j] = index % 256
			index++
		}
	}
	t.Log("ks:", ks)
	testTLockerRandCrazy(t, x, ks)
}

func testTLockerRandCrazyB(t *testing.T, x TLocker[int]) {
	t.Helper()
	var iter = [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 7, 0},
		{2, 3, 4, 5, 6, 7, 0, 1},
		{3, 4, 5, 6, 7, 0, 1, 2},
		{4, 5, 6, 7, 0, 1, 2, 3},
		{5, 6, 7, 0, 1, 2, 3, 4},
		{6, 7, 0, 1, 2, 3, 4, 5},
		{7, 0, 1, 2, 3, 4, 5, 6},
	}
	var ks = make([][]int, 32)
	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			ks[i*8+j] = make([]int, 8)
			copy(ks[i*8+j], iter[j])
			for k := 0; k < 8; k++ {
				ks[i*8+j][k] += int(i)
			}
		}
	}

	t.Log("ks:", ks)
	testTLockerRandCrazy(t, x, ks)
}

func testTLockerRandCrazyC(t *testing.T, x TLocker[int]) {
	t.Helper()
	var iter = [][]int{
		{0, 1, 2, 3, 4, 5, 6, 7},
		{1, 2, 3, 4, 5, 6, 7, 0},
		{2, 3, 4, 5, 6, 7, 0, 1},
		{3, 4, 5, 6, 7, 0, 1, 2},
		{4, 5, 6, 7, 0, 1, 2, 3},
		{5, 6, 7, 0, 1, 2, 3, 4},
		{6, 7, 0, 1, 2, 3, 4, 5},
		{7, 0, 1, 2, 3, 4, 5, 6},
	}
	var ks = make([][]int, 32)
	for i := 0; i < 4; i++ {
		for j := 0; j < 8; j++ {
			ks[i*8+j] = make([]int, 8)
			copy(ks[i*8+j], iter[j])
			for k := 0; k < 8; k++ {
				ks[i*8+j][k] += int(i)
			}
			sort.Ints(ks[i*8+j])
		}
	}

	t.Log("ks:", ks)
	testTLockerRandCrazy(t, x, ks)
}

func testTLockerRandCrazy(t *testing.T, x TLocker[int], ks [][]int) {
	t.Helper()
	var wg sync.WaitGroup
	wg.Add(200)

	var count = int(30000)
	t.Log(ks)
	var t1 = time.Now()
	//200 go routine
	//50 write go routine
	for i := 0; i < 50; i++ {
		go func() {
			defer wg.Done()
			for j := int(0); j < count; j++ {
				x.Locks(ks[j%32])
				x.Unlocks(ks[j%32])
			}
		}()
		go func() {
			defer wg.Done()
			for j := int(0); j < count; j++ {
				x.RLocks(ks[j%32])
				x.RUnlocks(ks[j%32])
			}
		}()
		go func() {
			defer wg.Done()
			for j := int(0); j < count; j++ {
				x.RLocks(ks[j%32])
				x.RUnlocks(ks[j%32])
			}
		}()
		go func() {
			defer wg.Done()
			for j := int(0); j < count; j++ {
				x.RLocks(ks[j%32])
				x.RUnlocks(ks[j%32])
			}
		}()
	}

	wg.Wait()
	var t2 = time.Now()
	var d = t2.Sub(t1)
	t.Log("use time:", d, "average:", d/(200*time.Duration(count)))
}
