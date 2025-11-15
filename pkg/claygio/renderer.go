package claygio

import (
	"gioui.org/op"
	"github.com/zodimo/clay-go/clay"
)

type RendererOptions struct {
	FontManager *FontManager
}

type RendererOption func(*RendererOptions)

func RendererWithFontManager(fontManager *FontManager) RendererOption {
	return func(o *RendererOptions) {
		o.FontManager = fontManager
	}
}

type Renderer interface {
	Render(ops *op.Ops, commands []clay.Clay_RenderCommand)
}
type renderer struct {
	fontManager       *FontManager
	clippingContainer *ClippingContainer
	clippingActive    bool
	opCallStack       []op.CallOp
}

func NewRenderer(opts ...RendererOption) Renderer {
	options := &RendererOptions{
		FontManager: NewFontManager(),
	}
	for _, opt := range opts {
		opt(options)
	}
	return &renderer{
		fontManager: options.FontManager,
		opCallStack: make([]op.CallOp, 0),
	}
}

func (r *renderer) Render(ops *op.Ops, commands []clay.Clay_RenderCommand) {
	for _, command := range commands {
		r.render(ops, command)
	}
}
