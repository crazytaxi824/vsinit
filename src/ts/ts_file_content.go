package ts

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

	//go:embed eslint/eslintrc-ts.json
	eslintrcJSON []byte
)

var mainTS = []byte(`main();

function main() {
  console.log('hello world');
}
`)

var filesAndContent = []util.FileContent{
	{Path: ".vscode/launch.json", Content: launchJSON},
	{Path: ".vscode/tasks.json", Content: tasksJSON},
	// {Path: ".vscode/settings.json", Content: settingsJSON},
	{Path: ".gitignore", Content: gitignore},
	{Path: "package.json", Content: packageJSON},
	{Path: "tsconfig.json", Content: tsconfigJSON},
	{Path: "src/main.ts", Content: mainTS},
}

func InitProject(tsjsSet *flag.FlagSet, jestflag, eslint, eslintLocal *bool) (suggs []*util.Suggestion, err error) {
	// parse arges first
	// nolint // flag.ExitOnError will do the os.Exit(2)
	tsjsSet.Parse(os.Args[2:])

	folders := createFolders
	files := filesAndContent

	if *jestflag {
		// 检查 npm 是否安装
		if sugg := util.CheckCMDInstall("npm"); sugg != nil {
			return nil, errors.New(sugg.String())
		}

		// add jest example test file
		folders = append(folders, testFolder)
		files = append(files, jestFileContent)
	}

	if *eslint && *eslintLocal {
		return nil, errors.New("can not setup eslint globally and locally at same time")
	} else if *eslint && !*eslintLocal {
		// 设置 global eslint
		fos, fis, sug, err := initProjectWithGlobalLint()
		if err != nil {
			return nil, err
		}
		folders = append(folders, fos...)
		files = append(files, fis...)
		suggs = sug
	} else if !*eslint && *eslintLocal {
		// 设置 local eslint
		fos, fis, sug, err := initProjectWithLocalLint()
		if err != nil {
			return nil, err
		}
		folders = append(folders, fos...)
		files = append(files, fis...)
		suggs = sug
	} else {
		// 不设置 eslint
		// 只需要设置 settings.json 文件
		files = append(files, initProjectWithoutLint())
	}

	// NOTE write project files first
	fmt.Println("init TypeScript project")
	if err := util.WriteFoldersAndFiles(folders, files); err != nil {
		return nil, err
	}

	// check and download dependencies
	if err := installAllDependencies(jestflag, eslint, eslintLocal); err != nil {
		return nil, err
	}

	// 检查返回是否为空
	if len(suggs) == 0 {
		return nil, nil
	}

	return suggs, nil
}

func installAllDependencies(jestflag, eslint, eslintLocal *bool) error {
	// NOTE 安装依赖, 必须放在后面，否则 package.json 需要改写。
	if *jestflag {
		// 设置 jest，检查依赖
		npmLibs, err := dependenciesNeedsToInstall(jestDependencies, "package.json")
		if err != nil {
			return err
		}

		// 下载依赖到项目中
		if err := util.NpmInstallDependencies("", false, npmLibs...); err != nil {
			return err
		}
	}

	// 下载 dependencies
	if *eslint { // global 情况
		vscDir, er := util.GetVscConfigDir()
		if er != nil {
			return er
		}

		eslintFolder := vscDir + util.EslintDirector
		pkgFilePath := eslintFolder + "/package.json"

		eslibs, err := dependenciesNeedsToInstall(eslintDependencies, pkgFilePath)
		if err != nil {
			return err
		}

		// 下载依赖到 ~/.vsc/eslint 中，
		if err := util.NpmInstallDependencies(eslintFolder, false, eslibs...); err != nil {
			return err
		}
	} else if *eslintLocal { // local 的情况
		eslibs, err := dependenciesNeedsToInstall(eslintDependencies, "package.json")
		if err != nil {
			return err
		}

		// 下载依赖到项目中
		if err := util.NpmInstallDependencies("", false, eslibs...); err != nil {
			return err
		}
	}

	return nil
}

// 不设置 eslint
func initProjectWithoutLint() (files util.FileContent) {
	// 不需要设置 cilint 的情况，直接写 setting
	settingJSON := genSettingsJSONwith("")
	files = util.FileContent{
		Path:    ".vscode/settings.json",
		Content: settingJSON,
	}
	return
}

// 设置 project eslint
// 需要写的文件:
// <project>/eslint/eslintrc-ts.json
// <project>/.vscode/settings.json, 替换 settings 中 -config 地址。
// npm install dependencies // FIXME
func initProjectWithLocalLint() (folders []string, files []util.FileContent, suggs []*util.Suggestion, err error) {
	// 获取绝对地址
	projectPath, er := filepath.Abs(".")
	if er != nil {
		return nil, nil, nil, er
	}
	// 添加 <project>/eslint 文件夹，添加 eslintrc-ts.json 文件
	esl := setupLocalEslint(projectPath)

	// setting.json 文件
	// 设置 settings.json 文件, 将 --config 设置为 cipath
	settingJSON, sug, er := _checkSettingJSON(esl.Espath)
	if er != nil {
		return nil, nil, nil, er
	}
	if sug != nil {
		suggs = append(suggs, sug)
	}
	if settingJSON != nil {
		// 添加 settings.json 文件到写入队列中
		esl.Files = append(esl.Files, util.FileContent{
			Path:    ".vscode/settings.json",
			Content: settingJSON,
		})
	}

	return esl.Folders, esl.Files, suggs, nil
}

// 设置 global golangci-lint
// 需要写的文件:
// ~/.vsc/golangci/dev-ci.yml, ~/.vsc/golangci/prod-ci.yml, 全局地址。
// ~/.vsc/vsc-config.json 全局配置文件。
// <project>/.vscode/settings.json, 替换 settings 中 -config 地址。
// npm install dependencies // FIXME
func initProjectWithGlobalLint() (folders []string, files []util.FileContent, suggs []*util.Suggestion, err error) {
	// 添加 ~/.vsc/golangci 文件夹，添加 dev-ci.yml, prod-ci.yml 文件
	// 添加 ~/.vsc/vsc-config.json 文件
	esl, err := setupGlobleEslint()
	if err != nil {
		return nil, nil, nil, err
	}

	// setting.json 文件
	// 设置 settings.json 文件, 将 --config 设置为 cipath
	settingJSON, sug, er := _checkSettingJSON(esl.Espath)
	if er != nil {
		return nil, nil, nil, er
	}
	if sug != nil {
		suggs = append(suggs, sug)
	}
	if settingJSON != nil {
		// 添加 settings.json 文件到写入队列中
		esl.Files = append(esl.Files, util.FileContent{
			Path:    ".vscode/settings.json",
			Content: settingJSON,
		})
	}

	return esl.Folders, esl.Files, suggs, nil
}

// 检查 .vscode/settings.json 是否存在
func _checkSettingJSON(esPath string) (newSetting []byte, sug *util.Suggestion, err error) {
	settingsPath, err := filepath.Abs(".vscode/settings.json")
	if err != nil {
		return nil, nil, err
	}

	sf, err := os.Open(settingsPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		// settings.json 不存在, 生成新的 settings.json 文件
		return genSettingsJSONwith(esPath), nil, nil
	}
	defer sf.Close()

	// 读取 settings.json 文件返回 eslint configFile 设置
	eslintConfigFile, err := _readSettingJSON(sf)
	if err != nil {
		return nil, nil, err
	}

	// 判断 --config 地址是否和要设置的 espath 相同, 如果相同则不更新 setting 文件。
	if eslintConfigFile == esPath { // 相同路径
		return nil, nil, nil
	}

	// 如果 settings.json 文件存在，而且 config != cipath, 则需要 suggestion
	// 建议手动添加设置到 .vscode/settings.json 中
	cilintConfig := bytes.ReplaceAll(eslintconfig, []byte(configPlaceHolder), []byte(esPath))
	return nil, &util.Suggestion{
		Problem:  "please add following in '.vscode/settings.json':",
		Solution: string(cilintConfig),
	}, nil
}

// 读取 setting.json 文件
func _readSettingJSON(file *os.File) (string, error) {
	// json 反序列化 settings.json
	jsonc, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	js, err := util.JSONCToJSON(jsonc)
	if err != nil {
		return "", err
	}

	type settingsStruct struct {
		EslintOption struct {
			ConfigFile string `json:"configFile,omitempty"`
		} `json:"eslint.options,omitempty"`
	}

	var settings settingsStruct
	err = json.Unmarshal(js, &settings)
	if err != nil {
		return "", err
	}

	return settings.EslintOption.ConfigFile, nil
}
