package main

import (
	"fmt"
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

var (
	renderer  *gioui.GioRenderer
	clayReady bool
)

func run(w *app.Window) error {
	var ops op.Ops

	clay.Clay_SetDebugModeEnabled(true)

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if renderer == nil {
				renderer = gioui.NewRenderer(gtx.Ops)
			}
			if !clayReady {
				clay.Clay_Initialize(1<<20,
					clay.Dimensions{
						Width:  float32(gtx.Constraints.Max.X),
						Height: float32(gtx.Constraints.Max.Y),
					},
					gioui.NewMeasurer(),
				)
				clayReady = true
			}

			log.Printf("window size: %v", gtx.Constraints.Max)
			clay.Clay_SetLayoutDimensions(
				clay.Dimensions{
					Width:  float32(gtx.Constraints.Max.X),
					Height: float32(gtx.Constraints.Max.Y),
				},
			)

			// 2. Build Clay layout
			clay.Clay_BeginLayout()
			main := clay.CLAY(clay.ElementDeclaration{
				// ID: clay.CLAY_ID("main"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.CLAY_SIZING_FIXED(400),
						Height: clay.CLAY_SIZING_FIXED(300),
					},
					Padding: clay.CLAY_PADDING_ALL(16),
				},
				BackgroundColor: clay.Color{R: 0.9, G: 0.2, B: 0.2, A: 1},
				CornerRadius:    clay.CLAY_CORNER_RADIUS(10),
			})

			inner := clay.CLAY(clay.ElementDeclaration{
				ID: clay.CLAY_ID("hello-text"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.CLAY_SIZING_FIXED(200),
						Height: clay.CLAY_SIZING_FIXED(50),
					},
				},
				BackgroundColor: clay.Color{R: 0.2, G: 0.2, B: 0.9, A: 1},
				CornerRadius:    clay.CLAY_CORNER_RADIUS(5),
			}).Text("H!", clay.TextElementConfig{
				FontSize: 24,
				// LineHeight: 48,
				Color:  clay.Color{R: 1, G: 1, B: 1, A: 1},
				FontID: 0,
			})
			inner.End()
			main.End()

			commands := clay.Clay_EndLayout()

			fmt.Printf("Generated %d render commands\n", len(commands))
			for i, cmd := range commands {
				fmt.Printf("Command %d: Type=%d, ID=%d, Bounds=%+v\n",
					i, cmd.CommandType, cmd.ID, cmd.BoundingBox)
			}

			renderer.SetViewport(clay.BoundingBox{
				X: 0, Y: 0,
				Width:  float32(gtx.Constraints.Max.X),
				Height: float32(gtx.Constraints.Max.Y),
			})
			if err := renderer.Render(commands); err != nil {
				log.Printf("render error: %v", err)
			}

			e.Frame(gtx.Ops)
		}
	}
}
