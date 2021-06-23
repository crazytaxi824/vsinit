package golang

import (
	"path/filepath"
	"testing"
)

func Test_AbsPath(t *testing.T) {
	fpath := "abc/def"
	t.Log(filepath.Abs(fpath))
}
