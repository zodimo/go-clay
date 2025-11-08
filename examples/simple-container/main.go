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
			app.Title("Clay Simple Container Example"),
			app.Size(unit.Dp(800), unit.Dp(600)),
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

			// Create the main container using engine methods directly
			engine.OpenElement(clay.ID("main"), clay.ElementConfig{
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.SizingFixed(400),
						Height: clay.SizingFixed(300),
					},
					Padding: clay.PaddingAll(16),
				},
				BackgroundColor: clay.Color{R: 0.9, G: 0.2, B: 0.2, A: 1.0}, // Bright red for visibility
				CornerRadius: clay.CornerRadius{
					TopLeft:     10.0,
					TopRight:    10.0,
					BottomLeft:  10.0,
					BottomRight: 10.0,
				},
			})

			// Add text element
			engine.OpenElement(clay.ID("hello-text"), clay.ElementConfig{
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.SizingFixed(200),
						Height: clay.SizingFixed(50),
					},
				},
				Text: &clay.TextConfig{
					FontSize:   24,
					LineHeight: 1.2,                                        // Set line height to 1.2 times font size
					Color:      clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0}, // White text
					FontID:     0,                                          // Default font
				},
				BackgroundColor: clay.Color{R: 0.2, G: 0.2, B: 0.9, A: 1.0}, // Blue background for text
				CornerRadius: clay.CornerRadius{
					TopLeft:     5.0,
					TopRight:    5.0,
					BottomLeft:  5.0,
					BottomRight: 5.0,
				},
			})
			engine.CloseElement() // Close text element
			engine.CloseElement() // Close main container

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
