package util

import (
	"errors"
	"fmt"
	"os"
	"strings"
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
	fmt.Println(CheckCMDInstall("vscode", "omg", "haha"))
}

func Test_MakeDirAlreadyExist(t *testing.T) {
	err := os.Mkdir("/Users/ray/Desktop/test", 0750)
	if err != nil && !errors.Is(err, os.ErrExist) {
		t.Error(err)
		return
	} else if errors.Is(err, os.ErrExist) {
		t.Log("ok")
	}
}

func Test_ReadJSON(t *testing.T) {
	var vs VscSetting
	err := vs.readJSON(strings.NewReader(`{"golangci":"abc","eslint":"def"}`))
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("%+v", vs)
}
