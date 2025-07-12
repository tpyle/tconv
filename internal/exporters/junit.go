package exporters

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

// JUnitTestSuites represents the root element for JUnit XML output
type JUnitTestSuites struct {
	XMLName    xml.Name         `xml:"testsuites"`
	Name       string           `xml:"name,attr,omitempty"`
	Tests      int              `xml:"tests,attr"`
	Failures   int              `xml:"failures,attr"`
	Errors     int              `xml:"errors,attr"`
	Skipped    int              `xml:"skipped,attr"`
	Time       float64          `xml:"time,attr"`
	Timestamp  string           `xml:"timestamp,attr,omitempty"`
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents a test suite in JUnit XML format
type JUnitTestSuite struct {
	XMLName   xml.Name        `xml:"testsuite"`
	Name      string          `xml:"name,attr"`
	Package   string          `xml:"package,attr,omitempty"`
	Tests     int             `xml:"tests,attr"`
	Failures  int             `xml:"failures,attr"`
	Errors    int             `xml:"errors,attr"`
	Skipped   int             `xml:"skipped,attr"`
	Time      float64         `xml:"time,attr"`
	Timestamp string          `xml:"timestamp,attr,omitempty"`
	TestCases []JUnitTestCase `xml:"testcase"`
}

// JUnitTestCase represents a test case in JUnit XML format
type JUnitTestCase struct {
	XMLName   xml.Name         `xml:"testcase"`
	Name      string           `xml:"name,attr"`
	ClassName string           `xml:"classname,attr,omitempty"`
	Time      float64          `xml:"time,attr"`
	Failure   *JUnitFailure    `xml:"failure,omitempty"`
	Error     *JUnitError      `xml:"error,omitempty"`
	Skipped   *JUnitSkipped    `xml:"skipped,omitempty"`
	SystemOut string           `xml:"system-out,omitempty"`
	SystemErr string           `xml:"system-err,omitempty"`
}

// JUnitFailure represents a test failure
type JUnitFailure struct {
	XMLName xml.Name `xml:"failure"`
	Type    string   `xml:"type,attr,omitempty"`
	Message string   `xml:"message,attr,omitempty"`
	Content string   `xml:",chardata"`
}

// JUnitError represents a test error
type JUnitError struct {
	XMLName xml.Name `xml:"error"`
	Type    string   `xml:"type,attr,omitempty"`
	Message string   `xml:"message,attr,omitempty"`
	Content string   `xml:",chardata"`
}

// JUnitSkipped represents a skipped test
type JUnitSkipped struct {
	XMLName xml.Name `xml:"skipped"`
	Message string   `xml:"message,attr,omitempty"`
}

// ExportJUnit converts a TikiTestResult to JUnit XML format
func ExportJUnit(result *models.UnifiedTestResult, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	return WriteJUnit(result, file)
}

// WriteJUnit writes a TikiTestResult as JUnit XML to the provided writer
func WriteJUnit(result *models.UnifiedTestResult, writer io.Writer) error {
	junitSuites := convertToJUnit(result)

	// Write XML header
	if _, err := writer.Write([]byte(xml.Header)); err != nil {
		return fmt.Errorf("failed to write XML header: %w", err)
	}

	// Create encoder with indentation
	encoder := xml.NewEncoder(writer)
	encoder.Indent("", "  ")

	if err := encoder.Encode(junitSuites); err != nil {
		return fmt.Errorf("failed to encode JUnit XML: %w", err)
	}

	return nil
}

// convertToJUnit converts TikiTestResult to JUnit XML structure
func convertToJUnit(result *models.UnifiedTestResult) *JUnitTestSuites {
	junitSuites := &JUnitTestSuites{
		Name:      "Combined Test Results",
		Tests:     result.Summary.Total,
		Failures:  result.Summary.Failed,
		Errors:    result.Summary.Errors,
		Skipped:   result.Summary.Skipped,
		Time:      result.Summary.Duration,
		Timestamp: result.Metadata.Timestamp.Format(time.RFC3339),
	}

	// Convert each test suite
	for _, suite := range result.TestSuites {
		junitSuite := JUnitTestSuite{
			Name:      suite.Name,
			Package:   suite.Package,
			Tests:     len(suite.Tests),
			Failures:  suite.Failures,
			Errors:    suite.Errors,
			Skipped:   suite.Skipped,
			Time:      suite.Time,
			Timestamp: result.Metadata.Timestamp.Format(time.RFC3339),
		}

		// Convert each test case
		for _, test := range suite.Tests {
			junitCase := JUnitTestCase{
				Name:      test.Name,
				ClassName: test.ClassName,
				Time:      test.Time,
				SystemOut: test.SystemOut,
				SystemErr: test.SystemErr,
			}

			// Handle different test statuses
			switch test.Status {
			case models.StatusFailed:
				junitCase.Failure = &JUnitFailure{
					Type:    test.ErrorType,
					Message: test.Message,
					Content: test.Details,
				}
			case models.StatusError:
				junitCase.Error = &JUnitError{
					Type:    test.ErrorType,
					Message: test.Message,
					Content: test.Details,
				}
			case models.StatusSkipped:
				junitCase.Skipped = &JUnitSkipped{
					Message: test.Message,
				}
			}

			junitSuite.TestCases = append(junitSuite.TestCases, junitCase)
		}

		junitSuites.TestSuites = append(junitSuites.TestSuites, junitSuite)
	}

	return junitSuites
}