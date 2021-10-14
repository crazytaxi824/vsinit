package ts

import (
	"bytes"
	"errors"
	"local/src/util"
	"os"
	"strings"
)

// eslint dependencies
var eslintDependencies = []string{
	"@typescript-eslint/eslint-plugin", // 必须
	"@typescript-eslint/parser",        // 必须
	"eslint-config-airbnb",             // 依赖
	"eslint-config-airbnb-typescript",  // ts 用
	"eslint-plugin-jest",               // jest unit test
	"eslint-plugin-promise",            // promise 用法
	"eslint-config-prettier",           // 解决 vscode 插件中 prettier 造成的代码问题
}

const (
	// 整个 ESLint 的设置占位符
	lintPlaceHolder = "${eslintPlaceHolder}"

	// "eslint.options{configFile}" 的占位符
	configPlaceHolder = "${configPlaceHolder}"

	// ESLint 文件夹名
	eslintDirector = "/eslint"

	// ESLint 配置文件名
	eslintFilePath = "/eslintrc-ts.json" // NOTE JS 要改
)

// ESLint setting
var eslintconfig = `  // 在 OUTPUT -> ESlint 频道打印 debug 信息. 用于配置 eslint.
  // 开启 eslint
  // "eslint.enable": false,

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

  // 自动修复 eslint rules
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },

  // 这里需要使用 class API, eslint.options 中才能够使用 new ESLint API - overrideConfigFile
  // https://eslint.org/docs/developer-guide/nodejs-api#eslint-class
  "eslint.useESLintClass": true,
  "eslint.options": {
    // eslint 配置文件地址
    "overrideConfigFile": "` + configPlaceHolder + `"
  },`

// 通过 vsi-config.json 获取 eslint.TS 配置文件地址.
//  - 如果 vsi-config.json 不存在, 则生成 vsi-config.json, eslintrc-ts.json 文件.
//  - 如果 vsi-config.json 存在，但是没有设置 eslint.TS 配置文件地址, 则 overwite vsi-config.json, eslintrc-ts.json 文件.
//  - 如果 vsi-config.json 存在，同时也设置了 eslint.TS 配置文件地址, 直接读取配置文件地址.
func readEslintPathFromVsiCfgJSON(ctx *util.VSContext, vsiDir string) error {
	// 读取 ~/.vsi/vsi-config.json 文件
	var vsiCfgJSON util.VsiConfigJSON
	err := vsiCfgJSON.ReadFromDir(vsiDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsi/vsi-config.json 文件不存在, 则生成该文件.
		return addVsiCfgJSON(ctx, vsiDir, vsiCfgJSON, false)
	}

	// 检查 eslint 设置情况
	if vsiCfgJSON.Eslint.TS == "" { // NOTE JS 要改
		// 没有设置 golangci-lint 的情况, //NOTE overwrite vsi-config.json 文件.
		return addVsiCfgJSON(ctx, vsiDir, vsiCfgJSON, true)
	}

	// 已经设置 eslint，直接返回已有的 eslint 配置文件地址
	ctx.SetLintPath(vsiCfgJSON.Eslint.TS) // NOTE JS 要改
	return nil
}

// 添加 ~/.vsi/vsi-config.json 文件
func addVsiCfgJSON(ctx *util.VSContext, vsiDir string, vsiCfgJSON util.VsiConfigJSON, overwrite bool) error {
	// 全局设置需要多添加多个 folder
	ctx.AddFolders(vsiDir, vsiDir+eslintDirector)

	// 设置 vsi-config 文件之前需要生成 eslint 配置文件, 并获取文件地址.
	ctx.AddLintConfigAndLintPath(vsiDir+eslintDirector+eslintFilePath, eslintrcJSON)

	// 设置 vsi-config.json 文件中的 ESLint 配置文件地址
	vsiCfgJSON.Eslint.TS = ctx.LintPath() // NOTE JS 要改

	b, er := vsiCfgJSON.JSONIndentFormat()
	if er != nil {
		return er
	}

	ctx.AddFiles(util.FileContent{
		Path:      vsiDir + util.VsiConfigFilePath,
		Content:   b,
		Overwrite: overwrite,
	})

	return nil
}

// 添加 .vscode/settings.json 文件，如果文件存在则给出建议
func addSettingJSON(ctx *util.VSContext) error {
	if ctx.LintPath() == "" {
		// 不设置 eslint 的情况
		ctx.AddFiles(newSettingsJSONwith(""))
		return nil
	}

	// 读取 .vscode/settings.json 文件, 获取 "eslint.options{configFile}" 的值
	// 只需要读取 configFile
	type settingsStruct struct {
		EslintOption struct {
			ConfigFile string `json:"configFile,omitempty"`
		} `json:"eslint.options,omitempty"`
	}

	var settings settingsStruct
	err := util.ReadSettingJSON(&settings)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// settings.json 不存在, 生成新的 settings.json 文件
		ctx.AddFiles(newSettingsJSONwith(ctx.LintPath()))
		return nil
	}

	// settings.json 存在的情况
	// 判断 configFile 地址是否和要设置的 espath 相同, 如果相同则不更新 setting 文件.
	if settings.EslintOption.ConfigFile == ctx.LintPath() {
		return nil
	}

	// 如果 settings.json 文件存在，而且 configFile != lintpath, 则需要 suggestion
	// 建议手动添加设置到 .vscode/settings.json 中
	lintConfig := strings.ReplaceAll(eslintconfig, configPlaceHolder, ctx.LintPath())
	ctx.AddSuggestions(&util.Suggestion{
		Problem:  "please add following in '.vscode/settings.json':",
		Solution: "{\n" + lintConfig + "\n}",
	})

	return nil
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
	r := "\n" + strings.ReplaceAll(eslintconfig, configPlaceHolder, esPath) + "\n"

	return util.FileContent{
		Path: util.SettingsJSONPath,
		// 将 setting template 中的 ${eslintPlaceHolder} 替换成整个 ESLint 的设置.
		Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), []byte(r)),
	}
}
