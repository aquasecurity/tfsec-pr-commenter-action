package commenter

import (
	"fmt"
	"strings"
)

type commitFileInfo struct {
	FileName  string
	hunkStart int
	hunkEnd   int
	sha       string
}

func getCommitFileInfo(ghConnector *connector) ([]*commitFileInfo, error) {

	prFiles, err := ghConnector.getFilesForPr()
	if err != nil {
		return nil, err
	}

	var (
		errs            []string
		commitFileInfos []*commitFileInfo
	)

	for _, file := range prFiles {
		info, err := getCommitInfo(file)
		if err != nil {
			errs = append(errs, err.Error())
			continue
		}
		commitFileInfos = append(commitFileInfos, info)
	}
	if len(errs) > 0 {
		return nil, fmt.Errorf("there were errors processing the PR files.\n%s", strings.Join(errs, "\n"))
	}
	return commitFileInfos, nil
}

func (cfi commitFileInfo) calculatePosition(line int) *int {
	position := line - cfi.hunkStart
	return &position
}
