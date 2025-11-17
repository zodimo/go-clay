package claygio

import (
	"image"
	"math"
	"runtime"
	"time"

	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"
	"gioui.org/widget"
	"github.com/zodimo/clay-go/clay"
	"github.com/zodimo/clay-go/pkg/claygio/fling"
)

const touchSlop = unit.Dp(3)

// Clickable represents a click tracker
type GioInput struct {
	pointerClick gesture.Click
	clickHistory []widget.Press

	requestClicks int
	pressedKey    key.Name

	pointerMove   PointerMove
	pointerScroll PointerScroll

	// clickedThisFrame tracks if a click occurred in the current frame
	clickedThisFrame bool
}

func NewGioInput() *GioInput {
	return &GioInput{
		pointerClick:  gesture.Click{},
		clickHistory:  make([]widget.Press, 0),
		pointerMove:   PointerMove{},
		pointerScroll: PointerScroll{},
	}
}

func (b *GioInput) GetPointerPosition() clay.Clay_Vector2 {
	return clay.Clay_Vector2{
		X: float32(b.pointerMove.Position.X),
		Y: float32(b.pointerMove.Position.Y),
	}
}

// Click executes a simple programmatic click.
func (b *GioInput) Click() {
	b.requestClicks++
}

// Clicked reports whether a click was registered in the current frame.
// This should be called after Update() has been called for the frame.
func (b *GioInput) Clicked() bool {
	return b.clickedThisFrame
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

func (b *GioInput) add(gtx layout.Context) {
	// Create a clip rectangle covering the entire window.
	// This makes the entire page clickable, matching the pattern from widget.Clickable.
	// The clip defines the hit area for pointer events.
	// Rectangle from (0,0) to window size
	clickArea := image.Rectangle{
		Min: image.Point{},
		Max: gtx.Constraints.Max,
	}
	defer clip.Rect(clickArea).Push(gtx.Ops).Pop()

	// Add the click gesture handler. This registers the gesture to receive
	// pointer events within the clip area above. The gesture.Click.Add() method
	// automatically handles event routing through the pointer queue.
	// gesture.Click.Add() internally calls event.Op(ops, c) to tag events with itself.
	b.pointerClick.Add(gtx.Ops)

	// Tag events with this GioInput instance. This is needed for:
	// 1. Keyboard event handling (if we add it in the future)
	// 2. Proper event routing in the Gio event system
	// This matches the pattern from widget.Clickable.layout()
	event.Op(gtx.Ops, b)

	// Add pointer move handler for tracking mouse position anywhere on the page
	b.pointerMove.Add(gtx.Ops)
	b.pointerScroll.Add(gtx.Ops)
}

// Update processes input events and updates the state.
// This must be called every frame before checking Clicked().
func (b *GioInput) Update(gtx layout.Context) {
	// Reset click state for this frame
	b.clickedThisFrame = false

	// Set up ops for the next frame (this registers the gesture handlers)
	b.add(gtx)

	// Update pointer move tracking
	b.pointerMove.Update(gtx)
	b.pointerScroll.Update(gtx.Metric, gtx.Source, gtx.Now, pointer.ScrollRange{Min: -1000, Max: 1000}, pointer.ScrollRange{Min: -1000, Max: 1000})
	// Process all click events - if any click occurred, set clickedThisFrame to true
	// This matches the pattern from widget.Clickable.layout() which processes
	// all pending clicks in a loop
	for {
		_, clicked := b.clickableUpdate(gtx)
		if !clicked {
			break
		}
		b.clickedThisFrame = true
	}
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

type PointerScroll struct {
	ScrollX   float32
	ScrollY   float32
	DeltaX    float32 // Accumulated integer delta for this frame
	DeltaY    float32 // Accumulated integer delta for this frame
	Time      time.Time
	DeltaTime time.Duration

	dragging   bool
	estimatorX fling.Extrapolation
	estimatorY fling.Extrapolation

	pid      pointer.ID
	lastX    int
	lastY    int
	flingerX fling.Animation
	flingerY fling.Animation
}

func (ps *PointerScroll) Add(ops *op.Ops) {
	event.Op(ops, ps)
}

// Stop any remaining fling movement.
func (s *PointerScroll) Stop() {
	s.flingerX = fling.Animation{}
	s.flingerY = fling.Animation{}
}

// Update state and report the scroll distance along axis.
func (s *PointerScroll) Update(cfg unit.Metric, q input.Source, t time.Time, scrollx, scrolly pointer.ScrollRange) {
	f := pointer.Filter{
		Target:  s,
		Kinds:   pointer.Press | pointer.Drag | pointer.Release | pointer.Scroll | pointer.Cancel,
		ScrollX: scrollx,
		ScrollY: scrolly,
	}

	s.DeltaTime = t.Sub(s.Time)
	s.Time = t

	// Reset delta for this frame
	s.DeltaX = 0
	s.DeltaY = 0

	for {
		evt, ok := q.Event(f)
		if !ok {
			break
		}
		e, ok := evt.(pointer.Event)
		if !ok {
			continue
		}
		switch e.Kind {
		case pointer.Press:
			if s.dragging {
				break
			}
			// Only scroll on touch drags, or on Android where mice
			// drags also scroll by convention.
			if e.Source != pointer.Touch && runtime.GOOS != "android" {
				break
			}
			s.Stop()
			s.estimatorX = fling.Extrapolation{}
			s.estimatorY = fling.Extrapolation{}
			vX := s.val(gesture.Horizontal, e.Position)
			vY := s.val(gesture.Vertical, e.Position)
			s.lastX = int(math.Round(float64(vX)))
			s.lastY = int(math.Round(float64(vY)))
			s.estimatorX.Sample(e.Time, vX)
			s.estimatorY.Sample(e.Time, vY)
			s.dragging = true
			s.pid = e.PointerID
		case pointer.Release:
			if s.pid != e.PointerID {
				break
			}
			flingX := s.estimatorX.Estimate()
			flingY := s.estimatorY.Estimate()
			if slop, d := float32(cfg.Dp(touchSlop)), flingX.Distance; d < -slop || d > slop {
				s.flingerX.Start(cfg, t, flingX.Velocity)
			}

			if slop, d := float32(cfg.Dp(touchSlop)), flingY.Distance; d < -slop || d > slop {
				s.flingerY.Start(cfg, t, flingY.Velocity)
			}
			fallthrough
		case pointer.Cancel:
			s.dragging = false
		case pointer.Scroll:
			s.ScrollX += e.Scroll.X
			s.ScrollY += e.Scroll.Y

			// Extract integer part as delta, keep fractional part for accumulation
			// Use math.Trunc (same as int() but explicit) to truncate towards zero
			// This matches Gio's gesture.Scroll pattern
			iscrollX := int(math.Trunc(float64(s.ScrollX)))
			iscrollY := int(math.Trunc(float64(s.ScrollY)))
			s.DeltaX += float32(iscrollX)
			s.DeltaY += float32(iscrollY)
			s.ScrollX -= float32(iscrollX)
			s.ScrollY -= float32(iscrollY)
		case pointer.Drag:
			if !s.dragging || s.pid != e.PointerID {
				continue
			}
			valX := s.val(gesture.Horizontal, e.Position)
			valY := s.val(gesture.Vertical, e.Position)
			s.estimatorX.Sample(e.Time, valX)
			s.estimatorY.Sample(e.Time, valY)
			vX := int(math.Round(float64(valX)))
			vY := int(math.Round(float64(valY)))

			distX := s.lastX - vX
			distY := s.lastY - vY
			if e.Priority < pointer.Grabbed {
				slop := cfg.Dp(touchSlop)
				if distX >= slop || -slop >= distX {
					q.Execute(pointer.GrabCmd{Tag: s, ID: e.PointerID})
				}
				if distY >= slop || -slop >= distY {
					q.Execute(pointer.GrabCmd{Tag: s, ID: e.PointerID})
				}
			} else {
				s.lastX = vX
				s.lastY = vY
			}
		}
	}

	// Update fling animations and add their contribution to scroll delta
	if s.flingerX.Active() {
		flingDelta := s.flingerX.Tick(t)
		s.DeltaX += float32(flingDelta)
	}
	if s.flingerY.Active() {
		flingDelta := s.flingerY.Tick(t)
		s.DeltaY += float32(flingDelta)
	}

	if s.flingerX.Active() || s.flingerY.Active() {
		q.Execute(op.InvalidateCmd{})
	}
}

func (s *PointerScroll) val(axis gesture.Axis, p f32.Point) float32 {
	switch axis {
	case gesture.Horizontal:
		return p.X
	case gesture.Vertical:
		return p.Y
	case gesture.Both:
		return p.X + p.Y
	default:
		return 0.0
	}
}
