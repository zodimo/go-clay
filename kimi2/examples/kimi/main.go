package main

import (
	"github.com/zodimo/go-clay/kimi2/clay"
)

type measurer struct{}

func (measurer) MeasureText(s string, cfg clay.TextElementConfig) clay.Dimensions {
	// dummy: 0.6 em per character
	w := float32(len(s)) * cfg.FontSize * 0.6
	h := cfg.FontSize * cfg.LineHeight
	return clay.Dimensions{Width: w, Height: h}
}

func main() {
	clay.Clay_Initialize(1<<20, clay.Dimensions{Width: 800, Height: 600}, measurer{})

	clay.Clay_BeginLayout()

	clay.CLAY(clay.ElementDeclaration{
		Layout: clay.LayoutConfig{
			Sizing: clay.Sizing{
				Width:  clay.CLAY_SIZING_GROW(1),
				Height: clay.CLAY_SIZING_GROW(1),
			},
			Padding: clay.CLAY_PADDING_ALL(20),
		},
		BackgroundColor: clay.Color{240, 240, 240, 255},
	}).Text("Hello from full Go Clay!", clay.TextElementConfig{
		FontSize:   24,
		LineHeight: 1.2,
		Color:      clay.Color{0, 0, 0, 255},
	}).End()

	cmds := clay.Clay_EndLayout()

	// draw cmds...
	for _, c := range cmds {
		switch d := c.Data.(type) {
		case clay.RectangleRenderData:
			println("RECT", c.BoundingBox.X, c.BoundingBox.Y, c.BoundingBox.Width, c.BoundingBox.Height)
		case clay.TextRenderData:
			println("TEXT", d.StringContents, c.BoundingBox.X, c.BoundingBox.Y)
		}
	}
}
