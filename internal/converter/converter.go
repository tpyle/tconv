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

type Converter struct{}

func New() *Converter {
	return &Converter{}
}

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

// LoadUnified loads a unified test result from a JSON file
func (c *Converter) LoadUnified(inputFile string) (*models.UnifiedTestResult, error) {
	file, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open input file: %w", err)
	}
	defer file.Close()

	var result models.UnifiedTestResult
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode unified test result: %w", err)
	}

	return &result, nil
}