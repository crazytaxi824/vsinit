package util

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

var src = ErrorMsg{
	"you haven't install 'code'",
	[]string{"solutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolution", "solutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolutionsolution"},
}

func Test_ErrorMsg(t *testing.T) {
	fmt.Println(src)
}

func Benchmark_ErrorMsg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		src.Error()
	}
	b.ReportAllocs()
}

func Benchmark_StringPlus(b *testing.B) {
	for i := 0; i < b.N; i++ {
		d := errors.New(src.Problem + "\n" + strings.Join(src.Solution, "\n"))
		_ = d
	}
	b.ReportAllocs()
}
