package python

import (
	"local/src/util"
)

var extensions = []string{"ms-python.python", "ms-python.vscode-pylance"}

func checkPython() ([]*util.Suggestion, error) {
	var suggs []*util.Suggestion

	// 检查 go 是否安装 // 检查 code 安装
	sug := util.CheckCMDInstall("python3")
	if sug != nil {
		suggs = append(suggs, sug)
	}

	// 检查 vscode and extensions,
	su, err := util.CheckVscodeAndExtensions(extensions)
	if err != nil {
		return nil, err
	}
	if su != nil {
		suggs = append(suggs, su...)
	}

	// 检查返回是否为空
	if len(suggs) == 0 {
		return nil, nil
	}

	return suggs, nil
}
