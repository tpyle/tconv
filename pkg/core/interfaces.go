// Package core defines the core interfaces and types for the tconv library.
// This package provides the fundamental abstractions that make the library
// extensible and suitable for consumption by other programs.
package core

import (
	"io"

	"github.com/tpyle/tconv/pkg/models"
)

// TestConverter defines the interface for converting test files to tiki format.
// Implementations should be thread-safe and stateless.
type TestConverter interface {
	// Convert processes a single test file and returns a tiki result.
	Convert(input io.Reader, sourceType string) (*models.TikiTestResult, error)

	// ConvertFile processes a test file by path and returns a tiki result.
	ConvertFile(filePath, sourceType string) (*models.TikiTestResult, error)

	// SupportedTypes returns a list of supported source types.
	SupportedTypes() []string
}

// TypeDetector defines the interface for automatically detecting test file types.
type TypeDetector interface {
	// DetectType analyzes the content and returns the detected type.
	DetectType(content io.Reader) (string, error)

	// DetectTypeFromFile analyzes a file and returns the detected type.
	DetectTypeFromFile(filePath string) (string, error)

	// SupportedExtensions returns a map of file extensions to likely types.
	SupportedExtensions() map[string]string
}

// ResultCombiner defines the interface for combining multiple test results.
type ResultCombiner interface {
	// Combine merges multiple tiki test results into a single result.
	Combine(results []*models.TikiTestResult, options CombineOptions) (*models.TikiTestResult, error)
}

// TestResultExporter defines the interface for exporting test results.
type TestResultExporter interface {
	// Export writes the tiki test result to the specified output.
	Export(result *models.TikiTestResult, output io.Writer, options ExportOptions) error
}

// Parser defines the interface that all format-specific parsers must implement.
type Parser interface {
	// Parse processes input data and returns a tiki test result.
	Parse(input io.Reader) (*models.TikiTestResult, error)

	// Validate checks if the input is valid for this parser.
	Validate(input io.Reader) error

	// Type returns the type identifier for this parser.
	Type() string
}

// Type aliases for backward compatibility - these reference the canonical types in models package

// TikiTestResult is an alias to the canonical type in models package
type TikiTestResult = models.TikiTestResult

// UnifiedTestResult is a deprecated alias for TikiTestResult
type UnifiedTestResult = models.TikiTestResult

// TestMetadata is an alias to the canonical type in models package
type TestMetadata = models.TestMetadata

// TestSuite is an alias to the canonical type in models package  
type TestSuite = models.TestSuite

// TestCase is an alias to the canonical type in models package
type TestCase = models.TestCase

// TestStatus is an alias to the canonical type in models package
type TestStatus = models.TestStatus

// Assertion is an alias to the canonical type in models package
type Assertion = models.Assertion

// TestSummary is an alias to the canonical type in models package
type TestSummary = models.TestSummary

// Status constants referencing the canonical constants in models package
const (
	StatusPassed  = models.StatusPassed
	StatusFailed  = models.StatusFailed 
	StatusSkipped = models.StatusSkipped
	StatusError   = models.StatusError
)

// CombineOptions provides configuration for result combination.
type CombineOptions struct {
	// Source name for the combined result
	Source string

	// PreserveSuiteNames indicates whether to keep original suite names
	PreserveSuiteNames bool

	// NamePrefix to add to suite names when combining
	NamePrefix string

	// MergeEnvironment indicates whether to merge environment variables
	MergeEnvironment bool
}

// ExportOptions provides configuration for result export.
type ExportOptions struct {
	// Format specifies the output format ("json", "yaml", "xml")
	Format string

	// Indent specifies the indentation for formatted output
	Indent string

	// IncludeMetadata indicates whether to include metadata in output
	IncludeMetadata bool
}