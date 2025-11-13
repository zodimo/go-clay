package mem

import (
	"fmt"
	"unsafe"
)

type MemArray[T any] struct {
	isHashmap       bool
	ZeroValue       T
	ZeroValuePtr    *T
	internalArray   []T
	hashInitialised bool
}

func (m *MemArray[T]) Length() int32 {
	return int32(len(m.internalArray))
}

func (m *MemArray[T]) initHashMap() {
	if m.hashInitialised {
		return
	}
	m.hashInitialised = true
	for i := int32(0); i < m.Capacity()-2; i++ {
		m.Add(m.ZeroValue)
	}

}
func (m *MemArray[T]) Capacity() int32 {
	return int32(cap(m.internalArray))
}

func (m *MemArray[T]) Add(item T) *T {
	if m.isFull() {
		panic(fmt.Sprintf("MemArray.Add capacity exceeded: %d + 1 > %d", m.Length(), m.Capacity()))
	}
	m.internalArray = m.internalArray[:m.Length()+1]
	m.internalArray[m.Length()-1] = item
	return &m.internalArray[m.Length()-1]
}

func (m *MemArray[T]) Get(index int32) *T {

	if !rangeCheck(index, m.Length()) {
		message := fmt.Sprintf("MemArray.Get index out of bounds: %d, length: %d\n", index, m.Length())
		panic(message)
	}

	internalArray := m.InternalArray()
	return &internalArray[index]
}

func (m *MemArray[T]) GetUnsafe(index int32) *T {

	if !rangeCheck(index, m.Capacity()) {
		message := fmt.Sprintf("MemArray.GetUnsafe index out of bounds: %d, capacity: %d\n", index, m.Capacity())
		panic(message)
	}

	if !rangeCheck(index, m.Length()) {
		return m.ZeroValuePtr
	}

	internalArray := m.InternalArray()
	return &internalArray[index]
}

func (m *MemArray[T]) GetValueUnsafe(index int32) T {

	if !rangeCheck(index, m.Capacity()) {
		message := fmt.Sprintf("MemArray.GetValueUnsafe index out of bounds: %d, capacity: %d\n", index, m.Capacity())
		panic(message)
	}

	if !rangeCheck(index, m.Length()) {
		return m.ZeroValue
	}

	internalArray := m.InternalArray()
	return internalArray[index]
}

func (m *MemArray[T]) Set(index int32, item T) {
	if !rangeCheck(index, m.Length()) {
		message := fmt.Sprintf("MemArray.Set index out of bounds: %d, length: %d\n", index, m.Length())
		panic(message)
	}
	internalArray := m.InternalArray()
	internalArray[index] = item
}

func (m *MemArray[T]) GetValue(index int32) T {
	if !rangeCheck(index, m.Length()) {
		message := fmt.Sprintf("MemArray.GetValue index out of bounds: %d, length: %d\n", index, m.Length())
		panic(message)
	}
	internalArray := m.InternalArray()
	return internalArray[index]
}

func (m *MemArray[T]) RemoveSwapback(index int32) T {
	if !rangeCheck(index, m.Length()) {
		message := fmt.Sprintf("MemArray.RemoveSwapback index out of bounds: %d, length: %d\n", index, m.Length())
		panic(message)
	}

	removed := m.internalArray[index]
	m.internalArray[index] = m.internalArray[m.Length()-1]
	m.internalArray = m.internalArray[:m.Length()-1]

	return removed
}

func (m *MemArray[T]) Reset() {
	if m.isHashmap {
		for i := int32(0); i < m.Capacity()-2; i++ {
			m.Set(i, m.ZeroValue)
		}
	} else {
		m.internalArray = m.internalArray[:0]

	}
}

func (m *MemArray[T]) Shrink(length int32) {
	if m.Length() < length {
		panic(fmt.Sprintf("MemArray.Shrink length is greater than the array length: %d < %d", m.Length(), length))
	}
	m.internalArray = m.internalArray[:m.Length()-length]
}

func (m *MemArray[T]) Grow(length int32) {
	if m.Length()+length > m.Capacity() {
		panic(fmt.Sprintf("MemArray.Grow capacity exceeded: %d + %d > %d", m.Length(), length, m.Capacity()))
	}
	m.internalArray = m.internalArray[:m.Length()+length]
}

func (m *MemArray[T]) isFull() bool {
	return m.Length() == m.Capacity()
}

func (c *MemArray[T]) InternalArray() []T {
	return c.internalArray
}

type MemArrayOptions[T any] struct {
	Arena        *Arena
	IsHashmap    bool
	ZeroValue    T
	ZeroValuePtr *T
}
type MemArrayOption[T any] func(*MemArrayOptions[T])

func MemArrayWithArena[T any](arena *Arena) MemArrayOption[T] {
	return func(o *MemArrayOptions[T]) {
		o.Arena = arena
	}
}

func MemArrayWithIsHashmap[T any]() MemArrayOption[T] {
	return func(o *MemArrayOptions[T]) {
		o.IsHashmap = true
	}
}

func MemArrayWithZeroValue[T any](value T) MemArrayOption[T] {
	return func(o *MemArrayOptions[T]) {
		o.ZeroValue = value
	}
}

func MemArrayWithZeroValuePtr[T any](value *T) MemArrayOption[T] {
	return func(o *MemArrayOptions[T]) {
		o.ZeroValuePtr = value
	}
}
func defaultMemArrayOptions[T any](size int) MemArrayOptions[T] {
	zero := new(T)
	return MemArrayOptions[T]{
		Arena:        NewArenaWithSizeUnsafe(size),
		IsHashmap:    false,
		ZeroValue:    *zero,
		ZeroValuePtr: zero,
	}
}

func NewMemArray[T any](capacity int32, options ...MemArrayOption[T]) MemArray[T] {

	sizeOfArray := int(capacity) * int(unsafe.Sizeof(new(T)))

	opts := defaultMemArrayOptions[T](sizeOfArray)
	for _, option := range options {
		option(&opts)
	}

	// zero := new(T)
	// size := unsafe.Sizeof(*zero)
	// internalArrayAddress, err := opts.Arena.Array_Allocate_Arena(capacity, uint32(size))
	// if err != nil {
	// 	panic(err)
	// }
	// // Convert the address to a pointer to the first element
	// firstElementPtr := UintptrToPtr[T](uintptr(unsafe.Pointer(opts.Arena.basePtr)), internalArrayAddress)
	// // Create a slice from the pointer with the correct length and capacity
	// internalArraySlice := unsafe.Slice(firstElementPtr, capacity)
	// internalArray := &internalArraySlice

	m := MemArray[T]{
		isHashmap:     opts.IsHashmap,
		ZeroValue:     opts.ZeroValue,
		ZeroValuePtr:  opts.ZeroValuePtr,
		internalArray: make([]T, 0, capacity),
	}
	if opts.IsHashmap {
		m.initHashMap()
	}
	return m
}

func rangeCheck(index int32, length int32) bool {
	return index < length && index >= 0
}

// can only get existing values, index < length
func MArray_Get[T any](array *MemArray[T], index int32) *T {
	return array.Get(index)
}

// can only get existing values, index < capacity
func MArray_GetUnsafe[T any](array *MemArray[T], index int32) *T {
	return array.GetUnsafe(index)
}

// can only get existing values, index < length
func MArray_GetValue[T any](array *MemArray[T], index int32) T {
	return array.GetValue(index)
}

// can only add new values up to the capacity
func MArray_Add[T any](array *MemArray[T], item T) *T {
	return array.Add(item)
}

// can only overwrite existing values
func MArray_Set[T any](array *MemArray[T], index int32, item T) {
	array.Set(index, item)
}

// can only remove existing values, index < length
func MArray_RemoveSwapback[T any](array *MemArray[T], index int32) T {
	return array.RemoveSwapback(index)
}

func MArray_Reset[T any](array *MemArray[T]) {
	array.Reset()
}
func MArray_GetAll[T any](array *MemArray[T]) []T {
	return array.internalArray
}

func MArray_Shrink[T any](array *MemArray[T], length int32) {
	array.Shrink(length)
}

func MArray_Grow[T any](array *MemArray[T], length int32) {
	array.Grow(length)
}
