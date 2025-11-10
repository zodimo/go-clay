package clay

import "hash/fnv"

// CLAY_DLL_EXPORT Clay_Context* Clay_Initialize(Clay_Arena arena, Clay_Dimensions layoutDimensions, Clay_ErrorHandler errorHandler);

var Clay__currentContext *Clay_Context

type Clay_ElementId = uint32
type Clay_String = string

type Clay_Arena struct{}
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

type Clay_LayoutElement struct{}
type Clay_RenderCommand struct{}

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
type Clay_CornerRadius struct{}
type Clay_AspectRatioElementConfig struct{}
type Clay_ImageElementConfig struct{}
type Clay_FloatingElementConfig struct{}
type Clay_CustomElementConfig struct{}
type Clay_ClipElementConfig struct{}
type Clay_BorderElementConfig struct{}

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

func Clay__OpenElementWithId(id Clay_ElementId) {}

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

func Clay__ConfigureOpenElement(elementDeclaration Clay_ElementDeclaration) {}

func Clay_SetLayoutDimensions(dimensions Clay_Dimensions) {
	currentContext := Clay_GetCurrentContext()
	currentContext.LayoutDimensions = dimensions
}

// #define CLAY_TEXT(text, textConfig) Clay__OpenTextElement(text, textConfig)
func CLAY_TEXT(text Clay_String, textConfig Clay_TextElementConfig) {
	Clay__OpenTextElement(text, textConfig)
}

func Clay__OpenTextElement(text Clay_String, textConfig Clay_TextElementConfig) {}

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
