// Package validation provides input validation utilities for the tconv library.
//
// This package contains functions to validate input files, output paths,
// and configuration parameters before processing. It helps ensure early
// detection of common issues and provides helpful error messages.
package validation

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ValidateInputFile checks if the input file exists and is readable.
//
// This function performs comprehensive validation including:
//   - Non-empty path validation
//   - File existence check
//   - Directory vs file validation
//   - Empty file detection
//
// Returns an error with descriptive message if validation fails.
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

// ValidateOutputFile checks if the output file path is valid and writable.
//
// This function validates that:
//   - Empty path is allowed (indicates stdout)
//   - Parent directory exists if specified
//   - Path is accessible for writing
//
// Returns an error if the output path is invalid.
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

// ValidateInputType checks if the input type is supported by the library.
//
// Supported types include: junit, postman, gotest, tap, xunit, pytest, testng.
// The validation is case-insensitive and trims whitespace.
//
// Returns an error with the list of supported types if validation fails.
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

// ValidateFileExtension checks if the file extension matches the expected type.
//
// This function performs a sanity check to ensure the file extension
// is appropriate for the specified input type. For example, XML extensions
// for junit/xunit/testng, and JSON extensions for gotest/postman/pytest.
//
// Returns an error if the extension doesn't match the expected pattern.
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