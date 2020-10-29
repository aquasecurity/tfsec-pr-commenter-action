package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"os"
)

type Range struct {
	Filename  string `json:"filename"`
	StartLine int    `json:"start_line"`
	EndLine   int    `json:"end_line"`
}

type Result struct {
	RuleID          string `json:"rule_id"`
	RuleDescription string `json:"rule_description"`
	RuleProvider    string `json:"rule_provider"`
	Link            string `json:"link"`
	Range           Range  `json:"location"`
	Description     string `json:"description"`
	RangeAnnotation string `json:"-"`
	Severity        string `json:"severity"`
}

type githubConnector struct {
	ctx      context.Context
	client   *github.Client
	owner    string
	repo     string
	prNumber int
}

func newConnector(owner, repo string, prNo int) *githubConnector {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &githubConnector{
		ctx:      ctx,
		client:   client,
		owner:    owner,
		repo:     repo,
		prNumber: prNo,
	}
}

func (gc *githubConnector) getPrFiles() ([]string, error) {
	files, _, err := gc.client.PullRequests.ListFiles(gc.ctx, gc.owner, gc.repo, gc.prNumber, nil)
	if err != nil {
		return nil, err
	}
	var filepaths []string
	for _, file := range files {
		if *file.Status != "deleted" {
			filepaths = append(filepaths, *file.Filename)
		}
	}
	return filepaths, nil
}

func main() {
	owner := ""
	repo := ""
	prNo := 1

	gc := newConnector(owner, repo, prNo)
	files, err := gc.getPrFiles()
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, file := range files {
		println(file)
	}
}
