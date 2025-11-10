



typedef struct Clay_ElementDeclaration {
    // Controls various settings that affect the size and position of an element, as well as the sizes and positions of any child elements.
    Clay_LayoutConfig layout;
    // Controls the background color of the resulting element.
    // By convention specified as 0-255, but interpretation is up to the renderer.
    // If no other config is specified, .backgroundColor will generate a RECTANGLE render command, otherwise it will be passed as a property to IMAGE or CUSTOM render commands.
    Clay_Color backgroundColor;
    // Controls the "radius", or corner rounding of elements, including rectangles, borders and images.
    Clay_CornerRadius cornerRadius;
    // Controls settings related to aspect ratio scaling.
    Clay_AspectRatioElementConfig aspectRatio;
    // Controls settings related to image elements.
    Clay_ImageElementConfig image;
    // Controls whether and how an element "floats", which means it layers over the top of other elements in z order, and doesn't affect the position and size of siblings or parent elements.
    // Note: in order to activate floating, .floating.attachTo must be set to something other than the default value.
    Clay_FloatingElementConfig floating;
    // Used to create CUSTOM render commands, usually to render element types not supported by Clay.
    Clay_CustomElementConfig custom;
    // Controls whether an element should clip its contents, as well as providing child x,y offset configuration for scrolling.
    Clay_ClipElementConfig clip;
    // Controls settings related to element borders, and will generate BORDER render commands.
    Clay_BorderElementConfig border;
    // A pointer that will be transparently passed through to resulting render commands.
    void *userData;
} Clay_ElementDeclaration;
