package gioui

import (
	"gioui.org/op"
	"github.com/zodimo/clay-go/clay"
	"github.com/zodimo/clay-go/pkg/claygio"
)

type GioRenderer struct {
	ops *op.Ops
}

func NewRenderer(ops *op.Ops) *GioRenderer {
	return &GioRenderer{
		ops: ops,
	}
}

func (r *GioRenderer) Render(commands []clay.Clay_RenderCommand) {
	for _, command := range commands {
		claygio.Render(r.ops, command)
	}
}
