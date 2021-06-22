package golang

import (
	_ "embed" // for go:embed file use

	"local/src/util"
)

var CreateFolders = []string{".vscode", "src"}

var (
	//go:embed cfgfiles/launch.json
	launchJSON []byte

	//go:embed cfgfiles/settings.json
	settingsJSON []byte

	//go:embed cfgfiles/gitignore
	gitignore []byte
)

var mainGO = []byte(`package main

import "fmt"

func main() {
	fmt.Println("hello world")
    // need to run "go mod init" first.
}
`)

var FilesAndContent = []util.FileContent{
	{
		Path:    ".vscode/launch.json",
		Content: launchJSON,
	},
	{
		Path:    ".vscode/settings.json",
		Content: settingsJSON,
	},
	{
		Path:    ".gitignore",
		Content: gitignore,
	},
	{
		Path:    "src/main.go",
		Content: mainGO,
	},
}
