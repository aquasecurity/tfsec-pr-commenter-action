package main

import (
	"encoding/json"
	"io/ioutil"
)

type checkRange struct {
	Filename  string `json:"filename"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type result struct {
	RuleID          string      `json:"long_id"`
	RuleDescription string      `json:"rule_description"`
	RuleProvider    string      `json:"rule_provider"`
	Links           []string    `json:"links"`
	Range           *checkRange `json:"location"`
	Description     string      `json:"description"`
	RangeAnnotation string      `json:"-"`
	Severity        string      `json:"severity"`
}

const resultsFile = "results.json"

func loadResultsFile() ([]result, error) {
	results := struct{ Results []result }{}

	file, err := ioutil.ReadFile(resultsFile)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &results)
	if err != nil {
		return nil, err
	}
	return results.Results, nil
}
