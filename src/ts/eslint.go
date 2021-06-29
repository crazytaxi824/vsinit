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

	// vscode workspace
	vsWorkspace = "${workspaceRoot}"
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
    // 这里是全局 eslint 配置文件的固定地址
    // ts lint config file
    "configFile": "` + configPlaceHolder + `"
  },
`)
)

// eslint 配置文件的位置，和文件夹和文件
type esLintStruct struct {
	Folders []string
	Files   []util.FileContent
	Espath  string // dev-ci.yml 的文件地址
}

// 设置项目 eslint, 写入 eslintrc-json 文件，返回 eslint config 的文件地址.
func setupLocalEslint(projectPath string) *esLintStruct {
	// 生成 eslintrc-ts.json 文件，返回文件地址。
	esl := _genEslintCfgFilesAndEspath(projectPath)

	// 使用 ${workspaceRoot} 替代绝对路径
	return &esLintStruct{esl.Folders, esl.Files, vsWorkspace + eslintDirector + eslintFilePath}
}

// 设置全局 eslint, 如果第一次写入，则生成新文件 eslintrc-ts.yml
// 如果之前已经设置过，则直接返回 eslint config 的文件地址.
func setupGlobleEslint() (*esLintStruct, error) {
	// 获取 .vsc 文件夹地址
	vscDir, err := util.GetVscConfigDir()
	if err != nil {
		return nil, err
	}

	// 读取 ~/.vsc/vsc-config.yml 文件
	var vscCfgYML util.VscConfigYML
	err = vscCfgYML.ReadFromFile(vscDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsc/vsc-config 文件不存在，
		// 生成 eslintrc-ts.json 文件
		return _newGlobalEslintSetup(vscDir)
	}

	// ~/.vsc/vsc-config 文件存在
	// 生成 folders & files
	esl := _genEslintCfgFilesAndEspath(vscDir)

	// 检查 eslint 设置
	// 没有设置 eslint 的情况
	// TODO JS 记得要改
	if vscCfgYML.Eslint.TS == "" {
		// 设置 eslint 配置文件地址
		vscCfgYML.Eslint.TS = esl.Espath

		// json 格式化
		b, er := vscCfgYML.JSONIndentFormat()
		if er != nil {
			return nil, er
		}

		// NOTE vsc-config.json 标记 overwrite, 否则不会重写文件
		esl.Files = append(esl.Files, util.FileContent{
			Path:      vscDir + util.VscConfigFilePath,
			Content:   b,
			Overwrite: true,
		})

		return &esl, nil
	}

	// 已经设置 eslint，直接返回已有的 eslint 配置文件地址
	esl.Espath = vscCfgYML.Eslint.TS // TODO JS 记得要改
	return &esl, nil
}

// 新写入 global eslint 配置
func _newGlobalEslintSetup(vscDir string) (*esLintStruct, error) {
	// 生成 eslintrc-ts.json 文件，返回文件地址。
	esl := _genEslintCfgFilesAndEspath(vscDir)

	// 设置 global cilint 配置文件的地址
	// TODO JS 记得要改
	var vscCfgYML util.VscConfigYML
	vscCfgYML.Eslint.TS = esl.Espath

	// json 格式化
	b, er := vscCfgYML.JSONIndentFormat()
	if er != nil {
		return nil, er
	}

	// 将 vsc-config.json 文件加入创建队列
	esl.Files = append(esl.Files, util.FileContent{
		Path:    vscDir + util.VscConfigFilePath,
		Content: b,
	})

	return &esl, nil
}

// 生成 eslintrc-ts.json 文件，返回文件地址。
func _genEslintCfgFilesAndEspath(dir string) esLintStruct {
	var esl esLintStruct

	// 创建 <dir>/eslint 文件夹，用于存放 eslintrc-ts.json 文件
	esl.Folders = append(esl.Folders, dir, dir+eslintDirector)

	// 创建 eslintrc-ts.json 文件
	esl.Files = append(esl.Files, util.FileContent{
		Path:    dir + eslintDirector + eslintFilePath,
		Content: eslintrcJSON,
	})

	// eslintrc-ts.json 的文件路径
	esl.Espath = dir + eslintDirector + eslintFilePath

	return esl
}

// 生成一个 settings.json 文件, 填入设置的 eslint path
func genSettingsJSONwith(esPath string) []byte {
	if esPath == "" {
		// 如果 espath 为空，则不设置 eslint 到 settings.json 中
		return bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), nil)
	}

	// 设置 eslint 到 settings.json 中，同时添加 espath
	r := bytes.ReplaceAll(eslintconfig, []byte(configPlaceHolder), []byte(esPath))
	return bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), r)
}
