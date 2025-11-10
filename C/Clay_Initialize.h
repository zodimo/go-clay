

CLAY_WASM_EXPORT("Clay_Initialize")
Clay_Context* Clay_Initialize(Clay_Arena arena, Clay_Dimensions layoutDimensions, Clay_ErrorHandler errorHandler) {
    // Cacheline align memory passed in
    uintptr_t baseOffset = 64 - ((uintptr_t)arena.memory % 64);
    baseOffset = baseOffset == 64 ? 0 : baseOffset;
    arena.memory += baseOffset;
    Clay_Context *context = Clay__Context_Allocate_Arena(&arena);
    if (context == NULL) return NULL;
    // DEFAULTS
    Clay_Context *oldContext = Clay_GetCurrentContext();
    *context = CLAY__INIT(Clay_Context) {
        .maxElementCount = oldContext ? oldContext->maxElementCount : Clay__defaultMaxElementCount,
        .maxMeasureTextCacheWordCount = oldContext ? oldContext->maxMeasureTextCacheWordCount : Clay__defaultMaxMeasureTextWordCacheCount,
        .errorHandler = errorHandler.errorHandlerFunction ? errorHandler : CLAY__INIT(Clay_ErrorHandler) { Clay__ErrorHandlerFunctionDefault, 0 },
        .layoutDimensions = layoutDimensions,
        .internalArena = arena,
    };
    Clay_SetCurrentContext(context);
    Clay__InitializePersistentMemory(context);
    Clay__InitializeEphemeralMemory(context);
    for (int32_t i = 0; i < context->layoutElementsHashMap.capacity; ++i) {
        context->layoutElementsHashMap.internalArray[i] = -1;
    }
    for (int32_t i = 0; i < context->measureTextHashMap.capacity; ++i) {
        context->measureTextHashMap.internalArray[i] = 0;
    }
    context->measureTextHashMapInternal.length = 1; // Reserve the 0 value to mean "no next element"
    context->layoutDimensions = layoutDimensions;
    return context;
}
