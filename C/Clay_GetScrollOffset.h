Clay_Vector2 Clay_GetScrollOffset(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
    }
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    // If the element has no id attached at this point, we need to generate one
    if (openLayoutElement->id == 0) {
        Clay__GenerateIdForAnonymousElement(openLayoutElement);
    }
    for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
        Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
        if (mapping->layoutElement == openLayoutElement) {
            return mapping->scrollPosition;
        }
    }
    return CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
}