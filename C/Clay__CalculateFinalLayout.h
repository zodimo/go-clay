
void Clay__CalculateFinalLayout(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    // Calculate sizing along the X axis
    Clay__SizeContainersAlongAxis(true);

    // Wrap text
    for (int32_t textElementIndex = 0; textElementIndex < context->textElementData.length; ++textElementIndex) {
        Clay__TextElementData *textElementData = Clay__TextElementDataArray_Get(&context->textElementData, textElementIndex);
        textElementData->wrappedLines = CLAY__INIT(Clay__WrappedTextLineArraySlice) { .length = 0, .internalArray = &context->wrappedTextLines.internalArray[context->wrappedTextLines.length] };
        Clay_LayoutElement *containerElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)textElementData->elementIndex);
        Clay_TextElementConfig *textConfig = Clay__FindElementConfigWithType(containerElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).textElementConfig;
        Clay__MeasureTextCacheItem *measureTextCacheItem = Clay__MeasureTextCached(&textElementData->text, textConfig);
        float lineWidth = 0;
        float lineHeight = textConfig->lineHeight > 0 ? (float)textConfig->lineHeight : textElementData->preferredDimensions.height;
        int32_t lineLengthChars = 0;
        int32_t lineStartOffset = 0;
        if (!measureTextCacheItem->containsNewlines && textElementData->preferredDimensions.width <= containerElement->dimensions.width) {
            Clay__WrappedTextLineArray_Add(&context->wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { containerElement->dimensions,  textElementData->text });
            textElementData->wrappedLines.length++;
            continue;
        }
        float spaceWidth = Clay__MeasureText(CLAY__INIT(Clay_StringSlice) { .length = 1, .chars = CLAY__SPACECHAR.chars, .baseChars = CLAY__SPACECHAR.chars }, textConfig, context->measureTextUserData).width;
        int32_t wordIndex = measureTextCacheItem->measuredWordsStartIndex;
        while (wordIndex != -1) {
            if (context->wrappedTextLines.length > context->wrappedTextLines.capacity - 1) {
                break;
            }
            Clay__MeasuredWord *measuredWord = Clay__MeasuredWordArray_Get(&context->measuredWords, wordIndex);
            // Only word on the line is too large, just render it anyway
            if (lineLengthChars == 0 && lineWidth + measuredWord->width > containerElement->dimensions.width) {
                Clay__WrappedTextLineArray_Add(&context->wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { { measuredWord->width, lineHeight }, { .length = measuredWord->length, .chars = &textElementData->text.chars[measuredWord->startOffset] } });
                textElementData->wrappedLines.length++;
                wordIndex = measuredWord->next;
                lineStartOffset = measuredWord->startOffset + measuredWord->length;
            }
            // measuredWord->length == 0 means a newline character
            else if (measuredWord->length == 0 || lineWidth + measuredWord->width > containerElement->dimensions.width) {
                // Wrapped text lines list has overflowed, just render out the line
                bool finalCharIsSpace = textElementData->text.chars[CLAY__MAX(lineStartOffset + lineLengthChars - 1, 0)] == ' ';
                Clay__WrappedTextLineArray_Add(&context->wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { { lineWidth + (finalCharIsSpace ? -spaceWidth : 0), lineHeight }, { .length = lineLengthChars + (finalCharIsSpace ? -1 : 0), .chars = &textElementData->text.chars[lineStartOffset] } });
                textElementData->wrappedLines.length++;
                if (lineLengthChars == 0 || measuredWord->length == 0) {
                    wordIndex = measuredWord->next;
                }
                lineWidth = 0;
                lineLengthChars = 0;
                lineStartOffset = measuredWord->startOffset;
            } else {
                lineWidth += measuredWord->width + textConfig->letterSpacing;
                lineLengthChars += measuredWord->length;
                wordIndex = measuredWord->next;
            }
        }
        if (lineLengthChars > 0) {
            Clay__WrappedTextLineArray_Add(&context->wrappedTextLines, CLAY__INIT(Clay__WrappedTextLine) { { lineWidth - textConfig->letterSpacing, lineHeight }, {.length = lineLengthChars, .chars = &textElementData->text.chars[lineStartOffset] } });
            textElementData->wrappedLines.length++;
        }
        containerElement->dimensions.height = lineHeight * (float)textElementData->wrappedLines.length;
    }

    // Scale vertical heights according to aspect ratio
    for (int32_t i = 0; i < context->aspectRatioElementIndexes.length; ++i) {
        Clay_LayoutElement* aspectElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->aspectRatioElementIndexes, i));
        Clay_AspectRatioElementConfig *config = Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig;
        aspectElement->dimensions.height = (1 / config->aspectRatio) * aspectElement->dimensions.width;
        aspectElement->layoutConfig->sizing.height.size.minMax.max = aspectElement->dimensions.height;
    }

    // Propagate effect of text wrapping, aspect scaling etc. on height of parents
    Clay__LayoutElementTreeNodeArray dfsBuffer = context->layoutElementTreeNodeArray1;
    dfsBuffer.length = 0;
    for (int32_t i = 0; i < context->layoutElementTreeRoots.length; ++i) {
        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, i);
        context->treeNodeVisited.internalArray[dfsBuffer.length] = false;
        Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)root->layoutElementIndex) });
    }
    while (dfsBuffer.length > 0) {
        Clay__LayoutElementTreeNode *currentElementTreeNode = Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, (int)dfsBuffer.length - 1);
        Clay_LayoutElement *currentElement = currentElementTreeNode->layoutElement;
        if (!context->treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
            context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;
            // If the element has no children or is the container for a text element, don't bother inspecting it
            if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement->childrenOrTextContent.children.length == 0) {
                dfsBuffer.length--;
                continue;
            }
            // Add the children to the DFS buffer (needs to be pushed in reverse so that stack traversal is in correct layout order)
            for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; i++) {
                context->treeNodeVisited.internalArray[dfsBuffer.length] = false;
                Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]) });
            }
            continue;
        }
        dfsBuffer.length--;

        // DFS node has been visited, this is on the way back up to the root
        Clay_LayoutConfig *layoutConfig = currentElement->layoutConfig;
        if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
            // Resize any parent containers that have grown in height along their non layout axis
            for (int32_t j = 0; j < currentElement->childrenOrTextContent.children.length; ++j) {
                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[j]);
                float childHeightWithPadding = CLAY__MAX(childElement->dimensions.height + layoutConfig->padding.top + layoutConfig->padding.bottom, currentElement->dimensions.height);
                currentElement->dimensions.height = CLAY__MIN(CLAY__MAX(childHeightWithPadding, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
            }
        } else if (layoutConfig->layoutDirection == CLAY_TOP_TO_BOTTOM) {
            // Resizing along the layout axis
            float contentHeight = (float)(layoutConfig->padding.top + layoutConfig->padding.bottom);
            for (int32_t j = 0; j < currentElement->childrenOrTextContent.children.length; ++j) {
                Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[j]);
                contentHeight += childElement->dimensions.height;
            }
            contentHeight += (float)(CLAY__MAX(currentElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
            currentElement->dimensions.height = CLAY__MIN(CLAY__MAX(contentHeight, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
        }
    }

    // Calculate sizing along the Y axis
    Clay__SizeContainersAlongAxis(false);

    // Scale horizontal widths according to aspect ratio
    for (int32_t i = 0; i < context->aspectRatioElementIndexes.length; ++i) {
        Clay_LayoutElement* aspectElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->aspectRatioElementIndexes, i));
        Clay_AspectRatioElementConfig *config = Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig;
        aspectElement->dimensions.width = config->aspectRatio * aspectElement->dimensions.height;
    }

    // Sort tree roots by z-index
    int32_t sortMax = context->layoutElementTreeRoots.length - 1;
    while (sortMax > 0) { // todo dumb bubble sort
        for (int32_t i = 0; i < sortMax; ++i) {
            Clay__LayoutElementTreeRoot current = *Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, i);
            Clay__LayoutElementTreeRoot next = *Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, i + 1);
            if (next.zIndex < current.zIndex) {
                Clay__LayoutElementTreeRootArray_Set(&context->layoutElementTreeRoots, i, next);
                Clay__LayoutElementTreeRootArray_Set(&context->layoutElementTreeRoots, i + 1, current);
            }
        }
        sortMax--;
    }

    // Calculate final positions and generate render commands
    context->renderCommands.length = 0;
    dfsBuffer.length = 0;
    for (int32_t rootIndex = 0; rootIndex < context->layoutElementTreeRoots.length; ++rootIndex) {
        dfsBuffer.length = 0;
        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, rootIndex);
        Clay_LayoutElement *rootElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)root->layoutElementIndex);
        Clay_Vector2 rootPosition = CLAY__DEFAULT_STRUCT;
        Clay_LayoutElementHashMapItem *parentHashMapItem = Clay__GetHashMapItem(root->parentId);
        // Position root floating containers
        if (Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) && parentHashMapItem) {
            Clay_FloatingElementConfig *config = Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
            Clay_Dimensions rootDimensions = rootElement->dimensions;
            Clay_BoundingBox parentBoundingBox = parentHashMapItem->boundingBox;
            // Set X position
            Clay_Vector2 targetAttachPosition = CLAY__DEFAULT_STRUCT;
            switch (config->attachPoints.parent) {
                case CLAY_ATTACH_POINT_LEFT_TOP:
                case CLAY_ATTACH_POINT_LEFT_CENTER:
                case CLAY_ATTACH_POINT_LEFT_BOTTOM: targetAttachPosition.x = parentBoundingBox.x; break;
                case CLAY_ATTACH_POINT_CENTER_TOP:
                case CLAY_ATTACH_POINT_CENTER_CENTER:
                case CLAY_ATTACH_POINT_CENTER_BOTTOM: targetAttachPosition.x = parentBoundingBox.x + (parentBoundingBox.width / 2); break;
                case CLAY_ATTACH_POINT_RIGHT_TOP:
                case CLAY_ATTACH_POINT_RIGHT_CENTER:
                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.x = parentBoundingBox.x + parentBoundingBox.width; break;
            }
            switch (config->attachPoints.element) {
                case CLAY_ATTACH_POINT_LEFT_TOP:
                case CLAY_ATTACH_POINT_LEFT_CENTER:
                case CLAY_ATTACH_POINT_LEFT_BOTTOM: break;
                case CLAY_ATTACH_POINT_CENTER_TOP:
                case CLAY_ATTACH_POINT_CENTER_CENTER:
                case CLAY_ATTACH_POINT_CENTER_BOTTOM: targetAttachPosition.x -= (rootDimensions.width / 2); break;
                case CLAY_ATTACH_POINT_RIGHT_TOP:
                case CLAY_ATTACH_POINT_RIGHT_CENTER:
                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.x -= rootDimensions.width; break;
            }
            switch (config->attachPoints.parent) { // I know I could merge the x and y switch statements, but this is easier to read
                case CLAY_ATTACH_POINT_LEFT_TOP:
                case CLAY_ATTACH_POINT_RIGHT_TOP:
                case CLAY_ATTACH_POINT_CENTER_TOP: targetAttachPosition.y = parentBoundingBox.y; break;
                case CLAY_ATTACH_POINT_LEFT_CENTER:
                case CLAY_ATTACH_POINT_CENTER_CENTER:
                case CLAY_ATTACH_POINT_RIGHT_CENTER: targetAttachPosition.y = parentBoundingBox.y + (parentBoundingBox.height / 2); break;
                case CLAY_ATTACH_POINT_LEFT_BOTTOM:
                case CLAY_ATTACH_POINT_CENTER_BOTTOM:
                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.y = parentBoundingBox.y + parentBoundingBox.height; break;
            }
            switch (config->attachPoints.element) {
                case CLAY_ATTACH_POINT_LEFT_TOP:
                case CLAY_ATTACH_POINT_RIGHT_TOP:
                case CLAY_ATTACH_POINT_CENTER_TOP: break;
                case CLAY_ATTACH_POINT_LEFT_CENTER:
                case CLAY_ATTACH_POINT_CENTER_CENTER:
                case CLAY_ATTACH_POINT_RIGHT_CENTER: targetAttachPosition.y -= (rootDimensions.height / 2); break;
                case CLAY_ATTACH_POINT_LEFT_BOTTOM:
                case CLAY_ATTACH_POINT_CENTER_BOTTOM:
                case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.y -= rootDimensions.height; break;
            }
            targetAttachPosition.x += config->offset.x;
            targetAttachPosition.y += config->offset.y;
            rootPosition = targetAttachPosition;
        }
        if (root->clipElementId) {
            Clay_LayoutElementHashMapItem *clipHashMapItem = Clay__GetHashMapItem(root->clipElementId);
            if (clipHashMapItem) {
                // Floating elements that are attached to scrolling contents won't be correctly positioned if external scroll handling is enabled, fix here
                if (context->externalScrollHandlingEnabled) {
                    Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(clipHashMapItem->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
                    if (clipConfig->horizontal) {
                        rootPosition.x += clipConfig->childOffset.x;
                    }
                    if (clipConfig->vertical) {
                        rootPosition.y += clipConfig->childOffset.y;
                    }
                }
                Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                    .boundingBox = clipHashMapItem->boundingBox,
                    .userData = 0,
                    .id = Clay__HashNumber(rootElement->id, rootElement->childrenOrTextContent.children.length + 10).id, // TODO need a better strategy for managing derived ids
                    .zIndex = root->zIndex,
                    .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START,
                });
            }
        }
        Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = rootElement, .position = rootPosition, .nextChildOffset = { .x = (float)rootElement->layoutConfig->padding.left, .y = (float)rootElement->layoutConfig->padding.top } });

        context->treeNodeVisited.internalArray[0] = false;
        while (dfsBuffer.length > 0) {
            Clay__LayoutElementTreeNode *currentElementTreeNode = Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, (int)dfsBuffer.length - 1);
            Clay_LayoutElement *currentElement = currentElementTreeNode->layoutElement;
            Clay_LayoutConfig *layoutConfig = currentElement->layoutConfig;
            Clay_Vector2 scrollOffset = CLAY__DEFAULT_STRUCT;

            // This will only be run a single time for each element in downwards DFS order
            if (!context->treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
                context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;

                Clay_BoundingBox currentElementBoundingBox = { currentElementTreeNode->position.x, currentElementTreeNode->position.y, currentElement->dimensions.width, currentElement->dimensions.height };
                if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)) {
                    Clay_FloatingElementConfig *floatingElementConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
                    Clay_Dimensions expand = floatingElementConfig->expand;
                    currentElementBoundingBox.x -= expand.width;
                    currentElementBoundingBox.width += expand.width * 2;
                    currentElementBoundingBox.y -= expand.height;
                    currentElementBoundingBox.height += expand.height * 2;
                }

                Clay__ScrollContainerDataInternal *scrollContainerData = CLAY__NULL;
                // Apply scroll offsets to container
                if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP)) {
                    Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;

                    // This linear scan could theoretically be slow under very strange conditions, but I can't imagine a real UI with more than a few 10's of scroll containers
                    for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
                        Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
                        if (mapping->layoutElement == currentElement) {
                            scrollContainerData = mapping;
                            mapping->boundingBox = currentElementBoundingBox;
                            scrollOffset = clipConfig->childOffset;
                            if (context->externalScrollHandlingEnabled) {
                                scrollOffset = CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
                            }
                            break;
                        }
                    }
                }

                Clay_LayoutElementHashMapItem *hashMapItem = Clay__GetHashMapItem(currentElement->id);
                if (hashMapItem) {
                    hashMapItem->boundingBox = currentElementBoundingBox;
                }

                int32_t sortedConfigIndexes[20];
                for (int32_t elementConfigIndex = 0; elementConfigIndex < currentElement->elementConfigs.length; ++elementConfigIndex) {
                    sortedConfigIndexes[elementConfigIndex] = elementConfigIndex;
                }
                sortMax = currentElement->elementConfigs.length - 1;
                while (sortMax > 0) { // todo dumb bubble sort
                    for (int32_t i = 0; i < sortMax; ++i) {
                        int32_t current = sortedConfigIndexes[i];
                        int32_t next = sortedConfigIndexes[i + 1];
                        Clay__ElementConfigType currentType = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, current)->type;
                        Clay__ElementConfigType nextType = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, next)->type;
                        if (nextType == CLAY__ELEMENT_CONFIG_TYPE_CLIP || currentType == CLAY__ELEMENT_CONFIG_TYPE_BORDER) {
                            sortedConfigIndexes[i] = next;
                            sortedConfigIndexes[i + 1] = current;
                        }
                    }
                    sortMax--;
                }

                bool emitRectangle = false;
                // Create the render commands for this element
                Clay_SharedElementConfig *sharedConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig;
                if (sharedConfig && sharedConfig->backgroundColor.a > 0) {
                   emitRectangle = true;
                }
                else if (!sharedConfig) {
                    emitRectangle = false;
                    sharedConfig = &Clay_SharedElementConfig_DEFAULT;
                }
                for (int32_t elementConfigIndex = 0; elementConfigIndex < currentElement->elementConfigs.length; ++elementConfigIndex) {
                    Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, sortedConfigIndexes[elementConfigIndex]);
                    Clay_RenderCommand renderCommand = {
                        .boundingBox = currentElementBoundingBox,
                        .userData = sharedConfig->userData,
                        .id = currentElement->id,
                    };

                    bool offscreen = Clay__ElementIsOffscreen(&currentElementBoundingBox);
                    // Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
                    bool shouldRender = !offscreen;
                    switch (elementConfig->type) {
                        case CLAY__ELEMENT_CONFIG_TYPE_ASPECT:
                        case CLAY__ELEMENT_CONFIG_TYPE_FLOATING:
                        case CLAY__ELEMENT_CONFIG_TYPE_SHARED:
                        case CLAY__ELEMENT_CONFIG_TYPE_BORDER: {
                            shouldRender = false;
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_CLIP: {
                            renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START;
                            renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
                                .clip = {
                                    .horizontal = elementConfig->config.clipElementConfig->horizontal,
                                    .vertical = elementConfig->config.clipElementConfig->vertical,
                                }
                            };
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: {
                            renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_IMAGE;
                            renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
                                .image = {
                                    .backgroundColor = sharedConfig->backgroundColor,
                                    .cornerRadius = sharedConfig->cornerRadius,
                                    .imageData = elementConfig->config.imageElementConfig->imageData,
                               }
                            };
                            emitRectangle = false;
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_TEXT: {
                            if (!shouldRender) {
                                break;
                            }
                            shouldRender = false;
                            Clay_ElementConfigUnion configUnion = elementConfig->config;
                            Clay_TextElementConfig *textElementConfig = configUnion.textElementConfig;
                            float naturalLineHeight = currentElement->childrenOrTextContent.textElementData->preferredDimensions.height;
                            float finalLineHeight = textElementConfig->lineHeight > 0 ? (float)textElementConfig->lineHeight : naturalLineHeight;
                            float lineHeightOffset = (finalLineHeight - naturalLineHeight) / 2;
                            float yPosition = lineHeightOffset;
                            for (int32_t lineIndex = 0; lineIndex < currentElement->childrenOrTextContent.textElementData->wrappedLines.length; ++lineIndex) {
                                Clay__WrappedTextLine *wrappedLine = Clay__WrappedTextLineArraySlice_Get(&currentElement->childrenOrTextContent.textElementData->wrappedLines, lineIndex);
                                if (wrappedLine->line.length == 0) {
                                    yPosition += finalLineHeight;
                                    continue;
                                }
                                float offset = (currentElementBoundingBox.width - wrappedLine->dimensions.width);
                                if (textElementConfig->textAlignment == CLAY_TEXT_ALIGN_LEFT) {
                                    offset = 0;
                                }
                                if (textElementConfig->textAlignment == CLAY_TEXT_ALIGN_CENTER) {
                                    offset /= 2;
                                }
                                Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                                    .boundingBox = { currentElementBoundingBox.x + offset, currentElementBoundingBox.y + yPosition, wrappedLine->dimensions.width, wrappedLine->dimensions.height },
                                    .renderData = { .text = {
                                        .stringContents = CLAY__INIT(Clay_StringSlice) { .length = wrappedLine->line.length, .chars = wrappedLine->line.chars, .baseChars = currentElement->childrenOrTextContent.textElementData->text.chars },
                                        .textColor = textElementConfig->textColor,
                                        .fontId = textElementConfig->fontId,
                                        .fontSize = textElementConfig->fontSize,
                                        .letterSpacing = textElementConfig->letterSpacing,
                                        .lineHeight = textElementConfig->lineHeight,
                                    }},
                                    .userData = textElementConfig->userData,
                                    .id = Clay__HashNumber(lineIndex, currentElement->id).id,
                                    .zIndex = root->zIndex,
                                    .commandType = CLAY_RENDER_COMMAND_TYPE_TEXT,
                                });
                                yPosition += finalLineHeight;

                                if (!context->disableCulling && (currentElementBoundingBox.y + yPosition > context->layoutDimensions.height)) {
                                    break;
                                }
                            }
                            break;
                        }
                        case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM: {
                            renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_CUSTOM;
                            renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
                                .custom = {
                                    .backgroundColor = sharedConfig->backgroundColor,
                                    .cornerRadius = sharedConfig->cornerRadius,
                                    .customData = elementConfig->config.customElementConfig->customData,
                                }
                            };
                            emitRectangle = false;
                            break;
                        }
                        default: break;
                    }
                    if (shouldRender) {
                        Clay__AddRenderCommand(renderCommand);
                    }
                    if (offscreen) {
                        // NOTE: You may be tempted to try an early return / continue if an element is off screen. Why bother calculating layout for its children, right?
                        // Unfortunately, a FLOATING_CONTAINER may be defined that attaches to a child or grandchild of this element, which is large enough to still
                        // be on screen, even if this element isn't. That depends on this element and it's children being laid out correctly (even if they are entirely off screen)
                    }
                }

                if (emitRectangle) {
                    Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                        .boundingBox = currentElementBoundingBox,
                        .renderData = { .rectangle = {
                                .backgroundColor = sharedConfig->backgroundColor,
                                .cornerRadius = sharedConfig->cornerRadius,
                        }},
                        .userData = sharedConfig->userData,
                        .id = currentElement->id,
                        .zIndex = root->zIndex,
                        .commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
                    });
                }

                // Setup initial on-axis alignment
                if (!Clay__ElementHasConfig(currentElementTreeNode->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
                    Clay_Dimensions contentSize = {0,0};
                    if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
                        for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                            Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                            contentSize.width += childElement->dimensions.width;
                            contentSize.height = CLAY__MAX(contentSize.height, childElement->dimensions.height);
                        }
                        contentSize.width += (float)(CLAY__MAX(currentElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
                        float extraSpace = currentElement->dimensions.width - (float)(layoutConfig->padding.left + layoutConfig->padding.right) - contentSize.width;
                        switch (layoutConfig->childAlignment.x) {
                            case CLAY_ALIGN_X_LEFT: extraSpace = 0; break;
                            case CLAY_ALIGN_X_CENTER: extraSpace /= 2; break;
                            default: break;
                        }
                        currentElementTreeNode->nextChildOffset.x += extraSpace;
                        extraSpace = CLAY__MAX(0, extraSpace);
                    } else {
                        for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                            Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                            contentSize.width = CLAY__MAX(contentSize.width, childElement->dimensions.width);
                            contentSize.height += childElement->dimensions.height;
                        }
                        contentSize.height += (float)(CLAY__MAX(currentElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
                        float extraSpace = currentElement->dimensions.height - (float)(layoutConfig->padding.top + layoutConfig->padding.bottom) - contentSize.height;
                        switch (layoutConfig->childAlignment.y) {
                            case CLAY_ALIGN_Y_TOP: extraSpace = 0; break;
                            case CLAY_ALIGN_Y_CENTER: extraSpace /= 2; break;
                            default: break;
                        }
                        extraSpace = CLAY__MAX(0, extraSpace);
                        currentElementTreeNode->nextChildOffset.y += extraSpace;
                    }

                    if (scrollContainerData) {
                        scrollContainerData->contentSize = CLAY__INIT(Clay_Dimensions) { contentSize.width + (float)(layoutConfig->padding.left + layoutConfig->padding.right), contentSize.height + (float)(layoutConfig->padding.top + layoutConfig->padding.bottom) };
                    }
                }
            }
            else {
                // DFS is returning upwards backwards
                bool closeClipElement = false;
                Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
                if (clipConfig) {
                    closeClipElement = true;
                    for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
                        Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
                        if (mapping->layoutElement == currentElement) {
                            scrollOffset = clipConfig->childOffset;
                            if (context->externalScrollHandlingEnabled) {
                                scrollOffset = CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
                            }
                            break;
                        }
                    }
                }

                if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER)) {
                    Clay_LayoutElementHashMapItem *currentElementData = Clay__GetHashMapItem(currentElement->id);
                    Clay_BoundingBox currentElementBoundingBox = currentElementData->boundingBox;

                    // Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
                    if (!Clay__ElementIsOffscreen(&currentElementBoundingBox)) {
                        Clay_SharedElementConfig *sharedConfig = Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED) ? Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig : &Clay_SharedElementConfig_DEFAULT;
                        Clay_BorderElementConfig *borderConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER).borderElementConfig;
                        Clay_RenderCommand renderCommand = {
                                .boundingBox = currentElementBoundingBox,
                                .renderData = { .border = {
                                    .color = borderConfig->color,
                                    .cornerRadius = sharedConfig->cornerRadius,
                                    .width = borderConfig->width
                                }},
                                .userData = sharedConfig->userData,
                                .id = Clay__HashNumber(currentElement->id, currentElement->childrenOrTextContent.children.length).id,
                                .commandType = CLAY_RENDER_COMMAND_TYPE_BORDER,
                        };
                        Clay__AddRenderCommand(renderCommand);
                        if (borderConfig->width.betweenChildren > 0 && borderConfig->color.a > 0) {
                            float halfGap = layoutConfig->childGap / 2;
                            Clay_Vector2 borderOffset = { (float)layoutConfig->padding.left - halfGap, (float)layoutConfig->padding.top - halfGap };
                            if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
                                for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                                    if (i > 0) {
                                        Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                                            .boundingBox = { currentElementBoundingBox.x + borderOffset.x + scrollOffset.x, currentElementBoundingBox.y + scrollOffset.y, (float)borderConfig->width.betweenChildren, currentElement->dimensions.height },
                                            .renderData = { .rectangle = {
                                                .backgroundColor = borderConfig->color,
                                            } },
                                            .userData = sharedConfig->userData,
                                            .id = Clay__HashNumber(currentElement->id, currentElement->childrenOrTextContent.children.length + 1 + i).id,
                                            .commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
                                        });
                                    }
                                    borderOffset.x += (childElement->dimensions.width + (float)layoutConfig->childGap);
                                }
                            } else {
                                for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                                    if (i > 0) {
                                        Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                                            .boundingBox = { currentElementBoundingBox.x + scrollOffset.x, currentElementBoundingBox.y + borderOffset.y + scrollOffset.y, currentElement->dimensions.width, (float)borderConfig->width.betweenChildren },
                                            .renderData = { .rectangle = {
                                                    .backgroundColor = borderConfig->color,
                                            } },
                                            .userData = sharedConfig->userData,
                                            .id = Clay__HashNumber(currentElement->id, currentElement->childrenOrTextContent.children.length + 1 + i).id,
                                            .commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
                                        });
                                    }
                                    borderOffset.y += (childElement->dimensions.height + (float)layoutConfig->childGap);
                                }
                            }
                        }
                    }
                }
                // This exists because the scissor needs to end _after_ borders between elements
                if (closeClipElement) {
                    Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
                        .id = Clay__HashNumber(currentElement->id, rootElement->childrenOrTextContent.children.length + 11).id,
                        .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_END,
                    });
                }

                dfsBuffer.length--;
                continue;
            }

            // Add children to the DFS buffer
            if (!Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
                dfsBuffer.length += currentElement->childrenOrTextContent.children.length;
                for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
                    Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
                    // Alignment along non layout axis
                    if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
                        currentElementTreeNode->nextChildOffset.y = currentElement->layoutConfig->padding.top;
                        float whiteSpaceAroundChild = currentElement->dimensions.height - (float)(layoutConfig->padding.top + layoutConfig->padding.bottom) - childElement->dimensions.height;
                        switch (layoutConfig->childAlignment.y) {
                            case CLAY_ALIGN_Y_TOP: break;
                            case CLAY_ALIGN_Y_CENTER: currentElementTreeNode->nextChildOffset.y += whiteSpaceAroundChild / 2; break;
                            case CLAY_ALIGN_Y_BOTTOM: currentElementTreeNode->nextChildOffset.y += whiteSpaceAroundChild; break;
                        }
                    } else {
                        currentElementTreeNode->nextChildOffset.x = currentElement->layoutConfig->padding.left;
                        float whiteSpaceAroundChild = currentElement->dimensions.width - (float)(layoutConfig->padding.left + layoutConfig->padding.right) - childElement->dimensions.width;
                        switch (layoutConfig->childAlignment.x) {
                            case CLAY_ALIGN_X_LEFT: break;
                            case CLAY_ALIGN_X_CENTER: currentElementTreeNode->nextChildOffset.x += whiteSpaceAroundChild / 2; break;
                            case CLAY_ALIGN_X_RIGHT: currentElementTreeNode->nextChildOffset.x += whiteSpaceAroundChild; break;
                        }
                    }

                    Clay_Vector2 childPosition = {
                        currentElementTreeNode->position.x + currentElementTreeNode->nextChildOffset.x + scrollOffset.x,
                        currentElementTreeNode->position.y + currentElementTreeNode->nextChildOffset.y + scrollOffset.y,
                    };

                    // DFS buffer elements need to be added in reverse because stack traversal happens backwards
                    uint32_t newNodeIndex = dfsBuffer.length - 1 - i;
                    dfsBuffer.internalArray[newNodeIndex] = CLAY__INIT(Clay__LayoutElementTreeNode) {
                        .layoutElement = childElement,
                        .position = { childPosition.x, childPosition.y },
                        .nextChildOffset = { .x = (float)childElement->layoutConfig->padding.left, .y = (float)childElement->layoutConfig->padding.top },
                    };
                    context->treeNodeVisited.internalArray[newNodeIndex] = false;

                    // Update parent offsets
                    if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
                        currentElementTreeNode->nextChildOffset.x += childElement->dimensions.width + (float)layoutConfig->childGap;
                    } else {
                        currentElementTreeNode->nextChildOffset.y += childElement->dimensions.height + (float)layoutConfig->childGap;
                    }
                }
            }
        }

        if (root->clipElementId) {
            Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) { .id = Clay__HashNumber(rootElement->id, rootElement->childrenOrTextContent.children.length + 11).id, .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_END });
        }
    }
}
