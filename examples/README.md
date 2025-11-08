# Clay Examples

This directory contains working examples of the Clay UI library with the Gio UI renderer.

## Available Examples

### Basic Examples

1. **simple-container** - A basic example showing a colored container with text
2. **sidebar-layout** - A two-column layout with sidebar navigation and main content
3. **responsive-grid** - A three-column grid layout with colored items

### Advanced Examples

4. **form-layout** - A form with input fields and a submit button

## Running the Examples

Each example is in its own directory and can be run independently:

```bash
# Run the simple container example
cd simple-container
go run main.go

# Run the sidebar layout example
cd sidebar-layout
go run main.go

# Run the responsive grid example
cd responsive-grid
go run main.go

# Run the form layout example
cd form-layout
go run main.go
```

## Current Implementation Status

### ✅ Working Features
- **Layout Engine**: Properly calculates element positions and sizes
- **Rectangle Rendering**: Colored backgrounds with proper bounds
- **Layout Types**: Fixed sizing, grow sizing, percentage sizing
- **Layout Direction**: Left-to-right and top-to-bottom layouts
- **Padding and Margins**: Proper spacing between elements
- **Child Alignment**: Center, left, right alignment
- **Borders**: Basic border rendering (partial)

### 🚧 Partially Working Features
- **Text Rendering**: Layout calculation works, but text display is stubbed
- **Border Rendering**: Basic support, but advanced features may not work

### ❌ Not Yet Implemented
- **Interactive Elements**: Buttons, hover states, click handling
- **Scrolling**: Scrollable containers
- **Images**: Image rendering
- **Advanced Text**: Font loading, text wrapping, text selection
- **Animations**: Transitions and animations

## Architecture Notes

### Current Workaround
The examples use a custom `renderRectangleWithBounds` function to properly render rectangles with bounds, as the current Gio renderer doesn't use the bounding box information from Clay's render commands.

### Key API Patterns
1. **Engine Setup**: Create engine, set dimensions, begin layout
2. **Element Creation**: Use `engine.OpenElement()` and `engine.CloseElement()`
3. **Text Elements**: Must set `LineHeight` in `TextConfig` for proper sizing
4. **Layout Sizing**: Use `SizingFixed()`, `SizingGrow()`, `SizingPercent()`
5. **Render Loop**: Process render commands and dispatch to appropriate renderers

## Example Structure

Each example follows this pattern:

```go
func main() {
    // Create Gio window
    w := &app.Window{}
    w.Option(app.Title("Example"), app.Size(...))
    
    // Run main loop
    go func() {
        run(w)
    }()
    app.Main()
}

func run(w *app.Window) error {
    var ops op.Ops
    
    for {
        switch e := w.Event().(type) {
        case app.FrameEvent:
            // Create Clay engine and Gio renderer
            engine := clay.NewLayoutEngine()
            renderer := gioui.NewRenderer(&ops)
            
            // Set up layout
            engine.SetLayoutDimensions(...)
            engine.BeginLayout()
            
            // Build UI
            createLayout(engine)
            
            // Render (using correct Clay C architecture)
            commands := engine.EndLayout()
            renderer.BeginFrame()
            renderer.Render(commands)  // Single unified call with bounds
            renderer.EndFrame()
            
            e.Frame(gtx.Ops)
        }
    }
}
```

## Contributing

When adding new examples:

1. Create a new directory under `examples/`
2. Follow the existing pattern and API usage
3. Include the `renderRectangleWithBounds` helper function
4. Set `LineHeight` for all text elements
5. Use proper `OpenElement`/`CloseElement` pairs
6. Test that the example runs without errors

## Dependencies

- **Clay**: `github.com/zodimo/go-clay`
- **Gio UI**: `gioui.org` (for windowing and rendering)
- **Go**: 1.23+ (as specified in go.mod)
