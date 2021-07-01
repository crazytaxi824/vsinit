package util

import (
	"fmt"
	"unsafe"
)

const (
	// internal error
	InternalErrMsg = "CMD is not in the list, please contact author"
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

// nolint // unsafe
func StringToBytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

// nolint // unsafe
func BytesToString(b []byte) string {
	x := (*[3]uintptr)(unsafe.Pointer(&b))
	s := [2]uintptr{x[0], x[1]}
	return *(*string)(unsafe.Pointer(&s))
}
