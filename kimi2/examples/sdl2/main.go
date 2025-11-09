package main

import (
	"time"

	"github.com/zodimo/go-clay/kimi2/clay"
	"github.com/zodimo/go-clay/kimi2/renderer"
)

type measurer struct{}

func (measurer) MeasureText(text string, cfg clay.TextElementConfig) clay.Dimensions {
	return clay.Dimensions{Width: float32(len(text)) * cfg.FontSize * 0.6, Height: cfg.FontSize * cfg.LineHeight}
}

func main() {
	r, err := renderer.NewSDL2Renderer("Clay Go Example", 800, 600)
	if err != nil {
		panic(err)
	}
	defer r.Close()

	measurer := measurer{}

	clay.Clay_Initialize(1<<20, clay.Dimensions{Width: 800, Height: 600}, measurer)

	running := true
	for running {
		clay.Clay_BeginLayout()

		clay.CLAY(clay.ElementDeclaration{
			ID: clay.CLAY_ID("Background"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.CLAY_SIZING_FIXED(800),
					Height: clay.CLAY_SIZING_FIXED(600),
				},
			},
			BackgroundColor: clay.Color{R: 0.2, G: 0.2, B: 0.2, A: 1.0},
		})

		clay.CLAY(clay.ElementDeclaration{
			ID: clay.CLAY_ID("HelloText"),
			Layout: clay.LayoutConfig{
				Sizing: clay.Sizing{
					Width:  clay.CLAY_SIZING_FIXED(200),
					Height: clay.CLAY_SIZING_FIXED(40),
				},
			},
		}).Text("Hello from Clay in Go!", clay.TextElementConfig{
			FontSize: 20,
			Color:    clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
		}).End()

		commands := clay.Clay_EndLayout()
		r.Render(commands)

		time.Sleep(16 * time.Millisecond)
	}
}
