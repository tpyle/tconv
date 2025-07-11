// This file demonstrates how external libraries can use tconv
// +build ignore

package main

import (
	"fmt"
	"log"

	"github.com/tpyle/tconv"
	"github.com/tpyle/tconv/pkg/models"
)

func main() {
	// Example 1: Convert a single file
	fmt.Println("=== Single File Conversion ===")
	err := tconv.Convert("testdata/sample_junit.xml", "junit", "example_output.json")
	if err != nil {
		log.Printf("Conversion failed: %v", err)
	} else {
		fmt.Println("✓ Successfully converted JUnit file")
	}

	// Example 2: Auto-detect and combine multiple files
	fmt.Println("\n=== Multiple File Conversion (Auto-detect) ===")
	files := []string{
		"testdata/sample_junit.xml",
		"testdata/sample_gotest.json",
		"testdata/sample_tap.tap",
	}
	
	err = tconv.ConvertMultiple(files, "", "example_combined.json")
	if err != nil {
		log.Printf("Combination failed: %v", err)
	} else {
		fmt.Println("✓ Successfully combined multiple file types")
	}

	// Example 3: Detect file type
	fmt.Println("\n=== File Type Detection ===")
	fileType, err := tconv.DetectFileType("testdata/sample_junit.xml")
	if err != nil {
		log.Printf("Detection failed: %v", err)
	} else {
		fmt.Printf("✓ Detected file type: %s\n", fileType)
	}

	// Example 4: Work with unified data structures
	fmt.Println("\n=== Working with Unified Data ===")
	
	// These types are now available as public API
	var result *tconv.UnifiedTestResult
	
	// You can also use the models package directly
	var modelResult *models.UnifiedTestResult
	
	// Status constants are available
	fmt.Printf("Available statuses: %v, %v, %v, %v\n", 
		tconv.StatusPassed, tconv.StatusFailed, tconv.StatusSkipped, tconv.StatusError)
	
	// Show supported formats
	fmt.Println("\n=== Supported Formats ===")
	formats := tconv.SupportedFormats()
	for _, format := range formats {
		fmt.Printf("- %s\n", format)
	}
	
	_ = result      // Prevent unused variable error
	_ = modelResult // Prevent unused variable error
}