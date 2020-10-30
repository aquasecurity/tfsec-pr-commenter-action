package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Starting the github commenter...")

	gc, err := newConnector()
	if err != nil {
		fail(err)
	}

	files, err := gc.getPrFiles()
	if err != nil {
		fail(err)
	}

	fmt.Println("Getting relevant results for the PR...")
	results, err := getRelevantResults(files)
	if err != nil {
		fail(err)
	}

	existingComments, err := gc.getExistingComments()
	if err != nil {
		fail(err)
	}

	for _, result := range results {
		fmt.Printf("Processing %s\n", result.fileName)
		gc.commentOnPrResult(result, existingComments)
	}
}

func fail(err error) {
	fmt.Printf("The commenter failed with the following error: %s\n", err)
	os.Exit(-1)
}
