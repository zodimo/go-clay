package clay

import "errors"

// ClaySlice represents the non-owning reference structure (arrayName##Slice)
// used throughout Clay to point to a sub-range of any backing array.
type Clay__Slice[T any] struct {
	// length (int32_t): The number of elements in the slice.
	Length int32

	// internalArray (Pointer to typeName): A pointer/slice to the first element
	// in the underlying memory block, signifying the non-owning reference.
	// This is the view into the specific segment.
	InternalArray []T
}

func NewClay__Slice[T any](length int32) Clay__Slice[T] {
	return Clay__Slice[T]{
		Length:        length,
		InternalArray: make([]T, length),
	}
}

func Clay__Slice_Get[T any](slice *Clay__Slice[T], index int32) *T {
	if !Clay__Array_RangeCheck(index, slice.Length) {
		return nil
	}
	return &slice.InternalArray[index]
}

// CreateSliceFromRange simulates the process of creating a non-owning slice
// reference (e.g., ElementConfigs slice) from a larger ClayArray.
// It performs bounds checking and returns the non-owning ClaySlice.
func CreateSliceFromRange[T any](baseArray *Clay__Array[T], startOffset int32, segmentLength int32) (*Clay__Slice[T], error) {
	if startOffset < 0 || startOffset+segmentLength > baseArray.Length {
		return nil, errors.New("slice range exceeds the bounds of the base array")
	}

	// The key implementation: the Go slice expression creates a reference (internalArray)
	// that points to the segment of the base array's memory, achieving the "non-owning reference" pattern.
	// The internalArray field in the ClaySlice now holds a pointer to the start of the segment.
	segmentView := baseArray.InternalArray[startOffset : startOffset+segmentLength]

	return &Clay__Slice[T]{
		Length:        segmentLength,
		InternalArray: segmentView,
	}, nil
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
	Chars     string
	BaseChars string
}
