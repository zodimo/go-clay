package mem

import "errors"

// ClaySlice represents the non-owning reference structure (arrayName##Slice)
// used throughout Clay to point to a sub-range of any backing array.
type MemSlice[T any] struct {
	// length (int32_t): The number of elements in the slice.
	Length int32

	// internalArray (Pointer to typeName): A pointer/slice to the first element
	// in the underlying memory block, signifying the non-owning reference.
	// This is the view into the specific segment.
	InternalArray []T
}

func NewMemSlice[T any](length int32) MemSlice[T] {
	return MemSlice[T]{
		Length:        length,
		InternalArray: make([]T, length),
	}
}

func MSlice_Set[T any](slice *MemSlice[T], index int32, item T) {
	if !rangeCheck(index, slice.Length) {
		return
	}
	slice.InternalArray[index] = item
}

func MSlice_Get[T any](slice *MemSlice[T], index int32) *T {
	if !rangeCheck(index, slice.Length) {
		return nil
	}
	return &slice.InternalArray[index]
}

func MSlice_GetValue[T any](slice *MemSlice[T], index int32) T {
	if !rangeCheck(index, slice.Length) {
		zero := new(T)
		return *zero
	}
	return slice.InternalArray[index]
}

func MArray_GetSlice[T any](array *MemArray[T], start int32, end int32) []T {
	slice, err := CreateSliceFromRange[T](array, start, end)
	if err != nil {
		panic(err)
	}
	return slice.InternalArray
}

// CreateSliceFromRange simulates the process of creating a non-owning slice
// reference (e.g., ElementConfigs slice) from a larger ClayArray.
// It performs bounds checking and returns the non-owning ClaySlice.
func CreateSliceFromRange[T any](baseArray *MemArray[T], startOffset int32, segmentLength int32) (MemSlice[T], error) {
	if startOffset < 0 || startOffset+segmentLength > baseArray.Capacity {
		return MemSlice[T]{}, errors.New("slice range exceeds the bounds of the base array")
	}

	// The key implementation: the Go slice expression creates a reference (internalArray)
	// that points to the segment of the base array's memory, achieving the "non-owning reference" pattern.
	// The internalArray field in the ClaySlice now holds a pointer to the start of the segment.
	segmentView := (*baseArray.InternalArray)[startOffset : startOffset+segmentLength]

	return MemSlice[T]{
		Length:        segmentLength,
		InternalArray: segmentView,
	}, nil
}
