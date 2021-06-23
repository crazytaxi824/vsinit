package golang

import (
	"os"
	"path/filepath"
	"testing"
)

func Test_AbsPath(t *testing.T) {
	fpath := "abc/def"
	t.Log(filepath.Abs(fpath))
}

func Test_CheckShell(t *testing.T) {
	t.Log(os.Getenv("SHELL"))
	t.Log(os.Getenv("GOBIN"))
	if os.Getenv("abcdef") != "" {
		t.Fail()
	}
}
