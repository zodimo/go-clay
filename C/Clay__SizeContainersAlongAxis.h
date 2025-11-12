void Clay__SizeContainersAlongAxis(bool xAxis) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay__int32_tArray bfsBuffer = context->layoutElementChildrenBuffer;
    Clay__int32_tArray resizableContainerBuffer = context->openLayoutElementStack;
    for (int32_t rootIndex = 0; rootIndex < context->layoutElementTreeRoots.length; ++rootIndex) {
        bfsBuffer.length = 0;
        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, rootIndex);
        Clay_LayoutElement *rootElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)root->layoutElementIndex);
        Clay__int32_tArray_Add(&bfsBuffer, (int32_t)root->layoutElementIndex);

        // Size floating containers to their parents
        if (Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)) {
            Clay_FloatingElementConfig *floatingElementConfig = Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
            Clay_LayoutElementHashMapItem *parentItem = Clay__GetHashMapItem(floatingElementConfig->parentId);
            if (parentItem && parentItem != &Clay_LayoutElementHashMapItem_DEFAULT) {
                Clay_LayoutElement *parentLayoutElement = parentItem->layoutElement;
                switch (rootElement->layoutConfig->sizing.width.type) {
                    case CLAY__SIZING_TYPE_GROW: {
                        rootElement->dimensions.width = parentLayoutElement->dimensions.width;
                        break;
                    }
                    case CLAY__SIZING_TYPE_PERCENT: {
                        rootElement->dimensions.width = parentLayoutElement->dimensions.width * rootElement->layoutConfig->sizing.width.size.percent;
                        break;
                    }
                    default: break;
                }
                switch (rootElement->layoutConfig->sizing.height.type) {
                    case CLAY__SIZING_TYPE_GROW: {
                        rootElement->dimensions.height = parentLayoutElement->dimensions.height;
                        break;
                    }
                    case CLAY__SIZING_TYPE_PERCENT: {
                        rootElement->dimensions.height = parentLayoutElement->dimensions.height * rootElement->layoutConfig->sizing.height.size.percent;
                        break;
                    }
                    default: break;
                }
            }
        }

        if (rootElement->layoutConfig->sizing.width.type != CLAY__SIZING_TYPE_PERCENT) {
            rootElement->dimensions.width = CLAY__MIN(CLAY__MAX(rootElement->dimensions.width, rootElement->layoutConfig->sizing.width.size.minMax.min), rootElement->layoutConfig->sizing.width.size.minMax.max);
        }
        if (rootElement->layoutConfig->sizing.height.type != CLAY__SIZING_TYPE_PERCENT) {
            rootElement->dimensions.height = CLAY__MIN(CLAY__MAX(rootElement->dimensions.height, rootElement->layoutConfig->sizing.height.size.minMax.min), rootElement->layoutConfig->sizing.height.size.minMax.max);
        }

        for (int32_t i = 0; i < bfsBuffer.length; ++i) {
            int32_t parentIndex = Clay__int32_tArray_GetValue(&bfsBuffer, i);
            Clay_LayoutElement *parent = Clay_LayoutElementArray_Get(&context->layoutElements, parentIndex);
            Clay_LayoutConfig *parentStyleConfig = parent->layoutConfig;
            int32_t growContainerCount = 0;
            float parentSize = xAxis ? parent->dimensions.width : parent->dimensions.height;
            float parentPadding = (float)(xAxis ? (parent->layoutConfig->padding.left + parent->layoutConfig->padding.right) : (parent->layoutConfig->padding.top + parent->layoutConfig->padding.bottom));
            float innerContentSize = 0, totalPaddingAndChildGaps = parentPadding;
            bool sizingAlongAxis = (xAxis && parentStyleConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) || (!xAxis && parentStyleConfig->layoutDirection == CLAY_TOP_TO_BOTTOM);
            resizableContainerBuffer.length = 0;
            float parentChildGap = parentStyleConfig->childGap;

            for (int32_t childOffset = 0; childOffset < parent->childrenOrTextContent.children.length; childOffset++) {
                int32_t childElementIndex = parent->childrenOrTextContent.children.elements[childOffset];
                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, childElementIndex);
                Clay_SizingAxis childSizing = xAxis ? childElement->layoutConfig->sizing.width : childElement->layoutConfig->sizing.height;
                float childSize = xAxis ? childElement->dimensions.width : childElement->dimensions.height;

                if (!Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) && childElement->childrenOrTextContent.children.length > 0) {
                    Clay__int32_tArray_Add(&bfsBuffer, childElementIndex);
                }

                if (childSizing.type != CLAY__SIZING_TYPE_PERCENT
                    && childSizing.type != CLAY__SIZING_TYPE_FIXED
                    && (!Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || (Clay__FindElementConfigWithType(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).textElementConfig->wrapMode == CLAY_TEXT_WRAP_WORDS)) // todo too many loops
//                    && (xAxis || !Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT))
                ) {
                    Clay__int32_tArray_Add(&resizableContainerBuffer, childElementIndex);
                }

                if (sizingAlongAxis) {
                    innerContentSize += (childSizing.type == CLAY__SIZING_TYPE_PERCENT ? 0 : childSize);
                    if (childSizing.type == CLAY__SIZING_TYPE_GROW) {
                        growContainerCount++;
                    }
                    if (childOffset > 0) {
                        innerContentSize += parentChildGap; // For children after index 0, the childAxisOffset is the gap from the previous child
                        totalPaddingAndChildGaps += parentChildGap;
                    }
                } else {
                    innerContentSize = CLAY__MAX(childSize, innerContentSize);
                }
            }

            // Expand percentage containers to size
            for (int32_t childOffset = 0; childOffset < parent->childrenOrTextContent.children.length; childOffset++) {
                int32_t childElementIndex = parent->childrenOrTextContent.children.elements[childOffset];
                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, childElementIndex);
                Clay_SizingAxis childSizing = xAxis ? childElement->layoutConfig->sizing.width : childElement->layoutConfig->sizing.height;
                float *childSize = xAxis ? &childElement->dimensions.width : &childElement->dimensions.height;
                if (childSizing.type == CLAY__SIZING_TYPE_PERCENT) {
                    *childSize = (parentSize - totalPaddingAndChildGaps) * childSizing.size.percent;
                    if (sizingAlongAxis) {
                        innerContentSize += *childSize;
                    }
                    Clay__UpdateAspectRatioBox(childElement);
                }
            }

            if (sizingAlongAxis) {
                float sizeToDistribute = parentSize - parentPadding - innerContentSize;
                // The content is too large, compress the children as much as possible
                if (sizeToDistribute < 0) {
                    // If the parent clips content in this axis direction, don't compress children, just leave them alone
                    Clay_ClipElementConfig *clipElementConfig = Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
                    if (clipElementConfig) {
                        if (((xAxis && clipElementConfig->horizontal) || (!xAxis && clipElementConfig->vertical))) {
                            continue;
                        }
                    }
                    // Scrolling containers preferentially compress before others
                    while (sizeToDistribute < -CLAY__EPSILON && resizableContainerBuffer.length > 0) {
                        float largest = 0;
                        float secondLargest = 0;
                        float widthToAdd = sizeToDistribute;
                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                            float childSize = xAxis ? child->dimensions.width : child->dimensions.height;
                            if (Clay__FloatEqual(childSize, largest)) { continue; }
                            if (childSize > largest) {
                                secondLargest = largest;
                                largest = childSize;
                            }
                            if (childSize < largest) {
                                secondLargest = CLAY__MAX(secondLargest, childSize);
                                widthToAdd = secondLargest - largest;
                            }
                        }

                        widthToAdd = CLAY__MAX(widthToAdd, sizeToDistribute / resizableContainerBuffer.length);

                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                            float *childSize = xAxis ? &child->dimensions.width : &child->dimensions.height;
                            float minSize = xAxis ? child->minDimensions.width : child->minDimensions.height;
                            float previousWidth = *childSize;
                            if (Clay__FloatEqual(*childSize, largest)) {
                                *childSize += widthToAdd;
                                if (*childSize <= minSize) {
                                    *childSize = minSize;
                                    Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex--);
                                }
                                sizeToDistribute -= (*childSize - previousWidth);
                            }
                        }
                    }
                // The content is too small, allow SIZING_GROW containers to expand
                } else if (sizeToDistribute > 0 && growContainerCount > 0) {
                    for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                        Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                        Clay__SizingType childSizing = xAxis ? child->layoutConfig->sizing.width.type : child->layoutConfig->sizing.height.type;
                        if (childSizing != CLAY__SIZING_TYPE_GROW) {
                            Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex--);
                        }
                    }
                    while (sizeToDistribute > CLAY__EPSILON && resizableContainerBuffer.length > 0) {
                        float smallest = CLAY__MAXFLOAT;
                        float secondSmallest = CLAY__MAXFLOAT;
                        float widthToAdd = sizeToDistribute;
                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                            float childSize = xAxis ? child->dimensions.width : child->dimensions.height;
                            if (Clay__FloatEqual(childSize, smallest)) { continue; }
                            if (childSize < smallest) {
                                secondSmallest = smallest;
                                smallest = childSize;
                            }
                            if (childSize > smallest) {
                                secondSmallest = CLAY__MIN(secondSmallest, childSize);
                                widthToAdd = secondSmallest - smallest;
                            }
                        }

                        widthToAdd = CLAY__MIN(widthToAdd, sizeToDistribute / resizableContainerBuffer.length);

                        for (int childIndex = 0; childIndex < resizableContainerBuffer.length; childIndex++) {
                            Clay_LayoutElement *child = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childIndex));
                            float *childSize = xAxis ? &child->dimensions.width : &child->dimensions.height;
                            float maxSize = xAxis ? child->layoutConfig->sizing.width.size.minMax.max : child->layoutConfig->sizing.height.size.minMax.max;
                            float previousWidth = *childSize;
                            if (Clay__FloatEqual(*childSize, smallest)) {
                                *childSize += widthToAdd;
                                if (*childSize >= maxSize) {
                                    *childSize = maxSize;
                                    Clay__int32_tArray_RemoveSwapback(&resizableContainerBuffer, childIndex--);
                                }
                                sizeToDistribute -= (*childSize - previousWidth);
                            }
                        }
                    }
                }
            // Sizing along the non layout axis ("off axis")
            } else {
                for (int32_t childOffset = 0; childOffset < resizableContainerBuffer.length; childOffset++) {
                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&resizableContainerBuffer, childOffset));
                    Clay_SizingAxis childSizing = xAxis ? childElement->layoutConfig->sizing.width : childElement->layoutConfig->sizing.height;
                    float minSize = xAxis ? childElement->minDimensions.width : childElement->minDimensions.height;
                    float *childSize = xAxis ? &childElement->dimensions.width : &childElement->dimensions.height;

                    float maxSize = parentSize - parentPadding;
                    // If we're laying out the children of a scroll panel, grow containers expand to the size of the inner content, not the outer container
                    if (Clay__ElementHasConfig(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP)) {
                        Clay_ClipElementConfig *clipElementConfig = Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
                        if (((xAxis && clipElementConfig->horizontal) || (!xAxis && clipElementConfig->vertical))) {
                            maxSize = CLAY__MAX(maxSize, innerContentSize);
                        }
                    }
                    if (childSizing.type == CLAY__SIZING_TYPE_GROW) {
                        *childSize = CLAY__MIN(maxSize, childSizing.size.minMax.max);
                    }
                    *childSize = CLAY__MAX(minSize, CLAY__MIN(*childSize, maxSize));
                }
            }
        }
    }
}
