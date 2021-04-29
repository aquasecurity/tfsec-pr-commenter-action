package commenter

import (
	"context"

	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
)

type connector struct {
	prs      *github.PullRequestsService
	comments *github.IssuesService
	owner    string
	repo     string
	prNumber int
}

type existingComment struct {
	filename  *string
	comment   *string
	commentId *int64
}

func createConnector(token, owner, repo string, prNumber int) *connector {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	return &connector{
		prs:      client.PullRequests,
		comments: client.Issues,
		owner:    owner,
		repo:     repo,
		prNumber: prNumber,
	}
}

func (c *connector) writeReviewComment(block *github.PullRequestComment, commentId *int64) error {
	ctx := context.Background()

	if commentId != nil {
		var _, err = c.prs.DeleteComment(ctx, c.owner, c.repo, *commentId)
		if err != nil {
			return err
		}
	}
	var _, _, err = c.prs.CreateComment(ctx, c.owner, c.repo, c.prNumber, block)
	if err != nil {
		return err
	}
	return nil
}

func (c *connector) writeGeneralComment(comment *github.IssueComment) error {
	ctx := context.Background()

	var _, _, err = c.comments.CreateComment(ctx, c.owner, c.repo, c.prNumber, comment)
	if err != nil {
		return err
	}

	return nil
}

func (c *connector) getFilesForPr() ([]*github.CommitFile, error) {
	files, _, err := c.prs.ListFiles(context.Background(), c.owner, c.repo, c.prNumber, nil)
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

func (c *connector) getExistingComments() ([]*existingComment, error) {
	ctx := context.Background()

	comments, _, err := c.prs.ListComments(ctx, c.owner, c.repo, c.prNumber, &github.PullRequestListCommentsOptions{})
	if err != nil {
		return nil, err
	}

	var existingComments []*existingComment
	for _, comment := range comments {
		existingComments = append(existingComments, &existingComment{
			filename:  comment.Path,
			comment:   comment.Body,
			commentId: comment.ID,
		})
	}
	return existingComments, nil
}

func (c *connector) prExists() bool {
	ctx := context.Background()

	_, _, err := c.prs.Get(ctx, c.owner, c.repo, c.prNumber)
	return err == nil
}
