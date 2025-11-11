package clay

import (
	"bytes"
	"unsafe"
)

func Clay__MemCmp(s1 unsafe.Pointer, s2 unsafe.Pointer, length int32) bool {
	// Convert the struct pointers to []byte slices
	// Note: This is a common but *unsafe* pattern in Go and should be used cautiously.
	s1Bytes := unsafe.Slice((*byte)(s1), length)
	s2Bytes := unsafe.Slice((*byte)(s2), length)

	return bytes.Equal(s1Bytes, s2Bytes)
}

func Clay__MemCmpTyped[T any](s1 *T, s2 *T) bool {
	size := unsafe.Sizeof(new(T))
	return Clay__MemCmp(unsafe.Pointer(s1), unsafe.Pointer(s2), int32(size))
}
