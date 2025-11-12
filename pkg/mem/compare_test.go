package mem

import (
	"testing"
	"unsafe"
)

func TestCompare(t *testing.T) {
	t.Run("compares equal byte arrays", func(t *testing.T) {
		arr1 := [5]byte{1, 2, 3, 4, 5}
		arr2 := [5]byte{1, 2, 3, 4, 5}

		result := Compare(unsafe.Pointer(&arr1[0]), unsafe.Pointer(&arr2[0]), 5)
		if !result {
			t.Error("expected Compare to return true for equal arrays")
		}
	})

	t.Run("compares different byte arrays", func(t *testing.T) {
		arr1 := [5]byte{1, 2, 3, 4, 5}
		arr2 := [5]byte{1, 2, 3, 4, 6}

		result := Compare(unsafe.Pointer(&arr1[0]), unsafe.Pointer(&arr2[0]), 5)
		if result {
			t.Error("expected Compare to return false for different arrays")
		}
	})

	t.Run("compares arrays with different first byte", func(t *testing.T) {
		arr1 := [5]byte{1, 2, 3, 4, 5}
		arr2 := [5]byte{0, 2, 3, 4, 5}

		result := Compare(unsafe.Pointer(&arr1[0]), unsafe.Pointer(&arr2[0]), 5)
		if result {
			t.Error("expected Compare to return false for arrays with different first byte")
		}
	})

	t.Run("compares arrays with different last byte", func(t *testing.T) {
		arr1 := [5]byte{1, 2, 3, 4, 5}
		arr2 := [5]byte{1, 2, 3, 4, 0}

		result := Compare(unsafe.Pointer(&arr1[0]), unsafe.Pointer(&arr2[0]), 5)
		if result {
			t.Error("expected Compare to return false for arrays with different last byte")
		}
	})

	t.Run("compares zero-length arrays", func(t *testing.T) {
		arr1 := [0]byte{}
		arr2 := [0]byte{}

		result := Compare(unsafe.Pointer(&arr1), unsafe.Pointer(&arr2), 0)
		if !result {
			t.Error("expected Compare to return true for zero-length arrays")
		}
	})

	t.Run("compares single byte", func(t *testing.T) {
		var b1 byte = 42
		var b2 byte = 42
		var b3 byte = 43

		result1 := Compare(unsafe.Pointer(&b1), unsafe.Pointer(&b2), 1)
		if !result1 {
			t.Error("expected Compare to return true for equal single bytes")
		}

		result2 := Compare(unsafe.Pointer(&b1), unsafe.Pointer(&b3), 1)
		if result2 {
			t.Error("expected Compare to return false for different single bytes")
		}
	})

	t.Run("compares partial arrays", func(t *testing.T) {
		arr1 := [5]byte{1, 2, 3, 4, 5}
		arr2 := [5]byte{1, 2, 3, 9, 9}

		// Compare first 3 bytes
		result := Compare(unsafe.Pointer(&arr1[0]), unsafe.Pointer(&arr2[0]), 3)
		if !result {
			t.Error("expected Compare to return true for equal partial arrays")
		}
	})
}

func TestCompareTyped(t *testing.T) {
	t.Run("compares equal integers", func(t *testing.T) {
		var i1 int32 = 42
		var i2 int32 = 42

		result := CompareTyped(&i1, &i2)
		if !result {
			t.Error("expected CompareTyped to return true for equal integers")
		}
	})

	t.Run("compares different integers", func(t *testing.T) {
		var i1 int32 = 42
		var i2 int32 = 43

		result := CompareTyped(&i1, &i2)
		if result {
			t.Error("expected CompareTyped to return false for different integers")
		}
	})

	t.Run("compares equal structs", func(t *testing.T) {
		type TestStruct struct {
			X int32
			Y int32
		}

		s1 := TestStruct{X: 10, Y: 20}
		s2 := TestStruct{X: 10, Y: 20}

		result := CompareTyped(&s1, &s2)
		if !result {
			t.Error("expected CompareTyped to return true for equal structs")
		}
	})

	t.Run("compares different structs", func(t *testing.T) {
		type TestStruct struct {
			X int32
			Y int32
		}

		s1 := TestStruct{X: 10, Y: 20}
		s2 := TestStruct{X: 10, Y: 21}

		result := CompareTyped(&s1, &s2)
		if result {
			t.Error("expected CompareTyped to return false for different structs")
		}
	})

	t.Run("compares structs with different first field", func(t *testing.T) {
		type TestStruct struct {
			X int32
			Y int32
		}

		s1 := TestStruct{X: 10, Y: 20}
		s2 := TestStruct{X: 11, Y: 20}

		result := CompareTyped(&s1, &s2)
		if result {
			t.Error("expected CompareTyped to return false for structs with different first field")
		}
	})

	t.Run("compares equal strings", func(t *testing.T) {
		s1 := "hello"
		s2 := "hello"

		result := CompareTyped(&s1, &s2)
		if !result {
			t.Error("expected CompareTyped to return true for equal strings")
		}
	})

	t.Run("compares different strings", func(t *testing.T) {
		s1 := "hello"
		s2 := "world"

		result := CompareTyped(&s1, &s2)
		if result {
			t.Error("expected CompareTyped to return false for different strings")
		}
	})

	t.Run("compares equal arrays", func(t *testing.T) {
		arr1 := [3]int32{1, 2, 3}
		arr2 := [3]int32{1, 2, 3}

		result := CompareTyped(&arr1, &arr2)
		if !result {
			t.Error("expected CompareTyped to return true for equal arrays")
		}
	})

	t.Run("compares different arrays", func(t *testing.T) {
		arr1 := [3]int32{1, 2, 3}
		arr2 := [3]int32{1, 2, 4}

		result := CompareTyped(&arr1, &arr2)
		if result {
			t.Error("expected CompareTyped to return false for different arrays")
		}
	})

	t.Run("compares pointers to same value", func(t *testing.T) {
		var i int32 = 42
		p1 := &i
		p2 := &i

		result := CompareTyped(p1, p2)
		if !result {
			t.Error("expected CompareTyped to return true for pointers to same value")
		}
	})

	t.Run("compares zero values", func(t *testing.T) {
		type TestStruct struct {
			X int32
			Y int32
		}

		s1 := TestStruct{}
		s2 := TestStruct{}

		result := CompareTyped(&s1, &s2)
		if !result {
			t.Error("expected CompareTyped to return true for zero-value structs")
		}
	})

	t.Run("compares complex structs", func(t *testing.T) {
		type ComplexStruct struct {
			A int64
			B float64
			C [4]byte
			D bool
		}

		s1 := ComplexStruct{
			A: 100,
			B: 3.14,
			C: [4]byte{1, 2, 3, 4},
			D: true,
		}
		s2 := ComplexStruct{
			A: 100,
			B: 3.14,
			C: [4]byte{1, 2, 3, 4},
			D: true,
		}

		result := CompareTyped(&s1, &s2)
		if !result {
			t.Error("expected CompareTyped to return true for equal complex structs")
		}

		s3 := ComplexStruct{
			A: 100,
			B: 3.14,
			C: [4]byte{1, 2, 3, 5},
			D: true,
		}

		result2 := CompareTyped(&s1, &s3)
		if result2 {
			t.Error("expected CompareTyped to return false for different complex structs")
		}
	})
}

