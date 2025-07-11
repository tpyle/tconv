package models

type TAPResult struct {
	Version     string    `json:"version,omitempty"`
	Plan        *TAPPlan  `json:"plan,omitempty"`
	TestLines   []TAPTest `json:"test_lines"`
	Diagnostics []string  `json:"diagnostics,omitempty"`
	BailOut     *string   `json:"bail_out,omitempty"`
}

type TAPPlan struct {
	Start int    `json:"start"`
	End   int    `json:"end"`
	Skip  string `json:"skip,omitempty"`
}

type TAPTest struct {
	Number      int               `json:"number"`
	OK          bool              `json:"ok"`
	Description string            `json:"description,omitempty"`
	Directive   *TAPDirective     `json:"directive,omitempty"`
	Diagnostics []string          `json:"diagnostics,omitempty"`
	YAML        map[string]interface{} `json:"yaml,omitempty"`
	Raw         string            `json:"raw"`
}

type TAPDirective struct {
	Type   string `json:"type"`   // "SKIP" or "TODO"
	Reason string `json:"reason,omitempty"`
}