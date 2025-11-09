package gioui

import (
	"fmt"
	"sync"

	"gioui.org/op"
	"github.com/zodimo/go-clay/kimi2/clay"
)

// CustomCommandType represents different types of custom commands
type CustomCommandType string

const (
	CustomCommandTypePlugin    CustomCommandType = "plugin"
	CustomCommandTypeCallback  CustomCommandType = "callback"
	CustomCommandTypeOperation CustomCommandType = "operation"
)

// CustomCommandHandler defines the interface for handling custom commands
type CustomCommandHandler interface {
	// HandleCustomCommand processes a custom command and returns an error if failed
	HandleCustomCommand(ops *op.Ops, cmd clay.RenderCommand) error

	// GetCommandType returns the type of custom command this handler supports
	GetCommandType() CustomCommandType

	// GetCommandID returns a unique identifier for this command handler
	GetCommandID() string
}

// CustomCommandRegistry manages registered custom command handlers
type CustomCommandRegistry struct {
	handlers map[string]CustomCommandHandler
	mutex    sync.RWMutex
}

// NewCustomCommandRegistry creates a new custom command registry
func NewCustomCommandRegistry() *CustomCommandRegistry {
	return &CustomCommandRegistry{
		handlers: make(map[string]CustomCommandHandler),
	}
}

// RegisterHandler registers a custom command handler
func (ccr *CustomCommandRegistry) RegisterHandler(handler CustomCommandHandler) error {
	if handler == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"RegisterHandler",
			"Handler cannot be nil",
		)
	}

	commandID := handler.GetCommandID()
	if commandID == "" {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"RegisterHandler",
			"Handler must have a non-empty command ID",
		)
	}

	ccr.mutex.Lock()
	defer ccr.mutex.Unlock()

	if _, exists := ccr.handlers[commandID]; exists {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"RegisterHandler",
			fmt.Sprintf("Handler with ID '%s' already registered", commandID),
		).WithContext("commandID", commandID)
	}

	ccr.handlers[commandID] = handler
	return nil
}

// UnregisterHandler removes a custom command handler
func (ccr *CustomCommandRegistry) UnregisterHandler(commandID string) error {
	ccr.mutex.Lock()
	defer ccr.mutex.Unlock()

	if _, exists := ccr.handlers[commandID]; !exists {
		return NewRenderError(
			ErrorTypeResourceNotFound,
			"UnregisterHandler",
			fmt.Sprintf("Handler with ID '%s' not found", commandID),
		).WithContext("commandID", commandID)
	}

	delete(ccr.handlers, commandID)
	return nil
}

// GetHandler retrieves a custom command handler by ID
func (ccr *CustomCommandRegistry) GetHandler(commandID string) (CustomCommandHandler, error) {
	ccr.mutex.RLock()
	defer ccr.mutex.RUnlock()

	handler, exists := ccr.handlers[commandID]
	if !exists {
		return nil, NewRenderError(
			ErrorTypeResourceNotFound,
			"GetHandler",
			fmt.Sprintf("Handler with ID '%s' not found", commandID),
		).WithContext("commandID", commandID)
	}

	return handler, nil
}

// ListHandlers returns a list of all registered handler IDs
func (ccr *CustomCommandRegistry) ListHandlers() []string {
	ccr.mutex.RLock()
	defer ccr.mutex.RUnlock()

	ids := make([]string, 0, len(ccr.handlers))
	for id := range ccr.handlers {
		ids = append(ids, id)
	}
	return ids
}

// ExecuteCustomCommand executes a custom command using the appropriate handler
func (ccr *CustomCommandRegistry) ExecuteCustomCommand(ops *op.Ops, cmd clay.RenderCommand) error {
	customData, ok := cmd.Data.(CustomCommandData)
	if !ok {
		return NewRenderError(
			ErrorTypeInvalidData,
			"ExecuteCustomCommand",
			"Invalid custom command data",
		)
	}
	// Extract command ID from custom data
	commandID, err := ccr.extractCommandID(customData)
	if err != nil {
		return err
	}

	handler, err := ccr.GetHandler(commandID)
	if err != nil {
		return err
	}

	return handler.HandleCustomCommand(ops, cmd)
}

// extractCommandID extracts the command ID from custom data
func (ccr *CustomCommandRegistry) extractCommandID(customData interface{}) (string, error) {
	switch data := customData.(type) {
	case map[string]interface{}:
		if id, exists := data["commandID"]; exists {
			if idStr, ok := id.(string); ok {
				return idStr, nil
			}
		}
	case CustomCommandData:
		return data.CommandID, nil
	case string:
		// Assume the string itself is the command ID
		return data, nil
	}

	return "", NewRenderError(
		ErrorTypeInvalidInput,
		"extractCommandID",
		"Unable to extract command ID from custom data",
	).WithContext("customDataType", fmt.Sprintf("%T", customData))
}

// CustomCommandData represents structured custom command data
type CustomCommandData struct {
	CommandID  string                 `json:"commandID"`
	Parameters map[string]interface{} `json:"parameters"`
	Metadata   map[string]string      `json:"metadata"`
}

// CallbackCommandHandler implements custom commands using callback functions
type CallbackCommandHandler struct {
	commandID string
	callback  func(*op.Ops, clay.RenderCommand) error
}

// NewCallbackCommandHandler creates a new callback-based custom command handler
func NewCallbackCommandHandler(commandID string, callback func(*op.Ops, clay.RenderCommand) error) *CallbackCommandHandler {
	return &CallbackCommandHandler{
		commandID: commandID,
		callback:  callback,
	}
}

// HandleCustomCommand implements CustomCommandHandler
func (cch *CallbackCommandHandler) HandleCustomCommand(ops *op.Ops, cmd clay.RenderCommand) error {
	if cch.callback == nil {
		return NewRenderError(
			ErrorTypeRenderingFailed,
			"HandleCustomCommand",
			"Callback function is nil",
		)
	}

	return cch.callback(ops, cmd)
}

// GetCommandType implements CustomCommandHandler
func (cch *CallbackCommandHandler) GetCommandType() CustomCommandType {
	return CustomCommandTypeCallback
}

// GetCommandID implements CustomCommandHandler
func (cch *CallbackCommandHandler) GetCommandID() string {
	return cch.commandID
}

// OperationCommandHandler implements custom commands using Gio operations
type OperationCommandHandler struct {
	commandID string
	operation func(*op.Ops, map[string]interface{}) error
}

// NewOperationCommandHandler creates a new operation-based custom command handler
func NewOperationCommandHandler(commandID string, operation func(*op.Ops, map[string]interface{}) error) *OperationCommandHandler {
	return &OperationCommandHandler{
		commandID: commandID,
		operation: operation,
	}
}

// HandleCustomCommand implements CustomCommandHandler
func (och *OperationCommandHandler) HandleCustomCommand(ops *op.Ops, cmd clay.RenderCommand) error {
	if och.operation == nil {
		return NewRenderError(
			ErrorTypeRenderingFailed,
			"HandleCustomCommand",
			"Operation function is nil",
		)
	}

	// Extract parameters from custom data
	var params map[string]interface{}
	switch data := cmd.Data.(type) {
	case map[string]interface{}:
		params = data
	case CustomCommandData:
		params = data.Parameters
	default:
		params = make(map[string]interface{})
	}

	return och.operation(ops, params)
}

// GetCommandType implements CustomCommandHandler
func (och *OperationCommandHandler) GetCommandType() CustomCommandType {
	return CustomCommandTypeOperation
}

// GetCommandID implements CustomCommandHandler
func (och *OperationCommandHandler) GetCommandID() string {
	return och.commandID
}

// PluginCommandHandler implements custom commands using external plugins
type PluginCommandHandler struct {
	commandID  string
	pluginPath string
	plugin     CustomPlugin
}

// CustomPlugin defines the interface for external plugins
type CustomPlugin interface {
	// Initialize initializes the plugin
	Initialize() error

	// Execute executes the plugin with given parameters
	Execute(ops *op.Ops, params map[string]interface{}) error

	// Cleanup cleans up plugin resources
	Cleanup() error

	// GetVersion returns the plugin version
	GetVersion() string
}

// NewPluginCommandHandler creates a new plugin-based custom command handler
func NewPluginCommandHandler(commandID, pluginPath string) *PluginCommandHandler {
	return &PluginCommandHandler{
		commandID:  commandID,
		pluginPath: pluginPath,
	}
}

// LoadPlugin loads the plugin (placeholder for actual plugin loading)
func (pch *PluginCommandHandler) LoadPlugin() error {
	// This would implement actual plugin loading from the plugin path
	// For now, return an error indicating plugins are not yet supported
	return NewRenderError(
		ErrorTypeUnsupportedOperation,
		"LoadPlugin",
		"Plugin loading not yet implemented",
	).WithContext("pluginPath", pch.pluginPath)
}

// HandleCustomCommand implements CustomCommandHandler
func (pch *PluginCommandHandler) HandleCustomCommand(ops *op.Ops, cmd clay.RenderCommand) error {
	if pch.plugin == nil {
		if err := pch.LoadPlugin(); err != nil {
			return err
		}
	}

	// Extract parameters from custom data
	var params map[string]interface{}
	switch data := cmd.Data.(type) {
	case map[string]interface{}:
		params = data
	case CustomCommandData:
		params = data.Parameters
	default:
		params = make(map[string]interface{})
	}

	return pch.plugin.Execute(ops, params)
}

// GetCommandType implements CustomCommandHandler
func (pch *PluginCommandHandler) GetCommandType() CustomCommandType {
	return CustomCommandTypePlugin
}

// GetCommandID implements CustomCommandHandler
func (pch *PluginCommandHandler) GetCommandID() string {
	return pch.commandID
}

// Built-in custom command examples

// CreateDebugOverlayHandler creates a debug overlay custom command handler
func CreateDebugOverlayHandler() CustomCommandHandler {
	return NewCallbackCommandHandler("debug_overlay", func(ops *op.Ops, cmd clay.RenderCommand) error {
		// This would implement a debug overlay showing render information
		// For now, it's a placeholder
		return nil
	})
}

// CreatePerformanceProfilerHandler creates a performance profiler custom command handler
func CreatePerformanceProfilerHandler() CustomCommandHandler {
	return NewOperationCommandHandler("performance_profiler", func(ops *op.Ops, params map[string]interface{}) error {
		// This would implement performance profiling operations
		// For now, it's a placeholder
		return nil
	})
}

// CreateCustomShapeHandler creates a custom shape rendering command handler
func CreateCustomShapeHandler() CustomCommandHandler {
	return NewCallbackCommandHandler("custom_shape", func(ops *op.Ops, cmd clay.RenderCommand) error {
		// This would implement custom shape rendering
		// For now, it's a placeholder
		return nil
	})
}
