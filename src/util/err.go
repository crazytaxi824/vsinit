package util

import (
	"strings"
)

const (
	InternalErrMsg = "CMD is not in the list, please contact author"

	GolintciCmd = "vsc setup go -golangci <path>"
)

type ErrorMsg struct {
	Problem  string
	Solution []string
}

func (e ErrorMsg) Error() string {
	return Warn(">>> "+e.Problem) + "\n" + strings.Join(e.Solution, "\n") + "\n\n"
}

type Erros []error

func (es Erros) Error() string {
	var builder strings.Builder
	for _, err := range es {
		builder.WriteString(err.Error())
	}
	return builder.String()
}
