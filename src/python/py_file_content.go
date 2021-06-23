package python

import (
	_ "embed" // for go:embed file use
	"fmt"

	"local/src/util"
)

var createFolders = []string{".vscode", "src"}

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

var filesAndContent = []util.FileContent{
	{Path: ".vscode/launch.json", Content: launchJSON},
	{Path: ".vscode/settings.json", Content: settingsJSON},
	{Path: ".gitignore", Content: gitignore},
	{Path: "src/main.py", Content: mainPY},
}

func InitProject() error {
	folders := createFolders
	files := filesAndContent

	fmt.Println("init Python project")
	return util.WriteFoldersAndFiles(folders, files)
}
