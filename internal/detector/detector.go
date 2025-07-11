// Package detector provides file type detection capabilities for test result files.
// It can automatically identify various test formats based on file extensions and content analysis.
package detector

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DetectFileType attempts to determine the test file format based on extension and content.
// It returns one of: "junit", "testng", "xunit", "postman", "pytest", "gotest", "tap"
func DetectFileType(filePath string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", filePath)
	}

	// First try to detect by file extension
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".xml":
		return detectXMLType(filePath)
	case ".json":
		return detectJSONType(filePath)
	case ".tap":
		return "tap", nil
	}

	// If extension doesn't help, try content detection
	return detectByContent(filePath)
}

// detectXMLType analyzes XML content to determine the specific XML test format.
func detectXMLType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read first few KB to detect XML type
	buf := make([]byte, 4096)
	n, err := file.Read(buf)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	content := strings.ToLower(string(buf[:n]))

	// Check for specific XML types in order of specificity
	if strings.Contains(content, "<testng-results") {
		return "testng", nil
	}
	if strings.Contains(content, "<assemblies") || strings.Contains(content, "<assembly") {
		return "xunit", nil
	}
	if strings.Contains(content, "<testsuite") || strings.Contains(content, "<testsuites") {
		return "junit", nil
	}
	if strings.Contains(content, "<test-results") || strings.Contains(content, "<test-run") {
		return "xunit", nil
	}

	// Default to junit for unrecognized XML files
	return "junit", nil
}

// detectJSONType analyzes JSON content to determine the specific JSON test format.
func detectJSONType(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Try to parse as JSON object first
	var data map[string]interface{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		// If single object parsing fails, try line-by-line (Go test format)
		file.Seek(0, 0)
		return detectGoTestFormat(file)
	}

	// Check for postman collection structure
	if info, ok := data["info"]; ok {
		if infoMap, ok := info.(map[string]interface{}); ok {
			if _, hasName := infoMap["name"]; hasName {
				if _, hasSchema := infoMap["schema"]; hasSchema {
					return "postman", nil
				}
			}
		}
	}

	// Check for pytest structure
	if _, hasTests := data["tests"]; hasTests {
		if _, hasSummary := data["summary"]; hasSummary {
			return "pytest", nil
		}
	}

	// Check for go test structure (single JSON object with Action field)
	if _, hasAction := data["Action"]; hasAction {
		return "gotest", nil
	}

	return "", fmt.Errorf("unable to determine JSON test format")
}

// detectGoTestFormat checks for Go test JSON format (line-delimited JSON).
func detectGoTestFormat(file *os.File) (string, error) {
	scanner := bufio.NewScanner(file)
	lineCount := 0
	gotestLines := 0

	for scanner.Scan() && lineCount < 10 {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var lineData map[string]interface{}
		if json.Unmarshal([]byte(line), &lineData) == nil {
			if _, hasAction := lineData["Action"]; hasAction {
				gotestLines++
			}
		}
		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan file: %w", err)
	}

	if gotestLines > 0 {
		return "gotest", nil
	}

	return "", fmt.Errorf("unable to determine JSON test format")
}

// detectByContent analyzes content without extension hints.
func detectByContent(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0

	for scanner.Scan() && lineCount < 10 {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		lineLower := strings.ToLower(line)

		// Check for TAP format
		if strings.HasPrefix(lineLower, "tap version") ||
			strings.HasPrefix(lineLower, "1..") ||
			(strings.HasPrefix(lineLower, "ok ") || strings.HasPrefix(lineLower, "not ok ")) {
			return "tap", nil
		}

		// Check for XML
		if strings.HasPrefix(lineLower, "<?xml") || strings.Contains(lineLower, "<testsuite") {
			// Reset and analyze as XML
			file.Seek(0, 0)
			return detectXMLType(filePath)
		}

		// Check for JSON
		if strings.HasPrefix(line, "{") || strings.HasPrefix(line, "[") {
			// Reset and analyze as JSON
			file.Seek(0, 0)
			return detectJSONType(filePath)
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("failed to scan file: %w", err)
	}

	return "", fmt.Errorf("unable to detect file format for %s", filePath)
}