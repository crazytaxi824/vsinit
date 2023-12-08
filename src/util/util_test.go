package util

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestFilePath(t *testing.T) {
	fp := "a/b/c.ext"
	t.Log(filepath.Base(fp))
	t.Log(filepath.Ext(fp))
	t.Log(filepath.Dir(fp))

	fp2 := "a/b/c/"
	t.Log(filepath.Base(fp2))
	t.Log(filepath.Ext(fp2))
	t.Log(filepath.Dir(fp2))

	fp3 := "c.ext"
	t.Log(filepath.Base(fp3))
	t.Log(filepath.Ext(fp3))
	t.Log(filepath.Dir(fp3))

	fp4 := "/a/b/c.ext"
	t.Log(filepath.Base(fp4))
	t.Log(filepath.Ext(fp4))
	t.Log(filepath.Dir(fp4))

	t.Log(filepath.Abs("~/Desktop/"))
}

func TestCreateFile(t *testing.T) {
	f, err := os.OpenFile("/Users/ray/Desktop/ttttt/", os.O_CREATE, 0600)
	if err != nil {
		t.Error(err)
		return
	}
	f.Close()
}

func TestTrimSpace(t *testing.T) {
	s := " ab c\n "
	t.Log(s)
	t.Log(strings.TrimSpace(s))
}

func TestSplit(t *testing.T) {
	s := "a"
	t.Log(strings.Split(s, ","))
}

func TestStrconv(t *testing.T) {
	s := "1 "
	t.Log(strconv.Atoi(s))
}

func TestToLower(t *testing.T) {
	s := "1"
	t.Log(strings.ToLower(s))
}
