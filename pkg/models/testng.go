package models

import "encoding/xml"

type TestNGResults struct {
	XMLName xml.Name     `xml:"testng-results"`
	Total   int          `xml:"total,attr"`
	Passed  int          `xml:"passed,attr"`
	Failed  int          `xml:"failed,attr"`
	Skipped int          `xml:"skipped,attr"`
	Suites  []TestNGSuite `xml:"suite"`
}

type TestNGSuite struct {
	XMLName       xml.Name      `xml:"suite"`
	Name          string        `xml:"name,attr"`
	DurationMs    string        `xml:"duration-ms,attr,omitempty"`
	StartedAt     string        `xml:"started-at,attr,omitempty"`
	FinishedAt    string        `xml:"finished-at,attr,omitempty"`
	Groups        *TestNGGroups `xml:"groups,omitempty"`
	Tests         []TestNGTest  `xml:"test"`
}

type TestNGGroups struct {
	XMLName xml.Name     `xml:"groups"`
	Groups  []TestNGGroup `xml:"group"`
}

type TestNGGroup struct {
	XMLName xml.Name `xml:"group"`
	Name    string   `xml:"name,attr"`
	Methods []TestNGGroupMethod `xml:"method"`
}

type TestNGGroupMethod struct {
	XMLName   xml.Name `xml:"method"`
	Signature string   `xml:"signature,attr"`
	Name      string   `xml:"name,attr"`
	Class     string   `xml:"class,attr"`
}

type TestNGTest struct {
	XMLName       xml.Name       `xml:"test"`
	Name          string         `xml:"name,attr"`
	DurationMs    string         `xml:"duration-ms,attr,omitempty"`
	StartedAt     string         `xml:"started-at,attr,omitempty"`
	FinishedAt    string         `xml:"finished-at,attr,omitempty"`
	Classes       []TestNGClass  `xml:"class"`
}

type TestNGClass struct {
	XMLName xml.Name        `xml:"class"`
	Name    string          `xml:"name,attr"`
	Methods []TestNGMethod  `xml:"test-method"`
}

type TestNGMethod struct {
	XMLName       xml.Name           `xml:"test-method"`
	Signature     string             `xml:"signature,attr"`
	Name          string             `xml:"name,attr"`
	ClassName     string             `xml:"class-name,attr"`
	Groups        string             `xml:"groups,attr,omitempty"`
	Description   string             `xml:"description,attr,omitempty"`
	DurationMs    string             `xml:"duration-ms,attr,omitempty"`
	StartedAt     string             `xml:"started-at,attr,omitempty"`
	FinishedAt    string             `xml:"finished-at,attr,omitempty"`
	Status        string             `xml:"status,attr"`
	IsConfig      string             `xml:"is-config,attr,omitempty"`
	Exception     *TestNGException   `xml:"exception,omitempty"`
	Reporter      []TestNGReporter   `xml:"reporter-output,omitempty"`
	Params        []TestNGParam      `xml:"params>param,omitempty"`
}

type TestNGException struct {
	XMLName    xml.Name `xml:"exception"`
	Class      string   `xml:"class,attr"`
	Message    string   `xml:"message"`
	FullStacktrace string `xml:"full-stacktrace"`
}

type TestNGReporter struct {
	XMLName xml.Name `xml:"reporter-output"`
	Output  string   `xml:",chardata"`
}

type TestNGParam struct {
	XMLName xml.Name `xml:"param"`
	Index   string   `xml:"index,attr"`
	Value   string   `xml:"value"`
}