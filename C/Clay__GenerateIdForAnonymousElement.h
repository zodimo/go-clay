Clay_ElementId Clay__GenerateIdForAnonymousElement(Clay_LayoutElement *openLayoutElement) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay_LayoutElement *parentElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 2));
    uint32_t offset = parentElement->childrenOrTextContent.children.length + parentElement->floatingChildrenCount;
    Clay_ElementId elementId = Clay__HashNumber(offset, parentElement->id);
    openLayoutElement->id = elementId.id;
    Clay__AddHashMapItem(elementId, openLayoutElement);
    Clay__StringArray_Add(&context->layoutElementIdStrings, elementId.stringId);
    return elementId;
}