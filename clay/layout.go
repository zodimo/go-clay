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

// Controls the alignment along the x axis (horizontal) of child elements.
type Clay_LayoutAlignmentX uint8

const (
	// (Default) Aligns child elements to the left hand side of this element, offset by padding.width.left
	CLAY_ALIGN_X_LEFT Clay_LayoutAlignmentX = iota
	// Aligns child elements to the right hand side of this element, offset by padding.width.right
	CLAY_ALIGN_X_RIGHT
	// Aligns child elements horizontally to the center of this element
	CLAY_ALIGN_X_CENTER
)

// Controls the alignment along the y axis (vertical) of child elements.
type Clay_LayoutAlignmentY uint8

const (
	// (Default) Aligns child elements to the top of this element, offset by padding.width.top
	CLAY_ALIGN_Y_TOP Clay_LayoutAlignmentY = iota
	// Aligns child elements to the bottom of this element, offset by padding.width.bottom
	CLAY_ALIGN_Y_BOTTOM
	// Aligns child elements vertically to the center of this element
	CLAY_ALIGN_Y_CENTER
)

// Controls how child elements are aligned on each axis.
type Clay_ChildAlignment struct {
	X Clay_LayoutAlignmentX
	Y Clay_LayoutAlignmentY
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
