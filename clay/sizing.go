package clay

// Controls how the element takes up space inside its parent container.
type Clay__SizingType uint8

const (
	// (default) Wraps tightly to the size of the element's contents.
	CLAY__SIZING_TYPE_FIT Clay__SizingType = iota
	// Expands along this axis to fill available space in the parent element, sharing it with other GROW elements.
	CLAY__SIZING_TYPE_GROW
	// Expects 0-1 range. Clamps the axis size to a percent of the parent container's axis size minus padding and child gaps.
	CLAY__SIZING_TYPE_PERCENT
	// Clamps the axis size to an exact size in pixels.
	CLAY__SIZING_TYPE_FIXED
)

// Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to,
// overriding sizing types such as FIT or GROW.
type Clay_SizingMinMax struct {
	Min float32
	Max float32
}

type Clay_SizingAxisSize struct {
	MinMax  Clay_SizingMinMax // Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to, overriding sizing types such as FIT or GROW.
	Percent float32           // Expects 0-1 range. Clamps the axis size to a percent of the parent container's axis size minus padding and child gaps.
}

// Controls the sizing of this element along one axis inside its parent container.
type Clay_SizingAxis struct {
	Type Clay__SizingType
	Size Clay_SizingAxisSize
}
