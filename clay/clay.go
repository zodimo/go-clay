// Package clay provides a high-performance, render engine agnostic UI layout library for Go.
//
// Go-Clay is inspired by the original Clay C library and provides:
// - Microsecond layout performance
// - Flexbox-like layout model
// - Render engine agnostic design
// - Memory efficient arena-based allocation
// - Declarative API similar to React
//
// Example usage:
//
//	engine := clay.NewLayoutEngine()
//	engine.BeginLayout()
//
//	clay.Container("main", clay.ElementConfig{
//		Layout: clay.LayoutConfig{
//			Sizing: clay.Sizing{
//				Width:  clay.SizingGrow(0),
//				Height: clay.SizingGrow(0),
//			},
//			Padding: clay.PaddingAll(16),
//		},
//		BackgroundColor: clay.Color{R: 0.9, G: 0.9, B: 0.9, A: 1.0},
//	}).Text("Hello, World!", clay.TextConfig{
//		FontSize: 24,
//		Color:    clay.Color{R: 0, G: 0, B: 0, A: 1.0},
//	})
//
//	commands := engine.EndLayout()
//	renderer.Render(commands)
package clay

import (
	"hash/fnv"
	"sync"
)

// Version information
const (
	Version = "0.1.0"
)

// ElementID represents a unique identifier for UI elements
type ElementID uint32

// ID creates an ElementID from a string
func ID(str string) ElementID {
	h := fnv.New32a()
	h.Write([]byte(str))
	return ElementID(h.Sum32())
}

// IDWithIndex creates an ElementID from a string and index
func IDWithIndex(str string, index int) ElementID {
	h := fnv.New32a()
	h.Write([]byte(str))
	h.Write([]byte{byte(index), byte(index >> 8), byte(index >> 16), byte(index >> 24)})
	return ElementID(h.Sum32())
}

// Dimensions represents size in pixels
type Dimensions struct {
	Width, Height float32
}

// Vector2 represents a 2D position or offset
type Vector2 struct {
	X, Y float32
}

// Color represents RGBA color with values 0.0-1.0
type Color struct {
	R, G, B, A float32
}

// ColorRGB creates a color with RGB values (0.0-1.0) and alpha 1.0
func ColorRGB(r, g, b float32) Color {
	return Color{R: r, G: g, B: b, A: 1.0}
}

// ColorRGBA creates a color with RGBA values (0.0-1.0)
func ColorRGBA(r, g, b, a float32) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// BoundingBox represents rectangle bounds for positioning and clipping
type BoundingBox struct {
	X, Y, Width, Height float32
}

// SizingType controls how elements take up space
type SizingType int

const (
	SizingFitToContent             SizingType = iota // Wrap to content
	SizingGrowToFillAvailableSpace                   // Fill available space
	SizingPercentOfParent                            // Percentage of parent
	SizingFixedPixelSize                             // Fixed pixel size
)

// SizingAxis controls sizing along one axis
type SizingAxis struct {
	Type     SizingType
	Min, Max float32 // Min/max constraints
	Percent  float32 // Percentage (0.0-1.0)
}

// Sizing controls element sizing
type Sizing struct {
	Width, Height SizingAxis
}

// SizingFit creates a fit sizing axis
func SizingFit() SizingAxis {
	return SizingAxis{Type: SizingFitToContent}
}

// SizingGrow creates a grow sizing axis with weight
func SizingGrow(weight float32) SizingAxis {
	return SizingAxis{Type: SizingGrowToFillAvailableSpace, Min: weight}
}

// SizingPercent creates a percentage sizing axis
func SizingPercent(percent float32) SizingAxis {
	return SizingAxis{Type: SizingPercentOfParent, Percent: percent}
}

// SizingFixed creates a fixed sizing axis
func SizingFixed(size float32) SizingAxis {
	return SizingAxis{Type: SizingFixedPixelSize, Min: size, Max: size}
}

// LayoutDirection controls child arrangement
type LayoutDirection int

const (
	LeftToRight LayoutDirection = iota
	TopToBottom
)

// AlignmentX controls horizontal alignment
type AlignmentX int

const (
	AlignXLeft AlignmentX = iota
	AlignXCenter
	AlignXRight
)

// AlignmentY controls vertical alignment
type AlignmentY int

const (
	AlignYTop AlignmentY = iota
	AlignYCenter
	AlignYBottom
)

// ChildAlignment controls child positioning
type ChildAlignment struct {
	X AlignmentX
	Y AlignmentY
}

// Padding represents space between element border and children
type Padding struct {
	Left, Right, Top, Bottom float32
}

// PaddingAll creates padding with same value on all sides
func PaddingAll(value float32) Padding {
	return Padding{Left: value, Right: value, Top: value, Bottom: value}
}

// PaddingHorizontal creates padding with different left/right values
func PaddingHorizontal(left, right float32) Padding {
	return Padding{Left: left, Right: right}
}

// PaddingVertical creates padding with different top/bottom values
func PaddingVertical(top, bottom float32) Padding {
	return Padding{Top: top, Bottom: bottom}
}

// LayoutConfig controls element layout properties
type LayoutConfig struct {
	Sizing         Sizing
	Padding        Padding
	ChildGap       float32
	Direction      LayoutDirection
	ChildAlignment ChildAlignment
}

// CornerRadius controls corner rounding
type CornerRadius struct {
	TopLeft, TopRight, BottomLeft, BottomRight float32
}

// CornerRadiusAll creates corner radius with same value on all corners
func CornerRadiusAll(radius float32) CornerRadius {
	return CornerRadius{
		TopLeft:     radius,
		TopRight:    radius,
		BottomLeft:  radius,
		BottomRight: radius,
	}
}

// BorderWidth controls border thickness
type BorderWidth struct {
	Left, Right, Top, Bottom float32
}

// BorderWidthAll creates border width with same value on all sides
func BorderWidthAll(width float32) BorderWidth {
	return BorderWidth{Left: width, Right: width, Top: width, Bottom: width}
}

// BorderConfig controls border styling
type BorderConfig struct {
	Width BorderWidth
	Color Color
}

// TextWrapMode controls text wrapping behavior
type TextWrapMode int

const (
	WrapWords    TextWrapMode = iota // Break on whitespace
	WrapNewlines                     // Break only on newlines
	WrapNone                         // No wrapping
)

// TextAlignment controls text alignment
type TextAlignment int

const (
	TextAlignLeft TextAlignment = iota
	TextAlignCenter
	TextAlignRight
)

// TextConfig controls text rendering
type TextConfig struct {
	FontSize      float32
	Color         Color
	FontID        uint16
	LineHeight    float32
	LetterSpacing float32
	WrapMode      TextWrapMode
	Alignment     TextAlignment
}

// ImageConfig controls image rendering
type ImageConfig struct {
	ImageData interface{}
	TintColor Color
}

// FloatingConfig controls floating element behavior
type FloatingConfig struct {
	AttachTo    ElementID
	ZIndex      int16
	Offset      Vector2
	AttachPoint AttachPoint
}

// AttachPoint controls floating element attachment
type AttachPoint int

const (
	AttachPointLeftTop AttachPoint = iota
	AttachPointLeftCenter
	AttachPointLeftBottom
	AttachPointCenterTop
	AttachPointCenterCenter
	AttachPointCenterBottom
	AttachPointRightTop
	AttachPointRightCenter
	AttachPointRightBottom
)

// ClipConfig controls clipping behavior
type ClipConfig struct {
	Horizontal, Vertical bool
	ChildOffset          Vector2
}

// ElementDeclaration represents an element's configuration
type ElementDeclaration struct {
	ID              ElementID
	Layout          LayoutConfig
	BackgroundColor Color
	Border          *BorderConfig
	CornerRadius    CornerRadius
	Text            *TextConfig
	Image           *ImageConfig
	Floating        *FloatingConfig
	Clip            *ClipConfig
	AspectRatio     float32
	UserData        interface{}
}

// ElementConfig is a shorthand for ElementDeclaration
type ElementConfig = ElementDeclaration

// RenderCommandType represents the type of render command
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

// RenderCommand represents a single render operation
type RenderCommand struct {
	BoundingBox BoundingBox
	CommandType RenderCommandType
	ZIndex      int16
	ID          ElementID
	Data        interface{} // Command-specific data
}

// RectangleCommand represents rectangle rendering data
type RectangleCommand struct {
	Color        Color
	CornerRadius CornerRadius
}

// TextCommand represents text rendering data
type TextCommand struct {
	Text          string
	FontID        uint16
	FontSize      float32
	Color         Color
	LineHeight    float32
	LetterSpacing float32
	Alignment     TextAlignment
}

// ImageCommand represents image rendering data
type ImageCommand struct {
	ImageData    interface{}
	TintColor    Color
	CornerRadius CornerRadius
}

// BorderCommand represents border rendering data
type BorderCommand struct {
	Color        Color
	Width        BorderWidth
	CornerRadius CornerRadius
}

// ClipStartCommand represents clipping start data
type ClipStartCommand struct {
	Horizontal, Vertical bool
}

// ClipEndCommand represents clipping end data
type ClipEndCommand struct {
	// No additional data needed
}

// CustomCommand represents custom rendering data
type CustomCommand struct {
	CustomData interface{}
}

// LayoutEngine manages the layout computation
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

	// Debug
	SetDebugMode(enabled bool)
	GetStats() LayoutStats
}

// LayoutStats provides performance metrics
type LayoutStats struct {
	ElementCount   int
	RenderCommands int
	LayoutTime     int64 // nanoseconds
	MemoryUsed     int   // bytes
}

// TextMeasurer interface for text measurement
type TextMeasurer interface {
	MeasureText(text string, config TextConfig) Dimensions
	GetTextMetrics(fontID uint16, fontSize float32) TextMetrics
}

// TextMetrics provides text measurement information
type TextMetrics struct {
	Ascent  float32
	Descent float32
	Height  float32
}

// Arena provides arena-based memory allocation
type Arena struct {
	memory   []byte
	offset   int
	capacity int
	mutex    sync.Mutex
}

// NewArena creates a new arena with the specified capacity
func NewArena(capacity int) *Arena {
	return &Arena{
		memory:   make([]byte, capacity),
		capacity: capacity,
	}
}

// Allocate allocates memory from the arena
func (a *Arena) Allocate(size int) []byte {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.offset+size > a.capacity {
		return nil // Arena full
	}

	ptr := a.memory[a.offset : a.offset+size]
	a.offset += size
	return ptr
}

// Reset resets the arena for reuse
func (a *Arena) Reset() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.offset = 0
}

// Used returns the amount of memory used
func (a *Arena) Used() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.offset
}

// Available returns the amount of memory available
func (a *Arena) Available() int {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.capacity - a.offset
}

// ContainerBuilder provides a fluent API for building containers
type ContainerBuilder struct {
	engine *layoutEngine
	id     ElementID
	config ElementDeclaration
}

// Text adds a text element to the container
func (c *ContainerBuilder) Text(text string, config TextConfig) *ContainerBuilder {
	c.engine.OpenElement(c.engine.generateID(), ElementDeclaration{
		Text: &config,
	})
	c.engine.CloseElement()
	return c
}

// Image adds an image element to the container
func (c *ContainerBuilder) Image(config ImageConfig) *ContainerBuilder {
	c.engine.OpenElement(c.engine.generateID(), ElementDeclaration{
		Image: &config,
	})
	c.engine.CloseElement()
	return c
}

// Container adds a child container
func (c *ContainerBuilder) Container(id ElementID, config ElementDeclaration) *ContainerBuilder {
	c.engine.OpenElement(id, config)
	return &ContainerBuilder{
		engine: c.engine,
		id:     id,
		config: config,
	}
}

// End closes the container
func (c *ContainerBuilder) End() {
	c.engine.CloseElement()
}

// Container creates a new container with the given ID and configuration
func Container(id ElementID, config ElementDeclaration) *ContainerBuilder {
	engine := getCurrentEngine()
	if engine == nil {
		panic("No active layout engine. Call BeginLayout() first.")
	}

	engine.OpenElement(id, config)
	return &ContainerBuilder{
		engine: engine,
		id:     id,
		config: config,
	}
}

// Text creates a text element
func Text(text string, config TextConfig) {
	engine := getCurrentEngine()
	if engine == nil {
		panic("No active layout engine. Call BeginLayout() first.")
	}

	engine.OpenElement(engine.generateID(), ElementDeclaration{
		Text: &config,
	})
	engine.CloseElement()
}

// Image creates an image element
func Image(config ImageConfig) {
	engine := getCurrentEngine()
	if engine == nil {
		panic("No active layout engine. Call BeginLayout() first.")
	}

	engine.OpenElement(engine.generateID(), ElementDeclaration{
		Image: &config,
	})
	engine.CloseElement()
}

// Global engine state
var (
	currentEngine *layoutEngine
	engineMutex   sync.RWMutex
)

// getCurrentEngine returns the current layout engine
func getCurrentEngine() *layoutEngine {
	engineMutex.RLock()
	defer engineMutex.RUnlock()
	return currentEngine
}

// setCurrentEngine sets the current layout engine
func setCurrentEngine(engine *layoutEngine) {
	engineMutex.Lock()
	defer engineMutex.Unlock()
	currentEngine = engine
}
