package semap

import "testing"

func TestWeightedW(t *testing.T) {
	testWeighted(t, NewWideSemMap(1, 1))
}

func TestWeightedPanicW(t *testing.T) {
	sem := NewWideSemMap(1, 5)
	testWeightedPanic(t, sem)
}

func TestLockW(t *testing.T) {
	var sem = NewWideSemMap(2, 5)
	testLock(sem)
}
