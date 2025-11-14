package clay

import "unsafe"

func Clay__HashString(key Clay_String, seed uint32) Clay_ElementId {
	hash := seed

	charsBytes := []byte(key.Chars)

	for _, charByte := range charsBytes {
		hash += uint32(charByte)
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}

	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return Clay_ElementId{
		Id:       hash + 1,
		Offset:   0,
		BaseId:   hash + 1,
		StringId: key,
	} // Reserve the hash result of zero as "null id"
}

func Clay__HashNumber(offset uint32, seed uint32) Clay_ElementId {
	hash := seed
	hash += (offset + 48)
	hash += (hash << 10)
	hash ^= (hash >> 6)
	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return Clay_ElementId{Id: hash + 1, Offset: offset, BaseId: hash + 1, StringId: CLAY__STRING_DEFAULT}
}

func Clay__HashData(data []byte, length int32) uint32 {
	hash := uint32(0)

	for _, charByte := range data[:length] {
		hash += uint32(charByte)
		hash += (hash << 10)
		hash ^= (hash >> 6)
	}
	return hash
}

func Clay__HashStringContentsWithConfig(text *Clay_String, config *Clay_TextElementConfig) uint32 {
	hash := uint32(0)

	if text.IsStaticallyAllocated {
		hash += uint32(uintptr(unsafe.Pointer(&text.Chars)))
		hash += (hash << 10)
		hash ^= (hash >> 6)
		hash += uint32(text.Length)
		hash += (hash << 10)
		hash ^= (hash >> 6)
	} else {
		hash = Clay__HashData(text.Chars, text.Length) % UINT32_MAX
	}

	hash += uint32(config.FontId)
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += uint32(config.FontSize)
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += uint32(config.LetterSpacing)
	hash += (hash << 10)
	hash ^= (hash >> 6)

	hash += (hash << 3)
	hash ^= (hash >> 11)
	hash += (hash << 15)
	return hash + 1 // Reserve the hash result of zero as "null id"
}
