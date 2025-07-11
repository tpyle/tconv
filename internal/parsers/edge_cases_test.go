package parsers

import (
	"os"
	"testing"
)

func TestParseJUnitMalformedXML(t *testing.T) {
	// Create temporary malformed XML file
	tempFile, err := os.CreateTemp("", "malformed_junit_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	malformedXML := `<?xml version="1.0"?>
<testsuites>
  <testsuite name="Test">
    <testcase name="test1"
  </testsuite>
</testsuites>`
	
	if _, err := tempFile.WriteString(malformedXML); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	_, err = ParseJUnit(tempFile.Name())
	if err == nil {
		t.Error("Expected error for malformed XML")
	}
}

func TestParsePostmanMalformedJSON(t *testing.T) {
	// Create temporary malformed JSON file
	tempFile, err := os.CreateTemp("", "malformed_postman_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	malformedJSON := `{"info": {"name": "test", "run": {`
	
	if _, err := tempFile.WriteString(malformedJSON); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	_, err = ParsePostman(tempFile.Name())
	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestParseGoTestMalformedJSON(t *testing.T) {
	// Create temporary malformed JSON file
	tempFile, err := os.CreateTemp("", "malformed_gotest_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	malformedJSON := `{"Time": "invalid-time", "Action": "run"`
	
	if _, err := tempFile.WriteString(malformedJSON); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	// Should not error on malformed lines, should skip them
	result, err := ParseGoTest(tempFile.Name())
	if err != nil {
		t.Errorf("Should handle malformed lines gracefully: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
}

func TestParseTAPEmptyFile(t *testing.T) {
	// Create temporary empty file
	tempFile, err := os.CreateTemp("", "empty_*.tap")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	tempFile.Close()
	
	result, err := ParseTAP(tempFile.Name())
	if err != nil {
		t.Errorf("Should handle empty file gracefully: %v", err)
	}
	if result == nil {
		t.Error("Result should not be nil")
	}
	if len(result.TestSuites) != 1 || len(result.TestSuites[0].Tests) != 0 {
		t.Error("Expected empty test suite")
	}
}

func TestParseXUnitMalformedXML(t *testing.T) {
	// Create temporary malformed XML file
	tempFile, err := os.CreateTemp("", "malformed_xunit_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	malformedXML := `<?xml version="1.0"?>
<assemblies>
  <assembly name="Test">
    <collection name="TestCollection"
  </assembly>
</assemblies>`
	
	if _, err := tempFile.WriteString(malformedXML); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	_, err = ParseXUnit(tempFile.Name())
	if err == nil {
		t.Error("Expected error for malformed XML")
	}
}

func TestParsePytestMalformedJSON(t *testing.T) {
	// Create temporary malformed JSON file
	tempFile, err := os.CreateTemp("", "malformed_pytest_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	malformedJSON := `{"report": {"pytest_version": "7.1.2"}, "tests": [{"nodeid": "test", "outcome"`
	
	if _, err := tempFile.WriteString(malformedJSON); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	_, err = ParsePytest(tempFile.Name())
	if err == nil {
		t.Error("Expected error for malformed JSON")
	}
}

func TestParseTestNGMalformedXML(t *testing.T) {
	// Create temporary malformed XML file
	tempFile, err := os.CreateTemp("", "malformed_testng_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	malformedXML := `<?xml version="1.0"?>
<testng-results>
  <suite name="Test">
    <test name="TestName"
  </suite>
</testng-results>`
	
	if _, err := tempFile.WriteString(malformedXML); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	_, err = ParseTestNG(tempFile.Name())
	if err == nil {
		t.Error("Expected error for malformed XML")
	}
}

func TestParseTAPWithBailOut(t *testing.T) {
	// Create temporary TAP file with bail out
	tempFile, err := os.CreateTemp("", "bailout_*.tap")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	tapContent := `TAP version 13
1..3
ok 1 - First test
Bail out! Something went wrong
not ok 2 - This test should not run`
	
	if _, err := tempFile.WriteString(tapContent); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	result, err := ParseTAP(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse TAP with bail out: %v", err)
	}
	
	// Should have the first test plus a bail out test
	if len(result.TestSuites[0].Tests) != 3 {
		t.Errorf("Expected 3 tests (including bail out), got %d", len(result.TestSuites[0].Tests))
	}
	
	// Check for bail out test
	found := false
	for _, test := range result.TestSuites[0].Tests {
		if test.Name == "Bail Out" {
			found = true
			if test.Status != "error" {
				t.Errorf("Expected bail out test to have error status, got %s", test.Status)
			}
		}
	}
	if !found {
		t.Error("Expected to find bail out test")
	}
}

func TestParseJUnitWithSystemOutput(t *testing.T) {
	// Create temporary JUnit file with system output
	tempFile, err := os.CreateTemp("", "junit_sysout_*.xml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	
	junitXML := `<?xml version="1.0"?>
<testsuite name="TestSuite" tests="1" failures="0" errors="0" time="0.5">
  <testcase name="testWithOutput" classname="TestClass" time="0.5">
    <system-out>Standard output content</system-out>
    <system-err>Error output content</system-err>
  </testcase>
</testsuite>`
	
	if _, err := tempFile.WriteString(junitXML); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}
	tempFile.Close()
	
	result, err := ParseJUnit(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to parse JUnit with system output: %v", err)
	}
	
	if len(result.TestSuites[0].Tests) != 1 {
		t.Errorf("Expected 1 test, got %d", len(result.TestSuites[0].Tests))
	}
	
	test := result.TestSuites[0].Tests[0]
	if test.SystemOut != "Standard output content" {
		t.Errorf("Expected system out, got '%s'", test.SystemOut)
	}
	if test.SystemErr != "Error output content" {
		t.Errorf("Expected system err, got '%s'", test.SystemErr)
	}
}