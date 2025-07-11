package core

import (
	"bufio"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DefaultTypeDetector implements the TypeDetector interface.
type DefaultTypeDetector struct {
	supportedExtensions map[string]string
}

// NewDefaultTypeDetector creates a new instance of DefaultTypeDetector.
func NewDefaultTypeDetector() *DefaultTypeDetector {
	return &DefaultTypeDetector{
		supportedExtensions: map[string]string{
			".xml":  "xml", // Will be refined by content analysis
			".json": "json", // Will be refined by content analysis
			".tap":  "tap",
		},
	}
}

// DetectType analyzes the content and returns the detected type.
func (d *DefaultTypeDetector) DetectType(content io.Reader) (string, error) {
	// Read the beginning of the content for analysis
	buf := make([]byte, 4096)
	n, err := content.Read(buf)
	if err != nil && err != io.EOF {
		return "", NewParsingError("unknown", err)
	}

	if n == 0 {
		return "", NewInvalidInputError("empty content")
	}

	data := string(buf[:n])
	return d.analyzeContent(data)
}

// DetectTypeFromFile analyzes a file and returns the detected type.
func (d *DefaultTypeDetector) DetectTypeFromFile(filePath string) (string, error) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", NewFileNotFoundError(filePath)
	}

	// First try extension-based detection
	ext := strings.ToLower(filepath.Ext(filePath))
	if baseType, exists := d.supportedExtensions[ext]; exists {
		if baseType == "tap" {
			return "tap", nil // TAP files are unambiguous
		}
		// For XML and JSON, we need content analysis
		return d.detectFromFileContent(filePath, baseType)
	}

	// If no extension match, try content analysis
	return d.detectFromFileContent(filePath, "")
}

// SupportedExtensions returns a map of file extensions to likely types.
func (d *DefaultTypeDetector) SupportedExtensions() map[string]string {
	// Return a copy to prevent external modification
	result := make(map[string]string)
	for k, v := range d.supportedExtensions {
		result[k] = v
	}
	return result
}

// detectFromFileContent analyzes file content to determine the specific type.
func (d *DefaultTypeDetector) detectFromFileContent(filePath, baseType string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", NewFileNotFoundError(filePath)
	}
	defer file.Close()

	switch baseType {
	case "xml":
		return d.detectXMLType(file)
	case "json":
		return d.detectJSONType(file)
	default:
		return d.detectByContent(file)
	}
}

// detectXMLType analyzes XML content to determine the specific XML test format.
func (d *DefaultTypeDetector) detectXMLType(reader io.Reader) (string, error) {
	buf := make([]byte, 4096)
	n, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		return "", NewParsingError("xml", err)
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
func (d *DefaultTypeDetector) detectJSONType(reader io.Reader) (string, error) {
	// Reset reader if it's a file
	if seeker, ok := reader.(io.Seeker); ok {
		seeker.Seek(0, 0)
	}

	// Try to parse as JSON object first
	var data map[string]interface{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		// If single object parsing fails, try line-by-line (Go test format)
		if seeker, ok := reader.(io.Seeker); ok {
			seeker.Seek(0, 0)
		}
		return d.detectGoTestFormat(reader)
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

	return "", NewUnsupportedFormatError("unknown JSON format", []string{"postman", "pytest", "gotest"})
}

// detectGoTestFormat checks for Go test JSON format (line-delimited JSON).
func (d *DefaultTypeDetector) detectGoTestFormat(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
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
		return "", NewParsingError("json", err)
	}

	if gotestLines > 0 {
		return "gotest", nil
	}

	return "", NewUnsupportedFormatError("unknown JSON format", []string{"postman", "pytest", "gotest"})
}

// detectByContent analyzes content without extension hints.
func (d *DefaultTypeDetector) detectByContent(reader io.Reader) (string, error) {
	scanner := bufio.NewScanner(reader)
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
			// Reset reader and analyze as XML
			if seeker, ok := reader.(io.Seeker); ok {
				seeker.Seek(0, 0)
			}
			return d.detectXMLType(reader)
		}

		// Check for JSON
		if strings.HasPrefix(line, "{") || strings.HasPrefix(line, "[") {
			// Reset reader and analyze as JSON
			if seeker, ok := reader.(io.Seeker); ok {
				seeker.Seek(0, 0)
			}
			return d.detectJSONType(reader)
		}

		lineCount++
	}

	if err := scanner.Err(); err != nil {
		return "", NewParsingError("unknown", err)
	}

	return "", NewUnsupportedFormatError("unknown format", []string{"junit", "testng", "xunit", "postman", "pytest", "gotest", "tap"})
}

// analyzeContent is a helper method for analyzing content from a byte slice.
func (d *DefaultTypeDetector) analyzeContent(content string) (string, error) {
	reader := strings.NewReader(content)
	return d.detectByContent(reader)
}