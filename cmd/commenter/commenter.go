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
		fail(err.Error())
	}

	c, err := commenter.NewCommenter(token, owner, repo, prNo)
	if err != nil {
		fail(err.Error())
	}
	results, err := loadResultsFile()
	if err != nil {
		fail(err.Error())
	}

	var errMessages []string
	workspacePath := fmt.Sprintf("%s/", os.Getenv("GITHUB_WORKSPACE"))
	for _, result := range results {
		result.Range.Filename = strings.ReplaceAll(result.Range.Filename, workspacePath, "")
		comment := generateErrorMessage(result)
		err := c.WriteMultiLineComment(result.Range.Filename, comment, result.Range.StartLine, result.Range.EndLine)
		if err != nil {
			// don't error if its simply that the comments aren't valid for the PR
			switch err.(type) {
			case commenter.CommentAlreadyWrittenError:
				fmt.Println("Comment already written so not writing")
			case commenter.CommentNotValidError:
				fmt.Println("Comment not written, not part of the current PR")
				continue
			default:
				errMessages = append(errMessages, err.Error())
			}
		} else {
			fmt.Printf("Writing comment to %s:%d:%d", result.Range.Filename, result.Range.StartLine, result.Range.EndLine)
		}
	}

	if len(errMessages) > 0 {
		fmt.Printf("There were %d errors:\n", len(errMessages))
		for _, err := range errMessages {
			fmt.Println(err)
		}
	}
}

func generateErrorMessage(result result) string {
	return fmt.Sprintf(`tfsec check %s failed. 

%s

For more information, see https://tfsec.dev/docs/%s/%s/`,
		result.RuleID, result.Description, strings.ToLower(result.RuleProvider), result.RuleID)
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

func fail(err string) {
	fmt.Printf("The commenter failed with the following error:\n%s\n", err)
	os.Exit(-1)
}
