package mem

import (
	"errors"
	"fmt"
)

// ClaySlice represents the non-owning reference structure (arrayName##Slice)
// used throughout Clay to point to a sub-range of any backing array.
type MemSlice[T any] struct {
	internalArray []T
}

func (c *MemSlice[T]) InternalArray() []T {
	return c.internalArray
}

func (s *MemSlice[T]) Capacity() int32 {
	return int32(cap(s.internalArray))
}
func (s *MemSlice[T]) Grow(length int32) {
	if s.Length()+length > s.Capacity() {
		panic(fmt.Sprintf("MemSlice.Grow capacity exceeded: %d + %d > %d", s.Length(), length, s.Capacity()))
	}
	s.internalArray = s.internalArray[:s.Length()+length]
}

func (s *MemSlice[T]) Shrink(length int32) {
	if s.Length() < length {
		panic(fmt.Sprintf("MemSlice.Shrink length is greater than the array length: %d < %d", s.Length(), length))
	}
	s.internalArray = s.internalArray[:s.Length()-length]
}
func (s *MemSlice[T]) Length() int32 {
	return int32(len(s.internalArray))
}

func NewMemSlice[T any](length int32) MemSlice[T] {
	initialSlice := make([]T, length)
	return MemSlice[T]{
		internalArray: initialSlice,
	}
}

func NewMemSliceWithData[T any](array []T) MemSlice[T] {
	return MemSlice[T]{
		internalArray: array,
	}
}

func MSlice_Set[T any](slice *MemSlice[T], index int32, item T) {
	if !rangeCheck(index, slice.Length()) {
		message := fmt.Sprintf("MemSlice.MSlice_Set index: %d, slice.Length(): %d\n", index, slice.Length())
		// fmt.Println(message)
		panic(message)
		// return
	}
	internalArray := slice.InternalArray()
	internalArray[index] = item
}

func MSlice_Get[T any](slice *MemSlice[T], index int32) *T {
	if !rangeCheck(index, slice.Length()) {
		message := fmt.Sprintf("MemSlice.MSlice_Get index: %d, slice.Length(): %d\n", index, slice.Length())
		// fmt.Println(message)
		panic(message)
		// return nil
	}
	internalArray := slice.InternalArray()
	return &internalArray[index]
}

func MSlice_GetValue[T any](slice *MemSlice[T], index int32) T {
	if !rangeCheck(index, slice.Length()) {
		message := fmt.Sprintf("MemSlice.MSlice_GetValue index: %d, slice.Length(): %d\n", index, slice.Length())
		// fmt.Println(message)
		panic(message)
		// zero := new(T)
		// return *zero
	}
	internalArray := slice.InternalArray()
	return internalArray[index]
}

func MArray_GetSlice[T any](array *MemArray[T], start int32, end int32) []T {
	// Convert end (exclusive) to segmentLength
	segmentLength := end - start
	slice, err := CreateSliceFromRange[T](array, start, segmentLength)
	if err != nil {
		panic(err)
	}
	internalArray := slice.InternalArray()
	return internalArray
}

// CreateSliceFromRange simulates the process of creating a non-owning slice
// reference (e.g., ElementConfigs slice) from a larger ClayArray.
// It performs bounds checking and returns the non-owning ClaySlice.
// startOffset: the starting index in the base array
// segmentLength: the length of the slice to create (not an end index)
func CreateSliceFromRange[T any](baseArray *MemArray[T], startOffset int32, segmentLength int32) (MemSlice[T], error) {
	if segmentLength < 0 {
		return MemSlice[T]{}, errors.New("segmentLength cannot be negative")
	}
	if startOffset < 0 {
		return MemSlice[T]{}, errors.New("startOffset cannot be negative")
	}
	// Check against capacity since slices can reference uninitialized capacity
	if startOffset+segmentLength > baseArray.Capacity() {
		return MemSlice[T]{}, fmt.Errorf("slice range exceeds the bounds of the base array: startOffset %d + segmentLength %d > capacity %d", startOffset, segmentLength, baseArray.Capacity())
	}

	// The key implementation: the Go slice expression creates a reference (internalArray)
	// that points to the segment of the base array's memory, achieving the "non-owning reference" pattern.
	// The internalArray field in the ClaySlice now holds a pointer to the start of the segment.
	segmentView := baseArray.InternalArray()[startOffset : startOffset+segmentLength]

	return MemSlice[T]{
		internalArray: segmentView,
	}, nil
}

func (slice MemSlice[T]) Get(index int32) T {
	if !rangeCheck(index, slice.Length()) {
		// message := fmt.Sprintf("MemSlice.Get index: %d, slice.Length: %d\n", index, slice.Length)
		// fmt.Println(message)
		// panic(message)
		zero := new(T)
		return *zero
	}
	internalArray := slice.InternalArray()
	return internalArray[index]
}

func MSlice_Grow[T any](slice *MemSlice[T], length int32) {
	slice.Grow(length)
}

func MSlice_Shrink[T any](slice *MemSlice[T], length int32) {
	slice.Shrink(length)
}
