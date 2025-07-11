package models

import "encoding/xml"

type XUnitAssemblies struct {
	XMLName    xml.Name       `xml:"assemblies"`
	Timestamp  string         `xml:"timestamp,attr,omitempty"`
	Assemblies []XUnitAssembly `xml:"assembly"`
}

type XUnitAssembly struct {
	XMLName        xml.Name          `xml:"assembly"`
	Name           string            `xml:"name,attr"`
	Environment    string            `xml:"environment,attr,omitempty"`
	TestFramework  string            `xml:"test-framework,attr,omitempty"`
	RunDate        string            `xml:"run-date,attr,omitempty"`
	RunTime        string            `xml:"run-time,attr,omitempty"`
	ConfigFile     string            `xml:"config-file,attr,omitempty"`
	Time           string            `xml:"time,attr,omitempty"`
	Total          int               `xml:"total,attr"`
	Passed         int               `xml:"passed,attr"`
	Failed         int               `xml:"failed,attr"`
	Skipped        int               `xml:"skipped,attr"`
	Errors         int               `xml:"errors,attr,omitempty"`
	Collections    []XUnitCollection `xml:"collection"`
}

type XUnitCollection struct {
	XMLName xml.Name    `xml:"collection"`
	Name    string      `xml:"name,attr"`
	Time    string      `xml:"time,attr,omitempty"`
	Total   int         `xml:"total,attr"`
	Passed  int         `xml:"passed,attr"`
	Failed  int         `xml:"failed,attr"`
	Skipped int         `xml:"skipped,attr"`
	Tests   []XUnitTest `xml:"test"`
}

type XUnitTest struct {
	XMLName xml.Name      `xml:"test"`
	Name    string        `xml:"name,attr"`
	Type    string        `xml:"type,attr,omitempty"`
	Method  string        `xml:"method,attr,omitempty"`
	Time    string        `xml:"time,attr,omitempty"`
	Result  string        `xml:"result,attr"`
	Output  string        `xml:"output,omitempty"`
	Reason  *XUnitReason  `xml:"reason,omitempty"`
	Failure *XUnitFailure `xml:"failure,omitempty"`
	Traits  *XUnitTraits  `xml:"traits,omitempty"`
}

type XUnitReason struct {
	XMLName xml.Name `xml:"reason"`
	Content string   `xml:",chardata"`
}

type XUnitFailure struct {
	XMLName      xml.Name `xml:"failure"`
	ExceptionType string  `xml:"exception-type,attr,omitempty"`
	Message      string   `xml:"message"`
	StackTrace   string   `xml:"stack-trace"`
}

type XUnitTraits struct {
	XMLName xml.Name     `xml:"traits"`
	Traits  []XUnitTrait `xml:"trait"`
}

type XUnitTrait struct {
	XMLName xml.Name `xml:"trait"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}