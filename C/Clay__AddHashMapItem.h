





Clay_LayoutElementHashMapItem *Clay__GetHashMapItem(uint32_t id) {
    Clay_Context* context = Clay_GetCurrentContext();
    uint32_t hashBucket = id % context->layoutElementsHashMap.capacity;
    int32_t elementIndex = context->layoutElementsHashMap.internalArray[hashBucket];
    while (elementIndex != -1) {
        Clay_LayoutElementHashMapItem *hashEntry = Clay__LayoutElementHashMapItemArray_Get(&context->layoutElementsHashMapInternal, elementIndex);
        if (hashEntry->elementId.id == id) {
            return hashEntry;
        }
        elementIndex = hashEntry->nextIndex;
    }
    return &Clay_LayoutElementHashMapItem_DEFAULT;
}




Clay_LayoutElementHashMapItem* Clay__AddHashMapItem(Clay_ElementId elementId, Clay_LayoutElement* layoutElement) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->layoutElementsHashMapInternal.length == context->layoutElementsHashMapInternal.capacity - 1) {
        return NULL;
    }
    
    Clay_LayoutElementHashMapItem item = { 
        .elementId = elementId, 
        .layoutElement = layoutElement,
         .nextIndex = -1, 
         .generation = context->generation + 1
         };
    uint32_t hashBucket = elementId.id % context->layoutElementsHashMap.capacity;
    int32_t hashItemPrevious = -1;
    int32_t hashItemIndex = context->layoutElementsHashMap.internalArray[hashBucket];

    while (hashItemIndex != -1) { // Just replace collision, not a big deal - leave it up to the end user
        Clay_LayoutElementHashMapItem *hashItem = Clay__LayoutElementHashMapItemArray_Get(&context->layoutElementsHashMapInternal, hashItemIndex);
        if (hashItem->elementId.id == elementId.id) { // Collision - resolve based on generation
            item.nextIndex = hashItem->nextIndex;
            if (hashItem->generation <= context->generation) { // First collision - assume this is the "same" element
                hashItem->elementId = elementId; // Make sure to copy this across. If the stringId reference has changed, we should update the hash item to use the new one.
                hashItem->generation = context->generation + 1;
                hashItem->layoutElement = layoutElement;
                hashItem->debugData->collision = false;
                hashItem->onHoverFunction = NULL;
                hashItem->hoverFunctionUserData = 0;
            } else { // Multiple collisions this frame - two elements have the same ID
                context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                    .errorType = CLAY_ERROR_TYPE_DUPLICATE_ID,
                    .errorText = CLAY_STRING("An element with this ID was already previously declared during this layout."),
                    .userData = context->errorHandler.userData });
                if (context->debugModeEnabled) {
                    hashItem->debugData->collision = true;
                }
            }
            return hashItem;
        }
        hashItemPrevious = hashItemIndex;
        hashItemIndex = hashItem->nextIndex;
    }
    Clay_LayoutElementHashMapItem *hashItem = Clay__LayoutElementHashMapItemArray_Add(&context->layoutElementsHashMapInternal, item);
    hashItem->debugData = Clay__DebugElementDataArray_Add(&context->debugElementData, CLAY__INIT(Clay__DebugElementData) CLAY__DEFAULT_STRUCT);
    if (hashItemPrevious != -1) {
        Clay__LayoutElementHashMapItemArray_Get(&context->layoutElementsHashMapInternal, hashItemPrevious)->nextIndex = (int32_t)context->layoutElementsHashMapInternal.length - 1;
    } else {
        context->layoutElementsHashMap.internalArray[hashBucket] = (int32_t)context->layoutElementsHashMapInternal.length - 1;
    }
    return hashItem;
}


