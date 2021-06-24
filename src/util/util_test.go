package util

import (
	"fmt"
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
