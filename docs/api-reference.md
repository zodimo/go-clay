# API Reference

## Core Types

### ElementID
```go
type ElementID uint32
```
Unique identifier for UI elements. Used for element lookup and interaction.

### Dimensions
```go
type Dimensions struct {
    Width, Height float32
}
```
Represents size in pixels.

### Vector2
```go
type Vector2 struct {
    X, Y float32
}
```
2D position or offset.

### Color
```go
type Color struct {
    R, G, B, A float32
}
```
RGBA color with values 0.0-1.0.

### BoundingBox
```go
type BoundingBox struct {
    X, Y, Width, Height float32
}
```
Rectangle bounds for positioning and clipping.

## Layout System

### SizingType
```go
type SizingType int

const (
    SizingFit     SizingType = iota // Wrap to content
    SizingGrow                      // Fill available space
    SizingPercent                   // Percentage of parent
    SizingFixed                     // Fixed pixel size
)
```

### SizingAxis
```go
type SizingAxis struct {
    Type     SizingType
    Min, Max float32  // Min/max constraints
    Percent  float32  // Percentage (0.0-1.0)
}
```

### Sizing
```go
type Sizing struct {
    Width, Height SizingAxis
}
```

### LayoutDirection
```go
type LayoutDirection int

const (
    LeftToRight LayoutDirection = iota
    TopToBottom
)
```

### Alignment
```go
type AlignmentX int
const (
    AlignXLeft   AlignmentX = iota
    AlignXCenter
    AlignXRight
)

type AlignmentY int
const (
    AlignYTop    AlignmentY = iota
    AlignYCenter
    AlignYBottom
)

type ChildAlignment struct {
    X, Y AlignmentX, AlignmentY
}
```

### Padding
```go
type Padding struct {
    Left, Right, Top, Bottom float32
}

func PaddingAll(value float32) Padding
func PaddingHorizontal(left, right float32) Padding
func PaddingVertical(top, bottom float32) Padding
```

### LayoutConfig
```go
type LayoutConfig struct {
    Sizing         Sizing
    Padding        Padding
    ChildGap       float32
    Direction      LayoutDirection
    ChildAlignment ChildAlignment
}
```

## Element System

### ElementDeclaration
```go
type ElementDeclaration struct {
    ID              ElementID
    Layout          LayoutConfig
    BackgroundColor Color
    Border          *BorderConfig
    CornerRadius    CornerRadius
    // ... other styling properties
}
```

### TextConfig
```go
type TextConfig struct {
    FontSize     float32
    Color        Color
    FontID       uint16
    LineHeight   float32
    LetterSpacing float32
    WrapMode     TextWrapMode
    Alignment    TextAlignment
}
```

### ImageConfig
```go
type ImageConfig struct {
    ImageData interface{}
    TintColor Color
}
```

## Layout Engine

### LayoutEngine
```go
type LayoutEngine interface {
    // Layout lifecycle
    BeginLayout()
    EndLayout() []RenderCommand
    
    // Element management
    OpenElement(id ElementID, config ElementDeclaration)
    CloseElement()
    
    // State management
    SetPointerState(pos Vector2, pressed bool)
    SetLayoutDimensions(dimensions Dimensions)
    SetScrollOffset(offset Vector2)
    
    // Element queries
    GetElementBounds(id ElementID) (BoundingBox, bool)
    IsPointerOver(id ElementID) bool
    GetScrollOffset(id ElementID) Vector2
}
```

### NewLayoutEngine
```go
func NewLayoutEngine() LayoutEngine
```
Creates a new layout engine instance.

## Render Commands

### RenderCommandType
```go
type RenderCommandType int

const (
    CommandRectangle RenderCommandType = iota
    CommandText
    CommandImage
    CommandBorder
    CommandClipStart
    CommandClipEnd
    CommandCustom
)
```

### RenderCommand
```go
type RenderCommand struct {
    BoundingBox BoundingBox
    CommandType RenderCommandType
    ZIndex      int16
    ID          ElementID
    Data        interface{} // Command-specific data
}
```

## Builder API

### Container
```go
func Container(id ElementID, config ElementDeclaration) *ContainerBuilder
```

### ContainerBuilder
```go
type ContainerBuilder struct {
    // ... internal fields
}

func (c *ContainerBuilder) Text(text string, config TextConfig) *ContainerBuilder
func (c *ContainerBuilder) Image(config ImageConfig) *ContainerBuilder
func (c *ContainerBuilder) Container(id ElementID, config ElementDeclaration) *ContainerBuilder
func (c *ContainerBuilder) End()
```

## Utility Functions

### Sizing Helpers
```go
func SizingFit() SizingAxis
func SizingGrow(weight float32) SizingAxis
func SizingPercent(percent float32) SizingAxis
func SizingFixed(size float32) SizingAxis
```

### Color Helpers
```go
func ColorRGB(r, g, b float32) Color
func ColorRGBA(r, g, b, a float32) Color
func ColorHex(hex string) Color
```

### ID Helpers
```go
func ID(str string) ElementID
func IDWithIndex(str string, index int) ElementID
```

## Error Handling

### LayoutError
```go
type LayoutError struct {
    Type    ErrorType
    Message string
    Element ElementID
}

type ErrorType int

const (
    ErrorInvalidSizing ErrorType = iota
    ErrorCircularReference
    ErrorInvalidElement
    ErrorMemoryExhausted
)
```

## Performance

### Memory Management
```go
type Arena struct {
    // ... internal fields
}

func NewArena(capacity int) *Arena
func (a *Arena) Allocate(size int) []byte
func (a *Arena) Reset()
```

### Profiling
```go
type LayoutStats struct {
    ElementCount    int
    RenderCommands  int
    LayoutTime      time.Duration
    MemoryUsed      int
}

func (e *LayoutEngine) GetStats() LayoutStats
```
