package gioui

import (
	"image"
	"image/color"
	"io"
	"log"
	"os"
	"testing"

	"gioui.org/op"
	"gioui.org/op/clip"
	"github.com/zodimo/go-clay/clay"
)

func TestGioRenderer_RenderImage(t *testing.T) {
	renderer := createTestRenderer()
	silentRenderer := createSilentTestRenderer()

	t.Run("Valid image rendering", func(t *testing.T) {
		// Create a test image
		img := image.NewRGBA(image.Rect(0, 0, 100, 100))

		// Fill with red color for testing
		for y := 0; y < 100; y++ {
			for x := 0; x < 100; x++ {
				img.Set(x, y, color.RGBA{R: 255, G: 0, B: 0, A: 255})
			}
		}

		cmd := clay.ImageCommand{
			ImageData: img,
			TintColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0}, // White tint (no change)
			CornerRadius: clay.CornerRadius{
				TopLeft:     5.0,
				TopRight:    5.0,
				BottomLeft:  5.0,
				BottomRight: 5.0,
			},
		}

		err := renderer.RenderImage(cmd)
		if err != nil {
			t.Errorf("RenderImage failed: %v", err)
		}
	})

	t.Run("Image with tint color", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 50, 50))

		cmd := clay.ImageCommand{
			ImageData:    img,
			TintColor:    clay.Color{R: 1.0, G: 0.5, B: 0.0, A: 0.8}, // Orange tint
			CornerRadius: clay.CornerRadius{},
		}

		err := renderer.RenderImage(cmd)
		if err != nil {
			t.Errorf("RenderImage with tint failed: %v", err)
		}
	})

	t.Run("Nil image data", func(t *testing.T) {
		cmd := clay.ImageCommand{
			ImageData:    nil,
			TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			CornerRadius: clay.CornerRadius{},
		}

		err := silentRenderer.RenderImage(cmd)
		if err == nil {
			t.Error("Expected error for nil image data")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	})

	t.Run("Invalid corner radius", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 50, 50))

		cmd := clay.ImageCommand{
			ImageData: img,
			TintColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			CornerRadius: clay.CornerRadius{
				TopLeft:     -5.0, // Invalid negative value
				TopRight:    5.0,
				BottomLeft:  5.0,
				BottomRight: 5.0,
			},
		}

		err := silentRenderer.RenderImage(cmd)
		if err == nil {
			t.Error("Expected error for negative corner radius")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	})

	t.Run("Invalid color values", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 50, 50))

		cmd := clay.ImageCommand{
			ImageData:    img,
			TintColor:    clay.Color{R: 2.0, G: 1.0, B: 1.0, A: 1.0}, // Invalid R > 1.0
			CornerRadius: clay.CornerRadius{},
		}

		err := silentRenderer.RenderImage(cmd)
		if err == nil {
			t.Error("Expected error for invalid color values")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	})

	t.Run("Nil operations context", func(t *testing.T) {
		// Create renderer with nil ops
		logger := log.New(os.Stdout, "test: ", log.LstdFlags)
		nilOpsRenderer := &GioRenderer{
			ops:              nil, // Nil ops context
			cache:            NewResourceCache(1024 * 1024),
			operationBuilder: NewOperationBuilder(nil),
			errorHandler:     NewErrorHandler(logger, false),
		}

		img := image.NewRGBA(image.Rect(0, 0, 50, 50))
		cmd := clay.ImageCommand{
			ImageData:    img,
			TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			CornerRadius: clay.CornerRadius{},
		}

		err := nilOpsRenderer.RenderImage(cmd)
		if err == nil {
			t.Error("Expected error for nil operations context")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	})
}

func TestImageCaching(t *testing.T) {
	renderer := createTestRenderer()

	// Create two identical images
	img1 := image.NewRGBA(image.Rect(0, 0, 100, 100))
	img2 := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Fill both with the same pattern
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			c := color.RGBA{R: uint8(x), G: uint8(y), B: 128, A: 255}
			img1.Set(x, y, c)
			img2.Set(x, y, c)
		}
	}

	cmd1 := clay.ImageCommand{
		ImageData:    img1,
		TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		CornerRadius: clay.CornerRadius{},
	}

	cmd2 := clay.ImageCommand{
		ImageData:    img2,
		TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		CornerRadius: clay.CornerRadius{},
	}

	// Render first image
	err := renderer.RenderImage(cmd1)
	if err != nil {
		t.Fatalf("First RenderImage failed: %v", err)
	}

	// Render second image (should use cache if implemented correctly)
	err = renderer.RenderImage(cmd2)
	if err != nil {
		t.Fatalf("Second RenderImage failed: %v", err)
	}

	// Both should succeed - the caching behavior is tested in cache_test.go
}

func TestImageFormats(t *testing.T) {
	renderer := createTestRenderer()

	t.Run("RGBA image", func(t *testing.T) {
		img := image.NewRGBA(image.Rect(0, 0, 50, 50))
		cmd := clay.ImageCommand{
			ImageData:    img,
			TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			CornerRadius: clay.CornerRadius{},
		}

		err := renderer.RenderImage(cmd)
		if err != nil {
			t.Errorf("RGBA image rendering failed: %v", err)
		}
	})

	t.Run("Gray image", func(t *testing.T) {
		img := image.NewGray(image.Rect(0, 0, 50, 50))
		cmd := clay.ImageCommand{
			ImageData:    img,
			TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			CornerRadius: clay.CornerRadius{},
		}

		err := renderer.RenderImage(cmd)
		if err != nil {
			t.Errorf("Gray image rendering failed: %v", err)
		}
	})

	t.Run("Paletted image", func(t *testing.T) {
		palette := color.Palette{
			color.RGBA{0, 0, 0, 255},
			color.RGBA{255, 0, 0, 255},
			color.RGBA{0, 255, 0, 255},
			color.RGBA{0, 0, 255, 255},
		}
		img := image.NewPaletted(image.Rect(0, 0, 50, 50), palette)

		cmd := clay.ImageCommand{
			ImageData:    img,
			TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			CornerRadius: clay.CornerRadius{},
		}

		err := renderer.RenderImage(cmd)
		if err != nil {
			t.Errorf("Paletted image rendering failed: %v", err)
		}
	})
}

func BenchmarkRenderImage(b *testing.B) {
	renderer := createTestRenderer()

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 256, 256))
	for y := 0; y < 256; y++ {
		for x := 0; x < 256; x++ {
			img.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 128, A: 255})
		}
	}

	cmd := clay.ImageCommand{
		ImageData:    img,
		TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		CornerRadius: clay.CornerRadius{TopLeft: 10, TopRight: 10, BottomLeft: 10, BottomRight: 10},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Reset ops for each iteration to simulate real usage
		renderer.ops.Reset()
		err := renderer.RenderImage(cmd)
		if err != nil {
			b.Fatalf("RenderImage failed: %v", err)
		}
	}
}

func BenchmarkRenderImageCached(b *testing.B) {
	renderer := createTestRenderer()

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	cmd := clay.ImageCommand{
		ImageData:    img,
		TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		CornerRadius: clay.CornerRadius{},
	}

	// Pre-cache the image
	renderer.RenderImage(cmd)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		renderer.ops.Reset()
		err := renderer.RenderImage(cmd)
		if err != nil {
			b.Fatalf("RenderImage failed: %v", err)
		}
	}
}

// Helper function to create a test renderer
func createTestRenderer() *GioRenderer {
	ops := &op.Ops{}
	logger := log.New(os.Stdout, "test: ", log.LstdFlags)
	return &GioRenderer{
		ops:              ops,
		cache:            NewResourceCache(1024 * 1024),
		operationBuilder: NewOperationBuilder(ops),
		errorHandler:     NewErrorHandler(logger, false),
		viewport: clay.BoundingBox{
			X:      0,
			Y:      0,
			Width:  800,
			Height: 600,
		},
		clipStack:    make([]clip.Stack, 0),
		maxClipDepth: 100,
	}
}

// createSilentTestRenderer creates a test renderer with silent error handling
func createSilentTestRenderer() *GioRenderer {
	ops := &op.Ops{}
	silentLogger := log.New(io.Discard, "", 0)
	return &GioRenderer{
		ops:              ops,
		cache:            NewResourceCache(1024 * 1024),
		operationBuilder: NewOperationBuilder(ops),
		errorHandler:     NewErrorHandler(silentLogger, false),
		viewport: clay.BoundingBox{
			X:      0,
			Y:      0,
			Width:  800,
			Height: 600,
		},
		clipStack:    make([]clip.Stack, 0),
		maxClipDepth: 100,
	}
}
