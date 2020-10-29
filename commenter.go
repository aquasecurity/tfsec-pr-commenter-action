package main

import (
	"context"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"os"
)

func main() {

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	issueComment := &github.IssueComment{}

}
