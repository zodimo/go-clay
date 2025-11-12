package mem

import (
	"testing"
)

func TestNewMemArray(t *testing.T) {
	t.Run("creates array with specified capacity", func(t *testing.T) {
		capacity := int32(10)
		arr := NewMemArray[int](capacity)

		if arr.Capacity != capacity {
			t.Errorf("expected Capacity = %d, got %d", capacity, arr.Capacity)
		}
		if arr.Length != 0 {
			t.Errorf("expected Length = 0, got %d", arr.Length)
		}
		if len(*arr.InternalArray) != int(capacity) {
			t.Errorf("expected InternalArray length = %d, got %d", capacity, len(*arr.InternalArray))
		}
	})

	t.Run("creates array with specified capacity and arena", func(t *testing.T) {
		capacity := int32(10)
		arr := NewMemArray[int](capacity, MemArrayWithArena(NewArenaWithSizeUnsafe(1024)))

		if arr.Capacity != capacity {
			t.Errorf("expected Capacity = %d, got %d", capacity, arr.Capacity)
		}
		if arr.Length != 0 {
			t.Errorf("expected Length = 0, got %d", arr.Length)
		}
		if len(*arr.InternalArray) != int(capacity) {
			t.Errorf("expected InternalArray length = %d, got %d", capacity, len(*arr.InternalArray))
		}
	})

	t.Run("creates array with zero capacity", func(t *testing.T) {
		arr := NewMemArray[int](0)

		if arr.Capacity != 0 {
			t.Errorf("expected Capacity = 0, got %d", arr.Capacity)
		}
		if arr.Length != 0 {
			t.Errorf("expected Length = 0, got %d", arr.Length)
		}
		if len(*arr.InternalArray) != 0 {
			t.Errorf("expected InternalArray length = 0, got %d", len(*arr.InternalArray))
		}
	})

	t.Run("creates array with different types", func(t *testing.T) {
		intArr := NewMemArray[int](5)
		if intArr.Capacity != 5 {
			t.Errorf("expected int array Capacity = 5, got %d", intArr.Capacity)
		}

		stringArr := NewMemArray[string](3)
		if stringArr.Capacity != 3 {
			t.Errorf("expected string array Capacity = 3, got %d", stringArr.Capacity)
		}

		type TestStruct struct {
			X int
			Y string
		}
		structArr := NewMemArray[TestStruct](2)
		if structArr.Capacity != 2 {
			t.Errorf("expected struct array Capacity = 2, got %d", structArr.Capacity)
		}
	})
}

func TestMemArray_Get(t *testing.T) {
	t.Run("gets valid index", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 42)
		arr.Length = 1

		ptr := MArray_Get(&arr, 0)
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *ptr != 42 {
			t.Errorf("expected *ptr = 42, got %d", *ptr)
		}
	})

	t.Run("returns nil for negative index", func(t *testing.T) {
		arr := NewMemArray[int](5)
		arr.Length = 3

		ptr := MArray_Get(&arr, -1)
		if ptr != nil {
			t.Error("expected nil pointer for negative index")
		}
	})

	t.Run("returns nil for index >= array capacity", func(t *testing.T) {
		arr := NewMemArray[int](5)
		arr.Length = 3

		// MemArray_Get checks against capacity (len(InternalArray)), not Length
		ptr := MArray_Get(&arr, 3)
		if ptr == nil {
			t.Error("expected non-nil pointer for index 3 (within capacity)")
		}

		ptr2 := MArray_Get(&arr, 5)
		if ptr2 != nil {
			t.Error("expected nil pointer for index >= capacity")
		}
	})

	t.Run("gets multiple valid indices", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		MArray_Set(&arr, 1, 20)
		MArray_Set(&arr, 2, 30)
		arr.Length = 3

		ptr0 := MArray_Get(&arr, 0)
		ptr1 := MArray_Get(&arr, 1)
		ptr2 := MArray_Get(&arr, 2)

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

	t.Run("returns pointer that can modify array", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 42)
		arr.Length = 1

		ptr := MArray_Get(&arr, 0)
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}

		*ptr = 100
		if MArray_GetValue(&arr, 0) != 100 {
			t.Errorf("expected InternalArray[0] = 100 after modification, got %d", (*arr.InternalArray)[0])
		}
	})
}

func TestMArray_GetValue(t *testing.T) {
	t.Run("gets value at valid index", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 42)
		arr.Length = 1

		value := MArray_GetValue(&arr, 0)
		if value != 42 {
			t.Errorf("expected value = 42, got %d", value)
		}
	})

	t.Run("returns zero value for negative index", func(t *testing.T) {
		arr := NewMemArray[int](5)
		arr.Length = 3

		value := MArray_GetValue(&arr, -1)
		if value != 0 {
			t.Errorf("expected zero value for negative index, got %d", value)
		}
	})

	t.Run("returns zero value for index >= array length", func(t *testing.T) {
		arr := NewMemArray[int](5)
		arr.Length = 3

		value := MArray_GetValue(&arr, 3)
		if value != 0 {
			t.Errorf("expected zero value for index >= length, got %d", value)
		}

		value2 := MArray_GetValue(&arr, 10)
		if value2 != 0 {
			t.Errorf("expected zero value for index >= capacity, got %d", value2)
		}
	})

	t.Run("gets zero value for string type", func(t *testing.T) {
		arr := NewMemArray[string](5)
		arr.Length = 3

		value := MArray_GetValue(&arr, 10)
		if value != "" {
			t.Errorf("expected empty string for invalid index, got %q", value)
		}
	})

	t.Run("gets multiple values", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		MArray_Set(&arr, 1, 20)
		MArray_Set(&arr, 2, 30)
		arr.Length = 3

		if MArray_GetValue(&arr, 0) != 10 {
			t.Error("expected value at index 0 = 10")
		}
		if MArray_GetValue(&arr, 1) != 20 {
			t.Error("expected value at index 1 = 20")
		}
		if MArray_GetValue(&arr, 2) != 30 {
			t.Error("expected value at index 2 = 30")
		}
	})
}

func TestMArray_Add(t *testing.T) {
	t.Run("adds item to empty array", func(t *testing.T) {
		arr := NewMemArray[int](5)

		ptr := MArray_Add(&arr, 42)
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *ptr != 42 {
			t.Errorf("expected *ptr = 42, got %d", *ptr)
		}
		if arr.Length != 1 {
			t.Errorf("expected Length = 1, got %d", arr.Length)
		}
		if MArray_GetValue(&arr, 0) != 42 {
			t.Errorf("expected InternalArray[0] = 42, got %d", (*arr.InternalArray)[0])
		}
	})

	t.Run("adds multiple items", func(t *testing.T) {
		arr := NewMemArray[int](5)

		MArray_Add(&arr, 10)
		MArray_Add(&arr, 20)
		MArray_Add(&arr, 30)

		if arr.Length != 3 {
			t.Errorf("expected Length = 3, got %d", arr.Length)
		}
		if MArray_GetValue(&arr, 0) != 10 {
			t.Errorf("expected InternalArray[0] = 10, got %d", MArray_GetValue(&arr, 0))
		}
		if MArray_GetValue(&arr, 1) != 20 {
			t.Errorf("expected InternalArray[1] = 20, got %d", MArray_GetValue(&arr, 1))
		}
		if MArray_GetValue(&arr, 2) != 30 {
			t.Errorf("expected InternalArray[2] = 30, got %d", MArray_GetValue(&arr, 2))
		}
	})

	t.Run("returns nil when capacity-1 reached", func(t *testing.T) {
		arr := NewMemArray[int](3)

		// Add items until capacity-1
		ptr1 := MArray_Add(&arr, 10)
		if ptr1 == nil {
			t.Error("expected non-nil pointer for first add")
		}

		ptr2 := MArray_Add(&arr, 20)
		if ptr2 == nil {
			t.Error("expected non-nil pointer for second add")
		}

		// This should return nil because Length (2) == Capacity-1 (2)
		ptr3 := MArray_Add(&arr, 30)
		if ptr3 != nil {
			t.Error("expected nil pointer when capacity-1 reached")
		}
		if arr.Length != 2 {
			t.Errorf("expected Length to remain 2, got %d", arr.Length)
		}
	})

	t.Run("adds different types", func(t *testing.T) {
		stringArr := NewMemArray[string](3)
		ptr := MArray_Add(&stringArr, "hello")
		if ptr == nil || *ptr != "hello" {
			t.Error("failed to add string")
		}

		type TestStruct struct {
			X int
		}
		structArr := NewMemArray[TestStruct](3)
		ptr2 := MArray_Add(&structArr, TestStruct{X: 42})
		if ptr2 == nil || ptr2.X != 42 {
			t.Error("failed to add struct")
		}
	})
}

func TestMArray_Set(t *testing.T) {
	t.Run("sets value at valid index", func(t *testing.T) {
		arr := NewMemArray[int](5)
		arr.Length = 3

		MArray_Set(&arr, 1, 100)
		if MArray_GetValue(&arr, 1) != 100 {
			t.Errorf("expected InternalArray[1] = 100, got %d", MArray_GetValue(&arr, 1))
		}
	})

	t.Run("sets value at index 0", func(t *testing.T) {
		arr := NewMemArray[int](5)
		arr.Length = 3

		MArray_Set(&arr, 0, 42)
		if MArray_GetValue(&arr, 0) != 42 {
			t.Errorf("expected InternalArray[0] = 42, got %d", MArray_GetValue(&arr, 0))
		}
	})

	t.Run("sets value at last valid index", func(t *testing.T) {
		arr := NewMemArray[int](5)
		arr.Length = 3

		MArray_Set(&arr, 4, 99)
		if MArray_GetValue(&arr, 4) != 99 {
			t.Errorf("expected InternalArray[4] = 99, got %d", MArray_GetValue(&arr, 4))
		}
	})

	t.Run("does not set value for negative index", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		arr.Length = 3

		MArray_Set(&arr, -1, 999)
		if MArray_GetValue(&arr, 0) != 10 {
			t.Error("expected InternalArray[0] to remain unchanged")
		}
	})

	t.Run("does not set value for index >= capacity", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		arr.Length = 3

		MArray_Set(&arr, 5, 999)
		MArray_Set(&arr, 10, 999)
		// Should not panic or modify anything
	})

	t.Run("sets multiple values", func(t *testing.T) {
		arr := NewMemArray[int](5)
		arr.Length = 5

		MArray_Set(&arr, 0, 1)
		MArray_Set(&arr, 1, 2)
		MArray_Set(&arr, 2, 3)

		if MArray_GetValue(&arr, 0) != 1 {
			t.Error("failed to set value at index 0")
		}
		if MArray_GetValue(&arr, 1) != 2 {
			t.Error("failed to set value at index 1")
		}
		if MArray_GetValue(&arr, 2) != 3 {
			t.Error("failed to set value at index 2")
		}
	})

	t.Run("overwrites existing value", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 2, 100)
		arr.Length = 3

		MArray_Set(&arr, 2, 200)
		if MArray_GetValue(&arr, 2) != 200 {
			t.Errorf("expected InternalArray[2] = 200, got %d", MArray_GetValue(&arr, 2))
		}
	})
}

func TestMArray_RemoveSwapback(t *testing.T) {
	t.Run("removes item at valid index", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		MArray_Set(&arr, 1, 20)
		MArray_Set(&arr, 2, 30)
		(*arr.InternalArray)[2] = 30
		arr.Length = 3

		removed := MArray_RemoveSwapback(&arr, 1)
		if removed != 20 {
			t.Errorf("expected removed = 20, got %d", removed)
		}
		if arr.Length != 2 {
			t.Errorf("expected Length = 2, got %d", arr.Length)
		}
		// Last element should be swapped to index 1
		if MArray_GetValue(&arr, 1) != 30 {
			t.Errorf("expected InternalArray[1] = 30 (swapped), got %d", MArray_GetValue(&arr, 1))
		}
	})

	t.Run("removes last item", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		MArray_Set(&arr, 1, 20)
		arr.Length = 2

		removed := MArray_RemoveSwapback(&arr, 1)
		if removed != 20 {
			t.Errorf("expected removed = 20, got %d", removed)
		}
		if arr.Length != 1 {
			t.Errorf("expected Length = 1, got %d", arr.Length)
		}
	})

	t.Run("removes first item", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		MArray_Set(&arr, 1, 20)
		MArray_Set(&arr, 2, 30)
		arr.Length = 3

		removed := MArray_RemoveSwapback(&arr, 0)
		if removed != 10 {
			t.Errorf("expected removed = 10, got %d", removed)
		}
		if arr.Length != 2 {
			t.Errorf("expected Length = 2, got %d", arr.Length)
		}
		// Last element (30) should be swapped to index 0
		if MArray_GetValue(&arr, 0) != 30 {
			t.Errorf("expected InternalArray[0] = 30 (swapped), got %d", MArray_GetValue(&arr, 0))
		}
	})

	t.Run("removes from single element array", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 42)
		arr.Length = 1

		removed := MArray_RemoveSwapback(&arr, 0)
		if removed != 42 {
			t.Errorf("expected removed = 42, got %d", removed)
		}
		if arr.Length != 0 {
			t.Errorf("expected Length = 0, got %d", arr.Length)
		}
	})

	t.Run("returns zero value for negative index", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		arr.Length = 1

		removed := MArray_RemoveSwapback(&arr, -1)
		if removed != 0 {
			t.Errorf("expected zero value for negative index, got %d", removed)
		}
		if arr.Length != 1 {
			t.Errorf("expected Length to remain 1, got %d", arr.Length)
		}
	})

	t.Run("returns zero value for index >= length", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		arr.Length = 1

		removed := MArray_RemoveSwapback(&arr, 1)
		if removed != 0 {
			t.Errorf("expected zero value for invalid index, got %d", removed)
		}
		if arr.Length != 1 {
			t.Errorf("expected Length to remain 1, got %d", arr.Length)
		}
	})

	t.Run("removes multiple items sequentially", func(t *testing.T) {
		arr := NewMemArray[int](5)
		MArray_Set(&arr, 0, 10)
		MArray_Set(&arr, 1, 20)
		MArray_Set(&arr, 2, 30)
		MArray_Set(&arr, 3, 40)
		arr.Length = 4

		// Remove index 1 (20)
		removed1 := MArray_RemoveSwapback(&arr, 1)
		if removed1 != 20 || arr.Length != 3 {
			t.Error("failed to remove first item")
		}

		// Remove index 0 (10, but 40 was swapped to index 1, so we remove what's at 0)
		// Actually, after first removal: [10, 40, 30, ...] with length 3
		// Removing index 0 removes 10, swaps 30 to index 0: [30, 40, ...] with length 2
		removed2 := MArray_RemoveSwapback(&arr, 0)
		if arr.Length != 2 {
			t.Errorf("expected Length = 2 after second removal, got %d", arr.Length)
		}
		_ = removed2
	})

	t.Run("returns zero value for string type", func(t *testing.T) {
		arr := NewMemArray[string](5)
		MArray_Set(&arr, 0, "hello")
		arr.Length = 1

		removed := MArray_RemoveSwapback(&arr, 10)
		if removed != "" {
			t.Errorf("expected empty string for invalid index, got %q", removed)
		}
	})
}

func TestMemArray_Integration(t *testing.T) {
	t.Run("full workflow", func(t *testing.T) {
		arr := NewMemArray[int](10)

		// Add items
		MArray_Add(&arr, 10)
		MArray_Add(&arr, 20)
		MArray_Add(&arr, 30)

		if arr.Length != 3 {
			t.Errorf("expected Length = 3 after adds, got %d", arr.Length)
		}

		// Get values
		if MArray_GetValue(&arr, 0) != 10 {
			t.Error("failed to get value at index 0")
		}

		// Set value
		MArray_Set(&arr, 1, 25)
		if MArray_GetValue(&arr, 1) != 25 {
			t.Error("failed to set value at index 1")
		}

		// Remove item
		removed := MArray_RemoveSwapback(&arr, 0)
		if removed != 10 {
			t.Errorf("expected removed = 10, got %d", removed)
		}
		if arr.Length != 2 {
			t.Errorf("expected Length = 2 after removal, got %d", arr.Length)
		}
	})
}

func TestMemArray_WithArena(t *testing.T) {
	t.Run("creates array with arena and verifies basic operations", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)
		capacity := int32(10)
		arr := NewMemArray[int](capacity, MemArrayWithArena(arena))

		if arr.Capacity != capacity {
			t.Errorf("expected Capacity = %d, got %d", capacity, arr.Capacity)
		}
		if arr.Length != 0 {
			t.Errorf("expected Length = 0, got %d", arr.Length)
		}
		if len(*arr.InternalArray) != int(capacity) {
			t.Errorf("expected InternalArray length = %d, got %d", capacity, len(*arr.InternalArray))
		}
	})

	t.Run("MArray_Get with arena", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)
		arr := NewMemArray[int](5, MemArrayWithArena(arena))
		MArray_Set(&arr, 0, 42)
		arr.Length = 1

		ptr := MArray_Get(&arr, 0)
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *ptr != 42 {
			t.Errorf("expected *ptr = 42, got %d", *ptr)
		}

		// Test multiple gets
		MArray_Set(&arr, 1, 100)
		MArray_Set(&arr, 2, 200)
		arr.Length = 3

		ptr1 := MArray_Get(&arr, 1)
		ptr2 := MArray_Get(&arr, 2)
		if ptr1 == nil || *ptr1 != 100 {
			t.Errorf("expected ptr1 = 100, got %v", ptr1)
		}
		if ptr2 == nil || *ptr2 != 200 {
			t.Errorf("expected ptr2 = 200, got %v", ptr2)
		}

		// Test invalid index
		ptr3 := MArray_Get(&arr, -1)
		if ptr3 != nil {
			t.Error("expected nil pointer for negative index")
		}
	})

	t.Run("MArray_GetValue with arena", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)
		arr := NewMemArray[int](5, MemArrayWithArena(arena))
		MArray_Set(&arr, 0, 42)
		arr.Length = 1

		value := MArray_GetValue(&arr, 0)
		if value != 42 {
			t.Errorf("expected value = 42, got %d", value)
		}

		// Test multiple values
		MArray_Set(&arr, 1, 10)
		MArray_Set(&arr, 2, 20)
		arr.Length = 3

		if MArray_GetValue(&arr, 1) != 10 {
			t.Error("expected value at index 1 = 10")
		}
		if MArray_GetValue(&arr, 2) != 20 {
			t.Error("expected value at index 2 = 20")
		}

		// Test invalid index returns zero value
		zeroValue := MArray_GetValue(&arr, -1)
		if zeroValue != 0 {
			t.Errorf("expected zero value for negative index, got %d", zeroValue)
		}
	})

	t.Run("MArray_Add with arena", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)
		arr := NewMemArray[int](5, MemArrayWithArena(arena))

		ptr := MArray_Add(&arr, 42)
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}
		if *ptr != 42 {
			t.Errorf("expected *ptr = 42, got %d", *ptr)
		}
		if arr.Length != 1 {
			t.Errorf("expected Length = 1, got %d", arr.Length)
		}
		if MArray_GetValue(&arr, 0) != 42 {
			t.Errorf("expected InternalArray[0] = 42, got %d", MArray_GetValue(&arr, 0))
		}

		// Add multiple items
		MArray_Add(&arr, 10)
		MArray_Add(&arr, 20)
		MArray_Add(&arr, 30)

		if arr.Length != 4 {
			t.Errorf("expected Length = 4, got %d", arr.Length)
		}
		if MArray_GetValue(&arr, 1) != 10 {
			t.Errorf("expected InternalArray[1] = 10, got %d", MArray_GetValue(&arr, 1))
		}
		if MArray_GetValue(&arr, 2) != 20 {
			t.Errorf("expected InternalArray[2] = 20, got %d", MArray_GetValue(&arr, 2))
		}
		if MArray_GetValue(&arr, 3) != 30 {
			t.Errorf("expected InternalArray[3] = 30, got %d", MArray_GetValue(&arr, 3))
		}
	})

	t.Run("MArray_Set with arena", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)
		arr := NewMemArray[int](5, MemArrayWithArena(arena))
		arr.Length = 3

		MArray_Set(&arr, 1, 100)
		if MArray_GetValue(&arr, 1) != 100 {
			t.Errorf("expected InternalArray[1] = 100, got %d", MArray_GetValue(&arr, 1))
		}

		// Set multiple values
		MArray_Set(&arr, 0, 1)
		MArray_Set(&arr, 2, 3)

		if MArray_GetValue(&arr, 0) != 1 {
			t.Error("failed to set value at index 0")
		}
		if MArray_GetValue(&arr, 2) != 3 {
			t.Error("failed to set value at index 2")
		}

		// Test overwrite
		MArray_Set(&arr, 1, 200)
		if MArray_GetValue(&arr, 1) != 200 {
			t.Errorf("expected InternalArray[1] = 200 after overwrite, got %d", MArray_GetValue(&arr, 1))
		}
	})

	t.Run("MArray_RemoveSwapback with arena", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)
		arr := NewMemArray[int](5, MemArrayWithArena(arena))
		MArray_Set(&arr, 0, 10)
		MArray_Set(&arr, 1, 20)
		MArray_Set(&arr, 2, 30)
		arr.Length = 3

		removed := MArray_RemoveSwapback(&arr, 1)
		if removed != 20 {
			t.Errorf("expected removed = 20, got %d", removed)
		}
		if arr.Length != 2 {
			t.Errorf("expected Length = 2, got %d", arr.Length)
		}
		// Last element should be swapped to index 1
		if MArray_GetValue(&arr, 1) != 30 {
			t.Errorf("expected InternalArray[1] = 30 (swapped), got %d", MArray_GetValue(&arr, 1))
		}

		// Remove first item
		removed2 := MArray_RemoveSwapback(&arr, 0)
		if removed2 != 10 {
			t.Errorf("expected removed = 10, got %d", removed2)
		}
		if arr.Length != 1 {
			t.Errorf("expected Length = 1, got %d", arr.Length)
		}
		// Last element (30) should be swapped to index 0
		if MArray_GetValue(&arr, 0) != 30 {
			t.Errorf("expected InternalArray[0] = 30 (swapped), got %d", MArray_GetValue(&arr, 0))
		}
	})

	t.Run("full workflow with arena", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)
		arr := NewMemArray[int](10, MemArrayWithArena(arena))

		// Add items
		MArray_Add(&arr, 10)
		MArray_Add(&arr, 20)
		MArray_Add(&arr, 30)

		if arr.Length != 3 {
			t.Errorf("expected Length = 3 after adds, got %d", arr.Length)
		}

		// Get values
		if MArray_GetValue(&arr, 0) != 10 {
			t.Error("failed to get value at index 0")
		}

		// Set value
		MArray_Set(&arr, 1, 25)
		if MArray_GetValue(&arr, 1) != 25 {
			t.Error("failed to set value at index 1")
		}

		// Get pointer and modify
		ptr := MArray_Get(&arr, 2)
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}
		*ptr = 35
		if MArray_GetValue(&arr, 2) != 35 {
			t.Error("failed to modify value via pointer")
		}

		// Remove item
		removed := MArray_RemoveSwapback(&arr, 0)
		if removed != 10 {
			t.Errorf("expected removed = 10, got %d", removed)
		}
		if arr.Length != 2 {
			t.Errorf("expected Length = 2 after removal, got %d", arr.Length)
		}
	})

	t.Run("different types with arena", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)

		// Test string array
		stringArr := NewMemArray[string](3, MemArrayWithArena(arena))
		ptr := MArray_Add(&stringArr, "hello")
		if ptr == nil || *ptr != "hello" {
			t.Error("failed to add string to arena array")
		}

		// Test struct array
		type TestStruct struct {
			X int
			Y string
		}
		structArr := NewMemArray[TestStruct](3, MemArrayWithArena(arena))
		ptr2 := MArray_Add(&structArr, TestStruct{X: 42, Y: "test"})
		if ptr2 == nil || ptr2.X != 42 || ptr2.Y != "test" {
			t.Error("failed to add struct to arena array")
		}
	})

	t.Run("capacity limit with arena", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)
		arr := NewMemArray[int](3, MemArrayWithArena(arena))

		// Add items until capacity-1
		ptr1 := MArray_Add(&arr, 10)
		if ptr1 == nil {
			t.Error("expected non-nil pointer for first add")
		}

		ptr2 := MArray_Add(&arr, 20)
		if ptr2 == nil {
			t.Error("expected non-nil pointer for second add")
		}

		// This should return nil because Length (2) == Capacity-1 (2)
		ptr3 := MArray_Add(&arr, 30)
		if ptr3 != nil {
			t.Error("expected nil pointer when capacity-1 reached")
		}
		if arr.Length != 2 {
			t.Errorf("expected Length to remain 2, got %d", arr.Length)
		}
	})

	t.Run("boundary conditions with arena", func(t *testing.T) {
		arena := NewArenaWithSizeUnsafe(2048)
		arr := NewMemArray[int](5, MemArrayWithArena(arena))
		arr.Length = 3

		// Test setting at boundaries
		MArray_Set(&arr, 0, 1)
		MArray_Set(&arr, 4, 5)

		if MArray_GetValue(&arr, 0) != 1 {
			t.Error("failed to set value at index 0")
		}
		if MArray_GetValue(&arr, 4) != 5 {
			t.Error("failed to set value at last index")
		}

		// Test invalid indices
		MArray_Set(&arr, -1, 999)
		MArray_Set(&arr, 5, 999)
		// Should not modify anything
		if MArray_GetValue(&arr, 0) != 1 {
			t.Error("value at index 0 should remain unchanged")
		}
	})
}
