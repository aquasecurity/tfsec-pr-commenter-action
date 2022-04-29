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
	fmt.Println("Starting the github commenter")

	token := os.Getenv("INPUT_GITHUB_TOKEN")
	if len(token) == 0 {
		fail("the INPUT_GITHUB_TOKEN has not been set")
	}

	githubRepository := os.Getenv("GITHUB_REPOSITORY")
	split := strings.Split(githubRepository, "/")
	if len(split) != 2 {
		fail(fmt.Sprintf("unexpected value for GITHUB_REPOSITORY. Expected <organisation/name>, found %v", split))
	}
	owner := split[0]
	repo := split[1]

	fmt.Printf("Working in repository %s\n", repo)

	prNo, err := extractPullRequestNumber()
	if err != nil {
		fmt.Println("Not a PR, nothing to comment on, exiting")
		return
	}
	fmt.Printf("Working in PR %v\n", prNo)

	results, err := loadResultsFile()
	if err != nil {
		fail(fmt.Sprintf("failed to load results. %s", err.Error()))
	}

	if len(results) == 0 {
		fmt.Println("No issues found.")
		os.Exit(0)
	}
	fmt.Printf("TFSec found %v issues\n", len(results))

	c, err := commenter.NewCommenter(token, owner, repo, prNo)
	if err != nil {
		fail(fmt.Sprintf("could not connect to GitHub (%s)", err.Error()))
	}

	workspacePath := fmt.Sprintf("%s/", os.Getenv("GITHUB_WORKSPACE"))
	fmt.Printf("Working in GITHUB_WORKSPACE %s\n", workspacePath)

	var errMessages []string
	var validCommentWritten bool
	for _, result := range results {
		result.Range.Filename = strings.ReplaceAll(result.Range.Filename, workspacePath, "")
		comment := generateErrorMessage(result)
		fmt.Printf("Preparing comment for violation of rule %v in %v\n", result.RuleID, result.Range.Filename)
		err := c.WriteMultiLineComment(result.Range.Filename, comment, result.Range.StartLine, result.Range.EndLine)
		if err != nil {
			// don't error if its simply that the comments aren't valid for the PR
			switch err.(type) {
			case commenter.CommentAlreadyWrittenError:
				fmt.Println("Ignoring - comment already written")
				validCommentWritten = true
			case commenter.CommentNotValidError:
				fmt.Println("Ignoring - change not part of the current PR")
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
	return fmt.Sprintf(`tfsec check %s failed. 

Description: %s

Severity: %s

For more information, see:

%s`,
		result.RuleID, result.Description, result.Severity, formatUrls(result.Links))
}

func extractPullRequestNumber() (int, error) {
	github_event_file := "/github/workflow/event.json"
	file, err := ioutil.ReadFile(github_event_file)
	if err != nil {
		fail(fmt.Sprintf("GitHub event payload not found in %s", github_event_file))
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
		urlList += fmt.Sprintf("- %s\n", url)
	}
	return urlList
}

func fail(err string) {
	fmt.Printf("Error: %s\n", err)
	os.Exit(-1)
}
