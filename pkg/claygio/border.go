package claygio

import (
	"image"
	"image/color"

	"gioui.org/f32"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/zodimo/clay-go/clay"
)

// RenderBorderWithBounds renders a border using bounds from RenderCommand
func RenderBorderWithBounds(cmd clay.Clay_RenderCommand) op.CallOp {
	var ops op.Ops
	opRecord := op.Record(&ops)
	borderData := cmd.RenderData.Border
	boundingBox := cmd.BoundingBox

	// Convert bounds to image rectangle
	rect := image.Rect(
		int(boundingBox.X),
		int(boundingBox.Y),
		int(boundingBox.X+boundingBox.Width),
		int(boundingBox.Y+boundingBox.Height),
	)

	// Convert Clay color to Gio color
	gioColor := ClayToGioColor(borderData.Color)

	// Render border using bounds
	callOp := RenderBorderSides(rect, borderData.Width, gioColor, borderData.CornerRadius)
	callOp.Add(&ops)

	return opRecord.Stop()
}

// RenderBorderSides renders borders with potentially different widths per side.
// For uniform-width borders, it uses Gio's clip.Stroke with RRect for optimal rendering.
// For per-side widths, it renders each side individually.
func RenderBorderSides(bounds image.Rectangle, width clay.Clay_BorderWidth, color color.NRGBA, cornerRadius clay.Clay_CornerRadius) op.CallOp {
	var ops op.Ops
	opRecord := op.Record(&ops)

	// Check if all border widths are uniform
	uniformWidth := width.Top == width.Right && width.Right == width.Bottom && width.Bottom == width.Left && width.Top > 0

	if uniformWidth {
		// Use Gio's built-in stroke rendering for uniform-width borders
		RenderUniformBorder(&ops, bounds, int(width.Top), color, cornerRadius)
	} else {
		// Render each side individually for different widths
		RenderPerSideBorders(&ops, bounds, width, color, cornerRadius)
	}

	return opRecord.Stop()
}

// RenderUniformBorder renders a border with uniform width using Gio's clip.Stroke and RRect.
// This is the optimal path for borders where all sides have the same width.
func RenderUniformBorder(ops *op.Ops, bounds image.Rectangle, width int, color color.NRGBA, cornerRadius clay.Clay_CornerRadius) {
	// Create RRect with per-corner radii
	// Gio's RRect uses: SE (bottom-right), SW (bottom-left), NW (top-left), NE (top-right)
	rrect := clip.RRect{
		Rect: bounds,
		SE:   int(cornerRadius.BottomRight), // South-East (bottom-right)
		SW:   int(cornerRadius.BottomLeft),  // South-West (bottom-left)
		NW:   int(cornerRadius.TopLeft),     // North-West (top-left)
		NE:   int(cornerRadius.TopRight),    // North-East (top-right)
	}

	// Use Gio's built-in stroke rendering
	paint.FillShape(ops,
		color,
		clip.Stroke{
			Path:  rrect.Path(ops),
			Width: float32(width),
		}.Op(),
	)
}

// RenderPerSideBorders renders borders where each side can have a different width.
// This uses individual path rendering for each side.
func RenderPerSideBorders(ops *op.Ops, bounds image.Rectangle, width clay.Clay_BorderWidth, color color.NRGBA, cornerRadius clay.Clay_CornerRadius) {
	// Top border
	if width.Top > 0 {
		RenderBorderSide(ops, bounds, width.Top, color, cornerRadius.TopLeft, cornerRadius.TopRight, true, false)
	}

	// Right border
	if width.Right > 0 {
		RenderBorderSide(ops, bounds, width.Right, color, cornerRadius.TopRight, cornerRadius.BottomRight, false, true)
	}

	// Bottom border
	if width.Bottom > 0 {
		RenderBorderSide(ops, bounds, width.Bottom, color, cornerRadius.BottomLeft, cornerRadius.BottomRight, true, true)
	}

	// Left border
	if width.Left > 0 {
		RenderBorderSide(ops, bounds, width.Left, color, cornerRadius.TopLeft, cornerRadius.BottomLeft, false, false)
	}
}

// RenderBorderSide renders a single border side with corner radius support.
// horizontal: true for top/bottom borders, false for left/right borders
// isBottomOrRight: true for bottom/right borders, false for top/left borders
func RenderBorderSide(ops *op.Ops, bounds image.Rectangle, sideWidth uint16, color color.NRGBA, startRadius, endRadius float32, horizontal, isBottomOrRight bool) {
	var path clip.Path
	path.Begin(ops)

	width := float32(sideWidth)
	minX := float32(bounds.Min.X)
	minY := float32(bounds.Min.Y)
	maxX := float32(bounds.Max.X)
	maxY := float32(bounds.Max.Y)

	if horizontal {
		// Horizontal border (top or bottom)
		if isBottomOrRight {
			// Bottom border
			y := maxY - width
			path.MoveTo(f32.Pt(minX+startRadius, y))
			path.LineTo(f32.Pt(maxX-endRadius, y))
			if endRadius > 0 {
				// Bottom-right corner
				path.QuadTo(f32.Pt(maxX, y), f32.Pt(maxX, y+endRadius))
			} else {
				path.LineTo(f32.Pt(maxX, y))
			}
			path.LineTo(f32.Pt(maxX, maxY))
			path.LineTo(f32.Pt(minX, maxY))
			if startRadius > 0 {
				// Bottom-left corner
				path.LineTo(f32.Pt(minX, y+startRadius))
				path.QuadTo(f32.Pt(minX, y), f32.Pt(minX+startRadius, y))
			} else {
				path.LineTo(f32.Pt(minX, y))
			}
		} else {
			// Top border
			y := minY
			path.MoveTo(f32.Pt(minX+startRadius, y))
			path.LineTo(f32.Pt(maxX-endRadius, y))
			if endRadius > 0 {
				// Top-right corner
				path.QuadTo(f32.Pt(maxX, y), f32.Pt(maxX, y+endRadius))
			} else {
				path.LineTo(f32.Pt(maxX, y))
			}
			path.LineTo(f32.Pt(maxX, y+width))
			path.LineTo(f32.Pt(minX, y+width))
			if startRadius > 0 {
				// Top-left corner
				path.LineTo(f32.Pt(minX, y+startRadius))
				path.QuadTo(f32.Pt(minX, y), f32.Pt(minX+startRadius, y))
			} else {
				path.LineTo(f32.Pt(minX, y))
			}
		}
	} else {
		// Vertical border (left or right)
		if isBottomOrRight {
			// Right border
			x := maxX - width
			path.MoveTo(f32.Pt(x, minY+startRadius))
			path.LineTo(f32.Pt(x, maxY-endRadius))
			if endRadius > 0 {
				// Bottom-right corner
				path.QuadTo(f32.Pt(x, maxY), f32.Pt(x+endRadius, maxY))
			} else {
				path.LineTo(f32.Pt(x, maxY))
			}
			path.LineTo(f32.Pt(maxX, maxY))
			path.LineTo(f32.Pt(maxX, minY))
			if startRadius > 0 {
				// Top-right corner
				path.LineTo(f32.Pt(x+startRadius, minY))
				path.QuadTo(f32.Pt(x, minY), f32.Pt(x, minY+startRadius))
			} else {
				path.LineTo(f32.Pt(x, minY))
			}
		} else {
			// Left border
			x := minX
			path.MoveTo(f32.Pt(x, minY+startRadius))
			path.LineTo(f32.Pt(x, maxY-endRadius))
			if endRadius > 0 {
				// Bottom-left corner
				path.QuadTo(f32.Pt(x, maxY), f32.Pt(x+endRadius, maxY))
			} else {
				path.LineTo(f32.Pt(x, maxY))
			}
			path.LineTo(f32.Pt(x+width, maxY))
			path.LineTo(f32.Pt(x+width, minY))
			if startRadius > 0 {
				// Top-left corner
				path.LineTo(f32.Pt(x+startRadius, minY))
				path.QuadTo(f32.Pt(x, minY), f32.Pt(x, minY+startRadius))
			} else {
				path.LineTo(f32.Pt(x, minY))
			}
		}
	}

	path.Close()

	// Fill the border area
	clipStack := clip.Outline{Path: path.End()}.Op().Push(ops)
	paint.ColorOp{Color: color}.Add(ops)
	paint.PaintOp{}.Add(ops)
	clipStack.Pop()
}
