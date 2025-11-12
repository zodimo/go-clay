package claygio

import (
	"gioui.org/op"
	"github.com/zodimo/clay-go/clay"
)

// map all render variations to op.ops

type Renderer struct {
}

func (r *Renderer) Render(renderCommand clay.Clay_RenderCommand) {
	switch renderCommand.CommandType {
	case clay.CLAY_RENDER_COMMAND_TYPE_RECTANGLE:
		RenderRectangle(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_BORDER:
		RenderBorder(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_TEXT:
		RenderText(renderCommand)
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

func RenderRectangle(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Rectangle
	boundingBox := renderCommand.BoundingBox
	_ = spec
	_ = boundingBox
	return ops
}

func RenderBorder(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Border
	_ = spec
	return ops
}

func RenderText(renderCommand clay.Clay_RenderCommand) op.Ops {
	var ops op.Ops
	spec := renderCommand.RenderData.Text
	_ = spec
	return ops
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
