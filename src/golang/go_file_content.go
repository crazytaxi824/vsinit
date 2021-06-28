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

	//go:embed golangci-lint/prod-ci.yml
	prodci []byte
)

var mainGO = []byte(`package main

import "fmt"

func main() {
	fmt.Println("hello world")
    // need to run "go mod init <name>" first.
}
`)

var filesAndContent = []util.FileContent{
	{Path: ".vscode/launch.json", Content: launchJSON},
	{Path: ".gitignore", Content: gitignore},
	{Path: "src/main.go", Content: mainGO},
}

func InitProject(goSet *flag.FlagSet, cilintflag, cilintProjflag *bool) (suggs []*util.Suggestion, err error) {
	// nolint // flag.ExitOnError will do the os.Exit(2)
	goSet.Parse(os.Args[2:])

	if *cilintflag && *cilintProjflag {
		// 如果两个选项都有，则报错
		return nil, errors.New("can not setup golangci-lint globally and locally at same time")
	} else if *cilintflag && !*cilintProjflag {
		// 设置 global golangci-lint
		return initProject(true, false)
	} else if !*cilintflag && *cilintProjflag {
		// 设置 project golangci-lint
		return initProject(true, true)
	}

	// 不设置 golangci-lint
	return initProject(false, false)
}

func initProject(cilint, local bool) (suggs []*util.Suggestion, err error) {
	folders := createFolders
	files := filesAndContent

	if !cilint {
		// 不需要设置 cilint 的情况，直接写 setting
		settingJSON := genNewSettingsFile("")
		files = append(files, util.FileContent{
			Path:    ".vscode/settings.json",
			Content: settingJSON,
		})
	} else {
		// 需要设置 cilint 的情况
		var gls *golangciLintStruct

		if local { // 设置 project golangci lint
			projectPath, er := filepath.Abs(".")
			if er != nil {
				return nil, er
			}
			// 添加 <project>/golangci 文件夹，添加 dev-ci.yml, prod-ci.yml 文件
			gls = setupLocalCilint(projectPath)
		} else { // 设置 global lint
			// 添加 ~/.vsc/golangci 文件夹，添加 dev-ci.yml, prod-ci.yml 文件
			gls, err = setupGlobleCilint()
			if err != nil {
				return nil, err
			}
		}

		// 将 dev-ci.yml prod-ci.yml 配置文件都设为需要创建和写入
		folders = append(folders, gls.Folders...)
		files = append(files, gls.Files...)

		// 设置 settings.json 文件, 将 --config 设置为 cipath
		settingJSON, sug, er := checkSettingJSON(gls.Cipath)
		if er != nil {
			return nil, er
		}
		if sug != nil {
			suggs = append(suggs, sug)
		}
		if settingJSON != nil {
			// 添加 settings.json 文件到写入队列中
			files = append(files, util.FileContent{
				Path:    ".vscode/settings.json",
				Content: settingJSON,
			})
		}
	}

	fmt.Println("init Golang project")
	err = util.WriteFoldersAndFiles(folders, files)
	if err != nil {
		return nil, err
	}

	// 检查返回是否为空
	if len(suggs) == 0 {
		return nil, nil
	}

	return suggs, nil
}

// 检查 .vscode/settings.json 是否存在
func checkSettingJSON(ciPath string) (newSetting []byte, sug *util.Suggestion, err error) {
	settingsPath, err := filepath.Abs(".vscode/settings.json")
	if err != nil {
		return nil, nil, err
	}

	sf, err := os.Open(settingsPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		// settings.json 不存在, 生成新的 settings.json 文件
		return genNewSettingsFile(ciPath), nil, nil
	}
	defer sf.Close()

	// json 反序列化 settings.json
	golingFlags, err := readSettingJSON(sf)
	if err != nil {
		return nil, nil, err
	}

	// 判断 --config 地址是否和要设置的 cipath 相同, 如果相同则不更新 setting 文件。
	for _, v := range golingFlags {
		if v == "--config="+ciPath { // 相同的路径
			return nil, nil, nil
		}
	}

	// 如果 settings.json 文件存在，而且 config != cipath, 则需要 suggestion
	// 建议手动添加设置到 .vscode/settings.json 中
	cilintConfig := bytes.ReplaceAll(golangcilintconfig, []byte(configPlaceHolder), []byte(ciPath))
	return nil, &util.Suggestion{
		Problem:  "please add following in '.vscode/settings.json':",
		Solution: string(cilintConfig),
	}, nil
}

func readSettingJSON(file *os.File) ([]string, error) {
	// json 反序列化 settings.json
	jsonc, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	js, err := util.JSONCToJSON(jsonc)
	if err != nil {
		return nil, err
	}

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
