package main

import (
	"log"
	"os"

	"github.com/zodimo/go-clay/kimi2/clay"
	"github.com/zodimo/go-clay/kimi2/gioui"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"
)

func main() {
	go func() {
		w := &app.Window{}
		w.Option(
			app.Title("Clay hierarchy + resize"),
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
	r := gioui.NewRenderer(&ops)

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			clay.Clay_Initialize(1<<20, clay.Dimensions{Width: float32(gtx.Constraints.Max.X), Height: float32(gtx.Constraints.Max.Y)}, gioui.NewMeasurer())

			// ---- layout starts here ----
			clay.Clay_BeginLayout()

			// ROOT: fill window
			clay.CLAY(clay.ElementDeclaration{
				ID: clay.CLAY_ID("Root"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.CLAY_SIZING_GROW(1),
						Height: clay.CLAY_SIZING_GROW(1),
					},
					Direction: clay.CLAY_LEFT_TO_RIGHT,
					ChildGap:  0,
				},
				BackgroundColor: clay.Color{R: 0.2, G: 0.2, B: 0.2, A: 1.0},
			})

			// LEFT PANEL: fixed 200 px wide, full height
			clay.CLAY(clay.ElementDeclaration{
				ID: clay.CLAY_ID("LeftPanel"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.CLAY_SIZING_FIXED(200),
						Height: clay.CLAY_SIZING_GROW(1),
					},
					Padding: clay.CLAY_PADDING_ALL(20),
				},
				Text: &clay.TextElementConfig{
					Color:    clay.Color{R: 0.8, G: 0.8, B: 0.8, A: 1.0},
					FontSize: 18,
				},
			})

			// RIGHT AREA: takes remaining width
			clay.CLAY(clay.ElementDeclaration{
				ID: clay.CLAY_ID("RightArea"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.CLAY_SIZING_GROW(1),
						Height: clay.CLAY_SIZING_GROW(1),
					},
					Direction: clay.CLAY_TOP_TO_BOTTOM,
					Padding:   clay.CLAY_PADDING_ALL(0),
				},
				BackgroundColor: clay.Color{R: 0.2, G: 0.2, B: 0.2, A: 1.0},
			})

			// HEADER inside right area
			clay.CLAY(clay.ElementDeclaration{
				ID: clay.CLAY_ID("Header"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.CLAY_SIZING_GROW(1),
						Height: clay.CLAY_SIZING_FIXED(60),
					},
					Padding: clay.CLAY_PADDING_ALL(15),
				},
			}).Text("Header 60 px tall", clay.TextElementConfig{
				Color:    clay.Color{R: 1.0, G: 0.8, B: 0.4, A: 1.0},
				FontSize: 22,
			}).End()

			// CONTENT BLOCK: fills remaining vertical space
			clay.CLAY(clay.ElementDeclaration{
				ID: clay.CLAY_ID("Content"),
				Layout: clay.LayoutConfig{
					Sizing: clay.Sizing{
						Width:  clay.CLAY_SIZING_GROW(1),
						Height: clay.CLAY_SIZING_GROW(1),
					},
					Padding: clay.Padding{Top: 30, Left: 40},
				},
			}).Text("This block grows/shrinks\nwhen you resize the window.", clay.TextElementConfig{
				Color:    clay.Color{R: 0.6, G: 1.0, B: 0.6, A: 1.0},
				FontSize: 20,
			}).End()

			commands := clay.Clay_EndLayout()
			r.Render(commands)
			e.Frame(gtx.Ops)
		}
	}
}
