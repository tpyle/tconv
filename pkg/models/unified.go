// Package models provides the core data structures for test result representation.
// This package is designed to be importable by other programs that want to use
// the tiki test result format.
package models

import (
	"fmt"
	"strings"
	"time"
)

// TikiTestResult represents the standardized tiki test result format that all
// parsers convert to. This is the primary export type for library consumers.
//
// The structure is designed to accommodate test results from various testing
// frameworks while providing a consistent interface for consumers.
type TikiTestResult struct {
	// Metadata contains information about the test run source and environment
	Metadata TestMetadata `json:"metadata" yaml:"metadata" xml:"metadata"`

	// TestSuites contains all test suites from the converted results
	TestSuites []TestSuite `json:"test_suites" yaml:"test_suites" xml:"test_suites>suite"`

	// Summary provides aggregate statistics across all test suites
	Summary TestSummary `json:"summary" yaml:"summary" xml:"summary"`
}

// TestMetadata contains contextual information about the test run.
type TestMetadata struct {
	// Source identifies the original format (e.g., "junit", "gotest", "combined-mixed")
	Source string `json:"source" yaml:"source" xml:"source,attr"`

	// Timestamp indicates when the test was run or when this result was created
	Timestamp time.Time `json:"timestamp" yaml:"timestamp" xml:"timestamp,attr"`

	// Version of the unified format schema
	Version string `json:"version" yaml:"version" xml:"version,attr"`

	// Framework identifies the testing framework used (e.g., "JUnit 5", "Go test")
	Framework string `json:"framework,omitempty" yaml:"framework,omitempty" xml:"framework,attr,omitempty"`

	// Environment contains key-value pairs of environment information
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty" xml:"environment,omitempty"`

	// SourceFile contains the original file path if converted from a file
	SourceFile string `json:"source_file,omitempty" yaml:"source_file,omitempty" xml:"source_file,attr,omitempty"`
}

// TestSuite represents a logical grouping of related tests.
type TestSuite struct {
	// Name is the identifier for this test suite
	Name string `json:"name" yaml:"name" xml:"name,attr"`

	// Package represents the namespace or package containing this suite
	Package string `json:"package,omitempty" yaml:"package,omitempty" xml:"package,attr,omitempty"`

	// Tests contains all individual test cases in this suite
	Tests []TestCase `json:"tests" yaml:"tests" xml:"tests>test"`

	// Errors is the count of tests that had runtime errors
	Errors int `json:"errors" yaml:"errors" xml:"errors,attr"`

	// Failures is the count of tests that failed assertions
	Failures int `json:"failures" yaml:"failures" xml:"failures,attr"`

	// Skipped is the count of tests that were skipped
	Skipped int `json:"skipped" yaml:"skipped" xml:"skipped,attr"`

	// Time is the total execution time for this suite in seconds
	Time float64 `json:"time" yaml:"time" xml:"time,attr"`

	// Properties contains additional metadata specific to this suite
	Properties map[string]string `json:"properties,omitempty" yaml:"properties,omitempty" xml:"properties,omitempty"`
}

// TestCase represents an individual test execution.
type TestCase struct {
	// Name is the test case identifier
	Name string `json:"name" yaml:"name" xml:"name,attr"`

	// ClassName represents the class or module containing this test
	ClassName string `json:"class_name,omitempty" yaml:"class_name,omitempty" xml:"class_name,attr,omitempty"`

	// Time is the execution duration in seconds
	Time float64 `json:"time" yaml:"time" xml:"time,attr"`

	// Status indicates the test outcome
	Status TestStatus `json:"status" yaml:"status" xml:"status,attr"`

	// Message provides a brief description of the test result
	Message string `json:"message,omitempty" yaml:"message,omitempty" xml:"message,omitempty"`

	// Details contains extended information such as stack traces
	Details string `json:"details,omitempty" yaml:"details,omitempty" xml:"details,omitempty"`

	// ErrorType specifies the type of error if Status is error or failed
	ErrorType string `json:"error_type,omitempty" yaml:"error_type,omitempty" xml:"error_type,attr,omitempty"`

	// SystemOut contains standard output captured during test execution
	SystemOut string `json:"system_out,omitempty" yaml:"system_out,omitempty" xml:"system_out,omitempty"`

	// SystemErr contains standard error output captured during test execution
	SystemErr string `json:"system_err,omitempty" yaml:"system_err,omitempty" xml:"system_err,omitempty"`

	// Assertions contains individual assertion results for frameworks that support it
	Assertions []Assertion `json:"assertions,omitempty" yaml:"assertions,omitempty" xml:"assertions>assertion,omitempty"`

	// Properties contains additional test-specific metadata
	Properties map[string]string `json:"properties,omitempty" yaml:"properties,omitempty" xml:"properties,omitempty"`
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
	Name string `json:"name" yaml:"name" xml:"name,attr"`

	// Expected contains the expected value for the assertion
	Expected string `json:"expected,omitempty" yaml:"expected,omitempty" xml:"expected,attr,omitempty"`

	// Actual contains the actual value that was compared
	Actual string `json:"actual,omitempty" yaml:"actual,omitempty" xml:"actual,attr,omitempty"`

	// Passed indicates whether this assertion succeeded
	Passed bool `json:"passed" yaml:"passed" xml:"passed,attr"`

	// Message provides details about the assertion result
	Message string `json:"message,omitempty" yaml:"message,omitempty" xml:"message,omitempty"`
}

// TestSummary provides aggregate statistics for a test result.
type TestSummary struct {
	// Total number of tests executed
	Total int `json:"total" yaml:"total" xml:"total,attr"`

	// Passed number of tests that completed successfully
	Passed int `json:"passed" yaml:"passed" xml:"passed,attr"`

	// Failed number of tests that failed assertions
	Failed int `json:"failed" yaml:"failed" xml:"failed,attr"`

	// Skipped number of tests that were not executed
	Skipped int `json:"skipped" yaml:"skipped" xml:"skipped,attr"`

	// Errors number of tests that encountered runtime errors
	Errors int `json:"errors" yaml:"errors" xml:"errors,attr"`

	// Duration is the total execution time in seconds
	Duration float64 `json:"duration" yaml:"duration" xml:"duration,attr"`
}

// PassRate calculates the percentage of tests that passed.
func (ts *TestSummary) PassRate() float64 {
	if ts.Total == 0 {
		return 0.0
	}
	return float64(ts.Passed) / float64(ts.Total) * 100.0
}

// CombineResults combines multiple TikiTestResult instances into a single result.
// This is the backward-compatible API for result combination.
func CombineResults(results []*TikiTestResult, combinedSource string) *TikiTestResult {
	if len(results) == 0 {
		return nil
	}

	if len(results) == 1 {
		return results[0]
	}

	combined := &TikiTestResult{
		Metadata: TestMetadata{
			Source:      combinedSource,
			Timestamp:   time.Now(),
			Version:     "2.0", // Updated version for improved schema
			Framework:   "Combined",
			Environment: make(map[string]string),
		},
		TestSuites: make([]TestSuite, 0, len(results)*2), // Pre-allocate with reasonable capacity
		Summary:    TestSummary{},
	}

	// Collect environment variables and frameworks from all results
	frameworks := make(map[string]bool)
	for _, result := range results {
		if result == nil {
			continue // Skip nil results
		}

		if result.Metadata.Framework != "" {
			frameworks[result.Metadata.Framework] = true
		}

		for key, value := range result.Metadata.Environment {
			combined.Metadata.Environment[key] = value
		}
	}

	// Build framework list
	var frameworkList []string
	for framework := range frameworks {
		frameworkList = append(frameworkList, framework)
	}
	if len(frameworkList) > 0 {
		combined.Metadata.Framework = fmt.Sprintf("Combined (%s)", strings.Join(frameworkList, ", "))
	}

	// Combine all test suites and calculate summary
	for i, result := range results {
		if result == nil {
			continue
		}

		for _, suite := range result.TestSuites {
			// Add a prefix to suite names to avoid conflicts
			combinedSuite := suite
			if len(results) > 1 {
				combinedSuite.Name = fmt.Sprintf("File%d_%s", i+1, suite.Name)
				if combinedSuite.Package != "" {
					combinedSuite.Package = fmt.Sprintf("file%d.%s", i+1, suite.Package)
				}
			}
			combined.TestSuites = append(combined.TestSuites, combinedSuite)
		}

		// Aggregate summary statistics
		combined.Summary.Total += result.Summary.Total
		combined.Summary.Passed += result.Summary.Passed
		combined.Summary.Failed += result.Summary.Failed
		combined.Summary.Skipped += result.Summary.Skipped
		combined.Summary.Errors += result.Summary.Errors
		combined.Summary.Duration += result.Summary.Duration
	}

	return combined
}

// UnifiedTestResult is a deprecated alias for TikiTestResult.
// Use TikiTestResult instead.
type UnifiedTestResult = TikiTestResult