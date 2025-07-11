package parsers

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/tpyle/tconv/pkg/models"
)

var (
	tapVersionRegex = regexp.MustCompile(`^TAP version (\d+)`)
	tapPlanRegex    = regexp.MustCompile(`^(\d+)\.\.(\d+)(?:\s*#\s*SKIP\s*(.*))?`)
	tapTestRegex    = regexp.MustCompile(`^(not\s+)?ok\s+(\d+)?\s*([^#]*?)(?:\s*#\s*(SKIP|TODO|skip|todo)\s*(.*))?$`)
	tapBailOutRegex = regexp.MustCompile(`^Bail out!\s*(.*)`)
	tapDiagRegex    = regexp.MustCompile(`^#\s*(.*)`)
)

func ParseTAP(filePath string) (*models.UnifiedTestResult, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	tapResult, err := parseTAPContent(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TAP content: %w", err)
	}

	return convertTAPToUnified(tapResult), nil
}

func parseTAPContent(reader io.Reader) (*models.TAPResult, error) {
	scanner := bufio.NewScanner(reader)
	result := &models.TAPResult{
		TestLines:   make([]models.TAPTest, 0),
		Diagnostics: make([]string, 0),
	}

	var currentDiagnostics []string
	lines := make([]string, 0)

	// First pass: collect all lines
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	// Second pass: parse with look-ahead for diagnostics
	for i, line := range lines {
		// Version line
		if matches := tapVersionRegex.FindStringSubmatch(line); matches != nil {
			result.Version = matches[1]
			continue
		}

		// Plan line
		if matches := tapPlanRegex.FindStringSubmatch(line); matches != nil {
			start, _ := strconv.Atoi(matches[1])
			end, _ := strconv.Atoi(matches[2])
			plan := &models.TAPPlan{
				Start: start,
				End:   end,
			}
			if len(matches) > 3 && matches[3] != "" {
				plan.Skip = matches[3]
			}
			result.Plan = plan
			continue
		}

		// Bail out line
		if matches := tapBailOutRegex.FindStringSubmatch(line); matches != nil {
			reason := matches[1]
			result.BailOut = &reason
			continue
		}

		// Test line
		if matches := tapTestRegex.FindStringSubmatch(line); matches != nil {
			test := models.TAPTest{
				OK:          matches[1] == "",
				Raw:         line,
				Diagnostics: make([]string, len(currentDiagnostics)),
			}
			copy(test.Diagnostics, currentDiagnostics)
			currentDiagnostics = currentDiagnostics[:0]

			if matches[2] != "" {
				test.Number, _ = strconv.Atoi(matches[2])
			}

			if matches[3] != "" {
				desc := strings.TrimSpace(matches[3])
				// Remove common TAP prefixes like "- "
				if strings.HasPrefix(desc, "- ") {
					desc = strings.TrimSpace(desc[2:])
				}
				test.Description = desc
			}

			if matches[4] != "" {
				directive := &models.TAPDirective{
					Type: strings.ToUpper(matches[4]),
				}
				if matches[5] != "" {
					directive.Reason = matches[5]
				}
				test.Directive = directive
			}

			// Look ahead for diagnostics that follow this test
			j := i + 1
			var followingDiagnostics []string
			for j < len(lines) {
				nextLine := lines[j]
				if matches := tapDiagRegex.FindStringSubmatch(nextLine); matches != nil {
					followingDiagnostics = append(followingDiagnostics, matches[1])
					j++
				} else if tapTestRegex.MatchString(nextLine) || tapPlanRegex.MatchString(nextLine) || tapBailOutRegex.MatchString(nextLine) {
					// Stop if we hit another test or structural element
					break
				} else {
					// Unrecognized line - treat as diagnostic
					followingDiagnostics = append(followingDiagnostics, nextLine)
					j++
				}
			}

			// Combine pre and post diagnostics
			if len(followingDiagnostics) > 0 {
				allDiagnostics := make([]string, len(test.Diagnostics)+len(followingDiagnostics))
				copy(allDiagnostics, test.Diagnostics)
				copy(allDiagnostics[len(test.Diagnostics):], followingDiagnostics)
				test.Diagnostics = allDiagnostics
			}

			result.TestLines = append(result.TestLines, test)
			continue
		}

		// Diagnostic line (only if not following a test)
		if matches := tapDiagRegex.FindStringSubmatch(line); matches != nil {
			diagnostic := matches[1]
			currentDiagnostics = append(currentDiagnostics, diagnostic)
			result.Diagnostics = append(result.Diagnostics, diagnostic)
			continue
		}

		// Unrecognized line - treat as diagnostic (only if not following a test)
		currentDiagnostics = append(currentDiagnostics, line)
		result.Diagnostics = append(result.Diagnostics, line)
	}

	return result, nil
}

func convertTAPToUnified(tapResult *models.TAPResult) *models.UnifiedTestResult {
	result := &models.UnifiedTestResult{
		Metadata: models.TestMetadata{
			Source:    "tap",
			Timestamp: time.Now(),
			Version:   "1.0",
			Framework: "TAP",
		},
		TestSuites: make([]models.TestSuite, 1),
	}

	if tapResult.Version != "" {
		result.Metadata.Environment = map[string]string{
			"tap_version": tapResult.Version,
		}
	}

	suite := models.TestSuite{
		Name:    "TAP Tests",
		Package: "tap",
		Tests:   make([]models.TestCase, 0, len(tapResult.TestLines)),
	}

	var summary models.TestSummary

	for i, tapTest := range tapResult.TestLines {
		testCase := models.TestCase{
			Name:      tapTest.Description,
			ClassName: "tap.test",
		}

		if testCase.Name == "" {
			testCase.Name = fmt.Sprintf("Test %d", i+1)
		}

		// Determine status
		if tapTest.Directive != nil {
			switch tapTest.Directive.Type {
			case "SKIP":
				testCase.Status = models.StatusSkipped
				testCase.Message = tapTest.Directive.Reason
				suite.Skipped++
				summary.Skipped++
			case "TODO":
				if tapTest.OK {
					testCase.Status = models.StatusPassed
					testCase.Message = "TODO test unexpectedly passed: " + tapTest.Directive.Reason
					summary.Passed++
				} else {
					testCase.Status = models.StatusSkipped
					testCase.Message = tapTest.Directive.Reason
					suite.Skipped++
					summary.Skipped++
				}
			}
		} else if tapTest.OK {
			testCase.Status = models.StatusPassed
			summary.Passed++
		} else {
			testCase.Status = models.StatusFailed
			suite.Failures++
			summary.Failed++
		}

		// Add diagnostics as details
		if len(tapTest.Diagnostics) > 0 {
			testCase.Details = strings.Join(tapTest.Diagnostics, "\n")
		}

		// Add properties
		testCase.Properties = map[string]string{
			"tap_number": fmt.Sprintf("%d", tapTest.Number),
			"tap_raw":    tapTest.Raw,
		}

		suite.Tests = append(suite.Tests, testCase)
		summary.Total++
	}

	// Handle bail out
	if tapResult.BailOut != nil {
		bailOutTest := models.TestCase{
			Name:      "Bail Out",
			ClassName: "tap.bailout",
			Status:    models.StatusError,
			Message:   *tapResult.BailOut,
		}
		suite.Tests = append(suite.Tests, bailOutTest)
		suite.Errors++
		summary.Errors++
		summary.Total++
	}

	result.TestSuites[0] = suite
	result.Summary = summary

	return result
}