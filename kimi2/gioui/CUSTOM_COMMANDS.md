# Custom Commands API Documentation

## Overview

The Gio renderer for Clay provides a comprehensive custom command system that allows developers to extend the rendering capabilities beyond the built-in Clay commands. This system supports three types of custom command handlers: **Callback**, **Operation**, and **Plugin** (plugin support is planned for future implementation).

## Architecture

### Core Components

1. **CustomCommandRegistry**: Manages registration and execution of custom command handlers
2. **CustomCommandHandler Interface**: Defines the contract for all custom command handlers
3. **Handler Types**: Three different handler implementations for various use cases
4. **CustomCommandData**: Structured data format for complex custom commands

### Handler Types

#### 1. Callback Handlers (`CustomCommandTypeCallback`)
Simple function-based handlers that receive the full `CustomCommand` and `op.Ops`.

```go
handler := NewCallbackCommandHandler("my_command", func(ops *op.Ops, cmd clay.CustomCommand) error {
    // Custom rendering logic here
    return nil
})
```

#### 2. Operation Handlers (`CustomCommandTypeOperation`)
Parameter-based handlers that work with extracted parameters from the command data.

```go
handler := NewOperationCommandHandler("my_operation", func(ops *op.Ops, params map[string]interface{}) error {
    // Use params for rendering
    if value, ok := params["color"]; ok {
        // Handle color parameter
    }
    return nil
})
```

#### 3. Plugin Handlers (`CustomCommandTypePlugin`)
External plugin-based handlers (planned for future implementation).

```go
handler := NewPluginCommandHandler("my_plugin", "/path/to/plugin")
// Plugin loading not yet implemented
```

## Usage Guide

### Basic Setup

1. **Create a Renderer**:
```go
ops := new(op.Ops)
renderer := NewRenderer(ops)
```

2. **Register Custom Handlers**:
```go
// Register a callback handler
handler := NewCallbackCommandHandler("debug_overlay", func(ops *op.Ops, cmd clay.CustomCommand) error {
    // Draw debug overlay
    return nil
})

err := renderer.RegisterCustomHandler(handler)
if err != nil {
    log.Fatal(err)
}
```

3. **Use Custom Commands in Clay Layout**:
```go
// Method 1: Simple string command ID
customCmd := clay.CustomCommand{
    CustomData: "debug_overlay",
}

// Method 2: Map with parameters
customCmd := clay.CustomCommand{
    CustomData: map[string]interface{}{
        "commandID": "debug_overlay",
        "opacity":   0.5,
        "color":     "red",
    },
}

// Method 3: Structured data
customCmd := clay.CustomCommand{
    CustomData: CustomCommandData{
        CommandID: "debug_overlay",
        Parameters: map[string]interface{}{
            "opacity": 0.5,
            "color":   "red",
        },
        Metadata: map[string]string{
            "version": "1.0",
            "author":  "developer",
        },
    },
}

// Render the custom command
err := renderer.RenderCustom(customCmd)
```

### Advanced Usage Patterns

#### 1. Complex Shape Rendering

```go
shapeHandler := NewCallbackCommandHandler("custom_shape", func(ops *op.Ops, cmd clay.CustomCommand) error {
    // Extract shape parameters
    data, ok := cmd.CustomData.(map[string]interface{})
    if !ok {
        return errors.New("invalid shape data")
    }
    
    shapeType := data["type"].(string)
    
    switch shapeType {
    case "star":
        return renderStar(ops, data)
    case "hexagon":
        return renderHexagon(ops, data)
    default:
        return fmt.Errorf("unknown shape type: %s", shapeType)
    }
})
```

#### 2. Performance Profiling

```go
profilerHandler := NewOperationCommandHandler("profiler", func(ops *op.Ops, params map[string]interface{}) error {
    // Start timing
    start := time.Now()
    
    // Execute profiled operations
    if operations, ok := params["operations"].([]func(*op.Ops) error); ok {
        for _, op := range operations {
            if err := op(ops); err != nil {
                return err
            }
        }
    }
    
    // Log timing
    duration := time.Since(start)
    log.Printf("Operations took: %v", duration)
    
    return nil
})
```

#### 3. Dynamic Content Loading

```go
contentHandler := NewCallbackCommandHandler("dynamic_content", func(ops *op.Ops, cmd clay.CustomCommand) error {
    data := cmd.CustomData.(CustomCommandData)
    
    // Load content based on parameters
    contentType := data.Parameters["type"].(string)
    contentID := data.Parameters["id"].(string)
    
    switch contentType {
    case "image":
        return loadAndRenderImage(ops, contentID)
    case "text":
        return loadAndRenderText(ops, contentID)
    default:
        return fmt.Errorf("unsupported content type: %s", contentType)
    }
})
```

## Built-in Example Handlers

The system provides several example handlers that demonstrate common patterns:

### 1. Debug Overlay Handler

```go
handler := CreateDebugOverlayHandler()
renderer.RegisterCustomHandler(handler)

// Usage
cmd := clay.CustomCommand{
    CustomData: "debug_overlay",
}
```

### 2. Performance Profiler Handler

```go
handler := CreatePerformanceProfilerHandler()
renderer.RegisterCustomHandler(handler)

// Usage
cmd := clay.CustomCommand{
    CustomData: CustomCommandData{
        CommandID: "performance_profiler",
        Parameters: map[string]interface{}{
            "enabled": true,
            "level":   "detailed",
        },
    },
}
```

### 3. Custom Shape Handler

```go
handler := CreateCustomShapeHandler()
renderer.RegisterCustomHandler(handler)

// Usage
cmd := clay.CustomCommand{
    CustomData: map[string]interface{}{
        "commandID": "custom_shape",
        "type":      "star",
        "points":    5,
        "radius":    50,
    },
}
```

## Error Handling

The custom command system provides comprehensive error handling:

### Error Types

- **ErrorTypeInvalidInput**: Invalid handler or command data
- **ErrorTypeResourceNotFound**: Handler not found for command ID
- **ErrorTypeRenderingFailed**: Handler execution failed
- **ErrorTypeUnsupportedOperation**: Operation not supported (e.g., plugins)

### Error Handling Patterns

```go
err := renderer.RenderCustom(cmd)
if err != nil {
    if renderErr, ok := err.(*RenderError); ok {
        switch renderErr.Type {
        case ErrorTypeResourceNotFound:
            log.Printf("Handler not found: %s", renderErr.Message)
        case ErrorTypeRenderingFailed:
            log.Printf("Rendering failed: %s", renderErr.Message)
        default:
            log.Printf("Render error: %s", renderErr.Message)
        }
    }
}
```

## Registry Management

### Handler Registration

```go
// Register handler
err := renderer.RegisterCustomHandler(handler)
if err != nil {
    // Handle registration error
}

// List registered handlers
handlers := renderer.customRegistry.ListHandlers()
fmt.Printf("Registered handlers: %v", handlers)

// Get specific handler
handler, err := renderer.customRegistry.GetHandler("my_command")
if err != nil {
    // Handler not found
}

// Unregister handler
err = renderer.customRegistry.UnregisterHandler("my_command")
```

### Thread Safety

The `CustomCommandRegistry` is thread-safe and can be used concurrently:

```go
// Safe to call from multiple goroutines
go func() {
    renderer.RegisterCustomHandler(handler1)
}()

go func() {
    renderer.RegisterCustomHandler(handler2)
}()
```

## Best Practices

### 1. Command ID Naming

Use descriptive, namespaced command IDs:

```go
// Good
"app.ui.debug_overlay"
"game.effects.particle_system"
"chart.renderer.custom_axis"

// Avoid
"debug"
"custom"
"handler1"
```

### 2. Parameter Validation

Always validate parameters in your handlers:

```go
handler := NewOperationCommandHandler("validated_command", func(ops *op.Ops, params map[string]interface{}) error {
    // Validate required parameters
    color, ok := params["color"].(string)
    if !ok {
        return fmt.Errorf("missing or invalid 'color' parameter")
    }
    
    opacity, ok := params["opacity"].(float64)
    if !ok || opacity < 0 || opacity > 1 {
        return fmt.Errorf("invalid 'opacity' parameter: must be between 0 and 1")
    }
    
    // Proceed with validated parameters
    return nil
})
```

### 3. Resource Management

Clean up resources in your handlers:

```go
handler := NewCallbackCommandHandler("resource_handler", func(ops *op.Ops, cmd clay.CustomCommand) error {
    // Acquire resources
    resource := acquireResource()
    defer resource.Release() // Always clean up
    
    // Use resource for rendering
    return resource.Render(ops)
})
```

### 4. Error Context

Provide meaningful error messages with context:

```go
handler := NewCallbackCommandHandler("contextual_handler", func(ops *op.Ops, cmd clay.CustomCommand) error {
    data, ok := cmd.CustomData.(map[string]interface{})
    if !ok {
        return fmt.Errorf("contextual_handler: expected map data, got %T", cmd.CustomData)
    }
    
    if value, exists := data["required_param"]; !exists {
        return fmt.Errorf("contextual_handler: missing required parameter 'required_param'")
    } else if _, ok := value.(string); !ok {
        return fmt.Errorf("contextual_handler: parameter 'required_param' must be string, got %T", value)
    }
    
    return nil
})
```

## Integration with Clay Layout

Custom commands integrate seamlessly with Clay's layout system:

```go
// In your Clay layout code
func buildLayout() []clay.RenderCommand {
    commands := []clay.RenderCommand{
        // Regular Clay commands
        clay.RectangleCommand{...},
        clay.TextCommand{...},
        
        // Custom commands
        clay.CustomCommand{
            CustomData: "debug_overlay",
        },
        
        clay.CustomCommand{
            CustomData: map[string]interface{}{
                "commandID": "performance_marker",
                "label":     "Layout End",
            },
        },
    }
    
    return commands
}
```

## Future Enhancements

### Plugin System (Planned)

The plugin system will allow loading external custom command handlers:

```go
// Future plugin support
handler := NewPluginCommandHandler("external_renderer", "/path/to/plugin.so")
renderer.RegisterCustomHandler(handler)
```

### Command Composition

Future versions may support command composition and chaining:

```go
// Planned: Command composition
compositeCmd := clay.CustomCommand{
    CustomData: CompositeCommandData{
        Commands: []CustomCommandData{
            {CommandID: "setup_context", ...},
            {CommandID: "render_content", ...},
            {CommandID: "cleanup_context", ...},
        },
    },
}
```

This documentation provides a comprehensive guide to using the custom command system in the Gio renderer for Clay. The system is designed to be extensible, type-safe, and easy to use while providing powerful customization capabilities.
