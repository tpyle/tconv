package parsers

import (
	"testing"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func TestParseJUnit(t *testing.T) {
	result, err := ParseJUnit("../../testdata/sample_junit.xml")
	if err != nil {
		t.Fatalf("Failed to parse JUnit file: %v", err)
	}

	if result.Metadata.Source != "junit" {
		t.Errorf("Expected source to be 'junit', got '%s'", result.Metadata.Source)
	}

	if len(result.TestSuites) != 1 {
		t.Errorf("Expected 1 test suite, got %d", len(result.TestSuites))
	}

	suite := result.TestSuites[0]
	if suite.Name != "TestSuite1" {
		t.Errorf("Expected suite name 'TestSuite1', got '%s'", suite.Name)
	}

	if suite.Package != "com.example.tests" {
		t.Errorf("Expected package 'com.example.tests', got '%s'", suite.Package)
	}

	if len(suite.Tests) != 4 {
		t.Errorf("Expected 4 test cases, got %d", len(suite.Tests))
	}

	if suite.Failures != 1 {
		t.Errorf("Expected 1 failure, got %d", suite.Failures)
	}

	if suite.Errors != 1 {
		t.Errorf("Expected 1 error, got %d", suite.Errors)
	}

	if suite.Skipped != 0 {
		t.Errorf("Expected 0 skipped, got %d", suite.Skipped)
	}

	testCases := make(map[string]models.TestCase)
	for _, tc := range suite.Tests {
		testCases[tc.Name] = tc
	}

	if tc, exists := testCases["testSuccess"]; exists {
		if tc.Status != models.StatusPassed {
			t.Errorf("Expected testSuccess to be passed, got %s", tc.Status)
		}
		if tc.Time != 100*time.Millisecond {
			t.Errorf("Expected testSuccess time to be 100ms, got %v", tc.Time)
		}
	} else {
		t.Error("testSuccess not found")
	}

	if tc, exists := testCases["testFailure"]; exists {
		if tc.Status != models.StatusFailed {
			t.Errorf("Expected testFailure to be failed, got %s", tc.Status)
		}
		if tc.Message != "Expected 5 but was 4" {
			t.Errorf("Expected failure message 'Expected 5 but was 4', got '%s'", tc.Message)
		}
		if tc.ErrorType != "AssertionError" {
			t.Errorf("Expected error type 'AssertionError', got '%s'", tc.ErrorType)
		}
	} else {
		t.Error("testFailure not found")
	}

	if tc, exists := testCases["testError"]; exists {
		if tc.Status != models.StatusError {
			t.Errorf("Expected testError to be error, got %s", tc.Status)
		}
		if tc.Message != "Null pointer exception" {
			t.Errorf("Expected error message 'Null pointer exception', got '%s'", tc.Message)
		}
	} else {
		t.Error("testError not found")
	}

	if tc, exists := testCases["testSkipped"]; exists {
		if tc.Status != models.StatusSkipped {
			t.Errorf("Expected testSkipped to be skipped, got %s", tc.Status)
		}
		if tc.Message != "Test disabled" {
			t.Errorf("Expected skip message 'Test disabled', got '%s'", tc.Message)
		}
	} else {
		t.Error("testSkipped not found")
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

	if result.Summary.Errors != 1 {
		t.Errorf("Expected errors 1, got %d", result.Summary.Errors)
	}

	if result.Summary.Skipped != 1 {
		t.Errorf("Expected skipped 1, got %d", result.Summary.Skipped)
	}
}

func TestParseJUnitNonExistentFile(t *testing.T) {
	_, err := ParseJUnit("nonexistent.xml")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}