# tconv - Test Output Converter

A Go CLI application that converts multiple popular test output formats into a unified custom JSON format.

## Features

- **Multi-format Support**: Converts 7 popular test formats:
  - JUnit XML
  - Postman collection runner JSON
  - Go test JSON output
  - TAP (Test Anything Protocol)
  - xUnit XML
  - pytest JSON
  - TestNG XML
- **Unified Schema**: Creates a superset JSON format that encompasses all input types
- **CLI Interface**: Simple command-line interface for easy integration
- **Comprehensive Testing**: High test coverage across all parsers and converters

## Installation

```bash
go build -o tconv main.go
```

## Usage

```bash
# Convert JUnit XML output
./tconv -input junit_results.xml -type junit -output unified_results.json

# Convert Postman collection runner output
./tconv -input postman_results.json -type postman -output unified_results.json

# Convert Go test JSON output
./tconv -input gotest_results.json -type gotest -output unified_results.json

# Convert TAP output
./tconv -input tap_results.tap -type tap -output unified_results.json

# Convert xUnit XML output
./tconv -input xunit_results.xml -type xunit -output unified_results.json

# Convert pytest JSON output
./tconv -input pytest_results.json -type pytest -output unified_results.json

# Convert TestNG XML output
./tconv -input testng_results.xml -type testng -output unified_results.json

# Output to stdout (omit -output flag)
./tconv -input junit_results.xml -type junit
```

## Supported Input Formats

### JUnit XML
Standard JUnit XML format with support for:
- Test suites and test cases
- Failures, errors, and skipped tests
- System output and error capture
- Test properties and metadata

### Postman Collection Runner Output
Newman/Postman collection runner JSON output with support for:
- Request/response details
- Assertion results
- Timing information
- Collection metadata

### Go Test JSON Output
Go test `-json` output format with support for:
- Package-level test organization
- Test pass/fail/skip status
- Test output capture
- Timing information

### TAP (Test Anything Protocol)
TAP format with support for:
- Plan lines and version information
- Test results with descriptions
- TODO and SKIP directives
- Diagnostic messages
- Bail out handling

### xUnit XML
xUnit.net XML format with support for:
- Assembly and collection organization
- Test traits and categories
- Detailed failure information
- Timing and environment data

### pytest JSON
pytest JSON report format with support for:
- Test phases (setup, call, teardown)
- Detailed traceback information
- Test metadata and keywords
- Multiple outcome types

### TestNG XML
TestNG XML results format with support for:
- Suite and test organization
- Test groups and parameters
- Method-level timing
- Reporter output and exceptions

## Unified Output Format

The unified JSON format includes:

```json
{
  "metadata": {
    "source": "junit|postman|gotest",
    "timestamp": "2023-01-01T10:00:00Z",
    "version": "1.0",
    "framework": "test framework name",
    "environment": {}
  },
  "test_suites": [
    {
      "name": "suite name",
      "package": "package name",
      "tests": [
        {
          "name": "test name",
          "class_name": "class or package",
          "time": 0.123,
          "status": "passed|failed|skipped|error",
          "message": "test message",
          "details": "detailed output",
          "error_type": "error type",
          "system_out": "stdout",
          "system_err": "stderr",
          "assertions": [
            {
              "name": "assertion name",
              "expected": "expected value",
              "actual": "actual value", 
              "passed": true,
              "message": "assertion message"
            }
          ],
          "properties": {}
        }
      ],
      "errors": 0,
      "failures": 0,
      "skipped": 0,
      "time": 0.456,
      "properties": {}
    }
  ],
  "summary": {
    "total": 10,
    "passed": 8,
    "failed": 1,
    "skipped": 1,
    "errors": 0,
    "duration": 2.34
  }
}
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## Project Structure

```
├── main.go                 # CLI entry point
├── pkg/
│   ├── models/             # Data models for all formats
│   │   ├── unified.go      # Unified output format
│   │   ├── junit.go        # JUnit XML structures
│   │   ├── postman.go      # Postman JSON structures
│   │   └── gotest.go       # Go test JSON structures
│   ├── parsers/            # Format-specific parsers
│   │   ├── junit.go        # JUnit XML parser
│   │   ├── postman.go      # Postman JSON parser
│   │   ├── gotest.go       # Go test JSON parser
│   │   └── *_test.go       # Parser unit tests
│   └── converter/          # Main conversion logic
│       ├── converter.go    # Converter implementation
│       └── converter_test.go # Converter tests
├── testdata/               # Sample input files for testing
│   ├── sample_junit.xml
│   ├── sample_postman.json
│   ├── sample_gotest.json
│   ├── sample_tap.tap
│   ├── sample_xunit.xml
│   ├── sample_pytest.json
│   └── sample_testng.xml
└── README.md
```

## Test Coverage

The project includes comprehensive unit tests covering:

- **Parser Tests**: Each of the 7 parsers has dedicated tests with sample data
- **Integration Tests**: End-to-end testing through the main CLI
- **Error Handling**: Tests for file errors, parsing errors, and invalid inputs
- **Edge Cases**: Tests for various input formats and edge conditions

Coverage targets:
- Models: 100% (structure definitions for all 7 formats)
- Parsers: >95% (all parsing logic and error paths)
- Converter: >90% (main conversion flow)
- Main: >80% (CLI interface and integration)

Supported test formats with comprehensive parsing:
- JUnit XML (standard and custom variants)
- Postman/Newman JSON output
- Go test JSON streams
- TAP protocol (versions 12-14)
- xUnit.net XML results
- pytest JSON reports
- TestNG XML results

## Example Usage

Convert a JUnit XML file:
```bash
echo '<?xml version="1.0"?>
<testsuites>
  <testsuite name="MyTests" tests="1" failures="0" errors="0" time="0.5">
    <testcase name="testExample" classname="MyClass" time="0.5"/>
  </testsuite>
</testsuites>' > example_junit.xml

./tconv -input example_junit.xml -type junit -output result.json
```

Convert a TAP file:
```bash
echo 'TAP version 13
1..2
ok 1 - Basic test
not ok 2 - Failed test' > example.tap

./tconv -input example.tap -type tap -output result.json
```

The output will be a unified JSON format that can be consumed by various reporting tools regardless of the original test framework.

## Supported Frameworks

This tool supports output from these popular testing frameworks:
- **Java**: JUnit, TestNG
- **JavaScript/Node.js**: Jest (via TAP), Mocha (via TAP)  
- **Python**: pytest
- **Go**: go test
- **C#/.NET**: xUnit.net, NUnit (via JUnit XML)
- **API Testing**: Postman/Newman
- **Universal**: Any framework that outputs TAP format