package gioui

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

// ErrorType represents different categories of rendering errors
type ErrorType int

const (
	ErrorTypeUnknown ErrorType = iota
	ErrorTypeInvalidInput
	ErrorTypeInvalidData
	ErrorTypeInvalidState
	ErrorTypeResourceNotFound
	ErrorTypeRenderingFailed
	ErrorTypeMemoryExhausted
	ErrorTypeUnsupportedOperation
	ErrorTypeClipStackOverflow
	ErrorTypeCoordinateTransform
)

// String returns a string representation of the error type
func (et ErrorType) String() string {
	switch et {
	case ErrorTypeInvalidInput:
		return "InvalidInput"
	case ErrorTypeInvalidData:
		return "InvalidData"
	case ErrorTypeInvalidState:
		return "InvalidState"
	case ErrorTypeResourceNotFound:
		return "ResourceNotFound"
	case ErrorTypeRenderingFailed:
		return "RenderingFailed"
	case ErrorTypeMemoryExhausted:
		return "MemoryExhausted"
	case ErrorTypeUnsupportedOperation:
		return "UnsupportedOperation"
	case ErrorTypeClipStackOverflow:
		return "ClipStackOverflow"
	case ErrorTypeCoordinateTransform:
		return "CoordinateTransform"
	default:
		return "Unknown"
	}
}

// RenderError represents a detailed rendering error with context
type RenderError struct {
	Type        ErrorType
	Operation   string
	Message     string
	Context     map[string]interface{}
	Cause       error
	StackTrace  string
}

// Error implements the error interface
func (re *RenderError) Error() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("[%s]", re.Type))
	
	if re.Operation != "" {
		parts = append(parts, fmt.Sprintf("Operation: %s", re.Operation))
	}
	
	parts = append(parts, re.Message)
	
	if re.Cause != nil {
		parts = append(parts, fmt.Sprintf("Cause: %v", re.Cause))
	}
	
	return strings.Join(parts, " - ")
}

// Unwrap returns the underlying cause error
func (re *RenderError) Unwrap() error {
	return re.Cause
}

// NewRenderError creates a new render error with context
func NewRenderError(errorType ErrorType, operation, message string) *RenderError {
	return &RenderError{
		Type:       errorType,
		Operation:  operation,
		Message:    message,
		Context:    make(map[string]interface{}),
		StackTrace: getStackTrace(),
	}
}

// NewRenderErrorWithCause creates a new render error wrapping another error
func NewRenderErrorWithCause(errorType ErrorType, operation, message string, cause error) *RenderError {
	return &RenderError{
		Type:       errorType,
		Operation:  operation,
		Message:    message,
		Context:    make(map[string]interface{}),
		Cause:      cause,
		StackTrace: getStackTrace(),
	}
}

// WithContext adds context information to the error
func (re *RenderError) WithContext(key string, value interface{}) *RenderError {
	re.Context[key] = value
	return re
}

// GetContext retrieves context information from the error
func (re *RenderError) GetContext(key string) (interface{}, bool) {
	value, exists := re.Context[key]
	return value, exists
}

// getStackTrace captures the current stack trace
func getStackTrace() string {
	buf := make([]byte, 1024)
	for {
		n := runtime.Stack(buf, false)
		if n < len(buf) {
			return string(buf[:n])
		}
		buf = make([]byte, 2*len(buf))
	}
}

// ErrorHandler manages error logging and recovery
type ErrorHandler struct {
	logger        *log.Logger
	debugMode     bool
	errorCallback func(*RenderError)
}

// NewErrorHandler creates a new error handler
func NewErrorHandler(logger *log.Logger, debugMode bool) *ErrorHandler {
	return &ErrorHandler{
		logger:    logger,
		debugMode: debugMode,
	}
}

// SetErrorCallback sets a callback function for error notifications
func (eh *ErrorHandler) SetErrorCallback(callback func(*RenderError)) {
	eh.errorCallback = callback
}

// HandleError processes and logs a render error
func (eh *ErrorHandler) HandleError(err *RenderError) {
	if eh.logger != nil {
		if eh.debugMode {
			eh.logger.Printf("Render Error: %s\nContext: %+v\nStack Trace:\n%s",
				err.Error(), err.Context, err.StackTrace)
		} else {
			eh.logger.Printf("Render Error: %s", err.Error())
		}
	}
	
	if eh.errorCallback != nil {
		eh.errorCallback(err)
	}
}

// RecoverFromPanic recovers from panics and converts them to render errors
func (eh *ErrorHandler) RecoverFromPanic(operation string) *RenderError {
	if r := recover(); r != nil {
		var err error
		switch v := r.(type) {
		case error:
			err = v
		case string:
			err = fmt.Errorf("%s", v)
		default:
			err = fmt.Errorf("panic: %v", v)
		}
		
		renderErr := NewRenderErrorWithCause(
			ErrorTypeRenderingFailed,
			operation,
			"Panic occurred during rendering",
			err,
		)
		
		eh.HandleError(renderErr)
		return renderErr
	}
	return nil
}

// Validation functions for common error conditions

// ValidateImageData validates image data before processing
func ValidateImageData(data interface{}) error {
	if data == nil {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"ValidateImageData",
			"Image data cannot be nil",
		)
	}
	
	// Add more specific validation based on data type
	return nil
}

// ValidateBounds validates bounding box coordinates
func ValidateBounds(x, y, width, height int) error {
	if width <= 0 || height <= 0 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"ValidateBounds",
			fmt.Sprintf("Invalid bounds: width=%d, height=%d", width, height),
		).WithContext("x", x).WithContext("y", y).WithContext("width", width).WithContext("height", height)
	}
	return nil
}

// ValidateColor validates color values
func ValidateColor(r, g, b, a float32) error {
	if r < 0 || r > 1 || g < 0 || g > 1 || b < 0 || b > 1 || a < 0 || a > 1 {
		return NewRenderError(
			ErrorTypeInvalidInput,
			"ValidateColor",
			fmt.Sprintf("Color values must be between 0 and 1: r=%.3f, g=%.3f, b=%.3f, a=%.3f", r, g, b, a),
		).WithContext("r", r).WithContext("g", g).WithContext("b", b).WithContext("a", a)
	}
	return nil
}

// ValidateClipStackDepth validates clip stack depth to prevent overflow
func ValidateClipStackDepth(currentDepth, maxDepth int) error {
	if currentDepth >= maxDepth {
		return NewRenderError(
			ErrorTypeClipStackOverflow,
			"ValidateClipStackDepth",
			fmt.Sprintf("Clip stack overflow: current depth %d exceeds maximum %d", currentDepth, maxDepth),
		).WithContext("currentDepth", currentDepth).WithContext("maxDepth", maxDepth)
	}
	return nil
}

// ErrorRecovery provides mechanisms for recovering from errors
type ErrorRecovery struct {
	fallbackOperations map[ErrorType]func() error
}

// NewErrorRecovery creates a new error recovery system
func NewErrorRecovery() *ErrorRecovery {
	return &ErrorRecovery{
		fallbackOperations: make(map[ErrorType]func() error),
	}
}

// RegisterFallback registers a fallback operation for a specific error type
func (er *ErrorRecovery) RegisterFallback(errorType ErrorType, fallback func() error) {
	er.fallbackOperations[errorType] = fallback
}

// AttemptRecovery attempts to recover from an error using registered fallbacks
func (er *ErrorRecovery) AttemptRecovery(renderErr *RenderError) error {
	if fallback, exists := er.fallbackOperations[renderErr.Type]; exists {
		return fallback()
	}
	return renderErr // No recovery available
}

// Common error creation helpers

// ErrInvalidImageData creates an invalid image data error
func ErrInvalidImageData(operation string, data interface{}) *RenderError {
	return NewRenderError(
		ErrorTypeInvalidInput,
		operation,
		fmt.Sprintf("Invalid image data type: %T", data),
	).WithContext("dataType", fmt.Sprintf("%T", data))
}

// ErrUnsupportedImageFormat creates an unsupported image format error
func ErrUnsupportedImageFormat(operation, format string) *RenderError {
	return NewRenderError(
		ErrorTypeUnsupportedOperation,
		operation,
		fmt.Sprintf("Unsupported image format: %s", format),
	).WithContext("format", format)
}

// ErrMemoryExhausted creates a memory exhausted error
func ErrMemoryExhausted(operation string, requested, available int) *RenderError {
	return NewRenderError(
		ErrorTypeMemoryExhausted,
		operation,
		fmt.Sprintf("Memory exhausted: requested %d bytes, available %d bytes", requested, available),
	).WithContext("requested", requested).WithContext("available", available)
}

// ErrCoordinateTransform creates a coordinate transformation error
func ErrCoordinateTransform(operation string, cause error) *RenderError {
	return NewRenderErrorWithCause(
		ErrorTypeCoordinateTransform,
		operation,
		"Coordinate transformation failed",
		cause,
	)
}
