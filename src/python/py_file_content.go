package python

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

var mainPY = []byte(`def main():
    print("hello world")


main()
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
		Path:    "src/main.py",
		Content: mainPY,
	},
}
