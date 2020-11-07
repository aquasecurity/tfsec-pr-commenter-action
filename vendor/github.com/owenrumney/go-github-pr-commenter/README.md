# go-github-pr-commenter

## What is it?

A convenience libary that wraps the [go-github](https://github.com/google/go-github) library and allows you to quickly add comments to the lines of changes in a comment.

The intention is this is used with CI tools to automatically comment on new Github pull requests when static analysis checks are failing.

For an example of this in use, see the [tfsec-pr-commenter-action](https://github.com/tfsec/tfsec-pr-commenter-action). This Github action will run against your Terraform and report any security issues that are present in the code before it is committed.

## How do I use it?

The intention is to keep the interface as clean as possible; steps are 

- create a commenter for a repo and PR
- write comments to the commenter
    - comments which exist will not be written
    - comments that aren't appropriate (not part of the PR) will not be written
    
### Expected Errors

The following errors can be handled - I hope these are self explanatory

```go
type PrDoesNotExistError struct {
	owner    string
	repo     string
	prNumber int
}

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
```

### Basic Usage Example

```go
package main

import (
    commenter "github.com/owenrumney/go-github-pr-commenter"
    log "github.com/sirupsen/logrus"
    "os"
)


// Create the commenter
token := os.Getenv("GITHUB_TOKEN")

c, err := commenter.NewCommenter(token, "tfsec", "tfsec-example-project", 8)
if err != nil {
    fmt.Println(err.Error())
}

// process whatever static analysis results you've gathered
for _, result := myResults {
    err = c.WriteMultiLineComment(result.Path, result.Comment, result.StartLine, result.EndLine)
    if err != nil {
        if errors.Is(err, commenter.CommentNotValidError{}) {
            log.Debugf("result not relevant for commit. %s", err.Error())
        } else {
            log.Errorf("an error occurred writing the comment: %s", err.Error())
        }
    }
}
```
