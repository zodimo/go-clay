package claygio

import (
	"gioui.org/font/gofont"
	"gioui.org/io/system"
	"gioui.org/text"
	"github.com/zodimo/clay-go/clay"
	"golang.org/x/image/math/fixed"
)

type TextMeasurer interface {
	MeasureText(text clay.Clay_StringSlice, config *clay.Clay_TextElementConfig, userData interface{}) clay.Clay_Dimensions
}

type measurer struct {
	shaper  *text.Shaper
	options *MeasurerOptions
}

type MeasurerOptions struct {
	Collection      []text.FontFace
	LineHeightScale float32
}

type MeasurerOption func(*MeasurerOptions)

func MeasurerWithCollection(collection []text.FontFace) MeasurerOption {
	return func(o *MeasurerOptions) {
		o.Collection = collection
	}
}
func MeasurerWithLineHeightScale(lineHeightScale float32) MeasurerOption {
	return func(o *MeasurerOptions) {
		o.LineHeightScale = lineHeightScale
	}
}

// NewMeasurer initializes and returns a clay.TextMeasurer using Gioui's engine.
func NewMeasurer(opts ...MeasurerOption) TextMeasurer {
	options := &MeasurerOptions{
		Collection:      gofont.Collection(),
		LineHeightScale: 1.0,
	}
	for _, opt := range opts {
		opt(options)
	}
	// Create the shaper, passing the font collection.
	return &measurer{
		shaper:  text.NewShaper(text.WithCollection(options.Collection)),
		options: options,
	}
}
func (m *measurer) MeasureText(textToMeasureSlice clay.Clay_StringSlice, cfg *clay.Clay_TextElementConfig, userData interface{}) clay.Clay_Dimensions {

	textToMeasure := textToMeasureSlice.String()

	// Map clay config to Gioui's layout parameters.
	params := text.Parameters{
		// Use a default font face (e.g., the first in the gofont collection).
		Font:            m.options.Collection[0].Font,
		Alignment:       text.Start,
		LineHeightScale: m.options.LineHeightScale,

		// Set a very large MaxWidth to ensure the text is measured as a single line
		// (unless the caller specifies max width via a different clay config field).
		MaxWidth: 1000000,
		PxPerEm:  fixed.Int26_6(cfg.FontSize),

		// Locale is needed for correct text direction/shaping (Bidi, complex scripts).
		Locale:     system.Locale{},
		WrapPolicy: text.WrapGraphemes,
	}

	// Perform the text layout.
	m.shaper.LayoutString(params, textToMeasure)

	var (
		lineStartX fixed.Int26_6
		maxWidth   fixed.Int26_6 = 0
		maxHeight  fixed.Int26_6 = 0
		lineCount  int           = 1 // The number of lines.

		isFirstGlyph = true
	)

	// Iterate through the shaped glyphs to find the dimensions.
	for {
		g, ok := m.shaper.NextGlyph()
		if !ok {
			break
		}
		if isFirstGlyph {
			lineStartX = g.X
			isFirstGlyph = false
		}

		// The width is the position of the glyph (g.X) plus its advance width (g.Advance).
		currentLineEnd := g.X + g.Advance
		if currentLineEnd > maxWidth {
			maxWidth = currentLineEnd
		}
		// fmt.Printf("g: %+v\n", g)
		height := g.Ascent + g.Descent + g.Offset.Y
		if height > maxHeight {
			maxHeight = height
		}
	}

	scaledMaxHeight := fixed.Int26_6(maxHeight) * fixed.Int26_6(m.options.LineHeightScale)
	if scaledMaxHeight <= 0 {
		// fmt.Printf("scaledMaxHeight <= 0, using default: %d, lineHeightScale: %f\n", cfg.FontSize, m.options.LineHeightScale)
		scaledMaxHeight = fixed.Int26_6(cfg.FontSize) * fixed.Int26_6(m.options.LineHeightScale)
	}
	totalHeight := float64(scaledMaxHeight) * float64(lineCount)
	totalWidth := float64(maxWidth - lineStartX)

	// fmt.Printf("totalWidth: %f, totalHeight: %f\n", totalWidth, totalHeight)

	// Convert the dimensions from fixed-point (1/64th units) back to float32 (DP).
	return clay.Clay_Dimensions{
		Width:  float32(totalWidth),
		Height: float32(totalHeight),
	}
}
