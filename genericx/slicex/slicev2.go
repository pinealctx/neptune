package slicex

// Overwrite2FrontInPlace overwrite(copy) the elements of the slice to the front in place.
// The front elements will be overwritten, the other elements will remain unchanged.
// s - the slice to overwrite elements in
// startPos - the starting position of the elements to move
// length - the number of elements to move, if length is -1, it means to the end of the slice
// Example:
// s := []int{1, 2, 3, 4, 5, 6}
// Move2FrontInPlace(s, 2, 3)
// fmt.Println(s) // Output: [3, 4, 5, 4, 5, 6]
func Overwrite2FrontInPlace[T any](s []T, startPos int, length int) {
	// no need panic here because if something wrong, it will panic in copy function
	if length == -1 {
		copy(s, s[startPos:])
		return
	}
	copy(s[0:length], s[startPos:startPos+length])
}

// MergeSlice merges two slices into one slice.
// The order of elements in the merged slice is the same as the order of the input slices.
// If both input slices are nil or empty, it returns nil.
func MergeSlice[T any](xs ...[]T) []T {
	return MergeSlices(xs)
}

// MergeSlices merges multiple slices into one slice.
// The order of elements in the merged slice is the same as the order of the input slices.
// If all input slices are nil or empty, it returns nil.
func MergeSlices[T any](slices [][]T) []T {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	if totalLen == 0 {
		return nil
	}
	result := make([]T, 0, totalLen)
	for _, s := range slices {
		result = append(result, s...)
	}
	return result
}

// DeduplicationSlice : deduplicate a slice of comparable
func DeduplicationSlice[T comparable](slice []T) []T {
	size := len(slice)
	if size == 0 {
		return slice
	}
	s := make(map[T]struct{}, size)
	r := make([]T, 0, size)

	for _, item := range slice {
		_, ok := s[item]
		if !ok {
			s[item] = struct{}{}
			r = append(r, item)
		}
	}
	return r
}

// AppendByteSlice efficiently appends src to dst, automatically growing capacity if needed
// This replaces the manual capacity calculation and copy operations
// Returns the updated slice with src data appended
func AppendByteSlice[T any](dst, src []T) []T {
	if len(src) == 0 {
		return dst
	}

	currentLen := len(dst)
	requiredLen := currentLen + len(src)

	// Check if we need to grow capacity
	if cap(dst) < requiredLen {
		// Calculate new capacity (grow by 2x or required size, whichever is larger)
		newCap := cap(dst) * 2
		if newCap < requiredLen {
			newCap = requiredLen
		}

		// Allocate new slice and copy existing data
		newSlice := make([]T, currentLen, newCap)
		copy(newSlice, dst)
		dst = newSlice
	}

	// Extend slice length and copy new data
	dst = dst[:requiredLen]
	copy(dst[currentLen:], src)
	return dst
}
