package gioui

import (
	"image"
	"image/color"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/f32"
)

// ImageFilter represents image scaling filter options
type ImageFilter int

const (
	FilterNearest ImageFilter = iota
	FilterLinear
)

// OperationBuilder provides utilities for building Gio operations efficiently
type OperationBuilder struct {
	ops *op.Ops
}

// NewOperationBuilder creates a new operation builder
func NewOperationBuilder(ops *op.Ops) *OperationBuilder {
	return &OperationBuilder{ops: ops}
}

// BuildImageOperation creates an image rendering operation with proper clipping and scaling
func (ob *OperationBuilder) BuildImageOperation(img image.Image, bounds image.Rectangle, filter ImageFilter, tintColor color.NRGBA) {
	// Set up clipping for image bounds
	clipStack := clip.Rect(bounds).Push(ob.ops)
	defer clipStack.Pop()

	// Create image operation with appropriate filter
	imageOp := paint.NewImageOp(img)
	switch filter {
	case FilterLinear:
		imageOp.Filter = paint.FilterLinear
	case FilterNearest:
		imageOp.Filter = paint.FilterNearest
	}

	// Apply tint color if specified
	if tintColor.A > 0 {
		paint.ColorOp{Color: tintColor}.Add(ob.ops)
	}

	// Add image operation and paint
	imageOp.Add(ob.ops)
	paint.PaintOp{}.Add(ob.ops)
}

// BuildBorderOperation creates a border rendering operation using stroke
func (ob *OperationBuilder) BuildBorderOperation(bounds image.Rectangle, width float32, borderColor color.NRGBA, cornerRadius float32) {
	// Create path for border
	var path clip.Path
	path.Begin(ob.ops)

	// Convert bounds to float32 for path operations
	minPt := f32.Pt(float32(bounds.Min.X), float32(bounds.Min.Y))
	maxPt := f32.Pt(float32(bounds.Max.X), float32(bounds.Max.Y))

	if cornerRadius > 0 {
		// Rounded rectangle path
		path.MoveTo(f32.Pt(minPt.X+cornerRadius, minPt.Y))
		path.LineTo(f32.Pt(maxPt.X-cornerRadius, minPt.Y))
		path.QuadTo(f32.Pt(maxPt.X, minPt.Y), f32.Pt(maxPt.X, minPt.Y+cornerRadius))
		path.LineTo(f32.Pt(maxPt.X, maxPt.Y-cornerRadius))
		path.QuadTo(f32.Pt(maxPt.X, maxPt.Y), f32.Pt(maxPt.X-cornerRadius, maxPt.Y))
		path.LineTo(f32.Pt(minPt.X+cornerRadius, maxPt.Y))
		path.QuadTo(f32.Pt(minPt.X, maxPt.Y), f32.Pt(minPt.X, maxPt.Y-cornerRadius))
		path.LineTo(f32.Pt(minPt.X, minPt.Y+cornerRadius))
		path.QuadTo(f32.Pt(minPt.X, minPt.Y), f32.Pt(minPt.X+cornerRadius, minPt.Y))
	} else {
		// Simple rectangle path
		path.MoveTo(minPt)
		path.LineTo(f32.Pt(maxPt.X, minPt.Y))
		path.LineTo(maxPt)
		path.LineTo(f32.Pt(minPt.X, maxPt.Y))
		path.Close()
	}

	// Create stroke operation
	strokeOp := clip.Stroke{
		Path:  path.End(),
		Width: width,
	}.Op()

	// Apply stroke with color
	clipStack := strokeOp.Push(ob.ops)
	paint.ColorOp{Color: borderColor}.Add(ob.ops)
	paint.PaintOp{}.Add(ob.ops)
	clipStack.Pop()
}

// BuildGradientOperation creates a linear gradient operation
func (ob *OperationBuilder) BuildGradientOperation(bounds image.Rectangle, startColor, endColor color.NRGBA, vertical bool) {
	clipStack := clip.Rect(bounds).Push(ob.ops)
	defer clipStack.Pop()

	// Set up gradient stops
	var stop1, stop2 f32.Point
	if vertical {
		stop1 = f32.Pt(0, float32(bounds.Min.Y))
		stop2 = f32.Pt(0, float32(bounds.Max.Y))
	} else {
		stop1 = f32.Pt(float32(bounds.Min.X), 0)
		stop2 = f32.Pt(float32(bounds.Max.X), 0)
	}

	gradientOp := paint.LinearGradientOp{
		Stop1:  stop1,
		Stop2:  stop2,
		Color1: startColor,
		Color2: endColor,
	}

	gradientOp.Add(ob.ops)
	paint.PaintOp{}.Add(ob.ops)
}

// BatchOperations groups similar operations for performance
type BatchOperations struct {
	rectangles []RectangleBatch
	images     []ImageBatch
	texts      []TextBatch
}

// RectangleBatch represents a batch of rectangle operations
type RectangleBatch struct {
	Bounds image.Rectangle
	Color  color.NRGBA
}

// ImageBatch represents a batch of image operations
type ImageBatch struct {
	Image  image.Image
	Bounds image.Rectangle
	Filter ImageFilter
	Tint   color.NRGBA
}

// TextBatch represents a batch of text operations
type TextBatch struct {
	Text   string
	Bounds image.Rectangle
	Color  color.NRGBA
}

// NewBatchOperations creates a new batch operations container
func NewBatchOperations() *BatchOperations {
	return &BatchOperations{
		rectangles: make([]RectangleBatch, 0),
		images:     make([]ImageBatch, 0),
		texts:      make([]TextBatch, 0),
	}
}

// AddRectangle adds a rectangle to the batch
func (bo *BatchOperations) AddRectangle(bounds image.Rectangle, color color.NRGBA) {
	bo.rectangles = append(bo.rectangles, RectangleBatch{
		Bounds: bounds,
		Color:  color,
	})
}

// AddImage adds an image to the batch
func (bo *BatchOperations) AddImage(img image.Image, bounds image.Rectangle, filter ImageFilter, tint color.NRGBA) {
	bo.images = append(bo.images, ImageBatch{
		Image:  img,
		Bounds: bounds,
		Filter: filter,
		Tint:   tint,
	})
}

// AddText adds text to the batch
func (bo *BatchOperations) AddText(text string, bounds image.Rectangle, color color.NRGBA) {
	bo.texts = append(bo.texts, TextBatch{
		Text:   text,
		Bounds: bounds,
		Color:  color,
	})
}

// Execute renders all batched operations
func (bo *BatchOperations) Execute(ops *op.Ops) {
	builder := NewOperationBuilder(ops)

	// Render all rectangles
	for _, rect := range bo.rectangles {
		clipStack := clip.Rect(rect.Bounds).Push(ops)
		paint.ColorOp{Color: rect.Color}.Add(ops)
		paint.PaintOp{}.Add(ops)
		clipStack.Pop()
	}

	// Render all images
	for _, img := range bo.images {
		builder.BuildImageOperation(img.Image, img.Bounds, img.Filter, img.Tint)
	}

	// Clear batches after execution
	bo.rectangles = bo.rectangles[:0]
	bo.images = bo.images[:0]
	bo.texts = bo.texts[:0]
}
