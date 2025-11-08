package gioui

import (
	"image"
	"testing"

	"gioui.org/op"
	"github.com/zodimo/go-clay/clay"
)

func TestGioRenderer_RenderClipStart_Basic(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Set a viewport for clipping bounds
	err := renderer.SetViewport(clay.BoundingBox{
		X: 0, Y: 0, Width: 100, Height: 100,
	})
	if err != nil {
		t.Fatalf("SetViewport failed: %v", err)
	}

	// Test basic clipping start
	cmd := clay.ClipStartCommand{
		Horizontal: true,
		Vertical:   true,
	}

	err = renderer.RenderClipStart(cmd)
	if err != nil {
		t.Errorf("RenderClipStart failed: %v", err)
	}

	// Verify clip stack has one entry
	if len(renderer.clipStack) != 1 {
		t.Errorf("Expected clip stack length 1, got %d", len(renderer.clipStack))
	}
}

func TestGioRenderer_RenderClipStart_DirectionalClipping(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Set a viewport for clipping bounds
	err := renderer.SetViewport(clay.BoundingBox{
		X: 10, Y: 20, Width: 200, Height: 150,
	})
	if err != nil {
		t.Fatalf("SetViewport failed: %v", err)
	}

	tests := []struct {
		name       string
		horizontal bool
		vertical   bool
	}{
		{"Both directions", true, true},
		{"Horizontal only", true, false},
		{"Vertical only", false, true},
		{"Neither direction", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear any existing clips
			renderer.ClearClipStack()

			cmd := clay.ClipStartCommand{
				Horizontal: tt.horizontal,
				Vertical:   tt.vertical,
			}

			err := renderer.RenderClipStart(cmd)
			if err != nil {
				t.Errorf("RenderClipStart failed: %v", err)
			}

			// Verify clip stack has one entry
			if len(renderer.clipStack) != 1 {
				t.Errorf("Expected clip stack length 1, got %d", len(renderer.clipStack))
			}
		})
	}
}

func TestGioRenderer_RenderClipEnd_Basic(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Set a viewport for clipping bounds
	err := renderer.SetViewport(clay.BoundingBox{
		X: 0, Y: 0, Width: 100, Height: 100,
	})
	if err != nil {
		t.Fatalf("SetViewport failed: %v", err)
	}

	// Start a clip operation
	startCmd := clay.ClipStartCommand{
		Horizontal: true,
		Vertical:   true,
	}

	err = renderer.RenderClipStart(startCmd)
	if err != nil {
		t.Fatalf("RenderClipStart failed: %v", err)
	}

	// End the clip operation
	endCmd := clay.ClipEndCommand{}

	err = renderer.RenderClipEnd(endCmd)
	if err != nil {
		t.Errorf("RenderClipEnd failed: %v", err)
	}

	// Verify clip stack is empty
	if len(renderer.clipStack) != 0 {
		t.Errorf("Expected clip stack length 0, got %d", len(renderer.clipStack))
	}
}

func TestGioRenderer_RenderClipEnd_EmptyStack(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Try to end a clip operation without starting one
	endCmd := clay.ClipEndCommand{}

	err := renderer.RenderClipEnd(endCmd)
	if err == nil {
		t.Error("Expected RenderClipEnd to fail with empty stack")
	}

	// Verify error type
	if renderErr, ok := err.(*RenderError); ok {
		if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	} else {
		t.Errorf("Expected RenderError, got %T", err)
	}
}

func TestGioRenderer_HierarchicalClipping(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Set a viewport for clipping bounds
	err := renderer.SetViewport(clay.BoundingBox{
		X: 0, Y: 0, Width: 100, Height: 100,
	})
	if err != nil {
		t.Fatalf("SetViewport failed: %v", err)
	}

	// Start multiple nested clip operations
	for i := 0; i < 5; i++ {
		startCmd := clay.ClipStartCommand{
			Horizontal: true,
			Vertical:   true,
		}

		err = renderer.RenderClipStart(startCmd)
		if err != nil {
			t.Fatalf("RenderClipStart %d failed: %v", i, err)
		}

		// Verify clip stack depth
		expectedDepth := i + 1
		if len(renderer.clipStack) != expectedDepth {
			t.Errorf("Expected clip stack depth %d, got %d", expectedDepth, len(renderer.clipStack))
		}
	}

	// End clip operations in reverse order (LIFO)
	for i := 4; i >= 0; i-- {
		endCmd := clay.ClipEndCommand{}

		err = renderer.RenderClipEnd(endCmd)
		if err != nil {
			t.Fatalf("RenderClipEnd %d failed: %v", i, err)
		}

		// Verify clip stack depth
		expectedDepth := i
		if len(renderer.clipStack) != expectedDepth {
			t.Errorf("Expected clip stack depth %d, got %d", expectedDepth, len(renderer.clipStack))
		}
	}
}

func TestGioRenderer_ClipStackOverflow(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Set a viewport for clipping bounds
	err := renderer.SetViewport(clay.BoundingBox{
		X: 0, Y: 0, Width: 100, Height: 100,
	})
	if err != nil {
		t.Fatalf("SetViewport failed: %v", err)
	}

	// Fill up the clip stack to maximum depth
	maxDepth := renderer.maxClipDepth
	for i := 0; i < maxDepth; i++ {
		startCmd := clay.ClipStartCommand{
			Horizontal: true,
			Vertical:   true,
		}

		err = renderer.RenderClipStart(startCmd)
		if err != nil {
			t.Fatalf("RenderClipStart %d failed: %v", i, err)
		}
	}

	// Try to exceed maximum depth
	startCmd := clay.ClipStartCommand{
		Horizontal: true,
		Vertical:   true,
	}

	err = renderer.RenderClipStart(startCmd)
	if err == nil {
		t.Error("Expected RenderClipStart to fail with stack overflow")
	}

	// Verify error type
	if renderErr, ok := err.(*RenderError); ok {
		if renderErr.Type != ErrorTypeClipStackOverflow {
			t.Errorf("Expected ErrorTypeClipStackOverflow, got %v", renderErr.Type)
		}
	} else {
		t.Errorf("Expected RenderError, got %T", err)
	}
}

func TestGioRenderer_RenderClipStart_NilOps(t *testing.T) {
	renderer := NewRenderer(nil)

	startCmd := clay.ClipStartCommand{
		Horizontal: true,
		Vertical:   true,
	}

	err := renderer.RenderClipStart(startCmd)
	if err == nil {
		t.Error("Expected RenderClipStart to fail with nil ops")
	}

	// Verify error type
	if renderErr, ok := err.(*RenderError); ok {
		if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	} else {
		t.Errorf("Expected RenderError, got %T", err)
	}
}

func TestGioRenderer_RenderClipEnd_NilOps(t *testing.T) {
	renderer := NewRenderer(nil)

	endCmd := clay.ClipEndCommand{}

	err := renderer.RenderClipEnd(endCmd)
	if err == nil {
		t.Error("Expected RenderClipEnd to fail with nil ops")
	}

	// Verify error type
	if renderErr, ok := err.(*RenderError); ok {
		if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	} else {
		t.Errorf("Expected RenderError, got %T", err)
	}
}

func TestGioRenderer_CreateComplexClip(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Set a viewport for clipping bounds
	err := renderer.SetViewport(clay.BoundingBox{
		X: 0, Y: 0, Width: 100, Height: 100,
	})
	if err != nil {
		t.Fatalf("SetViewport failed: %v", err)
	}

	tests := []struct {
		name         string
		cornerRadius clay.CornerRadius
		expectError  bool
	}{
		{
			name: "No corner radius",
			cornerRadius: clay.CornerRadius{
				TopLeft: 0, TopRight: 0, BottomLeft: 0, BottomRight: 0,
			},
			expectError: false,
		},
		{
			name: "Uniform corner radius",
			cornerRadius: clay.CornerRadius{
				TopLeft: 10, TopRight: 10, BottomLeft: 10, BottomRight: 10,
			},
			expectError: false,
		},
		{
			name: "Different corner radii",
			cornerRadius: clay.CornerRadius{
				TopLeft: 5, TopRight: 10, BottomLeft: 15, BottomRight: 20,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bounds := image.Rectangle{
				Min: image.Point{X: int(renderer.viewport.X), Y: int(renderer.viewport.Y)},
				Max: image.Point{
					X: int(renderer.viewport.X + renderer.viewport.Width),
					Y: int(renderer.viewport.Y + renderer.viewport.Height),
				},
			}

			clipOp, err := renderer.CreateComplexClip(bounds, tt.cornerRadius)

			if tt.expectError {
				if err == nil {
					t.Error("Expected CreateComplexClip to fail")
				}
			} else {
				if err != nil {
					t.Errorf("CreateComplexClip failed: %v", err)
				}

				// Clean up the clip operation
				clipOp.Pop()
			}
		})
	}
}

func TestGioRenderer_ClearClipStack(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Set a viewport for clipping bounds
	err := renderer.SetViewport(clay.BoundingBox{
		X: 0, Y: 0, Width: 100, Height: 100,
	})
	if err != nil {
		t.Fatalf("SetViewport failed: %v", err)
	}

	// Add several clip operations
	for i := 0; i < 3; i++ {
		startCmd := clay.ClipStartCommand{
			Horizontal: true,
			Vertical:   true,
		}

		err = renderer.RenderClipStart(startCmd)
		if err != nil {
			t.Fatalf("RenderClipStart %d failed: %v", i, err)
		}
	}

	// Verify clip stack has entries
	if len(renderer.clipStack) != 3 {
		t.Errorf("Expected clip stack length 3, got %d", len(renderer.clipStack))
	}

	// Clear the clip stack
	err = renderer.ClearClipStack()
	if err != nil {
		t.Errorf("ClearClipStack failed: %v", err)
	}

	// Verify clip stack is empty
	if len(renderer.clipStack) != 0 {
		t.Errorf("Expected clip stack length 0 after clear, got %d", len(renderer.clipStack))
	}
}

func TestGioRenderer_ClippingWithRendering(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Set a viewport for clipping bounds
	err := renderer.SetViewport(clay.BoundingBox{
		X: 0, Y: 0, Width: 100, Height: 100,
	})
	if err != nil {
		t.Fatalf("SetViewport failed: %v", err)
	}

	// Start frame
	err = renderer.BeginFrame()
	if err != nil {
		t.Fatalf("BeginFrame failed: %v", err)
	}

	// Start clipping
	startCmd := clay.ClipStartCommand{
		Horizontal: true,
		Vertical:   true,
	}

	err = renderer.RenderClipStart(startCmd)
	if err != nil {
		t.Fatalf("RenderClipStart failed: %v", err)
	}

	// Render a rectangle within the clip
	rectCmd := clay.RectangleCommand{
		Color: clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
		CornerRadius: clay.CornerRadius{
			TopLeft: 5, TopRight: 5, BottomLeft: 5, BottomRight: 5,
		},
	}

	err = renderer.RenderRectangle(rectCmd)
	if err != nil {
		t.Errorf("RenderRectangle failed: %v", err)
	}

	// End clipping
	endCmd := clay.ClipEndCommand{}

	err = renderer.RenderClipEnd(endCmd)
	if err != nil {
		t.Errorf("RenderClipEnd failed: %v", err)
	}

	// End frame
	err = renderer.EndFrame()
	if err != nil {
		t.Errorf("EndFrame failed: %v", err)
	}

	// Verify clip stack is empty
	if len(renderer.clipStack) != 0 {
		t.Errorf("Expected clip stack length 0 after frame, got %d", len(renderer.clipStack))
	}
}
