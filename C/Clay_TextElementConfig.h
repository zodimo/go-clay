

void Clay__OpenTextElement(Clay_String text, Clay_TextElementConfig *textConfig) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->layoutElements.length == context->layoutElements.capacity - 1 || context->booleanWarnings.maxElementsExceeded) {
        context->booleanWarnings.maxElementsExceeded = true;
        return;
    }
    Clay_LayoutElement *parentElement = Clay__GetOpenLayoutElement();

    Clay_LayoutElement layoutElement = CLAY__DEFAULT_STRUCT;
    Clay_LayoutElement *textElement = Clay_LayoutElementArray_Add(&context->layoutElements, layoutElement);
    if (context->openClipElementStack.length > 0) {
        Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, Clay__int32_tArray_GetValue(&context->openClipElementStack, (int)context->openClipElementStack.length - 1));
    } else {
        Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, 0);
    }

    Clay__int32_tArray_Add(&context->layoutElementChildrenBuffer, context->layoutElements.length - 1);
    Clay__MeasureTextCacheItem *textMeasured = Clay__MeasureTextCached(&text, textConfig);
    Clay_ElementId elementId = Clay__HashNumber(parentElement->childrenOrTextContent.children.length, parentElement->id);
    textElement->id = elementId.id;
    Clay__AddHashMapItem(elementId, textElement);
    Clay__StringArray_Add(&context->layoutElementIdStrings, elementId.stringId);
    Clay_Dimensions textDimensions = { .width = textMeasured->unwrappedDimensions.width, .height = textConfig->lineHeight > 0 ? (float)textConfig->lineHeight : textMeasured->unwrappedDimensions.height };
    textElement->dimensions = textDimensions;
    textElement->minDimensions = CLAY__INIT(Clay_Dimensions) { .width = textMeasured->minWidth, .height = textDimensions.height };
    textElement->childrenOrTextContent.textElementData = Clay__TextElementDataArray_Add(&context->textElementData, CLAY__INIT(Clay__TextElementData) { .text = text, .preferredDimensions = textMeasured->unwrappedDimensions, .elementIndex = context->layoutElements.length - 1 });
    textElement->elementConfigs = CLAY__INIT(Clay__ElementConfigArraySlice) {
            .length = 1,
            .internalArray = Clay__ElementConfigArray_Add(&context->elementConfigs, CLAY__INIT(Clay_ElementConfig) { .type = CLAY__ELEMENT_CONFIG_TYPE_TEXT, .config = { .textElementConfig = textConfig }})
    };
    textElement->layoutConfig = &CLAY_LAYOUT_DEFAULT;
    parentElement->childrenOrTextContent.children.length++;
}
