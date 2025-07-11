package models

import (
	"encoding/json"
	"encoding/xml"
	"testing"
	"time"
)

func TestUnifiedTestResultSerialization(t *testing.T) {
	result := UnifiedTestResult{
		Metadata: TestMetadata{
			Source:    "test",
			Timestamp: time.Now(),
			Version:   "1.0",
			Framework: "TestFramework",
			Environment: map[string]string{
				"key": "value",
			},
		},
		TestSuites: []TestSuite{
			{
				Name:    "TestSuite",
				Package: "test.package",
				Tests: []TestCase{
					{
						Name:      "TestCase",
						ClassName: "TestClass",
						Time:      1.5,
						Status:    StatusPassed,
						Message:   "Test passed",
						Details:   "Test details",
						ErrorType: "NoError",
						SystemOut: "Output",
						SystemErr: "Error",
						Assertions: []Assertion{
							{
								Name:     "AssertionName",
								Expected: "expected",
								Actual:   "actual",
								Passed:   true,
								Message:  "Assertion passed",
							},
						},
						Properties: map[string]string{
							"prop": "value",
						},
					},
				},
				Errors:   0,
				Failures: 0,
				Skipped:  0,
				Time:     1.5,
				Properties: map[string]string{
					"suite_prop": "value",
				},
			},
		},
		Summary: TestSummary{
			Total:    1,
			Passed:   1,
			Failed:   0,
			Skipped:  0,
			Errors:   0,
			Duration: 1.5,
		},
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	var unmarshaled UnifiedTestResult
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	if unmarshaled.Metadata.Source != result.Metadata.Source {
		t.Errorf("Expected source %s, got %s", result.Metadata.Source, unmarshaled.Metadata.Source)
	}
}

func TestJUnitTestSuitesSerialization(t *testing.T) {
	testSuites := JUnitTestSuites{
		Name:     "AllTests",
		Tests:    2,
		Failures: 1,
		Errors:   0,
		Time:     "1.5",
		TestSuites: []JUnitTestSuite{
			{
				Name:     "TestSuite",
				Package:  "test.package",
				Tests:    2,
				Failures: 1,
				Errors:   0,
				Skipped:  0,
				Time:     "1.5",
				TestCases: []JUnitTestCase{
					{
						Name:      "test1",
						ClassName: "TestClass",
						Time:      "0.5",
					},
					{
						Name:      "test2",
						ClassName: "TestClass",
						Time:      "1.0",
						Failure: &JUnitFailure{
							Message: "Test failed",
							Type:    "AssertionError",
							Content: "Expected: true, Actual: false",
						},
					},
				},
			},
		},
	}

	// Test XML serialization
	xmlData, err := xml.Marshal(testSuites)
	if err != nil {
		t.Fatalf("Failed to marshal to XML: %v", err)
	}

	var unmarshaled JUnitTestSuites
	if err := xml.Unmarshal(xmlData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal from XML: %v", err)
	}

	if unmarshaled.Name != testSuites.Name {
		t.Errorf("Expected name %s, got %s", testSuites.Name, unmarshaled.Name)
	}
}

func TestTAPResultSerialization(t *testing.T) {
	tapResult := TAPResult{
		Version: "13",
		Plan: &TAPPlan{
			Start: 1,
			End:   2,
			Skip:  "",
		},
		TestLines: []TAPTest{
			{
				Number:      1,
				OK:          true,
				Description: "Test 1",
				Directive: &TAPDirective{
					Type:   "SKIP",
					Reason: "Not implemented",
				},
				Diagnostics: []string{"Diagnostic line"},
				Raw:         "ok 1 - Test 1",
			},
		},
		Diagnostics: []string{"Global diagnostic"},
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(tapResult)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	var unmarshaled TAPResult
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	if unmarshaled.Version != tapResult.Version {
		t.Errorf("Expected version %s, got %s", tapResult.Version, unmarshaled.Version)
	}
}

func TestStatusConstants(t *testing.T) {
	statuses := []TestStatus{
		StatusPassed,
		StatusFailed,
		StatusSkipped,
		StatusError,
	}

	expectedValues := []string{"passed", "failed", "skipped", "error"}

	for i, status := range statuses {
		if string(status) != expectedValues[i] {
			t.Errorf("Expected status %s, got %s", expectedValues[i], string(status))
		}
	}
}

func TestGoTestActionConstants(t *testing.T) {
	actions := []string{
		ActionRun,
		ActionPause,
		ActionCont,
		ActionPass,
		ActionBench,
		ActionFail,
		ActionOutput,
		ActionSkip,
	}

	expectedValues := []string{"run", "pause", "cont", "pass", "bench", "fail", "output", "skip"}

	for i, action := range actions {
		if action != expectedValues[i] {
			t.Errorf("Expected action %s, got %s", expectedValues[i], action)
		}
	}
}

func TestXUnitAssembliesSerialization(t *testing.T) {
	assemblies := XUnitAssemblies{
		Timestamp: "2023-01-01T10:00:00Z",
		Assemblies: []XUnitAssembly{
			{
				Name:          "TestAssembly",
				Environment:   "Test",
				TestFramework: "xUnit",
				Total:         1,
				Passed:        1,
				Failed:        0,
				Skipped:       0,
				Collections: []XUnitCollection{
					{
						Name:   "TestCollection",
						Total:  1,
						Passed: 1,
						Failed: 0,
						Skipped: 0,
						Tests: []XUnitTest{
							{
								Name:   "TestMethod",
								Result: "Pass",
								Time:   "0.5",
							},
						},
					},
				},
			},
		},
	}

	// Test XML serialization
	xmlData, err := xml.Marshal(assemblies)
	if err != nil {
		t.Fatalf("Failed to marshal to XML: %v", err)
	}

	var unmarshaled XUnitAssemblies
	if err := xml.Unmarshal(xmlData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal from XML: %v", err)
	}

	if unmarshaled.Timestamp != assemblies.Timestamp {
		t.Errorf("Expected timestamp %s, got %s", assemblies.Timestamp, unmarshaled.Timestamp)
	}
}

func TestPytestReportSerialization(t *testing.T) {
	report := PytestReport{
		Report: PytestReportInfo{
			PytestVersion: "7.1.2",
			Outcome:       "passed",
			Created:       1672574400.0,
			Duration:      1.5,
		},
		Tests: []PytestTest{
			{
				NodeID:  "test_module.py::test_function",
				Outcome: "passed",
				Setup: PytestPhase{
					Outcome:  "passed",
					Duration: 0.1,
				},
				Call: PytestPhase{
					Outcome:  "passed",
					Duration: 0.5,
				},
				Teardown: PytestPhase{
					Outcome:  "passed",
					Duration: 0.1,
				},
			},
		},
		Summary: PytestSummary{
			Total:    1,
			Passed:   1,
			Failed:   0,
			Skipped:  0,
			Error:    0,
			Duration: 1.5,
			Outcome:  "passed",
		},
	}

	// Test JSON serialization
	jsonData, err := json.Marshal(report)
	if err != nil {
		t.Fatalf("Failed to marshal to JSON: %v", err)
	}

	var unmarshaled PytestReport
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal from JSON: %v", err)
	}

	if unmarshaled.Report.PytestVersion != report.Report.PytestVersion {
		t.Errorf("Expected pytest version %s, got %s", report.Report.PytestVersion, unmarshaled.Report.PytestVersion)
	}
}

func TestTestNGResultsSerialization(t *testing.T) {
	results := TestNGResults{
		Total:   1,
		Passed:  1,
		Failed:  0,
		Skipped: 0,
		Suites: []TestNGSuite{
			{
				Name:       "TestSuite",
				DurationMs: "1500",
				Tests: []TestNGTest{
					{
						Name:       "TestName",
						DurationMs: "1500",
						Classes: []TestNGClass{
							{
								Name: "TestClass",
								Methods: []TestNGMethod{
									{
										Name:       "testMethod",
										ClassName:  "TestClass",
										Status:     "PASS",
										DurationMs: "1500",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Test XML serialization
	xmlData, err := xml.Marshal(results)
	if err != nil {
		t.Fatalf("Failed to marshal to XML: %v", err)
	}

	var unmarshaled TestNGResults
	if err := xml.Unmarshal(xmlData, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal from XML: %v", err)
	}

	if unmarshaled.Total != results.Total {
		t.Errorf("Expected total %d, got %d", results.Total, unmarshaled.Total)
	}
}