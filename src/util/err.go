package util

import "fmt"

const (
	InternalErrMsg = "CMD is not in the list, please contact author"

	GolintciCmd = "vsc setup go -golangci <path>"
)

type Suggestion struct {
	Problem  string
	Solution string
}

func (e *Suggestion) String() string {
	return warn(">>> "+e.Problem) + "\n" + e.Solution + "\n\n"
}

func warn(m string) string {
	return fmt.Sprintf("\033[0;37;41m%s\033[0m", m)
}
