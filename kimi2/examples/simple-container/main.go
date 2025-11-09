package main

import (
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/go-clay/kimi2/clay"
	"github.com/zodimo/go-clay/kimi2/gioui"
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
			clay.Clay_Initialize(1<<20, clay.Dimensions{Width: float32(gtx.Constraints.Max.X), Height: float32(gtx.Constraints.Max.Y)}, gioui.NewMeasurer())

			// Begin layout
			clay.Clay_BeginLayout()

			// Create the main container using engine methods directly
			main := clay.CLAY(clay.ElementDeclaration{
				ID: clay.CLAY_ID("main"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.CLAY_SIZING_FIXED(400),
						Height: clay.CLAY_SIZING_FIXED(300),
					},
					Padding: clay.CLAY_PADDING_ALL(16),
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
			clay.CLAY(clay.ElementDeclaration{
				ID: clay.CLAY_ID("hello-text"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.CLAY_SIZING_FIXED(200),
						Height: clay.CLAY_SIZING_FIXED(50),
					},
				},
				BackgroundColor: clay.Color{R: 0.2, G: 0.2, B: 0.9, A: 1.0}, // Blue background for text
				CornerRadius: clay.CornerRadius{
					TopLeft:     5.0,
					TopRight:    5.0,
					BottomLeft:  5.0,
					BottomRight: 5.0,
				},
			}).Text("Hello, world!", clay.TextElementConfig{
				FontSize:   24,
				LineHeight: 1.2,
				Color:      clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
				FontID:     0,
			}).End()

			main.End()

			// End layout and get render commands
			commands := clay.Clay_EndLayout()

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
