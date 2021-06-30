package golang

import (
	"bytes"
	_ "embed" // for go:embed file use
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
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

	//go:embed golangci-lint/dev-ci.yml
	devci []byte
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

func InitProject(goSet *flag.FlagSet, cilintflag, cilintProjflag *bool) (suggs []*util.Suggestion, err error) {
	// nolint // flag.ExitOnError will do the os.Exit(2)
	goSet.Parse(os.Args[2:])

	ff := initFoldersAndFiles(createFolders, filesAndContent)

	if *cilintflag && *cilintProjflag {
		// 如果两个选项都有，则报错
		return nil, errors.New("can not setup golangci-lint globally and locally at same time")
	} else if *cilintflag && !*cilintProjflag {
		// 设置 global golangci-lint
		err = ff.initProjectWithGlobalLint()
	} else if !*cilintflag && *cilintProjflag {
		// 设置 project golangci-lint
		err = ff.initLocalCiLint()
	} else {
		// 不设置 golangci-lint
		err = ff.initProjectWithoutLint()
	}

	if err != nil {
		return nil, err
	}

	// 写入所需文件
	fmt.Println("init Golang project")
	if err := ff.writeAllFiles(); err != nil {
		return nil, err
	}

	// 检查返回是否为空
	if len(ff.suggestions) != 0 {
		return ff.suggestions, nil
	}

	return nil, nil
}

// 不设置 golangci-lint
func (ff *foldersAndFiles) initProjectWithoutLint() error {
	// 不需要设置 cilint 的情况，直接写 setting
	err := ff.addSettingJSON()
	if err != nil {
		return err
	}
	return nil
}

// 设置 local golangci-lint:
//  - 写入 <project>/golangci/dev-ci.yml
//  - 写入 <project>/.vscode/settings.json 文件
func (ff *foldersAndFiles) initLocalCiLint() error {
	// 获取本项目的绝对地址
	projectPath, err := filepath.Abs(".")
	if err != nil {
		return err
	}

	// 添加 <project>/golangci 文件夹，添加 dev-ci.yml 文件
	ff.addCilintYMLAndCipath(projectPath)

	// 设置 settings.json 文件, 将 config 设置为 cilint 配置文件地址
	err = ff.addSettingJSON()
	if err != nil {
		return err
	}

	return nil
}

// 设置 global golangci-lint:
//  - 写入 ~/.vsc/golangci/dev-ci.yml 文件.
//  - 写入 <project>/.vscode/settings.json 文件.
//  - 写入 ~/.vsc/vsc-config.json 全局配置文件.
func (ff *foldersAndFiles) initProjectWithGlobalLint() error {
	// 获取 .vsc 文件夹地址
	vscDir, err := util.GetVscConfigDir()
	if err != nil {
		return err
	}

	// 从 vsc-config.json 文件获取 golangci 配置文件的地址。
	err = ff.readCilintPathFromVscCfgJSON(vscDir)
	if err != nil {
		return err
	}

	// 设置 settings.json 文件, 将 --config 设置为 cipath
	err = ff.addSettingJSON()
	if err != nil {
		return err
	}
	return nil
}

// 添加 .vscode/settings.json 文件，如果文件存在则给出建议
func (ff *foldersAndFiles) addSettingJSON() error {
	if ff.cipath == "" {
		// 不设置 golangci-lint 的情况
		ff._addFiles(newSettingsJSONwith(""))
		return nil
	}

	// 读取 .vscode/settings.json, 获取 "go.lintFlags" 的值
	golingFlags, err := _readSettingJSON()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// settings.json 不存在, 生成新的 settings.json 文件
		ff._addFiles(newSettingsJSONwith(ff.cipath))
		return nil
	}

	// 判断 --config 地址是否和要设置的 cipath 相同, 如果相同则不更新 setting 文件。
	for _, v := range golingFlags {
		if v == "--config="+ff.cipath { // 相同的路径
			return nil
		}
	}

	// 如果 settings.json 文件存在，而且 config != cipath, 则需要 suggestion
	// 建议手动添加设置到 .vscode/settings.json 中
	cilintConfig := bytes.ReplaceAll(golangcilintconfig, []byte(configPlaceHolder), []byte(ff.cipath))
	ff._addSuggestion(&util.Suggestion{
		Problem:  "please add following in '.vscode/settings.json':",
		Solution: string(cilintConfig),
	})

	return nil
}

// 读取 .vscode/settings.json, 获取 "go.lintFlags" 的值
func _readSettingJSON() ([]string, error) {
	// 读取 .vscode/settings.json
	settingsPath, err := filepath.Abs(util.SettingsJSONPath)
	if err != nil {
		return nil, err
	}

	sf, err := os.Open(settingsPath)
	if err != nil {
		return nil, err
	}
	defer sf.Close()

	// json 反序列化 settings.json
	jsonc, err := io.ReadAll(sf)
	if err != nil {
		return nil, err
	}

	js, err := util.JSONCToJSON(jsonc)
	if err != nil {
		return nil, err
	}

	// 只需要读取 go.lintFlags
	type settingsStruct struct {
		GolingFlags []string `json:"go.lintFlags,omitempty"`
	}

	var settings settingsStruct
	err = json.Unmarshal(js, &settings)
	if err != nil {
		return nil, err
	}

	return settings.GolingFlags, nil
}

// 写入所有文件
func (ff *foldersAndFiles) writeAllFiles() error {
	return util.WriteFoldersAndFiles(ff.folders, ff.files)
}
