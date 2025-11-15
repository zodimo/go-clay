package claygio

import (
	"fmt"
	"image"

	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/zodimo/clay-go/clay"
)

// renderTextWithBounds renders text using bounds from RenderCommand
func RenderTextWithBounds(ops *op.Ops, renderCommand clay.Clay_RenderCommand) error {
	bounds := renderCommand.BoundingBox
	cmd := renderCommand.RenderData.Text
	// Validate bounds
	if bounds.Width <= 0 || bounds.Height <= 0 {
		panic(fmt.Sprintf("Invalid bounds: width=%f, height=%f", bounds.Width, bounds.Height))
	}

	// Validate text content
	if cmd.StringContents.Length == 0 {
		// Empty text is valid, just return without rendering
		panic("Empty text")
	}

	fontManager := NewFontManager()
	// Get font and shaper from font manager
	fontObj := fontManager.GetFont(cmd.FontId)
	shaper := fontManager.GetShaper()

	// Create color operation
	colorMacro := op.Record(ops)
	paint.ColorOp{Color: ClayToGioColor(cmd.TextColor)}.Add(ops)
	colorCallOp := colorMacro.Stop()

	// Create label with Clay parameters
	label := widget.Label{
		Alignment:  text.Start, // positions
		MaxLines:   1,          // Unlimited - can be set from Clay command
		LineHeight: unit.Sp(cmd.LineHeight),
		WrapPolicy: text.WrapGraphemes,
	}

	// Position within bounds
	stack := op.Offset(image.Pt(int(bounds.X), int(bounds.Y))).Push(ops)
	defer stack.Pop()

	// Create layout context with bounds constraints
	gtx := layout.Context{
		Ops: ops,
		Constraints: layout.Constraints{
			Max: image.Pt(int(bounds.Width), int(bounds.Height)),
		},
		Metric: unit.Metric{PxPerDp: 1, PxPerSp: 1},
	}

	// Render text using gio pattern
	label.Layout(gtx, shaper, fontObj, unit.Sp(cmd.FontSize), cmd.StringContents.String(), colorCallOp)

	return nil
}
