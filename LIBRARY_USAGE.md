# tconv Library Usage Guide

The `tconv` project provides both a command-line tool and a Go library for converting test results from various formats into a tiki JSON representation.

## Library Overview

The library is designed to be modular and extensible, with clear interfaces that make it suitable for integration into other programs.

### Key Packages

- **`pkg/models`**: Core data structures and the tiki test result format
- **`pkg/detector`**: File type detection for automatic format identification  
- **`pkg/parsers`**: Format-specific parsers for different test frameworks
- **`pkg/converter`**: High-level conversion and combination logic
- **`pkg/core`**: Advanced interfaces and implementations (optional)
- **`pkg/validation`**: Input validation utilities

## Basic Usage

### Import the Library

```go
import (
    "github.com/tpyle/tconv/pkg/models"
    "github.com/tpyle/tconv/pkg/detector"
    "github.com/tpyle/tconv/pkg/converter"
    "github.com/tpyle/tconv/pkg/parsers"
)
```

### Convert a Single File

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/tpyle/tconv/pkg/converter"
)

func main() {
    conv := converter.New()
    
    // Convert a JUnit XML file to tiki format
    err := conv.Convert("test-results.xml", "junit", "output.json")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Conversion completed!")
}
```

### Auto-detect File Type

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/tpyle/tconv/pkg/detector"
    "github.com/tpyle/tconv/pkg/parsers"
)

func main() {
    // Detect file type automatically
    fileType, err := detector.DetectFileType("unknown-format.xml")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Detected file type: %s\n", fileType)
    
    // Parse using the detected type
    var result *models.UnifiedTestResult
    switch fileType {
    case "junit":
        result, err = parsers.ParseJUnit("unknown-format.xml")
    case "testng":
        result, err = parsers.ParseTestNG("unknown-format.xml")
    // ... handle other types
    }
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Parsed %d test suites\n", len(result.TestSuites))
}
```

### Combine Multiple Files

```go
package main

import (
    "log"
    
    "github.com/tpyle/tconv/pkg/converter"
)

func main() {
    conv := converter.New()
    
    // Combine multiple files of different types
    files := []string{
        "junit-results.xml",
        "gotest-output.json", 
        "pytest-report.json",
    }
    
    // Auto-detect types and combine
    err := conv.ConvertMultiple(files, "", "combined-output.json")
    if err != nil {
        log.Fatal(err)
    }
}
```

### Working with the Tiki Format

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
    
    "github.com/tpyle/tconv/pkg/models"
    "github.com/tpyle/tconv/pkg/parsers"
)

func main() {
    // Parse a test file
    result, err := parsers.ParseJUnit("test-results.xml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Access the tiki data
    fmt.Printf("Test run from: %s\n", result.Metadata.Source)
    fmt.Printf("Framework: %s\n", result.Metadata.Framework)
    fmt.Printf("Total tests: %d\n", result.Summary.Total)
    fmt.Printf("Pass rate: %.1f%%\n", result.Summary.PassRate())
    
    // Iterate through test suites
    for _, suite := range result.TestSuites {
        fmt.Printf("Suite: %s (%d tests)\n", suite.Name, len(suite.Tests))
        
        for _, test := range suite.Tests {
            fmt.Printf("  - %s: %s\n", test.Name, test.Status)
            if test.Status == models.StatusFailed {
                fmt.Printf("    Error: %s\n", test.Message)
            }
        }
    }
    
    // Export to JSON with custom formatting
    file, err := os.Create("pretty-output.json")
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    
    encoder := json.NewEncoder(file)
    encoder.SetIndent("", "  ")
    encoder.Encode(result)
}
```

## Advanced Usage

### Using Core Interfaces (Advanced)

The `pkg/core` package provides more advanced interfaces for extensibility:

```go
package main

import (
    "bytes"
    "log"
    
    "github.com/tpyle/tconv/pkg/core"
)

func main() {
    // Use the type detector interface
    detector := core.NewDefaultTypeDetector()
    
    // Detect from file
    fileType, err := detector.DetectTypeFromFile("test.xml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Combine results with custom options
    combiner := core.NewDefaultResultCombiner()
    options := core.CombineOptions{
        Source:             "my-combined-tests",
        PreserveSuiteNames: true,
        MergeEnvironment:   true,
        NamePrefix:         "batch1",
    }
    
    // Assuming you have multiple results...
    // combined, err := combiner.Combine(results, options)
    
    // Export with custom format
    exporter := core.NewDefaultExporter()
    exportOptions := core.ExportOptions{
        Format:          "yaml", // json, yaml, xml
        Indent:          "  ",
        IncludeMetadata: true,
    }
    
    var buf bytes.Buffer
    // err = exporter.Export(combined, &buf, exportOptions)
}
```

### Custom Error Handling

```go
package main

import (
    "errors"
    "fmt"
    "log"
    
    "github.com/tpyle/tconv/pkg/core"
    "github.com/tpyle/tconv/pkg/detector"
)

func main() {
    _, err := detector.DetectFileType("nonexistent.xml")
    if err != nil {
        // Check if it's a structured tconv error
        var tconvErr *core.TconvError
        if errors.As(err, &tconvErr) {
            fmt.Printf("Error type: %s\n", tconvErr.Type)
            fmt.Printf("Message: %s\n", tconvErr.Message)
            
            // Print context information
            for key, value := range tconvErr.Context {
                fmt.Printf("Context %s: %s\n", key, value)
            }
            
            // Print suggestions
            for _, suggestion := range tconvErr.Suggestions {
                fmt.Printf("Suggestion: %s\n", suggestion)
            }
        } else {
            log.Fatal(err)
        }
    }
}
```

## Data Structure Reference

### UnifiedTestResult

The main data structure representing converted test results:

```go
type UnifiedTestResult struct {
    Metadata   TestMetadata `json:"metadata"`    // Test run information
    TestSuites []TestSuite  `json:"test_suites"` // All test suites
    Summary    TestSummary  `json:"summary"`     // Aggregate statistics
}
```

### TestCase Status Values

```go
const (
    StatusPassed  TestStatus = "passed"  // Test completed successfully
    StatusFailed  TestStatus = "failed"  // Test failed assertion(s)  
    StatusSkipped TestStatus = "skipped" // Test was not executed
    StatusError   TestStatus = "error"   // Test encountered a runtime error
)
```

### Helper Methods

The `TestSummary` type provides useful calculated fields:

```go
summary := result.Summary
passRate := summary.PassRate() // Returns percentage as float64
```

The `TestStatus` type provides validation:

```go
status := models.StatusPassed
if status.IsValid() {
    fmt.Println("Status is valid")
}
```

## Supported Formats

- **JUnit XML**: Standard JUnit test reports
- **TestNG XML**: TestNG framework reports  
- **xUnit XML**: .NET xUnit framework reports
- **Go Test JSON**: Go's `go test -json` output
- **Postman JSON**: Postman collection runner results
- **pytest JSON**: pytest with JSON reporting plugin
- **TAP**: Test Anything Protocol format

## Error Types

The library provides structured error handling with these error types:

- `unsupported_format`: Unknown or unsupported file format
- `invalid_input`: Malformed or invalid input data
- `parsing_failed`: Error during format-specific parsing
- `file_not_found`: Input file does not exist
- `validation_failed`: Input validation errors
- `conversion_failed`: General conversion errors
- `combination_failed`: Error combining multiple results
- `export_failed`: Error during output generation

## Best Practices

1. **Always check error types** using `errors.As()` for better error handling
2. **Use auto-detection** when processing files of unknown formats
3. **Validate inputs** before processing with the validation package
4. **Handle nil results** when combining files (some may fail to parse)
5. **Set appropriate timeouts** for file operations in production use
6. **Use structured logging** to track conversion operations

## Thread Safety

All exported functions and types in the library are designed to be thread-safe and stateless, making them suitable for concurrent use in web servers and batch processing applications.