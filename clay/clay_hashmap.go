package clay

import "fmt"

func Clay__GetHashMapItem(id uint32) *Clay_LayoutElementHashMapItem {
	currentContext := Clay_GetCurrentContext()
	// Perform modulo with uint32 first to avoid negative results, then cast to int32
	hashBucket := int32(id % uint32(currentContext.LayoutElementsHashMap.Capacity()))

	elementIndex := Clay__Array_GetValue(&currentContext.LayoutElementsHashMap, hashBucket)
	for elementIndex != -1 {
		hashEntry := Clay__Array_Get(&currentContext.LayoutElementsHashMapInternal, elementIndex)
		if hashEntry.ElementId.Id == id {
			return hashEntry
		}
		elementIndex = hashEntry.NextIndex
	}

	return &Clay_LayoutElementHashMapItem_DEFAULT
}

func Clay__AddHashMapItem(elementId Clay_ElementId, layoutElement *Clay_LayoutElement) *Clay_LayoutElementHashMapItem {
	currentContext := Clay_GetCurrentContext()
	if currentContext.LayoutElementsHashMapInternal.Length() == currentContext.LayoutElementsHashMapInternal.Capacity()-1 {
		return nil
	}
	item := Clay_LayoutElementHashMapItem{
		ElementId:     elementId,
		LayoutElement: layoutElement,
		NextIndex:     -1,
		Generation:    currentContext.Generation + 1,
	}

	// Perform modulo with uint32 first to avoid negative results, then cast to int32
	hashBucket := int32(elementId.Id % uint32(currentContext.LayoutElementsHashMap.Capacity()))
	hashItemPrevious := int32(-1)
	hashItemIndex := Clay__Array_GetValue(&currentContext.LayoutElementsHashMap, hashBucket)
	for hashItemIndex != -1 { // Just replace collision, not a big deal - leave it up to the end user
		hashItem := Clay__Array_GetUnsafe[Clay_LayoutElementHashMapItem](&currentContext.LayoutElementsHashMapInternal, hashItemIndex)
		fmt.Println("hashItem", hashItem)
		if hashItem == &Clay_LayoutElementHashMapItem_DEFAULT {
			panic("hashItem is default value")
		}
		if hashItem.ElementId.Id == elementId.Id { // Collision - resolve based on generation
			item.NextIndex = hashItem.NextIndex
			if hashItem.Generation <= currentContext.Generation { // First collision - assume this is the "same" element
				hashItem.ElementId = elementId // Make sure to copy this across. If the stringId reference has changed, we should update the hash item to use the new one.
				hashItem.Generation = currentContext.Generation + 1
				hashItem.LayoutElement = layoutElement
				hashItem.DebugData.Collision = false
				hashItem.OnHoverFunction = nil
				hashItem.HoverFunctionUserData = 0
			} else { // Multiple collisions this frame - two elements have the same ID
				currentContext.ErrorHandler.ErrorHandlerFunction(Clay_ErrorData{
					ErrorType: CLAY_ERROR_TYPE_DUPLICATE_ID,
					ErrorText: CLAY_STRING("An element with this ID was already previously declared during this layout."),
					UserData:  currentContext.ErrorHandler.UserData,
				})
				if currentContext.DebugModeEnabled {
					hashItem.DebugData.Collision = true
				}
			}
			return hashItem
		}
		hashItemPrevious = hashItemIndex
		hashItemIndex = hashItem.NextIndex
	}

	hashItem := Clay__Array_Add(&currentContext.LayoutElementsHashMapInternal, item)
	hashItem.DebugData = Clay__Array_Add(&currentContext.DebugElementData, Clay__DebugElementData{})
	if hashItemPrevious != -1 {
		Clay__Array_Get[Clay_LayoutElementHashMapItem](&currentContext.LayoutElementsHashMapInternal, hashItemPrevious).NextIndex = currentContext.LayoutElementsHashMapInternal.Length() - 1
	} else {
		Clay__Array_Set(&currentContext.LayoutElementsHashMap, hashBucket, currentContext.LayoutElementsHashMapInternal.Length()-1)
	}
	return hashItem
}
