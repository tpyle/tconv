package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidateInputFile(t *testing.T) {
	// Test with empty path
	err := ValidateInputFile("")
	if err == nil {
		t.Error("Expected error for empty file path")
	}

	// Test with non-existent file
	err = ValidateInputFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}

	// Test with directory
	tempDir := t.TempDir()
	err = ValidateInputFile(tempDir)
	if err == nil {
		t.Error("Expected error for directory path")
	}

	// Test with empty file
	emptyFile := filepath.Join(tempDir, "empty.txt")
	if file, err := os.Create(emptyFile); err != nil {
		t.Fatalf("Failed to create empty file: %v", err)
	} else {
		file.Close()
	}
	err = ValidateInputFile(emptyFile)
	if err == nil {
		t.Error("Expected error for empty file")
	}

	// Test with valid file
	validFile := filepath.Join(tempDir, "valid.txt")
	if err := os.WriteFile(validFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create valid file: %v", err)
	}
	err = ValidateInputFile(validFile)
	if err != nil {
		t.Errorf("Expected no error for valid file, got: %v", err)
	}
}

func TestValidateOutputFile(t *testing.T) {
	// Test with empty path (stdout)
	err := ValidateOutputFile("")
	if err != nil {
		t.Errorf("Expected no error for empty output file (stdout), got: %v", err)
	}

	// Test with non-existent directory
	err = ValidateOutputFile("/nonexistent/dir/file.json")
	if err == nil {
		t.Error("Expected error for non-existent output directory")
	}

	// Test with valid path
	tempDir := t.TempDir()
	validPath := filepath.Join(tempDir, "output.json")
	err = ValidateOutputFile(validPath)
	if err != nil {
		t.Errorf("Expected no error for valid output path, got: %v", err)
	}
}

func TestValidateInputType(t *testing.T) {
	validTypes := []string{
		"junit", "postman", "gotest", "tap", "xunit", "pytest", "testng",
	}

	// Test valid types
	for _, inputType := range validTypes {
		err := ValidateInputType(inputType)
		if err != nil {
			t.Errorf("Expected no error for valid type %s, got: %v", inputType, err)
		}

		// Test case insensitive
		err = ValidateInputType(strings.ToUpper(inputType))
		if err != nil {
			t.Errorf("Expected no error for uppercase type %s, got: %v", inputType, err)
		}
	}

	// Test empty type
	err := ValidateInputType("")
	if err == nil {
		t.Error("Expected error for empty input type")
	}

	// Test invalid type
	err = ValidateInputType("invalid")
	if err == nil {
		t.Error("Expected error for invalid input type")
	}
}

func TestValidateFileExtension(t *testing.T) {
	testCases := []struct {
		filePath    string
		inputType   string
		expectError bool
	}{
		{"test.xml", "junit", false},
		{"test.json", "postman", false},
		{"test.json", "gotest", false},
		{"test.tap", "tap", false},
		{"test.txt", "tap", false},
		{"test.xml", "xunit", false},
		{"test.json", "pytest", false},
		{"test.xml", "testng", false},
		{"test.txt", "junit", true},   // Wrong extension
		{"test.xml", "postman", true}, // Wrong extension
		{"test.json", "junit", true},  // Wrong extension
	}

	for _, tc := range testCases {
		err := ValidateFileExtension(tc.filePath, tc.inputType)
		if tc.expectError && err == nil {
			t.Errorf("Expected error for %s with type %s", tc.filePath, tc.inputType)
		}
		if !tc.expectError && err != nil {
			t.Errorf("Expected no error for %s with type %s, got: %v", tc.filePath, tc.inputType, err)
		}
	}

	// Test unknown type (should not error)
	err := ValidateFileExtension("test.xyz", "unknown")
	if err != nil {
		t.Errorf("Expected no error for unknown type, got: %v", err)
	}
}

func TestValidateFileExtensionEdgeCases(t *testing.T) {
	// Test file without extension
	err := ValidateFileExtension("filename", "junit")
	if err == nil {
		t.Error("Expected error for file without extension")
	}

	// Test file with dot but no extension
	err = ValidateFileExtension("filename.", "junit")
	if err == nil {
		t.Error("Expected error for file with dot but no extension")
	}

	// Test case sensitivity
	err = ValidateFileExtension("test.XML", "junit")
	if err != nil {
		t.Errorf("Expected no error for uppercase extension, got: %v", err)
	}
}