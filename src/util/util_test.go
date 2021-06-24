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
	var errs Erros
	errs = append(errs, ErrorMsg{"a", []string{"a", "b"}}, ErrorMsg{"a", []string{"a", "b"}})

	var errs2 Erros
	errs2 = append(errs2, errs, ErrorMsg{"b", []string{"a", "b"}}, ErrorMsg{"b", []string{"a", "b"}})

	fmt.Println(errs2)
}

func Test_CheckCmdInstall(t *testing.T) {
	fmt.Println(CheckCMDInstall("vscode", "omg", "haha"))
}
