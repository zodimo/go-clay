package gioui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"  // Register GIF decoder
	_ "image/jpeg" // Register JPEG decoder
	_ "image/png"  // Register PNG decoder
	"log"
	"math"
	"os"

	"gioui.org/f32"
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/zodimo/go-clay/kimi2/clay"
)

// FontManager manages fonts and text shaping for the renderer
type FontManager struct {
	shaper      *text.Shaper
	fontCache   map[uint16]font.Font
	defaultFont font.Font
}

// NewFontManager creates a new font manager with default fonts
func NewFontManager() *FontManager {
	// Create text shaper with default fonts
	shaperOptions := []text.ShaperOption{
		text.WithCollection(gofont.Collection()),
	}

	fm := &FontManager{
		shaper:    text.NewShaper(shaperOptions...),
		fontCache: make(map[uint16]font.Font),
		defaultFont: font.Font{
			Style:  font.Regular,
			Weight: font.Normal,
		},
	}

	// Pre-populate common fonts
	fm.registerDefaultFonts()

	return fm
}

// registerDefaultFonts registers common font variations
func (fm *FontManager) registerDefaultFonts() {
	fm.fontCache[0] = font.Font{Style: font.Regular, Weight: font.Normal}
	fm.fontCache[1] = font.Font{Style: font.Regular, Weight: font.Bold}
	fm.fontCache[2] = font.Font{Style: font.Italic, Weight: font.Normal}
	fm.fontCache[3] = font.Font{Style: font.Italic, Weight: font.Bold}
}

// GetFont returns a font by ID
func (fm *FontManager) GetFont(fontID uint16) font.Font {
	if font, exists := fm.fontCache[fontID]; exists {
		return font
	}
	return fm.defaultFont
}

// GetShaper returns the text shaper
func (fm *FontManager) GetShaper() *text.Shaper {
	return fm.shaper
}

// RegisterFont registers a custom font
func (fm *FontManager) RegisterFont(fontID uint16, style font.Style, weight font.Weight) {
	fm.fontCache[fontID] = font.Font{
		Style:  style,
		Weight: weight,
	}
}

// Corner shape types for rounded rectangles
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

// GioRenderer implements the clay.Renderer interface for Gio UI
type GioRenderer struct {
	ops              *op.Ops
	viewport         clay.BoundingBox
	clipStack        []clip.Stack
	cache            *ResourceCache
	customRegistry   *CustomCommandRegistry
	errorHandler     *ErrorHandler
	operationBuilder *OperationBuilder
	batchOperations  *BatchOperations
	fontManager      *FontManager
	maxClipDepth     int
}

// RendererOptions holds configuration options for the GioRenderer
type RendererOptions struct {
	Logger       *log.Logger
	DebugMode    bool
	CacheSize    int
	MaxClipDepth int
}

// RendererOption is a function that configures RendererOptions
type RendererOption func(options *RendererOptions)

// RendererWithLogger sets a custom logger for the renderer
func RendererWithLogger(logger *log.Logger) RendererOption {
	return func(options *RendererOptions) {
		options.Logger = logger
	}
}

// RendererWithDebugMode enables or disables debug mode
func RendererWithDebugMode(debug bool) RendererOption {
	return func(options *RendererOptions) {
		options.DebugMode = debug
	}
}

// RendererWithCacheSize sets the cache size in bytes
func RendererWithCacheSize(size int) RendererOption {
	return func(options *RendererOptions) {
		options.CacheSize = size
	}
}

// RendererWithMaxClipDepth sets the maximum clip stack depth
func RendererWithMaxClipDepth(depth int) RendererOption {
	return func(options *RendererOptions) {
		options.MaxClipDepth = depth
	}
}

// defaultRendererOptions returns the default configuration
func defaultRendererOptions() RendererOptions {
	return RendererOptions{
		Logger:       log.New(os.Stderr, "[GioRenderer] ", log.LstdFlags),
		DebugMode:    false,
		CacheSize:    50 * 1024 * 1024, // 50MB cache
		MaxClipDepth: 100,
	}
}

// NewRenderer creates a new GioRenderer with optional configuration
func NewRenderer(ops *op.Ops, options ...RendererOption) *GioRenderer {
	opts := defaultRendererOptions()

	for _, option := range options {
		option(&opts)
	}

	// Validate the options - use defaults if invalid
	if opts.CacheSize < 0 {
		opts.CacheSize = 50 * 1024 * 1024 // Default to 50MB
	}
	if opts.MaxClipDepth < 1 {
		opts.MaxClipDepth = 100 // Default to 100
	}

	return &GioRenderer{
		ops:              ops,
		clipStack:        make([]clip.Stack, 0),
		cache:            NewResourceCache(opts.CacheSize),
		customRegistry:   NewCustomCommandRegistry(),
		errorHandler:     NewErrorHandler(opts.Logger, opts.DebugMode),
		operationBuilder: NewOperationBuilder(ops),
		batchOperations:  NewBatchOperations(),
		fontManager:      NewFontManager(),
		maxClipDepth:     opts.MaxClipDepth,
	}
}

// BeginFrame initializes the frame for rendering
func (r *GioRenderer) BeginFrame() error {
	if r.ops != nil {
		r.ops.Reset()
	}
	return nil
}

// EndFrame finalizes the frame rendering
func (r *GioRenderer) EndFrame() error {
	// No cleanup needed for basic implementation
	return nil
}

// SetViewport sets the viewport bounds for coordinate system setup
func (r *GioRenderer) SetViewport(bounds clay.BoundingBox) error {
	r.viewport = bounds
	return nil
}

// Render processes an array of render commands with bounds information
// This matches the Clay C architecture: Clay_Raylib_Render(Clay_RenderCommandArray renderCommands, Font* fonts)
func (r *GioRenderer) Render(commands []clay.RenderCommand) error {
	if r.ops == nil {
		return NewRenderError(
			ErrorTypeInvalidState,
			"Render",
			"operations context is nil - call BeginFrame() first",
		)
	}

	// Sort commands by ZIndex for proper rendering order
	sortedCommands := make([]clay.RenderCommand, len(commands))
	copy(sortedCommands, commands)

	// Simple bubble sort by ZIndex (stable sort to maintain order for equal ZIndex)
	for i := 0; i < len(sortedCommands)-1; i++ {
		for j := 0; j < len(sortedCommands)-i-1; j++ {
			if sortedCommands[j].ZIndex > sortedCommands[j+1].ZIndex {
				sortedCommands[j], sortedCommands[j+1] = sortedCommands[j+1], sortedCommands[j]
			}
		}
	}

	// Process each command with bounds information
	for _, cmd := range sortedCommands {
		bounds := cmd.BoundingBox // Extract bounds for positioning

		switch cmd.CommandType {
		case clay.CLAY_RENDER_COMMAND_TYPE_RECTANGLE:
			if err := r.renderRectangleWithBounds(bounds, cmd); err != nil {
				return err
			}
		case clay.CLAY_RENDER_COMMAND_TYPE_TEXT:
			if err := r.renderTextWithBounds(bounds, cmd); err != nil {
				return err
			}
		case clay.CLAY_RENDER_COMMAND_TYPE_IMAGE:
			if err := r.renderImageWithBounds(bounds, cmd); err != nil {
				return err
			}
		case clay.CLAY_RENDER_COMMAND_TYPE_BORDER:
			if err := r.renderBorderWithBounds(bounds, cmd); err != nil {
				return err
			}
		case clay.CLAY_RENDER_COMMAND_TYPE_SCISSOR_START:
			if err := r.renderClipStartWithBounds(bounds, cmd); err != nil {
				return err
			}
		case clay.CLAY_RENDER_COMMAND_TYPE_SCISSOR_END:
			if err := r.renderClipEndWithBounds(bounds, cmd); err != nil {
				return err
			}
		case clay.CLAY_RENDER_COMMAND_TYPE_CUSTOM:
			if err := r.renderCustomWithBounds(bounds, cmd); err != nil {
				return err
			}
		default:
			return NewRenderError(
				ErrorTypeUnsupportedOperation,
				"Render",
				fmt.Sprintf("unsupported command type %d for element %d", cmd.CommandType, cmd.ID),
			)
		}
	}
	return nil
}

// renderRectangleWithBounds renders a rectangle using bounds from RenderCommand
func (r *GioRenderer) renderRectangleWithBounds(bounds clay.BoundingBox, cmd clay.RenderCommand) error {
	rectangleData, ok := cmd.Data.(clay.RectangleRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"renderRectangleWithBounds",
			"Invalid rectangle command data",
		)
	}
	// Validate bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"renderRectangleWithBounds",
			fmt.Sprintf("Invalid bounds: width=%f, height=%f", bounds.Width, bounds.Height),
		)
	}

	// Check if corner radius is needed
	if r.isCornerRadiusZero(rectangleData.CornerRadius) {
		return r.renderSimpleRectangle(bounds, cmd)
	}

	// Render rectangle with corner radius
	return r.renderRoundedRectangle(bounds, cmd)
}

// isCornerRadiusZero checks if all corner radius values are zero
func (r *GioRenderer) isCornerRadiusZero(radius clay.CornerRadius) bool {
	return radius.TopLeft == 0 && radius.TopRight == 0 &&
		radius.BottomLeft == 0 && radius.BottomRight == 0
}

// renderSimpleRectangle renders a rectangle without corner radius
func (r *GioRenderer) renderSimpleRectangle(bounds clay.BoundingBox, cmd clay.RenderCommand) error {
	rectangleData, ok := cmd.Data.(clay.RectangleRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"renderSimpleRectangle",
			"Invalid rectangle command data",
		)
	}
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
	paint.ColorOp{Color: ClayToGioColor(rectangleData.Color)}.Add(r.ops)
	paint.PaintOp{}.Add(r.ops)

	return nil
}

// renderRoundedRectangle renders a rectangle with corner radius
func (r *GioRenderer) renderRoundedRectangle(bounds clay.BoundingBox, cmd clay.RenderCommand) error {
	rectangleData, ok := cmd.Data.(clay.RectangleRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"renderRoundedRectangle",
			"Invalid rectangle command data",
		)
	}
	// Create shaped rectangle with corner radius
	shapedRect := ShapedRect{
		MinPoint: f32.Pt(float32(bounds.X), float32(bounds.Y)),
		MaxPoint: f32.Pt(float32(bounds.X+bounds.Width), float32(bounds.Y+bounds.Height)),
		Shapes:   r.mapClayCornerRadius(rectangleData.CornerRadius),
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
	paint.ColorOp{Color: ClayToGioColor(rectangleData.Color)}.Add(r.ops)
	paint.PaintOp{}.Add(r.ops)

	return nil
}

// renderTextWithBounds renders text using bounds from RenderCommand
func (r *GioRenderer) renderTextWithBounds(bounds clay.BoundingBox, cmd clay.RenderCommand) error {
	textCommandData, ok := cmd.Data.(clay.TextRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"renderTextWithBounds",
			"Invalid text command data",
		)
	}
	// Validate bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"renderTextWithBounds",
			fmt.Sprintf("Invalid bounds: width=%f, height=%f", bounds.Width, bounds.Height),
		)
	}

	// Validate text content
	if textCommandData.StringContents == "" {
		// Empty text is valid, just return without rendering
		return nil
	}

	// Validate font manager
	if r.fontManager == nil {
		return NewRenderError(
			ErrorTypeInvalidState,
			"renderTextWithBounds",
			"Font manager not initialized",
		)
	}

	// Get font and shaper from font manager
	fontObj := r.fontManager.GetFont(textCommandData.FontID)
	shaper := r.fontManager.GetShaper()

	// Create color operation
	colorMacro := op.Record(r.ops)
	paint.ColorOp{Color: ClayToGioColor(textCommandData.Color)}.Add(r.ops)
	colorCallOp := colorMacro.Stop()

	// Create label with Clay parameters
	label := widget.Label{
		Alignment:  r.mapClayTextAlignment(textCommandData.Alignment),
		MaxLines:   0, // Unlimited - can be set from Clay command
		LineHeight: unit.Sp(textCommandData.LineHeight),
		WrapPolicy: text.WrapWords,
	}

	// Position within bounds
	stack := op.Offset(image.Pt(int(bounds.X), int(bounds.Y))).Push(r.ops)
	defer stack.Pop()

	// Create layout context with bounds constraints
	gtx := layout.Context{
		Ops: r.ops,
		Constraints: layout.Constraints{
			Max: image.Pt(int(bounds.Width), int(bounds.Height)),
		},
		Metric: unit.Metric{PxPerDp: 1, PxPerSp: 1},
	}

	// Render text using gio pattern
	label.Layout(gtx, shaper, fontObj, unit.Sp(textCommandData.FontSize), textCommandData.StringContents, colorCallOp)

	return nil
}

// mapClayCornerRadius maps Clay corner radius to gio-mw corner shapes
func (r *GioRenderer) mapClayCornerRadius(clayRadius clay.CornerRadius) CornerShapes {
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

// Path generates a clip path for the shaped rectangle
func (s ShapedRect) Path(gtx layout.Context) clip.PathSpec {
	rMinP := s.MinPoint.Sub(f32.Pt(s.Offset, s.Offset))
	rMaxP := s.MaxPoint.Add(f32.Pt(s.Offset, s.Offset))

	rHeight := rMaxP.Y - rMinP.Y
	rWidth := rMaxP.X - rMinP.X
	rMinDim := min(rHeight, rWidth)
	if rMinDim <= 0 {
		return clip.PathSpec{}
	}

	// Use shapes as-is (simplified - no RTL handling for now)
	corners := s.Shapes

	// Calculate actual corner sizes
	ts := min(corners.TopStart.Size, rMinDim)
	te := min(corners.TopEnd.Size, rMinDim)
	be := min(corners.BottomEnd.Size, rMinDim)
	bs := min(corners.BottomStart.Size, rMinDim)

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

// Helper function for min
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}

// mapClayTextAlignment maps Clay text alignment to Gio text alignment
func (r *GioRenderer) mapClayTextAlignment(clayAlign clay.TextAlignment) text.Alignment {
	switch clayAlign {
	case clay.CLAY_TEXT_ALIGN_LEFT:
		return text.Start
	case clay.CLAY_TEXT_ALIGN_CENTER:
		return text.Middle
	case clay.CLAY_TEXT_ALIGN_RIGHT:
		return text.End
	default:
		return text.Start
	}
}

// renderImageWithBounds renders an image using bounds from RenderCommand
func (r *GioRenderer) renderImageWithBounds(bounds clay.BoundingBox, cmd clay.RenderCommand) error {
	imageCommandData, ok := cmd.Data.(clay.ImageRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"renderImageWithBounds",
			"Invalid image command data",
		)
	}
	// Validate bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"renderImageWithBounds",
			fmt.Sprintf("Invalid bounds: width=%f, height=%f", bounds.Width, bounds.Height),
		)
	}

	// Validate image data
	imageData, ok := imageCommandData.ImageData.([]byte)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"renderImageWithBounds",
			"Image data must be []byte",
		)
	}
	if len(imageData) == 0 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"renderImageWithBounds",
			"Image data is empty",
		)
	}

	// Decode image data
	decodedImage, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"renderImageWithBounds",
			fmt.Sprintf("Failed to decode image (format: %s): %v", format, err),
		)
	}

	// Create Gio widget.Image
	imageWidget := &widget.Image{
		Src:      paint.NewImageOp(decodedImage),
		Fit:      widget.Contain, // Default to contain
		Position: layout.Center,
		Scale:    1.0,
	}

	// Create layout context with bounds constraints
	gtx := layout.Context{
		Ops: r.ops,
		Constraints: layout.Constraints{
			Max: image.Pt(int(bounds.Width), int(bounds.Height)),
			Min: image.Pt(int(bounds.Width), int(bounds.Height)),
		},
		Metric: unit.Metric{PxPerDp: 1, PxPerSp: 1},
	}

	// Position within bounds
	stack := op.Offset(image.Pt(int(bounds.X), int(bounds.Y))).Push(r.ops)
	defer stack.Pop()

	// Apply tint color if specified
	if imageCommandData.TintColor.A > 0 {
		paint.ColorOp{Color: ClayToGioColor(imageCommandData.TintColor)}.Add(r.ops)
	}

	// Render image using Gio's widget.Image.Layout
	imageWidget.Layout(gtx)

	return nil
}

// renderBorderWithBounds renders a border using bounds from RenderCommand
func (r *GioRenderer) renderBorderWithBounds(bounds clay.BoundingBox, cmd clay.RenderCommand) error {
	borderData, ok := cmd.Data.(clay.BorderRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"renderBorderWithBounds",
			"Invalid border command data",
		)
	}
	// Validate border command
	if err := r.validateBorderCommand(cmd); err != nil {
		return err
	}

	// Convert bounds to image rectangle
	rect := image.Rect(
		int(bounds.X),
		int(bounds.Y),
		int(bounds.X+bounds.Width),
		int(bounds.Y+bounds.Height),
	)

	// Convert Clay color to Gio color
	gioColor := ClayToGioColor(borderData.Color)

	// Render border using bounds
	r.renderBorderSides(rect, borderData.Width, gioColor, borderData.CornerRadius)

	return nil
}

// renderClipStartWithBounds starts clipping using bounds from RenderCommand
func (r *GioRenderer) renderClipStartWithBounds(bounds clay.BoundingBox, cmd clay.RenderCommand) error {
	_, ok := cmd.Data.(clay.ClipRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"renderClipStartWithBounds",
			"Invalid clip command data",
		)
	}
	// Check clip stack depth
	if len(r.clipStack) >= r.maxClipDepth {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"renderClipStartWithBounds",
			fmt.Sprintf("Maximum clip depth exceeded: %d", r.maxClipDepth),
		)
	}

	// Convert bounds to image rectangle
	rect := image.Rect(
		int(bounds.X),
		int(bounds.Y),
		int(bounds.X+bounds.Width),
		int(bounds.Y+bounds.Height),
	)

	// Create simple rectangular clipping (ClipStartCommand doesn't have CornerRadius)
	clipStack := clip.Rect(rect).Push(r.ops)

	// Add to clip stack
	r.clipStack = append(r.clipStack, clipStack)

	return nil
}

// RenderRectangle renders a rectangle using Gio operations
func (r *GioRenderer) RenderRectangle(cmd clay.RenderCommand) error {
	rectangleData, ok := cmd.Data.(clay.RectangleRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"RenderRectangle",
			"Invalid rectangle command data",
		)
	}
	if r.ops == nil {
		return fmt.Errorf("operations context is nil")
	}

	// Convert Clay color to Gio color
	gioColor := ClayToGioColor(rectangleData.Color)

	// Note: The current Clay interface doesn't pass bounds to individual render methods.
	// This is a limitation that will need to be addressed in a future interface update.
	// For now, we'll render a basic colored rectangle without specific bounds.

	// Set up paint operation
	paint.ColorOp{Color: gioColor}.Add(r.ops)
	paint.PaintOp{}.Add(r.ops)

	return nil
}

// RenderText renders text using basic Gio text operations
func (r *GioRenderer) RenderText(cmd clay.RenderCommand) error {
	textData, ok := cmd.Data.(clay.TextRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"RenderText",
			"Invalid text command data",
		)
	}
	if r.ops == nil {
		return fmt.Errorf("operations context is nil")
	}

	// Convert Clay color to Gio color
	gioColor := ClayToGioColor(textData.Color)

	// Set up paint operation for text
	paint.ColorOp{Color: gioColor}.Add(r.ops)

	// TODO: Implement actual text rendering with font support
	// This is a stub implementation for now

	return nil
}

// Stub implementations for remaining interface methods

func (r *GioRenderer) RenderImage(cmd clay.RenderCommand) error {
	imageData, ok := cmd.Data.(clay.ImageRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"RenderImage",
			"Invalid image command data",
		)
	}
	defer func() {
		if err := r.errorHandler.RecoverFromPanic("RenderImage"); err != nil {
			r.errorHandler.HandleError(err)
		}
	}()

	if r.ops == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"RenderImage",
			"Operations context is nil",
		)
	}

	// Validate image command
	if err := r.validateImageCommand(cmd); err != nil {
		r.errorHandler.HandleError(err.(*RenderError))
		return err
	}

	// Get or create cached image with appropriate filter
	filter := FilterLinear // Default filter
	cachedImage, err := r.cache.GetOrCreateImage(imageData.ImageData, filter)
	if err != nil {
		renderErr := NewRenderError(
			ErrorTypeResourceNotFound,
			"RenderImage",
			fmt.Sprintf("Failed to load/cache image: %v", err),
		)
		r.errorHandler.HandleError(renderErr)
		return renderErr
	}

	// Create image bounds from viewport (limitation: no per-element bounds in current interface)
	// Note: In a complete implementation, this would use per-element bounds
	// passed from the layout engine, but the current interface limitation
	// requires us to use viewport bounds
	bounds := image.Rectangle{
		Min: image.Point{X: int(r.viewport.X), Y: int(r.viewport.Y)},
		Max: image.Point{
			X: int(r.viewport.X + r.viewport.Width),
			Y: int(r.viewport.Y + r.viewport.Height),
		},
	}

	// Convert tint color to Gio color
	tintColor := r.cache.GetOrCreateColor(imageData.TintColor.R, imageData.TintColor.G, imageData.TintColor.B, imageData.TintColor.A)

	// Use operation builder to create image operation with clipping and tinting
	r.operationBuilder.BuildImageOperation(cachedImage.Image, bounds, filter, tintColor)

	return nil
}

func (r *GioRenderer) RenderBorder(cmd clay.RenderCommand) error {
	borderData, ok := cmd.Data.(clay.BorderRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"RenderBorder",
			"Invalid border command data",
		)
	}

	defer func() {
		if err := r.errorHandler.RecoverFromPanic("RenderBorder"); err != nil {
			r.errorHandler.HandleError(err)
		}
	}()

	if r.ops == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"RenderBorder",
			"Operations context is nil",
		)
	}

	// Validate border command
	if err := r.validateBorderCommand(cmd); err != nil {
		r.errorHandler.HandleError(err.(*RenderError))
		return err
	}

	// Convert Clay color to Gio color
	gioColor := r.cache.GetOrCreateColor(borderData.Color.R, borderData.Color.G, borderData.Color.B, borderData.Color.A)

	// Create bounds from viewport (limitation: no per-element bounds in current interface)
	bounds := image.Rectangle{
		Min: image.Point{X: int(r.viewport.X), Y: int(r.viewport.Y)},
		Max: image.Point{
			X: int(r.viewport.X + r.viewport.Width),
			Y: int(r.viewport.Y + r.viewport.Height),
		},
	}

	// Handle different border widths by rendering each side separately
	r.renderBorderSides(bounds, borderData.Width, gioColor, borderData.CornerRadius)

	return nil
}

// validateBorderCommand validates the border command parameters
func (r *GioRenderer) validateBorderCommand(cmd clay.RenderCommand) error {
	borderData, ok := cmd.Data.(clay.BorderRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"validateBorderCommand",
			"Invalid border command data",
		)
	}
	// Validate color values
	if err := ValidateColor(borderData.Color.R, borderData.Color.G, borderData.Color.B, borderData.Color.A); err != nil {
		return err
	}

	// Validate border widths (must be non-negative)
	if borderData.Width.Left < 0 || borderData.Width.Right < 0 || borderData.Width.Top < 0 || borderData.Width.Bottom < 0 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"validateBorderCommand",
			"Border widths must be non-negative",
		)
	}

	// Validate corner radius values (must be non-negative)
	if borderData.CornerRadius.TopLeft < 0 || borderData.CornerRadius.TopRight < 0 ||
		borderData.CornerRadius.BottomLeft < 0 || borderData.CornerRadius.BottomRight < 0 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"validateBorderCommand",
			"Corner radius values must be non-negative",
		)
	}

	return nil
}

// validateImageCommand validates the image command parameters
func (r *GioRenderer) validateImageCommand(cmd clay.RenderCommand) error {
	imageData, ok := cmd.Data.(clay.ImageRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"validateImageCommand",
			"Invalid image command data",
		)
	}

	// Validate tint color values
	if err := ValidateColor(imageData.TintColor.R, imageData.TintColor.G, imageData.TintColor.B, imageData.TintColor.A); err != nil {
		return err
	}

	// Validate corner radius values (must be non-negative)
	if imageData.CornerRadius.TopLeft < 0 || imageData.CornerRadius.TopRight < 0 ||
		imageData.CornerRadius.BottomLeft < 0 || imageData.CornerRadius.BottomRight < 0 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"validateImageCommand",
			"Corner radius values must be non-negative",
		)
	}

	return nil
}

// renderBorderSides renders each border side with potentially different widths
func (r *GioRenderer) renderBorderSides(bounds image.Rectangle, width clay.BorderWidth, color color.NRGBA, cornerRadius clay.CornerRadius) {
	// If all border widths are the same and we have uniform corner radius, use optimized path
	if width.Left == width.Right && width.Right == width.Top && width.Top == width.Bottom &&
		cornerRadius.TopLeft == cornerRadius.TopRight && cornerRadius.TopRight == cornerRadius.BottomLeft && cornerRadius.BottomLeft == cornerRadius.BottomRight {

		r.operationBuilder.BuildBorderOperation(bounds, width.Left, color, cornerRadius.TopLeft)
		return
	}

	// Render each side individually for different widths
	minPt := f32.Pt(float32(bounds.Min.X), float32(bounds.Min.Y))
	maxPt := f32.Pt(float32(bounds.Max.X), float32(bounds.Max.Y))

	// Top border
	if width.Top > 0 {
		r.renderBorderSide(minPt, f32.Pt(maxPt.X, minPt.Y+width.Top), color, cornerRadius.TopLeft, cornerRadius.TopRight, true)
	}

	// Right border
	if width.Right > 0 {
		r.renderBorderSide(f32.Pt(maxPt.X-width.Right, minPt.Y), maxPt, color, cornerRadius.TopRight, cornerRadius.BottomRight, false)
	}

	// Bottom border
	if width.Bottom > 0 {
		r.renderBorderSide(f32.Pt(minPt.X, maxPt.Y-width.Bottom), maxPt, color, cornerRadius.BottomLeft, cornerRadius.BottomRight, true)
	}

	// Left border
	if width.Left > 0 {
		r.renderBorderSide(minPt, f32.Pt(minPt.X+width.Left, maxPt.Y), color, cornerRadius.TopLeft, cornerRadius.BottomLeft, false)
	}
}

// renderBorderSide renders a single border side with corner radius support
func (r *GioRenderer) renderBorderSide(min, max f32.Point, color color.NRGBA, startRadius, endRadius float32, horizontal bool) {
	var path clip.Path
	path.Begin(r.ops)

	if horizontal {
		// Horizontal border (top/bottom)
		if startRadius > 0 {
			// Start with rounded corner
			path.MoveTo(f32.Pt(min.X+startRadius, min.Y))
		} else {
			path.MoveTo(min)
		}

		if endRadius > 0 {
			// End with rounded corner
			path.LineTo(f32.Pt(max.X-endRadius, min.Y))
			path.QuadTo(f32.Pt(max.X, min.Y), f32.Pt(max.X, min.Y+endRadius))
			path.LineTo(f32.Pt(max.X, max.Y))
			path.LineTo(f32.Pt(min.X, max.Y))
		} else {
			path.LineTo(f32.Pt(max.X, min.Y))
			path.LineTo(max)
			path.LineTo(f32.Pt(min.X, max.Y))
		}

		if startRadius > 0 {
			path.LineTo(f32.Pt(min.X, min.Y+startRadius))
			path.QuadTo(min, f32.Pt(min.X+startRadius, min.Y))
		} else {
			path.LineTo(min)
		}
	} else {
		// Vertical border (left/right)
		if startRadius > 0 {
			path.MoveTo(f32.Pt(min.X, min.Y+startRadius))
		} else {
			path.MoveTo(min)
		}

		if endRadius > 0 {
			path.LineTo(f32.Pt(min.X, max.Y-endRadius))
			path.QuadTo(f32.Pt(min.X, max.Y), f32.Pt(min.X+endRadius, max.Y))
			path.LineTo(f32.Pt(max.X, max.Y))
			path.LineTo(f32.Pt(max.X, min.Y))
		} else {
			path.LineTo(f32.Pt(min.X, max.Y))
			path.LineTo(max)
			path.LineTo(f32.Pt(max.X, min.Y))
		}

		if startRadius > 0 {
			path.LineTo(f32.Pt(min.X+startRadius, min.Y))
			path.QuadTo(min, f32.Pt(min.X, min.Y+startRadius))
		} else {
			path.LineTo(min)
		}
	}

	path.Close()

	// Fill the border area
	clipStack := clip.Outline{Path: path.End()}.Op().Push(r.ops)
	paint.ColorOp{Color: color}.Add(r.ops)
	paint.PaintOp{}.Add(r.ops)
	clipStack.Pop()
}

func (r *GioRenderer) RenderClipStart(cmd clay.RenderCommand) error {
	clipData, ok := cmd.Data.(clay.ClipRenderData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"RenderClipStart",
			"Invalid clip start command data",
		)
	}

	defer func() {
		if err := r.errorHandler.RecoverFromPanic("RenderClipStart"); err != nil {
			r.errorHandler.HandleError(err)
		}
	}()

	if r.ops == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"RenderClipStart",
			"Operations context is nil",
		)
	}

	// Check for clip stack overflow protection
	if len(r.clipStack) >= r.maxClipDepth {
		return NewRenderError(
			ErrorTypeClipStackOverflow,
			"RenderClipStart",
			fmt.Sprintf("Maximum clip depth of %d exceeded", r.maxClipDepth),
		)
	}

	// Create clipping bounds from current viewport
	// Note: In a complete implementation, this would use per-element bounds
	// passed from the layout engine, but the current interface limitation
	// requires us to use viewport bounds
	bounds := image.Rectangle{
		Min: image.Point{X: int(r.viewport.X), Y: int(r.viewport.Y)},
		Max: image.Point{
			X: int(r.viewport.X + r.viewport.Width),
			Y: int(r.viewport.Y + r.viewport.Height),
		},
	}

	// Apply directional clipping constraints
	if !clipData.Horizontal {
		// If horizontal clipping is disabled, extend bounds horizontally
		bounds.Min.X = -1000000 // Large negative value
		bounds.Max.X = 1000000  // Large positive value
	}
	if !clipData.Vertical {
		// If vertical clipping is disabled, extend bounds vertically
		bounds.Min.Y = -1000000 // Large negative value
		bounds.Max.Y = 1000000  // Large positive value
	}

	// Create and push clip operation
	clipOp := clip.Rect(bounds).Push(r.ops)

	// Add to clip stack for proper nesting
	r.clipStack = append(r.clipStack, clipOp)

	return nil
}

func (r *GioRenderer) RenderClipEnd(cmd clay.RenderCommand) error {
	if _, ok := cmd.Data.(clay.ClipRenderData); ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"RenderClipEnd",
			"Invalid clip end command data",
		)
	}

	defer func() {
		if err := r.errorHandler.RecoverFromPanic("RenderClipEnd"); err != nil {
			r.errorHandler.HandleError(err)
		}
	}()

	if r.ops == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"RenderClipEnd",
			"Operations context is nil",
		)
	}

	// Check if there are any clip operations to pop
	if len(r.clipStack) == 0 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"RenderClipEnd",
			"No clip operations to end - clip stack is empty",
		)
	}

	// Pop the most recent clip operation from the stack
	clipStackIndex := len(r.clipStack) - 1
	clipOp := r.clipStack[clipStackIndex]

	// Remove from our stack
	r.clipStack = r.clipStack[:clipStackIndex]

	// Pop the Gio clip operation
	clipOp.Pop()

	return nil
}

// CreateComplexClip creates a complex clipping shape with corner radius support
func (r *GioRenderer) CreateComplexClip(bounds image.Rectangle, cornerRadius clay.CornerRadius) (clip.Stack, error) {
	if r.ops == nil {
		return clip.Stack{}, NewRenderError(
			ErrorTypeInvalidInput,
			"CreateComplexClip",
			"Operations context is nil",
		)
	}

	// Check for clip stack overflow protection
	if len(r.clipStack) >= r.maxClipDepth {
		return clip.Stack{}, NewRenderError(
			ErrorTypeClipStackOverflow,
			"CreateComplexClip",
			fmt.Sprintf("Maximum clip depth of %d exceeded", r.maxClipDepth),
		)
	}

	// If no corner radius, use simple rectangle clipping
	if cornerRadius.TopLeft == 0 && cornerRadius.TopRight == 0 &&
		cornerRadius.BottomLeft == 0 && cornerRadius.BottomRight == 0 {
		return clip.Rect(bounds).Push(r.ops), nil
	}

	// Create rounded rectangle path for complex clipping
	var path clip.Path
	path.Begin(r.ops)

	minPt := f32.Pt(float32(bounds.Min.X), float32(bounds.Min.Y))
	maxPt := f32.Pt(float32(bounds.Max.X), float32(bounds.Max.Y))

	// Start from top-left corner (after radius)
	path.MoveTo(f32.Pt(minPt.X+cornerRadius.TopLeft, minPt.Y))

	// Top edge to top-right corner
	if cornerRadius.TopRight > 0 {
		path.LineTo(f32.Pt(maxPt.X-cornerRadius.TopRight, minPt.Y))
		// Top-right corner arc
		path.QuadTo(f32.Pt(maxPt.X, minPt.Y), f32.Pt(maxPt.X, minPt.Y+cornerRadius.TopRight))
	} else {
		path.LineTo(f32.Pt(maxPt.X, minPt.Y))
	}

	// Right edge to bottom-right corner
	if cornerRadius.BottomRight > 0 {
		path.LineTo(f32.Pt(maxPt.X, maxPt.Y-cornerRadius.BottomRight))
		// Bottom-right corner arc
		path.QuadTo(f32.Pt(maxPt.X, maxPt.Y), f32.Pt(maxPt.X-cornerRadius.BottomRight, maxPt.Y))
	} else {
		path.LineTo(f32.Pt(maxPt.X, maxPt.Y))
	}

	// Bottom edge to bottom-left corner
	if cornerRadius.BottomLeft > 0 {
		path.LineTo(f32.Pt(minPt.X+cornerRadius.BottomLeft, maxPt.Y))
		// Bottom-left corner arc
		path.QuadTo(f32.Pt(minPt.X, maxPt.Y), f32.Pt(minPt.X, maxPt.Y-cornerRadius.BottomLeft))
	} else {
		path.LineTo(f32.Pt(minPt.X, maxPt.Y))
	}

	// Left edge to top-left corner
	if cornerRadius.TopLeft > 0 {
		path.LineTo(f32.Pt(minPt.X, minPt.Y+cornerRadius.TopLeft))
		// Top-left corner arc
		path.QuadTo(f32.Pt(minPt.X, minPt.Y), f32.Pt(minPt.X+cornerRadius.TopLeft, minPt.Y))
	} else {
		path.LineTo(f32.Pt(minPt.X, minPt.Y))
	}

	path.Close()

	// Create clip operation from path
	clipOp := clip.Outline{Path: path.End()}.Op().Push(r.ops)

	return clipOp, nil
}

// ClearClipStack safely clears all clip operations (for error recovery)
func (r *GioRenderer) ClearClipStack() error {
	// Pop all remaining clip operations
	for len(r.clipStack) > 0 {
		clipStackIndex := len(r.clipStack) - 1
		clipOp := r.clipStack[clipStackIndex]
		r.clipStack = r.clipStack[:clipStackIndex]
		clipOp.Pop()
	}
	return nil
}

func (r *GioRenderer) RenderCustom(cmd clay.RenderCommand) error {
	return r.customRegistry.ExecuteCustomCommand(r.ops, cmd)
}

// Performance and Resource Management Methods

// RegisterCustomHandler registers a custom command handler
func (r *GioRenderer) RegisterCustomHandler(handler CustomCommandHandler) error {
	return r.customRegistry.RegisterHandler(handler)
}

// UnregisterCustomHandler removes a custom command handler
func (r *GioRenderer) UnregisterCustomHandler(commandID string) error {
	return r.customRegistry.UnregisterHandler(commandID)
}

// GetCacheStats returns current cache statistics
func (r *GioRenderer) GetCacheStats() CacheStats {
	return r.cache.GetStats()
}

// ClearCache clears all cached resources
func (r *GioRenderer) ClearCache() {
	r.cache.Clear()
}

// SetDebugMode enables or disables debug mode for error handling
func (r *GioRenderer) SetDebugMode(enabled bool) {
	r.errorHandler.debugMode = enabled
}

// SetErrorCallback sets a callback function for error notifications
func (r *GioRenderer) SetErrorCallback(callback func(*RenderError)) {
	r.errorHandler.SetErrorCallback(callback)
}

// GetClipStackDepth returns the current clip stack depth
func (r *GioRenderer) GetClipStackDepth() int {
	return len(r.clipStack)
}

// SetMaxClipDepth sets the maximum allowed clip stack depth
func (r *GioRenderer) SetMaxClipDepth(maxDepth int) {
	r.maxClipDepth = maxDepth
}

// BatchRenderOperations enables batched rendering for performance
func (r *GioRenderer) BatchRenderOperations(enable bool) {
	if enable {
		// Initialize batch operations if not already done
		if r.batchOperations == nil {
			r.batchOperations = NewBatchOperations()
		}
	}
}

// ExecuteBatchedOperations executes all batched operations
func (r *GioRenderer) ExecuteBatchedOperations() error {
	if r.batchOperations != nil {
		r.batchOperations.Execute(r.ops)
	}
	return nil
}

// RenderGradient renders a linear gradient (advanced color feature)
func (r *GioRenderer) RenderGradient(startColor, endColor clay.Color, vertical bool) error {
	defer func() {
		if err := r.errorHandler.RecoverFromPanic("RenderGradient"); err != nil {
			r.errorHandler.HandleError(err)
		}
	}()

	// Validate color values
	if err := ValidateColor(startColor.R, startColor.G, startColor.B, startColor.A); err != nil {
		r.errorHandler.HandleError(err.(*RenderError))
		return err
	}
	if err := ValidateColor(endColor.R, endColor.G, endColor.B, endColor.A); err != nil {
		r.errorHandler.HandleError(err.(*RenderError))
		return err
	}

	// Convert Clay colors to Gio colors
	gioStartColor := r.cache.GetOrCreateColor(startColor.R, startColor.G, startColor.B, startColor.A)
	gioEndColor := r.cache.GetOrCreateColor(endColor.R, endColor.G, endColor.B, endColor.A)

	// Create gradient bounds from viewport
	bounds := image.Rectangle{
		Min: image.Point{X: int(r.viewport.X), Y: int(r.viewport.Y)},
		Max: image.Point{
			X: int(r.viewport.X + r.viewport.Width),
			Y: int(r.viewport.Y + r.viewport.Height),
		},
	}

	// Use operation builder to render gradient
	r.operationBuilder.BuildGradientOperation(bounds, gioStartColor, gioEndColor, vertical)

	return nil
}

// ValidateRenderer performs comprehensive validation of renderer state
func (r *GioRenderer) ValidateRenderer() error {
	if r.ops == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"ValidateRenderer",
			"Operations context is nil",
		)
	}

	if r.cache == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"ValidateRenderer",
			"Resource cache is nil",
		)
	}

	if r.customRegistry == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"ValidateRenderer",
			"Custom command registry is nil",
		)
	}

	if r.errorHandler == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"ValidateRenderer",
			"Error handler is nil",
		)
	}

	if r.operationBuilder == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"ValidateRenderer",
			"Operation builder is nil",
		)
	}

	return nil
}

// Cleanup performs cleanup of resources and resets state
func (r *GioRenderer) Cleanup() error {
	// Clear all clips from stack
	for len(r.clipStack) > 0 {
		lastIndex := len(r.clipStack) - 1
		clipStack := r.clipStack[lastIndex]
		r.clipStack = r.clipStack[:lastIndex]
		clipStack.Pop()
	}

	// Clear cache
	if r.cache != nil {
		r.cache.Clear()
	}

	// Reset viewport
	r.viewport = clay.BoundingBox{}

	return nil
}

// renderClipEndWithBounds ends clipping using bounds from RenderCommand
func (r *GioRenderer) renderClipEndWithBounds(bounds clay.BoundingBox, cmd clay.RenderCommand) error {

	if _, ok := cmd.Data.(clay.ClipRenderData); ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"renderClipEndWithBounds",
			"Invalid clip end command data",
		)
	}

	if len(r.clipStack) == 0 {
		return NewRenderError(
			ErrorTypeInvalidState,
			"renderClipEndWithBounds",
			"No clip region to end - clip stack is empty",
		)
	}

	// Pop the last clip from the stack
	lastIndex := len(r.clipStack) - 1
	clipStack := r.clipStack[lastIndex]
	r.clipStack = r.clipStack[:lastIndex]
	clipStack.Pop()

	return nil
}

// renderCustomWithBounds renders a custom command using bounds from RenderCommand
func (r *GioRenderer) renderCustomWithBounds(bounds clay.BoundingBox, cmd clay.RenderCommand) error {
	// if _, ok := cmd.Data.(clay.CustomRenderData); ok {
	// 	return NewRenderError(
	// 		ErrorTypeInvalidData,
	// 		"renderCustomWithBounds",
	// 		"Invalid custom command data",
	// 	)
	// }

	if r.customRegistry == nil {
		return NewRenderError(
			ErrorTypeInvalidState,
			"renderCustomWithBounds",
			"Custom command registry is not initialized",
		)
	}

	// For now, we'll need to determine the command type from the CustomData
	// This is a simplified implementation - in a real scenario, you'd need
	// a way to identify which handler to use based on the CustomData

	// TODO: Implement proper custom command type identification
	// For now, just return an error indicating custom commands need proper setup
	return NewRenderError(
		ErrorTypeUnsupportedOperation,
		"renderCustomWithBounds",
		"Custom command rendering requires proper handler registration - not yet fully implemented",
	)
}

// GetRendererInfo returns information about the renderer capabilities
func (r *GioRenderer) GetRendererInfo() RendererInfo {
	return RendererInfo{
		Name:               "Gio UI Renderer",
		Version:            "1.2.0",
		SupportsImages:     true,
		SupportsBorders:    true,
		SupportsClipping:   true,
		SupportsCustom:     true,
		SupportsGradients:  true,
		SupportsCaching:    true,
		MaxClipDepth:       r.maxClipDepth,
		CacheSize:          r.cache.maxSize,
		CustomHandlerCount: len(r.customRegistry.handlers),
	}
}

// GetFontManager returns the FontManager instance for font registration and management
func (r *GioRenderer) GetFontManager() *FontManager {
	return r.fontManager
}

// RendererInfo provides information about renderer capabilities
type RendererInfo struct {
	Name               string
	Version            string
	SupportsImages     bool
	SupportsBorders    bool
	SupportsClipping   bool
	SupportsCustom     bool
	SupportsGradients  bool
	SupportsCaching    bool
	MaxClipDepth       int
	CacheSize          int
	CustomHandlerCount int
}
