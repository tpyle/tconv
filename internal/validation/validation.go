package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateInputFile checks if the input file exists and is readable
func ValidateInputFile(filePath string) error {
	if filePath == "" {
		return fmt.Errorf("input file path cannot be empty")
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", filePath)
		}
		return fmt.Errorf("cannot access input file: %w", err)
	}

	if info.IsDir() {
		return fmt.Errorf("input path is a directory, not a file: %s", filePath)
	}

	if info.Size() == 0 {
		return fmt.Errorf("input file is empty: %s", filePath)
	}

	return nil
}

// ValidateOutputFile checks if the output file path is valid
func ValidateOutputFile(filePath string) error {
	if filePath == "" {
		// Empty output file means stdout, which is valid
		return nil
	}

	dir := filepath.Dir(filePath)
	if dir != "." {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("output directory does not exist: %s", dir)
		}
	}

	return nil
}

// ValidateInputType checks if the input type is supported
func ValidateInputType(inputType string) error {
	if inputType == "" {
		return fmt.Errorf("input type cannot be empty")
	}

	supportedTypes := []string{
		"junit", "postman", "gotest", "tap", "xunit", "pytest", "testng",
	}

	inputType = strings.ToLower(strings.TrimSpace(inputType))
	for _, supported := range supportedTypes {
		if inputType == supported {
			return nil
		}
	}

	return fmt.Errorf("unsupported input type: %s (supported types: %s)", 
		inputType, strings.Join(supportedTypes, ", "))
}

// ValidateFileExtension checks if the file extension matches the expected type
func ValidateFileExtension(filePath, inputType string) error {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	expectedExtensions := map[string][]string{
		"junit":   {".xml"},
		"postman": {".json"},
		"gotest":  {".json"},
		"tap":     {".tap", ".txt"},
		"xunit":   {".xml"},
		"pytest":  {".json"},
		"testng":  {".xml"},
	}

	expected, exists := expectedExtensions[inputType]
	if !exists {
		return nil // Skip validation for unknown types
	}

	for _, expectedExt := range expected {
		if ext == expectedExt {
			return nil
		}
	}

	return fmt.Errorf("file extension %s does not match expected extensions for %s: %s", 
		ext, inputType, strings.Join(expected, ", "))
}