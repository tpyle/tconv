package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/tpyle/tconv/internal/converter"
	"github.com/tpyle/tconv/internal/exporters"
)

func main() {
	var inputFiles []string
	var inputType, outputFile, exportFormat string
	var showHelp, showVersion, exportMode bool

	pflag.StringSliceVarP(&inputFiles, "input", "i", []string{}, "Input file paths (can specify multiple)")
	pflag.StringVarP(&inputType, "type", "t", "", "Input type (junit, postman, gotest, tap, xunit, pytest, testng) - optional for auto-detection")
	pflag.StringVarP(&outputFile, "output", "o", "", "Output file path")
	pflag.StringVarP(&exportFormat, "export", "e", "", "Export tiki JSON to specified format (junit, tap, gotest)")
	pflag.BoolVar(&exportMode, "export-mode", false, "Enable export mode: convert FROM tiki JSON TO other formats")
	pflag.BoolVarP(&showHelp, "help", "h", false, "Show help message")
	pflag.BoolVarP(&showVersion, "version", "v", false, "Show version information")
	pflag.Parse()

	// Also accept positional arguments as input files
	if len(pflag.Args()) > 0 {
		inputFiles = append(inputFiles, pflag.Args()...)
	}

	if showHelp {
		fmt.Printf("tconv - Test Converter Tool\n\n")
		fmt.Printf("Usage: %s [OPTIONS] [files...]\n\n", os.Args[0])
		fmt.Printf("Convert test output files between different formats into a tiki JSON format.\n")
		fmt.Printf("Multiple input files of different types will be combined into a single output.\n")
		fmt.Printf("File types are auto-detected when --type is not specified.\n")
		fmt.Printf("Can also export tiki JSON back to specific test formats.\n\n")
		fmt.Printf("Options:\n")
		pflag.PrintDefaults()
		fmt.Printf("\nSupported input types:\n")
		fmt.Printf("  junit    - JUnit XML test results\n")
		fmt.Printf("  postman  - Postman collection runner JSON\n")
		fmt.Printf("  gotest   - Go test JSON output\n")
		fmt.Printf("  tap      - Test Anything Protocol\n")
		fmt.Printf("  xunit    - xUnit.net XML test results\n")
		fmt.Printf("  pytest   - pytest JSON report\n")
		fmt.Printf("  testng   - TestNG XML test results\n")
		fmt.Printf("\nSupported export formats:\n")
		fmt.Printf("  junit    - JUnit XML test results\n")
		fmt.Printf("  tap      - Test Anything Protocol\n")
		fmt.Printf("  gotest   - Go test JSON output\n")
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  # Import: Single file\n")
		fmt.Printf("  %s --input test-results.xml --type junit --output unified.json\n", os.Args[0])
		fmt.Printf("  %s -i results.tap -t tap -o output.json\n", os.Args[0])
		fmt.Printf("\n  # Import: Multiple files of different types (auto-detected)\n")
		fmt.Printf("  %s --output combined.json results.xml gotest.json postman.json\n", os.Args[0])
		fmt.Printf("  %s -i junit.xml -i pytest.json -i tap.tap -o combined.json\n", os.Args[0])
		fmt.Printf("\n  # Export: Convert tiki JSON to other formats\n")
		fmt.Printf("  %s --export junit --input tiki.json --output results.xml\n", os.Args[0])
		fmt.Printf("  %s -e tap -i combined.json -o output.tap\n", os.Args[0])
		fmt.Printf("  %s --export gotest --input tiki.json --output test.json\n", os.Args[0])
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("tconv version 1.0.0\n")
		fmt.Printf("Test format converter supporting JUnit, Postman, Go test, TAP, xUnit, pytest, and TestNG\n")
		fmt.Printf("Now supports combining multiple input files into a single output\n")
		os.Exit(0)
	}

	// Check for export mode
	if exportFormat != "" {
		// Export mode: convert FROM tiki JSON TO specific format
		if len(inputFiles) != 1 {
			fmt.Fprintf(os.Stderr, "Error: Export mode requires exactly one input file (tiki JSON)\n")
			os.Exit(1)
		}
		if outputFile == "" {
			fmt.Fprintf(os.Stderr, "Error: Output file is required for export\n")
			os.Exit(1)
		}

		// Load tiki JSON
		conv := converter.New()
		result, err := conv.LoadTiki(inputFiles[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error loading tiki JSON: %v\n", err)
			os.Exit(1)
		}

		// Export to target format
		manager := exporters.NewExportManager()
		if err := manager.Export(result, exportFormat, outputFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error exporting to %s: %v\n", exportFormat, err)
			os.Exit(1)
		}

		fmt.Printf("Successfully exported %s to %s (%s format)\n", inputFiles[0], outputFile, exportFormat)
		return
	}

	// Regular conversion mode: convert TO tiki JSON
	if len(inputFiles) == 0 || outputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Missing required arguments\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s --input <file(s)> [--type <type>] --output <file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s [--type <type>] --output <file> file1 file2 ...\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s --export <format> --input <tiki.json> --output <file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Use --help for more information\n")
		os.Exit(1)
	}

	conv := converter.New()
	if err := conv.ConvertMultiple(inputFiles, inputType, outputFile); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(inputFiles) == 1 {
		fmt.Printf("Successfully converted %s to %s\n", inputFiles[0], outputFile)
	} else {
		fmt.Printf("Successfully combined %d files into %s\n", len(inputFiles), outputFile)
	}
}