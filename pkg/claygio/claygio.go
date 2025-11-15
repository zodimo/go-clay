package claygio

import (
	"gioui.org/layout"
	"gioui.org/op"
	"github.com/zodimo/clay-go/clay"
)

var _ Renderer = (*ClayGioEngine)(nil)
var _ TextMeasurer = (*ClayGioEngine)(nil)

type ClayGioEngine struct {
	config *ClaygioConfig
	input  *GioInput
}

type ClaygioConfig struct {
	fontManager  *FontManager
	renderer     Renderer
	textMeasurer TextMeasurer
}
type ClaygioConfigOption func(*ClaygioConfig)

func ClaygioWithFontManager(fontManager *FontManager) ClaygioConfigOption {
	return func(c *ClaygioConfig) {
		c.fontManager = fontManager
	}
}

func ClaygioWithRenderer(renderer Renderer) ClaygioConfigOption {
	return func(c *ClaygioConfig) {
		c.renderer = renderer
	}
}

func ClaygioWithTextMeasurer(textMeasurer TextMeasurer) ClaygioConfigOption {
	return func(c *ClaygioConfig) {
		c.textMeasurer = textMeasurer
	}
}
func NewClayGioEngine(opts ...ClaygioConfigOption) *ClayGioEngine {
	defaultFontManager := NewFontManager()
	defaultRenderer := NewRenderer(RendererWithFontManager(defaultFontManager))
	defaultTextMeasurer := NewMeasurer(MeasurerWithFontManager(defaultFontManager))
	config := &ClaygioConfig{
		fontManager:  defaultFontManager,
		renderer:     defaultRenderer,
		textMeasurer: defaultTextMeasurer,
	}
	for _, opt := range opts {
		opt(config)
	}
	return &ClayGioEngine{
		config: config,
		input:  NewGioInput(),
	}
}

func (c *ClayGioEngine) Render(ops *op.Ops, commands []clay.Clay_RenderCommand) {
	c.config.renderer.Render(ops, commands)
}

func (c *ClayGioEngine) MeasureText(text clay.Clay_StringSlice, config *clay.Clay_TextElementConfig, userData interface{}) clay.Clay_Dimensions {
	return c.config.textMeasurer.MeasureText(text, config, userData)
}

func (c *ClayGioEngine) GetMousePosition() clay.Clay_Vector2 {
	return c.input.GetPointerPosition()
}

func (c *ClayGioEngine) GetMouseScrollDelta() clay.Clay_Vector2 {
	return clay.Clay_Vector2{}
}

func (c *ClayGioEngine) UpdateInput(gtx layout.Context) {
	c.input.Update(gtx)
}
