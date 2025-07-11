// Package core defines the core interfaces and types for the tconv library.
// This package provides the fundamental abstractions that make the library
// extensible and suitable for consumption by other programs.
package core

import (
	"io"
	"time"
)

// TestConverter defines the interface for converting test files to unified format.
// Implementations should be thread-safe and stateless.
type TestConverter interface {
	// Convert processes a single test file and returns a unified result.
	Convert(input io.Reader, sourceType string) (*UnifiedTestResult, error)

	// ConvertFile processes a test file by path and returns a unified result.
	ConvertFile(filePath, sourceType string) (*UnifiedTestResult, error)

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
	// Combine merges multiple unified test results into a single result.
	Combine(results []*UnifiedTestResult, options CombineOptions) (*UnifiedTestResult, error)
}

// TestResultExporter defines the interface for exporting test results.
type TestResultExporter interface {
	// Export writes the unified test result to the specified output.
	Export(result *UnifiedTestResult, output io.Writer, options ExportOptions) error
}

// Parser defines the interface that all format-specific parsers must implement.
type Parser interface {
	// Parse processes input data and returns a unified test result.
	Parse(input io.Reader) (*UnifiedTestResult, error)

	// Validate checks if the input is valid for this parser.
	Validate(input io.Reader) error

	// Type returns the type identifier for this parser.
	Type() string
}

// UnifiedTestResult represents the standardized test result format that all
// parsers convert to. This is the primary export type for library consumers.
type UnifiedTestResult struct {
	// Metadata contains information about the test run source and environment
	Metadata TestMetadata `json:"metadata" yaml:"metadata"`

	// TestSuites contains all test suites from the converted results
	TestSuites []TestSuite `json:"test_suites" yaml:"test_suites"`

	// Summary provides aggregate statistics across all test suites
	Summary TestSummary `json:"summary" yaml:"summary"`
}

// TestMetadata contains contextual information about the test run.
type TestMetadata struct {
	// Source identifies the original format (e.g., "junit", "gotest", "combined-mixed")
	Source string `json:"source" yaml:"source"`

	// Timestamp indicates when the test was run or when this result was created
	Timestamp time.Time `json:"timestamp" yaml:"timestamp"`

	// Version of the unified format schema
	Version string `json:"version" yaml:"version"`

	// Framework identifies the testing framework used (e.g., "JUnit 5", "Go test")
	Framework string `json:"framework,omitempty" yaml:"framework,omitempty"`

	// Environment contains key-value pairs of environment information
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`

	// SourceFile contains the original file path if converted from a file
	SourceFile string `json:"source_file,omitempty" yaml:"source_file,omitempty"`
}

// TestSuite represents a logical grouping of related tests.
type TestSuite struct {
	// Name is the identifier for this test suite
	Name string `json:"name" yaml:"name"`

	// Package represents the namespace or package containing this suite
	Package string `json:"package,omitempty" yaml:"package,omitempty"`

	// Tests contains all individual test cases in this suite
	Tests []TestCase `json:"tests" yaml:"tests"`

	// Errors is the count of tests that had runtime errors
	Errors int `json:"errors" yaml:"errors"`

	// Failures is the count of tests that failed assertions
	Failures int `json:"failures" yaml:"failures"`

	// Skipped is the count of tests that were skipped
	Skipped int `json:"skipped" yaml:"skipped"`

	// Time is the total execution time for this suite in seconds
	Time float64 `json:"time" yaml:"time"`

	// Properties contains additional metadata specific to this suite
	Properties map[string]string `json:"properties,omitempty" yaml:"properties,omitempty"`
}

// TestCase represents an individual test execution.
type TestCase struct {
	// Name is the test case identifier
	Name string `json:"name" yaml:"name"`

	// ClassName represents the class or module containing this test
	ClassName string `json:"class_name,omitempty" yaml:"class_name,omitempty"`

	// Time is the execution duration in seconds
	Time float64 `json:"time" yaml:"time"`

	// Status indicates the test outcome
	Status TestStatus `json:"status" yaml:"status"`

	// Message provides a brief description of the test result
	Message string `json:"message,omitempty" yaml:"message,omitempty"`

	// Details contains extended information such as stack traces
	Details string `json:"details,omitempty" yaml:"details,omitempty"`

	// ErrorType specifies the type of error if Status is error or failed
	ErrorType string `json:"error_type,omitempty" yaml:"error_type,omitempty"`

	// SystemOut contains standard output captured during test execution
	SystemOut string `json:"system_out,omitempty" yaml:"system_out,omitempty"`

	// SystemErr contains standard error output captured during test execution
	SystemErr string `json:"system_err,omitempty" yaml:"system_err,omitempty"`

	// Assertions contains individual assertion results for frameworks that support it
	Assertions []Assertion `json:"assertions,omitempty" yaml:"assertions,omitempty"`

	// Properties contains additional test-specific metadata
	Properties map[string]string `json:"properties,omitempty" yaml:"properties,omitempty"`
}

// TestStatus represents the possible outcomes of a test execution.
type TestStatus string

// Standard test status values that all parsers should map to.
const (
	StatusPassed  TestStatus = "passed"  // Test completed successfully
	StatusFailed  TestStatus = "failed"  // Test failed assertion(s)
	StatusSkipped TestStatus = "skipped" // Test was not executed
	StatusError   TestStatus = "error"   // Test encountered a runtime error
)

// String implements fmt.Stringer for TestStatus.
func (ts TestStatus) String() string {
	return string(ts)
}

// IsValid checks if the TestStatus is one of the defined values.
func (ts TestStatus) IsValid() bool {
	switch ts {
	case StatusPassed, StatusFailed, StatusSkipped, StatusError:
		return true
	default:
		return false
	}
}

// Assertion represents an individual assertion within a test case.
type Assertion struct {
	// Name is the assertion identifier or description
	Name string `json:"name" yaml:"name"`

	// Expected contains the expected value for the assertion
	Expected string `json:"expected,omitempty" yaml:"expected,omitempty"`

	// Actual contains the actual value that was compared
	Actual string `json:"actual,omitempty" yaml:"actual,omitempty"`

	// Passed indicates whether this assertion succeeded
	Passed bool `json:"passed" yaml:"passed"`

	// Message provides details about the assertion result
	Message string `json:"message,omitempty" yaml:"message,omitempty"`
}

// TestSummary provides aggregate statistics for a test result.
type TestSummary struct {
	// Total number of tests executed
	Total int `json:"total" yaml:"total"`

	// Passed number of tests that completed successfully
	Passed int `json:"passed" yaml:"passed"`

	// Failed number of tests that failed assertions
	Failed int `json:"failed" yaml:"failed"`

	// Skipped number of tests that were not executed
	Skipped int `json:"skipped" yaml:"skipped"`

	// Errors number of tests that encountered runtime errors
	Errors int `json:"errors" yaml:"errors"`

	// Duration is the total execution time in seconds
	Duration float64 `json:"duration" yaml:"duration"`
}

// PassRate calculates the percentage of tests that passed.
func (ts *TestSummary) PassRate() float64 {
	if ts.Total == 0 {
		return 0.0
	}
	return float64(ts.Passed) / float64(ts.Total) * 100.0
}

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