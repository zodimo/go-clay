package mem

import (
	"errors"
	"fmt"
	"unsafe"
)

// ClaySlice represents the non-owning reference structure (arrayName##Slice)
// used throughout Clay to point to a sub-range of any backing array.
type MemSlice[T any] struct {
	// length (int32_t): The number of elements in the slice.
	length int32

	// internalArray (Pointer to typeName): A pointer/slice to the first element
	// in the underlying memory block, signifying the non-owning reference.
	// This is the view into the specific segment.
	InternalAddress uintptr
	BaseAddress     uintptr
}

func (c *MemSlice[T]) InternalArray() *[]T {
	ptr := UintptrToPtr[[]T](c.BaseAddress, c.InternalAddress)
	return ptr
}

func (s *MemSlice[T]) Grow(length int32) {
	s.length += length
}

func (s *MemSlice[T]) Shrink(length int32) {
	if s.length < length {
		s.length = 0
		return
	}
	s.length -= length
}
func (s *MemSlice[T]) Length() int32 {
	return s.length
}

func NewMemSlice[T any](length int32) MemSlice[T] {
	initialSlice := make([]T, length)
	return MemSlice[T]{
		length:          length,
		InternalAddress: uintptr(unsafe.Pointer(&initialSlice[0])),
		BaseAddress:     uintptr(unsafe.Pointer(&initialSlice[0])),
	}
}

func NewMemSliceWithData[T any](array []T) MemSlice[T] {
	initialSlice := make([]T, len(array))
	copy(initialSlice, array)
	return MemSlice[T]{
		length:          int32(len(array)),
		InternalAddress: uintptr(unsafe.Pointer(&initialSlice[0])),
		BaseAddress:     uintptr(unsafe.Pointer(&initialSlice[0])),
	}
}

func MSlice_Set[T any](slice *MemSlice[T], index int32, item T) {
	if !rangeCheck(index, slice.length) {
		message := fmt.Sprintf("MemSlice.MSlice_Set index: %d, slice.length: %d\n", index, slice.length)
		// fmt.Println(message)
		panic(message)
		// return
	}
	internalArray := slice.InternalArray()
	(*internalArray)[index] = item
}

func MSlice_Get[T any](slice *MemSlice[T], index int32) *T {
	if !rangeCheck(index, slice.length) {
		message := fmt.Sprintf("MemSlice.MSlice_Get index: %d, slice.length: %d\n", index, slice.length)
		// fmt.Println(message)
		panic(message)
		// return nil
	}
	internalArray := slice.InternalArray()
	return &(*internalArray)[index]
}

func MSlice_GetValue[T any](slice *MemSlice[T], index int32) T {
	if !rangeCheck(index, slice.length) {
		message := fmt.Sprintf("MemSlice.MSlice_GetValue index: %d, slice.length: %d\n", index, slice.length)
		// fmt.Println(message)
		panic(message)
		// zero := new(T)
		// return *zero
	}
	internalArray := slice.InternalArray()
	return (*internalArray)[index]
}

func MArray_GetSlice[T any](array *MemArray[T], start int32, end int32) []T {
	slice, err := CreateSliceFromRange[T](array, start, end)
	if err != nil {
		panic(err)
	}
	internalArray := slice.InternalArray()
	return *internalArray
}

type SafeMemoryPointer[T any] struct {
	BaseAddress     uintptr
	InternalAddress uintptr
}

func MArray_GetIndexMemory[T any](array *MemArray[T], index int32) SafeMemoryPointer[[]T] {
	return SafeMemoryPointer[[]T]{
		BaseAddress:     array.BaseAddress,
		InternalAddress: array.InternalAddress + uintptr(index)*unsafe.Sizeof(new(T)),
	}
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
	segmentView := baseArray.UnWrap()[startOffset : startOffset+segmentLength]

	return MemSlice[T]{
		length:          segmentLength,
		InternalAddress: uintptr(unsafe.Pointer(&segmentView[0])),
	}, nil
}

func (slice MemSlice[T]) Get(index int32) T {
	if !rangeCheck(index, slice.length) {
		// message := fmt.Sprintf("MemSlice.Get index: %d, slice.Length: %d\n", index, slice.Length)
		// fmt.Println(message)
		// panic(message)
		zero := new(T)
		return *zero
	}
	internalArray := slice.InternalArray()
	return (*internalArray)[index]
}

func MSlice_Grow[T any](slice *MemSlice[T], length int32) {
	slice.Grow(length)
}

func MSlice_Shrink[T any](slice *MemSlice[T], length int32) {
	slice.Shrink(length)
}
