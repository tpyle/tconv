package core

import (
	"fmt"
	"strings"
	"time"
)

// DefaultResultCombiner implements the ResultCombiner interface.
type DefaultResultCombiner struct{}

// NewDefaultResultCombiner creates a new instance of DefaultResultCombiner.
func NewDefaultResultCombiner() *DefaultResultCombiner {
	return &DefaultResultCombiner{}
}

// Combine merges multiple unified test results into a single result.
func (c *DefaultResultCombiner) Combine(results []*UnifiedTestResult, options CombineOptions) (*UnifiedTestResult, error) {
	if len(results) == 0 {
		return nil, fmt.Errorf("no results provided for combination")
	}

	if len(results) == 1 {
		return results[0], nil
	}

	// Initialize the combined result
	combined := &UnifiedTestResult{
		Metadata: TestMetadata{
			Source:      options.Source,
			Timestamp:   time.Now(),
			Version:     "2.0", // Updated version for new schema
			Framework:   "Combined",
			Environment: make(map[string]string),
		},
		TestSuites: make([]TestSuite, 0, len(results)*2), // Pre-allocate with reasonable capacity
		Summary:    TestSummary{},
	}

	if combined.Metadata.Source == "" {
		combined.Metadata.Source = "combined-mixed"
	}

	// Collect frameworks and environment variables from all results
	frameworks := make(map[string]bool)
	for _, result := range results {
		if result == nil {
			continue // Skip nil results
		}

		if result.Metadata.Framework != "" {
			frameworks[result.Metadata.Framework] = true
		}

		if options.MergeEnvironment {
			for key, value := range result.Metadata.Environment {
				combined.Metadata.Environment[key] = value
			}
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
			combinedSuite := suite

			// Handle suite naming based on options
			if !options.PreserveSuiteNames && len(results) > 1 {
				if options.NamePrefix != "" {
					combinedSuite.Name = fmt.Sprintf("%s_%s", options.NamePrefix, suite.Name)
				} else {
					combinedSuite.Name = fmt.Sprintf("File%d_%s", i+1, suite.Name)
				}

				if combinedSuite.Package != "" {
					if options.NamePrefix != "" {
						combinedSuite.Package = fmt.Sprintf("%s.%s", strings.ToLower(options.NamePrefix), suite.Package)
					} else {
						combinedSuite.Package = fmt.Sprintf("file%d.%s", i+1, suite.Package)
					}
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

	return combined, nil
}

// CombineResults is a convenience function that uses default options.
// This maintains backward compatibility with the existing API.
func CombineResults(results []*UnifiedTestResult, combinedSource string) *UnifiedTestResult {
	combiner := NewDefaultResultCombiner()
	options := CombineOptions{
		Source:             combinedSource,
		PreserveSuiteNames: false,
		MergeEnvironment:   true,
	}

	result, err := combiner.Combine(results, options)
	if err != nil {
		return nil // Maintain backward compatibility
	}

	return result
}