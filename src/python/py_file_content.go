package python

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

var mainFileContent = []byte(`def main():
    print("hello world")


main()
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
		Path:    "src/main.py",
		Content: mainFileContent,
	},
}
