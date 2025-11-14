package clay

import "github.com/zodimo/clay-go/pkg/mem"

func Clay__CalculateFinalLayout() {
	currentContext := Clay_GetCurrentContext()
	// Calculate sizing along the X axis
	Clay__SizeContainersAlongAxis(true)

	// Wrap text
	for textElementIndex := int32(0); textElementIndex < currentContext.TextElementData.Length(); textElementIndex++ {
		textElementData := Clay__Array_Get(&currentContext.TextElementData, textElementIndex)
		// Initialize wrappedLines to point to the current end of the WrappedTextLines array
		// This allows each text element to have its own slice that grows as lines are added
		wrappedLinesStartIndex := currentContext.WrappedTextLines.Length()
		wrappedLinesData := mem.MArray_GetSlice(&currentContext.WrappedTextLines, wrappedLinesStartIndex, wrappedLinesStartIndex)
		wrappedLines := NewClay__Slice[Clay__WrappedTextLine](wrappedLinesData)
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
			Clay__Slice_Grow(&textElementData.WrappedLines, 1)
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
			if currentContext.WrappedTextLines.Length() > currentContext.WrappedTextLines.Capacity()-1 {
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
				Clay__Slice_Grow(&textElementData.WrappedLines, 1)
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

				Clay__Slice_Grow(&textElementData.WrappedLines, 1)
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
			Clay__Slice_Grow(&textElementData.WrappedLines, 1)
		}

		containerElement.Dimensions.Height = lineHeight * float32(textElementData.WrappedLines.Length())
	}

	// Scale vertical heights according to aspect ratio
	for aspectRatioElementIndex := int32(0); aspectRatioElementIndex < currentContext.AspectRatioElementIndexes.Length(); aspectRatioElementIndex++ {
		aspectElement := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&currentContext.AspectRatioElementIndexes, aspectRatioElementIndex))
		aspectRatioElementConfig := Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).AspectRatioElementConfig
		aspectElement.Dimensions.Height = (1 / aspectRatioElementConfig.AspectRatio) * aspectElement.Dimensions.Width
		aspectElement.LayoutConfig.Sizing.Height.Size.MinMax.Max = aspectElement.Dimensions.Height
	}

	// Propagate effect of text wrapping, aspect scaling etc. on height of parents
	dfsBuffer := currentContext.LayoutElementTreeNodeArray1
	Clay__Array_Reset(&dfsBuffer)
	for layoutElementTreeRootIndex := int32(0); layoutElementTreeRootIndex < currentContext.LayoutElementTreeRoots.Length(); layoutElementTreeRootIndex++ {
		layoutElementTreeRoot := Clay__Array_Get(&currentContext.LayoutElementTreeRoots, layoutElementTreeRootIndex)
		if currentContext.TreeNodeVisited.Length() <= dfsBuffer.Length() {
			Clay__Array_Add(&currentContext.TreeNodeVisited, false)
		} else {
			Clay__Array_Set(&currentContext.TreeNodeVisited, dfsBuffer.Length(), false)
		}
		Clay__Array_Add(&dfsBuffer, Clay__LayoutElementTreeNode{LayoutElement: Clay__Array_Get(&currentContext.LayoutElements, layoutElementTreeRoot.LayoutElementIndex)})
	}

	for dfsBuffer.Length() > 0 {
		currentElementTreeNode := Clay__Array_Get(&dfsBuffer, dfsBuffer.Length()-1)
		currentElement := currentElementTreeNode.LayoutElement
		if !Clay__Array_GetValue(&currentContext.TreeNodeVisited, dfsBuffer.Length()-1) {
			Clay__Array_Set(&currentContext.TreeNodeVisited, dfsBuffer.Length()-1, true)
			// If the element has no children or is the container for a text element, don't bother inspecting it
			if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) || currentElement.ChildrenOrTextContent.Children.Length == 0 {
				Clay__Array_Shrink(&dfsBuffer, 1)
				continue
			}
			// Add the children to the DFS buffer (needs to be pushed in reverse so that stack traversal is in correct layout order)
			for childIndex := int32(0); childIndex < int32(currentElement.ChildrenOrTextContent.Children.Length); childIndex++ {

				if currentContext.TreeNodeVisited.Length() <= dfsBuffer.Length() {
					Clay__Array_Add(&currentContext.TreeNodeVisited, false)
				} else {
					Clay__Array_Set(&currentContext.TreeNodeVisited, dfsBuffer.Length(), false)
				}
				Clay__Array_Add(&dfsBuffer, Clay__LayoutElementTreeNode{
					LayoutElement: Clay__Array_Get(
						&currentContext.LayoutElements,
						currentElement.ChildrenOrTextContent.Children.Elements[childIndex],
					),
				},
				)
			}
			continue
		}
		Clay__Array_Shrink(&dfsBuffer, 1)

		// DFS node has been visited, this is on the way back up to the root
		layoutConfig := currentElement.LayoutConfig
		if layoutConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT {
			// Resize any parent containers that have grown in height along their non layout axis
			for childIndex := int32(0); childIndex < int32(currentElement.ChildrenOrTextContent.Children.Length); childIndex++ {
				childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements[childIndex])
				childHeightWithPadding := CLAY__MAX(childElement.Dimensions.Height+float32(layoutConfig.Padding.Top)+float32(layoutConfig.Padding.Bottom), currentElement.Dimensions.Height)
				currentElement.Dimensions.Height = CLAY__MIN(CLAY__MAX(childHeightWithPadding, layoutConfig.Sizing.Height.Size.MinMax.Min), layoutConfig.Sizing.Height.Size.MinMax.Max)
			}
		} else if layoutConfig.LayoutDirection == CLAY_TOP_TO_BOTTOM {
			// Resizing along the layout axis
			contentHeight := float32(layoutConfig.Padding.Top + layoutConfig.Padding.Bottom)
			for childIndex := int32(0); childIndex < int32(currentElement.ChildrenOrTextContent.Children.Length); childIndex++ {
				childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements[childIndex])
				contentHeight += childElement.Dimensions.Height
			}
			contentHeight += float32(CLAY__MAX(int32(currentElement.ChildrenOrTextContent.Children.Length)-1, 0) * int32(layoutConfig.ChildGap))
			currentElement.Dimensions.Height = CLAY__MIN(CLAY__MAX(contentHeight, layoutConfig.Sizing.Height.Size.MinMax.Min), layoutConfig.Sizing.Height.Size.MinMax.Max)
		}
	}

	// Calculate sizing along the Y axis
	Clay__SizeContainersAlongAxis(false)

	// Scale horizontal widths according to aspect ratio
	for aspectRatioElementIndex := int32(0); aspectRatioElementIndex < currentContext.AspectRatioElementIndexes.Length(); aspectRatioElementIndex++ {
		aspectElement := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&currentContext.AspectRatioElementIndexes, aspectRatioElementIndex))
		aspectRatioElementConfig := Clay__FindElementConfigWithType(aspectElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).AspectRatioElementConfig
		aspectElement.Dimensions.Width = aspectRatioElementConfig.AspectRatio * aspectElement.Dimensions.Height
	}

	// Sort tree roots by z-index
	sortMax := currentContext.LayoutElementTreeRoots.Length() - 1
	for sortMax > 0 { // todo dumb bubble sort
		for i := int32(0); i < sortMax; i++ {
			current := Clay__Array_GetValue(&currentContext.LayoutElementTreeRoots, i)
			next := Clay__Array_GetValue(&currentContext.LayoutElementTreeRoots, i+1)
			if next.ZIndex < current.ZIndex {
				Clay__Array_Set(&currentContext.LayoutElementTreeRoots, i, next)
				Clay__Array_Set(&currentContext.LayoutElementTreeRoots, i+1, current)
			}
		}
		sortMax--
	}

	// Calculate final positions and generate render commands
	Clay__Array_Reset(&currentContext.RenderCommands)
	Clay__Array_Reset(&dfsBuffer)
	for rootIndex := int32(0); rootIndex < currentContext.LayoutElementTreeRoots.Length(); rootIndex++ {
		Clay__Array_Reset(&dfsBuffer)
		root := Clay__Array_GetValue(&currentContext.LayoutElementTreeRoots, rootIndex)
		rootElement := Clay__Array_Get(&currentContext.LayoutElements, root.LayoutElementIndex)
		rootPosition := Clay_Vector2{}
		parentHashMapItem := Clay__GetHashMapItem(root.ParentId)
		// Position root floating containers
		if Clay__ElementHasConfig(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) && parentHashMapItem != nil {
			config := Clay__FindElementConfigWithType(rootElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).FloatingElementConfig
			rootDimensions := rootElement.Dimensions
			parentBoundingBox := parentHashMapItem.BoundingBox
			// Set X position
			targetAttachPosition := Clay_Vector2{}
			switch config.AttachPoints.Parent {
			case CLAY_ATTACH_POINT_LEFT_TOP:
			case CLAY_ATTACH_POINT_LEFT_CENTER:
			case CLAY_ATTACH_POINT_LEFT_BOTTOM:
				targetAttachPosition.X = parentBoundingBox.X
				break
			case CLAY_ATTACH_POINT_CENTER_TOP:
			case CLAY_ATTACH_POINT_CENTER_CENTER:
			case CLAY_ATTACH_POINT_CENTER_BOTTOM:
				targetAttachPosition.X = parentBoundingBox.X + (parentBoundingBox.Width / 2)
				break
			case CLAY_ATTACH_POINT_RIGHT_TOP:
			case CLAY_ATTACH_POINT_RIGHT_CENTER:
			case CLAY_ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.X = parentBoundingBox.X + parentBoundingBox.Width
				break
			}
			switch config.AttachPoints.Element {
			case CLAY_ATTACH_POINT_LEFT_TOP:
			case CLAY_ATTACH_POINT_LEFT_CENTER:
			case CLAY_ATTACH_POINT_LEFT_BOTTOM:
				break
			case CLAY_ATTACH_POINT_CENTER_TOP:
			case CLAY_ATTACH_POINT_CENTER_CENTER:
			case CLAY_ATTACH_POINT_CENTER_BOTTOM:
				targetAttachPosition.X -= (rootDimensions.Width / 2)
				break
			case CLAY_ATTACH_POINT_RIGHT_TOP:
			case CLAY_ATTACH_POINT_RIGHT_CENTER:
			case CLAY_ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.X -= rootDimensions.Width
				break
			}
			switch config.AttachPoints.Parent { // I know I could merge the x and y switch statements, but this is easier to read
			case CLAY_ATTACH_POINT_LEFT_TOP:
			case CLAY_ATTACH_POINT_RIGHT_TOP:
			case CLAY_ATTACH_POINT_CENTER_TOP:
				targetAttachPosition.Y = parentBoundingBox.Y
				break
			case CLAY_ATTACH_POINT_LEFT_CENTER:
			case CLAY_ATTACH_POINT_CENTER_CENTER:
			case CLAY_ATTACH_POINT_RIGHT_CENTER:
				targetAttachPosition.Y = parentBoundingBox.Y + (parentBoundingBox.Height / 2)
				break
			case CLAY_ATTACH_POINT_LEFT_BOTTOM:
			case CLAY_ATTACH_POINT_CENTER_BOTTOM:
			case CLAY_ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.Y = parentBoundingBox.Y + parentBoundingBox.Height
				break
			}
			switch config.AttachPoints.Element {
			case CLAY_ATTACH_POINT_LEFT_TOP:
			case CLAY_ATTACH_POINT_RIGHT_TOP:
			case CLAY_ATTACH_POINT_CENTER_TOP:
				break
			case CLAY_ATTACH_POINT_LEFT_CENTER:
			case CLAY_ATTACH_POINT_CENTER_CENTER:
			case CLAY_ATTACH_POINT_RIGHT_CENTER:
				targetAttachPosition.Y -= (rootDimensions.Height / 2)
				break
			case CLAY_ATTACH_POINT_LEFT_BOTTOM:
			case CLAY_ATTACH_POINT_CENTER_BOTTOM:
			case CLAY_ATTACH_POINT_RIGHT_BOTTOM:
				targetAttachPosition.Y -= rootDimensions.Height
				break
			}
			targetAttachPosition.X += config.Offset.X
			targetAttachPosition.Y += config.Offset.Y
			rootPosition = targetAttachPosition
		}
		if root.ClipElementId != 0 {
			clipHashMapItem := Clay__GetHashMapItem(root.ClipElementId)
			if clipHashMapItem != nil {
				// Floating elements that are attached to scrolling contents won't be correctly positioned if external scroll handling is enabled, fix here
				if currentContext.ExternalScrollHandlingEnabled {
					clipConfig := Clay__FindElementConfigWithType(clipHashMapItem.LayoutElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).ClipElementConfig
					if clipConfig.Horizontal {
						rootPosition.X += clipConfig.ChildOffset.X
					}
					if clipConfig.Vertical {
						rootPosition.Y += clipConfig.ChildOffset.Y
					}
				}
				Clay__AddRenderCommand(Clay_RenderCommand{
					BoundingBox: clipHashMapItem.BoundingBox,
					UserData:    0,
					Id:          Clay__HashNumber(rootElement.Id, uint32(rootElement.ChildrenOrTextContent.Children.Length+10)).Id, // TODO need a better strategy for managing derived ids
					ZIndex:      root.ZIndex,
					CommandType: CLAY_RENDER_COMMAND_TYPE_SCISSOR_START,
				})
			}
		}
		Clay__Array_Add(&dfsBuffer, Clay__LayoutElementTreeNode{
			LayoutElement: rootElement,
			Position:      rootPosition,
			NextChildOffset: Clay_Vector2{
				X: float32(rootElement.LayoutConfig.Padding.Left),
				Y: float32(rootElement.LayoutConfig.Padding.Top),
			},
		},
		)

		Clay__Array_Set(&currentContext.TreeNodeVisited, 0, false)
		for dfsBuffer.Length() > 0 {
			currentElementTreeNode := Clay__Array_GetValue(&dfsBuffer, dfsBuffer.Length()-1)
			currentElement := currentElementTreeNode.LayoutElement
			layoutConfig := currentElement.LayoutConfig
			scrollOffset := Clay_Vector2{}

			// This will only be run a single time for each element in downwards DFS order
			if !Clay__Array_GetValue(&currentContext.TreeNodeVisited, dfsBuffer.Length()-1) {
				Clay__Array_Set(&currentContext.TreeNodeVisited, dfsBuffer.Length()-1, true)

				currentElementBoundingBox := Clay_BoundingBox{
					X:      currentElementTreeNode.Position.X,
					Y:      currentElementTreeNode.Position.Y,
					Width:  currentElement.Dimensions.Width,
					Height: currentElement.Dimensions.Height,
				}
				if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING) {
					floatingElementConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_FLOATING).FloatingElementConfig
					expand := floatingElementConfig.Expand
					currentElementBoundingBox.X -= expand.Width
					currentElementBoundingBox.Width += expand.Width * 2
					currentElementBoundingBox.Y -= expand.Height
					currentElementBoundingBox.Height += expand.Height * 2
				}

				scrollContainerData := new(Clay__ScrollContainerDataInternal)
				// Apply scroll offsets to container
				if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP) {
					clipConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).ClipElementConfig

					// This linear scan could theoretically be slow under very strange conditions, but I can't imagine a real UI with more than a few 10's of scroll containers
					for i := int32(0); i < currentContext.ScrollContainerDatas.Length(); i++ {
						mapping := Clay__Array_Get(&currentContext.ScrollContainerDatas, i)
						if mapping.LayoutElement == currentElement {
							scrollContainerData = mapping
							mapping.BoundingBox = currentElementBoundingBox
							scrollOffset = clipConfig.ChildOffset
							if currentContext.ExternalScrollHandlingEnabled {
								scrollOffset = Clay_Vector2{}
							}
							break
						}
					}
				}

				hashMapItem := Clay__GetHashMapItem(currentElement.Id)
				if hashMapItem != nil {
					hashMapItem.BoundingBox = currentElementBoundingBox
				}

				sortedConfigIndexes := make([]int32, currentElement.ElementConfigs.Length())
				for elementConfigIndex := int32(0); elementConfigIndex < currentElement.ElementConfigs.Length(); elementConfigIndex++ {
					sortedConfigIndexes[elementConfigIndex] = elementConfigIndex
				}
				sortMax := currentElement.ElementConfigs.Length() - 1
				for sortMax > 0 { // todo dumb bubble sort
					for i := int32(0); i < sortMax; i++ {
						current := sortedConfigIndexes[i]
						next := sortedConfigIndexes[i+1]
						currentType := Clay__Slice_Get(&currentElement.ElementConfigs, current).Type
						nextType := Clay__Slice_Get(&currentElement.ElementConfigs, next).Type
						if nextType == CLAY__ELEMENT_CONFIG_TYPE_CLIP || currentType == CLAY__ELEMENT_CONFIG_TYPE_BORDER {
							sortedConfigIndexes[i] = next
							sortedConfigIndexes[i+1] = current
						}
					}
					sortMax--
				}

				emitRectangle := false
				// Create the render commands for this element
				sharedConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).SharedElementConfig
				if sharedConfig != nil && sharedConfig.BackgroundColor.A > 0 {
					emitRectangle = true
				} else if sharedConfig == nil {
					emitRectangle = false
					sharedConfig = &Clay_SharedElementConfig_DEFAULT
				}
				for elementConfigIndex := int32(0); elementConfigIndex < currentElement.ElementConfigs.Length(); elementConfigIndex++ {
					elementConfig := Clay__Slice_Get(&currentElement.ElementConfigs, sortedConfigIndexes[elementConfigIndex])
					renderCommand := Clay_RenderCommand{
						BoundingBox: currentElementBoundingBox,
						UserData:    sharedConfig.UserData,
						Id:          currentElement.Id,
					}

					offscreen := Clay__ElementIsOffscreen(&currentElementBoundingBox)
					// Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
					shouldRender := !offscreen
					switch elementConfig.Type {
					case CLAY__ELEMENT_CONFIG_TYPE_ASPECT:
						shouldRender = false
						break
					case CLAY__ELEMENT_CONFIG_TYPE_FLOATING:
						shouldRender = false
						break
					case CLAY__ELEMENT_CONFIG_TYPE_SHARED:
						shouldRender = false
						break
					case CLAY__ELEMENT_CONFIG_TYPE_BORDER:
						shouldRender = false
						break
					case CLAY__ELEMENT_CONFIG_TYPE_CLIP:
						{
							renderCommand.CommandType = CLAY_RENDER_COMMAND_TYPE_SCISSOR_START
							renderCommand.RenderData = Clay_RenderData{
								Clip: Clay_ClipRenderData{
									Horizontal: elementConfig.Config.ClipElementConfig.Horizontal,
									Vertical:   elementConfig.Config.ClipElementConfig.Vertical,
								},
							}
							break
						}
					case CLAY__ELEMENT_CONFIG_TYPE_IMAGE:
						{
							renderCommand.CommandType = CLAY_RENDER_COMMAND_TYPE_IMAGE
							renderCommand.RenderData = Clay_RenderData{
								Image: Clay_ImageRenderData{
									BackgroundColor: sharedConfig.BackgroundColor,
									CornerRadius:    sharedConfig.CornerRadius,
									ImageData:       elementConfig.Config.ImageElementConfig.ImageData,
								},
							}
							emitRectangle = false
							break
						}
					case CLAY__ELEMENT_CONFIG_TYPE_TEXT:
						{
							if !shouldRender {
								break
							}
							if currentElement.ChildrenOrTextContent.TextElementData == nil {
								// TextElementData should always be set for text elements, but add safety check
								break
							}
							shouldRender = false
							configUnion := elementConfig.Config
							textElementConfig := configUnion.TextElementConfig
							naturalLineHeight := currentElement.ChildrenOrTextContent.TextElementData.PreferredDimensions.Height

							var finalLineHeight float32
							if textElementConfig.LineHeight > 0 {
								finalLineHeight = float32(textElementConfig.LineHeight)
							} else {
								finalLineHeight = naturalLineHeight
							}
							lineHeightOffset := (finalLineHeight - naturalLineHeight) / 2
							yPosition := lineHeightOffset
							for lineIndex := int32(0); lineIndex < currentElement.ChildrenOrTextContent.TextElementData.WrappedLines.Length(); lineIndex++ {
								wrappedLine := Clay__Slice_Get(&currentElement.ChildrenOrTextContent.TextElementData.WrappedLines, lineIndex)
								if wrappedLine.Line.Length == 0 {
									yPosition += finalLineHeight
									continue
								}
								offset := (currentElementBoundingBox.Width - wrappedLine.Dimensions.Width)
								if textElementConfig.TextAlignment == CLAY_TEXT_ALIGN_LEFT {
									offset = 0
								}
								if textElementConfig.TextAlignment == CLAY_TEXT_ALIGN_CENTER {
									offset /= 2
								}
								Clay__AddRenderCommand(Clay_RenderCommand{
									BoundingBox: Clay_BoundingBox{
										X:      currentElementBoundingBox.X + offset,
										Y:      currentElementBoundingBox.Y + yPosition,
										Width:  wrappedLine.Dimensions.Width,
										Height: wrappedLine.Dimensions.Height,
									},
									RenderData: Clay_RenderData{
										Text: Clay_TextRenderData{
											StringContents: Clay_StringSlice{
												Length:    wrappedLine.Line.Length,
												Chars:     wrappedLine.Line.Chars,
												BaseChars: currentElement.ChildrenOrTextContent.TextElementData.Text.Chars,
											},
											TextColor:     textElementConfig.TextColor,
											FontId:        textElementConfig.FontId,
											FontSize:      textElementConfig.FontSize,
											LetterSpacing: textElementConfig.LetterSpacing,
											LineHeight:    textElementConfig.LineHeight,
										}},
									UserData:    textElementConfig.UserData,
									Id:          Clay__HashNumber(uint32(lineIndex), currentElement.Id).Id,
									ZIndex:      root.ZIndex,
									CommandType: CLAY_RENDER_COMMAND_TYPE_TEXT,
								})
								yPosition += finalLineHeight

								if !currentContext.DisableCulling && (currentElementBoundingBox.Y+yPosition > currentContext.LayoutDimensions.Height) {
									break
								}
							}
							break
						}
					case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM:
						{
							renderCommand.CommandType = CLAY_RENDER_COMMAND_TYPE_CUSTOM
							renderCommand.RenderData = Clay_RenderData{
								Custom: Clay_CustomRenderData{
									BackgroundColor: sharedConfig.BackgroundColor,
									CornerRadius:    sharedConfig.CornerRadius,
									CustomData:      elementConfig.Config.CustomElementConfig.CustomData,
								},
							}
							emitRectangle = false
							break
						}
					default:
						break
					}
					if shouldRender {
						Clay__AddRenderCommand(renderCommand)
					}
					if offscreen {
						// NOTE: You may be tempted to try an early return / continue if an element is off screen. Why bother calculating layout for its children, right?
						// Unfortunately, a FLOATING_CONTAINER may be defined that attaches to a child or grandchild of this element, which is large enough to still
						// be on screen, even if this element isn't. That depends on this element and it's children being laid out correctly (even if they are entirely off screen)
					}
				}

				if emitRectangle {
					Clay__AddRenderCommand(Clay_RenderCommand{
						BoundingBox: currentElementBoundingBox,
						RenderData: Clay_RenderData{
							Rectangle: Clay_RectangleRenderData{
								BackgroundColor: sharedConfig.BackgroundColor,
								CornerRadius:    sharedConfig.CornerRadius,
							},
						},
						UserData:    sharedConfig.UserData,
						Id:          currentElement.Id,
						ZIndex:      root.ZIndex,
						CommandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
					})
				}

				// Setup initial on-axis alignment
				if !Clay__ElementHasConfig(currentElementTreeNode.LayoutElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) {
					contentSize := Clay_Dimensions{Width: 0, Height: 0}
					if layoutConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT {
						for i := int32(0); i < int32(currentElement.ChildrenOrTextContent.Children.Length); i++ {
							childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements[i])
							contentSize.Width += childElement.Dimensions.Width
							contentSize.Height = CLAY__MAX(contentSize.Height, childElement.Dimensions.Height)
						}
						contentSize.Width += float32(CLAY__MAX(currentElement.ChildrenOrTextContent.Children.Length-1, 0) * layoutConfig.ChildGap)
						extraSpace := currentElement.Dimensions.Width - float32(layoutConfig.Padding.Left+layoutConfig.Padding.Right) - contentSize.Width
						switch layoutConfig.ChildAlignment.X {
						case CLAY_ALIGN_X_LEFT:
							extraSpace = 0
						case CLAY_ALIGN_X_CENTER:
							extraSpace /= 2
							break
						default:
							break
						}
						currentElementTreeNode.NextChildOffset.X += extraSpace
						extraSpace = CLAY__MAX(0, extraSpace)
					} else {
						for i := int32(0); i < int32(currentElement.ChildrenOrTextContent.Children.Length); i++ {
							childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements[i])
							contentSize.Width = CLAY__MAX(contentSize.Width, childElement.Dimensions.Width)
							contentSize.Height += childElement.Dimensions.Height
						}
						contentSize.Height += float32(CLAY__MAX(currentElement.ChildrenOrTextContent.Children.Length-1, 0) * layoutConfig.ChildGap)
						extraSpace := currentElement.Dimensions.Height - float32(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom) - contentSize.Height
						switch layoutConfig.ChildAlignment.Y {
						case CLAY_ALIGN_Y_TOP:
							extraSpace = 0
							break
						case CLAY_ALIGN_Y_CENTER:
							extraSpace /= 2
							break
						default:
							break
						}
						extraSpace = CLAY__MAX(0, extraSpace)
						currentElementTreeNode.NextChildOffset.Y += extraSpace
					}

					if scrollContainerData != nil {
						scrollContainerData.ContentSize = Clay_Dimensions{
							Width:  contentSize.Width + float32(layoutConfig.Padding.Left+layoutConfig.Padding.Right),
							Height: contentSize.Height + float32(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom),
						}
					}
				}
			} else {
				// DFS is returning upwards backwards
				var scrollOffset Clay_Vector2
				closeClipElement := false
				clipConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_CLIP).ClipElementConfig
				if clipConfig != nil {
					closeClipElement = true
					for i := int32(0); i < int32(currentContext.ScrollContainerDatas.Length()); i++ {
						mapping := Clay__Array_Get(&currentContext.ScrollContainerDatas, i)
						if mapping.LayoutElement == currentElement {
							scrollOffset = clipConfig.ChildOffset
							if currentContext.ExternalScrollHandlingEnabled {
								scrollOffset = Clay_Vector2{0, 0}
							}
							break
						}
					}
				}

				if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER) {
					currentElementData := Clay__GetHashMapItem(currentElement.Id)
					currentElementBoundingBox := currentElementData.BoundingBox

					// Culling - Don't bother to generate render commands for rectangles entirely outside the screen - this won't stop their children from being rendered if they overflow
					if !Clay__ElementIsOffscreen(&currentElementBoundingBox) {

						var sharedConfig *Clay_SharedElementConfig
						if Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED) {
							sharedConfig = Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_SHARED).SharedElementConfig
						} else {
							sharedConfig = &Clay_SharedElementConfig_DEFAULT
						}

						borderConfig := Clay__FindElementConfigWithType(currentElement, CLAY__ELEMENT_CONFIG_TYPE_BORDER).BorderElementConfig
						renderCommand := Clay_RenderCommand{
							BoundingBox: currentElementBoundingBox,
							RenderData: Clay_RenderData{
								Border: Clay_BorderRenderData{
									Color:        borderConfig.Color,
									CornerRadius: sharedConfig.CornerRadius,
									Width:        borderConfig.Width,
								},
							},
							UserData:    sharedConfig.UserData,
							Id:          Clay__HashNumber(currentElement.Id, uint32(int32(currentElement.ChildrenOrTextContent.Children.Length))).Id,
							CommandType: CLAY_RENDER_COMMAND_TYPE_BORDER,
						}
						Clay__AddRenderCommand(renderCommand)
						if borderConfig.Width.BetweenChildren > 0 && borderConfig.Color.A > 0 {
							halfGap := float32(layoutConfig.ChildGap / 2)
							borderOffset := Clay_Vector2{
								X: float32(layoutConfig.Padding.Left) - halfGap,
								Y: float32(layoutConfig.Padding.Top) - halfGap,
							}
							if layoutConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT {
								for i := int32(0); i < int32(currentElement.ChildrenOrTextContent.Children.Length); i++ {
									childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements[i])
									if i > 0 {
										Clay__AddRenderCommand(Clay_RenderCommand{
											BoundingBox: Clay_BoundingBox{
												X:      currentElementBoundingBox.X + borderOffset.X + scrollOffset.X,
												Y:      currentElementBoundingBox.Y + scrollOffset.Y,
												Width:  float32(borderConfig.Width.BetweenChildren),
												Height: currentElement.Dimensions.Height,
											},
											RenderData: Clay_RenderData{
												Rectangle: Clay_RectangleRenderData{
													BackgroundColor: borderConfig.Color,
												},
											},
											UserData:    sharedConfig.UserData,
											Id:          Clay__HashNumber(currentElement.Id, uint32(int32(currentElement.ChildrenOrTextContent.Children.Length)+1+i)).Id,
											CommandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
										})
									}
									borderOffset.X += (childElement.Dimensions.Width + float32(layoutConfig.ChildGap))
								}
							} else {
								for i := int32(0); i < int32(currentElement.ChildrenOrTextContent.Children.Length); i++ {
									childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements[i])
									if i > 0 {
										Clay__AddRenderCommand(Clay_RenderCommand{
											BoundingBox: Clay_BoundingBox{
												X:      currentElementBoundingBox.X + scrollOffset.X,
												Y:      currentElementBoundingBox.Y + borderOffset.Y + scrollOffset.Y,
												Width:  currentElement.Dimensions.Width,
												Height: float32(borderConfig.Width.BetweenChildren),
											},
											RenderData: Clay_RenderData{
												Rectangle: Clay_RectangleRenderData{
													BackgroundColor: borderConfig.Color,
												},
											},
											UserData:    sharedConfig.UserData,
											Id:          Clay__HashNumber(currentElement.Id, uint32(int32(currentElement.ChildrenOrTextContent.Children.Length)+1+i)).Id,
											CommandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
										})
									}
									borderOffset.Y += (childElement.Dimensions.Height + float32(layoutConfig.ChildGap))
								}
							}
						}
					}
				}
				// This exists because the scissor needs to end _after_ borders between elements
				if closeClipElement {
					Clay__AddRenderCommand(Clay_RenderCommand{
						Id:          Clay__HashNumber(currentElement.Id, uint32(int32(rootElement.ChildrenOrTextContent.Children.Length)+11)).Id,
						CommandType: CLAY_RENDER_COMMAND_TYPE_SCISSOR_END,
					})
				}

				Clay__Array_Shrink(&dfsBuffer, 1)
				continue
			}

			// Add children to the DFS buffer
			if !Clay__ElementHasConfig(currentElement, CLAY__ELEMENT_CONFIG_TYPE_TEXT) {
				Clay__Array_Grow(&dfsBuffer, int32(currentElement.ChildrenOrTextContent.Children.Length))
				for i := int32(0); i < int32(currentElement.ChildrenOrTextContent.Children.Length); i++ {
					childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements[i])
					// Alignment along non layout axis
					if layoutConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT {
						currentElementTreeNode.NextChildOffset.Y = float32(currentElement.LayoutConfig.Padding.Top)
						whiteSpaceAroundChild := currentElement.Dimensions.Height - float32(layoutConfig.Padding.Top+layoutConfig.Padding.Bottom) - childElement.Dimensions.Height
						switch layoutConfig.ChildAlignment.Y {
						case CLAY_ALIGN_Y_TOP:
							break
						case CLAY_ALIGN_Y_CENTER:
							currentElementTreeNode.NextChildOffset.Y += whiteSpaceAroundChild / 2
							break
						case CLAY_ALIGN_Y_BOTTOM:
							currentElementTreeNode.NextChildOffset.Y += whiteSpaceAroundChild
							break
						}
					} else {
						currentElementTreeNode.NextChildOffset.X = float32(currentElement.LayoutConfig.Padding.Left)
						whiteSpaceAroundChild := currentElement.Dimensions.Width - float32(layoutConfig.Padding.Left+layoutConfig.Padding.Right) - childElement.Dimensions.Width
						switch layoutConfig.ChildAlignment.X {
						case CLAY_ALIGN_X_LEFT:
							break
						case CLAY_ALIGN_X_CENTER:
							currentElementTreeNode.NextChildOffset.X += whiteSpaceAroundChild / 2
							break
						case CLAY_ALIGN_X_RIGHT:
							currentElementTreeNode.NextChildOffset.X += whiteSpaceAroundChild
							break
						}
					}

					childPosition := Clay_Vector2{
						X: currentElementTreeNode.Position.X + currentElementTreeNode.NextChildOffset.X + scrollOffset.X,
						Y: currentElementTreeNode.Position.Y + currentElementTreeNode.NextChildOffset.Y + scrollOffset.Y,
					}

					// DFS buffer elements need to be added in reverse because stack traversal happens backwards
					newNodeIndex := dfsBuffer.Length() - 1 - i
					Clay__Array_Set(&dfsBuffer, newNodeIndex, Clay__LayoutElementTreeNode{
						LayoutElement: childElement,
						Position:      Clay_Vector2{X: childPosition.X, Y: childPosition.Y},
						NextChildOffset: Clay_Vector2{
							X: float32(childElement.LayoutConfig.Padding.Left),
							Y: float32(childElement.LayoutConfig.Padding.Top),
						},
					})
					Clay__Array_Set(&currentContext.TreeNodeVisited, newNodeIndex, false)

					// Update parent offsets
					if layoutConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT {
						currentElementTreeNode.NextChildOffset.X += childElement.Dimensions.Width + float32(layoutConfig.ChildGap)
					} else {
						currentElementTreeNode.NextChildOffset.Y += childElement.Dimensions.Height + float32(layoutConfig.ChildGap)
					}
				}
			}
		}

		if root.ClipElementId != 0 {
			Clay__AddRenderCommand(Clay_RenderCommand{
				Id:          Clay__HashNumber(rootElement.Id, uint32(int32(rootElement.ChildrenOrTextContent.Children.Length)+11)).Id,
				CommandType: CLAY_RENDER_COMMAND_TYPE_SCISSOR_END,
			})
		}
	}

}
