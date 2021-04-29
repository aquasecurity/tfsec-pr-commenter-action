package commenter

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/go-github/v32/github"
)

// Commenter is the main commenter struct
type Commenter struct {
	pr               *connector
	existingComments []*existingComment
	files            []*commitFileInfo
	loaded           bool
}

var patchRegex *regexp.Regexp
var commitRefRegex *regexp.Regexp

// NewCommenter creates a Commenter for updating PR with comments
func NewCommenter(token, owner, repo string, prNumber int) (*Commenter, error) {
	regex, err := regexp.Compile(`^@@.*\+(\d+),(\d+).+?@@`)
	if err != nil {
		return nil, err
	}
	patchRegex = regex

	regex, err = regexp.Compile(".+ref=(.+)")
	if err != nil {
		return nil, err
	}
	commitRefRegex = regex

	if len(token) == 0 {
		return nil, errors.New("the GITHUB_TOKEN has not been set")
	}

	connector := createConnector(token, owner, repo, prNumber)

	if !connector.prExists() {
		return nil, newPrDoesNotExistError(connector)
	}

	c := &Commenter{
		pr: connector,
	}
	return c, nil
}

// WriteMultiLineComment writes a multiline review on a file in the github PR
func (c *Commenter) WriteMultiLineComment(file, comment string, startLine, endLine int) error {
	if !c.loaded {
		err := c.loadPr()
		if err != nil {
			return err
		}
	}

	if !c.checkCommentRelevant(file, startLine) || !c.checkCommentRelevant(file, endLine) {
		return newCommentNotValidError(file, startLine)
	}

	if startLine == endLine {
		return c.WriteLineComment(file, comment, endLine)
	}

	info, err := c.getFileInfo(file, endLine)
	if err != nil {
		return err
	}

	prComment := buildComment(file, comment, endLine, *info)
	prComment.StartLine = &startLine
	return c.writeCommentIfRequired(prComment)
}

// WriteLineComment writes a single review line on a file of the github PR
func (c *Commenter) WriteLineComment(file, comment string, line int) error {
	if !c.loaded {
		err := c.loadPr()
		if err != nil {
			return err
		}
	}

	if !c.checkCommentRelevant(file, line) {
		return newCommentNotValidError(file, line)
	}

	info, err := c.getFileInfo(file, line)
	if err != nil {
		return err
	}
	prComment := buildComment(file, comment, line, *info)
	return c.writeCommentIfRequired(prComment)
}

func (c *Commenter) WriteGeneralComment(comment string) error {
	if !c.loaded {
		err := c.loadPr()
		if err != nil {
			return err
		}
	}
	issueComment := &github.IssueComment{
		Body: &comment,
	}
	return c.pr.writeGeneralComment(issueComment)
}

func (c *Commenter) writeCommentIfRequired(prComment *github.PullRequestComment) error {
	var commentId *int64
	for _, existing := range c.existingComments {
		commentId = func(ec *existingComment) *int64 {
			if *ec.filename == *prComment.Path && *ec.comment == *prComment.Body {
				return ec.commentId
			}
			return nil
		}(existing)
	}
	return c.pr.writeReviewComment(prComment, commentId)
}

func (c *Commenter) getCommitFileInfo() error {
	prFiles, err := c.pr.getFilesForPr()
	if err != nil {
		return err
	}
	var errs []string
	for _, file := range prFiles {
		info, err := getCommitInfo(file)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		c.files = append(c.files, info)
	}
	if len(errs) > 0 {
		return fmt.Errorf("there were errors processing the PR files.\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func (c *Commenter) loadPr() error {
	err := c.getCommitFileInfo()
	if err != nil {
		return err
	}

	c.existingComments, err = c.pr.getExistingComments()
	if err != nil {
		return err
	}
	c.loaded = true
	return nil
}

func (c *Commenter) checkCommentRelevant(filename string, line int) bool {
	for _, file := range c.files {
		if relevant := func(file *commitFileInfo) bool {
			if file.FileName == filename {
				if line >= file.hunkStart && line <= file.hunkEnd {
					return true
				}
			}
			return false
		}(file); relevant {
			return true
		}
	}
	return false
}

func (c *Commenter) getFileInfo(file string, line int) (*commitFileInfo, error) {
	for _, info := range c.files {
		if info.FileName == file {
			if line >= info.hunkStart && line <= info.hunkEnd {
				return info, nil
			}
		}
	}
	return nil, errors.New("file not found, shouldn't have got to here")
}

func buildComment(file, comment string, line int, info commitFileInfo) *github.PullRequestComment {
	return &github.PullRequestComment{
		Line:     &line,
		Path:     &file,
		CommitID: &info.sha,
		Body:     &comment,
		Position: info.calculatePosition(line),
	}
}

func getCommitInfo(file *github.CommitFile) (*commitFileInfo, error) {
	groups := patchRegex.FindAllStringSubmatch(file.GetPatch(), -1)
	if len(groups) < 1 {
		return nil, errors.New("the patch details could not be resolved")
	}
	hunkStart, _ := strconv.Atoi(groups[0][1])
	hunkEnd, _ := strconv.Atoi(groups[0][2])

	shaGroups := commitRefRegex.FindAllStringSubmatch(file.GetContentsURL(), -1)
	if len(shaGroups) < 1 {
		return nil, errors.New("the sha details could not be resolved")
	}
	sha := shaGroups[0][1]

	return &commitFileInfo{
		FileName:  *file.Filename,
		hunkStart: hunkStart,
		hunkEnd:   hunkStart + (hunkEnd - 1),
		sha:       sha,
	}, nil
}
