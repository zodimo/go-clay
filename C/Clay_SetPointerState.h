void Clay_SetPointerState(Clay_Vector2 position, bool isPointerDown) {
    Clay_Context* context = Clay_GetCurrentContext();
    if (context->booleanWarnings.maxElementsExceeded) {
        return;
    }
    context->pointerInfo.position = position;
    context->pointerOverIds.length = 0;
    Clay__int32_tArray dfsBuffer = context->layoutElementChildrenBuffer;
    for (int32_t rootIndex = context->layoutElementTreeRoots.length - 1; rootIndex >= 0; --rootIndex) {
        dfsBuffer.length = 0;
        Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, rootIndex);
        Clay__int32_tArray_Add(&dfsBuffer, (int32_t)root->layoutElementIndex);
        context->treeNodeVisited.internalArray[0] = false;
        bool found = false;
        while (dfsBuffer.length > 0) {
            if (context->treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
                dfsBuffer.length--;
                continue;
            }
            context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;
            Clay_LayoutElement *currentElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&dfsBuffer, (int)dfsBuffer.length - 1));
            Clay_LayoutElementHashMapItem *mapItem = Clay__GetHashMapItem(currentElement->id); // TODO think of a way around this, maybe the fact that it's essentially a binary tree limits the cost, but the worst case is not great
            int32_t clipElementId = Clay__int32_tArray_GetValue(
                &context->layoutElementClipElementIds, 
                (int32_t)(currentElement - context->layoutElements.internalArray)
            );
            Clay_LayoutElementHashMapItem *clipItem = Clay__GetHashMapItem(clipElementId);
            if (mapItem) {
                Clay_BoundingBox elementBox = mapItem->boundingBox;
                elementBox.x -= root->pointerOffset.x;
                elementBox.y -= root->pointerOffset.y;
                if ((Clay__PointIsInsideRect(position, elementBox)) && (clipElementId == 0 || (Clay__PointIsInsideRect(position, clipItem->boundingBox)) || context->externalScrollHandlingEnabled)) {
                    if (mapItem->onHoverFunction) {
                        mapItem->onHoverFunction(mapItem->elementId, context->pointerInfo, mapItem->hoverFunctionUserData);
                    }
                    Clay_ElementIdArray_Add(&context->pointerOverIds, mapItem->elementId);
                    found = true;
                }
                if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
                    dfsBuffer.length--;
                    continue;
                }
                for (int32_t i = currentElement->childrenOrTextContent.children.length - 1; i >= 0; --i) {
                    Clay__int32_tArray_Add(&dfsBuffer, currentElement->childrenOrTextContent.children.elements[i]);
                    context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = false; // TODO needs to be ranged checked
                }
            } else {
                dfsBuffer.length--;
            }
        }

        Clay_LayoutElement *rootElement = Clay_LayoutElementArray_Get(&context->layoutElements, root->layoutElementIndex);
        if (found && Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) &&
                Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig->pointerCaptureMode == CLAY_POINTER_CAPTURE_MODE_CAPTURE) {
            break;
        }
    }

    if (isPointerDown) {
        if (context->pointerInfo.state == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
            context->pointerInfo.state = CLAY_POINTER_DATA_PRESSED;
        } else if (context->pointerInfo.state != CLAY_POINTER_DATA_PRESSED) {
            context->pointerInfo.state = CLAY_POINTER_DATA_PRESSED_THIS_FRAME;
        }
    } else {
        if (context->pointerInfo.state == CLAY_POINTER_DATA_RELEASED_THIS_FRAME) {
            context->pointerInfo.state = CLAY_POINTER_DATA_RELEASED;
        } else if (context->pointerInfo.state != CLAY_POINTER_DATA_RELEASED)  {
            context->pointerInfo.state = CLAY_POINTER_DATA_RELEASED_THIS_FRAME;
        }
    }
}