package ts

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

	//go:embed cfgfiles/tasks.json
	tasksJSON []byte

	//go:embed cfgfiles/gitignore
	gitignore []byte

	//go:embed cfgfiles/package.json
	packageJSON []byte

	//go:embed cfgfiles/tsconfig.json
	tsconfigJSON []byte

	//go:embed cfgfiles/example.test.ts
	exampleTestTS []byte
)

var mainTS = []byte(`main();

function main() {
  console.log('hello world');
}
`)

var filesAndContent = []util.FileContent{
	{Path: ".vscode/launch.json", Content: launchJSON},
	{Path: ".vscode/tasks.json", Content: tasksJSON},
	{Path: ".vscode/settings.json", Content: settingsJSON},
	{Path: ".gitignore", Content: gitignore},
	{Path: "package.json", Content: packageJSON},
	{Path: "tsconfig.json", Content: tsconfigJSON},
	{Path: "src/main.ts", Content: mainTS},
}

func InitProject(tsjsSet *flag.FlagSet, jestflag *bool) error {
	// 必须 node 和 typescript 都安装了.
	if err := util.CheckCMDInstall("node", "npm", "tsc"); err != nil {
		return err
	}

	// parse arges first
	// nolint // flag.ExitOnError will do the os.Exit(2)
	tsjsSet.Parse(os.Args[2:])

	folders := createFolders
	files := filesAndContent

	var npmLibs []string // Dependencies needs to be downloaded

	if *jestflag {
		// 检查 jest 是否安装
		if err := util.CheckCMDInstall("jest"); err != nil {
			return err
		}

		// add jest example test file
		folders = append(folders, testFolder)
		files = append(files, jestFileContent)

		// 设置 jest
		var err error
		npmLibs, err = setupJest()
		if err != nil {
			return err
		}
	}

	// NOTE write project files first
	fmt.Println("init TypeScript project")
	if err := util.WriteFoldersAndFiles(folders, files); err != nil {
		return err
	}

	// then npm install after wirte package.json file
	if err := util.NpmInstallDependencies("", npmLibs...); err != nil {
		return err
	}

	return nil
}
