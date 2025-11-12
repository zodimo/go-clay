package clay

func Clay__GetHashMapItem(id uint32) *Clay_LayoutElementHashMapItem {
	currentContext := Clay_GetCurrentContext()
	// Perform modulo with uint32 first to avoid negative results, then cast to int32
	hashBucket := int32(id % uint32(currentContext.LayoutElementsHashMap.Capacity))

	elementIndex := Clay__Array_GetValue(&currentContext.LayoutElementsHashMap, hashBucket)
	for elementIndex != -1 {
		hashEntry := Clay__Array_Get(&currentContext.LayoutElementsHashMapInternal, elementIndex)
		if hashEntry.ElementId.Id == id {
			return hashEntry
		}
		elementIndex = hashEntry.NextIndex
	}

	return &Clay_LayoutElementHashMapItem{}
}
