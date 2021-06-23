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
	// parse arges first
	// nolint // flag.ExitOnError will do the os.Exit(2)
	tsjsSet.Parse(os.Args[2:])

	folders := createFolders
	files := filesAndContent

	if *jestflag {
		// 检查 npm 是否安装
		if err := util.CheckCMDInstall("npm"); err != nil {
			return err
		}

		// add jest example test file
		folders = append(folders, testFolder)
		files = append(files, jestFileContent)
	}

	// NOTE write project files first
	fmt.Println("init TypeScript project")
	if err := util.WriteFoldersAndFiles(folders, files); err != nil {
		return err
	}

	// 安装依赖
	if *jestflag {
		// 设置 jest，检查依赖
		npmLibs, err := dependenciesNeedsToInstall()
		if err != nil {
			return err
		}

		// 下载依赖到项目中
		if err := util.NpmInstallDependencies("", npmLibs...); err != nil {
			return err
		}
	}

	return nil
}
