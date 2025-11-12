package clay

type Clay_Dimensions struct {
	Width  float32
	Height float32
}

type Clay_Padding struct {
	Left   uint16
	Right  uint16
	Top    uint16
	Bottom uint16
}

type Clay_LayoutDirection uint8

const (
	// (Default) Lays out child elements from left to right with increasing x.
	CLAY_LEFT_TO_RIGHT Clay_LayoutDirection = iota
	// Lays out child elements from top to bottom with increasing y.
	CLAY_TOP_TO_BOTTOM
)

type Clay_ChildAlignment struct {
	TopLeft     float32
	TopRight    float32
	BottomLeft  float32
	BottomRight float32
}

// Controls various settings that affect the size and position of an element, as well as the sizes and positions
// of any child elements.
type Clay_LayoutConfig struct {
	Sizing          Clay_Sizing
	Padding         Clay_Padding
	ChildGap        uint16
	ChildAlignment  Clay_ChildAlignment
	LayoutDirection Clay_LayoutDirection
}

// Controls the sizing of this element along one axis inside its parent container.
type Clay_Sizing struct {
	// Controls the width sizing of the element, along the x axis.
	Width Clay_SizingAxis
	// Controls the height sizing of the element, along the y axis.
	Height Clay_SizingAxis
}
