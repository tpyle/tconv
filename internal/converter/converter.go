// Package converter provides the core conversion functionality for transforming
// test results from various formats into the standardized tiki format.
//
// This package orchestrates the conversion process by:
//   - Detecting file types when not explicitly specified
//   - Routing files to appropriate format-specific parsers
//   - Handling multiple file conversions and result combination
//   - Managing output to files or writers
//
// The Converter type is the main entry point and is designed to be thread-safe
// and stateless, making it suitable for concurrent use.
package converter

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/tpyle/tconv/internal/detector"
	"github.com/tpyle/tconv/internal/parsers"
	"github.com/tpyle/tconv/pkg/models"
)

// Converter orchestrates the conversion of test files from various formats
// to the standardized tiki format. It supports single file conversion,
// multiple file conversion with result combination, and automatic type detection.
type Converter struct{}

// New creates a new Converter instance.
// The converter is stateless and thread-safe.
func New() *Converter {
	return &Converter{}
}

// Convert processes a single test file and converts it to tiki JSON format.
//
// Parameters:
//   - inputFile: Path to the input test file
//   - inputType: Format type ("junit", "testng", "xunit", "gotest", "postman", "pytest", "tap")
//   - outputFile: Path for output (empty string writes to stdout)
//
// Returns an error if conversion fails.
func (c *Converter) Convert(inputFile, inputType, outputFile string) error {

	var result *models.UnifiedTestResult
	var err error

	switch inputType {
	case "junit":
		result, err = parsers.ParseJUnit(inputFile)
	case "postman":
		result, err = parsers.ParsePostman(inputFile)
	case "gotest":
		result, err = parsers.ParseGoTest(inputFile)
	case "tap":
		result, err = parsers.ParseTAP(inputFile)
	case "xunit":
		result, err = parsers.ParseXUnit(inputFile)
	case "pytest":
		result, err = parsers.ParsePytest(inputFile)
	case "testng":
		result, err = parsers.ParseTestNG(inputFile)
	default:
		return fmt.Errorf("unsupported input type: %s", inputType)
	}

	if err != nil {
		return fmt.Errorf("failed to parse input file: %w", err)
	}

	var output io.Writer
	if outputFile == "" {
		output = os.Stdout
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		output = file
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	return nil
}

// ConvertMultiple processes multiple test files and combines them into a single tiki result.
//
// If inputType is empty, file types will be auto-detected for each file.
// If inputType is specified, all files will be treated as that type.
//
// Parameters:
//   - inputFiles: Slice of paths to input test files
//   - inputType: Format type (empty for auto-detection)
//   - outputFile: Path for output (empty string writes to stdout)
//
// Returns an error if any conversion or combination fails.
func (c *Converter) ConvertMultiple(inputFiles []string, inputType, outputFile string) error {
	if len(inputFiles) == 0 {
		return fmt.Errorf("no input files provided")
	}

	// For single file only (not multiple files), use the original Convert method
	if len(inputFiles) == 1 && inputType != "" {
		return c.Convert(inputFiles[0], inputType, outputFile)
	}

	// Parse all input files
	results := make([]*models.UnifiedTestResult, 0, len(inputFiles))
	detectedTypes := make([]string, 0, len(inputFiles))
	
	for _, inputFile := range inputFiles {
		var result *models.UnifiedTestResult
		var err error
		var fileType string

		// Determine file type - use specified type or auto-detect
		if inputType != "" {
			fileType = inputType
		} else {
			fileType, err = detector.DetectFileType(inputFile)
			if err != nil {
				return fmt.Errorf("failed to detect file type for %s: %w", inputFile, err)
			}
		}

		// Parse based on detected/specified type
		switch fileType {
		case "junit":
			result, err = parsers.ParseJUnit(inputFile)
		case "postman":
			result, err = parsers.ParsePostman(inputFile)
		case "gotest":
			result, err = parsers.ParseGoTest(inputFile)
		case "tap":
			result, err = parsers.ParseTAP(inputFile)
		case "xunit":
			result, err = parsers.ParseXUnit(inputFile)
		case "pytest":
			result, err = parsers.ParsePytest(inputFile)
		case "testng":
			result, err = parsers.ParseTestNG(inputFile)
		default:
			return fmt.Errorf("unsupported input type: %s", fileType)
		}

		if err != nil {
			return fmt.Errorf("failed to parse input file %s (type: %s): %w", inputFile, fileType, err)
		}

		results = append(results, result)
		detectedTypes = append(detectedTypes, fileType)
	}

	// Determine combined source name
	var combinedSource string
	if inputType != "" {
		combinedSource = fmt.Sprintf("combined-%s", inputType)
	} else {
		combinedSource = "combined-mixed"
	}

	// Combine all results
	combinedResult := models.CombineResults(results, combinedSource)
	if combinedResult == nil {
		return fmt.Errorf("failed to combine results")
	}

	// Write output
	var output io.Writer
	if outputFile == "" {
		output = os.Stdout
	} else {
		file, err := os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("failed to create output file: %w", err)
		}
		defer file.Close()
		output = file
	}

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(combinedResult); err != nil {
		return fmt.Errorf("failed to encode result: %w", err)
	}

	return nil
}

// LoadTiki loads a tiki test result from a JSON file.
//
// This function can be used to read previously converted results
// for further processing or export to other formats.
//
// Parameters:
//   - inputFile: Path to the tiki JSON file
//
// Returns the loaded TikiTestResult or an error if loading fails.
func (c *Converter) LoadTiki(inputFile string) (*models.TikiTestResult, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	var result models.TikiTestResult
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode tiki test result: %w", err)
	}

	return &result, nil
}

// LoadUnified is deprecated. Use LoadTiki instead.
//
// This method is maintained for backward compatibility and simply
// delegates to LoadTiki.
func (c *Converter) LoadUnified(inputFile string) (*models.TikiTestResult, error) {
	return c.LoadTiki(inputFile)
}