package core

import (
	"fmt"
)

// Error types for better error categorization and handling.
const (
	ErrTypeUnsupportedFormat = "unsupported_format"
	ErrTypeInvalidInput      = "invalid_input"
	ErrTypeParsingFailed     = "parsing_failed"
	ErrTypeFileNotFound      = "file_not_found"
	ErrTypeValidationFailed  = "validation_failed"
	ErrTypeConversionFailed  = "conversion_failed"
	ErrTypeCombinationFailed = "combination_failed"
	ErrTypeExportFailed      = "export_failed"
)

// TconvError represents a structured error with context and categorization.
type TconvError struct {
	Type        string            // Error category
	Message     string            // Human-readable error message
	Cause       error             // Underlying error if any
	Context     map[string]string // Additional context information
	Suggestions []string          // Helpful suggestions for resolution
}

// Error implements the error interface.
func (e *TconvError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Type, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error for error chain compatibility.
func (e *TconvError) Unwrap() error {
	return e.Cause
}

// WithContext adds context information to the error.
func (e *TconvError) WithContext(key, value string) *TconvError {
	if e.Context == nil {
		e.Context = make(map[string]string)
	}
	e.Context[key] = value
	return e
}

// WithSuggestion adds a helpful suggestion to the error.
func (e *TconvError) WithSuggestion(suggestion string) *TconvError {
	e.Suggestions = append(e.Suggestions, suggestion)
	return e
}

// NewTconvError creates a new structured error.
func NewTconvError(errorType, message string) *TconvError {
	return &TconvError{
		Type:        errorType,
		Message:     message,
		Context:     make(map[string]string),
		Suggestions: make([]string, 0),
	}
}

// NewTconvErrorWithCause creates a new structured error with an underlying cause.
func NewTconvErrorWithCause(errorType, message string, cause error) *TconvError {
	return &TconvError{
		Type:        errorType,
		Message:     message,
		Cause:       cause,
		Context:     make(map[string]string),
		Suggestions: make([]string, 0),
	}
}

// Common error constructors for frequently used error types.

// NewUnsupportedFormatError creates an error for unsupported file formats.
func NewUnsupportedFormatError(format string, supportedFormats []string) *TconvError {
	err := NewTconvError(ErrTypeUnsupportedFormat,
		fmt.Sprintf("unsupported format: %s", format))
	err.WithContext("requested_format", format)
	if len(supportedFormats) > 0 {
		err.WithContext("supported_formats", fmt.Sprintf("%v", supportedFormats))
		err.WithSuggestion(fmt.Sprintf("Use one of the supported formats: %v", supportedFormats))
	}
	return err
}

// NewInvalidInputError creates an error for invalid input data.
func NewInvalidInputError(message string) *TconvError {
	return NewTconvError(ErrTypeInvalidInput, message).
		WithSuggestion("Verify that the input file is valid and properly formatted")
}

// NewParsingError creates an error for parsing failures.
func NewParsingError(format string, cause error) *TconvError {
	err := NewTconvErrorWithCause(ErrTypeParsingFailed,
		fmt.Sprintf("failed to parse %s format", format), cause)
	err.WithContext("format", format)
	err.WithSuggestion("Check if the file is corrupted or in the wrong format")
	return err
}

// NewFileNotFoundError creates an error for missing files.
func NewFileNotFoundError(filePath string) *TconvError {
	return NewTconvError(ErrTypeFileNotFound,
		fmt.Sprintf("file not found: %s", filePath)).
		WithContext("file_path", filePath).
		WithSuggestion("Verify that the file path is correct and the file exists")
}

// NewValidationError creates an error for validation failures.
func NewValidationError(message string, cause error) *TconvError {
	err := NewTconvErrorWithCause(ErrTypeValidationFailed, message, cause)
	err.WithSuggestion("Check the input data format and structure")
	return err
}

// NewConversionError creates an error for conversion failures.
func NewConversionError(sourceFormat, targetFormat string, cause error) *TconvError {
	err := NewTconvErrorWithCause(ErrTypeConversionFailed,
		fmt.Sprintf("failed to convert from %s to %s", sourceFormat, targetFormat), cause)
	err.WithContext("source_format", sourceFormat)
	err.WithContext("target_format", targetFormat)
	return err
}

// NewCombinationError creates an error for result combination failures.
func NewCombinationError(reason string) *TconvError {
	return NewTconvError(ErrTypeCombinationFailed,
		fmt.Sprintf("failed to combine results: %s", reason)).
		WithSuggestion("Ensure all input results are valid and compatible")
}

// NewExportError creates an error for export failures.
func NewExportError(format string, cause error) *TconvError {
	err := NewTconvErrorWithCause(ErrTypeExportFailed,
		fmt.Sprintf("failed to export to %s format", format), cause)
	err.WithContext("export_format", format)
	return err
}