package claygio

import (
	"fmt"
	"image"

	"gioui.org/op"
	"gioui.org/op/clip"
	"github.com/zodimo/clay-go/clay"
)

// map all render variations to op.ops

func (r *renderer) render(ops *op.Ops, renderCommand clay.Clay_RenderCommand) {

	var callOp op.CallOp
	switch renderCommand.CommandType {
	case clay.CLAY_RENDER_COMMAND_TYPE_RECTANGLE:
		callOp = RenderRectangle(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_BORDER:
		callOp = RenderBorder(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_TEXT:
		callOp = r.RenderText(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_IMAGE:
		callOp = RenderImage(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_SCISSOR_START:
		r.RenderScissorStart(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_SCISSOR_END:
		r.RenderScissorEnd(renderCommand)
	case clay.CLAY_RENDER_COMMAND_TYPE_CUSTOM:
		callOp = RenderCustom(renderCommand)
	default:
		panic(fmt.Sprintf("Unknown render command type: %s", renderCommand.CommandType))
	}
	if r.clippingActive {
		fmt.Printf("clipping active, pushing to op call stack: %s\n", renderCommand.DebugString())
		r.opCallStack = append(r.opCallStack, callOp)
	} else {
		if len(r.opCallStack) > 0 {
			fmt.Printf("flushing op call stack, in clipping area :%s\n", r.clippingContainer.String())
			//setup container constraints
			// offsetOp := op.Offset(image.Pt(int(r.clippingContainer.X), int(r.clippingContainer.Y))).Push(ops)
			clipOp := clip.Rect{Max: image.Pt(int(r.clippingContainer.Width), int(r.clippingContainer.Height))}.Push(ops)
			for _, callOp := range r.opCallStack {
				callOp.Add(ops)
			}
			clipOp.Pop()
			// offsetOp.Pop()
			r.opCallStack = make([]op.CallOp, 0)
		} else {
			callOp.Add(ops)
		}
	}
}

func RenderRectangle(renderCommand clay.Clay_RenderCommand) op.CallOp {
	return RenderRectangleWithBounds(renderCommand)
}

func RenderBorder(renderCommand clay.Clay_RenderCommand) op.CallOp {
	return RenderBorderWithBounds(renderCommand)
}

func (r *renderer) RenderText(renderCommand clay.Clay_RenderCommand) op.CallOp {
	return r.RenderTextWithBounds(renderCommand)
}

func RenderImage(renderCommand clay.Clay_RenderCommand) op.CallOp {
	var ops op.Ops
	opRecord := op.Record(&ops)
	spec := renderCommand.RenderData.Image
	_ = spec
	return opRecord.Stop()
}

func (r *renderer) RenderScissorStart(renderCommand clay.Clay_RenderCommand) {
	if r.clippingActive {
		panic("Scissor start called while already clipping")
	}
	r.clippingContainer = &ClippingContainer{
		X:          renderCommand.BoundingBox.X,
		Y:          renderCommand.BoundingBox.Y,
		Width:      renderCommand.BoundingBox.Width,
		Height:     renderCommand.BoundingBox.Height,
		Horizontal: renderCommand.RenderData.Clip.Horizontal,
		Vertical:   renderCommand.RenderData.Clip.Vertical,
	}
	r.clippingActive = true

}

func (r *renderer) RenderScissorEnd(renderCommand clay.Clay_RenderCommand) {
	if !r.clippingActive {
		panic("Scissor end called while not clipping")
	}
	r.clippingActive = false
}

func RenderCustom(renderCommand clay.Clay_RenderCommand) op.CallOp {
	var ops op.Ops
	opRecord := op.Record(&ops)
	spec := renderCommand.RenderData.Custom
	_ = spec
	return opRecord.Stop()
}
