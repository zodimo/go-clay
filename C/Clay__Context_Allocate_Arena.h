

Clay_Context* Clay__Context_Allocate_Arena(Clay_Arena *arena) {
    size_t totalSizeBytes = sizeof(Clay_Context);
    if (totalSizeBytes > arena->capacity)
    {
        return NULL;
    }
    arena->nextAllocation += totalSizeBytes;
    return (Clay_Context*)(arena->memory);
}