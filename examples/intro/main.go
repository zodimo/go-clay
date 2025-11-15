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

type PanelConfig struct {
	Color        clay.Clay_Color
	CornerRadius clay.Clay_CornerRadius
}

type Document struct {
	Title   string
	Content string
}

var COLOR_WHITE = clay.CLAY_RGBA(255, 255, 255, 255)

func RenderDocumentTitles() []clay.ClayContainer {
	containers := []clay.ClayContainer{}
	for _, document := range Documents {
		button := clay.CLAY(
			"",
			clay.Clay_ElementDeclaration{
				Layout: clay.Clay_LayoutConfig{
					Padding: clay.Clay_Padding{
						Left:  16,
						Right: 16,
					},
				},
				CornerRadius: clay.CLAY_CORNER_RADIUS(5),
			},
			clay.CLAY_TEXT(document.Title,
				clay.TextWithFontSize(16),
				clay.TextWithColor(COLOR_WHITE),
			))
		containers = append(containers, button)
	}
	return containers
}

func RenderHeaderButton(text string) clay.ClayContainer {
	return clay.CLAY(
		"",
		clay.Clay_ElementDeclaration{
			Layout: clay.Clay_LayoutConfig{
				Padding: clay.Clay_Padding{
					Left:   16,
					Right:  16,
					Top:    8,
					Bottom: 8,
				},
			},
			BackgroundColor: clay.CLAY_RGBA(140, 140, 150, 255), //grey
			CornerRadius:    clay.CLAY_CORNER_RADIUS(5),
		},
		clay.CLAY_TEXT(text,
			clay.TextWithFontSize(16),
			clay.TextWithColor(clay.CLAY_RGBA(255, 255, 255, 255)),
			clay.TextWithTextAlignment(clay.CLAY_TEXT_ALIGN_CENTER),
			clay.TextWithFontId(1),
		),
	)
}

var SelectedDocumentIndex int = 0

func RenderDocumentContent(index int) []clay.ClayContainer {
	containers := []clay.ClayContainer{}
	document := Documents[index]
	containers = append(containers, clay.CLAY_TEXT(document.Title,
		clay.TextWithFontSize(24),
		clay.TextWithColor(COLOR_WHITE),
		clay.TextWithFontId(1),
	))
	containers = append(containers, clay.CLAY_TEXT(document.Content,
		clay.TextWithFontSize(24),
		clay.TextWithColor(COLOR_WHITE),
		clay.TextWithFontId(1),
	))

	return containers
}

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

	layoutExpand := clay.Clay_Sizing{
		Width:  clay.CLAY_SIZING_GROW(clay.Clay_SizingMinMax{Min: 0}),
		Height: clay.CLAY_SIZING_GROW(clay.Clay_SizingMinMax{Min: 0}),
	}
	panelConfig := PanelConfig{
		Color:        clay.CLAY_RGBA(90, 90, 90, 255), //grey
		CornerRadius: clay.CLAY_CORNER_RADIUS(8),
	}

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
					clay.NewErrorHandler(func(errorData clay.Clay_ErrorData) {
						fmt.Printf("ErrorMessage: %s\n", errorData.ErrorText.String())
						fmt.Printf("ErrorType: %s\n", errorData.ErrorType.String())
						fmt.Printf("UserData: %v\n", errorData.UserData)
					}, nil),
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

			// // 2. Build Clay layout
			clay.Clay_BeginLayout()
			clay.CLAY_ROOT(
				"OuterContainer",
				clay.Clay_ElementDeclaration{
					Layout: clay.Clay_LayoutConfig{
						Sizing:          layoutExpand,
						LayoutDirection: clay.CLAY_TOP_TO_BOTTOM,
						Padding:         clay.CLAY_PADDING_ALL(16),
						ChildGap:        16,
					},
					BackgroundColor: clay.CLAY_RGBA(153, 153, 153, 255), //dark grey for main container
				},
				/////////////////////// Children goes here
				clay.CLAY(
					"HeaderBar",
					clay.Clay_ElementDeclaration{
						Layout: clay.Clay_LayoutConfig{
							Sizing: clay.Clay_Sizing{
								Width:  clay.CLAY_SIZING_PERCENT(1),
								Height: clay.CLAY_SIZING_FIXED(60),
							},
							ChildGap:       16,
							ChildAlignment: clay.Clay_ChildAlignment{Y: clay.CLAY_ALIGN_Y_CENTER},
							Padding: clay.Clay_Padding{
								Left:  16,
								Right: 16,
							},
						},
						BackgroundColor: panelConfig.Color,
						CornerRadius:    panelConfig.CornerRadius,
					},
					RenderHeaderButton("File"),
					RenderHeaderButton("Edit"),
					clay.CLAY("spacer",
						clay.Clay_ElementDeclaration{
							Layout: clay.Clay_LayoutConfig{
								Sizing: clay.Clay_Sizing{
									Width: clay.CLAY_SIZING_GROW(clay.Clay_SizingMinMax{Min: 0}),
								},
							},
						},
					),
					RenderHeaderButton("Upload"),
					RenderHeaderButton("Media"),
					RenderHeaderButton("Support"),
				),
				clay.CLAY(
					"LowerContent",
					clay.Clay_ElementDeclaration{
						Layout: clay.Clay_LayoutConfig{
							Sizing:          layoutExpand,
							ChildGap:        16,
							LayoutDirection: clay.CLAY_LEFT_TO_RIGHT,
						},
					},
					clay.CLAY(
						"Sidebar",
						clay.Clay_ElementDeclaration{
							Layout: clay.Clay_LayoutConfig{
								Sizing: clay.Clay_Sizing{
									Width:  clay.CLAY_SIZING_FIXED(250),
									Height: clay.CLAY_SIZING_GROW(clay.Clay_SizingMinMax{Min: 0}),
								},
								LayoutDirection: clay.CLAY_TOP_TO_BOTTOM,
								ChildGap:        8,
								Padding:         clay.CLAY_PADDING_ALL(16),
							},
							BackgroundColor: panelConfig.Color,
							CornerRadius:    panelConfig.CornerRadius,
						},

						// Render document title here
						RenderDocumentTitles()...,
					),
					clay.CLAY(
						"MainContent",
						clay.Clay_ElementDeclaration{
							Layout: clay.Clay_LayoutConfig{
								LayoutDirection: clay.CLAY_TOP_TO_BOTTOM,
								Sizing:          layoutExpand,
								ChildGap:        16,
								Padding:         clay.CLAY_PADDING_ALL(16),
							},
							BackgroundColor: panelConfig.Color,
							CornerRadius:    panelConfig.CornerRadius,
							Clip: clay.Clay_ClipElementConfig{
								Vertical:    true,
								ChildOffset: clay.Clay_GetScrollOffset(),
							},
						},
						RenderDocumentContent(SelectedDocumentIndex)...,
					),
				),
			)

			commands := clay.Clay_EndLayout()

			fmt.Printf("Generated %d render commands\n", len(commands))
			for i, cmd := range commands {
				switch cmd.CommandType {
				case clay.CLAY_RENDER_COMMAND_TYPE_RECTANGLE:
				case clay.CLAY_RENDER_COMMAND_TYPE_TEXT:
				default:
					printCommand(i, cmd)
				}
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
