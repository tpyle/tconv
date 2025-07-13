package parsers

import (
	"testing"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func TestParseTestNG(t *testing.T) {
	result, err := ParseTestNG("../../testdata/sample_testng.xml")
	if err != nil {
		t.Fatalf("Failed to parse TestNG file: %v", err)
	}

	if result.Metadata.Source != "testng" {
		t.Errorf("Expected source to be 'testng', got '%s'", result.Metadata.Source)
	}

	if result.Metadata.Framework != "TestNG" {
		t.Errorf("Expected framework to be 'TestNG', got '%s'", result.Metadata.Framework)
	}

	if len(result.TestSuites) != 1 {
		t.Errorf("Expected 1 test suite, got %d", len(result.TestSuites))
	}

	suite := result.TestSuites[0]
	if suite.Name != "MathTests" {
		t.Errorf("Expected suite name 'MathTests', got '%s'", suite.Name)
	}

	if suite.Package != "TestSuite" {
		t.Errorf("Expected package 'TestSuite', got '%s'", suite.Package)
	}

	if len(suite.Tests) != 4 {
		t.Errorf("Expected 4 test cases, got %d", len(suite.Tests))
	}

	testCases := make(map[string]models.TestCase)
	for _, tc := range suite.Tests {
		testCases[tc.Name] = tc
	}

	if tc, exists := testCases["testAddition"]; exists {
		if tc.Status != models.StatusPassed {
			t.Errorf("Expected 'testAddition' to be passed, got %s", tc.Status)
		}
		if tc.ClassName != "TestClass" {
			t.Errorf("Expected className 'TestClass', got '%s'", tc.ClassName)
		}
		if tc.Time != 150*time.Millisecond {
			t.Errorf("Expected time 150ms, got %v", tc.Time)
		}
		if tc.Properties["signature"] != "TestClass.testAddition()" {
			t.Errorf("Expected signature 'TestClass.testAddition()', got '%s'", tc.Properties["signature"])
		}
		if tc.Properties["groups"] != "unit" {
			t.Errorf("Expected groups 'unit', got '%s'", tc.Properties["groups"])
		}
		if tc.Properties["description"] != "Test basic addition" {
			t.Errorf("Expected description 'Test basic addition', got '%s'", tc.Properties["description"])
		}
		if tc.SystemOut != "Addition test passed successfully" {
			t.Errorf("Expected SystemOut 'Addition test passed successfully', got '%s'", tc.SystemOut)
		}
	} else {
		t.Error("'testAddition' not found")
	}

	if tc, exists := testCases["testSubtraction"]; exists {
		if tc.Status != models.StatusFailed {
			t.Errorf("Expected 'testSubtraction' to be failed, got %s", tc.Status)
		}
		if tc.Message != "Expected 5 but was 4" {
			t.Errorf("Expected failure message 'Expected 5 but was 4', got '%s'", tc.Message)
		}
		if tc.ErrorType != "java.lang.AssertionError" {
			t.Errorf("Expected error type 'java.lang.AssertionError', got '%s'", tc.ErrorType)
		}
		if tc.Details == "" {
			t.Error("Expected details to contain stack trace")
		}
	} else {
		t.Error("'testSubtraction' not found")
	}

	if tc, exists := testCases["testMultiplication"]; exists {
		if tc.Status != models.StatusPassed {
			t.Errorf("Expected 'testMultiplication' to be passed, got %s", tc.Status)
		}
		if tc.Properties["parameters"] != "param[0]=3, param[1]=4" {
			t.Errorf("Expected parameters 'param[0]=3, param[1]=4', got '%s'", tc.Properties["parameters"])
		}
	} else {
		t.Error("'testMultiplication' not found")
	}

	if tc, exists := testCases["testDivision"]; exists {
		if tc.Status != models.StatusSkipped {
			t.Errorf("Expected 'testDivision' to be skipped, got %s", tc.Status)
		}
		if tc.SystemOut != "Division test skipped - feature not implemented" {
			t.Errorf("Expected SystemOut 'Division test skipped - feature not implemented', got '%s'", tc.SystemOut)
		}
	} else {
		t.Error("'testDivision' not found")
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

func TestParseTestNGNonExistentFile(t *testing.T) {
	_, err := ParseTestNG("nonexistent.xml")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}