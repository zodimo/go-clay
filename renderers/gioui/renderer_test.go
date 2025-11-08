package gioui

import (
	"io"
	"log"
	"testing"

	"gioui.org/op"
	"github.com/zodimo/go-clay/clay"
)

// newSilentRenderer creates a renderer with silent error handling for tests
func newSilentRenderer(ops *op.Ops) *GioRenderer {
	// Create a logger that discards output (silent)
	silentLogger := log.New(io.Discard, "", 0)

	return NewRenderer(ops, RendererWithLogger(silentLogger))
}

func TestNewRenderer(t *testing.T) {
	ops := &op.Ops{}

	t.Run("Default options", func(t *testing.T) {
		renderer := NewRenderer(ops)
		if renderer == nil {
			t.Fatal("NewRenderer() returned nil")
		}
		if renderer.ops != ops {
			t.Error("NewRenderer() did not set ops correctly")
		}
	})

	t.Run("Custom logger", func(t *testing.T) {
		customLogger := log.New(io.Discard, "[CustomRenderer] ", log.LstdFlags)
		renderer := NewRenderer(ops, RendererWithLogger(customLogger))

		if renderer == nil {
			t.Fatal("NewRenderer() returned nil")
		}
		if renderer.ops != ops {
			t.Error("NewRenderer() did not set ops correctly")
		}
	})

	t.Run("Debug mode enabled", func(t *testing.T) {
		renderer := NewRenderer(ops, RendererWithDebugMode(true))

		if renderer == nil {
			t.Fatal("NewRenderer() returned nil")
		}
	})

	t.Run("Custom cache size", func(t *testing.T) {
		customCacheSize := 10 * 1024 * 1024 // 10MB
		renderer := NewRenderer(ops, RendererWithCacheSize(customCacheSize))

		if renderer == nil {
			t.Fatal("NewRenderer() returned nil")
		}
	})

	t.Run("Multiple options", func(t *testing.T) {
		customLogger := log.New(io.Discard, "[Test] ", log.LstdFlags)
		renderer := NewRenderer(ops,
			RendererWithLogger(customLogger),
			RendererWithDebugMode(true),
			RendererWithCacheSize(5*1024*1024),
			RendererWithMaxClipDepth(50),
		)

		if renderer == nil {
			t.Fatal("NewRenderer() returned nil")
		}
	})

	t.Run("Invalid options use defaults", func(t *testing.T) {
		// Test that invalid values are corrected to defaults
		renderer := NewRenderer(ops,
			RendererWithCacheSize(-1000), // Invalid negative size
			RendererWithMaxClipDepth(0),  // Invalid zero depth
		)

		if renderer == nil {
			t.Fatal("NewRenderer() returned nil")
		}
		// Should not panic or fail, defaults should be used
	})
}

func TestGioRenderer_BeginFrame(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	err := renderer.BeginFrame()
	if err != nil {
		t.Errorf("BeginFrame() returned error: %v", err)
	}
}

func TestGioRenderer_EndFrame(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	err := renderer.EndFrame()
	if err != nil {
		t.Errorf("EndFrame() returned error: %v", err)
	}
}

func TestGioRenderer_SetViewport(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	viewport := clay.BoundingBox{
		X:      0,
		Y:      0,
		Width:  800,
		Height: 600,
	}

	err := renderer.SetViewport(viewport)
	if err != nil {
		t.Errorf("SetViewport() returned error: %v", err)
	}

	if renderer.viewport != viewport {
		t.Errorf("SetViewport() did not set viewport correctly. Got %v, want %v", renderer.viewport, viewport)
	}
}

func TestGioRenderer_RenderRectangle(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	cmd := clay.RectangleCommand{
		Color: clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
		CornerRadius: clay.CornerRadius{
			TopLeft:     5,
			TopRight:    5,
			BottomLeft:  5,
			BottomRight: 5,
		},
	}

	err := renderer.RenderRectangle(cmd)
	if err != nil {
		t.Errorf("RenderRectangle() returned error: %v", err)
	}
}

func TestGioRenderer_RenderText(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	cmd := clay.TextCommand{
		Text:       "Hello, World!",
		FontID:     1,
		FontSize:   16.0,
		Color:      clay.Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
		LineHeight: 1.2,
	}

	err := renderer.RenderText(cmd)
	if err != nil {
		t.Errorf("RenderText() returned error: %v", err)
	}
}

func TestGioRenderer_RenderRectangle_NilOps(t *testing.T) {
	renderer := NewRenderer(nil)

	cmd := clay.RectangleCommand{
		Color: clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
	}

	err := renderer.RenderRectangle(cmd)
	if err == nil {
		t.Error("RenderRectangle() should return error when ops is nil")
	}
}

func TestGioRenderer_RenderText_NilOps(t *testing.T) {
	renderer := NewRenderer(nil)

	cmd := clay.TextCommand{
		Text:     "Hello",
		FontSize: 16.0,
		Color:    clay.Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
	}

	err := renderer.RenderText(cmd)
	if err == nil {
		t.Error("RenderText() should return error when ops is nil")
	}
}

func TestGioRenderer_StubMethods(t *testing.T) {
	ops := &op.Ops{}
	renderer := newSilentRenderer(ops)

	// Set viewport for border rendering
	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	// Test that stub methods return appropriate errors
	tests := []struct {
		name        string
		fn          func() error
		shouldError bool
	}{
		{
			name: "RenderImage",
			fn: func() error {
				return renderer.RenderImage(clay.ImageCommand{})
			},
			shouldError: true, // Should error due to nil image data
		},
		{
			name: "RenderBorder",
			fn: func() error {
				return renderer.RenderBorder(clay.BorderCommand{
					Color:        clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
					Width:        clay.BorderWidthAll(2.0),
					CornerRadius: clay.CornerRadiusAll(5.0),
				})
			},
			shouldError: false, // RenderBorder is now implemented
		},
		{
			name: "RenderClipStart",
			fn: func() error {
				return renderer.RenderClipStart(clay.ClipStartCommand{})
			},
			shouldError: false,
		},
		{
			name: "RenderClipEnd",
			fn: func() error {
				return renderer.RenderClipEnd(clay.ClipEndCommand{})
			},
			shouldError: true, // Should error because no clip operations to end (empty stack)
		},
		{
			name: "RenderCustom",
			fn: func() error {
				return renderer.RenderCustom(clay.CustomCommand{})
			},
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear clip stack before each test to ensure independence
			renderer.ClearClipStack()

			err := tt.fn()
			if tt.shouldError && err == nil {
				t.Errorf("%s should return an error (not yet implemented)", tt.name)
			} else if !tt.shouldError && err != nil {
				t.Errorf("%s should not return an error, got: %v", tt.name, err)
			}
		})
	}
}

// Test that GioRenderer implements clay.Renderer interface
func TestGioRenderer_ImplementsRenderer(t *testing.T) {
	ops := &op.Ops{}
	var renderer clay.Renderer = NewRenderer(ops)

	if renderer == nil {
		t.Error("GioRenderer should implement clay.Renderer interface")
	}
}
