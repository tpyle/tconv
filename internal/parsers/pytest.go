package parsers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func ParsePytest(filePath string) (*models.UnifiedTestResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var pytestReport models.PytestReport
	if err := json.Unmarshal(data, &pytestReport); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	result := &models.UnifiedTestResult{
		Metadata: models.TestMetadata{
			Source:    "pytest",
			Timestamp: time.Now(),
			Version:   "1.0",
			Framework: "pytest",
		},
		TestSuites: make([]models.TestSuite, 0),
	}

	if pytestReport.Report.PytestVersion != "" {
		result.Metadata.Environment = map[string]string{
			"pytest_version": pytestReport.Report.PytestVersion,
			"outcome":        pytestReport.Report.Outcome,
		}
	}

	// Group tests by module/file
	testsByModule := make(map[string][]models.PytestTest)
	for _, test := range pytestReport.Tests {
		// Extract module name from nodeid (e.g., "tests/test_example.py::test_function" -> "tests/test_example.py")
		module := extractModuleFromNodeID(test.NodeID)
		testsByModule[module] = append(testsByModule[module], test)
	}

	var summary models.TestSummary
	summary.Total = pytestReport.Summary.Total
	summary.Passed = pytestReport.Summary.Passed
	summary.Failed = pytestReport.Summary.Failed
	summary.Skipped = pytestReport.Summary.Skipped
	summary.Errors = pytestReport.Summary.Error
	summary.Duration = pytestReport.Summary.Duration

	for module, tests := range testsByModule {
		suite := models.TestSuite{
			Name:    filepath.Base(module),
			Package: module,
			Tests:   make([]models.TestCase, 0, len(tests)),
		}

		for _, pytestTest := range tests {
			testCase := models.TestCase{
				Name:      extractTestNameFromNodeID(pytestTest.NodeID),
				ClassName: module,
				Properties: map[string]string{
					"nodeid": pytestTest.NodeID,
				},
			}

			// Calculate total duration from all phases
			var totalDuration float64
			var setupDuration, callDuration, teardownDuration float64

			if pytestTest.Setup.Duration > 0 {
				setupDuration = pytestTest.Setup.Duration
				totalDuration += setupDuration
			}
			if pytestTest.Call.Duration > 0 {
				callDuration = pytestTest.Call.Duration
				totalDuration += callDuration
			}
			if pytestTest.Teardown.Duration > 0 {
				teardownDuration = pytestTest.Teardown.Duration
				totalDuration += teardownDuration
			}

			testCase.Time = totalDuration

			// Determine status based on outcome
			switch pytestTest.Outcome {
			case "passed":
				testCase.Status = models.StatusPassed
			case "failed":
				testCase.Status = models.StatusFailed
				testCase.Message = extractFailureMessage(pytestTest)
				testCase.Details = extractFailureDetails(pytestTest)
				suite.Failures++
			case "skipped":
				testCase.Status = models.StatusSkipped
				testCase.Message = extractSkipMessage(pytestTest)
				suite.Skipped++
			case "error":
				testCase.Status = models.StatusError
				testCase.Message = extractErrorMessage(pytestTest)
				testCase.Details = extractErrorDetails(pytestTest)
				suite.Errors++
			}

			// Collect output from all phases
			var outputs []string
			if pytestTest.Setup.Stdout != "" {
				outputs = append(outputs, "Setup stdout:\n"+pytestTest.Setup.Stdout)
			}
			if pytestTest.Call.Stdout != "" {
				outputs = append(outputs, "Call stdout:\n"+pytestTest.Call.Stdout)
			}
			if pytestTest.Teardown.Stdout != "" {
				outputs = append(outputs, "Teardown stdout:\n"+pytestTest.Teardown.Stdout)
			}
			if len(outputs) > 0 {
				testCase.SystemOut = strings.Join(outputs, "\n")
			}

			// Collect stderr from all phases
			var errors []string
			if pytestTest.Setup.Stderr != "" {
				errors = append(errors, "Setup stderr:\n"+pytestTest.Setup.Stderr)
			}
			if pytestTest.Call.Stderr != "" {
				errors = append(errors, "Call stderr:\n"+pytestTest.Call.Stderr)
			}
			if pytestTest.Teardown.Stderr != "" {
				errors = append(errors, "Teardown stderr:\n"+pytestTest.Teardown.Stderr)
			}
			if len(errors) > 0 {
				testCase.SystemErr = strings.Join(errors, "\n")
			}

			// Add keywords and location as properties
			if len(pytestTest.Keywords) > 0 {
				testCase.Properties["keywords"] = strings.Join(pytestTest.Keywords, ",")
			}
			if len(pytestTest.Location) > 0 {
				locationStr := fmt.Sprintf("%v", pytestTest.Location)
				testCase.Properties["location"] = locationStr
			}

			suite.Tests = append(suite.Tests, testCase)
		}

		suite.Time = summary.Duration
		result.TestSuites = append(result.TestSuites, suite)
	}

	result.Summary = summary
	return result, nil
}

func extractModuleFromNodeID(nodeID string) string {
	parts := strings.Split(nodeID, "::")
	if len(parts) > 0 {
		return parts[0]
	}
	return "unknown"
}

func extractTestNameFromNodeID(nodeID string) string {
	parts := strings.Split(nodeID, "::")
	if len(parts) > 1 {
		return strings.Join(parts[1:], "::")
	}
	return nodeID
}

func extractFailureMessage(test models.PytestTest) string {
	if test.Call.Crash != nil {
		return test.Call.Crash.Message
	}
	if test.Setup.Crash != nil {
		return test.Setup.Crash.Message
	}
	if test.Teardown.Crash != nil {
		return test.Teardown.Crash.Message
	}
	
	// If no crash, try to get message from traceback
	if test.Call.Traceback != nil && len(test.Call.Traceback.Entries) > 0 {
		// Return the last traceback entry message which typically contains the assertion
		lastEntry := test.Call.Traceback.Entries[len(test.Call.Traceback.Entries)-1]
		return lastEntry.Message
	}
	
	return "Test failed"
}

func extractFailureDetails(test models.PytestTest) string {
	var details []string

	if test.Call.Traceback != nil && len(test.Call.Traceback.Entries) > 0 {
		details = append(details, "Call traceback:")
		for _, entry := range test.Call.Traceback.Entries {
			details = append(details, fmt.Sprintf("%s:%d - %s", entry.Path, entry.LineNo, entry.Message))
		}
	}

	if test.Setup.Traceback != nil && len(test.Setup.Traceback.Entries) > 0 {
		details = append(details, "Setup traceback:")
		for _, entry := range test.Setup.Traceback.Entries {
			details = append(details, fmt.Sprintf("%s:%d - %s", entry.Path, entry.LineNo, entry.Message))
		}
	}

	if test.Teardown.Traceback != nil && len(test.Teardown.Traceback.Entries) > 0 {
		details = append(details, "Teardown traceback:")
		for _, entry := range test.Teardown.Traceback.Entries {
			details = append(details, fmt.Sprintf("%s:%d - %s", entry.Path, entry.LineNo, entry.Message))
		}
	}

	return strings.Join(details, "\n")
}

func extractSkipMessage(test models.PytestTest) string {
	// Check for skip reason in various phases
	if test.Call.Log != "" {
		return test.Call.Log
	}
	if test.Setup.Log != "" {
		return test.Setup.Log
	}
	return "Test skipped"
}

func extractErrorMessage(test models.PytestTest) string {
	if test.Call.Crash != nil {
		return test.Call.Crash.Message
	}
	if test.Setup.Crash != nil {
		return test.Setup.Crash.Message
	}
	if test.Teardown.Crash != nil {
		return test.Teardown.Crash.Message
	}
	return "Test error"
}

func extractErrorDetails(test models.PytestTest) string {
	var details []string

	if test.Call.Crash != nil {
		details = append(details, fmt.Sprintf("Call error at %s:%d - %s", 
			test.Call.Crash.Path, test.Call.Crash.LineNo, test.Call.Crash.Message))
	}

	if test.Setup.Crash != nil {
		details = append(details, fmt.Sprintf("Setup error at %s:%d - %s", 
			test.Setup.Crash.Path, test.Setup.Crash.LineNo, test.Setup.Crash.Message))
	}

	if test.Teardown.Crash != nil {
		details = append(details, fmt.Sprintf("Teardown error at %s:%d - %s", 
			test.Teardown.Crash.Path, test.Teardown.Crash.LineNo, test.Teardown.Crash.Message))
	}

	return strings.Join(details, "\n")
}