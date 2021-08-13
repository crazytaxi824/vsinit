package ts

import (
	_ "embed" // for go:embed file use
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"local/src/util"
)

var createFolders = []string{".vscode", "src"}

var (
	//go:embed cfgfiles/launch.json
	launchJSON []byte

	//go:embed cfgfiles/settings_template.txt
	settingTemplate []byte

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

	//go:embed cfgfiles/eslintrc-ts.json
	eslintrcJSON []byte
)

var mainTS = []byte(`main();

function main() {
  console.log('hello world');
}
`)

var filesAndContent = []util.FileContent{
	{Path: util.LaunchJSONPath, Content: launchJSON},
	{Path: util.TasksJSONPath, Content: tasksJSON},
	{Path: util.GitignorePath, Content: gitignore},
	{Path: "package.json", Content: packageJSON},
	{Path: "tsconfig.json", Content: tsconfigJSON},
	{Path: "src/main.ts", Content: mainTS},
}

func InitProject(tsjs util.TSJSFlags) (suggs []*util.Suggestion, err error) {
	// parse arges first
	// nolint // flag.ExitOnError will do the os.Exit(2)
	tsjs.FlagSet.Parse(os.Args[3:])

	// 初始化
	ctx := util.InitFoldersAndFiles(createFolders, filesAndContent)

	// 写入 test 相关文件
	if *tsjs.Jest {
		err = initJest(ctx)
		if err != nil {
			return nil, err
		}
	}

	if *tsjs.ESLint && *tsjs.ESLintLocal {
		// 如果两个选项都有，则报错
		return nil, errors.New("can not setup eslint globally and locally at same time")
	} else if *tsjs.ESLint && !*tsjs.ESLintLocal {
		// 设置 global eslint
		err = initGlobalEslint(ctx)
	} else if !*tsjs.ESLint && *tsjs.ESLintLocal {
		// 设置 local eslint
		err = initLocalEslint(ctx)
	} else {
		// 不设置 eslint, 只需要设置 settings.json 文件
		err = initWithoutEslint(ctx)
	}

	if err != nil {
		return nil, err
	}

	// 写入所需文件
	fmt.Println("init TypeScript project")
	if err := ctx.WriteAllFiles(); err != nil {
		return nil, err
	}

	// 安装所有缺失的依赖
	if err := ctx.InstallMissingDependencies(); err != nil {
		return nil, err
	}

	return ctx.Suggestions(), nil
}

// 不设置 ESLint, 写入 <project>/.vscode/settings.json 文件.
func initWithoutEslint(ctx *util.VSContext) error {
	// 直接写 settings.json 文件
	err := addSettingJSON(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 设置 local ESLint:
//  - 写入 <project>/eslint/eslintrc-ts.json 本地配置文件.
//  - 写入 <project>/.vscode/settings.json 文件.
//  - 安装 ESLint 缺失的本地依赖.
func initLocalEslint(ctx *util.VSContext) error {
	// 检查 npm 是否安装，把 suggestion 当 error 返回，因为必须要安装依赖
	if sugg := util.CheckCMDInstall("npm"); sugg != nil {
		return errors.New(sugg.String())
	}

	// 获取项目的绝对地址
	projectPath, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	// 添加 <project>/eslint 文件夹，添加 eslintrc-ts.json 文件
	// ctx.addEslintJSONAndEspath(projectPath + eslintFilePath)
	ctx.AddLintConfigAndLintPath(projectPath+eslintFilePath, eslintrcJSON)

	// 设置 settings.json 文件, 将 config 设置为 eslint 配置文件地址
	err = addSettingJSON(ctx)
	if err != nil {
		return err
	}

	// 添加 ESLint 缺失的本地依赖
	return ctx.AddMissingDependencies(eslintDependencies, "package.json", "")
}

// 设置 global ESLint:
//  - 写入 ~/.vsi/eslint/eslintrc-ts.json 全局配置文件.
//  - 写入 ~/.vsi/vsi-config.json 全局配置文件.
//  - 写入 <project>/.vscode/settings.json 文件.
//  - 安装 ESLint 缺失的全局依赖.
func initGlobalEslint(ctx *util.VSContext) error {
	// 检查 npm 是否安装，把 suggestion 当 error 返回，因为必须要安装依赖
	if sugg := util.CheckCMDInstall("npm"); sugg != nil {
		return errors.New(sugg.String())
	}

	// 获取 ~/.vsi 文件夹地址
	vsiDir, err := util.GetVsiConfigDir()
	if err != nil {
		return err
	}

	// 通过 vsi-config.json 获取 eslint.TS 配置文件地址.
	err = readEslintPathFromVsiCfgJSON(ctx, vsiDir)
	if err != nil {
		return err
	}

	// 设置 settings.json 文件, 将 configFile 设置为 eslint 配置文件地址
	err = addSettingJSON(ctx)
	if err != nil {
		return err
	}

	// 添加 ESLint 缺失的全局依赖
	eslintFolder := vsiDir + eslintDirector
	pkgFilePath := eslintFolder + "/package.json"
	return ctx.AddMissingDependencies(eslintDependencies, pkgFilePath, eslintFolder)
}
