package tconv

import (
	"os"
	"testing"
)

func TestPublicAPI(t *testing.T) {
	// Test the public API functions
	
	t.Run("SupportedFormats", func(t *testing.T) {
		formats := SupportedFormats()
		if len(formats) == 0 {
			t.Error("Should return supported formats")
		}
		
		expectedFormats := map[string]bool{
			"junit": true, "testng": true, "xunit": true,
			"gotest": true, "postman": true, "pytest": true, "tap": true,
		}
		
		for _, format := range formats {
			if !expectedFormats[format] {
				t.Errorf("Unexpected format: %s", format)
			}
		}
	})
	
	t.Run("DetectFileType", func(t *testing.T) {
		// Test with a real sample file
		fileType, err := DetectFileType("testdata/sample_junit.xml")
		if err != nil {
			t.Fatalf("Detection failed: %v", err)
		}
		if fileType != "junit" {
			t.Errorf("Expected junit, got %s", fileType)
		}
	})
	
	t.Run("Convert", func(t *testing.T) {
		tempDir := t.TempDir()
		outputFile := tempDir + "/output.json"
		
		err := Convert("testdata/sample_junit.xml", "junit", outputFile)
		if err != nil {
			t.Fatalf("Conversion failed: %v", err)
		}
		
		// Verify output file exists and has content
		content, err := os.ReadFile(outputFile)
		if err != nil {
			t.Fatalf("Failed to read output: %v", err)
		}
		if len(content) == 0 {
			t.Error("Output file is empty")
		}
	})
}