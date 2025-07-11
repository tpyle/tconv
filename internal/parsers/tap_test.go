package parsers

import (
	"testing"

	"github.com/tpyle/tconv/pkg/models"
)

func TestParseTAP(t *testing.T) {
	result, err := ParseTAP("../../testdata/sample_tap.tap")
	if err != nil {
		t.Fatalf("Failed to parse TAP file: %v", err)
	}

	if result.Metadata.Source != "tap" {
		t.Errorf("Expected source to be 'tap', got '%s'", result.Metadata.Source)
	}

	if result.Metadata.Framework != "TAP" {
		t.Errorf("Expected framework to be 'TAP', got '%s'", result.Metadata.Framework)
	}

	if len(result.TestSuites) != 1 {
		t.Errorf("Expected 1 test suite, got %d", len(result.TestSuites))
	}

	suite := result.TestSuites[0]
	if suite.Name != "TAP Tests" {
		t.Errorf("Expected suite name 'TAP Tests', got '%s'", suite.Name)
	}

	if len(suite.Tests) != 4 {
		t.Errorf("Expected 4 test cases, got %d", len(suite.Tests))
	}

	testCases := make(map[string]models.TestCase)
	for _, tc := range suite.Tests {
		testCases[tc.Name] = tc
	}

	if tc, exists := testCases["Basic addition test"]; exists {
		if tc.Status != models.StatusPassed {
			t.Errorf("Expected 'Basic addition test' to be passed, got %s", tc.Status)
		}
		if tc.Properties["tap_number"] != "1" {
			t.Errorf("Expected tap_number '1', got '%s'", tc.Properties["tap_number"])
		}
	} else {
		t.Error("'Basic addition test' not found")
	}

	if tc, exists := testCases["Subtraction test"]; exists {
		if tc.Status != models.StatusFailed {
			t.Errorf("Expected 'Subtraction test' to be failed, got %s", tc.Status)
		}
		if tc.Properties["tap_number"] != "2" {
			t.Errorf("Expected tap_number '2', got '%s'", tc.Properties["tap_number"])
		}
		expectedDetails := "Expected: 5\nActual: 3"
		if tc.Details != expectedDetails {
			t.Errorf("Expected details '%s', got '%s'", expectedDetails, tc.Details)
		}
	} else {
		t.Error("'Subtraction test' not found")
	}

	if tc, exists := testCases["Multiplication test"]; exists {
		if tc.Status != models.StatusSkipped {
			t.Errorf("Expected 'Multiplication test' to be skipped, got %s", tc.Status)
		}
		if tc.Message != "Not implemented yet" {
			t.Errorf("Expected message 'Not implemented yet', got '%s'", tc.Message)
		}
	} else {
		t.Error("'Multiplication test' not found")
	}

	if tc, exists := testCases["Division test"]; exists {
		if tc.Status != models.StatusSkipped {
			t.Errorf("Expected 'Division test' to be skipped, got %s", tc.Status)
		}
		if tc.Message != "Fix division by zero handling" {
			t.Errorf("Expected message 'Fix division by zero handling', got '%s'", tc.Message)
		}
	} else {
		t.Error("'Division test' not found")
	}

	if result.Summary.Total != 4 {
		t.Errorf("Expected total 4, got %d", result.Summary.Total)
	}

	if result.Summary.Passed != 1 {
		t.Errorf("Expected passed 1, got %d", result.Summary.Passed)
	}

	if result.Summary.Failed != 1 {
		t.Errorf("Expected failed 1, got %d", result.Summary.Failed)
	}

	if result.Summary.Skipped != 2 {
		t.Errorf("Expected skipped 2, got %d", result.Summary.Skipped)
	}
}

func TestParseTAPNonExistentFile(t *testing.T) {
	_, err := ParseTAP("nonexistent.tap")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}