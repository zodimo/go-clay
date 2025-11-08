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
			app.Title("Clay Responsive Grid Example"),
			app.Size(unit.Dp(900), unit.Dp(600)),
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

			// Create responsive grid
			createResponsiveGrid(engine)

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

func createResponsiveGrid(engine clay.LayoutEngine) {
	// Grid container
	engine.OpenElement(clay.ID("grid"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingGrow(0),
				Height: clay.SizingGrow(0),
			},
			Direction: clay.LeftToRight,
			ChildGap:  8,
			Padding:   clay.PaddingAll(20),
		},
	})

	// Grid item 1
	engine.OpenElement(clay.ID("item1"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingPercent(0.33),
				Height: clay.SizingFixed(100),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.AlignXCenter,
				Y: clay.AlignYCenter,
			},
		},
		BackgroundColor: clay.Color{R: 0.9, G: 0.5, B: 0.5, A: 1.0},
		CornerRadius: clay.CornerRadius{
			TopLeft:     8.0,
			TopRight:    8.0,
			BottomLeft:  8.0,
			BottomRight: 8.0,
		},
	})

	// Item 1 text
	engine.OpenElement(clay.ID("item1-text"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()
	engine.CloseElement() // Close item1

	// Grid item 2
	engine.OpenElement(clay.ID("item2"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingPercent(0.33),
				Height: clay.SizingFixed(100),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.AlignXCenter,
				Y: clay.AlignYCenter,
			},
		},
		BackgroundColor: clay.Color{R: 0.5, G: 0.9, B: 0.5, A: 1.0},
		CornerRadius: clay.CornerRadius{
			TopLeft:     8.0,
			TopRight:    8.0,
			BottomLeft:  8.0,
			BottomRight: 8.0,
		},
	})

	// Item 2 text
	engine.OpenElement(clay.ID("item2-text"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()
	engine.CloseElement() // Close item2

	// Grid item 3
	engine.OpenElement(clay.ID("item3"), clay.ElementConfig{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.SizingPercent(0.34),
				Height: clay.SizingFixed(100),
			},
			ChildAlignment: clay.ChildAlignment{
				X: clay.AlignXCenter,
				Y: clay.AlignYCenter,
			},
		},
		BackgroundColor: clay.Color{R: 0.5, G: 0.5, B: 0.9, A: 1.0},
		CornerRadius: clay.CornerRadius{
			TopLeft:     8.0,
			TopRight:    8.0,
			BottomLeft:  8.0,
			BottomRight: 8.0,
		},
	})

	// Item 3 text
	engine.OpenElement(clay.ID("item3-text"), clay.ElementConfig{
		Text: &clay.TextConfig{
			FontSize:   16,
			LineHeight: 1.2,
			Color:      clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
			FontID:     0, // Regular font
		},
	})
	engine.CloseElement()
	engine.CloseElement() // Close item3

	engine.CloseElement() // Close grid container
}
