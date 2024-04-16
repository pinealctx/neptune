package slicex

import (
	"cmp"
	"math/rand"
	"time"
)

// Insert 往slice里面插入元素，如果index大于slice的长度，则追加在最后面，如果index小于0，则panic
func Insert[T any](s *[]T, index int, v T) {
	if index >= len(*s) {
		*s = append(*s, v)
		return
	}
	*s = append(*s, v)
	copy((*s)[index+1:], (*s)[index:])
	(*s)[index] = v
}

// Remove 移除index位置的元素，如果index小于0，则panic
func Remove[T any](s *[]T, index int) {
	if len(*s) == 0 {
		return
	}
	if index >= len(*s) {
		*s = (*s)[:len(*s)-1]
		return
	}
	copy((*s)[index:], (*s)[index+1:])
	*s = (*s)[:len(*s)-1]
}

// RemoveElem 删除第一个匹配的元素
func RemoveElem[T comparable](s *[]T, v T) {
	index := FindIndex[T](*s, v)
	if index == -1 {
		return
	}
	Remove(s, index)
}

// RemoveElems 删除第一个匹配的元素
func RemoveElems[T comparable](s *[]T, beRemoved []T) {
	for _, v := range beRemoved {
		RemoveElem(s, v)
	}
}

// FindIndex 查找第一个匹配的元素所在的index，返回-1代表没有找到
func FindIndex[T comparable](s []T, v T) int {
	for i, sv := range s {
		if sv == v {
			return i
		}
	}
	return -1
}

// Contain 是否包含指定元素
func Contain[T comparable](s []T, v T) bool {
	return FindIndex[T](s, v) != -1
}

// ContainFunc 是否包含指定元素
func ContainFunc[T any](s []T, f func(T) bool) bool {
	for _, v := range s {
		if f(v) {
			return true
		}
	}
	return false
}

// Clone 克隆一个slice
func Clone[T any](list []T) []T {
	return append([]T(nil), list...)
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Sum 求和
func Sum[T Number](list []T) T {
	var sum T
	for _, v := range list {
		sum += v
	}
	return sum
}

// Max 求最大值
func Max[T cmp.Ordered](list []T) T {
	l := len(list)
	if l == 0 {
		panic("empty list")
	}
	var maxV = list[0]
	for i := 1; i < l; i++ {
		if list[i] > maxV {
			maxV = list[i]
		}
	}
	return maxV
}

// Min 求最小值
func Min[T cmp.Ordered](list []T) T {
	l := len(list)
	if l == 0 {
		panic("empty list")
	}
	var minV = list[0]
	for i := 1; i < l; i++ {
		if list[i] < minV {
			minV = list[i]
		}
	}
	return minV
}

var (
	shuffleRand *rand.Rand
)

// Shuffle shuffles a list
func Shuffle[T any](list []T) []T {
	if shuffleRand == nil {
		shuffleRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
	for i := len(list) - 1; i > 0; i-- {
		j := shuffleRand.Intn(i + 1)
		list[i], list[j] = list[j], list[i]
	}
	return list
}
