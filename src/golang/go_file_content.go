package golang

import _ "embed" // for go:embed file use

var (
	//go:embed cfgfiles/launch.json
	launchJSONContent string

	//go:embed cfgfiles/settings.json
	settingsJSONContent string

	//go:embed cfgfiles/gitignore
	gitignoreContent string
)

var FilesAndContent = map[string]string{
	".vscode/launch.json":   launchJSONContent,
	".vscode/settings.json": settingsJSONContent,
	"src/main.go":           mainFileContent,
	".gitignore":            gitignoreContent,
}

const mainFileContent = `package main

import "fmt"

func main() {
	fmt.Println("hello world")
}
`
