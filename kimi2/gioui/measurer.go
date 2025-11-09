package gioui

import "github.com/zodimo/go-clay/kimi2/clay"

type measurer struct{}

func NewMeasurer() clay.TextMeasurer {
	return &measurer{}
}

func (measurer) MeasureText(text string, cfg clay.TextElementConfig) clay.Dimensions {
	return clay.Dimensions{Width: float32(len(text)) * cfg.FontSize * 0.6, Height: cfg.FontSize * cfg.LineHeight}
}
