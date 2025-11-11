
Clay_ElementConfig Clay__AttachElementConfig(Clay_ElementConfigUnion config, Clay__ElementConfigType type) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return CLAY__INIT(Clay_ElementConfig) CLAY__DEFAULT_STRUCT;
    }
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    openLayoutElement->elementConfigs.length++;
    return *Clay__ElementConfigArray_Add(&context->elementConfigs, CLAY__INIT(Clay_ElementConfig) { .type = type, .config = config });
}