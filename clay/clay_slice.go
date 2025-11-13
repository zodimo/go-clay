package clay

import (
	"github.com/zodimo/clay-go/pkg/mem"
)

type Clay__Slice[T any] = mem.MemSlice[T]

// ClaySlice represents the non-owning reference structure (arrayName##Slice)
// used throughout Clay to point to a sub-range of any backing array.
// type Clay__Slice[T any] struct {
// 	// length (int32_t): The number of elements in the slice.
// 	Length int32

// 	// internalArray (Pointer to typeName): A pointer/slice to the first element
// 	// in the underlying memory block, signifying the non-owning reference.
// 	// This is the view into the specific segment.
// 	InternalArray []T
// }

func NewClay__Slice[T any](data []T) Clay__Slice[T] {
	return mem.NewMemSliceWithData[T](data)
}

func Clay__Slice_Get[T any](slice *Clay__Slice[T], index int32) *T {
	return mem.MSlice_Get(slice, index)
}

func Clay__Slice_Grow[T any](slice *Clay__Slice[T], length int32) {
	mem.MSlice_Grow(slice, length)
}

func Clay__Slice_Shrink[T any](slice *Clay__Slice[T], length int32) {
	mem.MSlice_Shrink(slice, length)
}

// CreateSliceFromRange simulates the process of creating a non-owning slice
// reference (e.g., ElementConfigs slice) from a larger ClayArray.
// It performs bounds checking and returns the non-owning ClaySlice.
func CreateSliceFromRange[T any](baseArray *Clay__Array[T], startOffset int32, segmentLength int32) (Clay__Slice[T], error) {
	return mem.CreateSliceFromRange[T](baseArray, startOffset, segmentLength)
}

// Clay_StringSlice is used to represent non owning string slices, and includes
// a baseChars field which points to the string this slice is derived from.
//
//	typedef struct Clay_StringSlice {
//	    int32_t length;
//	    const char *chars;
//	    const char *baseChars; // The source string / char* that this slice was derived from
//	} Clay_StringSlice;
type Clay_StringSlice struct {
	Length    int32
	Chars     []byte
	BaseChars []byte
}
