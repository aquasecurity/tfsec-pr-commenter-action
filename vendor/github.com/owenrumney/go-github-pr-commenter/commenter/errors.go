package commenter

import "fmt"

// CommentAlreadyWrittenError returned when the error can't be written as it already exists
type CommentAlreadyWrittenError struct {
	filepath string
	comment  string
}

// CommentNotValidError returned when the comment is for a file or line not in the pr
type CommentNotValidError struct {
	filepath string
	lineNo   int
}

// PrDoesNotExistError returned when the PR can't be found, either as 401 or not existing
type PrDoesNotExistError struct {
	owner    string
	repo     string
	prNumber int
}

// AbuseRateLimitError return when the GitHub abuse rate limit is hit
type AbuseRateLimitError struct {
	owner            string
	repo             string
	prNumber         int
	BackoffInSeconds int
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

func newPrDoesNotExistError(owner, repo string, prNumber int) PrDoesNotExistError {
	return PrDoesNotExistError{
		owner:    owner,
		repo:     repo,
		prNumber: prNumber,
	}
}

func (e PrDoesNotExistError) Error() string {
	return fmt.Sprintf("PR number [%d] not found for %s/%s", e.prNumber, e.owner, e.repo)
}

func newAbuseRateLimitError(owner, repo string, prNumber, backoffInSeconds int) AbuseRateLimitError {
	return AbuseRateLimitError{
		owner:            owner,
		repo:             repo,
		prNumber:         prNumber,
		BackoffInSeconds: backoffInSeconds,
	}
}

func (e AbuseRateLimitError) Error() string {
	return fmt.Sprintf("Abuse limit reached on PR [%d] not found for %s/%s", e.prNumber, e.owner, e.repo)
}
