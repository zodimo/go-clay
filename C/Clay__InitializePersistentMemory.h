





void Clay__InitializePersistentMemory(Clay_Context* context) {
    // Persistent memory - initialized once and not reset
    int32_t maxElementCount = context->maxElementCount;
    int32_t maxMeasureTextCacheWordCount = context->maxMeasureTextCacheWordCount;
    Clay_Arena *arena = &context->internalArena;

    context->scrollContainerDatas = Clay__ScrollContainerDataInternalArray_Allocate_Arena(100, arena);
    context->layoutElementsHashMapInternal = Clay__LayoutElementHashMapItemArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementsHashMap = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->measureTextHashMapInternal = Clay__MeasureTextCacheItemArray_Allocate_Arena(maxElementCount, arena);
    context->measureTextHashMapInternalFreeList = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->measuredWordsFreeList = Clay__int32_tArray_Allocate_Arena(maxMeasureTextCacheWordCount, arena);
    context->measureTextHashMap = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->measuredWords = Clay__MeasuredWordArray_Allocate_Arena(maxMeasureTextCacheWordCount, arena);
    context->pointerOverIds = Clay_ElementIdArray_Allocate_Arena(maxElementCount, arena);
    context->debugElementData = Clay__DebugElementDataArray_Allocate_Arena(maxElementCount, arena);
    context->arenaResetOffset = arena->nextAllocation;
}
