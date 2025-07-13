package parsers

import (
	"testing"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func TestParseXUnit(t *testing.T) {
	result, err := ParseXUnit("../../testdata/sample_xunit.xml")
	if err != nil {
		t.Fatalf("Failed to parse xUnit file: %v", err)
	}

	if result.Metadata.Source != "xunit" {
		t.Errorf("Expected source to be 'xunit', got '%s'", result.Metadata.Source)
	}

	if result.Metadata.Framework != "xUnit" {
		t.Errorf("Expected framework to be 'xUnit', got '%s'", result.Metadata.Framework)
	}

	if len(result.TestSuites) != 1 {
		t.Errorf("Expected 1 test suite, got %d", len(result.TestSuites))
	}

	suite := result.TestSuites[0]
	if suite.Name != "TestCollection" {
		t.Errorf("Expected suite name 'TestCollection', got '%s'", suite.Name)
	}

	if suite.Package != "TestAssembly" {
		t.Errorf("Expected package 'TestAssembly', got '%s'", suite.Package)
	}

	if len(suite.Tests) != 4 {
		t.Errorf("Expected 4 test cases, got %d", len(suite.Tests))
	}

	testCases := make(map[string]models.TestCase)
	for _, tc := range suite.Tests {
		testCases[tc.Name] = tc
	}

	if tc, exists := testCases["TestClass.TestSuccess"]; exists {
		if tc.Status != models.StatusPassed {
			t.Errorf("Expected 'TestClass.TestSuccess' to be passed, got %s", tc.Status)
		}
		if tc.ClassName != "TestClass" {
			t.Errorf("Expected className 'TestClass', got '%s'", tc.ClassName)
		}
		if tc.Time != 100*time.Millisecond {
			t.Errorf("Expected time 100ms, got %v", tc.Time)
		}
		if tc.SystemOut != "Test output from TestSuccess" {
			t.Errorf("Expected SystemOut 'Test output from TestSuccess', got '%s'", tc.SystemOut)
		}
		if tc.Properties["method"] != "TestSuccess" {
			t.Errorf("Expected method 'TestSuccess', got '%s'", tc.Properties["method"])
		}
	} else {
		t.Error("'TestClass.TestSuccess' not found")
	}

	if tc, exists := testCases["TestClass.TestFailure"]; exists {
		if tc.Status != models.StatusFailed {
			t.Errorf("Expected 'TestClass.TestFailure' to be failed, got %s", tc.Status)
		}
		if tc.Message != "Expected: 5, Actual: 4" {
			t.Errorf("Expected failure message 'Expected: 5, Actual: 4', got '%s'", tc.Message)
		}
		if tc.ErrorType != "Xunit.Sdk.XunitException" {
			t.Errorf("Expected error type 'Xunit.Sdk.XunitException', got '%s'", tc.ErrorType)
		}
	} else {
		t.Error("'TestClass.TestFailure' not found")
	}

	if tc, exists := testCases["TestClass.TestSkipped"]; exists {
		if tc.Status != models.StatusSkipped {
			t.Errorf("Expected 'TestClass.TestSkipped' to be skipped, got %s", tc.Status)
		}
		if tc.Message != "Test disabled for maintenance" {
			t.Errorf("Expected skip message 'Test disabled for maintenance', got '%s'", tc.Message)
		}
		if tc.Properties["trait_Category"] != "Integration" {
			t.Errorf("Expected trait Category 'Integration', got '%s'", tc.Properties["trait_Category"])
		}
		if tc.Properties["trait_Priority"] != "Low" {
			t.Errorf("Expected trait Priority 'Low', got '%s'", tc.Properties["trait_Priority"])
		}
	} else {
		t.Error("'TestClass.TestSkipped' not found")
	}

	if tc, exists := testCases["TestClass.TestWithTraits"]; exists {
		if tc.Status != models.StatusPassed {
			t.Errorf("Expected 'TestClass.TestWithTraits' to be passed, got %s", tc.Status)
		}
		if tc.Properties["trait_Category"] != "Unit" {
			t.Errorf("Expected trait Category 'Unit', got '%s'", tc.Properties["trait_Category"])
		}
		if tc.Properties["trait_Priority"] != "High" {
			t.Errorf("Expected trait Priority 'High', got '%s'", tc.Properties["trait_Priority"])
		}
	} else {
		t.Error("'TestClass.TestWithTraits' not found")
	}

	if result.Summary.Total != 4 {
		t.Errorf("Expected total 4, got %d", result.Summary.Total)
	}

	if result.Summary.Passed != 2 {
		t.Errorf("Expected passed 2, got %d", result.Summary.Passed)
	}

	if result.Summary.Failed != 1 {
		t.Errorf("Expected failed 1, got %d", result.Summary.Failed)
	}

	if result.Summary.Skipped != 1 {
		t.Errorf("Expected skipped 1, got %d", result.Summary.Skipped)
	}
}

func TestParseXUnitNonExistentFile(t *testing.T) {
	_, err := ParseXUnit("nonexistent.xml")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}