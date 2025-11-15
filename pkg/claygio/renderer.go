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

type Renderer struct {
	fontManager *FontManager
}

func NewRenderer(opts ...RendererOption) *Renderer {
	options := &RendererOptions{
		FontManager: NewFontManager(),
	}
	for _, opt := range opts {
		opt(options)
	}
	return &Renderer{
		fontManager: options.FontManager,
	}
}

func (r *Renderer) Render(ops *op.Ops, commands []clay.Clay_RenderCommand) {
	for _, command := range commands {
		r.render(ops, command)
	}
}
