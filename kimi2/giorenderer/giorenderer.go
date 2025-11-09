package giorenderer

// import (
// 	"image"

// 	"github.com/zodimo/go-clay/kimi2/clay"

// 	"gioui.org/font/gofont"
// 	"gioui.org/layout"
// 	"gioui.org/text"
// 	"gioui.org/widget/material"
// )

// type Renderer struct {
// 	theme *material.Theme
// }

// func NewRenderer() *Renderer {

// 	th := material.NewTheme()
// 	th.Shaper = text.NewShaper(text.WithCollection(gofont.Collection()))
// 	return &Renderer{
// 		theme: th,
// 	}
// }

// func (r *Renderer) Layout(gtx layout.Context, commands []clay.RenderCommand) layout.Dimensions {
// 	var max image.Point
// 	for _, cmd := range commands {
// 		r.draw(gtx, cmd)
// 	}
// 	return layout.Dimensions{Size: max}
// }
