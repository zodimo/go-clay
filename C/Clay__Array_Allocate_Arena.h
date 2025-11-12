


void* Clay__Array_Allocate_Arena(int32_t capacity, uint32_t itemSize, Clay_Arena *arena)
{
    size_t totalSizeBytes = capacity * itemSize;
    uintptr_t nextAllocOffset = arena->nextAllocation + ((64 - (arena->nextAllocation % 64)) & 63);
    if (nextAllocOffset + totalSizeBytes <= arena->capacity) {
        arena->nextAllocation = nextAllocOffset + totalSizeBytes;
        return (void*)((uintptr_t)arena->memory + (uintptr_t)nextAllocOffset);
    }
    else {
        Clay__currentContext->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                .errorType = CLAY_ERROR_TYPE_ARENA_CAPACITY_EXCEEDED,
                .errorText = CLAY_STRING("Clay attempted to allocate memory in its arena, but ran out of capacity. Try increasing the capacity of the arena passed to Clay_Initialize()"),
                .userData = Clay__currentContext->errorHandler.userData });
    }
    return CLAY__NULL;
}
