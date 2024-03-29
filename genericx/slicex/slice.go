package slicex

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
