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

	// vscode workspace
	vsWorkspace = "${workspaceRoot}"
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

type golangciLintStruct struct {
	Folders []string
	Files   []util.FileContent
	Cipath  string
}

// 设置全局 golangci-lint, 如果第一次写入，则生成新文件，
// 如果之前已经设置过，则直接返回 golangci lint config 的文件地址.
func setupGlobleCilint() (*golangciLintStruct, error) {
	// 获取 .vsc 文件夹地址
	vscDir, err := util.GetVscConfigDir()
	if err != nil {
		return nil, err
	}

	// read vsc config file
	var vscSetting util.VscSetting
	err = vscSetting.ReadFromFile(vscDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsc/vsc-config 文件不存在，创建文件夹，创建文件
		return newGlobalCilintSetup(vscDir)
	}

	// 检查 golangci 设置
	if vscSetting.Golangci == "" {
		// 没有设置 golangci-lint 的情况
		gls := writeCilintFiles(vscDir)

		vscSetting.Golangci = gls.Cipath

		// json 格式化
		b, er := vscSetting.JSONIndentFormat()
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

	// 已经设置 golangci-lint
	return &golangciLintStruct{Cipath: vscSetting.Golangci}, nil
}

// 新写入 global golangci lint 设置
func newGlobalCilintSetup(vscDir string) (*golangciLintStruct, error) {
	// 生成 folders, files 加入创建队列
	gls := writeCilintFiles(vscDir)

	// 设置 global cilint 配置文件的地址
	vscSetting := util.VscSetting{
		Golangci: gls.Cipath,
	}

	// json 格式化
	b, er := vscSetting.JSONIndentFormat()
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

// 设置项目 golangci-lint, 写入文件，返回 golangci lint config 的文件地址.
func setupLocalCilint(projectPath string) *golangciLintStruct {
	gls := writeCilintFiles(projectPath)
	return &golangciLintStruct{gls.Folders, gls.Files, vsWorkspace + golangciDirector + devciFilePath}
}

// 在指定路径下写入 dev-ci.yml 和 prod-ci.yml 文件.
func writeCilintFiles(dir string) golangciLintStruct {
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

// 生成一个 settings.json 文件
func genNewSettingsFile(ciPath string) []byte {
	if ciPath == "" {
		// 如果 cipath 为空，则不设置 go.lint 到 settings.json 中
		return bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), nil)
	}

	// 设置 go.lint 到 settings.json 中，同时添加 cipath
	r := bytes.ReplaceAll(golangcilintconfig, []byte(configPlaceHolder), []byte(ciPath))
	return bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), r)
}
