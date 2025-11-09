package gioui

import (
	"image/color"

	"gioui.org/f32"
	"github.com/zodimo/go-clay/kimi2/clay"
)

// ClayToGioColor converts Clay's float32 RGBA color (0.0-1.0) to Gio's NRGBA color (0-255)
func ClayToGioColor(c clay.Color) color.NRGBA {
	return color.NRGBA{
		R: uint8(c.R * 255),
		G: uint8(c.G * 255),
		B: uint8(c.B * 255),
		A: uint8(c.A * 255),
	}
}

// ClayToGioPoint converts Clay's integer coordinates to Gio's float32 point
func ClayToGioPoint(x, y float32) f32.Point {
	return f32.Point{
		X: x,
		Y: y,
	}
}

// ClayBoundsToGioRect converts Clay's BoundingBox to Gio's rectangle points
func ClayBoundsToGioRect(bounds clay.BoundingBox) (min, max f32.Point) {
	min = f32.Point{
		X: bounds.X,
		Y: bounds.Y,
	}
	max = f32.Point{
		X: bounds.X + bounds.Width,
		Y: bounds.Y + bounds.Height,
	}
	return min, max
}
