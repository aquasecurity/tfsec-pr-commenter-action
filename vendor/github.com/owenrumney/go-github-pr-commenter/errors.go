package go_github_pr_commenter

import "fmt"

type NotPartOfPrError struct {
	filepath string
}

type CommentAlreadyWrittenError struct {
	filepath string
	comment  string
}

type CommentNotValidError struct {
	filepath string
	lineNo   int
}

type PrDoesNotExistError struct {
	owner    string
	repo     string
	prNumber int
}

func newNotPartOfPrError(filepath string) NotPartOfPrError {
	return NotPartOfPrError{
		filepath: filepath,
	}
}

func (e NotPartOfPrError) Error() string {
	return fmt.Sprintf("The file [%s] provided is not part of the PR", e.filepath)
}

func newCommentAlreadyWrittenError(filepath, comment string) CommentAlreadyWrittenError {
	return CommentAlreadyWrittenError{
		filepath: filepath,
		comment:  comment,
	}
}

func (e CommentAlreadyWrittenError) Error() string {
	return fmt.Sprintf("The file [%s] already has the comment written [%s]", e.filepath, e.comment)
}

func newCommentNotValidError(filepath string, line int) CommentNotValidError {
	return CommentNotValidError{
		filepath: filepath,
		lineNo:   line,
	}
}

func (e CommentNotValidError) Error() string {
	return fmt.Sprintf("There is nothing to comment on at line [%d] in file [%s]", e.lineNo, e.filepath)
}

func newPrDoesNotExistError(c *connector) PrDoesNotExistError {
	return PrDoesNotExistError{
		owner:    c.owner,
		repo:     c.repo,
		prNumber: c.prNumber,
	}
}

func (e PrDoesNotExistError) Error() string {
	return fmt.Sprintf("PR number [%d] not found for %s/%s", e.prNumber, e.owner, e.repo)
}
