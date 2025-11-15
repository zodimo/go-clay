package claygio

import (
	"image"
	"time"

	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/widget"
	"github.com/zodimo/clay-go/clay"
)

// Clickable represents a click tracker
type GioInput struct {
	pointerClick gesture.Click
	clickHistory []widget.Press

	requestClicks int
	pressedKey    key.Name

	pointerMove PointerMove
	// pointerScroll gesture.Scroll
}

func NewGioInput() *GioInput {
	return &GioInput{
		pointerClick: gesture.Click{},
		clickHistory: make([]widget.Press, 0),
		pointerMove:  PointerMove{},
	}
}

func (b *GioInput) GetPointerPosition() clay.Clay_Vector2 {
	return clay.Clay_Vector2{}
}

// Click executes a simple programmatic click.
func (b *GioInput) Click() {
	b.requestClicks++
}

// Clicked calls Update and reports whether a click was registered.
func (b *GioInput) Clicked(gtx layout.Context) bool {
	return b.clicked(gtx)
}

func (b *GioInput) clicked(gtx layout.Context) bool {
	_, clicked := b.clickableUpdate(gtx)
	return clicked
}

// Hovered reports whether a pointer is over the element.
func (b *GioInput) Hovered() bool {
	return b.pointerClick.Hovered()
}

// Pressed reports whether a pointer is pressing the element.
func (b *GioInput) Pressed() bool {
	return b.pointerClick.Pressed()
}

// History is the past pointer presses useful for drawing markers.
// History is retained for a short duration (about a second).
func (b *GioInput) History() []widget.Press {
	return b.clickHistory
}

// Layout and update the button state.
// func (b *GioInput) Layout(gtx layout.Context) {
// 	b.layout(gtx)
// }

// func (b *GioInput) layout(gtx layout.Context) {
// 	for {
// 		_, ok := b.clickableUpdate(gtx)
// 		if !ok {
// 			break
// 		}
// 	}

// 	defer clip.Rect(image.Rectangle{Max: gtx.Constraints.Max}).Push(gtx.Ops).Pop()
// 	semantic.EnabledOp(gtx.Enabled()).Add(gtx.Ops)
// 	b.pointerClick.Add(gtx.Ops)

// }

func (b *GioInput) add(ops *op.Ops) {
	b.pointerClick.Add(ops)
	b.pointerMove.Add(ops)
	// b.pointerScroll.Add(ops)

	event.Op(ops, b) // i don't know if we need this

}

// Update the button state by processing events, and return the next
// click, if any.
func (b *GioInput) Update(gtx layout.Context) {
	b.add(gtx.Ops)
	b.pointerMove.Update(gtx)
	// b.pointerScroll.Update(gtx)
	b.clickableUpdate(gtx)
}

func (b *GioInput) clickableUpdate(gtx layout.Context) (widget.Click, bool) {
	for len(b.clickHistory) > 0 {
		c := b.clickHistory[0]
		if c.End.IsZero() || gtx.Now.Sub(c.End) < 1*time.Second {
			break
		}
		n := copy(b.clickHistory, b.clickHistory[1:])
		b.clickHistory = b.clickHistory[:n]
	}
	if c := b.requestClicks; c > 0 {
		b.requestClicks = 0
		return widget.Click{
			NumClicks: c,
		}, true
	}
	for {
		e, ok := b.pointerClick.Update(gtx.Source)
		if !ok {
			break
		}
		switch e.Kind {
		case gesture.KindClick:
			if l := len(b.clickHistory); l > 0 {
				b.clickHistory[l-1].End = gtx.Now
			}
			return widget.Click{
				Modifiers: e.Modifiers,
				NumClicks: e.NumClicks,
			}, true
		case gesture.KindCancel:
			for i := range b.clickHistory {
				b.clickHistory[i].Cancelled = true
				if b.clickHistory[i].End.IsZero() {
					b.clickHistory[i].End = gtx.Now
				}
			}
		case gesture.KindPress:
			b.clickHistory = append(b.clickHistory, widget.Press{
				Position: e.Position,
				Start:    gtx.Now,
			})
		}
	}

	return widget.Click{}, false
}

type PointerMove struct {
	Position image.Point
	Time     time.Time
}

// Add the handler to the operation list to receive drag events.
func (pm *PointerMove) Add(ops *op.Ops) {
	event.Op(ops, pm)
}

func (pm *PointerMove) Update(gtx layout.Context) {
	for {
		evt, ok := gtx.Source.Event(pointer.Filter{
			Target: pm,
			Kinds:  pointer.Move,
		})
		if !ok {
			break
		}
		e, ok := evt.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Kind {
		case pointer.Move:
			pm.Position = e.Position.Round()
			pm.Time = gtx.Now
		}
	}
}
