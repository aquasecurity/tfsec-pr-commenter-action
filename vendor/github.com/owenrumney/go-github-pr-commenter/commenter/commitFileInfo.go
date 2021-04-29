package commenter

type commitFileInfo struct {
	FileName  string
	hunkStart int
	hunkEnd   int
	sha       string
}

func (cfi commitFileInfo) calculatePosition(line int) *int {
	position := line - cfi.hunkStart
	return &position
}
