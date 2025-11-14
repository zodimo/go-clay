package claygio

import (
	"image/color"

	"github.com/zodimo/clay-go/clay"
)

// ClayToGioColor converts Clay's float32 RGBA color (0.0-1.0) to Gio's NRGBA color (0-255)
func ClayToGioColor(c clay.Clay_Color) color.NRGBA {
	return color.NRGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: uint8(c.A * 255),
	}
}

// mapClayCornerRadius maps Clay corner radius to gio-mw corner shapes
func MapClayCornerRadius(clayRadius clay.Clay_CornerRadius) CornerShapes {
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
