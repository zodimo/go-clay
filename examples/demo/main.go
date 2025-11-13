package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/clay-go/clay"
	"github.com/zodimo/clay-go/pkg/mem"
	"github.com/zodimo/clay-go/renderers/gioui"
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
	memory := make([]byte, 100*1024*1024)
	arena, err := mem.NewArena(memory)
	if err != nil {
		return err
	}

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
				clay.Clay_Initialize(
					*arena,
					clay.Clay_Dimensions{
						Width:  float32(gtx.Constraints.Max.X),
						Height: float32(gtx.Constraints.Max.Y),
					},
					// gioui.NewMeasurer(),
					clay.NewErrorHandler(func(errorData clay.Clay_ErrorData) {
						fmt.Printf("ErrorMessage: %s\n", errorData.ErrorText.String())
						fmt.Printf("ErrorType: %s\n", errorData.ErrorType.String())
						fmt.Printf("UserData: %v\n", errorData.UserData)
					}, nil),
				)
				// clay.Clay_SetDebugModeEnabled(true)
				clayReady = true
			}

			// log.Printf("window size: %v", gtx.Constraints.Max)
			clay.Clay_SetLayoutDimensions(
				clay.Clay_Dimensions{
					Width:  float32(gtx.Constraints.Max.X),
					Height: float32(gtx.Constraints.Max.Y),
				},
			)

			// // 2. Build Clay layout
			clay.Clay_BeginLayout()
			clay.CLAY(
				clay.CLAY_ID("main"),
				clay.Clay_ElementDeclaration{
					Layout: clay.Clay_LayoutConfig{
						Sizing: clay.Clay_Sizing{
							Width:  clay.CLAY_SIZING_FIXED(400),
							Height: clay.CLAY_SIZING_FIXED(300),
						},
						Padding: clay.CLAY_PADDING_ALL(16),
					},
					BackgroundColor: clay.Clay_Color{R: 0.9, G: 0.2, B: 0.2, A: 1},
					CornerRadius:    clay.CLAY_CORNER_RADIUS(40),
				},
			)

			// clay.Clay__CloseElement()

			// how to let the error show an element is not closed?

			// main := clay.CLAY(clay.ElementDeclaration{
			// 	// ID: clay.CLAY_ID("main"),
			// 	Layout: clay.LayoutConfig{
			// 		Sizing: clay.Sizing{
			// 			Width:  clay.CLAY_SIZING_FIXED(400),
			// 			Height: clay.CLAY_SIZING_FIXED(300),
			// 		},
			// 		Padding: clay.CLAY_PADDING_ALL(16),
			// 	},
			// 	BackgroundColor: clay.Color{R: 0.9, G: 0.2, B: 0.2, A: 1},
			// 	CornerRadius:    clay.CLAY_CORNER_RADIUS(10),
			// })

			// inner := clay.CLAY(clay.ElementDeclaration{
			// 	ID: clay.CLAY_ID("hello-text"),
			// 	Layout: clay.LayoutConfig{
			// 		Sizing: clay.Sizing{
			// 			Width:  clay.CLAY_SIZING_FIXED(200),
			// 			Height: clay.CLAY_SIZING_FIXED(50),
			// 		},
			// 	},
			// 	BackgroundColor: clay.Color{R: 0.2, G: 0.2, B: 0.9, A: 1},
			// 	CornerRadius:    clay.CLAY_CORNER_RADIUS(5),
			// }).Text("H!", clay.TextElementConfig{
			// 	FontSize: 24,
			// 	// LineHeight: 48,
			// 	Color:  clay.Color{R: 1, G: 1, B: 1, A: 1},
			// 	FontID: 0,
			// })
			// inner.End()
			// main.End()

			commands := clay.Clay_EndLayout()

			fmt.Printf("Generated %d render commands\n", len(commands))
			for i, cmd := range commands {
				// fmt.Printf("Command %d: Type=%d, ID=%d, Bounds=%+v\n",
				// 	i, cmd.CommandType, cmd.Id, cmd.BoundingBox)
				// fmt.Printf("Command %d: %s\n", i, cmd.String())
				printCommand(i, cmd)
			}

			renderer.Render(commands)

			// renderer.SetViewport(clay.BoundingBox{
			// 	X: 0, Y: 0,
			// 	Width:  float32(gtx.Constraints.Max.X),
			// 	Height: float32(gtx.Constraints.Max.Y),
			// })
			// if err := renderer.Render(commands); err != nil {
			// 	log.Printf("render error: %v", err)
			// }

			e.Frame(gtx.Ops)
		}
	}
}

func printCommand(i int, cmd clay.Clay_RenderCommand) {
	fmt.Println("--------------------------------")
	fmt.Printf("CommandType: %s\n", cmd.CommandType.String())
	fmt.Printf("Id: %d\n", cmd.Id)
	fmt.Printf("ZIndex: %d\n", cmd.ZIndex)
	fmt.Printf("BoundingBox: %s\n", cmd.BoundingBox.String())
	fmt.Printf("RenderData: %s\n", cmd.RenderData.String(cmd.CommandType))
	fmt.Printf("UserData: %+v\n", cmd.UserData)

	// commandJSON, err := json.MarshalIndent(cmd, "", "  ")
	// if err != nil {
	// 	log.Printf("error marshalling command: %v", err)
	// }
	// fmt.Printf("Command %d:\n%s\n", i, string(commandJSON))
}
