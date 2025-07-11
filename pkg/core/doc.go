// Package core provides the fundamental interfaces, types, and implementations
// for the tconv (test converter) library.
//
// The core package is designed to be the primary export point for other
// programs that want to use tconv as a library. It provides clean abstractions
// and interfaces that make the library extensible and easy to integrate.
//
// # Key Interfaces
//
// The package defines several key interfaces:
//
//   - TestConverter: For converting test files to unified format
//   - TypeDetector: For automatically detecting test file types
//   - ResultCombiner: For merging multiple test results
//   - TestResultExporter: For exporting results in various formats
//   - Parser: For format-specific parsing implementations
//
// # Primary Types
//
// The UnifiedTestResult type is the central data structure that represents
// test results in a standardized format. All parsers convert their specific
// formats to this unified representation.
//
// # Error Handling
//
// The package provides structured error handling through the TconvError type,
// which categorizes errors and provides helpful context and suggestions.
//
// # Usage Examples
//
// Basic conversion:
//
//	detector := core.NewDefaultTypeDetector()
//	fileType, err := detector.DetectTypeFromFile("test-results.xml")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Parse using appropriate parser...
//	var result *core.UnifiedTestResult
//	// ... parsing logic
//
//	// Export to JSON
//	var buf bytes.Buffer
//	err = core.ExportToJSON(result, &buf)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// Combining multiple results:
//
//	combiner := core.NewDefaultResultCombiner()
//	options := core.CombineOptions{
//	    Source:             "combined-tests",
//	    PreserveSuiteNames: false,
//	    MergeEnvironment:   true,
//	}
//	combined, err := combiner.Combine(results, options)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
// # Thread Safety
//
// All implementations in this package are designed to be thread-safe and
// stateless, making them suitable for concurrent use.
package core