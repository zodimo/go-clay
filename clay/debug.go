package clay

type Clay__DebugElementData struct {
	Collision bool
	Collapsed bool
}

var CLAY__DEBUGVIEW_COLOR_1 = Clay_Color{R: 58, G: 56, B: 52, A: 255}
var CLAY__DEBUGVIEW_COLOR_2 = Clay_Color{R: 62, G: 60, B: 58, A: 255}
var CLAY__DEBUGVIEW_COLOR_3 = Clay_Color{R: 141, G: 133, B: 135, A: 255}
var CLAY__DEBUGVIEW_COLOR_4 = Clay_Color{R: 238, G: 226, B: 231, A: 255}
var CLAY__DEBUGVIEW_COLOR_SELECTED_ROW = Clay_Color{R: 102, G: 80, B: 78, A: 255}

const CLAY__DEBUGVIEW_ROW_HEIGHT = 30
const CLAY__DEBUGVIEW_OUTER_PADDING = 10
const CLAY__DEBUGVIEW_INDENT_WIDTH = 16

var Clay__debugViewWidth int32 = 400

var Clay__debugViewHighlightColor = Clay_Color{R: 168, G: 66, B: 28, A: 100}

type Clay__RenderDebugLayoutData struct {
	RowCount                int32
	SelectedElementRowIndex int32
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

// 	CLAY("Clay__DebugView",
// 		Clay_ElementDeclaration{
// 			Layout: Clay_LayoutConfig{
// 				Sizing: Clay_Sizing{
// 					Width:  CLAY_SIZING_FIXED(float32(Clay__debugViewWidth)),
// 					Height: CLAY_SIZING_FIXED(currentContext.LayoutDimensions.Height),
// 				},
// 				LayoutDirection: CLAY_TOP_TO_BOTTOM,
// 			},
// 			Floating: Clay_FloatingElementConfig{
// 				ZIndex: 32765,
// 				AttachPoints: Clay_FloatingAttachPoints{
// 					Element: CLAY_ATTACH_POINT_LEFT_CENTER,
// 					Parent:  CLAY_ATTACH_POINT_RIGHT_CENTER},
// 				AttachTo: CLAY_ATTACH_TO_ROOT,
// 				ClipTo:   CLAY_CLIP_TO_ATTACHED_PARENT,
// 			},
// 			Border: Clay_BorderElementConfig{
// 				Color: CLAY__DEBUGVIEW_COLOR_3,
// 				Width: Clay_BorderWidth{Bottom: 1},
// 			},
// 		},

// 		CLAY_AUTO_ID(
// 			Clay_ElementDeclaration{
// 				Layout: Clay_LayoutConfig{
// 					Sizing: Clay_Sizing{
// 						Width:  CLAY_SIZING_GROW(0),
// 						Height: CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT),
// 					},
// 					Padding: Clay_Padding{
// 						Left:   CLAY__DEBUGVIEW_OUTER_PADDING,
// 						Right:  CLAY__DEBUGVIEW_OUTER_PADDING,
// 						Top:    0,
// 						Bottom: 0,
// 					},
// 					ChildAlignment: Clay_ChildAlignment{
// 						Y: CLAY_ALIGN_Y_CENTER,
// 					},
// 					BackgroundColor: CLAY__DEBUGVIEW_COLOR_2,
// 				},
// 			},

// 			CLAY_TEXT("Clay Debug Tools", infoTextConfig),
// 			CLAY_AUTO_ID(
// 				Clay_ElementDeclaration{
// 					Layout: Clay_LayoutConfig{
// 						Sizing: Clay_Sizing{
// 							Width: CLAY_SIZING_GROW(0),
// 						},
// 					},
// 				},
// 			),

// 			// Close button
// 			CLAY_AUTO_ID(
// 				Clay_ElementDeclaration{
// 					Layout: Clay_LayoutConfig{
// 						Sizing: Clay_Sizing{
// 							Width:  CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10),
// 							Height: CLAY_SIZING_FIXED(CLAY__DEBUGVIEW_ROW_HEIGHT - 10),
// 						},
// 						ChildAlignment: Clay_ChildAlignment{
// 							X: CLAY_ALIGN_X_CENTER,
// 							Y: CLAY_ALIGN_Y_CENTER,
// 						},
// 					},
// 					BackgroundColor: CLAY_RGBA(217, 91, 67, 80),
// 					CornerRadius: Clay_CornerRadius{
// 						TopLeft:     4,
// 						TopRight:    4,
// 						BottomLeft:  4,
// 						BottomRight: 4,
// 					},
// 					Border: Clay_BorderElementConfig{
// 						Color: CLAY_RGBA(217, 91, 67, 255),
// 						Width: Clay_BorderWidth{
// 							Bottom: 1,
// 						},
// 					},
// 				},
// 				CLAY_ON_HOVER(HandleDebugViewCloseButtonInteraction, 0),
// 				CLAY_TEXT("x", infoTextConfig),
// 			),
// 		),
// 	)
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
