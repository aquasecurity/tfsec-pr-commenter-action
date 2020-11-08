package commenter

type CommitFileInfo struct {
	FileName  string
	hunkStart int
	hunkEnd   int
	sha       string
}

func (cfi CommitFileInfo) CalculatePosition(line int) *int {
	position := line - cfi.hunkStart
	return &position
}
