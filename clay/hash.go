package clay

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
