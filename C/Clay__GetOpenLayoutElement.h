

Clay_LayoutElement* Clay__GetOpenLayoutElement(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    return Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 1));
}