package clay

// Used by renderers to determine specific handling for each render command.
type Clay_RenderCommandType uint8

const (
	CLAY_RENDER_COMMAND_TYPE_NONE Clay_RenderCommandType = iota
	CLAY_RENDER_COMMAND_TYPE_RECTANGLE
	CLAY_RENDER_COMMAND_TYPE_BORDER
	CLAY_RENDER_COMMAND_TYPE_TEXT
	CLAY_RENDER_COMMAND_TYPE_IMAGE
	CLAY_RENDER_COMMAND_TYPE_SCISSOR_START
	CLAY_RENDER_COMMAND_TYPE_SCISSOR_END
	CLAY_RENDER_COMMAND_TYPE_CUSTOM
)

type Clay_RenderCommand struct {
	// A rectangular box that fully encloses this UI element, with the position relative to the root of the layout.
	BoundingBox Clay_BoundingBox
	// A struct union containing data specific to this command's commandType.
	RenderData Clay_RenderData
	// A pointer transparently passed through from the original element declaration.
	UserData interface{}
	// The id of this element, transparently passed through from the original element declaration.
	Id uint32
	// The z order required for drawing this command correctly.
	// Note: the render command array is already sorted in ascending order, and will produce correct results if drawn in naive order.
	// This field is intended for use in batching renderers for improved performance.
	ZIndex int16
	// Specifies how to handle rendering of this command.
	// CLAY_RENDER_COMMAND_TYPE_RECTANGLE - The renderer should draw a solid color rectangle.
	// CLAY_RENDER_COMMAND_TYPE_BORDER - The renderer should draw a colored border inset into the bounding box.
	// CLAY_RENDER_COMMAND_TYPE_TEXT - The renderer should draw text.
	// CLAY_RENDER_COMMAND_TYPE_IMAGE - The renderer should draw an image.
	// CLAY_RENDER_COMMAND_TYPE_SCISSOR_START - The renderer should begin clipping all future draw commands, only rendering content that falls within the provided boundingBox.
	// CLAY_RENDER_COMMAND_TYPE_SCISSOR_END - The renderer should finish any previously active clipping, and begin rendering elements in full again.
	// CLAY_RENDER_COMMAND_TYPE_CUSTOM - The renderer should provide a custom implementation for handling this render command based on its .customData
	CommandType Clay_RenderCommandType
}

type Clay_BoundingBox struct {
	X      float32
	Y      float32
	Width  float32
	Height float32
}

type Clay_RenderData struct {
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_RECTANGLE
	Rectangle Clay_RectangleRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_TEXT
	Text Clay_TextRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_IMAGE
	Image Clay_ImageRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_CUSTOM
	Custom Clay_CustomRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_BORDER
	Border Clay_BorderRenderData
	// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_START|END
	Clip Clay_ClipRenderData
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_RECTANGLE
type Clay_RectangleRenderData struct {
	// The solid background color to fill this rectangle with. Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	BackgroundColor Clay_Color
	// Controls the "radius", or corner rounding of elements, including rectangles, borders and images.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	CornerRadius Clay_CornerRadius
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_TEXT
type Clay_TextRenderData struct {
	// A string slice containing the text to be rendered.
	// Note: this is not guaranteed to be null terminated.
	StringContents Clay_StringSlice
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	TextColor Clay_Color
	// An integer representing the font to use to render this text, transparently passed through from the text declaration.
	FontId   uint16
	FontSize uint16
	// Specifies the extra whitespace gap in pixels between each character.
	LetterSpacing uint16
	// The height of the bounding box for this line of text.
	LineHeight uint16
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_IMAGE
type Clay_ImageRenderData struct {
	// The tint color for this image. Note that the default value is 0,0,0,0 and should likely be interpreted
	// as "untinted".
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	BackgroundColor Clay_Color
	// Controls the "radius", or corner rounding of this image.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	CornerRadius Clay_CornerRadius
	// A pointer transparently passed through from the original element definition, typically used to represent image data.
	ImageData interface{}
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_CUSTOM
type Clay_CustomRenderData struct {
	// Passed through from .backgroundColor in the original element declaration.
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	BackgroundColor Clay_Color
	// Controls the "radius", or corner rounding of this custom element.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	CornerRadius Clay_CornerRadius
	// A pointer transparently passed through from the original element definition.
	CustomData interface{}
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_BORDER
type Clay_BorderRenderData struct {
	// Controls a shared color for all this element's borders.
	// Conventionally represented as 0-255 for each channel, but interpretation is up to the renderer.
	Color Clay_Color
	// Specifies the "radius", or corner rounding of this border element.
	// The rounding is determined by drawing a circle inset into the element corner by (radius, radius) pixels.
	CornerRadius Clay_CornerRadius
	// Controls individual border side widths.
	Width Clay_BorderWidth
}

// Render command data when commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_START || commandType == CLAY_RENDER_COMMAND_TYPE_SCISSOR_END
type Clay_ClipRenderData struct {
	Horizontal bool
	Vertical   bool
}
