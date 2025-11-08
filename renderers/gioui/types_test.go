package gioui

import (
	"image/color"
	"testing"
	"gioui.org/f32"
	"github.com/zodimo/go-clay/clay"
)

func TestClayToGioColor(t *testing.T) {
	tests := []struct {
		name     string
		input    clay.Color
		expected color.NRGBA
	}{
		{
			name:     "Black color",
			input:    clay.Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
			expected: color.NRGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:     "White color",
			input:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			expected: color.NRGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:     "Red color with transparency",
			input:    clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 0.5},
			expected: color.NRGBA{R: 255, G: 0, B: 0, A: 127},
		},
		{
			name:     "Mid-range color",
			input:    clay.Color{R: 0.5, G: 0.25, B: 0.75, A: 0.8},
			expected: color.NRGBA{R: 127, G: 63, B: 191, A: 204},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClayToGioColor(tt.input)
			if result != tt.expected {
				t.Errorf("ClayToGioColor() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClayToGioPoint(t *testing.T) {
	tests := []struct {
		name     string
		x, y     float32
		expected f32.Point
	}{
		{
			name:     "Origin point",
			x:        0.0,
			y:        0.0,
			expected: f32.Point{X: 0.0, Y: 0.0},
		},
		{
			name:     "Positive coordinates",
			x:        100.5,
			y:        200.25,
			expected: f32.Point{X: 100.5, Y: 200.25},
		},
		{
			name:     "Negative coordinates",
			x:        -50.0,
			y:        -75.5,
			expected: f32.Point{X: -50.0, Y: -75.5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ClayToGioPoint(tt.x, tt.y)
			if result != tt.expected {
				t.Errorf("ClayToGioPoint() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestClayBoundsToGioRect(t *testing.T) {
	tests := []struct {
		name        string
		bounds      clay.BoundingBox
		expectedMin f32.Point
		expectedMax f32.Point
	}{
		{
			name:        "Origin rectangle",
			bounds:      clay.BoundingBox{X: 0, Y: 0, Width: 100, Height: 50},
			expectedMin: f32.Point{X: 0, Y: 0},
			expectedMax: f32.Point{X: 100, Y: 50},
		},
		{
			name:        "Offset rectangle",
			bounds:      clay.BoundingBox{X: 10, Y: 20, Width: 200, Height: 150},
			expectedMin: f32.Point{X: 10, Y: 20},
			expectedMax: f32.Point{X: 210, Y: 170},
		},
		{
			name:        "Negative position",
			bounds:      clay.BoundingBox{X: -10, Y: -5, Width: 30, Height: 25},
			expectedMin: f32.Point{X: -10, Y: -5},
			expectedMax: f32.Point{X: 20, Y: 20},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			min, max := ClayBoundsToGioRect(tt.bounds)
			if min != tt.expectedMin {
				t.Errorf("ClayBoundsToGioRect() min = %v, want %v", min, tt.expectedMin)
			}
			if max != tt.expectedMax {
				t.Errorf("ClayBoundsToGioRect() max = %v, want %v", max, tt.expectedMax)
			}
		})
	}
}
