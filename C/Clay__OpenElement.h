

void Clay__OpenElement(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->layoutElements.length == context->layoutElements.capacity - 1 || context->booleanWarnings.maxElementsExceeded) {
        context->booleanWarnings.maxElementsExceeded = true;
        return;
    }
    Clay_LayoutElement layoutElement = CLAY__DEFAULT_STRUCT;
    Clay_LayoutElement* openLayoutElement = Clay_LayoutElementArray_Add(&context->layoutElements, layoutElement);
    Clay__int32_tArray_Add(&context->openLayoutElementStack, context->layoutElements.length - 1);
    Clay__GenerateIdForAnonymousElement(openLayoutElement);
    if (context->openClipElementStack.length > 0) {
        Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, Clay__int32_tArray_GetValue(&context->openClipElementStack, (int)context->openClipElementStack.length - 1));
    } else {
        Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, 0);
    }
}
