package golang

import _ "embed" // for go:embed file use

var (
	//go:embed cfgfiles/launch.json
	launchJSONContent []byte

	//go:embed cfgfiles/settings.json
	settingsJSONContent []byte

	//go:embed cfgfiles/gitignore
	gitignoreContent []byte
)

var FilesAndContent = map[string][]byte{
	".vscode/launch.json":   launchJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"src/main.go":           mainFileContent,
	".gitignore":            gitignoreContent,
}

var mainFileContent = []byte(`package main

import "fmt"

func main() {
	fmt.Println("hello world")
  // need to run "go mod init" first.
}
`)
