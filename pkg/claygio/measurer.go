package claygio

import (
	"fmt"

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

	if cfg.LineHeight <= 0 {
		cfg.LineHeight = uint16(float32(cfg.FontSize) * m.options.LineHeightScale)
	}

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
		maxWidth  fixed.Int26_6 // The total width of the text (in 1/64th units).
		lineCount int           // The number of lines.
	)

	// Iterate through the shaped glyphs to find the dimensions.
	for {
		g, ok := m.shaper.NextGlyph()
		if !ok {
			break
		}

		if lineCount == 0 || g.Flags&text.FlagParagraphStart != 0 {
			lineCount++
		}

		// The width is the position of the glyph (g.X) plus its advance width (g.Advance).
		currentLineEnd := g.X + g.Advance
		if currentLineEnd > maxWidth {
			maxWidth = currentLineEnd
		}
	}

	// Calculate the accurate line height based on PxPerEm and LineHeightScale.
	// This is typically done outside of the loop.
	lineHeightFixed := float64(cfg.LineHeight) * float64(m.options.LineHeightScale)

	// Total height is the number of lines multiplied by the line height.
	if lineCount == 0 {
		lineCount = 1
	}
	totalHeightFixed := lineHeightFixed * float64(lineCount)

	fmt.Printf("maxWidth: %f, totalHeightFixed: %f\n", maxWidth, totalHeightFixed)

	// Convert the dimensions from fixed-point (1/64th units) back to float32 (DP).
	return clay.Clay_Dimensions{
		Width:  float32(maxWidth),
		Height: float32(totalHeightFixed),
	}
}
