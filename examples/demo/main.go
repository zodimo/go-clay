package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/unit"

	"github.com/zodimo/clay-go/clay"
	"github.com/zodimo/clay-go/pkg/claygio"
	"github.com/zodimo/clay-go/pkg/mem"
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
	clayReady bool
)

func run(w *app.Window) error {
	var ops op.Ops
	memory := make([]byte, 25*1024*1024)
	arena, err := mem.NewArena(memory)
	if err != nil {
		return err
	}

	fontCollection := gofont.Collection()
	fontManager := claygio.NewFontManager(claygio.FontManagerWithFontCollection(fontCollection))

	clayGioEngine := claygio.NewClayGioEngine(claygio.ClaygioWithFontManager(fontManager))

	for {
		switch e := w.Event().(type) {
		case app.DestroyEvent:
			return e.Err
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			if !clayReady {
				clay.Clay_Initialize(
					*arena,
					clay.Clay_Dimensions{
						Width:  float32(gtx.Constraints.Max.X),
						Height: float32(gtx.Constraints.Max.Y),
					},
					// gioui.NewMeasurer(),
					clay.Clay_WithErrorHandler(clay.NewErrorHandler(func(errorData clay.Clay_ErrorData) {
						fmt.Printf("ErrorMessage: %s\n", errorData.ErrorText.String())
						fmt.Printf("ErrorType: %s\n", errorData.ErrorType.String())
						fmt.Printf("UserData: %v\n", errorData.UserData)
					}, nil)),
				)
				// clay.Clay_SetDebugModeEnabled(true)
				clay.Clay_SetMeasureTextFunction(clayGioEngine.MeasureText, gtx)
				clayReady = true
			}

			// log.Printf("window size: %v", gtx.Constraints.Max)
			clay.Clay_SetLayoutDimensions(
				clay.Clay_Dimensions{
					Width:  float32(gtx.Constraints.Max.X),
					Height: float32(gtx.Constraints.Max.Y),
				},
			)

			clayGioEngine.UpdateInput(gtx)
			clay.Clay_SetPointerState(clayGioEngine.GetMousePosition(), true)

			// // 2. Build Clay layout
			clay.Clay_BeginLayout()
			clay.CLAY_ROOT(
				"main",
				clay.Clay_ElementDeclaration{
					Layout: clay.Clay_LayoutConfig{
						Sizing: clay.Clay_Sizing{
							Width:  clay.CLAY_SIZING_PERCENT(1),
							Height: clay.CLAY_SIZING_PERCENT(1),
						},
						Padding:         clay.CLAY_PADDING_ALL(40),
						LayoutDirection: clay.CLAY_TOP_TO_BOTTOM,
						ChildGap:        10,
						ChildAlignment:  clay.Clay_ChildAlignment{X: clay.CLAY_ALIGN_X_CENTER, Y: clay.CLAY_ALIGN_Y_CENTER},
					},
					Border: clay.Clay_BorderElementConfig{
						Color: clay.CLAY_RGBA(0, 0, 0, 255), // black
						Width: clay.Clay_BorderWidth{
							Left:   5,
							Right:  5,
							Top:    5,
							Bottom: 5,
						},
					},
					BackgroundColor: clay.CLAY_RGBA(153, 153, 153, 255), //dark grey for main container
					CornerRadius:    clay.CLAY_CORNER_RADIUS(40),
				},
				clay.CLAY_TEXT("Hello, world!", clay.TextWithFontSize(48)),
				clay.CLAY(
					"inner",
					clay.Clay_ElementDeclaration{
						Layout: clay.Clay_LayoutConfig{
							Sizing: clay.Clay_Sizing{
								Width:  clay.CLAY_SIZING_FIXED(200),
								Height: clay.CLAY_SIZING_FIXED(50),
							},
							Padding: clay.CLAY_PADDING_ALL(16),
						},
						BackgroundColor: clay.CLAY_RGBA(255, 0, 0, 255), // red
						CornerRadius:    clay.CLAY_CORNER_RADIUS(15),
					},
				),
				clay.CLAY(
					"inner2",
					clay.Clay_ElementDeclaration{
						Layout: clay.Clay_LayoutConfig{
							Sizing: clay.Clay_Sizing{
								Width:  clay.CLAY_SIZING_FIXED(200),
								Height: clay.CLAY_SIZING_FIXED(100),
							},
							Padding:         clay.CLAY_PADDING_ALL(16),
							LayoutDirection: clay.CLAY_TOP_TO_BOTTOM,
						},
						BackgroundColor: clay.CLAY_RGBA(0, 255, 0, 255), // green
						CornerRadius:    clay.CLAY_CORNER_RADIUS(15),
					},
					clay.CLAY_TEXT("good, morning!"),
					clay.CLAY(
						"inner3",
						clay.Clay_ElementDeclaration{
							Layout: clay.Clay_LayoutConfig{
								Sizing: clay.Clay_Sizing{
									Width:  clay.CLAY_SIZING_PERCENT(1),
									Height: clay.CLAY_SIZING_PERCENT(1),
								},
								Padding: clay.CLAY_PADDING_ALL(16),
							},
							BackgroundColor: clay.CLAY_RGBA(0, 0, 255, 255), // blue
							CornerRadius:    clay.CLAY_CORNER_RADIUS(15),
							Border: clay.Clay_BorderElementConfig{
								Color: clay.CLAY_RGBA(255, 0, 255, 255), // purple
								Width: clay.Clay_BorderWidth{
									Left:   2,
									Right:  2,
									Top:    2,
									Bottom: 2,
								},
							},
						},
						clay.CLAY_TEXT("good, morning!", clay.TextWithFontSize(8)),
					),
				),
			)

			commands := clay.Clay_EndLayout()

			fmt.Printf("Generated %d render commands\n", len(commands))
			for i, cmd := range commands {
				printCommand(i, cmd)
			}

			clayGioEngine.Render(gtx.Ops, commands)
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
