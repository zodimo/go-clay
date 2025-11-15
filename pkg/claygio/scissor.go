package claygio

import (
	"encoding/json"
	"fmt"
)

// ScissorStart renders a scissor start operation
//

type ClippingContainer struct {
	X      float32
	Y      float32
	Width  float32
	Height float32

	Horizontal bool
	Vertical   bool
}

func (c *ClippingContainer) String() string {
	jsonString, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Sprintf("Error marshalling clipping container: %v", err)
	}
	return string(jsonString)
}

// // renderClipStartWithBounds starts clipping using bounds from RenderCommand
// func (r *renderer) renderClipStartWithBounds(ops *op.Ops, bounds clay.Clay_BoundingBox, cmd clay.Clay_RenderCommand) error {

// 	// Convert bounds to image rectangle
// 	rect := image.Rect(
// 		int(bounds.X),
// 		int(bounds.Y),
// 		int(bounds.X+bounds.Width),
// 		int(bounds.Y+bounds.Height),
// 	)

// 	// Create simple rectangular clipping (ClipStartCommand doesn't have CornerRadius)
// 	clipStack := clip.Rect(rect).Push(r.ops)

// 	// Add to clip stack
// 	r.clipStack = append(r.clipStack, clipStack)

// 	return nil
// }

// func (r *renderer) RenderScissorStart(ops *op.Ops, cmd clay.Clay_RenderCommand) error {

// 	// Check for clip stack overflow protection
// 	if len(r.clipStack) >= r.maxClipDepth {
// 		return NewRenderError(
// 			ErrorTypeClipStackOverflow,
// 			"RenderClipStart",
// 			fmt.Sprintf("Maximum clip depth of %d exceeded", r.maxClipDepth),
// 		)
// 	}

// 	// Create clipping bounds from current viewport
// 	// Note: In a complete implementation, this would use per-element bounds
// 	// passed from the layout engine, but the current interface limitation
// 	// requires us to use viewport bounds
// 	bounds := image.Rectangle{
// 		Min: image.Point{X: int(r.viewport.X), Y: int(r.viewport.Y)},
// 		Max: image.Point{
// 			X: int(r.viewport.X + r.viewport.Width),
// 			Y: int(r.viewport.Y + r.viewport.Height),
// 		},
// 	}

// 	// Apply directional clipping constraints
// 	if !cmd.Horizontal {
// 		// If horizontal clipping is disabled, extend bounds horizontally
// 		bounds.Min.X = -1000000 // Large negative value
// 		bounds.Max.X = 1000000  // Large positive value
// 	}
// 	if !cmd.Vertical {
// 		// If vertical clipping is disabled, extend bounds vertically
// 		bounds.Min.Y = -1000000 // Large negative value
// 		bounds.Max.Y = 1000000  // Large positive value
// 	}

// 	// Create and push clip operation
// 	clipOp := clip.Rect(bounds).Push(r.ops)

// 	// Add to clip stack for proper nesting
// 	r.clipStack = append(r.clipStack, clipOp)

// 	return nil
// }

// func (r *GioRenderer) RenderClipEnd(cmd clay.ClipEndCommand) error {
// 	defer func() {
// 		if err := r.errorHandler.RecoverFromPanic("RenderClipEnd"); err != nil {
// 			r.errorHandler.HandleError(err)
// 		}
// 	}()

// 	if r.ops == nil {
// 		return NewRenderError(
// 			ErrorTypeInvalidInput,
// 			"RenderClipEnd",
// 			"Operations context is nil",
// 		)
// 	}

// 	// Check if there are any clip operations to pop
// 	if len(r.clipStack) == 0 {
// 		return NewRenderError(
// 			ErrorTypeInvalidInput,
// 			"RenderClipEnd",
// 			"No clip operations to end - clip stack is empty",
// 		)
// 	}

// 	// Pop the most recent clip operation from the stack
// 	clipStackIndex := len(r.clipStack) - 1
// 	clipOp := r.clipStack[clipStackIndex]

// 	// Remove from our stack
// 	r.clipStack = r.clipStack[:clipStackIndex]

// 	// Pop the Gio clip operation
// 	clipOp.Pop()

// 	return nil
// }

// // CreateComplexClip creates a complex clipping shape with corner radius support
// func (r *GioRenderer) CreateComplexClip(bounds image.Rectangle, cornerRadius clay.CornerRadius) (clip.Stack, error) {
// 	if r.ops == nil {
// 		return clip.Stack{}, NewRenderError(
// 			ErrorTypeInvalidInput,
// 			"CreateComplexClip",
// 			"Operations context is nil",
// 		)
// 	}

// 	// Check for clip stack overflow protection
// 	if len(r.clipStack) >= r.maxClipDepth {
// 		return clip.Stack{}, NewRenderError(
// 			ErrorTypeClipStackOverflow,
// 			"CreateComplexClip",
// 			fmt.Sprintf("Maximum clip depth of %d exceeded", r.maxClipDepth),
// 		)
// 	}

// 	// If no corner radius, use simple rectangle clipping
// 	if cornerRadius.TopLeft == 0 && cornerRadius.TopRight == 0 &&
// 		cornerRadius.BottomLeft == 0 && cornerRadius.BottomRight == 0 {
// 		return clip.Rect(bounds).Push(r.ops), nil
// 	}

// 	// Create rounded rectangle path for complex clipping
// 	var path clip.Path
// 	path.Begin(r.ops)

// 	minPt := f32.Pt(float32(bounds.Min.X), float32(bounds.Min.Y))
// 	maxPt := f32.Pt(float32(bounds.Max.X), float32(bounds.Max.Y))

// 	// Start from top-left corner (after radius)
// 	path.MoveTo(f32.Pt(minPt.X+cornerRadius.TopLeft, minPt.Y))

// 	// Top edge to top-right corner
// 	if cornerRadius.TopRight > 0 {
// 		path.LineTo(f32.Pt(maxPt.X-cornerRadius.TopRight, minPt.Y))
// 		// Top-right corner arc
// 		path.QuadTo(f32.Pt(maxPt.X, minPt.Y), f32.Pt(maxPt.X, minPt.Y+cornerRadius.TopRight))
// 	} else {
// 		path.LineTo(f32.Pt(maxPt.X, minPt.Y))
// 	}

// 	// Right edge to bottom-right corner
// 	if cornerRadius.BottomRight > 0 {
// 		path.LineTo(f32.Pt(maxPt.X, maxPt.Y-cornerRadius.BottomRight))
// 		// Bottom-right corner arc
// 		path.QuadTo(f32.Pt(maxPt.X, maxPt.Y), f32.Pt(maxPt.X-cornerRadius.BottomRight, maxPt.Y))
// 	} else {
// 		path.LineTo(f32.Pt(maxPt.X, maxPt.Y))
// 	}

// 	// Bottom edge to bottom-left corner
// 	if cornerRadius.BottomLeft > 0 {
// 		path.LineTo(f32.Pt(minPt.X+cornerRadius.BottomLeft, maxPt.Y))
// 		// Bottom-left corner arc
// 		path.QuadTo(f32.Pt(minPt.X, maxPt.Y), f32.Pt(minPt.X, maxPt.Y-cornerRadius.BottomLeft))
// 	} else {
// 		path.LineTo(f32.Pt(minPt.X, maxPt.Y))
// 	}

// 	// Left edge to top-left corner
// 	if cornerRadius.TopLeft > 0 {
// 		path.LineTo(f32.Pt(minPt.X, minPt.Y+cornerRadius.TopLeft))
// 		// Top-left corner arc
// 		path.QuadTo(f32.Pt(minPt.X, minPt.Y), f32.Pt(minPt.X+cornerRadius.TopLeft, minPt.Y))
// 	} else {
// 		path.LineTo(f32.Pt(minPt.X, minPt.Y))
// 	}

// 	path.Close()

// 	// Create clip operation from path
// 	clipOp := clip.Outline{Path: path.End()}.Op().Push(r.ops)

// 	return clipOp, nil
// }

// // ClearClipStack safely clears all clip operations (for error recovery)
// func (r *GioRenderer) ClearClipStack() error {
// 	// Pop all remaining clip operations
// 	for len(r.clipStack) > 0 {
// 		clipStackIndex := len(r.clipStack) - 1
// 		clipOp := r.clipStack[clipStackIndex]
// 		r.clipStack = r.clipStack[:clipStackIndex]
// 		clipOp.Pop()
// 	}
// 	return nil
// }
