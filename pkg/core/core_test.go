package core

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestDefaultResultCombiner(t *testing.T) {
	combiner := NewDefaultResultCombiner()

	t.Run("CombineEmptyResults", func(t *testing.T) {
		options := CombineOptions{Source: "test"}
		result, err := combiner.Combine([]*UnifiedTestResult{}, options)
		if err == nil {
			t.Error("Expected error for empty results")
		}
		if result != nil {
			t.Error("Expected nil result for empty input")
		}
	})

	t.Run("CombineSingleResult", func(t *testing.T) {
		singleResult := &UnifiedTestResult{
			Metadata: TestMetadata{Source: "test"},
			Summary:  TestSummary{Total: 1},
		}
		options := CombineOptions{Source: "combined"}
		result, err := combiner.Combine([]*UnifiedTestResult{singleResult}, options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}
		if result != singleResult {
			t.Error("Single result should be returned as-is")
		}
	})

	t.Run("CombineMultipleResults", func(t *testing.T) {
		result1 := &UnifiedTestResult{
			Metadata: TestMetadata{
				Source:    "junit",
				Framework: "JUnit 5",
				Environment: map[string]string{
					"java_version": "11",
				},
			},
			TestSuites: []TestSuite{
				{Name: "Suite1", Tests: []TestCase{{Name: "Test1", Status: StatusPassed}}},
			},
			Summary: TestSummary{Total: 1, Passed: 1},
		}
		result2 := &UnifiedTestResult{
			Metadata: TestMetadata{
				Source:    "gotest",
				Framework: "Go test",
				Environment: map[string]string{
					"go_version": "1.21",
				},
			},
			TestSuites: []TestSuite{
				{Name: "Suite2", Tests: []TestCase{{Name: "Test2", Status: StatusFailed}}},
			},
			Summary: TestSummary{Total: 1, Failed: 1},
		}

		options := CombineOptions{
			Source:           "combined-test",
			MergeEnvironment: true,
		}
		result, err := combiner.Combine([]*UnifiedTestResult{result1, result2}, options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if result.Metadata.Source != "combined-test" {
			t.Errorf("Expected source 'combined-test', got '%s'", result.Metadata.Source)
		}

		if result.Summary.Total != 2 {
			t.Errorf("Expected total 2, got %d", result.Summary.Total)
		}

		if result.Summary.Passed != 1 {
			t.Errorf("Expected passed 1, got %d", result.Summary.Passed)
		}

		if result.Summary.Failed != 1 {
			t.Errorf("Expected failed 1, got %d", result.Summary.Failed)
		}

		if len(result.TestSuites) != 2 {
			t.Errorf("Expected 2 test suites, got %d", len(result.TestSuites))
		}

		// Check environment merging
		if result.Metadata.Environment["java_version"] != "11" {
			t.Error("Environment should include java_version")
		}
		if result.Metadata.Environment["go_version"] != "1.21" {
			t.Error("Environment should include go_version")
		}
	})
}

func TestDefaultTypeDetector(t *testing.T) {
	detector := NewDefaultTypeDetector()

	t.Run("SupportedExtensions", func(t *testing.T) {
		extensions := detector.SupportedExtensions()
		if extensions[".xml"] != "xml" {
			t.Error("Expected .xml extension to map to xml")
		}
		if extensions[".json"] != "json" {
			t.Error("Expected .json extension to map to json")
		}
		if extensions[".tap"] != "tap" {
			t.Error("Expected .tap extension to map to tap")
		}
	})

	t.Run("DetectFromContent", func(t *testing.T) {
		tests := []struct {
			name     string
			content  string
			expected string
		}{
			{
				name:     "TAP format",
				content:  "TAP version 13\n1..2\nok 1 - test passes\nnot ok 2 - test fails",
				expected: "tap",
			},
			{
				name:     "JUnit XML",
				content:  "<?xml version=\"1.0\"?>\n<testsuite name=\"test\">",
				expected: "junit",
			},
			{
				name:     "TestNG XML",
				content:  "<?xml version=\"1.0\"?>\n<testng-results total=\"5\">",
				expected: "testng",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				reader := strings.NewReader(tt.content)
				result, err := detector.DetectType(reader)
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if result != tt.expected {
					t.Errorf("Expected %s, got %s", tt.expected, result)
				}
			})
		}
	})
}

func TestDefaultExporter(t *testing.T) {
	exporter := NewDefaultExporter()
	testResult := &UnifiedTestResult{
		Metadata: TestMetadata{
			Source:    "test",
			Timestamp: time.Now(),
			Version:   "2.0",
		},
		TestSuites: []TestSuite{
			{
				Name: "TestSuite",
				Tests: []TestCase{
					{Name: "Test1", Status: StatusPassed, Time: 1.0},
				},
			},
		},
		Summary: TestSummary{Total: 1, Passed: 1},
	}

	t.Run("ExportJSON", func(t *testing.T) {
		var buf bytes.Buffer
		options := ExportOptions{
			Format:          "json",
			Indent:          "  ",
			IncludeMetadata: true,
		}
		err := exporter.Export(testResult, &buf, options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "\"metadata\"") {
			t.Error("JSON output should contain metadata")
		}
		if !strings.Contains(output, "\"test_suites\"") {
			t.Error("JSON output should contain test_suites")
		}
	})

	t.Run("ExportYAML", func(t *testing.T) {
		var buf bytes.Buffer
		options := ExportOptions{
			Format:          "yaml",
			IncludeMetadata: true,
		}
		err := exporter.Export(testResult, &buf, options)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		output := buf.String()
		if !strings.Contains(output, "metadata:") {
			t.Error("YAML output should contain metadata")
		}
		if !strings.Contains(output, "test_suites:") {
			t.Error("YAML output should contain test_suites")
		}
	})

	t.Run("UnsupportedFormat", func(t *testing.T) {
		var buf bytes.Buffer
		options := ExportOptions{Format: "unsupported"}
		err := exporter.Export(testResult, &buf, options)
		if err == nil {
			t.Error("Expected error for unsupported format")
		}
	})
}

func TestTconvError(t *testing.T) {
	t.Run("BasicError", func(t *testing.T) {
		err := NewTconvError("test_type", "test message")
		if err.Type != "test_type" {
			t.Errorf("Expected type 'test_type', got '%s'", err.Type)
		}
		if err.Message != "test message" {
			t.Errorf("Expected message 'test message', got '%s'", err.Message)
		}
	})

	t.Run("ErrorWithContext", func(t *testing.T) {
		err := NewTconvError("test_type", "test message")
		err.WithContext("key", "value")
		if err.Context["key"] != "value" {
			t.Error("Context should be set")
		}
	})

	t.Run("ErrorWithSuggestion", func(t *testing.T) {
		err := NewTconvError("test_type", "test message")
		err.WithSuggestion("try this")
		if len(err.Suggestions) != 1 || err.Suggestions[0] != "try this" {
			t.Error("Suggestion should be added")
		}
	})
}

func TestTestStatus(t *testing.T) {
	t.Run("ValidStatuses", func(t *testing.T) {
		validStatuses := []TestStatus{StatusPassed, StatusFailed, StatusSkipped, StatusError}
		for _, status := range validStatuses {
			if !status.IsValid() {
				t.Errorf("Status %s should be valid", status)
			}
		}
	})

	t.Run("InvalidStatus", func(t *testing.T) {
		invalidStatus := TestStatus("invalid")
		if invalidStatus.IsValid() {
			t.Error("Invalid status should not be valid")
		}
	})

	t.Run("StringRepresentation", func(t *testing.T) {
		if StatusPassed.String() != "passed" {
			t.Errorf("Expected 'passed', got '%s'", StatusPassed.String())
		}
	})
}

func TestTestSummary(t *testing.T) {
	t.Run("PassRate", func(t *testing.T) {
		summary := &TestSummary{Total: 10, Passed: 8}
		expected := 80.0
		actual := summary.PassRate()
		if actual != expected {
			t.Errorf("Expected pass rate %.1f, got %.1f", expected, actual)
		}
	})

	t.Run("PassRateZeroTotal", func(t *testing.T) {
		summary := &TestSummary{Total: 0, Passed: 0}
		expected := 0.0
		actual := summary.PassRate()
		if actual != expected {
			t.Errorf("Expected pass rate %.1f for zero total, got %.1f", expected, actual)
		}
	})
}