bool Clay_Hovered(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return false;
    }
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    // If the element has no id attached at this point, we need to generate one
    if (openLayoutElement->id == 0) {
        Clay__GenerateIdForAnonymousElement(openLayoutElement);
    }
    for (int32_t i = 0; i < context->pointerOverIds.length; ++i) {
        if (Clay_ElementIdArray_Get(&context->pointerOverIds, i)->id == openLayoutElement->id) {
            return true;
        }
    }
    return false;
}