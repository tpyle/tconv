package parsers

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func ParsePostman(filePath string) (*models.TikiTestResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var postmanRun models.PostmanRun
	if err := json.Unmarshal(data, &postmanRun); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	result := &models.TikiTestResult{
		Metadata: models.TestMetadata{
			Source:    "postman",
			Timestamp: time.Now(),
			Version:   "1.0",
			Framework: "Newman",
		},
		TestSuites: make([]models.TestSuite, 0),
	}

	if postmanRun.Info.Name != "" {
		result.Metadata.Environment = map[string]string{
			"collection_name": postmanRun.Info.Name,
			"collection_id":   postmanRun.Info.PostmanID,
		}
	}

	testSuite := models.TestSuite{
		Name:     postmanRun.Info.Name,
		Package:  "postman",
		Tests:    make([]models.TestCase, 0, len(postmanRun.Run.Executions)),
		Errors:   0,
		Failures: 0,
		Skipped:  0,
		Time:     time.Duration(float64(postmanRun.Run.Timings.Completed-postmanRun.Run.Timings.Started) * float64(time.Millisecond)),
	}

	var summary models.TestSummary
	summary.Duration = testSuite.Time

	for _, execution := range postmanRun.Run.Executions {
		testCase := models.TestCase{
			Name:        execution.Item.Name,
			ClassName:   "postman.request",
			Time:        time.Duration(execution.Response.ResponseTime * float64(time.Millisecond)),
			SystemOut:   execution.Response.Body,
			Properties:  make(map[string]string),
			Assertions:  make([]models.Assertion, 0),
		}

		testCase.Properties["method"] = execution.Request.Method
		testCase.Properties["url"] = execution.Request.URL.Raw
		testCase.Properties["status_code"] = fmt.Sprintf("%d", execution.Response.Code)
		testCase.Properties["status"] = execution.Response.Status

		hasFailures := false
		hasErrors := execution.RequestError != nil

		for _, assertion := range execution.Assertions {
			unifiedAssertion := models.Assertion{
				Name:   assertion.Assertion,
				Passed: assertion.Error == nil,
			}

			if assertion.Error != nil {
				if errorStr, ok := assertion.Error.(string); ok {
					unifiedAssertion.Message = errorStr
				} else if errorMap, ok := assertion.Error.(map[string]interface{}); ok {
					if message, exists := errorMap["message"]; exists {
						unifiedAssertion.Message = fmt.Sprintf("%v", message)
					}
				}
				hasFailures = true
			}

			testCase.Assertions = append(testCase.Assertions, unifiedAssertion)
		}

		if hasErrors {
			testCase.Status = models.StatusError
			testCase.Message = "Request failed"
			if execution.RequestError != nil {
				testCase.Details = fmt.Sprintf("%v", execution.RequestError)
			}
			testSuite.Errors++
			summary.Errors++
		} else if hasFailures {
			testCase.Status = models.StatusFailed
			var failedAssertions []string
			for _, assertion := range testCase.Assertions {
				if !assertion.Passed {
					failedAssertions = append(failedAssertions, assertion.Name)
				}
			}
			testCase.Message = fmt.Sprintf("Assertions failed: %s", strings.Join(failedAssertions, ", "))
			testSuite.Failures++
			summary.Failed++
		} else {
			testCase.Status = models.StatusPassed
			summary.Passed++
		}

		testSuite.Tests = append(testSuite.Tests, testCase)
		summary.Total++
	}

	result.TestSuites = append(result.TestSuites, testSuite)
	result.Summary = summary

	return result, nil
}