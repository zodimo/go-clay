package clay

import (
	"errors"
	"unsafe"
)

// According to the 1999 ISO C standard (C99), size_t is an unsigned integer type of at least 16 bit (see sections 7.17 and 7.18.3).

// typedef struct Clay_Arena {
//     uintptr_t nextAllocation;
//     size_t capacity;
//     char *memory;
// } Clay_Arena;

// Clay_Arena is a memory arena structure that is used by clay to manage its internal allocations.
// Rather than creating it by hand, it's easier to use Clay_CreateArenaWithCapacityAndMemory()

// Clay_Arena represents the Clay_Arena structure for memory management.
// It acts as a bump-pointer allocator over a pre-allocated memory block.
type Clay_Arena struct {
	// NextAllocation (uintptr_t): Tracks the offset in 'Memory' where the next allocation will begin.
	NextAllocation uintptr

	// Capacity (size_t): The total size (in bytes) of the memory block.
	Capacity uintptr

	// Memory (char *): The actual contiguous memory block managed by the arena.
	Memory []byte

	// ArenaResetOffset is the boundary between Persistent and Ephemeral memory.
	// This field is crucial for the O(1) frame reset mechanism described in the training module.
	ArenaResetOffset uintptr
}

// --- 1. Initialization and Setup ---

// NewClayArena initializes the Arena structure with a pre-allocated byte slice.
func NewClayArena(memory []byte) (*Clay_Arena, error) {
	a := &Clay_Arena{
		Memory:           memory,
		Capacity:         uintptr(len(memory)),
		NextAllocation:   0,
		ArenaResetOffset: 0,
	}

	// In a complete Clay implementation (Clay_Initialize), the memory pointer
	// would often be cacheline-aligned (e.g., 64 bytes) before the first allocation
	// (the Clay_Context). We'll simulate that initial alignment padding here.
	const cacheLineSize = 64

	// Calculate current memory start address and needed padding
	if len(a.Memory) > 0 {
		currentPtr := uintptr(unsafe.Pointer(&a.Memory[0]))
		padding := (cacheLineSize - (currentPtr % cacheLineSize)) % cacheLineSize

		// "Allocate" the padding space
		if a.NextAllocation+padding <= a.Capacity {
			a.NextAllocation += padding
		} else {
			// Handle error: Not enough space even for initial alignment
			return nil, errors.New("arena too small for initial alignment")
		}
	}

	return a, nil
}

// --- 2. Allocation (The Bump-Pointer Mechanism) ---

// Allocate attempts to allocate a block of memory of the given size from the arena.
// NOTE: This implementation does not currently handle allocation alignment for the
// *data structure* being placed, which is required for types larger than a byte.
func (a *Clay_Arena) Allocate(size uintptr) []byte {
	// A full implementation would first align NextAllocation for the requested size/type

	// Check if the allocation fits
	if a.NextAllocation+size > a.Capacity {
		// Triggers the CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED scenario
		panic("Arena capacity exceeded: Cannot allocate required memory.")
	}

	start := a.NextAllocation
	end := start + size

	// Bump the pointer for the next allocation
	a.NextAllocation = end

	// Return the slice representing the allocated region
	return a.Memory[start:end]
}

// --- 3. Persistent and Ephemeral Memory Management ---

// InitializePersistentMemory marks the end of the persistent region.
// All future allocations up to this point are considered persistent (retained across frames).
func (a *Clay_Arena) InitializePersistentMemory() {
	a.ArenaResetOffset = a.NextAllocation
}

// ResetEphemeralMemory executes Clay_BeginLayout's reset mechanism.
// It instantly "frees" all transient data by moving the allocation pointer back
// to the boundary, achieving O(1) performance for frame-to-frame reset.
func (a *Clay_Arena) ResetEphemeralMemory() {
	a.NextAllocation = a.ArenaResetOffset

	// In a production system, you might optionally zero out the memory from
	// the reset offset to the current end to clear stale data, though this
	// would trade speed for safety/cleanness.
}

// AllocateStruct allocates space for a single instance of type T from the arena
// and returns a pointer (*T) to that memory location.
// This method relies on Go's 'unsafe' package to type-cast the memory address.
func AllocateStruct[T any](a *Clay_Arena) (*T, error) {
	// 1. Determine the size and alignment requirements for the type T
	var zero T
	size := unsafe.Sizeof(zero)
	alignment := unsafe.Alignof(zero)

	// 2. Calculate necessary padding for alignment
	// Get the starting address of the arena's memory block
	memStartPtr := uintptr(unsafe.Pointer(&a.Memory[0]))

	// Calculate the actual current memory address for the next allocation
	currentAddress := memStartPtr + a.NextAllocation

	// Calculate padding needed to align currentAddress to alignment
	padding := (alignment - (currentAddress % alignment)) % alignment

	// 3. Check capacity including padding
	allocationSize := size + padding
	if a.NextAllocation+allocationSize > a.Capacity {
		return nil, errors.New("arena capacity exceeded: cannot allocate struct")
	}

	// 4. Apply Padding and Bump the pointer
	a.NextAllocation += padding

	// The aligned start address for the struct
	structStartOffset := a.NextAllocation

	// Bump the pointer past the struct size
	a.NextAllocation += size

	// 5. Unsafe Type-Casting (The core C-like functionality)
	// a. Get a pointer to the start of the memory block
	memPointer := unsafe.Pointer(&a.Memory[0])

	// b. Add the offset to get the address of the newly allocated struct's memory
	structAddress := uintptr(memPointer) + structStartOffset

	// c. Type-cast the raw address (uintptr) back into a Go pointer to type T (*T)
	return (*T)(unsafe.Pointer(structAddress)), nil
}
