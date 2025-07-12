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
package tconv

import (
	"io"

	"github.com/tpyle/tconv/internal/converter"
	"github.com/tpyle/tconv/internal/detector"
	"github.com/tpyle/tconv/internal/exporters"
	"github.com/tpyle/tconv/pkg/models"
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

type Format string

const (
	FormatJUnit   Format = "junit"   // JUnit XML test results
	FormatTestNG  Format = "testng"  // TestNG XML test results
	FormatXUnit   Format = "xunit"   // xUnit.net XML test results
	FormatGoTest  Format = "gotest"  // Go test JSON output
	FormatPostman Format = "postman" // Postman collection runner JSON
	FormatPyTest  Format = "pytest"  // pytest JSON report
	FormatTAP     Format = "tap"     // Test Anything Protocol
)

// SupportedFormats returns a list of all supported test file formats for import.
func SupportedFormats() []Format {
	return []Format{
		FormatJUnit,   // JUnit XML test results
		FormatTestNG,  // TestNG XML test results
		FormatXUnit,   // xUnit.net XML test results
		FormatGoTest,  // Go test JSON output
		FormatPostman, // Postman collection runner JSON
		FormatPyTest,  // pytest JSON report
		FormatTAP,     // Test Anything Protocol
	}
}

// Export converts a unified test result to the specified format and writes it to a file.
//
// This function performs the reverse conversion - from unified format to specific formats.
//
// Parameters:
//   - result: The unified test result to export
//   - format: Target format ("junit", "tap", "gotest")
//   - outputFile: Path where the converted file will be written
//
// Returns an error if the export fails.
func Export(result *models.UnifiedTestResult, format, outputFile string) error {
	manager := exporters.NewExportManager()
	return manager.Export(result, format, outputFile)
}

// ExportToWriter converts a unified test result to the specified format and writes it to a writer.
//
// This function performs the reverse conversion using an io.Writer instead of a file.
//
// Parameters:
//   - result: The unified test result to export
//   - format: Target format ("junit", "tap", "gotest")
//   - writer: Writer where the converted content will be written
//
// Returns an error if the export fails.
func ExportToWriter(result *models.UnifiedTestResult, format string, writer io.Writer) error {
	manager := exporters.NewExportManager()
	return manager.Write(result, format, writer)
}

// SupportedExportFormats returns a list of all supported export formats.
func SupportedExportFormats() []string {
	manager := exporters.NewExportManager()
	return manager.SupportedExportFormats()
}

// LoadUnified loads a unified test result from a JSON file.
//
// This function can be used to read previously converted results for export to other formats.
//
// Parameters:
//   - inputFile: Path to the unified JSON file
//
// Returns the loaded UnifiedTestResult or an error if loading fails.
func LoadUnified(inputFile string) (*models.UnifiedTestResult, error) {
	conv := converter.New()
	return conv.LoadUnified(inputFile)
}
