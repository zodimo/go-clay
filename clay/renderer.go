package clay

// Renderer interface for rendering commands
// This interface matches the Clay C architecture where renderers
// receive complete command arrays with bounds information
type Renderer interface {
	// Unified rendering method that processes command arrays with bounds
	// This matches the Clay C pattern: Clay_Raylib_Render(Clay_RenderCommandArray renderCommands, Font* fonts)
	Render(commands []RenderCommand) error

	// Frame lifecycle management
	BeginFrame() error
	EndFrame() error
	SetViewport(bounds BoundingBox) error
}

// Legacy interface for backward compatibility (deprecated)
// TODO: Remove after all renderers are updated to use unified Render method
type LegacyRenderer interface {
	// Individual render methods (deprecated - use Render method instead)
	RenderRectangle(cmd RectangleCommand) error
	RenderText(cmd TextCommand) error
	RenderImage(cmd ImageCommand) error
	RenderBorder(cmd BorderCommand) error
	RenderClipStart(cmd ClipStartCommand) error
	RenderClipEnd(cmd ClipEndCommand) error
	RenderCustom(cmd CustomCommand) error

	// State management
	BeginFrame() error
	EndFrame() error
	SetViewport(bounds BoundingBox) error
}
