
Clay__MeasureTextCacheItem *Clay__MeasureTextCached(Clay_String *text, Clay_TextElementConfig *config) {
    Clay_Context* context = Clay_GetCurrentContext();
    #ifndef CLAY_WASM
    if (!Clay__MeasureText) {
        if (!context->booleanWarnings.textMeasurementFunctionNotSet) {
            context->booleanWarnings.textMeasurementFunctionNotSet = true;
            context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                    .errorType = CLAY_ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED,
                    .errorText = CLAY_STRING("Clay's internal MeasureText function is null. You may have forgotten to call Clay_SetMeasureTextFunction(), or passed a NULL function pointer by mistake."),
                    .userData = context->errorHandler.userData });
        }
        return &Clay__MeasureTextCacheItem_DEFAULT;
    }
    #endif
    uint32_t id = Clay__HashStringContentsWithConfig(text, config);
    uint32_t hashBucket = id % (context->maxMeasureTextCacheWordCount / 32);
    int32_t elementIndexPrevious = 0;
    int32_t elementIndex = context->measureTextHashMap.internalArray[hashBucket];
    while (elementIndex != 0) {
        Clay__MeasureTextCacheItem *hashEntry = Clay__MeasureTextCacheItemArray_Get(&context->measureTextHashMapInternal, elementIndex);
        if (hashEntry->id == id) {
            hashEntry->generation = context->generation;
            return hashEntry;
        }
        // This element hasn't been seen in a few frames, delete the hash map item
        if (context->generation - hashEntry->generation > 2) {
            // Add all the measured words that were included in this measurement to the freelist
            int32_t nextWordIndex = hashEntry->measuredWordsStartIndex;
            while (nextWordIndex != -1) {
                Clay__MeasuredWord *measuredWord = Clay__MeasuredWordArray_Get(&context->measuredWords, nextWordIndex);
                Clay__int32_tArray_Add(&context->measuredWordsFreeList, nextWordIndex);
                nextWordIndex = measuredWord->next;
            }

            int32_t nextIndex = hashEntry->nextIndex;
            Clay__MeasureTextCacheItemArray_Set(&context->measureTextHashMapInternal, elementIndex, CLAY__INIT(Clay__MeasureTextCacheItem) { .measuredWordsStartIndex = -1 });
            Clay__int32_tArray_Add(&context->measureTextHashMapInternalFreeList, elementIndex);
            if (elementIndexPrevious == 0) {
                context->measureTextHashMap.internalArray[hashBucket] = nextIndex;
            } else {
                Clay__MeasureTextCacheItem *previousHashEntry = Clay__MeasureTextCacheItemArray_Get(&context->measureTextHashMapInternal, elementIndexPrevious);
                previousHashEntry->nextIndex = nextIndex;
            }
            elementIndex = nextIndex;
        } else {
            elementIndexPrevious = elementIndex;
            elementIndex = hashEntry->nextIndex;
        }
    }

    int32_t newItemIndex = 0;
    Clay__MeasureTextCacheItem newCacheItem = { .measuredWordsStartIndex = -1, .id = id, .generation = context->generation };
    Clay__MeasureTextCacheItem *measured = NULL;
    if (context->measureTextHashMapInternalFreeList.length > 0) {
        newItemIndex = Clay__int32_tArray_GetValue(&context->measureTextHashMapInternalFreeList, context->measureTextHashMapInternalFreeList.length - 1);
        context->measureTextHashMapInternalFreeList.length--;
        Clay__MeasureTextCacheItemArray_Set(&context->measureTextHashMapInternal, newItemIndex, newCacheItem);
        measured = Clay__MeasureTextCacheItemArray_Get(&context->measureTextHashMapInternal, newItemIndex);
    } else {
        if (context->measureTextHashMapInternal.length == context->measureTextHashMapInternal.capacity - 1) {
            if (!context->booleanWarnings.maxTextMeasureCacheExceeded) {
                context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                        .errorType = CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
                        .errorText = CLAY_STRING("Clay ran out of capacity while attempting to measure text elements. Try using Clay_SetMaxElementCount() with a higher value."),
                        .userData = context->errorHandler.userData });
                context->booleanWarnings.maxTextMeasureCacheExceeded = true;
            }
            return &Clay__MeasureTextCacheItem_DEFAULT;
        }
        measured = Clay__MeasureTextCacheItemArray_Add(&context->measureTextHashMapInternal, newCacheItem);
        newItemIndex = context->measureTextHashMapInternal.length - 1;
    }

    int32_t start = 0;
    int32_t end = 0;
    float lineWidth = 0;
    float measuredWidth = 0;
    float measuredHeight = 0;
    float spaceWidth = Clay__MeasureText(CLAY__INIT(Clay_StringSlice) { .length = 1, .chars = CLAY__SPACECHAR.chars, .baseChars = CLAY__SPACECHAR.chars }, config, context->measureTextUserData).width;
    Clay__MeasuredWord tempWord = { .next = -1 };
    Clay__MeasuredWord *previousWord = &tempWord;
    while (end < text->length) {
        if (context->measuredWords.length == context->measuredWords.capacity - 1) {
            if (!context->booleanWarnings.maxTextMeasureCacheExceeded) {
                context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                    .errorType = CLAY_ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED,
                    .errorText = CLAY_STRING("Clay has run out of space in it's internal text measurement cache. Try using Clay_SetMaxMeasureTextCacheWordCount() (default 16384, with 1 unit storing 1 measured word)."),
                    .userData = context->errorHandler.userData });
                context->booleanWarnings.maxTextMeasureCacheExceeded = true;
            }
            return &Clay__MeasureTextCacheItem_DEFAULT;
        }
        char current = text->chars[end];
        if (current == ' ' || current == '\n') {
            int32_t length = end - start;
            Clay_Dimensions dimensions = CLAY__DEFAULT_STRUCT;
            if (length > 0) {
                dimensions = Clay__MeasureText(CLAY__INIT(Clay_StringSlice) {.length = length, .chars = &text->chars[start], .baseChars = text->chars}, config, context->measureTextUserData);
            }
            measured->minWidth = CLAY__MAX(dimensions.width, measured->minWidth);
            measuredHeight = CLAY__MAX(measuredHeight, dimensions.height);
            if (current == ' ') {
                dimensions.width += spaceWidth;
                previousWord = Clay__AddMeasuredWord(CLAY__INIT(Clay__MeasuredWord) { .startOffset = start, .length = length + 1, .width = dimensions.width, .next = -1 }, previousWord);
                lineWidth += dimensions.width;
            }
            if (current == '\n') {
                if (length > 0) {
                    previousWord = Clay__AddMeasuredWord(CLAY__INIT(Clay__MeasuredWord) { .startOffset = start, .length = length, .width = dimensions.width, .next = -1 }, previousWord);
                }
                previousWord = Clay__AddMeasuredWord(CLAY__INIT(Clay__MeasuredWord) { .startOffset = end + 1, .length = 0, .width = 0, .next = -1 }, previousWord);
                lineWidth += dimensions.width;
                measuredWidth = CLAY__MAX(lineWidth, measuredWidth);
                measured->containsNewlines = true;
                lineWidth = 0;
            }
            start = end + 1;
        }
        end++;
    }
    if (end - start > 0) {
        Clay_Dimensions dimensions = Clay__MeasureText(CLAY__INIT(Clay_StringSlice) { .length = end - start, .chars = &text->chars[start], .baseChars = text->chars }, config, context->measureTextUserData);
        Clay__AddMeasuredWord(CLAY__INIT(Clay__MeasuredWord) { .startOffset = start, .length = end - start, .width = dimensions.width, .next = -1 }, previousWord);
        lineWidth += dimensions.width;
        measuredHeight = CLAY__MAX(measuredHeight, dimensions.height);
        measured->minWidth = CLAY__MAX(dimensions.width, measured->minWidth);
    }
    measuredWidth = CLAY__MAX(lineWidth, measuredWidth) - config->letterSpacing;

    measured->measuredWordsStartIndex = tempWord.next;
    measured->unwrappedDimensions.width = measuredWidth;
    measured->unwrappedDimensions.height = measuredHeight;

    if (elementIndexPrevious != 0) {
        Clay__MeasureTextCacheItemArray_Get(&context->measureTextHashMapInternal, elementIndexPrevious)->nextIndex = newItemIndex;
    } else {
        context->measureTextHashMap.internalArray[hashBucket] = newItemIndex;
    }
    return measured;
}
