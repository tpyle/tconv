package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
	"github.com/tpyle/tconv/internal/converter"
)

func main() {
	var inputFiles []string
	var inputType, outputFile string
	var showHelp, showVersion bool

	pflag.StringSliceVarP(&inputFiles, "input", "i", []string{}, "Input file paths (can specify multiple)")
	pflag.StringVarP(&inputType, "type", "t", "", "Input type (junit, postman, gotest, tap, xunit, pytest, testng) - optional for auto-detection")
	pflag.StringVarP(&outputFile, "output", "o", "", "Output file path")
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
		fmt.Printf("Convert test output files between different formats into a unified JSON format.\n")
		fmt.Printf("Multiple input files of different types will be combined into a single output.\n")
		fmt.Printf("File types are auto-detected when --type is not specified.\n\n")
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
		fmt.Printf("\nExamples:\n")
		fmt.Printf("  # Single file\n")
		fmt.Printf("  %s --input test-results.xml --type junit --output unified.json\n", os.Args[0])
		fmt.Printf("  %s -i results.tap -t tap -o output.json\n", os.Args[0])
		fmt.Printf("\n  # Multiple files of same type\n")
		fmt.Printf("  %s --input file1.xml --input file2.xml --type junit --output combined.json\n", os.Args[0])
		fmt.Printf("  %s -i file1.xml -i file2.xml -t junit -o combined.json\n", os.Args[0])
		fmt.Printf("  %s --type junit --output combined.json file1.xml file2.xml file3.xml\n", os.Args[0])
		fmt.Printf("\n  # Multiple files of different types (auto-detected)\n")
		fmt.Printf("  %s --output combined.json results.xml gotest.json postman.json\n", os.Args[0])
		fmt.Printf("  %s -i junit.xml -i pytest.json -i tap.tap -o combined.json\n", os.Args[0])
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("tconv version 1.0.0\n")
		fmt.Printf("Test format converter supporting JUnit, Postman, Go test, TAP, xUnit, pytest, and TestNG\n")
		fmt.Printf("Now supports combining multiple input files into a single output\n")
		os.Exit(0)
	}

	if len(inputFiles) == 0 || outputFile == "" {
		fmt.Fprintf(os.Stderr, "Error: Missing required arguments\n\n")
		fmt.Fprintf(os.Stderr, "Usage: %s --input <file(s)> [--type <type>] --output <file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "       %s [--type <type>] --output <file> file1 file2 ...\n", os.Args[0])
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