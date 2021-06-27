package golang

import (
	"bytes"
	"local/src/util"
)

const (
	lintPlaceHolder   = "${golangcilintPlaceHolder}"
	configPlaceHolder = "//vsc:cilint"
)

// setup golangci-lint.yml config file.
func SetupLocalGolangciLint(ciDir string) (folders []string, files []util.FileContent) {
	folders = []string{ciDir, ciDir + util.GolangciDirector}
	files = []util.FileContent{
		{Path: ciDir + util.GolangciDirector + "/dev-ci.yml", Content: devci},
		{Path: ciDir + util.GolangciDirector + "/prod-ci.yml", Content: prodci},
	}
	return
}

var golangciConfig = []byte(`
  // golangci-lint 设置
  "go.lintTool": "golangci-lint",

  // NOTE save 时 golangci-lint 整个 package，使用 'file' 时，
  // 如果变量定义在别的文件中会造成 undeclared 错误。
  "go.lintOnSave": "package",

  "go.lintFlags": [
    "--fast", // without --fast can freeze your editor.

    // golangci-lint 配置文件地址
    // "--config=${workspaceRoot}/golangci.yml" // 本地
    "--config=xxx" ` + configPlaceHolder + ` DON'T EDIT
  ],
`)

// 替换 settings_template.txt 模板中的 place holder
func replaceLintPlaceHolder(setting []byte) []byte {
	return bytes.ReplaceAll(settingTemplate, []byte(lintPlaceHolder), setting)
}

// 替换模板中的 --config 行
func replaceLintConfig(file []byte, cilintCfgPath string) (settings []byte, sug *util.Suggestion, err error) {
	var (
		multiComment bool
		er           error
		found        bool // 是否找到了设置
		cfg          = "\"--config=" + cilintCfgPath + "\" " + configPlaceHolder + " DON'T EDIT"
	)

	lines := bytes.Split(file, []byte("\n"))
	for i := range lines {
		start := 0
		var buf bytes.Buffer
		if multiComment {
			ci := bytes.Index(lines[i], []byte("*/"))
			if ci == -1 {
				continue
			} else {
				start = ci + 2
			}
		}

		multiComment, er = util.JsoncLineTojson(lines[i], start, &buf)
		if er != nil {
			return nil, nil, er
		}

		if bytes.Contains(buf.Bytes(), []byte("\"--config=")) &&
			bytes.Contains(lines[i], []byte(configPlaceHolder)) {
			space := bytes.Index(lines[i], []byte("\"")) // 计算空格数量
			lines[i] = append(lines[i][:space], cfg...)
			found = true
			break
		}
	}

	if !found { // 如果没有找到 cilint 设置
		return nil, &util.Suggestion{
			Problem:  "can't find golangci-lint config, please add following settings to 'go.lintFlags'",
			Solution: cfg,
		}, nil
	}

	return bytes.Join(lines, []byte("\n")), nil, nil
}
