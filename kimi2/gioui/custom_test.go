package gioui

import (
	"testing"

	"gioui.org/op"
	"github.com/zodimo/go-clay/kimi2/clay"
)

func TestNewCustomCommandRegistry(t *testing.T) {
	registry := NewCustomCommandRegistry()

	if registry == nil {
		t.Fatal("Expected registry to be created, got nil")
	}

	if registry.handlers == nil {
		t.Fatal("Expected handlers map to be initialized")
	}

	if len(registry.handlers) != 0 {
		t.Errorf("Expected empty handlers map, got %d handlers", len(registry.handlers))
	}
}

func TestCustomCommandRegistry_RegisterHandler(t *testing.T) {
	registry := NewCustomCommandRegistry()

	t.Run("Valid handler", func(t *testing.T) {
		handler := NewCallbackCommandHandler("test_handler", func(ops *op.Ops, cmd clay.CustomCommand) error {
			return nil
		})

		err := registry.RegisterHandler(handler)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify handler was registered
		handlers := registry.ListHandlers()
		if len(handlers) != 1 || handlers[0] != "test_handler" {
			t.Errorf("Expected handler to be registered with ID 'test_handler', got: %v", handlers)
		}
	})

	t.Run("Nil handler", func(t *testing.T) {
		err := registry.RegisterHandler(nil)
		if err == nil {
			t.Error("Expected error for nil handler")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	})

	t.Run("Duplicate handler ID", func(t *testing.T) {
		handler1 := NewCallbackCommandHandler("duplicate_id", func(ops *op.Ops, cmd clay.CustomCommand) error {
			return nil
		})
		handler2 := NewCallbackCommandHandler("duplicate_id", func(ops *op.Ops, cmd clay.CustomCommand) error {
			return nil
		})

		// Register first handler
		err := registry.RegisterHandler(handler1)
		if err != nil {
			t.Errorf("Expected no error for first handler, got: %v", err)
		}

		// Try to register duplicate
		err = registry.RegisterHandler(handler2)
		if err == nil {
			t.Error("Expected error for duplicate handler ID")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	})
}

func TestCustomCommandRegistry_UnregisterHandler(t *testing.T) {
	registry := NewCustomCommandRegistry()

	t.Run("Existing handler", func(t *testing.T) {
		handler := NewCallbackCommandHandler("test_handler", func(ops *op.Ops, cmd clay.CustomCommand) error {
			return nil
		})

		// Register handler first
		registry.RegisterHandler(handler)

		// Unregister it
		err := registry.UnregisterHandler("test_handler")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify handler was removed
		handlers := registry.ListHandlers()
		if len(handlers) != 0 {
			t.Errorf("Expected no handlers after unregistering, got: %v", handlers)
		}
	})

	t.Run("Non-existing handler", func(t *testing.T) {
		err := registry.UnregisterHandler("non_existing")
		if err == nil {
			t.Error("Expected error for non-existing handler")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeResourceNotFound {
			t.Errorf("Expected ErrorTypeResourceNotFound, got %v", renderErr.Type)
		}
	})
}

func TestCustomCommandRegistry_GetHandler(t *testing.T) {
	registry := NewCustomCommandRegistry()

	t.Run("Existing handler", func(t *testing.T) {
		handler := NewCallbackCommandHandler("test_handler", func(ops *op.Ops, cmd clay.CustomCommand) error {
			return nil
		})

		registry.RegisterHandler(handler)

		retrieved, err := registry.GetHandler("test_handler")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if retrieved != handler {
			t.Error("Expected to retrieve the same handler instance")
		}
	})

	t.Run("Non-existing handler", func(t *testing.T) {
		_, err := registry.GetHandler("non_existing")
		if err == nil {
			t.Error("Expected error for non-existing handler")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeResourceNotFound {
			t.Errorf("Expected ErrorTypeResourceNotFound, got %v", renderErr.Type)
		}
	})
}

func TestCustomCommandRegistry_ExecuteCustomCommand(t *testing.T) {
	registry := NewCustomCommandRegistry()
	ops := new(op.Ops)

	t.Run("Valid command with string ID", func(t *testing.T) {
		executed := false
		handler := NewCallbackCommandHandler("test_command", func(ops *op.Ops, cmd clay.CustomCommand) error {
			executed = true
			return nil
		})

		registry.RegisterHandler(handler)

		cmd := clay.CustomCommand{
			CustomData: "test_command",
		}

		err := registry.ExecuteCustomCommand(ops, cmd)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !executed {
			t.Error("Expected handler to be executed")
		}
	})

	t.Run("Valid command with map data", func(t *testing.T) {
		executed := false
		handler := NewCallbackCommandHandler("map_command", func(ops *op.Ops, cmd clay.CustomCommand) error {
			executed = true
			return nil
		})

		registry.RegisterHandler(handler)

		cmd := clay.CustomCommand{
			CustomData: map[string]interface{}{
				"commandID": "map_command",
				"param1":    "value1",
			},
		}

		err := registry.ExecuteCustomCommand(ops, cmd)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !executed {
			t.Error("Expected handler to be executed")
		}
	})

	t.Run("Valid command with CustomCommandData", func(t *testing.T) {
		executed := false
		handler := NewCallbackCommandHandler("struct_command", func(ops *op.Ops, cmd clay.CustomCommand) error {
			executed = true
			return nil
		})

		registry.RegisterHandler(handler)

		cmd := clay.CustomCommand{
			CustomData: CustomCommandData{
				CommandID: "struct_command",
				Parameters: map[string]interface{}{
					"param1": "value1",
				},
			},
		}

		err := registry.ExecuteCustomCommand(ops, cmd)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !executed {
			t.Error("Expected handler to be executed")
		}
	})

	t.Run("Invalid custom data", func(t *testing.T) {
		cmd := clay.CustomCommand{
			CustomData: 12345, // Invalid type
		}

		err := registry.ExecuteCustomCommand(ops, cmd)
		if err == nil {
			t.Error("Expected error for invalid custom data")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	})

	t.Run("Non-existing handler", func(t *testing.T) {
		cmd := clay.CustomCommand{
			CustomData: "non_existing_command",
		}

		err := registry.ExecuteCustomCommand(ops, cmd)
		if err == nil {
			t.Error("Expected error for non-existing handler")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeResourceNotFound {
			t.Errorf("Expected ErrorTypeResourceNotFound, got %v", renderErr.Type)
		}
	})
}

func TestCallbackCommandHandler(t *testing.T) {
	t.Run("Valid callback", func(t *testing.T) {
		executed := false
		handler := NewCallbackCommandHandler("test_callback", func(ops *op.Ops, cmd clay.CustomCommand) error {
			executed = true
			return nil
		})

		if handler.GetCommandID() != "test_callback" {
			t.Errorf("Expected command ID 'test_callback', got '%s'", handler.GetCommandID())
		}

		if handler.GetCommandType() != CustomCommandTypeCallback {
			t.Errorf("Expected command type %v, got %v", CustomCommandTypeCallback, handler.GetCommandType())
		}

		ops := new(op.Ops)
		cmd := clay.CustomCommand{}

		err := handler.HandleCustomCommand(ops, cmd)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !executed {
			t.Error("Expected callback to be executed")
		}
	})

	t.Run("Nil callback", func(t *testing.T) {
		handler := &CallbackCommandHandler{
			commandID: "nil_callback",
			callback:  nil,
		}

		ops := new(op.Ops)
		cmd := clay.CustomCommand{}

		err := handler.HandleCustomCommand(ops, cmd)
		if err == nil {
			t.Error("Expected error for nil callback")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeRenderingFailed {
			t.Errorf("Expected ErrorTypeRenderingFailed, got %v", renderErr.Type)
		}
	})
}

func TestOperationCommandHandler(t *testing.T) {
	t.Run("Valid operation", func(t *testing.T) {
		executed := false
		handler := NewOperationCommandHandler("test_operation", func(ops *op.Ops, params map[string]interface{}) error {
			executed = true
			return nil
		})

		if handler.GetCommandID() != "test_operation" {
			t.Errorf("Expected command ID 'test_operation', got '%s'", handler.GetCommandID())
		}

		if handler.GetCommandType() != CustomCommandTypeOperation {
			t.Errorf("Expected command type %v, got %v", CustomCommandTypeOperation, handler.GetCommandType())
		}

		ops := new(op.Ops)
		cmd := clay.CustomCommand{
			CustomData: map[string]interface{}{
				"param1": "value1",
			},
		}

		err := handler.HandleCustomCommand(ops, cmd)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !executed {
			t.Error("Expected operation to be executed")
		}
	})

	t.Run("Nil operation", func(t *testing.T) {
		handler := &OperationCommandHandler{
			commandID: "nil_operation",
			operation: nil,
		}

		ops := new(op.Ops)
		cmd := clay.CustomCommand{}

		err := handler.HandleCustomCommand(ops, cmd)
		if err == nil {
			t.Error("Expected error for nil operation")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeRenderingFailed {
			t.Errorf("Expected ErrorTypeRenderingFailed, got %v", renderErr.Type)
		}
	})
}

func TestPluginCommandHandler(t *testing.T) {
	t.Run("Plugin loading not implemented", func(t *testing.T) {
		handler := NewPluginCommandHandler("test_plugin", "/path/to/plugin")

		if handler.GetCommandID() != "test_plugin" {
			t.Errorf("Expected command ID 'test_plugin', got '%s'", handler.GetCommandID())
		}

		if handler.GetCommandType() != CustomCommandTypePlugin {
			t.Errorf("Expected command type %v, got %v", CustomCommandTypePlugin, handler.GetCommandType())
		}

		ops := new(op.Ops)
		cmd := clay.CustomCommand{}

		err := handler.HandleCustomCommand(ops, cmd)
		if err == nil {
			t.Error("Expected error for unimplemented plugin loading")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeUnsupportedOperation {
			t.Errorf("Expected ErrorTypeUnsupportedOperation, got %v", renderErr.Type)
		}
	})
}

func TestGioRenderer_RenderCustom(t *testing.T) {
	t.Run("Valid custom command", func(t *testing.T) {
		ops := new(op.Ops)
		renderer := NewRenderer(ops)

		// Register a test handler
		executed := false
		handler := NewCallbackCommandHandler("test_render", func(ops *op.Ops, cmd clay.CustomCommand) error {
			executed = true
			return nil
		})

		err := renderer.RegisterCustomHandler(handler)
		if err != nil {
			t.Fatalf("Failed to register handler: %v", err)
		}

		// Execute custom command
		cmd := clay.CustomCommand{
			CustomData: "test_render",
		}

		err = renderer.RenderCustom(cmd)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		if !executed {
			t.Error("Expected custom handler to be executed")
		}
	})

	t.Run("Nil custom registry", func(t *testing.T) {
		renderer := &GioRenderer{
			customRegistry: nil,
		}

		cmd := clay.CustomCommand{
			CustomData: "test_command",
		}

		err := renderer.RenderCustom(cmd)
		if err == nil {
			t.Error("Expected error for nil custom registry")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeInvalidInput {
			t.Errorf("Expected ErrorTypeInvalidInput, got %v", renderErr.Type)
		}
	})

	t.Run("Non-existing handler", func(t *testing.T) {
		ops := new(op.Ops)
		renderer := NewRenderer(ops)

		cmd := clay.CustomCommand{
			CustomData: "non_existing_handler",
		}

		err := renderer.RenderCustom(cmd)
		if err == nil {
			t.Error("Expected error for non-existing handler")
		}

		renderErr, ok := err.(*RenderError)
		if !ok {
			t.Errorf("Expected RenderError, got %T", err)
		} else if renderErr.Type != ErrorTypeResourceNotFound {
			t.Errorf("Expected ErrorTypeResourceNotFound, got %v", renderErr.Type)
		}
	})
}

func TestGioRenderer_RegisterCustomHandler(t *testing.T) {
	ops := new(op.Ops)
	renderer := NewRenderer(ops)

	t.Run("Valid handler registration", func(t *testing.T) {
		handler := NewCallbackCommandHandler("test_handler", func(ops *op.Ops, cmd clay.CustomCommand) error {
			return nil
		})

		err := renderer.RegisterCustomHandler(handler)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}

		// Verify handler was registered
		handlers := renderer.customRegistry.ListHandlers()
		if len(handlers) != 1 || handlers[0] != "test_handler" {
			t.Errorf("Expected handler to be registered, got: %v", handlers)
		}
	})
}

func TestExampleCustomHandlers(t *testing.T) {
	t.Run("Debug overlay handler", func(t *testing.T) {
		handler := CreateDebugOverlayHandler()

		if handler.GetCommandID() != "debug_overlay" {
			t.Errorf("Expected command ID 'debug_overlay', got '%s'", handler.GetCommandID())
		}

		if handler.GetCommandType() != CustomCommandTypeCallback {
			t.Errorf("Expected command type %v, got %v", CustomCommandTypeCallback, handler.GetCommandType())
		}

		// Test execution (should not error)
		ops := new(op.Ops)
		cmd := clay.CustomCommand{}

		err := handler.HandleCustomCommand(ops, cmd)
		if err != nil {
			t.Errorf("Expected no error from debug overlay handler, got: %v", err)
		}
	})

	t.Run("Performance profiler handler", func(t *testing.T) {
		handler := CreatePerformanceProfilerHandler()

		if handler.GetCommandID() != "performance_profiler" {
			t.Errorf("Expected command ID 'performance_profiler', got '%s'", handler.GetCommandID())
		}

		if handler.GetCommandType() != CustomCommandTypeOperation {
			t.Errorf("Expected command type %v, got %v", CustomCommandTypeOperation, handler.GetCommandType())
		}
	})

	t.Run("Custom shape handler", func(t *testing.T) {
		handler := CreateCustomShapeHandler()

		if handler.GetCommandID() != "custom_shape" {
			t.Errorf("Expected command ID 'custom_shape', got '%s'", handler.GetCommandID())
		}

		if handler.GetCommandType() != CustomCommandTypeCallback {
			t.Errorf("Expected command type %v, got %v", CustomCommandTypeCallback, handler.GetCommandType())
		}

		// Test execution (should not error)
		ops := new(op.Ops)
		cmd := clay.CustomCommand{}

		err := handler.HandleCustomCommand(ops, cmd)
		if err != nil {
			t.Errorf("Expected no error from custom shape handler, got: %v", err)
		}
	})
}
