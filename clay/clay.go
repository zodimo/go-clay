package clay

import (
	"hash/fnv"

	"github.com/zodimo/go-arena-memory/mem"
)

// CLAY_DLL_EXPORT Clay_Context* Clay_Initialize(Clay_Arena arena, Clay_Dimensions layoutDimensions, Clay_ErrorHandler errorHandler);

type Clay_Arena = mem.Arena

var Clay__currentContext *Clay_Context

type Clay_ElementId = uint32
type Clay_String = string

type Clay_Dimensions struct {
	Width  float32
	Height float32
}
type Clay_ErrorHandler struct {
	ErrorHandlerFunction func(errorData Clay_ErrorData)
	UserData             interface{}
}

// Controls how text "wraps", that is how it is broken into multiple lines when there is insufficient horizontal space.
type Clay_TextElementConfigWrapMode uint8

const (
	// (default) breaks on whitespace characters.
	CLAY_TEXT_WRAP_WORDS Clay_TextElementConfigWrapMode = iota
	// Don't break on space characters, only on newlines.
	CLAY_TEXT_WRAP_NEWLINES
	// Disable text wrapping entirely.
	CLAY_TEXT_WRAP_NONE
)

// Controls how wrapped lines of text are horizontally aligned within the outer text bounding box.
type Clay_TextAlignment uint8

const (
	// (default) Horizontally aligns wrapped lines of text to the left hand side of their bounding box.
	CLAY_TEXT_ALIGN_LEFT Clay_TextAlignment = iota
	// Horizontally aligns wrapped lines of text to the center of their bounding box.
	CLAY_TEXT_ALIGN_CENTER
	// Horizontally aligns wrapped lines of text to the right hand side of their bounding box.
	CLAY_TEXT_ALIGN_RIGHT
)

// Controls various functionality related to text elements.
type Clay_TextElementConfig struct {
	// A pointer that will be transparently passed through to the resulting render command.
	UserData interface{}
	// The RGBA color of the font to render, conventionally specified as 0-255.
	TextColor Clay_Color
	// An integer transparently passed to Clay_MeasureText to identify the font to use.
	// The debug view will pass fontId = 0 for its internal text.
	FontId uint16
	// Controls the size of the font. Handled by the function provided to Clay_MeasureText.
	FontSize uint16
	// Controls extra horizontal spacing between characters. Handled by the function provided to Clay_MeasureText.
	LetterSpacing uint16
	// Controls additional vertical space between wrapped lines of text.
	LineHeight uint16
	// Controls how text "wraps", that is how it is broken into multiple lines when there is insufficient horizontal space.
	// CLAY_TEXT_WRAP_WORDS (default) breaks on whitespace characters.
	// CLAY_TEXT_WRAP_NEWLINES doesn't break on space characters, only on newlines.
	// CLAY_TEXT_WRAP_NONE disables wrapping entirely.
	WrapMode Clay_TextElementConfigWrapMode
	// Controls how wrapped lines of text are horizontally aligned within the outer text bounding box.
	// CLAY_TEXT_ALIGN_LEFT (default) - Horizontally aligns wrapped lines of text to the left hand side of their bounding box.
	// CLAY_TEXT_ALIGN_CENTER - Horizontally aligns wrapped lines of text to the center of their bounding box.
	// CLAY_TEXT_ALIGN_RIGHT - Horizontally aligns wrapped lines of text to the right hand side of their bounding box.
	TextAlignment Clay_TextAlignment
}

const Clay__defaultMaxElementCount int32 = 8192
const Clay__defaultMaxMeasureTextWordCacheCount int32 = 16384

type Clay_ErrorType uint8

// Represents the type of error clay encountered while computing layout.
const (
	// A text measurement function wasn't provided using Clay_SetMeasureTextFunction(), or the provided function was null.
	CLAY_ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED Clay_ErrorType = iota
	// Clay attempted to allocate its internal data structures but ran out of space.
	// The arena passed to Clay_Initialize was created with a capacity smaller than that required by Clay_MinMemorySize().
	CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED
	// Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxElementCount().
	CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED
	// Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxMeasureTextCacheWordCount().
	CLAY_ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED
	// Two elements were declared with exactly the same ID within one layout.
	CLAY_ERROR_TYPE_DUPLICATE_ID
	// A floating element was declared using CLAY_ATTACH_TO_ELEMENT_ID and either an invalid .parentId was provided or no element with the provided .parentId was found.
	CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND
	// An element was declared that using CLAY_SIZING_PERCENT but the percentage value was over 1. Percentage values are expected to be in the 0-1 range.
	CLAY_ERROR_TYPE_PERCENTAGE_OVER_1
	// Clay encountered an internal error. It would be wonderful if you could report this so we can fix it!
	CLAY_ERROR_TYPE_INTERNAL_ERROR
	// Clay__OpenElement was called more times than Clay__CloseElement, so there were still remaining open elements when the layout ended.
	CLAY_ERROR_TYPE_UNBALANCED_OPEN_CLOSE
)

type Clay_ErrorData struct {
	// Represents the type of error clay encountered while computing layout.
	// CLAY_ERROR_TYPE_TEXT_MEASUREMENT_FUNCTION_NOT_PROVIDED - A text measurement function wasn't provided using Clay_SetMeasureTextFunction(), or the provided function was null.
	// CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED - Clay attempted to allocate its internal data structures but ran out of space. The arena passed to Clay_Initialize was created with a capacity smaller than that required by Clay_MinMemorySize().
	// CLAY_ERROR_TYPE_ELEMENTS_CAPACITY_EXCEEDED - Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxElementCount().
	// CLAY_ERROR_TYPE_TEXT_MEASUREMENT_CAPACITY_EXCEEDED - Clay ran out of capacity in its internal array for storing elements. This limit can be increased with Clay_SetMaxMeasureTextCacheWordCount().
	// CLAY_ERROR_TYPE_DUPLICATE_ID - Two elements were declared with exactly the same ID within one layout.
	// CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND - A floating element was declared using CLAY_ATTACH_TO_ELEMENT_ID and either an invalid .parentId was provided or no element with the provided .parentId was found.
	// CLAY_ERROR_TYPE_PERCENTAGE_OVER_1 - An element was declared that using CLAY_SIZING_PERCENT but the percentage value was over 1. Percentage values are expected to be in the 0-1 range.
	// CLAY_ERROR_TYPE_INTERNAL_ERROR - Clay encountered an internal error. It would be wonderful if you could report this so we can fix it!
	ErrorType Clay_ErrorType

	// A string containing human-readable error text that explains the error in more detail.
	ErrorText Clay_String
	// A transparent pointer passed through from when the error handler was first provided.
	UserData interface{}
}

type Clay_Vector2 struct {
	X float32
	Y float32
}

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

var Clay__ErrorHandlerFunctionDefault = Clay_ErrorHandler{
	ErrorHandlerFunction: func(errorData Clay_ErrorData) {
		// Do nothing
	},
	UserData: nil,
}

type Clay_ArraySlice[T any] struct {
	Length        int32
	InternalArray []T
}

type Clay__Array[T any] struct {
	Capacity      int32
	Length        int32
	InternalArray []T
}

func Clay__Array_Allocate_Arena[T any](capacity int32) Clay__Array[T] {
	return Clay__Array[T]{
		Capacity:      capacity,
		Length:        0,
		InternalArray: make([]T, capacity),
	}
}

type Clay_LayoutElement struct {
	// union {
	//     Clay__LayoutElementChildren children;
	//     Clay__TextElementData *textElementData;
	// } childrenOrTextContent;
	Dimensions            Clay_Dimensions
	MinDimensions         Clay_Dimensions
	LayoutConfig          Clay_LayoutConfig
	ElementConfigs        Clay__Array[Clay_ElementConfig] // slice
	Id                    uint32
	FloatingChildrenCount uint16
}

// Used by renderers to determine specific handling for each render command.
type Clay_RenderCommandType uint8

const (
	CLAY_RENDER_COMMAND_TYPE_NONE Clay_RenderCommandType = iota
	CLAY_RENDER_COMMAND_TYPE_RECTANGLE
	CLAY_RENDER_COMMAND_TYPE_BORDER
	CLAY_RENDER_COMMAND_TYPE_TEXT
	CLAY_RENDER_COMMAND_TYPE_IMAGE
	CLAY_RENDER_COMMAND_TYPE_SCISSOR_START
	CLAY_RENDER_COMMAND_TYPE_SCISSOR_END
	CLAY_RENDER_COMMAND_TYPE_CUSTOM
)

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_RECTANGLE
type Clay_RectangleRenderData struct {
	// The solid background color to fill this rectangle with. Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	BackgroundColor Clay_Color
	// Controls the "radius", or corner rounding of elements, including rectangles, borders and images.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	CornerRadius Clay_CornerRadius
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_TEXT
type Clay_TextRenderData struct {
	// A string slice containing the text to be rendered.
	// Note: this is not guaranteed to be null terminated.
	StringContents Clay_String // Slice
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	TextColor Clay_Color
	// An integer representing the font to use to render this text, transparently passed through from the text declaration.
	FontId   uint16
	FontSize uint16
	// Specifies the extra whitespace gap in pixels between each character.
	LetterSpacing uint16
	// The height of the bounding box for this line of text.
	LineHeight uint16
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_IMAGE
type Clay_ImageRenderData struct {
	// The tint color for this image. Note that the default value is 0,0,0,0 and should likely be interpreted
	// as "untinted".
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	BackgroundColor Clay_Color
	// Controls the "radius", or corner rounding of this image.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	CornerRadius Clay_CornerRadius
	// A pointer transparently passed through from the original element definition, typically used to represent image data.
	ImageData interface{}
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_CUSTOM
type Clay_CustomRenderData struct {
	// Passed through from .backgroundColor in the original element declaration.
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	BackgroundColor Clay_Color
	// Controls the "radius", or corner rounding of this custom element.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	CornerRadius Clay_CornerRadius
	// A pointer transparently passed through from the original element definition.
	CustomData interface{}
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_BORDER
type Clay_BorderRenderData struct {
	// Controls a shared color for all this element's borders.
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	Color Clay_Color
	// Specifies the "radius", or corner rounding of this border element.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	CornerRadius Clay_CornerRadius
	// Controls individual border side widths.
	Width Clay_BorderWidth
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_START || commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_END
type Clay_ClipRenderData struct {
	Horizontal bool
	Vertical   bool
}

type Clay_RenderData struct {
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_RECTANGLE
	Rectangle Clay_RectangleRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_TEXT
	Text Clay_TextRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_IMAGE
	Image Clay_ImageRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_CUSTOM
	Custom Clay_CustomRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_BORDER
	Border Clay_BorderRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_START|END
	Clip Clay_ClipRenderData
}

type Clay_RenderCommand struct {
	// A rectangular box that fully encloses this UI element, with the position relative to the root of the layout.
	BoundingBox Clay_BoundingBox
	// A struct union containing data specific to this command's commandType.
	RenderData Clay_RenderData
	// A pointer transparently passed through from the original element declaration.
	UserData interface{}
	// The id of this element, transparently passed through from the original element declaration.
	Id uint32
	// The z order required for drawing this command correctly.
	// Note: the render command array is already sorted in ascending order, and will produce correct results if drawn in naive order.
	// This field is intended for use in batching renderers for improved performance.
	ZIndex int16
	// Specifies how to handle rendering of this command.
	// CLAY_RENDER_COMMAND_TYPE_RECTANGLE - The renderer should draw a solid color rectangle.
	// CLAY_RENDER_COMMAND_TYPE_BORDER - The renderer should draw a colored border inset into the bounding box.
	// CLAY_RENDER_COMMAND_TYPE_TEXT - The renderer should draw text.
	// CLAY_RENDER_COMMAND_TYPE_IMAGE - The renderer should draw an image.
	// CLAY_RENDER_COMMAND_TYPE_SCISSOR_START - The renderer should begin clipping all future draw commands, only rendering content that falls within the provided boundingBox.
	// CLAY_RENDER_COMMAND_TYPE_SCISSOR_END - The renderer should finish any previously active clipping, and begin rendering elements in full again.
	// CLAY_RENDER_COMMAND_TYPE_CUSTOM - The renderer should provide a custom implementation for handling this render command based on its .customData
	CommandType Clay_RenderCommandType
}

func (a *Clay__Array[T]) Get(index int32) T {
	return a.InternalArray[index]
}
func (a *Clay__Array[T]) Add(value T) {
	a.InternalArray[a.Length] = value
	a.Length++
}
func (a *Clay__Array[T]) Set(index int32, value T) {
	if index < a.Length {
		a.InternalArray[index] = value
	}
}

type Clay_BooleanWarnings struct {
	MaxElementsExceeded           bool
	MaxRenderCommandsExceeded     bool
	MaxTextMeasureCacheExceeded   bool
	TextMeasurementFunctionNotSet bool
}

type Clay__Warning struct {
	BaseMessage    Clay_String
	DynamicMessage Clay_String
}

type Clay__WrappedTextLine struct {
	Line       Clay_String
	Dimensions Clay_Dimensions
}
type Clay__TextElementData struct {
	Text                Clay_String
	PreferredDimensions Clay_Dimensions
	ElementIndex        int32
	WrappedLines        Clay_ArraySlice[Clay__WrappedTextLine]
}

type Clay_SharedElementConfig struct {
	BackgroundColor Clay_Color
	CornerRadius    Clay_CornerRadius
}

type Clay_ElementConfigType uint8

const (
	CLAY__ELEMENT_CONFIG_TYPE_NONE Clay_ElementConfigType = iota
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
	Type   Clay_ElementConfigType
	Config Clay_ElementConfigUnion
}

type Clay_LayoutElementTreeNode struct {
	LayoutElement   Clay_LayoutElement
	Position        Clay_Vector2
	NextChildOffset Clay_Vector2
}

type Clay_LayoutElementTreeRoot struct {
	LayoutElementIndex int32
	ParentId           uint32 // This can be zero in the case of the root layout tree
	ClipElementId      uint32 // This can be zero if there is no clip element
	ZIndex             int16
	PointerOffset      Clay_Vector2 // Only used when scroll containers are managed externally
}

type Clay_BoundingBox struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

type Clay__DebugElementData struct {
	Collision bool
	Collapsed bool
}

type Clay_LayoutElementHashMapItem struct { // todo get this struct into a single cache line
	BoundingBox           Clay_BoundingBox
	ElementId             Clay_ElementId
	LayoutElement         Clay_LayoutElement
	OnHoverFunction       func(elementId Clay_ElementId, pointerInfo Clay_PointerData, userData any)
	HoverFunctionUserData any
	NextIndex             int32
	Generation            uint32
	DebugData             Clay__DebugElementData
}

type Clay__MeasuredWord struct {
	StartOffset int32
	Length      int32
	Width       float32
	Next        int32
}

type Clay__ScrollContainerDataInternal struct {
	LayoutElement       Clay_LayoutElement
	BoundingBox         Clay_BoundingBox
	ContentSize         Clay_Dimensions
	ScrollOrigin        Clay_Vector2
	PointerOrigin       Clay_Vector2
	ScrollMomentum      Clay_Vector2
	ScrollPosition      Clay_Vector2
	PreviousDelta       Clay_Vector2
	MomentumTime        float32
	ElementId           Clay_ElementId
	OpenThisFrame       bool
	PointerScrollActive bool
}

type Clay_Context struct {
	MaxElementCount              int32
	MaxMeasureTextCacheWordCount int32
	WarningsEnabled              bool
	ErrorHandler                 Clay_ErrorHandler

	BooleanWarnings Clay_BooleanWarnings
	Warnings        Clay__Array[Clay__Warning]

	PointerInfo      Clay_PointerData
	LayoutDimensions Clay_Dimensions

	DynamicElementIndexBaseHash Clay_ElementId
	DynamicElementIndex         int32

	DebugModeEnabled              bool
	DisableCulling                bool
	ExternalScrollHandlingEnabled bool

	DebugSelectedElementId uint32
	Generation             uint32
	ArenaResetOffset       uintptr

	MeasureTextUserData       interface{}
	QueryScrollOffsetUserData interface{}

	InternalArena Clay_Arena

	// Layout Elements / Render Commands
	LayoutElements              Clay__Array[Clay_LayoutElement]
	RenderCommands              Clay__Array[Clay_RenderCommand]
	OpenLayoutElementStack      Clay__Array[int32]
	LayoutElementChildren       Clay__Array[int32]
	LayoutElementChildrenBuffer Clay__Array[int32]
	TextElementData             Clay__Array[Clay__TextElementData]
	AspectRatioElementIndexes   Clay__Array[int32]
	ReusableElementIndexBuffer  Clay__Array[int32]
	LayoutElementClipElementIds Clay__Array[int32]

	// Configs
	LayoutConfigs             Clay__Array[Clay_LayoutConfig]
	ElementConfigs            Clay__Array[Clay_ElementConfig]
	TextElementConfigs        Clay__Array[Clay_TextElementConfig]
	AspectRatioElementConfigs Clay__Array[Clay_AspectRatioElementConfig]
	ImageElementConfigs       Clay__Array[Clay_ImageElementConfig]
	FloatingElementConfigs    Clay__Array[Clay_FloatingElementConfig]
	ClipElementConfigs        Clay__Array[Clay_ClipElementConfig]
	CustomElementConfigs      Clay__Array[Clay_CustomElementConfig]
	BorderElementConfigs      Clay__Array[Clay_BorderElementConfig]
	SharedElementConfigs      Clay__Array[Clay_SharedElementConfig]

	// Misc Data Structures
	LayoutElementIdStrings             Clay__Array[Clay_String]
	WrappedTextLines                   Clay__Array[Clay__WrappedTextLine]
	LayoutElementTreeNodeArray1        Clay__Array[Clay_LayoutElementTreeNode]
	LayoutElementTreeRoots             Clay__Array[Clay_LayoutElementTreeRoot]
	LayoutElementsHashMapInternal      Clay__Array[Clay_LayoutElementHashMapItem]
	LayoutElementsHashMap              Clay__Array[int32]
	MeasureTextHashMapInternal         Clay__Array[int32]
	MeasureTextHashMapInternalFreeList Clay__Array[int32]
	MeasureTextHashMap                 Clay__Array[int32]
	MeasuredWords                      Clay__Array[Clay__MeasuredWord]
	MeasuredWordsFreeList              Clay__Array[int32]
	OpenClipElementStack               Clay__Array[int32]
	PointerOverIds                     Clay__Array[Clay_ElementId]
	ScrollContainerDatas               Clay__Array[Clay__ScrollContainerDataInternal]
	TreeNodeVisited                    Clay__Array[bool]
	DynamicStringData                  Clay__Array[byte] // char
	DebugElementData                   Clay__Array[Clay__DebugElementData]
}

type Clay_Padding struct {
	Left   uint16
	Right  uint16
	Top    uint16
	Bottom uint16
}
type Clay_ChildAlignment struct {
	TopLeft     float32
	TopRight    float32
	BottomLeft  float32
	BottomRight float32
}
type Clay_LayoutDirection uint8

const (
	// (Default) Lays out child elements from left to right with increasing x.
	CLAY_LEFT_TO_RIGHT Clay_LayoutDirection = iota
	// Lays out child elements from top to bottom with increasing y.
	CLAY_TOP_TO_BOTTOM
)

// Controls various settings that affect the size and position of an element, as well as the sizes and positions
// of any child elements.
type Clay_LayoutConfig struct {
	Sizing          Clay_Sizing
	Padding         Clay_Padding
	ChildGap        uint16
	ChildAlignment  Clay_ChildAlignment
	LayoutDirection Clay_LayoutDirection
}

type Clay_Color struct {
	R float32 // range between 0 and 1
	G float32 // range between 0 and 1
	B float32 // range between 0 and 1
	A float32 // range between 0 and 1
}
type Clay_CornerRadius struct {
	TopLeft     float32
	TopRight    float32
	BottomLeft  float32
	BottomRight float32
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

// Controls how mouse pointer events like hover and click are captured or passed through to elements underneath a floating element.
type Clay_PointerCaptureMode uint8

const (
	// (default) "Capture" the pointer event and don't allow events like hover and click to pass through to elements underneath.

	CLAY_POINTER_CAPTURE_MODE_CAPTURE Clay_PointerCaptureMode = iota
	//    CLAY_POINTER_CAPTURE_MODE_PARENT, TODO pass pointer through to attached parent

	// Transparently pass through pointer events like hover and click to elements underneath the floating element.

	CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH
)

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

// Controls settings related to element borders.
type Clay_BorderElementConfig struct {
	Color Clay_Color       // Controls the color of all borders with width > 0. Conventionally represented as 0-255, but interpretation is up to the renderer.
	Width Clay_BorderWidth // Controls the widths of individual borders. At least one of these should be > 0 for a BORDER render command to be generated.

}

type Clay_BorderWidth struct {
	Left   uint16
	Right  uint16
	Top    uint16
	Bottom uint16
	// Creates borders between each child element, depending on the .layoutDirection.
	// e.g. for LEFT_TO_RIGHT, borders will be vertical lines, and for TOP_TO_BOTTOM borders will be horizontal lines.
	// .betweenChildren borders will result in individual RECTANGLE render commands being generated.
	BetweenChildren uint16
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
		newContext.LayoutElementsHashMap.InternalArray[i] = -1
	}
	for i := int32(0); i < newContext.MeasureTextHashMap.Capacity; i++ {
		newContext.MeasureTextHashMap.InternalArray[i] = 0
	}

	return newContext
}

func Clay_SetCurrentContext(context *Clay_Context) {
	Clay__currentContext = context
}

func Clay_GetCurrentContext() *Clay_Context {
	return Clay__currentContext
}

func Clay_BeginLayout() {
}

func Clay_EndLayout() Clay__Array[Clay_RenderCommand] {
	return Clay__Array[Clay_RenderCommand]{
		Capacity:      Clay__defaultMaxElementCount,
		Length:        0,
		InternalArray: make([]Clay_RenderCommand, Clay__defaultMaxElementCount),
	}
}

var CLAY__DEFAULT_STRUCT = Clay_LayoutElement{}

func Clay__OpenElementWithId(elementId Clay_ElementId) {
	context := Clay_GetCurrentContext()
	if context.LayoutElements.Length == context.LayoutElements.Capacity-1 || context.BooleanWarnings.MaxElementsExceeded {
		context.BooleanWarnings.MaxElementsExceeded = true
		return
	}
	layoutElement := CLAY__DEFAULT_STRUCT
	layoutElement.Id = elementId
	// Clay_LayoutElement * openLayoutElement = Clay_LayoutElementArray_Add(&context->layoutElements, layoutElement);
	// Clay__int32_tArray_Add(&context->openLayoutElementStack, context->layoutElements.length - 1);
	// Clay__AddHashMapItem(elementId, openLayoutElement);
	// Clay__StringArray_Add(&context->layoutElementIdStrings, elementId.stringId);
	// if (context->openClipElementStack.length > 0) {
	//     Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, Clay__int32_tArray_GetValue(&context->openClipElementStack, (int)context->openClipElementStack.length - 1));
	// } else {
	//     Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, 0);
	// }
}

func Clay__HashString(key Clay_String) Clay_ElementId {
	h := fnv.New32a()
	h.Write([]byte(key))
	return h.Sum32()
}

func CLAY_STRING(label string) Clay_String {
	return label
}

func CLAY(elementID Clay_ElementId, elementDeclaration Clay_ElementDeclaration) {
	Clay__OpenElementWithId(elementID)
	Clay__ConfigureOpenElement(elementDeclaration)
}

func Clay__ConfigureOpenElement(elementDeclaration Clay_ElementDeclaration) {
	// Clay_Context* context = Clay_GetCurrentContext();
	// Clay_LayoutElement *openLayoutElement = Clay__GetOpenLayoutElement();
	// openLayoutElement->layoutConfig = Clay__StoreLayoutConfig(declaration->layout);
	// if ((declaration->layout.sizing.width.type == CLAY__SIZING_TYPE_PERCENT && declaration->layout.sizing.width.size.percent > 1) || (declaration->layout.sizing.height.type == CLAY__SIZING_TYPE_PERCENT && declaration->layout.sizing.height.size.percent > 1)) {
	//     context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
	//             .errorType = CLAY_ERROR_TYPE_PERCENTAGE_OVER_1,
	//             .errorText = CLAY_STRING("An element was configured with CLAY_SIZING_PERCENT, but the provided percentage value was over 1.0. Clay expects a value between 0 and 1, i.e. 20% is 0.2."),
	//             .userData = context->errorHandler.userData });
	// }

	// openLayoutElement->elementConfigs.internalArray = &context->elementConfigs.internalArray[context->elementConfigs.length];
	// Clay_SharedElementConfig *sharedConfig = NULL;
	// if (declaration->backgroundColor.a > 0) {
	//     sharedConfig = Clay__StoreSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .backgroundColor = declaration->backgroundColor });
	//     Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
	// }
	// if (!Clay__MemCmp((char *)(&declaration->cornerRadius), (char *)(&Clay__CornerRadius_DEFAULT), sizeof(Clay_CornerRadius))) {
	//     if (sharedConfig) {
	//         sharedConfig->cornerRadius = declaration->cornerRadius;
	//     } else {
	//         sharedConfig = Clay__StoreSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .cornerRadius = declaration->cornerRadius });
	//         Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
	//     }
	// }
	// if (declaration->userData != 0) {
	//     if (sharedConfig) {
	//         sharedConfig->userData = declaration->userData;
	//     } else {
	//         sharedConfig = Clay__StoreSharedElementConfig(CLAY__INIT(Clay_SharedElementConfig) { .userData = declaration->userData });
	//         Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .sharedElementConfig = sharedConfig }, CLAY__ELEMENT_CONFIG_TYPE_SHARED);
	//     }
	// }
	// if (declaration->image.imageData) {
	//     Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .imageElementConfig = Clay__StoreImageElementConfig(declaration->image) }, CLAY__ELEMENT_CONFIG_TYPE_IMAGE);
	// }
	// if (declaration->aspectRatio.aspectRatio > 0) {
	//     Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .aspectRatioElementConfig = Clay__StoreAspectRatioElementConfig(declaration->aspectRatio) }, CLAY__ELEMENT_CONFIG_TYPE_ASPECT);
	//     Clay__int32_tArray_Add(&context->aspectRatioElementIndexes, context->layoutElements.length - 1);
	// }
	// if (declaration->floating.attachTo != CLAY_ATTACH_TO_NONE) {
	//     Clay_FloatingElementConfig floatingConfig = declaration->floating;
	//     // This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here
	//     Clay_LayoutElement *hierarchicalParent = Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 2));
	//     if (hierarchicalParent) {
	//         uint32_t clipElementId = 0;
	//         if (declaration->floating.attachTo == CLAY_ATTACH_TO_PARENT) {
	//             // Attach to the element's direct hierarchical parent
	//             floatingConfig.parentId = hierarchicalParent->id;
	//             if (context->openClipElementStack.length > 0) {
	//                 clipElementId = Clay__int32_tArray_GetValue(&context->openClipElementStack, (int)context->openClipElementStack.length - 1);
	//             }
	//         } else if (declaration->floating.attachTo == CLAY_ATTACH_TO_ELEMENT_WITH_ID) {
	//             Clay_LayoutElementHashMapItem *parentItem = Clay__GetHashMapItem(floatingConfig.parentId);
	//             if (parentItem == &Clay_LayoutElementHashMapItem_DEFAULT) {
	//                 context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
	//                         .errorType = CLAY_ERROR_TYPE_FLOATING_CONTAINER_PARENT_NOT_FOUND,
	//                         .errorText = CLAY_STRING("A floating element was declared with a parentId, but no element with that ID was found."),
	//                         .userData = context->errorHandler.userData });
	//             } else {
	//                 clipElementId = Clay__int32_tArray_GetValue(&context->layoutElementClipElementIds, (int32_t)(parentItem->layoutElement - context->layoutElements.internalArray));
	//             }
	//         } else if (declaration->floating.attachTo == CLAY_ATTACH_TO_ROOT) {
	//             floatingConfig.parentId = Clay__HashString(CLAY_STRING("Clay__RootContainer"), 0).id;
	//         }
	//         if (declaration->floating.clipTo == CLAY_CLIP_TO_NONE) {
	//             clipElementId = 0;
	//         }
	//         int32_t currentElementIndex = Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 1);
	//         Clay__int32_tArray_Set(&context->layoutElementClipElementIds, currentElementIndex, clipElementId);
	//         Clay__int32_tArray_Add(&context->openClipElementStack, clipElementId);
	//         Clay__LayoutElementTreeRootArray_Add(&context->layoutElementTreeRoots, CLAY__INIT(Clay__LayoutElementTreeRoot) {
	//                 .layoutElementIndex = Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 1),
	//                 .parentId = floatingConfig.parentId,
	//                 .clipElementId = clipElementId,
	//                 .zIndex = floatingConfig.zIndex,
	//         });
	//         Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .floatingElementConfig = Clay__StoreFloatingElementConfig(floatingConfig) }, CLAY__ELEMENT_CONFIG_TYPE_FLOATING);
	//     }
	// }
	// if (declaration->custom.customData) {
	//     Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .customElementConfig = Clay__StoreCustomElementConfig(declaration->custom) }, CLAY__ELEMENT_CONFIG_TYPE_CUSTOM);
	// }

	// if (declaration->clip.horizontal | declaration->clip.vertical) {
	//     Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .clipElementConfig = Clay__StoreClipElementConfig(declaration->clip) }, CLAY__ELEMENT_CONFIG_TYPE_CLIP);
	//     Clay__int32_tArray_Add(&context->openClipElementStack, (int)openLayoutElement->id);
	//     // Retrieve or create cached data to track scroll position across frames
	//     Clay__ScrollContainerDataInternal *scrollOffset = CLAY__NULL;
	//     for (int32_t i = 0; i < context->scrollContainerDatas.length; i++) {
	//         Clay__ScrollContainerDataInternal *mapping = Clay__ScrollContainerDataInternalArray_Get(&context->scrollContainerDatas, i);
	//         if (openLayoutElement->id == mapping->elementId) {
	//             scrollOffset = mapping;
	//             scrollOffset->layoutElement = openLayoutElement;
	//             scrollOffset->openThisFrame = true;
	//         }
	//     }
	//     if (!scrollOffset) {
	//         scrollOffset = Clay__ScrollContainerDataInternalArray_Add(&context->scrollContainerDatas, CLAY__INIT(Clay__ScrollContainerDataInternal){.layoutElement = openLayoutElement, .scrollOrigin = {-1,-1}, .elementId = openLayoutElement->id, .openThisFrame = true});
	//     }
	//     if (context->externalScrollHandlingEnabled) {
	//         scrollOffset->scrollPosition = Clay__QueryScrollOffset(scrollOffset->elementId, context->queryScrollOffsetUserData);
	//     }
	// }
	// if (!Clay__MemCmp((char *)(&declaration->border.width), (char *)(&Clay__BorderWidth_DEFAULT), sizeof(Clay_BorderWidth))) {
	//     Clay__AttachElementConfig(CLAY__INIT(Clay_ElementConfigUnion) { .borderElementConfig = Clay__StoreBorderElementConfig(declaration->border) }, CLAY__ELEMENT_CONFIG_TYPE_BORDER);
	// }
}

func Clay_SetLayoutDimensions(dimensions Clay_Dimensions) {
	currentContext := Clay_GetCurrentContext()
	currentContext.LayoutDimensions = dimensions
}

// #define CLAY_TEXT(text, textConfig) Clay__OpenTextElement(text, textConfig)
func CLAY_TEXT(text Clay_String, textConfig Clay_TextElementConfig) {
	Clay__OpenTextElement(text, textConfig)
}

func Clay__OpenTextElement(text Clay_String, textConfig Clay_TextElementConfig) {
	// Clay_Context* context = Clay_GetCurrentContext();
	// if (context->layoutElements.length == context->layoutElements.capacity - 1 || context->booleanWarnings.maxElementsExceeded) {
	//     context->booleanWarnings.maxElementsExceeded = true;
	//     return;
	// }
	// Clay_LayoutElement *parentElement = Clay__GetOpenLayoutElement();

	// Clay_LayoutElement layoutElement = CLAY__DEFAULT_STRUCT;
	// Clay_LayoutElement *textElement = Clay_LayoutElementArray_Add(&context->layoutElements, layoutElement);
	// if (context->openClipElementStack.length > 0) {
	//     Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, Clay__int32_tArray_GetValue(&context->openClipElementStack, (int)context->openClipElementStack.length - 1));
	// } else {
	//     Clay__int32_tArray_Set(&context->layoutElementClipElementIds, context->layoutElements.length - 1, 0);
	// }

	// Clay__int32_tArray_Add(&context->layoutElementChildrenBuffer, context->layoutElements.length - 1);
	// Clay__MeasureTextCacheItem *textMeasured = Clay__MeasureTextCached(&text, textConfig);
	// Clay_ElementId elementId = Clay__HashNumber(parentElement->childrenOrTextContent.children.length, parentElement->id);
	// textElement->id = elementId.id;
	// Clay__AddHashMapItem(elementId, textElement);
	// Clay__StringArray_Add(&context->layoutElementIdStrings, elementId.stringId);
	// Clay_Dimensions textDimensions = { .width = textMeasured->unwrappedDimensions.width, .height = textConfig->lineHeight > 0 ? (float)textConfig->lineHeight : textMeasured->unwrappedDimensions.height };
	// textElement->dimensions = textDimensions;
	// textElement->minDimensions = CLAY__INIT(Clay_Dimensions) { .width = textMeasured->minWidth, .height = textDimensions.height };
	// textElement->childrenOrTextContent.textElementData = Clay__TextElementDataArray_Add(&context->textElementData, CLAY__INIT(Clay__TextElementData) { .text = text, .preferredDimensions = textMeasured->unwrappedDimensions, .elementIndex = context->layoutElements.length - 1 });
	// textElement->elementConfigs = CLAY__INIT(Clay__ElementConfigArraySlice) {
	//         .length = 1,
	//         .internalArray = Clay__ElementConfigArray_Add(&context->elementConfigs, CLAY__INIT(Clay_ElementConfig) { .type = CLAY__ELEMENT_CONFIG_TYPE_TEXT, .config = { .textElementConfig = textConfig }})
	// };
	// textElement->layoutConfig = &CLAY_LAYOUT_DEFAULT;
	// parentElement->childrenOrTextContent.children.length++;
}

func Clay__InitializePersistentMemory(context *Clay_Context) {
	// Persistent memory - initialized once and not reset
	// maxElementCount := context.MaxElementCount;
	// maxMeasureTextCacheWordCount := context.MaxMeasureTextCacheWordCount;
	// arena = &context.internalArena;

	// context->scrollContainerDatas = Clay__ScrollContainerDataInternalArray_Allocate_Arena(100, arena);
	// context->layoutElementsHashMapInternal = Clay__LayoutElementHashMapItemArray_Allocate_Arena(maxElementCount, arena);
	// context->layoutElementsHashMap = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->measureTextHashMapInternal = Clay__MeasureTextCacheItemArray_Allocate_Arena(maxElementCount, arena);
	// context->measureTextHashMapInternalFreeList = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->measuredWordsFreeList = Clay__int32_tArray_Allocate_Arena(maxMeasureTextCacheWordCount, arena);
	// context->measureTextHashMap = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->measuredWords = Clay__MeasuredWordArray_Allocate_Arena(maxMeasureTextCacheWordCount, arena);
	// context->pointerOverIds = Clay_ElementIdArray_Allocate_Arena(maxElementCount, arena);
	// context->debugElementData = Clay__DebugElementDataArray_Allocate_Arena(maxElementCount, arena);
	// context->arenaResetOffset = arena->nextAllocation;
}

func Clay__InitializeEphemeralMemory(context *Clay_Context) {
	// int32_t maxElementCount = context->maxElementCount;
	// // Ephemeral Memory - reset every frame
	// Clay_Arena *arena = &context->internalArena;
	// arena->nextAllocation = context->arenaResetOffset;

	// context->layoutElementChildrenBuffer = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->layoutElements = Clay_LayoutElementArray_Allocate_Arena(maxElementCount, arena);
	// context->warnings = Clay__WarningArray_Allocate_Arena(100, arena);

	// context->layoutConfigs = Clay__LayoutConfigArray_Allocate_Arena(maxElementCount, arena);
	// context->elementConfigs = Clay__ElementConfigArray_Allocate_Arena(maxElementCount, arena);
	// context->textElementConfigs = Clay__TextElementConfigArray_Allocate_Arena(maxElementCount, arena);
	// context->aspectRatioElementConfigs = Clay__AspectRatioElementConfigArray_Allocate_Arena(maxElementCount, arena);
	// context->imageElementConfigs = Clay__ImageElementConfigArray_Allocate_Arena(maxElementCount, arena);
	// context->floatingElementConfigs = Clay__FloatingElementConfigArray_Allocate_Arena(maxElementCount, arena);
	// context->clipElementConfigs = Clay__ClipElementConfigArray_Allocate_Arena(maxElementCount, arena);
	// context->customElementConfigs = Clay__CustomElementConfigArray_Allocate_Arena(maxElementCount, arena);
	// context->borderElementConfigs = Clay__BorderElementConfigArray_Allocate_Arena(maxElementCount, arena);
	// context->sharedElementConfigs = Clay__SharedElementConfigArray_Allocate_Arena(maxElementCount, arena);

	// context->layoutElementIdStrings = Clay__StringArray_Allocate_Arena(maxElementCount, arena);
	// context->wrappedTextLines = Clay__WrappedTextLineArray_Allocate_Arena(maxElementCount, arena);
	// context->layoutElementTreeNodeArray1 = Clay__LayoutElementTreeNodeArray_Allocate_Arena(maxElementCount, arena);
	// context->layoutElementTreeRoots = Clay__LayoutElementTreeRootArray_Allocate_Arena(maxElementCount, arena);
	// context->layoutElementChildren = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->openLayoutElementStack = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->textElementData = Clay__TextElementDataArray_Allocate_Arena(maxElementCount, arena);
	// context->aspectRatioElementIndexes = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->renderCommands = Clay_RenderCommandArray_Allocate_Arena(maxElementCount, arena);
	// context->treeNodeVisited = Clay__boolArray_Allocate_Arena(maxElementCount, arena);
	// context->treeNodeVisited.length = context->treeNodeVisited.capacity; // This array is accessed directly rather than behaving as a list
	// context->openClipElementStack = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->reusableElementIndexBuffer = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->layoutElementClipElementIds = Clay__int32_tArray_Allocate_Arena(maxElementCount, arena);
	// context->dynamicStringData = Clay__charArray_Allocate_Arena(maxElementCount, arena);
}

func Clay__Context_Allocate_Arena(arena *Clay_Arena) *Clay_Context {
	clay_Context, err := mem.AllocateStruct[Clay_Context](arena)
	if err != nil {
		return nil
	}
	return clay_Context
}
