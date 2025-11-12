package clay

import (
	"github.com/zodimo/go-arena-memory/mem"
)

// CLAY_DLL_EXPORT Clay_Context* Clay_Initialize(Clay_Arena arena, Clay_Dimensions layoutDimensions, Clay_ErrorHandler errorHandler);

type Clay_Arena = mem.Arena

type Clay__QueryScrollOffsetFunction func(elementId uint32, userData interface{}) Clay_Vector2

// Primarily created via the CLAY_ID(), CLAY_IDI(), CLAY_ID_LOCAL() and CLAY_IDI_LOCAL() macros.
// Represents a hashed string ID used for identifying and finding specific clay UI elements, required
// by functions such as Clay_PointerOver() and Clay_GetElementData().
type Clay_ElementId struct {
	Id       uint32      // The resulting hash generated from the other fields.
	Offset   uint32      // A numerical offset applied after computing the hash from stringId.
	BaseId   uint32      // A base hash value to start from, for example the parent element ID is used when calculating CLAY_ID_LOCAL().
	StringId Clay_String // The string id to hash.
}

// Note: Clay_String is not guaranteed to be null terminated. It may be if created from a literal C string,
// but it is also used to represent slices.
type Clay_String struct {
	// Set this boolean to true if the char* data underlying this string will live for the entire lifetime of the program.
	// This will automatically be set for strings created with CLAY_STRING, as the macro requires a string literal.
	IsStaticallyAllocated bool
	Length                int32
	// The underlying character memory. Note: this will not be copied and will not extend the lifetime of the underlying memory.
	Chars string
}

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

type Clay_ArraySlice[T any] struct {
	Length        int32
	InternalArray []T
}

func Clay__Array_Allocate_Arena[T any](capacity int32, arena *Clay_Arena) *Clay__Array[T] {
	outArr, err := mem.AllocateStructObject[Clay__Array[T]](arena, NewClay__Array[T](capacity))
	if err != nil {
		panic(err)
	}
	return outArr
}

type Clay__LayoutElementChildren struct {
	Elements []int32
	Length   uint16
}

type Clay__LayoutElementChildrenOrTextContent struct {
	Children        Clay__LayoutElementChildren
	TextElementData *Clay__TextElementData
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
	StringContents Clay_StringSlice
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

type Clay__LayoutElementTreeRoot struct {
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
	LayoutElement         *Clay_LayoutElement
	OnHoverFunction       func(elementId Clay_ElementId, pointerInfo Clay_PointerData, userData any)
	HoverFunctionUserData any
	NextIndex             int32
	Generation            uint32
	DebugData             *Clay__DebugElementData
}

type Clay__MeasuredWord struct {
	StartOffset int32
	Length      int32
	Width       float32
	Next        int32
}

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

type Clay_Context struct {
	MaxElementCount              int32
	MaxMeasureTextCacheWordCount int32
	WarningsEnabled              bool
	ErrorHandler                 Clay_ErrorHandler

	BooleanWarnings Clay_BooleanWarnings
	Warnings        *Clay__Array[Clay__Warning]

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
	LayoutElements              *Clay__Array[Clay_LayoutElement]
	RenderCommands              *Clay__Array[Clay_RenderCommand]
	OpenLayoutElementStack      *Clay__Array[int32]
	LayoutElementChildren       *Clay__Array[int32]
	LayoutElementChildrenBuffer *Clay__Array[int32]
	TextElementData             *Clay__Array[Clay__TextElementData]
	AspectRatioElementIndexes   *Clay__Array[int32]
	ReusableElementIndexBuffer  *Clay__Array[int32]
	LayoutElementClipElementIds *Clay__Array[int32]

	// Configs*
	LayoutConfigs             *Clay__Array[Clay_LayoutConfig]
	ElementConfigs            *Clay__Array[Clay_ElementConfig]
	TextElementConfigs        *Clay__Array[Clay_TextElementConfig]
	AspectRatioElementConfigs *Clay__Array[Clay_AspectRatioElementConfig]
	ImageElementConfigs       *Clay__Array[Clay_ImageElementConfig]
	FloatingElementConfigs    *Clay__Array[Clay_FloatingElementConfig]
	ClipElementConfigs        *Clay__Array[Clay_ClipElementConfig]
	CustomElementConfigs      *Clay__Array[Clay_CustomElementConfig]
	BorderElementConfigs      *Clay__Array[Clay_BorderElementConfig]
	SharedElementConfigs      *Clay__Array[Clay_SharedElementConfig]

	// Misc Data Structures
	LayoutElementIdStrings             *Clay__Array[Clay_String]
	WrappedTextLines                   *Clay__Array[Clay__WrappedTextLine]
	LayoutElementTreeNodeArray1        *Clay__Array[Clay_LayoutElementTreeNode]
	LayoutElementTreeRoots             *Clay__Array[Clay__LayoutElementTreeRoot]
	LayoutElementsHashMapInternal      *Clay__Array[Clay_LayoutElementHashMapItem]
	LayoutElementsHashMap              *Clay__Array[int32]
	MeasureTextHashMapInternal         *Clay__Array[Clay__MeasureTextCacheItem]
	MeasureTextHashMapInternalFreeList *Clay__Array[int32]
	MeasureTextHashMap                 *Clay__Array[int32]
	MeasuredWords                      *Clay__Array[Clay__MeasuredWord]
	MeasuredWordsFreeList              *Clay__Array[int32]
	OpenClipElementStack               *Clay__Array[int32]
	PointerOverIds                     *Clay__Array[Clay_ElementId]
	ScrollContainerDatas               *Clay__Array[Clay__ScrollContainerDataInternal]
	TreeNodeVisited                    *Clay__Array[bool]
	DynamicStringData                  *Clay__Array[byte] // char
	DebugElementData                   *Clay__Array[Clay__DebugElementData]
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

	Clay__Array_Add(currentContext.OpenLayoutElementStack, 0)
	Clay__Array_Add(currentContext.LayoutElementTreeRoots, Clay__LayoutElementTreeRoot{
		LayoutElementIndex: 0,
	})

}

func CLAY__MAX(x, y float32) float32 {
	if x > y {
		return x
	}
	return y
}
func CLAY__MIN(x, y float32) float32 {
	if x < y {
		return x
	}
	return y
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

	//this may not be right
	openLayoutElement.ChildrenOrTextContent.Children.Elements = currentContext.LayoutElementChildren.InternalArray //[currentContext.LayoutElementChildren.Length]
	if layoutConfig.LayoutDirection == CLAY_LEFT_TO_RIGHT {
		openLayoutElement.Dimensions.Width = leftRightPadding
		openLayoutElement.MinDimensions.Width = leftRightPadding
		for i := uint16(0); i < openLayoutElement.ChildrenOrTextContent.Children.Length; i++ {
			childIndex := Clay__Array_GetValue(currentContext.LayoutElementChildrenBuffer, currentContext.LayoutElementChildrenBuffer.Length-int32(openLayoutElement.ChildrenOrTextContent.Children.Length)+int32(i))
			child := Clay__Array_Get(currentContext.LayoutElements, childIndex)
			openLayoutElement.Dimensions.Width += child.Dimensions.Width
			openLayoutElement.Dimensions.Height = CLAY__MAX(openLayoutElement.Dimensions.Height, child.Dimensions.Height+topBottomPadding)

			// Minimum size of child elements doesn't matter to clip containers as they can shrink and hide their contents
			if !elementHasClipHorizontal {
				openLayoutElement.MinDimensions.Width += child.MinDimensions.Width
			}
			if !elementHasClipVertical {
				openLayoutElement.MinDimensions.Height = CLAY__MAX(openLayoutElement.MinDimensions.Height, child.MinDimensions.Height+topBottomPadding)
			}
			Clay__Array_Add(currentContext.LayoutElementChildren, childIndex)
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
			childIndex := Clay__Array_GetValue(currentContext.LayoutElementChildrenBuffer, currentContext.LayoutElementChildrenBuffer.Length-int32(openLayoutElement.ChildrenOrTextContent.Children.Length)+int32(i))
			child := Clay__Array_Get(currentContext.LayoutElements, childIndex)
			openLayoutElement.Dimensions.Height += child.Dimensions.Height
			openLayoutElement.Dimensions.Width = CLAY__MAX(openLayoutElement.Dimensions.Width, child.Dimensions.Width+leftRightPadding)
			if !elementHasClipVertical {
				openLayoutElement.MinDimensions.Height += child.MinDimensions.Height
			}
			if !elementHasClipHorizontal {
				openLayoutElement.MinDimensions.Width = CLAY__MAX(openLayoutElement.MinDimensions.Width, child.MinDimensions.Width+leftRightPadding)
			}
			Clay__Array_Add(currentContext.LayoutElementChildren, childIndex)
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
	closingElementIndex := Clay__Array_RemoveSwapback(currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-1)

	// Get the currently open parent
	openLayoutElement = Clay__GetOpenLayoutElement()

	if currentContext.OpenLayoutElementStack.Length > 1 {
		if elementIsFloating {
			openLayoutElement.FloatingChildrenCount++
			return
		}
		openLayoutElement.ChildrenOrTextContent.Children.Length++
		Clay__Array_Add(currentContext.LayoutElementChildrenBuffer, closingElementIndex)
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
	// Clay__CalculateFinalLayout();
	return currentContext.RenderCommands.InternalArray[:currentContext.RenderCommands.Length]

}

func Clay__AddRenderCommand(renderCommand Clay_RenderCommand) {
	currentContext := Clay_GetCurrentContext()
	if currentContext.RenderCommands.Length < currentContext.RenderCommands.Capacity-1 {
		Clay__Array_Add(currentContext.RenderCommands, renderCommand)
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
	openLayoutElement := Clay__Array_Add(currentContext.LayoutElements, layoutElement)
	Clay__Array_Add(currentContext.OpenLayoutElementStack, currentContext.LayoutElements.Length-1)
	Clay__AddHashMapItem(elementId, openLayoutElement)
	Clay__Array_Add(currentContext.LayoutElementIdStrings, elementId.StringId)
	if currentContext.OpenClipElementStack.Length > 0 {
		Clay__Array_Set(currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length-1, Clay__Array_GetValue(currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length-1))
	} else {
		Clay__Array_Set(currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length-1, 0)
	}
}

func CLAY_STRING(label string) Clay_String {
	return Clay_String{
		IsStaticallyAllocated: true,
		Length:                int32(len(label)),
		Chars:                 label,
	}
}

func CLAY(elementID Clay_ElementId, elementDeclaration Clay_ElementDeclaration) {
	Clay__OpenElementWithId(elementID)
	Clay__ConfigureOpenElement(elementDeclaration)
}

func Clay__StoreLayoutConfig(config Clay_LayoutConfig) *Clay_LayoutConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_LayoutConfig{}
	}
	return Clay__Array_Add(currentContext.LayoutConfigs, config)

}

func Clay__AttachElementConfig(config Clay_ElementConfigUnion, configType Clay__ElementConfigType) Clay_ElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return Clay_ElementConfig{}
	}
	openLayoutElement := Clay__GetOpenLayoutElement()
	openLayoutElement.ElementConfigs.Length++
	return *Clay__Array_Add(currentContext.ElementConfigs, Clay_ElementConfig{Type: configType, Config: config})
}

func Clay__StoreSharedElementConfig(config Clay_SharedElementConfig) *Clay_SharedElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_SharedElementConfig{}
	}
	return Clay__Array_Add(currentContext.SharedElementConfigs, config)
}

func Clay__StoreImageElementConfig(config Clay_ImageElementConfig) *Clay_ImageElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_ImageElementConfig{}
	}
	return Clay__Array_Add(currentContext.ImageElementConfigs, config)
}

func Clay__StoreAspectRatioElementConfig(config Clay_AspectRatioElementConfig) *Clay_AspectRatioElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_AspectRatioElementConfig{}
	}
	return Clay__Array_Add(currentContext.AspectRatioElementConfigs, config)
}
func Clay__StoreFloatingElementConfig(config Clay_FloatingElementConfig) *Clay_FloatingElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_FloatingElementConfig{}
	}
	return Clay__Array_Add(currentContext.FloatingElementConfigs, config)
}

func Clay__StoreCustomElementConfig(config Clay_CustomElementConfig) *Clay_CustomElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_CustomElementConfig{}
	}
	return Clay__Array_Add(currentContext.CustomElementConfigs, config)
}

func Clay__StoreClipElementConfig(config Clay_ClipElementConfig) *Clay_ClipElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_ClipElementConfig{}
	}
	return Clay__Array_Add(currentContext.ClipElementConfigs, config)
}

func Clay__StoreBorderElementConfig(config Clay_BorderElementConfig) *Clay_BorderElementConfig {
	currentContext := Clay_GetCurrentContext()
	if currentContext.BooleanWarnings.MaxElementsExceeded {
		return &Clay_BorderElementConfig{}
	}
	return Clay__Array_Add(currentContext.BorderElementConfigs, config)
}

func Clay__GetHashMapItem(id uint32) *Clay_LayoutElementHashMapItem {
	currentContext := Clay_GetCurrentContext()
	// Perform modulo with uint32 first to avoid negative results, then cast to int32
	hashBucket := int32(id % uint32(currentContext.LayoutElementsHashMap.Capacity))

	elementIndex := currentContext.LayoutElementsHashMap.InternalArray[hashBucket]
	for elementIndex != -1 {
		hashEntry := Clay__Array_Get(currentContext.LayoutElementsHashMapInternal, elementIndex)
		if hashEntry.ElementId.Id == id {
			return hashEntry
		}
		elementIndex = hashEntry.NextIndex
	}

	return &Clay_LayoutElementHashMapItem{}
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
	elementConfigsSegmentView := currentContext.ElementConfigs.InternalArray[nextAvailableElementConfigIndex : nextAvailableElementConfigIndex+1]
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
		Clay__Array_Add(currentContext.AspectRatioElementIndexes, currentContext.LayoutElements.Length-1)
	}

	if elementDeclaration.Floating.AttachTo != CLAY_ATTACH_TO_NONE {
		floatingConfig := elementDeclaration.Floating
		// This looks dodgy but because of the auto generated root element the depth of the tree will always be at least 2 here

		hierarchicalParent := Clay__Array_Get[Clay_LayoutElement](currentContext.LayoutElements, Clay__Array_GetValue[int32](currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-2))
		if hierarchicalParent != nil {
			var clipElementId int32 = 0
			if elementDeclaration.Floating.AttachTo == CLAY_ATTACH_TO_PARENT {
				// Attach to the element's direct hierarchical parent
				floatingConfig.ParentId = hierarchicalParent.Id
				if currentContext.OpenClipElementStack.Length > 0 {
					clipElementId = Clay__Array_GetValue(currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length-1)
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
						for i, elem := range currentContext.LayoutElements.InternalArray {
							if &elem == parentItem.LayoutElement {
								clipItemIndex = int32(i)
								break
							}
						}
						if clipItemIndex != -1 {
							clipElementId = Clay__Array_GetValue[int32](currentContext.LayoutElementClipElementIds, clipItemIndex)
						}
					}
				} else if elementDeclaration.Floating.AttachTo == CLAY_ATTACH_TO_ROOT {
					floatingConfig.ParentId = Clay__HashString(CLAY_STRING("Clay__RootContainer"), 0).Id
				}

				if elementDeclaration.Floating.ClipTo == CLAY_CLIP_TO_NONE {
					clipElementId = 0
				}
				currentElementIndex := Clay__Array_GetValue[int32](currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-1)
				Clay__Array_Set(currentContext.LayoutElementClipElementIds, currentElementIndex, clipElementId)
				Clay__Array_Add(currentContext.OpenClipElementStack, clipElementId)
				Clay__Array_Add(currentContext.LayoutElementTreeRoots, Clay__LayoutElementTreeRoot{
					LayoutElementIndex: Clay__Array_GetValue[int32](currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-1),
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

		Clay__Array_Add(currentContext.OpenClipElementStack, int32(openLayoutElement.Id))
		// Retrieve or create cached data to track scroll position across frames
		var scrollOffset *Clay__ScrollContainerDataInternal = nil
		for i := int32(0); i < currentContext.ScrollContainerDatas.Length; i++ {
			mapping := Clay__Array_Get[Clay__ScrollContainerDataInternal](currentContext.ScrollContainerDatas, i)
			if openLayoutElement.Id == mapping.ElementId {
				scrollOffset = mapping
				scrollOffset.LayoutElement = openLayoutElement
				scrollOffset.OpenThisFrame = true
			}
		}
		if scrollOffset == nil {
			scrollOffset = Clay__Array_Add(currentContext.ScrollContainerDatas, Clay__ScrollContainerDataInternal{
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

func Clay_SetLayoutDimensions(dimensions Clay_Dimensions) {
	currentContext := Clay_GetCurrentContext()
	currentContext.LayoutDimensions = dimensions
}

// #define CLAY_TEXT(text, textConfig) Clay__OpenTextElement(text, textConfig)
func CLAY_TEXT(text Clay_String, textConfig Clay_TextElementConfig) {
	Clay__OpenTextElement(text, textConfig)
}

func Clay__GetOpenLayoutElement() *Clay_LayoutElement {
	currentContext := Clay_GetCurrentContext()
	return Clay__Array_Get[Clay_LayoutElement](currentContext.LayoutElements, Clay__Array_GetValue[int32](currentContext.OpenLayoutElementStack, currentContext.OpenLayoutElementStack.Length-1))

	// Clay_LayoutElement* Clay__GetOpenLayoutElement(void) {
	//     Clay_Context* context = Clay_GetCurrentContext();
	//     return Clay_LayoutElementArray_Get(&context->layoutElements, Clay__int32_tArray_GetValue(&context->openLayoutElementStack, context->openLayoutElementStack.length - 1));
	// }

}
func Clay__MeasureTextCached(text *Clay_String, textConfig Clay_TextElementConfig) *Clay__MeasureTextCacheItem {
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
	hashItemIndex := currentContext.LayoutElementsHashMap.InternalArray[hashBucket]
	for hashItemIndex != -1 { // Just replace collision, not a big deal - leave it up to the end user
		hashItem := Clay__Array_Get[Clay_LayoutElementHashMapItem](currentContext.LayoutElementsHashMapInternal, hashItemIndex)
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

	hashItem := Clay__Array_Add(currentContext.LayoutElementsHashMapInternal, item)
	hashItem.DebugData = Clay__Array_Add(currentContext.DebugElementData, Clay__DebugElementData{})
	if hashItemPrevious != -1 {
		Clay__Array_Get[Clay_LayoutElementHashMapItem](currentContext.LayoutElementsHashMapInternal, hashItemPrevious).NextIndex = currentContext.LayoutElementsHashMapInternal.Length - 1
	} else {
		currentContext.LayoutElementsHashMap.InternalArray[hashBucket] = currentContext.LayoutElementsHashMapInternal.Length - 1
	}
	return hashItem
}

func Clay__OpenTextElement(text Clay_String, textConfig Clay_TextElementConfig) {
	currentContext := Clay_GetCurrentContext()
	if currentContext.LayoutElements.Length == currentContext.LayoutElements.Capacity-1 || currentContext.BooleanWarnings.MaxElementsExceeded {
		currentContext.BooleanWarnings.MaxElementsExceeded = true
		return
	}
	parentElement := Clay__GetOpenLayoutElement()

	layoutElement := Clay_LayoutElement{}

	textElement := Clay__Array_Add[Clay_LayoutElement](currentContext.LayoutElements, layoutElement)

	if currentContext.OpenClipElementStack.Length > 0 {
		Clay__Array_Set(currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length-1, Clay__Array_GetValue[int32](currentContext.OpenClipElementStack, currentContext.OpenClipElementStack.Length-1))
	} else {
		Clay__Array_Set(currentContext.LayoutElementClipElementIds, currentContext.LayoutElements.Length-1, 0)
	}

	Clay__Array_Add(currentContext.LayoutElementChildrenBuffer, currentContext.LayoutElements.Length-1)

	textMeasured := Clay__MeasureTextCached(&text, textConfig)

	elementId := Clay__HashNumber(uint32(parentElement.ChildrenOrTextContent.Children.Length), parentElement.Id)

	textElement.Id = elementId.Id

	Clay__AddHashMapItem(elementId, textElement)
	Clay__Array_Add(currentContext.LayoutElementIdStrings, elementId.StringId)

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

	textElement.ChildrenOrTextContent.TextElementData = Clay__Array_Add(currentContext.TextElementData, Clay__TextElementData{
		Text:                text,
		PreferredDimensions: textMeasured.UnwrappedDimensions,
		ElementIndex:        currentContext.LayoutElements.Length - 1,
	})

	// add config to element configs

	config := Clay__Array_Add(currentContext.ElementConfigs, Clay_ElementConfig{
		Type:   CLAY__ELEMENT_CONFIG_TYPE_TEXT,
		Config: Clay_ElementConfigUnion{TextElementConfig: &textConfig},
	})
	if config != nil {
		configIndex := currentContext.ElementConfigs.Length - 1

		segmentView := currentContext.ElementConfigs.InternalArray[configIndex : configIndex+1]
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
	context.LayoutElementTreeNodeArray1 = Clay__Array_Allocate_Arena[Clay_LayoutElementTreeNode](maxElementCount, arena)
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
