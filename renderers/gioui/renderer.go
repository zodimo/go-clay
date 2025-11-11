package gioui

import (
	"gioui.org/op"
	"github.com/zodimo/clay-go/clay"
)

type GioRenderer struct {
	ops *op.Ops
}

func NewRenderer(ops *op.Ops) *GioRenderer {
	return &GioRenderer{
		ops: ops,
	}
}

func (r *GioRenderer) Render(commands clay.Clay__Array[clay.Clay_RenderCommand]) {}
