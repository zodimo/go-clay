Clay_LayoutConfig * Clay__StoreLayoutConfig(Clay_LayoutConfig config) {  
    return Clay_GetCurrentContext()->booleanWarnings.maxElementsExceeded ? &CLAY_LAYOUT_DEFAULT : Clay__LayoutConfigArray_Add(&Clay_GetCurrentContext()->layoutConfigs, config); 
}
