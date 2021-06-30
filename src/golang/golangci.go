package golang

import (
	"bytes"
	"errors"
	"local/src/util"
	"os"
)

const (
	// 整个 golangci-lint 的设置占位符
	lintPlaceHolder = "${golangcilintPlaceHolder}"

	// go.lintFlags --config 的占位符
	configPlaceHolder = "${configPlaceHolder}"

	// golangci 文件夹
	golangciDirector = "/golangci"

	// golangci-lint config file path
	devciFilePath  = "/dev-ci.yml"
	prodciFilePath = "/prod-ci.yml"

	// setting.json 文件地址
	settingJSONPath = ".vscode/settings.json"
)

// golangci-lint setting
var (
	golangcilintconfig = []byte(`
  // golangci-lint 设置
  "go.lintTool": "golangci-lint",

  // NOTE save 时 golangci-lint 整个 package，使用 'file' 时，
  // 如果变量定义在别的文件中会造成 undeclared 错误。
  "go.lintOnSave": "package",

  "go.lintFlags": [
    "--fast", // without --fast can freeze your editor.

    // golangci-lint 配置文件地址
    "--config=` + configPlaceHolder + `"
  ],
`)
)

// 生成 dev-ci.yml 和 prod-ci.yml 文件，返回配置文件地址。
func (ff *foldersAndFiles) addCilintYMLAndCipath(dir string) {
	// 创建 <dir>/golangci 文件夹，用于存放 dev-ci.yml, prod-ci.yml 文件
	ff._addFolders(dir, dir+golangciDirector)

	// 创建 dev-ci.yml, prod-ci.yml 文件
	ff._addFiles(util.FileContent{
		Path:    dir + golangciDirector + devciFilePath,
		Content: devci,
	}, util.FileContent{
		Path:    dir + golangciDirector + prodciFilePath,
		Content: prodci,
	})

	// ci.yml 的文件路径
	ff.cipath = dir + golangciDirector + devciFilePath
}

// 通过 vsc-config.json 获取 golangci 配置文件地址.
// 如果 vsc-config.json 不存在，生成 vsc-config.json, dev-ci.yml, prod-ci.yml 文件
// 如果 vsc-config.json 存在，但是没有设置过 golangci 配置文件地址，
// 则 overwite vsc-config.json, dev-ci.yml, prod-ci.yml 文件.
// 如果 vsc-config.json 存在，同时也设置了 golangci 配置文件地址，直接读取配置文件地址。
func (ff *foldersAndFiles) readCilintPathFromVscCfgJSON(vscDir string) error {
	// 读取 ~/.vsc/vsc-config.yml 文件
	var vscCfgJSON util.VscConfigJSON
	err := vscCfgJSON.ReadFromDir(vscDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsc/vsc-config.json 文件不存在
		return ff.addVscCfgJSON(vscDir, vscCfgJSON, false)
	}

	// 检查 golangci 设置
	// 没有设置 golangci-lint 的情况
	if vscCfgJSON.Golangci == "" {
		// overwrite vsc-config.json 文件
		return ff.addVscCfgJSON(vscDir, vscCfgJSON, true)
	}

	// 已经设置 golangci-lint，直接返回已有的 golangci lint 配置文件地址
	ff.cipath = vscCfgJSON.Golangci
	return nil
}

// 写 vsc-config.json 文件,
func (ff *foldersAndFiles) addVscCfgJSON(vscDir string, vscCfgJSON util.VscConfigJSON, overwrite bool) error {
	// 设置 vsc-config 文件之前需要生成 dev-ci.yml prod-ci.yml 文件
	// 并获取 cipath 地址.
	ff.addCilintYMLAndCipath(vscDir)

	// 设置 vsc-config.json 文件
	vscCfgJSON.Golangci = ff.cipath

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

// 生成一个 settings.json 文件, 填入设置的 golangci lint path
func newSettingsJSONwith(ciPath string) util.FileContent {
	if ciPath == "" {
		// 如果 cipath 为空，则不设置 go.lint 到 settings.json 中
		return util.FileContent{
			Path:    settingJSONPath,
			Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), nil),
		}
	}

	// 设置 go.lint 到 settings.json 中，同时添加 cipath
	r := bytes.ReplaceAll(golangcilintconfig, []byte(configPlaceHolder), []byte(ciPath))
	return util.FileContent{
		Path:    settingJSONPath,
		Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), r),
	}
}
