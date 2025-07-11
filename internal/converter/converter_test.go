package converter

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConverterJUnit(t *testing.T) {
	converter := New()
	
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.json")
	
	err := converter.Convert("../../testdata/sample_junit.xml", "junit", outputFile)
	if err != nil {
		t.Fatalf("Failed to convert JUnit file: %v", err)
	}
	
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Output file is empty")
	}
}

func TestConverterPostman(t *testing.T) {
	converter := New()
	
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.json")
	
	err := converter.Convert("../../testdata/sample_postman.json", "postman", outputFile)
	if err != nil {
		t.Fatalf("Failed to convert Postman file: %v", err)
	}
	
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Output file is empty")
	}
}

func TestConverterGoTest(t *testing.T) {
	converter := New()
	
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.json")
	
	err := converter.Convert("../../testdata/sample_gotest.json", "gotest", outputFile)
	if err != nil {
		t.Fatalf("Failed to convert Go test file: %v", err)
	}
	
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Output file is empty")
	}
}

func TestConverterUnsupportedType(t *testing.T) {
	converter := New()
	
	err := converter.Convert("../../testdata/sample_junit.xml", "unsupported", "")
	if err == nil {
		t.Error("Expected error for unsupported input type")
	}
	
	expectedError := "unsupported input type: unsupported"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestConverterNonExistentFile(t *testing.T) {
	converter := New()
	
	err := converter.Convert("nonexistent.xml", "junit", "")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}

func TestConverterTAP(t *testing.T) {
	converter := New()
	
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.json")
	
	err := converter.Convert("../../testdata/sample_tap.tap", "tap", outputFile)
	if err != nil {
		t.Fatalf("Failed to convert TAP file: %v", err)
	}
	
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Output file is empty")
	}
}

func TestConverterXUnit(t *testing.T) {
	converter := New()
	
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.json")
	
	err := converter.Convert("../../testdata/sample_xunit.xml", "xunit", outputFile)
	if err != nil {
		t.Fatalf("Failed to convert xUnit file: %v", err)
	}
	
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Output file is empty")
	}
}

func TestConverterPytest(t *testing.T) {
	converter := New()
	
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.json")
	
	err := converter.Convert("../../testdata/sample_pytest.json", "pytest", outputFile)
	if err != nil {
		t.Fatalf("Failed to convert pytest file: %v", err)
	}
	
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Output file is empty")
	}
}

func TestConverterTestNG(t *testing.T) {
	converter := New()
	
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "output.json")
	
	err := converter.Convert("../../testdata/sample_testng.xml", "testng", outputFile)
	if err != nil {
		t.Fatalf("Failed to convert TestNG file: %v", err)
	}
	
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Output file is empty")
	}
}

func TestConverterAllNewFormats(t *testing.T) {
	converter := New()
	
	testCases := []struct {
		name      string
		file      string
		inputType string
	}{
		{"TAP", "../../testdata/sample_tap.tap", "tap"},
		{"xUnit", "../../testdata/sample_xunit.xml", "xunit"},
		{"pytest", "../../testdata/sample_pytest.json", "pytest"},
		{"TestNG", "../../testdata/sample_testng.xml", "testng"},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tempDir := t.TempDir()
			outputFile := filepath.Join(tempDir, "output.json")
			
			err := converter.Convert(tc.file, tc.inputType, outputFile)
			if err != nil {
				t.Fatalf("Failed to convert %s file: %v", tc.name, err)
			}
			
			if _, err := os.Stat(outputFile); os.IsNotExist(err) {
				t.Errorf("Output file was not created for %s", tc.name)
			}
			
			content, err := os.ReadFile(outputFile)
			if err != nil {
				t.Fatalf("Failed to read output file for %s: %v", tc.name, err)
			}
			
			if len(content) == 0 {
				t.Errorf("Output file is empty for %s", tc.name)
			}
		})
	}
}

func TestConverterValidation(t *testing.T) {
	converter := New()
	
	// Test invalid input type
	err := converter.Convert("../../testdata/sample_junit.xml", "invalid", "")
	if err == nil {
		t.Error("Expected error for invalid input type")
	}
	
	// Test non-existent file
	err = converter.Convert("nonexistent.xml", "junit", "")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
	
	// Test wrong file extension
	err = converter.Convert("../../testdata/sample_junit.xml", "postman", "")
	if err == nil {
		t.Error("Expected error for wrong file extension")
	}
}


func TestConverterInvalidOutputPath(t *testing.T) {
	converter := New()
	
	err := converter.Convert("../../testdata/sample_junit.xml", "junit", "/invalid/path/output.json")
	if err == nil {
		t.Error("Expected error for invalid output path")
	}
}

func TestConverterMultipleFiles(t *testing.T) {
	converter := New()
	
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "combined_output.json")
	
	inputFiles := []string{
		"../../testdata/sample_junit.xml",
		"../../testdata/sample_junit.xml",
	}
	
	err := converter.ConvertMultiple(inputFiles, "junit", outputFile)
	if err != nil {
		t.Fatalf("Failed to convert multiple files: %v", err)
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Output file is empty")
	}
	
	// Should contain combined source indicator
	if !strings.Contains(string(content), "combined-junit") {
		t.Error("Output should indicate combined source")
	}
}

func TestConverterMultipleFilesSingleFile(t *testing.T) {
	converter := New()
	
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "single_output.json")
	
	inputFiles := []string{"../../testdata/sample_junit.xml"}
	
	err := converter.ConvertMultiple(inputFiles, "junit", outputFile)
	if err != nil {
		t.Fatalf("Failed to convert single file through multiple interface: %v", err)
	}
	
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}
	
	if len(content) == 0 {
		t.Error("Output file is empty")
	}
}

func TestConverterMultipleFilesEmpty(t *testing.T) {
	converter := New()
	
	err := converter.ConvertMultiple([]string{}, "junit", "")
	if err == nil {
		t.Error("Expected error for empty file list")
	}
}

func TestConverterMultipleFilesNonExistent(t *testing.T) {
	converter := New()
	
	inputFiles := []string{
		"../../testdata/sample_junit.xml",
		"nonexistent.xml",
	}
	
	err := converter.ConvertMultiple(inputFiles, "junit", "")
	if err == nil {
		t.Error("Expected error for non-existent file in list")
	}
}