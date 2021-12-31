package golang

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

	//go:embed cfgfiles/gitignore
	gitignore []byte

	//go:embed cfgfiles/golangci.yml
	golangciYML []byte
)

var mainGO = []byte(`package main

import "fmt"

func main() {
	fmt.Println("hello world")
	// need to run "go mod init <name>" first.
}
`)

var filesAndContent = []util.FileContent{
	{Path: util.LaunchJSONPath, Content: launchJSON},
	{Path: util.GitignorePath, Content: gitignore},
	{Path: "src/main.go", Content: mainGO},
}

func InitProject(gofs util.GoFlags) (suggs []*util.Suggestion, err error) {
	// nolint // flag.ExitOnError will do the os.Exit(2)
	gofs.FlagSet.Parse(os.Args[3:])

	ctx := util.InitFoldersAndFiles(createFolders, filesAndContent)

	if *gofs.Cilint && *gofs.CilintLocal {
		// 如果两个选项都有，则报错
		return nil, errors.New("can not setup golangci-lint globally and locally at same time")
	} else if *gofs.Cilint && !*gofs.CilintLocal {
		// 设置 global golangci-lint
		err = initGlobalLint(ctx)
	} else if !*gofs.Cilint && *gofs.CilintLocal {
		// 设置 project golangci-lint
		err = initLocalCiLint(ctx)
	} else {
		// 不设置 golangci-lint
		err = initWithoutLint(ctx)
	}

	if err != nil {
		return nil, err
	}

	// 写入所需文件
	fmt.Println("init Golang project")
	if err := ctx.WriteAllFiles(); err != nil {
		return nil, err
	}

	return ctx.Suggestions(), nil
}

// 不设置 golangci-lint
func initWithoutLint(ctx *util.VSContext) error {
	// 不需要设置 cilint 的情况，直接写 setting
	err := addSettingJSON(ctx)
	if err != nil {
		return err
	}
	return nil
}

// 设置 local golangci-lint:
//  - 写入 <project>/golangci/golangci.yml
//  - 写入 <project>/.vscode/settings.json 文件
func initLocalCiLint(ctx *util.VSContext) error {
	// 获取本项目的绝对地址
	projectPath, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	// 添加 <project>/golangci 文件夹，添加 golangci.yml 文件
	ctx.AddLintConfigAndLintPath(projectPath+cilintFilePath, golangciYML)

	// 设置 settings.json 文件, 将 config 设置为 cilint 配置文件地址
	err = addSettingJSON(ctx)
	if err != nil {
		return err
	}

	return nil
}

// 设置 global golangci-lint:
//  - 写入 ~/.vsi/golangci/golangci.yml 文件.
//  - 写入 <project>/.vscode/settings.json 文件.
//  - 写入 ~/.vsi/vsi-config.json 全局配置文件.
func initGlobalLint(ctx *util.VSContext) error {
	// 获取 .vsi 文件夹地址
	vsiDir, err := util.GetVsiConfigDir()
	if err != nil {
		return err
	}

	// 从 vsi-config.json 文件获取 golangci 配置文件的地址。
	err = readCilintPathFromVsiCfgJSON(ctx, vsiDir)
	if err != nil {
		return err
	}

	// 设置 settings.json 文件, 将 --config 设置为 cipath
	err = addSettingJSON(ctx)
	if err != nil {
		return err
	}
	return nil
}
