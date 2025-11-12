package mem

import (
	"bytes"
	"unsafe"
)

func Compare(s1 unsafe.Pointer, s2 unsafe.Pointer, length int32) bool {
	// Convert the struct pointers to []byte slices
	// Note: This is a common but *unsafe* pattern in Go and should be used cautiously.
	s1Bytes := unsafe.Slice((*byte)(s1), length)
	s2Bytes := unsafe.Slice((*byte)(s2), length)

	return bytes.Equal(s1Bytes, s2Bytes)
}

func CompareTyped[T any](s1 *T, s2 *T) bool {
	var zero T
	size := unsafe.Sizeof(zero)
	return Compare(unsafe.Pointer(s1), unsafe.Pointer(s2), int32(size))
}
