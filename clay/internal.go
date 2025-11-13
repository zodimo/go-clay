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

func Clay__Array_Allocate_Arena[T any](capacity int32, arena *Clay_Arena, options ...mem.MemArrayOption[T]) Clay__Array[T] {
	// var zero T
	// typeT := reflect.TypeOf(zero).String()
	// fmt.Println("typeT", typeT, "capacity", capacity)
	return mem.NewMemArray[T](capacity, append(options, mem.MemArrayWithArena[T](arena))...)
}

type Clay__LayoutElementChildren struct {
	ElementsPtr mem.SafeMemoryPointer[[]int32] // *elements, this translation is not working as expected
	Length      uint16
}

func (c *Clay__LayoutElementChildren) Elements() []int32 {
	ptr := mem.UintptrToPtr[[]int32](c.ElementsPtr.BaseAddress, c.ElementsPtr.InternalAddress)
	return *ptr
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

	for i := int32(0); i < openLayoutElement.ElementConfigs.Length(); i++ {
		config := Clay__Slice_Get(&openLayoutElement.ElementConfigs, i)
		if config == nil {
			panic("config is nil")
		}
		if config.Type == CLAY__ELEMENT_CONFIG_TYPE_CLIP {
			elementHasClipHorizontal = config.Config.ClipElementConfig.Horizontal
			elementHasClipVertical = config.Config.ClipElementConfig.Vertical
			Clay__Array_Pop(&currentContext.OpenClipElementStack)
			break
		} else if config.Type == CLAY__ELEMENT_CONFIG_TYPE_FLOATING {
			Clay__Array_Pop(&currentContext.OpenClipElementStack)
		}
	}

	leftRightPadding := float32(layoutConfig.Padding.Left + layoutConfig.Padding.Right)
	topBottomPadding := float32(layoutConfig.Padding.Top + layoutConfig.Padding.Bottom)

	// Attach children to the current open element

	//attach to the unallocated slice at the end of the array from length to capacity
	remainderSlicePointer := mem.MArray_GetIndexMemory(
		&currentContext.LayoutElementChildren,
		currentContext.LayoutElementChildren.Length())

	openLayoutElement.ChildrenOrTextContent.Children.ElementsPtr = remainderSlicePointer

	if layoutConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT {
		openLayoutElement.Dimensions.Width = leftRightPadding
		openLayoutElement.MinDimensions.Width = leftRightPadding
		for i := uint16(0); i < openLayoutElement.ChildrenOrTextContent.Children.Length; i++ {
			childIndex := Clay__Array_GetValue(&currentContext.LayoutElementChildrenBuffer, currentContext.LayoutElementChildrenBuffer.Length()-int32(openLayoutElement.ChildrenOrTextContent.Children.Length)+int32(i))
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
			childIndex := Clay__Array_GetValue(&currentContext.LayoutElementChildrenBuffer, currentContext.LayoutElementChildrenBuffer.Length()-int32(openLayoutElement.ChildrenOrTextContent.Children.Length)+int32(i))
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

	Clay__Array_Shrink(&currentContext.LayoutElementChildrenBuffer, int32(openLayoutElement.ChildrenOrTextContent.Children.Length))

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
	closingElementIndex := Clay__Array_RemoveSwapback(&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length()-1)

	// Get the currently open parent
	openLayoutElement = Clay__GetOpenLayoutElement()

	if currentContext.OpenLayoutElementStack.Length() > 1 {
		if elementIsFloating {
			openLayoutElement.FloatingChildrenCount++
			return
		}
		openLayoutElement.ChildrenOrTextContent.Children.Length++
		Clay__Array_Add(&currentContext.LayoutElementChildrenBuffer, closingElementIndex)
	}

}

func Clay__ElementHasConfig(layoutElement *Clay_LayoutElement, configType Clay__ElementConfigType) bool {
	for i := int32(0); i < layoutElement.ElementConfigs.Length(); i++ {
		if Clay__Slice_Get(&layoutElement.ElementConfigs, i).Type == configType {
			return true
		}
	}
	return false
}

func Clay__UpdateAspectRatioBox(layoutElement *Clay_LayoutElement) {
	for j := int32(0); j < layoutElement.ElementConfigs.Length(); j++ {
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
	panic("Clay__RenderDebugView not implemented")
}

// func Clay__RenderDebugView_() {
// 	currentContext := Clay_GetCurrentContext()
// 	closeButtonId := Clay__HashString(CLAY_STRING("Clay__DebugViewTopHeaderCloseButtonOuter"), 0)
// 	if currentContext.PointerInfo.State == CLAY_POINTER_DATA_PRESSED_THIS_FRAME {
// 		for i := int32(0); i < currentContext.PointerOverIds.Length(); i++ {
// 			elementId := Clay__Array_Get(&currentContext.PointerOverIds, i)
// 			if elementId.Id == closeButtonId.Id {
// 				currentContext.DebugModeEnabled = false
// 				return
// 			}
// 		}
// 	}

// 	initialRootsLength := currentContext.LayoutElementTreeRoots.Length
// 	initialElementsLength := currentContext.LayoutElements.Length
// 	infoTextConfig := Clay_TextElementConfig{
// 		TextColor: CLAY__DEBUGVIEW_COLOR_4,
// 		FontSize:  16,
// 		WrapMode:  CLAY_TEXT_WRAP_NONE,
// 	}
// 	infoTitleConfig := Clay_TextElementConfig{
// 		TextColor: CLAY__DEBUGVIEW_COLOR_3,
// 		FontSize:  16,
// 		WrapMode:  CLAY_TEXT_WRAP_NONE,
// 	}
// 	scrollId := Clay__HashString(CLAY_STRING("Clay__DebugViewOuterScrollPane"), 0)
// 	scrollYOffset := float32(0.0)

// 	pointerInDebugView := currentContext.PointerInfo.Position.Y < currentContext.LayoutDimensions.Height-300
// 	for i := int32(0); i < currentContext.ScrollContainerDatas.Length(); i++ {
// 		scrollContainerData := Clay__Array_Get(&currentContext.ScrollContainerDatas, i)
// 		if scrollContainerData.ElementId == scrollId.Id {
// 			if !currentContext.ExternalScrollHandlingEnabled {
// 				scrollYOffset = scrollContainerData.ScrollPosition.Y
// 			} else {
// 				pointerInDebugView = currentContext.PointerInfo.Position.Y+scrollContainerData.ScrollPosition.Y < currentContext.LayoutDimensions.Height-300
// 			}
// 			break
// 		}
// 	}
// 	var highlightedRow int32
// 	if pointerInDebugView {
// 		highlightedRow = int32((currentContext.PointerInfo.Position.Y-scrollYOffset)/float32(CLAY__DEBUGVIEW_ROW_HEIGHT)) - 1
// 	} else {
// 		highlightedRow = -1
// 	}

// 	layoutData := Clay__RenderDebugLayoutData{}

// 	CLAY(CLAY_ID("Clay__DebugView"), Clay_ElementDeclaration{
// 		Layout: Clay_LayoutConfig{
// 			Sizing: Clay_Sizing{
// 				Width:  CLAY_SIZING_FIXED(float32(Clay__debugViewWidth)),
// 				Height: CLAY_SIZING_FIXED(currentContext.LayoutDimensions.Height),
// 			},
// 			LayoutDirection: CLAY_TOP_TO_BOTTOM,
// 		},
// 		Floating: Clay_FloatingElementConfig{
// 			ZIndex: 32765,
// 			AttachPoints: Clay_FloatingAttachPoints{
// 				Element: CLAY_ATTACH_POINT_LEFT_CENTER,
// 				Parent:  CLAY_ATTACH_POINT_RIGHT_CENTER},
// 			AttachTo: CLAY_ATTACH_TO_ROOT,
// 			ClipTo:   CLAY_CLIP_TO_ATTACHED_PARENT,
// 		},
// 		Border: Clay_BorderElementConfig{
// 			Color: CLAY__DEBUGVIEW_COLOR_3,
// 			Width: Clay_BorderWidth{Bottom: 1},
// 		},
// 	})
// 	//     {
// 	// 	   CLAY_AUTO_ID({
// 	// 		.layout : {
// 	// 			.sizing = {
// 	// 				CLAY_SIZING_GROW(0),
// 	// 				CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT),
// 	// 				},
// 	// 				.padding = {
// 	// 					CLAY__DEBUGVIEW_OUTER_PADDING,
// 	// 					CLAY__DEBUGVIEW_OUTER_PADDING,
// 	// 					0,
// 	// 					0,
// 	// 					 },
// 	// 					 .childAlignment = {
// 	// 						.y = CLAY_ALIGN_Y_CENTER,
// 	// 						},
// 	// 						 },
// 	// 						 .backgroundColor = CLAY__DEBUGVIEW_COLOR_2,
// 	// 						  },
// 	// 						  ) {
// 	// 		   CLAY_TEXT(CLAY_STRING("Clay Debug Tools"), infoTextConfig);
// 	// 		   CLAY_AUTO_ID({
// 	// 			.layout = {
// 	// 				.sizing = {
// 	// 					.width = CLAY_SIZING_GROW(0),
// 	// 					 },
// 	// 					  },
// 	// 					   },
// 	// 					   ) {}
// 	// 		   // Close button
// 	// 		   CLAY_AUTO_ID({
// 	// 			   .layout = { .sizing = {CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10)}, .childAlignment = {CLAY_ALIGN_X_CENTER, CLAY_ALIGN_Y_CENTER} },
// 	// 			   .backgroundColor = {217,91,67,80},
// 	// 			   .cornerRadius = CLAY_CORNER_RADIUS(4),
// 	// 			   .border = { .color = { 217,91,67,255 }, .width = { 1, 1, 1, 1, 0 } },
// 	// 		   }) {
// 	// 			   Clay_OnHover(HandleDebugViewCloseButtonInteraction, 0);
// 	// 			   CLAY_TEXT(CLAY_STRING("x"), CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16 }));
// 	// 		   }
// 	// 	   }
// 	// 	   CLAY_AUTO_ID({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(1)} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_3 } ) {}
// 	// 	   CLAY(scrollId, { .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .clip = { .horizontal = true, .vertical = true, .childOffset = Clay_GetScrollOffset() } }) {
// 	// 		   CLAY_AUTO_ID({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)}, .layoutDirection = CLAY_TOP_TO_BOTTOM }, .backgroundColor = ((initialElementsLength + initialRootsLength) & 1) == 0 ? CLAY__DEBUGVIEW_COLOR_2 : CLAY__DEBUGVIEW_COLOR_1 }) {
// 	// 			   Clay_ElementId panelContentsId = Clay__HashString(CLAY_STRING("Clay__DebugViewPaneOuter"), 0);
// 	// 			   // Element list
// 	// 			   CLAY(panelContentsId, { .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)} }, .floating = { .zIndex = 32766, .pointerCaptureMode = CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH, .attachTo = CLAY_ATTACH_TO_PARENT, .clipTo = CLAY_CLIP_TO_ATTACHED_PARENT } }) {
// 	// 				   CLAY_AUTO_ID({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_GROW(0)}, .padding = { CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
// 	// 					   layoutData = Clay__RenderDebugLayoutElementsList((int32_t)initialRootsLength, highlightedRow);
// 	// 				   }
// 	// 			   }
// 	// 			   float contentWidth = Clay__GetHashMapItem(panelContentsId.id)->layoutElement->dimensions.width;
// 	// 			   CLAY_AUTO_ID({ .layout = { .sizing = {.width = CLAY_SIZING_FIXED(contentWidth) }, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {}
// 	// 			   for (int32_t i = 0; i < layoutData.rowCount; i++) {
// 	// 				   Clay_Color rowColor = (i & 1) == 0 ? CLAY__DEBUGVIEW_COLOR_2 : CLAY__DEBUGVIEW_COLOR_1;
// 	// 				   if (i == layoutData.selectedElementRowIndex) {
// 	// 					   rowColor = CLAY__DEBUGVIEW_COLOR_SELECTED_ROW;
// 	// 				   }
// 	// 				   if (i == highlightedRow) {
// 	// 					   rowColor.r *= 1.25f;
// 	// 					   rowColor.g *= 1.25f;
// 	// 					   rowColor.b *= 1.25f;
// 	// 				   }
// 	// 				   CLAY_AUTO_ID({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .layoutDirection = CLAY_TOP_TO_BOTTOM }, .backgroundColor = rowColor } ) {}
// 	// 			   }
// 	// 		   }
// 	// 	   }
// 	// 	   CLAY_AUTO_ID({ .layout = { .sizing = {.width = CLAY_SIZING_GROW(0), .height = CLAY_SIZING_FIXED(1)} }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_3 }) {}
// 	// 	   if (context->debugSelectedElementId != 0) {
// 	// 		   Clay_LayoutElementHashMapItem *selectedItem = Clay__GetHashMapItem(context->debugSelectedElementId);
// 	// 		   CLAY_AUTO_ID({
// 	// 			   .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(300)}, .layoutDirection = CLAY_TOP_TO_BOTTOM },
// 	// 			   .backgroundColor = CLAY__DEBUGVIEW_COLOR_2 ,
// 	// 			   .clip = { .vertical = true, .childOffset = Clay_GetScrollOffset() },
// 	// 			   .border = { .color = CLAY__DEBUGVIEW_COLOR_3, .width = { .betweenChildren = 1 } }
// 	// 		   }) {
// 	// 			   CLAY_AUTO_ID({ .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT + 8)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
// 	// 				   CLAY_TEXT(CLAY_STRING("Layout Config"), infoTextConfig);
// 	// 				   CLAY_AUTO_ID({ .layout = { .sizing = { .width = CLAY_SIZING_GROW(0) } } }) {}
// 	// 				   if (selectedItem->elementId.stringId.Length()!= 0) {
// 	// 					   CLAY_TEXT(selectedItem->elementId.stringId, infoTitleConfig);
// 	// 					   if (selectedItem->elementId.offset != 0) {
// 	// 						   CLAY_TEXT(CLAY_STRING(" ("), infoTitleConfig);
// 	// 						   CLAY_TEXT(Clay__IntToString(selectedItem->elementId.offset), infoTitleConfig);
// 	// 						   CLAY_TEXT(CLAY_STRING(")"), infoTitleConfig);
// 	// 					   }
// 	// 				   }
// 	// 			   }
// 	// 			   Clay_Padding attributeConfigPadding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 8, 8};
// 	// 			   // Clay_LayoutConfig debug info
// 	// 			   CLAY_AUTO_ID({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
// 	// 				   // .boundingBox
// 	// 				   CLAY_TEXT(CLAY_STRING("Bounding Box"), infoTitleConfig);
// 	// 				   CLAY_AUTO_ID({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
// 	// 					   CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig);
// 	// 					   CLAY_TEXT(Clay__IntToString(selectedItem->boundingBox.x), infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig);
// 	// 					   CLAY_TEXT(Clay__IntToString(selectedItem->boundingBox.y), infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(", width: "), infoTextConfig);
// 	// 					   CLAY_TEXT(Clay__IntToString(selectedItem->boundingBox.width), infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(", height: "), infoTextConfig);
// 	// 					   CLAY_TEXT(Clay__IntToString(selectedItem->boundingBox.height), infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
// 	// 				   }
// 	// 				   // .layoutDirection
// 	// 				   CLAY_TEXT(CLAY_STRING("Layout Direction"), infoTitleConfig);
// 	// 				   Clay_LayoutConfig *layoutConfig = selectedItem->layoutElement->layoutConfig;
// 	// 				   CLAY_TEXT(layoutConfig->layoutDirection == CLAY_TOP_TO_BOTTOM ? CLAY_STRING("TOP_TO_BOTTOM") : CLAY_STRING("LEFT_TO_RIGHT"), infoTextConfig);
// 	// 				   // .sizing
// 	// 				   CLAY_TEXT(CLAY_STRING("Sizing"), infoTitleConfig);
// 	// 				   CLAY_AUTO_ID({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
// 	// 					   CLAY_TEXT(CLAY_STRING("width: "), infoTextConfig);
// 	// 					   Clay__RenderDebugLayoutSizing(layoutConfig->sizing.width, infoTextConfig);
// 	// 				   }
// 	// 				   CLAY_AUTO_ID({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
// 	// 					   CLAY_TEXT(CLAY_STRING("height: "), infoTextConfig);
// 	// 					   Clay__RenderDebugLayoutSizing(layoutConfig->sizing.height, infoTextConfig);
// 	// 				   }
// 	// 				   // .padding
// 	// 				   CLAY_TEXT(CLAY_STRING("Padding"), infoTitleConfig);
// 	// 				   CLAY(CLAY_ID("Clay__DebugViewElementInfoPadding"), { }) {
// 	// 					   CLAY_TEXT(CLAY_STRING("{ left: "), infoTextConfig);
// 	// 					   CLAY_TEXT(Clay__IntToString(layoutConfig->padding.left), infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(", right: "), infoTextConfig);
// 	// 					   CLAY_TEXT(Clay__IntToString(layoutConfig->padding.right), infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(", top: "), infoTextConfig);
// 	// 					   CLAY_TEXT(Clay__IntToString(layoutConfig->padding.top), infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(", bottom: "), infoTextConfig);
// 	// 					   CLAY_TEXT(Clay__IntToString(layoutConfig->padding.bottom), infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
// 	// 				   }
// 	// 				   // .childGap
// 	// 				   CLAY_TEXT(CLAY_STRING("Child Gap"), infoTitleConfig);
// 	// 				   CLAY_TEXT(Clay__IntToString(layoutConfig->childGap), infoTextConfig);
// 	// 				   // .childAlignment
// 	// 				   CLAY_TEXT(CLAY_STRING("Child Alignment"), infoTitleConfig);
// 	// 				   CLAY_AUTO_ID({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
// 	// 					   CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig);
// 	// 					   Clay_String alignX = CLAY_STRING("LEFT");
// 	// 					   if (layoutConfig->childAlignment.x == CLAY_ALIGN_X_CENTER) {
// 	// 						   alignX = CLAY_STRING("CENTER");
// 	// 					   } else if (layoutConfig->childAlignment.x == CLAY_ALIGN_X_RIGHT) {
// 	// 						   alignX = CLAY_STRING("RIGHT");
// 	// 					   }
// 	// 					   CLAY_TEXT(alignX, infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig);
// 	// 					   Clay_String alignY = CLAY_STRING("TOP");
// 	// 					   if (layoutConfig->childAlignment.y == CLAY_ALIGN_Y_CENTER) {
// 	// 						   alignY = CLAY_STRING("CENTER");
// 	// 					   } else if (layoutConfig->childAlignment.y == CLAY_ALIGN_Y_BOTTOM) {
// 	// 						   alignY = CLAY_STRING("BOTTOM");
// 	// 					   }
// 	// 					   CLAY_TEXT(alignY, infoTextConfig);
// 	// 					   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
// 	// 				   }
// 	// 			   }
// 	// 			   for (int32_t elementConfigIndex = 0; elementConfigIndex < selectedItem->layoutElement->elementConfigs.Length(); ++elementConfigIndex) {
// 	// 				   Clay_ElementConfig *elementConfig = Clay__ElementConfigArraySlice_Get(&selectedItem->layoutElement->elementConfigs, elementConfigIndex);
// 	// 				   Clay__RenderDebugViewElementConfigHeader(selectedItem->elementId.stringId, elementConfig->type);
// 	// 				   switch (elementConfig->type) {
// 	// 					   case CLAY__ELEMENT_CONFIG_TYPE_SHARED: {
// 	// 						   Clay_SharedElementConfig *sharedConfig = elementConfig->config.sharedElementConfig;
// 	// 						   CLAY_AUTO_ID({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM }}) {
// 	// 							   // .backgroundColor
// 	// 							   CLAY_TEXT(CLAY_STRING("Background Color"), infoTitleConfig);
// 	// 							   Clay__RenderDebugViewColor(sharedConfig->backgroundColor, infoTextConfig);
// 	// 							   // .cornerRadius
// 	// 							   CLAY_TEXT(CLAY_STRING("Corner Radius"), infoTitleConfig);
// 	// 							   Clay__RenderDebugViewCornerRadius(sharedConfig->cornerRadius, infoTextConfig);
// 	// 						   }
// 	// 						   break;
// 	// 					   }
// 	// 					   case CLAY__ELEMENT_CONFIG_TYPE_TEXT: {
// 	// 						   Clay_TextElementConfig *textConfig = elementConfig->config.textElementConfig;
// 	// 						   CLAY_AUTO_ID({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
// 	// 							   // .fontSize
// 	// 							   CLAY_TEXT(CLAY_STRING("Font Size"), infoTitleConfig);
// 	// 							   CLAY_TEXT(Clay__IntToString(textConfig->fontSize), infoTextConfig);
// 	// 							   // .fontId
// 	// 							   CLAY_TEXT(CLAY_STRING("Font ID"), infoTitleConfig);
// 	// 							   CLAY_TEXT(Clay__IntToString(textConfig->fontId), infoTextConfig);
// 	// 							   // .lineHeight
// 	// 							   CLAY_TEXT(CLAY_STRING("Line Height"), infoTitleConfig);
// 	// 							   CLAY_TEXT(textConfig->lineHeight == 0 ? CLAY_STRING("auto") : Clay__IntToString(textConfig->lineHeight), infoTextConfig);
// 	// 							   // .letterSpacing
// 	// 							   CLAY_TEXT(CLAY_STRING("Letter Spacing"), infoTitleConfig);
// 	// 							   CLAY_TEXT(Clay__IntToString(textConfig->letterSpacing), infoTextConfig);
// 	// 							   // .wrapMode
// 	// 							   CLAY_TEXT(CLAY_STRING("Wrap Mode"), infoTitleConfig);
// 	// 							   Clay_String wrapMode = CLAY_STRING("WORDS");
// 	// 							   if (textConfig->wrapMode == CLAY_TEXT_WRAP_NONE) {
// 	// 								   wrapMode = CLAY_STRING("NONE");
// 	// 							   } else if (textConfig->wrapMode == CLAY_TEXT_WRAP_NEWLINES) {
// 	// 								   wrapMode = CLAY_STRING("NEWLINES");
// 	// 							   }
// 	// 							   CLAY_TEXT(wrapMode, infoTextConfig);
// 	// 							   // .textAlignment
// 	// 							   CLAY_TEXT(CLAY_STRING("Text Alignment"), infoTitleConfig);
// 	// 							   Clay_String textAlignment = CLAY_STRING("LEFT");
// 	// 							   if (textConfig->textAlignment == CLAY_TEXT_ALIGN_CENTER) {
// 	// 								   textAlignment = CLAY_STRING("CENTER");
// 	// 							   } else if (textConfig->textAlignment == CLAY_TEXT_ALIGN_RIGHT) {
// 	// 								   textAlignment = CLAY_STRING("RIGHT");
// 	// 							   }
// 	// 							   CLAY_TEXT(textAlignment, infoTextConfig);
// 	// 							   // .textColor
// 	// 							   CLAY_TEXT(CLAY_STRING("Text Color"), infoTitleConfig);
// 	// 							   Clay__RenderDebugViewColor(textConfig->textColor, infoTextConfig);
// 	// 						   }
// 	// 						   break;
// 	// 					   }
// 	// 					   case CLAY__ELEMENT_CONFIG_TYPE_ASPECT: {
// 	// 						   Clay_AspectRatioElementConfig *aspectRatioConfig = elementConfig->config.aspectRatioElementConfig;
// 	// 						   CLAY(CLAY_ID("Clay__DebugViewElementInfoAspectRatioBody"), { .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
// 	// 							   CLAY_TEXT(CLAY_STRING("Aspect Ratio"), infoTitleConfig);
// 	// 							   // Aspect Ratio
// 	// 							   CLAY(CLAY_ID("Clay__DebugViewElementInfoAspectRatio"), { }) {
// 	// 								   CLAY_TEXT(Clay__IntToString(aspectRatioConfig->aspectRatio), infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING("."), infoTextConfig);
// 	// 								   float frac = aspectRatioConfig->aspectRatio - (int)(aspectRatioConfig->aspectRatio);
// 	// 								   frac *= 100;
// 	// 								   if ((int)frac < 10) {
// 	// 									   CLAY_TEXT(CLAY_STRING("0"), infoTextConfig);
// 	// 								   }
// 	// 								   CLAY_TEXT(Clay__IntToString(frac), infoTextConfig);
// 	// 							   }
// 	// 						   }
// 	// 						   break;
// 	// 					   }
// 	// 					   case CLAY__ELEMENT_CONFIG_TYPE_IMAGE: {
// 	// 						   Clay_ImageElementConfig *imageConfig = elementConfig->config.imageElementConfig;
// 	// 						   Clay_AspectRatioElementConfig aspectConfig = { 1 };
// 	// 						   if (Clay__ElementHasConfig(selectedItem->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT)) {
// 	// 							   aspectConfig = *Clay__FindElementConfigWithType(selectedItem->layoutElement, CLAY__ELEMENT_CONFIG_TYPE_ASPECT).aspectRatioElementConfig;
// 	// 						   }
// 	// 						   CLAY(CLAY_ID("Clay__DebugViewElementInfoImageBody"), { .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
// 	// 							   // Image Preview
// 	// 							   CLAY_TEXT(CLAY_STRING("Preview"), infoTitleConfig);
// 	// 							   CLAY_AUTO_ID({ .layout = { .sizing = { .width = CLAY_SIZING_GROW(64, 128), .height = CLAY_SIZING_GROW(64, 128) }}, .aspectRatio = aspectConfig, .image = *imageConfig }) {}
// 	// 						   }
// 	// 						   break;
// 	// 					   }
// 	// 					   case CLAY__ELEMENT_CONFIG_TYPE_CLIP: {
// 	// 						   Clay_ClipElementConfig *clipConfig = elementConfig->config.clipElementConfig;
// 	// 						   CLAY_AUTO_ID({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
// 	// 							   // .vertical
// 	// 							   CLAY_TEXT(CLAY_STRING("Vertical"), infoTitleConfig);
// 	// 							   CLAY_TEXT(clipConfig->vertical ? CLAY_STRING("true") : CLAY_STRING("false") , infoTextConfig);
// 	// 							   // .horizontal
// 	// 							   CLAY_TEXT(CLAY_STRING("Horizontal"), infoTitleConfig);
// 	// 							   CLAY_TEXT(clipConfig->horizontal ? CLAY_STRING("true") : CLAY_STRING("false") , infoTextConfig);
// 	// 						   }
// 	// 						   break;
// 	// 					   }
// 	// 					   case CLAY__ELEMENT_CONFIG_TYPE_FLOATING: {
// 	// 						   Clay_FloatingElementConfig *floatingConfig = elementConfig->config.floatingElementConfig;
// 	// 						   CLAY_AUTO_ID({ .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
// 	// 							   // .offset
// 	// 							   CLAY_TEXT(CLAY_STRING("Offset"), infoTitleConfig);
// 	// 							   CLAY_AUTO_ID({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
// 	// 								   CLAY_TEXT(CLAY_STRING("{ x: "), infoTextConfig);
// 	// 								   CLAY_TEXT(Clay__IntToString(floatingConfig->offset.x), infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING(", y: "), infoTextConfig);
// 	// 								   CLAY_TEXT(Clay__IntToString(floatingConfig->offset.y), infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
// 	// 							   }
// 	// 							   // .expand
// 	// 							   CLAY_TEXT(CLAY_STRING("Expand"), infoTitleConfig);
// 	// 							   CLAY_AUTO_ID({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
// 	// 								   CLAY_TEXT(CLAY_STRING("{ width: "), infoTextConfig);
// 	// 								   CLAY_TEXT(Clay__IntToString(floatingConfig->expand.width), infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING(", height: "), infoTextConfig);
// 	// 								   CLAY_TEXT(Clay__IntToString(floatingConfig->expand.height), infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
// 	// 							   }
// 	// 							   // .zIndex
// 	// 							   CLAY_TEXT(CLAY_STRING("z-index"), infoTitleConfig);
// 	// 							   CLAY_TEXT(Clay__IntToString(floatingConfig->zIndex), infoTextConfig);
// 	// 							   // .parentId
// 	// 							   CLAY_TEXT(CLAY_STRING("Parent"), infoTitleConfig);
// 	// 							   Clay_LayoutElementHashMapItem *hashItem = Clay__GetHashMapItem(floatingConfig->parentId);
// 	// 							   CLAY_TEXT(hashItem->elementId.stringId, infoTextConfig);
// 	// 							   // .attachPoints
// 	// 							   CLAY_TEXT(CLAY_STRING("Attach Points"), infoTitleConfig);
// 	// 							   CLAY_AUTO_ID({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
// 	// 								   CLAY_TEXT(CLAY_STRING("{ element: "), infoTextConfig);
// 	// 								   Clay_String attachPointElement = CLAY_STRING("LEFT_TOP");
// 	// 								   if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_LEFT_CENTER) {
// 	// 									   attachPointElement = CLAY_STRING("LEFT_CENTER");
// 	// 								   } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_LEFT_BOTTOM) {
// 	// 									   attachPointElement = CLAY_STRING("LEFT_BOTTOM");
// 	// 								   } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_CENTER_TOP) {
// 	// 									   attachPointElement = CLAY_STRING("CENTER_TOP");
// 	// 								   } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_CENTER_CENTER) {
// 	// 									   attachPointElement = CLAY_STRING("CENTER_CENTER");
// 	// 								   } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_CENTER_BOTTOM) {
// 	// 									   attachPointElement = CLAY_STRING("CENTER_BOTTOM");
// 	// 								   } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_RIGHT_TOP) {
// 	// 									   attachPointElement = CLAY_STRING("RIGHT_TOP");
// 	// 								   } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_RIGHT_CENTER) {
// 	// 									   attachPointElement = CLAY_STRING("RIGHT_CENTER");
// 	// 								   } else if (floatingConfig->attachPoints.element == CLAY_ATTACH_POINT_RIGHT_BOTTOM) {
// 	// 									   attachPointElement = CLAY_STRING("RIGHT_BOTTOM");
// 	// 								   }
// 	// 								   CLAY_TEXT(attachPointElement, infoTextConfig);
// 	// 								   Clay_String attachPointParent = CLAY_STRING("LEFT_TOP");
// 	// 								   if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_LEFT_CENTER) {
// 	// 									   attachPointParent = CLAY_STRING("LEFT_CENTER");
// 	// 								   } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_LEFT_BOTTOM) {
// 	// 									   attachPointParent = CLAY_STRING("LEFT_BOTTOM");
// 	// 								   } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_CENTER_TOP) {
// 	// 									   attachPointParent = CLAY_STRING("CENTER_TOP");
// 	// 								   } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_CENTER_CENTER) {
// 	// 									   attachPointParent = CLAY_STRING("CENTER_CENTER");
// 	// 								   } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_CENTER_BOTTOM) {
// 	// 									   attachPointParent = CLAY_STRING("CENTER_BOTTOM");
// 	// 								   } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_RIGHT_TOP) {
// 	// 									   attachPointParent = CLAY_STRING("RIGHT_TOP");
// 	// 								   } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_RIGHT_CENTER) {
// 	// 									   attachPointParent = CLAY_STRING("RIGHT_CENTER");
// 	// 								   } else if (floatingConfig->attachPoints.parent == CLAY_ATTACH_POINT_RIGHT_BOTTOM) {
// 	// 									   attachPointParent = CLAY_STRING("RIGHT_BOTTOM");
// 	// 								   }
// 	// 								   CLAY_TEXT(CLAY_STRING(", parent: "), infoTextConfig);
// 	// 								   CLAY_TEXT(attachPointParent, infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
// 	// 							   }
// 	// 							   // .pointerCaptureMode
// 	// 							   CLAY_TEXT(CLAY_STRING("Pointer Capture Mode"), infoTitleConfig);
// 	// 							   Clay_String pointerCaptureMode = CLAY_STRING("NONE");
// 	// 							   if (floatingConfig->pointerCaptureMode == CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH) {
// 	// 								   pointerCaptureMode = CLAY_STRING("PASSTHROUGH");
// 	// 							   }
// 	// 							   CLAY_TEXT(pointerCaptureMode, infoTextConfig);
// 	// 							   // .attachTo
// 	// 							   CLAY_TEXT(CLAY_STRING("Attach To"), infoTitleConfig);
// 	// 							   Clay_String attachTo = CLAY_STRING("NONE");
// 	// 							   if (floatingConfig->attachTo == CLAY_ATTACH_TO_PARENT) {
// 	// 								   attachTo = CLAY_STRING("PARENT");
// 	// 							   } else if (floatingConfig->attachTo == CLAY_ATTACH_TO_ELEMENT_WITH_ID) {
// 	// 								   attachTo = CLAY_STRING("ELEMENT_WITH_ID");
// 	// 							   } else if (floatingConfig->attachTo == CLAY_ATTACH_TO_ROOT) {
// 	// 								   attachTo = CLAY_STRING("ROOT");
// 	// 							   }
// 	// 							   CLAY_TEXT(attachTo, infoTextConfig);
// 	// 							   // .clipTo
// 	// 							   CLAY_TEXT(CLAY_STRING("Clip To"), infoTitleConfig);
// 	// 							   Clay_String clipTo = CLAY_STRING("ATTACHED_PARENT");
// 	// 							   if (floatingConfig->clipTo == CLAY_CLIP_TO_NONE) {
// 	// 								   clipTo = CLAY_STRING("NONE");
// 	// 							   }
// 	// 							   CLAY_TEXT(clipTo, infoTextConfig);
// 	// 						   }
// 	// 						   break;
// 	// 					   }
// 	// 					   case CLAY__ELEMENT_CONFIG_TYPE_BORDER: {
// 	// 						   Clay_BorderElementConfig *borderConfig = elementConfig->config.borderElementConfig;
// 	// 						   CLAY(CLAY_ID("Clay__DebugViewElementInfoBorderBody"), { .layout = { .padding = attributeConfigPadding, .childGap = 8, .layoutDirection = CLAY_TOP_TO_BOTTOM } }) {
// 	// 							   CLAY_TEXT(CLAY_STRING("Border Widths"), infoTitleConfig);
// 	// 							   CLAY_AUTO_ID({ .layout = { .layoutDirection = CLAY_LEFT_TO_RIGHT } }) {
// 	// 								   CLAY_TEXT(CLAY_STRING("{ left: "), infoTextConfig);
// 	// 								   CLAY_TEXT(Clay__IntToString(borderConfig->width.left), infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING(", right: "), infoTextConfig);
// 	// 								   CLAY_TEXT(Clay__IntToString(borderConfig->width.right), infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING(", top: "), infoTextConfig);
// 	// 								   CLAY_TEXT(Clay__IntToString(borderConfig->width.top), infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING(", bottom: "), infoTextConfig);
// 	// 								   CLAY_TEXT(Clay__IntToString(borderConfig->width.bottom), infoTextConfig);
// 	// 								   CLAY_TEXT(CLAY_STRING(" }"), infoTextConfig);
// 	// 							   }
// 	// 							   // .textColor
// 	// 							   CLAY_TEXT(CLAY_STRING("Border Color"), infoTitleConfig);
// 	// 							   Clay__RenderDebugViewColor(borderConfig->color, infoTextConfig);
// 	// 						   }
// 	// 						   break;
// 	// 					   }
// 	// 					   case CLAY__ELEMENT_CONFIG_TYPE_CUSTOM:
// 	// 					   default: break;
// 	// 				   }
// 	// 			   }
// 	// 		   }
// 	// 	   } else {
// 	// 		   CLAY(CLAY_ID("Clay__DebugViewWarningsScrollPane"), { .layout = { .sizing = {CLAY_SIZING_GROW(0), CLAY_SIZING_FIXED(300)}, .childGap = 6, .layoutDirection = CLAY_TOP_TO_BOTTOM }, .backgroundColor = CLAY__DEBUGVIEW_COLOR_2, .clip = { .horizontal = true, .vertical = true, .childOffset = Clay_GetScrollOffset() } }) {
// 	// 			   Clay_TextElementConfig *warningConfig = CLAY_TEXT_CONFIG({ .textColor = CLAY__DEBUGVIEW_COLOR_4, .fontSize = 16, .wrapMode = CLAY_TEXT_WRAP_NONE });
// 	// 			   CLAY(CLAY_ID("Clay__DebugViewWarningItemHeader"), { .layout = { .sizing = {.height = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childGap = 8, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
// 	// 				   CLAY_TEXT(CLAY_STRING("Warnings"), warningConfig);
// 	// 			   }
// 	// 			   CLAY(CLAY_ID("Clay__DebugViewWarningsTopBorder"), { .layout = { .sizing = { .width = CLAY_SIZING_GROW(0), .height = CLAY_SIZING_FIXED(1)} }, .backgroundColor = {200, 200, 200, 255} }) {}
// 	// 			   int32_t previousWarningsLength = context->warnings.Length();
// 	// 			   for (int32_t i = 0; i < previousWarningsLength; i++) {
// 	// 				   Clay__Warning warning = context->warnings.internalArray[i];
// 	// 				   CLAY(CLAY_IDI("Clay__DebugViewWarningItem", i), { .layout = { .sizing = {.height = CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT)}, .padding = {CLAY__DEBUGVIEW_OUTER_PADDING, CLAY__DEBUGVIEW_OUTER_PADDING, 0, 0 }, .childGap = 8, .childAlignment = {.y = CLAY_ALIGN_Y_CENTER} } }) {
// 	// 					   CLAY_TEXT(warning.baseMessage, warningConfig);
// 	// 					   if (warning.dynamicMessage.Length()> 0) {
// 	// 						   CLAY_TEXT(warning.dynamicMessage, warningConfig);
// 	// 					   }
// 	// 				   }
// 	// 			   }
// 	// 		   }
// 	// 	   }
// 	//    }

// }
func Clay__FindElementConfigWithType(element *Clay_LayoutElement, configType Clay__ElementConfigType) Clay_ElementConfigUnion {
	for i := int32(0); i < element.ElementConfigs.Length(); i++ {
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
				fmt.Printf("Clay__SizeContainersAlongAxis childOffset: %d\n", childOffset)
				fmt.Printf("Clay__SizeContainersAlongAxis parent.ChildrenOrTextContent.Children.Length: %d\n", parent.ChildrenOrTextContent.Children.Length)
				childElementIndex := mem.NewMemSliceWithData(parent.ChildrenOrTextContent.Children.Elements()).Get(childOffset)
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
				childElementIndex := mem.NewMemSliceWithData(parent.ChildrenOrTextContent.Children.Elements()).Get(childOffset)
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
				}
			}

		}

	}
}
func Clay__ElementIsOffscreen(boundingBox *Clay_BoundingBox) bool {
	context := Clay_GetCurrentContext()
	if context.DisableCulling {
		return false
	}

	return boundingBox.X > float32(context.LayoutDimensions.Width) ||
		boundingBox.Y > float32(context.LayoutDimensions.Height) ||
		boundingBox.X+boundingBox.Width < 0 ||
		boundingBox.Y+boundingBox.Height < 0
}

func Clay__CalculateFinalLayout() {
	currentContext := Clay_GetCurrentContext()
	// Calculate sizing along the X axis
	Clay__SizeContainersAlongAxis(true)

	// Wrap text
	for textElementIndex := int32(0); textElementIndex < currentContext.TextElementData.Length(); textElementIndex++ {
		textElementData := Clay__Array_Get(&currentContext.TextElementData, textElementIndex)
		wrappedLinesData := mem.MArray_GetSlice(&currentContext.WrappedTextLines, 0, currentContext.WrappedTextLines.Length())
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
			if currentContext.WrappedTextLines.Length() > currentContext.WrappedTextLines.Capacity-1 {
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
						currentElement.ChildrenOrTextContent.Children.Elements()[childIndex],
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
				childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements()[childIndex])
				childHeightWithPadding := CLAY__MAX(childElement.Dimensions.Height+float32(layoutConfig.Padding.Top)+float32(layoutConfig.Padding.Bottom), currentElement.Dimensions.Height)
				currentElement.Dimensions.Height = CLAY__MIN(CLAY__MAX(childHeightWithPadding, layoutConfig.Sizing.Height.Size.MinMax.Min), layoutConfig.Sizing.Height.Size.MinMax.Max)
			}
		} else if layoutConfig.LayoutDirection == CLAY_TOP_TO_BOTTOM {
			// Resizing along the layout axis
			contentHeight := float32(layoutConfig.Padding.Top + layoutConfig.Padding.Bottom)
			for childIndex := int32(0); childIndex < int32(currentElement.ChildrenOrTextContent.Children.Length); childIndex++ {
				childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements()[childIndex])
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

				sortedConfigIndexes := make([]int32, 20)
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
					case CLAY__ELEMENT_CONFIG_TYPE_FLOATING:
					case CLAY__ELEMENT_CONFIG_TYPE_SHARED:
					case CLAY__ELEMENT_CONFIG_TYPE_BORDER:
						{
							shouldRender = false
							break
						}
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
							childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements()[i])
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
							childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements()[i])
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
									childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements()[i])
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
									childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements()[i])
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
					childElement := Clay__Array_Get(&currentContext.LayoutElements, currentElement.ChildrenOrTextContent.Children.Elements()[i])
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

func Clay__AddRenderCommand(renderCommand Clay_RenderCommand) {
	currentContext := Clay_GetCurrentContext()
	if currentContext.RenderCommands.Length() < currentContext.RenderCommands.Capacity-1 {
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
func Clay__GenerateIdForAnonymousElement(openLayoutElement *Clay_LayoutElement) Clay_ElementId {
	currentContext := Clay_GetCurrentContext()
	parentElement := Clay__Array_Get(&currentContext.LayoutElements, Clay__Array_GetValue(&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length()-2))
	offset := uint32(parentElement.ChildrenOrTextContent.Children.Length + parentElement.FloatingChildrenCount)
	elementId := Clay__HashNumber(offset, parentElement.Id)
	openLayoutElement.Id = elementId.Id
	Clay__AddHashMapItem(elementId, openLayoutElement)
	Clay__Array_Add(&currentContext.LayoutElementIdStrings, elementId.StringId)
	return elementId
}

func Clay__OpenElement() {
	currentContext := Clay_GetCurrentContext()
	if currentContext.LayoutElements.Length() == currentContext.LayoutElements.Capacity-1 || currentContext.BooleanWarnings.MaxElementsExceeded {
		currentContext.BooleanWarnings.MaxElementsExceeded = true
		return
	}
	layoutElement := Clay_LayoutElement{}
	openLayoutElement := Clay__Array_Add(&currentContext.LayoutElements, layoutElement)
	Clay__Array_Add(&currentContext.OpenLayoutElementStack, currentContext.LayoutElements.Length()-1)
	Clay__GenerateIdForAnonymousElement(openLayoutElement)
	if currentContext.OpenClipElementStack.Length() > 0 {
		Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length()-1, Clay__Array_GetValue(&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length()-1))
	} else {
		Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length()-1, 0)
	}
}
func Clay__OpenElementWithId(elementId Clay_ElementId) {
	currentContext := Clay_GetCurrentContext()
	if currentContext.LayoutElements.Length() == currentContext.LayoutElements.Capacity-1 || currentContext.BooleanWarnings.MaxElementsExceeded {
		currentContext.BooleanWarnings.MaxElementsExceeded = true
		return
	}
	layoutElement := Clay_LayoutElement{}
	layoutElement.Id = elementId.Id
	openLayoutElement := Clay__Array_Add(&currentContext.LayoutElements, layoutElement)
	Clay__Array_Add(&currentContext.OpenLayoutElementStack, currentContext.LayoutElements.Length()-1) // add the index of the new layout element to the open layout element stack
	Clay__AddHashMapItem(elementId, openLayoutElement)
	Clay__Array_Add(&currentContext.LayoutElementIdStrings, elementId.StringId)
	if currentContext.OpenClipElementStack.Length() > 0 {
		if currentContext.LayoutElementClipElementIds.Length() == 0 {
			Clay__Array_Add(&currentContext.LayoutElementClipElementIds, Clay__Array_GetValue(&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length()-1))

		} else {
			Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length()-1, Clay__Array_GetValue(&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length()-1))
		}
	} else {
		if currentContext.LayoutElementClipElementIds.Length() == 0 {
			Clay__Array_Add(&currentContext.LayoutElementClipElementIds, 0)
		} else {
			if currentContext.LayoutElementClipElementIds.Length() <= currentContext.LayoutElements.Length()-1 {
				Clay__Array_Add(&currentContext.LayoutElementClipElementIds, 0)
			} else {
				Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length()-1, 0)
			}
		}
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
	Clay__Slice_Grow(&openLayoutElement.ElementConfigs, 1)
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
	Clay__ConfigureOpenElementPtr(&elementDeclaration)
}
func Clay__ConfigureOpenElementPtr(elementDeclaration *Clay_ElementDeclaration) {

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
	nextAvailableElementConfigIndex := currentContext.ElementConfigs.Length()
	fmt.Println("nextAvailableElementConfigIndex", nextAvailableElementConfigIndex)

	elementConfigsPointer := mem.MArray_GetIndexMemory(&currentContext.ElementConfigs, nextAvailableElementConfigIndex)
	openLayoutElement.ElementConfigs.BaseAddress = elementConfigsPointer.BaseAddress
	openLayoutElement.ElementConfigs.InternalAddress = elementConfigsPointer.InternalAddress

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
		Clay__Array_Add(&currentContext.AspectRatioElementIndexes, currentContext.LayoutElements.Length()-1)
	}

	if elementDeclaration.Floating.AttachTo != CLAY_ATTACH_TO_NONE {
		floatingConfig := elementDeclaration.Floating
		// This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here

		hierarchicalParent := Clay__Array_Get[Clay_LayoutElement](&currentContext.LayoutElements, Clay__Array_GetValue[int32](&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length()-2))
		if hierarchicalParent != nil {
			var clipElementId int32 = 0
			if elementDeclaration.Floating.AttachTo == CLAY_ATTACH_TO_PARENT {
				// Attach to the element's direct hierarchical parent
				floatingConfig.ParentId = hierarchicalParent.Id
				if currentContext.OpenClipElementStack.Length() > 0 {
					clipElementId = Clay__Array_GetValue(&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length()-1)
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
				currentElementIndex := Clay__Array_GetValue[int32](&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length()-1)
				Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentElementIndex, clipElementId)
				Clay__Array_Add(&currentContext.OpenClipElementStack, clipElementId)
				Clay__Array_Add(&currentContext.LayoutElementTreeRoots, Clay__LayoutElementTreeRoot{
					LayoutElementIndex: Clay__Array_GetValue[int32](&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length()-1),
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
		for i := int32(0); i < currentContext.ScrollContainerDatas.Length(); i++ {
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
	return Clay__Array_Get[Clay_LayoutElement](&currentContext.LayoutElements, Clay__Array_GetValue[int32](&currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length()-1))

	// Clay_LayoutElement* Clay__GetOpenLayoutElement(void) {
	//     Clay_Context* context = Clay_GetCurrentContext();
	//     return Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.Length()- 1));
	// }

}
func Clay__MeasureTextCached(text *Clay_String, textConfig *Clay_TextElementConfig) *Clay__MeasureTextCacheItem {
	panic("not implemented")
}

func Clay__OpenTextElement(text Clay_String, textConfig *Clay_TextElementConfig) {
	currentContext := Clay_GetCurrentContext()
	if currentContext.LayoutElements.Length() == currentContext.LayoutElements.Capacity-1 || currentContext.BooleanWarnings.MaxElementsExceeded {
		currentContext.BooleanWarnings.MaxElementsExceeded = true
		return
	}
	parentElement := Clay__GetOpenLayoutElement()

	layoutElement := Clay_LayoutElement{}

	textElement := Clay__Array_Add[Clay_LayoutElement](&currentContext.LayoutElements, layoutElement)

	if currentContext.OpenClipElementStack.Length() > 0 {
		Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length()-1, Clay__Array_GetValue[int32](&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length()-1))
	} else {
		Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length()-1, 0)
	}

	Clay__Array_Add(&currentContext.LayoutElementChildrenBuffer, currentContext.LayoutElements.Length()-1)

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
		ElementIndex:        currentContext.LayoutElements.Length() - 1,
	})

	// add config to element configs

	config := Clay__Array_Add(&currentContext.ElementConfigs, Clay_ElementConfig{
		Type:   CLAY__ELEMENT_CONFIG_TYPE_TEXT,
		Config: Clay_ElementConfigUnion{TextElementConfig: textConfig},
	})
	if config != nil {
		configIndex := currentContext.ElementConfigs.Length() - 1

		segmentView := mem.MArray_GetSlice(&currentContext.ElementConfigs, configIndex, configIndex+1)
		textElement.ElementConfigs = NewClay__Slice[Clay_ElementConfig](segmentView)
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
	context.LayoutElementsHashMapInternal = Clay__Array_Allocate_Arena[Clay_LayoutElementHashMapItem](maxElementCount, arena, mem.MemArrayWithZeroValuePtr[Clay_LayoutElementHashMapItem](&Clay_LayoutElementHashMapItem_DEFAULT))
	context.LayoutElementsHashMap = Clay__Array_Allocate_Arena[int32](maxElementCount, arena, mem.MemArrayWithIsHashmap[int32](), mem.MemArrayWithZeroValue[int32](-1))
	context.MeasureTextHashMapInternal = Clay__Array_Allocate_Arena[Clay__MeasureTextCacheItem](maxElementCount, arena)
	context.MeasureTextHashMapInternalFreeList = Clay__Array_Allocate_Arena[int32](maxElementCount, arena)
	context.MeasuredWordsFreeList = Clay__Array_Allocate_Arena[int32](maxMeasureTextCacheWordCount, arena)
	context.MeasureTextHashMap = Clay__Array_Allocate_Arena[int32](maxElementCount, arena, mem.MemArrayWithIsHashmap[int32](), mem.MemArrayWithZeroValue[int32](0))
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
	// context.TreeNodeVisited.Length() = context.TreeNodeVisited.Capacity // This array is accessed directly rather than behaving as a list
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
