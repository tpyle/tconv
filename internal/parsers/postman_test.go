package parsers

import (
	"testing"

	"github.com/tpyle/tconv/pkg/models"
)

func TestParsePostman(t *testing.T) {
	result, err := ParsePostman("../../testdata/sample_postman.json")
	if err != nil {
		t.Fatalf("Failed to parse Postman file: %v", err)
	}

	if result.Metadata.Source != "postman" {
		t.Errorf("Expected source to be 'postman', got '%s'", result.Metadata.Source)
	}

	if result.Metadata.Framework != "Newman" {
		t.Errorf("Expected framework to be 'Newman', got '%s'", result.Metadata.Framework)
	}

	if len(result.TestSuites) != 1 {
		t.Errorf("Expected 1 test suite, got %d", len(result.TestSuites))
	}

	suite := result.TestSuites[0]
	if suite.Name != "Test Collection" {
		t.Errorf("Expected suite name 'Test Collection', got '%s'", suite.Name)
	}

	if suite.Package != "postman" {
		t.Errorf("Expected package 'postman', got '%s'", suite.Package)
	}

	if len(suite.Tests) != 2 {
		t.Errorf("Expected 2 test cases, got %d", len(suite.Tests))
	}

	testCases := make(map[string]models.TestCase)
	for _, tc := range suite.Tests {
		testCases[tc.Name] = tc
	}

	if tc, exists := testCases["Get Users"]; exists {
		if tc.Status != models.StatusPassed {
			t.Errorf("Expected 'Get Users' to be passed, got %s", tc.Status)
		}
		if tc.ClassName != "postman.request" {
			t.Errorf("Expected className 'postman.request', got '%s'", tc.ClassName)
		}
		if tc.Time != 0.15 {
			t.Errorf("Expected time 0.15, got %f", tc.Time)
		}
		if tc.Properties["method"] != "GET" {
			t.Errorf("Expected method 'GET', got '%s'", tc.Properties["method"])
		}
		if tc.Properties["status_code"] != "200" {
			t.Errorf("Expected status_code '200', got '%s'", tc.Properties["status_code"])
		}
		if len(tc.Assertions) != 2 {
			t.Errorf("Expected 2 assertions, got %d", len(tc.Assertions))
		}
		for _, assertion := range tc.Assertions {
			if !assertion.Passed {
				t.Errorf("Expected assertion '%s' to pass", assertion.Name)
			}
		}
	} else {
		t.Error("'Get Users' test case not found")
	}

	if tc, exists := testCases["Create User"]; exists {
		if tc.Status != models.StatusFailed {
			t.Errorf("Expected 'Create User' to be failed, got %s", tc.Status)
		}
		if tc.Properties["method"] != "POST" {
			t.Errorf("Expected method 'POST', got '%s'", tc.Properties["method"])
		}
		if tc.Properties["status_code"] != "201" {
			t.Errorf("Expected status_code '201', got '%s'", tc.Properties["status_code"])
		}
		if len(tc.Assertions) != 2 {
			t.Errorf("Expected 2 assertions, got %d", len(tc.Assertions))
		}
		
		passedCount := 0
		failedCount := 0
		for _, assertion := range tc.Assertions {
			if assertion.Passed {
				passedCount++
			} else {
				failedCount++
			}
		}
		if passedCount != 1 || failedCount != 1 {
			t.Errorf("Expected 1 passed and 1 failed assertion, got %d passed, %d failed", passedCount, failedCount)
		}
	} else {
		t.Error("'Create User' test case not found")
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

	if result.Summary.Errors != 0 {
		t.Errorf("Expected errors 0, got %d", result.Summary.Errors)
	}

	if result.Summary.Skipped != 0 {
		t.Errorf("Expected skipped 0, got %d", result.Summary.Skipped)
	}
}

func TestParsePostmanNonExistentFile(t *testing.T) {
	_, err := ParsePostman("nonexistent.json")
	if err == nil {
		t.Error("Expected error for nonexistent file")
	}
}