package gioui

import (
	"image"
	"image/color"
	"io"
	"log"
	"os"

	"gioui.org/font"
	"gioui.org/op"

	"github.com/zodimo/go-clay/clay"
)

// ExampleNewRenderer demonstrates custom renderer configuration
func ExampleNewRenderer() {
	ops := &op.Ops{}

	// Create renderer with custom logger (silent for tests)
	silentLogger := log.New(io.Discard, "", 0)
	renderer := NewRenderer(ops, RendererWithLogger(silentLogger))

	// Create renderer with multiple custom options
	customLogger := log.New(os.Stdout, "[MyApp] ", log.LstdFlags)
	rendererWithOptions := NewRenderer(ops,
		RendererWithLogger(customLogger),
		RendererWithDebugMode(true),
		RendererWithCacheSize(10*1024*1024), // 10MB cache
		RendererWithMaxClipDepth(50),
	)

	// Use the renderers...
	_ = renderer
	_ = rendererWithOptions
	// Output:
}

// ExampleRendererWithLogger demonstrates custom logger configuration
func ExampleRendererWithLogger() {
	ops := &op.Ops{}

	// Create a custom logger
	customLogger := log.New(os.Stderr, "[CustomRenderer] ", log.LstdFlags|log.Lshortfile)

	// Create renderer with custom logger
	renderer := NewRenderer(ops, RendererWithLogger(customLogger))

	// Use the renderer...
	_ = renderer
	// Output:
}

// ExampleRenderer_UnifiedRender demonstrates the unified Render method with multiple command types
func ExampleRenderer_UnifiedRender() {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	// Begin frame
	renderer.BeginFrame()

	// Create render commands with bounds
	commands := []clay.RenderCommand{
		// Rectangle with corner radius
		{
			BoundingBox: clay.BoundingBox{X: 10, Y: 10, Width: 200, Height: 100},
			CommandType: clay.CommandRectangle,
			ZIndex:      0,
			ID:          clay.ID("rounded-rect"),
			Data: clay.RectangleCommand{
				Color: clay.Color{R: 0.2, G: 0.6, B: 1.0, A: 1.0}, // Blue
				CornerRadius: clay.CornerRadius{
					TopLeft:     10.0,
					TopRight:    10.0,
					BottomLeft:  10.0,
					BottomRight: 10.0,
				},
			},
		},
		// Text with font support
		{
			BoundingBox: clay.BoundingBox{X: 30, Y: 50, Width: 160, Height: 30},
			CommandType: clay.CommandText,
			ZIndex:      1,
			ID:          clay.ID("hello-text"),
			Data: clay.TextCommand{
				Text:       "Hello, Clay!",
				FontSize:   24.0,
				FontID:     0,                                          // Default font
				Color:      clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0}, // White
				LineHeight: 1.2,
			},
		},
	}

	// Render all commands at once using unified method
	err := renderer.Render(commands)
	if err != nil {
		log.Printf("Render error: %v", err)
	}

	// End frame
	renderer.EndFrame()
	// Output:
}

// ExampleRenderer_CornerRadius demonstrates rendering rectangles with different corner radius values
func ExampleRenderer_CornerRadius() {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.BeginFrame()

	commands := []clay.RenderCommand{
		// Rectangle with all corners rounded
		{
			BoundingBox: clay.BoundingBox{X: 10, Y: 10, Width: 100, Height: 60},
			CommandType: clay.CommandRectangle,
			ZIndex:      0,
			Data: clay.RectangleCommand{
				Color: clay.Color{R: 1.0, G: 0.0, B: 0.0, A: 1.0}, // Red
				CornerRadius: clay.CornerRadius{
					TopLeft:     15.0,
					TopRight:    15.0,
					BottomLeft:  15.0,
					BottomRight: 15.0,
				},
			},
		},
		// Rectangle with only top corners rounded
		{
			BoundingBox: clay.BoundingBox{X: 120, Y: 10, Width: 100, Height: 60},
			CommandType: clay.CommandRectangle,
			ZIndex:      0,
			Data: clay.RectangleCommand{
				Color: clay.Color{R: 0.0, G: 1.0, B: 0.0, A: 1.0}, // Green
				CornerRadius: clay.CornerRadius{
					TopLeft:     20.0,
					TopRight:    20.0,
					BottomLeft:  0.0,
					BottomRight: 0.0,
				},
			},
		},
		// Simple rectangle without corner radius
		{
			BoundingBox: clay.BoundingBox{X: 230, Y: 10, Width: 100, Height: 60},
			CommandType: clay.CommandRectangle,
			ZIndex:      0,
			Data: clay.RectangleCommand{
				Color:        clay.Color{R: 0.0, G: 0.0, B: 1.0, A: 1.0}, // Blue
				CornerRadius: clay.CornerRadius{},                        // All zeros
			},
		},
	}

	renderer.Render(commands)
	renderer.EndFrame()
	// Output:
}

// ExampleRenderer_TextRendering demonstrates text rendering with FontManager
func ExampleRenderer_TextRendering() {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.BeginFrame()

	// Access FontManager to register custom fonts
	fontManager := renderer.GetFontManager()
	fontManager.RegisterFont(10, font.Regular, font.Bold) // Custom font ID 10

	commands := []clay.RenderCommand{
		// Regular text
		{
			BoundingBox: clay.BoundingBox{X: 10, Y: 10, Width: 200, Height: 30},
			CommandType: clay.CommandText,
			ZIndex:      0,
			Data: clay.TextCommand{
				Text:       "Regular Text",
				FontSize:   16.0,
				FontID:     0, // Default regular font
				Color:      clay.Color{R: 0.0, G: 0.0, B: 0.0, A: 1.0},
				LineHeight: 1.2,
			},
		},
		// Bold text
		{
			BoundingBox: clay.BoundingBox{X: 10, Y: 50, Width: 200, Height: 30},
			CommandType: clay.CommandText,
			ZIndex:      0,
			Data: clay.TextCommand{
				Text:       "Bold Text",
				FontSize:   18.0,
				FontID:     1, // Bold font
				Color:      clay.Color{R: 0.2, G: 0.2, B: 0.2, A: 1.0},
				LineHeight: 1.3,
			},
		},
		// Custom font
		{
			BoundingBox: clay.BoundingBox{X: 10, Y: 90, Width: 200, Height: 30},
			CommandType: clay.CommandText,
			ZIndex:      0,
			Data: clay.TextCommand{
				Text:       "Custom Font",
				FontSize:   20.0,
				FontID:     10, // Custom registered font
				Color:      clay.Color{R: 0.4, G: 0.4, B: 0.4, A: 1.0},
				LineHeight: 1.4,
			},
		},
	}

	renderer.Render(commands)
	renderer.EndFrame()
	// Output:
}

// ExampleRenderer_ImageRendering demonstrates image rendering with tint colors and corner radius
func ExampleRenderer_ImageRendering() {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.BeginFrame()

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, color.RGBA{R: 100, G: 150, B: 200, A: 255})
		}
	}

	commands := []clay.RenderCommand{
		// Image without tint (original colors)
		{
			BoundingBox: clay.BoundingBox{X: 10, Y: 10, Width: 100, Height: 100},
			CommandType: clay.CommandImage,
			ZIndex:      0,
			Data: clay.ImageCommand{
				ImageData:    img,
				TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0}, // White = no tint
				CornerRadius: clay.CornerRadius{},
			},
		},
		// Image with tint color
		{
			BoundingBox: clay.BoundingBox{X: 120, Y: 10, Width: 100, Height: 100},
			CommandType: clay.CommandImage,
			ZIndex:      0,
			Data: clay.ImageCommand{
				ImageData: img,
				TintColor: clay.Color{R: 1.0, G: 0.5, B: 0.0, A: 0.8}, // Orange tint with transparency
				CornerRadius: clay.CornerRadius{
					TopLeft:     10.0,
					TopRight:    10.0,
					BottomLeft:  10.0,
					BottomRight: 10.0,
				},
			},
		},
		// Image with rounded corners
		{
			BoundingBox: clay.BoundingBox{X: 230, Y: 10, Width: 100, Height: 100},
			CommandType: clay.CommandImage,
			ZIndex:      0,
			Data: clay.ImageCommand{
				ImageData:    img,
				TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
				CornerRadius: clay.CornerRadius{TopLeft: 20.0, TopRight: 20.0, BottomLeft: 0.0, BottomRight: 0.0},
			},
		},
	}

	renderer.Render(commands)
	renderer.EndFrame()
	// Output:
}

// ExampleRenderer_ComplexLayout demonstrates a complex layout with multiple features
func ExampleRenderer_ComplexLayout() {
	ops := &op.Ops{}
	renderer := NewRenderer(ops)

	renderer.BeginFrame()
	renderer.SetViewport(clay.BoundingBox{X: 0, Y: 0, Width: 800, Height: 600})

	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 80, 80))
	for y := 0; y < 80; y++ {
		for x := 0; x < 80; x++ {
			img.Set(x, y, color.RGBA{R: 200, G: 100, B: 50, A: 255})
		}
	}

	commands := []clay.RenderCommand{
		// Background container with rounded corners
		{
			BoundingBox: clay.BoundingBox{X: 50, Y: 50, Width: 300, Height: 200},
			CommandType: clay.CommandRectangle,
			ZIndex:      0,
			Data: clay.RectangleCommand{
				Color: clay.Color{R: 0.95, G: 0.95, B: 0.95, A: 1.0},
				CornerRadius: clay.CornerRadius{
					TopLeft:     12.0,
					TopRight:    12.0,
					BottomLeft:  12.0,
					BottomRight: 12.0,
				},
			},
		},
		// Header text
		{
			BoundingBox: clay.BoundingBox{X: 70, Y: 70, Width: 260, Height: 30},
			CommandType: clay.CommandText,
			ZIndex:      1,
			Data: clay.TextCommand{
				Text:       "Card Title",
				FontSize:   20.0,
				FontID:     1, // Bold
				Color:      clay.Color{R: 0.1, G: 0.1, B: 0.1, A: 1.0},
				LineHeight: 1.2,
			},
		},
		// Body text
		{
			BoundingBox: clay.BoundingBox{X: 70, Y: 110, Width: 260, Height: 60},
			CommandType: clay.CommandText,
			ZIndex:      1,
			Data: clay.TextCommand{
				Text:       "This is a complex layout example showing text, images, and rounded rectangles working together.",
				FontSize:   14.0,
				FontID:     0, // Regular
				Color:      clay.Color{R: 0.3, G: 0.3, B: 0.3, A: 1.0},
				LineHeight: 1.4,
			},
		},
		// Image with rounded corners
		{
			BoundingBox: clay.BoundingBox{X: 70, Y: 180, Width: 80, Height: 80},
			CommandType: clay.CommandImage,
			ZIndex:      1,
			Data: clay.ImageCommand{
				ImageData:    img,
				TintColor:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
				CornerRadius: clay.CornerRadius{TopLeft: 8.0, TopRight: 8.0, BottomLeft: 8.0, BottomRight: 8.0},
			},
		},
	}

	renderer.Render(commands)
	renderer.EndFrame()
	// Output:
}
