package gioui

import (
	"testing"

	"gioui.org/op"
	"github.com/zodimo/go-clay/kimi2/clay"
)

func TestGioRenderer_RenderBorder_Basic(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Set viewport for border rendering
	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	cmd := clay.BorderCommand{
		Color:        clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
		Width:        clay.BorderWidthAll(2.0),
		CornerRadius: clay.CornerRadiusAll(0.0),
	}

	err := renderer.RenderBorder(cmd)
	if err != nil {
		t.Errorf("RenderBorder() returned error: %v", err)
	}
}

func TestGioRenderer_RenderBorder_WithCornerRadius(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	cmd := clay.BorderCommand{
		Color:        clay.Color{R: 0.0, G: 1.0, B: 0.0, A: 1.0},
		Width:        clay.BorderWidthAll(3.0),
		CornerRadius: clay.CornerRadiusAll(10.0),
	}

	err := renderer.RenderBorder(cmd)
	if err != nil {
		t.Errorf("RenderBorder() with corner radius returned error: %v", err)
	}
}

func TestGioRenderer_RenderBorder_DifferentWidths(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	cmd := clay.BorderCommand{
		Color:        clay.Color{R: 0.0, G: 0.0, B: 1.0, A: 1.0},
		Width:        clay.BorderWidth{Left: 1.0, Right: 2.0, Top: 3.0, Bottom: 4.0},
		CornerRadius: clay.CornerRadiusAll(5.0),
	}

	err := renderer.RenderBorder(cmd)
	if err != nil {
		t.Errorf("RenderBorder() with different widths returned error: %v", err)
	}
}

func TestGioRenderer_RenderBorder_DifferentCornerRadius(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	cmd := clay.BorderCommand{
		Color: clay.Color{R: 1.0, G: 1.0, B: 0.0, A: 1.0},
		Width: clay.BorderWidthAll(2.0),
		CornerRadius: clay.CornerRadius{
			TopLeft:     5.0,
			TopRight:    10.0,
			BottomLeft:  15.0,
			BottomRight: 20.0,
		},
	}

	err := renderer.RenderBorder(cmd)
	if err != nil {
		t.Errorf("RenderBorder() with different corner radius returned error: %v", err)
	}
}

func TestGioRenderer_RenderBorder_ZeroWidth(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	cmd := clay.BorderCommand{
		Color:        clay.Color{R: 1.0, G: 0.0, B: 1.0, A: 1.0},
		Width:        clay.BorderWidthAll(0.0),
		CornerRadius: clay.CornerRadiusAll(5.0),
	}

	err := renderer.RenderBorder(cmd)
	if err != nil {
		t.Errorf("RenderBorder() with zero width returned error: %v", err)
	}
}

func TestGioRenderer_RenderBorder_PartialZeroWidths(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	cmd := clay.BorderCommand{
		Color:        clay.Color{R: 0.5, G: 0.5, B: 0.5, A: 1.0},
		Width:        clay.BorderWidth{Left: 2.0, Right: 0.0, Top: 3.0, Bottom: 0.0},
		CornerRadius: clay.CornerRadiusAll(0.0),
	}

	err := renderer.RenderBorder(cmd)
	if err != nil {
		t.Errorf("RenderBorder() with partial zero widths returned error: %v", err)
	}
}

func TestGioRenderer_RenderBorder_Transparency(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	cmd := clay.BorderCommand{
		Color:        clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 0.5}, // Semi-transparent
		Width:        clay.BorderWidthAll(2.0),
		CornerRadius: clay.CornerRadiusAll(5.0),
	}

	err := renderer.RenderBorder(cmd)
	if err != nil {
		t.Errorf("RenderBorder() with transparency returned error: %v", err)
	}
}

func TestGioRenderer_RenderBorder_NilOps(t *testing.T) {
	renderer := NewRenderer(nil)

	cmd := clay.BorderCommand{
		Color:        clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
		Width:        clay.BorderWidthAll(2.0),
		CornerRadius: clay.CornerRadiusAll(5.0),
	}

	err := renderer.RenderBorder(cmd)
	if err == nil {
		t.Error("RenderBorder() should return error when ops is nil")
	}
}

func TestGioRenderer_RenderBorder_InvalidColor(t *testing.T) {
	ops := &op.Ops{}
	renderer := newSilentRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	tests := []struct {
		name  string
		color clay.Color
	}{
		{
			name:  "Negative red",
			color: clay.Color{R: -0.1, G: 0.0, B: 0.0, A: 1.0},
		},
		{
			name:  "Red too large",
			color: clay.Color{R: 1.1, G: 0.0, B: 0.0, A: 1.0},
		},
		{
			name:  "Negative alpha",
			color: clay.Color{R: 1.0, G: 0.0, B: 0.0, A: -0.1},
		},
		{
			name:  "Alpha too large",
			color: clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := clay.BorderCommand{
				Color:        tt.color,
				Width:        clay.BorderWidthAll(2.0),
				CornerRadius: clay.CornerRadiusAll(5.0),
			}

			err := renderer.RenderBorder(cmd)
			if err == nil {
				t.Errorf("RenderBorder() should return error for invalid color: %v", tt.color)
			}
		})
	}
}

func TestGioRenderer_RenderBorder_NegativeWidths(t *testing.T) {
	ops := &op.Ops{}
	renderer := newSilentRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	tests := []struct {
		name  string
		width clay.BorderWidth
	}{
		{
			name:  "Negative left",
			width: clay.BorderWidth{Left: -1.0, Right: 2.0, Top: 2.0, Bottom: 2.0},
		},
		{
			name:  "Negative right",
			width: clay.BorderWidth{Left: 2.0, Right: -1.0, Top: 2.0, Bottom: 2.0},
		},
		{
			name:  "Negative top",
			width: clay.BorderWidth{Left: 2.0, Right: 2.0, Top: -1.0, Bottom: 2.0},
		},
		{
			name:  "Negative bottom",
			width: clay.BorderWidth{Left: 2.0, Right: 2.0, Top: 2.0, Bottom: -1.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := clay.BorderCommand{
				Color:        clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
				Width:        tt.width,
				CornerRadius: clay.CornerRadiusAll(5.0),
			}

			err := renderer.RenderBorder(cmd)
			if err == nil {
				t.Errorf("RenderBorder() should return error for negative width: %v", tt.width)
			}
		})
	}
}

func TestGioRenderer_RenderBorder_NegativeCornerRadius(t *testing.T) {
	ops := &op.Ops{}
	renderer := newSilentRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	tests := []struct {
		name         string
		cornerRadius clay.CornerRadius
	}{
		{
			name:         "Negative top-left",
			cornerRadius: clay.CornerRadius{TopLeft: -1.0, TopRight: 5.0, BottomLeft: 5.0, BottomRight: 5.0},
		},
		{
			name:         "Negative top-right",
			cornerRadius: clay.CornerRadius{TopLeft: 5.0, TopRight: -1.0, BottomLeft: 5.0, BottomRight: 5.0},
		},
		{
			name:         "Negative bottom-left",
			cornerRadius: clay.CornerRadius{TopLeft: 5.0, TopRight: 5.0, BottomLeft: -1.0, BottomRight: 5.0},
		},
		{
			name:         "Negative bottom-right",
			cornerRadius: clay.CornerRadius{TopLeft: 5.0, TopRight: 5.0, BottomLeft: 5.0, BottomRight: -1.0},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := clay.BorderCommand{
				Color:        clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
				Width:        clay.BorderWidthAll(2.0),
				CornerRadius: tt.cornerRadius,
			}

			err := renderer.RenderBorder(cmd)
			if err == nil {
				t.Errorf("RenderBorder() should return error for negative corner radius: %v", tt.cornerRadius)
			}
		})
	}
}

func TestGioRenderer_RenderBorder_LargeValues(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 1000, Height: 1000})

	cmd := clay.BorderCommand{
		Color:        clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
		Width:        clay.BorderWidthAll(50.0),
		CornerRadius: clay.CornerRadiusAll(100.0),
	}

	err := renderer.RenderBorder(cmd)
	if err != nil {
		t.Errorf("RenderBorder() with large values returned error: %v", err)
	}
}

func TestGioRenderer_RenderBorder_SmallViewport(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 10, Height: 10})

	cmd := clay.BorderCommand{
		Color:        clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
		Width:        clay.BorderWidthAll(2.0),
		CornerRadius: clay.CornerRadiusAll(3.0),
	}

	err := renderer.RenderBorder(cmd)
	if err != nil {
		t.Errorf("RenderBorder() with small viewport returned error: %v", err)
	}
}

func TestGioRenderer_RenderBorder_MultipleRenders(t *testing.T) {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 100})

	// Render multiple borders to test state management
	borders := []clay.BorderCommand{
		{
			Color:        clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0},
			Width:        clay.BorderWidthAll(1.0),
			CornerRadius: clay.CornerRadiusAll(0.0),
		},
		{
			Color:        clay.Color{R: 0.0, G: 1.0, B: 0.0, A: 1.0},
			Width:        clay.BorderWidthAll(2.0),
			CornerRadius: clay.CornerRadiusAll(5.0),
		},
		{
			Color:        clay.Color{R: 0.0, G: 0.0, B: 1.0, A: 1.0},
			Width:        clay.BorderWidth{Left: 1.0, Right: 2.0, Top: 3.0, Bottom: 4.0},
			CornerRadius: clay.CornerRadius{TopLeft: 2.0, TopRight: 4.0, BottomLeft: 6.0, BottomRight: 8.0},
		},
	}

	for i, cmd := range borders {
		err := renderer.RenderBorder(cmd)
		if err != nil {
			t.Errorf("RenderBorder() call %d returned error: %v", i+1, err)
		}
	}
}
