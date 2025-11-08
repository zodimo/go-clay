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
			app.Title("Clay Sidebar Layout Example"),
			app.Size(unit.Dp(1000), unit.Dp(700)),
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

			// Create sidebar layout
			createSidebarLayout(engine)

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

func createSidebarLayout(engine clay.LayoutEngine) {
	// Main container
	engine.OpenElement(clay.ID("main"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingGrow(0),
			},
			Direction: clay.LeftToRight,
			ChildGap:  16,
		},
	})

	// Sidebar
	engine.OpenElement(clay.ID("sidebar"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingFixed(300),
				Height: clay.SizingGrow(0),
			},
			Direction: clay.TopToBottom,
			Padding:   clay.PaddingAll(16),
		},
		BackgroundColor: clay.Color{R: 0.8, G: 0.8, B: 0.9, A: 1.0},
		CornerRadius: clay.CornerRadius{
			TopRight:    8.0,
			BottomRight: 8.0,
		},
	})

	// Sidebar title
	engine.OpenElement(clay.ID("sidebar-title"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   20,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0, G: 0, B: 0, A: 1.0},
			FontID:     1, // Bold font
		},
	})
	engine.CloseElement() // Close sidebar title

	// Navigation container
	engine.OpenElement(clay.ID("nav"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingGrow(0),
			},
			Direction: clay.TopToBottom,
			ChildGap:  8,
		},
	})

	// Navigation items
	engine.OpenElement(clay.ID("nav-home"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0.2, G: 0.2, B: 0.2, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()

	engine.OpenElement(clay.ID("nav-about"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0.2, G: 0.2, B: 0.2, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()

	engine.OpenElement(clay.ID("nav-contact"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0.2, G: 0.2, B: 0.2, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()

	engine.CloseElement() // Close nav container
	engine.CloseElement() // Close sidebar

	// Main content
	engine.OpenElement(clay.ID("content"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingGrow(0),
			},
			Padding: clay.PaddingAll(16),
		},
		BackgroundColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
	})

	// Main content title
	engine.OpenElement(clay.ID("content-title"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   24,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0, G: 0, B: 0, A: 1.0},
			FontID:     1, // Bold font
		},
	})
	engine.CloseElement()

	// Main content description
	engine.OpenElement(clay.ID("content-desc"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 0.3, G: 0.3, B: 0.3, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()

	engine.CloseElement() // Close content
	engine.CloseElement() // Close main container
}
