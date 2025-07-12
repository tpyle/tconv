package exporters

import (
	"fmt"
	"io"
	"strings"

	"github.com/tpyle/tconv/pkg/models"
)

// Exporter defines the interface for format-specific exporters
type Exporter interface {
	Export(result *models.UnifiedTestResult, outputPath string) error
	Write(result *models.UnifiedTestResult, writer io.Writer) error
	FormatName() string
}

// ExportManager handles export operations for different formats
type ExportManager struct {
	exporters map[string]Exporter
}

// NewExportManager creates a new export manager with all supported exporters
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

// RegisterExporter adds a new exporter for a specific format
func (em *ExportManager) RegisterExporter(format string, exporter Exporter) {
	em.exporters[strings.ToLower(format)] = exporter
}

// Export converts and writes the tiki result to the specified format
func (em *ExportManager) Export(result *models.UnifiedTestResult, format, outputPath string) error {
	exporter, exists := em.exporters[strings.ToLower(format)]
	if !exists {
		return fmt.Errorf("unsupported export format: %s", format)
	}

	return exporter.Export(result, outputPath)
}

// Write converts and writes the tiki result to the specified format using a writer
func (em *ExportManager) Write(result *models.UnifiedTestResult, format string, writer io.Writer) error {
	exporter, exists := em.exporters[strings.ToLower(format)]
	if !exists {
		return fmt.Errorf("unsupported export format: %s", format)
	}

	return exporter.Write(result, writer)
}

// SupportedExportFormats returns a list of all supported export formats
func (em *ExportManager) SupportedExportFormats() []string {
	formats := make([]string, 0, len(em.exporters))
	for format := range em.exporters {
		formats = append(formats, format)
	}
	return formats
}

// Concrete exporter implementations

// JUnitExporter exports to JUnit XML format
type JUnitExporter struct{}

func (e *JUnitExporter) Export(result *models.UnifiedTestResult, outputPath string) error {
	return ExportJUnit(result, outputPath)
}

func (e *JUnitExporter) Write(result *models.UnifiedTestResult, writer io.Writer) error {
	return WriteJUnit(result, writer)
}

func (e *JUnitExporter) FormatName() string {
	return "junit"
}

// TAPExporter exports to TAP format
type TAPExporter struct{}

func (e *TAPExporter) Export(result *models.UnifiedTestResult, outputPath string) error {
	return ExportTAP(result, outputPath)
}

func (e *TAPExporter) Write(result *models.UnifiedTestResult, writer io.Writer) error {
	return WriteTAP(result, writer)
}

func (e *TAPExporter) FormatName() string {
	return "tap"
}

// GoTestExporter exports to Go test JSON format
type GoTestExporter struct{}

func (e *GoTestExporter) Export(result *models.UnifiedTestResult, outputPath string) error {
	return ExportGoTest(result, outputPath)
}

func (e *GoTestExporter) Write(result *models.UnifiedTestResult, writer io.Writer) error {
	return WriteGoTest(result, writer)
}

func (e *GoTestExporter) FormatName() string {
	return "gotest"
}