package js

import (
	"bytes"
	"errors"
	"local/src/util"
	"os"
)

// eslint dependencies
var eslintDependencies = []string{
	"eslint-plugin-jest",        // jest unit test
	"eslint-plugin-promise",     // promise 用法
	"eslint-config-prettier",    // 解决 vscode 插件中 prettier 造成的代码问题
	"eslint-config-airbnb-base", // js 专用 lint
}

const (
	// 整个 ESLint 的设置占位符
	lintPlaceHolder = "${eslintPlaceHolder}"

	// "eslint.options{configFile}" 的占位符
	configPlaceHolder = "${configPlaceHolder}"

	// ESLint 文件夹名
	eslintDirector = "/eslint"

	// ESLint 配置文件名
	eslintFilePath = "/eslintrc-js.json"
)

// ESLint setting
var eslintconfig = []byte(`
  // 在 OUTPUT -> ESlint 频道打印 debug 信息. 用于配置 eslint.
  "eslint.debug": true,

  // save 的时候运行 eslint
  "eslint.run": "onSave",

  // eslint 检查文件类型
  "eslint.validate": [
    "javascript"
  ],

  // 单独设置 eslint 配置文件
  "eslint.options": {
    // NOTE eslint(cmd)<=v7.x 可以工作，但是 CLIEngine 已经弃用.
    // https://eslint.org/docs/developer-guide/nodejs-api#cliengine
    // eslint config file
    "configFile": "` + configPlaceHolder + `"
  },
`)

// 通过 vsc-config.json 获取 eslint.JS 配置文件地址.
//  - 如果 vsc-config.json 不存在, 则生成 vsc-config.json, eslintrc-js.json 文件.
//  - 如果 vsc-config.json 存在，但是没有设置 eslint.JS 配置文件地址, 则 overwite vsc-config.json, eslintrc-js.json 文件.
//  - 如果 vsc-config.json 存在，同时也设置了 eslint.JS 配置文件地址, 直接读取配置文件地址.
func (ff *foldersAndFiles) readEslintPathFromVscCfgJSON(vscDir string) error {
	// 读取 ~/.vsc/vsc-config.json 文件
	var vscCfgJSON util.VscConfigJSON
	err := vscCfgJSON.ReadFromDir(vscDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsc/vsc-config.json 文件不存在, 则生成该文件.
		return ff.addVscCfgJSON(vscDir, vscCfgJSON, false)
	}

	// 检查 eslint 设置情况
	if vscCfgJSON.Eslint.JS == "" { // TODO JS 记得要改
		// 没有设置 golangci-lint 的情况, //NOTE overwrite vsc-config.json 文件.
		return ff.addVscCfgJSON(vscDir, vscCfgJSON, true)
	}

	// 已经设置 eslint，直接返回已有的 eslint 配置文件地址
	ff.espath = vscCfgJSON.Eslint.JS // TODO JS 记得要改
	return nil
}

// 添加 ~/.vsc/vsc-config.json 文件
func (ff *foldersAndFiles) addVscCfgJSON(vscDir string, vscCfgJSON util.VscConfigJSON, overwrite bool) error {
	// 设置 vsc-config 文件之前需要生成 eslint 配置文件, 并获取文件地址.
	ff.addEslintJSONAndEspath(vscDir, true)

	// 设置 vsc-config.json 文件中的 ESLint 配置文件地址
	vscCfgJSON.Eslint.JS = ff.espath // TODO JS 要改

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

// 生成 eslintrc-js.json 文件，记录 eslint 配置文件地址。
//  - 如果是 global 设置需要多添加一个 folder.
func (ff *foldersAndFiles) addEslintJSONAndEspath(dir string, global bool) {
	if global {
		// 创建 <dir>/eslint 文件夹，用于存放 eslintrc-ts.json 文件
		ff._addFolders(dir, dir+eslintDirector)

		// 创建 eslintrc-ts.json 文件
		ff._addFiles(util.FileContent{
			Path:    dir + eslintDirector + eslintFilePath,
			Content: eslintrcJSON,
		})

		// eslintrc-ts.json 的文件路径
		ff.espath = dir + eslintDirector + eslintFilePath
		return
	}

	ff._addFiles(util.FileContent{
		Path:    dir + eslintFilePath,
		Content: eslintrcJSON,
	})

	// eslintrc-ts.json 的文件路径
	ff.espath = dir + eslintFilePath
}

// 生成一个新的 settings.json 文件, 填入设置的 ESLint 配置文件地址
func newSettingsJSONwith(esPath string) util.FileContent {
	if esPath == "" {
		// 如果 espath 为空，则不设置 eslint 到 settings.json 中
		return util.FileContent{
			Path:    util.SettingsJSONPath,
			Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), nil),
		}
	}

	// ESLint 中的 ${configPlaceHolder} 替换成 ESLint 配置文件的地址
	r := bytes.ReplaceAll(eslintconfig, []byte(configPlaceHolder), []byte(esPath))
	return util.FileContent{
		Path: util.SettingsJSONPath,
		// 将 setting template 中的 ${eslintPlaceHolder} 替换成整个 ESLint 的设置.
		Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), r),
	}
}
