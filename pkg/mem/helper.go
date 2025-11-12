package mem

import "unsafe"

func UintptrToPtr[T any](basePtr *byte, address uintptr) *T {
	basePtrUintptr := uintptr(unsafe.Pointer(basePtr))
	offset := address - basePtrUintptr
	return (*T)(unsafe.Add(unsafe.Pointer(basePtr), offset))
}
