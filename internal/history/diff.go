package diff

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffFile compares the changes between two files
func DiffFile(oldContent, newContent string) []diffmatchpatch.Diff {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(oldContent, newContent, false)
	dmp.DiffCleanupSemantic(diffs)
	return diffs
}
