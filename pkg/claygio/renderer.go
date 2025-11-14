package claygio

import (
	"gioui.org/op"
	"github.com/zodimo/clay-go/clay"
)

type Renderer struct {
	ops *op.Ops
}

func NewRenderer(ops *op.Ops) *Renderer {
	return &Renderer{
		ops: ops,
	}
}

func (r *Renderer) Render(commands []clay.Clay_RenderCommand) {
	for _, command := range commands {
		Render(r.ops, command)
	}
}
