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
	{Path: ".gitignore", Content: gitignore},
	{Path: "package.json", Content: packageJSON},
	{Path: "tsconfig.json", Content: tsconfigJSON},
	{Path: "src/main.ts", Content: mainTS},
}

func InitProject(tsjsSet *flag.FlagSet, jestflag, eslint, eslintLocal *bool) (suggs []*util.Suggestion, err error) {
	// parse arges first
	// nolint // flag.ExitOnError will do the os.Exit(2)
	tsjsSet.Parse(os.Args[2:])

	// 初始化
	ff := initFoldersAndFiles(createFolders, filesAndContent)

	// 写入 test 相关文件
	if *jestflag {
		err = ff.writeTestFile()
		if err != nil {
			return nil, err
		}
	}

	if *eslint && *eslintLocal {
		return nil, errors.New("can not setup eslint globally and locally at same time")
	} else if *eslint && !*eslintLocal {
		// 设置 global eslint
		err = ff.initProjectWithGlobalLint()
	} else if !*eslint && *eslintLocal {
		// 设置 local eslint
		err = ff.initProjectWithLocalLint()
	} else {
		// 不设置 eslint, 只需要设置 settings.json 文件
		err = ff.initProjectWithoutLint()
	}

	if err != nil {
		return nil, err
	}

	// NOTE write project files first
	fmt.Println("init TypeScript project")
	if err := util.WriteFoldersAndFiles(ff.folders, ff.files); err != nil {
		return nil, err
	}

	// check and download dependencies
	if err := installMissingDependencies(jestflag, eslint, eslintLocal); err != nil {
		return nil, err
	}

	// 检查返回是否为空
	if len(ff.suggestions) != 0 {
		return ff.suggestions, nil
	}

	return nil, nil
}

// 不设置 eslint
func (ff *foldersAndFiles) initProjectWithoutLint() error {
	// 直接写 settings.json 文件
	err := ff.writeSettingJSON()
	if err != nil {
		return err
	}
	return nil
}

// 设置 project eslint
// 需要写的文件:
// <project>/eslint/eslintrc-ts.json
// <project>/.vscode/settings.json, 替换 settings 中 -config 地址。
// npm install dependencies
func (ff *foldersAndFiles) initProjectWithLocalLint() error {
	// 获取绝对地址
	projectPath, err := filepath.Abs(".")
	if err != nil {
		return err
	}
	// 添加 <project>/eslint 文件夹，添加 eslintrc-ts.json 文件
	// esl := setupLocalEslint(projectPath)
	ff.writeEslintJSONAndEspath(projectPath)

	// setting.json 文件
	// 设置 settings.json 文件, 将 --config 设置为 cipath
	err = ff.writeSettingJSON()
	if err != nil {
		return err
	}

	return nil
}

// 设置 global golangci-lint
// 需要写的文件:
// ~/.vsc/golangci/dev-ci.yml, ~/.vsc/golangci/prod-ci.yml, 全局地址。
// ~/.vsc/vsc-config.json 全局配置文件。
// <project>/.vscode/settings.json, 替换 settings 中 -config 地址。
// npm install dependencies
func (ff *foldersAndFiles) initProjectWithGlobalLint() error {
	// 获取 .vsc 文件夹地址
	vscDir, err := util.GetVscConfigDir()
	if err != nil {
		return err
	}

	// 通过 vsc-config.json 获取 eslint.TS 配置文件地址.
	// 如果 vsc-config.json 不存在，生成 vsc-config.json, eslintrc-ts.json 文件
	// 如果 vsc-config.json 存在，但是没有设置过 eslint.TS 配置文件地址，
	// 则 overwite vsc-config.json, eslintrc-ts.json 文件.
	// 如果 vsc-config.json 存在，同时也设置了 eslint.TS 配置文件地址，直接读取配置文件地址。
	err = ff.readEslintPathFromVscCfgJSON(vscDir)
	if err != nil {
		return err
	}

	// // setting.json 文件
	// // 设置 settings.json 文件, 将 --config 设置为 cipath
	err = ff.writeSettingJSON()
	if err != nil {
		return err
	}

	return nil
}

// 检查 .vscode/settings.json 是否存在
func (ff *foldersAndFiles) writeSettingJSON() error {
	if ff.espath == "" {
		// 不设置 eslint 的情况
		ff.addFiles(genSettingsJSONwith(""))
		return nil
	}

	eslintConfigFile, err := _readSettingJSON()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// settings.json 不存在, 生成新的 settings.json 文件
		ff.addFiles(genSettingsJSONwith(ff.espath))
		return nil
	}

	// 判断 --config 地址是否和要设置的 espath 相同, 如果相同则不更新 setting 文件。
	if eslintConfigFile == ff.espath { // 相同路径
		return nil
	}

	// 如果 settings.json 文件存在，而且 config != espath, 则需要 suggestion
	// 建议手动添加设置到 .vscode/settings.json 中
	cilintConfig := bytes.ReplaceAll(eslintconfig, []byte(configPlaceHolder), []byte(ff.espath))
	ff.addSuggestion(&util.Suggestion{
		Problem:  "please add following in '.vscode/settings.json':",
		Solution: string(cilintConfig),
	})

	return nil
}

// 读取 setting.json 文件
func _readSettingJSON() (string, error) {
	// 读取 .vscode/settings.json
	settingsPath, err := filepath.Abs(".vscode/settings.json")
	if err != nil {
		return "", err
	}

	sf, err := os.Open(settingsPath)
	if err != nil {
		return "", err
	}
	defer sf.Close()

	// json 反序列化 settings.json
	jsonc, err := io.ReadAll(sf)
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

// 安装依赖
func installMissingDependencies(jestflag, eslint, eslintLocal *bool) error {
	// NOTE 安装依赖, 必须放在后面，否则 package.json 需要改写。
	if *jestflag {
		// 检查本地 package.json 文件
		err := _checkAndInstallMissingDependencies("package.json", "", jestDependencies)
		if err != nil {
			return err
		}
	}

	// 下载 dependencies
	if *eslint { // global 情况
		vscDir, err := util.GetVscConfigDir()
		if err != nil {
			return err
		}

		eslintFolder := vscDir + eslintDirector
		pkgFilePath := eslintFolder + "/package.json"

		// NOTE 读取 ~/.vsc/eslint/package.json 文件
		err = _checkAndInstallMissingDependencies(pkgFilePath, eslintFolder, eslintDependencies)
		if err != nil {
			return err
		}

		return nil // NOTE 这里不需要再继续匹配了
	} else if *eslintLocal { // local 的情况
		// 检查本地 package.json 文件
		err := _checkAndInstallMissingDependencies("package.json", "", eslintDependencies)
		if err != nil {
			return err
		}
	}

	return nil
}

// 检查哪些依赖没有装，然后安装依赖。
func _checkAndInstallMissingDependencies(pkgJSONPath, prefix string, checkDeps []string) error {
	npmLibs, err := dependenciesNeedsToInstall(checkDeps, pkgJSONPath)
	if err != nil {
		return err
	}

	err = util.NpmInstallDependencies(prefix, false, npmLibs...)
	if err != nil {
		return err
	}

	return nil
}
