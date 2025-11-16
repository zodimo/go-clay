package clay

type Clay_PointerDataInteractionState uint8

const (
	CLAY_POINTER_DATA_PRESSED_THIS_FRAME Clay_PointerDataInteractionState = iota
	CLAY_POINTER_DATA_PRESSED
	CLAY_POINTER_DATA_RELEASED_THIS_FRAME
	CLAY_POINTER_DATA_RELEASED
)

// Information on the current state of pointer interactions this frame.
type Clay_PointerData struct {
	// The position of the mouse / touch / pointer relative to the root of the layout.
	Position Clay_Vector2
	// Represents the current state of interaction with clay this frame.
	// CLAY_POINTER_DATA_PRESSED_THIS_FRAME - A left mouse click, or touch occurred this frame.
	// CLAY_POINTER_DATA_PRESSED - The left mouse button click or touch happened at some point in the past, and is still currently held down this frame.
	// CLAY_POINTER_DATA_RELEASED_THIS_FRAME - The left mouse button click or touch was released this frame.
	// CLAY_POINTER_DATA_RELEASED - The left mouse button click or touch is not currently down / was released at some point in the past.
	State Clay_PointerDataInteractionState
}

// Controls how mouse pointer events like hover and click are captured or passed through to elements underneath a floating element.
type Clay_PointerCaptureMode uint8

const (
	// (default) "Capture" the pointer event and don't allow events like hover and click to pass through to elements underneath.
	CLAY_POINTER_CAPTURE_MODE_CAPTURE Clay_PointerCaptureMode = iota
	// CLAY_POINTER_CAPTURE_MODE_PARENT, TODO pass pointer through to attached parent
	// Transparently pass through pointer events like hover and click to elements underneath the floating element.
	CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH
)

func Clay__PointIsInsideRect(point Clay_Vector2, rect Clay_BoundingBox) bool {
	return point.X >= rect.X && point.X <= rect.X+rect.Width && point.Y >= rect.Y && point.Y <= rect.Y+rect.Height
}

func Clay_SetPointerState(position Clay_Vector2, isPointerDown bool) {
	context := Clay_GetCurrentContext()
	if context.BooleanWarnings.MaxElementsExceeded {
		return
	}
	context.PointerInfo.Position = position
	Clay__Array_Reset(&context.PointerOverIds)
	dfsBuffer := context.LayoutElementChildrenBuffer
	for rootIndex := int32(context.LayoutElementTreeRoots.Length() - 1); rootIndex >= 0; rootIndex-- {
		Clay__Array_Reset(&dfsBuffer)
		root := Clay__Array_Get(&context.LayoutElementTreeRoots, rootIndex)
		Clay__Array_Add(&dfsBuffer, root.LayoutElementIndex)
		Clay__Array_Set(&context.TreeNodeVisited, 0, false)
		found := false
		for dfsBuffer.Length() > 0 {
			if Clay__Array_GetValue(&context.TreeNodeVisited, dfsBuffer.Length()-1) {
				Clay__Array_Shrink(&dfsBuffer, 1)
				continue
			}
			Clay__Array_Set(&context.TreeNodeVisited, dfsBuffer.Length()-1, true)
			currentElement := Clay__Array_Get(&context.LayoutElements, Clay__Array_GetValue(&dfsBuffer, dfsBuffer.Length()-1))
			mapItem := Clay__GetHashMapItem(currentElement.Id)
			clipElementIdIndex := Clay__Array_IndexOf(&context.LayoutElements, currentElement)              // TODO think of a way around this, maybe the fact that it's essentially a binary tree limits the cost, but the worst case is not great
			clipElementId := Clay__Array_GetValue(&context.LayoutElementClipElementIds, clipElementIdIndex) // pointer arithmetic
			clipItem := Clay__GetHashMapItem(uint32(clipElementId))
			if mapItem != nil {
				elementBox := mapItem.BoundingBox
				elementBox.X -= root.PointerOffset.X
				elementBox.Y -= root.PointerOffset.Y
				if (Clay__PointIsInsideRect(position, elementBox)) && (clipElementId == 0 || (Clay__PointIsInsideRect(position, clipItem.BoundingBox)) || context.ExternalScrollHandlingEnabled) {
					if mapItem.OnHoverFunction != nil {
						mapItem.OnHoverFunction(mapItem.ElementId, context.PointerInfo, mapItem.HoverFunctionUserData)
					}
					Clay__Array_Add(&context.PointerOverIds, mapItem.ElementId)
					found = true
				}
				if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) {
					Clay__Array_Shrink(&dfsBuffer, 1)
					continue
				}
				for i := int32(currentElement.ChildrenOrTextContent.Children.Length) - 1; i >= 0; i-- {
					// Ensure TreeNodeVisited is large enough before adding to dfsBuffer
					if context.TreeNodeVisited.Length() < dfsBuffer.Length() {
						panic("treeNodeVisited[] is not the same length as dfsBuffer")
					}
					if context.TreeNodeVisited.Length() == dfsBuffer.Length() {
						Clay__Array_Add(&context.TreeNodeVisited, false)
					} else {
						Clay__Array_Set(&context.TreeNodeVisited, dfsBuffer.Length(), false)
					}
					Clay__Array_Add(&dfsBuffer, currentElement.ChildrenOrTextContent.Children.Elements[i])
				}
			} else {
				Clay__Array_Shrink(&dfsBuffer, 1)
			}
		}

		rootElement := Clay__Array_Get(&context.LayoutElements, root.LayoutElementIndex)
		if found && Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) &&
			Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).FloatingElementConfig.PointerCaptureMode == CLAY_POINTER_CAPTURE_MODE_CAPTURE {
			break
		}
	}

	if isPointerDown {
		if context.PointerInfo.State == CLAY_POINTER_DATA_PRESSED_THIS_FRAME {
			context.PointerInfo.State = CLAY_POINTER_DATA_PRESSED
		} else if context.PointerInfo.State != CLAY_POINTER_DATA_PRESSED {
			context.PointerInfo.State = CLAY_POINTER_DATA_PRESSED_THIS_FRAME
		}
	} else {
		if context.PointerInfo.State == CLAY_POINTER_DATA_RELEASED_THIS_FRAME {
			context.PointerInfo.State = CLAY_POINTER_DATA_RELEASED
		} else if context.PointerInfo.State != CLAY_POINTER_DATA_RELEASED {
			context.PointerInfo.State = CLAY_POINTER_DATA_RELEASED_THIS_FRAME
		}
	}
}

func Clay_PointerOver(elementId Clay_ElementId) bool { // TODO return priority for separating multiple results
	context := Clay_GetCurrentContext()
	for i := int32(0); i < context.PointerOverIds.Length(); i++ {
		if Clay__Array_GetValue(&context.PointerOverIds, i).Id == elementId.Id {
			return true
		}
	}
	return false
}

func Clay_OnHover(onHoverFunction Clay_OnHoverFunction, userData any) {
	context := Clay_GetCurrentContext()
	if context.BooleanWarnings.MaxElementsExceeded {
		return
	}
	openLayoutElement := Clay__GetOpenLayoutElement()
	if openLayoutElement.Id == 0 {
		Clay__GenerateIdForAnonymousElement(openLayoutElement)
	}
	hashMapItem := Clay__GetHashMapItem(openLayoutElement.Id)
	hashMapItem.OnHoverFunction = onHoverFunction
	hashMapItem.HoverFunctionUserData = userData
}

func Clay_Hovered() bool {
	context := Clay_GetCurrentContext()
	if context.BooleanWarnings.MaxElementsExceeded {
		return false
	}
	openLayoutElement := Clay__GetOpenLayoutElement()
	// If the element has no id attached at this point, we need to generate one
	if openLayoutElement.Id == 0 {
		Clay__GenerateIdForAnonymousElement(openLayoutElement)
	}
	for i := int32(0); i < context.PointerOverIds.Length(); i++ {
		if Clay__Array_Get(&context.PointerOverIds, i).Id == openLayoutElement.Id {
			return true
		}
	}
	return false
}
