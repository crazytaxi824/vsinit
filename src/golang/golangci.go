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

// golangci lint 配置文件的位置，和文件夹和文件
type golangciLintStruct struct {
	Folders []string
	Files   []util.FileContent
	Cipath  string // dev-ci.yml 的文件地址
}

// 设置 local golangci-lint, 生成文件 dev-ci.yml prod-ci.yml，
// 返回 golangci lint config 的文件地址.
func setupLocalCilint(projectPath string) *golangciLintStruct {
	// 生成 dev-ci.yml 和 prod-ci.yml 文件，返回文件地址。
	gls := _genCilintCfgFilesAndCipath(projectPath)
	return &gls
}

// 设置全局 golangci-lint, 如果第一次写入，则生成新文件 dev-ci.yml prod-ci.yml
// 如果之前已经设置过，则直接返回 golangci lint config 的文件地址.
func setupGlobleCilint() (*golangciLintStruct, error) {
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
		// 生成 dev-ci.yml, prod-ci.yml,vsc-config.yml 文件
		return _newGlobalCilintSetup(vscDir)
	}

	// ~/.vsc/vsc-config 文件存在
	// 生成 folders & files
	gls := _genCilintCfgFilesAndCipath(vscDir)

	// 检查 golangci 设置
	// 没有设置 golangci-lint 的情况
	if vscCfgYML.Golangci == "" {
		// 设置 golangci lint 配置文件地址
		vscCfgYML.Golangci = gls.Cipath

		// json 格式化
		b, er := vscCfgYML.JSONIndentFormat()
		if er != nil {
			return nil, er
		}

		// NOTE vsc-config.json 标记 overwrite, 否则不会重写文件
		gls.Files = append(gls.Files, util.FileContent{
			Path:      vscDir + util.VscConfigFilePath,
			Content:   b,
			Overwrite: true,
		})

		return &gls, nil
	}

	// 已经设置 golangci-lint，直接返回已有的 golangci lint 配置文件地址
	gls.Cipath = vscCfgYML.Golangci
	return &gls, nil
}

// 新写入 global golangci lint 配置
func _newGlobalCilintSetup(vscDir string) (*golangciLintStruct, error) {
	// 生成 dev-ci.yml 和 prod-ci.yml 文件，返回文件地址。
	gls := _genCilintCfgFilesAndCipath(vscDir)

	// 设置 global cilint 配置文件的地址
	vscCfgYML := util.VscConfigYML{
		Golangci: gls.Cipath,
	}

	// json 格式化
	b, er := vscCfgYML.JSONIndentFormat()
	if er != nil {
		return nil, er
	}

	// 将 vsc-config.json 文件加入创建队列
	gls.Files = append(gls.Files, util.FileContent{
		Path:    vscDir + util.VscConfigFilePath,
		Content: b,
	})

	return &gls, nil
}

// 生成 dev-ci.yml 和 prod-ci.yml 文件，返回文件地址。
func _genCilintCfgFilesAndCipath(dir string) golangciLintStruct {
	var gls golangciLintStruct

	// 创建 <dir>/golangci 文件夹，用于存放 dev-ci.yml, prod-ci.yml 文件
	gls.Folders = append(gls.Folders, dir, dir+golangciDirector)

	// 创建 dev-ci.yml, prod-ci.yml 文件
	gls.Files = append(gls.Files, util.FileContent{
		Path:    dir + golangciDirector + devciFilePath,
		Content: devci,
	}, util.FileContent{
		Path:    dir + golangciDirector + prodciFilePath,
		Content: prodci,
	})

	// ci.yml 的文件路径
	gls.Cipath = dir + golangciDirector + devciFilePath

	return gls
}

// 生成一个 settings.json 文件, 填入设置的 golangci lint path
func genSettingsJSONwith(ciPath string) []byte {
	if ciPath == "" {
		// 如果 cipath 为空，则不设置 go.lint 到 settings.json 中
		return bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), nil)
	}

	// 设置 go.lint 到 settings.json 中，同时添加 cipath
	r := bytes.ReplaceAll(golangcilintconfig, []byte(configPlaceHolder), []byte(ciPath))
	return bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), r)
}
