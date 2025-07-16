package slicex

import (
	"fmt"
	"testing"
)

func TestInsert(_ *testing.T) {
	var a []int
	Insert[int](&a, 1, 10)
	Insert[int](&a, 1, 11)
	Insert[int](&a, 1, 12)
	Insert[int](&a, 0, 13)
	Insert[int](&a, 1, 15)
	Insert[int](&a, 100, 16)
	Insert[int](&a, 100, 17)
	Insert[int](&a, 100, 18)
	fmt.Printf("%+v\n", a)

	Remove[int](&a, 100)
	fmt.Printf("%+v\n", a)

	Remove[int](&a, 0)
	fmt.Printf("%+v\n", a)

	Remove[int](&a, 1)
	fmt.Printf("%+v\n", a)

	RemoveElem[int](&a, 100)
	fmt.Printf("%+v\n", a)

	RemoveElem[int](&a, 17)
	fmt.Printf("%+v\n", a)

	fmt.Printf("%+v\n", Contain[int](a, 12))
	fmt.Printf("%+v\n", Contain[int](a, 1))

	fmt.Printf("%+v\n", FindIndex[int](a, 16))
	fmt.Printf("%+v\n", FindIndex[int](a, 160))
}
