






void Clay__InitializeEphemeralMemory(Clay_Context* context) {
    int32_t maxElementCount = context->maxElementCount;
    // Ephemeral Memory - reset every frame
    Clay_Arena *arena = &context->internalArena;
    arena->nextAllocation = context->arenaResetOffset;

    context->layoutElementChildrenBuffer = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElements = Clay_LayoutElementArray_Allocate_Arena(maxElementCount, arena);
    context->warnings = Clay__WarningArray_Allocate_Arena(100, arena);

    context->layoutConfigs = Clay__LayoutConfigArray_Allocate_Arena(maxElementCount, arena);
    context->elementConfigs = Clay__ElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->textElementConfigs = Clay__TextElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->aspectRatioElementConfigs = Clay__AspectRatioElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->imageElementConfigs = Clay__ImageElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->floatingElementConfigs = Clay__FloatingElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->clipElementConfigs = Clay__ClipElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->customElementConfigs = Clay__CustomElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->borderElementConfigs = Clay__BorderElementConfigArray_Allocate_Arena(maxElementCount, arena);
    context->sharedElementConfigs = Clay__SharedElementConfigArray_Allocate_Arena(maxElementCount, arena);

    context->layoutElementIdStrings = Clay__StringArray_Allocate_Arena(maxElementCount, arena);
    context->wrappedTextLines = Clay__WrappedTextLineArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementTreeNodeArray1 = Clay__LayoutElementTreeNodeArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementTreeRoots = Clay__LayoutElementTreeRootArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementChildren = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->openLayoutElementStack = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->textElementData = Clay__TextElementDataArray_Allocate_Arena(maxElementCount, arena);
    context->aspectRatioElementIndexes = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->renderCommands = Clay_RenderCommandArray_Allocate_Arena(maxElementCount, arena);
    context->treeNodeVisited = Clay__boolArray_Allocate_Arena(maxElementCount, arena);
    context->treeNodeVisited.length = context->treeNodeVisited.capacity; // This array is accessed directly rather than behaving as a list
    context->openClipElementStack = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->reusableElementIndexBuffer = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->layoutElementClipElementIds = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
    context->dynamicStringData = Clay__charArray_Allocate_Arena(maxElementCount, arena);
}
