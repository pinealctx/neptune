package slicex

import (
	"fmt"
	"testing"
)

func TestInsert(t *testing.T) {
	var a []int
	a = Insert[int](a, 1, 10)
	a = Insert[int](a, 1, 11)
	a = Insert[int](a, 1, 12)
	a = Insert[int](a, 0, 13)
	a = Insert[int](a, 1, 15)
	a = Insert[int](a, 100, 16)
	a = Insert[int](a, 100, 17)
	a = Insert[int](a, 100, 18)
	fmt.Printf("%+v\n", a)

	a = Remove[int](a, 100)
	fmt.Printf("%+v\n", a)

	a = Remove[int](a, 0)
	fmt.Printf("%+v\n", a)

	a = Remove[int](a, 1)
	fmt.Printf("%+v\n", a)

	a = RemoveElem[int](a, 100)
	fmt.Printf("%+v\n", a)

	a = RemoveElem[int](a, 17)
	fmt.Printf("%+v\n", a)

	fmt.Printf("%+v\n", Contain[int](a, 12))
	fmt.Printf("%+v\n", Contain[int](a, 1))

	fmt.Printf("%+v\n", FindIndex[int](a, 16))
	fmt.Printf("%+v\n", FindIndex[int](a, 160))
}
