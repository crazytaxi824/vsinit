package util

import (
	"fmt"
	"path/filepath"
	"testing"
)

func Test_NilSlice(t *testing.T) {
	var ers []error
	if ers != nil {
		t.Log(ers)
	}

	ers = append(ers, nil)
	if ers != nil {
		t.Log(ers)
	}
}

func Test_Errors(t *testing.T) {
	var errs []*Suggestion
	errs = append(errs, &Suggestion{"a", "a\nb"}, &Suggestion{"a", "a\nb"})

	var errs2 []*Suggestion
	errs2 = append(errs2, errs...)
	errs2 = append(errs2, &Suggestion{"b", "a\nb"}, &Suggestion{"b", "a\nb"})

	fmt.Println(errs2)
}

func Test_CheckCmdInstall(t *testing.T) {
	fmt.Println(CheckCMDInstall("vscode"))
	fmt.Println(CheckCMDInstall("code"))
}

func Test_StringToByte(t *testing.T) {
	s := "abc"
	t.Log(StringToBytes(s))

	b := []byte("abc")
	t.Log(BytesToString(b))
}

func Test_IsPath(t *testing.T) {
	fs := []string{"-jest", "--jest", "-cilint", "--cilint", "-eslint", "--eslint"}

	args1 := []string{"vs", "init", "go", "-cilint"}
	args2 := []string{"vs", "init", "go", "-cilint", "-jest"}
	args3 := []string{"vs", "init", "go", "-cilint", "abc"}

	lastArg := func(args []string) string {
		if len(args) <= 3 {
			return "."
		}
		return args[len(args)-1]
	}

	inside := func(flags []string, arg string) bool {
		for _, v := range flags {
			if v == arg {
				return true
			}
		}
		return false
	}

	t.Log(inside(fs, lastArg(args1)))
	t.Log(inside(fs, lastArg(args2)))
	t.Log(inside(fs, lastArg(args3)))
}

func Test_filePathAbs(t *testing.T) {
	t.Log(filepath.Abs("."))
	t.Log(filepath.Abs("./"))
}
