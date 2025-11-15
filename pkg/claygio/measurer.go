package claygio

import (
	"image"

	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/zodimo/clay-go/clay"
	"golang.org/x/image/math/fixed"
)

type TextMeasurer interface {
	MeasureText(text clay.Clay_StringSlice, config *clay.Clay_TextElementConfig, userData interface{}) clay.Clay_Dimensions
}

type measurer struct {
	options *MeasurerOptions
}

type MeasurerOptions struct {
	FontManager *FontManager
}

type MeasurerOption func(*MeasurerOptions)

func MeasurerWithFontManager(fontManager *FontManager) MeasurerOption {
	return func(o *MeasurerOptions) {
		o.FontManager = fontManager
	}
}

// NewMeasurer initializes and returns a clay.TextMeasurer using Gioui's engine.
func NewMeasurer(opts ...MeasurerOption) TextMeasurer {
	options := &MeasurerOptions{
		FontManager: NewFontManager(),
	}
	for _, opt := range opts {
		opt(options)
	}

	return &measurer{
		options: options,
	}
}

func (m *measurer) MeasureText(textToMeasureSlice clay.Clay_StringSlice, cfg *clay.Clay_TextElementConfig, userData interface{}) clay.Clay_Dimensions {
	textToMeasure := textToMeasureSlice.String()

	gtx, ok := userData.(layout.Context)
	if !ok {
		panic("userData is not a layout.Context")
	}

	// Get font using the same method as the renderer (FontManager)
	fontManager := m.options.FontManager
	fontObj := fontManager.GetFont(cfg.FontId)

	textSize := fixed.I(gtx.Sp(unit.Sp(cfg.FontSize)))
	lineHeight := fixed.I(gtx.Sp(unit.Sp(cfg.LineHeight)))

	// Map clay wrap mode to Gioui wrap policy
	var wrapPolicy text.WrapPolicy
	switch cfg.WrapMode {
	case clay.CLAY_TEXT_WRAP_NONE:
		wrapPolicy = text.WrapWords // For measurement, we still need some wrapping, but use large MaxWidth
	case clay.CLAY_TEXT_WRAP_NEWLINES:
		wrapPolicy = text.WrapWords
	case clay.CLAY_TEXT_WRAP_WORDS:
		wrapPolicy = text.WrapGraphemes
	default:
		wrapPolicy = text.WrapGraphemes
	}

	// Map clay config to Gioui's layout parameters
	// For measurement, we want to measure the full text without truncation
	// so the layout system can allocate enough space. The renderer uses MaxLines: 1
	// but if we allocate enough space, it won't truncate.
	params := text.Parameters{
		Font:       fontObj,
		Alignment:  text.Start,
		LineHeight: lineHeight,
		MaxLines:   0, // 0 means unlimited - measure full text
		// Set a very large MaxWidth to ensure the text is measured as a single line
		MaxWidth:   1000000,
		PxPerEm:    textSize,
		Locale:     system.Locale{},
		WrapPolicy: wrapPolicy,
	}

	// Perform the text layout - same as Label.LayoutDetailed
	fontManager.GetShaper().LayoutString(params, textToMeasure)

	// Calculate bounds the same way as textIterator.processGlyph in Label
	var bounds image.Rectangle
	var first bool = true

	// Iterate through all glyphs until NextGlyph returns false (matching Label's iteration)
	for g, ok := fontManager.GetShaper().NextGlyph(); ok; g, ok = fontManager.GetShaper().NextGlyph() {
		// Calculate logical bounds for this glyph - same as Label's processGlyph
		logicalBounds := image.Rectangle{
			Min: image.Pt(g.X.Floor(), int(g.Y)-g.Ascent.Ceil()),
			Max: image.Pt((g.X + g.Advance).Ceil(), int(g.Y)+g.Descent.Ceil()),
		}

		if first {
			first = false
			bounds = logicalBounds
		} else {
			// Accumulate bounds by taking min/max - same as Label
			bounds.Min.X = min(bounds.Min.X, logicalBounds.Min.X)
			bounds.Min.Y = min(bounds.Min.Y, logicalBounds.Min.Y)
			bounds.Max.X = max(bounds.Max.X, logicalBounds.Max.X)
			bounds.Max.Y = max(bounds.Max.Y, logicalBounds.Max.Y)
		}
	}

	// Get the size from bounds - same as Label returns it.bounds.Size()
	size := bounds.Size()

	// fmt.Printf("bounds: %+v, size: %+v\n", bounds, size)

	// Convert to float32 dimensions
	return clay.Clay_Dimensions{
		Width:  float32(size.X),
		Height: float32(size.Y),
	}
}
