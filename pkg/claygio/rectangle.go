package claygio

import (
	"image"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/zodimo/clay-go/clay"
)

type ShapedRect struct {
	MinPoint f32.Point
	MaxPoint f32.Point
	Offset   float32
	Shapes   CornerShapes
}

// renderRectangleWithBounds renders a rectangle using bounds from RenderCommand
func RenderRectangleWithBounds(ops *op.Ops, cmd clay.Clay_RenderCommand) error {

	rectangleData := cmd.RenderData.Rectangle
	boundingBox := cmd.BoundingBox
	// Check if corner radius is needed
	if IsCornerRadiusZero(rectangleData.CornerRadius) {
		return RenderSimpleRectangle(ops, boundingBox, cmd)
	}

	// Render rectangle with corner radius
	return RenderRoundedRectangle(ops, boundingBox, cmd)
}

// renderSimpleRectangle renders a rectangle without corner radius
func RenderSimpleRectangle(ops *op.Ops, bounds clay.Clay_BoundingBox, cmd clay.Clay_RenderCommand) error {
	rectangleData := cmd.RenderData.Rectangle
	// Convert bounds to Gio rectangle
	rect := image.Rect(
		int(bounds.X),
		int(bounds.Y),
		int(bounds.X+bounds.Width),
		int(bounds.Y+bounds.Height),
	)

	// Create clipping region
	clipOp := clip.Rect(rect).Push(ops)
	defer clipOp.Pop()

	// Apply color and paint
	paint.ColorOp{Color: ClayToGioColor(rectangleData.BackgroundColor)}.Add(ops)
	paint.PaintOp{}.Add(ops)

	return nil
}

// Path generates a clip path for the shaped rectangle
func (s ShapedRect) Path(gtx layout.Context) clip.PathSpec {
	rMinP := s.MinPoint.Sub(f32.Pt(s.Offset, s.Offset))
	rMaxP := s.MaxPoint.Add(f32.Pt(s.Offset, s.Offset))

	rHeight := rMaxP.Y - rMinP.Y
	rWidth := rMaxP.X - rMinP.X
	rMinDim := min(rHeight, rWidth)
	if rMinDim <= 0 {
		return clip.PathSpec{}
	}

	// Use shapes as-is (simplified - no RTL handling for now)
	corners := s.Shapes

	// Calculate actual corner sizes
	ts := min(corners.TopStart.Size, rMinDim)
	te := min(corners.TopEnd.Size, rMinDim)
	be := min(corners.BottomEnd.Size, rMinDim)
	bs := min(corners.BottomStart.Size, rMinDim)

	// Build path with rounded corners
	var path clip.Path
	path.Begin(gtx.Ops)

	// Top edge
	path.MoveTo(f32.Point{X: rMinP.X + ts, Y: rMinP.Y})
	path.LineTo(f32.Point{X: rMaxP.X - te, Y: rMinP.Y})

	// Top-right corner
	if te > 0 && corners.TopEnd.Kind == CornerKindRound {
		fPoint := f32.Point{X: rMaxP.X - te, Y: rMinP.Y + te}
		path.ArcTo(fPoint, fPoint, math.Pi/2)
	} else if te > 0 {
		path.LineTo(f32.Point{X: rMaxP.X, Y: rMinP.Y + te})
	}

	// Right edge
	path.LineTo(f32.Point{X: rMaxP.X, Y: rMaxP.Y - be})

	// Bottom-right corner
	if be > 0 && corners.BottomEnd.Kind == CornerKindRound {
		fPoint := f32.Point{X: rMaxP.X - be, Y: rMaxP.Y - be}
		path.ArcTo(fPoint, fPoint, math.Pi/2)
	} else if be > 0 {
		path.LineTo(f32.Point{X: rMaxP.X - be, Y: rMaxP.Y})
	}

	// Bottom edge
	path.LineTo(f32.Point{X: rMinP.X + bs, Y: rMaxP.Y})

	// Bottom-left corner
	if bs > 0 && corners.BottomStart.Kind == CornerKindRound {
		fPoint := f32.Point{X: rMinP.X + bs, Y: rMaxP.Y - bs}
		path.ArcTo(fPoint, fPoint, math.Pi/2)
	} else if bs > 0 {
		path.LineTo(f32.Point{X: rMinP.X, Y: rMaxP.Y - bs})
	}

	// Left edge
	path.LineTo(f32.Point{X: rMinP.X, Y: rMinP.Y + ts})

	// Top-left corner
	if ts > 0 && corners.TopStart.Kind == CornerKindRound {
		fPoint := f32.Point{X: rMinP.X + ts, Y: rMinP.Y + ts}
		path.ArcTo(fPoint, fPoint, math.Pi/2)
	} else if ts > 0 {
		path.LineTo(f32.Point{X: rMinP.X + ts, Y: rMinP.Y})
	}

	path.Close()
	return path.End()
}

// renderRoundedRectangle renders a rectangle with corner radius
func RenderRoundedRectangle(ops *op.Ops, bounds clay.Clay_BoundingBox, cmd clay.Clay_RenderCommand) error {
	rectangleData := cmd.RenderData.Rectangle
	// Create shaped rectangle with corner radius
	shapedRect := ShapedRect{
		MinPoint: f32.Pt(float32(bounds.X), float32(bounds.Y)),
		MaxPoint: f32.Pt(float32(bounds.X+bounds.Width), float32(bounds.Y+bounds.Height)),
		Shapes:   MapClayCornerRadius(rectangleData.CornerRadius),
	}

	// Create layout context for path generation
	gtx := layout.Context{
		Ops: ops,
		Constraints: layout.Constraints{
			Max: image.Pt(int(bounds.Width), int(bounds.Height)),
		},
	}

	// Create clipping path with rounded corners
	pathSpec := shapedRect.Path(gtx)
	clipOp := clip.Outline{Path: pathSpec}.Op().Push(ops)
	defer clipOp.Pop()

	// Apply color and paint
	paint.ColorOp{Color: ClayToGioColor(rectangleData.BackgroundColor)}.Add(ops)
	paint.PaintOp{}.Add(ops)

	return nil
}
