package mem

import "unsafe"

func UintptrToPtr[T any](basePtr uintptr, address uintptr) *T {
	offset := address - basePtr
	return (*T)(unsafe.Add(unsafe.Pointer(basePtr), offset))
}

func ArenaPtrToPtr[T any](arena Arena, address uintptr) *T {
	return UintptrToPtr[T](uintptr(unsafe.Pointer(arena.basePtr)), address)
}

// func UintptrToPtrSliceUnsafe[T any](address uintptr) *[]T {
// 	return (*[]T)(unsafe.Pointer(address))
// }
