package clay

func Clay__SizeContainersAlongAxis(xAxis bool) {
	currentContext := Clay_GetCurrentContext()
	bfsBuffer := currentContext.LayoutElementChildrenBuffer
	resizableContainerBuffer := currentContext.OpenLayoutElementStack

	for rootIndex := int32(0); rootIndex < currentContext.LayoutElementTreeRoots.Length(); rootIndex++ {
		Clay__Array_Reset(&bfsBuffer)

		root := Clay__Array_Get(&currentContext.LayoutElementTreeRoots, rootIndex)
		rootElement := Clay__Array_Get(&currentContext.LayoutElements, root.LayoutElementIndex)
		Clay__Array_Add(&bfsBuffer, root.LayoutElementIndex)

		// Size floating containers to their parents
		if Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) {
			floatingElementConfig := Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).FloatingElementConfig
			parentItem := Clay__GetHashMapItem(floatingElementConfig.ParentId)
			if parentItem != nil && parentItem != &Clay_LayoutElementHashMapItem_DEFAULT {
				parentLayoutElement := parentItem.LayoutElement
				switch rootElement.LayoutConfig.Sizing.Width.Type {
				case CLAY__SIZING_TYPE_GROW:
					{
						rootElement.Dimensions.Width = parentLayoutElement.Dimensions.Width
						break
					}
				case CLAY__SIZING_TYPE_PERCENT:
					{
						rootElement.Dimensions.Width = parentLayoutElement.Dimensions.Width * rootElement.LayoutConfig.Sizing.Width.Size.Percent
						break
					}
				default:
					break
				}
				switch rootElement.LayoutConfig.Sizing.Height.Type {
				case CLAY__SIZING_TYPE_GROW:
					{
						rootElement.Dimensions.Height = parentLayoutElement.Dimensions.Height
						break
					}
				case CLAY__SIZING_TYPE_PERCENT:
					{
						rootElement.Dimensions.Height = parentLayoutElement.Dimensions.Height * rootElement.LayoutConfig.Sizing.Height.Size.Percent
						break
					}
				default:
					break
				}
			}
		}

		if rootElement.LayoutConfig.Sizing.Width.Type != CLAY__SIZING_TYPE_PERCENT {
			rootElement.Dimensions.Width = CLAY__MIN(CLAY__MAX(rootElement.Dimensions.Width, rootElement.LayoutConfig.Sizing.Width.Size.MinMax.Min), rootElement.LayoutConfig.Sizing.Width.Size.MinMax.Max)
		}
		if rootElement.LayoutConfig.Sizing.Height.Type != CLAY__SIZING_TYPE_PERCENT {
			rootElement.Dimensions.Height = CLAY__MIN(CLAY__MAX(rootElement.Dimensions.Height, rootElement.LayoutConfig.Sizing.Height.Size.MinMax.Min), rootElement.LayoutConfig.Sizing.Height.Size.MinMax.Max)
		}

		for i := int32(0); i < bfsBuffer.Length(); i++ {
			parentIndex := Clay__Array_GetValue(&bfsBuffer, i)
			parent := Clay__Array_Get(&currentContext.LayoutElements, parentIndex)
			parentStyleConfig := parent.LayoutConfig
			growContainerCount := 0
			var parentSize float32
			if xAxis {
				parentSize = parent.Dimensions.Width
			} else {
				parentSize = parent.Dimensions.Height
			}
			var parentPadding float32
			if xAxis {
				parentPadding = float32(parent.LayoutConfig.Padding.Left + parent.LayoutConfig.Padding.Right)
			} else {
				parentPadding = float32(parent.LayoutConfig.Padding.Top + parent.LayoutConfig.Padding.Bottom)
			}
			var innerContentSize float32
			totalPaddingAndChildGaps := parentPadding
			sizingAlongAxis := (xAxis && parentStyleConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT) || (!xAxis && parentStyleConfig.LayoutDirection == CLAY_TOP_TO_BOTTOM)
			Clay__Array_Reset(&resizableContainerBuffer)
			parentChildGap := parentStyleConfig.ChildGap

			for childOffset := int32(0); childOffset < int32(parent.ChildrenOrTextContent.Children.Length); childOffset++ {
				childElementIndex := parent.ChildrenOrTextContent.Children.Elements[childOffset]
				childElement := Clay__Array_Get(&currentContext.LayoutElements, childElementIndex)
				var childSizing Clay_SizingAxis
				if xAxis {
					childSizing = childElement.LayoutConfig.Sizing.Width
				} else {
					childSizing = childElement.LayoutConfig.Sizing.Height
				}
				var childSize float32
				if xAxis {
					childSize = childElement.Dimensions.Width
				} else {
					childSize = childElement.Dimensions.Height
				}

				if !Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) && childElement.ChildrenOrTextContent.Children.Length > 0 {
					Clay__Array_Add(&bfsBuffer, childElementIndex)
				}

				if childSizing.Type != CLAY__SIZING_TYPE_PERCENT && childSizing.Type != CLAY__SIZING_TYPE_FIXED && (!Clay__ElementHasConfig(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || (Clay__FindElementConfigWithType(childElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).TextElementConfig.WrapMode == CLAY_TEXT_WRAP_WORDS)) { // todo too many loops
					Clay__Array_Add(&resizableContainerBuffer, childElementIndex)
				}

				if sizingAlongAxis {
					if childSizing.Type == CLAY__SIZING_TYPE_PERCENT {
						innerContentSize += 0
					} else {
						innerContentSize += childSize
					}

					if childSizing.Type == CLAY__SIZING_TYPE_GROW {
						growContainerCount++
					}
					if childOffset > 0 {
						innerContentSize += float32(parentChildGap) // For children after index 0, the childAxisOffset is the gap from the previous child
						totalPaddingAndChildGaps += float32(parentChildGap)
					}
				} else {
					innerContentSize = CLAY__MAX(childSize, innerContentSize)
				}

			}

			// Expand percentage containers to size
			for childOffset := int32(0); childOffset < int32(parent.ChildrenOrTextContent.Children.Length); childOffset++ {
				childElementIndex := parent.ChildrenOrTextContent.Children.Elements[childOffset]
				childElement := Clay__Array_Get(&currentContext.LayoutElements, childElementIndex)
				var childSizing Clay_SizingAxis
				var childSize float32

				if xAxis {
					childSizing = childElement.LayoutConfig.Sizing.Width
					childSize = childElement.Dimensions.Width
				} else {
					childSizing = childElement.LayoutConfig.Sizing.Height
					childSize = childElement.Dimensions.Height

				}

				if childSizing.Type == CLAY__SIZING_TYPE_PERCENT {
					childSize = (parentSize - totalPaddingAndChildGaps) * childSizing.Size.Percent
					if xAxis {
						childElement.Dimensions.Width = childSize
					} else {
						childElement.Dimensions.Height = childSize
					}
					if sizingAlongAxis {
						innerContentSize += childSize
					}
					Clay__UpdateAspectRatioBox(childElement)
				}
			}
			if sizingAlongAxis {
				sizeToDistribute := parentSize - parentPadding - innerContentSize
				// The content is too large, compress the children as much as possible
				if sizeToDistribute < 0 {
					// If the parent clips content in this axis direction, don't compress children, just leave them alone
					clipElementConfig := Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP).ClipElementConfig

					if clipElementConfig != nil {
						if (xAxis && clipElementConfig.Horizontal) || (!xAxis && clipElementConfig.Vertical) {
							continue
						}
					}
					// Scrolling containers preferentially compress before others
					for sizeToDistribute < -CLAY__EPSILON && resizableContainerBuffer.Length() > 0 {
						var largest float32 = 0
						var secondLargest float32 = 0
						var widthToAdd float32 = sizeToDistribute
						for childIndex := int32(0); childIndex < resizableContainerBuffer.Length(); childIndex++ {
							child := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&resizableContainerBuffer, childIndex))
							var childSize float32
							if xAxis {
								childSize = child.Dimensions.Width
							} else {
								childSize = child.Dimensions.Height
							}

							if Clay__FloatEqual(childSize, largest) {
								continue
							}
							if childSize > largest {
								secondLargest = largest
								largest = childSize
							}
							if childSize < largest {
								secondLargest = CLAY__MAX(secondLargest, childSize)
								widthToAdd = secondLargest - largest
							}
						}

						widthToAdd = CLAY__MAX(widthToAdd, sizeToDistribute/float32(resizableContainerBuffer.Length()))

						for childIndex := int32(0); childIndex < resizableContainerBuffer.Length(); childIndex++ {
							child := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&resizableContainerBuffer, childIndex))
							var childSize float32
							var minSize float32
							if xAxis {
								childSize = child.Dimensions.Width
								minSize = child.MinDimensions.Width
							} else {
								childSize = child.Dimensions.Height
								minSize = child.MinDimensions.Height
							}

							previousWidth := childSize
							if Clay__FloatEqual(childSize, largest) {
								childSize += widthToAdd
								if childSize <= minSize {
									childSize = minSize
									Clay__Array_RemoveSwapback(&resizableContainerBuffer, childIndex)
									childIndex--
								}
								if xAxis {
									child.Dimensions.Width = childSize
								} else {
									child.Dimensions.Height = childSize
								}
								sizeToDistribute -= (childSize - previousWidth)
							}
						}
					}
					// The content is too small, allow SIZING_GROW containers to expand
				} else if sizeToDistribute > 0 && growContainerCount > 0 {
					for childIndex := int32(0); childIndex < resizableContainerBuffer.Length(); childIndex++ {
						child := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&resizableContainerBuffer, childIndex))

						var childSizing Clay__SizingType
						if xAxis {
							childSizing = child.LayoutConfig.Sizing.Width.Type
						} else {
							childSizing = child.LayoutConfig.Sizing.Height.Type
						}
						if childSizing != CLAY__SIZING_TYPE_GROW {
							Clay__Array_RemoveSwapback(&resizableContainerBuffer, childIndex)
							childIndex--
						}
					}
					for sizeToDistribute > CLAY__EPSILON && resizableContainerBuffer.Length() > 0 {
						var smallest float32 = CLAY__MAXFLOAT
						var secondSmallest float32 = CLAY__MAXFLOAT
						widthToAdd := sizeToDistribute
						for childIndex := int32(0); childIndex < resizableContainerBuffer.Length(); childIndex++ {
							child := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&resizableContainerBuffer, childIndex))
							var childSize float32
							if xAxis {
								childSize = child.Dimensions.Width
							} else {
								childSize = child.Dimensions.Height
							}
							if Clay__FloatEqual(childSize, smallest) {
								continue
							}
							if childSize < smallest {
								secondSmallest = smallest
								smallest = childSize
							}
							if childSize > smallest {
								secondSmallest = CLAY__MIN(secondSmallest, childSize)
								widthToAdd = secondSmallest - smallest
							}
						}

						widthToAdd = CLAY__MIN(widthToAdd, sizeToDistribute/float32(resizableContainerBuffer.Length()))

						for childIndex := int32(0); childIndex < resizableContainerBuffer.Length(); childIndex++ {
							child := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&resizableContainerBuffer, childIndex))
							var childSize float32
							var maxSize float32
							if xAxis {
								childSize = child.Dimensions.Width
								maxSize = child.LayoutConfig.Sizing.Width.Size.MinMax.Max
							} else {
								childSize = child.Dimensions.Height
								maxSize = child.LayoutConfig.Sizing.Height.Size.MinMax.Max
							}
							previousWidth := childSize
							if Clay__FloatEqual(childSize, smallest) {
								childSize += widthToAdd
								if childSize >= maxSize {
									childSize = maxSize
									Clay__Array_RemoveSwapback(&resizableContainerBuffer, childIndex)
									childIndex--
								}
								if xAxis {
									child.Dimensions.Width = childSize
								} else {
									child.Dimensions.Height = childSize
								}
								sizeToDistribute -= (childSize - previousWidth)
							}
						}
					}
				}
				// Sizing along the non layout axis ("off axis")
			} else {
				for childOffset := int32(0); childOffset < resizableContainerBuffer.Length(); childOffset++ {

					childElement := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&resizableContainerBuffer, childOffset))
					var childSizing Clay_SizingAxis
					var minSize float32
					var childSize float32
					if xAxis {
						childSizing = childElement.LayoutConfig.Sizing.Width
						minSize = childElement.MinDimensions.Width
						childSize = childElement.Dimensions.Width
					} else {
						childSizing = childElement.LayoutConfig.Sizing.Height
						minSize = childElement.MinDimensions.Height
						childSize = childElement.Dimensions.Height
					}

					var maxSize float32 = parentSize - parentPadding
					// If we're laying out the children of a scroll panel, grow containers expand to the size of the inner content, not the outer container
					if Clay__ElementHasConfig(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP) {
						clipElementConfig := Clay__FindElementConfigWithType(parent, CLAY__ELEMENT_CONFIG_TYPE_CLIP).ClipElementConfig
						if (xAxis && clipElementConfig.Horizontal) || (!xAxis && clipElementConfig.Vertical) {
							maxSize = CLAY__MAX(maxSize, innerContentSize)
						}
					}
					if childSizing.Type == CLAY__SIZING_TYPE_GROW {
						childSize = CLAY__MIN(maxSize, childSizing.Size.MinMax.Max)
					}
					childSize = CLAY__MAX(minSize, CLAY__MIN(childSize, maxSize))
					if xAxis {
						childElement.Dimensions.Width = childSize
					} else {
						childElement.Dimensions.Height = childSize
					}
				}
			}

		}

	}
}
