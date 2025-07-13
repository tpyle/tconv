package parsers

import (
	"testing"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func TestParsePytest(t *testing.T) {
	result, err := ParsePytest("../../testdata/sample_pytest.json")
	if err != nil {
		t.Fatalf("Failed to parse pytest file: %v", err)
	}

	if result.Metadata.Source != "pytest" {
		t.Errorf("Expected source to be 'pytest', got '%s'", result.Metadata.Source)
	}

	if result.Metadata.Framework != "pytest" {
		t.Errorf("Expected framework to be 'pytest', got '%s'", result.Metadata.Framework)
	}

	if result.Metadata.Environment["pytest_version"] != "7.1.2" {
		t.Errorf("Expected pytest_version '7.1.2', got '%s'", result.Metadata.Environment["pytest_version"])
	}

	if len(result.TestSuites) != 1 {
		t.Errorf("Expected 1 test suite, got %d", len(result.TestSuites))
	}

	suite := result.TestSuites[0]
	if suite.Name != "test_math.py" {
		t.Errorf("Expected suite name 'test_math.py', got '%s'", suite.Name)
	}

	if suite.Package != "tests/test_math.py" {
		t.Errorf("Expected package 'tests/test_math.py', got '%s'", suite.Package)
	}

	if len(suite.Tests) != 3 {
		t.Errorf("Expected 3 test cases, got %d", len(suite.Tests))
	}

	testCases := make(map[string]models.TestCase)
	for _, tc := range suite.Tests {
		testCases[tc.Name] = tc
	}

	if tc, exists := testCases["test_addition"]; exists {
		if tc.Status != models.StatusPassed {
			t.Errorf("Expected 'test_addition' to be passed, got %s", tc.Status)
		}
		if tc.ClassName != "tests/test_math.py" {
			t.Errorf("Expected className 'tests/test_math.py', got '%s'", tc.ClassName)
		}
		expectedTime := time.Duration((0.001 + 0.045 + 0.001) * float64(time.Second)) // setup + call + teardown
		if tc.Time != expectedTime {
			t.Errorf("Expected time %v, got %v", expectedTime, tc.Time)
		}
		if tc.Properties["nodeid"] != "tests/test_math.py::test_addition" {
			t.Errorf("Expected nodeid 'tests/test_math.py::test_addition', got '%s'", tc.Properties["nodeid"])
		}
		if tc.Properties["keywords"] != "test_addition,math,unit" {
			t.Errorf("Expected keywords 'test_addition,math,unit', got '%s'", tc.Properties["keywords"])
		}
		if tc.SystemOut != "Call stdout:\nRunning addition test\n" {
			t.Errorf("Expected SystemOut 'Call stdout:\\nRunning addition test\\n', got '%s'", tc.SystemOut)
		}
	} else {
		t.Error("'test_addition' not found")
	}

	if tc, exists := testCases["test_subtraction"]; exists {
		if tc.Status != models.StatusFailed {
			t.Errorf("Expected 'test_subtraction' to be failed, got %s", tc.Status)
		}
		if tc.Message != "assert 2 == 3" {
			t.Errorf("Expected failure message 'assert 2 == 3', got '%s'", tc.Message)
		}
		if tc.Details == "" {
			t.Error("Expected details to contain traceback information")
		}
		if tc.SystemErr != "Call stderr:\nAssertionError: Expected 2, got 3\n" {
			t.Errorf("Expected SystemErr to contain stderr, got '%s'", tc.SystemErr)
		}
	} else {
		t.Error("'test_subtraction' not found")
	}

	if tc, exists := testCases["test_division"]; exists {
		if tc.Status != models.StatusSkipped {
			t.Errorf("Expected 'test_division' to be skipped, got %s", tc.Status)
		}
		if tc.Message != "Division test disabled pending fix" {
			t.Errorf("Expected skip message 'Division test disabled pending fix', got '%s'", tc.Message)
		}
	} else {
		t.Error("'test_division' not found")
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

	if result.Summary.Skipped != 1 {
		t.Errorf("Expected skipped 1, got %d", result.Summary.Skipped)
	}

	if result.Summary.Duration != 1234*time.Millisecond {
		t.Errorf("Expected duration 1.234s, got %v", result.Summary.Duration)
	}
}

func TestParsePytestNonExistentFile(t *testing.T) {
	_, err := ParsePytest("nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}