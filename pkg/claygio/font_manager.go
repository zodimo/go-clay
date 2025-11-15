package claygio

import (
	"gioui.org/font"
	"gioui.org/font/gofont"
	"gioui.org/text"
)

// FontManager manages fonts and text shaping for the renderer

type FontManagerOptions struct {
	FontCollection []font.FontFace
}

type FontManagerOption func(*FontManagerOptions)

func FontManagerWithFontCollection(fontCollection []font.FontFace) FontManagerOption {
	return func(o *FontManagerOptions) {
		o.FontCollection = fontCollection
	}
}

type FontManager struct {
	shaper      *text.Shaper
	fontCache   map[uint16]font.Font
	defaultFont font.Font
}

// NewFontManager creates a new font manager with default fonts
func NewFontManager(opts ...FontManagerOption) *FontManager {
	options := &FontManagerOptions{
		FontCollection: gofont.Collection(),
	}
	for _, opt := range opts {
		opt(options)
	}

	// Create text shaper with default fonts
	shaperOptions := []text.ShaperOption{
		text.WithCollection(options.FontCollection),
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
