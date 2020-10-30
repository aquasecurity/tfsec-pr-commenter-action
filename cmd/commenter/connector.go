package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type githubConnector struct {
	ctx      context.Context
	client   *github.Client
	owner    string
	repo     string
	prNumber int
	commitId string
}

func newConnector() (*githubConnector, error) {

	token := os.Getenv("INPUT_GITHUB_TOKEN")
	if len(token) == 0 {
		return nil, errors.New("the INPUT_GITHUB_TOKEN has not been set")
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("INPUT_GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	split := strings.Split(githubRepository, "/")
	if len(split) != 2 {
		panic(fmt.Sprintf("Expected value for split not found. Expected 2 in %v", split))
	}
	owner := split[0]
	repo := split[1]
	commitId := os.Getenv("GITHUB_SHA")
	prNo, err := extractPullRequestNumber()
	if err != nil {
		panic("unable to get the PR number, can't continue")
	}

	return &githubConnector{
		ctx:      ctx,
		client:   client,
		owner:    owner,
		repo:     repo,
		prNumber: prNo,
		commitId: commitId,
	}, nil
}

func (gc *githubConnector) getPrFiles() ([]*github.CommitFile, error) {
	files, _, err := gc.client.PullRequests.ListFiles(gc.ctx, gc.owner, gc.repo, gc.prNumber, nil)
	if err != nil {
		return nil, err
	}
	var commitFiles []*github.CommitFile
	for _, file := range files {
		if *file.Status != "deleted" {
			commitFiles = append(commitFiles, file)
		}

	}
	return commitFiles, nil
}

func (gc *githubConnector) getExistingComments() ([]string, error) {
	var bodies []string
	comments, _, err := gc.client.PullRequests.ListComments(gc.ctx, gc.owner, gc.repo, gc.prNumber, &github.PullRequestListCommentsOptions{})
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		bodies = append(bodies, comment.GetBody())
	}
	return bodies, nil
}

func (gc *githubConnector) commentOnPrResult(result *commentBlock, existingComments []string) {
	errorMessage := fmt.Sprintf(`tfsec check %s failed. 

%s

For more information, see https://tfsec.dev/docs/%s/%s/`, result.code, result.description, strings.ToLower(result.provider), result.code)

	for _, existingComment := range existingComments {
		if errorMessage == existingComment {
			// don't create the comment, its already there for this block
			return
		}
	}

	comment := buildComment(result, errorMessage)
	fmt.Printf("%+v\n", comment)
	_, _, err := gc.client.PullRequests.CreateComment(gc.ctx, gc.owner, gc.repo, gc.prNumber, comment)
	if err != nil {
		fmt.Println("Error occurred %s", err.Error())
	}
}

func buildComment(result *commentBlock, errorMessage string) *github.PullRequestComment {
	comment := &github.PullRequestComment{
		Line:     &result.startLine,
		Path:     &result.fileName,
		CommitID: &result.sha,
		Body:     &errorMessage,
		Position: &result.position,
	}

	if result.startLine != result.endLine {
		comment.StartLine = &result.startLine
		comment.Line = &result.endLine
	}
	return comment
}

func extractPullRequestNumber() (int, error) {
	file, err := ioutil.ReadFile("/github/workflow/event.json")
	if err != nil {
		return -1, err
	}

	var data interface{}
	err = json.Unmarshal(file, &data)
	if err != nil {
		return -1, err
	}
	payload := data.(map[string]interface{})

	return strconv.Atoi(fmt.Sprintf("%v", payload["number"]))
}
