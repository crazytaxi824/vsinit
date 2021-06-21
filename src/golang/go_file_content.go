package golang

import (
	_ "embed" // for go:embed file use

	"local/src/util"
)

var CreateFolders = []string{".vscode", "src"}

var (
	//go:embed cfgfiles/launch.json
	launchJSONContent []byte

	//go:embed cfgfiles/settings.json
	settingsJSONContent []byte

	//go:embed cfgfiles/gitignore
	gitignoreContent []byte
)

var mainFileContent = []byte(`package main

import "fmt"

func main() {
	fmt.Println("hello world")
    // need to run "go mod init" first.
}
`)

var FilesAndContent = []util.FileContent{
	{
		Path:    ".vscode/launch.json",
		Content: launchJSONContent,
	},
	{
		Path:    ".vscode/settings.json",
		Content: settingsJSONContent,
	},
	{
		Path:    ".gitignore",
		Content: gitignoreContent,
	},
	{
		Path:    "src/main.go",
		Content: mainFileContent,
	},
}
