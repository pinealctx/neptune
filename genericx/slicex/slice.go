package slicex

// Insert 往slice里面插入元素，如果index大于slice的长度，则追加在最后面，如果index小于0，则panic
func Insert[T any](s []T, index int, v T) []T {
	if index >= len(s) {
		return append(s, v)
	}
	s = append(s, v)
	copy(s[index+1:], s[index:])
	s[index] = v
	return s
}

// Remove 移除index位置的元素，如果index小于0，则panic
func Remove[T any](s []T, index int) []T {
	if len(s) == 0 {
		return s
	}
	if index >= len(s) {
		return s[:len(s)-1]
	}
	copy(s[index:], s[index+1:])
	return s[:len(s)-1]
}

// RemoveElem 删除第一个匹配的元素
func RemoveElem[T comparable](s []T, v T) []T {
	index := FindIndex[T](s, v)
	if index == -1 {
		return s
	}
	return Remove(s, index)
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
