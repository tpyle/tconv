// Package exporters provides functionality to export tiki test results
// to various target formats.
//
// This package supports exporting to multiple formats including JUnit XML,
// TAP (Test Anything Protocol), and Go test JSON. The ExportManager
// coordinates format-specific exporters and provides a unified interface.
//
// Each format-specific exporter implements the Exporter interface,
// allowing for easy extension with additional export formats.
package exporters

import (
	"fmt"
	"io"
	"strings"

	"github.com/tpyle/tconv/pkg/models"
)

// Exporter defines the interface for format-specific exporters.
//
// All format-specific exporters must implement this interface to provide
// consistent export functionality. Implementations should be thread-safe.
type Exporter interface {
	// Export writes the tiki result to a file in the exporter's format
	Export(result *models.TikiTestResult, outputPath string) error
	
	// Write outputs the tiki result to a writer in the exporter's format
	Write(result *models.TikiTestResult, writer io.Writer) error
	
	// FormatName returns the name of the export format
	FormatName() string
}

// ExportManager handles export operations for different formats.
//
// It maintains a registry of format-specific exporters and provides
// a unified interface for exporting tiki results to various formats.
type ExportManager struct {
	exporters map[string]Exporter
}

// NewExportManager creates a new export manager with all supported exporters.
//
// The manager is pre-configured with exporters for JUnit XML, TAP, and Go test JSON formats.
// Additional exporters can be registered using RegisterExporter.
func NewExportManager() *ExportManager {
	manager := &ExportManager{
		exporters: make(map[string]Exporter),
	}

	// Register all available exporters
	manager.RegisterExporter("junit", &JUnitExporter{})
	manager.RegisterExporter("tap", &TAPExporter{})
	manager.RegisterExporter("gotest", &GoTestExporter{})

	return manager
}

// RegisterExporter adds a new exporter for a specific format.
//
// The format name is converted to lowercase for consistent lookup.
// If an exporter already exists for the format, it will be replaced.
func (em *ExportManager) RegisterExporter(format string, exporter Exporter) {
	em.exporters[strings.ToLower(format)] = exporter
}

// Export converts and writes the tiki result to the specified format.
//
// The format lookup is case-insensitive. Returns an error if the format
// is not supported or if the export operation fails.
func (em *ExportManager) Export(result *models.TikiTestResult, format, outputPath string) error {
	exporter, exists := em.exporters[strings.ToLower(format)]
	if !exists {
		return fmt.Errorf("unsupported export format: %s", format)
	}

	return exporter.Export(result, outputPath)
}

// Write converts and writes the tiki result to the specified format using a writer.
//
// This method is useful for writing to stdout, buffers, or other io.Writer implementations.
// The format lookup is case-insensitive.
func (em *ExportManager) Write(result *models.TikiTestResult, format string, writer io.Writer) error {
	exporter, exists := em.exporters[strings.ToLower(format)]
	if !exists {
		return fmt.Errorf("unsupported export format: %s", format)
	}

	return exporter.Write(result, writer)
}

// SupportedExportFormats returns a list of all supported export formats.
//
// The returned slice contains format names in lowercase as they are stored internally.
func (em *ExportManager) SupportedExportFormats() []string {
	formats := make([]string, 0, len(em.exporters))
	for format := range em.exporters {
		formats = append(formats, format)
	}
	return formats
}

// Concrete exporter implementations

// JUnitExporter exports tiki test results to JUnit XML format.
//
// The generated XML follows the JUnit XML schema and is compatible
// with most CI/CD systems and test result viewers.
type JUnitExporter struct{}

func (e *JUnitExporter) Export(result *models.TikiTestResult, outputPath string) error {
	return ExportJUnit(result, outputPath)
}

func (e *JUnitExporter) Write(result *models.TikiTestResult, writer io.Writer) error {
	return WriteJUnit(result, writer)
}

func (e *JUnitExporter) FormatName() string {
	return "junit"
}

// TAPExporter exports tiki test results to TAP (Test Anything Protocol) format.
//
// TAP is a simple text-based interface between testing modules and
// test harnesses, widely supported by testing frameworks.
type TAPExporter struct{}

func (e *TAPExporter) Export(result *models.TikiTestResult, outputPath string) error {
	return ExportTAP(result, outputPath)
}

func (e *TAPExporter) Write(result *models.TikiTestResult, writer io.Writer) error {
	return WriteTAP(result, writer)
}

func (e *TAPExporter) FormatName() string {
	return "tap"
}

// GoTestExporter exports tiki test results to Go test JSON format.
//
// The output format matches the JSON structure produced by 'go test -json',
// making it compatible with Go testing tools and CI systems.
type GoTestExporter struct{}

func (e *GoTestExporter) Export(result *models.TikiTestResult, outputPath string) error {
	return ExportGoTest(result, outputPath)
}

func (e *GoTestExporter) Write(result *models.TikiTestResult, writer io.Writer) error {
	return WriteGoTest(result, writer)
}

func (e *GoTestExporter) FormatName() string {
	return "gotest"
}