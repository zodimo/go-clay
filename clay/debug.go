package clay

type Clay__DebugElementData struct {
	Collision bool
	Collapsed bool
}

var CLAY__DEBUGVIEW_COLOR_1 = Clay_Color{R: 58, G: 56, B: 52, A: 255}
var CLAY__DEBUGVIEW_COLOR_2 = Clay_Color{R: 62, G: 60, B: 58, A: 255}
var CLAY__DEBUGVIEW_COLOR_3 = Clay_Color{R: 141, G: 133, B: 135, A: 255}
var CLAY__DEBUGVIEW_COLOR_4 = Clay_Color{R: 238, G: 226, B: 231, A: 255}
var CLAY__DEBUGVIEW_COLOR_SELECTED_ROW = Clay_Color{R: 102, G: 80, B: 78, A: 255}

const CLAY__DEBUGVIEW_ROW_HEIGHT = 30
const CLAY__DEBUGVIEW_OUTER_PADDING = 10
const CLAY__DEBUGVIEW_INDENT_WIDTH = 16

var Clay__debugViewWidth int32 = 400

var Clay__debugViewHighlightColor = Clay_Color{R: 168, G: 66, B: 28, A: 100}

type Clay__RenderDebugLayoutData struct {
	RowCount                int32
	SelectedElementRowIndex int32
}
