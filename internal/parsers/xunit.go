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

func ParseXUnit(filePath string) (*models.TikiTestResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var assemblies models.XUnitAssemblies
	if err := xml.Unmarshal(data, &assemblies); err != nil {
		var assembly models.XUnitAssembly
		if err := xml.Unmarshal(data, &assembly); err != nil {
			return nil, fmt.Errorf("failed to unmarshal XML: %w", err)
		}
		assemblies.Assemblies = []models.XUnitAssembly{assembly}
	}

	result := &models.TikiTestResult{
		Metadata: models.TestMetadata{
			Source:    "xunit",
			Timestamp: time.Now(),
			Version:   "1.0",
			Framework: "xUnit",
		},
		TestSuites: make([]models.TestSuite, 0),
	}

	if assemblies.Timestamp != "" {
		result.Metadata.Environment = map[string]string{
			"timestamp": assemblies.Timestamp,
		}
	}

	var summary models.TestSummary

	for _, assembly := range assemblies.Assemblies {
		if assembly.Environment != "" || assembly.TestFramework != "" {
			if result.Metadata.Environment == nil {
				result.Metadata.Environment = make(map[string]string)
			}
			if assembly.Environment != "" {
				result.Metadata.Environment["environment"] = assembly.Environment
			}
			if assembly.TestFramework != "" {
				result.Metadata.Environment["test_framework"] = assembly.TestFramework
			}
		}

		for _, collection := range assembly.Collections {
			suite := models.TestSuite{
				Name:     collection.Name,
				Package:  assembly.Name,
				Tests:    make([]models.TestCase, 0, len(collection.Tests)),
				Errors:   0,
				Failures: collection.Failed,
				Skipped:  collection.Skipped,
			}

			if timeValue, err := strconv.ParseFloat(collection.Time, 64); err == nil {
				suite.Time = timeValue
			}

			for _, xunitTest := range collection.Tests {
				testCase := models.TestCase{
					Name:      xunitTest.Name,
					ClassName: xunitTest.Type,
					SystemOut: xunitTest.Output,
				}

				if xunitTest.Method != "" {
					testCase.Properties = map[string]string{
						"method": xunitTest.Method,
					}
				}

				if timeValue, err := strconv.ParseFloat(xunitTest.Time, 64); err == nil {
					testCase.Time = timeValue
				}

				// Determine status and details
				switch xunitTest.Result {
				case "Pass":
					testCase.Status = models.StatusPassed
					summary.Passed++
				case "Fail":
					testCase.Status = models.StatusFailed
					if xunitTest.Failure != nil {
						testCase.Message = xunitTest.Failure.Message
						testCase.Details = xunitTest.Failure.StackTrace
						testCase.ErrorType = xunitTest.Failure.ExceptionType
					}
					summary.Failed++
				case "Skip":
					testCase.Status = models.StatusSkipped
					if xunitTest.Reason != nil {
						testCase.Message = xunitTest.Reason.Content
					}
					summary.Skipped++
				default:
					testCase.Status = models.StatusError
					testCase.Message = "Unknown result: " + xunitTest.Result
					suite.Errors++
					summary.Errors++
				}

				// Add traits as properties
				if xunitTest.Traits != nil {
					if testCase.Properties == nil {
						testCase.Properties = make(map[string]string)
					}
					for _, trait := range xunitTest.Traits.Traits {
						testCase.Properties["trait_"+trait.Name] = trait.Value
					}
				}

				suite.Tests = append(suite.Tests, testCase)
				summary.Total++
				summary.Duration += testCase.Time
			}

			result.TestSuites = append(result.TestSuites, suite)
		}
	}

	result.Summary = summary
	return result, nil
}