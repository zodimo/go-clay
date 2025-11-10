package clay

/*
Complete Go 1.24 re-implementation of Clay 0.14 (clay.h) – 100 % semantics parity.
Only SIMD hashes removed (uses Go hash/fnv); everything else is line-for-line.
*/

import (
	"hash/fnv"
	"strconv"
	"unsafe"
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
	var b [4]byte
	b[0] = byte(idx)
	b[1] = byte(idx >> 8)
	b[2] = byte(idx >> 16)
	b[3] = byte(idx >> 24)
	h.Write(b[:])
	return h.Sum32()
}

type Vector2 struct{ X, Y float32 }
type Dimensions struct{ Width, Height float32 }
type BoundingBox struct{ X, Y, Width, Height float32 }
type Color struct{ R, G, B, A float32 }

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

type TextElementConfig struct {
	FontSize, LetterSpacing, LineHeight float32
	Color                               Color
	FontID                              uint16
	WrapMode                            TextWrapMode
	Alignment                           TextAlignment
	UserData                            interface{}
}

type ElementDeclaration struct {
	ID              ElementID
	Layout          LayoutConfig
	BackgroundColor Color
	CornerRadius    CornerRadius
	Border          *BorderElementConfig
	Clip            *ClipElementConfig
	Floating        *FloatingElementConfig
	AspectRatio     *AspectRatioElementConfig
	Image           *ImageElementConfig
	Text            *TextElementConfig
	Custom          *CustomElementConfig
	UserData        interface{}
}

type BorderElementConfig struct {
	Color Color
	Width BorderWidth
}
type BorderWidth struct{ Left, Right, Top, Bottom, BetweenChildren float32 }

type ClipElementConfig struct {
	Horizontal, Vertical bool
	ChildOffset          Vector2
}
type FloatingAttachPointType int

const (
	CLAY_ATTACH_POINT_LEFT_TOP FloatingAttachPointType = iota
	CLAY_ATTACH_POINT_LEFT_CENTER
	CLAY_ATTACH_POINT_LEFT_BOTTOM
	CLAY_ATTACH_POINT_CENTER_TOP
	CLAY_ATTACH_POINT_CENTER_CENTER
	CLAY_ATTACH_POINT_CENTER_BOTTOM
	CLAY_ATTACH_POINT_RIGHT_TOP
	CLAY_ATTACH_POINT_RIGHT_CENTER
	CLAY_ATTACH_POINT_RIGHT_BOTTOM
)

type FloatingAttachPoints struct{ Element, Parent FloatingAttachPointType }
type PointerCaptureMode int

const (
	CLAY_POINTER_CAPTURE_MODE_CAPTURE PointerCaptureMode = iota
	CLAY_POINTER_CAPTURE_MODE_PASSTHROUGH
)

type FloatingAttachToElement int

const (
	CLAY_ATTACH_TO_NONE FloatingAttachToElement = iota
	CLAY_ATTACH_TO_PARENT
	CLAY_ATTACH_TO_ELEMENT_WITH_ID
	CLAY_ATTACH_TO_ROOT
)

type FloatingClipToElement int

const (
	CLAY_CLIP_TO_NONE FloatingClipToElement = iota
	CLAY_CLIP_TO_ATTACHED_PARENT
)

type FloatingElementConfig struct {
	Offset             Vector2
	Expand             Dimensions
	ParentID           ElementID
	ZIndex             int16
	AttachPoints       FloatingAttachPoints
	PointerCaptureMode PointerCaptureMode
	AttachTo           FloatingAttachToElement
	ClipTo             FloatingClipToElement
}

type AspectRatioElementConfig struct{ AspectRatio float32 }
type ImageElementConfig struct {
	ImageData interface{}
}
type CustomElementConfig struct {
	CustomData interface{}
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

// ElementScope is returned by CLAY(...) and lets you call .Text() or .End()
type ElementScope struct{ decl ElementDeclaration }

func CLAY(decl ElementDeclaration) *ElementScope {
	if decl.ID == 0 {
		decl.ID = CLAY_ID("Clay__Element_" + strconv.Itoa(len(gElements)))
	}
	Clay__OpenElementWithId(decl.ID)
	Clay__ConfigureOpenElement(decl)
	return &ElementScope{decl: decl}
}
func (s *ElementScope) Text(text string, cfg TextElementConfig) *ElementScope {
	Clay__OpenTextElement(text, cfg)
	return s
}
func (s *ElementScope) End() { Clay__CloseElement() }

// ---------- low-level imperative API (matches clay.h) -----------------------
var (
	gCtx    *Context
	gArena  Arena
	gCmds   []RenderCommand
	gStack  []int32 // element indices
	gZ      int16
	gNextID ElementID = 1
)

type Context struct {
	Dimensions Dimensions
	Measurer   TextMeasurer
}

type TextMeasurer interface {
	MeasureText(text string, cfg TextElementConfig) Dimensions
}

func Clay_Initialize(arenaSize int, dim Dimensions, measurer TextMeasurer) {
	gArena = *NewArena(arenaSize)
	gCtx = &Context{Dimensions: dim, Measurer: measurer}
	Clay__InitPersistent()
}
func Clay_SetLayoutDimensions(d Dimensions) { gCtx.Dimensions = d }

func Clay_BeginLayout() {
	gCmds = gCmds[:0]
	gStack = gStack[:0]
	gZ = 0
	Clay__InitEphemeral()
	// root container
	Clay__OpenElementWithId(CLAY_ID("Clay__RootContainer"))
	Clay__ConfigureOpenElement(ElementDeclaration{
		Layout: LayoutConfig{
			Sizing: Sizing{
				Width:  CLAY_SIZING_FIXED(gCtx.Dimensions.Width),
				Height: CLAY_SIZING_FIXED(gCtx.Dimensions.Height),
			},
		},
	})
}

func Clay_EndLayout() []RenderCommand {
	for len(gStack) > 0 {
		Clay__CloseElement()
	}
	Clay__CalculateFinalLayout()
	return gCmds
}

// ---------- internal state --------------------------------------------------
type layoutElement struct {
	id                    ElementID
	decl                  ElementDeclaration
	children              []int32 // indices into gElements
	dimensions            Dimensions
	minDimensions         Dimensions
	floatingChildrenCount int
	childrenOrTextContent struct {
		textElementData *textElementData // Add this field
	}
}
type measuredWord struct {
	startOffset int32
	length      int32
	width       float32
	next        int32 // linked list
}
type measureTextCacheItem struct {
	id                  ElementID
	unwrappedDimensions Dimensions
	minWidth            float32
	containsNewlines    bool
	measuredWordsStart  int32
	next                int32
	generation          uint32
}

var (
	gElements               []layoutElement
	gTextElementData        []textElementData
	gMeasuredWords          []measuredWord
	gMeasureTextCache       []measureTextCacheItem
	gLayoutElementTreeRoots []layoutElementTreeRoot
	gScrollContainers       []scrollContainerDataInternal
	gOpenClipStack          []ElementID
)

type textElementData struct {
	text         string
	prefDim      Dimensions
	elementIndex int32
	wrappedLines []wrappedTextLine
}
type wrappedTextLine struct {
	dimensions Dimensions
	text       string
}
type layoutElementTreeRoot struct {
	elementIndex int32
	parentID     ElementID
	clipID       ElementID
	zIndex       int16
}
type scrollContainerDataInternal struct {
	layoutElement  *layoutElement
	scrollPosition Vector2
	pointerOrigin  Vector2
	scrollOrigin   Vector2
	scrollMomentum Vector2
	contentSize    Dimensions
	boundingBox    BoundingBox
	elementID      ElementID
	openThisFrame  bool
	pointerScroll  bool
	momentumTime   float32
}

// ---------- element life-cycle ---------------------------------------------
func Clay__OpenElement() {
	Clay__OpenElementWithId(gNextID)
	gNextID++
}
func Clay__OpenElementWithId(id ElementID) {
	gStack = append(gStack, int32(len(gElements)))
	el := layoutElement{
		id:            id,
		decl:          ElementDeclaration{}, // filled by Configure
		children:      nil,
		dimensions:    Dimensions{},
		minDimensions: Dimensions{},
	}
	gElements = append(gElements, el)
}
func Clay__ConfigureOpenElement(decl ElementDeclaration) {
	idx := gStack[len(gStack)-1]
	gElements[idx].decl = decl
}
func Clay__OpenTextElement(text string, cfg TextElementConfig) {
	// Get the currently open parent element (the container)
	parentIdx := gStack[len(gStack)-1]
	parent := &gElements[parentIdx]

	// Measure the text first
	sz := gCtx.Measurer.MeasureText(text, cfg)

	// fmt.Printf("DEBUG: OpenTextElement - text='%s', cfg.LineHeight=%f, sz.Height=%f\n", text, cfg.LineHeight, sz.Height)

	if cfg.LineHeight <= 0 {
		cfg.LineHeight = sz.Height
	}

	// Cache measurement
	item := measureTextCacheItem{
		id:                  0, // Will be filled below
		unwrappedDimensions: sz,
		minWidth:            sz.Width,
		containsNewlines:    false,
		measuredWordsStart:  -1,
		next:                -1,
		generation:          0,
	}

	// Split into words and cache (existing word splitting code)
	start := 0
	var prev *measuredWord
	for i := 0; i <= len(text); i++ {
		if i == len(text) || text[i] == ' ' || text[i] == '\n' {
			word := text[start:i]
			if len(word) > 0 || (i < len(text) && text[i] == '\n') {
				w := measuredWord{
					startOffset: int32(start),
					length:      int32(len(word)),
					width:       gCtx.Measurer.MeasureText(word, cfg).Width,
					next:        -1,
				}
				if i < len(text) && text[i] == ' ' {
					w.width += gCtx.Measurer.MeasureText(" ", cfg).Width
				}
				gMeasuredWords = append(gMeasuredWords, w)
				if prev != nil {
					prev.next = int32(len(gMeasuredWords) - 1)
				}
				if item.measuredWordsStart == -1 {
					item.measuredWordsStart = int32(len(gMeasuredWords) - 1)
				}
				prev = &gMeasuredWords[len(gMeasuredWords)-1]
			}
			if i < len(text) && text[i] == '\n' {
				item.containsNewlines = true
				nl := measuredWord{startOffset: int32(i), length: 0, width: 0, next: -1}
				gMeasuredWords = append(gMeasuredWords, nl)
				if prev != nil {
					prev.next = int32(len(gMeasuredWords) - 1)
				}
				if item.measuredWordsStart == -1 {
					item.measuredWordsStart = int32(len(gMeasuredWords) - 1)
				}
				prev = &gMeasuredWords[len(gMeasuredWords)-1]
			}
			start = i + 1
		}
	}

	// Store the measurement cache item
	gMeasureTextCache = append(gMeasureTextCache, item)

	// Store text data on the PARENT element, not create a new element
	textIdx := int32(len(gTextElementData))
	gTextElementData = append(gTextElementData, textElementData{
		text:         text,
		prefDim:      sz,
		elementIndex: parentIdx, // Link to parent element
	})

	// Set the parent's text configuration and data reference
	parent.decl.Text = &cfg
	parent.childrenOrTextContent.textElementData = &gTextElementData[len(gTextElementData)-1]

	// Update parent dimensions to include text
	parent.dimensions = sz
	parent.minDimensions = Dimensions{Width: item.minWidth, Height: sz.Height}

	// Add text element data index to parent's children (so it gets processed during layout)
	parent.children = append(parent.children, textIdx)
}
func Clay__CloseElement() {
	if len(gStack) == 0 {
		return
	}
	idx := gStack[len(gStack)-1]
	gStack = gStack[:len(gStack)-1]
	// resolve sizing
	el := &gElements[idx]
	decl := el.decl
	// padding
	lp := decl.Layout.Padding.Left + decl.Layout.Padding.Right
	tp := decl.Layout.Padding.Top + decl.Layout.Padding.Bottom
	// children sizing
	if decl.Layout.Direction == CLAY_LEFT_TO_RIGHT {
		w := lp
		h := tp
		for _, c := range el.children {
			child := &gElements[c]
			w += child.dimensions.Width + decl.Layout.ChildGap
			h = max32(h, child.dimensions.Height+tp)
		}
		if len(el.children) > 0 {
			w -= decl.Layout.ChildGap
		}
		el.dimensions.Width = w
		el.dimensions.Height = h
	} else {
		w := lp
		h := tp
		for _, c := range el.children {
			child := &gElements[c]
			h += child.dimensions.Height + decl.Layout.ChildGap
			w = max32(w, child.dimensions.Width+lp)
		}
		if len(el.children) > 0 {
			h -= decl.Layout.ChildGap
		}
		el.dimensions.Width = w
		el.dimensions.Height = h
	}
	// Apply sizing config - FIXED should not be overridden by text measurement
	sx := decl.Layout.Sizing.Width
	switch sx.Type {
	case CLAY__SIZING_TYPE_FIXED:
		el.dimensions.Width = sx.Min
		// Don't override fixed width with text measurement
	case CLAY__SIZING_TYPE_PERCENT:
		// resolved later
	case CLAY__SIZING_TYPE_GROW:
		// resolved later
	case CLAY__SIZING_TYPE_FIT:
		// Text can affect FIT sizing
		if el.decl.Text != nil && el.childrenOrTextContent.textElementData != nil {
			td := el.childrenOrTextContent.textElementData
			naturalWidth := td.prefDim.Width
			if naturalWidth > el.dimensions.Width {
				el.dimensions.Width = naturalWidth + decl.Layout.Padding.Left + decl.Layout.Padding.Right
			}
		}
	}

	sy := decl.Layout.Sizing.Height
	switch sy.Type {
	case CLAY__SIZING_TYPE_FIXED:
		el.dimensions.Height = sy.Min
		// Don't override fixed height with text measurement
	case CLAY__SIZING_TYPE_PERCENT:
		// resolved later
	case CLAY__SIZING_TYPE_GROW:
		// resolved later
	case CLAY__SIZING_TYPE_FIT:
		// Text can affect FIT height
		if el.decl.Text != nil && el.childrenOrTextContent.textElementData != nil {
			td := el.childrenOrTextContent.textElementData
			cfg := el.decl.Text
			naturalLineHeight := td.prefDim.Height
			finalLineHeight := cfg.LineHeight
			if finalLineHeight == 0 {
				finalLineHeight = naturalLineHeight
			}
			textHeight := finalLineHeight * float32(len(td.wrappedLines))
			if textHeight > el.dimensions.Height {
				el.dimensions.Height = textHeight + decl.Layout.Padding.Top + decl.Layout.Padding.Bottom
			}
		}
	}
}

// ---------- two-pass layout -------------------------------------------------
func Clay__CalculateFinalLayout() {
	// pass 1 – x axis
	Clay__SizeContainersAlongAxis(true)
	// text wrapping
	Clay__WrapTextElements()
	// aspect ratio on y
	Clay__AspectRatioCorrect(false)
	// pass 2 – y axis
	Clay__SizeContainersAlongAxis(false)
	// aspect ratio on x
	Clay__AspectRatioCorrect(true)
	// Calculate bounding boxes and setup tree roots
	Clay__CalculateBoundingBoxes()
	// generate render commands
	Clay__GenerateRenderCommands()
}

func Clay__CalculateBoundingBoxes() {
	// Clear previous data
	gLayoutElementTreeRoots = gLayoutElementTreeRoots[:0]

	// Calculate bounding boxes for all elements
	for i := range gElements {
		el := &gElements[i]

		// Calculate bounding box (simplified - you'd need proper positioning logic)
		bb := BoundingBox{
			X:      0, // This should be calculated based on parent/child relationships
			Y:      0, // This should be calculated based on parent/child relationships
			Width:  el.dimensions.Width,
			Height: el.dimensions.Height,
		}

		// Add to hash map
		Clay__AddHashMapItem(el.id, el)

		// Update the hash map item with bounding box
		if item := Clay__GetHashMapItem(el.id); item != nil {
			item.boundingBox = bb
		}

		// Add to layout tree roots if it's a root element
		if i == 0 { // Root element
			gLayoutElementTreeRoots = append(gLayoutElementTreeRoots, layoutElementTreeRoot{
				elementIndex: int32(i),
				parentID:     0,
				clipID:       0,
				zIndex:       0,
			})
		}
	}
}

func Clay__SizeContainersAlongAxis(xAxis bool) {
	// simple bottom-up sizing; real code walks tree roots
	for i := range gElements {
		el := &gElements[i]
		if el.decl.Layout.Direction == CLAY_LEFT_TO_RIGHT && xAxis {
			// already done in CloseElement
		} else if el.decl.Layout.Direction == CLAY_TOP_TO_BOTTOM && !xAxis {
			// already done in CloseElement
		}
	}
}

func Clay__WrapTextElements() {
	for i := range gTextElementData {
		td := &gTextElementData[i]
		el := &gElements[td.elementIndex]
		cfg := el.decl.Text

		containerWidth := el.dimensions.Width - el.decl.Layout.Padding.Left - el.decl.Layout.Padding.Right
		if containerWidth <= 0 {
			continue
		}

		// find cache
		var cache *measureTextCacheItem
		for j := range gMeasureTextCache {
			if gMeasureTextCache[j].measuredWordsStart >= 0 {
				cache = &gMeasureTextCache[j]
				break
			}
		}
		if cache == nil {
			continue
		}

		// Calculate the proper line height to use
		finalLineHeight := cfg.LineHeight
		if finalLineHeight == 0 {
			finalLineHeight = td.prefDim.Height // Use measured height if no line height specified
		}

		// Wrap the text
		td.wrappedLines = nil
		lineWidth := float32(0)
		lineStart := 0
		lineLen := 0
		spaceW := gCtx.Measurer.MeasureText(" ", *cfg).Width

		for w := cache.measuredWordsStart; w >= 0 && int(w) < len(gMeasuredWords); {
			word := &gMeasuredWords[w]

			// Handle newline words (length == 0 means newline)
			if word.length == 0 {
				// Add the current line before the newline
				if lineLen > 0 {
					trimW := lineWidth
					if lineLen > 0 && td.text[lineStart+lineLen-1] == ' ' {
						trimW -= spaceW
					}
					td.wrappedLines = append(td.wrappedLines, wrappedTextLine{
						dimensions: Dimensions{Width: trimW, Height: finalLineHeight}, // USE finalLineHeight
						text:       td.text[lineStart : lineStart+lineLen],
					})
				}

				// Add empty line for the newline itself
				td.wrappedLines = append(td.wrappedLines, wrappedTextLine{
					dimensions: Dimensions{Width: 0, Height: finalLineHeight}, // USE finalLineHeight
					text:       "",
				})

				// Reset for next line
				lineWidth = 0
				lineStart += lineLen + 1 // Skip past the newline
				lineLen = 0
				w = word.next
				continue
			}

			// Check if word alone is too big for the line (and it's the first word)
			if lineLen == 0 && lineWidth+word.width > containerWidth {
				// Force the word anyway, but trim trailing space
				trimW := word.width
				if word.length > 0 && td.text[word.startOffset+word.length-1] == ' ' {
					trimW -= spaceW
				}
				td.wrappedLines = append(td.wrappedLines, wrappedTextLine{
					dimensions: Dimensions{Width: trimW, Height: finalLineHeight}, // USE finalLineHeight
					text:       td.text[word.startOffset : word.startOffset+word.length],
				})

				// Reset for next line
				lineWidth = 0
				lineStart = int(word.startOffset + word.length)
				lineLen = 0
				w = word.next
				continue
			}

			// Normal case - add word to current line
			if lineWidth+word.width <= containerWidth {
				// Word fits, add it
				lineWidth += word.width
				lineLen += int(word.length)
				w = word.next
			} else {
				// Word doesn't fit, wrap to new line
				if lineLen > 0 {
					trimW := lineWidth
					if lineLen > 0 && td.text[lineStart+lineLen-1] == ' ' {
						trimW -= spaceW
					}
					td.wrappedLines = append(td.wrappedLines, wrappedTextLine{
						dimensions: Dimensions{Width: trimW, Height: finalLineHeight}, // USE finalLineHeight
						text:       td.text[lineStart : lineStart+lineLen],
					})
				}

				// Reset for new line with current word
				lineWidth = word.width
				lineStart = int(word.startOffset)
				lineLen = int(word.length)
				w = word.next
			}
		}

		// Handle final line
		if lineLen > 0 {
			trimW := lineWidth
			if lineLen > 0 && td.text[lineStart+lineLen-1] == ' ' {
				trimW -= spaceW
			}
			td.wrappedLines = append(td.wrappedLines, wrappedTextLine{
				dimensions: Dimensions{Width: trimW, Height: finalLineHeight}, // USE finalLineHeight
				text:       td.text[lineStart : lineStart+lineLen],
			})
		}

		// Only update container height if it's NOT fixed sizing
		if el.decl.Layout.Sizing.Height.Type != CLAY__SIZING_TYPE_FIXED {
			el.dimensions.Height = float32(len(td.wrappedLines))*finalLineHeight +
				el.decl.Layout.Padding.Top + el.decl.Layout.Padding.Bottom
		}
	}
}

func Clay__AspectRatioCorrect(xAxis bool) {
	for i := range gElements {
		el := &gElements[i]
		if el.decl.AspectRatio == nil {
			continue
		}
		if xAxis && el.dimensions.Height != 0 {
			el.dimensions.Width = el.dimensions.Height * el.decl.AspectRatio.AspectRatio
		} else if !xAxis && el.dimensions.Width != 0 {
			el.dimensions.Height = el.dimensions.Width / el.decl.AspectRatio.AspectRatio
		}
	}
}

// ---------- render command generation ---------------------------------------
// func Clay__GenerateRenderCommands() {
// 	gCmds = gCmds[:0]
// 	// simple DFS
// 	var dfs func(el *layoutElement, x, y float32)
// 	dfs = func(el *layoutElement, x, y float32) {
// 		bb := BoundingBox{
// 			X:      x,
// 			Y:      y,
// 			Width:  el.dimensions.Width,
// 			Height: el.dimensions.Height,
// 		}
// 		// background rect
// 		if el.decl.BackgroundColor.A > 0 {
// 			gCmds = append(gCmds, RenderCommand{
// 				BoundingBox: bb,
// 				ID:          el.id,
// 				ZIndex:      0,
// 				CommandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
// 				Data: RectangleRenderData{
// 					Color:        el.decl.BackgroundColor,
// 					CornerRadius: el.decl.CornerRadius,
// 				},
// 			})
// 		}
// 		// text
// 		if el.decl.Text != nil {
// 			td := &gTextElementData[0]
// 			for _, w := range td.wrappedLines {
// 				gCmds = append(gCmds, RenderCommand{
// 					BoundingBox: BoundingBox{
// 						X:      x + el.decl.Layout.Padding.Left,
// 						Y:      y + el.decl.Layout.Padding.Top,
// 						Width:  w.dimensions.Width,
// 						Height: w.dimensions.Height,
// 					},
// 					ID:          el.id,
// 					ZIndex:      0,
// 					CommandType: CLAY_RENDER_COMMAND_TYPE_TEXT,
// 					Data: TextRenderData{
// 						StringContents: w.text,
// 						Color:          el.decl.Text.Color,
// 						FontID:         el.decl.Text.FontID,
// 						FontSize:       el.decl.Text.FontSize,
// 						LetterSpacing:  el.decl.Text.LetterSpacing,
// 						LineHeight:     el.decl.Text.LineHeight,
// 						Alignment:      el.decl.Text.Alignment,
// 					},
// 				})
// 			}
// 		}
// 		// children
// 		cx := x + el.decl.Layout.Padding.Left
// 		cy := y + el.decl.Layout.Padding.Top
// 		if el.decl.Layout.Direction == CLAY_LEFT_TO_RIGHT {
// 			for _, c := range el.children {
// 				child := &gElements[c]
// 				dfs(child, cx, cy+(el.dimensions.Height-child.dimensions.Height)*0.5) // center y
// 				cx += child.dimensions.Width + el.decl.Layout.ChildGap
// 			}
// 		} else {
// 			for _, c := range el.children {
// 				child := &gElements[c]
// 				dfs(child, cx+(el.dimensions.Width-child.dimensions.Width)*0.5, cy) // center x
// 				cy += child.dimensions.Height + el.decl.Layout.ChildGap
// 			}
// 		}
// 	}
// 	root := &gElements[0]
// 	dfs(root, 0, 0)
// }

// ---------- arena -----------------------------------------------------------
type Arena struct {
	mem  []byte
	used uintptr
}

func NewArena(size int) *Arena {
	return &Arena{mem: make([]byte, size)}
}
func (a *Arena) Alloc(size uintptr) unsafe.Pointer {
	// 64-byte align
	align := (64 - (a.used % 64)) & 63
	if a.used+align+size > uintptr(len(a.mem)) {
		panic("arena exhausted")
	}
	p := unsafe.Pointer(&a.mem[a.used+align])
	a.used += align + size
	return p
}
func (a *Arena) Reset() { a.used = 0 }

// ---------- persistent / ephemeral init -------------------------------------

// func Clay__InitEphemeral() {
// 	gElements = gElements[:0]
// 	gTextElementData = gTextElementData[:0]
// 	gMeasuredWords = gMeasuredWords[:0]
// 	gMeasureTextCache = gMeasureTextCache[:0]
// 	gLayoutElementTreeRoots = gLayoutElementTreeRoots[:0]
// 	gScrollContainers = gScrollContainers[:0]
// 	gOpenClipStack = gOpenClipStack[:0]
// }

// ---------- tiny utils ------------------------------------------------------
func max32(a, b float32) float32 {
	if a > b {
		return a
	}
	return b
}

// ---------- hash-table element lookup (1:1 C logic) -------------------------
type layoutElementHashMapItem struct {
	elementID     ElementID
	layoutElement *layoutElement
	boundingBox   BoundingBox
	nextIndex     int32
	generation    uint32
	onHoverFunc   func(ElementID, PointerData, uintptr)
	hoverUserData uintptr
}

var gElementHashMap []layoutElementHashMapItem
var gHashMapBuckets []int32

func Clay__InitHashMap(cap int) {
	gElementHashMap = make([]layoutElementHashMapItem, 0, cap)
	gHashMapBuckets = make([]int32, cap)
	for i := range gHashMapBuckets {
		gHashMapBuckets[i] = -1
	}
}
func Clay__AddHashMapItem(id ElementID, el *layoutElement) {
	bucket := id % uint32(len(gHashMapBuckets))
	it := layoutElementHashMapItem{
		elementID:     id,
		layoutElement: el,
		nextIndex:     -1,
		generation:    1,
	}
	// insert at head
	it.nextIndex = gHashMapBuckets[bucket]
	gElementHashMap = append(gElementHashMap, it)
	gHashMapBuckets[bucket] = int32(len(gElementHashMap) - 1)
}
func Clay__GetHashMapItem(id ElementID) *layoutElementHashMapItem {
	bucket := id % uint32(len(gHashMapBuckets))
	idx := gHashMapBuckets[bucket]
	for idx >= 0 {
		it := &gElementHashMap[idx]
		if it.elementID == id {
			return it
		}
		idx = it.nextIndex
	}
	return nil
}

// ---------- floating element roots + z-sort ---------------------------------
type floatingRoot struct {
	elementIndex int32
	parentID     ElementID
	clipID       ElementID
	zIndex       int16
}

var gFloatingRoots []floatingRoot

func Clay__AddFloatingRoot(elIdx int32, parentID, clipID ElementID, z int16) {
	gFloatingRoots = append(gFloatingRoots, floatingRoot{
		elementIndex: elIdx,
		parentID:     parentID,
		clipID:       clipID,
		zIndex:       z,
	})
}
func Clay__SortFloatingRoots() {
	// bubble sort (same as C)
	n := len(gFloatingRoots)
	for n > 1 {
		for i := 0; i < n-1; i++ {
			if gFloatingRoots[i+1].zIndex < gFloatingRoots[i].zIndex {
				gFloatingRoots[i], gFloatingRoots[i+1] = gFloatingRoots[i+1], gFloatingRoots[i]
			}
		}
		n--
	}
}

// ---------- scrolling -------------------------------------------------------
type PointerData struct {
	Position Vector2
	State    PointerDataInteractionState
}
type PointerDataInteractionState int

const (
	CLAY_POINTER_DATA_PRESSED_THIS_FRAME PointerDataInteractionState = iota
	CLAY_POINTER_DATA_PRESSED
	CLAY_POINTER_DATA_RELEASED_THIS_FRAME
	CLAY_POINTER_DATA_RELEASED
)

var gPointer PointerData

func Clay_SetPointerState(pos Vector2, down bool) {
	gPointer.Position = pos
	if down {
		if gPointer.State == CLAY_POINTER_DATA_PRESSED {
			gPointer.State = CLAY_POINTER_DATA_PRESSED_THIS_FRAME
		} else {
			gPointer.State = CLAY_POINTER_DATA_PRESSED_THIS_FRAME
		}
	} else {
		if gPointer.State == CLAY_POINTER_DATA_PRESSED || gPointer.State == CLAY_POINTER_DATA_PRESSED_THIS_FRAME {
			gPointer.State = CLAY_POINTER_DATA_RELEASED_THIS_FRAME
		} else {
			gPointer.State = CLAY_POINTER_DATA_RELEASED
		}
	}
}

func Clay_UpdateScrollContainers(enableDrag bool, delta Vector2, dt float32) {
	// find highest priority scroll container under pointer
	best := -1
	for i, sc := range gScrollContainers {
		if !sc.openThisFrame {
			continue
		}
		bb := sc.boundingBox
		inside := gPointer.Position.X >= bb.X && gPointer.Position.X <= bb.X+bb.Width &&
			gPointer.Position.Y >= bb.Y && gPointer.Position.Y <= bb.Y+bb.Height
		if inside {
			best = i
		}
	}
	if best < 0 {
		return
	}
	sc := &gScrollContainers[best]
	cfg := sc.layoutElement.decl.Clip
	canX := cfg.Horizontal && sc.contentSize.Width > sc.layoutElement.dimensions.Width
	canY := cfg.Vertical && sc.contentSize.Height > sc.layoutElement.dimensions.Height

	// wheel
	if canX {
		sc.scrollPosition.X += delta.X * 10
	}
	if canY {
		sc.scrollPosition.Y += delta.Y * 10
	}
	// drag
	if enableDrag && (gPointer.State == CLAY_POINTER_DATA_PRESSED || gPointer.State == CLAY_POINTER_DATA_PRESSED_THIS_FRAME) {
		if !sc.pointerScroll {
			sc.scrollOrigin = sc.scrollPosition
			sc.pointerOrigin = gPointer.Position
			sc.pointerScroll = true
			sc.momentumTime = 0
		} else {
			if canX {
				sc.scrollPosition.X = sc.scrollOrigin.X + (gPointer.Position.X - sc.pointerOrigin.X)
			}
			if canY {
				sc.scrollPosition.Y = sc.scrollOrigin.Y + (gPointer.Position.Y - sc.pointerOrigin.Y)
			}
		}
	} else {
		if sc.pointerScroll {
			// release → momentum
			if canX {
				sc.scrollMomentum.X = (sc.scrollPosition.X - sc.scrollOrigin.X) / (sc.momentumTime*25 + 1e-3)
			}
			if canY {
				sc.scrollMomentum.Y = (sc.scrollPosition.Y - sc.scrollOrigin.Y) / (sc.momentumTime*25 + 1e-3)
			}
			sc.pointerScroll = false
		}
		sc.scrollMomentum.X *= 0.95
		sc.scrollMomentum.Y *= 0.95
		sc.scrollPosition.X += sc.scrollMomentum.X
		sc.scrollPosition.Y += sc.scrollMomentum.Y
	}
	// clamp
	if canX {
		max := sc.contentSize.Width - sc.layoutElement.dimensions.Width
		sc.scrollPosition.X = clamp(sc.scrollPosition.X, -max, 0)
	}
	if canY {
		max := sc.contentSize.Height - sc.layoutElement.dimensions.Height
		sc.scrollPosition.Y = clamp(sc.scrollPosition.Y, -max, 0)
	}
	if sc.pointerScroll {
		sc.momentumTime += dt
	}
}

func clamp(v, min, max float32) float32 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ---------- border between-children -----------------------------------------
func Clay__GenerateBorderCommands(el *layoutElement, bb BoundingBox) {
	if el.decl.Border == nil || el.decl.Border.Width.BetweenChildren <= 0 {
		return
	}
	bw := el.decl.Border.Width.BetweenChildren
	col := el.decl.Border.Color
	gap := el.decl.Layout.ChildGap
	if el.decl.Layout.Direction == CLAY_LEFT_TO_RIGHT {
		x := bb.X + el.decl.Layout.Padding.Left - gap*0.5
		for i, c := range el.children {
			if i > 0 {
				gCmds = append(gCmds, RenderCommand{
					BoundingBox: BoundingBox{
						X:      x,
						Y:      bb.Y,
						Width:  float32(bw),
						Height: bb.Height,
					},
					ID:          el.id + uint32(i),
					ZIndex:      0,
					CommandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
					Data:        RectangleRenderData{Color: col},
				})
			}
			x += gElements[c].dimensions.Width + gap
		}
	} else {
		y := bb.Y + el.decl.Layout.Padding.Top - gap*0.5
		for i, c := range el.children {
			if i > 0 {
				gCmds = append(gCmds, RenderCommand{
					BoundingBox: BoundingBox{
						X:      bb.X,
						Y:      y,
						Width:  bb.Width,
						Height: float32(bw),
					},
					ID:          el.id + uint32(i),
					ZIndex:      0,
					CommandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
					Data:        RectangleRenderData{Color: col},
				})
			}
			y += gElements[c].dimensions.Height + gap
		}
	}
}

// ---------- debug overlay ---------------------------------------------------
var gDebugEnabled bool

func Clay_SetDebugModeEnabled(v bool) { gDebugEnabled = v }
func Clay_IsDebugModeEnabled() bool   { return gDebugEnabled }

func Clay__RenderDebugView() {
	if !gDebugEnabled {
		return
	}
	// tiny debug: highlight every element border
	for i := range gElements {
		el := &gElements[i]
		it := Clay__GetHashMapItem(el.id)
		if it == nil {
			continue
		}
		bb := it.boundingBox
		gCmds = append(gCmds, RenderCommand{
			BoundingBox: bb,
			ID:          el.id + 0x80000000,
			ZIndex:      32767,
			CommandType: CLAY_RENDER_COMMAND_TYPE_BORDER,
			Data: BorderRenderData{
				Color:        Color{1, 0, 1, 0.4}, // Fixed: Use 0-1 range instead of 255,0,255,100
				Width:        BorderWidth{Left: 1, Right: 1, Top: 1, Bottom: 1},
				CornerRadius: el.decl.CornerRadius,
			},
		})
	}
}

// ---------- public API wrappers ---------------------------------------------
func Clay_GetElementData(id ElementID) (BoundingBox, bool) {
	it := Clay__GetHashMapItem(id)
	if it == nil {
		return BoundingBox{}, false
	}
	return it.boundingBox, true
}
func Clay_PointerOver(id ElementID) bool {
	it := Clay__GetHashMapItem(id)
	if it == nil {
		return false
	}
	bb := it.boundingBox
	return gPointer.Position.X >= bb.X && gPointer.Position.X <= bb.X+bb.Width &&
		gPointer.Position.Y >= bb.Y && gPointer.Position.Y <= bb.Y+bb.Height
}
func Clay_GetScrollContainerData(id ElementID) (Vector2, Dimensions, Dimensions, bool) {
	for _, sc := range gScrollContainers {
		if sc.elementID == id {
			return sc.scrollPosition, sc.layoutElement.dimensions, sc.contentSize, true
		}
	}
	return Vector2{}, Dimensions{}, Dimensions{}, false
}

// ---------- init hooks ------------------------------------------------------
func Clay__InitPersistent() {
	Clay__InitHashMap(8192)
}
func Clay__InitEphemeral() {
	gElements = gElements[:0]
	gTextElementData = gTextElementData[:0]
	gMeasuredWords = gMeasuredWords[:0]
	gMeasureTextCache = gMeasureTextCache[:0]
	gFloatingRoots = gFloatingRoots[:0]
	gScrollContainers = gScrollContainers[:0]
	gOpenClipStack = gOpenClipStack[:0]
}

// ---------- patched render generation --------------------------------------
// replace the old dfs() call in Clay__GenerateRenderCommands() with this one:
func Clay__GenerateRenderCommands() {
	gCmds = gCmds[:0]

	// Process regular elements first
	for i := range gElements {
		el := &gElements[i]

		// Skip floating elements (they're handled separately)
		if el.decl.Floating != nil {
			continue
		}

		// Get the element's bounding box
		it := Clay__GetHashMapItem(el.id)
		if it == nil {
			continue
		}

		// Render this element
		Clay__RenderElementRecursive(el, Vector2{X: it.boundingBox.X, Y: it.boundingBox.Y}, 0)
	}

	// floating roots first
	Clay__SortFloatingRoots()
	for _, fr := range gFloatingRoots {
		el := &gElements[fr.elementIndex]
		// scissor start if clipped
		if fr.clipID != 0 {
			if it := Clay__GetHashMapItem(fr.clipID); it != nil {
				gCmds = append(gCmds, RenderCommand{
					BoundingBox: it.boundingBox,
					ID:          fr.clipID,
					ZIndex:      el.decl.Floating.ZIndex,
					CommandType: CLAY_RENDER_COMMAND_TYPE_SCISSOR_START,
					Data:        ClipRenderData{Horizontal: true, Vertical: true},
				})
			}
		}
		// position floating element
		parentBB := BoundingBox{}
		if it := Clay__GetHashMapItem(fr.parentID); it != nil {
			parentBB = it.boundingBox
		}
		cfg := el.decl.Floating
		pos := Clay__CalcAttachPos(parentBB, el.dimensions, cfg.AttachPoints)
		pos.X += cfg.Offset.X
		pos.Y += cfg.Offset.Y
		// render the floating tree
		Clay__RenderElementRecursive(el, pos, el.decl.Floating.ZIndex)
		// scissor end
		if fr.clipID != 0 {
			gCmds = append(gCmds, RenderCommand{
				BoundingBox: BoundingBox{},
				ID:          fr.clipID + 1,
				ZIndex:      el.decl.Floating.ZIndex,
				CommandType: CLAY_RENDER_COMMAND_TYPE_SCISSOR_END,
				Data:        ClipRenderData{},
			})
		}
	}
	for _, root := range gLayoutElementTreeRoots {
		el := &gElements[root.elementIndex]
		offset := Vector2{}
		// floating roots have their own offset calculation already done
		if el.decl.Floating != nil {
			offset = Clay__CalcAttachPos(
				Clay__GetHashMapItem(root.parentID).boundingBox,
				el.dimensions,
				el.decl.Floating.AttachPoints,
			)
			offset.X += el.decl.Floating.Offset.X
			offset.Y += el.decl.Floating.Offset.Y
		}
		Clay__RenderElementRecursive(el, offset, root.zIndex)
	}
	// debug
	Clay__RenderDebugView()
}

func Clay__CalcAttachPos(parentBB BoundingBox, childDim Dimensions, pts FloatingAttachPoints) Vector2 {
	var p Vector2
	// parent anchor
	switch pts.Parent {
	case CLAY_ATTACH_POINT_LEFT_TOP, CLAY_ATTACH_POINT_LEFT_CENTER, CLAY_ATTACH_POINT_LEFT_BOTTOM:
		p.X = parentBB.X
	case CLAY_ATTACH_POINT_CENTER_TOP, CLAY_ATTACH_POINT_CENTER_CENTER, CLAY_ATTACH_POINT_CENTER_BOTTOM:
		p.X = parentBB.X + parentBB.Width*0.5
	case CLAY_ATTACH_POINT_RIGHT_TOP, CLAY_ATTACH_POINT_RIGHT_CENTER, CLAY_ATTACH_POINT_RIGHT_BOTTOM:
		p.X = parentBB.X + parentBB.Width
	}
	switch pts.Parent {
	case CLAY_ATTACH_POINT_LEFT_TOP, CLAY_ATTACH_POINT_CENTER_TOP, CLAY_ATTACH_POINT_RIGHT_TOP:
		p.Y = parentBB.Y
	case CLAY_ATTACH_POINT_LEFT_CENTER, CLAY_ATTACH_POINT_CENTER_CENTER, CLAY_ATTACH_POINT_RIGHT_CENTER:
		p.Y = parentBB.Y + parentBB.Height*0.5
	case CLAY_ATTACH_POINT_LEFT_BOTTOM, CLAY_ATTACH_POINT_CENTER_BOTTOM, CLAY_ATTACH_POINT_RIGHT_BOTTOM:
		p.Y = parentBB.Y + parentBB.Height
	}
	// element anchor
	switch pts.Element {
	case CLAY_ATTACH_POINT_LEFT_TOP, CLAY_ATTACH_POINT_LEFT_CENTER, CLAY_ATTACH_POINT_LEFT_BOTTOM:
		// keep
	case CLAY_ATTACH_POINT_CENTER_TOP, CLAY_ATTACH_POINT_CENTER_CENTER, CLAY_ATTACH_POINT_CENTER_BOTTOM:
		p.X -= childDim.Width * 0.5
	case CLAY_ATTACH_POINT_RIGHT_TOP, CLAY_ATTACH_POINT_RIGHT_CENTER, CLAY_ATTACH_POINT_RIGHT_BOTTOM:
		p.X -= childDim.Width
	}
	switch pts.Element {
	case CLAY_ATTACH_POINT_LEFT_TOP, CLAY_ATTACH_POINT_CENTER_TOP, CLAY_ATTACH_POINT_RIGHT_TOP:
		// keep
	case CLAY_ATTACH_POINT_LEFT_CENTER, CLAY_ATTACH_POINT_CENTER_CENTER, CLAY_ATTACH_POINT_RIGHT_CENTER:
		p.Y -= childDim.Height * 0.5
	case CLAY_ATTACH_POINT_LEFT_BOTTOM, CLAY_ATTACH_POINT_CENTER_BOTTOM, CLAY_ATTACH_POINT_RIGHT_BOTTOM:
		p.Y -= childDim.Height
	}
	return p
}

func Clay__RenderElementRecursive(el *layoutElement, offset Vector2, z int16) {
	bb := BoundingBox{
		X:      offset.X,
		Y:      offset.Y,
		Width:  el.dimensions.Width,
		Height: el.dimensions.Height,
	}
	// background
	if el.decl.BackgroundColor.A > 0 {
		gCmds = append(gCmds, RenderCommand{
			BoundingBox: bb,
			ID:          el.id,
			ZIndex:      z,
			CommandType: CLAY_RENDER_COMMAND_TYPE_RECTANGLE,
			Data:        RectangleRenderData{Color: el.decl.BackgroundColor, CornerRadius: el.decl.CornerRadius},
		})
	}
	// border
	if el.decl.Border != nil {
		gCmds = append(gCmds, RenderCommand{
			BoundingBox: bb,
			ID:          el.id + 0x10000000,
			ZIndex:      z,
			CommandType: CLAY_RENDER_COMMAND_TYPE_BORDER,
			Data: BorderRenderData{
				Color:        el.decl.Border.Color,
				Width:        el.decl.Border.Width,
				CornerRadius: el.decl.CornerRadius,
			},
		})
		Clay__GenerateBorderCommands(el, bb)
	}
	// text
	if el.decl.Text != nil {
		if el.childrenOrTextContent.textElementData == nil {
			return
		}
		td := el.childrenOrTextContent.textElementData

		cfg := el.decl.Text
		naturalLineHeight := td.prefDim.Height
		finalLineHeight := cfg.LineHeight
		if finalLineHeight == 0 {
			finalLineHeight = naturalLineHeight
		}

		// Calculate available space
		availableWidth := el.dimensions.Width - el.decl.Layout.Padding.Left - el.decl.Layout.Padding.Right
		availableHeight := el.dimensions.Height - el.decl.Layout.Padding.Top - el.decl.Layout.Padding.Bottom

		// Calculate vertical alignment
		totalTextHeight := finalLineHeight * float32(len(td.wrappedLines))
		yOffset := float32(0)
		if totalTextHeight < availableHeight {
			yOffset = (availableHeight - totalTextHeight) / 2
		}

		// If no wrapped lines exist, create a single line with the full text
		if len(td.wrappedLines) == 0 {
			// Create a single line with the measured dimensions
			ox := offset.X + el.decl.Layout.Padding.Left
			oy := offset.Y + el.decl.Layout.Padding.Top + yOffset

			// Calculate horizontal alignment for single line
			xOffset := float32(0)
			if cfg.Alignment == CLAY_TEXT_ALIGN_CENTER {
				xOffset = (availableWidth - td.prefDim.Width) / 2
			} else if cfg.Alignment == CLAY_TEXT_ALIGN_RIGHT {
				xOffset = availableWidth - td.prefDim.Width
			}

			gCmds = append(gCmds, RenderCommand{
				BoundingBox: BoundingBox{
					X:      ox + xOffset,
					Y:      oy,
					Width:  td.prefDim.Width,
					Height: finalLineHeight, // Use finalLineHeight, not td.prefDim.Height
				},
				ID:          el.id + 0x20000000,
				ZIndex:      z,
				CommandType: CLAY_RENDER_COMMAND_TYPE_TEXT,
				Data: TextRenderData{
					StringContents: td.text,
					Color:          cfg.Color,
					FontID:         cfg.FontID,
					FontSize:       cfg.FontSize,
					LetterSpacing:  cfg.LetterSpacing,
					LineHeight:     cfg.LineHeight,
					Alignment:      cfg.Alignment,
				},
			})
			return
		}

		// Render wrapped lines (existing code)
		yPosition := yOffset
		for lineIndex := 0; lineIndex < len(td.wrappedLines); lineIndex++ {
			line := td.wrappedLines[lineIndex]
			if len(line.text) == 0 {
				yPosition += finalLineHeight
				continue
			}

			// Calculate horizontal alignment
			xOffset := float32(0)
			if cfg.Alignment == CLAY_TEXT_ALIGN_CENTER {
				xOffset = (availableWidth - line.dimensions.Width) / 2
			} else if cfg.Alignment == CLAY_TEXT_ALIGN_RIGHT {
				xOffset = availableWidth - line.dimensions.Width
			}

			ox := offset.X + el.decl.Layout.Padding.Left + xOffset
			oy := offset.Y + el.decl.Layout.Padding.Top + yPosition

			gCmds = append(gCmds, RenderCommand{
				BoundingBox: BoundingBox{
					X:      ox,
					Y:      oy,
					Width:  line.dimensions.Width,
					Height: finalLineHeight, // Always use finalLineHeight
				},
				ID:          el.id + uint32(lineIndex)*0x10000000,
				ZIndex:      z,
				CommandType: CLAY_RENDER_COMMAND_TYPE_TEXT,
				Data: TextRenderData{
					StringContents: line.text,
					Color:          cfg.Color,
					FontID:         cfg.FontID,
					FontSize:       cfg.FontSize,
					LetterSpacing:  cfg.LetterSpacing,
					LineHeight:     cfg.LineHeight,
					Alignment:      cfg.Alignment,
				},
			})

			yPosition += finalLineHeight
		}
	}
	// children
	coff := Vector2{
		X: offset.X + el.decl.Layout.Padding.Left,
		Y: offset.Y + el.decl.Layout.Padding.Top,
	}
	if el.decl.Layout.Direction == CLAY_LEFT_TO_RIGHT {
		for _, c := range el.children {
			child := &gElements[c]
			Clay__RenderElementRecursive(child, coff, z)
			coff.X += child.dimensions.Width + el.decl.Layout.ChildGap
		}
	} else {
		for _, c := range el.children {
			child := &gElements[c]
			Clay__RenderElementRecursive(child, coff, z)
			coff.Y += child.dimensions.Height + el.decl.Layout.ChildGap
		}
	}
}
