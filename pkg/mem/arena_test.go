// Copyright (c) 2024 Nic Barker
// Copyright (c) 2024 go-arena-memory contributors
//
// This software is provided 'as-is', without any express or implied warranty.
// See LICENSE file for full license text.
package mem

import (
	"testing"
	"unsafe"
)

// uintptrToPtr safely converts a uintptr address back to a pointer by computing
// it relative to the original slice's base pointer. This is necessary for
// race detector's checkptr validation, which requires pointers to be derived
// from valid Go objects.
func uintptrToPtr[T any](base []byte, address uintptr) *T {
	basePtr := uintptr(unsafe.Pointer(&base[0]))
	offset := address - basePtr
	return (*T)(unsafe.Add(unsafe.Pointer(&base[0]), offset))
}

func TestNewArena(t *testing.T) {
	t.Run("creates arena with valid memory", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, err := NewArena(memory)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if arena == nil {
			t.Fatal("expected arena to be non-nil")
		}
		if arena.Capacity != 1024 {
			t.Errorf("expected capacity 1024, got %d", arena.Capacity)
		}
	})

	t.Run("creates arena with empty memory", func(t *testing.T) {
		memory := make([]byte, 0)
		arena, err := NewArena(memory)
		if err == nil {
			t.Fatalf("expected error, got nil")
		}
		if arena != nil {
			t.Fatalf("expected arena to be nil, got %v", arena)
		}
	})

}

func TestArena_Allocate(t *testing.T) {
	t.Run("allocates memory successfully", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		address, err := arena.Allocate(100)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if address == 0 {
			t.Errorf("expected address to be non-zero, got %d", address)
		}

		expectedOffset := ((arena.CacheLineSize - ((arena.NextAllocation + 100) % arena.CacheLineSize)) & (arena.CacheLineSize - 1)) + 100

		if arena.NextAllocation != expectedOffset {
			t.Errorf("expected NextAllocation to be %d, got %d", expectedOffset, arena.NextAllocation)
		}
	})

	t.Run("allocates multiple blocks sequentially", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		initialOffset := arena.NextAllocation

		block1, err1 := arena.Allocate(50)
		if err1 != nil {
			t.Fatalf("expected no error on first allocation, got %v", err1)
		}

		block2, err2 := arena.Allocate(100)
		if err2 != nil {
			t.Fatalf("expected no error on second allocation, got %v", err2)
		}

		if block1 == 0 {
			t.Errorf("expected block1 to be non-zero, got %d", block1)
		}
		if block2 == 0 {
			t.Errorf("expected block2 to be non-zero, got %d", block2)
		}

		expectedOffset := initialOffset + 64 + 128

		if arena.NextAllocation != expectedOffset {
			t.Errorf("expected NextAllocation %d, got %d", expectedOffset, arena.NextAllocation)
		}
	})

	t.Run("returns error when capacity exceeded", func(t *testing.T) {
		memory := make([]byte, 100)
		arena, _ := NewArena(memory)

		// Try to allocate more than available
		_, err := arena.Allocate(200)
		if err == nil {
			t.Fatal("expected error when capacity exceeded, got nil")
		}
		if err.Error() != "arena capacity exceeded: cannot allocate required memory" {
			t.Errorf("expected specific error message, got %v", err)
		}
	})

	t.Run("allocates up to capacity limit", func(t *testing.T) {
		memory := make([]byte, 100)
		arena, _ := NewArena(memory)

		initialOffset := arena.NextAllocation
		available := arena.Capacity - initialOffset

		address, err := arena.Allocate(available)
		if err != nil {
			t.Fatalf("expected no error allocating up to capacity, got %v", err)
		}
		if address == 0 {
			t.Errorf("expected address to be non-zero, got %d", address)
		}

		// Next allocation should fail
		_, err2 := arena.Allocate(1)
		if err2 == nil {
			t.Fatal("expected error when allocating beyond capacity, got nil")
		}
	})

	t.Run("allocated memory is writable", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		address, err := arena.Allocate(10)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		buffer := uintptrToPtr[[10]byte](memory, address)
		// Write to allocated memory
		for i := range buffer {
			buffer[i] = byte(i)
		}

		buffer2 := uintptrToPtr[[10]byte](memory, address)

		// Verify the data was written
		for i := range buffer {
			if buffer2[i] != byte(i) {
				t.Errorf("expected buffer2[%d] = %d, got %d", i, i, buffer2[i])
			}
		}
	})
}

func TestArena_AllocateStruct(t *testing.T) {
	type TestStruct struct {
		X int64
		Y int64
	}

	t.Run("allocates struct successfully", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		ptr, err := AllocateStruct[TestStruct](arena)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}

		// Verify we can write to the struct
		ptr.X = 42
		ptr.Y = 100

		if ptr.X != 42 {
			t.Errorf("expected X = 42, got %d", ptr.X)
		}
		if ptr.Y != 100 {
			t.Errorf("expected Y = 100, got %d", ptr.Y)
		}
	})

	t.Run("allocates multiple structs with proper alignment", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		ptr1, err1 := AllocateStruct[TestStruct](arena)
		if err1 != nil {
			t.Fatalf("expected no error on first allocation, got %v", err1)
		}

		ptr2, err2 := AllocateStruct[TestStruct](arena)
		if err2 != nil {
			t.Fatalf("expected no error on second allocation, got %v", err2)
		}

		// Verify structs are distinct
		if ptr1 == ptr2 {
			t.Fatal("expected different pointers for different allocations")
		}

		// Verify we can write to both independently
		ptr1.X = 1
		ptr1.Y = 2
		ptr2.X = 3
		ptr2.Y = 4

		if ptr1.X != 1 || ptr1.Y != 2 {
			t.Errorf("ptr1 values corrupted: X=%d, Y=%d", ptr1.X, ptr1.Y)
		}
		if ptr2.X != 3 || ptr2.Y != 4 {
			t.Errorf("ptr2 values corrupted: X=%d, Y=%d", ptr2.X, ptr2.Y)
		}
	})

	t.Run("structs are properly aligned", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		ptr, err := AllocateStruct[TestStruct](arena)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Check alignment
		address := uintptr(unsafe.Pointer(ptr))
		alignment := unsafe.Alignof(TestStruct{})
		if address%alignment != 0 {
			t.Errorf("struct not properly aligned: address %d, alignment %d", address, alignment)
		}
	})

	t.Run("returns error when capacity exceeded", func(t *testing.T) {
		// Use a size that's large enough to pass NewArena but too small for the struct
		// TestStruct is 16 bytes (2 int64s), plus alignment padding
		memory := make([]byte, 100)
		arena, err := NewArena(memory)
		if err != nil {
			t.Fatalf("expected no error creating arena, got %v", err)
		}

		// Allocate most of the memory to leave just a small amount
		available := arena.Capacity - arena.NextAllocation
		_, err1 := arena.Allocate(available - 5) // Leave only 5 bytes
		if err1 != nil {
			t.Fatalf("expected no error allocating memory, got %v", err1)
		}

		// Now try to allocate struct - should fail
		_, err2 := AllocateStruct[TestStruct](arena)
		if err2 == nil {
			t.Fatal("expected error when capacity exceeded, got nil")
		}
		if err2.Error() != "arena capacity exceeded: cannot allocate required memory" {
			t.Errorf("expected specific error message, got %v", err2)
		}
	})

	t.Run("allocates different struct types", func(t *testing.T) {
		type SmallStruct struct {
			X int8
		}
		type LargeStruct struct {
			X [100]int64
		}

		memory := make([]byte, 2048)
		arena, _ := NewArena(memory)

		small, err1 := AllocateStruct[SmallStruct](arena)
		if err1 != nil {
			t.Fatalf("expected no error allocating SmallStruct, got %v", err1)
		}

		large, err2 := AllocateStruct[LargeStruct](arena)
		if err2 != nil {
			t.Fatalf("expected no error allocating LargeStruct, got %v", err2)
		}

		small.X = 42
		large.X[0] = 100

		if small.X != 42 {
			t.Errorf("expected small.X = 42, got %d", small.X)
		}
		if large.X[0] != 100 {
			t.Errorf("expected large.X[0] = 100, got %d", large.X[0])
		}
	})
}

func TestArena_PersistentEphemeralMemory(t *testing.T) {
	t.Run("initializes persistent memory boundary", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		// Allocate some memory
		_, err := arena.Allocate(100)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Mark persistent boundary
		arena.InitializePersistentMemory()

		if arena.ArenaResetOffset != arena.NextAllocation {
			t.Errorf("expected ArenaResetOffset = NextAllocation (%d), got %d",
				arena.NextAllocation, arena.ArenaResetOffset)
		}
	})

	t.Run("reset ephemeral memory preserves persistent region", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		// Allocate persistent memory
		persistent, err1 := arena.Allocate(100)
		if err1 != nil {
			t.Fatalf("expected no error, got %v", err1)
		}

		persistentBuffer := uintptrToPtr[[100]byte](memory, persistent)

		// Write to persistent memory
		for i := range persistentBuffer {
			persistentBuffer[i] = byte(i)
		}

		// Mark persistent boundary
		arena.InitializePersistentMemory()
		resetOffset := arena.NextAllocation

		// Allocate ephemeral memory
		ephemeral, err2 := arena.Allocate(50)
		if err2 != nil {
			t.Fatalf("expected no error, got %v", err2)
		}

		ephemeralBuffer := uintptrToPtr[[50]byte](memory, ephemeral)
		// Write to ephemeral memory
		for i := range ephemeralBuffer {
			ephemeralBuffer[i] = 0xFF
		}

		// Reset ephemeral memory
		arena.ResetEphemeralMemory()

		// Verify NextAllocation was reset
		if arena.NextAllocation != resetOffset {
			t.Errorf("expected NextAllocation = %d after reset, got %d",
				resetOffset, arena.NextAllocation)
		}

		// Verify persistent memory is still intact
		for i := range persistentBuffer {
			if persistentBuffer[i] != byte(i) {
				t.Errorf("persistent memory corrupted at index %d: expected %d, got %d",
					i, i, persistentBuffer[i])
			}
		}
	})

	t.Run("can allocate after reset", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		// Allocate persistent memory
		_, err1 := arena.Allocate(100)
		if err1 != nil {
			t.Fatalf("expected no error, got %v", err1)
		}

		arena.InitializePersistentMemory()

		// Allocate and reset ephemeral memory multiple times
		for i := 0; i < 5; i++ {
			_, err2 := arena.Allocate(50)
			if err2 != nil {
				t.Fatalf("expected no error on allocation %d, got %v", i, err2)
			}
			arena.ResetEphemeralMemory()
		}

		// Should still be able to allocate
		_, err3 := arena.Allocate(50)
		if err3 != nil {
			t.Fatalf("expected no error after multiple resets, got %v", err3)
		}
	})

	t.Run("reset before initialization sets to zero", func(t *testing.T) {
		memory := make([]byte, 1024)
		arena, _ := NewArena(memory)

		// Allocate some memory
		_, err := arena.Allocate(100)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Reset without initializing persistent memory
		arena.ResetEphemeralMemory()

		// Should reset to 0 (or initial alignment offset)
		if arena.NextAllocation != arena.ArenaResetOffset {
			t.Errorf("expected NextAllocation = ArenaResetOffset (%d), got %d",
				arena.ArenaResetOffset, arena.NextAllocation)
		}
		if arena.ArenaResetOffset != 0 {
			t.Errorf("expected ArenaResetOffset = 0 before initialization, got %d",
				arena.ArenaResetOffset)
		}
	})
}

func TestArena_Integration(t *testing.T) {
	t.Run("complex allocation scenario", func(t *testing.T) {
		memory := make([]byte, 2048)
		arena, err := NewArena(memory)
		if err != nil {
			t.Fatalf("expected no error creating arena, got %v", err)
		}

		// Allocate persistent data
		persistent1, err1 := arena.Allocate(200)
		if err1 != nil {
			t.Fatalf("expected no error, got %v", err1)
		}
		persistent1Buffer := uintptrToPtr[[200]byte](memory, persistent1)
		persistent1Buffer[0] = 0xAA

		persistent2, err2 := arena.Allocate(100)
		if err2 != nil {
			t.Fatalf("expected no error, got %v", err2)
		}
		persistent2Buffer := uintptrToPtr[[100]byte](memory, persistent2)
		persistent2Buffer[0] = 0xBB

		// Mark persistent boundary
		arena.InitializePersistentMemory()

		// Allocate ephemeral data
		ephemeral1, err3 := arena.Allocate(150)
		if err3 != nil {
			t.Fatalf("expected no error, got %v", err3)
		}
		ephemeral1Buffer := uintptrToPtr[[150]byte](memory, ephemeral1)
		ephemeral1Buffer[0] = 0xCC

		// Allocate struct
		type DataStruct struct {
			Value int64
		}
		structPtr, err4 := AllocateStruct[DataStruct](arena)
		if err4 != nil {
			t.Fatalf("expected no error, got %v", err4)
		}
		structPtr.Value = 12345

		// Reset ephemeral
		arena.ResetEphemeralMemory()

		// Verify persistent data intact
		if persistent1Buffer[0] != 0xAA {
			t.Error("persistent1 data corrupted")
		}
		if persistent2Buffer[0] != 0xBB {
			t.Error("persistent2 data corrupted")
		}

		// Allocate new ephemeral data
		ephemeral2, err5 := arena.Allocate(50)
		if err5 != nil {
			t.Fatalf("expected no error after reset, got %v", err5)
		}
		ephemeral2Buffer := uintptrToPtr[[50]byte](memory, ephemeral2)
		ephemeral2Buffer[0] = 0xDD

		// Verify we can still use the arena
		if arena.NextAllocation < arena.ArenaResetOffset {
			t.Error("NextAllocation should be >= ArenaResetOffset")
		}
	})
}

type MyArray[T any] struct {
	Capacity      int
	Length        int
	InternalArray []T
}

func NewMyArray[T any](capacity int) MyArray[T] {
	return MyArray[T]{
		Capacity:      capacity,
		Length:        0,
		InternalArray: make([]T, capacity),
	}
}

func TestAllocateStructObject_MyArray(t *testing.T) {

	t.Run("allocates MyArray with bounded slice", func(t *testing.T) {
		memory := make([]byte, 4096)
		arena, err := NewArena(memory)
		if err != nil {
			t.Fatalf("expected no error creating arena, got %v", err)
		}

		// Create a MyArray with capacity 10
		initialArray := NewMyArray[int](10)
		ptr, err := AllocateStructObject(arena, initialArray)
		if err != nil {
			t.Fatalf("expected no error allocating MyArray, got %v", err)
		}
		if ptr == nil {
			t.Fatal("expected non-nil pointer")
		}

		// Verify initial state
		if ptr.Capacity != 10 {
			t.Errorf("expected Capacity = 10, got %d", ptr.Capacity)
		}
		if ptr.Length != 0 {
			t.Errorf("expected Length = 0, got %d", ptr.Length)
		}
		if len(ptr.InternalArray) != 10 {
			t.Errorf("expected InternalArray length = 10, got %d", len(ptr.InternalArray))
		}
		if cap(ptr.InternalArray) != 10 {
			t.Errorf("expected InternalArray capacity = 10, got %d", cap(ptr.InternalArray))
		}
	})

	t.Run("test bounds - valid access within capacity", func(t *testing.T) {
		memory := make([]byte, 4096)
		arena, _ := NewArena(memory)

		initialArray := NewMyArray[int](5)
		ptr, err := AllocateStructObject(arena, initialArray)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Test valid access: indices 0 to Capacity-1
		for i := 0; i < ptr.Capacity; i++ {
			ptr.InternalArray[i] = i * 10
			ptr.Length++
		}

		// Verify all values were written correctly
		for i := 0; i < ptr.Capacity; i++ {
			if ptr.InternalArray[i] != i*10 {
				t.Errorf("expected InternalArray[%d] = %d, got %d", i, i*10, ptr.InternalArray[i])
			}
		}

		if ptr.Length != ptr.Capacity {
			t.Errorf("expected Length = %d, got %d", ptr.Capacity, ptr.Length)
		}
	})

	t.Run("test bounds - lower boundary (index 0)", func(t *testing.T) {
		memory := make([]byte, 4096)
		arena, _ := NewArena(memory)

		initialArray := NewMyArray[string](3)
		ptr, err := AllocateStructObject(arena, initialArray)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Test lower boundary: index 0
		ptr.InternalArray[0] = "first"
		ptr.Length = 1

		if ptr.InternalArray[0] != "first" {
			t.Errorf("expected InternalArray[0] = 'first', got %q", ptr.InternalArray[0])
		}
		if ptr.Length != 1 {
			t.Errorf("expected Length = 1, got %d", ptr.Length)
		}
	})

	t.Run("test bounds - upper boundary (index Capacity-1)", func(t *testing.T) {
		memory := make([]byte, 4096)
		arena, _ := NewArena(memory)

		capacity := 7
		initialArray := NewMyArray[float64](capacity)
		ptr, err := AllocateStructObject(arena, initialArray)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Test upper boundary: index Capacity-1
		lastIndex := ptr.Capacity - 1
		ptr.InternalArray[lastIndex] = 3.14159
		ptr.Length = ptr.Capacity

		if ptr.InternalArray[lastIndex] != 3.14159 {
			t.Errorf("expected InternalArray[%d] = 3.14159, got %f", lastIndex, ptr.InternalArray[lastIndex])
		}
		if ptr.Length != capacity {
			t.Errorf("expected Length = %d, got %d", capacity, ptr.Length)
		}
	})

	t.Run("test bounds - out of bounds access (index >= Capacity)", func(t *testing.T) {
		memory := make([]byte, 4096)
		arena, _ := NewArena(memory)

		capacity := 5
		initialArray := NewMyArray[int](capacity)
		ptr, err := AllocateStructObject(arena, initialArray)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Test that accessing index >= Capacity would panic
		// Note: Go slices don't prevent this at compile time, but we can test runtime behavior
		defer func() {
			if r := recover(); r == nil {
				// If we reach here, the panic didn't occur, which is expected
				// because Go slices allow access beyond length (up to capacity)
				// But we should verify the slice is bounded by the capacity we set
			}
		}()

		// Try to access beyond the constructed capacity
		// Since we used make([]T, capacity), the slice has length=capacity
		// Accessing beyond this should panic
		shouldPanic := func() {
			_ = ptr.InternalArray[capacity] // This should panic
		}

		panicked := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			shouldPanic()
		}()

		if !panicked {
			t.Error("expected panic when accessing index >= Capacity, but no panic occurred")
		}
	})

	t.Run("test bounds - negative index access", func(t *testing.T) {
		memory := make([]byte, 4096)
		arena, _ := NewArena(memory)

		initialArray := NewMyArray[int](5)
		ptr, err := AllocateStructObject(arena, initialArray)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Test that accessing negative index panics
		// Use a variable to compute negative index at runtime (Go doesn't allow negative literal indices)
		panicked := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			negIndex := -1
			_ = ptr.InternalArray[negIndex] // This should panic
		}()

		if !panicked {
			t.Error("expected panic when accessing negative index, but no panic occurred")
		}
	})

	t.Run("test bounds - multiple MyArray allocations", func(t *testing.T) {
		memory := make([]byte, 8192)
		arena, _ := NewArena(memory)

		// Allocate multiple MyArray instances with different capacities
		array1 := NewMyArray[int](3)
		ptr1, err1 := AllocateStructObject(arena, array1)
		if err1 != nil {
			t.Fatalf("expected no error allocating first array, got %v", err1)
		}

		array2 := NewMyArray[int](5)
		ptr2, err2 := AllocateStructObject(arena, array2)
		if err2 != nil {
			t.Fatalf("expected no error allocating second array, got %v", err2)
		}

		array3 := NewMyArray[int](2)
		ptr3, err3 := AllocateStructObject(arena, array3)
		if err3 != nil {
			t.Fatalf("expected no error allocating third array, got %v", err3)
		}

		// Verify they are distinct
		if ptr1 == ptr2 || ptr1 == ptr3 || ptr2 == ptr3 {
			t.Fatal("expected all pointers to be distinct")
		}

		// Test bounds on each array independently
		ptr1.InternalArray[0] = 100
		ptr1.InternalArray[ptr1.Capacity-1] = 200
		ptr1.Length = ptr1.Capacity

		ptr2.InternalArray[0] = 300
		ptr2.InternalArray[ptr2.Capacity-1] = 400
		ptr2.Length = ptr2.Capacity

		ptr3.InternalArray[0] = 500
		ptr3.InternalArray[ptr3.Capacity-1] = 600
		ptr3.Length = ptr3.Capacity

		// Verify values are independent
		if ptr1.InternalArray[0] != 100 || ptr1.InternalArray[ptr1.Capacity-1] != 200 {
			t.Error("ptr1 values corrupted")
		}
		if ptr2.InternalArray[0] != 300 || ptr2.InternalArray[ptr2.Capacity-1] != 400 {
			t.Error("ptr2 values corrupted")
		}
		if ptr3.InternalArray[0] != 500 || ptr3.InternalArray[ptr3.Capacity-1] != 600 {
			t.Error("ptr3 values corrupted")
		}
	})

	t.Run("test bounds - capacity zero", func(t *testing.T) {
		memory := make([]byte, 4096)
		arena, _ := NewArena(memory)

		initialArray := NewMyArray[int](0)
		ptr, err := AllocateStructObject(arena, initialArray)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if ptr.Capacity != 0 {
			t.Errorf("expected Capacity = 0, got %d", ptr.Capacity)
		}
		if len(ptr.InternalArray) != 0 {
			t.Errorf("expected InternalArray length = 0, got %d", len(ptr.InternalArray))
		}

		// Accessing any index should panic
		panicked := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			_ = ptr.InternalArray[0] // This should panic
		}()

		if !panicked {
			t.Error("expected panic when accessing index 0 on zero-capacity array, but no panic occurred")
		}
	})

	t.Run("test bounds - capacity one", func(t *testing.T) {
		memory := make([]byte, 4096)
		arena, _ := NewArena(memory)

		initialArray := NewMyArray[int](1)
		ptr, err := AllocateStructObject(arena, initialArray)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// Only index 0 should be valid
		ptr.InternalArray[0] = 42
		ptr.Length = 1

		if ptr.InternalArray[0] != 42 {
			t.Errorf("expected InternalArray[0] = 42, got %d", ptr.InternalArray[0])
		}

		// Index 1 should panic
		panicked := false
		func() {
			defer func() {
				if r := recover(); r != nil {
					panicked = true
				}
			}()
			_ = ptr.InternalArray[1] // This should panic
		}()

		if !panicked {
			t.Error("expected panic when accessing index 1 on capacity-1 array, but no panic occurred")
		}
	})

	t.Run("test bounds - generic type parameter", func(t *testing.T) {
		memory := make([]byte, 4096)
		arena, _ := NewArena(memory)

		// Test with different type parameters
		intArray := NewMyArray[int](3)
		ptrInt, err1 := AllocateStructObject(arena, intArray)
		if err1 != nil {
			t.Fatalf("expected no error, got %v", err1)
		}

		stringArray := NewMyArray[string](3)
		ptrString, err2 := AllocateStructObject(arena, stringArray)
		if err2 != nil {
			t.Fatalf("expected no error, got %v", err2)
		}

		// Test bounds on both
		ptrInt.InternalArray[0] = 1
		ptrInt.InternalArray[ptrInt.Capacity-1] = 3
		ptrInt.Length = ptrInt.Capacity

		ptrString.InternalArray[0] = "a"
		ptrString.InternalArray[ptrString.Capacity-1] = "c"
		ptrString.Length = ptrString.Capacity

		if ptrInt.InternalArray[0] != 1 || ptrInt.InternalArray[ptrInt.Capacity-1] != 3 {
			t.Error("int array bounds test failed")
		}
		if ptrString.InternalArray[0] != "a" || ptrString.InternalArray[ptrString.Capacity-1] != "c" {
			t.Error("string array bounds test failed")
		}
	})
}
