# Migration Guide: Clay C Architecture Implementation

## Overview

This guide covers migrating the Go Clay implementation from the incorrect individual render methods to the correct Clay C architecture with unified command processing.

## Problem Summary

The Go port incorrectly implemented the renderer interface with individual methods that lack bounds information:

```go
// WRONG: Current Go implementation
type Renderer interface {
    RenderRectangle(cmd RectangleCommand) error  // ❌ No bounds
    RenderText(cmd TextCommand) error            // ❌ No bounds
    // ...
}
```

The correct Clay C architecture uses a single method that receives commands with bounds:

```go
// CORRECT: Clay C architecture
type Renderer interface {
    Render(commands []RenderCommand) error  // ✅ Commands include bounds
    BeginFrame() error
    EndFrame() error
    SetViewport(bounds BoundingBox) error
}
```

## Migration Steps

### Step 1: Update Core Interface

**File**: `clay/renderer.go`

**Before**:
```go
type Renderer interface {
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

**After**:
```go
type Renderer interface {
    // Unified rendering method matching Clay C architecture
    Render(commands []RenderCommand) error
    
    // Lifecycle methods
    BeginFrame() error
    EndFrame() error
    SetViewport(bounds BoundingBox) error
}
```

### Step 2: Update Gio Renderer Implementation

**File**: `renderers/gioui/renderer.go`

**Replace all individual methods with unified Render method**:

```go
func (r *GioRenderer) Render(commands []RenderCommand) error {
    for _, cmd := range commands {
        bounds := cmd.BoundingBox  // ✅ Bounds available for positioning
        
        switch cmd.CommandType {
        case CommandRectangle:
            if err := r.renderRectangle(bounds, cmd.Data.(RectangleCommand)); err != nil {
                return err
            }
            
        case CommandText:
            if err := r.renderText(bounds, cmd.Data.(TextCommand)); err != nil {
                return err
            }
            
        case CommandImage:
            if err := r.renderImage(bounds, cmd.Data.(ImageCommand)); err != nil {
                return err
            }
            
        case CommandBorder:
            if err := r.renderBorder(bounds, cmd.Data.(BorderCommand)); err != nil {
                return err
            }
            
        case CommandClipStart:
            if err := r.renderClipStart(bounds, cmd.Data.(ClipStartCommand)); err != nil {
                return err
            }
            
        case CommandClipEnd:
            if err := r.renderClipEnd(bounds, cmd.Data.(ClipEndCommand)); err != nil {
                return err
            }
            
        case CommandCustom:
            if err := r.renderCustom(bounds, cmd.Data.(CustomCommand)); err != nil {
                return err
            }
        }
    }
    return nil
}
```

### Step 3: Implement Bounds-Based Rendering Methods

**Rectangle Rendering with Bounds**:
```go
func (r *GioRenderer) renderRectangle(bounds BoundingBox, cmd RectangleCommand) error {
    if r.ops == nil {
        return fmt.Errorf("operations context is nil")
    }
    
    // ✅ Use bounds for proper rectangle shape
    rect := clip.Rect{
        Min: f32.Pt(float32(bounds.X), float32(bounds.Y)),
        Max: f32.Pt(float32(bounds.X+bounds.Width), float32(bounds.Y+bounds.Height)),
    }
    
    // Handle corner radius if specified
    var clipOp clip.Stack
    if cmd.CornerRadius.TopLeft > 0 || cmd.CornerRadius.TopRight > 0 ||
       cmd.CornerRadius.BottomLeft > 0 || cmd.CornerRadius.BottomRight > 0 {
        clipOp = r.createRoundedRectClip(bounds, cmd.CornerRadius)
    } else {
        clipOp = rect.Push(r.ops)
    }
    defer clipOp.Pop()
    
    // ✅ Use command data for styling
    gioColor := ClayToGioColor(cmd.Color)
    paint.ColorOp{Color: gioColor}.Add(r.ops)
    paint.PaintOp{}.Add(r.ops)
    
    return nil
}
```

**Text Rendering with Bounds**:
```go
func (r *GioRenderer) renderText(bounds BoundingBox, cmd TextCommand) error {
    if r.ops == nil {
        return fmt.Errorf("operations context is nil")
    }
    
    // ✅ Use bounds for text positioning
    textPos := f32.Pt(float32(bounds.X), float32(bounds.Y))
    
    // Load font based on FontID
    font, err := r.fontManager.GetFont(cmd.FontID)
    if err != nil {
        return fmt.Errorf("failed to load font %d: %v", cmd.FontID, err)
    }
    
    // ✅ Use command data for content and styling
    gioColor := ClayToGioColor(cmd.Color)
    paint.ColorOp{Color: gioColor}.Add(r.ops)
    
    // Position and render text
    op.Offset(textPos).Add(r.ops)
    
    // TODO: Implement actual Gio text rendering
    // This requires proper Gio text operations
    
    return nil
}
```

### Step 4: Update Examples

**Remove Workaround Functions**:

**Before** (with workarounds):
```go
// Render with workarounds
renderer.BeginFrame()
for _, cmd := range commands {
    switch cmd.CommandType {
    case clay.CommandRectangle:
        // ❌ Custom workaround function
        renderRectangleWithBounds(renderer, cmd.Data.(clay.RectangleCommand), cmd.BoundingBox, &ops)
    case clay.CommandText:
        // ❌ Broken method without bounds
        renderer.RenderText(cmd.Data.(clay.TextCommand))
    }
}
renderer.EndFrame()
```

**After** (clean and correct):
```go
// Render with unified method
renderer.BeginFrame()
renderer.Render(commands)  // ✅ Single call, bounds included
renderer.EndFrame()
```

### Step 5: Update Tests

**Update all renderer tests to use unified method**:

```go
func TestGioRenderer_RenderRectangle_WithBounds(t *testing.T) {
    ops := &op.Ops{}
    renderer := NewRenderer(ops)
    
    commands := []clay.RenderCommand{
        {
            BoundingBox: clay.BoundingBox{X: 10, Y: 20, Width: 100, Height: 50},
            CommandType: clay.CommandRectangle,
            Data: clay.RectangleCommand{
                Color: clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
            },
        },
    }
    
    err := renderer.Render(commands)
    require.NoError(t, err)
    
    // Verify actual rectangle operations were added
    // TODO: Add visual validation tests
}
```

## Benefits of Migration

1. **Proper Positioning**: Rectangles and text render at correct positions
2. **Matches Clay C**: Consistent architecture across all Clay implementations
3. **Simplified Usage**: Single method call instead of complex switching logic
4. **No Workarounds**: Examples become clean and simple
5. **Extensible**: Easy to add new command types following the same pattern

## Breaking Changes

- **Renderer Interface**: Complete change from individual methods to unified method
- **Examples**: All examples need updates to remove workarounds
- **Tests**: All renderer tests need updates

## Backward Compatibility

A `LegacyRenderer` interface is provided temporarily for backward compatibility, but all new code should use the unified `Renderer` interface.

## Timeline

1. **Phase 1** (Day 1): Update core interface and Gio renderer structure
2. **Phase 2** (Day 2): Implement bounds-based rendering methods
3. **Phase 3** (Day 3): Update examples and remove workarounds
4. **Phase 4** (Day 4): Update all tests and documentation

## Validation

After migration, verify:
- [ ] Rectangles render as proper shapes at correct positions
- [ ] Text renders at correct positions (when text implementation is complete)
- [ ] Examples work without custom workaround functions
- [ ] All tests pass with new interface
- [ ] Architecture matches Clay C pattern exactly
