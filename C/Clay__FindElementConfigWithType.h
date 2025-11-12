Clay_ElementConfigUnion Clay__FindElementConfigWithType(Clay_LayoutElement *element, Clay__ElementConfigType type) {
    for (int32_t i = 0; i < element->elementConfigs.length; i++) {
        Clay_ElementConfig *config = Clay__ElementConfigArraySlice_Get(&element->elementConfigs, i);
        if (config->type == type) {
            return config->config;
        }
    }
    return CLAY__INIT(Clay_ElementConfigUnion) { NULL };
}