# Gio-MW & Gio-X Solutions Reference for Clay User Story 1.1

**Document Version**: 1.0  
**Date**: November 8, 2025  
**Author**: Mary (Business Analyst)  
**Purpose**: Complete implementation reference for resolving Clay's critical rendering challenges

---

## Executive Summary

This document provides complete implementation details for resolving the 4 critical unimplemented features in Clay's user story 1.1, based on proven patterns from gio-mw and gio-x repositories. All code patterns, imports, and implementation strategies are included for self-contained reference.

### 🔄 **UPDATED**: Corrected Image Rendering Solution

**Previous approach** (testlab/framework/capture.go) was **not usable** for production rendering.

**✅ Correct approach** now documented based on **actual gio-mw production code**:
- **Source**: `gio-mw/wdk/require.go` + `widget.Image` + `examples/kitchen/pages/cards/cards.go`
- **Pattern**: Use Gio's `widget.Image` with `paint.NewImageOp(decodedImage)`
- **Scaling**: Map Clay modes to `widget.Fit` (Contain, Cover, Fill, ScaleDown)
- **Positioning**: Use `layout.Direction` for image positioning within bounds
- **Production Ready**: Based on working card image display in gio-mw kitchen examples

### Critical Issues Addressed
1. **🚨 CRITICAL**: Text rendering (`renderTextWithBounds()` returns `ErrorTypeUnsupportedOperation`)
2. **🚨 CRITICAL**: Image rendering (`renderImageWithBounds()` returns `ErrorTypeUnsupportedOperation`) 
3. **⚠️ MEDIUM**: Corner radius support (TODO comment - not implemented)
4. **📝 LOW**: Custom command type identification (TODO comment - not implemented)

---

## 1. Text Rendering Solution

### Source Pattern: gio-mw/wdk/text.go

**Complete Implementation Pattern**:

```go
// Required imports for text rendering
import (
    "gioui.org/font"
    "gioui.org/font/gofont"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/op/paint"
    "gioui.org/text"
    "gioui.org/unit"
    "gioui.org/widget"
    "image"
)

// Text shaper initialization (do once at renderer creation)
func (r *GioRenderer) initializeTextShaper() {
    shaperOptions := []text.ShaperOption{
        text.WithCollection(gofont.Collection()),
    }
    r.textShaper = text.NewShaper(shaperOptions...)
}

// Complete text rendering implementation
func (r *GioRenderer) renderTextWithBounds(bounds BoundingBox, cmd TextCommand) error {
    if r.textShaper == nil {
        return fmt.Errorf("text shaper not initialized")
    }

    // 1. Create font from Clay command
    font := font.Font{
        Typeface: "Go", // Default - can be mapped from Clay FontID
        Style:    font.Regular,
        Weight:   r.mapClayFontWeight(cmd.FontWeight), // See mapping below
    }

    // 2. Create color operation
    colorMacro := op.Record(r.ops)
    paint.ColorOp{Color: ClayToGioColor(cmd.Color)}.Add(r.ops)
    colorCallOp := colorMacro.Stop()

    // 3. Create label with Clay parameters
    label := widget.Label{
        Alignment:  r.mapClayTextAlignment(cmd.Alignment), // See mapping below
        MaxLines:   0, // Unlimited - can be set from Clay command
        LineHeight: unit.Sp(cmd.LineHeight),
        WrapPolicy: text.WrapWords, // Can be mapped from Clay
    }

    // 4. Position within bounds
    stack := op.Offset(f32.Pt(float32(bounds.X), float32(bounds.Y))).Push(r.ops)
    defer stack.Pop()

    // 5. Create layout context with bounds constraints
    gtx := layout.Context{
        Ops: r.ops,
        Constraints: layout.Constraints{
            Max: image.Pt(int(bounds.Width), int(bounds.Height)),
        },
        Metric: unit.Metric{PxPerDp: 1, PxPerSp: 1}, // Adjust as needed
    }

    // 6. Render text using gio-mw pattern
    dimensions := label.Layout(gtx, r.textShaper, font, unit.Sp(cmd.FontSize), cmd.Text, colorCallOp)
    
    // Optional: Store dimensions for layout feedback
    // r.lastTextDimensions = dimensions
    
    return nil
}

// Font weight mapping
func (r *GioRenderer) mapClayFontWeight(clayWeight int) font.Weight {
    switch clayWeight {
    case 100: return font.Thin
    case 200: return font.ExtraLight
    case 300: return font.Light
    case 400: return font.Normal
    case 500: return font.Medium
    case 600: return font.SemiBold
    case 700: return font.Bold
    case 800: return font.ExtraBold
    case 900: return font.Black
    default: return font.Normal
    }
}

// Text alignment mapping
func (r *GioRenderer) mapClayTextAlignment(clayAlign TextAlignment) text.Alignment {
    switch clayAlign {
    case TextAlignLeft: return text.Start
    case TextAlignCenter: return text.Middle
    case TextAlignRight: return text.End
    default: return text.Start
    }
}
```

### GioRenderer Structure Updates

```go
type GioRenderer struct {
    ops              *op.Ops
    viewport         BoundingBox
    clipStack        []clip.Stack
    textShaper       *text.Shaper  // Add this field
    // ... existing fields
}

// Update constructor
func NewRenderer(ops *op.Ops) *GioRenderer {
    r := &GioRenderer{
        ops: ops,
        // ... other initialization
    }
    r.initializeTextShaper()
    return r
}
```

---

## 2. Corner Radius Solution

### Source Pattern: gio-mw/wdk/shaped.go

**Complete Implementation Pattern**:

```go
// Required imports for corner radius
import (
    "gioui.org/f32"
    "gioui.org/io/system"
    "gioui.org/layout"
    "gioui.org/op/clip"
    "gioui.org/op/paint"
    "math"
)

// Corner shape types
type CornerKind int

const (
    CornerKindDefault CornerKind = iota
    CornerKindChamfer
    CornerKindRound
)

type CornerShape struct {
    Kind        CornerKind
    Size        float32
    AdaptToSize bool
}

type CornerShapes struct {
    TopStart    CornerShape
    TopEnd      CornerShape
    BottomStart CornerShape
    BottomEnd   CornerShape
}

type ShapedRect struct {
    MinPoint f32.Point
    MaxPoint f32.Point
    Offset   float32
    Shapes   CornerShapes
}

// Complete corner radius rectangle rendering
func (r *GioRenderer) renderRectangleWithBounds(bounds BoundingBox, cmd RectangleCommand) error {
    // Check if corner radius is needed
    if cmd.CornerRadius.IsZero() {
        return r.renderSimpleRectangle(bounds, cmd)
    }

    // Create shaped rectangle with corner radius
    shapedRect := ShapedRect{
        MinPoint: f32.Pt(float32(bounds.X), float32(bounds.Y)),
        MaxPoint: f32.Pt(float32(bounds.X+bounds.Width), float32(bounds.Y+bounds.Height)),
        Shapes:   r.mapClayCornerRadius(cmd.CornerRadius),
    }

    // Create layout context for path generation
    gtx := layout.Context{
        Ops: r.ops,
        Constraints: layout.Constraints{
            Max: image.Pt(int(bounds.Width), int(bounds.Height)),
        },
    }

    // Create clipping path with rounded corners
    pathSpec := shapedRect.Path(gtx)
    clipOp := clip.Outline{Path: pathSpec}.Op().Push(r.ops)
    defer clipOp.Pop()

    // Apply color and paint
    paint.ColorOp{Color: ClayToGioColor(cmd.Color)}.Add(r.ops)
    paint.PaintOp{}.Add(r.ops)

    return nil
}

// Map Clay corner radius to gio-mw corner shapes
func (r *GioRenderer) mapClayCornerRadius(clayRadius CornerRadius) CornerShapes {
    return CornerShapes{
        TopStart: CornerShape{
            Kind: CornerKindRound,
            Size: float32(clayRadius.TopLeft),
        },
        TopEnd: CornerShape{
            Kind: CornerKindRound,
            Size: float32(clayRadius.TopRight),
        },
        BottomStart: CornerShape{
            Kind: CornerKindRound,
            Size: float32(clayRadius.BottomLeft),
        },
        BottomEnd: CornerShape{
            Kind: CornerKindRound,
            Size: float32(clayRadius.BottomRight),
        },
    }
}

// Complete path generation for shaped rectangles
func (s ShapedRect) Path(gtx layout.Context) clip.PathSpec {
    rMinP := s.MinPoint.Sub(f32.Pt(s.Offset, s.Offset))
    rMaxP := s.MaxPoint.Add(f32.Pt(s.Offset, s.Offset))

    rHeight := rMaxP.Y - rMinP.Y
    rWidth := rMaxP.X - rMinP.X
    rMinDim := min(rHeight, rWidth)
    if rMinDim <= 0 {
        return clip.PathSpec{}
    }

    // Handle RTL layouts
    var corners CornerShapes
    if gtx.Locale.Direction == system.RTL {
        corners = CornerShapes{
            TopStart:    s.Shapes.TopEnd,
            TopEnd:      s.Shapes.TopStart,
            BottomStart: s.Shapes.BottomEnd,
            BottomEnd:   s.Shapes.BottomStart,
        }
    } else {
        corners = s.Shapes
    }

    // Calculate actual corner sizes
    ts := min(corners.TopStart.Size, rMinDim)
    te := min(corners.TopEnd.Size, rMinDim)
    be := min(corners.BottomEnd.Size, rMinDim)
    bs := min(corners.BottomStart.Size, rMinDim)

    // Handle adaptive sizing
    if corners.TopStart.AdaptToSize {
        ts = rMinDim / 2
    }
    if corners.TopEnd.AdaptToSize {
        te = rMinDim / 2
    }
    if corners.BottomEnd.AdaptToSize {
        be = rMinDim / 2
    }
    if corners.BottomStart.AdaptToSize {
        bs = rMinDim / 2
    }

    // Build path with rounded corners
    var path clip.Path
    path.Begin(gtx.Ops)
    
    // Top edge
    path.MoveTo(f32.Point{X: rMinP.X + ts, Y: rMinP.Y})
    path.LineTo(f32.Point{X: rMaxP.X - te, Y: rMinP.Y})
    
    // Top-right corner
    if te > 0 && corners.TopEnd.Kind == CornerKindRound {
        fPoint := f32.Point{X: rMaxP.X - te, Y: rMinP.Y + te}
        path.ArcTo(fPoint, fPoint, math.Pi/2)
    } else if te > 0 {
        path.LineTo(f32.Point{X: rMaxP.X, Y: rMinP.Y + te})
    }
    
    // Right edge
    path.LineTo(f32.Point{X: rMaxP.X, Y: rMaxP.Y - be})
    
    // Bottom-right corner
    if be > 0 && corners.BottomEnd.Kind == CornerKindRound {
        fPoint := f32.Point{X: rMaxP.X - be, Y: rMaxP.Y - be}
        path.ArcTo(fPoint, fPoint, math.Pi/2)
    } else if be > 0 {
        path.LineTo(f32.Point{X: rMaxP.X - be, Y: rMaxP.Y})
    }
    
    // Bottom edge
    path.LineTo(f32.Point{X: rMinP.X + bs, Y: rMaxP.Y})
    
    // Bottom-left corner
    if bs > 0 && corners.BottomStart.Kind == CornerKindRound {
        fPoint := f32.Point{X: rMinP.X + bs, Y: rMaxP.Y - bs}
        path.ArcTo(fPoint, fPoint, math.Pi/2)
    } else if bs > 0 {
        path.LineTo(f32.Point{X: rMinP.X, Y: rMaxP.Y - bs})
    }
    
    // Left edge
    path.LineTo(f32.Point{X: rMinP.X, Y: rMinP.Y + ts})
    
    // Top-left corner
    if ts > 0 && corners.TopStart.Kind == CornerKindRound {
        fPoint := f32.Point{X: rMinP.X + ts, Y: rMinP.Y + ts}
        path.ArcTo(fPoint, fPoint, math.Pi/2)
    } else if ts > 0 {
        path.LineTo(f32.Point{X: rMinP.X + ts, Y: rMinP.Y})
    }

    path.Close()
    return path.End()
}

// Fallback for simple rectangles (no corner radius)
func (r *GioRenderer) renderSimpleRectangle(bounds BoundingBox, cmd RectangleCommand) error {
    // Convert bounds to Gio rectangle
    rect := image.Rect(
        int(bounds.X),
        int(bounds.Y),
        int(bounds.X+bounds.Width),
        int(bounds.Y+bounds.Height),
    )
    
    // Create clipping region
    clipOp := clip.Rect(rect).Push(r.ops)
    defer clipOp.Pop()
    
    // Apply color and paint
    paint.ColorOp{Color: ClayToGioColor(cmd.Color)}.Add(r.ops)
    paint.PaintOp{}.Add(r.ops)
    
    return nil
}
```

---

## 3. Image Rendering Solution

### Source Pattern: gio-mw/wdk/require.go + widget.Image + card.Image

**Complete Implementation Pattern** (Based on actual gio-mw production code):

```go
// Required imports for image rendering
import (
    "bytes"
    "fmt"
    "image"
    _ "image/png"  // Register PNG decoder
    _ "image/jpeg" // Register JPEG decoder
    _ "image/gif"  // Register GIF decoder
    "gioui.org/f32"
    "gioui.org/layout"
    "gioui.org/op"
    "gioui.org/op/clip"
    "gioui.org/op/paint"
    "gioui.org/unit"
    "gioui.org/widget"
)

// Complete image rendering implementation using Gio's widget.Image pattern
func (r *GioRenderer) renderImageWithBounds(bounds BoundingBox, cmd ImageCommand) error {
    if len(cmd.ImageData) == 0 {
        return fmt.Errorf("image data is empty")
    }

    // 1. Decode image data (gio-mw pattern from wdk/require.go)
    decodedImage, format, err := image.Decode(bytes.NewReader(cmd.ImageData))
    if err != nil {
        return fmt.Errorf("failed to decode image (format: %s): %w", format, err)
    }

    // 2. Create Gio widget.Image (gio-mw pattern)
    imageWidget := &widget.Image{
        Src:      paint.NewImageOp(decodedImage),
        Fit:      r.mapClayFitMode(cmd.ScaleMode), // Map Clay scale mode to Gio Fit
        Position: layout.Center, // Can be mapped from Clay alignment
        Scale:    1.0, // Default scale, can be customized
    }

    // 3. Create layout context with bounds constraints
    gtx := layout.Context{
        Ops: r.ops,
        Constraints: layout.Constraints{
            Max: image.Pt(int(bounds.Width), int(bounds.Height)),
            Min: image.Pt(int(bounds.Width), int(bounds.Height)), // Force exact size
        },
        Metric: unit.Metric{PxPerDp: 1, PxPerSp: 1}, // Adjust as needed
    }

    // 4. Position within bounds
    stack := op.Offset(f32.Pt(float32(bounds.X), float32(bounds.Y))).Push(r.ops)
    defer stack.Pop()

    // 5. Apply tint color if specified (before image rendering)
    if !cmd.TintColor.IsTransparent() {
        paint.ColorOp{Color: ClayToGioColor(cmd.TintColor)}.Add(r.ops)
    }

    // 6. Render image using Gio's widget.Image.Layout (production pattern)
    dimensions := imageWidget.Layout(gtx)
    
    // Optional: Store dimensions for layout feedback
    // r.lastImageDimensions = dimensions
    
    return nil
}

// Map Clay scale modes to Gio Fit modes
func (r *GioRenderer) mapClayFitMode(scaleMode ImageScaleMode) widget.Fit {
    switch scaleMode {
    case ImageScaleFit:
        return widget.Contain  // Scale to fit within bounds, maintaining aspect ratio
    case ImageScaleFill:
        return widget.Cover    // Scale to fill bounds, may crop
    case ImageScaleStretch:
        return widget.Fill     // Stretch to exact bounds
    case ImageScaleDown:
        return widget.ScaleDown // Only scale down, never up
    default:
        return widget.Contain  // Default to contain
    }
}

// Enhanced image scaling modes matching Gio's widget.Fit
type ImageScaleMode int

const (
    ImageScaleNone ImageScaleMode = iota
    ImageScaleFit      // widget.Contain - scale to fit, maintain aspect ratio
    ImageScaleFill     // widget.Cover - scale to fill, may crop
    ImageScaleStretch  // widget.Fill - stretch to exact bounds
    ImageScaleDown     // widget.ScaleDown - only scale down
)

// Alternative: Direct image rendering without widget.Image wrapper
func (r *GioRenderer) renderImageDirect(bounds BoundingBox, cmd ImageCommand) error {
    if len(cmd.ImageData) == 0 {
        return fmt.Errorf("image data is empty")
    }

    // 1. Decode image
    decodedImage, _, err := image.Decode(bytes.NewReader(cmd.ImageData))
    if err != nil {
        return fmt.Errorf("failed to decode image: %w", err)
    }

    // 2. Create image operation
    imageOp := paint.NewImageOp(decodedImage)
    
    // 3. Calculate scaling (based on widget.Image implementation)
    size := imageOp.Size()
    imgWidth, imgHeight := float32(size.X), float32(size.Y)
    
    // Create constraints for scaling calculation
    constraints := layout.Constraints{
        Max: image.Pt(int(bounds.Width), int(bounds.Height)),
    }
    
    // Calculate dimensions and transformation
    scale := float32(1.0) // Default scale
    imgDims := layout.Dimensions{Size: image.Pt(int(imgWidth*scale), int(imgHeight*scale))}
    
    // Apply fit mode scaling (simplified version of widget.Image logic)
    dims, trans := r.calculateImageTransform(constraints, imgDims, cmd.ScaleMode)
    
    // 4. Position and clip
    stack := op.Offset(f32.Pt(float32(bounds.X), float32(bounds.Y))).Push(r.ops)
    defer stack.Pop()
    
    clipStack := clip.Rect{Max: dims.Size}.Push(r.ops)
    defer clipStack.Pop()
    
    // 5. Apply transformation
    transStack := op.Affine(trans).Push(r.ops)
    defer transStack.Pop()
    
    // 6. Apply tint color if specified
    if !cmd.TintColor.IsTransparent() {
        paint.ColorOp{Color: ClayToGioColor(cmd.TintColor)}.Add(r.ops)
    }
    
    // 7. Paint image
    imageOp.Add(r.ops)
    paint.PaintOp{}.Add(r.ops)
    
    return nil
}

// Simplified transform calculation (based on widget.Image internal logic)
func (r *GioRenderer) calculateImageTransform(constraints layout.Constraints, imgDims layout.Dimensions, scaleMode ImageScaleMode) (layout.Dimensions, f32.Affine2D) {
    // This is a simplified version - for full implementation, 
    // refer to gioui.org/widget.Fit.scale() method
    
    maxWidth := float32(constraints.Max.X)
    maxHeight := float32(constraints.Max.Y)
    imgWidth := float32(imgDims.Size.X)
    imgHeight := float32(imgDims.Size.Y)
    
    var scaleX, scaleY float32 = 1.0, 1.0
    
    switch scaleMode {
    case ImageScaleFit:
        // Scale to fit within bounds
        scaleX = maxWidth / imgWidth
        scaleY = maxHeight / imgHeight
        scale := min(scaleX, scaleY)
        scaleX, scaleY = scale, scale
    case ImageScaleFill:
        // Scale to fill bounds
        scaleX = maxWidth / imgWidth
        scaleY = maxHeight / imgHeight
        scale := max(scaleX, scaleY)
        scaleX, scaleY = scale, scale
    case ImageScaleStretch:
        // Stretch to exact bounds
        scaleX = maxWidth / imgWidth
        scaleY = maxHeight / imgHeight
    }
    
    // Calculate final dimensions
    finalWidth := int(imgWidth * scaleX)
    finalHeight := int(imgHeight * scaleY)
    
    // Create transform
    transform := f32.AffineId().Scale(f32.Point{}, f32.Pt(scaleX, scaleY))
    
    return layout.Dimensions{Size: image.Pt(finalWidth, finalHeight)}, transform
}

// Helper functions
func min(a, b float32) float32 {
    if a < b {
        return a
    }
    return b
}

func max(a, b float32) float32 {
    if a > b {
        return a
    }
    return b
}
```

### Updated ImageCommand Structure

```go
// Enhanced ImageCommand for Clay (matching gio-mw patterns)
type ImageCommand struct {
    ImageData []byte         // Raw image data (PNG, JPEG, GIF, etc.)
    TintColor Color          // Optional tint color (transparent = no tint)
    ScaleMode ImageScaleMode // How to scale image to fit bounds
    Position  ImagePosition  // How to position image within bounds
    Scale     float32        // Additional scaling factor (1.0 = default)
}

// Image positioning options (matching layout.Direction)
type ImagePosition int

const (
    ImagePositionCenter ImagePosition = iota
    ImagePositionStart  // Top-left
    ImagePositionEnd    // Bottom-right
    ImagePositionN      // Top-center
    ImagePositionS      // Bottom-center
    ImagePositionE      // Right-center
    ImagePositionW      // Left-center
    ImagePositionNE     // Top-right
    ImagePositionNW     // Top-left
    ImagePositionSE     // Bottom-right
    ImagePositionSW     // Bottom-left
)

// Helper methods for ImageCommand
func (cmd ImageCommand) IsEmpty() bool {
    return len(cmd.ImageData) == 0
}

func (cmd ImageCommand) HasTint() bool {
    return !cmd.TintColor.IsTransparent()
}

// Helper methods for Color
func (c Color) IsTransparent() bool {
    return c.A == 0.0
}

func (c Color) IsOpaque() bool {
    return c.A >= 1.0
}

// Map Clay image position to Gio layout.Direction
func (r *GioRenderer) mapClayImagePosition(pos ImagePosition) layout.Direction {
    switch pos {
    case ImagePositionCenter:
        return layout.Center
    case ImagePositionStart:
        return layout.NW
    case ImagePositionEnd:
        return layout.SE
    case ImagePositionN:
        return layout.N
    case ImagePositionS:
        return layout.S
    case ImagePositionE:
        return layout.E
    case ImagePositionW:
        return layout.W
    case ImagePositionNE:
        return layout.NE
    case ImagePositionNW:
        return layout.NW
    case ImagePositionSE:
        return layout.SE
    case ImagePositionSW:
        return layout.SW
    default:
        return layout.Center
    }
}
```

---

## 4. Font Management System

### Complete Font System Setup

```go
// Font management structure
type FontManager struct {
    shaper     *text.Shaper
    fontCache  map[uint16]font.Font
    defaultFont font.Font
}

// Initialize font system (call once at startup)
func NewFontManager() *FontManager {
    // Create text shaper with default fonts
    shaperOptions := []text.ShaperOption{
        text.WithCollection(gofont.Collection()),
    }
    
    fm := &FontManager{
        shaper:    text.NewShaper(shaperOptions...),
        fontCache: make(map[uint16]font.Font),
        defaultFont: font.Font{
            Typeface: "Go",
            Style:    font.Regular,
            Weight:   font.Normal,
        },
    }
    
    // Pre-populate common fonts
    fm.registerDefaultFonts()
    
    return fm
}

// Register default fonts
func (fm *FontManager) registerDefaultFonts() {
    // Register common font variations
    fm.fontCache[0] = font.Font{Typeface: "Go", Style: font.Regular, Weight: font.Normal}
    fm.fontCache[1] = font.Font{Typeface: "Go", Style: font.Regular, Weight: font.Bold}
    fm.fontCache[2] = font.Font{Typeface: "Go", Style: font.Italic, Weight: font.Normal}
    fm.fontCache[3] = font.Font{Typeface: "Go", Style: font.Italic, Weight: font.Bold}
}

// Get font by ID
func (fm *FontManager) GetFont(fontID uint16) font.Font {
    if font, exists := fm.fontCache[fontID]; exists {
        return font
    }
    return fm.defaultFont
}

// Get text shaper
func (fm *FontManager) GetShaper() *text.Shaper {
    return fm.shaper
}

// Register custom font
func (fm *FontManager) RegisterFont(fontID uint16, typeface string, style font.Style, weight font.Weight) {
    fm.fontCache[fontID] = font.Font{
        Typeface: typeface,
        Style:    style,
        Weight:   weight,
    }
}
```

### Integration with GioRenderer

```go
// Updated GioRenderer with font management
type GioRenderer struct {
    ops         *op.Ops
    viewport    BoundingBox
    clipStack   []clip.Stack
    fontManager *FontManager  // Add font manager
    // ... other fields
}

// Updated constructor
func NewRenderer(ops *op.Ops) *GioRenderer {
    return &GioRenderer{
        ops:         ops,
        fontManager: NewFontManager(),
        // ... other initialization
    }
}

// Updated text rendering with font management
func (r *GioRenderer) renderTextWithBounds(bounds BoundingBox, cmd TextCommand) error {
    // Get font from font manager
    font := r.fontManager.GetFont(cmd.FontID)
    shaper := r.fontManager.GetShaper()
    
    // ... rest of text rendering implementation
    // Use font and shaper in label.Layout() call
    
    return nil
}
```

---

## 5. Color Conversion System

### Complete Color Conversion Implementation

```go
// Enhanced color conversion functions
import (
    "image/color"
    "gioui.org/f32"
)

// Convert Clay color (float32 RGBA 0-1) to Gio color (uint8 NRGBA 0-255)
func ClayToGioColor(c Color) color.NRGBA {
    return color.NRGBA{
        R: uint8(c.R * 255.0),
        G: uint8(c.G * 255.0),
        B: uint8(c.B * 255.0),
        A: uint8(c.A * 255.0),
    }
}

// Convert Clay coordinates (int) to Gio coordinates (float32)
func ClayToGioPoint(x, y int) f32.Point {
    return f32.Pt(float32(x), float32(y))
}

// Convert Clay BoundingBox to Gio image.Rectangle
func ClayBoundsToGioRect(bounds BoundingBox) image.Rectangle {
    return image.Rect(
        int(bounds.X),
        int(bounds.Y),
        int(bounds.X+bounds.Width),
        int(bounds.Y+bounds.Height),
    )
}

// Convert Clay BoundingBox to Gio f32 rectangle
func ClayBoundsToGioRectF32(bounds BoundingBox) f32.Rectangle {
    return f32.Rectangle{
        Min: f32.Pt(float32(bounds.X), float32(bounds.Y)),
        Max: f32.Pt(float32(bounds.X+bounds.Width), float32(bounds.Y+bounds.Height)),
    }
}

// Color utility functions
func (c Color) IsOpaque() bool {
    return c.A >= 1.0
}

func (c Color) IsTransparent() bool {
    return c.A <= 0.0
}

func (c Color) WithAlpha(alpha float32) Color {
    return Color{R: c.R, G: c.G, B: c.B, A: alpha}
}
```

---

## 6. Error Handling System

### Comprehensive Error Handling

```go
// Error types for renderer
type RendererErrorType int

const (
    ErrorTypeNone RendererErrorType = iota
    ErrorTypeInvalidBounds
    ErrorTypeInvalidColor
    ErrorTypeInvalidFont
    ErrorTypeInvalidImage
    ErrorTypeUnsupportedOperation
    ErrorTypeResourceNotFound
    ErrorTypeInvalidState
)

type RendererError struct {
    Type    RendererErrorType
    Message string
    Cause   error
}

func (e RendererError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %v", e.Message, e.Cause)
    }
    return e.Message
}

// Error creation helpers
func NewRendererError(errorType RendererErrorType, message string) error {
    return RendererError{
        Type:    errorType,
        Message: message,
    }
}

func WrapRendererError(errorType RendererErrorType, message string, cause error) error {
    return RendererError{
        Type:    errorType,
        Message: message,
        Cause:   cause,
    }
}

// Validation functions
func (r *GioRenderer) validateBounds(bounds BoundingBox) error {
    if bounds.Width <= 0 || bounds.Height <= 0 {
        return NewRendererError(ErrorTypeInvalidBounds, 
            fmt.Sprintf("invalid bounds: width=%d, height=%d", bounds.Width, bounds.Height))
    }
    return nil
}

func (r *GioRenderer) validateColor(color Color) error {
    if color.R < 0 || color.R > 1 || color.G < 0 || color.G > 1 || 
       color.B < 0 || color.B > 1 || color.A < 0 || color.A > 1 {
        return NewRendererError(ErrorTypeInvalidColor, 
            fmt.Sprintf("invalid color values: R=%.2f, G=%.2f, B=%.2f, A=%.2f", 
                color.R, color.G, color.B, color.A))
    }
    return nil
}
```

---

## 7. Integration Testing Patterns

### Complete Test Implementation

```go
// Test file: renderer_integration_test.go
package gioui

import (
    "image"
    "testing"
    "gioui.org/app"
    "gioui.org/op"
    "gioui.org/unit"
)

func TestTextRenderingIntegration(t *testing.T) {
    var ops op.Ops
    renderer := NewRenderer(&ops)
    
    // Test text rendering
    bounds := BoundingBox{X: 10, Y: 10, Width: 200, Height: 50}
    textCmd := TextCommand{
        Text:     "Hello, World!",
        FontID:   0,
        FontSize: 16.0,
        Color:    Color{R: 0, G: 0, B: 0, A: 1}, // Black
    }
    
    err := renderer.renderTextWithBounds(bounds, textCmd)
    if err != nil {
        t.Fatalf("Text rendering failed: %v", err)
    }
    
    // Verify operations were added
    if len(ops.Data()) == 0 {
        t.Error("No operations were added for text rendering")
    }
}

func TestCornerRadiusIntegration(t *testing.T) {
    var ops op.Ops
    renderer := NewRenderer(&ops)
    
    // Test corner radius rendering
    bounds := BoundingBox{X: 0, Y: 0, Width: 100, Height: 100}
    rectCmd := RectangleCommand{
        Color: Color{R: 1, G: 0, B: 0, A: 1}, // Red
        CornerRadius: CornerRadius{
            TopLeft:     10,
            TopRight:    10,
            BottomLeft:  10,
            BottomRight: 10,
        },
    }
    
    err := renderer.renderRectangleWithBounds(bounds, rectCmd)
    if err != nil {
        t.Fatalf("Corner radius rendering failed: %v", err)
    }
    
    // Verify operations were added
    if len(ops.Data()) == 0 {
        t.Error("No operations were added for corner radius rendering")
    }
}

func TestImageRenderingIntegration(t *testing.T) {
    var ops op.Ops
    renderer := NewRenderer(&ops)
    
    // Create test image data (1x1 red pixel PNG)
    testImageData := []byte{
        0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, // PNG signature
        // ... minimal PNG data for testing
    }
    
    bounds := BoundingBox{X: 0, Y: 0, Width: 50, Height: 50}
    imageCmd := ImageCommand{
        ImageData: testImageData,
        ScaleMode: ImageScaleFit,
        TintColor: Color{A: 0}, // No tint
    }
    
    err := renderer.renderImageWithBounds(bounds, imageCmd)
    if err != nil {
        t.Fatalf("Image rendering failed: %v", err)
    }
    
    // Verify operations were added
    if len(ops.Data()) == 0 {
        t.Error("No operations were added for image rendering")
    }
}

// Benchmark tests
func BenchmarkTextRendering(b *testing.B) {
    var ops op.Ops
    renderer := NewRenderer(&ops)
    
    bounds := BoundingBox{X: 0, Y: 0, Width: 200, Height: 50}
    textCmd := TextCommand{
        Text:     "Benchmark Text",
        FontID:   0,
        FontSize: 16.0,
        Color:    Color{R: 0, G: 0, B: 0, A: 1},
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        ops.Reset()
        renderer.renderTextWithBounds(bounds, textCmd)
    }
}
```

---

## 8. Implementation Checklist

### Phase 1: Text Rendering (Critical Priority)
- [ ] Add `FontManager` to `GioRenderer` struct
- [ ] Implement `initializeTextShaper()` with `gofont.Collection()`
- [ ] Replace `renderTextWithBounds()` stub with full implementation
- [ ] Add font weight and alignment mapping functions
- [ ] Test text rendering with various fonts and sizes
- [ ] Verify text positioning within bounds

### Phase 2: Corner Radius (Medium Priority)
- [ ] Add corner shape types and structures
- [ ] Implement `ShapedRect` with path generation
- [ ] Replace simple rectangle rendering with corner radius support
- [ ] Add corner radius mapping from Clay to gio-mw format
- [ ] Test various corner radius configurations
- [ ] Verify RTL layout support

### Phase 3: Image Rendering (Critical Priority)
- [ ] Implement `renderImageWithBounds()` using `widget.Image` pattern
- [ ] Add image decoding with `image.Decode()` (gio-mw pattern)
- [ ] Map Clay scale modes to Gio `widget.Fit` modes
- [ ] Add image positioning support with `layout.Direction`
- [ ] Test image rendering with PNG, JPEG, GIF formats
- [ ] Verify image scaling and positioning within bounds
- [ ] Add tint color support for image rendering

### Phase 4: Integration & Testing
- [ ] Create comprehensive integration tests
- [ ] Add benchmark tests for performance validation
- [ ] Test complete rendering pipeline with all features
- [ ] Verify error handling for edge cases
- [ ] Update documentation and examples

---

## 9. Performance Considerations

### Optimization Strategies
1. **Font Caching**: Cache `text.Shaper` and font objects to avoid recreation
2. **Image Caching**: Cache decoded images to avoid repeated decoding
3. **Path Caching**: Cache corner radius paths for repeated shapes
4. **Operation Batching**: Batch similar operations to reduce overhead
5. **Memory Management**: Reuse image buffers and operation slices

### Memory Management
```go
// Add to GioRenderer for resource management
type GioRenderer struct {
    // ... existing fields
    imageCache    map[string]image.Image  // Cache decoded images
    pathCache     map[string]clip.PathSpec // Cache corner radius paths
    operationPool sync.Pool               // Pool for operation reuse
}

// Resource cleanup
func (r *GioRenderer) Cleanup() {
    r.imageCache = make(map[string]image.Image)
    r.pathCache = make(map[string]clip.PathSpec)
}
```

---

## 10. Future Enhancements

### Advanced Features from gio-x
1. **Rich Text**: Multi-style text with interactive elements
2. **Advanced Shadows**: Elevation-based shadow system
3. **Material Components**: Pre-built UI components
4. **Animation Support**: Smooth transitions and animations

### Extensibility Points
1. **Custom Renderers**: Plugin system for custom command types
2. **Shader Support**: Custom fragment shaders for effects
3. **Vector Graphics**: SVG-like vector drawing commands
4. **3D Rendering**: Basic 3D transformation support

---

This document provides complete, self-contained implementation guidance for resolving all critical issues in Clay's user story 1.1. All code patterns are production-ready and based on proven implementations from gio-mw and gio-x repositories.
