package claygio

import (
	"gioui.org/op"
	"github.com/zodimo/clay-go/clay"
)

// map all render variations to op.ops

func Render(ops *op.Ops, renderCommand clay.Clay_RenderCommand) {
	switch renderCommand.CommandType {
	case clay.CLAY_RENDER_COMMAND_TYPE_RECTANGLE:
		RenderRectangle(ops, renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_BORDER:
		RenderBorder(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_TEXT:
		RenderText(ops, renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_IMAGE:
		RenderImage(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_SCISSOR_START:
		RenderScissorStart(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_SCISSOR_END:
		RenderScissorEnd(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_CUSTOM:
		RenderCustom(renderCommand)
	}
}

func RenderRectangle(ops *op.Ops, renderCommand clay.Clay_RenderCommand) {
	RenderRectangleWithBounds(ops, renderCommand)
}

func RenderBorder(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Border
	_ = spec
	return ops
}

func RenderText(ops *op.Ops, renderCommand clay.Clay_RenderCommand) {
	RenderTextWithBounds(ops, renderCommand)
}

func RenderImage(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Image
	_ = spec
	return ops
}

func RenderScissorStart(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Clip
	_ = spec
	return ops
}

func RenderScissorEnd(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Custom
	_ = spec
	return ops
}

func RenderCustom(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Custom
	_ = spec
	return ops

}
