package models

type PostmanRun struct {
	Info      PostmanInfo        `json:"info"`
	Event     []PostmanEvent     `json:"event,omitempty"`
	Values    []PostmanValue     `json:"values,omitempty"`
	Run       PostmanRunDetails  `json:"run"`
}

type PostmanInfo struct {
	PostmanID string `json:"postman_id"`
	Name      string `json:"name"`
	Schema    string `json:"schema"`
}

type PostmanEvent struct {
	Listen string            `json:"listen"`
	Script PostmanScript    `json:"script"`
}

type PostmanScript struct {
	Type string   `json:"type"`
	Exec []string `json:"exec"`
}

type PostmanValue struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
	Type  string      `json:"type"`
}

type PostmanRunDetails struct {
	Stats      PostmanStats      `json:"stats"`
	Timings    PostmanTimings    `json:"timings"`
	Executions []PostmanExecution `json:"executions"`
}

type PostmanStats struct {
	Iterations PostmanIterationStats `json:"iterations"`
	Items      PostmanItemStats      `json:"items"`
	Scripts    PostmanScriptStats    `json:"scripts"`
	Prerequests PostmanScriptStats   `json:"prerequests"`
	Requests   PostmanRequestStats   `json:"requests"`
	Tests      PostmanTestStats      `json:"tests"`
	Assertions PostmanAssertionStats `json:"assertions"`
	TestScripts PostmanScriptStats   `json:"testScripts"`
}

type PostmanIterationStats struct {
	Total   int `json:"total"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
}

type PostmanItemStats struct {
	Total   int `json:"total"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
}

type PostmanScriptStats struct {
	Total   int `json:"total"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
}

type PostmanRequestStats struct {
	Total   int `json:"total"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
}

type PostmanTestStats struct {
	Total   int `json:"total"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
}

type PostmanAssertionStats struct {
	Total   int `json:"total"`
	Pending int `json:"pending"`
	Failed  int `json:"failed"`
}

type PostmanTimings struct {
	ResponseAverage float64 `json:"responseAverage"`
	ResponseMin     float64 `json:"responseMin"`
	ResponseMax     float64 `json:"responseMax"`
	ResponseSd      float64 `json:"responseSd"`
	DnsAverage      float64 `json:"dnsAverage"`
	DnsMin          float64 `json:"dnsMin"`
	DnsMax          float64 `json:"dnsMax"`
	DnsSd           float64 `json:"dnsSd"`
	FirstByteAverage float64 `json:"firstByteAverage"`
	FirstByteMin    float64 `json:"firstByteMin"`
	FirstByteMax    float64 `json:"firstByteMax"`
	FirstByteSd     float64 `json:"firstByteSd"`
	Started         int64   `json:"started"`
	Completed       int64   `json:"completed"`
}

type PostmanExecution struct {
	Item      PostmanItem      `json:"item"`
	Request   PostmanRequest   `json:"request"`
	Response  PostmanResponse  `json:"response"`
	Assertions []PostmanAssertion `json:"assertions"`
	RequestError interface{}    `json:"requestError,omitempty"`
}

type PostmanItem struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type PostmanRequest struct {
	Method string            `json:"method"`
	URL    PostmanURL        `json:"url"`
	Header []PostmanHeader   `json:"header"`
	Body   PostmanBody       `json:"body,omitempty"`
}

type PostmanURL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Port     string   `json:"port,omitempty"`
	Path     []string `json:"path"`
	Query    []PostmanQuery `json:"query,omitempty"`
}

type PostmanHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PostmanBody struct {
	Mode string `json:"mode"`
	Raw  string `json:"raw,omitempty"`
}

type PostmanQuery struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PostmanResponse struct {
	ID            string          `json:"id"`
	OriginalRequest PostmanRequest `json:"originalRequest"`
	ResponseTime  float64         `json:"responseTime"`
	Timings       map[string]float64 `json:"timings"`
	Header        []PostmanHeader `json:"header"`
	Cookie        []interface{}   `json:"cookie"`
	Body          string          `json:"body"`
	Status        string          `json:"status"`
	Code          int             `json:"code"`
}

type PostmanAssertion struct {
	Assertion string      `json:"assertion"`
	Error     interface{} `json:"error,omitempty"`
}