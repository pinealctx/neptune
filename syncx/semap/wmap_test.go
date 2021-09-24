package semap

import "testing"

func TestWeightedW(t *testing.T) {
	testWeighted(t, NewWideSemMap(WithSize(1), WithRwRatio(1)))
}

func TestWeightedPanicW(t *testing.T) {
	sem := NewWideSemMap(WithSize(1), WithRwRatio(5))
	testWeightedPanic(t, sem)
}

func TestLockW(t *testing.T) {
	var sem = NewWideSemMap(WithSize(2), WithRwRatio(5))
	testLock(sem)
}
