# Go-Clay: UI Layout Library for Go

A high-performance, render engine agnostic UI layout library for Go, inspired by the original Clay C library.

## Features

- **Microsecond layout performance** - Optimized layout algorithms
- **Flexbox-like layout model** - Complex responsive layouts with text wrapping and scrolling
- **Render engine agnostic** - Works with any rendering backend
- **Memory efficient** - Arena-based allocation with predictable memory usage
- **Declarative API** - React-like nested component syntax
- **Zero dependencies** - Core library has no external dependencies

## Quick Start

```go
package main

import (
    "github.com/zodimo/go-clay/clay"
    "github.com/zodimo/go-clay/renderers/gioui"
)

func main() {
    // Create layout engine
    engine := clay.NewLayoutEngine()
    
    // Create Gio UI renderer
    renderer := gioui.NewRenderer()
    
    // Set up layout
    engine.BeginLayout()
    
    // Declare UI elements
    clay.Container("main", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingGrow(0),
                Height: clay.SizingGrow(0),
            },
            Padding: clay.PaddingAll(16),
        },
        BackgroundColor: clay.Color{R: 0.9, G: 0.9, B: 0.9, A: 1.0},
    }).Text("Hello, Go-Clay!", clay.TextConfig{
        FontSize: 24,
        Color:    clay.Color{R: 0, G: 0, B: 0, A: 1.0},
    })
    
    // Generate render commands
    commands := engine.EndLayout()
    
    // Render with Gio UI
    renderer.Render(commands)
}
```

## Architecture

Go-Clay separates layout computation from rendering:

1. **Layout Engine** - Computes element positions and sizes
2. **Render Commands** - Output primitives for any renderer
3. **Renderer Interface** - Pluggable rendering backends

## Documentation

### Core Documentation
- [API Reference](api-reference.md) - Complete API documentation
- [Layout System](layout-system.md) - Flexbox-like layout concepts
- [Renderer Guide](renderer-guide.md) - Creating custom renderers
- [Examples](examples.md) - Usage examples and patterns
- [Performance Guide](performance.md) - Optimization tips

### Implementation Analysis
- [Gio v0.9.0 Analysis](gio-analysis.md) - Gio rendering patterns and API structure

### Project Management
- [EPIC: Gio Renderer Implementation](EPIC-Gio-Renderer.md) - Current development epic
- [Project Brief](brief.md) - Project overview and goals
- [Product Requirements](prd.md) - Detailed requirements specification
- [Architecture](architecture.md) - System design and architecture

## Renderers

- [Gio UI Renderer](renderers/gioui.md) - Native Go UI framework
- [OpenGL Renderer](renderers/opengl.md) - OpenGL/WebGL rendering
- [SDL Renderer](renderers/sdl.md) - SDL2/SDL3 rendering
- [Terminal Renderer](renderers/terminal.md) - Terminal-based UI

## License

MIT License - see [LICENSE](../LICENSE) file for details.
