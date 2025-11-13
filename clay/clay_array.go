package clay

import "github.com/zodimo/clay-go/pkg/mem"

type Clay__Array[T any] = mem.MemArray[T]

// type Clay__Array[T any] struct {
// 	Capacity      int32
// 	Length        int32
// 	InternalArray []T
// }

// func NewClay__Array[T any](capacity int32) Clay__Array[T] {
// 	return Clay__Array[T]{
// 		Capacity:      capacity,
// 		Length:        0,
// 		InternalArray: make([]T, capacity),
// 	}
// }

// func Clay__Array_RangeCheck(index int32, length int32) bool {
// 	return index < length && index >= 0
// }

func Clay__Array_Pop[T any](array *Clay__Array[T]) T {
	return mem.MArray_RemoveSwapback(array, array.Length()-1)
}

func Clay__Array_Get[T any](array *Clay__Array[T], index int32) *T {
	return mem.MArray_Get(array, index)
}

func Clay__Array_GetUnsafe[T any](array *Clay__Array[T], index int32) *T {
	return mem.MArray_GetUnsafe(array, index)
}

func Clay__Array_GetValue[T any](array *Clay__Array[T], index int32) T {
	return mem.MArray_GetValue(array, index)
}

func Clay__Array_Add[T any](array *Clay__Array[T], item T) *T {
	return mem.MArray_Add(array, item)
}

func Clay__Array_Set[T any](array *Clay__Array[T], index int32, item T) {
	mem.MArray_Set(array, index, item)
}

func Clay__Array_RemoveSwapback[T any](array *Clay__Array[T], index int32) T {
	return mem.MArray_RemoveSwapback(array, index)
}

func Clay__Array_Shrink[T any](array *Clay__Array[T], length int32) {
	mem.MArray_Shrink(array, length)
}

func Clay__Array_Grow[T any](array *Clay__Array[T], length int32) {
	mem.MArray_Grow(array, length)
}

func Clay__Array_Reset[T any](array *Clay__Array[T]) {
	mem.MArray_Reset(array)
}

// typeName arrayName##_RemoveSwapback(arrayName *array, int32_t index) {                                          \
// 	if (Clay__Array_RangeCheck(index, array->length)) {                                                         \
// 		array->length--;                                                                                        \
// 		typeName removed = array->internalArray[index];                                                         \
// 		array->internalArray[index] = array->internalArray[array->length];                                      \
// 		return removed;                                                                                         \
// 	}                                                                                                           \
// 	return typeName##_DEFAULT;                                                                                  \
// }

// // The below functions define array bounds checking and convenience functions for a provided type.
// #define CLAY__ARRAY_DEFINE_FUNCTIONS(typeName, arrayName)                                                       \
//                                                                                                                 \
// typedef struct                                                                                                  \
// {                                                                                                               \
//     int32_t length;                                                                                             \
//     typeName *internalArray;                                                                                    \
// } arrayName##Slice;                                                                                             \
//                                                                                                                 \
// typeName typeName##_DEFAULT = CLAY__DEFAULT_STRUCT;                                                             \
//                                                                                                                 \
// arrayName arrayName##_Allocate_Arena(int32_t capacity, Clay_Arena *arena) {                                     \
//     return CLAY__INIT(arrayName){.capacity = capacity, .length = 0,                                             \
//         .internalArray = (typeName *)Clay__Array_Allocate_Arena(capacity, sizeof(typeName), arena)};            \
// }                                                                                                               \
//                                                                                                                 \
// typeName *arrayName##_Get(arrayName *array, int32_t index) {                                                    \
//     return Clay__Array_RangeCheck(index, array->length) ? &array->internalArray[index] : &typeName##_DEFAULT;   \
// }                                                                                                               \
//                                                                                                                 \
// typeName arrayName##_GetValue(arrayName *array, int32_t index) {                                                \
//     return Clay__Array_RangeCheck(index, array->length) ? array->internalArray[index] : typeName##_DEFAULT;     \
// }                                                                                                               \
//                                                                                                                 \
// typeName *arrayName##_Add(arrayName *array, typeName item) {                                                    \
//     if (Clay__Array_AddCapacityCheck(array->length, array->capacity)) {                                         \
//         array->internalArray[array->length++] = item;                                                           \
//         return &array->internalArray[array->length - 1];                                                        \
//     }                                                                                                           \
//     return &typeName##_DEFAULT;                                                                                 \
// }                                                                                                               \
//                                                                                                                 \
// typeName *arrayName##Slice_Get(arrayName##Slice *slice, int32_t index) {                                        \
//     return Clay__Array_RangeCheck(index, slice->length) ? &slice->internalArray[index] : &typeName##_DEFAULT;   \
// }                                                                                                               \
//                                                                                                                 \
//   \
//                                                                                                                 \
// void arrayName##_Set(arrayName *array, int32_t index, typeName value) {                                         \
// 	if (Clay__Array_RangeCheck(index, array->capacity)) {                                                       \
// 		array->internalArray[index] = value;                                                                    \
// 		array->length = index < array->length ? array->length : index + 1;                                      \
// 	}                                                                                                           \
// }                                                                                                               \
