package mem

import (
	"errors"
	"testing"
)

func TestNewMemSlice(t *testing.T) {
	t.Run("creates slice with specified length", func(t *testing.T) {
		length := int32(5)
		slice := NewMemSlice[int](length)

		if slice.Length != length {
			t.Errorf("expected Length = %d, got %d", length, slice.Length)
		}
		if len(slice.InternalArray) != int(length) {
			t.Errorf("expected InternalArray length = %d, got %d", length, len(slice.InternalArray))
		}
	})

	t.Run("creates slice with zero length", func(t *testing.T) {
		slice := NewMemSlice[int](0)

		if slice.Length != 0 {
			t.Errorf("expected Length = 0, got %d", slice.Length)
		}
		if len(slice.InternalArray) != 0 {
			t.Errorf("expected InternalArray length = 0, got %d", len(slice.InternalArray))
		}
	})

	t.Run("creates slice with different types", func(t *testing.T) {
		intSlice := NewMemSlice[int](3)
		if intSlice.Length != 3 {
			t.Errorf("expected int slice Length = 3, got %d", intSlice.Length)
		}

		stringSlice := NewMemSlice[string](2)
		if stringSlice.Length != 2 {
			t.Errorf("expected string slice Length = 2, got %d", stringSlice.Length)
		}

		type TestStruct struct {
			X int
		}
		structSlice := NewMemSlice[TestStruct](4)
		if structSlice.Length != 4 {
			t.Errorf("expected struct slice Length = 4, got %d", structSlice.Length)
		}
	})
}

func TestMemSlice_Get(t *testing.T) {
	t.Run("gets valid index", func(t *testing.T) {
		slice := NewMemSlice[int](5)
		slice.InternalArray[0] = 42

		ptr := MSlice_Get(&slice, 0)
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *ptr != 42 {
			t.Errorf("expected *ptr = 42, got %d", *ptr)
		}
	})

	t.Run("returns nil for negative index", func(t *testing.T) {
		slice := NewMemSlice[int](5)

		ptr := MSlice_Get(&slice, -1)
		if ptr != nil {
			t.Error("expected nil pointer for negative index")
		}
	})

	t.Run("returns nil for index >= slice length", func(t *testing.T) {
		slice := NewMemSlice[int](3)

		ptr := MSlice_Get(&slice, 3)
		if ptr != nil {
			t.Error("expected nil pointer for index >= length")
		}

		ptr2 := MSlice_Get(&slice, 5)
		if ptr2 != nil {
			t.Error("expected nil pointer for index > length")
		}
	})

	t.Run("gets multiple valid indices", func(t *testing.T) {
		slice := NewMemSlice[int](5)
		MSlice_Set(&slice, 0, 10)
		MSlice_Set(&slice, 1, 20)
		MSlice_Set(&slice, 2, 30)

		ptr0 := MSlice_Get(&slice, 0)
		ptr1 := MSlice_Get(&slice, 1)
		ptr2 := MSlice_Get(&slice, 2)

		if ptr0 == nil || *ptr0 != 10 {
			t.Errorf("expected ptr0 = 10, got %v", ptr0)
		}
		if ptr1 == nil || *ptr1 != 20 {
			t.Errorf("expected ptr1 = 20, got %v", ptr1)
		}
		if ptr2 == nil || *ptr2 != 30 {
			t.Errorf("expected ptr2 = 30, got %v", ptr2)
		}
	})

	t.Run("returns pointer that can modify slice", func(t *testing.T) {
		slice := NewMemSlice[int](5)
		MSlice_Set(&slice, 0, 42)

		ptr := MSlice_Get(&slice, 0)
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}

		*ptr = 100
		if MSlice_GetValue(&slice, 0) != 100 {
			t.Errorf("expected InternalArray[0] = 100 after modification, got %d", slice.InternalArray[0])
		}
	})

	t.Run("gets last valid index", func(t *testing.T) {
		slice := NewMemSlice[int](5)
		MSlice_Set(&slice, 4, 99)

		ptr := MSlice_Get(&slice, 4)
		if ptr == nil || *ptr != 99 {
			t.Errorf("expected ptr = 99, got %v", ptr)
		}
	})
}

func TestCreateSliceFromRange(t *testing.T) {
	t.Run("creates slice from valid range", func(t *testing.T) {
		arr := NewMemArray[int](10)
		MArray_Set(&arr, 0, 10)
		MArray_Set(&arr, 1, 20)
		MArray_Set(&arr, 2, 30)
		MArray_Set(&arr, 3, 40)
		arr.Length = 4

		slice, err := CreateSliceFromRange(&arr, 1, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if slice.Length != 2 {
			t.Errorf("expected slice Length = 2, got %d", slice.Length)
		}
		if len(slice.InternalArray) != 2 {
			t.Errorf("expected slice InternalArray length = 2, got %d", len(slice.InternalArray))
		}
		// Verify the slice points to the correct elements
		if MSlice_GetValue(&slice, 0) != 20 {
			t.Errorf("expected slice[0] = 20, got %d", MSlice_GetValue(&slice, 0))
		}
		if MSlice_GetValue(&slice, 1) != 30 {
			t.Errorf("expected slice[1] = 30, got %d", MSlice_GetValue(&slice, 1))
		}
	})

	t.Run("creates slice from start of array", func(t *testing.T) {
		arr := NewMemArray[int](10)
		MArray_Set(&arr, 0, 10)
		MArray_Set(&arr, 1, 20)
		arr.Length = 2

		slice, err := CreateSliceFromRange(&arr, 0, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if slice.Length != 2 {
			t.Errorf("expected slice Length = 2, got %d", slice.Length)
		}
		if MSlice_GetValue(&slice, 0) != 10 || MSlice_GetValue(&slice, 1) != 20 {
			t.Error("slice does not contain expected values")
		}
	})

	t.Run("creates slice from end of array", func(t *testing.T) {
		arr := NewMemArray[int](10)
		MArray_Set(&arr, 2, 30)
		MArray_Set(&arr, 3, 40)
		arr.Length = 4

		slice, err := CreateSliceFromRange(&arr, 2, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if slice.Length != 2 {
			t.Errorf("expected slice Length = 2, got %d", slice.Length)
		}
		if MSlice_GetValue(&slice, 0) != 30 || MSlice_GetValue(&slice, 1) != 40 {
			t.Error("slice does not contain expected values")
		}
	})

	t.Run("creates single element slice", func(t *testing.T) {
		arr := NewMemArray[int](10)
		MArray_Set(&arr, 1, 42)
		arr.Length = 3

		slice, err := CreateSliceFromRange(&arr, 1, 1)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if slice.Length != 1 {
			t.Errorf("expected slice Length = 1, got %d", slice.Length)
		}
		if MSlice_GetValue(&slice, 0) != 42 {
			t.Errorf("expected slice[0] = 42, got %d", slice.InternalArray[0])
		}
	})

	t.Run("returns error for negative start offset", func(t *testing.T) {
		arr := NewMemArray[int](10)
		arr.Length = 5

		_, err := CreateSliceFromRange(&arr, -1, 2)
		if err == nil {
			t.Fatal("expected error for negative start offset")
		}

		if !errors.Is(err, errors.New("slice range exceeds the bounds of the base array")) {
			// Check error message
			if err.Error() != "slice range exceeds the bounds of the base array" {
				t.Errorf("expected specific error message, got %v", err)
			}
		}
	})

	t.Run("returns error when start offset exceeds array length", func(t *testing.T) {
		arr := NewMemArray[int](10)
		arr.Length = 5

		_, err := CreateSliceFromRange(&arr, 5, 1)
		if err == nil {
			t.Fatal("expected error when start offset >= array length")
		}

	})

	t.Run("returns error when range exceeds array bounds", func(t *testing.T) {
		arr := NewMemArray[int](10)
		arr.Length = 5

		_, err := CreateSliceFromRange(&arr, 3, 3)
		if err == nil {
			t.Fatal("expected error when range exceeds array bounds")
		}

		// startOffset (3) + segmentLength (3) = 6 > arr.Length (5)
	})

	t.Run("returns error when start + length exceeds array length", func(t *testing.T) {
		arr := NewMemArray[int](10)
		arr.Length = 5

		_, err := CreateSliceFromRange(&arr, 0, 6)
		if err == nil {
			t.Fatal("expected error when start + length > array length")
		}

	})

	t.Run("slice modifications reflect in base array", func(t *testing.T) {
		arr := NewMemArray[int](10)
		MArray_Set(&arr, 1, 20)
		MArray_Set(&arr, 2, 30)
		arr.Length = 4

		slice, err := CreateSliceFromRange(&arr, 1, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Modify through slice
		MSlice_Set(&slice, 0, 200)
		MSlice_Set(&slice, 1, 300)

		// Verify base array is modified
		if MArray_GetValue(&arr, 1) != 200 {
			t.Errorf("expected arr[1] = 200, got %d", MArray_GetValue(&arr, 1))
		}
		if MArray_GetValue(&arr, 2) != 300 {
			t.Errorf("expected arr[2] = 300, got %d", MArray_GetValue(&arr, 2))
		}
	})

	t.Run("base array modifications reflect in slice", func(t *testing.T) {
		arr := NewMemArray[int](10)
		MArray_Set(&arr, 1, 20)
		MArray_Set(&arr, 2, 30)
		arr.Length = 4

		slice, err := CreateSliceFromRange(&arr, 1, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Modify base array
		MArray_Set(&arr, 1, 999)
		MArray_Set(&arr, 2, 888)

		// Verify slice reflects changes
		if MSlice_GetValue(&slice, 0) != 999 {
			t.Errorf("expected slice[0] = 999, got %d", MSlice_GetValue(&slice, 0))
		}
		if MSlice_GetValue(&slice, 1) != 888 {
			t.Errorf("expected slice[1] = 888, got %d", MSlice_GetValue(&slice, 1))
		}
	})

	t.Run("creates slice with zero length", func(t *testing.T) {
		arr := NewMemArray[int](10)
		arr.Length = 5

		slice, err := CreateSliceFromRange(&arr, 2, 0)
		if err != nil {
			t.Fatalf("expected no error for zero-length slice, got %v", err)
		}
		if slice.Length != 0 {
			t.Errorf("expected slice Length = 0, got %d", slice.Length)
		}
		if len(slice.InternalArray) != 0 {
			t.Errorf("expected slice InternalArray length = 0, got %d", len(slice.InternalArray))
		}
	})

	t.Run("creates multiple slices from same array", func(t *testing.T) {
		arr := NewMemArray[int](10)
		for i := int32(0); i < 6; i++ {
			MArray_Set(&arr, i, int(i*10))
		}
		arr.Length = 6

		slice1, err1 := CreateSliceFromRange(&arr, 0, 2)
		if err1 != nil {
			t.Fatalf("expected no error creating slice1, got %v", err1)
		}

		slice2, err2 := CreateSliceFromRange(&arr, 2, 2)
		if err2 != nil {
			t.Fatalf("expected no error creating slice2, got %v", err2)
		}

		slice3, err3 := CreateSliceFromRange(&arr, 4, 2)
		if err3 != nil {
			t.Fatalf("expected no error creating slice3, got %v", err3)
		}

		// Verify slices are independent views
		if MSlice_GetValue(&slice1, 0) != 0 || MSlice_GetValue(&slice1, 1) != 10 {
			t.Error("slice1 does not contain expected values")
		}
		if MSlice_GetValue(&slice2, 0) != 20 || MSlice_GetValue(&slice2, 1) != 30 {
			t.Error("slice2 does not contain expected values")
		}
		if MSlice_GetValue(&slice3, 0) != 40 || MSlice_GetValue(&slice3, 1) != 50 {
			t.Error("slice3 does not contain expected values")
		}
	})

	t.Run("creates slice with different types", func(t *testing.T) {
		stringArr := NewMemArray[string](10)
		MArray_Set(&stringArr, 0, "hello")
		MArray_Set(&stringArr, 1, "world")
		stringArr.Length = 2

		slice, err := CreateSliceFromRange(&stringArr, 0, 2)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if slice.Length != 2 {
			t.Errorf("expected slice Length = 2, got %d", slice.Length)
		}
		if MSlice_GetValue(&slice, 0) != "hello" || MSlice_GetValue(&slice, 1) != "world" {
			t.Error("slice does not contain expected string values")
		}
	})
}

func TestMemSlice_Integration(t *testing.T) {
	t.Run("full workflow with slice operations", func(t *testing.T) {
		arr := NewMemArray[int](10)
		for i := int32(0); i < 5; i++ {
			MArray_Set(&arr, i, int(i*10))
		}
		arr.Length = 5

		// Create slice
		slice, err := CreateSliceFromRange(&arr, 1, 3)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Get values from slice
		ptr := MSlice_Get(&slice, 0)
		if ptr == nil || *ptr != 10 {
			t.Error("failed to get value from slice")
		}

		// Modify through slice
		MSlice_Set(&slice, 1, 999)

		// Verify base array is modified
		if MArray_GetValue(&arr, 2) != 999 {
			t.Error("slice modification did not reflect in base array")
		}
	})
}
