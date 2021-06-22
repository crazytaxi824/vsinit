// 纯 js 项目用, 不包含 react lint

package js

import (
	_ "embed" // for go:embed file use
	"flag"
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

	//go:embed cfgfiles/package.json
	packageJSON []byte

	// for unit test 'jest' use
	//go:embed cfgfiles/example.test.js
	exampleTestJS []byte
)

// file content
var mainJS = []byte(`main();

function main() {
  console.log('hello world');
}
`)

// filesAndContent JS project files
var filesAndContent = []util.FileContent{
	{Path: ".vscode/launch.json", Content: launchJSON},
	{Path: ".vscode/settings.json", Content: settingsJSON},
	{Path: ".gitignore", Content: gitignore},
	{Path: "package.json", Content: packageJSON},
	{Path: "src/main.js", Content: mainJS},
}

// for jest use only

const testFolder = "test"

var jestFileContent = util.FileContent{
	Path:    testFolder + "/example.test.js",
	Content: exampleTestJS,
}

func InitProject(tsjsSet *flag.FlagSet, jestflag *bool) error {
	if err := util.CheckCMDInstall("node"); err != nil {
		return err
	}

	// parse arges first
	// nolint // flag.ExitOnError will do the os.Exit(2)
	tsjsSet.Parse(os.Args[2:])

	folders := createFolders
	files := filesAndContent
	if *jestflag {
		// add jest example test file
		folders = append(folders, testFolder)
		files = append(files, jestFileContent)
	}

	fmt.Println("init JavaScript project")
	return util.WriteFoldersAndFiles(folders, files)
}
