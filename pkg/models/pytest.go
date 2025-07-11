package models

type PytestReport struct {
	Report    PytestReportInfo `json:"report"`
	Tests     []PytestTest     `json:"tests"`
	Summary   PytestSummary    `json:"summary"`
	Collectors []PytestCollector `json:"collectors,omitempty"`
}

type PytestReportInfo struct {
	PytestVersion string                 `json:"pytest_version"`
	Outcome       string                 `json:"outcome"`
	Result        map[string]interface{} `json:"result,omitempty"`
	Node          map[string]interface{} `json:"node,omitempty"`
	Created       float64                `json:"created"`
	Duration      float64                `json:"duration,omitempty"`
}

type PytestTest struct {
	NodeID    string                 `json:"nodeid"`
	Outcome   string                 `json:"outcome"`
	Setup     PytestPhase            `json:"setup,omitempty"`
	Call      PytestPhase            `json:"call,omitempty"`
	Teardown  PytestPhase            `json:"teardown,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Keywords  []string               `json:"keywords,omitempty"`
	Location  []interface{}          `json:"location,omitempty"`
}

type PytestPhase struct {
	Outcome   string  `json:"outcome"`
	Duration  float64 `json:"duration,omitempty"`
	Crash     *PytestCrash `json:"crash,omitempty"`
	Traceback *PytestTraceback `json:"traceback,omitempty"`
	Stdout    string  `json:"stdout,omitempty"`
	Stderr    string  `json:"stderr,omitempty"`
	Log       string  `json:"log,omitempty"`
}

type PytestCrash struct {
	Path      string `json:"path"`
	LineNo    int    `json:"lineno"`
	Message   string `json:"message"`
}

type PytestTraceback struct {
	Entries []PytestTracebackEntry `json:"entries,omitempty"`
}

type PytestTracebackEntry struct {
	Path     string `json:"path"`
	LineNo   int    `json:"lineno"`
	Message  string `json:"message"`
}

type PytestSummary struct {
	Total    int            `json:"total"`
	Passed   int            `json:"passed"`
	Failed   int            `json:"failed"`
	Skipped  int            `json:"skipped"`
	Error    int            `json:"error"`
	Duration float64        `json:"duration"`
	Outcome  string         `json:"outcome"`
}

type PytestCollector struct {
	NodeID   string                 `json:"nodeid"`
	Outcome  string                 `json:"outcome"`
	Result   []map[string]interface{} `json:"result,omitempty"`
}