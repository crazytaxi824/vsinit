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
	// {Path: ".vscode/settings.json", Content: settingsJSON},
	{Path: ".gitignore", Content: gitignore},
	{Path: "src/main.go", Content: mainGO},
}

func InitProject(goSet *flag.FlagSet, cilintflag, cilintProjflag *bool) (suggs []*util.Suggestion, err error) {
	// nolint // flag.ExitOnError will do the os.Exit(2)
	goSet.Parse(os.Args[2:])

	if *cilintflag && *cilintProjflag {
		return nil, errors.New("can not setup golangci-lint globally and locally at same time")
	} else if *cilintflag && !*cilintProjflag {
		return initProject(true, false)
	} else if !*cilintflag && *cilintProjflag {
		return initProject(true, true)
	}
	return initProject(false, false)
}

func initProject(cilint, local bool) (suggs []*util.Suggestion, err error) {
	folders := createFolders
	files := filesAndContent

	if !cilint { // 如果不需要设置 cilint
		settingJSON := genNewSettingsFile("")
		files = append(files, util.FileContent{
			Path:    ".vscode/settings.json",
			Content: settingJSON,
		})
	} else {
		var (
			gls *golangciLintStruct
		)

		if local { // 设置项目 golangci lint
			projectPath, er := filepath.Abs(".")
			if er != nil {
				return nil, er
			}
			// 设置需要创建的文件夹和要写的 golangci.yml 文件
			gls = setupLocalCilint(projectPath)
		} else { // 设置 global lint
			// 设置需要创建的文件夹和要写的 golangci.yml 文件
			gls, err = setupGlobleCilint()
			if err != nil {
				return nil, err
			}
		}

		// 将 dev-ci.yml prod-ci.yml 配置文件都设为需要创建和写入
		folders = append(folders, gls.Folders...)
		files = append(files, gls.Files...)

		// 设置 cipath 到 setting.json 中
		// settingJSON, overwrite, sug, er := checkSettingsJSONfileExist(cipath)
		settingJSON, sug, er := checkSettingJSONExist(gls.Cipath)
		if er != nil {
			return nil, er
		}
		if sug != nil {
			suggs = append(suggs, sug)
		}
		// 添加 settings.json 文件
		if settingJSON != nil {
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

func checkSettingJSONExist(ciPath string) (newSetting []byte, sug *util.Suggestion, err error) {
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

	// 如果 settings.json 文件存在，需要 suggestion
	jsonc, err := io.ReadAll(sf)
	if err != nil {
		return nil, nil, err
	}

	js, err := util.JSONCToJSON(jsonc)
	if err != nil {
		return nil, nil, err
	}

	type settingsStruct struct {
		GolingFlags []string `json:"go.lintFlags,omitempty"`
	}

	var settings settingsStruct
	err = json.Unmarshal(js, &settings)
	if err != nil {
		return nil, nil, err
	}

	// 判断 --config 地址是否和要设置的 cipath 相同
	for _, v := range settings.GolingFlags {
		if v == "--config="+ciPath { // 相同的路径
			return nil, nil, nil
		}
	}

	r := bytes.ReplaceAll(golangcilintconfig, []byte(configPlaceHolder), []byte(ciPath))

	return nil, &util.Suggestion{
		Problem:  "please add following in '.vscode/settings.json':",
		Solution: string(r),
	}, nil
}
