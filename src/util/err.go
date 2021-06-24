package util

import "strings"

const ErrInternalMsg = "CMD is not in the list, please contact author"

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
