package parsers

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

func ParseJUnit(filePath string) (*models.TikiTestResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var testSuites models.JUnitTestSuites
	if err := xml.Unmarshal(data, &testSuites); err != nil {
		var testSuite models.JUnitTestSuite
		if err := xml.Unmarshal(data, &testSuite); err != nil {
			return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
		}
		testSuites.TestSuites = []models.JUnitTestSuite{testSuite}
	}

	result := &models.TikiTestResult{
		Metadata: models.TestMetadata{
			Source:    "junit",
			Timestamp: time.Now(),
			Version:   "1.0",
		},
		TestSuites: make([]models.TestSuite, 0, len(testSuites.TestSuites)),
	}

	var summary models.TestSummary

	for _, junitSuite := range testSuites.TestSuites {
		suite := models.TestSuite{
			Name:     junitSuite.Name,
			Package:  junitSuite.Package,
			Tests:    make([]models.TestCase, 0, len(junitSuite.TestCases)),
			Errors:   junitSuite.Errors,
			Failures: junitSuite.Failures,
			Skipped:  junitSuite.Skipped,
		}

		if timeValue, err := strconv.ParseFloat(junitSuite.Time, 64); err == nil {
			suite.Time = timeValue
		}

		if junitSuite.Properties != nil {
			suite.Properties = make(map[string]string)
			for _, prop := range junitSuite.Properties.Properties {
				suite.Properties[prop.Name] = prop.Value
			}
		}

		for _, junitCase := range junitSuite.TestCases {
			testCase := models.TestCase{
				Name:      junitCase.Name,
				ClassName: junitCase.ClassName,
				SystemOut: junitCase.SystemOut,
				SystemErr: junitCase.SystemErr,
			}

			if timeValue, err := strconv.ParseFloat(junitCase.Time, 64); err == nil {
				testCase.Time = timeValue
			}

			if junitCase.Failure != nil {
				testCase.Status = models.StatusFailed
				testCase.Message = junitCase.Failure.Message
				testCase.Details = junitCase.Failure.Content
				testCase.ErrorType = junitCase.Failure.Type
				summary.Failed++
			} else if junitCase.Error != nil {
				testCase.Status = models.StatusError
				testCase.Message = junitCase.Error.Message
				testCase.Details = junitCase.Error.Content
				testCase.ErrorType = junitCase.Error.Type
				summary.Errors++
			} else if junitCase.Skipped != nil {
				testCase.Status = models.StatusSkipped
				testCase.Message = junitCase.Skipped.Message
				summary.Skipped++
			} else {
				testCase.Status = models.StatusPassed
				summary.Passed++
			}

			suite.Tests = append(suite.Tests, testCase)
			summary.Total++
			summary.Duration += testCase.Time
		}

		result.TestSuites = append(result.TestSuites, suite)
	}

	result.Summary = summary
	return result, nil
}