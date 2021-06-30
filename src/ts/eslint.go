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
	// 整个 ESLint 的设置占位符
	lintPlaceHolder = "${eslintPlaceHolder}"

	// "eslint.options{configFile}" 的占位符
	configPlaceHolder = "${configPlaceHolder}"

	// ESLint 文件夹名
	eslintDirector = "/eslint"

	// ESLint 配置文件名
	eslintFilePath = "/eslintrc-ts.json" // TODO JS 需要更改
)

// ESLint setting
var eslintconfig = []byte(`
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

// 通过 vsc-config.json 获取 eslint.TS 配置文件地址.
//  - 如果 vsc-config.json 不存在, 则生成 vsc-config.json, eslintrc-ts.json 文件.
//  - 如果 vsc-config.json 存在，但是没有设置 eslint.TS 配置文件地址, 则 overwite vsc-config.json, eslintrc-ts.json 文件.
//  - 如果 vsc-config.json 存在，同时也设置了 eslint.TS 配置文件地址, 直接读取配置文件地址.
func readEslintPathFromVscCfgJSON(ff *util.FoldersAndFiles, vscDir string) error {
	// 读取 ~/.vsc/vsc-config.json 文件
	var vscCfgJSON util.VscConfigJSON
	err := vscCfgJSON.ReadFromDir(vscDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsc/vsc-config.json 文件不存在, 则生成该文件.
		return addVscCfgJSON(ff, vscDir, vscCfgJSON, false)
	}

	// 检查 eslint 设置情况
	if vscCfgJSON.Eslint.TS == "" { // TODO JS 记得要改
		// 没有设置 golangci-lint 的情况, //NOTE overwrite vsc-config.json 文件.
		return addVscCfgJSON(ff, vscDir, vscCfgJSON, true)
	}

	// 已经设置 eslint，直接返回已有的 eslint 配置文件地址
	ff.SetLintPath(vscCfgJSON.Eslint.TS) // TODO JS 记得要改
	return nil
}

// 添加 ~/.vsc/vsc-config.json 文件
func addVscCfgJSON(ff *util.FoldersAndFiles, vscDir string, vscCfgJSON util.VscConfigJSON, overwrite bool) error {
	// 全局设置需要多添加多个 folder
	ff.AddFolders(vscDir, vscDir+eslintDirector)

	// 设置 vsc-config 文件之前需要生成 eslint 配置文件, 并获取文件地址.
	// ff.addEslintJSONAndEspath(vscDir + eslintDirector + eslintFilePath)
	ff.AddLintConfigAndLintPath(vscDir+eslintDirector+eslintFilePath, eslintrcJSON)

	// 设置 vsc-config.json 文件中的 ESLint 配置文件地址
	vscCfgJSON.Eslint.TS = ff.LintPath() // TODO JS 要改

	b, er := vscCfgJSON.JSONIndentFormat()
	if er != nil {
		return er
	}

	ff.AddFiles(util.FileContent{
		Path:      vscDir + util.VscConfigFilePath,
		Content:   b,
		Overwrite: overwrite,
	})

	return nil
}

// 添加 .vscode/settings.json 文件，如果文件存在则给出建议
func addSettingJSON(ff *util.FoldersAndFiles) error {
	if ff.LintPath() == "" {
		// 不设置 eslint 的情况
		ff.AddFiles(newSettingsJSONwith(""))
		return nil
	}

	// 读取 .vscode/settings.json 文件, 获取 "eslint.options{configFile}" 的值
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
		ff.AddFiles(newSettingsJSONwith(ff.LintPath()))
		return nil
	}

	// settings.json 存在的情况
	// 判断 configFile 地址是否和要设置的 espath 相同, 如果相同则不更新 setting 文件.
	if settings.EslintOption.ConfigFile == ff.LintPath() {
		return nil
	}

	// 如果 settings.json 文件存在，而且 configFile != lintpath, 则需要 suggestion
	// 建议手动添加设置到 .vscode/settings.json 中
	cilintConfig := bytes.ReplaceAll(eslintconfig, []byte(configPlaceHolder), []byte(ff.LintPath()))
	ff.AddSuggestions(&util.Suggestion{
		Problem:  "please add following in '.vscode/settings.json':",
		Solution: string(cilintConfig),
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
	r := bytes.ReplaceAll(eslintconfig, []byte(configPlaceHolder), []byte(esPath))
	return util.FileContent{
		Path: util.SettingsJSONPath,
		// 将 setting template 中的 ${eslintPlaceHolder} 替换成整个 ESLint 的设置.
		Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), r),
	}
}
