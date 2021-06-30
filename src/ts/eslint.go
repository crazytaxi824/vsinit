package ts

import (
	"bytes"
	"errors"
	"local/src/util"
	"os"
)

// eslint dependencies
var eslintDependencies = []string{
	"eslint-plugin-import",
	"eslint-plugin-jsx-a11y",
	"eslint-plugin-react",
	"eslint-plugin-react-hooks",
	"@typescript-eslint/parser", // parser
	"@typescript-eslint/eslint-plugin",
	"eslint-plugin-jest",              // jest unit test
	"eslint-plugin-promise",           // promise 用法
	"eslint-config-airbnb-typescript", // ts 用
	"eslint-config-prettier",          // 解决 vscode 插件中 prettier 造成的代码问题
	"eslint-config-airbnb-base",       // js 专用 lint
}

const (
	// 整个 golangci-lint 的设置占位符
	lintPlaceHolder = "${eslintPlaceHolder}"

	// go.lintFlags --config 的占位符
	configPlaceHolder = "${configPlaceHolder}"

	// golangci 文件夹
	eslintDirector = "/eslint"

	// golangci-lint config file path
	eslintFilePath = "/eslintrc-ts.json"
)

// golangci-lint setting
var (
	eslintconfig = []byte(`
  // 在 OUTPUT -> ESlint 频道打印 debug 信息. 用于配置 eslint.
  "eslint.debug": true,

  // save 的时候运行 eslint
  "eslint.run": "onSave",

  // eslint 检查文件类型
  "eslint.validate": [
    "typescriptreact",
    "typescript",
    "javascriptreact",
    "javascript"
  ],

  // 单独设置 eslint 配置文件
  "eslint.options": {
    // NOTE eslint(cmd)<=v7.x 可以工作，但是 CLIEngine 已经弃用。
    // https://eslint.org/docs/developer-guide/nodejs-api#cliengine
    // eslint 配置文件地址
    "configFile": "` + configPlaceHolder + `"
  },
`)
)

func (ff *foldersAndFiles) addMissingGlobalEslintDependencies() error {
	vscDir, err := util.GetVscConfigDir()
	if err != nil {
		return err
	}

	eslintFolder := vscDir + eslintDirector
	pkgFilePath := eslintFolder + "/package.json"

	// NOTE 读取 ~/.vsc/eslint/package.json 文件
	libs, err := checkMissingdependencies(eslintDependencies, pkgFilePath)
	if err != nil {
		return err
	}

	if len(libs) > 0 {
		ff._addDependencies(util.DependenciesInstall{
			Dependencies: libs,
			Prefix:       eslintFolder,
			Global:       false,
		})
	}
	return nil
}

func (ff *foldersAndFiles) addMissingLocalEslintDependencies() error {
	// 检查本地 package.json 文件
	libs, err := checkMissingdependencies(eslintDependencies, "package.json")
	if err != nil {
		return err
	}

	if len(libs) > 0 {
		ff._addDependencies(util.DependenciesInstall{
			Dependencies: libs,
			Prefix:       "",
			Global:       false,
		})
	}

	return nil
}

// 通过 vsc-config.json 获取 eslint.TS 配置文件地址.
// 如果 vsc-config.json 不存在，生成 vsc-config.json, eslintrc-ts.json 文件
// 如果 vsc-config.json 存在，但是没有设置过 eslint.TS 配置文件地址，
// 则 overwite vsc-config.json, eslintrc-ts.json 文件.
// 如果 vsc-config.json 存在，同时也设置了 eslint.TS 配置文件地址，直接读取配置文件地址。
func (ff *foldersAndFiles) readEslintPathFromVscCfgJSON(vscDir string) error {
	// 读取 ~/.vsc/vsc-config.json 文件
	var vscCfgJSON util.VscConfigJSON
	err := vscCfgJSON.ReadFromDir(vscDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsc/vsc-config.json 文件不存在
		return ff.addVscCfgJSON(vscDir, vscCfgJSON, false)
	}

	// 检查 eslint 设置情况
	if vscCfgJSON.Eslint.TS == "" {
		// 没有设置 golangci-lint 的情况
		return ff.addVscCfgJSON(vscDir, vscCfgJSON, true)
	}

	// 已经设置 eslint，直接返回已有的 eslint 配置文件地址
	ff.espath = vscCfgJSON.Eslint.TS // TODO JS 记得要改
	return nil
}

// 写 vsc-config.json 文件,
func (ff *foldersAndFiles) addVscCfgJSON(vscDir string, vscCfgJSON util.VscConfigJSON, overwrite bool) error {
	// 设置 vsc-config 文件之前需要生成 dev-ci.yml prod-ci.yml 文件
	// 并获取 cipath 地址.
	ff.addEslintJSONAndEspath(vscDir)

	// 设置 vsc-config.json 文件
	vscCfgJSON.Eslint.TS = ff.espath // TODO JS 要改

	b, er := vscCfgJSON.JSONIndentFormat()
	if er != nil {
		return er
	}

	ff._addFiles(util.FileContent{
		Path:      vscDir + util.VscConfigFilePath,
		Content:   b,
		Overwrite: overwrite,
	})

	return nil
}

// 生成 eslintrc-ts.json 文件，返回文件地址。
func (ff *foldersAndFiles) addEslintJSONAndEspath(dir string) {
	// 创建 <dir>/eslint 文件夹，用于存放 eslintrc-ts.json 文件
	ff._addFolders(dir, dir+eslintDirector)

	// 创建 eslintrc-ts.json 文件
	ff._addFiles(util.FileContent{
		Path:    dir + eslintDirector + eslintFilePath,
		Content: eslintrcJSON,
	})

	// eslintrc-ts.json 的文件路径
	ff.espath = dir + eslintDirector + eslintFilePath
}

// 生成一个 settings.json 文件, 填入设置的 eslint path
func newSettingsJSONwith(esPath string) util.FileContent {
	if esPath == "" {
		// 如果 espath 为空，则不设置 eslint 到 settings.json 中
		return util.FileContent{
			Path:    util.SettingsJSONPath,
			Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), nil),
		}
	}

	// 设置 eslint 到 settings.json 中，同时添加 espath
	r := bytes.ReplaceAll(eslintconfig, []byte(configPlaceHolder), []byte(esPath))
	return util.FileContent{
		Path:    util.SettingsJSONPath,
		Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), r),
	}
}
