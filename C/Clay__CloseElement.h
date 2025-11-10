

void Clay__CloseElement(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return;
    }
    Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
    Clay_LayoutConfig *layoutConfig = openLayoutElement->layoutConfig;
    if (!layoutConfig) {
        openLayoutElement->layoutConfig = &Clay_LayoutConfig_DEFAULT;
        layoutConfig = &Clay_LayoutConfig_DEFAULT;
    }
    bool elementHasClipHorizontal = false;
    bool elementHasClipVertical = false;
    for (int32_t i = 0; i < openLayoutElement->elementConfigs.length; i++) {
        Clay_ElementConfig *config = Clay__ElementConfigArraySlice_Get(&openLayoutElement->elementConfigs, i);
        if (config->type == CLAY__ELEMENT_CONFIG_TYPE_CLIP) {
            elementHasClipHorizontal = config->config.clipElementConfig->horizontal;
            elementHasClipVertical = config->config.clipElementConfig->vertical;
            context->openClipElementStack.length--;
            break;
        } else if (config->type == CLAY__ELEMENT_CONFIG_TYPE_FLOATING) {
            context->openClipElementStack.length--;
        }
    }

    float leftRightPadding = (float)(layoutConfig->padding.left + layoutConfig->padding.right);
    float topBottomPadding = (float)(layoutConfig->padding.top + layoutConfig->padding.bottom);

    // Attach children to the current open element
    openLayoutElement->childrenOrTextContent.children.elements = &context->layoutElementChildren.internalArray[context->layoutElementChildren.length];
    if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
        openLayoutElement->dimensions.width = leftRightPadding;
        openLayoutElement->minDimensions.width = leftRightPadding;
        for (int32_t i = 0; i < openLayoutElement->childrenOrTextContent.children.length; i++) {
            int32_t childIndex = Clay__int32_tArray_GetValue(&context->layoutElementChildrenBuffer, (int)context->layoutElementChildrenBuffer.length - openLayoutElement->childrenOrTextContent.children.length + i);
            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, childIndex);
            openLayoutElement->dimensions.width += child->dimensions.width;
            openLayoutElement->dimensions.height = CLAY__MAX(openLayoutElement->dimensions.height, child->dimensions.height + topBottomPadding);
            // Minimum size of child elements doesn't matter to clip containers as they can shrink and hide their contents
            if (!elementHasClipHorizontal) {
                openLayoutElement->minDimensions.width += child->minDimensions.width;
            }
            if (!elementHasClipVertical) {
                openLayoutElement->minDimensions.height = CLAY__MAX(openLayoutElement->minDimensions.height, child->minDimensions.height + topBottomPadding);
            }
            Clay__int32_tArray_Add(&context->layoutElementChildren, childIndex);
        }
        float childGap = (float)(CLAY__MAX(openLayoutElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
        openLayoutElement->dimensions.width += childGap;
        if (!elementHasClipHorizontal) {
            openLayoutElement->minDimensions.width += childGap;
        }
    }
    else if (layoutConfig->layoutDirection == CLAY_TOP_TO_BOTTOM) {
        openLayoutElement->dimensions.height = topBottomPadding;
        openLayoutElement->minDimensions.height = topBottomPadding;
        for (int32_t i = 0; i < openLayoutElement->childrenOrTextContent.children.length; i++) {
            int32_t childIndex = Clay__int32_tArray_GetValue(&context->layoutElementChildrenBuffer, (int)context->layoutElementChildrenBuffer.length - openLayoutElement->childrenOrTextContent.children.length + i);
            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, childIndex);
            openLayoutElement->dimensions.height += child->dimensions.height;
            openLayoutElement->dimensions.width = CLAY__MAX(openLayoutElement->dimensions.width, child->dimensions.width + leftRightPadding);
            // Minimum size of child elements doesn't matter to clip containers as they can shrink and hide their contents
            if (!elementHasClipVertical) {
                openLayoutElement->minDimensions.height += child->minDimensions.height;
            }
            if (!elementHasClipHorizontal) {
                openLayoutElement->minDimensions.width = CLAY__MAX(openLayoutElement->minDimensions.width, child->minDimensions.width + leftRightPadding);
            }
            Clay__int32_tArray_Add(&context->layoutElementChildren, childIndex);
        }
        float childGap = (float)(CLAY__MAX(openLayoutElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
        openLayoutElement->dimensions.height += childGap;
        if (!elementHasClipVertical) {
            openLayoutElement->minDimensions.height += childGap;
        }
    }

    context->layoutElementChildrenBuffer.length -= openLayoutElement->childrenOrTextContent.children.length;

    // Clamp element min and max width to the values configured in the layout
    if (layoutConfig->sizing.width.type != CLAY__SIZING_TYPE_PERCENT) {
        if (layoutConfig->sizing.width.size.minMax.max <= 0) { // Set the max size if the user didn't specify, makes calculations easier
            layoutConfig->sizing.width.size.minMax.max = CLAY__MAXFLOAT;
        }
        openLayoutElement->dimensions.width = CLAY__MIN(CLAY__MAX(openLayoutElement->dimensions.width, layoutConfig->sizing.width.size.minMax.min), layoutConfig->sizing.width.size.minMax.max);
        openLayoutElement->minDimensions.width = CLAY__MIN(CLAY__MAX(openLayoutElement->minDimensions.width, layoutConfig->sizing.width.size.minMax.min), layoutConfig->sizing.width.size.minMax.max);
    } else {
        openLayoutElement->dimensions.width = 0;
    }

    // Clamp element min and max height to the values configured in the layout
    if (layoutConfig->sizing.height.type != CLAY__SIZING_TYPE_PERCENT) {
        if (layoutConfig->sizing.height.size.minMax.max <= 0) { // Set the max size if the user didn't specify, makes calculations easier
            layoutConfig->sizing.height.size.minMax.max = CLAY__MAXFLOAT;
        }
        openLayoutElement->dimensions.height = CLAY__MIN(CLAY__MAX(openLayoutElement->dimensions.height, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
        openLayoutElement->minDimensions.height = CLAY__MIN(CLAY__MAX(openLayoutElement->minDimensions.height, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
    } else {
        openLayoutElement->dimensions.height = 0;
    }

    Clay__UpdateAspectRatioBox(openLayoutElement);

    bool elementIsFloating = Clay__ElementHasConfig(openLayoutElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING);

    // Close the currently open element
    int32_t closingElementIndex = Clay__int32_tArray_RemoveSwapback(&context->openLayoutElementStack, (int)context->openLayoutElementStack.length - 1);

    // Get the currently open parent
    openLayoutElement = Clay__GetOpenLayoutElement();

    if (context->openLayoutElementStack.length > 1) {
        if(elementIsFloating) {
            openLayoutElement->floatingChildrenCount++;
            return;
        }
        openLayoutElement->childrenOrTextContent.children.length++;
        Clay__int32_tArray_Add(&context->layoutElementChildrenBuffer, closingElementIndex);
    }
}