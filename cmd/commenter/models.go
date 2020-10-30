package main

type commitFileInfo struct {
	fileName  string
	hunkStart int
	hunkEnd   int
	sha       string
}

type commentBlock struct {
	fileName    string
	startLine   int
	endLine     int
	position    int
	sha         string
	code        string
	description string
	provider    string
}

type checkRange struct {
	Filename  string `json:"filename"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type result struct {
	RuleID          string      `json:"rule_id"`
	RuleDescription string      `json:"rule_description"`
	RuleProvider    string      `json:"rule_provider"`
	Link            string      `json:"link"`
	Range           *checkRange `json:"location"`
	Description     string      `json:"description"`
	RangeAnnotation string      `json:"-"`
	Severity        string      `json:"severity"`
}
