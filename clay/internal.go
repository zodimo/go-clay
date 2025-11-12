package clay

import (
	"fmt"

	"github.com/zodimo/clay-go/pkg/mem"
)

type Clay__ScrollContainerDataInternal struct {
	LayoutElement       *Clay_LayoutElement
	BoundingBox         Clay_BoundingBox
	ContentSize         Clay_Dimensions
	ScrollOrigin        Clay_Vector2
	PointerOrigin       Clay_Vector2
	ScrollMomentum      Clay_Vector2
	ScrollPosition      Clay_Vector2
	PreviousDelta       Clay_Vector2
	MomentumTime        float32
	ElementId           uint32
	OpenThisFrame       bool
	PointerScrollActive bool
}

type Clay__LayoutElementTreeRoot struct {
	LayoutElementIndex int32
	ParentId           uint32 // This can be zero in the case of the root layout tree
	ClipElementId      uint32 // This can be zero if there is no clip element
	ZIndex             int16
	PointerOffset      Clay_Vector2 // Only used when scroll containers are managed externally
}

func Clay__Array_Allocate_Arena[T any](capacity int32, arena *Clay_Arena) Clay__Array[T] {
	// var zero T
	// typeT := reflect.TypeOf(zero).String()
	// fmt.Println("typeT", typeT, "capacity", capacity)
	return mem.NewMemArray[T](capacity, mem.MemArrayWithArena(arena))
}

type Clay__LayoutElementChildren struct {
	Elements []int32 // *elements
	Length   uint16
}

type Clay__LayoutElementChildrenOrTextContent struct {
	Children        Clay__LayoutElementChildren
	TextElementData *Clay__TextElementData
}

type Clay__LayoutElementTreeNode struct {
	LayoutElement   *Clay_LayoutElement
	Position        Clay_Vector2
	NextChildOffset Clay_Vector2
}

func Clay__CloseElement() {

	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return
	}
	openLayoutElement := Clay__GetOpenLayoutElement()
	layoutConfig := openLayoutElement.LayoutConfig
	if layoutConfig == nil {
		openLayoutElement.LayoutConfig = &Clay_LayoutConfig_DEFAULT
		layoutConfig = &Clay_LayoutConfig_DEFAULT
	}

	elementHasClipHorizontal := false
	elementHasClipVertical := false

	for i := int32(0); i < openLayoutElement.ElementConfigs.Length; i++ {
		config := Clay__Slice_Get(&openLayoutElement.ElementConfigs, i)
		if config.Type == CLAY__ELEMENT_CONFIG_TYPE_CLIP {
			elementHasClipHorizontal = config.Config.ClipElementConfig.Horizontal
			elementHasClipVertical = config.Config.ClipElementConfig.Vertical
			currentContext.OpenClipElementStack.Length--
			break
		} else if config.Type == CLAY__ELEMENT_CONFIG_TYPE_FLOATING {
			currentContext.OpenClipElementStack.Length--
		}
	}

	leftRightPadding := float32(layoutConfig.Padding.Left + layoutConfig.Padding.Right)
	topBottomPadding := float32(layoutConfig.Padding.Top + layoutConfig.Padding.Bottom)

	// Attach children to the current open element

	//attach to the unallocated slice at the end of the array from length to capacity
	openLayoutElement.ChildrenOrTextContent.Children.Elements = mem.MArray_GetSlice(&currentContext.LayoutElementChildren, currentContext.LayoutElementChildren.Length, currentContext.LayoutElementChildren.Capacity) //[currentContext.LayoutElementChildren.Length]
	if layoutConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT {
		openLayoutElement.Dimensions.Width = leftRightPadding
		openLayoutElement.MinDimensions.Width = leftRightPadding
		for i := uint16(0); i < openLayoutElement.ChildrenOrTextContent.Children.Length; i++ {
			childIndex := Clay__Array_GetValue(&currentContext.LayoutElementChildrenBuffer, currentContext.LayoutElementChildrenBuffer.Length-int32(openLayoutElement.ChildrenOrTextContent.Children.Length)+int32(i))
			child := Clay__Array_Get(&currentContext.LayoutElements, childIndex)
			openLayoutElement.Dimensions.Width += child.Dimensions.Width
			openLayoutElement.Dimensions.Height = CLAY__MAX(openLayoutElement.Dimensions.Height, child.Dimensions.Height+topBottomPadding)

			// Minimum size of child elements doesn't matter to clip containers as they can shrink and hide their contents
			if !elementHasClipHorizontal {
				openLayoutElement.MinDimensions.Width += child.MinDimensions.Width
			}
			if !elementHasClipVertical {
				openLayoutElement.MinDimensions.Height = CLAY__MAX(openLayoutElement.MinDimensions.Height, child.MinDimensions.Height+topBottomPadding)
			}
			Clay__Array_Add(&currentContext.LayoutElementChildren, childIndex)
		}

		childGap := float32(CLAY__MAX(float32(openLayoutElement.ChildrenOrTextContent.Children.Length-1), 0) * float32(layoutConfig.ChildGap))

		openLayoutElement.Dimensions.Width += childGap
		if !elementHasClipHorizontal {
			openLayoutElement.MinDimensions.Width += childGap
		}
	} else if layoutConfig.LayoutDirection == CLAY_TOP_TO_BOTTOM {
		openLayoutElement.Dimensions.Height = topBottomPadding
		openLayoutElement.MinDimensions.Height = topBottomPadding
		for i := uint16(0); i < openLayoutElement.ChildrenOrTextContent.Children.Length; i++ {
			childIndex := Clay__Array_GetValue(&currentContext.LayoutElementChildrenBuffer, currentContext.LayoutElementChildrenBuffer.Length-int32(openLayoutElement.ChildrenOrTextContent.Children.Length)+int32(i))
			child := Clay__Array_Get(&currentContext.LayoutElements, childIndex)
			openLayoutElement.Dimensions.Height += child.Dimensions.Height
			openLayoutElement.Dimensions.Width = CLAY__MAX(openLayoutElement.Dimensions.Width, child.Dimensions.Width+leftRightPadding)
			if !elementHasClipVertical {
				openLayoutElement.MinDimensions.Height += child.MinDimensions.Height
			}
			if !elementHasClipHorizontal {
				openLayoutElement.MinDimensions.Width = CLAY__MAX(openLayoutElement.MinDimensions.Width, child.MinDimensions.Width+leftRightPadding)
			}
			Clay__Array_Add(&currentContext.LayoutElementChildren, childIndex)
		}

		childGap := float32(CLAY__MAX(float32(openLayoutElement.ChildrenOrTextContent.Children.Length-1), 0) * float32(layoutConfig.ChildGap))

		openLayoutElement.Dimensions.Height += childGap
		if !elementHasClipVertical {
			openLayoutElement.MinDimensions.Height += childGap
		}
	}

	currentContext.LayoutElementChildrenBuffer.Length -= int32(openLayoutElement.ChildrenOrTextContent.Children.Length)

	// Clamp element min and max width to the values configured in the layout
	if layoutConfig.Sizing.Width.Type != CLAY__SIZING_TYPE_PERCENT {
		if layoutConfig.Sizing.Width.Size.MinMax.Max <= 0 { // Set the max size if the user didn't specify, makes calculations easier
			layoutConfig.Sizing.Width.Size.MinMax.Max = CLAY__MAXFLOAT
		}
		openLayoutElement.Dimensions.Width = CLAY__MIN(CLAY__MAX(openLayoutElement.Dimensions.Width, layoutConfig.Sizing.Width.Size.MinMax.Min), layoutConfig.Sizing.Width.Size.MinMax.Max)
		openLayoutElement.MinDimensions.Width = CLAY__MIN(CLAY__MAX(openLayoutElement.MinDimensions.Width, layoutConfig.Sizing.Width.Size.MinMax.Min), layoutConfig.Sizing.Width.Size.MinMax.Max)
	} else {
		openLayoutElement.Dimensions.Width = 0
	}

	// Clamp element min and max height to the values configured in the layout
	if layoutConfig.Sizing.Height.Type != CLAY__SIZING_TYPE_PERCENT {
		if layoutConfig.Sizing.Height.Size.MinMax.Max <= 0 { // Set the max size if the user didn't specify, makes calculations easier
			layoutConfig.Sizing.Height.Size.MinMax.Max = CLAY__MAXFLOAT
		}

		openLayoutElement.Dimensions.Height = CLAY__MIN(CLAY__MAX(openLayoutElement.Dimensions.Height, layoutConfig.Sizing.Height.Size.MinMax.Min), layoutConfig.Sizing.Height.Size.MinMax.Max)
		openLayoutElement.MinDimensions.Height = CLAY__MIN(CLAY__MAX(openLayoutElement.MinDimensions.Height, layoutConfig.Sizing.Height.Size.MinMax.Min), layoutConfig.Sizing.Height.Size.MinMax.Max)
	} else {
		openLayoutElement.Dimensions.Height = 0
	}

	Clay__UpdateAspectRatioBox(openLayoutElement)

	elementIsFloating := Clay__ElementHasConfig(openLayoutElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)

	// Close the currently open element
	closingElementIndex := Clay__Array_RemoveSwapback(&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-1)

	// Get the currently open parent
	openLayoutElement = Clay__GetOpenLayoutElement()

	if currentContext.OpenLayoutElementStack.Length > 1 {
		if elementIsFloating {
			openLayoutElement.FloatingChildrenCount++
			return
		}
		openLayoutElement.ChildrenOrTextContent.Children.Length++
		Clay__Array_Add(&currentContext.LayoutElementChildrenBuffer, closingElementIndex)
	}

}

func Clay__ElementHasConfig(layoutElement *Clay_LayoutElement, configType Clay__ElementConfigType) bool {
	for i := int32(0); i < layoutElement.ElementConfigs.Length; i++ {
		if Clay__Slice_Get(&layoutElement.ElementConfigs, i).Type == configType {
			return true
		}
	}
	return false
}

func Clay__UpdateAspectRatioBox(layoutElement *Clay_LayoutElement) {
	for j := int32(0); j < layoutElement.ElementConfigs.Length; j++ {
		config := Clay__Slice_Get(&layoutElement.ElementConfigs, j)
		if config.Type == CLAY__ELEMENT_CONFIG_TYPE_ASPECT {
			aspectConfig := config.Config.AspectRatioElementConfig
			if aspectConfig.AspectRatio == 0 {
				break
			}

			if layoutElement.Dimensions.Width == 0 && layoutElement.Dimensions.Height != 0 {
				layoutElement.Dimensions.Width = layoutElement.Dimensions.Height * aspectConfig.AspectRatio
			} else if layoutElement.Dimensions.Width != 0 && layoutElement.Dimensions.Height == 0 {
				layoutElement.Dimensions.Height = layoutElement.Dimensions.Width * (1 / aspectConfig.AspectRatio)
			}
			break
		}
	}
}

func Clay__RenderDebugView() {
	panic("Clay__RenderDebugViewElementConfigHeader not implemented")
}
func Clay__FindElementConfigWithType(element *Clay_LayoutElement, configType Clay__ElementConfigType) Clay_ElementConfigUnion {
	for i := int32(0); i < element.ElementConfigs.Length; i++ {
		config := Clay__Slice_Get(&element.ElementConfigs, i)
		if config.Type == configType {
			return config.Config
		}
	}
	return Clay_ElementConfigUnion{}
}

func Clay__SizeContainersAlongAxis(xAxis bool) {
	currentContext := Clay_GetCurrentContext()
	bfsBuffer := currentContext.LayoutElementChildrenBuffer
	resizableContainerBuffer := currentContext.OpenLayoutElementStack

	for rootIndex := int32(0); rootIndex < currentContext.LayoutElementTreeRoots.Length; rootIndex++ {
		bfsBuffer.Length = 0

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

		for i := int32(0); i < bfsBuffer.Length; i++ {
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
			resizableContainerBuffer.Length = 0
			parentChildGap := parentStyleConfig.ChildGap

			for childOffset := int32(0); childOffset < int32(parent.ChildrenOrTextContent.Children.Length); childOffset++ {
				fmt.Printf("Clay__SizeContainersAlongAxis childOffset: %d\n", childOffset)
				fmt.Printf("Clay__SizeContainersAlongAxis parent.ChildrenOrTextContent.Children.Length: %d\n", parent.ChildrenOrTextContent.Children.Length)
				childElementIndex := mem.NewMemSliceWithData(parent.ChildrenOrTextContent.Children.Elements).Get(childOffset)
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
				childElementIndex := mem.NewMemSliceWithData(parent.ChildrenOrTextContent.Children.Elements).Get(childOffset)
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
					for sizeToDistribute < -CLAY__EPSILON && resizableContainerBuffer.Length > 0 {
						var largest float32 = 0
						var secondLargest float32 = 0
						var widthToAdd float32 = sizeToDistribute
						for childIndex := int32(0); childIndex < resizableContainerBuffer.Length; childIndex++ {
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

						widthToAdd = CLAY__MAX(widthToAdd, sizeToDistribute/float32(resizableContainerBuffer.Length))

						for childIndex := int32(0); childIndex < resizableContainerBuffer.Length; childIndex++ {
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
								sizeToDistribute -= (childSize - previousWidth)
							}
						}
					}
					// The content is too small, allow SIZING_GROW containers to expand
				} else if sizeToDistribute > 0 && growContainerCount > 0 {
					for childIndex := int32(0); childIndex < resizableContainerBuffer.Length; childIndex++ {
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
					for sizeToDistribute > CLAY__EPSILON && resizableContainerBuffer.Length > 0 {
						var smallest float32 = CLAY__MAXFLOAT
						var secondSmallest float32 = CLAY__MAXFLOAT
						widthToAdd := sizeToDistribute
						for childIndex := int32(0); childIndex < resizableContainerBuffer.Length; childIndex++ {
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

						widthToAdd = CLAY__MIN(widthToAdd, sizeToDistribute/float32(resizableContainerBuffer.Length))

						for childIndex := int32(0); childIndex < resizableContainerBuffer.Length; childIndex++ {
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
								sizeToDistribute -= (childSize - previousWidth)
							}
						}
					}
				}
				// Sizing along the non layout axis ("off axis")
			} else {
				for childOffset := int32(0); childOffset < resizableContainerBuffer.Length; childOffset++ {

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
				}
			}

		}

	}
}

func Clay__CalculateFinalLayout() {
	currentContext := Clay_GetCurrentContext()
	// Calculate sizing along the X axis
	Clay__SizeContainersAlongAxis(true)

	// Wrap text
	for textElementIndex := int32(0); textElementIndex < currentContext.TextElementData.Length; textElementIndex++ {
		textElementData := Clay__Array_Get(&currentContext.TextElementData, textElementIndex)
		wrappedLinesData := mem.MArray_GetSlice(&currentContext.WrappedTextLines, 0, currentContext.WrappedTextLines.Length)
		wrappedLines := Clay__Slice[Clay__WrappedTextLine]{
			Length:        currentContext.WrappedTextLines.Length,
			InternalArray: wrappedLinesData,
		}
		textElementData.WrappedLines = wrappedLines
		containerElement := Clay__Array_Get(&currentContext.LayoutElements, textElementData.ElementIndex)
		textConfig := Clay__FindElementConfigWithType(containerElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT).TextElementConfig
		measureTextCacheItem := Clay__MeasureTextCached(&textElementData.Text, textConfig)
		var lineWidth float32 = 0
		var lineHeight float32 = 0
		if textConfig.LineHeight > 0 {
			lineHeight = float32(textConfig.LineHeight)
		} else {
			lineHeight = textElementData.PreferredDimensions.Height
		}
		var lineLengthChars int32 = 0
		var lineStartOffset int32 = 0

		if !measureTextCacheItem.ContainsNewlines && textElementData.PreferredDimensions.Width <= containerElement.Dimensions.Width {
			Clay__Array_Add(&currentContext.WrappedTextLines, Clay__WrappedTextLine{Dimensions: containerElement.Dimensions, Line: textElementData.Text})
			textElementData.WrappedLines.Length++
			continue
		}
		spaceWidth := Clay__MeasureText(
			Clay_StringSlice{
				Length:    1,
				Chars:     CLAY__SPACECHAR.Chars,
				BaseChars: CLAY__SPACECHAR.Chars,
			},
			textConfig,
			currentContext.MeasureTextUserData).Width
		wordIndex := measureTextCacheItem.MeasuredWordsStartIndex
		for wordIndex != -1 {
			if currentContext.WrappedTextLines.Length > currentContext.WrappedTextLines.Capacity-1 {
				break
			}
			measuredWord := Clay__Array_Get(&currentContext.MeasuredWords, wordIndex)
			// Only word on the line is too large, just render it anyway
			if lineLengthChars == 0 && lineWidth+measuredWord.Width > containerElement.Dimensions.Width {
				Clay__Array_Add(&currentContext.WrappedTextLines, Clay__WrappedTextLine{
					Dimensions: Clay_Dimensions{
						Width:  measuredWord.Width,
						Height: lineHeight,
					},
					Line: Clay_String{
						Length: measuredWord.Length,
						Chars:  textElementData.Text.Chars[measuredWord.StartOffset : measuredWord.StartOffset+measuredWord.Length],
					},
				},
				)
				textElementData.WrappedLines.Length++
				wordIndex = measuredWord.Next
				lineStartOffset = measuredWord.StartOffset + measuredWord.Length
			} else if measuredWord.Length == 0 || lineWidth+measuredWord.Width > containerElement.Dimensions.Width {
				// measuredWord->length == 0 means a newline character
				// Wrapped text lines list has overflowed, just render out the line
				maxIndex := CLAY__MAX(lineStartOffset+lineLengthChars-1, 0)
				finalCharIsSpace := textElementData.Text.Chars[maxIndex] == ' '

				var Dimensions Clay_Dimensions
				var Line Clay_String
				if finalCharIsSpace {
					Dimensions = Clay_Dimensions{
						Width:  lineWidth - spaceWidth,
						Height: lineHeight,
					}
					Line = Clay_String{
						Length: lineLengthChars - 1,
						Chars:  textElementData.Text.Chars[lineStartOffset : lineStartOffset+lineLengthChars-1],
					}
				} else {
					Dimensions = Clay_Dimensions{
						Width:  lineWidth,
						Height: lineHeight,
					}
					Line = Clay_String{
						Length: lineLengthChars,
						Chars:  textElementData.Text.Chars[lineStartOffset : lineStartOffset+lineLengthChars],
					}
				}
				Clay__Array_Add(&currentContext.WrappedTextLines,
					Clay__WrappedTextLine{
						Dimensions: Dimensions,
						Line:       Line,
					},
				)

				textElementData.WrappedLines.Length++
				if lineLengthChars == 0 || measuredWord.Length == 0 {
					wordIndex = measuredWord.Next
				}
				lineWidth = 0
				lineLengthChars = 0
				lineStartOffset = measuredWord.StartOffset
			} else {
				lineWidth += measuredWord.Width + float32(textConfig.LetterSpacing)
				lineLengthChars += measuredWord.Length
				wordIndex = measuredWord.Next
			}
		}
		if lineLengthChars > 0 {
			Clay__Array_Add(&currentContext.WrappedTextLines, Clay__WrappedTextLine{
				Dimensions: Clay_Dimensions{
					Width:  lineWidth - float32(textConfig.LetterSpacing),
					Height: lineHeight,
				},
				Line: Clay_String{
					Length: lineLengthChars,
					Chars:  textElementData.Text.Chars[lineStartOffset : lineStartOffset+lineLengthChars],
				},
			})
			textElementData.WrappedLines.Length++
		}
		containerElement.Dimensions.Height = lineHeight * float32(textElementData.WrappedLines.Length)
	}

	// Scale vertical heights according to aspect ratio
	for aspectRatioElementIndex := int32(0); aspectRatioElementIndex < currentContext.AspectRatioElementIndexes.Length; aspectRatioElementIndex++ {
		aspectElement := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&currentContext.AspectRatioElementIndexes, aspectRatioElementIndex))
		aspectRatioElementConfig := Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).AspectRatioElementConfig
		aspectElement.Dimensions.Height = (1 / aspectRatioElementConfig.AspectRatio) * aspectElement.Dimensions.Width
		aspectElement.LayoutConfig.Sizing.Height.Size.MinMax.Max = aspectElement.Dimensions.Height
	}

	// Propagate effect of text wrapping, aspect scaling etc. on height of parents
	dfsBuffer := currentContext.LayoutElementTreeNodeArray1
	dfsBuffer.Length = 0
	for layoutElementTreeRootIndex := int32(0); layoutElementTreeRootIndex < currentContext.LayoutElementTreeRoots.Length; layoutElementTreeRootIndex++ {
		layoutElementTreeRoot := Clay__Array_Get(&currentContext.LayoutElementTreeRoots, layoutElementTreeRootIndex)
		Clay__Array_Set(&currentContext.TreeNodeVisited, dfsBuffer.Length, false)
		Clay__Array_Add(&dfsBuffer, Clay__LayoutElementTreeNode{LayoutElement: Clay__Array_Get(&currentContext.LayoutElements, layoutElementTreeRoot.LayoutElementIndex)})
	}

	// while (dfsBuffer.length > 0) {
	// 	Clay__LayoutElementTreeNode *currentElementTreeNode = Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, (int)dfsBuffer.length - 1);
	// 	Clay_LayoutElement *currentElement = currentElementTreeNode->layoutElement;
	// 	if (!context->treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
	// 		context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;
	// 		// If the element has no children or is the container for a text element, don't bother inspecting it
	// 		if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement->childrenOrTextContent.children.length == 0) {
	// 			dfsBuffer.length--;
	// 			continue;
	// 		}
	// 		// Add the children to the DFS buffer (needs to be pushed in reverse so that stack traversal is in correct layout order)
	// 		for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; i++) {
	// 			context->treeNodeVisited.internalArray[dfsBuffer.length] = false;
	// 			Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]) });
	// 		}
	// 		continue;
	// 	}
	// 	dfsBuffer.length--;

	// 	// DFS node has been visited, this is on the way back up to the root
	// 	Clay_LayoutConfig *layoutConfig = currentElement->layoutConfig;
	// 	if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
	// 		// Resize any parent containers that have grown in height along their non layout axis
	// 		for (int32_t j = 0; j < currentElement->childrenOrTextContent.children.length; ++j) {
	// 			Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[j]);
	// 			float childHeightWithPadding = CLAY__MAX(childElement->dimensions.height + layoutConfig->padding.top + layoutConfig->padding.bottom, currentElement->dimensions.height);
	// 			currentElement->dimensions.height = CLAY__MIN(CLAY__MAX(childHeightWithPadding, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
	// 		}
	// 	} else if (layoutConfig->layoutDirection == CLAY_TOP_TO_BOTTOM) {
	// 		// Resizing along the layout axis
	// 		float contentHeight = (float)(layoutConfig->padding.top + layoutConfig->padding.bottom);
	// 		for (int32_t j = 0; j < currentElement->childrenOrTextContent.children.length; ++j) {
	// 			Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[j]);
	// 			contentHeight += childElement->dimensions.height;
	// 		}
	// 		contentHeight += (float)(CLAY__MAX(currentElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
	// 		currentElement->dimensions.height = CLAY__MIN(CLAY__MAX(contentHeight, layoutConfig->sizing.height.size.minMax.min), layoutConfig->sizing.height.size.minMax.max);
	// 	}
	// }

	// // Calculate sizing along the Y axis
	// Clay__SizeContainersAlongAxis(false);

	// // Scale horizontal widths according to aspect ratio
	// for (int32_t i = 0; i < context->aspectRatioElementIndexes.length; ++i) {
	// 	Clay_LayoutElement* aspectElement = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->aspectRatioElementIndexes, i));
	// 	Clay_AspectRatioElementConfig *config = Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig;
	// 	aspectElement->dimensions.width = config->aspectRatio * aspectElement->dimensions.height;
	// }

	// // Sort tree roots by z-index
	// int32_t sortMax = context->layoutElementTreeRoots.length - 1;
	// while (sortMax > 0) { // todo dumb bubble sort
	// 	for (int32_t i = 0; i < sortMax; ++i) {
	// 		Clay__LayoutElementTreeRoot current = *Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, i);
	// 		Clay__LayoutElementTreeRoot next = *Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, i + 1);
	// 		if (next.zIndex < current.zIndex) {
	// 			Clay__LayoutElementTreeRootArray_Set(&context->layoutElementTreeRoots, i, next);
	// 			Clay__LayoutElementTreeRootArray_Set(&context->layoutElementTreeRoots, i + 1, current);
	// 		}
	// 	}
	// 	sortMax--;
	// }

	// // Calculate final positions and generate render commands
	// context->renderCommands.length = 0;
	// dfsBuffer.length = 0;
	// for (int32_t rootIndex = 0; rootIndex < context->layoutElementTreeRoots.length; ++rootIndex) {
	// 	dfsBuffer.length = 0;
	// 	Clay__LayoutElementTreeRoot *root = Clay__LayoutElementTreeRootArray_Get(&context->layoutElementTreeRoots, rootIndex);
	// 	Clay_LayoutElement *rootElement = Clay_LayoutElementArray_Get(&context->layoutElements, (int)root->layoutElementIndex);
	// 	Clay_Vector2 rootPosition = CLAY__DEFAULT_STRUCT;
	// 	Clay_LayoutElementHashMapItem *parentHashMapItem = Clay__GetHashMapItem(root->parentId);
	// 	// Position root floating containers
	// 	if (Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) && parentHashMapItem) {
	// 		Clay_FloatingElementConfig *config = Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
	// 		Clay_Dimensions rootDimensions = rootElement->dimensions;
	// 		Clay_BoundingBox parentBoundingBox = parentHashMapItem->boundingBox;
	// 		// Set X position
	// 		Clay_Vector2 targetAttachPosition = CLAY__DEFAULT_STRUCT;
	// 		switch (config->attachPoints.parent) {
	// 			case CLAY_ATTACH_POINT_LEFT_TOP:
	// 			case CLAY_ATTACH_POINT_LEFT_CENTER:
	// 			case CLAY_ATTACH_POINT_LEFT_BOTTOM: targetAttachPosition.x = parentBoundingBox.x; break;
	// 			case CLAY_ATTACH_POINT_CENTER_TOP:
	// 			case CLAY_ATTACH_POINT_CENTER_CENTER:
	// 			case CLAY_ATTACH_POINT_CENTER_BOTTOM: targetAttachPosition.x = parentBoundingBox.x + (parentBoundingBox.width / 2); break;
	// 			case CLAY_ATTACH_POINT_RIGHT_TOP:
	// 			case CLAY_ATTACH_POINT_RIGHT_CENTER:
	// 			case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.x = parentBoundingBox.x + parentBoundingBox.width; break;
	// 		}
	// 		switch (config->attachPoints.element) {
	// 			case CLAY_ATTACH_POINT_LEFT_TOP:
	// 			case CLAY_ATTACH_POINT_LEFT_CENTER:
	// 			case CLAY_ATTACH_POINT_LEFT_BOTTOM: break;
	// 			case CLAY_ATTACH_POINT_CENTER_TOP:
	// 			case CLAY_ATTACH_POINT_CENTER_CENTER:
	// 			case CLAY_ATTACH_POINT_CENTER_BOTTOM: targetAttachPosition.x -= (rootDimensions.width / 2); break;
	// 			case CLAY_ATTACH_POINT_RIGHT_TOP:
	// 			case CLAY_ATTACH_POINT_RIGHT_CENTER:
	// 			case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.x -= rootDimensions.width; break;
	// 		}
	// 		switch (config->attachPoints.parent) { // I know I could merge the x and y switch statements, but this is easier to read
	// 			case CLAY_ATTACH_POINT_LEFT_TOP:
	// 			case CLAY_ATTACH_POINT_RIGHT_TOP:
	// 			case CLAY_ATTACH_POINT_CENTER_TOP: targetAttachPosition.y = parentBoundingBox.y; break;
	// 			case CLAY_ATTACH_POINT_LEFT_CENTER:
	// 			case CLAY_ATTACH_POINT_CENTER_CENTER:
	// 			case CLAY_ATTACH_POINT_RIGHT_CENTER: targetAttachPosition.y = parentBoundingBox.y + (parentBoundingBox.height / 2); break;
	// 			case CLAY_ATTACH_POINT_LEFT_BOTTOM:
	// 			case CLAY_ATTACH_POINT_CENTER_BOTTOM:
	// 			case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.y = parentBoundingBox.y + parentBoundingBox.height; break;
	// 		}
	// 		switch (config->attachPoints.element) {
	// 			case CLAY_ATTACH_POINT_LEFT_TOP:
	// 			case CLAY_ATTACH_POINT_RIGHT_TOP:
	// 			case CLAY_ATTACH_POINT_CENTER_TOP: break;
	// 			case CLAY_ATTACH_POINT_LEFT_CENTER:
	// 			case CLAY_ATTACH_POINT_CENTER_CENTER:
	// 			case CLAY_ATTACH_POINT_RIGHT_CENTER: targetAttachPosition.y -= (rootDimensions.height / 2); break;
	// 			case CLAY_ATTACH_POINT_LEFT_BOTTOM:
	// 			case CLAY_ATTACH_POINT_CENTER_BOTTOM:
	// 			case CLAY_ATTACH_POINT_RIGHT_BOTTOM: targetAttachPosition.y -= rootDimensions.height; break;
	// 		}
	// 		targetAttachPosition.x += config->offset.x;
	// 		targetAttachPosition.y += config->offset.y;
	// 		rootPosition = targetAttachPosition;
	// 	}
	// 	if (root->clipElementId) {
	// 		Clay_LayoutElementHashMapItem *clipHashMapItem = Clay__GetHashMapItem(root->clipElementId);
	// 		if (clipHashMapItem) {
	// 			// Floating elements that are attached to scrolling contents won't be correctly positioned if external scroll handling is enabled, fix here
	// 			if (context->externalScrollHandlingEnabled) {
	// 				Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(clipHashMapItem->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
	// 				if (clipConfig->horizontal) {
	// 					rootPosition.x += clipConfig->childOffset.x;
	// 				}
	// 				if (clipConfig->vertical) {
	// 					rootPosition.y += clipConfig->childOffset.y;
	// 				}
	// 			}
	// 			Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	// 				.boundingBox = clipHashMapItem->boundingBox,
	// 				.userData = 0,
	// 				.id = Clay__HashNumber(rootElement->id, rootElement->childrenOrTextContent.children.length + 10).id, // TODO need a better strategy for managing derived ids
	// 				.zIndex = root->zIndex,
	// 				.commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START,
	// 			});
	// 		}
	// 	}
	// 	Clay__LayoutElementTreeNodeArray_Add(&dfsBuffer, CLAY__INIT(Clay__LayoutElementTreeNode) { .layoutElement = rootElement, .position = rootPosition, .nextChildOffset = { .x = (float)rootElement->layoutConfig->padding.left, .y = (float)rootElement->layoutConfig->padding.top } });

	// 	context->treeNodeVisited.internalArray[0] = false;
	// 	while (dfsBuffer.length > 0) {
	// 		Clay__LayoutElementTreeNode *currentElementTreeNode = Clay__LayoutElementTreeNodeArray_Get(&dfsBuffer, (int)dfsBuffer.length - 1);
	// 		Clay_LayoutElement *currentElement = currentElementTreeNode->layoutElement;
	// 		Clay_LayoutConfig *layoutConfig = currentElement->layoutConfig;
	// 		Clay_Vector2 scrollOffset = CLAY__DEFAULT_STRUCT;

	// 		// This will only be run a single time for each element in downwards DFS order
	// 		if (!context->treeNodeVisited.internalArray[dfsBuffer.length - 1]) {
	// 			context->treeNodeVisited.internalArray[dfsBuffer.length - 1] = true;

	// 			Clay_BoundingBox currentElementBoundingBox = { currentElementTreeNode->position.x, currentElementTreeNode->position.y, currentElement->dimensions.width, currentElement->dimensions.height };
	// 			if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)) {
	// 				Clay_FloatingElementConfig *floatingElementConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).floatingElementConfig;
	// 				Clay_Dimensions expand = floatingElementConfig->expand;
	// 				currentElementBoundingBox.x -= expand.width;
	// 				currentElementBoundingBox.width += expand.width * 2;
	// 				currentElementBoundingBox.y -= expand.height;
	// 				currentElementBoundingBox.height += expand.height * 2;
	// 			}

	// 			Clay__ScrollContainerDataInternal *scrollContainerData = CLAY__NULL;
	// 			// Apply scroll offsets to container
	// 			if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP)) {
	// 				Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;

	// 				// This linear scan could theoretically be slow under very strange conditions, but I can't imagine a real UI with more than a few 10's of scroll containers
	// 				for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
	// 					Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
	// 					if (mapping->layoutElement == currentElement) {
	// 						scrollContainerData = mapping;
	// 						mapping->boundingBox = currentElementBoundingBox;
	// 						scrollOffset = clipConfig->childOffset;
	// 						if (context->externalScrollHandlingEnabled) {
	// 							scrollOffset = CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
	// 						}
	// 						break;
	// 					}
	// 				}
	// 			}

	// 			Clay_LayoutElementHashMapItem *hashMapItem = Clay__GetHashMapItem(currentElement->id);
	// 			if (hashMapItem) {
	// 				hashMapItem->boundingBox = currentElementBoundingBox;
	// 			}

	// 			int32_t sortedConfigIndexes[20];
	// 			for (int32_t elementConfigIndex = 0; elementConfigIndex < currentElement->elementConfigs.length; ++elementConfigIndex) {
	// 				sortedConfigIndexes[elementConfigIndex] = elementConfigIndex;
	// 			}
	// 			sortMax = currentElement->elementConfigs.length - 1;
	// 			while (sortMax > 0) { // todo dumb bubble sort
	// 				for (int32_t i = 0; i < sortMax; ++i) {
	// 					int32_t current = sortedConfigIndexes[i];
	// 					int32_t next = sortedConfigIndexes[i + 1];
	// 					Clay__ElementConfigType currentType = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, current)->type;
	// 					Clay__ElementConfigType nextType = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, next)->type;
	// 					if (nextType == CLAY__ELEMENT_CONFIG_TYPE_CLIP || currentType == CLAY__ELEMENT_CONFIG_TYPE_BORDER) {
	// 						sortedConfigIndexes[i] = next;
	// 						sortedConfigIndexes[i + 1] = current;
	// 					}
	// 				}
	// 				sortMax--;
	// 			}

	// 			bool emitRectangle = false;
	// 			// Create the render commands for this element
	// 			Clay_SharedElementConfig *sharedConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig;
	// 			if (sharedConfig && sharedConfig->backgroundColor.a > 0) {
	// 			   emitRectangle = true;
	// 			}
	// 			else if (!sharedConfig) {
	// 				emitRectangle = false;
	// 				sharedConfig = &Clay_SharedElementConfig_DEFAULT;
	// 			}
	// 			for (int32_t elementConfigIndex = 0; elementConfigIndex < currentElement->elementConfigs.length; ++elementConfigIndex) {
	// 				Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&currentElement->elementConfigs, sortedConfigIndexes[elementConfigIndex]);
	// 				Clay_RenderCommand renderCommand = {
	// 					.boundingBox = currentElementBoundingBox,
	// 					.userData = sharedConfig->userData,
	// 					.id = currentElement->id,
	// 				};

	// 				bool offscreen = Clay__ElementIsOffscreen(&currentElementBoundingBox);
	// 				// Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
	// 				bool shouldRender = !offscreen;
	// 				switch (elementConfig->type) {
	// 					case CLAY__ELEMENT_CONFIG_TYPE_ASPECT:
	// 					case CLAY__ELEMENT_CONFIG_TYPE_FLOATING:
	// 					case CLAY__ELEMENT_CONFIG_TYPE_SHARED:
	// 					case CLAY__ELEMENT_CONFIG_TYPE_BORDER: {
	// 						shouldRender = false;
	// 						break;
	// 					}
	// 					case CLAY__ELEMENT_CONFIG_TYPE_CLIP: {
	// 						renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START;
	// 						renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
	// 							.clip = {
	// 								.horizontal = elementConfig->config.clipElementConfig->horizontal,
	// 								.vertical = elementConfig->config.clipElementConfig->vertical,
	// 							}
	// 						};
	// 						break;
	// 					}
	// 					case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: {
	// 						renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_IMAGE;
	// 						renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
	// 							.image = {
	// 								.backgroundColor = sharedConfig->backgroundColor,
	// 								.cornerRadius = sharedConfig->cornerRadius,
	// 								.imageData = elementConfig->config.imageElementConfig->imageData,
	// 						   }
	// 						};
	// 						emitRectangle = false;
	// 						break;
	// 					}
	// 					case CLAY__ELEMENT_CONFIG_TYPE_TEXT: {
	// 						if (!shouldRender) {
	// 							break;
	// 						}
	// 						shouldRender = false;
	// 						Clay_ElementConfigUnion configUnion = elementConfig->config;
	// 						Clay_TextElementConfig *textElementConfig = configUnion.textElementConfig;
	// 						float naturalLineHeight = currentElement->childrenOrTextContent.textElementData->preferredDimensions.height;
	// 						float finalLineHeight = textElementConfig->lineHeight > 0 ? (float)textElementConfig->lineHeight : naturalLineHeight;
	// 						float lineHeightOffset = (finalLineHeight - naturalLineHeight) / 2;
	// 						float yPosition = lineHeightOffset;
	// 						for (int32_t lineIndex = 0; lineIndex < currentElement->childrenOrTextContent.textElementData->wrappedLines.length; ++lineIndex) {
	// 							Clay__WrappedTextLine *wrappedLine = Clay__WrappedTextLineArraySlice_Get(&currentElement->childrenOrTextContent.textElementData->wrappedLines, lineIndex);
	// 							if (wrappedLine->line.length == 0) {
	// 								yPosition += finalLineHeight;
	// 								continue;
	// 							}
	// 							float offset = (currentElementBoundingBox.width - wrappedLine->dimensions.width);
	// 							if (textElementConfig->textAlignment == CLAY_TEXT_ALIGN_LEFT) {
	// 								offset = 0;
	// 							}
	// 							if (textElementConfig->textAlignment == CLAY_TEXT_ALIGN_CENTER) {
	// 								offset /= 2;
	// 							}
	// 							Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	// 								.boundingBox = { currentElementBoundingBox.x + offset, currentElementBoundingBox.y + yPosition, wrappedLine->dimensions.width, wrappedLine->dimensions.height },
	// 								.renderData = { .text = {
	// 									.stringContents = CLAY__INIT(Clay_StringSlice) { .length = wrappedLine->line.length, .chars = wrappedLine->line.chars, .baseChars = currentElement->childrenOrTextContent.textElementData->text.chars },
	// 									.textColor = textElementConfig->textColor,
	// 									.fontId = textElementConfig->fontId,
	// 									.fontSize = textElementConfig->fontSize,
	// 									.letterSpacing = textElementConfig->letterSpacing,
	// 									.lineHeight = textElementConfig->lineHeight,
	// 								}},
	// 								.userData = textElementConfig->userData,
	// 								.id = Clay__HashNumber(lineIndex, currentElement->id).id,
	// 								.zIndex = root->zIndex,
	// 								.commandType = CLAY_RENDER_COMMAND_TYPE_TEXT,
	// 							});
	// 							yPosition += finalLineHeight;

	// 							if (!context->disableCulling && (currentElementBoundingBox.y + yPosition > context->layoutDimensions.height)) {
	// 								break;
	// 							}
	// 						}
	// 						break;
	// 					}
	// 					case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM: {
	// 						renderCommand.commandType = CLAY_RENDER_COMMAND_TYPE_CUSTOM;
	// 						renderCommand.renderData = CLAY__INIT(Clay_RenderData) {
	// 							.custom = {
	// 								.backgroundColor = sharedConfig->backgroundColor,
	// 								.cornerRadius = sharedConfig->cornerRadius,
	// 								.customData = elementConfig->config.customElementConfig->customData,
	// 							}
	// 						};
	// 						emitRectangle = false;
	// 						break;
	// 					}
	// 					default: break;
	// 				}
	// 				if (shouldRender) {
	// 					Clay__AddRenderCommand(renderCommand);
	// 				}
	// 				if (offscreen) {
	// 					// NOTE: You may be tempted to try an early return / continue if an element is off screen. Why bother calculating layout for its children, right?
	// 					// Unfortunately, a FLOATING_CONTAINER may be defined that attaches to a child or grandchild of this element, which is large enough to still
	// 					// be on screen, even if this element isn't. That depends on this element and it's children being laid out correctly (even if they are entirely off screen)
	// 				}
	// 			}

	// 			if (emitRectangle) {
	// 				Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	// 					.boundingBox = currentElementBoundingBox,
	// 					.renderData = { .rectangle = {
	// 							.backgroundColor = sharedConfig->backgroundColor,
	// 							.cornerRadius = sharedConfig->cornerRadius,
	// 					}},
	// 					.userData = sharedConfig->userData,
	// 					.id = currentElement->id,
	// 					.zIndex = root->zIndex,
	// 					.commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
	// 				});
	// 			}

	// 			// Setup initial on-axis alignment
	// 			if (!Clay__ElementHasConfig(currentElementTreeNode->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
	// 				Clay_Dimensions contentSize = {0,0};
	// 				if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
	// 					for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
	// 						Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
	// 						contentSize.width += childElement->dimensions.width;
	// 						contentSize.height = CLAY__MAX(contentSize.height, childElement->dimensions.height);
	// 					}
	// 					contentSize.width += (float)(CLAY__MAX(currentElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
	// 					float extraSpace = currentElement->dimensions.width - (float)(layoutConfig->padding.left + layoutConfig->padding.right) - contentSize.width;
	// 					switch (layoutConfig->childAlignment.x) {
	// 						case CLAY_ALIGN_X_LEFT: extraSpace = 0; break;
	// 						case CLAY_ALIGN_X_CENTER: extraSpace /= 2; break;
	// 						default: break;
	// 					}
	// 					currentElementTreeNode->nextChildOffset.x += extraSpace;
	// 					extraSpace = CLAY__MAX(0, extraSpace);
	// 				} else {
	// 					for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
	// 						Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
	// 						contentSize.width = CLAY__MAX(contentSize.width, childElement->dimensions.width);
	// 						contentSize.height += childElement->dimensions.height;
	// 					}
	// 					contentSize.height += (float)(CLAY__MAX(currentElement->childrenOrTextContent.children.length - 1, 0) * layoutConfig->childGap);
	// 					float extraSpace = currentElement->dimensions.height - (float)(layoutConfig->padding.top + layoutConfig->padding.bottom) - contentSize.height;
	// 					switch (layoutConfig->childAlignment.y) {
	// 						case CLAY_ALIGN_Y_TOP: extraSpace = 0; break;
	// 						case CLAY_ALIGN_Y_CENTER: extraSpace /= 2; break;
	// 						default: break;
	// 					}
	// 					extraSpace = CLAY__MAX(0, extraSpace);
	// 					currentElementTreeNode->nextChildOffset.y += extraSpace;
	// 				}

	// 				if (scrollContainerData) {
	// 					scrollContainerData->contentSize = CLAY__INIT(Clay_Dimensions) { contentSize.width + (float)(layoutConfig->padding.left + layoutConfig->padding.right), contentSize.height + (float)(layoutConfig->padding.top + layoutConfig->padding.bottom) };
	// 				}
	// 			}
	// 		}
	// 		else {
	// 			// DFS is returning upwards backwards
	// 			bool closeClipElement = false;
	// 			Clay_ClipElementConfig *clipConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).clipElementConfig;
	// 			if (clipConfig) {
	// 				closeClipElement = true;
	// 				for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
	// 					Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
	// 					if (mapping->layoutElement == currentElement) {
	// 						scrollOffset = clipConfig->childOffset;
	// 						if (context->externalScrollHandlingEnabled) {
	// 							scrollOffset = CLAY__INIT(Clay_Vector2) CLAY__DEFAULT_STRUCT;
	// 						}
	// 						break;
	// 					}
	// 				}
	// 			}

	// 			if (Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER)) {
	// 				Clay_LayoutElementHashMapItem *currentElementData = Clay__GetHashMapItem(currentElement->id);
	// 				Clay_BoundingBox currentElementBoundingBox = currentElementData->boundingBox;

	// 				// Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
	// 				if (!Clay__ElementIsOffscreen(&currentElementBoundingBox)) {
	// 					Clay_SharedElementConfig *sharedConfig = Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED) ? Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).sharedElementConfig : &Clay_SharedElementConfig_DEFAULT;
	// 					Clay_BorderElementConfig *borderConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER).borderElementConfig;
	// 					Clay_RenderCommand renderCommand = {
	// 							.boundingBox = currentElementBoundingBox,
	// 							.renderData = { .border = {
	// 								.color = borderConfig->color,
	// 								.cornerRadius = sharedConfig->cornerRadius,
	// 								.width = borderConfig->width
	// 							}},
	// 							.userData = sharedConfig->userData,
	// 							.id = Clay__HashNumber(currentElement->id, currentElement->childrenOrTextContent.children.length).id,
	// 							.commandType = CLAY_RENDER_COMMAND_TYPE_BORDER,
	// 					};
	// 					Clay__AddRenderCommand(renderCommand);
	// 					if (borderConfig->width.betweenChildren > 0 && borderConfig->color.a > 0) {
	// 						float halfGap = layoutConfig->childGap / 2;
	// 						Clay_Vector2 borderOffset = { (float)layoutConfig->padding.left - halfGap, (float)layoutConfig->padding.top - halfGap };
	// 						if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
	// 							for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
	// 								Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
	// 								if (i > 0) {
	// 									Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	// 										.boundingBox = { currentElementBoundingBox.x + borderOffset.x + scrollOffset.x, currentElementBoundingBox.y + scrollOffset.y, (float)borderConfig->width.betweenChildren, currentElement->dimensions.height },
	// 										.renderData = { .rectangle = {
	// 											.backgroundColor = borderConfig->color,
	// 										} },
	// 										.userData = sharedConfig->userData,
	// 										.id = Clay__HashNumber(currentElement->id, currentElement->childrenOrTextContent.children.length + 1 + i).id,
	// 										.commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
	// 									});
	// 								}
	// 								borderOffset.x += (childElement->dimensions.width + (float)layoutConfig->childGap);
	// 							}
	// 						} else {
	// 							for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
	// 								Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
	// 								if (i > 0) {
	// 									Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	// 										.boundingBox = { currentElementBoundingBox.x + scrollOffset.x, currentElementBoundingBox.y + borderOffset.y + scrollOffset.y, currentElement->dimensions.width, (float)borderConfig->width.betweenChildren },
	// 										.renderData = { .rectangle = {
	// 												.backgroundColor = borderConfig->color,
	// 										} },
	// 										.userData = sharedConfig->userData,
	// 										.id = Clay__HashNumber(currentElement->id, currentElement->childrenOrTextContent.children.length + 1 + i).id,
	// 										.commandType = CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
	// 									});
	// 								}
	// 								borderOffset.y += (childElement->dimensions.height + (float)layoutConfig->childGap);
	// 							}
	// 						}
	// 					}
	// 				}
	// 			}
	// 			// This exists because the scissor needs to end _after_ borders between elements
	// 			if (closeClipElement) {
	// 				Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) {
	// 					.id = Clay__HashNumber(currentElement->id, rootElement->childrenOrTextContent.children.length + 11).id,
	// 					.commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_END,
	// 				});
	// 			}

	// 			dfsBuffer.length--;
	// 			continue;
	// 		}

	// 		// Add children to the DFS buffer
	// 		if (!Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT)) {
	// 			dfsBuffer.length += currentElement->childrenOrTextContent.children.length;
	// 			for (int32_t i = 0; i < currentElement->childrenOrTextContent.children.length; ++i) {
	// 				Clay_LayoutElement *childElement = Clay_LayoutElementArray_Get(&context->layoutElements, currentElement->childrenOrTextContent.children.elements[i]);
	// 				// Alignment along non layout axis
	// 				if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
	// 					currentElementTreeNode->nextChildOffset.y = currentElement->layoutConfig->padding.top;
	// 					float whiteSpaceAroundChild = currentElement->dimensions.height - (float)(layoutConfig->padding.top + layoutConfig->padding.bottom) - childElement->dimensions.height;
	// 					switch (layoutConfig->childAlignment.y) {
	// 						case CLAY_ALIGN_Y_TOP: break;
	// 						case CLAY_ALIGN_Y_CENTER: currentElementTreeNode->nextChildOffset.y += whiteSpaceAroundChild / 2; break;
	// 						case CLAY_ALIGN_Y_BOTTOM: currentElementTreeNode->nextChildOffset.y += whiteSpaceAroundChild; break;
	// 					}
	// 				} else {
	// 					currentElementTreeNode->nextChildOffset.x = currentElement->layoutConfig->padding.left;
	// 					float whiteSpaceAroundChild = currentElement->dimensions.width - (float)(layoutConfig->padding.left + layoutConfig->padding.right) - childElement->dimensions.width;
	// 					switch (layoutConfig->childAlignment.x) {
	// 						case CLAY_ALIGN_X_LEFT: break;
	// 						case CLAY_ALIGN_X_CENTER: currentElementTreeNode->nextChildOffset.x += whiteSpaceAroundChild / 2; break;
	// 						case CLAY_ALIGN_X_RIGHT: currentElementTreeNode->nextChildOffset.x += whiteSpaceAroundChild; break;
	// 					}
	// 				}

	// 				Clay_Vector2 childPosition = {
	// 					currentElementTreeNode->position.x + currentElementTreeNode->nextChildOffset.x + scrollOffset.x,
	// 					currentElementTreeNode->position.y + currentElementTreeNode->nextChildOffset.y + scrollOffset.y,
	// 				};

	// 				// DFS buffer elements need to be added in reverse because stack traversal happens backwards
	// 				uint32_t newNodeIndex = dfsBuffer.length - 1 - i;
	// 				dfsBuffer.internalArray[newNodeIndex] = CLAY__INIT(Clay__LayoutElementTreeNode) {
	// 					.layoutElement = childElement,
	// 					.position = { childPosition.x, childPosition.y },
	// 					.nextChildOffset = { .x = (float)childElement->layoutConfig->padding.left, .y = (float)childElement->layoutConfig->padding.top },
	// 				};
	// 				context->treeNodeVisited.internalArray[newNodeIndex] = false;

	// 				// Update parent offsets
	// 				if (layoutConfig->layoutDirection == CLAY_LEFT_TO_RIGHT) {
	// 					currentElementTreeNode->nextChildOffset.x += childElement->dimensions.width + (float)layoutConfig->childGap;
	// 				} else {
	// 					currentElementTreeNode->nextChildOffset.y += childElement->dimensions.height + (float)layoutConfig->childGap;
	// 				}
	// 			}
	// 		}
	// 	}

	// 	if (root->clipElementId) {
	// 		Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand) { .id = Clay__HashNumber(rootElement->id, rootElement->childrenOrTextContent.children.length + 11).id, .commandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_END });
	// 	}
	// }

}

func Clay__AddRenderCommand(renderCommand Clay_RenderCommand) {
	currentContext := Clay_GetCurrentContext()
	if currentContext.RenderCommands.Length < currentContext.RenderCommands.Capacity-1 {
		Clay__Array_Add(&currentContext.RenderCommands, renderCommand)
	} else {
		if !currentContext.BooleanWarnings.MaxRenderCommandsExceeded {
			currentContext.BooleanWarnings.MaxRenderCommandsExceeded = true
			currentContext.ErrorHandler.ErrorHandlerFunction(Clay_ErrorData{
				ErrorType: CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
				ErrorText: CLAY_STRING("Clay ran out of capacity while attempting to create render commands. This is usually caused by a large amount of wrapping text elements while close to the max element capacity. Try using Clay_SetMaxElementCount() with a higher value."),
				UserData:  currentContext.ErrorHandler.UserData,
			})
		}
	}
}

func Clay__OpenElementWithId(elementId Clay_ElementId) {
	currentContext := Clay_GetCurrentContext()
	if currentContext.LayoutElements.Length == currentContext.LayoutElements.Capacity-1 || currentContext.BooleanWarnings.MaxElementsExceeded {
		currentContext.BooleanWarnings.MaxElementsExceeded = true
		return
	}
	layoutElement := Clay_LayoutElement{}
	layoutElement.Id = elementId.Id
	openLayoutElement := Clay__Array_Add(&currentContext.LayoutElements, layoutElement)
	Clay__Array_Add(&currentContext.OpenLayoutElementStack, currentContext.LayoutElements.Length-1) // add the index of the new layout element to the open layout element stack
	Clay__AddHashMapItem(elementId, openLayoutElement)
	Clay__Array_Add(&currentContext.LayoutElementIdStrings, elementId.StringId)
	if currentContext.OpenClipElementStack.Length > 0 {
		Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length-1, Clay__Array_GetValue(&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length-1))
	} else {
		Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length-1, 0)
	}
}

func Clay__StoreLayoutConfig(config Clay_LayoutConfig) *Clay_LayoutConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_LayoutConfig{}
	}
	return Clay__Array_Add(&currentContext.LayoutConfigs, config)

}

func Clay__AttachElementConfig(config Clay_ElementConfigUnion, configType Clay__ElementConfigType) Clay_ElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return Clay_ElementConfig{}
	}
	openLayoutElement := Clay__GetOpenLayoutElement()
	openLayoutElement.ElementConfigs.Length++
	return *Clay__Array_Add(&currentContext.ElementConfigs, Clay_ElementConfig{Type: configType, Config: config})
}

func Clay__StoreSharedElementConfig(config Clay_SharedElementConfig) *Clay_SharedElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_SharedElementConfig{}
	}
	return Clay__Array_Add(&currentContext.SharedElementConfigs, config)
}

func Clay__StoreImageElementConfig(config Clay_ImageElementConfig) *Clay_ImageElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_ImageElementConfig{}
	}
	return Clay__Array_Add(&currentContext.ImageElementConfigs, config)
}

func Clay__StoreAspectRatioElementConfig(config Clay_AspectRatioElementConfig) *Clay_AspectRatioElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_AspectRatioElementConfig{}
	}
	return Clay__Array_Add(&currentContext.AspectRatioElementConfigs, config)
}
func Clay__StoreFloatingElementConfig(config Clay_FloatingElementConfig) *Clay_FloatingElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_FloatingElementConfig{}
	}
	return Clay__Array_Add(&currentContext.FloatingElementConfigs, config)
}

func Clay__StoreCustomElementConfig(config Clay_CustomElementConfig) *Clay_CustomElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_CustomElementConfig{}
	}
	return Clay__Array_Add(&currentContext.CustomElementConfigs, config)
}

func Clay__StoreClipElementConfig(config Clay_ClipElementConfig) *Clay_ClipElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_ClipElementConfig{}
	}
	return Clay__Array_Add(&currentContext.ClipElementConfigs, config)
}

func Clay__StoreBorderElementConfig(config Clay_BorderElementConfig) *Clay_BorderElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_BorderElementConfig{}
	}
	return Clay__Array_Add(&currentContext.BorderElementConfigs, config)
}

func Clay__ConfigureOpenElement(elementDeclaration Clay_ElementDeclaration) {
	currentContext := Clay_GetCurrentContext()
	openLayoutElement := Clay__GetOpenLayoutElement()
	openLayoutElement.LayoutConfig = Clay__StoreLayoutConfig(elementDeclaration.Layout)

	if elementDeclaration.Layout.Sizing.Width.Type == CLAY__SIZING_TYPE_PERCENT && elementDeclaration.Layout.Sizing.Width.Size.Percent > 1 || elementDeclaration.Layout.Sizing.Height.Type == CLAY__SIZING_TYPE_PERCENT && elementDeclaration.Layout.Sizing.Height.Size.Percent > 1 {
		currentContext.ErrorHandler.ErrorHandlerFunction(Clay_ErrorData{
			ErrorType: CLAY_ERROR_TYPE_PERCENTAGE_OVER_1,
			ErrorText: CLAY_STRING("An element was configured with CLAY_SIZING_PERCENT, but the provided percentage value was over 1.0. Clay expects a value between 0 and 1, i.e. 20% is 0.2."),
			UserData:  currentContext.ErrorHandler.UserData,
		})
	}

	//get a lice of the next available slot in the element configs array
	nextAvailableElementConfigIndex := currentContext.ElementConfigs.Length
	fmt.Println("nextAvailableElementConfigIndex", nextAvailableElementConfigIndex)
	elementConfigsSegmentView := mem.MArray_GetSlice(&currentContext.ElementConfigs, nextAvailableElementConfigIndex, nextAvailableElementConfigIndex+1)
	openLayoutElement.ElementConfigs.InternalArray = elementConfigsSegmentView

	var sharedConfig *Clay_SharedElementConfig = nil

	if elementDeclaration.BackgroundColor.A > 0 {
		sharedConfig = new(Clay_SharedElementConfig)
		sharedConfig.BackgroundColor = elementDeclaration.BackgroundColor
		Clay__AttachElementConfig(Clay_ElementConfigUnion{SharedElementConfig: sharedConfig}, CLAY__ELEMENT_CONFIG_TYPE_SHARED)
	}
	if !Clay__MemCmpTyped(&elementDeclaration.CornerRadius, &Clay_CornerRadius{}) {
		if sharedConfig != nil {
			sharedConfig.CornerRadius = elementDeclaration.CornerRadius
		} else {
			sharedConfig = new(Clay_SharedElementConfig)
			sharedConfig.CornerRadius = elementDeclaration.CornerRadius
			Clay__AttachElementConfig(Clay_ElementConfigUnion{SharedElementConfig: sharedConfig}, CLAY__ELEMENT_CONFIG_TYPE_SHARED)
		}
	}

	if elementDeclaration.UserData != nil {
		if sharedConfig != nil {
			sharedConfig.UserData = elementDeclaration.UserData
		} else {
			sharedConfig = Clay__StoreSharedElementConfig(Clay_SharedElementConfig{UserData: elementDeclaration.UserData})
			Clay__AttachElementConfig(Clay_ElementConfigUnion{SharedElementConfig: sharedConfig}, CLAY__ELEMENT_CONFIG_TYPE_SHARED)
		}
	}

	if elementDeclaration.Image.ImageData != nil {
		Clay__AttachElementConfig(Clay_ElementConfigUnion{ImageElementConfig: Clay__StoreImageElementConfig(elementDeclaration.Image)}, CLAY__ELEMENT_CONFIG_TYPE_IMAGE)
	}
	if elementDeclaration.AspectRatio.AspectRatio > 0 {
		Clay__AttachElementConfig(Clay_ElementConfigUnion{AspectRatioElementConfig: Clay__StoreAspectRatioElementConfig(elementDeclaration.AspectRatio)}, CLAY__ELEMENT_CONFIG_TYPE_ASPECT)
		Clay__Array_Add(&currentContext.AspectRatioElementIndexes, currentContext.LayoutElements.Length-1)
	}

	if elementDeclaration.Floating.AttachTo != CLAY_ATTACH_TO_NONE {
		floatingConfig := elementDeclaration.Floating
		// This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here

		hierarchicalParent := Clay__Array_Get[Clay_LayoutElement](&currentContext.LayoutElements, Clay__Array_GetValue[int32](&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-2))
		if hierarchicalParent != nil {
			var clipElementId int32 = 0
			if elementDeclaration.Floating.AttachTo == CLAY_ATTACH_TO_PARENT {
				// Attach to the element's direct hierarchical parent
				floatingConfig.ParentId = hierarchicalParent.Id
				if currentContext.OpenClipElementStack.Length > 0 {
					clipElementId = Clay__Array_GetValue(&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length-1)
				} else if elementDeclaration.Floating.AttachTo == CLAY_ATTACH_TO_ELEMENT_WITH_ID {
					parentItem := Clay__GetHashMapItem(floatingConfig.ParentId)
					// check if parentItem is pointing to the default item
					defaultItem := &Clay_LayoutElementHashMapItem_DEFAULT
					if parentItem == defaultItem {
						currentContext.ErrorHandler.ErrorHandlerFunction(Clay_ErrorData{
							ErrorType: CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND,
							ErrorText: CLAY_STRING("A floating element was declared with a parentId, but no element with that ID was found."),
							UserData:  currentContext.ErrorHandler.UserData,
						})
					} else {
						var clipItemIndex int32 = -1
						for i, elem := range mem.MArray_GetAll(&currentContext.LayoutElements) {
							if &elem == parentItem.LayoutElement {
								clipItemIndex = int32(i)
								break
							}
						}
						if clipItemIndex != -1 {
							clipElementId = Clay__Array_GetValue[int32](&currentContext.LayoutElementClipElementIds, clipItemIndex)
						}
					}
				} else if elementDeclaration.Floating.AttachTo == CLAY_ATTACH_TO_ROOT {
					floatingConfig.ParentId = Clay__HashString(CLAY_STRING("Clay__RootContainer"), 0).Id
				}

				if elementDeclaration.Floating.ClipTo == CLAY_CLIP_TO_NONE {
					clipElementId = 0
				}
				currentElementIndex := Clay__Array_GetValue[int32](&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-1)
				Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentElementIndex, clipElementId)
				Clay__Array_Add(&currentContext.OpenClipElementStack, clipElementId)
				Clay__Array_Add(&currentContext.LayoutElementTreeRoots, Clay__LayoutElementTreeRoot{
					LayoutElementIndex: Clay__Array_GetValue[int32](&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-1),
					ParentId:           floatingConfig.ParentId,
					ClipElementId:      uint32(clipElementId),
					ZIndex:             floatingConfig.ZIndex,
				})
				Clay__AttachElementConfig(Clay_ElementConfigUnion{FloatingElementConfig: Clay__StoreFloatingElementConfig(floatingConfig)}, CLAY__ELEMENT_CONFIG_TYPE_FLOATING)
			}
		}
		if elementDeclaration.Custom.CustomData != nil {
			Clay__AttachElementConfig(Clay_ElementConfigUnion{
				CustomElementConfig: Clay__StoreCustomElementConfig(elementDeclaration.Custom),
			}, CLAY__ELEMENT_CONFIG_TYPE_CUSTOM)
		}
	}

	if elementDeclaration.Clip.Horizontal || elementDeclaration.Clip.Vertical {
		Clay__AttachElementConfig(Clay_ElementConfigUnion{
			ClipElementConfig: Clay__StoreClipElementConfig(elementDeclaration.Clip),
		}, CLAY__ELEMENT_CONFIG_TYPE_CLIP)

		Clay__Array_Add(&currentContext.OpenClipElementStack, int32(openLayoutElement.Id))
		// Retrieve or create cached data to track scroll position across frames
		var scrollOffset *Clay__ScrollContainerDataInternal = nil
		for i := int32(0); i < currentContext.ScrollContainerDatas.Length; i++ {
			mapping := Clay__Array_Get[Clay__ScrollContainerDataInternal](&currentContext.ScrollContainerDatas, i)
			if openLayoutElement.Id == mapping.ElementId {
				scrollOffset = mapping
				scrollOffset.LayoutElement = openLayoutElement
				scrollOffset.OpenThisFrame = true
			}
		}
		if scrollOffset == nil {
			scrollOffset = Clay__Array_Add(&currentContext.ScrollContainerDatas, Clay__ScrollContainerDataInternal{
				LayoutElement: openLayoutElement,
				ScrollOrigin:  Clay_Vector2{-1, -1},
				ElementId:     openLayoutElement.Id,
				OpenThisFrame: true})
		}
		if currentContext.ExternalScrollHandlingEnabled {
			scrollOffset.ScrollPosition = Clay__QueryScrollOffset(scrollOffset.ElementId, currentContext.QueryScrollOffsetUserData)
		}
	}
	if !Clay__MemCmpTyped(&elementDeclaration.Border.Width, &Clay_BorderWidth{}) {
		Clay__AttachElementConfig(Clay_ElementConfigUnion{
			BorderElementConfig: Clay__StoreBorderElementConfig(elementDeclaration.Border),
		}, CLAY__ELEMENT_CONFIG_TYPE_BORDER)
	}
}

func Clay__GetOpenLayoutElement() *Clay_LayoutElement {
	currentContext := Clay_GetCurrentContext()
	return Clay__Array_Get[Clay_LayoutElement](&currentContext.LayoutElements, Clay__Array_GetValue[int32](&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-1))

	// Clay_LayoutElement* Clay__GetOpenLayoutElement(void) {
	//     Clay_Context* context = Clay_GetCurrentContext();
	//     return Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 1));
	// }

}
func Clay__MeasureTextCached(text *Clay_String, textConfig *Clay_TextElementConfig) *Clay__MeasureTextCacheItem {
	panic("not implemented")
}

func Clay__AddHashMapItem(elementId Clay_ElementId, layoutElement *Clay_LayoutElement) *Clay_LayoutElementHashMapItem {
	currentContext := Clay_GetCurrentContext()
	if currentContext.LayoutElementsHashMapInternal.Length == currentContext.LayoutElementsHashMapInternal.Capacity-1 {
		return nil
	}
	item := Clay_LayoutElementHashMapItem{
		ElementId:     elementId,
		LayoutElement: layoutElement,
		NextIndex:     -1,
		Generation:    currentContext.Generation + 1,
	}

	// Perform modulo with uint32 first to avoid negative results, then cast to int32
	hashBucket := int32(elementId.Id % uint32(currentContext.LayoutElementsHashMap.Capacity))
	hashItemPrevious := int32(-1)
	hashItemIndex := Clay__Array_GetValue(&currentContext.LayoutElementsHashMap, hashBucket)
	for hashItemIndex != -1 { // Just replace collision, not a big deal - leave it up to the end user
		hashItem := Clay__Array_Get[Clay_LayoutElementHashMapItem](&currentContext.LayoutElementsHashMapInternal, hashItemIndex)
		if hashItem.ElementId.Id == elementId.Id { // Collision - resolve based on generation
			item.NextIndex = hashItem.NextIndex
			if hashItem.Generation <= currentContext.Generation { // First collision - assume this is the "same" element
				hashItem.ElementId = elementId // Make sure to copy this across. If the stringId reference has changed, we should update the hash item to use the new one.
				hashItem.Generation = currentContext.Generation + 1
				hashItem.LayoutElement = layoutElement
				hashItem.DebugData.Collision = false
				hashItem.OnHoverFunction = nil
				hashItem.HoverFunctionUserData = 0
			} else { // Multiple collisions this frame - two elements have the same ID
				currentContext.ErrorHandler.ErrorHandlerFunction(Clay_ErrorData{
					ErrorType: CLAY_ERROR_TYPE_DUPLICATE_ID,
					ErrorText: CLAY_STRING("An element with this ID was already previously declared during this layout."),
					UserData:  currentContext.ErrorHandler.UserData,
				})
				if currentContext.DebugModeEnabled {
					hashItem.DebugData.Collision = true
				}
			}
			return hashItem
		}
		hashItemPrevious = hashItemIndex
		hashItemIndex = hashItem.NextIndex
	}

	hashItem := Clay__Array_Add(&currentContext.LayoutElementsHashMapInternal, item)
	hashItem.DebugData = Clay__Array_Add(&currentContext.DebugElementData, Clay__DebugElementData{})
	if hashItemPrevious != -1 {
		Clay__Array_Get[Clay_LayoutElementHashMapItem](&currentContext.LayoutElementsHashMapInternal, hashItemPrevious).NextIndex = currentContext.LayoutElementsHashMapInternal.Length - 1
	} else {
		Clay__Array_Set(&currentContext.LayoutElementsHashMap, hashBucket, currentContext.LayoutElementsHashMapInternal.Length-1)
	}
	return hashItem
}

func Clay__OpenTextElement(text Clay_String, textConfig *Clay_TextElementConfig) {
	currentContext := Clay_GetCurrentContext()
	if currentContext.LayoutElements.Length == currentContext.LayoutElements.Capacity-1 || currentContext.BooleanWarnings.MaxElementsExceeded {
		currentContext.BooleanWarnings.MaxElementsExceeded = true
		return
	}
	parentElement := Clay__GetOpenLayoutElement()

	layoutElement := Clay_LayoutElement{}

	textElement := Clay__Array_Add[Clay_LayoutElement](&currentContext.LayoutElements, layoutElement)

	if currentContext.OpenClipElementStack.Length > 0 {
		Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length-1, Clay__Array_GetValue[int32](&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length-1))
	} else {
		Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length-1, 0)
	}

	Clay__Array_Add(&currentContext.LayoutElementChildrenBuffer, currentContext.LayoutElements.Length-1)

	textMeasured := Clay__MeasureTextCached(&text, textConfig)

	elementId := Clay__HashNumber(uint32(parentElement.ChildrenOrTextContent.Children.Length), parentElement.Id)

	textElement.Id = elementId.Id

	Clay__AddHashMapItem(elementId, textElement)
	Clay__Array_Add(&currentContext.LayoutElementIdStrings, elementId.StringId)

	// Clay_Dimensions textDimensions = { .width = textMeasured->unwrappedDimensions.width, .height = textConfig->lineHeight > 0 ? (float)textConfig->lineHeight : textMeasured->unwrappedDimensions.height };

	textDimensions := Clay_Dimensions{
		Width:  textMeasured.UnwrappedDimensions.Width,
		Height: textMeasured.UnwrappedDimensions.Height,
	}

	if textConfig.LineHeight > 0 {
		textDimensions.Height = float32(textConfig.LineHeight)
	}

	textElement.Dimensions = textDimensions

	textElement.MinDimensions = Clay_Dimensions{
		Width:  textMeasured.MinWidth,
		Height: textDimensions.Height,
	}

	textElement.ChildrenOrTextContent.TextElementData = Clay__Array_Add(&currentContext.TextElementData, Clay__TextElementData{
		Text:                text,
		PreferredDimensions: textMeasured.UnwrappedDimensions,
		ElementIndex:        currentContext.LayoutElements.Length - 1,
	})

	// add config to element configs

	config := Clay__Array_Add(&currentContext.ElementConfigs, Clay_ElementConfig{
		Type:   CLAY__ELEMENT_CONFIG_TYPE_TEXT,
		Config: Clay_ElementConfigUnion{TextElementConfig: textConfig},
	})
	if config != nil {
		configIndex := currentContext.ElementConfigs.Length - 1

		segmentView := mem.MArray_GetSlice(&currentContext.ElementConfigs, configIndex, configIndex+1)
		textElement.ElementConfigs = Clay__Slice[Clay_ElementConfig]{
			Length:        1,
			InternalArray: segmentView,
		}
	}
	textElement.LayoutConfig = &Clay_LayoutConfig{}
	parentElement.ChildrenOrTextContent.Children.Length++
}

type Clay__MeasureTextCacheItem struct {
	UnwrappedDimensions     Clay_Dimensions
	MeasuredWordsStartIndex int32
	MinWidth                float32
	ContainsNewlines        bool
	// Hash map data
	Id         uint32
	NextIndex  int32
	Generation uint32
}

func Clay__InitializePersistentMemory(context *Clay_Context) {
	// Persistent memory - initialized once and not reset
	maxElementCount := context.MaxElementCount
	maxMeasureTextCacheWordCount := context.MaxMeasureTextCacheWordCount
	arena := &context.InternalArena

	context.ScrollContainerDatas = Clay__Array_Allocate_Arena[Clay__ScrollContainerDataInternal](100, arena)
	context.LayoutElementsHashMapInternal = Clay__Array_Allocate_Arena[Clay_LayoutElementHashMapItem](maxElementCount, arena)
	context.LayoutElementsHashMap = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.MeasureTextHashMapInternal = Clay__Array_Allocate_Arena[Clay__MeasureTextCacheItem](maxElementCount, arena)
	context.MeasureTextHashMapInternalFreeList = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.MeasuredWordsFreeList = Clay__Array_Allocate_Arena[int32](maxMeasureTextCacheWordCount, arena)
	context.MeasureTextHashMap = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.MeasuredWords = Clay__Array_Allocate_Arena[Clay__MeasuredWord](maxMeasureTextCacheWordCount, arena)
	context.PointerOverIds = Clay__Array_Allocate_Arena[Clay_ElementId](maxElementCount, arena)
	context.DebugElementData = Clay__Array_Allocate_Arena[Clay__DebugElementData](maxElementCount, arena)
	context.ArenaResetOffset = arena.NextAllocation
}

func Clay__InitializeEphemeralMemory(context *Clay_Context) {
	maxElementCount := context.MaxElementCount
	// Ephemeral Memory - reset every frame
	arena := &context.InternalArena
	arena.NextAllocation = context.ArenaResetOffset

	context.LayoutElementChildrenBuffer = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.LayoutElements = Clay__Array_Allocate_Arena[Clay_LayoutElement](maxElementCount, arena)
	context.Warnings = Clay__Array_Allocate_Arena[Clay__Warning](100, arena)

	context.LayoutConfigs = Clay__Array_Allocate_Arena[Clay_LayoutConfig](maxElementCount, arena)
	context.ElementConfigs = Clay__Array_Allocate_Arena[Clay_ElementConfig](maxElementCount, arena)
	context.TextElementConfigs = Clay__Array_Allocate_Arena[Clay_TextElementConfig](maxElementCount, arena)
	context.AspectRatioElementConfigs = Clay__Array_Allocate_Arena[Clay_AspectRatioElementConfig](maxElementCount, arena)
	context.ImageElementConfigs = Clay__Array_Allocate_Arena[Clay_ImageElementConfig](maxElementCount, arena)
	context.FloatingElementConfigs = Clay__Array_Allocate_Arena[Clay_FloatingElementConfig](maxElementCount, arena)
	context.ClipElementConfigs = Clay__Array_Allocate_Arena[Clay_ClipElementConfig](maxElementCount, arena)
	context.CustomElementConfigs = Clay__Array_Allocate_Arena[Clay_CustomElementConfig](maxElementCount, arena)
	context.BorderElementConfigs = Clay__Array_Allocate_Arena[Clay_BorderElementConfig](maxElementCount, arena)
	context.SharedElementConfigs = Clay__Array_Allocate_Arena[Clay_SharedElementConfig](maxElementCount, arena)

	context.LayoutElementIdStrings = Clay__Array_Allocate_Arena[Clay_String](maxElementCount, arena)
	context.WrappedTextLines = Clay__Array_Allocate_Arena[Clay__WrappedTextLine](maxElementCount, arena)
	context.LayoutElementTreeNodeArray1 = Clay__Array_Allocate_Arena[Clay__LayoutElementTreeNode](maxElementCount, arena)
	context.LayoutElementTreeRoots = Clay__Array_Allocate_Arena[Clay__LayoutElementTreeRoot](maxElementCount, arena)
	context.LayoutElementChildren = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.OpenLayoutElementStack = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.TextElementData = Clay__Array_Allocate_Arena[Clay__TextElementData](maxElementCount, arena)
	context.AspectRatioElementIndexes = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.RenderCommands = Clay__Array_Allocate_Arena[Clay_RenderCommand](maxElementCount, arena)
	context.TreeNodeVisited = Clay__Array_Allocate_Arena[bool](maxElementCount, arena)
	context.TreeNodeVisited.Length = context.TreeNodeVisited.Capacity // This array is accessed directly rather than behaving as a list
	context.OpenClipElementStack = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.ReusableElementIndexBuffer = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.LayoutElementClipElementIds = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.DynamicStringData = Clay__Array_Allocate_Arena[byte](maxElementCount, arena)
}

func Clay__Context_Allocate_Arena(arena *Clay_Arena) *Clay_Context {
	clay_Context, err := mem.AllocateStruct[Clay_Context](arena)
	if err != nil {
		return nil
	}
	return clay_Context
}
