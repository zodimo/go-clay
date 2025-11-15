package clay

func Clay_GetScrollOffset() Clay_Vector2 {
	context := Clay_GetCurrentContext()
	if context.BooleanWarnings.MaxElementsExceeded {
		return Clay_Vector2{}
	}
	openLayoutElement := Clay__GetOpenLayoutElement()
	// If the element has no id attached at this point, we need to generate one
	if openLayoutElement.Id == 0 {
		Clay__GenerateIdForAnonymousElement(openLayoutElement)
	}
	for i := int32(0); i < context.ScrollContainerDatas.Length(); i++ {
		mapping := Clay__Array_Get(&context.ScrollContainerDatas, i)
		if mapping.LayoutElement == openLayoutElement {
			return mapping.ScrollPosition
		}
	}
	return Clay_Vector2{0, 0}
}

func Clay_UpdateScrollContainers(enableDragScrolling bool, scrollDelta Clay_Vector2, deltaTime float32) {
	context := Clay_GetCurrentContext()
	isPointerActive := enableDragScrolling && (context.PointerInfo.State == CLAY_POINTER_DATA_PRESSED || context.PointerInfo.State == CLAY_POINTER_DATA_PRESSED_THIS_FRAME)
	// Don't apply scroll events to ancestors of the inner element
	highestPriorityElementIndex := int32(-1)
	var highestPriorityScrollData *Clay__ScrollContainerDataInternal = nil
	for i := int32(0); i < context.ScrollContainerDatas.Length(); i++ {
		scrollData := Clay__Array_Get(&context.ScrollContainerDatas, i)
		if !scrollData.OpenThisFrame {
			Clay__Array_RemoveSwapback(&context.ScrollContainerDatas, i)
			continue
		}
		scrollData.OpenThisFrame = false
		hashMapItem := Clay__GetHashMapItem(scrollData.ElementId)
		// Element isn't rendered this frame but scroll offset has been retained
		if hashMapItem == nil {
			Clay__Array_RemoveSwapback(&context.ScrollContainerDatas, i)
			continue
		}

		// Touch / click is released
		if !isPointerActive && scrollData.PointerScrollActive {
			xDiff := scrollData.ScrollPosition.X - scrollData.ScrollOrigin.X
			if xDiff < -10 || xDiff > 10 {
				scrollData.ScrollMomentum.X = (scrollData.ScrollPosition.X - scrollData.ScrollOrigin.X) / (scrollData.MomentumTime * 25)
			}
			yDiff := scrollData.ScrollPosition.Y - scrollData.ScrollOrigin.Y
			if yDiff < -10 || yDiff > 10 {
				scrollData.ScrollMomentum.Y = (scrollData.ScrollPosition.Y - scrollData.ScrollOrigin.Y) / (scrollData.MomentumTime * 25)
			}
			scrollData.PointerScrollActive = false

			scrollData.PointerOrigin = Clay_Vector2{0, 0}
			scrollData.ScrollOrigin = Clay_Vector2{0, 0}
			scrollData.MomentumTime = 0
		}

		// Apply existing momentum
		scrollData.ScrollPosition.X += scrollData.ScrollMomentum.X
		scrollData.ScrollMomentum.X *= 0.95
		scrollOccurred := scrollDelta.X != 0 || scrollDelta.Y != 0
		if (scrollData.ScrollMomentum.X > -0.1 && scrollData.ScrollMomentum.X < 0.1) || scrollOccurred {
			scrollData.ScrollMomentum.X = 0
		}
		scrollData.ScrollPosition.X = CLAY__MIN(CLAY__MAX(scrollData.ScrollPosition.X, -(CLAY__MAX(scrollData.ContentSize.Width-scrollData.LayoutElement.Dimensions.Width, 0))), 0)

		scrollData.ScrollPosition.Y += scrollData.ScrollMomentum.Y
		scrollData.ScrollMomentum.Y *= 0.95
		if (scrollData.ScrollMomentum.Y > -0.1 && scrollData.ScrollMomentum.Y < 0.1) || scrollOccurred {
			scrollData.ScrollMomentum.Y = 0
		}
		scrollData.ScrollPosition.Y = CLAY__MIN(CLAY__MAX(scrollData.ScrollPosition.Y, -(CLAY__MAX(scrollData.ContentSize.Height-scrollData.LayoutElement.Dimensions.Height, 0))), 0)

		for j := int32(0); j < context.PointerOverIds.Length(); j++ { // TODO n & m are small here but this being n*m gives me the creeps
			if scrollData.LayoutElement.Id == Clay__Array_Get(&context.PointerOverIds, j).Id {
				highestPriorityElementIndex = j
				highestPriorityScrollData = scrollData
			}
		}
	}

	if highestPriorityElementIndex > -1 && highestPriorityScrollData != nil {
		scrollElement := highestPriorityScrollData.LayoutElement
		clipConfig := Clay__FindElementConfigWithType(scrollElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).ClipElementConfig
		canScrollVertically := clipConfig.Vertical && highestPriorityScrollData.ContentSize.Height > scrollElement.Dimensions.Height
		canScrollHorizontally := clipConfig.Horizontal && highestPriorityScrollData.ContentSize.Width > scrollElement.Dimensions.Width
		// Handle wheel scroll
		if canScrollVertically {
			highestPriorityScrollData.ScrollPosition.Y = highestPriorityScrollData.ScrollPosition.Y + scrollDelta.Y*10
		}
		if canScrollHorizontally {
			highestPriorityScrollData.ScrollPosition.X = highestPriorityScrollData.ScrollPosition.X + scrollDelta.X*10
		}
		// Handle click / touch scroll
		if isPointerActive {
			highestPriorityScrollData.ScrollMomentum = Clay_Vector2{0, 0}
			if !highestPriorityScrollData.PointerScrollActive {
				highestPriorityScrollData.PointerOrigin = context.PointerInfo.Position
				highestPriorityScrollData.ScrollOrigin = highestPriorityScrollData.ScrollPosition
				highestPriorityScrollData.PointerScrollActive = true
			} else {
				scrollDeltaX := 0.0
				scrollDeltaY := 0.0
				if canScrollHorizontally {
					oldXScrollPosition := highestPriorityScrollData.ScrollPosition.X
					highestPriorityScrollData.ScrollPosition.X = highestPriorityScrollData.ScrollOrigin.X + (context.PointerInfo.Position.X - highestPriorityScrollData.PointerOrigin.X)
					highestPriorityScrollData.ScrollPosition.X = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.ScrollPosition.X, 0), -(highestPriorityScrollData.ContentSize.Width - highestPriorityScrollData.BoundingBox.Width))
					scrollDeltaX = float64(highestPriorityScrollData.ScrollPosition.X - oldXScrollPosition)
				}
				if canScrollVertically {
					oldYScrollPosition := highestPriorityScrollData.ScrollPosition.Y
					highestPriorityScrollData.ScrollPosition.Y = highestPriorityScrollData.ScrollOrigin.Y + (context.PointerInfo.Position.Y - highestPriorityScrollData.PointerOrigin.Y)
					highestPriorityScrollData.ScrollPosition.Y = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.ScrollPosition.Y, 0), -(highestPriorityScrollData.ContentSize.Height - highestPriorityScrollData.BoundingBox.Height))
					scrollDeltaY = float64(highestPriorityScrollData.ScrollPosition.Y - oldYScrollPosition)
				}
				if scrollDeltaX > -0.1 && scrollDeltaX < 0.1 && scrollDeltaY > -0.1 && scrollDeltaY < 0.1 && highestPriorityScrollData.MomentumTime > 0.15 {
					highestPriorityScrollData.MomentumTime = 0
					highestPriorityScrollData.PointerOrigin = context.PointerInfo.Position
					highestPriorityScrollData.ScrollOrigin = highestPriorityScrollData.ScrollPosition
				} else {
					highestPriorityScrollData.MomentumTime += deltaTime
				}
			}
		}
		// Clamp any changes to scroll position to the maximum size of the contents
		if canScrollVertically {
			highestPriorityScrollData.ScrollPosition.Y = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.ScrollPosition.Y, 0), -(highestPriorityScrollData.ContentSize.Height - scrollElement.Dimensions.Height))
		}
		if canScrollHorizontally {
			highestPriorityScrollData.ScrollPosition.X = CLAY__MAX(CLAY__MIN(highestPriorityScrollData.ScrollPosition.X, 0), -(highestPriorityScrollData.ContentSize.Width - scrollElement.Dimensions.Width))
		}
	}
}
