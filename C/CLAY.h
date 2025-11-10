






#define CLAY(id, ...)                                                                                                                                               \
    for (                                                                                                                                                           \
        CLAY__ELEMENT_DEFINITION_LATCH = (Clay__OpenElementWithId(id), Clay__ConfigureOpenElement(CLAY__CONFIG_WRAPPER(Clay_ElementDeclaration, __VA_ARGS__)), 0);  \
        CLAY__ELEMENT_DEFINITION_LATCH < 1;                                                                                                                         \
        CLAY__ELEMENT_DEFINITION_LATCH=1, Clay__CloseElement()                                                                                                      \
    )
