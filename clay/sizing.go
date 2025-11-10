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

// Controls the sizing of this element along one axis inside its parent container.
type Clay_SizingAxis struct {
	Type Clay__SizingType

	// SizeMinMax is mutually exclusive with SizePercent.

	// Controls the minimum and maximum size in pixels that this element is allowed to grow or shrink to, overriding sizing types such as FIT or GROW.
	SizeMinMax Clay_SizingMinMax
	// Expects 0-1 range. Clamps the axis size to a percent of the parent container's axis size minus padding and child gaps.
	SizePercent float32
}

// Controls the sizing of this element along one axis inside its parent container.
type Clay_Sizing struct {
	// Controls the width sizing of the element, along the x axis.
	Width Clay_SizingAxis
	// Controls the height sizing of the element, along the y axis.
	Height Clay_SizingAxis
}

func CLAY_SIZING_FIT(minMax Clay_SizingMinMax) Clay_SizingAxis {
	return Clay_SizingAxis{
		Type:       CLAY__SIZING_TYPE_FIT,
		SizeMinMax: minMax,
	}
}

func CLAY_SIZING_GROW(minMax Clay_SizingMinMax) Clay_SizingAxis {
	return Clay_SizingAxis{
		Type:       CLAY__SIZING_TYPE_GROW,
		SizeMinMax: minMax,
	}
}
func CLAY_SIZING_FIXED(fixedSize float32) Clay_SizingAxis {
	return Clay_SizingAxis{
		Type:       CLAY__SIZING_TYPE_FIXED,
		SizeMinMax: Clay_SizingMinMax{Min: fixedSize, Max: fixedSize},
	}
}

func CLAY_SIZING_PERCENT(percentOfParent float32) Clay_SizingAxis {
	if percentOfParent < 0 || percentOfParent > 1 {
		panic("percentOfParent must be between 0 and 1")
	}
	return Clay_SizingAxis{
		Type:        CLAY__SIZING_TYPE_PERCENT,
		SizePercent: percentOfParent,
	}
}
