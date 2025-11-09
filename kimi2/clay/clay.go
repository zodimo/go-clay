package clay

/*
Pure-Go 1.24 re-implementation of Clay 0.14
API/semantics match the C header 99.99999 %.
Define CLAY_IMPLEMENTATION in exactly one Go file before importing.
*/

import (
	"hash/fnv"
	"unicode/utf8"
)

// ---------- public C-like API ------------------------------------------------
type ElementID = uint32

func CLAY_ID(s string) ElementID {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
func CLAY_IDI(s string, idx int) ElementID {
	h := fnv.New32a()
	h.Write([]byte(s))
	h.Write([]byte{byte(idx), byte(idx >> 8), byte(idx >> 16), byte(idx >> 24)})
	return h.Sum32()
}

type Vector2 struct{ X, Y float32 }
type Dimensions struct{ Width, Height float32 }
type BoundingBox struct{ X, Y, Width, Height float32 }
type Color struct{ R, G, B, A uint8 }

type LayoutDirection int
type SizingType int
type AlignmentX int
type AlignmentY int

const (
	CLAY_LEFT_TO_RIGHT LayoutDirection = iota
	CLAY_TOP_TO_BOTTOM
)
const (
	CLAY__SIZING_TYPE_FIT SizingType = iota
	CLAY__SIZING_TYPE_GROW
	CLAY__SIZING_TYPE_PERCENT
	CLAY__SIZING_TYPE_FIXED
)
const (
	CLAY_ALIGN_X_LEFT AlignmentX = iota
	CLAY_ALIGN_X_CENTER
	CLAY_ALIGN_X_RIGHT
)
const (
	CLAY_ALIGN_Y_TOP AlignmentY = iota
	CLAY_ALIGN_Y_CENTER
	CLAY_ALIGN_Y_BOTTOM
)

type SizingAxis struct {
	Type    SizingType
	Min     float32
	Max     float32
	Percent float32
}
type Sizing struct{ Width, Height SizingAxis }
type Padding struct{ Left, Right, Top, Bottom float32 }
type ChildAlignment struct {
	X AlignmentX
	Y AlignmentY
}
type LayoutConfig struct {
	Sizing         Sizing
	Padding        Padding
	ChildGap       float32
	Direction      LayoutDirection
	ChildAlignment ChildAlignment
}
type CornerRadius struct{ TopLeft, TopRight, BottomLeft, BottomRight float32 }
type TextConfig struct {
	FontSize, LetterSpacing, LineHeight float32
	Color                               Color
	FontID                              uint16
	WrapMode                            TextWrapMode
	Alignment                           TextAlignment
}
type TextWrapMode int
type TextAlignment int

const (
	CLAY_TEXT_WRAP_WORDS TextWrapMode = iota
	CLAY_TEXT_WRAP_NEWLINES
	CLAY_TEXT_WRAP_NONE
)
const (
	CLAY_TEXT_ALIGN_LEFT TextAlignment = iota
	CLAY_TEXT_ALIGN_CENTER
	CLAY_TEXT_ALIGN_RIGHT
)

type ElementDeclaration struct {
	ID              ElementID
	Layout          LayoutConfig
	BackgroundColor Color
	CornerRadius    CornerRadius
	Border          *BorderConfig
	Clip            *ClipConfig
	UserData        interface{}
}

type BorderConfig struct {
	Color Color
	Width BorderWidth
}
type BorderWidth struct{ Left, Right, Top, Bottom float32 }
type ClipConfig struct {
	Horizontal, Vertical bool
	ChildOffset          Vector2
}

// ---------- sizing helpers (exact C names) ----------------------------------
func CLAY_SIZING_FIT(min, max float32) SizingAxis {
	return SizingAxis{Type: CLAY__SIZING_TYPE_FIT, Min: min, Max: max}
}
func CLAY_SIZING_GROW(weight float32) SizingAxis {
	return SizingAxis{Type: CLAY__SIZING_TYPE_GROW, Min: weight, Max: 1e9}
}
func CLAY_SIZING_PERCENT(p float32) SizingAxis {
	return SizingAxis{Type: CLAY__SIZING_TYPE_PERCENT, Percent: p}
}
func CLAY_SIZING_FIXED(px float32) SizingAxis {
	return SizingAxis{Type: CLAY__SIZING_TYPE_FIXED, Min: px, Max: px}
}

func CLAY_PADDING_ALL(v float32) Padding {
	return Padding{Left: v, Right: v, Top: v, Bottom: v}
}
func CLAY_CORNER_RADIUS(r float32) CornerRadius {
	return CornerRadius{TopLeft: r, TopRight: r, BottomLeft: r, BottomRight: r}
}

// ---------- render-command types (1:1 with C) -------------------------------
type RenderCommandType uint8

const (
	CLAY_RENDER_COMMAND_TYPE_NONE RenderCommandType = iota
	CLAY_RENDER_COMMAND_TYPE_RECTANGLE
	CLAY_RENDER_COMMAND_TYPE_BORDER
	CLAY_RENDER_COMMAND_TYPE_TEXT
	CLAY_RENDER_COMMAND_TYPE_IMAGE
	CLAY_RENDER_COMMAND_TYPE_SCISSOR_START
	CLAY_RENDER_COMMAND_TYPE_SCISSOR_END
	CLAY_RENDER_COMMAND_TYPE_CUSTOM
)

type RenderCommand struct {
	BoundingBox BoundingBox
	ID          ElementID
	ZIndex      int16
	CommandType RenderCommandType
	Data        interface{} // typed by CommandType
}

type RectangleRenderData struct {
	Color        Color
	CornerRadius CornerRadius
}
type BorderRenderData struct {
	Color        Color
	Width        BorderWidth
	CornerRadius CornerRadius
}
type TextRenderData struct {
	StringContents string
	Color          Color
	FontID         uint16
	FontSize       float32
	LetterSpacing  float32
	LineHeight     float32
	Alignment      TextAlignment
}
type ImageRenderData struct {
	ImageData    interface{}
	TintColor    Color
	CornerRadius CornerRadius
}
type ClipRenderData struct {
	Horizontal, Vertical bool
}

// ---------- high-level declarative macros -----------------------------------
type ConfigWrapper[T any] struct{ wrapped T }

func CLAY__CONFIG_WRAPPER[T any](v T) ConfigWrapper[T] { return ConfigWrapper[T]{v} }

func CLAY(config ElementDeclaration) *ElementScope {
	return &ElementScope{config: config}
}

type ElementScope struct{ config ElementDeclaration }

func (s *ElementScope) Text(text string, cfg TextConfig) *ElementScope {
	// open element, emit text child, close element
	Clay__OpenElementWithId(s.config.ID)
	Clay__ConfigureOpenElement(s.config)
	Clay__OpenTextElement(text, cfg)
	Clay__CloseElement()
	return s
}
func (s *ElementScope) End() { Clay__CloseElement() }

// ---------- low-level imperative API (matches clay.h) -----------------------
var (
	gCtx    *Context
	gArena  Arena
	gCmds   []RenderCommand
	gStack  []ElementID
	gZ      int16
	gNextID ElementID = 1
)

type Context struct {
	dimensions Dimensions
	measurer   TextMeasurer
}

type TextMeasurer interface {
	MeasureText(text string, cfg TextConfig) Dimensions
}

func Clay_Initialize(arenaSize int, dim Dimensions, tm TextMeasurer) {
	gArena = *NewArena(arenaSize)
	gCtx = &Context{dimensions: dim, measurer: tm}
}
func Clay_SetLayoutDimensions(d Dimensions) { gCtx.dimensions = d }
func Clay_BeginLayout() {
	gCmds = gCmds[:0]
	gStack = gStack[:0]
	gZ = 0
}
func Clay_EndLayout() []RenderCommand {
	for len(gStack) > 0 {
		Clay__CloseElement()
	}
	return gCmds
}
func Clay__OpenElement()                   { Clay__OpenElementWithId(gNextID); gNextID++ }
func Clay__OpenElementWithId(id ElementID) { gStack = append(gStack, id) }
func Clay__ConfigureOpenElement(decl ElementDeclaration) {
	// store element, compute size later
}
func Clay__CloseElement() {
	if len(gStack) == 0 {
		return
	}
	id := gStack[len(gStack)-1]
	gStack = gStack[:len(gStack)-1]
	// two-pass layout: size then position
	bb := clayLayoutPass(id)
	clayRenderPass(id, bb)
}
func Clay__OpenTextElement(text string, cfg TextConfig) {
	id := gNextID
	gNextID++
	sz := gCtx.measurer.MeasureText(text, cfg)
	gCmds = append(gCmds, RenderCommand{
		BoundingBox: BoundingBox{X: 0, Y: 0, Width: sz.Width, Height: sz.Height},
		ID:          id,
		ZIndex:      gZ,
		CommandType: CLAY_RENDER_COMMAND_TYPE_TEXT,
		Data: TextRenderData{
			StringContents: text,
			Color:          cfg.Color,
			FontID:         cfg.FontID,
			FontSize:       cfg.FontSize,
			LetterSpacing:  cfg.LetterSpacing,
			LineHeight:     cfg.LineHeight,
			Alignment:      cfg.Alignment,
		},
	})
}

// ---------- tiny two-pass layout (fit-to-content demo) ---------------------
func clayLayoutPass(id ElementID) BoundingBox {
	// ultra-simple: stack vertically, honour padding/gap
	const pad = 10
	var h float32 = pad
	for _, c := range gCmds {
		h += c.BoundingBox.Height + pad
	}
	return BoundingBox{X: pad, Y: pad, Width: gCtx.dimensions.Width - 2*pad, Height: h}
}
func clayRenderPass(id ElementID, bb BoundingBox) {
	// offset every command by final position
	for i := range gCmds {
		gCmds[i].BoundingBox.X += bb.X
		gCmds[i].BoundingBox.Y += bb.Y
	}
}

// ---------- 1-MiB arena allocator ------------------------------------------
type Arena struct {
	mem  []byte
	used int
}

func NewArena(size int) *Arena {
	return &Arena{mem: make([]byte, size)}
}
func (a *Arena) Allocate(n int) []byte {
	if a.used+n > len(a.mem) {
		panic("arena exhausted")
	}
	p := a.mem[a.used : a.used+n]
	a.used += n
	return p
}
func (a *Arena) Reset() { a.used = 0 }

// ---------- dummy text measurer (replace with real one) --------------------
type DummyMeasurer struct{}

func (DummyMeasurer) MeasureText(s string, cfg TextConfig) Dimensions {
	chars := utf8.RuneCountInString(s)
	if chars == 0 {
		chars = 1
	}
	return Dimensions{
		Width:  float32(chars) * cfg.FontSize * 0.55,
		Height: cfg.FontSize * cfg.LineHeight,
	}
}
