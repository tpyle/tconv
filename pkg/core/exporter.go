package core

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"gopkg.in/yaml.v3"
)

// DefaultExporter implements the TestResultExporter interface.
type DefaultExporter struct{}

// NewDefaultExporter creates a new instance of DefaultExporter.
func NewDefaultExporter() *DefaultExporter {
	return &DefaultExporter{}
}

// Export writes the unified test result to the specified output.
func (e *DefaultExporter) Export(result *UnifiedTestResult, output io.Writer, options ExportOptions) error {
	if result == nil {
		return NewInvalidInputError("result cannot be nil")
	}

	if output == nil {
		return NewInvalidInputError("output writer cannot be nil")
	}

	// Set default format if not specified
	format := options.Format
	if format == "" {
		format = "json"
	}

	// Set default indent if not specified
	indent := options.Indent
	if indent == "" {
		indent = "  "
	}

	// Create a copy of the result if metadata should be excluded
	exportResult := result
	if !options.IncludeMetadata {
		exportResult = &UnifiedTestResult{
			TestSuites: result.TestSuites,
			Summary:    result.Summary,
		}
	}

	switch format {
	case "json":
		return e.exportJSON(exportResult, output, indent)
	case "yaml", "yml":
		return e.exportYAML(exportResult, output, indent)
	case "xml":
		return e.exportXML(exportResult, output, indent)
	default:
		return NewUnsupportedFormatError(format, []string{"json", "yaml", "xml"})
	}
}

// exportJSON exports the result as JSON.
func (e *DefaultExporter) exportJSON(result *UnifiedTestResult, output io.Writer, indent string) error {
	encoder := json.NewEncoder(output)
	encoder.SetIndent("", indent)
	if err := encoder.Encode(result); err != nil {
		return NewExportError("json", err)
	}
	return nil
}

// exportYAML exports the result as YAML.
func (e *DefaultExporter) exportYAML(result *UnifiedTestResult, output io.Writer, indent string) error {
	encoder := yaml.NewEncoder(output)
	defer encoder.Close()

	// Configure YAML encoder
	encoder.SetIndent(len(indent))

	if err := encoder.Encode(result); err != nil {
		return NewExportError("yaml", err)
	}
	return nil
}

// exportXML exports the result as XML.
func (e *DefaultExporter) exportXML(result *UnifiedTestResult, output io.Writer, indent string) error {
	// Write XML header
	if _, err := output.Write([]byte(xml.Header)); err != nil {
		return NewExportError("xml", err)
	}

	encoder := xml.NewEncoder(output)
	if indent != "" {
		encoder.Indent("", indent)
	}

	// Create a wrapper struct for better XML representation
	wrapper := struct {
		XMLName xml.Name `xml:"unified-test-result"`
		*UnifiedTestResult
	}{
		UnifiedTestResult: result,
	}

	if err := encoder.Encode(wrapper); err != nil {
		return NewExportError("xml", err)
	}

	return nil
}

// ExportToJSON is a convenience function for JSON export with default options.
func ExportToJSON(result *UnifiedTestResult, output io.Writer) error {
	exporter := NewDefaultExporter()
	options := ExportOptions{
		Format:          "json",
		Indent:          "  ",
		IncludeMetadata: true,
	}
	return exporter.Export(result, output, options)
}

// ExportToYAML is a convenience function for YAML export with default options.
func ExportToYAML(result *UnifiedTestResult, output io.Writer) error {
	exporter := NewDefaultExporter()
	options := ExportOptions{
		Format:          "yaml",
		Indent:          "  ",
		IncludeMetadata: true,
	}
	return exporter.Export(result, output, options)
}

// ExportToXML is a convenience function for XML export with default options.
func ExportToXML(result *UnifiedTestResult, output io.Writer) error {
	exporter := NewDefaultExporter()
	options := ExportOptions{
		Format:          "xml",
		Indent:          "  ",
		IncludeMetadata: true,
	}
	return exporter.Export(result, output, options)
}