package clay

import (
	"github.com/zodimo/clay-go/pkg/mem"
)

// CLAY_DLL_EXPORT Clay_Context* Clay_Initialize(Clay_Arena arena, Clay_Dimensions layoutDimensions, Clay_ErrorHandler errorHandler);

type Clay_Arena = mem.Arena

// Primarily created via the CLAY_ID(), CLAY_IDI(), CLAY_ID_LOCAL() and CLAY_IDI_LOCAL() macros.
// Represents a hashed string ID used for identifying and finding specific clay UI elements, required
// by functions such as Clay_PointerOver() and Clay_GetElementData().
type Clay_ElementId struct {
	Id       uint32      // The resulting hash generated from the other fields.
	Offset   uint32      // A numerical offset applied after computing the hash from stringId.
	BaseId   uint32      // A base hash value to start from, for example the parent element ID is used when calculating CLAY_ID_LOCAL().
	StringId Clay_String // The string id to hash.
}

type Clay_Vector2 struct {
	X float32
	Y float32
}

type Clay_LayoutElement struct {
	ChildrenOrTextContent Clay__LayoutElementChildrenOrTextContent
	Dimensions            Clay_Dimensions
	MinDimensions         Clay_Dimensions
	LayoutConfig          *Clay_LayoutConfig
	ElementConfigs        Clay__Slice[Clay_ElementConfig]
	Id                    uint32
	FloatingChildrenCount uint16
}

type Clay_SharedElementConfig struct {
	BackgroundColor Clay_Color
	CornerRadius    Clay_CornerRadius
	UserData        interface{}
}

type Clay__ElementConfigType uint8

const (
	CLAY__ELEMENT_CONFIG_TYPE_NONE Clay__ElementConfigType = iota
	CLAY__ELEMENT_CONFIG_TYPE_BORDER
	CLAY__ELEMENT_CONFIG_TYPE_FLOATING
	CLAY__ELEMENT_CONFIG_TYPE_CLIP
	CLAY__ELEMENT_CONFIG_TYPE_ASPECT
	CLAY__ELEMENT_CONFIG_TYPE_IMAGE
	CLAY__ELEMENT_CONFIG_TYPE_TEXT
	CLAY__ELEMENT_CONFIG_TYPE_CUSTOM
	CLAY__ELEMENT_CONFIG_TYPE_SHARED
)

type Clay_ElementConfigUnion struct {
	TextElementConfig        *Clay_TextElementConfig
	AspectRatioElementConfig *Clay_AspectRatioElementConfig
	ImageElementConfig       *Clay_ImageElementConfig
	FloatingElementConfig    *Clay_FloatingElementConfig
	CustomElementConfig      *Clay_CustomElementConfig
	ClipElementConfig        *Clay_ClipElementConfig
	BorderElementConfig      *Clay_BorderElementConfig
	SharedElementConfig      *Clay_SharedElementConfig
}

type Clay_ElementConfig struct {
	Type   Clay__ElementConfigType
	Config Clay_ElementConfigUnion
}

type Clay_LayoutElementTreeNode struct {
	LayoutElement   Clay_LayoutElement
	Position        Clay_Vector2
	NextChildOffset Clay_Vector2
}

type Clay_LayoutElementHashMapItem struct { // todo get this struct into a single cache line
	BoundingBox           Clay_BoundingBox
	ElementId             Clay_ElementId
	LayoutElement         *Clay_LayoutElement
	OnHoverFunction       func(elementId Clay_ElementId, pointerInfo Clay_PointerData, userData any)
	HoverFunctionUserData any
	NextIndex             int32
	Generation            uint32
	DebugData             *Clay__DebugElementData
}

type Clay_AspectRatioElementConfig struct {
	AspectRatio float32
}
type Clay_ImageElementConfig struct {
	ImageData interface{}
}

// Controls where a floating element is offset relative to its parent element.
// Note: see https://github.com/user-attachments/assets/b8c6dfaa-c1b1-41a4-be55-013473e4a6ce for a visual explanation.
type Clay_FloatingAttachPointType uint8

const (
	CLAY_ATTACH_POINT_LEFT_TOP Clay_FloatingAttachPointType = iota
	CLAY_ATTACH_POINT_LEFT_CENTER
	CLAY_ATTACH_POINT_LEFT_BOTTOM
	CLAY_ATTACH_POINT_CENTER_TOP
	CLAY_ATTACH_POINT_CENTER_CENTER
	CLAY_ATTACH_POINT_CENTER_BOTTOM
	CLAY_ATTACH_POINT_RIGHT_TOP
	CLAY_ATTACH_POINT_RIGHT_CENTER
	CLAY_ATTACH_POINT_RIGHT_BOTTOM
)

type Clay_FloatingAttachPoints struct {
	Element Clay_FloatingAttachPointType
	Parent  Clay_FloatingAttachPointType
}

// Controls which element a floating element is "attached" to (i.e. relative offset from).
type Clay_FloatingAttachToElement uint8

const (
	// (default) Disables floating for this element.
	CLAY_ATTACH_TO_NONE Clay_FloatingAttachToElement = iota
	// Attaches this floating element to its parent, positioned based on the .attachPoints and .offset fields.
	CLAY_ATTACH_TO_PARENT
	// Attaches this floating element to an element with a specific ID, specified with the .parentId field. positioned based on the .attachPoints and .offset fields.
	CLAY_ATTACH_TO_ELEMENT_WITH_ID
	// Attaches this floating element to the root of the layout, which combined with the .offset field provides functionality similar to "absolute positioning".
	CLAY_ATTACH_TO_ROOT
)

// Controls whether or not a floating element is clipped to the same clipping rectangle as the element it's attached to.
type Clay_FloatingClipToElement uint8

const (
	// (default) - The floating element does not inherit clipping.
	CLAY_CLIP_TO_NONE Clay_FloatingClipToElement = iota
	// The floating element is clipped to the same clipping rectangle as the element it's attached to.
	CLAY_CLIP_TO_ATTACHED_PARENT
)

type Clay_FloatingElementConfig struct {
	// Offsets this floating element by the provided x,y coordinates from its attachPoints.
	Offset Clay_Vector2
	// Expands the boundaries of the outer floating element without affecting its children.
	Expand Clay_Dimensions
	// When used in conjunction with .attachTo = CLAY_ATTACH_TO_ELEMENT_WITH_ID, attaches this floating element to the element in the hierarchy with the provided ID.
	// Hint: attach the ID to the other element with .id = CLAY_ID("yourId"), and specify the id the same way, with .parentId = CLAY_ID("yourId").id
	ParentId uint32
	// Controls the z index of this floating element and all its children. Floating elements are sorted in ascending z order before output.
	// zIndex is also passed to the renderer for all elements contained within this floating element.
	ZIndex int16
	// Controls how mouse pointer events like hover and click are captured or passed through to elements underneath / behind a floating element.
	// Enum is of the form CLAY_ATTACH_POINT_foo_bar. See Clay_FloatingAttachPoints for more details.
	// Note: see <img src="https://github.com/user-attachments/assets/b8c6dfaa-c1b1-41a4-be55-013473e4a6ce />
	// and <img src="https://github.com/user-attachments/assets/ebe75e0d-1904-46b0-982d-418f929d1516 /> for a visual explanation.
	AttachPoints Clay_FloatingAttachPoints
	// Controls how mouse pointer events like hover and click are captured or passed through to elements underneath a floating element.
	// CLAY_POINTER_CAPTURE_MODE_CAPTURE (default) - "Capture" the pointer event and don't allow events like hover and click to pass through to elements underneath.
	// CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH - Transparently pass through pointer events like hover and click to elements underneath the floating element.
	PointerCaptureMode Clay_PointerCaptureMode
	// Controls which element a floating element is "attached" to (i.e. relative offset from).
	// CLAY_ATTACH_TO_NONE (default) - Disables floating for this element.
	// CLAY_ATTACH_TO_PARENT - Attaches this floating element to its parent, positioned based on the .attachPoints and .offset fields.
	// CLAY_ATTACH_TO_ELEMENT_WITH_ID - Attaches this floating element to an element with a specific ID, specified with the .parentId field. positioned based on the .attachPoints and .offset fields.
	// CLAY_ATTACH_TO_ROOT - Attaches this floating element to the root of the layout, which combined with the .offset field provides functionality similar to "absolute positioning".
	AttachTo Clay_FloatingAttachToElement
	// Controls whether or not a floating element is clipped to the same clipping rectangle as the element it's attached to.
	// CLAY_CLIP_TO_NONE (default) - The floating element does not inherit clipping.
	// CLAY_CLIP_TO_ATTACHED_PARENT - The floating element is clipped to the same clipping rectangle as the element it's attached to.
	ClipTo Clay_FloatingClipToElement
}

// Controls various settings related to custom elements.
type Clay_CustomElementConfig struct {
	// A transparent pointer through which you can pass custom data to the renderer.
	// Generates CUSTOM render commands.
	CustomData interface{}
}

// Controls the axis on which an element switches to "scrolling", which clips the contents and allows scrolling in that direction.
type Clay_ClipElementConfig struct {
	Horizontal  bool         // Clip overflowing elements on the X axis.
	Vertical    bool         // Clip overflowing elements on the Y axis.
	ChildOffset Clay_Vector2 // Offsets the x,y positions of all child elements. Used primarily for scrolling containers.

}

type Clay_ElementDeclaration struct {
	Layout          Clay_LayoutConfig
	BackgroundColor Clay_Color
	CornerRadius    Clay_CornerRadius
	AspectRatio     Clay_AspectRatioElementConfig
	Image           Clay_ImageElementConfig
	Floating        Clay_FloatingElementConfig
	Custom          Clay_CustomElementConfig
	Clip            Clay_ClipElementConfig
	Border          Clay_BorderElementConfig
	UserData        interface{}
}

func Clay_Initialize(arena Clay_Arena, layoutDimensions Clay_Dimensions, errorHandler Clay_ErrorHandler) *Clay_Context {

	clay_Context := Clay__Context_Allocate_Arena(&arena)
	if clay_Context == nil {
		return nil
	}

	oldContext := Clay_GetCurrentContext()

	newContext := &Clay_Context{
		MaxElementCount:              Clay__defaultMaxElementCount,
		MaxMeasureTextCacheWordCount: Clay__defaultMaxMeasureTextWordCacheCount,
		ErrorHandler:                 Clay__ErrorHandlerFunctionDefault,
		LayoutDimensions:             layoutDimensions,
		InternalArena:                arena,
	}
	if oldContext != nil {
		if errorHandler.ErrorHandlerFunction != nil {
			newContext.ErrorHandler = errorHandler
		}

		newContext.MaxElementCount = oldContext.MaxElementCount
		newContext.MaxMeasureTextCacheWordCount = oldContext.MaxMeasureTextCacheWordCount

	}

	Clay_SetCurrentContext(newContext)
	Clay__InitializePersistentMemory(newContext)
	Clay__InitializeEphemeralMemory(newContext)

	// reset the hash maps
	for i := int32(0); i < newContext.LayoutElementsHashMap.Capacity; i++ {
		Clay__Array_Set(&newContext.LayoutElementsHashMap, i, -1)
	}
	for i := int32(0); i < newContext.MeasureTextHashMap.Capacity; i++ {
		Clay__Array_Set(&newContext.MeasureTextHashMap, i, 0)
	}

	return newContext
}

func Clay_BeginLayout() {

	currentContext := Clay_GetCurrentContext()
	Clay__InitializeEphemeralMemory(currentContext)
	currentContext.Generation++
	currentContext.DynamicElementIndex = 0
	// Set up the root container that covers the entire window
	rootDimensions := Clay_Dimensions{
		Width:  currentContext.LayoutDimensions.Width,
		Height: currentContext.LayoutDimensions.Height,
	}
	if currentContext.DebugModeEnabled {
		rootDimensions.Width -= float32(Clay__debugViewWidth)
	}

	currentContext.BooleanWarnings = Clay_BooleanWarnings{}
	Clay__OpenElementWithId(CLAY_ID("Clay__RootContainer"))
	Clay__ConfigureOpenElement(Clay_ElementDeclaration{
		Layout: Clay_LayoutConfig{
			Sizing: Clay_Sizing{
				Width:  CLAY_SIZING_FIXED(rootDimensions.Width),
				Height: CLAY_SIZING_FIXED(rootDimensions.Height),
			},
		},
	})

	Clay__Array_Add(&currentContext.OpenLayoutElementStack, 0)
	Clay__Array_Add(&currentContext.LayoutElementTreeRoots, Clay__LayoutElementTreeRoot{
		LayoutElementIndex: 0,
	})

}

func Clay_EndLayout() []Clay_RenderCommand {

	currentContext := Clay_GetCurrentContext()
	Clay__CloseElement()
	elementsExceededBeforeDebugView := currentContext.BooleanWarnings.MaxElementsExceeded
	if currentContext.DebugModeEnabled && !elementsExceededBeforeDebugView {
		currentContext.WarningsEnabled = false
		Clay__RenderDebugView()
		currentContext.WarningsEnabled = true
	}
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		var message Clay_String
		if !elementsExceededBeforeDebugView {
			message = CLAY_STRING("Clay Error: Layout elements exceeded Clay__maxElementCount after adding the debug-view to the layout.")
		} else {
			message = CLAY_STRING("Clay Error: Layout elements exceeded Clay__maxElementCount")
		}
		Clay__AddRenderCommand(Clay_RenderCommand{
			BoundingBox: Clay_BoundingBox{
				X:      currentContext.LayoutDimensions.Width/2 - 59*4,
				Y:      currentContext.LayoutDimensions.Height / 2,
				Width:  0,
				Height: 0,
			},
			RenderData: Clay_RenderData{
				Text: Clay_TextRenderData{
					StringContents: Clay_StringSlice{
						Length:    message.Length,
						Chars:     message.Chars,
						BaseChars: message.Chars,
					},
					TextColor: Clay_Color{R: 255, G: 0, B: 0, A: 255},
					FontSize:  16,
				},
			},
			CommandType: CLAY_RENDER_COMMAND_TYPE_TEXT,
		})
	}
	if currentContext.OpenLayoutElementStack.Length > 1 {
		currentContext.ErrorHandler.ErrorHandlerFunction(Clay_ErrorData{
			ErrorType: CLAY_ERROR_TYPE_UNBALANCED_OPEN_CLOSE,
			ErrorText: CLAY_STRING("There were still open layout elements when EndLayout was called. This results from an unequal number of calls to Clay__OpenElement and Clay__CloseElement."),
			UserData:  currentContext.ErrorHandler.UserData,
		})
	}
	Clay__CalculateFinalLayout()
	return mem.MArray_GetAll(&currentContext.RenderCommands)

}

func Clay_SetLayoutDimensions(dimensions Clay_Dimensions) {
	currentContext := Clay_GetCurrentContext()
	currentContext.LayoutDimensions = dimensions
}
