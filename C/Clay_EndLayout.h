Clay_RenderCommandArray Clay_EndLayout(void) {
    Clay_Context* context = Clay_GetCurrentContext();
    Clay__CloseElement();
    bool elementsExceededBeforeDebugView = context->booleanWarnings.maxElementsExceeded;
    if (context->debugModeEnabled && !elementsExceededBeforeDebugView) {
        context->warningsEnabled = false;
        Clay__RenderDebugView();
        context->warningsEnabled = true;
    }
    if (context->booleanWarnings.maxElementsExceeded) {
        Clay_String message;
        if (!elementsExceededBeforeDebugView) {
            message = CLAY_STRING("Clay Error: Layout elements exceeded Clay__maxElementCount after adding the debug-view to the layout.");
        } else {
            message = CLAY_STRING("Clay Error: Layout elements exceeded Clay__maxElementCount");
        }
        Clay__AddRenderCommand(CLAY__INIT(Clay_RenderCommand ) {
            .boundingBox = { 
                context->layoutDimensions.width / 2 - 59 * 4, 
                context->layoutDimensions.height / 2, 
                0, 
                0 
            },
            .renderData = {
                 .text = { 
                    .stringContents = CLAY__INIT(Clay_StringSlice) { 
                        .length = message.length, 
                        .chars = message.chars, 
                        .baseChars = message.chars 
                    }, 
                    .textColor = {255, 0, 0, 255}, 
                    .fontSize = 16 
                } 
            },
            .commandType = CLAY_RENDER_COMMAND_TYPE_TEXT
        });
    }
    if (context->openLayoutElementStack.length > 1) {
        context->errorHandler.errorHandlerFunction(CLAY__INIT(Clay_ErrorData) {
                .errorType = CLAY_ERROR_TYPE_UNBALANCED_OPEN_CLOSE,
                .errorText = CLAY_STRING("There were still open layout elements when EndLayout was called. This results from an unequal number of calls to Clay__OpenElement and Clay__CloseElement."),
                .userData = context->errorHandler.userData });
    }
    Clay__CalculateFinalLayout();
    return context->renderCommands;
}