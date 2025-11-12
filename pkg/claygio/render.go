package claygio

import (
	"image"
	"image/color"
	"math"

	"gioui.org/f32"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"github.com/zodimo/clay-go/clay"
)

// map all render variations to op.ops

func Render(ops *op.Ops, renderCommand clay.Clay_RenderCommand) {
	switch renderCommand.CommandType {
	case clay.CLAY_RENDER_COMMAND_TYPE_RECTANGLE:
		RenderRectangle(ops, renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_BORDER:
		RenderBorder(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_TEXT:
		RenderText(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_IMAGE:
		RenderImage(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_SCISSOR_START:
		RenderScissorStart(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_SCISSOR_END:
		RenderScissorEnd(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_CUSTOM:
		RenderCustom(renderCommand)
	}
}

func RenderRectangle(ops *op.Ops, renderCommand clay.Clay_RenderCommand) {
	RenderRectangleWithBounds(ops, renderCommand)
}

func RenderBorder(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Border
	_ = spec
	return ops
}

func RenderText(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Text
	_ = spec
	return ops
}

func RenderImage(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Image
	_ = spec
	return ops
}

func RenderScissorStart(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Clip
	_ = spec
	return ops
}

func RenderScissorEnd(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Custom
	_ = spec
	return ops
}

func RenderCustom(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Custom
	_ = spec
	return ops

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

// isCornerRadiusZero checks if all corner radius values are zero
func IsCornerRadiusZero(radius clay.Clay_CornerRadius) bool {
	return radius.TopLeft == 0 && radius.TopRight == 0 &&
		radius.BottomLeft == 0 && radius.BottomRight == 0
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

type ShapedRect struct {
	MinPoint f32.Point
	MaxPoint f32.Point
	Offset   float32
	Shapes   CornerShapes
}

type CornerShapes struct {
	TopStart    CornerShape
	TopEnd      CornerShape
	BottomStart CornerShape
	BottomEnd   CornerShape
}

// Corner shape types for rounded rectangles
type CornerKind int

const (
	CornerKindDefault CornerKind = iota
	CornerKindChamfer
	CornerKindRound
)

type CornerShape struct {
	Kind        CornerKind
	Size        float32
	AdaptToSize bool
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

// Helper function for min
func min(a, b float32) float32 {
	if a < b {
		return a
	}
	return b
}
