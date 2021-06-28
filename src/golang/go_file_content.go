package golang

import (
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

	//go:embed cfgfiles/settings.json
	settingsJSON []byte

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
    // need to run "go mod init" first.
}
`)

// var filesAndContent = []util.FileContent{
// 	{Path: ".vscode/launch.json", Content: launchJSON},
// 	{Path: ".vscode/settings.json", Content: settingsJSON},
// 	{Path: ".gitignore", Content: gitignore},
// 	{Path: "src/main.go", Content: mainGO},
// }

// func InitProject() error {
// 	folders := createFolders
// 	files := filesAndContent

// 	fmt.Println("init Golang project")
// 	return util.WriteFoldersAndFiles(folders, files)
// }

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
			fos    []string
			fis    []util.FileContent
			cipath string
		)

		if local { // 设置项目 golangci lint
			projectPath, er := filepath.Abs(".")
			if er != nil {
				return nil, er
			}
			fos, fis, cipath = setupLocalCilint(projectPath)
		} else { // 设置 global lint
			fos, fis, cipath, err = setupGlobleCilint()
			if err != nil {
				return nil, err
			}
		}

		folders = append(folders, fos...)
		files = append(files, fis...)

		// 设置 cipath 到 setting.json 中
		settingJSON, overwrite, sug, er := checkSettingsJSONfileExist(cipath)
		if er != nil {
			return nil, er
		}
		if sug != nil {
			suggs = append(suggs, sug)
		}

		files = append(files, util.FileContent{
			Path:      ".vscode/settings.json",
			Content:   settingJSON,
			Overwrite: overwrite,
		})
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

// 查看 .vscode/settings.json 文件是否存在
func checkSettingsJSONfileExist(cipath string) (newSetting []byte, overwrite bool, sug *util.Suggestion, err error) {
	settingsPath, err := filepath.Abs(".vscode/settings.json")
	if err != nil {
		return nil, false, nil, err
	}

	sf, err := os.Open(settingsPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, false, nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		// settings.json 不存在, 生成新的 settings.json 文件
		return genNewSettingsFile(cipath), false, nil, nil
	}
	defer sf.Close()

	// 如果 settings.json 文件存在，需要 overwrite
	overwrite = true

	// 查看 lint 设置是否存在
	jsonc, err := io.ReadAll(sf)
	if err != nil {
		return nil, false, nil, err
	}

	// 读取 setting.json 文件
	js, err := util.JSONCToJSON(jsonc)
	if err != nil {
		return nil, false, nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(js, &m)
	if err != nil {
		return nil, false, nil, err
	}

	// 查看 lintTool, lintOnSave, lintFlags 设置是否存在
	var lints [][]byte
	if _, ok := m["go.lintFlags"]; !ok {
		lints = append(lints, lintFlags)
	} else {
		// 如果已经设置了 go.lintFlags, 修改 --config 设置。
		// 如果 --config 不存在，通过 suggestion 提示。
		newSetting, sug, err = replaceCilintConfigPath(jsonc, cipath)
		if err != nil {
			return nil, false, nil, err
		}
	}

	if _, ok := m["go.lintTool"]; !ok {
		lints = append(lints, lintTool)
	}
	if _, ok := m["go.lintOnSave"]; !ok {
		lints = append(lints, lintOnSave)
	}

	// append 到 setting.json 中
	newSetting, err = util.AppendToJSONC(newSetting, golangciSettings(lints...))
	if err != nil {
		return nil, false, nil, err
	}

	return newSetting, overwrite, sug, nil
}
