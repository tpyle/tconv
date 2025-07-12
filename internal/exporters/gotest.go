package exporters

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

// GoTestEvent represents a single event in Go test JSON output
type GoTestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package,omitempty"`
	Test    string    `json:"Test,omitempty"`
	Elapsed float64   `json:"Elapsed,omitempty"`
	Output  string    `json:"Output,omitempty"`
}

// ExportGoTest converts a TikiTestResult to Go test JSON format
func ExportGoTest(result *models.TikiTestResult, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	return WriteGoTest(result, file)
}

// WriteGoTest writes a TikiTestResult as Go test JSON to the provided writer
func WriteGoTest(result *models.TikiTestResult, writer io.Writer) error {
	encoder := json.NewEncoder(writer)
	baseTime := result.Metadata.Timestamp

	// Generate events for each test suite
	for _, suite := range result.TestSuites {
		packageName := suite.Package
		if packageName == "" {
			packageName = "test/" + strings.ToLower(strings.ReplaceAll(suite.Name, " ", ""))
		}

		// Package start event
		if err := encoder.Encode(GoTestEvent{
			Time:    baseTime,
			Action:  "run",
			Package: packageName,
		}); err != nil {
			return fmt.Errorf("failed to write package start event: %w", err)
		}

		currentTime := baseTime

		// Generate events for each test
		for _, test := range suite.Tests {
			testName := cleanTestName(test.Name)
			
			// Test start event
			if err := encoder.Encode(GoTestEvent{
				Time:    currentTime,
				Action:  "run",
				Package: packageName,
				Test:    testName,
			}); err != nil {
				return fmt.Errorf("failed to write test start event: %w", err)
			}

			// Test output if available
			if test.SystemOut != "" {
				if err := encoder.Encode(GoTestEvent{
					Time:    currentTime,
					Action:  "output",
					Package: packageName,
					Test:    testName,
					Output:  test.SystemOut + "\n",
				}); err != nil {
					return fmt.Errorf("failed to write test output: %w", err)
				}
			}

			// Error output if available
			if test.SystemErr != "" {
				if err := encoder.Encode(GoTestEvent{
					Time:    currentTime,
					Action:  "output",
					Package: packageName,
					Test:    testName,
					Output:  test.SystemErr + "\n",
				}); err != nil {
					return fmt.Errorf("failed to write test error output: %w", err)
				}
			}

			// Test result event
			var action string
			switch test.Status {
			case models.StatusPassed:
				action = "pass"
			case models.StatusFailed, models.StatusError:
				action = "fail"
				// Add failure details as output
				if test.Message != "" || test.Details != "" {
					failureOutput := fmt.Sprintf("--- FAIL: %s\n", testName)
					if test.Message != "" {
						failureOutput += fmt.Sprintf("    %s\n", test.Message)
					}
					if test.Details != "" {
						failureOutput += fmt.Sprintf("    %s\n", test.Details)
					}
					
					if err := encoder.Encode(GoTestEvent{
						Time:    currentTime,
						Action:  "output",
						Package: packageName,
						Test:    testName,
						Output:  failureOutput,
					}); err != nil {
						return fmt.Errorf("failed to write failure output: %w", err)
					}
				}
			case models.StatusSkipped:
				action = "skip"
				// Add skip reason as output
				if test.Message != "" {
					skipOutput := fmt.Sprintf("--- SKIP: %s\n    %s\n", testName, test.Message)
					if err := encoder.Encode(GoTestEvent{
						Time:    currentTime,
						Action:  "output",
						Package: packageName,
						Test:    testName,
						Output:  skipOutput,
					}); err != nil {
						return fmt.Errorf("failed to write skip output: %w", err)
					}
				}
			}

			// Update time for test completion
			currentTime = currentTime.Add(time.Duration(test.Time * float64(time.Second)))

			if err := encoder.Encode(GoTestEvent{
				Time:    currentTime,
				Action:  action,
				Package: packageName,
				Test:    testName,
				Elapsed: test.Time,
			}); err != nil {
				return fmt.Errorf("failed to write test result event: %w", err)
			}
		}

		// Package result event
		packageAction := "pass"
		if suite.Failures > 0 || suite.Errors > 0 {
			packageAction = "fail"
		}

		// Update time for package completion
		currentTime = currentTime.Add(time.Duration(suite.Time * float64(time.Second)))

		if err := encoder.Encode(GoTestEvent{
			Time:    currentTime,
			Action:  packageAction,
			Package: packageName,
			Elapsed: suite.Time,
		}); err != nil {
			return fmt.Errorf("failed to write package result event: %w", err)
		}
	}

	return nil
}

// cleanTestName ensures the test name is valid for Go test format
func cleanTestName(name string) string {
	// Replace spaces and special characters with underscores
	name = strings.ReplaceAll(name, " ", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	
	// Ensure it starts with Test if it doesn't already
	if !strings.HasPrefix(name, "Test") {
		if len(name) > 0 {
			name = "Test" + strings.ToUpper(name[:1]) + name[1:]
		} else {
			name = "Test"
		}
	}
	
	return name
}