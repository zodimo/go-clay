package claygio

import (
	"gioui.org/gesture"
	"gioui.org/layout"
)

type GioInput struct {
	gtx   layout.Context
	click gesture.Click
}

func NewGioInput(gtx layout.Context) *GioInput {
	return &GioInput{
		gtx: gtx,
	}
}

// func (g *GioInput) GetMousePosition() clay.Clay_Vector2 {
// 	return clay.Clay_Vector2{
// 		X: float32(g.gtx.Queue.Input().Mouse.Position.X),
// 		Y: float32(g.gtx.Queue.Input().Mouse.Position.Y),
// 	}
// }
