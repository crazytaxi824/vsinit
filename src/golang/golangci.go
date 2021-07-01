package golang

import (
	"bytes"
	"errors"
	"local/src/util"
	"os"
	"strings"
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
var golangcilintconfig = `  // golangci-lint 设置
  "go.lintTool": "golangci-lint",

  // NOTE save 时 golangci-lint 整个 package，使用 'file' 时，
  // 如果变量定义在别的文件中会造成 undeclared 错误。
  "go.lintOnSave": "package",

  "go.lintFlags": [
    "--fast", // without --fast can freeze your editor.

    // golangci-lint 配置文件地址
    "--config=` + configPlaceHolder + `"
  ],`

// 通过 vsc-config.json 获取 golangci 配置文件地址.
//  - 如果 vsc-config.json 不存在，生成 vsc-config.json, golangci.yml 文件.
//  - 如果 vsc-config.json 存在，但是没有设置 golangci 配置文件地址，则 overwite vsc-config.json, golangci.yml 文件.
//  - 如果 vsc-config.json 存在，同时也设置了 golangci 配置文件地址，直接读取配置文件地址.
func readCilintPathFromVscCfgJSON(ff *util.FoldersAndFiles, vscDir string) error {
	// 读取 ~/.vsc/vsc-config.json 文件
	var vscCfgJSON util.VscConfigJSON
	err := vscCfgJSON.ReadFromDir(vscDir)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsc/vsc-config.json 文件不存在, 则生成该文件.
		return addVscCfgJSON(ff, vscDir, vscCfgJSON, false)
	}

	// 检查 golangci 设置
	if vscCfgJSON.Golangci == "" {
		// 没有设置 golangci-lint 的情况, //NOTE overwrite vsc-config.json 文件.
		return addVscCfgJSON(ff, vscDir, vscCfgJSON, true)
	}

	// 已经设置 golangci-lint，直接返回已有的 golangci lint 配置文件地址
	ff.SetLintPath(vscCfgJSON.Golangci)
	return nil
}

// 添加 ~/.vsc/vsc-config.json 文件
func addVscCfgJSON(ff *util.FoldersAndFiles, vscDir string, vscCfgJSON util.VscConfigJSON, overwrite bool) error {
	// 全局设置需要多添加多个 folder.
	ff.AddFolders(vscDir, vscDir+golangciDirector)

	// 设置 vsc-config 文件之前需要生成 golangci.yml 文件, 并获取文件地址.
	ff.AddLintConfigAndLintPath(vscDir+golangciDirector+cilintFilePath, golangciYML)

	// 设置 vsc-config.json 文件中的 golangci 配置文件地址
	vscCfgJSON.Golangci = ff.LintPath()

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
		// 不设置 golangci-lint 的情况
		ff.AddFiles(newSettingsJSONwith(""))
		return nil
	}

	// 读取 .vscode/settings.json, 获取 "go.lintFlags" 的值
	// 只需要读取 go.lintFlags
	type settingsStruct struct {
		GolingFlags []string `json:"go.lintFlags,omitempty"`
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

	// 判断 --config 地址是否和要设置的 cipath 相同, 如果相同则不更新 setting 文件。
	for _, v := range settings.GolingFlags {
		if v == "--config="+ff.LintPath() { // 相同的路径
			return nil
		}
	}

	// 如果 settings.json 文件存在，而且 config != cipath, 则需要 suggestion
	// 建议手动添加设置到 .vscode/settings.json 中
	lintConfig := strings.ReplaceAll(golangcilintconfig, configPlaceHolder, ff.LintPath())
	ff.AddSuggestions(&util.Suggestion{
		Problem:  "please add following in '.vscode/settings.json':",
		Solution: "{\n" + lintConfig + "\n}",
	})

	return nil
}

// 生成一个新的 settings.json 文件, 填入设置的 golangci 配置文件地址
func newSettingsJSONwith(ciPath string) util.FileContent {
	if ciPath == "" {
		// 如果 cipath 为空，则不设置 go.lint 到 settings.json 中
		return util.FileContent{
			Path:    util.SettingsJSONPath,
			Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), nil),
		}
	}

	// go.lint 中的 ${configPlaceHolder} 替换成 cilint 配置文件的地址
	r := "\n" + strings.ReplaceAll(golangcilintconfig, configPlaceHolder, ciPath) + "\n"
	return util.FileContent{
		Path: util.SettingsJSONPath,
		// 将 setting template 中的 ${golangcilintPlaceHolder} 替换成整个 cilint 的设置.
		Content: bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), []byte(r)),
	}
}
