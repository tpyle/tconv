package parsers

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

// ParseGoTest parses Go test JSON output and converts it to tiki format.
//
// This parser handles the JSON output format produced by 'go test -json'.
// It processes test events including pass/fail status, timing information,
// and output messages to create a comprehensive test result.
//
// Returns a TikiTestResult with converted test data or an error if parsing fails.
func ParseGoTest(filePath string) (*models.TikiTestResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	result := &models.TikiTestResult{
		Metadata: models.TestMetadata{
			Source:    "gotest",
			Timestamp: time.Now(),
			Version:   "1.0",
			Framework: "go test",
		},
		TestSuites: make([]models.TestSuite, 0),
	}

	packages := make(map[string]*models.TestSuite)
	tests := make(map[string]*models.TestCase)
	var summary models.TestSummary

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			continue
		}

		var event models.GoTestEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		packageName := event.Package
		if packageName == "" {
			packageName = "main"
		}

		if _, exists := packages[packageName]; !exists {
			packages[packageName] = &models.TestSuite{
				Name:    packageName,
				Package: packageName,
				Tests:   make([]models.TestCase, 0),
			}
		}

		suite := packages[packageName]

		switch event.Action {
		case models.ActionRun:
			if event.Test != "" {
				testKey := packageName + "::" + event.Test
				tests[testKey] = &models.TestCase{
					Name:      event.Test,
					ClassName: packageName,
					Status:    models.StatusPassed,
				}
			}

		case models.ActionOutput:
			if event.Test != "" {
				testKey := packageName + "::" + event.Test
				if test, exists := tests[testKey]; exists {
					if test.SystemOut == "" {
						test.SystemOut = event.Output
					} else {
						test.SystemOut += event.Output
					}
				}
			}

		case models.ActionPass:
			if event.Test != "" {
				testKey := packageName + "::" + event.Test
				if test, exists := tests[testKey]; exists {
					test.Status = models.StatusPassed
					test.Time = time.Duration(event.Elapsed * float64(time.Second))
					summary.Passed++
					summary.Total++
				}
			} else {
				suite.Time = time.Duration(event.Elapsed * float64(time.Second))
			}

		case models.ActionFail:
			if event.Test != "" {
				testKey := packageName + "::" + event.Test
				if test, exists := tests[testKey]; exists {
					test.Status = models.StatusFailed
					test.Time = time.Duration(event.Elapsed * float64(time.Second))
					test.Message = "Test failed"
					if test.SystemOut != "" {
						test.Details = strings.TrimSpace(test.SystemOut)
					}
					suite.Failures++
					summary.Failed++
					summary.Total++
				}
			} else {
				suite.Time = time.Duration(event.Elapsed * float64(time.Second))
			}

		case models.ActionSkip:
			if event.Test != "" {
				testKey := packageName + "::" + event.Test
				if test, exists := tests[testKey]; exists {
					test.Status = models.StatusSkipped
					test.Time = time.Duration(event.Elapsed * float64(time.Second))
					test.Message = "Test skipped"
					if test.SystemOut != "" {
						test.Details = strings.TrimSpace(test.SystemOut)
					}
					suite.Skipped++
					summary.Skipped++
					summary.Total++
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	for _, suite := range packages {
		for _, test := range tests {
			if test.ClassName == suite.Package {
				suite.Tests = append(suite.Tests, *test)
				summary.Duration += test.Time
			}
		}
		result.TestSuites = append(result.TestSuites, *suite)
	}

	result.Summary = summary
	return result, nil
}