package python

import (
	_ "embed" // for go:embed file use
	"fmt"
	"os"

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
	{Path: util.LaunchJSONPath, Content: launchJSON},
	{Path: util.SettingsJSONPath, Content: settingsJSON},
	{Path: util.GitignorePath, Content: gitignore},
	{Path: "src/main.py", Content: mainPY},
}

func InitProject() error {
	if len(os.Args) > 3 && (os.Args[3] == "-h" || os.Args[3] == "-help" || os.Args[3] == "--help") {
		fmt.Println("Usage of 'vs init py': no flags")
		os.Exit(0)
	} else if len(os.Args) > 3 {
		util.HelpMsg()
		os.Exit(2)
	}

	ff := util.InitFoldersAndFiles(createFolders, filesAndContent)

	fmt.Println("init Python project")
	return ff.WriteAllFiles()
}
