package util

const (
	InternalErrMsg = "CMD is not in the list, please contact author"

	GolintciCmd = "vsc setup go -golangci <path>"
)

type Suggestion struct {
	Problem  string
	Solution string
}

func (e *Suggestion) String() string {
	return Warn(">>> "+e.Problem) + "\n" + e.Solution + "\n\n"
}
