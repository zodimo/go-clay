# Renderer Guide

Go-Clay is designed to be render engine agnostic. This guide shows you how to create custom renderers and integrate with existing rendering systems.

## Clay Rendering Architecture

Clay uses a command-based rendering system where the layout engine produces an array of render commands that renderers process sequentially. This approach provides flexibility and performance benefits.

### Core Concepts

1. **Layout Phase**: Clay computes layout and generates a `Clay_RenderCommandArray`
2. **Render Phase**: Renderers process the command array, handling each command based on its type
3. **Command Types**: Different rendering operations (rectangle, text, image, border, clipping, custom)

## Renderer Interface

All renderers should implement a function that processes a `Clay_RenderCommandArray`:

```c
// C API Pattern
void MyRenderer_Render(Clay_RenderCommandArray renderCommands, MyRendererFonts* fonts) {
    for (int i = 0; i < renderCommands.length; i++) {
        Clay_RenderCommand *renderCommand = Clay_RenderCommandArray_Get(&renderCommands, i);
        Clay_BoundingBox boundingBox = renderCommand->boundingBox;
        
        switch (renderCommand->commandType) {
            case CLAY_RENDER_COMMAND_TYPE_RECTANGLE: {
                // Handle rectangle rendering
                break;
            }
            case CLAY_RENDER_COMMAND_TYPE_TEXT: {
                // Handle text rendering
                break;
            }
            // ... other command types
        }
    }
}
```

For Go-Clay, renderers implement the `Renderer` interface:

```go
type Renderer interface {
    // Unified rendering method that processes command arrays
    Render(commands []RenderCommand) error
    
    // Frame lifecycle management
    BeginFrame() error
    EndFrame() error
    SetViewport(bounds BoundingBox) error
}
```

## Render Command Structure

### Clay_RenderCommand (C API)

```c
typedef struct Clay_RenderCommand {
    // Bounding box for this UI element, relative to layout root
    Clay_BoundingBox boundingBox;
    
    // Command-specific data (union)
    Clay_RenderData renderData;
    
    // User data pointer from element declaration
    void *userData;
    
    // Element ID from element declaration
    uint32_t id;
    
    // Z-order for correct rendering (array is pre-sorted)
    int16_t zIndex;
    
    // Command type determining how to render
    Clay_RenderCommandType commandType;
} Clay_RenderCommand;
```

### RenderCommand (Go API)

```go
type RenderCommand struct {
    BoundingBox BoundingBox
    CommandType RenderCommandType
    ZIndex      int16
    ID          ElementID
    Data        interface{} // Command-specific data
}
```

## Command Types

### CLAY_RENDER_COMMAND_TYPE_RECTANGLE
Renders a solid color rectangle with optional corner radius.

**C Render Data:**
```c
typedef struct Clay_RectangleRenderData {
    Clay_Color backgroundColor;
    Clay_CornerRadius cornerRadius;
} Clay_RectangleRenderData;
```

**Go Render Data:**
```go
type RectangleCommand struct {
    Color        Color
    CornerRadius CornerRadius
}
```

### CLAY_RENDER_COMMAND_TYPE_TEXT
Renders text with specified font, size, and styling.

**C Render Data:**
```c
typedef struct Clay_TextRenderData {
    Clay_StringSlice stringContents;  // Not null-terminated
    Clay_Color textColor;
    uint16_t fontId;
    uint16_t fontSize;
    uint16_t letterSpacing;
    uint16_t lineHeight;
} Clay_TextRenderData;
```

**Go Render Data:**
```go
type TextCommand struct {
    Text          string
    FontID        uint16
    FontSize      float32
    Color         Color
    LineHeight    float32
    LetterSpacing float32
    Alignment     TextAlignment
}
```

### CLAY_RENDER_COMMAND_TYPE_IMAGE
Renders an image with optional tinting and corner radius.

**C Render Data:**
```c
typedef struct Clay_ImageRenderData {
    Clay_Color backgroundColor;  // Tint color (0,0,0,0 = no tint)
    Clay_CornerRadius cornerRadius;
    void* imageData;  // Renderer-specific image data
} Clay_ImageRenderData;
```

**Go Render Data:**
```go
type ImageCommand struct {
    ImageData    interface{}
    TintColor    Color
    CornerRadius CornerRadius
}
```

### CLAY_RENDER_COMMAND_TYPE_BORDER
Renders a border inset into the bounding box.

**C Render Data:**
```c
typedef struct Clay_BorderRenderData {
    Clay_Color color;
    Clay_CornerRadius cornerRadius;
    Clay_BorderWidth width;  // Individual side widths
} Clay_BorderRenderData;
```

**Go Render Data:**
```go
type BorderCommand struct {
    Color        Color
    Width        BorderWidth
    CornerRadius CornerRadius
}
```

### CLAY_RENDER_COMMAND_TYPE_SCISSOR_START
Begins clipping region - only render content within bounding box.

**C Render Data:**
```c
typedef struct Clay_ScrollRenderData {
    bool horizontal;
    bool vertical;
} Clay_ClipRenderData;
```

**Go Render Data:**
```go
type ClipStartCommand struct {
    Horizontal, Vertical bool
}
```

### CLAY_RENDER_COMMAND_TYPE_SCISSOR_END
Ends current clipping region.

**C Render Data:** None required

**Go Render Data:**
```go
type ClipEndCommand struct {
    // No additional data needed
}
```

### CLAY_RENDER_COMMAND_TYPE_CUSTOM
Custom renderer-specific operations.

**C Render Data:**
```c
typedef struct Clay_CustomRenderData {
    Clay_Color backgroundColor;
    Clay_CornerRadius cornerRadius;
    void* customData;  // Renderer-specific data
} Clay_CustomRenderData;
```

**Go Render Data:**
```go
type CustomCommand struct {
    CustomData interface{}
}
```

## Creating a Custom Renderer

### C Renderer Example

```c
#include "clay.h"

typedef struct MyRenderer {
    // Renderer state
    int width, height;
    void* renderTarget;
} MyRenderer;

void MyRenderer_Render(Clay_RenderCommandArray renderCommands, MyRenderer* renderer) {
    for (int i = 0; i < renderCommands.length; i++) {
        Clay_RenderCommand *cmd = Clay_RenderCommandArray_Get(&renderCommands, i);
        Clay_BoundingBox bounds = cmd->boundingBox;
        
        switch (cmd->commandType) {
            case CLAY_RENDER_COMMAND_TYPE_RECTANGLE: {
                Clay_RectangleRenderData *data = &cmd->renderData.rectangle;
                // Render rectangle using data->backgroundColor, data->cornerRadius
                MyRenderer_DrawRectangle(renderer, bounds, data->backgroundColor, data->cornerRadius);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_TEXT: {
                Clay_TextRenderData *data = &cmd->renderData.text;
                // Render text using data->stringContents, data->textColor, etc.
                MyRenderer_DrawText(renderer, bounds, data);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_IMAGE: {
                Clay_ImageRenderData *data = &cmd->renderData.image;
                // Render image using data->imageData, data->backgroundColor (tint)
                MyRenderer_DrawImage(renderer, bounds, data);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_BORDER: {
                Clay_BorderRenderData *data = &cmd->renderData.border;
                // Render border using data->color, data->width, data->cornerRadius
                MyRenderer_DrawBorder(renderer, bounds, data);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_SCISSOR_START: {
                // Begin clipping to bounds
                MyRenderer_BeginClip(renderer, bounds);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_SCISSOR_END: {
                // End clipping
                MyRenderer_EndClip(renderer);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_CUSTOM: {
                Clay_CustomRenderData *data = &cmd->renderData.custom;
                // Handle custom rendering based on data->customData
                MyRenderer_HandleCustom(renderer, bounds, data);
                break;
            }
        }
    }
}
```

### Go Renderer Example

```go
package main

import (
    "github.com/zodimo/go-clay"
)

type MyRenderer struct {
    // Your renderer state
    width, height int
    renderTarget  interface{}
}

func NewMyRenderer() *MyRenderer {
    return &MyRenderer{}
}

func (r *MyRenderer) BeginFrame() error {
    // Initialize frame rendering
    return nil
}

func (r *MyRenderer) EndFrame() error {
    // Finalize frame rendering
    return nil
}

func (r *MyRenderer) SetViewport(bounds clay.BoundingBox) error {
    // Set rendering viewport
    r.width = int(bounds.Width)
    r.height = int(bounds.Height)
    return nil
}

func (r *MyRenderer) Render(commands []clay.RenderCommand) error {
    for _, cmd := range commands {
        bounds := cmd.BoundingBox
        
        switch cmd.CommandType {
        case clay.CommandRectangle:
            if rectCmd, ok := cmd.Data.(clay.RectangleCommand); ok {
                // Render rectangle using rectCmd.Color, rectCmd.CornerRadius
                r.drawRectangle(bounds, rectCmd)
            }
            
        case clay.CommandText:
            if textCmd, ok := cmd.Data.(clay.TextCommand); ok {
                // Render text using textCmd fields
                r.drawText(bounds, textCmd)
            }
            
        case clay.CommandImage:
            if imageCmd, ok := cmd.Data.(clay.ImageCommand); ok {
                // Render image using imageCmd fields
                r.drawImage(bounds, imageCmd)
            }
            
        case clay.CommandBorder:
            if borderCmd, ok := cmd.Data.(clay.BorderCommand); ok {
                // Render border using borderCmd fields
                r.drawBorder(bounds, borderCmd)
            }
            
        case clay.CommandClipStart:
            if clipCmd, ok := cmd.Data.(clay.ClipStartCommand); ok {
                // Begin clipping
                r.beginClip(bounds, clipCmd)
            }
            
        case clay.CommandClipEnd:
            // End clipping
            r.endClip()
            
        case clay.CommandCustom:
            if customCmd, ok := cmd.Data.(clay.CustomCommand); ok {
                // Handle custom rendering
                r.handleCustom(bounds, customCmd)
            }
        }
    }
    return nil
}

// Helper methods
func (r *MyRenderer) drawRectangle(bounds clay.BoundingBox, cmd clay.RectangleCommand) {
    // Implementation specific to your renderer
}

func (r *MyRenderer) drawText(bounds clay.BoundingBox, cmd clay.TextCommand) {
    // Implementation specific to your renderer
}

// ... other helper methods
```

## Integration Examples

### OpenGL Renderer

```c
#include <GL/gl.h>
#include "clay.h"

typedef struct OpenGLRenderer {
    GLuint shaderProgram;
    GLuint vao, vbo;
} OpenGLRenderer;

void OpenGL_Clay_Render(Clay_RenderCommandArray renderCommands, OpenGLRenderer* renderer) {
    glUseProgram(renderer->shaderProgram);
    glBindVertexArray(renderer->vao);
    
    for (int i = 0; i < renderCommands.length; i++) {
        Clay_RenderCommand *cmd = Clay_RenderCommandArray_Get(&renderCommands, i);
        
        switch (cmd->commandType) {
            case CLAY_RENDER_COMMAND_TYPE_RECTANGLE: {
                Clay_RectangleRenderData *data = &cmd->renderData.rectangle;
                
                // Set up rectangle vertices
                float vertices[] = {
                    cmd->boundingBox.x, cmd->boundingBox.y,
                    cmd->boundingBox.x + cmd->boundingBox.width, cmd->boundingBox.y,
                    cmd->boundingBox.x, cmd->boundingBox.y + cmd->boundingBox.height,
                    cmd->boundingBox.x + cmd->boundingBox.width, cmd->boundingBox.y + cmd->boundingBox.height,
                };
                
                // Upload vertices
                glBindBuffer(GL_ARRAY_BUFFER, renderer->vbo);
                glBufferData(GL_ARRAY_BUFFER, sizeof(vertices), vertices, GL_DYNAMIC_DRAW);
                
                // Set color uniform
                glUniform4f(glGetUniformLocation(renderer->shaderProgram, "color"),
                    data->backgroundColor.r / 255.0f,
                    data->backgroundColor.g / 255.0f,
                    data->backgroundColor.b / 255.0f,
                    data->backgroundColor.a / 255.0f);
                
                // Draw
                glDrawArrays(GL_TRIANGLE_STRIP, 0, 4);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_SCISSOR_START: {
                glEnable(GL_SCISSOR_TEST);
                glScissor((int)cmd->boundingBox.x, (int)cmd->boundingBox.y,
                         (int)cmd->boundingBox.width, (int)cmd->boundingBox.height);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_SCISSOR_END: {
                glDisable(GL_SCISSOR_TEST);
                break;
            }
            
            // Handle other command types...
        }
    }
}
```

### Raylib Integration

```c
#include "raylib.h"
#include "clay.h"

void Clay_Raylib_Render(Clay_RenderCommandArray renderCommands, Font* fonts) {
    for (int i = 0; i < renderCommands.length; i++) {
        Clay_RenderCommand *cmd = Clay_RenderCommandArray_Get(&renderCommands, i);
        Clay_BoundingBox bounds = {
            roundf(cmd->boundingBox.x), roundf(cmd->boundingBox.y),
            roundf(cmd->boundingBox.width), roundf(cmd->boundingBox.height)
        };
        
        switch (cmd->commandType) {
            case CLAY_RENDER_COMMAND_TYPE_RECTANGLE: {
                Clay_RectangleRenderData *data = &cmd->renderData.rectangle;
                Color color = {
                    data->backgroundColor.r, data->backgroundColor.g,
                    data->backgroundColor.b, data->backgroundColor.a
                };
                
                if (data->cornerRadius.topLeft > 0) {
                    // Draw rounded rectangle (implementation needed)
                    DrawRectangleRounded((Rectangle){bounds.x, bounds.y, bounds.width, bounds.height},
                                       data->cornerRadius.topLeft / bounds.width, 0, color);
                } else {
                    DrawRectangle(bounds.x, bounds.y, bounds.width, bounds.height, color);
                }
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_TEXT: {
                Clay_TextRenderData *data = &cmd->renderData.text;
                Color color = {data->textColor.r, data->textColor.g, data->textColor.b, data->textColor.a};
                
                // Convert Clay_StringSlice to null-terminated string
                char* text = malloc(data->stringContents.length + 1);
                memcpy(text, data->stringContents.chars, data->stringContents.length);
                text[data->stringContents.length] = '\0';
                
                DrawTextEx(fonts[data->fontId], text, (Vector2){bounds.x, bounds.y}, 
                          data->fontSize, data->letterSpacing, color);
                
                free(text);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_SCISSOR_START: {
                BeginScissorMode(bounds.x, bounds.y, bounds.width, bounds.height);
                break;
            }
            
            case CLAY_RENDER_COMMAND_TYPE_SCISSOR_END: {
                EndScissorMode();
                break;
            }
            
            // Handle other command types...
        }
    }
}
```

## Performance Considerations

### Command Batching
Group similar commands for better performance:

```c
void OptimizedRenderer_Render(Clay_RenderCommandArray renderCommands) {
    // First pass: collect commands by type
    RectangleCommand rectangles[1000];
    TextCommand texts[1000];
    int rectCount = 0, textCount = 0;
    
    for (int i = 0; i < renderCommands.length; i++) {
        Clay_RenderCommand *cmd = Clay_RenderCommandArray_Get(&renderCommands, i);
        
        switch (cmd->commandType) {
            case CLAY_RENDER_COMMAND_TYPE_RECTANGLE:
                if (rectCount < 1000) {
                    rectangles[rectCount++] = (RectangleCommand){
                        .bounds = cmd->boundingBox,
                        .data = cmd->renderData.rectangle
                    };
                }
                break;
            case CLAY_RENDER_COMMAND_TYPE_TEXT:
                if (textCount < 1000) {
                    texts[textCount++] = (TextCommand){
                        .bounds = cmd->boundingBox,
                        .data = cmd->renderData.text
                    };
                }
                break;
        }
    }
    
    // Second pass: batch render by type
    BatchRenderRectangles(rectangles, rectCount);
    BatchRenderTexts(texts, textCount);
}
```

### Memory Management
Use object pools to reduce allocations:

```go
type RendererPool struct {
    rectanglePool sync.Pool
    textPool      sync.Pool
}

func (p *RendererPool) getRectangleData() *RectangleRenderData {
    if v := p.rectanglePool.Get(); v != nil {
        return v.(*RectangleRenderData)
    }
    return &RectangleRenderData{}
}

func (p *RendererPool) putRectangleData(data *RectangleRenderData) {
    // Reset data
    *data = RectangleRenderData{}
    p.rectanglePool.Put(data)
}
```

## Testing Renderers

### Unit Testing

```go
func TestMyRenderer(t *testing.T) {
    renderer := NewMyRenderer()
    
    commands := []clay.RenderCommand{
        {
            BoundingBox: clay.BoundingBox{X: 10, Y: 10, Width: 100, Height: 50},
            CommandType: clay.CommandRectangle,
            Data: clay.RectangleCommand{
                Color: clay.Color{R: 255, G: 0, B: 0, A: 255},
            },
        },
        {
            BoundingBox: clay.BoundingBox{X: 10, Y: 70, Width: 200, Height: 30},
            CommandType: clay.CommandText,
            Data: clay.TextCommand{
                Text:     "Hello World",
                FontSize: 16,
                Color:    clay.Color{R: 0, G: 0, B: 0, A: 255},
            },
        },
    }
    
    err := renderer.Render(commands)
    assert.NoError(t, err)
}
```

### Integration Testing

```c
void TestClayRenderer() {
    // Set up Clay context
    Clay_Initialize(Clay_CreateArenaWithCapacityAndMemory(1024 * 1024, malloc(1024 * 1024)));
    
    // Create layout
    Clay_BeginLayout();
    CLAY(CLAY_RECTANGLE({.color = {255, 0, 0, 255}})) {
        CLAY_TEXT(CLAY_STRING("Test"), CLAY_TEXT_CONFIG({.fontSize = 16}));
    }
    Clay_RenderCommandArray commands = Clay_EndLayout();
    
    // Test renderer
    MyRenderer renderer = {0};
    MyRenderer_Render(commands, &renderer);
    
    // Verify results
    assert(renderer.drawCallCount == 2); // Rectangle + Text
}
```

## Best Practices

1. **Process commands in order** - Clay pre-sorts commands by z-index
2. **Handle all command types** - Even if some are no-ops for your renderer
3. **Respect clipping regions** - Implement SCISSOR_START/END properly
4. **Optimize for your platform** - Batch similar operations when possible
5. **Handle errors gracefully** - Return meaningful error messages
6. **Test thoroughly** - Verify all command types and edge cases
7. **Document limitations** - Note any renderer-specific constraints

## Legacy Interface Support

Go-Clay also supports a legacy interface for backward compatibility:

```go
type LegacyRenderer interface {
    RenderRectangle(cmd RectangleCommand) error
    RenderText(cmd TextCommand) error
    RenderImage(cmd ImageCommand) error
    RenderBorder(cmd BorderCommand) error
    RenderClipStart(cmd ClipStartCommand) error
    RenderClipEnd(cmd ClipEndCommand) error
    RenderCustom(cmd CustomCommand) error
    
    BeginFrame() error
    EndFrame() error
    SetViewport(bounds BoundingBox) error
}
```

However, the unified `Render(commands []RenderCommand)` approach is preferred as it matches the Clay C API pattern and provides better performance opportunities.