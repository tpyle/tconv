package models

import "time"

type GoTestEvent struct {
	Time    time.Time `json:"Time"`
	Action  string    `json:"Action"`
	Package string    `json:"Package,omitempty"`
	Test    string    `json:"Test,omitempty"`
	Elapsed float64   `json:"Elapsed,omitempty"`
	Output  string    `json:"Output,omitempty"`
}

const (
	ActionRun    = "run"
	ActionPause  = "pause"
	ActionCont   = "cont"
	ActionPass   = "pass"
	ActionBench  = "bench"
	ActionFail   = "fail"
	ActionOutput = "output"
	ActionSkip   = "skip"
)