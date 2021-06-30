package util

import "fmt"

const (
	// internal error
	InternalErrMsg = "CMD is not in the list, please contact author"

	// FIXME
	GolintciCmd = "vsc setup go -golangci <path>"
)

// 固定文件路径
const (
	SettingsJSONPath = ".vscode/settings.json"
	LaunchJSONPath   = ".vscode/launch.json"
	TasksJSONPath    = ".vscode/tasks.json"
	GitignorePath    = ".gitignore"
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
