package parsers

import (
	"testing"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func TestParseGoTest(t *testing.T) {
	result, err := ParseGoTest("../../testdata/sample_gotest.json")
	if err != nil {
		t.Fatalf("Failed to parse Go test file: %v", err)
	}

	if result.Metadata.Source != "gotest" {
		t.Errorf("Expected source to be 'gotest', got '%s'", result.Metadata.Source)
	}

	if result.Metadata.Framework != "go test" {
		t.Errorf("Expected framework to be 'go test', got '%s'", result.Metadata.Framework)
	}

	if len(result.TestSuites) != 1 {
		t.Errorf("Expected 1 test suite, got %d", len(result.TestSuites))
	}

	suite := result.TestSuites[0]
	if suite.Name != "github.com/example/myapp" {
		t.Errorf("Expected suite name 'github.com/example/myapp', got '%s'", suite.Name)
	}

	if suite.Package != "github.com/example/myapp" {
		t.Errorf("Expected package 'github.com/example/myapp', got '%s'", suite.Package)
	}

	if len(suite.Tests) != 3 {
		t.Errorf("Expected 3 test cases, got %d", len(suite.Tests))
	}

	testCases := make(map[string]models.TestCase)
	for _, tc := range suite.Tests {
		testCases[tc.Name] = tc
	}

	if tc, exists := testCases["TestAddition"]; exists {
		if tc.Status != models.StatusPassed {
			t.Errorf("Expected TestAddition to be passed, got %s", tc.Status)
		}
		if tc.ClassName != "github.com/example/myapp" {
			t.Errorf("Expected className 'github.com/example/myapp', got '%s'", tc.ClassName)
		}
		if tc.Time != 150*time.Millisecond {
			t.Errorf("Expected time 150ms, got %v", tc.Time)
		}
		expectedOutput := "=== RUN   TestAddition\n    math_test.go:10: Testing addition\n"
		if tc.SystemOut != expectedOutput {
			t.Errorf("Expected SystemOut '%s', got '%s'", expectedOutput, tc.SystemOut)
		}
	} else {
		t.Error("TestAddition not found")
	}

	if tc, exists := testCases["TestSubtraction"]; exists {
		if tc.Status != models.StatusFailed {
			t.Errorf("Expected TestSubtraction to be failed, got %s", tc.Status)
		}
		if tc.Time != 200*time.Millisecond {
			t.Errorf("Expected time 200ms, got %v", tc.Time)
		}
		if tc.Message != "Test failed" {
			t.Errorf("Expected message 'Test failed', got '%s'", tc.Message)
		}
		expectedOutput := "=== RUN   TestSubtraction\n    math_test.go:20: Testing subtraction\n    math_test.go:22: Expected 2, got 3"
		if tc.Details != expectedOutput {
			t.Errorf("Expected Details '%s', got '%s'", expectedOutput, tc.Details)
		}
	} else {
		t.Error("TestSubtraction not found")
	}

	if tc, exists := testCases["TestDivision"]; exists {
		if tc.Status != models.StatusSkipped {
			t.Errorf("Expected TestDivision to be skipped, got %s", tc.Status)
		}
		if tc.Time != 100*time.Millisecond {
			t.Errorf("Expected time 100ms, got %v", tc.Time)
		}
		if tc.Message != "Test skipped" {
			t.Errorf("Expected message 'Test skipped', got '%s'", tc.Message)
		}
		expectedOutput := "=== RUN   TestDivision\n    math_test.go:30: Skipping division test"
		if tc.Details != expectedOutput {
			t.Errorf("Expected Details '%s', got '%s'", expectedOutput, tc.Details)
		}
	} else {
		t.Error("TestDivision not found")
	}

	if result.Summary.Total != 3 {
		t.Errorf("Expected total 3, got %d", result.Summary.Total)
	}

	if result.Summary.Passed != 1 {
		t.Errorf("Expected passed 1, got %d", result.Summary.Passed)
	}

	if result.Summary.Failed != 1 {
		t.Errorf("Expected failed 1, got %d", result.Summary.Failed)
	}

	if result.Summary.Errors != 0 {
		t.Errorf("Expected errors 0, got %d", result.Summary.Errors)
	}

	if result.Summary.Skipped != 1 {
		t.Errorf("Expected skipped 1, got %d", result.Summary.Skipped)
	}
}

func TestParseGoTestNonExistentFile(t *testing.T) {
	_, err := ParseGoTest("nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}