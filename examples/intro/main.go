package main

import (
	"fmt"
	"log"
	"os"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
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

func Call[T any](fn func() T) T {
	return fn()
}

func RenderDocumentTitles(gtx layout.Context) []clay.ClayContainer {
	containers := []clay.ClayContainer{}
	for i, document := range Documents {
		var button clay.ClayContainer
		if i == SelectedDocumentIndex {
			button = clay.CLAY(
				"",
				clay.Clay_ElementDeclaration{
					Layout: clay.Clay_LayoutConfig{
						Padding: clay.CLAY_PADDING_ALL(16),
						Sizing: clay.Clay_Sizing{
							Width: clay.CLAY_SIZING_GROW(clay.Clay_SizingMinMax{Min: 0}),
						},
					},
					CornerRadius:    clay.CLAY_CORNER_RADIUS(5),
					BackgroundColor: clay.CLAY_RGBA(120, 120, 120, 255), //grey
					// OnHover: clay.Clay_OnHoverConfig{
					// 	OnHoverFunction: handlerSidebarClick,
					// 	UserData:        SidebarClickData{DocumentIndex: i, Gtx: gtx},
					// },
				},
				clay.CLAY_TEXT(document.Title,
					clay.TextWithFontSize(16),
					clay.TextWithColor(COLOR_WHITE),
				))
		} else {

			button = clay.CLAY(
				"",
				clay.Clay_ElementDeclaration{
					Layout: clay.Clay_LayoutConfig{
						Padding: clay.CLAY_PADDING_ALL(16),
					},
					CornerRadius: clay.CLAY_CORNER_RADIUS(5),
				},
				//element is only open from here... so calling onHover in the declaration will be for the parent element...

				clay.CLAY_ON_HOVER(handlerSidebarClick, SidebarClickData{DocumentIndex: i, Gtx: gtx}),
				clay.CLAY_TEXT(document.Title,
					clay.TextWithFontSize(16),
					clay.TextWithColor(COLOR_WHITE),
				))
		}

		containers = append(containers, button)
	}
	return containers
}

func RenderFileMenu(show bool) clay.ClayContainer {
	if !show {
		return nil
	}

	// we need the container, using hover and offset creates a pointer over continuity problem
	return clay.CLAY("FileMenuContainer", clay.Clay_ElementDeclaration{
		Floating: clay.Clay_FloatingElementConfig{
			AttachTo: clay.CLAY_ATTACH_TO_PARENT,
			AttachPoints: clay.Clay_FloatingAttachPoints{
				// Element: clay.CLAY_ATTACH_POINT_LEFT_TOP,
				Parent: clay.CLAY_ATTACH_POINT_LEFT_BOTTOM,
			},
		},
		Layout: clay.Clay_LayoutConfig{
			LayoutDirection: clay.CLAY_TOP_TO_BOTTOM,
			Padding: clay.Clay_Padding{
				Top: 8,
			},
		},
	},
		clay.CLAY(
			"FileMenu",
			clay.Clay_ElementDeclaration{
				BackgroundColor: clay.CLAY_RGBA(40, 40, 40, 255), //grey
				CornerRadius:    clay.CLAY_CORNER_RADIUS(8),
				Layout: clay.Clay_LayoutConfig{
					LayoutDirection: clay.CLAY_TOP_TO_BOTTOM,
					ChildGap:        8,
					Padding:         clay.CLAY_PADDING_ALL(16),
					Sizing: clay.Clay_Sizing{
						Width: clay.CLAY_SIZING_FIXED(200),
					},
				},
			},
			RenderDropdownMenuButton("New"),
			RenderDropdownMenuButton("Open"),
			RenderDropdownMenuButton("Close"),
		),
	)
}

func RenderFileButton() clay.ClayContainer {
	fileMenuVisible := clay.Clay_PointerOver(clay.Clay_GetElementId("FileButton")) || clay.Clay_PointerOver(clay.Clay_GetElementId("FileMenuContainer")) || clay.Clay_PointerOver(clay.Clay_GetElementId("FileMenu"))
	return clay.CLAY(
		"FileButton",
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
		clay.CLAY_TEXT("File",
			clay.TextWithFontSize(16),
			clay.TextWithColor(clay.CLAY_RGBA(255, 255, 255, 255)),
			clay.TextWithTextAlignment(clay.CLAY_TEXT_ALIGN_CENTER),
			clay.TextWithFontId(1),
		),
		RenderFileMenu(fileMenuVisible),
	)
}

func RenderDropdownMenuButton(text string) clay.ClayContainer {
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
		},
		clay.CLAY_TEXT(text,
			clay.TextWithFontSize(16),
			clay.TextWithColor(clay.CLAY_RGBA(255, 255, 255, 255)),
			clay.TextWithTextAlignment(clay.CLAY_TEXT_ALIGN_CENTER),
			clay.TextWithFontId(1),
		),
	)
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

type SidebarClickData struct {
	DocumentIndex int
	Gtx           layout.Context
}

func handlerSidebarClick(elementId clay.Clay_ElementId, pointerInfo clay.Clay_PointerData, userData any) {
	clickData, ok := userData.(SidebarClickData)
	if !ok {
		panic("userData is not a SidebarClickData")
	}
	clicked := pointerInfo.State == clay.CLAY_POINTER_DATA_PRESSED_THIS_FRAME
	if clicked {
		fmt.Printf("Sidebar clicked: %v for element %d\n", clicked, clickData.DocumentIndex)
		SelectedDocumentIndex = clickData.DocumentIndex
	}
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
					clay.Clay_WithErrorHandler(clay.NewErrorHandler(func(errorData clay.Clay_ErrorData) {
						fmt.Printf("ErrorMessage: %s\n", errorData.ErrorText.String())
						fmt.Printf("ErrorType: %s\n", errorData.ErrorType.String())
						fmt.Printf("UserData: %v\n", errorData.UserData)
					}, nil)),
					clay.Clay_WithMaxMeasureTextCacheWordCount(32000),
				)
				// clay.Clay_SetDebugModeEnabled(true)
				clay.Clay_SetMeasureTextFunction(clayGioEngine.MeasureText, gtx)
				clay.Clay_SetMaxMeasureTextCacheWordCount(32000)
				clayReady = true
			}

			clayGioEngine.UpdateInput(gtx)
			clay.Clay_SetPointerState(clayGioEngine.GetMousePosition(), clayGioEngine.IsPointerDown())

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
					RenderFileButton(),
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
						RenderDocumentTitles(gtx)...,
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
							// Clip: clay.Clay_ClipElementConfig{
							// 	Vertical:    true,
							// 	ChildOffset: clay.Clay_GetScrollOffset(),
							// },
						},
						RenderDocumentContent(SelectedDocumentIndex)...,
					),
				),
			)

			commands := clay.Clay_EndLayout()

			// fmt.Printf("Generated %d render commands\n", len(commands))
			for i, cmd := range commands {
				switch cmd.CommandType {
				case clay.CLAY_RENDER_COMMAND_TYPE_RECTANGLE:
				case clay.CLAY_RENDER_COMMAND_TYPE_TEXT:
				default:
					printCommand(i, cmd)
				}
			}

			clayGioEngine.Render(gtx.Ops, commands)
			gtx.Execute(op.InvalidateCmd{At: gtx.Now})
			e.Frame(gtx.Ops)
		}
	}
}

func printCommand(i int, cmd clay.Clay_RenderCommand) {
	fmt.Printf("Command %d:\n%s\n", i, cmd.DebugString())

	// commandJSON, err := json.MarshalIndent(cmd, "", "  ")
	// if err != nil {
	// 	log.Printf("error marshalling command: %v", err)
	// }
	// fmt.Printf("Command %d:\n%s\n", i, string(commandJSON))
}
