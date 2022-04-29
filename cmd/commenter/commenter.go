package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/owenrumney/go-github-pr-commenter/commenter"
)

func main() {
	fmt.Println("Starting the github commenter...")

	token := os.Getenv("INPUT_GITHUB_TOKEN")
	if len(token) == 0 {
		fail("the INPUT_GITHUB_TOKEN has not been set")
	}

	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	split := strings.Split(githubRepository, "/")
	if len(split) != 2 {
		fail(fmt.Sprintf("Expected value for split not found. Expected 2 in %v", split))
	}
	owner := split[0]
	repo := split[1]

	prNo, err := extractPullRequestNumber()
	if err != nil {
		fmt.Println("Not a PR, nothing to comment on, exiting...")
		return
	}

	c, err := commenter.NewCommenter(token, owner, repo, prNo)
	if err != nil {
		fail(fmt.Sprintf("failed to create a new commenter. %s", err.Error()))
	}
	results, err := loadResultsFile()
	if err != nil {
		fail(fmt.Sprintf("failed to load results. %s", err.Error()))
	}

	if len(results) == 0 {
		fmt.Println("No issues found.")
		os.Exit(0)
	}

	var errMessages []string
	workspacePath := fmt.Sprintf("%s/", os.Getenv("GITHUB_WORKSPACE"))
	var validCommentWritten bool
	for _, result := range results {
		result.Range.Filename = strings.ReplaceAll(result.Range.Filename, workspacePath, "")
		comment := generateErrorMessage(result)
		err := c.WriteMultiLineComment(result.Range.Filename, comment, result.Range.StartLine, result.Range.EndLine)
		if err != nil {
			// don't error if its simply that the comments aren't valid for the PR
			switch err.(type) {
			case commenter.CommentAlreadyWrittenError:
				fmt.Println("Comment already written so not writing")
				validCommentWritten = true
			case commenter.CommentNotValidError:
				fmt.Printf("%s .... not writing as not part of the current PR\n", result.Description)
				continue
			default:
				errMessages = append(errMessages, err.Error())
			}
		} else {
			validCommentWritten = true
			fmt.Printf("Commenting for %s to %s:%d:%d\n", result.Description, result.Range.Filename, result.Range.StartLine, result.Range.EndLine)
		}
	}

	if len(errMessages) > 0 {
		fmt.Printf("There were %d errors:\n", len(errMessages))
		for _, err := range errMessages {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	if validCommentWritten || len(errMessages) > 0 {
		if softFail, ok := os.LookupEnv("INPUT_SOFT_FAIL_COMMENTER"); ok && strings.ToLower(softFail) == "true" {
			return
		}
		os.Exit(1)
	}
}

func generateErrorMessage(result result) string {
	return fmt.Sprintf(`:warning: tfsec found a **%s** severity issue from rule `+"`%s`"+`:
> %s

More information available %s`,
		result.Severity, result.RuleID, result.Description, formatUrls(result.Links))
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

	prNumber, err := strconv.Atoi(fmt.Sprintf("%v", payload["number"]))
	if err != nil {
		return 0, fmt.Errorf("not a valid PR")
	}
	return prNumber, nil
}

func formatUrls(urls []string) string {
	urlList := ""
	for _, url := range urls {
		if urlList != "" {
			urlList += fmt.Sprintf(" and ")
		}
		urlList += fmt.Sprintf("[here](%s)", url)
	}
	return urlList
}

func fail(err string) {
	fmt.Printf("The commenter failed with the following error:\n%s\n", err)
	os.Exit(-1)
}
