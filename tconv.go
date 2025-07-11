// Package tconv provides a unified test result conversion library.
//
// This package converts test results from various formats (JUnit, TestNG, xUnit,
// Go test, Postman, pytest, TAP) into a standardized JSON format.
//
// The main entry points for library consumers are:
//   - Convert(): Convert single files
//   - ConvertMultiple(): Convert and combine multiple files
//   - DetectFileType(): Auto-detect test file formats
//   - The models package: Access to unified data structures
//
// Basic usage:
//
//	import "github.com/tpyle/tconv"
//
//	// Convert a single file
//	err := tconv.Convert("test-results.xml", "junit", "output.json")
//
//	// Auto-detect and combine multiple files
//	files := []string{"junit.xml", "gotest.json", "pytest.json"}
//	err := tconv.ConvertMultiple(files, "", "combined.json")
//
//	// Detect file type
//	fileType, err := tconv.DetectFileType("unknown.xml")
//
package tconv

import (
	"github.com/tpyle/tconv/internal/converter"
	"github.com/tpyle/tconv/internal/detector"
)

// Convert processes a single test file and converts it to the unified JSON format.
//
// Parameters:
//   - inputFile: Path to the input test file
//   - inputType: Format type ("junit", "testng", "xunit", "gotest", "postman", "pytest", "tap")
//   - outputFile: Path where the converted JSON will be written
//
// Returns an error if the conversion fails.
func Convert(inputFile, inputType, outputFile string) error {
	conv := converter.New()
	return conv.Convert(inputFile, inputType, outputFile)
}

// ConvertMultiple processes multiple test files and combines them into a single unified result.
//
// If inputType is empty, file types will be auto-detected.
// If inputType is specified, all files will be treated as that type.
//
// Parameters:
//   - inputFiles: Slice of paths to input test files
//   - inputType: Format type (empty for auto-detection, or "junit", "testng", etc.)
//   - outputFile: Path where the combined JSON will be written
//
// Returns an error if any conversion fails.
func ConvertMultiple(inputFiles []string, inputType, outputFile string) error {
	conv := converter.New()
	return conv.ConvertMultiple(inputFiles, inputType, outputFile)
}

// DetectFileType automatically determines the test file format based on file extension and content.
//
// Supported formats: junit, testng, xunit, gotest, postman, pytest, tap
//
// Parameters:
//   - filePath: Path to the test file to analyze
//
// Returns the detected format type or an error if detection fails.
func DetectFileType(filePath string) (string, error) {
	return detector.DetectFileType(filePath)
}

// SupportedFormats returns a list of all supported test file formats.
func SupportedFormats() []string {
	return []string{
		"junit",   // JUnit XML test results
		"testng",  // TestNG XML test results  
		"xunit",   // xUnit.net XML test results
		"gotest",  // Go test JSON output
		"postman", // Postman collection runner JSON
		"pytest",  // pytest JSON report
		"tap",     // Test Anything Protocol
	}
}