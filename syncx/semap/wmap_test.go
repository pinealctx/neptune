package semap

import "testing"

func TestWeightedW(t *testing.T) {
	testWeighted(t, NewWideSemMap(WithRwRatio(1)))
}

func TestWeightedPanicW(t *testing.T) {
	sem := NewWideSemMap(WithRwRatio(5))
	testWeightedPanic(t, sem)
}

func TestLockW(t *testing.T) {
	var sem = NewWideSemMap(WithRwRatio(5))
	testLock(sem)
}
