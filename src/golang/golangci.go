package golang

import (
	"bytes"
	"encoding/json"
	"errors"
	"local/src/util"
	"os"
)

const (
	// 整个 golangci-lint 的设置占位符
	lintPlaceHolder = "${golangcilintPlaceHolder}"

	// go.lintFlags --config 的占位符
	configPlaceHolder = "${configPlaceHolder}"

	// go.lintFlags --config 的占位符
	configVSComments = "//vsc:cilint"

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
	lintTool = []byte(`
  // golangci-lint 设置
  "go.lintTool": "golangci-lint"`)

	lintOnSave = []byte(`
  // NOTE save 时 golangci-lint 整个 package，使用 'file' 时，
  // 如果变量定义在别的文件中会造成 undeclared 错误。
  "go.lintOnSave": "package"`)

	lintFlags = []byte(`
  "go.lintFlags": [
    "--fast", // without --fast can freeze your editor.

    // golangci-lint 配置文件地址
    "--config=` + configPlaceHolder + `" ` + configVSComments + ` DON'T EDIT
  ]`)
)

// TODO golangci setup

// 设置全局 golangci-lint, 如果第一次写入，则生成新文件，
// 如果之前已经设置过，则直接返回 golangci lint config 的文件地址.
func setupGlobleCilint() (folders []string, files []util.FileContent, cipath string, err error) {
	vscDir, err := util.GetVscConfigDir()
	if err != nil {
		return nil, nil, "", err
	}

	// read vsc config file
	f, err := os.Open(vscDir + util.VscConfigFilePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, nil, "", err
	} else if errors.Is(err, os.ErrNotExist) {
		// ~/.vsc/vsc-config 文件不存在，创建文件夹，创建文件
		folders, files, cipath = writeCilintFiles(vscDir)
		return folders, files, cipath, nil
	}
	defer f.Close()

	// ~/.vsc/vsc-config 文件存在, 读取文件
	var vsSetting util.VscSetting
	de := json.NewDecoder(f)
	err = de.Decode(&vsSetting)
	if err != nil {
		return nil, nil, "", err
	}

	// 检查 golangci 设置
	if vsSetting.Golangci == "" { // 没有设置 golangci-lint
		folders, files, cipath = writeCilintFiles(vscDir)
		return folders, files, cipath, nil
	}

	// 已经设置 golangci-lint
	return nil, nil, vsSetting.Golangci, nil
}

// 设置项目 golangci-lint, 写入文件，返回 golangci lint config 的文件地址.
func setupLocalCilint(projectPath string) (folders []string, files []util.FileContent, cipath string) {
	folders, files, _ = writeCilintFiles(projectPath)
	return folders, files, vsWorkspace + golangciDirector + devciFilePath
}

// 写入 dev-ci.yml 和 prod-ci.yml 文件
func writeCilintFiles(dir string) (folders []string, files []util.FileContent, cipath string) {
	folders = append(folders, dir, dir+golangciDirector)
	files = append(files, util.FileContent{
		Path:    dir + golangciDirector + devciFilePath,
		Content: devci,
	}, util.FileContent{
		Path:    dir + golangciDirector + prodciFilePath,
		Content: prodci,
	})
	return folders, files, dir + golangciDirector + devciFilePath
}

// 组合 lintTool, lintOnSave, lintFlags 设置
func golangciSettings(settings ...[]byte) []byte {
	if len(settings) == 0 {
		return nil
	}
	return bytes.Join(settings, []byte(",\n"))
}

// 生成一个 settings.json 文件
func genNewSettingsFile(ciPath string) []byte {
	if ciPath == "" {
		return replaceCilintPlaceHolder(nil)
	}

	golangciConfig := append(golangciSettings(lintTool, lintOnSave, lintFlags), ',', '\n')
	r := bytes.ReplaceAll(golangciConfig, []byte(configPlaceHolder), []byte(ciPath))
	return replaceCilintPlaceHolder(r)
}

// 替换 settings_template.txt 模板中的 place holder, ${golangcilintPlaceHolder}
// 添加整个 golangci-lint setting.
func replaceCilintPlaceHolder(content []byte) []byte {
	return bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), content)
}
