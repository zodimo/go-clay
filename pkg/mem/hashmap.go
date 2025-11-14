package mem

import "errors"

type HashElementId struct {
	Id       uint32
	Offset   uint32
	BaseId   uint32
	StringId string
}

type HashMapContext[T any] struct {
	HashMapInternal MemArray[HashMapItem[T]]
	HashMap         MemArray[int32]
	Generation      uint32
}

func HashString(key string, seed uint32) HashElementId {
	hash := seed

	charsBytes := []byte(key)

	for _, charByte := range charsBytes {
		hash += uint32(charByte)
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}

	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return HashElementId{
		Id:       hash + 1,
		Offset:   0,
		BaseId:   hash + 1,
		StringId: key,
	} // Reserve the hash result of zero as "null id"
}

func HashNumber(offset uint32, seed uint32) HashElementId {
	hash := seed
	hash += (offset + 48)
	hash += (hash << 10)
	hash ^= (hash >> 6)
	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return HashElementId{Id: hash + 1, Offset: offset, BaseId: hash + 1, StringId: ""}
}

func AddHashMapItem[T any](htx *HashMapContext[T], elementId HashElementId, element *T) (*HashMapItem[T], error) {
	if htx.HashMapInternal.Length() == htx.HashMapInternal.Capacity()-1 {
		return nil, errors.New("hashmap is full")
	}
	item := HashMapItem[T]{
		ElementId:  elementId,
		Element:    element,
		NextIndex:  -1,
		Generation: htx.Generation + 1,
	}

	// Perform modulo with uint32 first to avoid negative results, then cast to int32
	hashBucket := int32(elementId.Id % uint32(htx.HashMap.Capacity()))
	hashItemPrevious := int32(-1)
	hashItemIndex := MArray_GetValue(&htx.HashMap, hashBucket)

	// This loop is blocking
	for hashItemIndex != -1 { // Just replace collision, not a big deal - leave it up to the end user
		hashItem := MArray_Get[HashMapItem[T]](&htx.HashMapInternal, hashItemIndex)
		if hashItem.ElementId.Id == elementId.Id { // Collision - resolve based on generation
			item.NextIndex = hashItem.NextIndex
			if hashItem.Generation <= htx.Generation { // First collision - assume this is the "same" element
				hashItem.ElementId = elementId // Make sure to copy this across. If the stringId reference has changed, we should update the hash item to use the new one.
				hashItem.Generation = htx.Generation + 1
				hashItem.Element = element
			} else { // Multiple collisions this frame - two elements have the same ID
				return nil, errors.New("multiple collisions this frame - two elements have the same ID")
			}
			return hashItem, nil
		}
		hashItemPrevious = hashItemIndex
		hashItemIndex = hashItem.NextIndex
	}

	hashItem := MArray_Add(&htx.HashMapInternal, item)
	if hashItemPrevious != -1 {
		MArray_Get[HashMapItem[T]](&htx.HashMapInternal, hashItemPrevious).NextIndex = htx.HashMapInternal.Length() - 1
	} else {
		MArray_Set(&htx.HashMap, hashBucket, htx.HashMapInternal.Length()-1)
	}
	return hashItem, nil
}

type HashMapItem[T any] struct {
	ElementId  HashElementId
	Element    *T
	NextIndex  int32
	Generation uint32
}

func GetHashMapItem[T any](htx *HashMapContext[T], id uint32) (*HashMapItem[T], bool) {
	// Perform modulo with uint32 first to avoid negative results, then cast to int32
	hashBucket := int32(id % uint32(htx.HashMap.Capacity()))

	elementIndex := MArray_GetValue(&htx.HashMap, hashBucket)
	for elementIndex != -1 {
		hashEntry := MArray_Get[HashMapItem[T]](&htx.HashMapInternal, elementIndex)
		if hashEntry.ElementId.Id == id {
			return hashEntry, true
		}
		elementIndex = hashEntry.NextIndex
	}

	return nil, false
}
