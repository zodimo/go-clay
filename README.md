# Go-Clay: UI Layout Library for Go

A high-performance, render engine agnostic UI layout library for Go, inspired by the original Clay C library.

## Features

- **Microsecond layout performance** - Optimized layout algorithms
- **Flexbox-like layout model** - Complex responsive layouts with text wrapping and scrolling
- **Render engine agnostic** - Works with any rendering backend
- **Memory efficient** - Arena-based allocation with predictable memory usage
- **Declarative API** - React-like nested component syntax
- **Zero dependencies** - Core library has no external dependencies
- **Cross-platform** - Works on all Go-supported platforms

## Quick Start

```go
package main

import (
    "github.com/zodimo/go-clay"
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

## Installation

```bash
go get github.com/zodimo/go-clay
```

## Renderers

### Gio UI Renderer
Native Go UI framework for cross-platform applications.

```bash
go get github.com/zodimo/go-clay/renderers/gioui
```

**Version Compatibility:**
- **v0.8.x**: Use `gioui_renderer.go` (default)
- **v0.9.0**: Use `gioui_renderer_v0.9.0.go` for latest Gio UI compatibility

### OpenGL Renderer
OpenGL/WebGL rendering for games and graphics applications.

```bash
go get github.com/zodimo/go-clay/renderers/opengl
```

### SDL Renderer
SDL2/SDL3 rendering for cross-platform applications.

```bash
go get github.com/zodimo/go-clay/renderers/sdl
```

### Terminal Renderer
Terminal-based UI for command-line applications.

```bash
go get github.com/zodimo/go-clay/renderers/terminal
```

## Examples

### Basic Layout
```go
engine := clay.NewLayoutEngine()
engine.BeginLayout()

clay.Container("main", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingGrow(0),
            Height: clay.SizingGrow(0),
        },
        Direction: clay.LeftToRight,
        ChildGap:  16,
    },
}).
    Container("sidebar", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingFixed(200),
                Height: clay.SizingGrow(0),
            },
        },
        BackgroundColor: clay.Color{R: 0.8, G: 0.8, B: 0.8, A: 1.0},
    }).
        Text("Sidebar", clay.TextConfig{FontSize: 16}).
    Container("content", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingGrow(0),
                Height: clay.SizingGrow(0),
            },
        },
        BackgroundColor: clay.Color{R: 1.0, G: 1.0, B: 1.0, A: 1.0},
    }).
        Text("Main Content", clay.TextConfig{FontSize: 16})

commands := engine.EndLayout()
```

### Responsive Grid
```go
clay.Container("grid", clay.ElementConfig{
    Layout: clay.LayoutConfig{
        Sizing: clay.Sizing{
            Width:  clay.SizingGrow(0),
            Height: clay.SizingGrow(0),
        },
        Direction: clay.LeftToRight,
        ChildGap:  8,
    },
}).
    Container("item1", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingPercent(0.33),
                Height: clay.SizingFixed(100),
            },
        },
        BackgroundColor: clay.Color{R: 0.9, G: 0.5, B: 0.5, A: 1.0},
    }).
        Text("Item 1", clay.TextConfig{FontSize: 16}).
    Container("item2", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingPercent(0.33),
                Height: clay.SizingFixed(100),
            },
        },
        BackgroundColor: clay.Color{R: 0.5, G: 0.9, B: 0.5, A: 1.0},
    }).
        Text("Item 2", clay.TextConfig{FontSize: 16}).
    Container("item3", clay.ElementConfig{
        Layout: clay.LayoutConfig{
            Sizing: clay.Sizing{
                Width:  clay.SizingPercent(0.34),
                Height: clay.SizingFixed(100),
            },
        },
        BackgroundColor: clay.Color{R: 0.5, G: 0.5, B: 0.9, A: 1.0},
    }).
        Text("Item 3", clay.TextConfig{FontSize: 16})
```

## Layout System

### Sizing Types

- **SizingFit** - Wrap to content size
- **SizingGrow** - Fill available space
- **SizingPercent** - Percentage of parent
- **SizingFixed** - Fixed pixel size

### Layout Directions

- **LeftToRight** - Horizontal layout
- **TopToBottom** - Vertical layout

### Alignment

- **X**: Left, Center, Right
- **Y**: Top, Center, Bottom

## Performance

### Memory Management
```go
// Use arena-based allocation for large layouts
arena := clay.NewArena(10 * 1024 * 1024) // 10MB
engine := clay.NewLayoutEngineWithArena(arena)

// Reset arena between frames
arena.Reset()
```

### Performance Metrics
```go
stats := engine.GetStats()
fmt.Printf("Elements: %d\n", stats.ElementCount)
fmt.Printf("Render Commands: %d\n", stats.RenderCommands)
fmt.Printf("Layout Time: %v\n", stats.LayoutTime)
fmt.Printf("Memory Used: %d bytes\n", stats.MemoryUsed)
```

## Custom Renderers

Create custom renderers by implementing the `Renderer` interface (following Clay C architecture):

```go
type MyRenderer struct {
    // Your renderer state
}

func (r *MyRenderer) Render(commands []clay.RenderCommand) error {
    for _, cmd := range commands {
        bounds := cmd.BoundingBox  // Position and size information
        
        switch cmd.CommandType {
        case clay.CommandRectangle:
            rectCmd := cmd.Data.(clay.RectangleCommand)
            // Use bounds for positioning, rectCmd for styling
            // ... render rectangle at bounds position with rectCmd color
        case clay.CommandText:
            textCmd := cmd.Data.(clay.TextCommand)
            // Use bounds for positioning, textCmd for content and styling
            // ... render text at bounds position with textCmd properties
        }
    }
    return nil
}

func (r *MyRenderer) BeginFrame() error {
    // Initialize frame rendering
    return nil
}

func (r *MyRenderer) EndFrame() error {
    // Finalize frame rendering
    return nil
}

func (r *MyRenderer) RenderRectangle(cmd clay.RectangleCommand) error {
    // Render rectangle
    return nil
}

// ... implement other methods
```

## Documentation

### Project Documentation (BMad Method)
- [Project Brief](docs/brief.md) - Project overview, goals, and scope
- [Product Requirements](docs/prd.md) - Functional and non-functional requirements
- [Architecture Document](docs/architecture.md) - Technical architecture and design decisions

### Developer Documentation
- [API Reference](docs/api-reference.md) - Complete API documentation
- [Layout System](docs/layout-system.md) - Flexbox-like layout concepts
- [Renderer Guide](docs/renderer-guide.md) - Creating custom renderers
- [Examples](docs/examples.md) - Usage examples and patterns
- [Performance Guide](docs/performance.md) - Optimization tips

## Examples

See the [examples directory](examples/) for complete working examples:

- [Gio UI Example](examples/gioui_example/) - Native Go UI framework
- [OpenGL Example](examples/opengl_example/) - OpenGL/WebGL rendering
- [SDL Example](examples/sdl_example/) - SDL2/SDL3 rendering
- [Terminal Example](examples/terminal_example/) - Terminal-based UI

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by the original [Clay C library](https://github.com/nicbarker/clay)
- Built for the Go community
- Render engine agnostic design