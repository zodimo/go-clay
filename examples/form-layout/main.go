package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-clay/clay"
	"github.com/zodimo/go-clay/renderers/gioui"
)

func main() {
	go func() {
		w := &app.Window{}
		w.Option(
			app.Title("Clay Form Layout Example"),
			app.Size(unit.Dp(500), unit.Dp(600)),
		)

		if err := run(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func run(w *app.Window) error {
	var ops op.Ops

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			// Create Clay layout engine and Gio renderer
			engine := clay.NewLayoutEngine()
			renderer := gioui.NewRenderer(&ops)

			// Set viewport
			viewport := clay.BoundingBox{
				X:      0,
				Y:      0,
				Width:  float32(gtx.Constraints.Max.X),
				Height: float32(gtx.Constraints.Max.Y),
			}
			renderer.SetViewport(viewport)

			// Set layout dimensions
			engine.SetLayoutDimensions(clay.Dimensions{
				Width:  float32(gtx.Constraints.Max.X),
				Height: float32(gtx.Constraints.Max.Y),
			})

			// Begin layout
			engine.BeginLayout()

			// Create form layout
			createForm(engine)

			// End layout and get render commands
			commands := engine.EndLayout()

			// Render with Gio using unified Render method
			renderer.BeginFrame()
			if err := renderer.Render(commands); err != nil {
				log.Printf("Render error: %v", err)
			}
			renderer.EndFrame()

			e.Frame(gtx.Ops)
		}
	}
}

func createForm(engine clay.LayoutEngine) {
	// Form container
	engine.OpenElement(clay.ID("form"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingFixed(400),
				Height: clay.SizingGrow(0),
			},
			Direction: clay.TopToBottom,
			ChildGap:  16,
			Padding:  clay.PaddingAll(20),
			ChildAlignment: clay.ChildAlignment{
				X: clay.AlignXCenter,
				Y: clay.AlignYCenter,
			},
		},
		BackgroundColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		CornerRadius: clay.CornerRadius{
			TopLeft:     8.0,
			TopRight:    8.0,
			BottomLeft:  8.0,
			BottomRight: 8.0,
		},
	})

	// Form title
	engine.OpenElement(clay.ID("form-title"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   24,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0, G: 0, B: 0, A: 1.0},
			FontID:     1, // Bold font
		},
	})
	engine.CloseElement()

	// Name field container
	engine.OpenElement(clay.ID("name-field"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingFixed(60),
			},
			Direction: clay.TopToBottom,
			ChildGap:   4,
		},
	})

	// Name label
	engine.OpenElement(clay.ID("name-label"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   14,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0.3, G: 0.3, B: 0.3, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()

	// Name input
	engine.OpenElement(clay.ID("name-input"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingFixed(40),
			},
			Padding: clay.PaddingAll(8),
		},
		BackgroundColor: clay.Color{R: 0.95, G: 0.95, B: 0.95, A: 1.0},
		Border: &clay.BorderConfig{
			Width: clay.BorderWidthAll(1),
			Color: clay.Color{R: 0.8, G: 0.8, B: 0.8, A: 1.0},
		},
		CornerRadius: clay.CornerRadius{
			TopLeft:     4.0,
			TopRight:    4.0,
			BottomLeft:  4.0,
			BottomRight: 4.0,
		},
	})

	// Name input placeholder text
	engine.OpenElement(clay.ID("name-placeholder"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0.6, G: 0.6, B: 0.6, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()
	engine.CloseElement() // Close name input
	engine.CloseElement() // Close name field

	// Email field container
	engine.OpenElement(clay.ID("email-field"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingFixed(60),
			},
			Direction: clay.TopToBottom,
			ChildGap:   4,
		},
	})

	// Email label
	engine.OpenElement(clay.ID("email-label"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   14,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0.3, G: 0.3, B: 0.3, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()

	// Email input
	engine.OpenElement(clay.ID("email-input"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingFixed(40),
			},
			Padding: clay.PaddingAll(8),
		},
		BackgroundColor: clay.Color{R: 0.95, G: 0.95, B: 0.95, A: 1.0},
		Border: &clay.BorderConfig{
			Width: clay.BorderWidthAll(1),
			Color: clay.Color{R: 0.8, G: 0.8, B: 0.8, A: 1.0},
		},
		CornerRadius: clay.CornerRadius{
			TopLeft:     4.0,
			TopRight:    4.0,
			BottomLeft:  4.0,
			BottomRight: 4.0,
		},
	})

	// Email input placeholder text
	engine.OpenElement(clay.ID("email-placeholder"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0.6, G: 0.6, B: 0.6, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()
	engine.CloseElement() // Close email input
	engine.CloseElement() // Close email field

	// Submit button
	engine.OpenElement(clay.ID("submit-button"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingFixed(120),
				Height: clay.SizingFixed(40),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.AlignXCenter,
				Y: clay.AlignYCenter,
			},
		},
		BackgroundColor: clay.Color{R: 0.2, G: 0.6, B: 1.0, A: 1.0},
		CornerRadius: clay.CornerRadiusAll(4),
	})

	// Submit button text
	engine.OpenElement(clay.ID("submit-text"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			FontID:     1, // Bold font
		},
	})
	engine.CloseElement()
	engine.CloseElement() // Close submit button

	engine.CloseElement() // Close form container
}
