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
	Elements []int32
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

	openLayoutElement.ChildrenOrTextContent.Children.Elements = mem.MArray_GetSlice(
		&currentContext.LayoutElementChildren,
		currentContext.LayoutElementChildren.Length(),
		currentContext.LayoutElementChildren.Capacity(),
	)

	if layoutConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT {
		openLayoutElement.Dimensions.Width = leftRightPadding
		openLayoutElement.MinDimensions.Width = leftRightPadding
		if openLayoutElement.ChildrenOrTextContent.Children.Length > 0 && int32(openLayoutElement.ChildrenOrTextContent.Children.Length) <= currentContext.LayoutElementChildrenBuffer.Length() {
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
		}

		childGap := float32(CLAY__MAX(float32(openLayoutElement.ChildrenOrTextContent.Children.Length-1), 0) * float32(layoutConfig.ChildGap))

		openLayoutElement.Dimensions.Width += childGap
		if !elementHasClipHorizontal {
			openLayoutElement.MinDimensions.Width += childGap
		}
	} else if layoutConfig.LayoutDirection == CLAY_TOP_TO_BOTTOM {
		openLayoutElement.Dimensions.Height = topBottomPadding
		openLayoutElement.MinDimensions.Height = topBottomPadding
		if openLayoutElement.ChildrenOrTextContent.Children.Length > 0 && int32(openLayoutElement.ChildrenOrTextContent.Children.Length) <= currentContext.LayoutElementChildrenBuffer.Length() {
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
		}

		childGap := float32(CLAY__MAX(float32(openLayoutElement.ChildrenOrTextContent.Children.Length-1), 0) * float32(layoutConfig.ChildGap))

		openLayoutElement.Dimensions.Height += childGap
		if !elementHasClipVertical {
			openLayoutElement.MinDimensions.Height += childGap
		}
	}

	// Resize the Elements slice to match the Length field after adding children
	if openLayoutElement.ChildrenOrTextContent.Children.Length > 0 {
		if int32(openLayoutElement.ChildrenOrTextContent.Children.Length) <= currentContext.LayoutElementChildrenBuffer.Length() {
			openLayoutElement.ChildrenOrTextContent.Children.Elements = openLayoutElement.ChildrenOrTextContent.Children.Elements[:openLayoutElement.ChildrenOrTextContent.Children.Length]
			Clay__Array_Shrink(&currentContext.LayoutElementChildrenBuffer, int32(openLayoutElement.ChildrenOrTextContent.Children.Length))
		} else {
			// If buffer doesn't have enough elements, set Elements to an empty slice
			openLayoutElement.ChildrenOrTextContent.Children.Elements = []int32{}
			openLayoutElement.ChildrenOrTextContent.Children.Length = 0
		}
	} else {
		// If there are no children, set Elements to an empty slice
		openLayoutElement.ChildrenOrTextContent.Children.Elements = []int32{}
	}

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

	// Get the currently open parent (only if stack is not empty)
	if currentContext.OpenLayoutElementStack.Length() > 0 {
		openLayoutElement = Clay__GetOpenLayoutElement()
	}

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

func Clay__AddRenderCommand(renderCommand Clay_RenderCommand) {
	if renderCommand.CommandType == CLAY_RENDER_COMMAND_TYPE_NONE {
		//The NONE type does not have defined expectations from the renderer, so we can safely ignore it
		fmt.Println("Attempted to add a render command with a command type of NONE")
		return
	}
	currentContext := Clay_GetCurrentContext()
	if currentContext.RenderCommands.Length() < currentContext.RenderCommands.Capacity()-1 {
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
	if currentContext.LayoutElements.Length() == currentContext.LayoutElements.Capacity()-1 || currentContext.BooleanWarnings.MaxElementsExceeded {
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
	if currentContext.LayoutElements.Length() == currentContext.LayoutElements.Capacity()-1 || currentContext.BooleanWarnings.MaxElementsExceeded {
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
	// openLayoutElement := Clay__GetOpenLayoutElement()
	// Clay__Slice_Grow(&openLayoutElement.ElementConfigs, 1)
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

	// Capture the starting length before attaching configs
	elementConfigsStartLength := currentContext.ElementConfigs.Length()

	var sharedConfig *Clay_SharedElementConfig = nil

	if elementDeclaration.BackgroundColor.A > 0 {
		sharedConfig = new(Clay_SharedElementConfig)
		sharedConfig.BackgroundColor = elementDeclaration.BackgroundColor
		Clay__AttachElementConfig(Clay_ElementConfigUnion{SharedElementConfig: sharedConfig}, CLAY__ELEMENT_CONFIG_TYPE_SHARED)
	}
	if !Clay__MemCmpTyped(&elementDeclaration.CornerRadius, &Clay__CornerRadius_DEFAULT) {
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
		fmt.Println("Attaching image data", elementDeclaration.Image.ImageData)
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
	if !Clay__MemCmpTyped(&elementDeclaration.Border.Width, &Clay__BorderWidth_DEFAULT) {
		Clay__AttachElementConfig(Clay_ElementConfigUnion{
			BorderElementConfig: Clay__StoreBorderElementConfig(elementDeclaration.Border),
		}, CLAY__ELEMENT_CONFIG_TYPE_BORDER)
	}

	// Set the element configs slice AFTER all configs have been attached
	elementConfigs, err := mem.CreateSliceFromRange(&currentContext.ElementConfigs, elementConfigsStartLength, currentContext.ElementConfigs.Length())
	if err != nil {
		panic(err)
	}
	openLayoutElement.ElementConfigs = elementConfigs
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
	currentContext := Clay_GetCurrentContext()
	if MeasureTextFunction == nil {
		if !currentContext.BooleanWarnings.TextMeasurementFunctionNotSet {
			currentContext.BooleanWarnings.TextMeasurementFunctionNotSet = true
			currentContext.ErrorHandler.ErrorHandlerFunction(Clay_ErrorData{
				ErrorType: CLAY_ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED,
				ErrorText: CLAY_STRING("Clay's internal MeasureText function is null. You may have forgotten to call Clay_SetMeasureTextFunction(), or passed a NULL function pointer by mistake."),
				UserData:  currentContext.ErrorHandler.UserData})
		}
		return &Clay__MeasureTextCacheItem_DEFAULT
	}

	id := Clay__HashStringContentsWithConfig(text, textConfig)
	hashBucket := int32(id % (uint32(currentContext.MaxMeasureTextCacheWordCount) / 32))
	elementIndexPrevious := int32(0)
	elementIndex := Clay__Array_GetValue[int32](&currentContext.MeasureTextHashMap, hashBucket)
	for elementIndex != 0 {
		hashEntry := Clay__Array_Get[Clay__MeasureTextCacheItem](&currentContext.MeasureTextHashMapInternal, elementIndex)
		if hashEntry.Id == id {
			hashEntry.Generation = currentContext.Generation
			return hashEntry
		}
		// This element hasn't been seen in a few frames, delete the hash map item
		if currentContext.Generation-hashEntry.Generation > 2 {
			// Add all the measured words that were included in this measurement to the freelist
			nextWordIndex := hashEntry.MeasuredWordsStartIndex
			for nextWordIndex != -1 {
				measuredWord := Clay__Array_Get[Clay__MeasuredWord](&currentContext.MeasuredWords, nextWordIndex)
				Clay__Array_Add(&currentContext.MeasuredWordsFreeList, nextWordIndex)
				nextWordIndex = measuredWord.Next
			}

			nextIndex := hashEntry.NextIndex
			Clay__Array_Set(&currentContext.MeasureTextHashMapInternal, elementIndex, Clay__MeasureTextCacheItem{MeasuredWordsStartIndex: -1}) //@TODO review if -1 is correct
			Clay__Array_Add(&currentContext.MeasureTextHashMapInternalFreeList, elementIndex)
			if elementIndexPrevious == 0 {
				Clay__Array_Set(&currentContext.MeasureTextHashMap, hashBucket, nextIndex)
			} else {
				previousHashEntry := Clay__Array_Get[Clay__MeasureTextCacheItem](&currentContext.MeasureTextHashMapInternal, elementIndexPrevious)
				previousHashEntry.NextIndex = nextIndex
			}
			elementIndex = nextIndex
		} else {
			elementIndexPrevious = elementIndex
			elementIndex = hashEntry.NextIndex
		}
	}

	newItemIndex := int32(0)
	newCacheItem := Clay__MeasureTextCacheItem{MeasuredWordsStartIndex: -1, Id: id, Generation: currentContext.Generation}
	measured := new(Clay__MeasureTextCacheItem)
	if currentContext.MeasureTextHashMapInternalFreeList.Length() > 0 {
		newItemIndex = Clay__Array_GetValue[int32](&currentContext.MeasureTextHashMapInternalFreeList, currentContext.MeasureTextHashMapInternalFreeList.Length()-1)
		Clay__Array_Shrink(&currentContext.MeasureTextHashMapInternalFreeList, 1)
		Clay__Array_Set(&currentContext.MeasureTextHashMapInternal, newItemIndex, newCacheItem)
		measured = Clay__Array_Get[Clay__MeasureTextCacheItem](&currentContext.MeasureTextHashMapInternal, newItemIndex)
	} else {
		if currentContext.MeasureTextHashMapInternal.Length() == currentContext.MeasureTextHashMapInternal.Capacity()-1 {
			if !currentContext.BooleanWarnings.MaxTextMeasureCacheExceeded {
				currentContext.ErrorHandler.ErrorHandlerFunction(Clay_ErrorData{
					ErrorType: CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED,
					ErrorText: CLAY_STRING("Clay ran out of capacity while attempting to measure text elements. Try using Clay_SetMaxElementCount() with a higher value."),
					UserData:  currentContext.ErrorHandler.UserData})
				currentContext.BooleanWarnings.MaxTextMeasureCacheExceeded = true
			}
			return &Clay__MeasureTextCacheItem_DEFAULT
		}
		measured = Clay__Array_Add(&currentContext.MeasureTextHashMapInternal, newCacheItem)
		newItemIndex = currentContext.MeasureTextHashMapInternal.Length() - 1
	}

	start := int32(0)
	end := int32(0)
	lineWidth := float32(0)
	measuredWidth := float32(0)
	measuredHeight := float32(0)
	spaceWidth := Clay__MeasureText(Clay_StringSlice{
		Length:    1,
		Chars:     CLAY__SPACECHAR.Chars,
		BaseChars: CLAY__SPACECHAR.Chars,
	},
		textConfig,
		currentContext.MeasureTextUserData,
	).Width
	tempWord := Clay__MeasuredWord{Next: -1}
	previousWord := &tempWord
	for end < text.Length {
		if currentContext.MeasuredWords.Length() == currentContext.MeasuredWords.Capacity()-1 {
			if !currentContext.BooleanWarnings.MaxTextMeasureCacheExceeded {
				currentContext.ErrorHandler.ErrorHandlerFunction(Clay_ErrorData{
					ErrorType: CLAY_ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED,
					ErrorText: CLAY_STRING("Clay has run out of space in it's internal text measurement cache. Try using Clay_SetMaxMeasureTextCacheWordCount() (default 16384, with 1 unit storing 1 measured word)."),
					UserData:  currentContext.ErrorHandler.UserData,
				})
				currentContext.BooleanWarnings.MaxTextMeasureCacheExceeded = true
			}
			return &Clay__MeasureTextCacheItem_DEFAULT
		}
		current := text.Chars[end]
		if current == ' ' || current == '\n' {
			length := end - start
			dimensions := Clay_Dimensions{}
			if length > 0 {
				dimensions = Clay__MeasureText(Clay_StringSlice{
					Length:    length,
					Chars:     text.Chars[start:],
					BaseChars: text.Chars,
				},
					textConfig,
					currentContext.MeasureTextUserData,
				)
			}
			measured.MinWidth = CLAY__MAX(dimensions.Width, measured.MinWidth)
			measuredHeight = CLAY__MAX(measuredHeight, dimensions.Height)
			if current == ' ' {
				dimensions.Width += spaceWidth
				previousWord = Clay__AddMeasuredWord(Clay__MeasuredWord{
					StartOffset: start,
					Length:      length + 1,
					Width:       dimensions.Width,
					Next:        -1,
				}, previousWord)
				lineWidth += dimensions.Width
			}
			if current == '\n' {
				if length > 0 {
					previousWord = Clay__AddMeasuredWord(Clay__MeasuredWord{
						StartOffset: start,
						Length:      length,
						Width:       dimensions.Width,
						Next:        -1,
					}, previousWord)
				}
				previousWord = Clay__AddMeasuredWord(Clay__MeasuredWord{
					StartOffset: end + 1,
					Length:      0,
					Width:       0,
					Next:        -1,
				}, previousWord)
				lineWidth += dimensions.Width
				measuredWidth = CLAY__MAX(lineWidth, measuredWidth)
				measured.ContainsNewlines = true
				lineWidth = 0
			}
			start = end + 1
		}
		end++
	}

	if end-start > 0 {
		dimensions := Clay__MeasureText(Clay_StringSlice{
			Length:    end - start,
			Chars:     text.Chars[start:],
			BaseChars: text.Chars,
		},
			textConfig,
			currentContext.MeasureTextUserData,
		)
		Clay__AddMeasuredWord(Clay__MeasuredWord{
			StartOffset: start,
			Length:      end - start,
			Width:       dimensions.Width,
			Next:        -1,
		}, previousWord)
		lineWidth += dimensions.Width
		measuredHeight = CLAY__MAX(measuredHeight, dimensions.Height)
		measured.MinWidth = CLAY__MAX(dimensions.Width, measured.MinWidth)
	}
	measuredWidth = CLAY__MAX(lineWidth, measuredWidth) - float32(textConfig.LetterSpacing)

	measured.MeasuredWordsStartIndex = tempWord.Next
	measured.UnwrappedDimensions.Width = measuredWidth
	measured.UnwrappedDimensions.Height = measuredHeight

	if elementIndexPrevious != 0 {
		Clay__Array_Get[Clay__MeasureTextCacheItem](&currentContext.MeasureTextHashMapInternal, elementIndexPrevious).NextIndex = newItemIndex
	} else {
		Clay__Array_Set(&currentContext.MeasureTextHashMap, hashBucket, newItemIndex)
	}
	return measured

}

func Clay__AddMeasuredWord(word Clay__MeasuredWord, previousWord *Clay__MeasuredWord) *Clay__MeasuredWord {
	currentContext := Clay_GetCurrentContext()
	if currentContext.MeasuredWordsFreeList.Length() > 0 {
		newItemIndex := Clay__Array_GetValue[int32](&currentContext.MeasuredWordsFreeList, currentContext.MeasuredWordsFreeList.Length()-1)
		Clay__Array_Shrink(&currentContext.MeasuredWordsFreeList, 1)
		Clay__Array_Set(&currentContext.MeasuredWords, newItemIndex, word)
		previousWord.Next = newItemIndex
		return Clay__Array_Get[Clay__MeasuredWord](&currentContext.MeasuredWords, newItemIndex)
	} else {
		previousWord.Next = currentContext.MeasuredWords.Length()
		return Clay__Array_Add(&currentContext.MeasuredWords, word)
	}
}

func Clay__OpenTextElement(text Clay_String, textConfig *Clay_TextElementConfig) {
	currentContext := Clay_GetCurrentContext()
	if currentContext.LayoutElements.Length() == currentContext.LayoutElements.Capacity()-1 || currentContext.BooleanWarnings.MaxElementsExceeded {
		currentContext.BooleanWarnings.MaxElementsExceeded = true
		return
	}
	parentElement := Clay__GetOpenLayoutElement()

	layoutElement := Clay_LayoutElement{}

	textElement := Clay__Array_Add[Clay_LayoutElement](&currentContext.LayoutElements, layoutElement)

	if currentContext.OpenClipElementStack.Length() > 0 {
		if currentContext.LayoutElementClipElementIds.Length() == 0 {
			Clay__Array_Add(&currentContext.LayoutElementClipElementIds, Clay__Array_GetValue[int32](&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length()-1))
		} else {
			Clay__Array_Set(&currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length()-1, Clay__Array_GetValue[int32](&currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length()-1))
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

func Clay__InitializePersistentMemory(context *Clay_Context) {
	// Persistent memory - initialized once and not reset
	maxElementCount := context.MaxElementCount
	maxMeasureTextCacheWordCount := context.MaxMeasureTextCacheWordCount
	arena := &context.InternalArena

	context.ScrollContainerDatas = Clay__Array_Allocate_Arena[Clay__ScrollContainerDataInternal](100, arena)
	context.LayoutElementsHashMapInternal = Clay__Array_Allocate_Arena[Clay_LayoutElementHashMapItem](maxElementCount, arena, mem.MemArrayWithZeroValuePtr[Clay_LayoutElementHashMapItem](&Clay_LayoutElementHashMapItem_DEFAULT))
	context.LayoutElementsHashMap = Clay__Array_Allocate_Arena[int32](maxElementCount, arena, mem.MemArrayWithIsHashmap[int32](), mem.MemArrayWithZeroValue[int32](-1))
	context.MeasureTextHashMapInternal = Clay__Array_Allocate_Arena[Clay__MeasureTextCacheItem](maxElementCount, arena, mem.MemArrayWithZeroValuePtr[Clay__MeasureTextCacheItem](&Clay__MeasureTextCacheItem_DEFAULT))
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
	context.TreeNodeVisited = Clay__Array_Allocate_Arena[bool](maxElementCount, arena, mem.MemArrayWithInitialLength[bool](maxElementCount))
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
