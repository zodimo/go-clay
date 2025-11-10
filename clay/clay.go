package clay

import "hash/fnv"

// CLAY_DLL_EXPORT Clay_Context* Clay_Initialize(Clay_Arena arena, Clay_Dimensions layoutDimensions, Clay_ErrorHandler errorHandler);

type Clay_ElementId = uint32
type Clay_String = string

type Clay_Arena struct{}
type Clay_Dimensions struct {
	Width  float32
	Height float32
}
type Clay_ErrorHandler struct{}
type Clay_Context struct{}

type Clay_RenderCommand struct{}

type Clay_RenderCommandArray = []Clay_RenderCommand

type Clay_Padding struct{}
type Clay_ChildAlignment struct{}
type Clay_LayoutDirection struct{}

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
	return &Clay_Context{}
}

func Clay_SetCurrentContext(context *Clay_Context) {
}

func Clay_GetCurrentContext() *Clay_Context {
	return &Clay_Context{}
}

func Clay_BeginLayout() {
}

func Clay_EndLayout() Clay_RenderCommandArray {
	return Clay_RenderCommandArray{}
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
