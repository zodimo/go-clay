package clay

type Clay_BorderWidth struct {
	Left   uint16
	Right  uint16
	Top    uint16
	Bottom uint16
	// Creates borders between each child element, depending on the .layoutDirection.
	// e.g. for LEFT_TO_RIGHT, borders will be vertical lines, and for TOP_TO_BOTTOM borders will be horizontal lines.
	// .betweenChildren borders will result in individual RECTANGLE render commands being generated.
	BetweenChildren uint16
}

// Controls settings related to element borders.
type Clay_BorderElementConfig struct {
	Color Clay_Color       // Controls the color of all borders with width > 0. Conventionally represented as 0-255, but interpretation is up to the renderer.
	Width Clay_BorderWidth // Controls the widths of individual borders. At least one of these should be > 0 for a BORDER render command to be generated.

}
