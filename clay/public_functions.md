# Public functions in clay.h

```c
CLAY_DLL_EXPORT uint32_t Clay_MinMemorySize(void);
CLAY_DLL_EXPORT Clay_Arena Clay_CreateArenaWithCapacityAndMemory(size_t capacity, void *memory);
CLAY_DLL_EXPORT void Clay_SetPointerState(Clay_Vector2 position, bool pointerDown);
CLAY_DLL_EXPORT Clay_Context* Clay_Initialize(Clay_Arena arena, Clay_Dimensions layoutDimensions, Clay_ErrorHandler errorHandler);
CLAY_DLL_EXPORT Clay_Context* Clay_GetCurrentContext(void);
CLAY_DLL_EXPORT void Clay_SetCurrentContext(Clay_Context* context);
CLAY_DLL_EXPORT void Clay_UpdateScrollContainers(bool enableDragScrolling, Clay_Vector2 scrollDelta, float deltaTime);
CLAY_DLL_EXPORT Clay_Vector2 Clay_GetScrollOffset(void);
CLAY_DLL_EXPORT void Clay_SetLayoutDimensions(Clay_Dimensions dimensions);
CLAY_DLL_EXPORT void Clay_BeginLayout(void);
CLAY_DLL_EXPORT Clay_RenderCommandArray Clay_EndLayout(void);
CLAY_DLL_EXPORT Clay_ElementId Clay_GetElementId(Clay_String idString);
CLAY_DLL_EXPORT Clay_ElementId Clay_GetElementIdWithIndex(Clay_String idString, uint32_t index);
CLAY_DLL_EXPORT Clay_ElementData Clay_GetElementData(Clay_ElementId id);
CLAY_DLL_EXPORT bool Clay_Hovered(void);
CLAY_DLL_EXPORT void Clay_OnHover(void (*onHoverFunction)(Clay_ElementId elementId, Clay_PointerData pointerData, intptr_t userData), intptr_t userData);
CLAY_DLL_EXPORT bool Clay_PointerOver(Clay_ElementId elementId);
CLAY_DLL_EXPORT Clay_ElementIdArray Clay_GetPointerOverIds(void);
CLAY_DLL_EXPORT Clay_ScrollContainerData Clay_GetScrollContainerData(Clay_ElementId id);
CLAY_DLL_EXPORT void Clay_SetMeasureTextFunction(Clay_Dimensions (*measureTextFunction)(Clay_StringSlice text, Clay_TextElementConfig *config, void *userData), void *userData);
CLAY_DLL_EXPORT void Clay_SetQueryScrollOffsetFunction(Clay_Vector2 (*queryScrollOffsetFunction)(uint32_t elementId, void *userData), void *userData);
CLAY_DLL_EXPORT Clay_RenderCommand * Clay_RenderCommandArray_Get(Clay_RenderCommandArray* array, int32_t index);
CLAY_DLL_EXPORT void Clay_SetDebugModeEnabled(bool enabled);
CLAY_DLL_EXPORT bool Clay_IsDebugModeEnabled(void);
CLAY_DLL_EXPORT void Clay_SetCullingEnabled(bool enabled);
CLAY_DLL_EXPORT int32_t Clay_GetMaxElementCount(void);
CLAY_DLL_EXPORT void Clay_SetMaxElementCount(int32_t maxElementCount);
CLAY_DLL_EXPORT int32_t Clay_GetMaxMeasureTextCacheWordCount(void);
CLAY_DLL_EXPORT void Clay_SetMaxMeasureTextCacheWordCount(int32_t maxMeasureTextCacheWordCount);
CLAY_DLL_EXPORT void Clay_ResetMeasureTextCache(void);
CLAY_DLL_EXPORT void Clay__OpenElement(void);
CLAY_DLL_EXPORT void Clay__OpenElementWithId(Clay_ElementId elementId);
CLAY_DLL_EXPORT void Clay__ConfigureOpenElement(const Clay_ElementDeclaration config);
CLAY_DLL_EXPORT void Clay__ConfigureOpenElementPtr(const Clay_ElementDeclaration *config);
CLAY_DLL_EXPORT void Clay__CloseElement(void);
CLAY_DLL_EXPORT Clay_ElementId Clay__HashString(Clay_String key, uint32_t seed);
CLAY_DLL_EXPORT Clay_ElementId Clay__HashStringWithOffset(Clay_String key, uint32_t offset, uint32_t seed);
CLAY_DLL_EXPORT void Clay__OpenTextElement(Clay_String text, Clay_TextElementConfig *textConfig);
CLAY_DLL_EXPORT Clay_TextElementConfig *Clay__StoreTextElementConfig(Clay_TextElementConfig config);
CLAY_DLL_EXPORT uint32_t Clay__GetParentElementId(void);
CLAY_DLL_EXPORT Clay_ElementIdArray Clay_GetPointerOverIds(void) {

```