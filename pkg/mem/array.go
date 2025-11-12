package mem

import "unsafe"

type MemArray[T any] struct {
	Capacity      int32
	Length        int32
	InternalArray *[]T
}

type MemArrayOptions struct {
	Arena *Arena
}
type MemArrayOption func(*MemArrayOptions)

func MemArrayWithArena(arena *Arena) MemArrayOption {
	return func(o *MemArrayOptions) {
		o.Arena = arena
	}
}
func NewMemArray[T any](capacity int32, options ...MemArrayOption) MemArray[T] {
	opts := MemArrayOptions{
		Arena: nil,
	}
	for _, option := range options {
		option(&opts)
	}

	if opts.Arena != nil {
		zero := new(T)
		size := unsafe.Sizeof(*zero)
		internalArrayAddress, err := opts.Arena.Array_Allocate_Arena(capacity, uint32(size))
		if err != nil {
			panic(err)
		}
		// Convert the address to a pointer to the first element
		firstElementPtr := UintptrToPtr[T](opts.Arena.basePtr, internalArrayAddress)
		// Create a slice from the pointer with the correct length and capacity
		internalArraySlice := unsafe.Slice(firstElementPtr, capacity)
		internalArray := &internalArraySlice

		return MemArray[T]{
			Capacity:      capacity,
			Length:        0,
			InternalArray: internalArray,
		}
	} else {
		internalArray := make([]T, capacity)

		return MemArray[T]{
			Capacity:      capacity,
			Length:        0,
			InternalArray: &internalArray,
		}
	}

}

func rangeCheck(index int32, length int32) bool {
	return index < length && index >= 0
}

func MArray_Get[T any](array *MemArray[T], index int32) *T {
	if !rangeCheck(index, int32(len(*array.InternalArray))) {
		return nil
	}
	return &(*array.InternalArray)[index]
}
func MArray_GetValue[T any](array *MemArray[T], index int32) T {

	if !rangeCheck(index, int32(len(*array.InternalArray))) {
		zero := new(T)
		return *zero
	}
	return (*array.InternalArray)[index]
}

func MArray_Add[T any](array *MemArray[T], item T) *T {
	if array.Length == array.Capacity-1 {
		return nil
	}
	(*array.InternalArray)[array.Length] = item
	array.Length++
	return &(*array.InternalArray)[array.Length-1]
}

func MArray_Set[T any](array *MemArray[T], index int32, item T) {
	if index < 0 || index >= int32(len(*array.InternalArray)) {
		return
	}
	(*array.InternalArray)[index] = item
}

func MArray_RemoveSwapback[T any](array *MemArray[T], index int32) T {
	zero := new(T)
	if !rangeCheck(index, array.Length) {
		return *zero
	}
	array.Length--
	removed := (*array.InternalArray)[index]
	(*array.InternalArray)[index] = (*array.InternalArray)[array.Length]
	return removed
}
