package clay

var Clay__currentContext *Clay_Context

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
	LayoutElementTreeRoots             Clay__Array[Clay__LayoutElementTreeRoot]
	LayoutElementsHashMapInternal      Clay__Array[Clay_LayoutElementHashMapItem]
	LayoutElementsHashMap              Clay__Array[int32]
	MeasureTextHashMapInternal         Clay__Array[Clay__MeasureTextCacheItem]
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

func Clay_SetCurrentContext(context *Clay_Context) {
	Clay__currentContext = context
}

func Clay_GetCurrentContext() *Clay_Context {
	return Clay__currentContext
}
