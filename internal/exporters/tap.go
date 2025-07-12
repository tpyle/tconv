package exporters

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/tpyle/tconv/pkg/models"
)

// ExportTAP converts a UnifiedTestResult to TAP (Test Anything Protocol) format
func ExportTAP(result *models.UnifiedTestResult, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	return WriteTAP(result, file)
}

// WriteTAP writes a UnifiedTestResult as TAP format to the provided writer
func WriteTAP(result *models.UnifiedTestResult, writer io.Writer) error {
	// Write TAP version header
	if _, err := fmt.Fprintf(writer, "TAP version 13\n"); err != nil {
		return fmt.Errorf("failed to write TAP header: %w", err)
	}

	// Write plan
	if _, err := fmt.Fprintf(writer, "1..%d\n", result.Summary.Total); err != nil {
		return fmt.Errorf("failed to write TAP plan: %w", err)
	}

	testNumber := 1

	// Convert each test suite
	for _, suite := range result.TestSuites {
		// Add suite comment
		if _, err := fmt.Fprintf(writer, "# Test suite: %s\n", suite.Name); err != nil {
			return fmt.Errorf("failed to write suite comment: %w", err)
		}

		// Convert each test case
		for _, test := range suite.Tests {
			testName := formatTestName(suite.Name, test.Name)
			
			switch test.Status {
			case models.StatusPassed:
				if _, err := fmt.Fprintf(writer, "ok %d - %s\n", testNumber, testName); err != nil {
					return fmt.Errorf("failed to write test result: %w", err)
				}
			case models.StatusFailed, models.StatusError:
				if _, err := fmt.Fprintf(writer, "not ok %d - %s\n", testNumber, testName); err != nil {
					return fmt.Errorf("failed to write test result: %w", err)
				}
				// Add diagnostic information
				if test.Message != "" {
					if _, err := fmt.Fprintf(writer, "  ---\n  message: %s\n", escapeTAP(test.Message)); err != nil {
						return fmt.Errorf("failed to write test message: %w", err)
					}
				}
				if test.Details != "" {
					if _, err := fmt.Fprintf(writer, "  details: |\n"); err != nil {
						return fmt.Errorf("failed to write details header: %w", err)
					}
					// Indent details
					for _, line := range strings.Split(test.Details, "\n") {
						if _, err := fmt.Fprintf(writer, "    %s\n", line); err != nil {
							return fmt.Errorf("failed to write details: %w", err)
						}
					}
				}
				if test.Message != "" || test.Details != "" {
					if _, err := fmt.Fprintf(writer, "  ...\n"); err != nil {
						return fmt.Errorf("failed to write YAML end: %w", err)
					}
				}
			case models.StatusSkipped:
				if _, err := fmt.Fprintf(writer, "ok %d - %s # SKIP %s\n", testNumber, testName, escapeTAP(test.Message)); err != nil {
					return fmt.Errorf("failed to write skipped test: %w", err)
				}
			}

			// Add timing information as comment if available
			if test.Time > 0 {
				if _, err := fmt.Fprintf(writer, "# Duration: %.3fs\n", test.Time); err != nil {
					return fmt.Errorf("failed to write timing: %w", err)
				}
			}

			testNumber++
		}
	}

	// Write summary comments
	if _, err := fmt.Fprintf(writer, "# Summary:\n"); err != nil {
		return fmt.Errorf("failed to write summary header: %w", err)
	}
	if _, err := fmt.Fprintf(writer, "# Total: %d\n", result.Summary.Total); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}
	if _, err := fmt.Fprintf(writer, "# Passed: %d\n", result.Summary.Passed); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}
	if _, err := fmt.Fprintf(writer, "# Failed: %d\n", result.Summary.Failed); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}
	if _, err := fmt.Fprintf(writer, "# Skipped: %d\n", result.Summary.Skipped); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}
	if _, err := fmt.Fprintf(writer, "# Errors: %d\n", result.Summary.Errors); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}
	if _, err := fmt.Fprintf(writer, "# Duration: %.3fs\n", result.Summary.Duration); err != nil {
		return fmt.Errorf("failed to write summary: %w", err)
	}

	return nil
}

// formatTestName creates a readable test name for TAP output
func formatTestName(suiteName, testName string) string {
	if suiteName != "" {
		return fmt.Sprintf("%s::%s", suiteName, testName)
	}
	return testName
}

// escapeTAP escapes special characters for TAP format
func escapeTAP(text string) string {
	// Replace newlines with spaces and remove excessive whitespace
	text = strings.ReplaceAll(text, "\n", " ")
	text = strings.ReplaceAll(text, "\r", " ")
	text = strings.Join(strings.Fields(text), " ")
	return text
}