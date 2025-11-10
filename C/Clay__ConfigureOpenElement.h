



void Clay__ConfigureOpenElementPtr(const Clay_ElementDeclaration *declaration) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    openLayoutElement->layoutConfig = Clay__StoreLayoutConfig(declaration->layout);
    if ((declaration->layout.sizing.width.type == CLAY__SIZING_TYPE_PERCENT && declaration->layout.sizing.width.size.percent > 1) || (declaration->layout.sizing.height.type == CLAY__SIZING_TYPE_PERCENT && declaration->layout.sizing.height.size.percent > 1)) {
        context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                .errorType = CLAY_ERROR_TYPE_PERCENTAGE_OVER_1,
                .errorText = CLAY_STRING("An element was configured with CLAY_SIZING_PERCENT, but the provided percentage value was over 1.0. Clay expects a value between 0 and 1, i.e. 20% is 0.2."),
                .userData = context->errorHandler.userData });
    }

    openLayoutElement->elementConfigs.internalArray = &context->elementConfigs.internalArray[context->elementConfigs.length];
    Clay_SharedElementConfig *sharedConfig = NULL;
    if (declaration->backgroundColor.a > 0) {
        sharedConfig = Clay__StoreSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .backgroundColor = declaration->backgroundColor });
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
    }
    if (!Clay__MemCmp((char *)(&declaration->cornerRadius), (char *)(&Clay__CornerRadius_DEFAULT), sizeof(Clay_CornerRadius))) {
        if (sharedConfig) {
            sharedConfig->cornerRadius = declaration->cornerRadius;
        } else {
            sharedConfig = Clay__StoreSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .cornerRadius = declaration->cornerRadius });
            Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
        }
    }
    if (declaration->userData != 0) {
        if (sharedConfig) {
            sharedConfig->userData = declaration->userData;
        } else {
            sharedConfig = Clay__StoreSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .userData = declaration->userData });
            Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
        }
    }
    if (declaration->image.imageData) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .imageElementConfig = Clay__StoreImageElementConfig(declaration->image) }, CLAY__ELEMENT_CONFIG_TYPE_IMAGE);
    }
    if (declaration->aspectRatio.aspectRatio > 0) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .aspectRatioElementConfig = Clay__StoreAspectRatioElementConfig(declaration->aspectRatio) }, CLAY__ELEMENT_CONFIG_TYPE_ASPECT);
        Clay__int32_tArray_Add(&context->aspectRatioElementIndexes, context->layoutElements.length - 1);
    }
    if (declaration->floating.attachTo != CLAY_ATTACH_TO_NONE) {
        Clay_FloatingElementConfig floatingConfig = declaration->floating;
        // This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here
        Clay_LayoutElement *hierarchicalParent = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 2));
        if (hierarchicalParent) {
            uint32_t clipElementId = 0;
            if (declaration->floating.attachTo == CLAY_ATTACH_TO_PARENT) {
                // Attach to the element's direct hierarchical parent
                floatingConfig.parentId = hierarchicalParent->id;
                if (context->openClipElementStack.length > 0) {
                    clipElementId = Clay__int32_tArray_GetValue(&context->openClipElementStack, (int)context->openClipElementStack.length - 1);
                }
            } else if (declaration->floating.attachTo == CLAY_ATTACH_TO_ELEMENT_WITH_ID) {
                Clay_LayoutElementHashMapItem *parentItem = Clay__GetHashMapItem(floatingConfig.parentId);
                if (parentItem == &Clay_LayoutElementHashMapItem_DEFAULT) {
                    context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                            .errorType = CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND,
                            .errorText = CLAY_STRING("A floating element was declared with a parentId, but no element with that ID was found."),
                            .userData = context->errorHandler.userData });
                } else {
                    clipElementId = Clay__int32_tArray_GetValue(&context->layoutElementClipElementIds, (int32_t)(parentItem->layoutElement - context->layoutElements.internalArray));
                }
            } else if (declaration->floating.attachTo == CLAY_ATTACH_TO_ROOT) {
                floatingConfig.parentId = Clay__HashString(CLAY_STRING("Clay__RootContainer"), 0).id;
            }
            if (declaration->floating.clipTo == CLAY_CLIP_TO_NONE) {
                clipElementId = 0;
            }
            int32_t currentElementIndex = Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 1);
            Clay__int32_tArray_Set(&context->layoutElementClipElementIds, currentElementIndex, clipElementId);
            Clay__int32_tArray_Add(&context->openClipElementStack, clipElementId);
            Clay__LayoutElementTreeRootArray_Add(&context->layoutElementTreeRoots, CLAY__INIT(Clay__LayoutElementTreeRoot) {
                    .layoutElementIndex = Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 1),
                    .parentId = floatingConfig.parentId,
                    .clipElementId = clipElementId,
                    .zIndex = floatingConfig.zIndex,
            });
            Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .floatingElementConfig = Clay__StoreFloatingElementConfig(floatingConfig) }, CLAY__ELEMENT_CONFIG_TYPE_FLOATING);
        }
    }
    if (declaration->custom.customData) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .customElementConfig = Clay__StoreCustomElementConfig(declaration->custom) }, CLAY__ELEMENT_CONFIG_TYPE_CUSTOM);
    }

    if (declaration->clip.horizontal | declaration->clip.vertical) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .clipElementConfig = Clay__StoreClipElementConfig(declaration->clip) }, CLAY__ELEMENT_CONFIG_TYPE_CLIP);
        Clay__int32_tArray_Add(&context->openClipElementStack, (int)openLayoutElement->id);
        // Retrieve or create cached data to track scroll position across frames
        Clay__ScrollContainerDataInternal *scrollOffset = CLAY__NULL;
        for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
            Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
            if (openLayoutElement->id == mapping->elementId) {
                scrollOffset = mapping;
                scrollOffset->layoutElement = openLayoutElement;
                scrollOffset->openThisFrame = true;
            }
        }
        if (!scrollOffset) {
            scrollOffset = Clay__ScrollContainerDataInternalArray_Add(&context->scrollContainerDatas, CLAY__INIT(Clay__ScrollContainerDataInternal){.layoutElement = openLayoutElement, .scrollOrigin = {-1,-1}, .elementId = openLayoutElement->id, .openThisFrame = true});
        }
        if (context->externalScrollHandlingEnabled) {
            scrollOffset->scrollPosition = Clay__QueryScrollOffset(scrollOffset->elementId, context->queryScrollOffsetUserData);
        }
    }
    if (!Clay__MemCmp((char *)(&declaration->border.width), (char *)(&Clay__BorderWidth_DEFAULT), sizeof(Clay_BorderWidth))) {
        Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .borderElementConfig = Clay__StoreBorderElementConfig(declaration->border) }, CLAY__ELEMENT_CONFIG_TYPE_BORDER);
    }
}

void Clay__ConfigureOpenElement(const Clay_ElementDeclaration declaration) {
    Clay__ConfigureOpenElementPtr(&declaration);
}
