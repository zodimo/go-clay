// Copyright (c) 2024 Nic Barker
// Copyright (c) 2024 go-arena-memory contributors
//
// This software is provided 'as-is', without any express or implied warranty.
// See LICENSE file for full license text.
package mem

import (
	"errors"
	"unsafe"
)

// Arena represents the Arena structure for memory management.
// It acts as a bump-pointer allocator over a pre-allocated memory block.
type Arena struct {
	// NextAllocation (uintptr_t): Tracks the offset in 'Memory' where the next allocation will begin.
	NextAllocation uintptr

	// Capacity (size_t): The total size (in bytes) of the memory block.
	Capacity uintptr

	// Memory (char *): The actual contiguous memory block managed by the arena.
	Memory uintptr

	// basePtr is a pointer to the base of the memory block, used for safe pointer arithmetic
	// with the race detector. This keeps the connection to the original allocation.
	basePtr *byte

	// ArenaResetOffset is the boundary between Persistent and Ephemeral memory.
	// This field is crucial for the O(1) frame reset mechanism described in the training module.
	ArenaResetOffset uintptr

	// CacheLineSize is the size of the cache line to use for alignment.
	CacheLineSize uintptr
}
type ArenaOptions struct {
	CacheLineSize uintptr
}

type ArenaOption func(*ArenaOptions)

func ArenaWithCacheLineSize(size uintptr) ArenaOption {
	return func(o *ArenaOptions) {
		o.CacheLineSize = size
	}
}

func defaultArenaOptions() ArenaOptions {
	return ArenaOptions{
		CacheLineSize: 64,
	}
}

func NewArenaWithSize(size int) (*Arena, error) {
	memory := make([]byte, size)
	return NewArena(memory)
}

func NewArenaWithSizeUnsafe(size int) *Arena {
	memory := make([]byte, size)
	arena, err := NewArena(memory)
	if err != nil {
		panic(err)
	}
	return arena
}

// NewArena initializes the Arena structure with a pre-allocated byte slice.
func NewArena(memory []byte, options ...ArenaOption) (*Arena, error) {
	opts := defaultArenaOptions()
	for _, option := range options {
		option(&opts)
	}

	if len(memory) == 0 {
		return nil, errors.New("memory cannot be empty")
	}

	memStartPtr := uintptr(unsafe.Pointer(&memory[0]))

	a := &Arena{
		Memory:           memStartPtr,
		basePtr:          &memory[0],
		Capacity:         uintptr(len(memory)),
		NextAllocation:   0,
		ArenaResetOffset: 0,
		CacheLineSize:    opts.CacheLineSize,
	}

	return a, nil
}

// Allocate attempts to allocate a block of memory of the given size from the arena.
// It returns the address of the allocated memory and an error if the allocation fails.
func (a *Arena) Allocate(size uintptr) (uintptr, error) {
	nextAllocOffset := a.NextAllocation + ((a.CacheLineSize - ((a.NextAllocation + size) % a.CacheLineSize)) & (a.CacheLineSize - 1)) + size
	if a.NextAllocation+size <= a.Capacity {
		thisAllocationOffset := a.Memory + a.NextAllocation
		a.NextAllocation = nextAllocOffset
		return thisAllocationOffset, nil
	} else {
		return 0, errors.New("arena capacity exceeded: cannot allocate required memory")
	}
}

func (a *Arena) Array_Allocate_Arena(capacity int32, itemSize uint32) (uintptr, error) {
	totalSizeBytes := uintptr(capacity) * uintptr(itemSize)
	return a.Allocate(totalSizeBytes)
}

// InitializePersistentMemory marks the end of the persistent region.
// All future allocations up to this point are considered persistent (retained across frames).
func (a *Arena) InitializePersistentMemory() {
	a.ArenaResetOffset = a.NextAllocation
}

// It instantly "frees" all transient data by moving the allocation pointer back
// to the boundary, achieving O(1) performance for frame-to-frame reset.
func (a *Arena) ResetEphemeralMemory() {
	a.NextAllocation = a.ArenaResetOffset

	// In a production system, you might optionally zero out the memory from
	// the reset offset to the current end to clear stale data, though this
	// would trade speed for safety/cleanness.
}

// AllocateStruct allocates space for a single instance of type T from the arena
// and returns a pointer (*T) to that memory location.
// This method relies on Go's 'unsafe' package to type-cast the memory address.
func AllocateStruct[T any](a *Arena) (*T, error) {
	// 1. Determine the size and alignment requirements for the type T
	var zero T
	return AllocateStructObject(a, zero)

}
func AllocateStructObject[T any](a *Arena, obj T) (*T, error) {
	// 1. Determine the size and alignment requirements for the type T
	size := unsafe.Sizeof(obj)

	// b. Add the offset to get the address of the newly allocated struct's memory
	structAddress, err := a.Allocate(size)
	if err != nil {
		return nil, err
	}

	// c. Convert uintptr address to pointer using unsafe.Add to maintain connection
	// to the original allocation for race detector validation
	basePtr := uintptr(unsafe.Pointer(a.basePtr))
	offset := structAddress - basePtr
	ptr := (*T)(unsafe.Add(unsafe.Pointer(a.basePtr), offset))

	// Copy the obj data into the allocated struct
	*ptr = obj
	return ptr, nil
}
