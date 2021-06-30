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

	// golangci 文件夹名
	golangciDirector = "/golangci"

	// golangci-lint 配置文件名
	cilintFilePath = "/golangci.yml"
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

// 通过 vsc-config.json 获取 golangci 配置文件地址.
//  - 如果 vsc-config.json 不存在，生成 vsc-config.json, golangci.yml 文件.
//  - 如果 vsc-config.json 存在，但是没有设置 golangci 配置文件地址，则 overwite vsc-config.json, golangci.yml 文件.
//  - 如果 vsc-config.json 存在，同时也设置了 golangci 配置文件地址，直接读取配置文件地址.
func (ff *foldersAndFiles) readCilintPathFromVscCfgJSON(vscDir string) error {
	// 读取 ~/.vsc/vsc-config.json 文件
	var vscCfgJSON util.VscConfigJSON
	err := vscCfgJSON.ReadFromDir(vscDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsc/vsc-config.json 文件不存在, 则生成该文件.
		return ff.addVscCfgJSON(vscDir, vscCfgJSON, false)
	}

	// 检查 golangci 设置
	if vscCfgJSON.Golangci == "" {
		// 没有设置 golangci-lint 的情况, //NOTE overwrite vsc-config.json 文件.
		return ff.addVscCfgJSON(vscDir, vscCfgJSON, true)
	}

	// 已经设置 golangci-lint，直接返回已有的 golangci lint 配置文件地址
	ff.cipath = vscCfgJSON.Golangci
	return nil
}

// 添加 ~/.vsc/vsc-config.json 文件
func (ff *foldersAndFiles) addVscCfgJSON(vscDir string, vscCfgJSON util.VscConfigJSON, overwrite bool) error {
	// 全局设置需要多添加一个 folder.
	ff._addFolders(vscDir)

	// 设置 vsc-config 文件之前需要生成 golangci.yml 文件, 并获取文件地址.
	ff.addCilintYMLAndCipath(vscDir + golangciDirector)

	// 设置 vsc-config.json 文件中的 golangci 配置文件地址
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

// 生成 golangci.yml 文件，记录配置文件地址。
//  - 如果是 global 设置需要多添加一个 folder.
func (ff *foldersAndFiles) addCilintYMLAndCipath(dir string) {
	// 创建 golangci.yml 文件
	ff._addFiles(util.FileContent{
		Path:    dir + cilintFilePath,
		Content: golangciYML,
	})

	// golangci.yml 的文件路径
	ff.cipath = dir + cilintFilePath
}

// 生成一个 settings.json 文件, 填入设置的 golangci 配置文件地址
func newSettingsJSONwith(ciPath string) util.FileContent {
	if ciPath == "" {
		// 如果 cipath 为空，则不设置 go.lint 到 settings.json 中
		return util.FileContent{
			Path:    util.SettingsJSONPath,
			Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), nil),
		}
	}

	// go.lint 中的 ${configPlaceHolder} 替换成 cilint 配置文件的地址
	r := bytes.ReplaceAll(golangcilintconfig, []byte(configPlaceHolder), []byte(ciPath))
	return util.FileContent{
		Path: util.SettingsJSONPath,
		// 将 setting template 中的 ${golangcilintPlaceHolder} 替换成整个 cilint 的设置.
		Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), r),
	}
}
