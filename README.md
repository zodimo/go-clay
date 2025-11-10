# Go-Clay: Go Bindings for Clay UI Layout Library

[![Go Reference](https://pkg.go.dev/badge/github.com/zodimo/go-clay.svg)](https://pkg.go.dev/github.com/zodimo/go-clay)

**Go-Clay** provides Go bindings for [Clay](https://github.com/nicbarker/clay), a high-performance 2D UI layout library written in C.

## About Clay

**_Clay_** (short for **C Layout**) is a high-performance 2D UI layout library with the following features:

- **Microsecond layout performance** - Optimized layout algorithms
- **Flexbox-like layout model** - Complex, responsive layouts including text wrapping, scrolling containers and aspect ratio scaling
- **Single ~4k LOC header file** - `clay.h` with **zero** dependencies (including no standard library)
- **WASM support** - Compile with clang to a 15kb uncompressed `.wasm` file for use in the browser
- **Static arena-based memory** - No malloc/free, low total memory overhead (~3.5mb for 8192 layout elements)
- **React-like nested declarative syntax** - Familiar component-based API
- **Renderer agnostic** - Outputs a sorted list of rendering primitives that can be easily composited in any 3D engine

For more information about Clay, visit:
- [Clay Website](https://nicbarker.com/clay) - Interactive examples
- [Original C Library](https://github.com/nicbarker/clay) - Source repository
- [Introduction Video](https://youtu.be/DYWTw19_8r4) - Overview and demo

## Installation

```bash
go get github.com/zodimo/go-clay
```

## Requirements

- Go 1.23+ with CGO enabled
- C compiler (gcc, clang, or MSVC)
- The `C/clay.h` header file (included in this repository)

## Quick Start

```go
package main

import (
    "github.com/zodimo/go-clay/clay"
    "unsafe"
)

func main() {
    // Calculate minimum memory size required
    totalMemorySize := clay.MinMemorySize()
    
    // Create arena with required memory
    // Note: In production, use a proper memory allocator
    memory := make([]byte, totalMemorySize)
    arena := clay.CreateArenaWithCapacityAndMemory(
        uintptr(totalMemorySize),
        unsafe.Pointer(&memory[0]),
    )
    
    // Initialize Clay with screen dimensions
    dimensions := clay.Dimensions{
        Width:  1920,
        Height: 1080,
    }
    
    errorHandler := clay.ErrorHandler{
        ErrorHandlerFunction: handleClayErrors,
        UserData:             0,
    }
    
    context := clay.Initialize(arena, dimensions, errorHandler)
    if context == nil {
        panic("Failed to initialize Clay")
    }
    
    // Set layout dimensions (call on window resize)
    clay.SetLayoutDimensions(dimensions)
    
    // Set pointer state for mouse/touch interactions
    pointerPos := clay.Vector2{X: 100, Y: 200}
    clay.SetPointerState(pointerPos, false)
    
    // Begin layout declaration
    clay.BeginLayout()
    
    // Declare your UI hierarchy here
    // (See examples below)
    
    // End layout and get render commands
    renderCommands := clay.EndLayout()
    
    // Render the commands using your renderer
    for i := 0; i < int(renderCommands.Length); i {
        cmd := renderCommands.InternalArray[i]
        // Handle different command types
        switch cmd.CommandType {
        case clay.RenderCommandTypeRectangle:
            // Render rectangle
        case clay.RenderCommandTypeText:
            // Render text
        // ... other command types
        }
    }
}

func handleClayErrors(errorData clay.ErrorData) {
    // Handle Clay errors
    switch errorData.ErrorType {
    case clay.ErrorTypeArenaCapacityExceeded:
        // Handle memory issues
    case clay.ErrorTypeElementsCapacityExceeded:
        // Handle too many elements
    // ... other error types
    }
}
```

## Usage Pattern

The general order of operations when using Clay:

1. **Initialize** - Call `clay.Initialize()` with arena and dimensions
2. **Set Layout Dimensions** - Call `clay.SetLayoutDimensions()` (e.g., on window resize)
3. **Set Pointer State** - Call `clay.SetPointerState()` for mouse/touch input
4. **Update Scroll Containers** - Call `clay.UpdateScrollContainers()` for scrolling
5. **Begin Layout** - Call `clay.BeginLayout()`
6. **Declare UI** - Use Clay's declarative API to build your UI hierarchy
7. **End Layout** - Call `clay.EndLayout()` to get render commands
8. **Render** - Process the render commands with your renderer

## Building UI Hierarchies

Clay uses a declarative, React-like syntax. In Go, you'll work with the generated bindings:

```go
// Example: Creating a container with children
clay.BeginLayout()

// Use Clay's macros and functions through the bindings
// Note: The exact API depends on the generated bindings
// Refer to the generated clay package documentation

clay.EndLayout()
```

## Configuration

This project uses [c-for-go](https://github.com/xlab/c-for-go) to generate Go bindings from the C header file.

### Regenerating Bindings

If you need to regenerate the bindings (e.g., after updating `C/clay.h`):

```bash
make generate
```

Or manually:

```bash
c-for-go clay.yml
```

### Makefile Commands

- `make` or `make all` - Generate Go bindings
- `make generate` - Explicitly generate bindings
- `make clean` - Remove all generated files
- `make test` - Build the package to verify it compiles

## Project Structure

```
go-clay/
├── C/
│   └── clay.h          # Original Clay C header file
├── clay/               # Generated Go bindings (do not edit)
│   ├── clay.go
│   ├── const.go
│   ├── doc.go
│   └── cgo_helpers.h
├── clay.yml            # c-for-go configuration
├── Makefile            # Build automation
└── README.md           # This file
```

## Memory Management

Clay uses arena-based memory allocation. You must provide a memory buffer of at least `MinMemorySize()` bytes:

```go
// Calculate required memory
minSize := clay.MinMemorySize()

// Allocate memory (use your preferred allocator)
memory := make([]byte, minSize)

// Create arena
arena := clay.CreateArenaWithCapacityAndMemory(
    uintptr(minSize),
    unsafe.Pointer(&memory[0]),
)
```

## Error Handling

Clay provides error callbacks for various error conditions:

```go
errorHandler := clay.ErrorHandler{
    ErrorHandlerFunction: func(errorData clay.ErrorData) {
        switch errorData.ErrorType {
        case clay.ErrorTypeArenaCapacityExceeded:
            // Memory arena too small
        case clay.ErrorTypeElementsCapacityExceeded:
            // Too many UI elements
        case clay.ErrorTypeTextMeasurementFunctionNotProvided:
            // Text measurement function missing
        // ... handle other error types
        }
    },
    UserData: 0, // Optional user data
}
```

## Text Measurement

Clay requires a text measurement function for text rendering:

```go
// Set text measurement function
clay.SetMeasureTextFunction(
    func(text clay.StringSlice, config *clay.TextElementConfig, userData unsafe.Pointer) clay.Dimensions {
        // Measure text based on config (fontId, fontSize, etc.)
        // Return dimensions
        return clay.Dimensions{
            Width:  float32(text.Length) * float32(config.FontSize), // Simplified
            Height: float32(config.FontSize),
        }
    },
    nil, // userData
)
```

## Examples

For comprehensive examples, see:
- [Clay Examples Directory](https://github.com/nicbarker/clay/tree/main/examples) - Original C examples
- [Clay Website](https://nicbarker.com/clay) - Interactive browser examples

## Documentation

- [Clay C Library README](https://github.com/nicbarker/clay/blob/main/README.md) - Comprehensive documentation
- [Generated Go Documentation](https://pkg.go.dev/github.com/zodimo/go-clay/clay) - Go API reference (when published)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) file for details.

This project includes the Clay C library header file, which is also licensed under MIT.

## Acknowledgments

- [Clay](https://github.com/nicbarker/clay) by [nicbarker](https://github.com/nicbarker) - The original C library
- [c-for-go](https://github.com/xlab/c-for-go) - C to Go bindings generator

## Support

For questions about:
- **Clay library itself**: Join the [Clay Discord server](https://discord.gg/b4FTWkxdvT)
- **Go bindings**: Open an issue in this repository
