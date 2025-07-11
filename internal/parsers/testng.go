package parsers

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func ParseTestNG(filePath string) (*models.UnifiedTestResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var testngResults models.TestNGResults
	if err := xml.Unmarshal(data, &testngResults); err != nil {
		return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
	}

	result := &models.UnifiedTestResult{
		Metadata: models.TestMetadata{
			Source:    "testng",
			Timestamp: time.Now(),
			Version:   "1.0",
			Framework: "TestNG",
		},
		TestSuites: make([]models.TestSuite, 0),
	}

	var summary models.TestSummary
	summary.Total = testngResults.Total
	summary.Passed = testngResults.Passed
	summary.Failed = testngResults.Failed
	summary.Skipped = testngResults.Skipped

	for _, testngSuite := range testngResults.Suites {
		for _, testngTest := range testngSuite.Tests {
			suite := models.TestSuite{
				Name:    testngTest.Name,
				Package: testngSuite.Name,
				Tests:   make([]models.TestCase, 0),
			}

			// Calculate suite duration
			if testngTest.DurationMs != "" {
				if durationMs, err := strconv.ParseFloat(testngTest.DurationMs, 64); err == nil {
					suite.Time = durationMs / 1000.0 // Convert to seconds
				}
			}

			for _, testngClass := range testngTest.Classes {
				for _, method := range testngClass.Methods {
					testCase := models.TestCase{
						Name:      method.Name,
						ClassName: method.ClassName,
						Properties: map[string]string{
							"signature":   method.Signature,
							"is_config":   method.IsConfig,
						},
					}

					// Add groups if present
					if method.Groups != "" {
						testCase.Properties["groups"] = method.Groups
					}

					// Add description if present
					if method.Description != "" {
						testCase.Properties["description"] = method.Description
					}

					// Calculate test duration
					if method.DurationMs != "" {
						if durationMs, err := strconv.ParseFloat(method.DurationMs, 64); err == nil {
							testCase.Time = durationMs / 1000.0 // Convert to seconds
							summary.Duration += testCase.Time
						}
					}

					// Determine status
					switch strings.ToUpper(method.Status) {
					case "PASS":
						testCase.Status = models.StatusPassed
					case "FAIL":
						testCase.Status = models.StatusFailed
						if method.Exception != nil {
							testCase.Message = method.Exception.Message
							testCase.Details = method.Exception.FullStacktrace
							testCase.ErrorType = method.Exception.Class
						}
						suite.Failures++
					case "SKIP":
						testCase.Status = models.StatusSkipped
						suite.Skipped++
					default:
						testCase.Status = models.StatusError
						testCase.Message = "Unknown status: " + method.Status
						suite.Errors++
					}

					// Collect reporter output
					if len(method.Reporter) > 0 {
						var outputs []string
						for _, reporter := range method.Reporter {
							outputs = append(outputs, reporter.Output)
						}
						testCase.SystemOut = strings.Join(outputs, "\n")
					}

					// Add parameters if present
					if len(method.Params) > 0 {
						var params []string
						for _, param := range method.Params {
							params = append(params, fmt.Sprintf("param[%s]=%s", param.Index, param.Value))
						}
						testCase.Properties["parameters"] = strings.Join(params, ", ")
					}

					// Add timing information
					if method.StartedAt != "" {
						testCase.Properties["started_at"] = method.StartedAt
					}
					if method.FinishedAt != "" {
						testCase.Properties["finished_at"] = method.FinishedAt
					}

					suite.Tests = append(suite.Tests, testCase)
				}
			}

			result.TestSuites = append(result.TestSuites, suite)
		}
	}

	result.Summary = summary
	return result, nil
}