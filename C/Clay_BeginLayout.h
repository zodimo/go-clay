void Clay_BeginLayout(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay__InitializeEphemeralMemory(context);
    context->generation++;
    context->dynamicElementIndex = 0;
    // Set up the root container that covers the entire window
    Clay_Dimensions rootDimensions = {context->layoutDimensions.width, context->layoutDimensions.height};
    if (context->debugModeEnabled) {
        rootDimensions.width -= (float)Clay__debugViewWidth;
    }
    context->booleanWarnings = CLAY__INIT(Clay_BooleanWarnings) CLAY__DEFAULT_STRUCT;
    Clay__OpenElementWithId(CLAY_ID("Clay__RootContainer"));
    Clay__ConfigureOpenElement(CLAY__INIT(Clay_ElementDeclaration) {
        .layout = { .sizing = {CLAY_SIZING_FIXED((rootDimensions.width)), CLAY_SIZING_FIXED(rootDimensions.height)} }
    });
    Clay__int32_tArray_Add(&context->openLayoutElementStack, 0);
    Clay__LayoutElementTreeRootArray_Add(&context->layoutElementTreeRoots, CLAY__INIT(Clay__LayoutElementTreeRoot) { .layoutElementIndex = 0 });
}