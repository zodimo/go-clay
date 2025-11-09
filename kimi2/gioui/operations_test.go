package gioui

import (
	"image"
	"image/color"
	"testing"

	"gioui.org/op"
)

func TestOperationBuilder(t *testing.T) {
	var ops op.Ops
	builder := NewOperationBuilder(&ops)

	if builder == nil {
		t.Fatal("NewOperationBuilder returned nil")
	}

	if builder.ops != &ops {
		t.Error("OperationBuilder ops field not set correctly")
	}
}

func TestBuildImageOperation(t *testing.T) {
	var ops op.Ops
	builder := NewOperationBuilder(&ops)

	// Create a simple test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	bounds := image.Rect(10, 10, 110, 110)
	tintColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

	// Test with FilterLinear
	builder.BuildImageOperation(img, bounds, FilterLinear, tintColor)

	// Test with FilterNearest
	builder.BuildImageOperation(img, bounds, FilterNearest, tintColor)

	// Test with transparent tint (should not apply tint)
	transparentTint := color.NRGBA{R: 255, G: 0, B: 0, A: 0}
	builder.BuildImageOperation(img, bounds, FilterLinear, transparentTint)
}

func TestBuildBorderOperation(t *testing.T) {
	var ops op.Ops
	builder := NewOperationBuilder(&ops)

	bounds := image.Rect(0, 0, 100, 100)
	borderColor := color.NRGBA{R: 0, G: 0, B: 255, A: 255}

	// Test simple border (no corner radius)
	builder.BuildBorderOperation(bounds, 2.0, borderColor, 0)

	// Test rounded border
	builder.BuildBorderOperation(bounds, 3.0, borderColor, 10.0)

	// Test with different bounds
	largeBounds := image.Rect(50, 50, 200, 150)
	builder.BuildBorderOperation(largeBounds, 1.0, borderColor, 5.0)
}

func TestBuildGradientOperation(t *testing.T) {
	var ops op.Ops
	builder := NewOperationBuilder(&ops)

	bounds := image.Rect(0, 0, 100, 100)
	startColor := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	endColor := color.NRGBA{R: 0, G: 255, B: 0, A: 255}

	// Test vertical gradient
	builder.BuildGradientOperation(bounds, startColor, endColor, true)

	// Test horizontal gradient
	builder.BuildGradientOperation(bounds, startColor, endColor, false)
}

func TestBatchOperations(t *testing.T) {
	batch := NewBatchOperations()

	if batch == nil {
		t.Fatal("NewBatchOperations returned nil")
	}

	// Test adding rectangles
	bounds1 := image.Rect(0, 0, 50, 50)
	color1 := color.NRGBA{R: 255, G: 0, B: 0, A: 255}
	batch.AddRectangle(bounds1, color1)

	bounds2 := image.Rect(50, 50, 100, 100)
	color2 := color.NRGBA{R: 0, G: 255, B: 0, A: 255}
	batch.AddRectangle(bounds2, color2)

	if len(batch.rectangles) != 2 {
		t.Errorf("Expected 2 rectangles in batch, got %d", len(batch.rectangles))
	}

	// Test adding images
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	tint := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	batch.AddImage(img, bounds1, FilterLinear, tint)

	if len(batch.images) != 1 {
		t.Errorf("Expected 1 image in batch, got %d", len(batch.images))
	}

	// Test adding text
	batch.AddText("Test Text", bounds1, color1)

	if len(batch.texts) != 1 {
		t.Errorf("Expected 1 text in batch, got %d", len(batch.texts))
	}

	// Test execution
	var ops op.Ops
	batch.Execute(&ops)

	// After execution, batches should be cleared
	if len(batch.rectangles) != 0 {
		t.Errorf("Expected rectangles to be cleared after execution, got %d", len(batch.rectangles))
	}
	if len(batch.images) != 0 {
		t.Errorf("Expected images to be cleared after execution, got %d", len(batch.images))
	}
	if len(batch.texts) != 0 {
		t.Errorf("Expected texts to be cleared after execution, got %d", len(batch.texts))
	}
}

func TestImageFilter(t *testing.T) {
	tests := []struct {
		filter   ImageFilter
		expected ImageFilter
	}{
		{FilterNearest, FilterNearest},
		{FilterLinear, FilterLinear},
	}

	for _, test := range tests {
		if test.filter != test.expected {
			t.Errorf("ImageFilter value mismatch: got %v, expected %v", test.filter, test.expected)
		}
	}
}

func TestBatchTypes(t *testing.T) {
	// Test RectangleBatch
	rectBatch := RectangleBatch{
		Bounds: image.Rect(0, 0, 100, 100),
		Color:  color.NRGBA{R: 255, G: 0, B: 0, A: 255},
	}

	if rectBatch.Bounds.Empty() {
		t.Error("RectangleBatch bounds should not be empty")
	}

	// Test ImageBatch
	img := image.NewRGBA(image.Rect(0, 0, 50, 50))
	imgBatch := ImageBatch{
		Image:  img,
		Bounds: image.Rect(0, 0, 50, 50),
		Filter: FilterLinear,
		Tint:   color.NRGBA{R: 255, G: 255, B: 255, A: 255},
	}

	if imgBatch.Image == nil {
		t.Error("ImageBatch image should not be nil")
	}

	// Test TextBatch
	textBatch := TextBatch{
		Text:   "Test",
		Bounds: image.Rect(0, 0, 100, 20),
		Color:  color.NRGBA{R: 0, G: 0, B: 0, A: 255},
	}

	if textBatch.Text != "Test" {
		t.Errorf("TextBatch text mismatch: got %s, expected Test", textBatch.Text)
	}
}

func BenchmarkBuildImageOperation(b *testing.B) {
	var ops op.Ops
	builder := NewOperationBuilder(&ops)
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	bounds := image.Rect(0, 0, 100, 100)
	tint := color.NRGBA{R: 255, G: 255, B: 255, A: 255}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ops.Reset()
		builder.BuildImageOperation(img, bounds, FilterLinear, tint)
	}
}

func BenchmarkBuildBorderOperation(b *testing.B) {
	var ops op.Ops
	builder := NewOperationBuilder(&ops)
	bounds := image.Rect(0, 0, 100, 100)
	borderColor := color.NRGBA{R: 0, G: 0, B: 255, A: 255}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ops.Reset()
		builder.BuildBorderOperation(bounds, 2.0, borderColor, 5.0)
	}
}

func BenchmarkBatchOperations(b *testing.B) {
	var ops op.Ops
	batch := NewBatchOperations()
	bounds := image.Rect(0, 0, 100, 100)
	color1 := color.NRGBA{R: 255, G: 0, B: 0, A: 255}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Add operations to batch
		for j := 0; j < 10; j++ {
			batch.AddRectangle(bounds, color1)
		}
		
		// Execute batch
		batch.Execute(&ops)
		ops.Reset()
	}
}
