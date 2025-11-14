package claygio

import "github.com/zodimo/clay-go/clay"

func MeasureText(text clay.Clay_StringSlice, config *clay.Clay_TextElementConfig, userData interface{}) clay.Clay_Dimensions {
	// Clay_TextElementConfig contains members such as fontId, fontSize, letterSpacing etc
	// Note: Clay_String->chars is not guaranteed to be null terminated
	return clay.Clay_Dimensions{
		Width:  float32(text.Length) * float32(config.FontSize), // <- this will only work for monospace fonts, see the renderers/ directory for more advanced text measurement
		Height: float32(config.FontSize),
	}
}
