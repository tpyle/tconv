package tconv

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func TestExportFunctionality(t *testing.T) {
	// Create a sample tiki test result
	result := &models.TikiTestResult{
		Metadata: models.TestMetadata{
			Source:    "test",
			Timestamp: time.Now(),
			Version:   "2.0",
			Framework: "Test Framework",
		},
		TestSuites: []models.TestSuite{
			{
				Name: "TestSuite1",
				Tests: []models.TestCase{
					{
						Name:   "TestPassed",
						Status: models.StatusPassed,
						Time:   1500 * time.Millisecond,
					},
					{
						Name:    "TestFailed",
						Status:  models.StatusFailed,
						Time:    800 * time.Millisecond,
						Message: "Test failed",
						Details: "Assertion failed",
					},
					{
						Name:    "TestSkipped",
						Status:  models.StatusSkipped,
						Message: "Test skipped",
					},
				},
				Errors:   0,
				Failures: 1,
				Skipped:  1,
				Time:     2300 * time.Millisecond,
			},
		},
		Summary: models.TestSummary{
			Total:    3,
			Passed:   1,
			Failed:   1,
			Skipped:  1,
			Duration: 2300 * time.Millisecond,
		},
	}

	t.Run("SupportedExportFormats", func(t *testing.T) {
		formats := SupportedExportFormats()
		if len(formats) == 0 {
			t.Error("Should return supported export formats")
		}

		// Check for expected formats
		expectedFormats := []string{"junit", "tap", "gotest"}
		found := make(map[string]bool)
		for _, format := range formats {
			found[format] = true
		}

		for _, expected := range expectedFormats {
			if !found[expected] {
				t.Errorf("Expected format %s not found in supported formats", expected)
			}
		}
	})

	t.Run("ExportToJUnit", func(t *testing.T) {
		var buf bytes.Buffer
		err := ExportToWriter(result, "junit", &buf)
		if err != nil {
			t.Fatalf("Export to JUnit failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "<?xml") {
			t.Error("JUnit output should contain XML header")
		}
		if !strings.Contains(output, "<testsuites") {
			t.Error("JUnit output should contain testsuites element")
		}
		if !strings.Contains(output, "TestPassed") {
			t.Error("JUnit output should contain test names")
		}
	})

	t.Run("ExportToTAP", func(t *testing.T) {
		var buf bytes.Buffer
		err := ExportToWriter(result, "tap", &buf)
		if err != nil {
			t.Fatalf("Export to TAP failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "TAP version 13") {
			t.Error("TAP output should contain version header")
		}
		if !strings.Contains(output, "1..3") {
			t.Error("TAP output should contain plan")
		}
		if !strings.Contains(output, "ok 1") {
			t.Error("TAP output should contain test results")
		}
		if !strings.Contains(output, "not ok 2") {
			t.Error("TAP output should contain failed test")
		}
	})

	t.Run("ExportToGoTest", func(t *testing.T) {
		var buf bytes.Buffer
		err := ExportToWriter(result, "gotest", &buf)
		if err != nil {
			t.Fatalf("Export to Go test failed: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, `"Action":"run"`) {
			t.Error("Go test output should contain run actions")
		}
		if !strings.Contains(output, `"Action":"pass"`) {
			t.Error("Go test output should contain pass actions")
		}
		if !strings.Contains(output, `"Action":"fail"`) {
			t.Error("Go test output should contain fail actions")
		}
	})

	t.Run("UnsupportedFormat", func(t *testing.T) {
		var buf bytes.Buffer
		err := ExportToWriter(result, "unsupported", &buf)
		if err == nil {
			t.Error("Expected error for unsupported format")
		}
	})
}