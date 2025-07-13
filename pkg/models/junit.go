package models

import "encoding/xml"

// JUnitTestSuites represents the root element of JUnit XML containing multiple test suites.
// This structure is used when parsing JUnit XML files that have a <testsuites> root element.
type JUnitTestSuites struct {
	XMLName    xml.Name         `xml:"testsuites"`
	Name       string           `xml:"name,attr,omitempty"`
	Tests      int              `xml:"tests,attr,omitempty"`
	Failures   int              `xml:"failures,attr,omitempty"`
	Errors     int              `xml:"errors,attr,omitempty"`
	Time       string           `xml:"time,attr,omitempty"`
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

// JUnitTestSuite represents a single test suite in JUnit XML format.
// This structure maps to the <testsuite> element and contains test cases and metadata.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Name       string          `xml:"name,attr"`
	Package    string          `xml:"package,attr,omitempty"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Skipped    int             `xml:"skipped,attr,omitempty"`
	Time       string          `xml:"time,attr"`
	Timestamp  string          `xml:"timestamp,attr,omitempty"`
	Properties *JUnitProperties `xml:"properties,omitempty"`
	TestCases  []JUnitTestCase `xml:"testcase"`
	SystemOut  string          `xml:"system-out,omitempty"`
	SystemErr  string          `xml:"system-err,omitempty"`
}

// JUnitProperties represents the <properties> element containing key-value property pairs.
type JUnitProperties struct {
	Properties []JUnitProperty `xml:"property"`
}

// JUnitProperty represents a single property with name and value attributes.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitTestCase represents a single test case in JUnit XML format.
// This structure maps to the <testcase> element and contains test execution results.
type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Name      string        `xml:"name,attr"`
	ClassName string        `xml:"classname,attr,omitempty"`
	Time      string        `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
	Error     *JUnitError   `xml:"error,omitempty"`
	Skipped   *JUnitSkipped `xml:"skipped,omitempty"`
	SystemOut string        `xml:"system-out,omitempty"`
	SystemErr string        `xml:"system-err,omitempty"`
}

// JUnitFailure represents a test failure with message, type, and content details.
type JUnitFailure struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

// JUnitError represents a test error with message, type, and content details.
type JUnitError struct {
	Message string `xml:"message,attr,omitempty"`
	Type    string `xml:"type,attr,omitempty"`
	Content string `xml:",chardata"`
}

// JUnitSkipped represents a skipped test with an optional message.
type JUnitSkipped struct {
	Message string `xml:"message,attr,omitempty"`
}