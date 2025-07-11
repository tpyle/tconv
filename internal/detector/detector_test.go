package detector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectFileType(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()

	tests := []struct {
		name        string
		filename    string
		content     string
		expected    string
		expectError bool
	}{
		{
			name:     "JUnit XML file",
			filename: "junit.xml",
			content:  "<?xml version=\"1.0\"?>\n<testsuite name=\"test\" tests=\"1\">\n<testcase name=\"test1\"/>\n</testsuite>",
			expected: "junit",
		},
		{
			name:     "TestNG XML file",
			filename: "testng.xml",
			content:  "<?xml version=\"1.0\"?>\n<testng-results total=\"1\" passed=\"1\">\n<suite name=\"test\"/>\n</testng-results>",
			expected: "testng",
		},
		{
			name:     "xUnit XML file",
			filename: "xunit.xml",
			content:  "<?xml version=\"1.0\"?>\n<assemblies>\n<assembly name=\"test\">\n<collection name=\"tests\"/>\n</assembly>\n</assemblies>",
			expected: "xunit",
		},
		{
			name:     "Go test JSON file",
			filename: "gotest.json",
			content:  "{\"Time\":\"2023-01-01T10:00:00Z\",\"Action\":\"run\",\"Package\":\"test\",\"Test\":\"TestExample\"}\n{\"Time\":\"2023-01-01T10:00:01Z\",\"Action\":\"pass\",\"Package\":\"test\",\"Test\":\"TestExample\"}",
			expected: "gotest",
		},
		{
			name:     "Postman JSON file",
			filename: "postman.json",
			content:  "{\"info\":{\"name\":\"Test Collection\",\"schema\":\"https://schema.getpostman.com/json/collection/v2.1.0/collection.json\"},\"item\":[]}",
			expected: "postman",
		},
		{
			name:     "pytest JSON file",
			filename: "pytest.json",
			content:  "{\"tests\":[{\"nodeid\":\"test.py::test_example\",\"outcome\":\"passed\"}],\"summary\":{\"total\":1,\"passed\":1}}",
			expected: "pytest",
		},
		{
			name:     "TAP file",
			filename: "test.tap",
			content:  "TAP version 13\n1..2\nok 1 - test passes\nnot ok 2 - test fails",
			expected: "tap",
		},
		{
			name:        "Non-existent file",
			filename:    "nonexistent.xml",
			expectError: true,
		},
		{
			name:        "Unrecognized format",
			filename:    "unknown.txt",
			content:     "This is not a test file",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := filepath.Join(tempDir, tt.filename)

			// Create the test file if content is provided
			if tt.content != "" {
				err := os.WriteFile(filePath, []byte(tt.content), 0644)
				if err != nil {
					t.Fatalf("Failed to create test file: %v", err)
				}
			}

			// Test the detection
			result, err := DetectFileType(filePath)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			}
		})
	}
}

func TestDetectFileTypeWithRealSamples(t *testing.T) {
	// Test with actual sample files from testdata
	tests := []struct {
		filePath string
		expected string
	}{
		{"../../testdata/sample_junit.xml", "junit"},
		{"../../testdata/sample_testng.xml", "testng"},
		{"../../testdata/sample_xunit.xml", "xunit"},
		{"../../testdata/sample_gotest.json", "gotest"},
		{"../../testdata/sample_postman.json", "postman"},
		{"../../testdata/sample_pytest.json", "pytest"},
		{"../../testdata/sample_tap.tap", "tap"},
	}

	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			// Check if file exists first
			if _, err := os.Stat(tt.filePath); os.IsNotExist(err) {
				t.Skip("Sample file not found, skipping test")
			}

			result, err := DetectFileType(tt.filePath)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}